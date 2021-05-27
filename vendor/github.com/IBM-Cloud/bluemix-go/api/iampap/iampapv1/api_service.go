package iampapv1

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//IAMPAPAPI is the IAMpapv2 client ...
type IAMPAPAPI interface {
	IAMPolicy() IAMPolicy
	IAMService() IAMService
	AuthorizationPolicies() AuthorizationPolicyRepository
	V1Policy() V1PolicyRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//IamPapService holds the client
type iampapService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (IAMPAPAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.IAMPAPService)
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
		ep, err := config.EndpointLocator.IAMPAPEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}
	return &iampapService{
		Client: client.New(config, bluemix.IAMPAPService, tokenRefreher),
	}, nil
}

//IAMPolicy API
func (a *iampapService) IAMPolicy() IAMPolicy {
	return newIAMPolicyAPI(a.Client)
}

//IAMService API
func (a *iampapService) IAMService() IAMService {
	return newIAMServiceAPI(a.Client)
}

//AuthorizationPolicies API
func (a *iampapService) AuthorizationPolicies() AuthorizationPolicyRepository {
	return NewAuthorizationPolicyRepository(a.Client)
}

func (a *iampapService) V1Policy() V1PolicyRepository {
	return NewV1PolicyRepository(a.Client)
}
