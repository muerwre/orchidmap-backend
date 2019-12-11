package model

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"strings"
	"time"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Sticker struct {
	Angle   float64 `json:"angle"`
	Latlng  Point   `json:"latlng"`
	Sticker string  `json:"sticker"`
	Set     string  `json:"set"`
	Text    string  `json:"text"`
}

type PointArray []Point
type StickerArray []Sticker

type Route struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Address     string       `json:"address"`
	Title       string       `json:"title"`
	Version     int          `json:"-"`
	Route       PointArray   `sql:"route" gorm:"name:route;type:longtext" json:"route"`
	Stickers    StickerArray `gorm:"name:stickers;type:longtext" json:"stickers"`
	Distance    float64      `json:"distance"`
	IsPublic    bool         `json:"is_public"`
	IsPublished bool         `json:"is_published"`
	IsDeleted   bool         `json:"-"`
	Logo        string       `json:"logo"`
	Provider    string       `json:"provider"`
	Description string       `json:"description"`
	User        User         `gorm:"foreignkey:UserId" json:"-"`
	UserId      uint         `json:"-"`
}

func (p *PointArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &p)
}

func (p PointArray) Value() (driver.Value, error) {
	val, err := json.Marshal(p)
	return string(val), err
}

func (s *StickerArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s)
}

func (s StickerArray) Value() (driver.Value, error) {
	val, err := json.Marshal(s)
	return string(val), err
}

func (s *StickerArray) CleanForPost() {
	out := &StickerArray{}

	for _, b := range *s {
		if b.Latlng.Lat != 0 &&
			b.Latlng.Lng != 0 &&
			b.Sticker != "" &&
			b.Set != "" &&
			len(b.Text) <= 256 &&
			b.Angle >= -math.Pi &&
			b.Angle <= math.Pi {
			*out = append(*out, b)
		}
	}

	*s = *out
}

func (p *PointArray) CleanForPost() {
	out := &PointArray{}

	for _, b := range *p {
		if b.Lat != 0 && b.Lng != 0 {
			*out = append(*out, b)
		}
	}

	*p = *out
}

func (r *Route) CleanForPost() {
	r.Stickers.CleanForPost()
	r.Route.CleanForPost()

	if len(r.Title) > 100 {
		r.Title = r.Title[:100]
	}

	if len(r.Address) > 64 {
		r.Address = r.Title[:64]
	}

	if len(r.Description) > 256 {
		r.Description = r.Description[:256]
	}

	if len(r.Provider) > 16 {
		r.Provider = r.Provider[:16]
	}

	if len(r.Logo) > 16 {
		r.Logo = r.Logo[:16]
	}

	res := &Route{
		CreatedAt:   r.CreatedAt,
		Stickers:    r.Stickers,
		Route:       r.Route,
		Title:       strings.Trim(r.Title, ""),
		Description: r.Description,
		Distance:    r.Distance,
		Provider:    r.Provider,
		Logo:        r.Logo,
		Address:     r.Address,
		ID:          r.ID,
		User:        r.User,
	}

	*r = *res
}

func (r *Route) CanBeEditedBy(u *User) bool {
	return r.ID == 0 || r.UserId == u.ID || u.Role == "admin"
}
