package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

const powerSubnetTypeName = "powerSubnet"

// listPowerSubnets lists subnets in the Power Server.
func (o *ClusterUninstaller) listPowerSubnets() (cloudResources, error) {
	o.Logger.Debugf("Listing Power Server Subnets")

	if o.instanceClient == nil {
		o.Logger.Infof("Skipping deleting Power service subnets because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	networks, err := o.networkClient.GetAll()
	if err != nil {
		o.Logger.Warnf("Error networkClient.GetAll: %v", err)
		return nil, err
	}

	result := []cloudResource{}
	for _, network := range networks.Networks {
		if strings.Contains(*network.Name, o.InfraID) {
			o.Logger.Debugf("listPowerSubnets: FOUND: %s, %s", *network.NetworkID, *network.Name)
			result = append(result, cloudResource{
				key:      *network.NetworkID,
				name:     *network.Name,
				status:   "",
				typeName: powerSubnetTypeName,
				id:       *network.NetworkID,
			})
		}
	}
	if len(result) == 0 {
		o.Logger.Debugf("listPowerSubnets: NO matching subnet against: %s", o.InfraID)
		for _, network := range networks.Networks {
			o.Logger.Debugf("listPowerSubnets: network: %s", *network.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

// deleteNetworkInterfaces deletes all network interfaces attached to a subnet.
func (o *ClusterUninstaller) deleteNetworkInterfaces(subnetID string) error {
	interfaces, err := o.networkClient.GetAllNetworkInterfaces(subnetID)
	if err != nil {
		return fmt.Errorf("failed to list network interfaces: %w", err)
	}

	for _, nic := range interfaces.Interfaces {
		if nic.ID != nil {
			o.Logger.Debugf("Deleting network interface %q from subnet %q", *nic.ID, subnetID)
			if err := o.networkClient.DeleteNetworkInterface(subnetID, *nic.ID); err != nil {
				o.Logger.Warnf("Failed to delete network interface %q: %v", *nic.ID, err)
				// Continue trying to delete other interfaces
			}
		}
	}

	return nil
}

// isNetworkInterfaceError checks if an error indicates network interfaces are blocking deletion. (i.e 409 Conflict)
func isNetworkInterfaceError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "one or more network interfaces have an IP allocation") ||
		strings.Contains(errStr, "status 409") ||
		strings.Contains(errStr, "409")
}

func (o *ClusterUninstaller) deletePowerSubnet(item cloudResource) error {
	if _, err := o.networkClient.Get(item.id); err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Power Network %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting Power Network %q", item.name)

	if err := o.networkClient.Delete(item.id); err != nil {
		// If deletion failed due to attached network interfaces, delete them and retry
		if isNetworkInterfaceError(err) {
			o.Logger.Debugf("Subnet %q has attached network interfaces. Deleting them...", item.name)
			if nicErr := o.deleteNetworkInterfaces(item.id); nicErr != nil {
				o.Logger.Warnf("Failed to delete network interfaces for subnet %q: %v", item.name, nicErr)
			}
			// Return error to trigger retry after NIC deletion
			return fmt.Errorf("subnet deletion blocked by network interfaces: %w", err)
		}

		o.Logger.Infof("Error: o.networkClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted Power Network %q", item.name)

	return nil
}

// destroyPowerSubnets removes all subnet resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyPowerSubnets() error {
	firstPassList, err := o.listPowerSubnets()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(powerSubnetTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyPowerSubnets: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deletePowerSubnet(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyPowerSubnets: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(powerSubnetTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyPowerSubnets: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyPowerSubnets: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
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
			o.Logger.Debugf("destroyPowerSubnets: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyPowerSubnets: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
