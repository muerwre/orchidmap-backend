package model

import (
	"time"
)

type User struct {
	Uid     string `json:"uid"`
	Role    string `json:"role" gorm:"default:'guest'"`
	Token   string `json:"token"`
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	Version int    `json:"-"`

	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}
