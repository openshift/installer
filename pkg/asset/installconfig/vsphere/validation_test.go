package vsphere

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	vapitags "github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/mo"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"go.uber.org/mock/gomock"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	validCIDR     = "10.0.0.0/16"
	mu            sync.Mutex
	stopListening = false
)

const (
	tagTestCreateRegionCategory      = 0x01
	tagTestCreateZoneCategory        = 0x02
	tagTestAttachRegionTags          = 0x04
	tagTestAttachZoneTags            = 0x08
	tagTestNothingCreatedOrAttached  = 0x10
	tagTestAttachRegionOnClusterTags = 0x20
	tagTestAttachZoneOnHostsTags     = 0x40
)

func validIPIInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				APIVIPs:     []string{"192.168.111.0"},
				IngressVIPs: []string{"192.168.111.1"},
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
				Region: "test-region-east",
				Zone:   "test-zone-1a",
				Topology: vsphere.Topology{
					Datacenter:     "DC0",
					ComputeCluster: "/DC0/host/DC0_C0",
					ResourcePool:   "/DC0/host/DC0_C0/Resources/test-resourcepool",
					Folder:         "/DC0/vm",
					Networks: []string{
						"DC0_DVPG0",
					},
					Datastore: "/DC0/datastore/LocalDS_0",
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

func createTagAndAttachToType(ctx context.Context, restClient *rest.Client, finder Finder, objectType vsphereObjectType, tagName, categoryID string) error {
	tagMgr := vapitags.NewManager(restClient)
	var refs []mo.Reference

	switch objectType {
	case datacenter:
		objects, err := finder.DatacenterList(ctx, "/...")
		if err != nil {
			return err
		}
		for _, o := range objects {
			refs = append(refs, o.Reference())
		}
	case clusterComputeResource:
		objects, err := finder.ClusterComputeResourceList(ctx, "/...")
		if err != nil {
			return err
		}
		for _, o := range objects {
			refs = append(refs, o.Reference())
		}
	case hostSystem:
		objects, err := finder.HostSystemList(ctx, "/...")
		if err != nil {
			return err
		}
		for _, o := range objects {
			refs = append(refs, o.Reference())
		}
	default:
		return errors.New("unknown object type")
	}

	tagID, err := tagMgr.CreateTag(ctx, &vapitags.Tag{
		Name:       tagName,
		CategoryID: categoryID,
	})
	if err != nil {
		return err
	}

	return tagMgr.AttachTagToMultipleObjects(ctx, tagID, refs)
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
			if err := createTagAndAttachToType(ctx, restClient, finder, datacenter, "us-east", categoryID); err != nil {
				return nil, err
			}
		}
		if attachmentMask&tagTestAttachRegionOnClusterTags != 0 {
			if err := createTagAndAttachToType(ctx, restClient, finder, clusterComputeResource, "us-east", categoryID); err != nil {
				return nil, err
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
			if err := createTagAndAttachToType(ctx, restClient, finder, clusterComputeResource, "us-east-1a", categoryID); err != nil {
				return nil, err
			}
		}
		if attachmentMask&tagTestAttachZoneOnHostsTags != 0 {
			if err := createTagAndAttachToType(ctx, restClient, finder, hostSystem, "us-east-1a", categoryID); err != nil {
				return nil, err
			}
		}
	}

	return tagMgr, nil
}

