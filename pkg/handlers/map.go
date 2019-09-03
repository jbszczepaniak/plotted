package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
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
	var activitiesMutex sync.Mutex

	var summaryWg sync.WaitGroup
	summaryPool := make(chan int, 20)

	shouldFetch := true
	for i := 1; shouldFetch; i++ {
		summaryWg.Add(1)
		go func(i int) {
			defer summaryWg.Done()
			defer func() {
				<-summaryPool
			}()
			summaryPool <- 1
			fmt.Printf("Asking for page %d\n", i)

			opts.After = optional.NewInt32(int32(after.Unix()))
			opts.Before = optional.NewInt32(int32(before.Unix()))
			opts.Page = optional.NewInt32(int32(i))
			opts.PerPage = optional.NewInt32(200)
			summary, resp, err := client.ActivitiesApi.GetLoggedInAthleteActivities(reqCtx, &opts)
			if err != nil {
				fmt.Printf("error %v for page %d, do not fetch anymore\n	", err, i)
				shouldFetch = false
				return
			}
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("status code %d for page %d, do not fetch anymore\n", resp.StatusCode, i)
				shouldFetch = false
				return
			}
			if len(summary) == 0 {
				fmt.Printf("page %dreturned empty, do not fetch anymore\n", i)
				shouldFetch = false
				return

			}
			fmt.Printf("Page %d OK\n", i)

			activitiesMutex.Lock()
			activities = append(activities, summary...)
			activitiesMutex.Unlock()
		}(i)
		if i%5 == 0 {
			time.Sleep(1000 * time.Millisecond)
		}
	}
	summaryWg.Wait()

	var polylines [][][]float64
	var polylinesMutex sync.Mutex
	var wg sync.WaitGroup
	pool := make(chan int, 20)

	for _, activity := range activities {
		wg.Add(1)
		go func(activity swagger.SummaryActivity) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			pool <- 1

			var polyline []byte

			cachedPolyline := fmt.Sprintf("%d.polyline", activity.Id)
			exists, _ := m.Cache.Exists(reqCtx, cachedPolyline)
			if exists {
				polyline, _ = m.Cache.Get(reqCtx, cachedPolyline)
			} else {
				detailed, _, err := client.ActivitiesApi.GetActivityById(reqCtx, activity.Id, nil)
				if err != nil {
					log.Printf("err for activity %d, err: %v", activity.Id, err)
					return
				}
				if detailed.Map_.Polyline == "" {
					return
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
				polylinesMutex.Lock()
				polylines = append(polylines, polylineDecoded)
				polylinesMutex.Unlock()
			}
		}(activity)
	}
	wg.Wait()

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
