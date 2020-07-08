package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenter:          "test-vcenter",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *vsphere.Platform
		expectedError string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "missing vCenter name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = ""
				return p
			}(),
			expectedError: `^test-path\.vCenter: Required value: must specify the name of the vCenter$`,
		},
		{
			name: "missing username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Username = ""
				return p
			}(),
			expectedError: `^test-path\.username: Required value: must specify the username$`,
		},
		{
			name: "missing password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Password = ""
				return p
			}(),
			expectedError: `^test-path\.password: Required value: must specify the password$`,
		},
		{
			name: "missing datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Datacenter = ""
				return p
			}(),
			expectedError: `^test-path\.datacenter: Required value: must specify the datacenter$`,
		},
		{
			name: "missing default datastore",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DefaultDatastore = ""
				return p
			}(),
			expectedError: `^test-path\.defaultDatastore: Required value: must specify the default datastore$`,
		},
		{
			name: "valid VIPs",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = "192.168.111.3"
				return p
			}(),
			// expectedError: `^test-path\.apiVIP: Invalid value: "": "" is not a valid IP`,
		},
		{
			name: "missing API VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = ""
				p.IngressVIP = "192.168.111.3"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Required value: must specify a VIP for the API`,
		},
		{
			name: "missing Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = ""
				return p
			}(),
			expectedError: `^test-path\.ingressVIP: Required value: must specify a VIP for Ingress`,
		},
		{
			name: "Invalid API VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111"
				p.IngressVIP = "192.168.111.2"
				return p
			}(),
			expectedError: `^test-path.apiVIP: Invalid value: "192.168.111": "192.168.111" is not a valid IP`,
		},
		{
			name: "Invalid Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111"
				return p
			}(),
			expectedError: `^test-path.ingressVIP: Invalid value: "192.168.111": "192.168.111" is not a valid IP`,
		},
		{
			name: "Same API and Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111.1"
				return p
			}(),
			expectedError: `^test-path.apiVIP: Invalid value: "192.168.111.1": IPs for both API and Ingress should not be the same`,
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
