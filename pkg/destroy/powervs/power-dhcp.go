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

// listDHCPNetworks lists PowerVS DHCP instances matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listDHCPNetworks() (cloudResources, error) {
	var (
		dhcpIDs []string
		dhcpID  string
		ctx     context.Context
		cancel  context.CancelFunc
		result  = make([]cloudResource, 0, 1)
		// https://github.com/IBM-Cloud/power-go-client/blob/master/power/models/d_h_c_p_server_detail.go#L21
		dhcpServer *models.DHCPServerDetail
		// https://github.com/IBM-Cloud/power-go-client/blob/master/power/models/p_vm_instance.go#L22
		instance *models.PVMInstance
		err      error
	)

	if false { // @TODO o.searchByTag {
		// Should we list by tag matching?
		// @TODO dhcpIDs, err = o.listByTag(TagTypeDHCP)
		err = fmt.Errorf("listByTag(TagTypeDHCP) is not supported yet")
	} else {
		// Otherwise list will list by name matching.
		dhcpIDs, err = o.listDHCPNetworksByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, dhcpID = range dhcpIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listDHCPNetworks: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		o.Logger.Debugf("listDHCPNetworks: Getting DHCP %s %+v", dhcpID, dhcpServer)
		dhcpServer, err = o.dhcpClient.Get(dhcpID)
		if err != nil {
			if strings.Contains(err.Error(), "could not retrieve dhcp server") {
				continue
			}
			return nil, fmt.Errorf("listDHCPNetworks could not get DHCP %s: %w", dhcpID, err)
		}

		if dhcpServer.Network == nil {
			o.Logger.Debugf("listDHCPNetworks: DHCP has empty Network: %s", *dhcpServer.ID)
			continue
		}
		if dhcpServer.Network.Name == nil {
			o.Logger.Debugf("listDHCPNetworks: DHCP has empty Network.Name: %s", *dhcpServer.ID)

			instance, err = o.instanceClient.Get(*dhcpServer.ID)
			o.Logger.Debugf("listDHCPNetworks: Getting DHCP VM %s %+v", *dhcpServer.ID, instance)
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

			result = append(result, cloudResource{
				key:      *dhcpServer.ID,
				name:     *dhcpServer.ID,
				status:   "VM",
				typeName: dhcpTypeName,
				id:       *dhcpServer.ID,
			})
			continue
		}

		result = append(result, cloudResource{
			key:      *dhcpServer.ID,
			name:     *dhcpServer.Network.Name,
			status:   "DHCP",
			typeName: dhcpTypeName,
			id:       *dhcpServer.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listDHCPNetworksByName lists previously found DHCP networks in found instances in the vpc.
func (o *ClusterUninstaller) listDHCPNetworksByName() ([]string, error) {
	var (
		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/d_h_c_p_servers.go#L19
		dhcpServers models.DHCPServers
		// https://github.com/IBM-Cloud/power-go-client/blob/v1.8.3/power/models/d_h_c_p_server.go#L20-L33
		dhcpServer *models.DHCPServer
		ctx        context.Context
		cancel     context.CancelFunc
		result     = make([]string, 0, 1)
		foundOne   = false
		err        error
	)

	o.Logger.Debugf("Listing DHCP networks by NAME")

	if o.dhcpClient == nil {
		o.Logger.Infof("Skipping deleting DHCP servers because no service instance was found")
		return result, nil
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listDHCPNetworksByName: case <-ctx.Done()")
		return result, ctx.Err() // we're cancelled, abort
	default:
	}

	dhcpServers, err = o.dhcpClient.GetAll()
	if err != nil {
		o.Logger.Fatalf("Failed to list DHCP servers: %v", err)
	}

	for _, dhcpServer = range dhcpServers {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listDHCPNetworksByName: case <-ctx.Done()")
			return result, ctx.Err() // we're cancelled, abort
		default:
		}

		if dhcpServer.Network == nil {
			o.Logger.Debugf("listDHCPNetworksByName: DHCP has empty Network: %s", *dhcpServer.ID)
			continue
		}

		if dhcpServer.Network.Name == nil {
			result = append(result, *dhcpServer.ID)
			continue
		}

		if strings.Contains(*dhcpServer.Network.Name, o.InfraID) {
			o.Logger.Debugf("listDHCPNetworksByName: FOUND: %s (%s)", *dhcpServer.Network.Name, *dhcpServer.ID)
			foundOne = true
			result = append(result, *dhcpServer.ID)
		}
	}
	if !foundOne {
		o.Logger.Debugf("listDHCPNetworksByName: NO matching DHCP network found in:")
		for _, dhcpServer = range dhcpServers {
			select {
			case <-ctx.Done():
				o.Logger.Debugf("listDHCPNetworksByName: case <-ctx.Done()")
				return result, ctx.Err() // we're cancelled, abort
			default:
			}

			if dhcpServer.Network == nil {
				continue
			}
			if dhcpServer.Network.Name == nil {
				continue
			}
			o.Logger.Debugf("listDHCPNetworksByName: only found DHCP: %s", *dhcpServer.Network.Name)
		}
	}

	return result, nil
}

// extractNetworkIDFromError extracts network ID from error message if present.
// Error format: "network xxx-xxx-xxxxx still attached to pvm-instances".
func extractNetworkIDFromError(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()
	// Look for pattern "network <uuid> still attached"
	parts := strings.Split(errStr, "network ")
	if len(parts) > 1 {
		remaining := parts[1]
		// Extract UUID from error message (format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
		spaceIdx := strings.Index(remaining, " ")
		if spaceIdx > 0 {
			networkID := remaining[:spaceIdx]
			// UUID format validation (36 chars with dashes)
			if len(networkID) == 36 && strings.Count(networkID, "-") == 4 {
				return networkID
			}
		}
	}
	return ""
}

// findNetworkIDByName finds a network ID by matching the network name.
func (o *ClusterUninstaller) findNetworkIDByName(networkName string) string {
	if o.networkClient == nil {
		return ""
	}
	networks, err := o.networkClient.GetAll()
	if err != nil {
		o.Logger.Debugf("Failed to list networks to find ID for %q: %v", networkName, err)
		return ""
	}
	for _, network := range networks.Networks {
		if network.Name != nil && *network.Name == networkName {
			if network.NetworkID != nil {
				return *network.NetworkID
			}
		}
	}
	return ""
}

// isDHCPNetworkAttachedError checks if an error indicates the DHCP network is attached to PVM instances.
func isDHCPNetworkAttachedError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "still attached to pvm-instances") ||
		strings.Contains(errStr, "still attached to") ||
		strings.Contains(errStr, "pcloudDhcpDeleteBadRequest")
}

