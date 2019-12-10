package api

import (
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/controller"
)

type AuthRouter struct {
}

// AuthRouter for /api/auth/*
func (a *AuthRouter) Init(router *mux.Router, api *API) {
	ctrl := &controller.AuthController{}

	router.Handle("/", api.Handler(ctrl.CheckCredentials)).Methods("GET")
	router.Handle("/vk", api.Handler(ctrl.LoginVkUser)).Methods("GET")
	router.Handle("/guest", api.Handler(ctrl.GetGuestUser)).Methods("GET")
}
