package db

import (
	"errors"
	"fmt"

	"github.com/muerwre/orchidgo/model"
)

func (db *DB) AssumeUserExist(uid string, token string) (*model.User, error) {
	if uid == "" || token == "" {
		return nil, errors.New("Empty credentials providen")
	}

	user := &model.User{}

	db.Where("uid = ? AND token = ?", uid, token).Find(&user)

	fmt.Printf("%+v", user)

	if user.Role == "" {
		return nil, errors.New("Empty credentials providen")
	}

	return user, nil
}
