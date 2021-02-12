package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mniak/oauth2device"
	"github.com/mniak/oauth2device/googledevice"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/yaml.v2"
)

type GoogleAuthTokenSource struct {
	ctx    context.Context
	Config oauth2.Config
}

func NewGoogleTokenSource(ctx context.Context) GoogleAuthTokenSource {
	return GoogleAuthTokenSource{
		ctx: ctx,
		Config: oauth2.Config{
			ClientID:     viper.GetString("Youtube.ClientID"),
			ClientSecret: viper.GetString("Youtube.ClientSecret"),
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost",
			Scopes:       []string{"https://www.googleapis.com/auth/youtube"},
		},
	}
}

func (ts GoogleAuthTokenSource) Token() (token *oauth2.Token, err error) {
	deviceConfig := &oauth2device.Config{
		Config:         &ts.Config,
		DeviceEndpoint: googledevice.DeviceEndpoint,
	}
	httpClient := AddLoggingTransportIfNeeded(http.DefaultClient)
	codeReq, err := oauth2device.RequestDeviceCode(httpClient, deviceConfig)
	if err != nil {
		return
	}

	fmt.Printf("[Google Authentication] Navigate to %v type the following code: %v\n", codeReq.VerificationURL, codeReq.UserCode)

	token, err = oauth2device.WaitForDeviceAuthorization(httpClient, deviceConfig, codeReq)
	return
}

func GetConfig(clientId, clientSecret string) oauth2.Config {
	config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
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
