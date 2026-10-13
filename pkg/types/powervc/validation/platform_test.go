package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervc"
)

func validControlPlanePort() *powervc.PortTarget {
	return &powervc.PortTarget{
		Network: powervc.NetworkFilter{Name: "network"},
		FixedIPs: []powervc.FixedIP{
			{Subnet: powervc.SubnetFilter{Name: "subnet"}},
		},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *powervc.Platform
		expectedError string
	}{
		{
			name: "valid",
			platform: &powervc.Platform{
				Cloud:            "cloud",
				ClusterOSImage:   "rhcos-image",
				ControlPlanePort: validControlPlanePort(),
			},
		},
		{
			name: "missing controlPlanePort",
			platform: &powervc.Platform{
				Cloud:          "cloud",
				ClusterOSImage: "rhcos-image",
			},
			expectedError: `^test-path\.controlPlanePort: Required value: controlPlanePort is required with powervc$`,
		},
		{
			name: "external network set",
			platform: &powervc.Platform{
				Cloud:            "cloud",
				ClusterOSImage:   "rhcos-image",
				ExternalNetwork:  "external",
				ControlPlanePort: validControlPlanePort(),
			},
			expectedError: `^test-path\.externalNetwork: Invalid value: "external": Cannot set external network with powervc$`,
		},
		{
			name: "missing clusterOSImage",
			platform: &powervc.Platform{
				Cloud:            "cloud",
				ControlPlanePort: validControlPlanePort(),
			},
			expectedError: `^test-path\.clusterOSImage: Required value: clusterOSImage must be the name of an existing image with powervc$`,
		},
		{
			name: "clusterOSImage is an http URL",
			platform: &powervc.Platform{
				Cloud:            "cloud",
				ClusterOSImage:   "http://example.com/rhcos.ova.gz",
				ControlPlanePort: validControlPlanePort(),
			},
			expectedError: `^test-path\.clusterOSImage: Invalid value: "http://example.com/rhcos.ova.gz": clusterOSImage must be the name of an existing image, not a URL, with powervc$`,
		},
		{
			name: "clusterOSImage is a file URL",
			platform: &powervc.Platform{
				Cloud:            "cloud",
				ClusterOSImage:   "file:///home/user/rhcos.ova.gz",
				ControlPlanePort: validControlPlanePort(),
			},
			expectedError: `^test-path\.clusterOSImage: Invalid value: "file:///home/user/rhcos.ova.gz": clusterOSImage must be the name of an existing image, not a URL, with powervc$`,
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
