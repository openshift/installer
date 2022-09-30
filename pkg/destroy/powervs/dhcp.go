package powervs

import (
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
	"strings"
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
			o.Logger.Debugf("listDHCPNetworks: DHCP has empty Network.Name: %s", *dhcpServer.ID)
			continue
		}

		if strings.Contains(*dhcpServer.Network.Name, o.InfraID) {
			o.Logger.Debugf("listDHCPNetworks: FOUND: %s (%s)", *dhcpServer.Network.Name, *dhcpServer.ID)
			foundOne = true
			result = append(result, cloudResource{
				key:      *dhcpServer.ID,
				name:     *dhcpServer.Network.Name,
				status:   "",
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
		o.Logger.Infof("Deleted DHCP network %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting DHCP network %q", item.name)

	err = o.dhcpClient.Delete(item.id)
	if err != nil {
		o.Logger.Infof("Error: o.dhcpClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted DHCP network %q", item.name)

	return nil
}

// destroyDHCPNetworks searches for DHCP networks that are in a previous list
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyDHCPNetworks() error {
	found, err := o.listDHCPNetworks()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(dhcpTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyDHCPNetworks: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted DHCP network %q", item.name)
				continue
			}
			err := o.destroyDHCPNetwork(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(dhcpTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(dhcpTypeName); len(items) > 0 {
		return errors.Errorf("destroyDHCPNetworks: %d undeleted items pending", len(items))
	}
	return nil
}
