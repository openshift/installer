package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/swauth"
	osClient "github.com/gophercloud/utils/client"
	"github.com/gophercloud/utils/internal"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/gophercloud/utils/terraform/mutexkv"
)

type Config struct {
	CACertFile                  string
	ClientCertFile              string
	ClientKeyFile               string
	Cloud                       string
	DefaultDomain               string
	DomainID                    string
	DomainName                  string
	EndpointOverrides           map[string]interface{}
	EndpointType                string
	IdentityEndpoint            string
	Insecure                    *bool
	Password                    string
	ProjectDomainName           string
	ProjectDomainID             string
	Region                      string
	Swauth                      bool
	TenantID                    string
	TenantName                  string
	Token                       string
	UserDomainName              string
	UserDomainID                string
	Username                    string
	UserID                      string
	ApplicationCredentialID     string
	ApplicationCredentialName   string
	ApplicationCredentialSecret string
	UseOctavia                  bool
	MaxRetries                  int
	DisableNoCacheHeader        bool

	DelayedAuth   bool
	AllowReauth   bool
	OsClient      *gophercloud.ProviderClient
	authOpts      *gophercloud.AuthOptions
	authenticated bool
	authFailed    error
	swClient      *gophercloud.ServiceClient
	swAuthFailed  error

	TerraformVersion string
	SDKVersion       string

	mutexkv.MutexKV
}

// LoadAndValidate performs the authentication and initial configuration
// of an OpenStack Provider Client. This sets up the HTTP client and
// authenticates to an OpenStack cloud.
//
// Individual Service Clients are created later in this file.
func (c *Config) LoadAndValidate() error {
	// Make sure at least one of auth_url or cloud was specified.
	if c.IdentityEndpoint == "" && c.Cloud == "" {
		return fmt.Errorf("One of 'auth_url' or 'cloud' must be specified")
	}

	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	clientOpts := new(clientconfig.ClientOpts)

	// If a cloud entry was given, base AuthOptions on a clouds.yaml file.
	if c.Cloud != "" {
		clientOpts.Cloud = c.Cloud

		cloud, err := clientconfig.GetCloudFromYAML(clientOpts)
		if err != nil {
			return err
		}

		if c.Region == "" && cloud.RegionName != "" {
			c.Region = cloud.RegionName
		}

		if c.CACertFile == "" && cloud.CACertFile != "" {
			c.CACertFile = cloud.CACertFile
		}

		if c.ClientCertFile == "" && cloud.ClientCertFile != "" {
			c.ClientCertFile = cloud.ClientCertFile
		}

		if c.ClientKeyFile == "" && cloud.ClientKeyFile != "" {
			c.ClientKeyFile = cloud.ClientKeyFile
		}

		if c.Insecure == nil && cloud.Verify != nil {
			v := (!*cloud.Verify)
			c.Insecure = &v
		}
	} else {
		authInfo := &clientconfig.AuthInfo{
			AuthURL:                     c.IdentityEndpoint,
			DefaultDomain:               c.DefaultDomain,
			DomainID:                    c.DomainID,
			DomainName:                  c.DomainName,
			Password:                    c.Password,
			ProjectDomainID:             c.ProjectDomainID,
			ProjectDomainName:           c.ProjectDomainName,
			ProjectID:                   c.TenantID,
			ProjectName:                 c.TenantName,
			Token:                       c.Token,
			UserDomainID:                c.UserDomainID,
			UserDomainName:              c.UserDomainName,
			Username:                    c.Username,
			UserID:                      c.UserID,
			ApplicationCredentialID:     c.ApplicationCredentialID,
			ApplicationCredentialName:   c.ApplicationCredentialName,
			ApplicationCredentialSecret: c.ApplicationCredentialSecret,
		}
		clientOpts.AuthInfo = authInfo
	}

	ao, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] OpenStack allowReauth: %t", c.AllowReauth)
	ao.AllowReauth = c.AllowReauth

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(terraformUserAgent(c.TerraformVersion, c.SDKVersion))

	config, err := internal.PrepareTLSConfig(c.CACertFile, c.ClientCertFile, c.ClientKeyFile, c.Insecure)
	if err != nil {
		return err
	}

	var logger osClient.Logger
	// if OS_DEBUG is set, log the requests and responses
	if os.Getenv("OS_DEBUG") != "" {
		logger = &osClient.DefaultLogger{}
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &osClient.RoundTripper{
			Rt:         transport,
			MaxRetries: c.MaxRetries,
			Logger:     logger,
		},
	}

	if !c.DisableNoCacheHeader {
		extraHeaders := map[string][]string{
			"Cache-Control": {"no-cache"},
		}
		client.HTTPClient.Transport.(*osClient.RoundTripper).SetHeaders(extraHeaders)
	}

	if !c.DelayedAuth && !c.Swauth {
		err = openstack.Authenticate(client, *ao)
		if err != nil {
			return err
		}
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	c.authOpts = ao
	c.OsClient = client

	return nil
}

