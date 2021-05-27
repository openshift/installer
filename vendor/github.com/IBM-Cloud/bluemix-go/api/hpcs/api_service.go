package hpcs

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//HPCSV2 is the resource client ...
type HPCSV2 interface {
	Endpoint() EndpointRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//hpcsService holds the client
type hpcsService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (HPCSV2, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.HPCService)
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
		ep, err := config.EndpointLocator.HpcsEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &hpcsService{
		Client: client.New(config, bluemix.HPCService, tokenRefreher),
	}, nil
}

//Hpcs API
func (a *hpcsService) Endpoint() EndpointRepository {
	return NewHpcsEndpointRepository(a.Client)
}
