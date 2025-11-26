package validation

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	validZone   = "valid-zone"
	invalidZone = "invalid-zone"

	validCtrlPlaneFlavor = "valid-control-plane-flavor"
	validComputeFlavor   = "valid-compute-flavor"

	notExistFlavor = "non-existant-flavor"

	invalidComputeFlavor   = "invalid-compute-flavor"
	invalidCtrlPlaneFlavor = "invalid-control-plane-flavor"
	warningComputeFlavor   = "warning-compute-flavor"
	warningCtrlPlaneFlavor = "warning-control-plane-flavor"

	baremetalFlavor = "baremetal-flavor"

	invalidType      = "invalid-type"
	volumeSmallSize  = 10
	volumeMediumSize = 40
	volumeLargeSize  = 100
)

var volumeTypes = []string{"performance", "standard"}
var invalidVolumeTypes = []string{"performance", "invalid-type"}
var volumeType = volumeTypes[0]

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
			Size:  volumeSmallSize,
			Types: volumeTypes,
			Zones: []string{""},
		},
	}
}

func warningMachinePoolMediumVolume() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validCtrlPlaneFlavor,
		Zones:      []string{""},
		RootVolume: &openstack.RootVolume{
			Size:  volumeMediumSize,
			Types: volumeTypes,
			Zones: []string{""},
		},
	}
}

func validMachinePoolLargeVolume() *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: validCtrlPlaneFlavor,
		Zones:      []string{""},
		RootVolume: &openstack.RootVolume{
			Size:  volumeLargeSize,
			Types: volumeTypes,
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
					RAM:   16384,
					Disk:  100,
					VCPUs: 4,
				},
			},
			validComputeFlavor: {
				Flavor: flavors.Flavor{
					Name:  validComputeFlavor,
					RAM:   8192,
					Disk:  100,
					VCPUs: 2,
				},
			},
			invalidCtrlPlaneFlavor: {
				Flavor: flavors.Flavor{
					Name:  invalidCtrlPlaneFlavor,
					RAM:   8192, // too low
					Disk:  100,
					VCPUs: 2, // too low
				},
			},
			invalidComputeFlavor: {
				Flavor: flavors.Flavor{
					Name:  invalidComputeFlavor,
					RAM:   8192,
					Disk:  10, // too low
					VCPUs: 2,
				},
			},
			warningCtrlPlaneFlavor: {
				Flavor: flavors.Flavor{
					Name:  warningCtrlPlaneFlavor,
					RAM:   16384,
					Disk:  40, // not recommended
					VCPUs: 4,
				},
			},
			warningComputeFlavor: {
				Flavor: flavors.Flavor{
					Name:  invalidComputeFlavor,
					RAM:   8192,
					Disk:  40, // not recommended
					VCPUs: 2,
				},
			},
			baremetalFlavor: {
				Flavor: flavors.Flavor{
					Name:  baremetalFlavor,
					RAM:   8192, // too low
					Disk:  10,   // too low
					VCPUs: 2,    // too low
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
		VolumeTypes: volumeTypes,
	}
}

func TestOpenStackMachinepoolValidation(t *testing.T) {
	cases := []struct {
		name            string
		controlPlane    bool // only matters for flavor
		mpool           *openstack.MachinePool
		cloudInfo       *CloudInfo
		expectedError   bool
		expectedErrMsg  string // NOTE: this is a REGEXP
		expectedWarnMsg string //NOTE: this is a REGEXP
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
			expectedErrMsg: "controlPlane.platform.openstack.type: Invalid value: \"invalid-control-plane-flavor\": Flavor did not meet the following minimum requirements: Must have minimum of 16384 MB RAM, had 8192 MB; Must have minimum of 4 VCPUs, had 2",
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
			name:         "warning control plane flavorName",
			controlPlane: true,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = warningCtrlPlaneFlavor
				return mp
			}(),
			cloudInfo:       validMpoolCloudInfo(),
			expectedWarnMsg: `Flavor does not meet the following recommended requirements: It is recommended to have 100 GB Disk, had 40 GB`,
		},
		{
			name:         "warning compute flavorName",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePool()
				mp.FlavorName = warningComputeFlavor
				return mp
			}(),
			cloudInfo:       validMpoolCloudInfo(),
			expectedWarnMsg: `Flavor does not meet the following recommended requirements: It is recommended to have 100 GB Disk, had 40 GB`,
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
			expectedErrMsg: "Volume size must be greater than 25 GB to use root volumes, had 10 GB",
		},
		{
			name:         "volume not recommended",
			controlPlane: false,
			mpool: func() *openstack.MachinePool {
				mp := warningMachinePoolMediumVolume()
				mp.FlavorName = invalidCtrlPlaneFlavor
				return mp
			}(),
			cloudInfo:       validMpoolCloudInfo(),
			expectedWarnMsg: "Volume size is recommended to be greater than 100 GB to use root volumes, had 40 GB",
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
			expectedErrMsg: `compute\[0\].platform.openstack.rootVolume.zones.zone\[0\]: Invalid value: "AZ1": Zone either does not exist in this cloud, or is not available, compute\[0\].platform.openstack.rootVolume.zones.zone\[1\]: Invalid value: "AZ2": Zone either does not exist in this cloud, or is not available, compute\[0\].platform.openstack.rootVolume.zones: Invalid value: \["AZ1","AZ2"\]: there must be either just one volume availability zone common to all nodes or the number of compute and volume availability zones must be equal`,
		},
		{
			name:         "invalid volume types",
			controlPlane: true,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.Types = invalidVolumeTypes
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  true,
			expectedErrMsg: "controlPlane.platform.openstack.rootVolume.types: Invalid value: \"invalid-type\": Volume type either does not exist in this cloud, or is not available",
		},
		{
			name:         "valid volume type",
			controlPlane: true,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.DeprecatedType = volumeType
				mp.RootVolume.Types = []string{}
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:         "valid volume types",
			controlPlane: true,
			mpool: func() *openstack.MachinePool {
				mp := validMachinePoolLargeVolume()
				mp.RootVolume.Types = volumeTypes
				return mp
			}(),
			cloudInfo:      validMpoolCloudInfo(),
			expectedError:  false,
			expectedErrMsg: "",
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

			hook := logrusTest.NewGlobal()
			aggregatedErrors := ValidateMachinePool(tc.mpool, tc.cloudInfo, tc.controlPlane, fieldPath).ToAggregate()
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
			if len(tc.expectedWarnMsg) > 0 {
				assert.Regexp(t, tc.expectedWarnMsg, hook.LastEntry().Message)
			}
		})
	}
}
