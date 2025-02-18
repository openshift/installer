package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	dhcpTypeName = "dhcp"
)

// listDHCPNetworks lists previously found DHCP networks in found instances in the vpc.
func (o *ClusterUninstaller) listDHCPNetworks() (cloudResources, error) {
	// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/d_h_c_p_servers.go#L19
	var dhcpServers models.DHCPServers
	// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/d_h_c_p_server.go#L18-L31
	var dhcpServer *models.DHCPServer
	var err error

	o.Logger.Debugf("Listing DHCP networks")

	if o.dhcpClient == nil {
		o.Logger.Infof("Skipping deleting DHCP servers because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	dhcpServers, err = o.dhcpClient.GetAll()
	if err != nil {
		o.Logger.Fatalf("Failed to list DHCP servers: %v", err)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, dhcpServer = range dhcpServers {
		if dhcpServer.Network == nil {
			o.Logger.Debugf("listDHCPNetworks: DHCP has empty Network: %s", *dhcpServer.ID)
			continue
		}
		if dhcpServer.Network.Name == nil {
			// https://github.com/IBM-Cloud/power-go-client/blob/master/power/models/p_vm_instance.go#L22
			var instance *models.PVMInstance

			o.Logger.Debugf("listDHCPNetworks: DHCP has empty Network.Name: %s", *dhcpServer.ID)

			instance, err = o.instanceClient.Get(*dhcpServer.ID)
			o.Logger.Debugf("listDHCPNetworks: Getting instance %s %v", *dhcpServer.ID, err)
			if err != nil {
				continue
			}

			if instance.Status == nil {
				continue
			}
			// If there is a backing DHCP VM and it has a status, then check for an ERROR state
			o.Logger.Debugf("listDHCPNetworks: instance.Status: %s", *instance.Status)
			if *instance.Status != "ERROR" {
				continue
			}

			foundOne = true
			result = append(result, cloudResource{
				key:      *dhcpServer.ID,
				name:     *dhcpServer.ID,
				status:   "VM",
				typeName: dhcpTypeName,
				id:       *dhcpServer.ID,
			})
			continue
		}

		if strings.Contains(*dhcpServer.Network.Name, o.InfraID) {
			o.Logger.Debugf("listDHCPNetworks: FOUND: %s (%s)", *dhcpServer.Network.Name, *dhcpServer.ID)
			foundOne = true
			result = append(result, cloudResource{
				key:      *dhcpServer.ID,
				name:     *dhcpServer.Network.Name,
				status:   "DHCP",
				typeName: dhcpTypeName,
				id:       *dhcpServer.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listDHCPNetworks: NO matching DHCP network found in:")
		for _, dhcpServer = range dhcpServers {
			if dhcpServer.Network == nil {
				continue
			}
			if dhcpServer.Network.Name == nil {
				continue
			}
			o.Logger.Debugf("listDHCPNetworks: only found DHCP: %s", *dhcpServer.Network.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) destroyDHCPNetwork(item cloudResource) error {
	var err error

	_, err = o.dhcpClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted DHCP Network %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting DHCP network %q", item.name)

	err = o.dhcpClient.Delete(item.id)
	if err != nil {
		o.Logger.Infof("Error: o.dhcpClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted DHCP Network %q", item.name)

	return nil
}

func (o *ClusterUninstaller) destroyDHCPVM(item cloudResource) error {
	var err error

	_, err = o.instanceClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted DHCP VM %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting DHCP VM %q", item.name)

	err = o.instanceClient.Delete(item.id)
	if err != nil {
		o.Logger.Infof("Error: DHCP o.instanceClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted DHCP VM %q", item.name)

	return nil
}

// destroyDHCPNetworks searches for DHCP networks that are in a previous list
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyDHCPNetworks() error {
	firstPassList, err := o.listDHCPNetworks()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(dhcpTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyDHCPNetworks: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			var err2 error

			switch item.status {
			case "DHCP":
				err2 = o.destroyDHCPNetwork(item)
			case "VM":
				err2 = o.destroyDHCPVM(item)
			default:
				err2 = fmt.Errorf("unknown DHCP item status %s", item.status)
				return true, err2
			}
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyDHCPNetworks: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(dhcpTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyDHCPNetworks: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyDHCPNetworks: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listDHCPNetworks()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyDHCPNetworks: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyDHCPNetworks: ExponentialBackoffWithContext (list) returns ", err)
	}

	// PowerVS hack:
	// We were asked to query for the subnet still existing as a test for the DHCP network to be
	// finally destroyed.  Even though we can't list it anymore, it is still being destroyed. :(
	backoff = wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listPowerSubnets()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyDHCPNetworks: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyDHCPNetworks: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