// destroyDHCPNetwork deletes a PowerVS DHCP network.
func (o *ClusterUninstaller) destroyDHCPNetwork(item cloudResource) error {
	var err error

	dhcpServer, err := o.dhcpClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted DHCP Network %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting DHCP network %q", item.name)

	// Before deleting the DHCP server, check if its network has attached network interfaces
	// that need to be deleted first. This prevents errors like "network still attached to pvm-instances, that fail early before subnet network interfaces are deleted in the destroyPowerSubnets() stage."
	var networkID string
	if dhcpServer.Network != nil {
		// Try to find network ID by name
		if dhcpServer.Network.Name != nil {
			networkID = o.findNetworkIDByName(*dhcpServer.Network.Name)
			if networkID != "" {
				o.Logger.Debugf("Found network ID %s for DHCP subnet %q. Checking for network interfaces...", networkID, item.name)
				// Try to delete network interfaces from the subnet
				if nicErr := o.deleteNetworkInterfaces(networkID); nicErr != nil {
					o.Logger.Debugf("Note: Could not delete network interfaces for DHCP subnet %q: %v (will attempt DHCP deletion anyway)", item.name, nicErr)
				}
			}
		}
	}

	err = o.dhcpClient.Delete(item.id)
	if err != nil {
		// If deletion failed because network is still attached to instances, try deleting network interfaces
		if isDHCPNetworkAttachedError(err) {
			// Try to extract network ID from error message if we don't have it yet
			if networkID == "" {
				networkID = extractNetworkIDFromError(err)
			}
			// If still no network ID, try finding by name again
			if networkID == "" && dhcpServer.Network != nil && dhcpServer.Network.Name != nil {
				networkID = o.findNetworkIDByName(*dhcpServer.Network.Name)
			}

			if networkID != "" {
				o.Logger.Debugf("DHCP subnet %q is still attached to instances. Attempting to delete network interfaces from network %s...", item.name, networkID)
				if nicErr := o.deleteNetworkInterfaces(networkID); nicErr != nil {
					o.Logger.Warnf("Failed to delete network interfaces for DHCP subnet %q: %v", item.name, nicErr)
				}
				// Return error to trigger retry after NIC deletion
				return fmt.Errorf("DHCP server deletion blocked by attached network interfaces: %w", err)
			}
			o.Logger.Warnf("Could not determine network ID for DHCP subnet %q to delete network interfaces", item.name)
		}
		o.Logger.Infof("Error: o.dhcpClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted DHCP Network %q", item.name)

	return nil
}

// destroyDHCPVM deletes a PowerVS backing VM for a DHCP network.
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

	// Note: DHCP server subnets will be deleted in the destroyPowerSubnets() stage.
	// We no longer wait for them here since:
	// 1. Network interfaces are now properly deleted during DHCP deletion
	// 2. Subnet deletion happens in a later stage with its own retry logic
	// 3. Waiting here was causing timeouts since subnets are deleted in a different stage

	return nil
}
