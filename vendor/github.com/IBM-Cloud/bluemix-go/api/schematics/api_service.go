package schematics

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

//SchematicsServiceAPI is the Aramda K8s client ...
type SchematicsServiceAPI interface {
	Workspaces() Workspaces

	//TODO Add other services
}

//VpcContainerService holds the client
type scService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (SchematicsServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.SchematicsService)
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
		ep, err := config.EndpointLocator.SchematicsEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &scService{
		Client: client.New(config, bluemix.SchematicsService, tokenRefreher),
	}, nil
}

//Clusters implements Clusters API
func (c scService) Workspaces() Workspaces {
	return newWorkspaceAPI(c.Client)
}
