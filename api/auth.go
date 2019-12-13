package api

import (
	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidmap-backend/controller"
)

// AuthRouter for /api/auth/*
func AuthRouter(r *gin.RouterGroup, a *API) {
	r.GET("/", controller.Auth.CheckCredentials)
	r.GET("/vk", controller.Auth.LoginVkUser)
	r.GET("/guest", controller.Auth.GetGuestUser)
}
