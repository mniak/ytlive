package config

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

var Root *ConfigurationRoot = new(ConfigurationRoot)

func Load() (*ConfigurationRoot, error) {
	filename, err := getConfigFilename()
	if err != nil {
		return Root, errors.Wrap(err, "could not build configuration name")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Root, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return Root, errors.Wrap(err, "could not open configuration file")
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	err = decoder.Decode(Root)
	return Root, err
}

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

func Save() error {
	filename, err := getConfigFilename()
	if err != nil {
		return errors.Wrap(err, "could not build configuration name")
	}

	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "could not create configuration file")
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)

	err = encoder.Encode(Root)
	return errors.Wrap(err, "could not write to configuration file")
}
