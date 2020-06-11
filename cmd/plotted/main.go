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
	stravaClientID  string
	stravaSecret    string
	mapboxToken     string
	port            int
	environment     string
	projectID       string
	selfURL         string
	gaeUsed         bool
	gaeCredentials  []byte
	fileStoragePath string
	cache           storage.Storage
	stateStore      storage.Storage
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

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
	gaeCredentials = []byte(os.Getenv("GOOGLE_CREDENTIALS"))

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
	fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePath != "" {
		if _, err := os.Stat(fileStoragePath); os.IsNotExist(err) {
			panic(fmt.Sprintf("provided FILE_STORAGE_PATH=%s does not exist", fileStoragePath))
		}
	}

	if environment == "production" {
		ctx := appengine.BackgroundContext()
		projectID = appengine.AppID(ctx)

		selfURL = fmt.Sprintf("https://%s.appspot.com", projectID)
		cache, err = googleStorage.NewGoogleStorage(ctx, projectID, "store", gaeCredentials)
		if err != nil {
			panic(err)
		}
		stateStore, err = googleStorage.NewGoogleStorage(ctx, projectID, "state_to_token", gaeCredentials)
		if err != nil {
			panic(err)
		}
	} else {
		selfURL = fmt.Sprintf("http://localhost:%d", port)
		cache, err = fileStorage.NewFileStorage(fileStoragePath)
		if err != nil {
			panic(err)
		}
		stateStore, err = fileStorage.NewFileStorage(fileStoragePath)
		if err != nil {
			panic(err)
		}
	}

	conf := &oauth2.Config{
		ClientID:     stravaClientID,
		ClientSecret: stravaSecret,
		Scopes:       []string{"activity:read_all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
		RedirectURL: fmt.Sprintf("%s/auth_callback", selfURL),
	}

	mapServer := handlers.MapServer{conf, mapboxToken, cache, stateStore}
	authServer := handlers.AuthCallbackServer{conf, selfURL, stateStore}
	indexServer := handlers.IndexServer{conf, stateStore}

	http.Handle("/auth_callback", &authServer)
	http.Handle("/map", &mapServer)
	http.Handle("/", &indexServer)

	log.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
