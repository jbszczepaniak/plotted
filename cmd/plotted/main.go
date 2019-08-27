package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jedruniu/plotted/pkg/handlers"
	"github.com/jedruniu/plotted/pkg/storage"
	fileStorage "github.com/jedruniu/plotted/pkg/storage/file"
	googleStorage "github.com/jedruniu/plotted/pkg/storage/google"
	"google.golang.org/appengine"

	"golang.org/x/oauth2"

	"net/http"
)

var (
	stravaClientID string
	stravaSecret   string
	mapboxToken    string
	port           int
	environment    string
	projectID      string
	selfUrl        string
	cache          storage.Storage
	stateStore     storage.Storage
)

func main() {
	// Google App Engine specific environment variables
	var err error
	port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("PORT not provided, or not an integer")
	}
	environment = os.Getenv("NODE_ENV")
	if environment == "" {
		environment = "dev"
	}

	// Application specific environment variables
	stravaClientID = os.Getenv("STRAVA_CLIENT_ID")
	if stravaClientID == "" {
		panic("STRAVA_CLIENT_ID not provided")
	}
	stravaSecret = os.Getenv("STRAVA_SECRET")
	if stravaSecret == "" {
		panic("STRAVA_SECRET not provided")
	}
	mapboxToken = os.Getenv("MAPBOX_TOKEN")
	if mapboxToken == "" {
		panic("MAPBOX_TOKEN not provided")
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)

	ctx := appengine.BackgroundContext()
	projectID = appengine.AppID(ctx)

	if environment == "production" {
		selfUrl = fmt.Sprintf("https://%s.appspot.com", projectID)
	} else {
		selfUrl = fmt.Sprintf("http://localhost:%d", port)
	}

	if environment == "production" {
		cache, err = googleStorage.NewGoogleStorage(ctx, projectID, "store")
		stateStore, err = googleStorage.NewGoogleStorage(ctx, projectID, "state_to_token")
	} else {
		cache, err = fileStorage.NewFileStorage("../../store")
		stateStore, err = fileStorage.NewFileStorage("../../store")
	}
	if err != nil {
		panic(err)
	}

	conf := &oauth2.Config{
		ClientID:     stravaClientID,
		ClientSecret: stravaSecret,
		Scopes:       []string{"activity:read_all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
		RedirectURL: fmt.Sprintf("%s/auth_callback", selfUrl),
	}

	mapServer := handlers.MapServer{conf, mapboxToken, cache, stateStore}
	authServer := handlers.AuthCallbackServer{conf, selfUrl, stateStore}
	indexServer := handlers.IndexServer{conf, stateStore}

	http.Handle("/auth_callback", &authServer)
	http.Handle("/map", &mapServer)
	http.Handle("/", &indexServer)

	log.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
