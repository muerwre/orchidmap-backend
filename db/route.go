package db

import (
	"errors"

	"github.com/muerwre/orchidmap-backend/model"
)

func (db *DB) FindRouteByAddress(address string) (*model.Route, error) {
	if address == "" {
		return nil, errors.New("Name is empty")
	}

	route := &model.Route{}
	db.Where("address = ?", address).First(&route)

	if route.ID == 0 {
		return nil, errors.New("Route not found")
	}

	return route, nil
}
