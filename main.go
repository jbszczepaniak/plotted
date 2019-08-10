package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"

	swagger "github.com/jedruniu/plotted/swagger-generated"

	"github.com/antihax/optional"
	gopoly "github.com/twpayne/go-polyline"

	"golang.org/x/oauth2"

	"net/http"
)

var (
	stravaClientID = flag.String("strava_clientID", "", "Strava client ID")
	stravaSecret   = flag.String("strava_secret", "", "Strava Secret")
	mapBoxToken    = flag.String("mapbox", "", "Mapbox API Access token")

	layout = "02/01/2006"
)

func init() {
	flag.Parse()
}

var code string
var token string

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     *stravaClientID,
		ClientSecret: *stravaSecret,
		Scopes:       []string{"activity:write,activity:read_all,profile:read_all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
		RedirectURL: "http://localhost:8888/auth_callback",
	}

	http.HandleFunc("/auth_callback", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")

		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			fmt.Println("tutaj?", err)
		}
		token = tok.AccessToken

		http.Redirect(w, r, "http://localhost:8888/map?after=30/01/2018&before=30/09/2019", 302)
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		cfg := swagger.NewConfiguration()
		client := swagger.NewAPIClient(cfg)

		ctx = context.WithValue(ctx, swagger.ContextAccessToken, token)

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

			summary, _, _ := client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, &opts)
			if len(summary) == 0 {
				break
			}
			activities = append(activities, summary...)
		}

		var floatPolylines [][][]float64

		for _, activity := range activities {
			cachedFileName := fmt.Sprintf("cache/%d.cache", activity.Id)

			_, err := os.Stat(cachedFileName)

			cacheExists := false
			if err == nil {
				cacheExists = true
			}

			var polyline []byte

			if cacheExists {
				cacheContent, err := ioutil.ReadFile(cachedFileName)
				if err != nil {
					log.Printf("error when reading %s, err: %v", cachedFileName, err)
				}
				polyline = cacheContent
			} else {
				detailed, _, err := client.ActivitiesApi.GetActivityById(ctx, activity.Id, nil)
				if err != nil {
					log.Printf("err for activity %d, err: %v", activity.Id, err)
					continue
				}
				if detailed.Map_.Polyline == "" {
					continue // activity without a map
				}
				polyline = []byte(detailed.Map_.Polyline)

				file, err := os.Create(cachedFileName)
				if err != nil {
					log.Printf("error when creting %s, err: %v", cachedFileName, err)
				}
				file.Write(polyline)
				file.Close()
			}

			var polyDecoded [][]float64

			polyDecoded, _, err = gopoly.DecodeCoords(polyline)
			if err != nil {
				log.Printf("for file %d, err: %v", activity.Id, err)
			} else {
				floatPolylines = append(floatPolylines, polyDecoded)
			}

		}

		templ, _ := template.ParseFiles("index_tmpl.html")

		data := struct {
			EncodedRoutes [][][]float64
			MapboxToken   string
		}{
			floatPolylines,
			*mapBoxToken,
		}
		templ.Delims("", "")

		templ.Execute(w, data)

	})

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	templ, err := template.ParseFiles("static/index_tmpl.html")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)

	data := struct{ Auth string }{url}
	_ = templ.Execute(buf, data)
	file, _ := os.Create("static/index.html")
	defer file.Close()
	file.Write(buf.Bytes())

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":8888", nil))
}
