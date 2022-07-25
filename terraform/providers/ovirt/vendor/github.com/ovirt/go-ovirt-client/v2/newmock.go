package ovirtclient

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

// NewMock creates a new in-memory mock client. This client can be used as a testing facility for
// higher level code.
//goland:noinspection GoUnusedExportedFunction
func NewMock() MockClient {
	return NewMockWithLogger(&noopLogger{})
}

// NewMockWithLogger is identical to NewMock, but accepts a logger.
func NewMockWithLogger(logger Logger) MockClient {
	testCluster := generateTestCluster()
	testHost := generateTestHost(testCluster)
	testStorageDomain := generateTestStorageDomain()
	secondaryStorageDomain := generateTestStorageDomain()
	testDatacenter := generateTestDatacenter(testCluster)
	testNetwork := generateTestNetwork(testDatacenter)
	testVNICProfile := generateTestVNICProfile(testNetwork)
	blankTemplate := &template{
		nil,
		DefaultBlankTemplateID,
		"Blank",
		"Blank template",
		TemplateStatusOK,
		&vmCPU{
			&vmCPUTopo{
				cores:   1,
				threads: 1,
				sockets: 1,
			},
			nil,
		},
	}

	client := getClient(
		logger,
		testStorageDomain,
		secondaryStorageDomain,
		testCluster,
		testHost,
		blankTemplate,
		testVNICProfile,
		testNetwork,
		testDatacenter,
	)

	testCluster.client = client
	testHost.client = client
	blankTemplate.client = client
	testStorageDomain.client = client
	secondaryStorageDomain.client = client
	testDatacenter.client = client
	testNetwork.client = client
	testVNICProfile.client = client

	return client
}

func getClient(
	logger Logger,
	testStorageDomain *storageDomain,
	secondaryStorageDomain *storageDomain,
	testCluster *cluster,
	testHost *host,
	blankTemplate *template,
	testVNICProfile *vnicProfile,
	testNetwork *network,
	testDatacenter *datacenterWithClusters,
) *mockClient {
	client := &mockClient{
		ctx:             nil,
		logger:          logger,
		url:             "https://localhost/ovirt-engine/api",
		lock:            &sync.Mutex{},
		vms:             map[VMID]*vm{},
		tags:            map[TagID]*tag{},
		nonSecureRandom: rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
		storageDomains: map[StorageDomainID]*storageDomain{
			testStorageDomain.ID():      testStorageDomain,
			secondaryStorageDomain.ID(): secondaryStorageDomain,
		},
		disks: map[DiskID]*diskWithData{},
		clusters: map[ClusterID]*cluster{
			testCluster.ID(): testCluster,
		},
		hosts: map[HostID]*host{
			testHost.ID(): testHost,
		},
		templates: map[TemplateID]*template{
			blankTemplate.ID(): blankTemplate,
		},
		nics: map[NICID]*nic{},
		vnicProfiles: map[VNICProfileID]*vnicProfile{
			testVNICProfile.ID(): testVNICProfile,
		},
		networks: map[NetworkID]*network{
			testNetwork.ID(): testNetwork,
		},
		dataCenters: map[DatacenterID]*datacenterWithClusters{
			testDatacenter.ID(): testDatacenter,
		},
		vmDiskAttachmentsByVM:   map[VMID]map[DiskAttachmentID]*diskAttachment{},
		vmDiskAttachmentsByDisk: map[DiskID]*diskAttachment{},
		templateDiskAttachmentsByTemplate: map[TemplateID][]*templateDiskAttachment{
			blankTemplate.ID(): {},
		},
		templateDiskAttachmentsByDisk: map[DiskID]*templateDiskAttachment{},
		affinityGroups: map[ClusterID]map[AffinityGroupID]*affinityGroup{
			testCluster.ID(): {},
		},
		vmIPs:                map[VMID]map[string][]net.IP{},
		instanceTypes:        nil,
		graphicsConsolesByVM: map[VMID][]*vmGraphicsConsole{},
	}
	client.instanceTypes = getInstanceTypes(client)
	return client
}

func getInstanceTypes(client *mockClient) map[InstanceTypeID]*instanceType {
	instanceTypes := map[InstanceTypeID]*instanceType{
		"00000009-0009-0009-0009-0000000000f1": {
			client,
			"00000009-0009-0009-0009-0000000000f1",
			"Large",
		},
		"00000007-0007-0007-0007-00000000010a": {
			client,
			"00000007-0007-0007-0007-00000000010a",
			"Medium",
		},
		"00000005-0005-0005-0005-0000000002e6": {
			client,
			"00000005-0005-0005-0005-0000000002e6",
			"Small",
		},
		"00000003-0003-0003-0003-0000000000be": {
			client,
			"00000003-0003-0003-0003-0000000000be",
			"Tiny",
		},
		"0000000b-000b-000b-000b-00000000021f": {
			client,
			"0000000b-000b-000b-000b-00000000021f",
			"XLarge",
		},
	}
	return instanceTypes
}

func generateTestVNICProfile(testNetwork *network) *vnicProfile {
	return &vnicProfile{
		id:        VNICProfileID(uuid.NewString()),
		name:      "test",
		networkID: testNetwork.ID(),
	}
}

func generateTestNetwork(testDatacenter *datacenterWithClusters) *network {
	return &network{
		id:   NetworkID(uuid.NewString()),
		name: "test",
		dcID: testDatacenter.ID(),
	}
}

func generateTestDatacenter(testCluster *cluster) *datacenterWithClusters {
	return &datacenterWithClusters{
		datacenter: datacenter{
			id:   DatacenterID(uuid.NewString()),
			name: "test",
		},
		clusters: []ClusterID{
			testCluster.ID(),
		},
	}
}

func generateTestStorageDomain() *storageDomain {
	return &storageDomain{
		id:             StorageDomainID(uuid.NewString()),
		name:           "Test storage domain",
		available:      10 * 1024 * 1024 * 1024,
		status:         StorageDomainStatusActive,
		externalStatus: StorageDomainExternalStatusNA,
	}
}

func generateTestCluster() *cluster {
	return &cluster{
		id:   ClusterID(uuid.NewString()),
		name: "Test cluster",
	}
}

func generateTestHost(c *cluster) *host {
	return &host{
		id:        HostID(uuid.NewString()),
		clusterID: c.ID(),
		status:    HostStatusUp,
	}
}
