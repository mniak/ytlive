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
	if err := authenticate(config); err != nil {
		return errors.Wrap(err, "could not authenticate")
	}

	return nil
}

func authenticate(oauthConfig oauth2.Config) error {
	deviceConfig := &oauth2device.Config{
		Config:         &oauthConfig,
		DeviceEndpoint: googledevice.DeviceEndpoint,
	}
	httpClient := http.DefaultClient
	codeReq, err := oauth2device.RequestDeviceCode(httpClient, deviceConfig)
	if err != nil {
		return errors.Wrap(err, "failed to request device code")
	}

	fmt.Printf("[Google Authentication] Navigate to %v type the following code: %v\n", codeReq.VerificationURL, codeReq.UserCode)

	token, err := oauth2device.WaitForDeviceAuthorization(httpClient, deviceConfig, codeReq)
	if err != nil {
		return errors.Wrap(err, "could not get token")
	}
	config.Root.Token = *token
	return config.Save()
}
