package youtube

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mniak/oauth2device"
	"github.com/mniak/oauth2device/googledevice"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	httpClient := addLoggingTransportIfNeeded(http.DefaultClient)
	codeReq, err := oauth2device.RequestDeviceCode(httpClient, deviceConfig)
	if err != nil {
		return
	}

	fmt.Printf("[Google Authentication] Navigate to %v type the following code: %v\n", codeReq.VerificationURL, codeReq.UserCode)

	token, err = oauth2device.WaitForDeviceAuthorization(httpClient, deviceConfig, codeReq)
	return
}
