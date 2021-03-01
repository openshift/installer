package defaults

import (
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
)

// We explicitly disable reading auth data from env variables by setting an invalid EnvPrefix.
// By doing this, we make sure that the data from clouds.yaml is enough to authenticate.
// For more information: https://github.com/gophercloud/utils/blob/8677e053dcf1f05d0fa0a616094aace04690eb94/openstack/clientconfig/requests.go#L508
const envPrefix = "NO_ENV_VARIABLES_"

var serviceClientOpts *clientconfig.ClientOpts

// DefaultClientOpts generates default client opts based on cloud name. First, it tries to create
// and reuse an existing auth token to prevent multiple reauthentications during installation.
// If it's not possible the function returns static client opts with "password" authentication.
func DefaultClientOpts(cloudName string) *clientconfig.ClientOpts {
	if serviceClientOpts != nil {
		return serviceClientOpts
	}

	conn, err := clientconfig.NewServiceClient("identity", staticClientOpts(cloudName))
	if err != nil {
		logrus.Warnf("Cannot authenticate with given credentials: %v", err)
		return staticClientOpts(cloudName)
	}

	authResult := conn.GetAuthResult()
	auth, ok := authResult.(tokens.CreateResult)
	if !ok {
		logrus.Warn("Unable to find auth token in the response")
		return staticClientOpts(cloudName)
	}

	tokenID, err := auth.ExtractTokenID()
	if err != nil {
		logrus.Warnf("Unable to extract auth token: %v", err)
		return staticClientOpts(cloudName)
	}

	project, err := auth.ExtractProject()
	if err != nil {
		logrus.Warnf("Unable to extract project: %v", err)
		return staticClientOpts(cloudName)
	}

	serviceClientOpts = staticClientOpts(cloudName)
	serviceClientOpts.AuthType = clientconfig.AuthV3Token
	serviceClientOpts.AuthInfo = &clientconfig.AuthInfo{
		Token:     tokenID,
		AuthURL:   conn.IdentityEndpoint,
		ProjectID: project.ID,
	}

	return serviceClientOpts
}

func staticClientOpts(cloudName string) *clientconfig.ClientOpts {
	opts := new(clientconfig.ClientOpts)
	opts.Cloud = cloudName
	opts.EnvPrefix = envPrefix
	return opts
}
