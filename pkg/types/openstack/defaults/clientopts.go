package defaults

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
)

// DefaultClientOpts generates default client opts based on cloud name
func DefaultClientOpts(cloudName string) *clientconfig.ClientOpts {
	opts := new(clientconfig.ClientOpts)
	opts.Cloud = cloudName
	// We explicitly disable reading auth data from env variables by setting an invalid EnvPrefix.
	// By doing this, we make sure that the data from clouds.yaml is enough to authenticate.
	// For more information: https://github.com/gophercloud/utils/blob/8677e053dcf1f05d0fa0a616094aace04690eb94/openstack/clientconfig/requests.go#L508
	opts.EnvPrefix = "NO_ENV_VARIABLES_"
	return opts
}
