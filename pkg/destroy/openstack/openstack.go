// Package openstack provides a cluster-destroyer for openstack clusters.
package openstack

import (
	"os"
	"strings"
	"time"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/apiversions"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
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
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
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
	Logger logrus.FieldLogger
}

// New returns an OpenStack destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Cloud:  metadata.ClusterPlatformMetadata.OpenStack.Cloud,
		Filter: metadata.ClusterPlatformMetadata.OpenStack.Identifier,
		Logger: logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	deleteFuncs := map[string]deleteFunc{}
	populateDeleteFuncs(deleteFuncs)
	returnChannel := make(chan string)

	opts := &clientconfig.ClientOpts{
		Cloud: o.Cloud,
	}

	// launch goroutines
	for name, function := range deleteFuncs {
		go deleteRunner(name, function, opts, o.Filter, o.Logger, returnChannel)
	}

	// wait for them to finish
	for i := 0; i < len(deleteFuncs); i++ {
		select {
		case res := <-returnChannel:
			o.Logger.Debugf("goroutine %v complete", res)
		}
	}

	return nil
}

func deleteRunner(deleteFuncName string, dFunction deleteFunc, opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger, channel chan string) {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 10,
		Factor:   1.3,
		Steps:    100,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return dFunction(opts, filter, logger)
	})

	if err != nil {
		logger.Fatalf("Unrecoverable error/timed out: %v", err)
		os.Exit(1)
	}

	// record that the goroutine has run to completion
	channel <- deleteFuncName
	return
}

// populateDeleteFuncs is the list of functions that will be launched as
// goroutines.
func populateDeleteFuncs(funcs map[string]deleteFunc) {
	funcs["deleteServers"] = deleteServers
	funcs["deleteTrunks"] = deleteTrunks
	funcs["deleteLoadBalancers"] = deleteLoadBalancers
	funcs["deletePorts"] = deletePorts
	funcs["deleteSecurityGroups"] = deleteSecurityGroups
	funcs["deleteRouters"] = deleteRouters
	funcs["deleteSubnets"] = deleteSubnets
	funcs["deleteSubnetPools"] = deleteSubnetPools
	funcs["deleteNetworks"] = deleteNetworks
	funcs["deleteContainers"] = deleteContainers
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
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	listOpts := servers.ListOpts{
		// FIXME(shardy) when gophercloud supports tags we should
		// filter by tag here
		// https://github.com/gophercloud/gophercloud/pull/1115
		// and Nova doesn't seem to support filter by Metadata, so for
		// now we do client side filtering below based on the
		// Metadata key (which matches the server properties).
	}

	allPages, err := servers.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	serverObjects := []ObjectWithTags{}
	for _, server := range allServers {
		serverObjects = append(
			serverObjects, ObjectWithTags{
				ID:   server.ID,
				Tags: server.Metadata})
	}

	filteredServers := filterObjects(serverObjects, filter)
	for _, server := range filteredServers {
		logger.Debugf("Deleting Server: %+v", server.ID)
		err = servers.Delete(conn, server.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the server cannot be found and exit with an appropiate message if it's another type of error
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				logger.Fatalf("%v", err)
				os.Exit(1)
			}
			logger.Debugf("Cannot find server %v. It's probably already been deleted.", server.ID)
		}
	}
	return len(filteredServers) == 0, nil
}

func deletePorts(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack ports")
	defer logger.Debugf("Exiting deleting openstack ports")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	tags := filterTags(filter)
	listOpts := ports.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := ports.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, port := range allPorts {
		listOpts := floatingips.ListOpts{
			PortID: port.ID,
		}
		allPages, err := floatingips.List(conn, listOpts).AllPages()
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}
		allFIPs, err := floatingips.ExtractFloatingIPs(allPages)
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}
		for _, fip := range allFIPs {
			logger.Debugf("Deleting Floating IP: %+v", fip.ID)
			err = floatingips.Delete(conn, fip.ID).ExtractErr()
			if err != nil {
				// Ignore the error if the floating ip cannot be found and exit with an appropiate message if it's another type of error
				if _, ok := err.(gophercloud.ErrDefault404); !ok {
					logger.Fatalf("%v", err)
					os.Exit(1)
				}
				logger.Debugf("Cannot find floating ip %v. It's probably already been deleted.", fip.ID)
			}
		}

		logger.Debugf("Deleting Port: %+v", port.ID)
		err = ports.Delete(conn, port.ID).ExtractErr()
		if err != nil {
			// This can fail when port is still in use so return/retry
			logger.Debugf("Deleting Port failed: %v", err)
			return false, nil
		}
	}
	return len(allPorts) == 0, nil
}

func deleteSecurityGroups(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack security-groups")
	defer logger.Debugf("Exiting deleting openstack security-groups")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	tags := filterTags(filter)
	listOpts := sg.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := sg.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allGroups, err := sg.ExtractGroups(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, group := range allGroups {
		logger.Debugf("Deleting Security Group: %+v", group.ID)
		err = sg.Delete(conn, group.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the security group cannot be found
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				// This can fail when sg is still in use by servers
				logger.Debugf("Deleting Security Group failed: %v", err)
				return false, nil
			}
			logger.Debugf("Cannot find security group %v. It's probably already been deleted.", group.ID)
		}
	}
	return len(allGroups) == 0, nil
}

