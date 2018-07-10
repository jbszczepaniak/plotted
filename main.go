package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	strava "github.com/strava/go.strava"
)

var (
	start       = flag.String("start", "", "start date")
	end         = flag.String("end", "", "end date")
	extended    = flag.Bool("extended", false, "extended training information")
	stravaToken = flag.String("strava", "", "Strava API Access token")
	mapBoxToken = flag.String("mapbox", "", "Mapbox API Access token")
	layout      = "02/01/2006"
)

func init() {
	flag.Parse()
}

func main() {
	after, _ := time.Parse(layout, *start)
	after = after.AddDate(0, 0, -1)
	before, _ := time.Parse(layout, *end)
	before = before.AddDate(0, 0, 1)
	client := strava.NewClient(*stravaToken)
	currentAthleteService := strava.NewCurrentAthleteService(client)
	activities, err := currentAthleteService.ListActivities().Before(int(before.Unix())).After(int(after.Unix())).Do()
	if err != nil {
		panic(err)
	}

	ids := []int64{}

	for _, activity := range activities {
		y, m, d := activity.StartDate.Date()
		fmt.Printf("%02d/%02d/%d, %s, Distance: %v km", d, m, y, activity.Type, activity.Distance/1000)
		if activity.Commute {
			fmt.Printf(" (commute)")
		}
		fmt.Printf("\n")
		if *extended {
			fmt.Printf("Average heartrate: %v bps\n", activity.AverageHeartrate)
			fmt.Printf("Maximum speed: %v m/s\n", activity.MaximunSpeed)
			fmt.Printf("Average speed: %v m/s\n", activity.AverageSpeed)
			fmt.Printf("Elevation gain: %v m\n", activity.TotalElevationGain)
			fmt.Printf("\n")
		}
		ids = append(ids, activity.Id)
	}
	polylines := []string{}

	activitiesService := strava.NewActivitiesService(client)
	for _, id := range ids {
		cachedFileName := fmt.Sprintf("cache/%d.json", id)

		_, err := os.Stat(cachedFileName)

		cacheExists := false
		if err == nil {
			cacheExists = true
		}

		var activity *strava.ActivityDetailed

		if cacheExists {
			cacheContent, _ := ioutil.ReadFile(cachedFileName)
			err := json.Unmarshal(cacheContent, &activity)
			if err != nil {
				panic(err)
			}
		} else {
			activity, _ = activitiesService.Get(id).Do()
			serialized, err := json.Marshal(activity)
			if err != nil {
				panic(err)
			}
			file, err := os.Create(cachedFileName)
			if err != nil {
				panic(err)
			}
			file.Write(serialized)
			file.Close()
		}

		polylines = append(polylines, floatTuples(activity.Map.Polyline.Decode()).String())
		fmt.Println("id: ", id, "size of polylines:", len(polylines))
	}

	templ, _ := template.ParseFiles("index.html")
	buf := new(bytes.Buffer)
	data := struct {
		EncodedRoutes []string
		MapboxToken   string
	}{
		polylines,
		*mapBoxToken,
	}

	_ = templ.Execute(buf, data)

	file, _ := os.Create("page.html")
	defer file.Close()
	file.Write(buf.Bytes())

}

type floatTuples [][2]float64

func (ft floatTuples) String() string {
	ftAsStringList := []string{}
	for _, elem := range ft {
		elemStr := "[" + strconv.FormatFloat(elem[0], 'f', 6, 64) + "," + strconv.FormatFloat(elem[1], 'f', 6, 64) + "]"
		ftAsStringList = append(ftAsStringList, elemStr)
	}

	return "[" + strings.Join(ftAsStringList, ",") + "]"
}
