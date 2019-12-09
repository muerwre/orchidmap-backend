// Package vk provides constants for using OAuth2 to access VK.com.
package vk // import "golang.org/x/oauth2/vk"

import (
	"golang.org/x/oauth2"
)

// Endpoint is VK's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://oauth.vk.com/authorize",
	TokenURL: "https://oauth.vk.com/access_token",
}

type VkApiResponse struct {
	Response []struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"photo"`
	} `json:"response"`
}
