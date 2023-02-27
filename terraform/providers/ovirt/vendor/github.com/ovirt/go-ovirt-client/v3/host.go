package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "Host" -n "host" -T HostID

// HostID is the identifier for hosts.
type HostID string

// HostClient contains the API portion that deals with hosts.
type HostClient interface {
	ListHosts(retries ...RetryStrategy) ([]Host, error)
	GetHost(id HostID, retries ...RetryStrategy) (Host, error)
}

// HostData is the core of Host, providing only data access functions.
type HostData interface {
	// ID returns the identifier of the host in question.
	ID() HostID
	// ClusterID returns the ID of the cluster this host belongs to.
	ClusterID() ClusterID
	// Status returns the status of this host.
	Status() HostStatus
}

// Host is the representation of a host returned from the oVirt Engine API. Hosts, also known as hypervisors, are the
// physical servers on which virtual machines run. Full virtualization is provided by using a loadable Linux kernel
// module called Kernel-based Virtual Machine (KVM).
//
// See https://www.ovirt.org/documentation/administration_guide/#chap-Hosts for details.
type Host interface {
	HostData
}

// HostStatus represents the complex states an oVirt host can be in.
type HostStatus string

const (
	// HostStatusConnecting indicates that the engine cannot currently communicate with the host, so it
	// is trying to connect it before fencing it.
	HostStatusConnecting HostStatus = "connecting"
	// HostStatusDown indicates that the host in question is down.
	HostStatusDown HostStatus = "down"
	// HostStatusError indicates that the specified host failed to perform correctly. For example, it
	// failed to run a virtual machine several times without success.
	HostStatusError HostStatus = "error"
	// HostStatusInitializing indicates that the host is shortly before being in HostStatusUp.
	HostStatusInitializing HostStatus = "initializing"
	// HostStatusInstallFailed indicates that setting up the host failed. The administrator needs to look
	// in the log to discover the reason for failure.
	HostStatusInstallFailed HostStatus = "install_failed"
	// HostStatusInstalling indicates that the host is currently being set up.
	HostStatusInstalling HostStatus = "installing"
	// HostStatusInstallingOS indicates that the operating system is being isntalled using Satellite/Foreman.
	HostStatusInstallingOS HostStatus = "installing_os"
	// HostStatusKDumping indicates that the host kernel has crashed and is currently dumping memory.
	HostStatusKDumping HostStatus = "kdumping"
	// HostStatusMaintenance indicates that the host is currently under maintenance and can currently not run
	// virtual machines.
	HostStatusMaintenance HostStatus = "maintenance"
	// HostStatusNonOperational indicates that the host is currently not able to perform due to various reasons,
	// such as a storage not being connected, not supporting a mandatory network, etc.
	HostStatusNonOperational HostStatus = "non_operational"
	// HostStatusNonResponsive indicates that the host is not responding to the engine.
	HostStatusNonResponsive HostStatus = "non_responsive"
	// HostStatusPendingApproval a deprecated status that a vintage ovirt-node/RHV-H hos is pending
	// administrator approval. This status is no longer relevant as Vintage Nodes are no longer supported.
	HostStatusPendingApproval HostStatus = "pending_approval"
	// HostStatusPreparingForMaintenance indicates that the host is currently being drained of all virtual machines
	// via live migration to other hosts. Once the migration is complete the host will move into HostStatusMaintenance.
	HostStatusPreparingForMaintenance HostStatus = "preparing_for_maintenance"
	// HostStatusReboot indicates that the host is currently undergoing a reboot.
	HostStatusReboot HostStatus = "reboot"
	// HostStatusUnassigned indicates that the host is not yet assigned and is in the activation process.
	HostStatusUnassigned HostStatus = "unassigned"
	// HostStatusUp indicates that the host is operating normally.
	HostStatusUp HostStatus = "up"
)

// HostStatusList is a list of HostStatus.
type HostStatusList []HostStatus

// HostStatusValues returns all possible HostStatus values.
func HostStatusValues() HostStatusList {
	return []HostStatus{
		HostStatusConnecting,
		HostStatusDown,
		HostStatusError,
		HostStatusInitializing,
		HostStatusInstallFailed,
		HostStatusInstalling,
		HostStatusInstallingOS,
		HostStatusKDumping,
		HostStatusMaintenance,
		HostStatusNonOperational,
		HostStatusNonResponsive,
		HostStatusPendingApproval,
		HostStatusPreparingForMaintenance,
		HostStatusReboot,
		HostStatusUnassigned,
		HostStatusUp,
	}
}

// Strings creates a string list of the values.
func (l HostStatusList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

func convertSDKHost(sdkHost *ovirtsdk4.Host, client Client) (Host, error) {
	id, ok := sdkHost.Id()
	if !ok {
		return nil, newError(EFieldMissing, "returned host did not contain an ID")
	}
	status, ok := sdkHost.Status()
	if !ok {
		return nil, newError(EFieldMissing, "returned host did not contain a status")
	}
	sdkCluster, ok := sdkHost.Cluster()
	if !ok {
		return nil, newError(EFieldMissing, "returned host did not contain a cluster")
	}
	clusterID, ok := sdkCluster.Id()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch cluster ID from host %s", id)
	}
	return &host{
		client:    client,
		id:        HostID(id),
		status:    HostStatus(status),
		clusterID: ClusterID(clusterID),
	}, nil
}

type host struct {
	client Client

	id        HostID
	clusterID ClusterID
	status    HostStatus
}

func (h host) ID() HostID {
	return h.id
}

func (h host) ClusterID() ClusterID {
	return h.clusterID
}

func (h host) Status() HostStatus {
	return h.status
}
