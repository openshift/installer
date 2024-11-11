package openstack

import (
	"context"
	"testing"

	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func TestCloudProviderConfig(t *testing.T) {
	cases := []struct {
		name           string
		installConfig  *types.InstallConfig
		expectedConfig string
	}{
		{
			name: "default install config",
			installConfig: &types.InstallConfig{
				Networking: &types.Networking{},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{},
				},
			},
			expectedConfig: `[Global]
secret-name = openstack-credentials
secret-namespace = kube-system
region = my_region
`,
		},
	}

	cloud := clientconfig.Cloud{
		AuthInfo: &clientconfig.AuthInfo{
			Username:   "my_user",
			Password:   "my_secret_password",
			AuthURL:    "https://my_auth_url.com/v3/",
			ProjectID:  "f12f928576ae4d21bdb984da5dd1d3bf",
			DomainID:   "default",
			DomainName: "Default",
		},
		RegionName: "my_region",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualConfig, _, err := generateCloudProviderConfig(context.Background(), nil, &cloud, *tc.installConfig)
			assert.NoError(t, err, "unexpected error when generating cloud provider config")
			assert.Equal(t, tc.expectedConfig, actualConfig, "unexpected cloud provider config")
		})
	}
}