func deleteRouters(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack routers")
	defer logger.Debugf("Exiting deleting openstack routers")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	tags := filterTags(filter)
	listOpts := routers.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := routers.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allRouters, err := routers.ExtractRouters(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, router := range allRouters {
		// If a user provisioned floating ip was used, it needs to be dissociated
		// Any floating Ip's associated with routers that are going to be deleted will be dissociated
		fipOpts := floatingips.ListOpts{
			RouterID: router.ID,
		}

		fipPages, err := floatingips.List(conn, fipOpts).AllPages()
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}

		allFIPs, err := floatingips.ExtractFloatingIPs(fipPages)
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}

		for _, fip := range allFIPs {
			_, err := floatingips.Update(conn, fip.ID, floatingips.UpdateOpts{}).Extract()
			if err != nil {
				logger.Fatalf("%v", err)
			}
		}

		// Clean Gateway interface
		updateOpts := routers.UpdateOpts{
			GatewayInfo: &routers.GatewayInfo{},
		}

		_, err = routers.Update(conn, router.ID, updateOpts).Extract()
		if err != nil {
			logger.Fatalf("%v", err)
		}

		// Get router interface ports
		portListOpts := ports.ListOpts{
			DeviceID: router.ID,
		}
		allPagesPort, err := ports.List(conn, portListOpts).AllPages()
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}
		allPorts, err := ports.ExtractPorts(allPagesPort)
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}

		// map to keep track of whethere interface for subnet was already removed
		removedSubnets := make(map[string]bool)
		for _, port := range allPorts {
			for _, IP := range port.FixedIPs {
				if !removedSubnets[IP.SubnetID] {
					removeOpts := routers.RemoveInterfaceOpts{
						SubnetID: IP.SubnetID,
					}
					logger.Debugf("Removing Subnet %v from Router %v\n", IP.SubnetID, router.ID)
					_, err = routers.RemoveInterface(conn, router.ID, removeOpts).Extract()
					if err != nil {
						if _, ok := err.(gophercloud.ErrDefault404); !ok {
							// This can fail when subnet is still in use
							logger.Debugf("Removing Subnet from Router failed: %v", err)
							return false, nil
						}
						logger.Debugf("Cannot find subnet %v. It's probably already been removed from router %v.", IP.SubnetID, router.ID)
					}
					removedSubnets[IP.SubnetID] = true
				}
			}
		}
		logger.Debugf("Deleting Router: %+v\n", router.ID)
		err = routers.Delete(conn, router.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the router cannot be found and exit with an appropiate message if it's another type of error
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				logger.Fatalf("%v", err)
				os.Exit(1)
			}
			logger.Debugf("Cannot find router %v. It's probably already been deleted.", router.ID)
		}
	}
	return len(allRouters) == 0, nil
}

func deleteSubnets(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack subnets")
	defer logger.Debugf("Exiting deleting openstack subnets")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	tags := filterTags(filter)
	listOpts := subnets.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := subnets.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allSubnets, err := subnets.ExtractSubnets(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, subnet := range allSubnets {
		logger.Debugf("Deleting Subnet: %+v", subnet.ID)
		err = subnets.Delete(conn, subnet.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the subnet cannot be found
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				// This can fail when subnet is still in use
				logger.Debugf("Deleting Subnet failed: %v", err)
				return false, nil
			}
			logger.Debugf("Cannot find subnet %v. It's probably already been deleted.", subnet.ID)
		}
	}
	return len(allSubnets) == 0, nil
}

func deleteNetworks(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack networks")
	defer logger.Debugf("Exiting deleting openstack networks")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	tags := filterTags(filter)
	listOpts := networks.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := networks.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, network := range allNetworks {
		logger.Debugf("Deleting network: %+v", network.ID)
		err = networks.Delete(conn, network.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the network cannot be found
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				// This can fail when network is still in use
				logger.Debugf("Deleting Network failed: %v", err)
				return false, nil
			}
			logger.Debugf("Cannot find network %v. It's probably already been deleted.", network.ID)
		}
	}
	return len(allNetworks) == 0, nil
}

