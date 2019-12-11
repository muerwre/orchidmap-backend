package controller

import (
	"math"
	"net/http"
	"strings"
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

type BetweenRange struct {
	Min float64 `form:"min"`
	Max float64 `form:"max"`
}

type RouteShallow struct {
	Address     string  `json:"address" sql:"address"`
	Distance    float64 `json:"distance" sql:"distance"`
	Title       string  `json:"title" sql:"title"`
	IsPublished bool    `json:"is_published" sql:"is_published"`
	IsPublic    bool    `json:"is_public" sql:"is_public"`
}

type LimitRange struct {
	Min   float64 `gorm:"column:min" sql:"min" json:"min"`
	Max   float64 `gorm:"column:max" sql:"max" json:"max"`
	Count int     `gorm:"column:count" sql:"count" json:"count"`
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
		route.IsPublished = false
	}

	route.CleanForPost()

	if exist.ID != 0 {
		d.Model(&route).Updates(route)
	} else {
		d.Create(&route)
	}

	c.JSON(http.StatusBadRequest, gin.H{"route": route, "exist": exist.ID != 0})
}

func (a *RouteController) PatchRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	address := c.PostForm("address")
	title := strings.Trim(c.PostForm("title"), "")
	public := strings.Trim(c.PostForm("is_public"), "") == "true"

	route := &model.Route{}

	d.Where("address = ?", address).First(&route)

	if route.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	if !route.CanBeEditedBy(u) {
		c.JSON(http.StatusConflict, gin.H{"error": "Not an owner", "code": "not_an_owner"})
		return
	}

	if len(title) > 100 {
		title = title[:100]
	}

	d.Model(&route).Update(map[string]interface{}{"title": title, "is_public": public})

	c.JSON(http.StatusBadRequest, gin.H{"route": route})
}

func (a *RouteController) DeleteRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	address := c.PostForm("address")

	route := &model.Route{}

	d.Where("address = ?", address).First(&route)

	if route.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	if !route.CanBeEditedBy(u) {
		c.JSON(http.StatusConflict, gin.H{"error": "Not an owner", "code": "not_an_owner"})
		return
	}

	d.Model(&route).Update(map[string]interface{}{"deleted_at": time.Now().UTC().Truncate(time.Second)})

	c.JSON(http.StatusBadRequest, gin.H{"route": route})
}

func (a *RouteController) PublishRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	address := c.PostForm("address")
	published := c.PostForm("published")

	route := &model.Route{}

	d.Where("address = ?", address).First(&route)

	if route.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	if u.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can publish routes", "code": "insufficient_rights"})
		return
	}

	d.Model(&route).Update(map[string]interface{}{"is_published": published == "true"})

	c.JSON(http.StatusBadRequest, gin.H{"route": route})
}

func (a *RouteController) GetAllRoutes(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	tab := c.Param("tab")
	between := &BetweenRange{}

	if (tab != "my") &&
		tab != "all" &&
		tab != "starred" {
		tab = "all"
	}

	if tab == "my" && u.ID == 0 {
		c.JSON(
			http.StatusOK,
			gin.H{"tab": tab, "routes": &[]RouteShallow{}, "limits": &LimitRange{Min: 0, Max: 0}, "between": between},
		)
		return
	}

	routes := &[]RouteShallow{}

	err := c.ShouldBindQuery(&between)

	if err != nil {
		between.Min, between.Max = 0, 0
	} else {
		c.BindQuery(&between)
	}

	if between.Max >= 200 || between.Max <= 0 {
		between.Max = 1e4
	}

	if between.Min > between.Max || between.Min <= 0 {
		between.Min = 0
	}

	q := d.Model(&routes).Where("distance >= ? AND distance <= ?", between.Min, between.Max)

	if tab == "starred" {
		q = q.Where("is_public = ? AND is_published = ?", true, true)
	}

	if tab == "my" {
		q = q.Where("user_id = ?", u.ID)
	}

	limits := &LimitRange{}

	q.Select("min(distance) as min, max(distance) as max, count(*) as count").First(&model.Route{}).Scan(&limits)
	q.Find(&[]model.Route{}).Offset(0).Limit(20).Scan(&routes)

	limits.Min = math.Floor((limits.Min / 25)) * 25
	limits.Max = math.Ceil((limits.Max / 25)) * 25

	if limits.Max <= 0 || len(*routes) == 0 {
		limits.Min = 0
	} else if limits.Min == limits.Max {
		limits.Min = limits.Max - 25
	} else if limits.Max > 200 {
		limits.Max = 200
	}

	c.JSON(http.StatusOK, gin.H{"tab": tab, "routes": routes, "limits": limits, "between": between})
}
