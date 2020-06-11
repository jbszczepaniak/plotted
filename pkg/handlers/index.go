package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jedruniu/plotted/pkg/storage"
	"golang.org/x/oauth2"
)

type IndexServer struct {
	OauthConfig *oauth2.Config
	StateStore  storage.Storage
}

func (i *IndexServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()

	err := i.StateStore.Set(r.Context(), state, []byte{})
	if err != nil {
		http.Error(w, fmt.Sprintf("could not set state information in the store"), http.StatusBadRequest)
		return
	}

	url := i.OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	tmpl, err := template.New("").Parse(IndexHTML)
	if err != nil {
		http.Error(w, fmt.Sprintf("parsing html file failed"), http.StatusBadRequest)
		return
	}

	data := struct{ Auth string }{url}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template failed, err: %v", err)
	}

}
