package ovirtclient

import (
	"errors"
	"fmt"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) AddVMToAffinityGroup(
	clusterID ClusterID,
	vmID VMID,
	agID AffinityGroupID,
	retries ...RetryStrategy,
) error {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	vm, err := ovirtsdk4.NewVmBuilder().Id(string(vmID)).Build()
	if err != nil {
		return wrap(err, EBug, "Failed to build SDK VM object")
	}
	return retry(
		fmt.Sprintf("adding VM %s to affinity group %s", vmID, agID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.
				SystemService().
				ClustersService().
				ClusterService(string(clusterID)).
				AffinityGroupsService().
				GroupService(string(agID)).
				VmsService().
				Add().
				Vm(vm).
				Send()
			// Work around bug 1932320 on older oVirt versions.
			if err != nil && !errors.Is(err, ovirtsdk4.XMLTagNotMatchError{ActualTag: "action", ExpectedTag: "vm"}) {
				return err
			}
			return nil
		},
	)
}

func (m *mockClient) AddVMToAffinityGroup(
	clusterID ClusterID,
	vmID VMID,
	agID AffinityGroupID,
	_ ...RetryStrategy,
) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	clusterAGs, ok := m.affinityGroups[clusterID]
	if !ok {
		return newError(ENotFound, "Cluster %s not found", clusterID)
	}

	ag, ok := clusterAGs[agID]
	if !ok {
		return newError(ENotFound, "Affinity group %s not found", agID)
	}
	for _, agVMID := range ag.vmids {
		if vmID == agVMID {
			return newError(EConflict, "VM %s is already a member of affinity group %s", vmID, agID)
		}
	}

	ag.vmids = append(ag.vmids, vmID)
	return nil
}
