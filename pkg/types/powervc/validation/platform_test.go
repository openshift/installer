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
				ControlPlanePort: validControlPlanePort(),
			},
		},
		{
			name: "missing controlPlanePort",
			platform: &powervc.Platform{
				Cloud: "cloud",
			},
			expectedError: `^test-path\.controlPlanePort: Required value: controlPlanePort is required with powervc$`,
		},
		{
			name: "external network set",
			platform: &powervc.Platform{
				Cloud:            "cloud",
				ExternalNetwork:  "external",
				ControlPlanePort: validControlPlanePort(),
			},
			expectedError: `^test-path\.externalNetwork: Invalid value: "external": Cannot set external network with powervc$`,
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
