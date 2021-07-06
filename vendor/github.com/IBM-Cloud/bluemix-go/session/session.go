package session

import (
	"fmt"
	"strconv"
	"time"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/endpoints"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

//Session ...
type Session struct {
	Config *bluemix.Config
}

//New ...
func New(configs ...*bluemix.Config) (*Session, error) {
	var c *bluemix.Config

	if len(configs) == 0 {
		c = &bluemix.Config{}
	} else {
		c = configs[0]
	}
	sess := &Session{
		Config: c,
	}

	if len(c.IBMID) == 0 {
		c.IBMID = helpers.EnvFallBack([]string{"IBMID"}, "")
	}

	if len(c.IBMIDPassword) == 0 {
		c.IBMIDPassword = helpers.EnvFallBack([]string{"IBMID_PASSWORD"}, "")
	}

	if len(c.BluemixAPIKey) == 0 {
		c.BluemixAPIKey = helpers.EnvFallBack([]string{"IC_API_KEY", "IBMCLOUD_API_KEY", "BM_API_KEY", "BLUEMIX_API_KEY"}, "")
	}

	if len(c.IAMAccessToken) == 0 {
		c.IAMAccessToken = helpers.EnvFallBack([]string{"IC_IAM_TOKEN", "IBMCLOUD_IAM_TOKEN"}, "")
	}

	if len(c.IAMRefreshToken) == 0 {
		c.IAMRefreshToken = helpers.EnvFallBack([]string{"IC_IAM_REFRESH_TOKEN", "IBMCLOUD_IAM_REFRESH_TOKEN"}, "")
	}

	if len(c.Region) == 0 {
		c.Region = helpers.EnvFallBack([]string{"IC_REGION", "IBMCLOUD_REGION", "BM_REGION", "BLUEMIX_REGION"}, "us-south")
	}
	if c.MaxRetries == nil {
		c.MaxRetries = helpers.Int(3)
		retries := helpers.EnvFallBack([]string{"MAX_RETRIES"}, "3")
		i, err := strconv.Atoi(retries)
		if err != nil {
			fmt.Printf("MAX_RETRIES has invalid retries format. Default retries will be set to %q", *c.MaxRetries)
		}
		if i < 0 {
			fmt.Printf("MAX_RETRIES has invalid retries format. Default retries will be set to %q", *c.MaxRetries)
		}
		if err == nil && i >= 0 {
			c.MaxRetries = &i
		}
	}
	if c.HTTPTimeout == 0 {
		c.HTTPTimeout = 180 * time.Second
		timeout := helpers.EnvFallBack([]string{"IC_TIMEOUT", "IBMCLOUD_TIMEOUT", "BM_TIMEOUT", "BLUEMIX_TIMEOUT"}, "180")
		timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%ss", timeout))
		if err != nil {
			fmt.Printf("IC_TIMEOUT or IBMCLOUD_TIMEOUT has invalid time format. Default timeout will be set to %q", c.HTTPTimeout)
		}
		if err == nil {
			c.HTTPTimeout = timeoutDuration
		}
	}
	if len(c.Visibility) == 0 {
		c.Visibility = helpers.EnvFallBack([]string{"IC_VISIBILITY", "IBMCLOUD_VISIBILITY"}, "public")
	}
	if c.RetryDelay == nil {
		c.RetryDelay = helpers.Duration(30 * time.Second)
	}
	if c.EndpointLocator == nil {
		c.EndpointLocator = endpoints.NewEndpointLocator(c.Region, c.Visibility)
	}

	if c.Debug {
		trace.Logger = trace.NewLogger("true")
	}

	return sess, nil
}

//Copy allows sessions to create a copy of it and optionally override any defaults via the config
func (s *Session) Copy(mccpgs ...*bluemix.Config) *Session {
	return &Session{
		Config: s.Config.Copy(mccpgs...),
	}
}
