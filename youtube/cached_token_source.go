package youtube

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type CachedTokenSource struct {
	Inner  oauth2.TokenSource
	config oauth2.Config
	ctx    context.Context
}

func NewCachedTokenSource(ctx context.Context, inner oauth2.TokenSource, config oauth2.Config) CachedTokenSource {
	return CachedTokenSource{
		ctx:    ctx,
		Inner:  inner,
		config: config,
	}
}

const cacheFileName = ".youtube-token.cache"

func (ts CachedTokenSource) tryLoadToken() (*oauth2.Token, error) {
	token := &oauth2.Token{}

	file, err := os.Open(cacheFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (ts CachedTokenSource) saveToken(token *oauth2.Token) error {
	file, err := os.Create(cacheFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(token)
}

func (ts CachedTokenSource) Token() (*oauth2.Token, error) {

	token, err := ts.tryLoadToken()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to load token from cache"))
	}

	token, err = ts.config.TokenSource(ts.ctx, token).Token()
	if err != nil {
		return token, err
	}

	err = ts.saveToken(token)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to write token in cache"))
	}
	return token, nil
}
