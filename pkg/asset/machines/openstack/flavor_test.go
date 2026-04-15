package openstack

import (
	"testing"

	"github.com/openshift/installer/pkg/types/openstack"
)

func TestResolveBootstrapFlavor(t *testing.T) {
	tests := []struct {
		name               string
		platform           *openstack.Platform
		controlPlaneFlavor string
		want               string
	}{
		{
			name: "explicit bootstrap flavor is returned",
			platform: &openstack.Platform{
				BootstrapFlavor: "m1.xlarge",
			},
			controlPlaneFlavor: "m1.large",
			want:               "m1.xlarge",
		},
		{
			name: "empty bootstrap flavor falls back to control plane flavor",
			platform: &openstack.Platform{
				BootstrapFlavor: "",
			},
			controlPlaneFlavor: "m1.large",
			want:               "m1.large",
		},
		{
			name:               "nil platform falls back to control plane flavor",
			platform:           nil,
			controlPlaneFlavor: "m1.large",
			want:               "m1.large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveBootstrapFlavor(tt.platform, tt.controlPlaneFlavor)
			if got != tt.want {
				t.Errorf("ResolveBootstrapFlavor() = %q, want %q", got, tt.want)
			}
		})
	}
}
