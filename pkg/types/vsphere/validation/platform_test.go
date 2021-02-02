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
			expectedError: `^test-path\.vCenter: Invalid value: "": must specify the name of the vCenter$`,
		},
		{
			name: "vcenter name has transport",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "http://www.forest.com?param=value"
				return p
			}(),
			expectedError: `^test-path\.vCenter: Invalid value: "http:\/\/www\.forest\.com\?param=value": vCenter hostname cannot contain url scheme$`,
		},
		{
			name: "vcenter name has transport and port",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "https://www.forest.com:9089"
				return p
			}(),
			expectedError: `^test-path\.vCenter: Invalid value: "https:\/\/www\.forest\.com:9089": vCenter hostname cannot contain url scheme$`,
		},
		{
			name: "vcenter name has path and port",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "www.forest/tropical/rain.com:9089"
				return p
			}(),
			expectedError: `^test-path\.vCenter: Invalid value: "www\.forest\/tropical\/rain.com:9089": vCenter hostname is not valid`,
		},
		{
			name: "vcenter name has transport",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "scheme://"
				return p
			}(),
			expectedError: `^test-path\.vCenter: Invalid value: "scheme:\/\/": port can only contain numbers$`,
		},
		{
			name: "vcenter name with port",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "www.forest.tropical:9089"
				return p
			}(),
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
		{
			name: "Capital letters in vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "tEsT-vCenter"
				return p
			}(),
			expectedError: `^test-path.vCenter: Invalid value: "tEsT-vCenter": must be all lower case`,
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

func TestValidateVCenterAddress(t *testing.T) {

	cases := []struct {
		name          string
		vCenterName   string
		expectedError string
	}{
		{
			name:        "vcenter with port",
			vCenterName: "ww.forest.com:9090",
		},
		{
			name:        "vcenter without port",
			vCenterName: "ww.forest.com",
		},
		{
			name:        "vcenter with port and number",
			vCenterName: "ww.878forest.com:9090",
		},
		{
			name:        "vcenter with port and -",
			vCenterName: "ww.forest-hdlkd.c-om:9090",
		},
		{
			name:        "vcenter with port with all allowed charcaters",
			vCenterName: "w-987w.forest.com:9090",
		},
		{
			name:          "vcenter with port cannot start with a -",
			vCenterName:   "-w-987w.forest.com:9090",
			expectedError: "vCenter hostname is not valid",
		},
		{
			name:          "vcenter with port cannot start with a number",
			vCenterName:   "0-w-987w.forest.com:9090",
			expectedError: "vCenter hostname is not valid",
		},
		{
			name:          "vcenter with scheme",
			vCenterName:   "https://www.forest.com:9090",
			expectedError: "vCenter hostname cannot contain url scheme",
		},
		{
			name:          "vcenter with params",
			vCenterName:   "www.forest.com?param=name",
			expectedError: "vCenter hostname cannot contain request params",
		},
		{
			name:          "vcenter with path",
			vCenterName:   "www.forest.com/pathvariable:9090",
			expectedError: "vCenter hostname is not valid",
		},
		{
			name:          "vcenter with invalid port",
			vCenterName:   "www.forest.com:",
			expectedError: "port can only contain numbers",
		},
		{
			name:          "vcenter with invalid port",
			vCenterName:   "www.forest.com:987dsc",
			expectedError: "port can only contain numbers",
		},
		{
			name:          "vcenter with spaces",
			vCenterName:   "   ",
			expectedError: "must specify the name of the vCenter",
		},
		{
			name:          "vcenter empty",
			vCenterName:   "",
			expectedError: "must specify the name of the vCenter",
		},
		{
			name:          "vcenter only transport",
			vCenterName:   "https://",
			expectedError: "port can only contain numbers",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateVCenterAddress(tc.vCenterName)
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
