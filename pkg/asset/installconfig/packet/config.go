package packet

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var defaultPacketConfigEnvVar = "PACKET_CONFIG"

// TODO(displague) what is the preferred config for Packet projects? support
// both yaml and json?
var defaultPacketConfigPath = filepath.Join(os.Getenv("HOME"), ".packet-config.yaml")

// Config holds Packet api access details
type Config struct {
	// APIURL is the Base URL for accessing the Packet API (https://api.packet.com/)
	APIURL string `json:"api_url,omitempty"`

	// APIKey is the User or Project API Key used to authenticate requests to the Packet API
	APIKey string `json:"api_key,omitempty"`
}

// LoadPacketConfig from the following location (first wins):
// 1. PACKET_CONFIG env variable
// 2. $defaultPacketConfigPath
// See #@Config for the expected format
func LoadPacketConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(discoverPath())
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NewConfig will return an Config by loading
// the configuration from locations specified in @LoadPacketConfig
func NewConfig() (Config, error) {
	c := Config{}
	in, err := LoadPacketConfig()
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
	path, _ := os.LookupEnv(defaultPacketConfigEnvVar)
	if path != "" {
		return path
	}

	return defaultPacketConfigPath
}

// Save will serialize the config back into the locations
// specified in @LoadPacketConfig, first location with a file, wins.
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
