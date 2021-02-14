package config

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

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
