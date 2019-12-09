package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/model"
	"github.com/muerwre/orchidgo/utils/logger"
)

// Router creates new AuthRouter
func Router(router *mux.Router, logger *logger.Logger) {
	router.Handle("/", gziphandler.GzipHandler(logger.Log(CheckCredentials))).Methods("GET")
}

// CheckCredentials checks id and token and returns guest token if they're incorrect
func CheckCredentials(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	user, err := ctx.DB.AssumeUserExist(r.URL.Query()["id"][0], r.URL.Query()["token"][0])
	error := ""

	if err != nil {
		user = ctx.DB.GenerateGuestUser()
		error = "User not found, falling back to guest"
		ctx.DB.Create(&user)
	}

	random_url := ctx.DB.GenerateRandomUrl()

	if user == nil || random_url == "" {
		return errors.New("Failed to create reandom sequence")
	}

	err = json.NewEncoder(w).Encode(struct {
		*model.User
		RandomUrl string `json:"random_url"`
		Error     string `json:"error"`
	}{User: user, RandomUrl: random_url, Error: error})

	if err != nil {
		return err
	}

	return nil
}
