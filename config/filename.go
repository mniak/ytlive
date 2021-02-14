package config

import (
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

var configFilenameCache string

func getConfigFilename() (string, error) {
	if configFilenameCache != "" {
		return configFilenameCache, nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "could not find home directory")
	}

	configpath := path.Join(home, ".ytlive")
	return configpath, nil
}
