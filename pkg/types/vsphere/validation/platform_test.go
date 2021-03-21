package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenter:          "test-vcenter",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
		APIVIP:           "192.168.111.2",
		IngressVIP:       "192.168.111.3",
	}
}

func validNetwork() *types.Networking {
	n := types.Networking{}
	cidr := ipnet.MustParseCIDR("192.168.111.1/24")
	n.MachineNetwork = append(n.MachineNetwork, types.MachineNetworkEntry{
		CIDR: *cidr,
	})
	return &n
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *vsphere.Platform
		networking    *types.Networking
		expectedError string
	}{
		{
			name:       "minimal",
			platform:   validPlatform(),
			networking: validNetwork(),
		},
		{
			name: "missing vCenter name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = ""
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.vCenter: Required value: must specify the name of the vCenter$`,
		},
		{
			name: "missing username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Username = ""
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.username: Required value: must specify the username$`,
		},
		{
			name: "missing password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Password = ""
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.password: Required value: must specify the password$`,
		},
		{
			name: "missing datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Datacenter = ""
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.datacenter: Required value: must specify the datacenter$`,
		},
		{
			name: "missing default datastore",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DefaultDatastore = ""
				return p
			}(),
			networking:    validNetwork(),
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
			networking: validNetwork(),
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
			networking:    validNetwork(),
			expectedError: `^test-path\.apiVIP: Required value: must specify a VIP for both API and Ingress VIPs when specifying either`,
		},
		{
			name: "missing Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = ""
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.ingressVIP: Required value: must specify a VIP for both API and Ingress VIPs when specifying either`,
		},
		{
			name: "Invalid API VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111"
				p.IngressVIP = "192.168.111.2"
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path.apiVIP: Invalid value: "192.168.111": not a valid IP address$`,
		},
		{
			name: "Invalid Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111"
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path.ingressVIP: Invalid value: "192.168.111": not a valid IP address$`,
		},
		{
			name: "Same API and Ingress VIP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111.1"
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path.apiVIP: Invalid value: "192.168.111.1": IPs for both API and Ingress should not be the same`,
		},
		{
			name: "Capital letters in vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "tEsT-vCenter"
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.vCenter: Invalid value: "tEsT-vCenter": must be the domain name or IP address of the vCenter`,
		},
		{
			name: "URL as vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "https://test-center"
				return p
			}(),
			networking:    validNetwork(),
			expectedError: `^test-path\.vCenter: Invalid value: "https://test-center": must be the domain name or IP address of the vCenter$`,
		},
		{
			name: "APIVIP not in machine CIDR",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.1.1"
				p.IngressVIP = "192.168.0.1"
				return p
			}(),
			networking: func() *types.Networking {
				p := types.Networking{}
				value, err := ipnet.ParseCIDR("192.168.0.0/24")
				if err != nil {
					fmt.Println(err)
				}
				p.MachineNetwork = append(p.MachineNetwork, types.MachineNetworkEntry{
					CIDR: *value,
				})
				return &p
			}(),
			expectedError: `^test-path.apiVIP: Invalid value: "192.168.1.1": must be contained within one of the machine networks$`,
		},
		{
			name: "IngressVIP not in machine CIDR",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.0.1"
				p.IngressVIP = "192.168.1.1"
				return p
			}(),
			networking: func() *types.Networking {
				p := types.Networking{}
				value, err := ipnet.ParseCIDR("192.168.0.0/24")
				if err != nil {
					fmt.Println(err)
				}
				p.MachineNetwork = append(p.MachineNetwork, types.MachineNetworkEntry{
					CIDR: *value,
				})
				return &p
			}(),
			expectedError: `^test-path.ingressVIP: Invalid value: "192.168.1.1": must be contained within one of the machine networks$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, tc.networking, field.NewPath("test-path")).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
