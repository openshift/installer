package ovirt

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var defaultOvirtConfigEnvVar = "OVIRT_CONFIG"
var defaultOvirtConfigPath = filepath.Join(os.Getenv("HOME"), ".ovirt", "ovirt-config.yaml")

// ErrCanNotLoadOvirtConfig is returned when the config file fails to load.
var ErrCanNotLoadOvirtConfig error = errors.New("can not load ovirt config")

// Config holds oVirt api access details.
type Config struct {
	URL      string `yaml:"ovirt_url"`
	Username string `yaml:"ovirt_username"`
	Password string `yaml:"ovirt_password"`
	CAFile   string `yaml:"ovirt_cafile,omitempty"`
	Insecure bool   `yaml:"ovirt_insecure,omitempty"`
	CABundle string `yaml:"ovirt_ca_bundle,omitempty"`
}

// LoadOvirtConfig from the following location (first wins):
// 1. OVIRT_CONFIG env variable
// 2  $defaultOvirtConfigPath
func LoadOvirtConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(discoverPath())
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetOvirtConfig will return an Config by loading
// the configuration from locations specified in @LoadOvirtConfig
// error is return if the configuration could not be retained.
func GetOvirtConfig() (Config, error) {
	c := Config{}
	in, err := LoadOvirtConfig()
	if err != nil {
		return c, errors.Wrap(ErrCanNotLoadOvirtConfig, err.Error())
	}

	err = yaml.Unmarshal(in, &c)
	if err != nil {
		return c, errors.Wrap(ErrCanNotLoadOvirtConfig, err.Error())
	}

	return c, nil
}

func discoverPath() string {
	path, _ := os.LookupEnv(defaultOvirtConfigEnvVar)
	if path != "" {
		return path
	}

	return defaultOvirtConfigPath
}

// Save will serialize the config back into the locations
// specified in @LoadOvirtConfig, first location with a file, wins.
func (c *Config) Save() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	path := discoverPath()
	return ioutil.WriteFile(path, out, os.FileMode(0600))
}