func deleteContainers(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack containers")
	defer logger.Debugf("Exiting deleting openstack containers")

	conn, err := clientconfig.NewServiceClient("object-store", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	listOpts := containers.ListOpts{Full: false}

	allPages, err := containers.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allContainers, err := containers.ExtractNames(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, container := range allContainers {
		metadata, err := containers.Get(conn, container, nil).ExtractMetadata()
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}
		for key, val := range filter {
			// Swift mangles the case so openshiftClusterID becomes
			// Openshiftclusterid in the X-Container-Meta- HEAD output
			titlekey := strings.Title(strings.ToLower(key))
			if metadata[titlekey] == val {
				listOpts := objects.ListOpts{Full: false}
				allPages, err := objects.List(conn, container, listOpts).AllPages()
				if err != nil {
					logger.Fatalf("%v", err)
					os.Exit(1)
				}
				allObjects, err := objects.ExtractNames(allPages)
				if err != nil {
					logger.Fatalf("%v", err)
					os.Exit(1)
				}
				for _, object := range allObjects {
					logger.Debugf("Deleting object: %+v\n", object)
					_, err = objects.Delete(conn, container, object, nil).Extract()
					if err != nil {
						// Ignore the error if the object cannot be found and exit with an appropiate message if it's another type of error
						if _, ok := err.(gophercloud.ErrDefault404); !ok {
							logger.Fatalf("%v", err)
							os.Exit(1)
						}
						logger.Debugf("Cannot find object %v in container %v. It's probably already been deleted.", object, container)
					}
				}
				logger.Debugf("Deleting container: %+v\n", container)
				_, err = containers.Delete(conn, container).Extract()
				if err != nil {
					// Ignore the error if the container cannot be found and exit with an appropiate message if it's another type of error
					if _, ok := err.(gophercloud.ErrDefault404); !ok {
						logger.Fatalf("%v", err)
						os.Exit(1)
					}
					logger.Debugf("Cannot find container %v. It's probably already been deleted.", container)
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
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	tags := filterTags(filter)
	listOpts := trunks.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}
	allPages, err := trunks.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allTrunks, err := trunks.ExtractTrunks(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, trunk := range allTrunks {
		logger.Debugf("Deleting Trunk: %+v", trunk.ID)
		err = trunks.Delete(conn, trunk.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the trunk cannot be found
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				// This can fail when the trunk is still in use so return/retry
				logger.Debugf("Deleting Trunk failed: %v", err)
				return false, nil
			}
			logger.Debugf("Cannot find trunk %v. It's probably already been deleted.", trunk.ID)
		}
	}
	return len(allTrunks) == 0, nil
}

func deleteLoadBalancers(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack load balancers")
	defer logger.Debugf("Exiting deleting openstack load balancers")

	conn, err := clientconfig.NewServiceClient("load-balancer", opts)
	if err != nil {
		logger.Debugf("%v", err)
		return true, nil
	}

	newallPages, err := apiversions.List(conn).AllPages()
	if err != nil {
		logger.Debug("Unable to list api versions: %v", err)
	}

	allAPIVersions, err := apiversions.ExtractAPIVersions(newallPages)
	if err != nil {
		logger.Debug("Unable to extract api versions: %v", err)
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
			logger.Fatalf("%v", err)
			os.Exit(1)
		}

		allLoadBalancers, err = loadbalancers.ExtractLoadBalancers(allPages)
		if err != nil {
			logger.Fatalf("%v", err)
			os.Exit(1)
		}
	}

	listOpts := loadbalancers.ListOpts{
		Description: strings.Join(tags, ","),
	}

	allPages, err := loadbalancers.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allLoadBalancersWithTaggedDescription, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allLoadBalancers = append(allLoadBalancers, allLoadBalancersWithTaggedDescription...)
	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: true,
	}
	for _, loadbalancer := range allLoadBalancers {
		logger.Debugf("Deleting LoadBalancer: %+v", loadbalancer.ID)
		err = loadbalancers.Delete(conn, loadbalancer.ID, deleteOpts).ExtractErr()
		if err != nil {
			// Ignore the error if the load balancer cannot be found
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				// This can fail when the load balancer is still in use so return/retry
				logger.Debugf("Deleting load balancer failed: %v", err)
				return false, nil
			}
			logger.Debugf("Cannot find load balancer %v. It's probably already been deleted.", loadbalancer.ID)
		}
	}

	return len(allLoadBalancers) == 0, nil
}

func deleteSubnetPools(opts *clientconfig.ClientOpts, filter Filter, logger logrus.FieldLogger) (bool, error) {
	logger.Debug("Deleting openstack subnet-pools")
	defer logger.Debugf("Exiting deleting openstack subnet-pools")

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	tags := filterTags(filter)
	listOpts := subnetpools.ListOpts{
		TagsAny: strings.Join(tags, ","),
	}

	allPages, err := subnetpools.List(conn, listOpts).AllPages()
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}

	allSubnetPools, err := subnetpools.ExtractSubnetPools(allPages)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(1)
	}
	for _, subnetPool := range allSubnetPools {
		logger.Debugf("Deleting Subnet Pool: %+v", subnetPool.ID)
		err = subnetpools.Delete(conn, subnetPool.ID).ExtractErr()
		if err != nil {
			// Ignore the error if the subnet pool cannot be found
			if _, ok := err.(gophercloud.ErrDefault404); !ok {
				logger.Debugf("Deleting subnet pool failed: %v", err)
				return false, nil
			}
			logger.Debugf("Cannot find subnet pool %v. It's probably already been deleted.", subnetPool.ID)
		}
	}
	return len(allSubnetPools) == 0, nil
}
