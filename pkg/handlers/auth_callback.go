package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jedruniu/plotted/pkg/storage"
	"golang.org/x/oauth2"
)

type AuthCallbackServer struct {
	OauthConfig *oauth2.Config
	SelfURL     string
	StateStore  storage.Storage
}

func (a *AuthCallbackServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	code := r.URL.Query().Get("code")
	callbackState := r.URL.Query().Get("state")
	ok, err := a.StateStore.Exists(reqCtx, callbackState)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not retrieve state information from store"), http.StatusBadRequest)
		return
	}
	if !ok {
		http.Error(w, fmt.Sprintf("state verification failed"), http.StatusBadRequest)
		return
	}

	token, err := a.OauthConfig.Exchange(reqCtx, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not exchange ouath2 token, err: %v", err), http.StatusInternalServerError)
		return
	}
	err = a.StateStore.Set(reqCtx, callbackState, []byte(token.AccessToken))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not set state information in the store"), http.StatusBadRequest)
		return
	}
	after := time.Now().AddDate(0, -3, 0).Format(layout)
	before := time.Now().Format(layout)
	http.Redirect(w, r, fmt.Sprintf("%s/map?after=%s&before=%s&state=%s", a.SelfURL, after, before, callbackState), 302)
}
