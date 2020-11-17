package ovirt

import (
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"

	ovirtsdk "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var defaultOvirtConfigEnvVar = "OVIRT_CONFIG"
var defaultOvirtConfigPath = filepath.Join(os.Getenv("HOME"), ".ovirt", "ovirt-config.yaml")

// Config holds oVirt api access details
type Config struct {
	URL      string `yaml:"ovirt_url"`
	FQDN     string `yaml:"ovirt_fqdn"`
	PemURL   string `yaml:"ovirt_pem_url"`
	Username string `yaml:"ovirt_username"`
	Password string `yaml:"ovirt_password"`
	CAFile   string `yaml:"ovirt_cafile,omitempty"`
	Insecure bool   `yaml:"ovirt_insecure,omitempty"`
	CABundle string `yaml:"ovirt_ca_bundle,omitempty"`
}

// clientHTTP struct - Hold info about http calls
type clientHTTP struct {
	saveFilePath string // Path for saving file (GET method)
	urlAddr      string // URL or Address
	skipVerify   bool   // skipt cert validatin in the http call
	certPool     *x509.CertPool
}

// LoadOvirtConfig from the following location (first wins):
// 1. OVIRT_CONFIG env variable
// 2  $defaultOvirtConfigPath
// See #@Config for the expected format
func LoadOvirtConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(discoverPath())
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NewConfig will return an Config by loading
// the configuration from locations specified in @LoadOvirtConfig
func NewConfig() (Config, error) {
	c := Config{}
	in, err := LoadOvirtConfig()
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
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, 0600)
}

// getValidatedConnection will create a connection and validate it before returning.
func (c *Config) getValidatedConnection() (*ovirtsdk.Connection, error) {
	connection, err := getConnection(*c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build configuration for ovirt connection validation")
	}

	if err := connection.Test(); err != nil {
		_ = connection.Close()
		return nil, err
	}
	return connection, nil
}
