package validation

import (
	"crypto/rand"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/powervs"
)

func genRandNum(min int64, max int64) int64 {
	// calculate the max we will be using
	bg := big.NewInt(max - min)

	// get big.Int between 0 and bg
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	// add n to min to support the passed in range
	return n.Int64() + min
}

func validMinimalPlatform() *powervs.Platform {
	// Avoid lint error: G404: Use of weak random number generator (math/rand instead of crypto/rand) (gosec)
	zoneNames := powervs.ZoneNames()
	len64 := int64(len(zoneNames))
	idx := genRandNum(0, len64)
	if idx < 0 || idx > len64 {
		panic(errors.New("genRandNum out of bounds of zoneNames"))
	}
	zone := powervs.ZoneNames()[idx]
	return &powervs.Platform{
		Zone: zone,
	}
}

func validMachinePool() *powervs.MachinePool {
	return &powervs.MachinePool{}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *powervs.Platform
		valid    bool
	}{
		{
			name:     "Config: Minimal platform config",
			platform: validMinimalPlatform(),
			valid:    true,
		},
		{
			name: "Zone: Invalid zone",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.Zone = "invalid"
				return p
			}(),
			valid: false,
		},
		{
			name: "Region: Invalid region",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.Region = "invalid"
				return p
			}(),
			valid: false,
		},
		{
			name: "Region: Missing region",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.Region = ""
				return p
			}(),
			valid: true,
		},
		{
			name: "Machine Pool: Valid machine pool",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.DefaultMachinePlatform = validMachinePool()
				return p
			}(),
			valid: true,
		},
		{
			name: "ServiceInstanceID: Valid Power VS ServiceInstanceID empty",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceInstanceGUID = ""
				return p
			}(),
			valid: true,
		},
		{
			name: "ServiceInstanceID: Valid Power VS ServiceInstanceID",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceInstanceGUID = "05d5dbfd-2a62-4d01-b37b-71211be442f6"
				return p
			}(),
			valid: true,
		},
		{
			name: "ServiceInstanceID: Invalid Power VS ServiceInstanceID",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceInstanceGUID = "abc123"
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid url (no hostname) for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "/some/path",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid url (has path) for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "https://test-cos.random.local/some/path",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "valid url (has version path, no trailing '/') for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "https://test-cos.random.local/v2",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "valid url (has version path and trailing '/') for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "https://test-cos.random.local/v35/",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "invalid url (has request) for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "https://test-iam.random.local?foo=some",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "valid url (no scheme) for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "test-cos.random.local",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "valid url (with scheme) for service endpoint",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "https://test-cos.random.local",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "duplicate service endpoints",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "https://test-cos.random.local",
				}, {
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "test-cos.random.local",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "multiple valid service endpoints",
			platform: func() *powervs.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.PowerVSServiceEndpoint{{
					Name: string(configv1.IBMCloudServiceCOS),
					URL:  "test-cos.random.local",
				}, {
					Name: string(configv1.IBMCloudServiceDNSServices),
					URL:  "test-dns.random.local",
				}}
				return p
			}(),
			valid: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
