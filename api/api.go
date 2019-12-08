package api

import (
	"github.com/gorilla/mux"

	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/logger"
	"github.com/muerwre/orchidgo/router/auth"
)

type API struct {
	App    *app.App
	Config *Config
	Logger *logger.Logger
}

func New(a *app.App) (api *API, err error) {
	api = &API{App: a}

	api.Config, err = InitConfig()

	if err != nil {
		return nil, err
	}

	api.Logger, _ = logger.CreateLogger(a)

	return api, nil
}

func (a *API) Init(r *mux.Router) {
	// r.Handle("/hello", gziphandler.GzipHandler(a.Logger.Log(a.RootHandler))).Methods("GET")

	auth.Router(r.PathPrefix("/test").Subrouter(), a.Logger)
}

// func (a *API) RootHandler(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
// _, err := w.Write([]byte(`{"hello" : "world"}`))
// return err
// }
