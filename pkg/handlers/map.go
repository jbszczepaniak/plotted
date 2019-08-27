package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"github.com/jedruniu/plotted/pkg/storage"
	swagger "github.com/jedruniu/plotted/swagger-generated"
	gopoly "github.com/twpayne/go-polyline"
	"golang.org/x/oauth2"
)

var layout = "02/01/2006"

type MapServer struct {
	OauthConfig *oauth2.Config
	MapboxToken string
	Cache       storage.Storage
	StateStore  storage.Storage
}

func (m *MapServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	cfg := swagger.NewConfiguration()
	client := swagger.NewAPIClient(cfg)
	state := r.URL.Query().Get("state")

	token, err := m.StateStore.Get(reqCtx, state)
	if err != nil {
		log.Fatal(err)

	}

	reqCtx = context.WithValue(reqCtx, swagger.ContextAccessToken, string(token))

	opts := swagger.GetLoggedInAthleteActivitiesOpts{}

	unparsedAfter := r.URL.Query().Get("after")
	unparsedBefore := r.URL.Query().Get("before")

	after, _ := time.Parse(layout, unparsedAfter)
	after = after.AddDate(0, 0, -1)
	before, _ := time.Parse(layout, unparsedBefore)
	before = before.AddDate(0, 0, 1)

	var activities []swagger.SummaryActivity

	for i := 1; ; i++ {
		opts.After = optional.NewInt32(int32(after.Unix()))
		opts.Before = optional.NewInt32(int32(before.Unix()))
		opts.Page = optional.NewInt32(int32(i))
		opts.PerPage = optional.NewInt32(200)

		summary, resp, err := client.ActivitiesApi.GetLoggedInAthleteActivities(reqCtx, &opts)
		if err != nil {
			http.Error(w, err.Error(), resp.StatusCode)
			return
		}
		if len(summary) == 0 {
			break
		}
		activities = append(activities, summary...)
	}

	var polylines [][][]float64

	for _, activity := range activities {
		var polyline []byte

		cachedPolyline := fmt.Sprintf("%d.polyline", activity.Id)
		exists, _ := m.Cache.Exists(reqCtx, cachedPolyline)
		if exists {
			polyline, _ = m.Cache.Get(reqCtx, cachedPolyline)
		} else {
			detailed, _, err := client.ActivitiesApi.GetActivityById(reqCtx, activity.Id, nil)
			if err != nil {
				log.Printf("err for activity %d, err: %v", activity.Id, err)
				continue
			}
			if detailed.Map_.Polyline == "" {
				continue
			}
			polyline = []byte(detailed.Map_.Polyline)
			err = m.Cache.Set(reqCtx, cachedPolyline, polyline)
			if err != nil {
				panic(err)
			}
		}

		var polylineDecoded [][]float64

		polylineDecoded, _, err := gopoly.DecodeCoords(polyline)
		if err != nil {
			log.Printf("could not decode polyline from file %d, err: %v", activity.Id, err)
		} else {
			polylines = append(polylines, polylineDecoded)
		}

	}

	tmpl, _ := template.New("").Parse(MapHTML)

	data := struct {
		EncodedRoutes [][][]float64
		MapboxToken   string
	}{
		polylines,
		m.MapboxToken,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}
