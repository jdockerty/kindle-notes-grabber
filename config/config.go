package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Email    string `yaml:"email" env:"KNG_EMAIL"`
	Password string `yaml:"password" env:"KNG_PASSWORD"`
}

// New uses the 'cleanenv' package to read from the relevant configuration file,
// acting as a wrapper for the error and returning a correctly configured struct.
func New(path string) (*Config, error) {

	// Return an error when the provided path is not valid.
	if isValidPath := fs.ValidPath(path); !isValidPath {
		errorMsg := fmt.Sprintf("Invalid path provided: %s", path)
		return nil, errors.New(errorMsg)
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("Configuration is not set: %s", err)
	}

	return &cfg, nil
}
