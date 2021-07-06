package iampapv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//IAMPAPAPIV2 is the resource client ...
type IAMPAPAPIV2 interface {
	IAMRoles() RoleRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//iamService holds the client
type roleService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (IAMPAPAPIV2, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.IAMPAPServicev2)
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
		ep, err := config.EndpointLocator.IAMEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &roleService{
		Client: client.New(config, bluemix.IAMPAPServicev2, tokenRefreher),
	}, nil
}

//CustomRole API
func (a *roleService) IAMRoles() RoleRepository {
	return NewRoleRepository(a.Client)
}
