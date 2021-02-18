package config

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

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
