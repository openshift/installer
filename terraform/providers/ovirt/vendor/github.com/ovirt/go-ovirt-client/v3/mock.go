package ovirtclient

import (
	"context"
	"math/rand"
	"net"
	"sync"

	"github.com/google/uuid"
)

// MockClient provides in-memory client functions, and additionally provides the ability to inject
// information.
type MockClient interface {
	Client

	// GenerateUUID generates a UUID for testing purposes.
	GenerateUUID() string
}

type mockClient struct {
	ctx                               context.Context
	logger                            Logger
	url                               string
	lock                              *sync.Mutex
	nonSecureRandom                   *rand.Rand
	vms                               map[VMID]*vm
	storageDomains                    map[StorageDomainID]*storageDomain
	disks                             map[DiskID]*diskWithData
	clusters                          map[ClusterID]*cluster
	hosts                             map[HostID]*host
	templates                         map[TemplateID]*template
	nics                              map[NICID]*nic
	vnicProfiles                      map[VNICProfileID]*vnicProfile
	networks                          map[NetworkID]*network
	dataCenters                       map[DatacenterID]*datacenterWithClusters
	vmDiskAttachmentsByVM             map[VMID]map[DiskAttachmentID]*diskAttachment
	vmDiskAttachmentsByDisk           map[DiskID]*diskAttachment
	templateDiskAttachmentsByTemplate map[TemplateID][]*templateDiskAttachment
	templateDiskAttachmentsByDisk     map[DiskID]*templateDiskAttachment
	tags                              map[TagID]*tag
	affinityGroups                    map[ClusterID]map[AffinityGroupID]*affinityGroup
	vmIPs                             map[VMID]map[string][]net.IP
	instanceTypes                     map[InstanceTypeID]*instanceType
	graphicsConsolesByVM              map[VMID][]*vmGraphicsConsole
}

func (m *mockClient) WithContext(ctx context.Context) Client {
	return &mockClient{
		ctx,
		m.logger.WithContext(ctx),
		m.url,
		m.lock,
		m.nonSecureRandom,
		m.vms,
		m.storageDomains,
		m.disks,
		m.clusters,
		m.hosts,
		m.templates,
		m.nics,
		m.vnicProfiles,
		m.networks,
		m.dataCenters,
		m.vmDiskAttachmentsByVM,
		m.vmDiskAttachmentsByDisk,
		m.templateDiskAttachmentsByTemplate,
		m.templateDiskAttachmentsByDisk,
		m.tags,
		m.affinityGroups,
		m.vmIPs,
		m.instanceTypes,
		m.graphicsConsolesByVM,
	}
}

func (m *mockClient) GetContext() context.Context {
	return m.ctx
}

func (m *mockClient) Reconnect() (err error) {
	return nil
}

func (m *mockClient) GetURL() string {
	return m.url
}

func (m *mockClient) GenerateUUID() string {
	return uuid.NewString()
}
