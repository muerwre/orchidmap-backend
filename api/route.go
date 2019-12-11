package api

import (
	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidgo/controller"
)

// RouteRouter for /api/route/*
func RouteRouter(r *gin.RouterGroup, a *API) {
	r.GET("/", controller.Route.GetRoute)

	restricted := r.Group("/").Use(a.AuthRequired)
	{
		restricted.POST("/", controller.Route.SaveRoute)
		restricted.PATCH("/", controller.Route.PatchRoute)
		restricted.DELETE("/", controller.Route.DeleteRoute)
		restricted.POST("/publish", controller.Route.PublishRoute)
	}

	optional := r.Group("/").Use(a.AuthOptional)
	{
		optional.GET("/all/:tab", controller.Route.GetAllRoutes)
	}
	// router.get('/list', list);
}
