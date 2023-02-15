package vsphere

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	vapitags "github.com/vmware/govmomi/vapi/tags"
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

const (
	tagTestCreateRegionCategory     = 0x01
	tagTestCreateZoneCategory       = 0x02
	tagTestAttachRegionTags         = 0x04
	tagTestAttachZoneTags           = 0x08
	tagTestNothingCreatedOrAttached = 0x10
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

func teardownTagAttachmentTest(ctx context.Context, tagMgr *vapitags.Manager) error {
	tags, err := tagMgr.ListTags(ctx)
	if err != nil {
		return err
	}
	attachedMos, err := tagMgr.GetAttachedObjectsOnTags(ctx, tags)
	if err != nil {
		return err
	}

	for _, attachedMo := range attachedMos {
		for _, mo := range attachedMo.ObjectIDs {
			err := tagMgr.DetachTag(ctx, attachedMo.TagID, mo)
			if err != nil {
				return err
			}
		}
		err := tagMgr.DeleteTag(ctx, attachedMo.Tag)
		if err != nil {
			return err
		}
	}

	categories, err := tagMgr.GetCategories(ctx)
	if err != nil {
		return err
	}
	for _, category := range categories {
		cat := category
		err := tagMgr.DeleteCategory(ctx, &cat)
		if err != nil {
			return err
		}
	}
	return nil
}

func setupTagAttachmentTest(ctx context.Context, restClient *rest.Client, finder Finder, attachmentMask int64) (*vapitags.Manager, error) {
	tagMgr := vapitags.NewManager(restClient)

	if attachmentMask&tagTestCreateRegionCategory != 0 {
		categoryID, err := tagMgr.CreateCategory(ctx, &vapitags.Category{
			Name:        "openshift-region",
			Description: "region tag category",
		})
		if err != nil {
			return nil, err
		}

		if attachmentMask&tagTestAttachRegionTags != 0 {
			tagID, err := tagMgr.CreateTag(ctx, &vapitags.Tag{
				Name:       "us-east",
				CategoryID: categoryID,
			})
			if err != nil {
				return nil, err
			}
			datacenters, err := finder.DatacenterList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datacenter := range datacenters {
				err = tagMgr.AttachTag(ctx, tagID, datacenter)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	if attachmentMask&tagTestCreateZoneCategory != 0 {
		categoryID, err := tagMgr.CreateCategory(ctx, &vapitags.Category{
			Name:        "openshift-zone",
			Description: "zone tag category",
		})
		if err != nil {
			return nil, err
		}
		if attachmentMask&tagTestAttachZoneTags != 0 {
			tagID, err := tagMgr.CreateTag(ctx, &vapitags.Tag{
				Name:       "us-east-1a",
				CategoryID: categoryID,
			})
			if err != nil {
				return nil, err
			}
			clusters, err := finder.ClusterComputeResourceList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, cluster := range clusters {
				err = tagMgr.AttachTag(ctx, tagID, cluster)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return tagMgr, nil
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
		tagTestMask               int64
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
	}, {
		name:                      "multi-zone tag categories present and tags attached",
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		tagTestMask: tagTestCreateZoneCategory |
			tagTestCreateRegionCategory |
			tagTestAttachRegionTags |
			tagTestAttachZoneTags,
	}, {
		name:                      "multi-zone tag categories, missing zone tag attachment",
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		tagTestMask: tagTestCreateZoneCategory |
			tagTestCreateRegionCategory |
			tagTestAttachRegionTags,
		expectErr: "platform.vsphere.failureDomains.topology.computeCluster: Internal error: tag associated with tag category openshift-zone not attached to this resource or ancestor",
	}, {
		name:                      "multi-zone tag categories, missing zone and region tag categories",
		multiZoneValidationMethod: validateMultiZoneProvisioning,
		failureDomain:             &validMultiVCenterPlatform().FailureDomains[0],
		tagTestMask:               tagTestNothingCreatedOrAttached,
		expectErr:                 "platform.vsphere: Internal error: tag categories openshift-zone and openshift-region must be created",
	},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

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

	restClient := rest.NewClient(client)

	defer restClient.CloseIdleConnections()
	err = restClient.Login(context.TODO(), simulator.DefaultLogin)
	if err != nil {
		t.Error(err)
	}

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
				var tagMgr *vapitags.Manager
				if test.tagTestMask != 0 {
					tagMgr, err = setupTagAttachmentTest(ctx, restClient, finder, test.tagTestMask)
					if err != nil {
						assert.NoError(t, err)
					}
					validationCtx.zoneTagCategoryID = ""
					validationCtx.regionTagCategoryID = ""
					validationCtx.TagManager = tagMgr
				}
				err = test.multiZoneValidationMethod(validationCtx, test.failureDomain).ToAggregate()
				if test.tagTestMask != 0 {
					err := teardownTagAttachmentTest(ctx, tagMgr)
					if err != nil {
						assert.NoError(t, err)
					}
					validationCtx.zoneTagCategoryID = ""
					validationCtx.regionTagCategoryID = ""
					validationCtx.TagManager = nil
				}
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
