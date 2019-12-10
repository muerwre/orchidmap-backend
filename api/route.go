package api

import (
	"github.com/gorilla/mux"
)

// RouteRouter for /api/route/*
func RouteRouter(router *mux.Router, api *API) {
	// ctrl := &controller.AuthController{}

	// router.Handle("/", api.Handler(ctrl.CheckCredentials)).Methods("GET")
	// router.Handle("/vk", api.Handler(ctrl.LoginVkUser)).Methods("GET")
	// router.Handle("/guest", api.Handler(ctrl.GetGuestUser)).Methods("GET")

	// 	router.post('/star', star);
	// router.post('/', post);
	// router.get('/', get);
	// router.patch('/', patch);
	// router.delete('/', drop);
	// router.get('/list', list);

}