// simulatorHelper starts the govmomi simulator
// returning the simulator.Server so that we can defer closing later and
// shutdown the simulator to change versions.
func simulatorHelper(t *testing.T, vs *mock.VSphereSimulator) (*validationContext, *simulator.Server, *rest.Client, error) {
	t.Helper()

	if vs == nil {
		vs = mock.NewSimulator("", "", 0, 0)
	}
	server, err := vs.StartSimulator()
	if err != nil {
		return nil, nil, nil, err
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	finder, err := mock.GetFinder(server)
	if err != nil {
		return nil, nil, nil, err
	}

	client, _, err := mock.GetClient(server)
	if err != nil {
		return nil, nil, nil, err
	}

	restClient := rest.NewClient(client)

	defer restClient.CloseIdleConnections()
	err = restClient.Login(context.TODO(), simulator.DefaultLogin)
	if err != nil {
		return nil, nil, nil, err
	}

	vmFolder, err := finder.Folder(ctx, "/DC0/vm")
	if err != nil {
		return nil, nil, nil, err
	}

	myFolder, err := vmFolder.CreateFolder(ctx, "my-folder")
	if err != nil {
		return nil, nil, nil, err
	}
	var folder mo.Folder

	err = myFolder.Properties(ctx, myFolder.Reference(), nil, &folder)
	if err != nil {
		return nil, nil, nil, err
	}

	resourcePools, err := finder.ResourcePoolList(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = resourcePools[0].Create(ctx, "test-resourcepool", vim25types.DefaultResourceConfigSpec())
	if err != nil {
		return nil, nil, nil, err
	}

	sessionMgr := session.NewManager(client)
	userSession, err := sessionMgr.UserSession(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	username := userSession.UserName
	validPermissionsAuthManagerClient, err := buildAuthManagerClient(ctx, t, ctrl, finder, username, nil, nil, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	return &validationContext{
		AuthManager: validPermissionsAuthManagerClient,
		Finder:      finder,
		Client:      client,
	}, server, restClient, nil
}

func setupHostGroup(ctx context.Context, finder Finder, failureDomain *vsphere.FailureDomain) error {
	clusterObj, err := finder.ClusterComputeResource(ctx, failureDomain.Topology.ComputeCluster)
	if err != nil {
		return err
	}

	clusterConfigSpec := &vim25types.ClusterConfigSpecEx{
		GroupSpec: []vim25types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: vim25types.ArrayUpdateSpec{
					Operation: vim25types.ArrayUpdateOperation("add"),
				},
				Info: &vim25types.ClusterHostGroup{
					ClusterGroupInfo: vim25types.ClusterGroupInfo{
						Name: failureDomain.Topology.HostGroup,
					},
				},
			},
		},
	}

	task, err := clusterObj.Reconfigure(ctx, clusterConfigSpec, true)
	if err != nil {
		return err
	}

	return task.Wait(ctx)
}

func TestValidateFailureDomains(t *testing.T) {
	validationCtx, server, restClient, err := simulatorHelper(t, nil)
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(*validationContext, *vsphere.FailureDomain, bool) field.ErrorList
		failureDomain    *vsphere.FailureDomain
		tagTestMask      int64
		checkTags        bool
		expectErr        string
	}{
		{
			name:             "multi-zone validation",
			failureDomain:    &validMultiVCenterPlatform().FailureDomains[0],
			validationMethod: validateFailureDomain,
		}, {
			name: "multi-zone validation - invalid datacenter",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.Datacenter = "invalid-dc"
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `[platform.vsphere.failureDomains.topology.datacenter: Invalid value: "invalid-dc": datacenter 'invalid-dc' not found, platform.vsphere.failureDomains.topology.datastore: Invalid value: "invalid-dc": unable to find datacenter invalid-dc: datacenter 'invalid-dc' not found, platform.vsphere.failureDomains.topology: Invalid value: "invalid-dc": datacenter './invalid-dc' not found]`,
		},
		{
			name: "multi-zone validation - invalid cluster",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.ComputeCluster = "/DC0/host/invalid-cluster"
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `\[platform.vsphere.failureDomains.topology.computeCluster: Invalid value: "/DC0/host/invalid-cluster": cluster '/DC0/host/invalid-cluster' not found, platform.vsphere.failureDomains.topology: Invalid value: "DC0_DVPG0": could not find vSphere cluster at /DC0/host/invalid-cluster: cluster '/DC0/host/invalid-cluster' not found\]`,
		},
		{
			name: "multi-zone validation - missing cluster",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.ComputeCluster = ""
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `[platform.vsphere: Internal error: please specify a datacenter, platform.vsphere.failureDomains.topology.computeCluster: Required value: must specify the cluster]`,
		},
		{
			name: "multi-zone validation - missing datastore",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.Datastore = ""
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `^platform.vsphere.failureDomains.topology.datastore: Required value: must specify the datastore$`,
		},
		{
			name: "multi-zone validation - invalid resource pool",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.ResourcePool = "/DC0/host/DC0_C0/Resources/invalid-resourcepool"
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `^platform.vsphere.failureDomains.topology.resourcePool: Invalid value: "/DC0/host/DC0_C0/Resources/invalid-resourcepool": resource pool '/DC0/host/DC0_C0/Resources/invalid-resourcepool' not found$`,
		},
		{
			name: "multi-zone validation - invalid network",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.Networks = []string{
					"invalid-network",
				}
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `^platform.vsphere.failureDomains.topology: Invalid value: "invalid-network": unable to find network provided$`,
		},
		{
			name: "multi-zone validation - create missing folder",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.Folder = "/DC0/vm/create-missing-folder"
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        ``,
		},
		{
			name:             "multi-zone tag categories present and tags attached",
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiVCenterPlatform().FailureDomains[0],
			checkTags:        true,
			tagTestMask: tagTestCreateZoneCategory |
				tagTestCreateRegionCategory |
				tagTestAttachRegionTags |
				tagTestAttachZoneTags,
		}, {
			name:             "multi-zone tag categories, missing zone tag attachment",
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiVCenterPlatform().FailureDomains[0],
			checkTags:        true,
			tagTestMask: tagTestCreateZoneCategory |
				tagTestCreateRegionCategory |
				tagTestAttachRegionTags,
			expectErr: "platform.vsphere.failureDomains.topology.computeCluster: Invalid value: \"openshift-zone\": tag associated with tag category not attached to this resource or ancestor",
		}, {
			name:             "multi-zone tag categories, missing zone and region tag categories",
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiVCenterPlatform().FailureDomains[0],
			checkTags:        true,
			tagTestMask:      tagTestNothingCreatedOrAttached,
			expectErr:        "platform.vsphere: Internal error: tag categories openshift-zone and openshift-region must be created",
		},
		{
			name: "vm-host group zonal tag categories present and tags attached",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				// The Zone and Region names *MUST* match the tag names
				failureDomain.Zone = "us-east-1a"
				failureDomain.Region = "us-east"
				failureDomain.Topology.HostGroup = "us-east-1a-hostgroup"
				failureDomain.ZoneType = vsphere.HostGroupFailureDomain
				failureDomain.RegionType = vsphere.ComputeClusterFailureDomain
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        true,
			tagTestMask: tagTestCreateZoneCategory |
				tagTestCreateRegionCategory |
				tagTestAttachRegionOnClusterTags |
				tagTestAttachZoneOnHostsTags,
			expectErr: ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.validationMethod != nil {
				var tagMgr *vapitags.Manager

				if test.failureDomain.Topology.HostGroup != "" {
					if err := setupHostGroup(ctx, validationCtx.Finder, test.failureDomain); err != nil {
						assert.NoError(t, err)
					}
				}

				if test.tagTestMask != 0 {
					tagMgr, err = setupTagAttachmentTest(ctx, restClient, validationCtx.Finder, test.tagTestMask)
					if err != nil {
						assert.NoError(t, err)
					}
					validationCtx.zoneTagCategoryID = ""
					validationCtx.regionTagCategoryID = ""
					validationCtx.TagManager = tagMgr
				}

				err = test.validationMethod(validationCtx, test.failureDomain, test.checkTags).ToAggregate()
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

func Test_validateVCenterVersion(t *testing.T) {
	tests := []struct {
		name             string
		VSphereSimulator *mock.VSphereSimulator
		fldPath          *field.Path
		expectErr        string
	}{
		{
			name: "valid vcenter version",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "8.0.0",
				VCenterBuild:   20519528,
				EsxiVersion:    "8.0.0",
				EsxiBuild:      20519528,
			},
			fldPath:   field.NewPath("platform").Child("vsphere").Child("vcenters"),
			expectErr: ``,
		},
		{
			name: "vcf 9 is not supported",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "9.0.0",
				VCenterBuild:   24755230,
				// https://github.com/vmware/govmomi/blob/release-0.52/vim25/types/esxi_version.go
				// is not up-to-date with versions
				EsxiVersion: "8.0.2",
				EsxiBuild:   24859861,
			},
			fldPath:   field.NewPath("platform").Child("vsphere").Child("vcenters"),
			expectErr: `platform.vsphere.vcenters: Required value: Unsupported or untested version of vSphere. Current vCenter version: 9.0.0, build: 24755230`,
		},
		{
			name: "vsphere 7 eol but csi support",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "7.0.2",
				// vCenter Server 7.0 Update 2a
				VCenterBuild: 17920168,
				EsxiVersion:  "7.0.2",
				// ESXi 7.0.2 EP2
				EsxiBuild: 18538813,
			},
			fldPath:   field.NewPath("platform").Child("vsphere").Child("vcenters"),
			expectErr: ``,
		},
		{
			name: "vsphere 6 eol",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "6.7",
				// vCenter Server 6.7 Update 3w
				VCenterBuild: 24337536,
				EsxiVersion:  "6.7",
				// ESXi 6.7 P08
				EsxiBuild: 20497097,
			},
			fldPath:   field.NewPath("platform").Child("vsphere").Child("vcenters"),
			expectErr: `platform.vsphere.vcenters: Required value: Unsupported or untested version of vSphere. Current vCenter version: 6.7, build: 24337536`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationCtx, server, _, err := simulatorHelper(t, test.VSphereSimulator)

			if err != nil {
				t.Error(err)
			}
			defer server.Close()

			err = validateVCenterVersion(validationCtx, test.fldPath).ToAggregate()

			if test.expectErr != "" {
				assert.Regexp(t, test.expectErr, err)
			} else if err != nil {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_validateESXiVersion(t *testing.T) {
	vSphereFldPath := field.NewPath("platform").Child("vsphere")
	computeClusterFldPath := vSphereFldPath.Child("failureDomains").Child("topology").Child("computeCluster")
	platform := validMultiVCenterPlatform()

	tests := []struct {
		name               string
		VSphereSimulator   *mock.VSphereSimulator
		computeClusterPath string
		expectErr          string
	}{
		{
			name: "valid esxi version",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "8.0.0",
				VCenterBuild:   20519528,
				EsxiVersion:    "8.0.0",
				EsxiBuild:      20519528,
			},
			computeClusterPath: platform.FailureDomains[0].Topology.ComputeCluster,
			expectErr:          ``,
		},
		{
			name: "unsupported esxi version",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "6.7",
				// vCenter Server 6.7 Update 3w
				VCenterBuild: 24337536,
				EsxiVersion:  "6.7",
				// ESXi 6.7 P08
				EsxiBuild: 20497097,
			},
			computeClusterPath: platform.FailureDomains[0].Topology.ComputeCluster,
			expectErr:          `[platform.vsphere: Required value: Unsupported or untested version of vSphere.  The ESXi host: DC0_C0_H0 is version: 6.7 and build: 20497097, platform.vsphere: Required value: Unsupported or untested version of vSphere.  The ESXi host: DC0_C0_H1 is version: 6.7 and build: 20497097, platform.vsphere: Required value: Unsupported or untested version of vSphere.  The ESXi host: DC0_C0_H2 is version: 6.7 and build: 20497097, platform.vsphere: Required value: Unsupported or untested version of vSphere.  The ESXi host: DC0_C0_H3 is version: 6.7 and build: 20497097, platform.vsphere: Required value: Unsupported or untested version of vSphere.  The ESXi host: DC0_C0_H4 is version: 6.7 and build: 20497097, platform.vsphere: Required value: Unsupported or untested version of vSphere.  The ESXi host: DC0_C0_H5 is version: 6.7 and build: 20497097]`,
		},
		{
			name: "vsphere 7 eol but csi support",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "7.0.2",
				// vCenter Server 7.0 Update 2a
				VCenterBuild: 17920168,
				EsxiVersion:  "7.0.2",
				// ESXi 7.0.2 EP2
				EsxiBuild: 18538813,
			},
			computeClusterPath: platform.FailureDomains[0].Topology.ComputeCluster,
			expectErr:          ``,
		},
		{
			name: "computeCluster not found",
			VSphereSimulator: &mock.VSphereSimulator{
				VCenterVersion: "6.5.0",
				VCenterBuild:   5973321,
				EsxiVersion:    "6.5.0",
				EsxiBuild:      5973321,
			},
			computeClusterPath: "/DC0/host/invalid-cluster",
			expectErr:          `platform.vsphere.failureDomains.topology.computeCluster: Invalid value: "/DC0/host/invalid-cluster": cluster '/DC0/host/invalid-cluster' not found`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationCtx, server, _, err := simulatorHelper(t, test.VSphereSimulator)

			if err != nil {
				t.Error(err)
			}
			defer server.Close()

			err = validateESXiVersion(validationCtx, test.computeClusterPath, vSphereFldPath, computeClusterFldPath).ToAggregate()

			if test.expectErr != "" {
				assert.Regexp(t, test.expectErr, err)
			} else if err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_ensureDNS(t *testing.T) {
	platformFieldPath := field.NewPath("platform")

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	tests := []struct {
		name          string
		installConfig *types.InstallConfig
		expectErr     string
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ensureDNS(test.installConfig, platformFieldPath, resolver).ToAggregate()
			if test.expectErr != "" {
				assert.Regexp(t, test.expectErr, err)
			} else if err != nil {
				t.Error(err)
			}
		})
	}
}

// lbHelper creates a listening port of 6443 on loopback (127.0.0.1)
// If stopListening is true the loop is broken and the port closes.
func lbHelper(t *testing.T) {
	t.Helper()
	mu.Lock()
	stopListening = false
	mu.Unlock()
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:6443")
	if err != nil {
		t.Error(err)
	}

	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		t.Error(err)
	}
	for {
		mu.Lock()
		if stopListening {
			err = listen.Close()
			if err != nil {
				t.Logf("closing tcp listener: %s", err.Error())
			}
			break
		}
		mu.Unlock()

		_, err = listen.Accept()

		if err != nil {
			t.Error(err)
		}
	}
	mu.Unlock()
}

// Test_ensureLoadBalancer uses lbHelper to create an open
// on 127.0.0.1:6443. Also uses nip.io as the base domain
// and cluster name as 7f000001 which equals 127.0.0.1
// Examples for testing logrus output from here:
// https://maxchadwick.xyz/blog/testing-log-output-in-go-logrus
func Test_ensureLoadBalancer(t *testing.T) {
	go lbHelper(t)
	logger, hook := test.NewNullLogger()
	localLogger = logger

	tests := []struct {
		name          string
		installConfig *types.InstallConfig
		stopTCPListen bool
		expectLevel   string
		expectWarn    string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.stopTCPListen {
				mu.Lock()
				stopListening = true
				mu.Unlock()
			}
			mu.Lock()
			ensureLoadBalancer(test.installConfig)
			mu.Unlock()
			if test.expectWarn != "" {
				// there should be only one entry
				entries := hook.AllEntries()
				assert.NotEmpty(t, entries)
				for _, e := range entries {
					assert.Equal(t, test.expectLevel, e.Level.String())
					assert.Regexp(t, test.expectWarn, e.Message)
				}
			}
			hook.Reset()
		})
	}
}

