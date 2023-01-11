package vsphere

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	vim25types "github.com/vmware/govmomi/vim25/types"
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
				ResourcePool:     "/DC0/host/DC0_C0/Resources/test-resourcepool",
				Network:          fmt.Sprintf("%s_DVPG0", dcName),
				Password:         "valid_password",
				Username:         "valid_username",
				VCenter:          "valid-vcenter",
				APIVIPs:          []string{"192.168.111.0"},
				IngressVIPs:      []string{"192.168.111.1"},
			},
		},
	}
}

func validMultiVCenterPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenter{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test_username",
				Password: "test_password",
				Datacenters: []string{
					"DC0",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "test-east-1a",
				Region: "test-region--east",
				Zone:   "test-zone-1a",
				Topology: vsphere.Topology{
					Datacenter:     "DC0",
					ComputeCluster: "/DC0/host/DC0_C0",
					ResourcePool:   "/DC0/host/DC0_C0/Resources/test-resourcepool",
					Folder:         "/DC0/vm",
					Networks: []string{
						"DC0_DVPG0",
					},
					Datastore: "LocalDS_0",
				}},
		},
	}
}

func TestValidate(t *testing.T) {
	server := mock.StartSimulator()
	defer server.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dcName := "DC0"
	fName := "/F0"
	dcName1 := "DC1"
	tests := []struct {
		name                      string
		installConfig             *types.InstallConfig
		validationMethod          func(*validationContext, *types.InstallConfig) error
		multiZoneValidationMethod func(*validationContext, *vsphere.FailureDomain) field.ErrorList
		failureDomain             *vsphere.FailureDomain
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
		expectErr:        `^platform.vsphere.datacenter: Invalid value: "invalid_dc": datacenter 'invalid_dc' not found`,
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
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		multiZoneValidationMethod: validateMultiZoneProvisioning,
	}, {
		name: "multi-zone validation - invalid datacenter",
		failureDomain: func() *vsphere.FailureDomain {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.Datacenter = "invalid-dc"
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology.datacenter: Invalid value: "invalid-dc": datacenter 'invalid-dc' not found$`,
	}, {
		name: "multi-zone validation - invalid cluster",
		failureDomain: func() *vsphere.FailureDomain {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.ComputeCluster = "/DC0/host/invalid-cluster"
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology.computeCluster: Invalid value: "/DC0/host/invalid-cluster": cluster '/DC0/host/invalid-cluster' not found$`,
	}, {
		name: "multi-zone validation - invalid resource pool",
		failureDomain: func() *vsphere.FailureDomain {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.ResourcePool = "/DC0/host/DC0_C0/Resources/invalid-resourcepool"
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology.resourcePool: Invalid value: "/DC0/host/DC0_C0/Resources/invalid-resourcepool": resource pool '/DC0/host/DC0_C0/Resources/invalid-resourcepool' not found$`,
	}, {
		name: "multi-zone validation - invalid network",
		failureDomain: func() *vsphere.FailureDomain {
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
		failureDomain: func() *vsphere.FailureDomain {
			failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
			failureDomain.Topology.Folder = "/DC0/vm/invalid-folder"
			return failureDomain
		}(),
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		expectErr:                 `^platform.vsphere.failureDomains.topology.folder: Invalid value: "/DC0/vm/invalid-folder": folder '/DC0/vm/invalid-folder' not found$`,
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
	rootFolder := object.NewRootFolder(client)
	_, err = rootFolder.CreateFolder(ctx, "/DC0/vm/my-folder")
	if err != nil {
		t.Error(err)
	}

	resourcePools, err := finder.ResourcePoolList(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = resourcePools[0].Create(ctx, "test-resourcepool", vim25types.DefaultResourceConfigSpec())
	if err != nil {
		t.Error(err)
		return
	}

	sessionMgr := session.NewManager(client)
	userSession, err := sessionMgr.UserSession(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	username := userSession.UserName
	validPermissionsAuthManagerClient, err := buildAuthManagerClient(ctx, ctrl, finder, username, nil, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}

	validationCtx := &validationContext{
		AuthManager: validPermissionsAuthManagerClient,
		Finder:      finder,
		Client:      client,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.validationMethod != nil {
				err = test.validationMethod(validationCtx, test.installConfig)
			} else if test.multiZoneValidationMethod != nil {
				err = test.multiZoneValidationMethod(validationCtx, test.failureDomain).ToAggregate()
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