func (c *Config) Authenticate() error {
	if !c.DelayedAuth {
		return nil
	}

	c.MutexKV.Lock("auth")
	defer c.MutexKV.Unlock("auth")

	if c.authFailed != nil {
		return c.authFailed
	}

	if !c.authenticated {
		if err := openstack.Authenticate(c.OsClient, *c.authOpts); err != nil {
			c.authFailed = err
			return err
		}
		c.authenticated = true
	}

	return nil
}

// DetermineEndpoint is a helper method to determine if the user wants to
// override an endpoint returned from the catalog.
func (c *Config) DetermineEndpoint(client *gophercloud.ServiceClient, service string) *gophercloud.ServiceClient {
	finalEndpoint := client.ResourceBaseURL()

	if v, ok := c.EndpointOverrides[service]; ok {
		if endpoint, ok := v.(string); ok && endpoint != "" {
			finalEndpoint = endpoint
			client.Endpoint = endpoint
			client.ResourceBase = ""
		}
	}

	log.Printf("[DEBUG] OpenStack Endpoint for %s: %s", service, finalEndpoint)

	return client
}

// DetermineRegion is a helper method to determine the region based on
// the user's settings.
func (c *Config) DetermineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] OpenStack Region is: %s", region)
	return region
}

// The following methods assist with the creation of individual Service Clients
// which interact with the various OpenStack services.

type commonCommonServiceClientInitFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)

func (c *Config) CommonServiceClientInit(newClient commonCommonServiceClientInitFunc, region, service string) (*gophercloud.ServiceClient, error) {
	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	client, err := newClient(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.DetermineRegion(region),
		Availability: clientconfig.GetEndpointType(c.EndpointType),
	})

	if err != nil {
		return client, err
	}

	// Check if an endpoint override was specified for the volume service.
	client = c.DetermineEndpoint(client, service)

	return client, nil
}

func (c *Config) BlockStorageV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewBlockStorageV1, region, "volume")
}

func (c *Config) BlockStorageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewBlockStorageV2, region, "volumev2")
}

func (c *Config) BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewBlockStorageV3, region, "volumev3")
}

func (c *Config) ComputeV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewComputeV2, region, "compute")
}

func (c *Config) DNSV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewDNSV2, region, "dns")
}

func (c *Config) IdentityV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewIdentityV3, region, "identity")
}

func (c *Config) ImageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewImageServiceV2, region, "image")
}

func (c *Config) MessagingV2Client(region string) (*gophercloud.ServiceClient, error) {
	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	client, err := openstack.NewMessagingV2(c.OsClient, "", gophercloud.EndpointOpts{
		Region:       c.DetermineRegion(region),
		Availability: clientconfig.GetEndpointType(c.EndpointType),
	})

	if err != nil {
		return client, err
	}

	// Check if an endpoint override was specified for the messaging service.
	client = c.DetermineEndpoint(client, "message")

	return client, nil
}

func (c *Config) NetworkingV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewNetworkV2, region, "network")
}

func (c *Config) ObjectStorageV1Client(region string) (*gophercloud.ServiceClient, error) {
	var client *gophercloud.ServiceClient
	var err error

	// If Swift Authentication is being used, return a swauth client.
	// Otherwise, use a Keystone-based client.
	if c.Swauth {
		if !c.DelayedAuth {
			client, err = swauth.NewObjectStorageV1(c.OsClient, swauth.AuthOpts{
				User: c.Username,
				Key:  c.Password,
			})
		} else {
			c.MutexKV.Lock("SwAuth")
			defer c.MutexKV.Unlock("SwAuth")

			if c.swAuthFailed != nil {
				return nil, c.swAuthFailed
			}

			if c.swClient == nil {
				c.swClient, err = swauth.NewObjectStorageV1(c.OsClient, swauth.AuthOpts{
					User: c.Username,
					Key:  c.Password,
				})

				if err != nil {
					c.swAuthFailed = err
					return nil, err
				}
			}

			client = c.swClient
		}
	} else {
		if err := c.Authenticate(); err != nil {
			return nil, err
		}

		client, err = openstack.NewObjectStorageV1(c.OsClient, gophercloud.EndpointOpts{
			Region:       c.DetermineRegion(region),
			Availability: clientconfig.GetEndpointType(c.EndpointType),
		})

		if err != nil {
			return client, err
		}
	}

	// Check if an endpoint override was specified for the object-store service.
	client = c.DetermineEndpoint(client, "object-store")

	return client, nil
}

func (c *Config) OrchestrationV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewOrchestrationV1, region, "orchestration")
}

func (c *Config) LoadBalancerV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewLoadBalancerV2, region, "octavia")
}

func (c *Config) DatabaseV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewDBV1, region, "database")
}

func (c *Config) ContainerInfraV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewContainerInfraV1, region, "container-infra")
}

func (c *Config) SharedfilesystemV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewSharedFileSystemV2, region, "sharev2")
}

func (c *Config) KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewKeyManagerV1, region, "key-manager")
}
