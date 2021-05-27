package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jdockerty/kindle-notes-grabber/config"
)

const (
	testEmail    string = "test.email123@test.com"
	testPassword string = "superSecret123"
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

	// Create a temporary file to write our config into the current directory.
	testFile, err := ioutil.TempFile(".", "test-config*.yaml")
	if err != nil {
		suite.T().Fatalf("Error creating temporary config file: %s", err.Error())
	}

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

func (suite *ConfigSuite) TestInvalidConfigPath() {
	
	_, err := config.New("some_fake_path/")
	assert.Error(suite.T(), err, "Expected error return from function, got %s", err)
}
func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
