package functions

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

//FunctionServiceAPI ..
type FunctionServiceAPI interface {
	Namespaces() Functions
}

//fnService holds the client
type fnService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (FunctionServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.FunctionsService)
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
		ep, err := config.EndpointLocator.FunctionsEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &fnService{
		Client: client.New(config, bluemix.FunctionsService, tokenRefreher),
	}, nil
}

//NewCF ...
func NewCF(sess *session.Session) (FunctionServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.FunctionsService)
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

	if config.UAAAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}

	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.FunctionsEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &fnService{
		Client: client.New(config, bluemix.FunctionsService, tokenRefreher),
	}, nil
}

//Namespaces ..
func (ns *fnService) Namespaces() Functions {
	return newFunctionsAPI(ns.Client)
}
