package nutanix

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	foundation_central "github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

// Version represents api version
const Version = "3.1"

// Config ...
type Config struct {
	Endpoint           string
	Username           string
	Password           string
	Port               string
	Insecure           bool
	SessionAuth        bool
	WaitTimeout        int64
	ProxyURL           string
	FoundationEndpoint string              // Required field for connecting to foundation VM APIs
	FoundationPort     string              // Port for connecting to foundation VM APIs
	RequiredFields     map[string][]string // RequiredFields is client name to its required fields mapping for validations and usage in every client
}

// Client ...
func (c *Config) Client() (*Client, error) {
	configCreds := client.Credentials{
		URL:                fmt.Sprintf("%s:%s", c.Endpoint, c.Port),
		Endpoint:           c.Endpoint,
		Username:           c.Username,
		Password:           c.Password,
		Port:               c.Port,
		Insecure:           c.Insecure,
		SessionAuth:        c.SessionAuth,
		ProxyURL:           c.ProxyURL,
		FoundationEndpoint: c.FoundationEndpoint,
		FoundationPort:     c.FoundationPort,
		RequiredFields:     c.RequiredFields,
	}

	v3Client, err := v3.NewV3Client(configCreds)
	if err != nil {
		return nil, err
	}
	karbonClient, err := karbon.NewKarbonAPIClient(configCreds)
	if err != nil {
		return nil, err
	}
	foundationClient, err := foundation.NewFoundationAPIClient(configCreds)
	if err != nil {
		return nil, err
	}
	fcClient, err := foundation_central.NewFoundationCentralClient(configCreds)
	if err != nil {
		return nil, err
	}
	return &Client{
		WaitTimeout:         c.WaitTimeout,
		API:                 v3Client,
		KarbonAPI:           karbonClient,
		FoundationClientAPI: foundationClient,
		FoundationCentral:   fcClient,
	}, nil
}

// Client represents the nutanix API client
type Client struct {
	API                 *v3.Client
	KarbonAPI           *karbon.Client
	FoundationClientAPI *foundation.Client
	WaitTimeout         int64
	FoundationCentral   *foundation_central.Client
}
