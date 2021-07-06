package usermanagementv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//UserManagementAPI is the resource client ...
type UserManagementAPI interface {
	UserInvite() Users
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//userManagement holds the client
type userManagement struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (UserManagementAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.UserManagement)
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
		ep, err := config.EndpointLocator.UserManagementEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &userManagement{
		Client: client.New(config, bluemix.UserManagement, tokenRefreher),
	}, nil
}

// UserInvite API
func (a *userManagement) UserInvite() Users {
	return NewUserInviteHandler(a.Client)
}
