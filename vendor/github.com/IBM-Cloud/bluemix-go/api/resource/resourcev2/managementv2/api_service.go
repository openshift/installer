package managementv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ResourceManagementAPI is the resource client ...
type ResourceManagementAPIv2 interface {
	ResourceQuota() ResourceQuotaRepository
	ResourceGroup() ResourceGroupRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//resourceManagementService holds the client
type resourceManagementService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ResourceManagementAPIv2, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.ResourceManagementService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
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
		ep, err := config.EndpointLocator.ResourceManagementEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}
	return &resourceManagementService{
		Client: client.New(config, bluemix.ResourceManagementService, tokenRefreher),
	}, nil
}

//ResourceQuota API
func (a *resourceManagementService) ResourceQuota() ResourceQuotaRepository {
	return newResourceQuotaAPI(a.Client)
}

//ResourceGroup API
func (a *resourceManagementService) ResourceGroup() ResourceGroupRepository {
	return newResourceGroupAPI(a.Client)
}
