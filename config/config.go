package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Email    string `mapstructure:"email" env:"KNG_EMAIL"`
	Password string `mapstructure:"password" env:"KNG_PASSWORD"`
}

// New uses the 'cleanenv' package to read from the relevant configuration file,
// acting as a wrapper for the error and returning a correctly configured struct.
func New(path string) (*Config, error) {

	var cfg Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
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
