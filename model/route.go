package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type Point struct {
	Lat float64
	Lng float64
}

type Sticker struct {
	Angle   float64
	Latlng  Point
	Sticker string
	Set     string
	Text    string
}

type PointArray []Point
type StickerArray []Sticker

type Route struct {
	*gorm.Model

	Address     string
	Title       string
	Version     int          `json:"-"`
	Route       PointArray   `sql:"route" gorm:"name:route;type:longtext"`
	Stickers    StickerArray `gorm:"name:stickers;type:longtext"`
	Distance    float64
	IsPublic    bool
	IsStarred   bool
	IsDeleted   bool
	Logo        string
	Provider    string
	Description string
	User        User `gorm:"foreignkey:UserId"`
	UserId      uint
}

func (s *PointArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s)
}

func (s PointArray) Value() (driver.Value, error) {
	val, err := json.Marshal(s)
	return string(val), err
}

func (s *StickerArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s)
}

func (s StickerArray) Value() (driver.Value, error) {
	val, err := json.Marshal(s)
	return string(val), err
}
