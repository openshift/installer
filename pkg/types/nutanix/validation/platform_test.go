package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/nutanix"
)

func validPlatform() *nutanix.Platform {
	return &nutanix.Platform{
		PrismCentral: nutanix.PrismCentral{
			Endpoint: nutanix.PrismEndpoint{Address: "test-pc", Port: 8080},
			Username: "test-username-pc",
			Password: "test-password-pc",
		},
		PrismElements: []nutanix.PrismElement{{
			UUID:     "test-pe-uuid",
			Endpoint: nutanix.PrismEndpoint{Address: "test-pe", Port: 8081},
		}},
		SubnetUUIDs: []string{"b06179c8-dea3-4f8e-818a-b2e88fbc2201"},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *nutanix.Platform
		expectedError string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "missing Prism Central address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Required value: must specify the Prism Central endpoint address$`,
		},
		{
			name: "missing Prism Central username",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Username = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.username: Required value: must specify the Prism Central username$`,
		},
		{
			name: "missing Prism Central password",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Password = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.password: Required value: must specify the Prism Central password$`,
		},
		{
			name: "missing prism elements",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements = []nutanix.PrismElement{}
				return p
			}(),
			expectedError: `^test-path\.prismElements: Required value: must specify one Prism Element$`,
		},
		{
			name: "missing prism element uuid",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].UUID = ""
				return p
			}(),
			expectedError: `^test-path\.prismElements\.uuid: Required value: must specify the Prism Element UUID$`,
		},
		{
			name: "missing prism element address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Address = ""
				return p
			}(),
			expectedError: `^test-path\.prismElements\.endpoint\.address: Required value: must specify the Prism Element endpoint address$`,
		},
		{
			name: "missing subnet",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.SubnetUUIDs[0] = ""
				return p
			}(),
			expectedError: `^test-path\.subnet: Required value: must specify the subnet$`,
		},
		{
			name: "Capital letters in Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = "tEsT-PrismCentral"
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Invalid value: "tEsT-PrismCentral": must be the domain name or IP address of the Prism Central$`,
		},
		{
			name: "URL as Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = "https://test-pc"
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Invalid value: "https://test-pc": must be the domain name or IP address of the Prism Central$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
