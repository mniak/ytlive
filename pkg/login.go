package pkg

import (
	"fmt"
	"net/http"

	"github.com/mniak/oauth2device"
	"github.com/mniak/oauth2device/googledevice"
	"github.com/mniak/ytlive/config"
	"github.com/mniak/ytlive/internal"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func Login() error {

	config := internal.GetOAuthConfig()

	_, err := authenticate(config)
	if err != nil {
		return errors.Wrap(err, "could not authenticate")
	}

	return nil
}

func authenticate(config oauth2.Config) (token *oauth2.Token, err error) {
	deviceConfig := &oauth2device.Config{
		Config:         &config,
		DeviceEndpoint: googledevice.DeviceEndpoint,
	}
	httpClient := internal.AddLoggingTransportIfNeeded(http.DefaultClient)
	codeReq, err := oauth2device.RequestDeviceCode(httpClient, deviceConfig)
	if err != nil {
		err = errors.Wrap(err, "failed to request device code")
		return
	}

	fmt.Printf("[Google Authentication] Navigate to %v type the following code: %v\n", codeReq.VerificationURL, codeReq.UserCode)

	token, err = oauth2device.WaitForDeviceAuthorization(httpClient, deviceConfig, codeReq)
	if err != nil {
		err = errors.Wrap(err, "could not get token")
		return
	}
	return token, saveToken(token)
}

func saveToken(token *oauth2.Token) error {
	config.Root.Token = *token
	err := config.Save()
	if err != nil {
		return errors.Wrap(err, "could not save token")
	}
	return nil
}
