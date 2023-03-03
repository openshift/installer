package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// FilledInTechPreviewFields returns a slice of field paths that were set in
// install-config, and that are only accepted in the context of a
// TechPreviewNoUpgrade feature set.
func FilledInTechPreviewFields(installConfig *types.InstallConfig) []*field.Path {
	var fields []*field.Path

	if installConfig.OpenStack.LoadBalancer != nil {
		fields = append(fields, field.NewPath("platform", "openstack", "loadBalancer"))
	}

	return fields
}
