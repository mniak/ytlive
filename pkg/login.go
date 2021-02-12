package pkg

import (
	"fmt"
	"net/http"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mniak/oauth2device"
	"github.com/mniak/oauth2device/googledevice"
	"github.com/mniak/ytlive/internal"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func Login(clientId, clientSecret string) error {

	config := internal.GetConfig(clientId, clientSecret)

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
		return
	}

	fmt.Printf("[Google Authentication] Navigate to %v type the following code: %v\n", codeReq.VerificationURL, codeReq.UserCode)

	token, err = oauth2device.WaitForDeviceAuthorization(httpClient, deviceConfig, codeReq)
	return token, saveToken(token)
}

func saveToken(token *oauth2.Token) error {
	tokens := make([]oauth2.Token, 0)
	tokens = append(tokens, *token)
	viper.Set("Tokens", tokens)
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	return viper.WriteConfigAs(path.Join(home, ".ytlive.toml"))
}
