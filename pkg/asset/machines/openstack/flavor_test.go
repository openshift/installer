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
		{
			name: "explicit bootstrap flavor takes precedence even when control plane flavor is empty",
			platform: &openstack.Platform{
				BootstrapFlavor: "m1.xlarge",
			},
			controlPlaneFlavor: "",
			want:               "m1.xlarge",
		},
		{
			name: "control plane flavor empty and bootstrap flavor empty returns empty",
			platform: &openstack.Platform{
				BootstrapFlavor: "",
			},
			controlPlaneFlavor: "",
			want:               "",
		},
		{
			name:               "nil platform and empty control plane flavor returns empty",
			platform:           nil,
			controlPlaneFlavor: "",
			want:               "",
		},
		{
			name: "bootstrap flavor case is preserved exactly",
			platform: &openstack.Platform{
				BootstrapFlavor: "M1.XLARGE",
			},
			controlPlaneFlavor: "m1.large",
			want:               "M1.XLARGE",
		},
		{
			name: "control plane flavor used as fallback when bootstrap flavor omitted",
			platform: &openstack.Platform{
				Cloud: "mycloud",
				// BootstrapFlavor intentionally not set
			},
			controlPlaneFlavor: "m1.xlarge",
			want:               "m1.xlarge",
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
