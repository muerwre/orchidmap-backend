package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidgo/db"
	"github.com/muerwre/orchidgo/model"
)

type RouteController struct{}

type RoutePostInput struct {
	Route *model.Route `json:"route"`
	Force bool         `json:"force"`
}

var Route = &RouteController{}

func (a *RouteController) GetRoute(c *gin.Context) {
	address := c.Query("name")
	d := c.MustGet("DB").(*db.DB)

	if address == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Name is undefined"})
		return
	}

	route, err := d.FindRouteByAddress(address)

	if err != nil || route == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (a *RouteController) SaveRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	var post RoutePostInput

	err := c.BindJSON(&post)

	route := post.Route
	force := post.Force

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exist := &model.Route{}

	d.Where("address = ?", route.Address).First(&exist)

	if !exist.CanBeEditedBy(u) {
		c.JSON(http.StatusConflict, gin.H{"error": "Not an owner", "code": "conflict"})
		return
	}

	if exist.ID != 0 && !force {
		c.JSON(http.StatusConflict, gin.H{"error": "Overwrite confirmation needed", "code": "already_exist"})
		return
	}

	if exist.ID != 0 {
		route.ID = exist.ID
	} else {
		route.CreatedAt = time.Now().UTC().Truncate(time.Second)
		route.User = *u
		route.IsStarred = false
	}

	route.CleanForPost()

	if exist.ID != 0 {
		d.Model(&route).Updates(route)
	} else {
		d.Create(&route)
	}

	c.JSON(http.StatusBadRequest, gin.H{"route": route, "exist": exist.ID != 0})
}
