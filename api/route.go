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
	}

	// 	router.post('/star', star);
	// router.post('/', post);
	// router.get('/', get);
	// router.patch('/', patch);
	// router.delete('/', drop);
	// router.get('/list', list);

}
