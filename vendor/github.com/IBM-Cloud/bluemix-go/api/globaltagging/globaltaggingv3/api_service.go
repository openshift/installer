package globaltaggingv3

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
type GlobalTaggingServiceAPI interface {
	Tags() Tags
}

//ICDService holds the client
type globalTaggingService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (GlobalTaggingServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.GlobalTaggingService)
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
		ep, err := config.EndpointLocator.GlobalTaggingEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &globalTaggingService{
		Client: client.New(config, bluemix.GlobalTaggingService, tokenRefreher),
	}, nil
}

//Tagging implements the global tagging API
func (c *globalTaggingService) Tags() Tags {
	return newTaggingAPI(c.Client)
}
