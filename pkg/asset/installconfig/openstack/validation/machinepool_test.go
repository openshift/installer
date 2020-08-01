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
	validZone   = "valid-zone"

	validCtrlPlaneFlavor = "valid-control-plane-flavor"
	validComputeFlavor   = "valid-compute-flavor"

	notExistFlavor = "non-existant-flavor"

	invalidComputeFlavor   = "invalid-compute-flavor"
	invalidCtrlPlaneFlavor = "invalid-control-plane-flavor"
)

func validMachinePool() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validCtrlPlaneFlavor,
	}
}

func validMpoolCloudInfo() *CloudInfo {
	return &CloudInfo{
		Flavors: map[string]*flavors.Flavor{
			validCtrlPlaneFlavor: {
				Name:  validCtrlPlaneFlavor,
				RAM:   16,
				Disk:  25,
				VCPUs: 4,
			},
			validComputeFlavor: {
				Name:  validComputeFlavor,
				RAM:   8,
				Disk:  25,
				VCPUs: 2,
			},
			invalidCtrlPlaneFlavor: {
				Name:  invalidCtrlPlaneFlavor,
				RAM:   8, // too low
				Disk:  25,
				VCPUs: 2, // too low
			},
			invalidComputeFlavor: {
				Name:  invalidComputeFlavor,
				RAM:   8,
				Disk:  10, // too low
				VCPUs: 2,
			},
		},
		Zones: []string{
			validZone,
		},
	}
}

func TestOpenStackMachinepoolValidation(t *testing.T) {
	cases := []struct {
		name           string
		controlPlane   bool // only matters for flavor
		mpool          *openstack.MachinePool
		cloudInfo      *CloudInfo
		expectedError  bool
		expectedErrMsg string // NOTE: this is a REGEXP
	}{
		{
			name:           "valid control plane",
			controlPlane:   true,
			mpool:          validMachinePool(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name: "valid zone",
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.Zones = []string{validZone}
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name: "invalid zone",
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.Zones = []string{"invalid-zone"}
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: "Zone either does not exist in this cloud, or is not available",
		},
		{
			name:           "valid compute",
			controlPlane:   false,
			mpool:          validMachinePool(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:         "not found control plane flavorName",
			controlPlane: true,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = notExistFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: "controlPlane.platform.openstack.flavorName: Not found: \"non-existant-flavor\"",
		},
		{
			name: "not found compute flavorName",
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = notExistFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.flavorName: Not found: "non-existant-flavor"`,
		},
		{
			name:         "invalid control plane flavorName",
			controlPlane: true,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = invalidCtrlPlaneFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: "controlPlane.platform.openstack.flavorName: Invalid value: \"invalid-control-plane-flavor\": Flavor did not meet the following minimum requirements: Must have minimum of 16 GB RAM, had 8 GB; Must have minimum of 4 VCPUs, had 2",
		},
		{
			name:         "invalid compute flavorName",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = invalidComputeFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.flavorName: Invalid value: "invalid-compute-flavor": Flavor did not meet the following minimum requirements: Must have minimum of 25 GB Disk, had 10 GB`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var fieldPath *field.Path
			if tc.controlPlane {
				fieldPath = field.NewPath("controlPlane", "platform", "openstack")
			} else {
				fieldPath = field.NewPath("compute").Index(0).Child("platform", "openstack")
			}

			aggregatedErrors := ValidateMachinePool(tc.mpool, tc.cloudInfo, tc.controlPlane, fieldPath).ToAggregate()
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}
