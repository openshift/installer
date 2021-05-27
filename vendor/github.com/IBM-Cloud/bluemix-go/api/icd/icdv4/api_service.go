package icdv4

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

//ICDServiceAPI is the Cloud Internet Services API ...
type ICDServiceAPI interface {
	Cdbs() Cdbs
	Users() Users
	Whitelists() Whitelists
	Groups() Groups
	Tasks() Tasks
	Connections() Connections
	AutoScaling() AutoScaling
}

//ICDService holds the client
type icdService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ICDServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.ICDService)
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
		ep, err := config.EndpointLocator.ICDEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &icdService{
		Client: client.New(config, bluemix.ICDService, tokenRefreher),
	}, nil
}

//Cdbs implements deployments API
func (c *icdService) Cdbs() Cdbs {
	return newCdbAPI(c.Client)
}

//Users implements users API
func (c *icdService) Users() Users {
	return newUsersAPI(c.Client)
}

//Whilelists implements whitelists API
func (c *icdService) Whitelists() Whitelists {
	return newWhitelistAPI(c.Client)
}

//Groups implements groups API
func (c *icdService) Groups() Groups {
	return newGroupAPI(c.Client)
}

//Tasks implements tasks API
func (c *icdService) Tasks() Tasks {
	return newTaskAPI(c.Client)
}

//Tasks implements tasks API
func (c *icdService) Connections() Connections {
	return newConnectionAPI(c.Client)
}

//AutoScaling implements AutoScaling API
func (c *icdService) AutoScaling() AutoScaling {
	return newAutoScalingAPI(c.Client)
}
