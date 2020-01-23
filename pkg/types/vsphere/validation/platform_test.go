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
		Cluster:          "test-cluster",
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
			name: "missing cluster",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Cluster = ""
				return p
			}(),
			expectedError: `^test-path\.cluster: Required value: must specify the cluster$`,
		},
		{
			name: "valid VIPs",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = "192.168.111.3"
				p.DNSVIP = "192.168.111.4"
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
				p.DNSVIP = "192.168.111.4"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Invalid value: "": "" is not a valid IP`,
		},
		{
			name: "missing Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = ""
				p.DNSVIP = "192.168.111.4"
				return p
			}(),
			expectedError: `^test-path\.ingressVIP: Invalid value: "": "" is not a valid IP`,
		},
		{
			name: "missing DNS VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = "192.168.111.3"
				p.DNSVIP = ""
				return p
			}(),
			expectedError: `^test-path\.dnsVIP: Invalid value: "": "" is not a valid IP`,
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
