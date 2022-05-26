package vsphere

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/vim25"
	types2 "github.com/vmware/govmomi/vim25/types"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	validCIDR = "10.0.0.0/16"
)

func validIPIInstallConfig(dcName string, fName string) *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				Cluster:          fmt.Sprintf("%s/%s_C0", fName, dcName),
				Datacenter:       fmt.Sprintf("%s/%s", fName, dcName),
				DefaultDatastore: "LocalDS_0",
				Network:          fmt.Sprintf("%s_DVPG0", dcName),
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid-vcenter",
				APIVIP:           "192.168.111.0",
				IngressVIP:       "192.168.111.1",
			},
		},
	}
}

func validMultiVCenterPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenterSpec{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test-username",
				Password: "test-password",
				Datacenters: []string{
					"DC0",
				},
			},
		},
		DeploymentZones: []vsphere.DeploymentZoneSpec{
			{
				Name:          "test-dz-east-1a",
				Server:        "test-vcenter",
				FailureDomain: "test-east-1a",
				ControlPlane:  "Allowed",
				PlacementConstraint: vsphere.PlacementConstraint{
					ResourcePool: "/DC0/host/DC0_C0/Resources/test-resourcepool",
					Folder:       "/DC0/vm",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomainSpec{
			{
				Name: "test-east-1a",
				Region: vsphere.FailureDomain{
					Name:        "test-region-east",
					Type:        "Datacenter",
					TagCategory: "openshift-region",
				},
				Zone: vsphere.FailureDomain{
					Name:        "test-zone-1a",
					Type:        "ComputeCluster",
					TagCategory: "openshift-zone",
				},
				Topology: vsphere.Topology{
					Datacenter:     "DC0",
					ComputeCluster: "/DC0/host/DC0_C0",
					Hosts:          nil,
					Networks: []string{
						"DC0_DVPG0",
					},
					Datastore: "test-datastore",
				}},
		},
	}
}

func TestValidate(t *testing.T) {
	server := mock.StartSimulator()
	defer server.Close()
	dcName := "DC0"
	fName := "F0"
	dcName1 := "DC1"
	tests := []struct {
		name                      string
		installConfig             *types.InstallConfig
		validationMethod          func(*vim25.Client, Finder, *types.InstallConfig) error
		multiZoneValidationMethod func(*vim25.Client, Finder, *vsphere.FailureDomainSpec, *vsphere.DeploymentZoneSpec) field.ErrorList
		deploymentZone            *vsphere.DeploymentZoneSpec
		failureDomain             *vsphere.FailureDomainSpec
		expectErr                 string
	}{{
		name:             "valid IPI install config",
		installConfig:    validIPIInstallConfig(dcName, ""),
		validationMethod: validateProvisioning,
	}, {
		name:             "valid IPI install config - DC in folder",
		installConfig:    validIPIInstallConfig(dcName1, fName),
		validationMethod: validateProvisioning,
	}, {
		name: "invalid IPI - no network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Network = ""
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform\.vsphere\.network: Required value: must specify the network$`,
	}, {
		name: "invalid IPI - invalid datacenter",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Datacenter = "invalid_dc"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_dc": datacenter './invalid_dc' not found`,
	}, {
		name: "invalid IPI - invalid network",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Network = "invalid_network"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_network": unable to find network provided$`,
	}, {
		name: "invalid IPI - invalid network - DC in folder",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName1, fName)
			c.Platform.VSphere.Network = "invalid_network"
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform.vsphere.network: Invalid value: "invalid_network": unable to find network provided$`,
	}, {
		name: "invalid IPI - no cluster",
		installConfig: func() *types.InstallConfig {
			c := validIPIInstallConfig(dcName, "")
			c.Platform.VSphere.Cluster = ""
			return c
		}(),
		validationMethod: validateProvisioning,
		expectErr:        `^platform\.vsphere\.cluster: Required value: must specify the cluster$`,
	}, {
		name:                      "multi-zone validation",
		deploymentZone:            &validMultiVCenterPlatform().DeploymentZones[0],
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		multiZoneValidationMethod: validateMultiZoneProvisioning,
	}, {
		name:           "multi-zone validation - invalid datacenter",
		deploymentZone: &validMultiVCenterPlatform().DeploymentZones[0],
		failureDomain: func() *vsphere.FailureDomainSpec {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.Datacenter = "invalid-dc"
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology: Invalid value: "invalid-dc": datacenter './invalid-dc' not found$`,
	}, {
		name:           "multi-zone validation - invalid cluster",
		deploymentZone: &validMultiVCenterPlatform().DeploymentZones[0],
		failureDomain: func() *vsphere.FailureDomainSpec {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.ComputeCluster = "/DC0/host/invalid-cluster"
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology.computeCluster: Invalid value: "/DC0/host/invalid-cluster": cluster '/DC0/host/invalid-cluster' not found$`,
	}, {
		name: "multi-zone validation - invalid resource pool",
		deploymentZone: func() *vsphere.DeploymentZoneSpec {
			deploymentZones := &validMultiVCenterPlatform().DeploymentZones[0]
			deploymentZones.PlacementConstraint.ResourcePool = "/DC0/host/DC0_C0/Resources/invalid-resourcepool"
			return deploymentZones
		}(),
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.deploymentZones.placementConstraint.resourcePool: Invalid value: "/DC0/host/DC0_C0/Resources/invalid-resourcepool": resource pool '/DC0/host/DC0_C0/Resources/invalid-resourcepool' not found$`,
	}, {
		name:           "multi-zone validation - invalid network",
		deploymentZone: &validMultiVCenterPlatform().DeploymentZones[0],
		failureDomain: func() *vsphere.FailureDomainSpec {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.Networks = []string{
				"invalid-network",
			}
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology: Invalid value: "invalid-network": unable to find network provided$`,
	}, {
		name: "multi-zone validation - invalid folder",
		deploymentZone: func() *vsphere.DeploymentZoneSpec {
			deploymentZones := &validMultiVCenterPlatform().DeploymentZones[0]
			deploymentZones.PlacementConstraint.Folder = "/DC0/vm/invalid-folder"
			return deploymentZones
		}(),
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.deploymentZones.placementConstraint.folder: Invalid value: "/DC0/vm/invalid-folder": folder '/DC0/vm/invalid-folder' not found$`,
	}}

	finder, err := mock.GetFinder(server)
	if err != nil {
		t.Error(err)
		return
	}
	client, _, err := mock.GetClient(server)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.TODO()
	resourcePools, err := finder.ResourcePoolList(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = resourcePools[0].Create(ctx, "test-resourcepool", types2.DefaultResourceConfigSpec())
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.validationMethod != nil {
				err = test.validationMethod(client, finder, test.installConfig)
			} else if test.multiZoneValidationMethod != nil {
				err = test.multiZoneValidationMethod(client, finder, test.failureDomain, test.deploymentZone).ToAggregate()
			} else {
				err = errors.New("no test method defined")
			}
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
