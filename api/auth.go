package api

import (
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/controller"
)

// AuthRouter for /api/auth/*
func AuthRouter(r *mux.Router, a *API) {
	r.Handle("/", a.Handler(controller.Auth.CheckCredentials)).Methods("GET")
	r.Handle("/vk", a.Handler(controller.Auth.LoginVkUser)).Methods("GET")
	r.Handle("/guest", a.Handler(controller.Auth.GetGuestUser)).Methods("GET")
}
