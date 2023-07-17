package validation

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func TestFilledInControlPlanePortField(t *testing.T) {
	t.Run("control_plane_port_field", func(t *testing.T) {
		fixedIP := openstack.FixedIP{
			Subnet: openstack.SubnetFilter{ID: "031a5b9d-5a89-4465-8d54-3517ec2bad48"},
		}
		installConfig := types.InstallConfig{
			Platform: types.Platform{
				OpenStack: &openstack.Platform{
					ControlPlanePort: &openstack.PortTarget{
						FixedIPs: []openstack.FixedIP{fixedIP},
					},
				},
			},
		}

		expectedField := field.NewPath("platform", "openstack", "controlPlanePort")

		var found bool

		for _, f := range FilledInTechPreviewFields(&installConfig) {
			if f.String() == expectedField.String() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected field %q to be detected", expectedField)
		}
	})
}
