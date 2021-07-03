package config

import (
	"fmt"
	"os"
	"strings"

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

// IMAPServer is a struct for storing the relevant information for a
// providers IMAP server for accessing a mailbox.
type IMAPServer struct {
	ServiceName string
	Address     string
	Port        int
}

var serviceNameToIMAPServer map[string]string = map[string]string{
	"gmail":   "imap.gmail.com",
	"outlook": "imap-mail.outlook.com",
	"yahoo":   "imap.mail.yahoo.com",
	"aol":     "imap.aol.com",
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
		return nil, fmt.Errorf("a 'kindle-notes' directory does not exist at '%s' to write the completed notebooks save file", homeDir)
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

// Populate is used to fill in the relevant information for an IMAP server given a simple service
// name, regardless of case, such as 'gmail' or 'outlook'. Note, that there are limited service mappings at this time.
func (im *IMAPServer) Populate(serviceName string) {
	sanitisedServiceName := strings.ToLower(serviceName)

	imapServer := serviceNameToIMAPServer[sanitisedServiceName]

	im.ServiceName = sanitisedServiceName
	im.Address = imapServer
	im.Port = 993
}
