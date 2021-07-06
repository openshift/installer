package certificatemanager

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

//CertificateManagerServiceAPI is the Aramda K8s client ...
type CertificateManagerServiceAPI interface {
	Certificate() Certificate
}

//CertificateManager Service holds the client
type cmService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (CertificateManagerServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.CertificateManager)
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
		ep, err := config.EndpointLocator.CertificateManagerEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &cmService{
		Client: client.New(config, bluemix.CertificateManager, tokenRefreher),
	}, nil
}

func (c *cmService) Certificate() Certificate {
	return newCertificateAPI(c.Client)
}
