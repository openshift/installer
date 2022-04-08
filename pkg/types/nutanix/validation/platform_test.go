package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/nutanix"
)

func validPlatform() *nutanix.Platform {
	return &nutanix.Platform{
		PrismCentral:     "test-pc",
		PrismElementUUID: "12992bc3-e919-454b-980e-8b51e217c9bd",
		Username:         "test-username",
		Password:         "test-password",
		SubnetUUID:       "b06179c8-dea3-4f8e-818a-b2e88fbc2201",
		Port:             "8080",
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
			name: "missing Prism Central name",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral: Required value: must specify the Prism Central$`,
		},
		{
			name: "missing username",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.Username = ""
				return p
			}(),
			expectedError: `^test-path\.username: Required value: must specify the username$`,
		},
		{
			name: "missing password",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.Password = ""
				return p
			}(),
			expectedError: `^test-path\.password: Required value: must specify the password$`,
		},
		{
			name: "missing prism element",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElementUUID = ""
				return p
			}(),
			expectedError: `^test-path\.prismElement: Required value: must specify the Prism Element$`,
		},
		{
			name: "valid VIPs",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = "192.168.111.3"
				return p
			}(),
		},
		{
			name: "missing API VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = ""
				p.IngressVIP = "192.168.111.3"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Required value: must specify a VIP for the API$`,
		},
		{
			name: "missing Ingress VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = ""
				return p
			}(),
			expectedError: `^test-path\.ingressVIP: Required value: must specify a VIP for Ingress$`,
		},
		{
			name: "Invalid API VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111"
				p.IngressVIP = "192.168.111.2"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Invalid value: "192.168.111": "192.168.111" is not a valid IP$`,
		},
		{
			name: "Invalid Ingress VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111"
				return p
			}(),
			expectedError: `^test-path\.ingressVIP: Invalid value: "192.168.111": "192.168.111" is not a valid IP$`,
		},
		{
			name: "Same API and Ingress VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111.1"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Invalid value: "192.168.111.1": IPs for both API and Ingress should not be the same$`,
		},
		{
			name: "Capital letters in Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral = "tEsT-PrismCentral"
				return p
			}(),
			expectedError: `^test-path\.prismCentral: Invalid value: "tEsT-PrismCentral": must be the domain name or IP address of the Prism Central$`,
		},
		{
			name: "URL as Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral = "https://test-pc"
				return p
			}(),
			expectedError: `^test-path\.prismCentral: Invalid value: "https://test-pc": must be the domain name or IP address of the Prism Central$`,
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
