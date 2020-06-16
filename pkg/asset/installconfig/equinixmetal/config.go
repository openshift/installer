package equinixmetal

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var defaultEquinixMetalConfigEnvVar = "EQUINIXMETAL_CONFIG"

// TODO(displague) what is the preferred config for Equinix Metal projects? support
// both yaml and json?
var defaultEquinixMetalConfigPath = filepath.Join(os.Getenv("HOME"), ".equinixmetal-config.yaml")

// Config holds Equinix Metal api access details
type Config struct {
	// APIURL is the Base URL for accessing the Equinix Metal API (https://api.equinix.com/metal/v1)
	APIURL string `json:"api_url,omitempty"`

	// APIKey is the User or Project API Key used to authenticate requests to the Equinix Metal API
	APIKey string `json:"api_key,omitempty"`
}

// LoadEquinixMetalConfig from the following location (first wins):
// 1. EQUINIXMETAL_CONFIG env variable
// 2. $defaultEquinixMetalConfigPath
// See #@Config for the expected format
func LoadEquinixMetalConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(discoverPath())
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NewConfig will return an Config by loading
// the configuration from locations specified in @LoadEquinixMetalConfig
func NewConfig() (Config, error) {
	c := Config{}
	in, err := LoadEquinixMetalConfig()
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(in, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func discoverPath() string {
	path, _ := os.LookupEnv(defaultEquinixMetalConfigEnvVar)
	if path != "" {
		return path
	}

	return defaultEquinixMetalConfigPath
}

// Save will serialize the config back into the locations
// specified in @LoadEquinixMetalConfig, first location with a file, wins.
func (c *Config) Save() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	path := discoverPath()
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, 0600)
}
