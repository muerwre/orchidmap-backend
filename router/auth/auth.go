package auth

import (
	"encoding/json"
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

	if err != nil {
		user = ctx.DB.GenerateGuestUser()
		ctx.DB.Create(&user)
	}

	random_url := ctx.DB.GenerateRandomUrl()

	err = json.NewEncoder(w).Encode(struct {
		User      *model.User `json:"user"`
		RandomUrl string      `json:"random_url"`
	}{User: user, RandomUrl: random_url})

	if err != nil {
		return err
	}

	return nil
}
