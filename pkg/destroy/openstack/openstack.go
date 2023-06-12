package openstack

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/apiversions"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	sharesnapshots "github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/snapshots"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"github.com/openshift/installer/pkg/types/openstack/validation/networkextensions"
)

const (
	cinderCSIClusterIDKey           = "cinder.csi.openstack.org/cluster"
	manilaCSIClusterIDKey           = "manila.csi.openstack.org/cluster"
	minOctaviaVersionWithTagSupport = "v2.5"
)

// Filter holds the key/value pairs for the tags we will be matching
// against.
type Filter map[string]string

// ObjectWithTags is a generic way to represent an OpenStack object
// and its tags so that filtering objects client-side can be done in a generic
// way.
//
// Note we use UUID not Name as not all OpenStack services require a unique
// name.
type ObjectWithTags struct {
	ID   string
	Tags map[string]string
}

// deleteFunc type is the interface a function needs to implement to be called as a goroutine.
// The (bool, error) return type mimics wait.ExponentialBackoff where the bool indicates successful
// completion, and the error is for unrecoverable errors.
type deleteFunc func(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	// Cloud is the cloud name as set in clouds.yml
	Cloud string
	// Filter contains the openshiftClusterID to filter tags
	Filter Filter
	// InfraID contains unique cluster identifier
	InfraID string
	Logger  logrus.FieldLogger
}

// New returns an OpenStack destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Cloud:   metadata.ClusterPlatformMetadata.OpenStack.Cloud,
		Filter:  metadata.ClusterPlatformMetadata.OpenStack.Identifier,
		InfraID: metadata.InfraID,
		Logger:  logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	opts := openstackdefaults.DefaultClientOpts(o.Cloud)

	// Check that the cloud has the minimum requirements for the destroy
	// script to work properly.
	if err := validateCloud(opts, o.Logger); err != nil {
		return nil, err
	}

	// deleteFuncs contains the functions that will be launched as
	// goroutines.
	deleteFuncs := map[string]deleteFunc{
		"cleanVIPsPorts":        cleanVIPsPorts,
		"deleteServers":         deleteServers,
		"deleteServerGroups":    deleteServerGroups,
		"deleteTrunks":          deleteTrunks,
		"deleteLoadBalancers":   deleteLoadBalancers,
		"deletePorts":           deletePortsByFilter,
		"deleteSecurityGroups":  deleteSecurityGroups,
		"clearRouterInterfaces": clearRouterInterfaces,
		"deleteSubnets":         deleteSubnets,
		"deleteSubnetPools":     deleteSubnetPools,
		"deleteNetworks":        deleteNetworks,
		"deleteContainers":      deleteContainers,
		"deleteVolumes":         deleteVolumes,
		"deleteShares":          deleteShares,
		"deleteVolumeSnapshots": deleteVolumeSnapshots,
		"deleteFloatingIPs":     deleteFloatingIPs,
		"deleteImages":          deleteImages,
	}
	returnChannel := make(chan string)

	// launch goroutines
	for name, function := range deleteFuncs {
		go deleteRunner(name, function, opts, o.Filter, o.Logger, returnChannel)
	}

	// wait for them to finish
	for i := 0; i < len(deleteFuncs); i++ {
		res := <-returnChannel
		o.Logger.Debugf("goroutine %v complete", res)
	}

	// we want to remove routers as the last thing as it requires detaching the
	// FIPs and that will cause it impossible to track which FIPs are tied to
	// LBs being deleted.
	err := deleteRouterRunner(opts, o.Filter, o.Logger)
	if err != nil {
		return nil, err
	}

	// we need to untag the custom network if it was provided by the user
	err = untagRunner(opts, o.InfraID, o.Logger)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func deleteRunner(deleteFuncName string, dFunction deleteFunc, opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger, channel chan string) {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 15,
		Factor:   1.3,
		Steps:    25,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return dFunction(opts, filter, logger)
	})

	if err != nil {
		logger.Fatalf("Unrecoverable error/timed out: %v", err)
	}

	// record that the goroutine has run to completion
	channel <- deleteFuncName
}

