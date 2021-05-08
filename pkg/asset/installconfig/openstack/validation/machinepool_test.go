package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

const (
	validZone   = "valid-zone"
	invalidZone = "invalid-zone"

	validCtrlPlaneFlavor = "valid-control-plane-flavor"
	validComputeFlavor   = "valid-compute-flavor"

	notExistFlavor = "non-existant-flavor"

	invalidComputeFlavor   = "invalid-compute-flavor"
	invalidCtrlPlaneFlavor = "invalid-control-plane-flavor"

	baremetalFlavor = "baremetal-flavor"

	volumeType      = "performance"
	invalidType     = "invalid-type"
	volumeSmallSize = 10
	volumeLargeSize = 25
)

func validMachinePool() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validCtrlPlaneFlavor,
		Zones:      []string{""},
	}
}

func invalidMachinePoolSmallVolume() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validCtrlPlaneFlavor,
		Zones:      []string{""},
		RootVolume: &openstack.RootVolume{
			Type:  volumeType,
			Size:  volumeSmallSize,
			Zones: []string{""},
		},
	}
}

func validMachinePoolLargeVolume() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validCtrlPlaneFlavor,
		Zones:      []string{""},
		RootVolume: &openstack.RootVolume{
			Type:  volumeType,
			Size:  volumeLargeSize,
			Zones: []string{validZone},
		},
	}
}

func validMpoolCloudInfo() *CloudInfo {
	return &CloudInfo{
		Flavors: map[string]Flavor{
			validCtrlPlaneFlavor: {
				Flavor: flavors.Flavor{
					Name:  validCtrlPlaneFlavor,
					RAM:   16,
					Disk:  25,
					VCPUs: 4,
				},
			},
			validComputeFlavor: {
				Flavor: flavors.Flavor{
					Name:  validComputeFlavor,
					RAM:   8,
					Disk:  25,
					VCPUs: 2,
				},
			},
			invalidCtrlPlaneFlavor: {
				Flavor: flavors.Flavor{
					Name:  invalidCtrlPlaneFlavor,
					RAM:   8, // too low
					Disk:  25,
					VCPUs: 2, // too low
				},
			},
			invalidComputeFlavor: {
				Flavor: flavors.Flavor{
					Name:  invalidComputeFlavor,
					RAM:   8,
					Disk:  10, // too low
					VCPUs: 2,
				},
			},
			baremetalFlavor: {
				Flavor: flavors.Flavor{
					Name:  baremetalFlavor,
					RAM:   8,  // too low
					Disk:  10, // too low
					VCPUs: 2,  // too low
				},
				Baremetal: true,
			},
		},
		ComputeZones: []string{
			validZone,
		},
		VolumeZones: []string{
			validZone,
		},
		VolumeTypes: []string{
			volumeType,
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
			cloudInfo: func() *CloudInfo {
				ci := validMpoolCloudInfo()
				return ci
			}(),
			expectedError:  true,
			expectedErrMsg: "controlPlane.platform.openstack.type: Not found: \"non-existant-flavor\"",
		},
		{
			name: "not found compute flavorName",
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = notExistFlavor
				return mp
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validMpoolCloudInfo()
				return ci
			}(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.type: Not found: "non-existant-flavor"`,
		},
		{
			name: "no flavor name",
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = ""
				return mp
			}(),
			cloudInfo: func() *CloudInfo {
				ci := validMpoolCloudInfo()
				return ci
			}(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.type: Required value: Flavor name must be provided`,
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
			expectedErrMsg: "controlPlane.platform.openstack.type: Invalid value: \"invalid-control-plane-flavor\": Flavor did not meet the following minimum requirements: Must have minimum of 16 GB RAM, had 8 GB; Must have minimum of 4 VCPUs, had 2",
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
			expectedErrMsg: `compute\[0\].platform.openstack.type: Invalid value: "invalid-compute-flavor": Flavor did not meet the following minimum requirements: Must have minimum of 25 GB Disk, had 10 GB`,
		},
		{
			name:         "valid baremetal compute",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = baremetalFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:         "volume too small",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := invalidMachinePoolSmallVolume()
				mp.FlavorName = invalidCtrlPlaneFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: "Volume size must be greater than 25 to use root volumes, had 10",
		},
		{
			name:         "volume big enough",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.FlavorName = invalidCtrlPlaneFlavor
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:         "valid root volume az",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:         "invalid root volume az",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.Zones = []string{invalidZone}
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.rootVolume.zones.zone\[0\]: Invalid value: \"invalid-zone\": Zone either does not exist in this cloud, or is not available`,
		},
		{
			name:         "volume and compute zones number mismatch",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.Zones = []string{"AZ1", "AZ2"}
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.rootVolume.zones: Invalid value: \[\]string{"AZ1", "AZ2"}: there must be either just one volume availability zone common to all nodes or the number of compute and volume availability zones must be equal`,
		},
		{
			name:         "empty volume type",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.Type = ""
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.rootVolume.type: Invalid value: \"\": Volume type must be specified to use root volumes`,
		},
		{
			name:         "invalid volume type",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.Type = invalidType
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: `compute\[0\].platform.openstack.rootVolume.type: Invalid value: \"invalid-type\": Volume Type either does not exist in this cloud, or is not available`,
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
