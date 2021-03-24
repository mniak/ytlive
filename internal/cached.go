package internal

import (
	"context"
	"log"
	"time"

	"github.com/mniak/ytlive/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type CachedTokenSource struct {
	config oauth2.Config
	ctx    context.Context
}

func NewCachedTokenSource(ctx context.Context, config oauth2.Config) CachedTokenSource {
	return CachedTokenSource{
		ctx:    ctx,
		config: config,
	}
}

func (ts CachedTokenSource) saveToken(token *oauth2.Token) error {
	config.Root.Token = *token

	err := config.Save()
	if err != nil {
		return err
	}

	return nil
}

func (ts CachedTokenSource) Token() (*oauth2.Token, error) {

	if config.Root.Token.AccessToken != "" && config.Root.Token.Expiry.After(time.Now()) {
		return &config.Root.Token, nil
	}

	if config.Root.Token.RefreshToken == "" {
		return nil, errors.New("could not refresh the token because the client secret is empty on config")
	}
	if config.Root.Application.ClientID == "" {
		return nil, errors.New("could not refresh the token because the client id is empty on config")
	}
	if config.Root.Application.ClientSecret == "" {
		return nil, errors.New("could not refresh the token because the client secret is empty on config")
	}

	token, err := ts.config.TokenSource(ts.ctx, &config.Root.Token).Token()
	if err != nil {
		return token, err
	}

	err = ts.saveToken(token)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to write token in cache"))
	}
	return token, nil
}
