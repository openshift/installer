package cisv1

import (
	gohttp "net/http"
	"strconv"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//CisServiceAPI is the Cloud Internet Services API ...
type CisServiceAPI interface {
	Zones() Zones
	Monitors() Monitors
	Pools() Pools
	Glbs() Glbs
	Settings() Settings
	Ips() Ips
	Dns() Dns
	Firewall() Firewall
	RateLimit() RateLimit
}

//CisService holds the client
type cisService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (CisServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.CisService)
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
		ep, err := config.EndpointLocator.CisEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &cisService{
		Client: client.New(config, bluemix.CisService, tokenRefreher),
	}, nil
}

//Zones implement albs API
func (c *cisService) Zones() Zones {
	return newZoneAPI(c.Client)
}

//Monitors implements Monitors API
func (c *cisService) Monitors() Monitors {
	return newMonitorAPI(c.Client)
}

//Pools implements Pools API
func (c *cisService) Pools() Pools {
	return newPoolAPI(c.Client)
}

//Glbs implements Glbs API
func (c *cisService) Glbs() Glbs {
	return newGlbAPI(c.Client)
}

//Settings implements Settings API
func (c *cisService) Settings() Settings {
	return newSettingsAPI(c.Client)
}

//Settings implements Settings API
func (c *cisService) Ips() Ips {
	return newIpsAPI(c.Client)
}

//Settings implements DNS records API
func (c *cisService) Dns() Dns {
	return newDnsAPI(c.Client)
}

func (c *cisService) Firewall() Firewall {
	return newFirewallAPI(c.Client)
}

func (c *cisService) RateLimit() RateLimit {
	return newRateLimitAPI(c.Client)
}

func errorsToString(e []Error) string {

	var errMsg string
	for _, err := range e {
		errFrag := "Code: " + strconv.Itoa(err.Code) + " " + err.Msg
		errMsg = errMsg + errFrag
	}
	return errMsg
}
