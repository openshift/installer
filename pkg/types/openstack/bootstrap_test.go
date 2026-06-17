package openstack

import (
	"testing"
)

func TestResolveBootstrapFlavor(t *testing.T) {
	tests := []struct {
		name               string
		platform           *Platform
		controlPlaneFlavor string
		expected           string
	}{
		{
			name: "explicit bootstrapFlavor is returned",
			platform: &Platform{
				BootstrapFlavor: "m1.xlarge",
			},
			controlPlaneFlavor: "m1.large",
			expected:           "m1.xlarge",
		},
		{
			name:               "empty bootstrapFlavor returns controlPlaneFlavor",
			platform:           &Platform{},
			controlPlaneFlavor: "m1.large",
			expected:           "m1.large",
		},
		{
			name:               "nil platform returns controlPlaneFlavor",
			platform:           nil,
			controlPlaneFlavor: "m1.large",
			expected:           "m1.large",
		},
		{
			name:               "empty both returns empty string",
			platform:           &Platform{},
			controlPlaneFlavor: "",
			expected:           "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.platform.ResolveBootstrapFlavor(tc.controlPlaneFlavor)
			if got != tc.expected {
				t.Errorf("ResolveBootstrapFlavor(%q) = %q, want %q", tc.controlPlaneFlavor, got, tc.expected)
			}
		})
	}
}
