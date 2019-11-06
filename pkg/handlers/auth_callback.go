package handlers

import (
	"fmt"
	"net/http"

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
		panic(err)
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
		panic(err)
	}

	http.Redirect(w, r, fmt.Sprintf("%s/map?after=30/05/2019&before=30/09/2019&state=%s", a.SelfURL, callbackState), 302)
}
