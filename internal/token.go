package internal

import (
	"context"

	"github.com/mniak/ytlive/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetOAuthConfig() oauth2.Config {
	config := oauth2.Config{
		ClientID:     config.Root.Application.ClientID,
		ClientSecret: config.Root.Application.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://www.googleapis.com/auth/youtube"},
	}
	return config
}

func GetTokenSource(config oauth2.Config) (context.Context, CachedTokenSource) {
	ctx := context.Background()
	tokenSource := NewCachedTokenSource(ctx, config)
	return ctx, tokenSource
}
