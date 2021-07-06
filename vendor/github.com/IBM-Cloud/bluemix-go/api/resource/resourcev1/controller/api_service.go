package controller

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ResourceControllerAPI is the resource client ...
type ResourceControllerAPI interface {
	ResourceServiceInstance() ResourceServiceInstanceRepository
	ResourceServiceAlias() ResourceServiceAliasRepository
	ResourceServiceKey() ResourceServiceKeyRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//resourceControllerService holds the client
type resourceControllerService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ResourceControllerAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.ResourceControllerService)
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
		Client: client.New(config, bluemix.ResourceManagementService, tokenRefreher),
	}, nil
}

//ResourceController API
func (a *resourceControllerService) ResourceServiceInstance() ResourceServiceInstanceRepository {
	return newResourceServiceInstanceAPI(a.Client)
}

//ResourceController API
func (a *resourceControllerService) ResourceServiceKey() ResourceServiceKeyRepository {
	return newResourceServiceKeyAPI(a.Client)
}

//ResourceController API
func (a *resourceControllerService) ResourceServiceAlias() ResourceServiceAliasRepository {
	return newResourceServiceAliasRepository(a.Client)
}
