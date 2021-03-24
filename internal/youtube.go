package internal

import (
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CreateYoutubeClient() (*youtube.Service, error) {
	config := GetOAuthConfig()
	ctx, tokenSource := GetTokenSource(config)

	svc, err := youtube.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, errors.Wrap(err, "could not create Youtube API client")
	}
	return svc, nil
}
