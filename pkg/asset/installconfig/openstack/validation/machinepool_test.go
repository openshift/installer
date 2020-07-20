package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

var (
	validFlavor = "valid-flavor"
)

func validMachinePool() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validFlavor,
	}
}

func validCloudInfo() *CloudInfo {
	return &CloudInfo{
		PlatformFlavor: &flavors.Flavor{
			Name: validFlavor,
		},
	}
}

func TestOpenStackMachinepoolValidation(t *testing.T) {
	cases := []struct {
		name           string
		mpool          *openstack.MachinePool
		cloudInfo      *CloudInfo
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:           "valid base case",
			mpool:          validMachinePool(),
			cloudInfo:      validCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			aggregatedErrors := ValidateMachinePool(tc.mpool, tc.cloudInfo, field.NewPath("controlPlane", "platform", "openstack")).ToAggregate()
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}
