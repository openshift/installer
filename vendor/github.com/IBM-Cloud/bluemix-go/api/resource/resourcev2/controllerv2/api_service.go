package controllerv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ResourceControllerAPIV2 is the resource client ...
type ResourceControllerAPIV2 interface {
	ResourceServiceInstanceV2() ResourceServiceInstanceRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//resourceControllerService holds the client
type resourceControllerService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ResourceControllerAPIV2, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.ResourceControllerServicev2)
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
		ep, err := config.EndpointLocator.ResourceControllerEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}
	return &resourceControllerService{
		Client: client.New(config, bluemix.ResourceControllerServicev2, tokenRefreher),
	}, nil
}

//ResourceController API
func (a *resourceControllerService) ResourceServiceInstanceV2() ResourceServiceInstanceRepository {
	return newResourceServiceInstanceAPI(a.Client)
}
