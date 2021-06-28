package config

import (
	"fmt"
	"os"

	"github.com/jdockerty/kindle-notes-grabber/notes"
	"github.com/spf13/viper"
)

// Config is a struct for holding the exported configuration of the program,
// this provides access to the variables that are required to log into an
// email account.
type Config struct {
	Email    string `mapstructure:"email" env:"KNG_EMAIL"`
	Password string `mapstructure:"password" env:"KNG_PASSWORD"`
}

// New returns a Config struct with the relevant values populated. This leverages
// a configuration file by the name of 'kng-config.yaml', at a specified path, or by
// the two environment variables: KNG_EMAIL and KNG_PASSWORD. Setting either of these
// will satisfy the unmarshaling requirements into the struct.
func New(path string) (*Config, error) {

	var cfg Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dirPath := fmt.Sprintf("%s/%s", homeDir, "kindle-notes")
	dirExists, err := notes.Exists(dirPath)
	if err != nil {
		return nil, err
	}

	if !dirExists {
		return nil, fmt.Errorf("a 'kindle-notes' directory does not at '%s' to write the completed notebooks save file", homeDir)
	}

	viper.SetEnvPrefix("KNG")
	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.AddConfigPath(path)
	viper.AddConfigPath(homeDir)

	viper.SetConfigName("kng-config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; Attempt to read from environemnt
			viper.Set("EMAIL", viper.Get("KNG_EMAIL"))
			viper.Set("PASSWORD", viper.Get("KNG_PASSWORD"))

		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