func Test_compareCurrentToTemplate(t *testing.T) {
	logger, hook := test.NewNullLogger()
	localLogger = logger

	tests := []struct {
		name                   string
		rhcosReleaseVersion    string
		templateProductVersion string
		expectErr              string
		expectWarn             string
	}{
		{
			name:                   "same version",
			rhcosReleaseVersion:    "412.86.202303211731-0",
			templateProductVersion: "412.86.202303211731-0",
			expectErr:              ``,
			expectWarn:             ``,
		},
		{
			name:                   "template newer than current rhcos",
			rhcosReleaseVersion:    "412.86.202303211731-0",
			templateProductVersion: "414.92.202304252144-0",
			expectErr:              `rhcos version: 414.92.202304252144-0 is too many revisions ahead current version: 412.86.202303211731-0`,
			expectWarn:             ``,
		},
		{
			name:                   "rhcos one release newer than template",
			rhcosReleaseVersion:    "413.92.202304252144-0",
			templateProductVersion: "412.86.202303211731-0",
			expectErr:              ``,
			expectWarn:             `rhcos version: 412.86.202303211731-0 is behind current version: 413.92.202304252144-0, installation may fail`,
		},
		{
			name:                   "rhcos two release newer than template",
			rhcosReleaseVersion:    "414.92.202304252144-0",
			templateProductVersion: "412.86.202303211731-0",
			expectErr:              `rhcos version: 412.86.202303211731-0 is too many revisions behind current version: 414.92.202304252144-0`,
			expectWarn:             ``,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := compareCurrentToTemplate(tc.templateProductVersion, tc.rhcosReleaseVersion)

			if tc.expectWarn != "" {
				// there should be only one entry
				entries := hook.AllEntries()
				assert.NotEmpty(t, entries)
				for _, e := range entries {
					assert.Regexp(t, tc.expectWarn, e.Message)
				}
			}

			if tc.expectErr != "" {
				assert.Regexp(t, tc.expectErr, err)
			} else if err != nil {
				t.Error(err)
			}
			hook.Reset()
		})
	}
}
