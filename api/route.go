package api

import (
	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidmap-backend/controller"
)

// RouteRouter for /api/route/*
func RouteRouter(r *gin.RouterGroup, a *API) {
	restricted := r.Group("/").Use(a.AuthRequired)
	optional := r.Group("/").Use(a.AuthOptional)

	r.GET("/", controller.Route.GetRoute)

	{
		optional.GET("/list/:tab", controller.Route.GetAllRoutes)
	}

	{
		restricted.POST("/", controller.Route.SaveRoute)
		restricted.PATCH("/", controller.Route.PatchRoute)
		restricted.DELETE("/", controller.Route.DeleteRoute)
		restricted.POST("/publish", controller.Route.PublishRoute)
	}

	// router.get('/list', list);
}
