package validation

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func TestFilledInTechPreviewFields(t *testing.T) {
	t.Run("load_balancer", func(t *testing.T) {
		installConfig := types.InstallConfig{
			Platform: types.Platform{
				OpenStack: &openstack.Platform{
					LoadBalancer: &v1.OpenStackPlatformLoadBalancer{
						Type: v1.LoadBalancerTypeUserManaged,
					},
				},
			},
		}

		expectedField := field.NewPath("platform", "openstack", "loadBalancer")

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
