package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jdockerty/kindle-notes-grabber/config"
)

const (
	testEmail    string = "test.email123@test.com"
	testPassword string = "superSecret123"
	emailEnv     string = "KNG_EMAIL"
	passwordEnv  string = "KNG_PASSWORD"
)

type ConfigSuite struct {
	suite.Suite
	Conf *config.Config
}

func (suite *ConfigSuite) SetupTest() {

	var exampleConfig = []byte(
		`---
email: test.email123@test.com
password: superSecret123
`)

	homeDir, err := os.UserHomeDir()
	suite.Assertions.Nil(err)
	setupDir := fmt.Sprintf("%s/kindle-notes", homeDir)
	os.Create(setupDir)

	// Create a temporary file to write our config into the current directory.
	testFile, err := os.Create("kng-config.yaml")
	if err != nil {
		suite.T().Fatalf("Error creating temporary config file: %s", err.Error())
	}

	// Unset the environment so that we rely on the test config file
	os.Unsetenv(emailEnv)
	os.Unsetenv(passwordEnv)

	defer testFile.Close()
	defer os.Remove(testFile.Name())

	testFile.Write(exampleConfig)
	suite.Conf, err = config.New(testFile.Name())
	assert.Nil(suite.T(), err, "Temporary file for setup should be at a valid path.")
}

func (suite *ConfigSuite) TestConfig() {
	assert := assert.New(suite.T())

	assert.Equal(testEmail, suite.Conf.Email)
	assert.Equal(testPassword, suite.Conf.Password)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func TestProviderMappingIsCorrect(t *testing.T) {
	assert := assert.New(t)

	gmailVariants := []string{
		"gmail",
		"GMAIL",
		"GMaIl",
	}

	outlookVariants := []string{
		"outlook",
		"OUTLOOK",
		"OutLOok",
	}
	
	
	for _, variant := range gmailVariants {
		var ims config.IMAPServer

		ims.Populate(variant)

		assert.Equal("gmail", ims.ServiceName)
		assert.Equal(993, ims.Port)
		assert.Equal("imap.gmail.com", ims.Address)
	}

	for _, variant := range outlookVariants {
		var ims config.IMAPServer

		ims.Populate(variant)

		assert.Equal("outlook", ims.ServiceName)
		assert.Equal(993, ims.Port)
		assert.Equal("imap-mail.outlook.com", ims.Address)
	}
}
