package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/model"
	"github.com/muerwre/orchidgo/utils/logger"
	"github.com/muerwre/orchidgo/utils/vk"
	"golang.org/x/oauth2"
)

type AuthResponse struct {
	User      *model.User `json:"user"`
	RandomUrl string      `json:"random_url"`
	Error     string      `json:"error"`
	Success   bool        `json:"success"`
}

// Router creates new AuthRouter
func Router(router *mux.Router, logger *logger.Logger) {
	router.Handle("/", gziphandler.GzipHandler(logger.Log(CheckCredentials))).Methods("GET")
	router.Handle("/vk", gziphandler.GzipHandler(logger.Log(LoginVkUser))).Methods("GET")
	router.Handle("/guest", gziphandler.GzipHandler(logger.Log(GetGuestUser))).Methods("GET")
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

	err = json.NewEncoder(w).Encode(AuthResponse{User: user, RandomUrl: random_url, Error: error, Success: error == ""})

	if err != nil {
		return err
	}

	return nil
}

func GetGuestUser(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	user := ctx.DB.GenerateGuestUser()
	random_url := ctx.DB.GenerateRandomUrl()

	ctx.DB.Create(&user)

	if user == nil || random_url == "" {
		return errors.New("Failed to create reandom sequence")
	}

	err := json.NewEncoder(w).Encode(AuthResponse{User: user, RandomUrl: random_url, Success: true})

	if err != nil {
		return err
	}

	return nil
}

func LoginVkUser(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	context := context.Background()
	config := &oauth2.Config{
		ClientID:     ctx.Config.VkClientId,
		ClientSecret: ctx.Config.VkClientSecret,
		Scopes:       []string{},
		Endpoint:     vk.Endpoint,
		RedirectURL:  "http://localhost:7777/api/auth/vk",
	}

	code := r.URL.Query()["code"][0]

	if code == "" {
		return errors.New("Code is incorrect")
	}

	token, err := config.Exchange(context, code)

	if err != nil {
		return err
	}

	url := fmt.Sprintf(`https://api.vk.com/method/users.get?user_id=%s&fields=photo&v=5.67&access_token=%s`, fmt.Sprintf("%v", token.Extra("user_id")), token.AccessToken)

	response, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed read response: %s", err.Error())
	}

	var data vk.VkApiResponse

	json.Unmarshal(contents, &data)

	if data.Response == nil {
		return errors.New("Can't get user")
	}

	user, err := ctx.DB.FindOrCreateUser(
		&model.User{
			Uid:   fmt.Sprintf("vk:%d", data.Response[0].Id),
			Name:  fmt.Sprintf("%s %s", data.Response[0].FirstName, data.Response[0].LastName),
			Photo: fmt.Sprintf("%s", data.Response[0].Photo),
			Role:  "vk",
		},
	)

	if err != nil {
		return errors.New("Can't get user")
	}

	random_url := ctx.DB.GenerateRandomUrl()

	err = json.NewEncoder(w).Encode(
		AuthResponse{User: user, RandomUrl: random_url, Success: true},
	)

	if err != nil {
		return err
	}

	return nil
}
