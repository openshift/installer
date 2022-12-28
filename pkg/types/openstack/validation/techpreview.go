package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// FilledInTechPreviewFields returns a slice of field paths that were set in
// install-config, and that are only accepted in the context of a
// TechPreviewNoUpgrade feature set.
func FilledInTechPreviewFields(installConfig *types.InstallConfig) (fields []*field.Path) {
	if installConfig == nil {
		return nil
	}

	if installConfig.OpenStack.ControlPlanePort != nil && installConfig.OpenStack.DeprecatedMachinesSubnet == "" {
		fields = append(fields, field.NewPath("platform", "openstack", "controlPlanePort"))
	}

	return fields
}
