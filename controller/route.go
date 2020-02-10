package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidmap-backend/db"
	"github.com/muerwre/orchidmap-backend/model"
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

	url := d.GenerateRandomUrl()

	c.JSON(http.StatusOK, gin.H{"route": route, "random_url": url})
}

func (a *RouteController) GetRandomRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	r := &model.Route{}
	min, err := strconv.Atoi(c.Query("min"))

	if err != nil {
		min = 0
	}

	max, err := strconv.Atoi(c.Query("max"))

	if err != nil {
		max = 0
	}

	q := d.Where("is_public = ? AND is_published = ?", true, true).Order("RAND()")

	if min > 0 && max >= 0 && max < 400 && max > min {
		q = q.Where("distance >= ? AND distance <= ?", float32(min)*0.8, float32(max)*1.2)
	}

	if min == 0 && max >= 0 && max < 400 && max > min {
		q = q.Where("distance >= ? AND distance <= ?", float32(max)*0.7, float32(max)*1.3)
	}

	q.First(&r)

	if r.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
	}

	c.JSON(http.StatusOK, gin.H{"id": r.Address, "title": r.Title, "distance": r.Distance, "description": r.Description})
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
		c.JSON(http.StatusConflict, gin.H{
			"error": "Overwrite confirmation needed",
			"code":  "already_exist",
		})
		return
	}

	if exist.ID != 0 {
		route.ID = exist.ID
		route.User = *u
		route.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	} else {
		route.CreatedAt = time.Now().UTC().Truncate(time.Second)
		route.UpdatedAt = time.Now().UTC().Truncate(time.Second)
		route.User = *u
		route.IsPublished = false
	}

	route.CleanForPost()

	if exist.ID != 0 {
		d.Save(&route)
	} else {
		d.Create(&route)
	}

	c.JSON(http.StatusOK, gin.H{"route": route, "exist": exist.ID != 0})
}

func (a *RouteController) PatchRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	post := &struct {
		Address string
		Title   string
		Public  bool `json:"is_public"`
	}{}

	err := c.BindJSON(&post)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	route := &model.Route{}

	d.Where("address = ?", post.Address).First(&route)

	if route.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	if !route.CanBeEditedBy(u) {
		c.JSON(http.StatusConflict, gin.H{"error": "Not an owner", "code": "not_an_owner"})
		return
	}

	if len(post.Title) > 100 {
		post.Title = post.Title[:100]
	}

	d.Model(&route).Update(map[string]interface{}{"title": post.Title, "is_public": post.Public})

	c.JSON(http.StatusOK, gin.H{"route": route})
}

func (a *RouteController) DeleteRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	post := &struct {
		Address string
	}{}

	err := c.BindJSON(&post)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if post.Address == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	route := &model.Route{}

	d.Where("address = ?", post.Address).First(&route)

	if route.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	if !route.CanBeEditedBy(u) {
		c.JSON(http.StatusConflict, gin.H{"error": "Not an owner", "code": "not_an_owner"})
		return
	}

	d.Model(&route).Update(map[string]interface{}{"deleted_at": time.Now().UTC().Truncate(time.Second)})

	c.JSON(http.StatusOK, gin.H{"route": route})
}

func (a *RouteController) PublishRoute(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	post := &struct {
		Address   string
		Published bool `json:"is_published"`
	}{}

	err := c.BindJSON(&post)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if post.Address == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	route := &model.Route{}

	d.Where("address = ?", post.Address).First(&route)

	if route.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	if u.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can publish routes", "code": "insufficient_rights"})
		return
	}

	d.Model(&route).Update(map[string]interface{}{"is_published": post.Published})

	c.JSON(http.StatusOK, gin.H{"route": route})
}

func (a *RouteController) GetAllRoutes(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	u := c.MustGet("User").(*model.User)

	tab := c.Param("tab")
	filter := &model.FilterRange{}

	if (tab != "my") &&
		tab != "pending" &&
		tab != "starred" {
		tab = "all"
	}

	routes := &[]model.RouteShallow{}

	err := c.ShouldBindQuery(&filter)

	if err != nil {
		filter.Min, filter.Max = 0, 0
	} else {
		c.BindQuery(&filter)
	}

	if filter.Max >= 200 || filter.Max <= 0 {
		filter.Max = 1e4
	}

	if filter.Min > filter.Max || filter.Min <= 0 {
		filter.Min = 0
	}

	if (tab == "my" && u.ID == 0) || (filter.Search != "" && len(filter.Search) <= 3) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"tab":    tab,
				"routes": &[]model.RouteShallow{},
				"limits": &model.LimitRange{Min: 0, Max: 0},
				"filter": filter,
			},
		)
		return
	}

	q := d.Model(&routes)

	if tab == "pending" {
		q = q.Where("is_public = ? AND is_published = ?", true, false)
	}

	if tab == "starred" {
		q = q.Where("is_public = ? AND is_published = ?", true, true)
	}

	if tab == "my" {
		q = q.Where("user_id = ?", u.ID)
	}

	if filter.Search != "" && len(filter.Search) > 3 {
		q = q.Where("address RLIKE ? OR title RLIKE ?", filter.Search, filter.Search)
	}

	limits := &model.LimitRange{}

	q.Select("min(distance) as min, max(distance) as max, count(*) as count").
		First(&model.Route{}).
		Scan(&limits)

	q = q.Where("distance >= ? AND distance <= ?", filter.Min, filter.Max)

	q.Select("count(*) as count").
		First(&model.Route{}).
		Scan(&limits)

	q.Find(&[]model.Route{}).
		Offset(filter.Shift).
		Limit(filter.Step).
		Order("updated_at desc").
		Scan(&routes)

	limits.Normalize(len(*routes))

	c.JSON(http.StatusOK, gin.H{"tab": tab, "routes": routes, "limits": limits, "filter": filter})
}
