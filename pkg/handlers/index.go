package handlers

import (
	"html/template"
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
		panic(err)
	}

	url := i.OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	tmpl, err := template.New("").Parse(IndexHTML)
	if err != nil {
		panic(err)
	}

	data := struct{ Auth string }{url}
	_ = tmpl.Execute(w, data)

}
