package vsphere

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
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
	validCIDR     = "10.0.0.0/16"
	mu            sync.Mutex
	stopListening = false
)

const (
	tagTestCreateRegionCategory     = 0x01
	tagTestCreateZoneCategory       = 0x02
	tagTestAttachRegionTags         = 0x04
	tagTestAttachZoneTags           = 0x08
	tagTestNothingCreatedOrAttached = 0x10
)

const wildcardDNS = "nip.io"

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

// simulatorHelper starts the govmomi simulator
// returning the simulator.Server so that we can defer closing later and
// shutdown the simulator to change versions.
func simulatorHelper(t *testing.T, setVersionToSupported bool) (*validationContext, *simulator.Server, *rest.Client, error) {
	t.Helper()

	server, err := mock.StartSimulator(setVersionToSupported)
	if err != nil {
		return nil, nil, nil, err
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
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

	rootFolder := object.NewRootFolder(client)
	_, err = rootFolder.CreateFolder(ctx, "/DC0/vm/my-folder")
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

	// Create cluster w/ no hosts for the ComputeResource failure case.  Cluster will have name Cluster1
	dc, err := finder.Datacenter(ctx, "/DC0")
	if err != nil {
		return nil, nil, nil, err
	}
	folders, err := dc.Folders(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = folders.HostFolder.CreateCluster(ctx, "Cluster1", vim25types.ClusterConfigSpecEx{})
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

func TestValidateFailureDomains(t *testing.T) {
	validationCtx, server, restClient, err := simulatorHelper(t, true)
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
			expectErr:        `^platform.vsphere.failureDomains.topology.computeCluster: Invalid value: "": full path of cluster is required$`,
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
		}, {
			name: "multi-zone validation - invalid folder",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.Folder = "/DC0/vm/invalid-folder"
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `^platform.vsphere.failureDomains.topology.folder: Invalid value: "/DC0/vm/invalid-folder": folder '/DC0/vm/invalid-folder' not found$`,
		},

		{
			name: "multi-zone validation - invalid folder",
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := &validMultiVCenterPlatform().FailureDomains[0]
				failureDomain.Topology.Folder = "/DC0/vm/invalid-folder"
				return failureDomain
			}(),
			validationMethod: validateFailureDomain,
			checkTags:        false,
			expectErr:        `^platform.vsphere.failureDomains.topology.folder: Invalid value: "/DC0/vm/invalid-folder": folder '/DC0/vm/invalid-folder' not found$`,
		}, {
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
			expectErr: "platform.vsphere.failureDomains.topology.computeCluster: Internal error: tag associated with tag category openshift-zone not attached to this resource or ancestor",
		}, {
			name:             "multi-zone tag categories, missing zone and region tag categories",
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiVCenterPlatform().FailureDomains[0],
			checkTags:        true,
			tagTestMask:      tagTestNothingCreatedOrAttached,
			expectErr:        "platform.vsphere: Internal error: tag categories openshift-zone and openshift-region must be created",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.validationMethod != nil {
				var tagMgr *vapitags.Manager

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
		name                  string
		setVersionToSupported bool
		fldPath               *field.Path
		expectErr             string
	}{
		{
			name:                  "valid vcenter version",
			setVersionToSupported: true,
			fldPath:               field.NewPath("platform").Child("vsphere").Child("vcenters"),
			expectErr:             ``,
		},
		{
			name:                  "unsupported vcenter version",
			setVersionToSupported: false,
			fldPath:               field.NewPath("platform").Child("vsphere").Child("vcenters"),
			expectErr:             `platform.vsphere.vcenters: Required value: The vSphere storage driver requires a minimum of vSphere 7 Update 2. Current vCenter version: 6.5.0, build: 5973321`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationCtx, server, _, err := simulatorHelper(t, test.setVersionToSupported)

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
		name                  string
		setVersionToSupported bool
		computeClusterPath    string
		expectErr             string
	}{
		{
			name:                  "valid esxi version",
			setVersionToSupported: true,
			computeClusterPath:    platform.FailureDomains[0].Topology.ComputeCluster,
			expectErr:             ``,
		},
		{
			name:                  "unsupported esxi version",
			setVersionToSupported: false,
			computeClusterPath:    platform.FailureDomains[0].Topology.ComputeCluster,
			expectErr:             `[platform.vsphere.failureDomains.topology.computeCluster: Required value: The vSphere storage driver requires a minimum of vSphere 7 Update 2. The ESXi host: DC0_C0_H0 is version: 6.5.0 and build: 5969303, platform.vsphere.failureDomains.topology.computeCluster: Required value: The vSphere storage driver requires a minimum of vSphere 7 Update 2. The ESXi host: DC0_C0_H1 is version: 6.5.0 and build: 5969303, platform.vsphere.failureDomains.topology.computeCluster: Required value: The vSphere storage driver requires a minimum of vSphere 7 Update 2. The ESXi host: DC0_C0_H2 is version: 6.5.0 and build: 5969303]`,
		},
		{
			name:                  "computeCluster not found",
			setVersionToSupported: true,
			computeClusterPath:    "/DC0/host/invalid-cluster",
			expectErr:             `platform.vsphere.failureDomains.topology.computeCluster: Invalid value: "/DC0/host/invalid-cluster": cluster '/DC0/host/invalid-cluster' not found`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationCtx, server, _, err := simulatorHelper(t, test.setVersionToSupported)

			if err != nil {
				t.Error(err)
			}
			defer server.Close()

			err = validateClusterComputeResources(validationCtx, test.computeClusterPath, vSphereFldPath, computeClusterFldPath, false, false).ToAggregate()

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
	}{
		{
			name: "valid dns",
			installConfig: func() *types.InstallConfig {
				installConfig := validIPIInstallConfig()
				installConfig.ObjectMeta.Name = "0a000803"
				installConfig.BaseDomain = wildcardDNS

				return installConfig
			}(),
			expectErr: ``,
		},
	}
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
	}{
		{
			name: "valid lb",
			installConfig: func() *types.InstallConfig {
				installConfig := validIPIInstallConfig()
				installConfig.ObjectMeta.Name = "7f000001"
				installConfig.BaseDomain = wildcardDNS
				installConfig.VSphere.APIVIPs = []string{}
				installConfig.VSphere.IngressVIPs = []string{}

				return installConfig
			}(),
			stopTCPListen: false,
			expectLevel:   ``,
			expectWarn:    ``,
		},
		{
			name: "warn lb",
			installConfig: func() *types.InstallConfig {
				installConfig := validIPIInstallConfig()
				installConfig.ObjectMeta.Name = "7f000002"
				installConfig.BaseDomain = wildcardDNS
				installConfig.VSphere.APIVIPs = []string{}
				installConfig.VSphere.IngressVIPs = []string{}

				return installConfig
			}(),
			stopTCPListen: true,
			expectLevel:   `warning`,
			expectWarn:    `Installation may fail, load balancer not available: dial tcp 127.0.0.2:6443: connect: connection refused`,
		},
	}
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

// Test_validateCheckClusterComputeResource test the various paths for the
// validateClusterComputeResourceHasHosts method.
func Test_validateCheckClusterComputeResource(t *testing.T) {
	vSphereFldPath := field.NewPath("platform").Child("vsphere")
	computeClusterFldPath := vSphereFldPath.Child("failureDomains").Child("topology").Child("computeCluster")
	validationCtx, server, _, err := simulatorHelper(t, true)
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Close()

	tests := []struct {
		name          string
		installConfig *types.InstallConfig
		dataCenter    string
		cluster       string
		expectErr     string
	}{{
		name:       "validateClusterComputeResourceHasHosts valid",
		dataCenter: "DC0",
		cluster:    "DC0_C0",
	}, {
		name:       "validateClusterComputeResourceHasHosts no ComputeResources",
		dataCenter: "DC0",
		cluster:    "Cluster1",
		expectErr:  `^platform.vsphere.failureDomains.topology.computeCluster: Internal error: no hosts found in cluster /DC0/host/Cluster1$`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// err := validateClusterComputeResourceHasHosts(ctx, test.dataCenter, test.cluster, validationCtx.Finder)
			err := validateClusterComputeResources(validationCtx, fmt.Sprintf("/%v/host/%v", test.dataCenter, test.cluster), vSphereFldPath, computeClusterFldPath, false, false).ToAggregate()
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
