package model

import "github.com/jinzhu/gorm"

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

type Route struct {
	*gorm.Model

	Address     string
	Title       string
	Version     int        `json:"-"`
	RawRoute    string     `json:"-" sql:"route" gorm:"name:route;type:longtext"`
	Route       []*Point   `sql:"-"`
	RawStickers string     `json:"-" sql:"stickers" gorm:"name:stickers;type:longtext"`
	Stickers    []*Sticker `sql:"-"`
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
