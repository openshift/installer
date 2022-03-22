package nutanix

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

// Version represents api version
const Version = "3.1"

// Config ...
type Config struct {
	Endpoint    string
	Username    string
	Password    string
	Port        string
	Insecure    bool
	SessionAuth bool
	WaitTimeout int64
	ProxyURL    string
}

// Client ...
func (c *Config) Client() (*Client, error) {
	configCreds := client.Credentials{
		URL:         fmt.Sprintf("%s:%s", c.Endpoint, c.Port),
		Endpoint:    c.Endpoint,
		Username:    c.Username,
		Password:    c.Password,
		Port:        c.Port,
		Insecure:    c.Insecure,
		SessionAuth: c.SessionAuth,
		ProxyURL:    c.ProxyURL,
	}

	v3Client, err := v3.NewV3Client(configCreds)
	if err != nil {
		return nil, err
	}
	karbonClient, err := karbon.NewKarbonAPIClient(configCreds)
	if err != nil {
		return nil, err
	}

	return &Client{
		WaitTimeout: c.WaitTimeout,
		API:         v3Client,
		KarbonAPI:   karbonClient,
	}, nil
}

// Client represents the nutanix API client
type Client struct {
	API         *v3.Client
	KarbonAPI   *karbon.Client
	WaitTimeout int64
}
