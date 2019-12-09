package model

import "time"

type User struct {
	ID        uint `gorm:"primary_key;AUTO_INCREMENT;omitempty"`
	Role      string
	Token     string
	FirstName string
	LastName  string
	Photo     string
	Password  string `json:"-"`
	Version   int    `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
