package api

import (
	"github.com/gorilla/mux"

	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/utils/logger"
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
	AuthRouter(r.PathPrefix("/auth").Subrouter(), a.Logger)
}
