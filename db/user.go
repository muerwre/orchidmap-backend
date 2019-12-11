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

func (db *DB) FindOrCreateUser(u *model.User) (*model.User, error) {
	if u.Uid == "" {
		return nil, errors.New("User id is not set")
	}

	user := &model.User{}

	db.Where("uid = ?", u.Uid).Find(&user)

	if user.Uid == "" {
		user = u
		user.Token = fmt.Sprintf("seq:%s", GenerateSequence(32))
		db.Create(&user)
	}

	return user, nil
}

func (db *DB) GetUserByToken(token string) (*model.User, error) {
	if token == "" {
		return nil, errors.New("Credentials are empty")
	}

	user := &model.User{}

	db.Where("token = ?", token).Find(&user).First(&user)

	if user.ID == 0 {
		return nil, errors.New("User not found")
	}

	return user, nil
}
