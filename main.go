package main

import (
	"flag"
	"fmt"
	"time"

	strava "github.com/strava/go.strava"
)

var (
	start    = flag.String("start", "", "start date")
	end      = flag.String("end", "", "end date")
	extended = flag.Bool("extended", false, "extended training information")
	token    = flag.String("token", "", "API Access token")
	layout   = "02/01/2006"
)

func init() {
	flag.Parse()
}

func main() {
	after, _ := time.Parse(layout, *start)
	before, _ := time.Parse(layout, *end)
	client := strava.NewClient(*token)
	currentAthleteService := strava.NewCurrentAthleteService(client)
	activities, _ := currentAthleteService.ListActivities().Before(int(before.Unix())).After(int(after.Unix())).Do()

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
	}
}
