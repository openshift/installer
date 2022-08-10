package containerv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//ContainerServiceAPI is the Aramda K8s client ...
type ContainerServiceAPI interface {
	Monitoring() Monitoring
	Logging() Logging
	Clusters() Clusters
	VPCs() VPCs
	WorkerPools() WorkerPool
	Albs() Alb
	Workers() Workers
	Kms() Kms
	Ingresses() Ingress
	Subnets() Subnets
	NlbDns() Nlbdns
	Satellite() Satellite
	DedicatedHost() DedicatedHost
	DedicatedHostPool() DedicatedHostPool
	DedicatedHostFlavor() DedicatedHostFlavor

	//TODO Add other services
}

//VpcContainerService holds the client
type csService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ContainerServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.VpcContainerService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"X-Original-User-Agent": []string{config.UserAgent},
			"User-Agent":            []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	if config.IAMAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.ContainerEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &csService{
		Client: client.New(config, bluemix.VpcContainerService, tokenRefreher),
	}, nil
}

//Clusters implements Clusters API
func (c *csService) Satellite() Satellite {
	return newSatelliteAPI(c.Client)
}

//Clusters implements Clusters API
func (c *csService) Clusters() Clusters {
	return newClusterAPI(c.Client)
}

//VPCs implements Cluster VPCs API
func (c *csService) VPCs() VPCs {
	return newVPCsAPI(c.Client)
}

//Monitor implements Monitor API
func (c *csService) Monitoring() Monitoring {
	return newMonitoringAPI(c.Client)
}

//Logging implements Monitor API
func (c *csService) Logging() Logging {
	return newLoggingAPI(c.Client)
}

//WorkerPools implements Cluster WorkerPools API
func (c *csService) WorkerPools() WorkerPool {
	return newWorkerPoolAPI(c.Client)
}
func (c *csService) Albs() Alb {
	return newAlbAPI(c.Client)
}
func (c *csService) NlbDns() Nlbdns {
	return newNlbdnsAPI(c.Client)
}
func (c *csService) Ingresses() Ingress {
	return newIngressAPI(c.Client)
}

//Kms implements Cluster Kms API
func (c *csService) Kms() Kms {
	return newKmsAPI(c.Client)
}

//Workers implements Cluster Workers API
func (c *csService) Workers() Workers {
	return newWorkerAPI(c.Client)
}

//Subnets implements Cluster Subnets API
func (c *csService) Subnets() Subnets {
	return newSubnetsAPI(c.Client)
}

//DedicatedHost implements DedicatedHost API
func (c *csService) DedicatedHost() DedicatedHost {
	return newDedicatedHostAPI(c.Client)
}

//DedicatedHostPool implements DedicatedHostPool API
func (c *csService) DedicatedHostPool() DedicatedHostPool {
	return newDedicatedHostPoolAPI(c.Client)
}

//DedicatedHostFlavor implements DedicatedHostFlavor API
func (c *csService) DedicatedHostFlavor() DedicatedHostFlavor {
	return newDedicatedHostFlavorAPI(c.Client)
}
