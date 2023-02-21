package ovirtclient

import (
	"fmt"
	"net"
	"time"
)

func (o *oVirtClient) StartVM(id VMID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("starting VM %s", id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().VmsService().VmService(string(id)).Start().Send()
			return err
		})
	return
}

func (m *mockClient) StartVM(id VMID, _ ...RetryStrategy) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	item, ok := m.vms[id]
	if !ok {
		return newError(ENotFound, "vm with ID %s not found", id)
	}

	if item.Status() == VMStatusUp {
		return nil
	}

	hostID, err := m.findSuitableHost(id)
	if err != nil {
		return err
	}
	item.hostID = &hostID
	item.status = VMStatusWaitForLaunch
	go func() {
		time.Sleep(2 * time.Second)
		m.lock.Lock()
		if item.status != VMStatusWaitForLaunch {
			m.lock.Unlock()
			return
		}
		item.status = VMStatusPoweringUp
		m.lock.Unlock()
		time.Sleep(2 * time.Second)
		m.lock.Lock()
		if item.status != VMStatusPoweringUp {
			m.lock.Unlock()
			return
		}
		item.status = VMStatusUp
		m.lock.Unlock()
		time.Sleep(10 * time.Second)
		m.lock.Lock()
		if item.status == VMStatusUp {
			m.vmIPs[item.id] = map[string][]net.IP{
				"lo": {
					net.ParseIP("::1"),
					net.ParseIP("127.0.0.1"),
				},
			}
			i := 0
			for _, nic := range m.nics {
				if nic.vmid == item.id {
					m.vmIPs[item.id][fmt.Sprintf("eth%d", i)] = []net.IP{
						net.ParseIP("192.168.0.123"),
						net.ParseIP("fe80::123"),
					}
					i++
				}
			}
		}
		m.lock.Unlock()
	}()
	return nil
}

func (m *mockClient) findSuitableHost(vmID VMID) (HostID, error) {
	var affectedAffinityGroups []*affinityGroup
	for _, clusterAffinityGroups := range m.affinityGroups {
		for _, affinityGroup := range clusterAffinityGroups {
			if affinityGroup.hasVM(vmID) && (affinityGroup.Enforcing() || affinityGroup.vmsRule.Enforcing()) {
				affectedAffinityGroups = append(affectedAffinityGroups, affinityGroup)
			}
		}
	}
	// Try to find a host that is suitable.
	var foundHost *host
	for _, host := range m.hosts {
		hostSuitable := true
	loop:
		for _, vm := range m.vms {
			if vm.hostID != nil && *vm.hostID == host.id {
				// If the VM resides on the current host
				for _, ag := range affectedAffinityGroups {
					// Check if the VM is a member of the AGs we care about
					if ag.hasVM(vm.id) {
						hostSuitable = false
						break loop
					}
				}
			}
		}
		if hostSuitable {
			foundHost = host
			break
		}
	}
	if foundHost == nil {
		return "", newError(EConflict, "no suitable host found matching affinity group rules")
	}
	hostID := foundHost.ID()
	return hostID, nil
}
