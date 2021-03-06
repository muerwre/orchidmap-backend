package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidmap-backend/app"
	"github.com/muerwre/orchidmap-backend/db"
	"github.com/muerwre/orchidmap-backend/model"
	"github.com/muerwre/orchidmap-backend/utils/vk"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type AuthResponse struct {
	User      *model.User `json:"user"`
	RandomUrl string      `json:"random_url"`
}

type AuthController struct{}

var Auth = AuthController{}

// CheckCredentials checks id and token and returns guest token if they're incorrect
func (a *AuthController) CheckCredentials(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)

	user, err := d.AssumeUserExist(c.Query("token"))
	status := http.StatusOK

	if err != nil {
		user = d.GenerateGuestUser()
		status = http.StatusCreated
		d.Create(&user)
	}

	random_url := d.GenerateRandomUrl()

	if user == nil || random_url == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to create reandom sequence"})
		return
	}

	c.JSON(status, AuthResponse{User: user, RandomUrl: random_url})
}

func (a *AuthController) GetGuestUser(c *gin.Context) {
	d := c.MustGet("DB").(*db.DB)
	user := d.GenerateGuestUser()
	random_url := d.GenerateRandomUrl()

	d.Create(&user)

	if user == nil || random_url == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to create random user"})
		return
	}

	c.JSON(200, AuthResponse{User: user, RandomUrl: random_url})
}

func (a *AuthController) LoginVkUser(c *gin.Context) {
	context := context.Background()
	cf := c.MustGet("Config").(*app.Config)
	d := c.MustGet("DB").(*db.DB)

	config := &oauth2.Config{
		ClientID:     cf.VkClientId,
		ClientSecret: cf.VkClientSecret,
		Scopes:       []string{},
		Endpoint:     vk.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/api/auth/vk", cf.VkCallbackUrl),
	}

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Code is incorrect"})
		return
	}

	token, err := config.Exchange(context, code)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to get token"})
		logrus.Infof("Failed to get token: %v", err.Error())
		return
	}

	url := fmt.Sprintf(
		`https://api.vk.com/method/users.get?user_id=%s&fields=photo&v=5.67&access_token=%s`,
		fmt.Sprintf("%v", int(token.Extra("user_id").(float64))),
		token.AccessToken,
	)

	response, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed getting user info"})
		logrus.Infof("Failed getting user info: %v", err.Error())
		return
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to read response"})
		return
	}

	var data vk.VkApiResponse

	err = json.Unmarshal(contents, &data)

	if data.Response == nil || err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Can't get user"})

		return
	}

	user, err := d.FindOrCreateUser(
		&model.User{
			Uid:   fmt.Sprintf("vk:%d", data.Response[0].Id),
			Name:  fmt.Sprintf("%s %s", data.Response[0].FirstName, data.Response[0].LastName),
			Photo: fmt.Sprintf("%v", data.Response[0].Photo),
			Role:  "vk",
		},
	)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Can't get user"})
		return
	}

	random_url := d.GenerateRandomUrl()

	c.HTML(http.StatusOK, "social.html", AuthResponse{User: user, RandomUrl: random_url})
}
