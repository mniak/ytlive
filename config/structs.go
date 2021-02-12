package config

import "golang.org/x/oauth2"

type ConfigurationRoot struct {
	Application struct {
		ClientID     string
		ClientSecret string
	}
	Token oauth2.Token
}
