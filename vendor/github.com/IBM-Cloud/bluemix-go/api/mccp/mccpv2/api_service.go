package mccpv2

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

//MccpServiceAPI is the mccpv2 client ...
type MccpServiceAPI interface {
	Organizations() Organizations
	Spaces() Spaces
	ServiceInstances() ServiceInstances
	ServiceKeys() ServiceKeys
	ServicePlans() ServicePlans
	ServiceOfferings() ServiceOfferings
	SpaceQuotas() SpaceQuotas
	OrgQuotas() OrgQuotas
	Apps() Apps
	Routes() Routes
	SharedDomains() SharedDomains
	PrivateDomains() PrivateDomains
	ServiceBindings() ServiceBindings
	Regions() RegionRepository
}

//MccpService holds the client
type mccpService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (MccpServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.MccpService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewUAARepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	if config.UAAAccessToken == "" || config.UAARefreshToken == "" {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.MCCPAPIEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &mccpService{
		Client: client.New(config, bluemix.MccpService, tokenRefreher),
	}, nil
}

//Organizations implements Organizations APIs
func (c *mccpService) Organizations() Organizations {
	return newOrganizationAPI(c.Client)
}

//Spaces implements Spaces APIs
func (c *mccpService) Spaces() Spaces {
	return newSpacesAPI(c.Client)
}

//ServicePlans implements ServicePlans APIs
func (c *mccpService) ServicePlans() ServicePlans {
	return newServicePlanAPI(c.Client)
}

//ServiceOfferings implements ServiceOfferings APIs
func (c *mccpService) ServiceOfferings() ServiceOfferings {
	return newServiceOfferingAPI(c.Client)
}

//ServiceInstances implements ServiceInstances APIs
func (c *mccpService) ServiceInstances() ServiceInstances {
	return newServiceInstanceAPI(c.Client)
}

//ServiceKeys implements ServiceKey APIs
func (c *mccpService) ServiceKeys() ServiceKeys {
	return newServiceKeyAPI(c.Client)
}

//SpaceQuotas implements SpaceQuota APIs
func (c *mccpService) SpaceQuotas() SpaceQuotas {
	return newSpaceQuotasAPI(c.Client)
}

//OrgQuotas implements OrgQuota APIs
func (c *mccpService) OrgQuotas() OrgQuotas {
	return newOrgQuotasAPI(c.Client)
}

//ServiceBindings implements ServiceBindings APIs
func (c *mccpService) ServiceBindings() ServiceBindings {
	return newServiceBindingAPI(c.Client)
}

//Apps implements Apps APIs

func (c *mccpService) Apps() Apps {
	return newAppAPI(c.Client)
}

//Routes implements Route APIs

func (c *mccpService) Routes() Routes {
	return newRouteAPI(c.Client)
}

//SharedDomains implements SharedDomian APIs

func (c *mccpService) SharedDomains() SharedDomains {
	return newSharedDomainAPI(c.Client)
}

//PrivateDomains implements PrivateDomains APIs

func (c *mccpService) PrivateDomains() PrivateDomains {
	return newPrivateDomainAPI(c.Client)
}

//Regions implements Regions APIs

func (c *mccpService) Regions() RegionRepository {
	return newRegionRepositoryAPI(c.Client)
}