// filterObjects will do client-side filtering given an appropriately filled out
// list of ObjectWithTags.
func filterObjects(osObjects []ObjectWithTags, filters Filter) []ObjectWithTags {
	objectsWithTags := []ObjectWithTags{}
	filteredObjects := []ObjectWithTags{}

	// first find the objects that have all the desired tags
	for _, object := range osObjects {
		allTagsFound := true
		for key := range filters {
			if _, ok := object.Tags[key]; !ok {
				// doesn't have one of the tags we're looking for so skip it
				allTagsFound = false
				break
			}
		}
		if allTagsFound {
			objectsWithTags = append(objectsWithTags, object)
		}
	}

	// now check that the values match
	for _, object := range objectsWithTags {
		valuesMatch := true
		for key, val := range filters {
			if object.Tags[key] != val {
				valuesMatch = false
				break
			}
		}
		if valuesMatch {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}

func filterTags(filters Filter) []string {
	tags := []string{}
	for k, v := range filters {
		tags = append(tags, strings.Join([]string{k, v}, "="))
	}
	return tags
}

func deleteServers(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack servers")
	defer logger.Debugf("Exiting deleting openstack servers")

	conn, err := clientconfig.NewServiceClient("compute", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allPages, err := servers.List(conn, servers.ListOpts{}).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	serverObjects := []ObjectWithTags{}
	for _, server := range allServers {
		serverObjects = append(
			serverObjects, ObjectWithTags{
				ID:   server.ID,
				Tags: server.Metadata})
	}

	filteredServers := filterObjects(serverObjects, filter)
	numberToDelete := len(filteredServers)
	numberDeleted := 0
	for _, server := range filteredServers {
		logger.Debugf("Deleting Server %q", server.ID)
		err = servers.Delete(conn, server.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the server cannot be found and return with an appropriate message if it's another type of error
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next server
				logger.Errorf("Deleting server %q failed: %v", server.ID, err)
				continue
			}
			logger.Debugf("Cannot find server %q. It's probably already been deleted.", server.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deleteServerGroups(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack server groups")
	defer logger.Debugf("Exiting deleting openstack server groups")

	// We need to delete all server groups that have names with the cluster
	// ID as a prefix
	var clusterID string
	for k, v := range filter {
		if strings.ToLower(k) == "openshiftclusterid" {
			clusterID = v
			break
		}
	}

	conn, err := clientconfig.NewServiceClient("compute", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allPages, err := servergroups.List(conn, nil).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	filteredGroups := make([]servergroups.ServerGroup, 0, len(allServerGroups))
	for _, serverGroup := range allServerGroups {
		if strings.HasPrefix(serverGroup.Name, clusterID) {
			filteredGroups = append(filteredGroups, serverGroup)
		}
	}

	numberToDelete := len(filteredGroups)
	numberDeleted := 0
	for _, serverGroup := range filteredGroups {
		logger.Debugf("Deleting Server Group %q", serverGroup.ID)
		if err = servergroups.Delete(conn, serverGroup.ID).ExtractErr(); err != nil {
			// Ignore the error if the server cannot be found and
			// return with an appropriate message if it's another
			// type of error
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next server group
				logger.Errorf("Deleting server group %q failed: %v", serverGroup.ID, err)
				continue
			}
			logger.Debugf("Cannot find server group %q. It's probably already been deleted.", serverGroup.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deletePortsByNetwork(opts *clientconfig.ClientOpts, networkID string, logger logrus.FieldLogger) (bool, error) {

	listOpts := ports.ListOpts{
		NetworkID: networkID,
	}

	result, err := deletePorts(opts, listOpts, logger)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	return result, err
}

func deletePortsByFilter(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {

	tags := filterTags(filter)
	listOpts := ports.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	result, err := deletePorts(opts, listOpts, logger)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	return result, err
}

func getFIPsByPort(conn *gophercloud.ServiceClient, logger logrus.FieldLogger) (map[string]floatingips.FloatingIP, error) {
	// Prefetch list of FIPs to save list calls for each port
	fipByPort := make(map[string]floatingips.FloatingIP)
	allPages, err := floatingips.List(conn, floatingips.ListOpts{}).AllPages()
	if err != nil {
		logger.Error(err)
		return fipByPort, nil
	}
	allFIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		logger.Error(err)
		return fipByPort, nil
	}

	// Organize FIPs for easy lookup
	for _, fip := range allFIPs {
		fipByPort[fip.PortID] = fip
	}
	return fipByPort, err
}

func deletePorts(opts *clientconfig.ClientOpts, listOpts ports.ListOpts, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack ports")
	defer logger.Debugf("Exiting deleting openstack ports")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allPages, err := ports.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	numberToDelete := len(allPorts)
	numberDeleted := 0

	fipByPort, err := getFIPsByPort(conn, logger)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	deletePortsWorker := func(portsChannel <-chan ports.Port, deletedChannel chan<- int) {
		localDeleted := 0
		for port := range portsChannel {
			// If a user provisioned floating ip was used, it needs to be dissociated.
			// Any floating Ip's associated with ports that are going to be deleted will be dissociated.
			if fip, ok := fipByPort[port.ID]; ok {
				logger.Debugf("Dissociating Floating IP %q", fip.ID)
				_, err := floatingips.Update(conn, fip.ID, floatingips.UpdateOpts{}).Extract()
				if err != nil {
					// Ignore the error if the floating ip cannot be found and return with an appropriate message if it's another type of error
					var gerr gophercloud.ErrDefault404
					if !errors.As(err, &gerr) {
						// Just log the error and move on to the next port
						logger.Errorf("While deleting port %q, the update of the floating IP %q failed with error: %v", port.ID, fip.ID, err)
						continue
					}
					logger.Debugf("Cannot find floating ip %q. It's probably already been deleted.", fip.ID)
				}
			}

			logger.Debugf("Deleting Port %q", port.ID)
			err = ports.Delete(conn, port.ID).ExtractErr()
			if err != nil {
				// This can fail when port is still in use so return/retry
				// Just log the error and move on to the next port
				logger.Debugf("Deleting Port %q failed with error: %v", port.ID, err)
				// Try to delete associated trunk
				deleteAssociatedTrunk(conn, logger, port.ID)
				continue
			}
			localDeleted++
		}
		deletedChannel <- localDeleted
	}

	const workersNumber = 10
	portsChannel := make(chan ports.Port, workersNumber)
	deletedChannel := make(chan int, workersNumber)

	// start worker goroutines
	for i := 0; i < workersNumber; i++ {
		go deletePortsWorker(portsChannel, deletedChannel)
	}

	// feed worker goroutines with ports
	for _, port := range allPorts {
		portsChannel <- port
	}
	close(portsChannel)

	// wait for them to finish and accumulate number of ports deleted by each
	for i := 0; i < workersNumber; i++ {
		numberDeleted += <-deletedChannel
	}

	return numberDeleted == numberToDelete, nil
}

func getSecurityGroups(conn *gophercloud.ServiceClient, filter Filter) ([]sg.SecGroup, error) {
	var emptySecurityGroups []sg.SecGroup
	tags := filterTags(filter)
	listOpts := sg.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := sg.List(conn, listOpts).AllPages()
	if err != nil {
		return emptySecurityGroups, err
	}

	allGroups, err := sg.ExtractGroups(allPages)
	if err != nil {
		return emptySecurityGroups, err
	}
	return allGroups, nil
}

func deleteSecurityGroups(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack security-groups")
	defer logger.Debugf("Exiting deleting openstack security-groups")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allGroups, err := getSecurityGroups(conn, filter)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	numberToDelete := len(allGroups)
	numberDeleted := 0
	for _, group := range allGroups {
		logger.Debugf("Deleting Security Group: %q", group.ID)
		err = sg.Delete(conn, group.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the security group cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when sg is still in use by servers
				// Just log the error and move on to the next security group
				logger.Debugf("Deleting Security Group %q failed with error: %v", group.ID, err)
				continue
			}
			logger.Debugf("Cannot find security group %q. It's probably already been deleted.", group.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func updateFips(allFIPs []floatingips.FloatingIP, opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) error {
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return err
	}

	for _, fip := range allFIPs {
		logger.Debugf("Updating FIP %s", fip.ID)
		_, err := floatingips.Update(conn, fip.ID, floatingips.UpdateOpts{}).Extract()
		if err != nil {
			// Ignore the error if the resource cannot be found and return with an appropriate message if it's another type of error
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				logger.Errorf("Updating floating IP %q for Router failed: %v", fip.ID, err)
				return err
			}
			logger.Debugf("Cannot find floating ip %q. It's probably already been deleted.", fip.ID)
		}
	}
	return nil
}

// deletePortFIPs looks up FIPs associated to the port and attempts to delete them
func deletePortFIPs(portID string, opts *clientconfig.ClientOpts, logger logrus.FieldLogger) error {
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return err
	}

	fipPages, err := floatingips.List(conn, floatingips.ListOpts{PortID: portID}).AllPages()

	if err != nil {
		logger.Error(err)
		return err
	}

	fips, err := floatingips.ExtractFloatingIPs(fipPages)
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, fip := range fips {
		logger.Debugf("Deleting FIP %q", fip.ID)
		err = floatingips.Delete(conn, fip.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the FIP cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				logger.Errorf("Deleting FIP %q failed: %v", fip.ID, err)
				return err
			}
			logger.Debugf("Cannot find FIP %q. It's probably already been deleted.", fip.ID)
		}
	}
	return nil
}

func getRouters(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) ([]routers.Router, error) {
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return nil, err
	}
	tags := filterTags(filter)
	listOpts := routers.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := routers.List(conn, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allRouters, err := routers.ExtractRouters(allPages)
	if err != nil {
		return nil, err
	}
	return allRouters, nil
}

func deleteRouters(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack routers")
	defer logger.Debugf("Exiting deleting openstack routers")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allRouters, err := getRouters(opts, filter, logger)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allRouters)
	numberDeleted := 0
	for _, router := range allRouters {
		fipOpts := floatingips.ListOpts{
			RouterID: router.ID,
		}

		fipPages, err := floatingips.List(conn, fipOpts).AllPages()
		if err != nil {
			logger.Error(err)
			return false, nil
		}

		allFIPs, err := floatingips.ExtractFloatingIPs(fipPages)
		if err != nil {
			logger.Error(err)
			return false, nil
		}
		// If a user provisioned floating ip was used, it needs to be dissociated
		// Any floating Ip's associated with routers that are going to be deleted will be dissociated
		err = updateFips(allFIPs, opts, filter, logger)
		if err != nil {
			logger.Error(err)
			continue
		}
		// Clean Gateway interface
		updateOpts := routers.UpdateOpts{
			GatewayInfo: &routers.GatewayInfo{},
		}

		_, err = routers.Update(conn, router.ID, updateOpts).Extract()
		if err != nil {
			logger.Error(err)
		}

		logger.Debugf("Deleting Router %q", router.ID)
		err = routers.Delete(conn, router.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the router cannot be found and return with an appropriate message if it's another type of error
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next router
				logger.Errorf("Deleting router %q failed: %v", router.ID, err)
				continue
			}
			logger.Debugf("Cannot find router %q. It's probably already been deleted.", router.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func getRouterInterfaces(conn *gophercloud.ServiceClient, allNetworks []networks.Network, logger logrus.FieldLogger) ([]ports.Port, error) {
	var routerPorts []ports.Port
	for _, network := range allNetworks {
		if len(network.Subnets) == 0 {
			continue
		}
		subnet, err := subnets.Get(conn, network.Subnets[0]).Extract()
		if err != nil {
			logger.Debug(err)
			return routerPorts, nil
		}
		if subnet.GatewayIP == "" {
			continue
		}
		portListOpts := ports.ListOpts{
			FixedIPs: []ports.FixedIPOpts{
				{
					SubnetID: network.Subnets[0],
				},
				{
					IPAddress: subnet.GatewayIP,
				},
			},
		}

		allPagesPort, err := ports.List(conn, portListOpts).AllPages()
		if err != nil {
			logger.Error(err)
			return routerPorts, nil
		}

		routerPorts, err = ports.ExtractPorts(allPagesPort)
		if err != nil {
			logger.Error(err)
			return routerPorts, nil
		}

		if len(routerPorts) != 0 {
			logger.Debugf("Found Port %q connected to Router", routerPorts[0].ID)
			return routerPorts, nil
		}
	}
	return routerPorts, nil
}

func clearRouterInterfaces(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debugf("Removing interfaces from router")
	defer logger.Debug("Exiting removal of interfaces from router")
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	tags := filterTags(filter)
	networkListOpts := networks.ListOpts{
		Tags: strings.Join(tags, ","),
	}

	allNetworksPages, err := networks.List(conn, networkListOpts).AllPages()
	if err != nil {
		logger.Debug(err)
		return false, nil
	}

	allNetworks, err := networks.ExtractNetworks(allNetworksPages)
	if err != nil {
		logger.Debug(err)
		return false, nil
	}

	// Identify router by checking any tagged Network that has a Subnet
	// with GatewayIP set
	routerPorts, err := getRouterInterfaces(conn, allNetworks, logger)
	if err != nil {
		logger.Debug(err)
		return false, nil
	}

	if len(routerPorts) == 0 {
		return true, nil
	}

	routerID := routerPorts[0].DeviceID
	router, err := routers.Get(conn, routerID).Extract()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	removed, err := removeRouterInterfaces(conn, filter, *router, logger)
	if err != nil {
		logger.Debug(err)
		return false, nil
	}
	return removed, nil
}

func removeRouterInterfaces(client *gophercloud.ServiceClient, filter Filter, router routers.Router, logger logrus.FieldLogger) (bool, error) {
	// Get router interface ports
	portListOpts := ports.ListOpts{
		DeviceID: router.ID,
	}
	allPagesPort, err := ports.List(client, portListOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, fmt.Errorf("failed to get ports list: %w", err)
	}
	allPorts, err := ports.ExtractPorts(allPagesPort)
	if err != nil {
		logger.Error(err)
		return false, fmt.Errorf("failed to extract ports list: %w", err)
	}
	tags := filterTags(filter)
	SubnetlistOpts := subnets.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allSubnetsPage, err := subnets.List(client, SubnetlistOpts).AllPages()
	if err != nil {
		logger.Debug(err)
		return false, fmt.Errorf("failed to list subnets list: %w", err)
	}

	allSubnets, err := subnets.ExtractSubnets(allSubnetsPage)
	if err != nil {
		logger.Debug(err)
		return false, fmt.Errorf("failed to extract subnets list: %w", err)
	}

	clusterTag := "openshiftClusterID=" + filter["openshiftClusterID"]
	clusterRouter := isClusterRouter(clusterTag, router.Tags)

	numberToDelete := len(allPorts)
	numberDeleted := 0
	var customInterfaces []ports.Port
	// map to keep track of whether interface for subnet was already removed
	removedSubnets := make(map[string]bool)
	for _, port := range allPorts {
		for _, IP := range port.FixedIPs {
			// Skip removal if Router was not created by CNO or installer and
			// interface is not handled by the Cluster
			if !clusterRouter && !isClusterSubnet(allSubnets, IP.SubnetID) {
				logger.Debugf("Found custom interface %q on Router %q", port.ID, router.ID)
				customInterfaces = append(customInterfaces, port)
				continue
			}
			if !removedSubnets[IP.SubnetID] {
				removeOpts := routers.RemoveInterfaceOpts{
					SubnetID: IP.SubnetID,
				}
				logger.Debugf("Removing Subnet %q from Router %q", IP.SubnetID, router.ID)
				_, err := routers.RemoveInterface(client, router.ID, removeOpts).Extract()
				if err != nil {
					var gerr gophercloud.ErrDefault404
					if !errors.As(err, &gerr) {
						// This can fail when subnet is still in use
						logger.Debugf("Removing Subnet %q from Router %q failed: %v", IP.SubnetID, router.ID, err)
						return false, nil
					}
					logger.Debugf("Cannot find subnet %q. It's probably already been removed from router %q.", IP.SubnetID, router.ID)
				}
				removedSubnets[IP.SubnetID] = true
				numberDeleted++
			}
		}
	}
	numberToDelete -= len(customInterfaces)
	return numberToDelete == numberDeleted, nil
}

func isClusterRouter(clusterTag string, tags []string) bool {
	for _, tag := range tags {
		if clusterTag == tag {
			return true
		}
	}
	return false
}

func getRouterByPort(client *gophercloud.ServiceClient, allPorts []ports.Port) (routers.Router, error) {
	empty := routers.Router{}
	for _, port := range allPorts {
		if port.DeviceID != "" {
			page, err := routers.List(client, routers.ListOpts{ID: port.DeviceID}).AllPages()
			if err != nil {
				return empty, fmt.Errorf("failed to get router list: %w", err)
			}

			routerList, err := routers.ExtractRouters(page)
			if err != nil {
				return empty, fmt.Errorf("failed to extract routers list: %w", err)
			}

			if len(routerList) == 1 {
				return routerList[0], nil
			}
		}
	}
	return empty, nil
}

func deleteLeftoverLoadBalancers(opts *clientconfig.ClientOpts, logger logrus.FieldLogger, networkID string) error {
	conn, err := clientconfig.NewServiceClient("load-balancer", opts)
	if err != nil {
		// Ignore the error if Octavia is not available for the cloud
		var gerr *gophercloud.ErrEndpointNotFound
		if errors.As(err, &gerr) {
			logger.Debug("Skip load balancer deletion because Octavia endpoint is not found")
			return nil
		}
		logger.Error(err)
		return err
	}

	listOpts := loadbalancers.ListOpts{
		VipNetworkID: networkID,
	}
	allPages, err := loadbalancers.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return err
	}

	allLoadBalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		logger.Error(err)
		return err
	}
	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: true,
	}
	deleted := 0
	for _, loadbalancer := range allLoadBalancers {
		if !strings.HasPrefix(loadbalancer.Description, "Kubernetes external service") {
			logger.Debugf("Not deleting LoadBalancer %q with description %q", loadbalancer.ID, loadbalancer.Description)
			continue
		}
		logger.Debugf("Deleting LoadBalancer %q", loadbalancer.ID)

		// Cascade delete of an LB won't remove the associated FIP, we have to do it ourselves.
		err := deletePortFIPs(loadbalancer.VipPortID, opts, logger)
		if err != nil {
			// Go to the next LB, but do not delete current one or we'll lose reference to the FIP that failed deletion.
			continue
		}

		err = loadbalancers.Delete(conn, loadbalancer.ID, deleteOpts).ExtractErr()
		if err != nil {
			// Ignore the error if the load balancer cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when the load balancer is still in use so return/retry
				// Just log the error and move on to the next LB
				logger.Debugf("Deleting load balancer %q failed: %v", loadbalancer.ID, err)
				continue
			}
			logger.Debugf("Cannot find load balancer %q. It's probably already been deleted.", loadbalancer.ID)
		}
		deleted++
	}

	if deleted != len(allLoadBalancers) {
		return fmt.Errorf("only deleted %d of %d load balancers", deleted, len(allLoadBalancers))
	}
	return nil
}

func isClusterSubnet(subnets []subnets.Subnet, subnetID string) bool {
	for _, subnet := range subnets {
		if subnet.ID == subnetID {
			return true
		}
	}
	return false
}

func deleteSubnets(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack subnets")
	defer logger.Debugf("Exiting deleting openstack subnets")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	tags := filterTags(filter)
	listOpts := subnets.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := subnets.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allSubnets, err := subnets.ExtractSubnets(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allSubnets)
	numberDeleted := 0
	for _, subnet := range allSubnets {
		logger.Debugf("Deleting Subnet: %q", subnet.ID)
		err = subnets.Delete(conn, subnet.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the subnet cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when subnet is still in use
				// Just log the error and move on to the next subnet
				logger.Debugf("Deleting Subnet %q failed: %v", subnet.ID, err)
				continue
			}
			logger.Debugf("Cannot find subnet %q. It's probably already been deleted.", subnet.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deleteNetworks(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack networks")
	defer logger.Debugf("Exiting deleting openstack networks")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	tags := filterTags(filter)
	listOpts := networks.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := networks.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	numberToDelete := len(allNetworks)
	numberDeleted := 0
	for _, network := range allNetworks {
		logger.Debugf("Deleting network: %q", network.ID)
		err = networks.Delete(conn, network.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the network cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when network is still in use. Let's log an error and try to fix this.
				logger.Debugf("Deleting Network %q failed: %v", network.ID, err)

				// First try to delete eventual leftover load balancers
				// *This has to be done before attempt to remove ports or we'll delete LB ports!*
				err := deleteLeftoverLoadBalancers(opts, logger, network.ID)
				if err != nil {
					logger.Error(err)
					// Do not attempt to delete ports on LB removal problem or we'll lose FIP associations!
					continue
				}

				// Only then try to remove all the ports it may still contain (untagged as well).
				// *We cannot delete ports before LBs because we'll lose FIP associations!*
				_, err = deletePortsByNetwork(opts, network.ID, logger)
				if err != nil {
					logger.Error(err)
				}
				continue
			}
			logger.Debugf("Cannot find network %q. It's probably already been deleted.", network.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deleteContainers(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack containers")
	defer logger.Debugf("Exiting deleting openstack containers")

	conn, err := clientconfig.NewServiceClient("object-store", opts)
	if err != nil {
		// Ignore the error if Swift is not available for the cloud
		var gerr *gophercloud.ErrEndpointNotFound
		if errors.As(err, &gerr) {
			logger.Debug("Skip container deletion because Swift endpoint is not found")
			return true, nil
		}
		logger.Error(err)
		return false, nil
	}

	listOpts := containers.ListOpts{Full: false}

	allPages, err := containers.List(conn, listOpts).AllPages()
	if err != nil {
		// Ignore the error if the user doesn't have the swiftoperator role.
		// Depending on the configuration Swift returns different error codes:
		// 403 with Keystone and 401 with internal Swauth.
		// It means we have to catch them both.
		// More information about Swith auth: https://docs.openstack.org/swift/latest/overview_auth.html
		var gerr403 gophercloud.ErrDefault403
		if errors.As(err, &gerr403) {
			logger.Debug("Skip container deletion because the user doesn't have the `swiftoperator` role")
			return true, nil
		}
		var gerr401 gophercloud.ErrDefault401
		if errors.As(err, &gerr401) {
			logger.Debug("Skip container deletion because the user doesn't have the `swiftoperator` role")
			return true, nil
		}
		logger.Error(err)
		return false, nil
	}

	allContainers, err := containers.ExtractNames(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	for _, container := range allContainers {
		metadata, err := containers.Get(conn, container, nil).ExtractMetadata()
		if err != nil {
			// Some containers that we fetched previously can already be deleted in
			// runtime. We should ignore these cases and continue to iterate through
			// the remaining containers.
			var gerr gophercloud.ErrDefault404
			if errors.As(err, &gerr) {
				continue
			}
			logger.Error(err)
			return false, nil
		}
		for key, val := range filter {
			// Swift mangles the case so openshiftClusterID becomes
			// Openshiftclusterid in the X-Container-Meta- HEAD output
			titlekey := strings.Title(strings.ToLower(key))
			if metadata[titlekey] == val {
				queue := newSemaphore(3)
				errCh := make(chan error)
				err := objects.List(conn, container, nil).EachPage(func(page pagination.Page) (bool, error) {
					objectsOnPage, err := objects.ExtractNames(page)
					if err != nil {
						return false, err
					}
					queue.Add(func() {
						for len(objectsOnPage) > 0 {
							logger.Debugf("Initiating bulk deletion of %d objects in container %q", len(objectsOnPage), container)
							resp, err := objects.BulkDelete(conn, container, objectsOnPage).Extract()
							if err != nil {
								errCh <- err
								return
							}
							if len(resp.Errors) > 0 {
								// Convert resp.Errors to golang errors.
								// Each error is represented by a list of 2 strings, where the first one
								// is the object name, and the second one contains an error message.
								for _, objectError := range resp.Errors {
									errCh <- fmt.Errorf("cannot delete object %q: %s", objectError[0], objectError[1])
								}
								logger.Debugf("Terminating object deletion routine with error. Deleted %d objects out of %d.", resp.NumberDeleted, len(objectsOnPage))
							}

							// Some object-storage instances may be set to have a limit to the LIST operation
							// that is higher to the limit to the BULK DELETE operation. On those clouds, objects
							// in the BULK DELETE call beyond the limit are silently ignored. In this loop, after
							// checking that no errors were encountered, we reduce the BULK DELETE list by the
							// number of processed objects, and send it back to the server if it's not empty.
							objectsOnPage = objectsOnPage[resp.NumberDeleted+resp.NumberNotFound:]
						}
						logger.Debugf("Terminating object deletion routine.")
					})
					return true, nil
				})
				if err != nil {
					var gerr gophercloud.ErrDefault404
					if !errors.As(err, &gerr) {
						logger.Errorf("Bulk deletion of container %q objects failed: %v", container, err)
						return false, nil
					}
				}
				var errs []error
				go func() {
					for err := range errCh {
						errs = append(errs, err)
					}
				}()

				queue.Wait()
				close(errCh)
				if len(errs) > 0 {
					return false, fmt.Errorf("errors occurred during bulk deletion of the objects of container %q: %w", container, k8serrors.NewAggregate(errs))
				}
				logger.Debugf("Deleting container %q", container)
				_, err = containers.Delete(conn, container).Extract()
				if err != nil {
					// Ignore the error if the container cannot be found and return with an appropriate message if it's another type of error
					var gerr gophercloud.ErrDefault404
					if !errors.As(err, &gerr) {
						logger.Errorf("Deleting container %q failed: %v", container, err)
						return false, nil
					}
					logger.Debugf("Cannot find container %q. It's probably already been deleted.", container)
				}
				// If a metadata key matched, we're done so break from the loop
				break
			}
		}
	}
	return true, nil
}

func deleteTrunks(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack trunks")
	defer logger.Debugf("Exiting deleting openstack trunks")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	tags := filterTags(filter)
	listOpts := trunks.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}
	allPages, err := trunks.List(conn, listOpts).AllPages()
	if err != nil {
		var gerr gophercloud.ErrDefault404
		if errors.As(err, &gerr) {
			logger.Debug("Skip trunk deletion because the cloud doesn't support trunk ports")
			return true, nil
		}
		logger.Error(err)
		return false, nil
	}

	allTrunks, err := trunks.ExtractTrunks(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	numberToDelete := len(allTrunks)
	numberDeleted := 0
	for _, trunk := range allTrunks {
		logger.Debugf("Deleting Trunk %q", trunk.ID)
		err = trunks.Delete(conn, trunk.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the trunk cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when the trunk is still in use so return/retry
				// Just log the error and move on to the next trunk
				logger.Debugf("Deleting Trunk %q failed: %v", trunk.ID, err)
				continue
			}
			logger.Debugf("Cannot find trunk %q. It's probably already been deleted.", trunk.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deleteAssociatedTrunk(conn *gophercloud.ServiceClient, logger logrus.FieldLogger, portID string) {
	logger.Debug("Deleting associated trunk")
	defer logger.Debugf("Exiting deleting associated trunk")

	listOpts := trunks.ListOpts{
		PortID: portID,
	}
	allPages, err := trunks.List(conn, listOpts).AllPages()
	if err != nil {
		var gerr gophercloud.ErrDefault404
		if errors.As(err, &gerr) {
			logger.Debug("Skip trunk deletion because the cloud doesn't support trunk ports")
			return
		}
		logger.Error(err)
		return
	}

	allTrunks, err := trunks.ExtractTrunks(allPages)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, trunk := range allTrunks {
		logger.Debugf("Deleting Trunk %q", trunk.ID)
		err = trunks.Delete(conn, trunk.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the trunk cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when the trunk is still in use so return/retry
				// Just log the error and move on to the next trunk
				logger.Debugf("Deleting Trunk %q failed: %v", trunk.ID, err)
				continue
			}
			logger.Debugf("Cannot find trunk %q. It's probably already been deleted.", trunk.ID)
		}
	}
	return
}

func deleteLoadBalancers(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack load balancers")
	defer logger.Debugf("Exiting deleting openstack load balancers")

	conn, err := clientconfig.NewServiceClient("load-balancer", opts)
	if err != nil {
		// Ignore the error if Octavia is not available for the cloud
		var gerr *gophercloud.ErrEndpointNotFound
		if errors.As(err, &gerr) {
			logger.Debug("Skip load balancer deletion because Octavia endpoint is not found")
			return true, nil
		}
		logger.Error(err)
		return false, nil
	}

	newallPages, err := apiversions.List(conn).AllPages()
	if err != nil {
		logger.Errorf("Unable to list api versions: %v", err)
		return false, nil
	}

	allAPIVersions, err := apiversions.ExtractAPIVersions(newallPages)
	if err != nil {
		logger.Errorf("Unable to extract api versions: %v", err)
		return false, nil
	}

	var octaviaTagSupport bool
	octaviaTagSupport = false
	for _, apiVersion := range allAPIVersions {
		if apiVersion.ID >= minOctaviaVersionWithTagSupport {
			octaviaTagSupport = true
		}
	}

	tags := filterTags(filter)
	var allLoadBalancers []loadbalancers.LoadBalancer
	if octaviaTagSupport {
		listOpts := loadbalancers.ListOpts{
			TagsAny: tags,
		}
		allPages, err := loadbalancers.List(conn, listOpts).AllPages()
		if err != nil {
			logger.Error(err)
			return false, nil
		}

		allLoadBalancers, err = loadbalancers.ExtractLoadBalancers(allPages)
		if err != nil {
			logger.Error(err)
			return false, nil
		}
	}

	listOpts := loadbalancers.ListOpts{
		Description: strings.Join(tags, ","),
	}

	allPages, err := loadbalancers.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allLoadBalancersWithTaggedDescription, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allLoadBalancers = append(allLoadBalancers, allLoadBalancersWithTaggedDescription...)
	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: true,
	}
	numberToDelete := len(allLoadBalancers)
	numberDeleted := 0
	for _, loadbalancer := range allLoadBalancers {
		logger.Debugf("Deleting LoadBalancer %q", loadbalancer.ID)
		err = loadbalancers.Delete(conn, loadbalancer.ID, deleteOpts).ExtractErr()
		if err != nil {
			// Ignore the error if the load balancer cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// This can fail when the load balancer is still in use so return/retry
				// Just log the error and move on to the next port
				logger.Debugf("Deleting load balancer %q failed: %v", loadbalancer.ID, err)
				continue
			}
			logger.Debugf("Cannot find load balancer %q. It's probably already been deleted.", loadbalancer.ID)
		}
		numberDeleted++
	}

	return numberDeleted == numberToDelete, nil
}

func deleteSubnetPools(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack subnet-pools")
	defer logger.Debugf("Exiting deleting openstack subnet-pools")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	tags := filterTags(filter)
	listOpts := subnetpools.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := subnetpools.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allSubnetPools, err := subnetpools.ExtractSubnetPools(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	numberToDelete := len(allSubnetPools)
	numberDeleted := 0
	for _, subnetPool := range allSubnetPools {
		logger.Debugf("Deleting Subnet Pool %q", subnetPool.ID)
		err = subnetpools.Delete(conn, subnetPool.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the subnet pool cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next subnet pool
				logger.Debugf("Deleting subnet pool %q failed: %v", subnetPool.ID, err)
				continue
			}
			logger.Debugf("Cannot find subnet pool %q. It's probably already been deleted.", subnetPool.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deleteVolumes(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting OpenStack volumes")
	defer logger.Debugf("Exiting deleting OpenStack volumes")

	var clusterID string
	for k, v := range filter {
		if strings.ToLower(k) == "openshiftclusterid" {
			clusterID = v
			break
		}
	}

	conn, err := clientconfig.NewServiceClient("volume", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	listOpts := volumes.ListOpts{}

	allPages, err := volumes.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	volumeIDs := []string{}
	for _, volume := range allVolumes {
		// First, we need to delete all volumes that have names with the cluster ID as a prefix.
		// They are created by the in-tree Cinder provisioner.
		if strings.HasPrefix(volume.Name, clusterID) {
			volumeIDs = append(volumeIDs, volume.ID)
		}
		// Second, we need to delete volumes created by the CSI driver. They contain their cluster ID
		// in the metadata.
		if val, ok := volume.Metadata[cinderCSIClusterIDKey]; ok && val == clusterID {
			volumeIDs = append(volumeIDs, volume.ID)
		}
	}

	deleteOpts := volumes.DeleteOpts{
		Cascade: false,
	}

	numberToDelete := len(volumeIDs)
	numberDeleted := 0
	for _, volumeID := range volumeIDs {
		logger.Debugf("Deleting volume %q", volumeID)
		err = volumes.Delete(conn, volumeID, deleteOpts).ExtractErr()
		if err != nil {
			// Ignore the error if the volume cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next volume
				logger.Debugf("Deleting volume %q failed: %v", volumeID, err)
				continue
			}
			logger.Debugf("Cannot find volume %q. It's probably already been deleted.", volumeID)
		}
		numberDeleted++
	}

	return numberDeleted == numberToDelete, nil
}

func deleteVolumeSnapshots(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting OpenStack volume snapshots")
	defer logger.Debugf("Exiting deleting OpenStack volume snapshots")

	var clusterID string
	for k, v := range filter {
		if strings.ToLower(k) == "openshiftclusterid" {
			clusterID = v
			break
		}
	}

	conn, err := clientconfig.NewServiceClient("volume", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	listOpts := snapshots.ListOpts{}

	allPages, err := snapshots.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allSnapshots)
	numberDeleted := 0
	for _, snapshot := range allSnapshots {
		// Delete only those snapshots that contain cluster ID in the metadata
		if val, ok := snapshot.Metadata[cinderCSIClusterIDKey]; ok && val == clusterID {
			logger.Debugf("Deleting volume snapshot %q", snapshot.ID)
			err = snapshots.Delete(conn, snapshot.ID).ExtractErr()
			if err != nil {
				// Ignore the error if the server cannot be found
				var gerr gophercloud.ErrDefault404
				if !errors.As(err, &gerr) {
					// Just log the error and move on to the next volume snapshot
					logger.Debugf("Deleting volume snapshot %q failed: %v", snapshot.ID, err)
					continue
				}
				logger.Debugf("Cannot find volume snapshot %q. It's probably already been deleted.", snapshot.ID)
			}
		}
		numberDeleted++
	}

	return numberDeleted == numberToDelete, nil
}

func deleteShares(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting OpenStack shares")
	defer logger.Debugf("Exiting deleting OpenStack shares")

	var clusterID string
	for k, v := range filter {
		if strings.ToLower(k) == "openshiftclusterid" {
			clusterID = v
			break
		}
	}

	conn, err := clientconfig.NewServiceClient("sharev2", opts)
	if err != nil {
		// Ignore the error if Manila is not available in the cloud
		var gerr *gophercloud.ErrEndpointNotFound
		if errors.As(err, &gerr) {
			logger.Debug("Skip share deletion because Manila endpoint is not found")
			return true, nil
		}
		logger.Error(err)
		return false, nil
	}

	listOpts := shares.ListOpts{
		Metadata: map[string]string{manilaCSIClusterIDKey: clusterID},
	}

	allPages, err := shares.ListDetail(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allShares, err := shares.ExtractShares(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allShares)
	numberDeleted := 0
	for _, share := range allShares {
		deleted, err := deleteShareSnapshots(conn, share.ID, logger)
		if err != nil {
			return false, err
		}
		if !deleted {
			return false, nil
		}

		logger.Debugf("Deleting share %q", share.ID)
		err = shares.Delete(conn, share.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the share cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next share
				logger.Debugf("Deleting share %q failed: %v", share.ID, err)
				continue
			}
			logger.Debugf("Cannot find share %q. It's probably already been deleted.", share.ID)
		}
		numberDeleted++
	}

	return numberDeleted == numberToDelete, nil
}

func deleteShareSnapshots(conn *gophercloud.ServiceClient, shareID string, logger logrus.FieldLogger) (bool, error) {
	logger.Debugf("Deleting OpenStack snapshots for share %v", shareID)
	defer logger.Debugf("Exiting deleting OpenStack snapshots for share %v", shareID)

	listOpts := sharesnapshots.ListOpts{
		ShareID: shareID,
	}

	allPages, err := sharesnapshots.ListDetail(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allSnapshots, err := sharesnapshots.ExtractSnapshots(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allSnapshots)
	numberDeleted := 0
	for _, snapshot := range allSnapshots {
		logger.Debugf("Deleting share snapshot %q", snapshot.ID)
		err = sharesnapshots.Delete(conn, snapshot.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the share snapshot cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next share snapshot
				logger.Debugf("Deleting share snapshot %q failed: %v", snapshot.ID, err)
				continue
			}
			logger.Debugf("Cannot find share snapshot %q. It's probably already been deleted.", snapshot.ID)
		}
		numberDeleted++
	}

	return numberDeleted == numberToDelete, nil
}

func deleteFloatingIPs(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack floating ips")
	defer logger.Debugf("Exiting deleting openstack floating ips")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}
	tags := filterTags(filter)
	listOpts := floatingips.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := floatingips.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allFloatingIPs)
	numberDeleted := 0
	for _, floatingIP := range allFloatingIPs {
		logger.Debugf("Deleting Floating IP %q", floatingIP.ID)
		err = floatingips.Delete(conn, floatingIP.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the floating ip cannot be found
			var gerr gophercloud.ErrDefault404
			if !errors.As(err, &gerr) {
				// Just log the error and move on to the next floating IP
				logger.Debugf("Deleting floating ip %q failed: %v", floatingIP.ID, err)
				continue
			}
			logger.Debugf("Cannot find floating ip %q. It's probably already been deleted.", floatingIP.ID)
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func deleteImages(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack base image")
	defer logger.Debugf("Exiting deleting openstack base image")

	conn, err := clientconfig.NewServiceClient("image", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	listOpts := images.ListOpts{
		Tags: filterTags(filter),
	}

	allPages, err := images.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToDelete := len(allImages)
	numberDeleted := 0
	for _, image := range allImages {
		logger.Debugf("Deleting image: %+v", image.ID)
		err := images.Delete(conn, image.ID).ExtractErr()
		if err != nil {
			// This can fail if the image is still in use by other VMs
			// Just log the error and move on to the next image
			logger.Debugf("Deleting Image failed: %v", err)
			continue
		}
		numberDeleted++
	}
	return numberDeleted == numberToDelete, nil
}

func untagRunner(opts *clientconfig.ClientOpts, infraID string, logger logrus.FieldLogger) error {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 10,
		Steps:    25,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return untagPrimaryNetwork(opts, infraID, logger)
	})
	if err != nil {
		if err == wait.ErrWaitTimeout {
			return err
		}
		return fmt.Errorf("unrecoverable error: %w", err)
	}

	return nil
}

func deleteRouterRunner(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) error {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 15,
		Factor:   1.3,
		Steps:    25,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return deleteRouters(opts, filter, logger)
	})
	if err != nil {
		if err == wait.ErrWaitTimeout {
			return err
		}
		return fmt.Errorf("unrecoverable error: %w", err)
	}

	return nil
}

// untagNetwork removes the tag from the primary cluster network based on unfra id
func untagPrimaryNetwork(opts *clientconfig.ClientOpts, infraID string, logger logrus.FieldLogger) (bool, error) {
	networkTag := infraID + "-primaryClusterNetwork"

	logger.Debugf("Removing tag %v from openstack networks", networkTag)
	defer logger.Debug("Exiting untagging openstack networks")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Debug(err)
		return false, nil
	}

	listOpts := networks.ListOpts{
		Tags: networkTag,
	}

	allPages, err := networks.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Debug(err)
		return false, nil
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		logger.Debug(err)
		return false, nil
	}

	if len(allNetworks) > 1 {
		return false, fmt.Errorf("more than one network with tag %s", networkTag)
	}

	if len(allNetworks) == 0 {
		// The network has already been deleted.
		return true, nil
	}

	err = attributestags.Delete(conn, "networks", allNetworks[0].ID, networkTag).ExtractErr()
	if err != nil {
		return false, nil
	}

	return true, nil
}

// validateCloud checks that the target cloud fulfills the minimum requirements
// for destroy to function.
func validateCloud(opts *clientconfig.ClientOpts, logger logrus.FieldLogger) error {
	logger.Debug("Validating the cloud")

	// A lack of support for network tagging can lead the Installer to
	// delete unmanaged resources.
	//
	// See https://bugzilla.redhat.com/show_bug.cgi?id=2013877
	logger.Debug("Validating network extensions")
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return fmt.Errorf("failed to build the network client: %w", err)
	}

	availableExtensions, err := networkextensions.Get(conn)
	if err != nil {
		return fmt.Errorf("failed to fetch network extensions: %w", err)
	}

	return networkextensions.Validate(availableExtensions)
}

// cleanClusterSgs removes the installer security groups from the user provided Port.
func cleanClusterSgs(providedPortSGs []string, clusterSGs []sg.SecGroup) []string {
	var sgs []string
	for _, providedPortSG := range providedPortSGs {
		if !isClusterSG(providedPortSG, clusterSGs) {
			sgs = append(sgs, providedPortSG)
		}
	}
	return sgs
}

func isClusterSG(providedPortSG string, clusterSGs []sg.SecGroup) bool {
	for _, clusterSG := range clusterSGs {
		if providedPortSG == clusterSG.ID {
			return true
		}
	}
	return false
}

func cleanVIPsPorts(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Cleaning provided Ports for API and Ingress VIPs")
	defer logger.Debugf("Exiting clean of provided Ports for API and Ingress VIPs")
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	tag := filter["openshiftClusterID"] + openstackdefaults.DualStackVIPsPortTag
	PortlistOpts := ports.ListOpts{
		TagsAny: tag,
	}
	allPages, err := ports.List(conn, PortlistOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		logger.Error(err)
		return false, nil
	}

	numberToClean := len(allPorts)
	numberCleaned := 0

	// Updating user provided API and Ingress Ports
	if len(allPorts) > 0 {
		clusterSGs, err := getSecurityGroups(conn, filter)
		if err != nil {
			logger.Error(err)
			return false, nil
		}
		fipByPort, err := getFIPsByPort(conn, logger)
		if err != nil {
			logger.Error(err)
			return false, nil
		}
		for _, port := range allPorts {
			logger.Debugf("Updating security groups for Port: %q", port.ID)
			sgs := cleanClusterSgs(port.SecurityGroups, clusterSGs)
			_, err := ports.Update(conn, port.ID, ports.UpdateOpts{SecurityGroups: &sgs}).Extract()
			if err != nil {
				return false, nil
			}
			if fip, ok := fipByPort[port.ID]; ok {
				logger.Debugf("Dissociating Floating IP %q", fip.ID)
				_, err := floatingips.Update(conn, fip.ID, floatingips.UpdateOpts{}).Extract()
				if err != nil {
					// Ignore the error if the floating ip cannot be found and return with an appropriate message if it's another type of error
					var gerr gophercloud.ErrDefault404
					if !errors.As(err, &gerr) {
						return false, nil
					}
					logger.Debugf("Cannot find floating ip %q. It's probably already been deleted.", fip.ID)
				}
			}

			logger.Debugf("Deleting tag for Port: %q", port.ID)
			err = attributestags.Delete(conn, "ports", port.ID, tag).ExtractErr()
			if err != nil {
				return false, nil
			}
			numberCleaned++
		}
	}
	return numberCleaned == numberToClean, nil
}
