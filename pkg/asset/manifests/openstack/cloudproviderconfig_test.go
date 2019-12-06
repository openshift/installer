package openstack

import (
	"testing"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfigSecret(t *testing.T) {
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

	expectedConfig := `[Global]
auth-url = "https://my_auth_url.com/v3/"
username = "my_user"
password = "my_secret_password"
tenant-id = "f12f928576ae4d21bdb984da5dd1d3bf"
domain-id = "default"
domain-name = "Default"
region = "my_region"
`
	actualConfig, err := CloudProviderConfigSecret(&cloud)
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, string(actualConfig), "unexpected cloud provider config")
}

func TestCloudProviderConfigSecretUserDomain(t *testing.T) {
	cloud := clientconfig.Cloud{
		AuthInfo: &clientconfig.AuthInfo{
			Username:       "my_user",
			Password:       "my_secret_password",
			AuthURL:        "https://my_auth_url.com/v3/",
			ProjectID:      "f12f928576ae4d21bdb984da5dd1d3bf",
			UserDomainID:   "default",
			UserDomainName: "Default",
		},
		RegionName: "my_region",
	}

	expectedConfig := `[Global]
auth-url = "https://my_auth_url.com/v3/"
username = "my_user"
password = "my_secret_password"
tenant-id = "f12f928576ae4d21bdb984da5dd1d3bf"
domain-id = "default"
domain-name = "Default"
region = "my_region"
`
	actualConfig, err := CloudProviderConfigSecret(&cloud)
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, string(actualConfig), "unexpected cloud provider config")
}

func TestCloudProviderConfigSecretQuoting(t *testing.T) {
	passwords := map[string]string{
		"regular":        "regular",
		"with\\n":        "with\\\\n",
		"with#":          "with#",
		"with$":          "with$",
		"with;":          "with;",
		"with \n \" \\ ": "with \\n \\\" \\\\ ",
		"with!":          "with!",
		"with?":          "with?",
		"with`":          "with`",
	}

	for k, v := range passwords {
		cloud := clientconfig.Cloud{
			AuthInfo: &clientconfig.AuthInfo{
				Password: k,
			},
		}

		expectedConfig := `[Global]
password = "` + v + `"
`
		actualConfig, err := CloudProviderConfigSecret(&cloud)
		assert.NoError(t, err, "failed to create cloud provider config")
		assert.Equal(t, expectedConfig, string(actualConfig), "unexpected cloud provider config")
	}
}

func TestCloudProviderConfig(t *testing.T) {
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

	expectedConfig := `[Global]
secret-name = openstack-credentials
secret-namespace = kube-system
region = my_region
`
	actualConfig := CloudProviderConfig(&cloud)
	assert.Equal(t, expectedConfig, string(actualConfig), "unexpected cloud provider config")
}
