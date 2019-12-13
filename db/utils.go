package db

import (
	"math/rand"

	"github.com/muerwre/orchidmap-backend/model"
)

var symbols = "ABCDEFGHIJKLMOPQRSTUVXYZabcdefghijgmlopqrstuvxyz01234567890"

// GenerateSequence creates random letter sequence
func GenerateSequence(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		r := rand.Intn(len(symbols))
		result = result + string(symbols[r])
	}

	return result
}

// GenerateGuestUser creates random account
func (d *DB) GenerateGuestUser() *model.User {
	token := "seq:" + string(GenerateSequence(64))

	for i := 0; i < 255; i++ {
		id := "guest:" + string(GenerateSequence(16))
		var c int
		d.Model(&model.User{}).Where("uid = ?", id).Count(&c)

		if c == 0 {
			return &model.User{Uid: id, Token: token, Role: "guest"}
		}
	}

	return nil
}

// GenerateRandomUrl creates random account
func (d *DB) GenerateRandomUrl() string {
	for i := 0; i < 255; i++ {
		id := string(GenerateSequence(24))
		var c int
		d.Model(&model.Route{}).Where("address = ?", id).Count(&c)

		if c == 0 {
			return id
		}
	}

	return ""
}
