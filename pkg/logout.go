package pkg

import (
	"github.com/mniak/ytlive/config"
	"golang.org/x/oauth2"
)

func Logout() error {
	config.Root.Token = oauth2.Token{}
	return config.Save()
}
