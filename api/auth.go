package api

import (
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/controller"
	"github.com/muerwre/orchidgo/utils/logger"
)

// AuthRouter for /api/auth/*
func AuthRouter(router *mux.Router, logger *logger.Logger) {
	ctrl := &controller.AuthController{}

	router.Handle("/", gziphandler.GzipHandler(logger.Log(ctrl.CheckCredentials))).Methods("GET")
	router.Handle("/vk", gziphandler.GzipHandler(logger.Log(ctrl.LoginVkUser))).Methods("GET")
	router.Handle("/guest", gziphandler.GzipHandler(logger.Log(ctrl.GetGuestUser))).Methods("GET")
}
