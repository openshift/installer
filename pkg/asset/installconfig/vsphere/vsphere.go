// Package vsphere collects vSphere-specific configuration.
package vsphere

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

const root = "/..."

// vCenterClient contains the login info/creds and client for the vCenter.
// They are contained in a single struct to facilitate client creation
// serving as validation of the vCenter, username, and password fields.
type vCenterClient struct {
	VCenter    string
	Username   string
	Password   string
	Client     *vim25.Client
	RestClient *rest.Client
	Logout     ClientLogout
}

// Platform collects vSphere-specific configuration.
func Platform() (*vsphere.Platform, error) {
	vCenter, err := getClients()
	if err != nil {
		return nil, err
	}
	defer vCenter.Logout()

	finder := NewFinder(vCenter.Client)
	ctx := context.TODO()

	dc, dcPath, err := getDataCenter(ctx, finder, vCenter.Client)
	if err != nil {
		return nil, err
	}

	cluster, err := getCluster(ctx, dcPath, finder, vCenter.Client)
	if err != nil {
		return nil, err
	}

	datastore, err := getDataStore(ctx, dcPath, finder, vCenter.Client)
	if err != nil {
		return nil, err
	}

	network, err := getNetwork(ctx, dc, cluster, finder, vCenter.Client)
	if err != nil {
		return nil, err
	}

	apiVIP, ingressVIP, err := getVIPs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get VIPs")
	}

	failureDomain := vsphere.FailureDomain{
		Name:   "generated-failure-domain",
		Zone:   "generated-zone",
		Region: "generated-region",
		Server: vCenter.VCenter,
		Topology: vsphere.Topology{
			Datacenter:     dc,
			ComputeCluster: cluster,
			Datastore:      datastore,
			Networks:       []string{network},
		},
	}

	vcenter := vsphere.VCenter{
		Server:      vCenter.VCenter,
		Port:        443,
		Username:    vCenter.Username,
		Password:    vCenter.Password,
		Datacenters: []string{dc},
	}

	platform := &vsphere.Platform{
		VCenters:       []vsphere.VCenter{vcenter},
		FailureDomains: []vsphere.FailureDomain{failureDomain},
		APIVIPs:        []string{apiVIP},
		IngressVIPs:    []string{ingressVIP},
	}

	return platform, nil
}

// getClients() surveys the user for username, password, & vcenter.
// Validation on the three fields is performed by creating a client.
// If creating the client fails, an error is returned.
func getClients() (*vCenterClient, error) {
	var vcenter, username, password string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "vCenter",
				Help:    "The domain name or IP address of the vCenter to be used for installation.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.Host(ans.(string))
			}),
		},
	}, &vcenter); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Username",
				Help:    "The username to login to the vCenter.",
			},
			Validate: survey.Required,
		},
	}, &username); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Password",
				Help:    "The password to login to the vCenter.",
			},
			Validate: survey.Required,
		},
	}, &password); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	// There is a noticeable delay when creating the client, so let the user know what's going on.
	logrus.Infof("Connecting to vCenter %s", vcenter)
	vim25Client, restClient, logoutFunction, err := CreateVSphereClients(context.TODO(),
		vcenter,
		username,
		password)

	// Survey does not allow validation of groups of input
	// so we perform our own validation.
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to vCenter %s. Ensure provided information is correct and client certs have been added to system trust", vcenter)
	}

	return &vCenterClient{
		VCenter:    vcenter,
		Username:   username,
		Password:   password,
		Client:     vim25Client,
		RestClient: restClient,
		Logout:     logoutFunction,
	}, nil
}

// getDataCenter searches the root for all datacenters and, if there is more than one, lets the user select
// one to use for installation. Returns the name and path of the selected datacenter. The name is used
// to generate the install config and the path is used to determine the options for cluster, datastore and network.
func getDataCenter(ctx context.Context, finder Finder, client *vim25.Client) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	dataCenters, err := finder.DatacenterList(ctx, root)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to list datacenters")
	}

	// API returns an error when no results, but let's leave this in to be defensive.
	if len(dataCenters) == 0 {
		return "", "", errors.New("did not find any datacenters")
	}
	if len(dataCenters) == 1 {
		name := strings.TrimPrefix(dataCenters[0].InventoryPath, "/")
		logrus.Infof("Defaulting to only available datacenter: %s", name)
		return name, dataCenters[0].InventoryPath, nil
	}

	dataCenterPaths := make(map[string]string)
	dataCenterChoices := make([]string, 0, len(dataCenters))
	for _, dc := range dataCenters {
		name := strings.TrimPrefix(dc.InventoryPath, "/")
		dataCenterPaths[name] = dc.InventoryPath
		dataCenterChoices = append(dataCenterChoices, name)
	}
	sort.Strings(dataCenterChoices)

	var selectedDataCenter string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Datacenter",
				Options: dataCenterChoices,
				Help:    "The Datacenter to be used for installation.",
			},
			Validate: survey.Required,
		},
	}, &selectedDataCenter); err != nil {
		return "", "", errors.Wrap(err, "failed UserInput")
	}

	return selectedDataCenter, dataCenterPaths[selectedDataCenter], nil
}

func getCluster(ctx context.Context, path string, finder Finder, client *vim25.Client) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	clusters, err := finder.ClusterComputeResourceList(ctx, formatPath(path))
	if err != nil {
		return "", errors.Wrap(err, "unable to list clusters")
	}

	// API returns an error when no results, but let's leave this in to be defensive.
	if len(clusters) == 0 {
		return "", errors.New("did not find any clusters")
	}
	if len(clusters) == 1 {
		logrus.Infof("Defaulting to only available cluster: %s", clusters[0].InventoryPath)
		return clusters[0].InventoryPath, nil
	}

	clusterChoices := make([]string, 0, len(clusters))
	for _, c := range clusters {
		clusterChoices = append(clusterChoices, c.InventoryPath)
	}
	sort.Strings(clusterChoices)

	var selectedcluster string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Cluster",
				Options: clusterChoices,
				Help:    "The cluster to be used for installation.",
			},
			Validate: survey.Required,
		},
	}, &selectedcluster); err != nil {
		return "", errors.Wrap(err, "failed UserInput")
	}

	return selectedcluster, nil
}

func getDataStore(ctx context.Context, path string, finder Finder, client *vim25.Client) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	dataStores, err := finder.DatastoreList(ctx, formatPath(path))
	if err != nil {
		return "", errors.Wrap(err, "unable to list datastores")
	}

	// API returns an error when no results, but let's leave this in to be defensive.
	if len(dataStores) == 0 {
		return "", errors.New("did not find any datastores")
	}
	if len(dataStores) == 1 {
		logrus.Infof("Defaulting to only available datastore: %s", dataStores[0].InventoryPath)
		return dataStores[0].InventoryPath, nil
	}

	dataStoreChoices := make([]string, 0, len(dataStores))
	for _, ds := range dataStores {
		dataStoreChoices = append(dataStoreChoices, ds.InventoryPath)
	}
	sort.Strings(dataStoreChoices)

	var selectedDataStore string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Default Datastore",
				Options: dataStoreChoices,
				Help:    "The default datastore to be used for installation.",
			},
			Validate: survey.Required,
		},
	}, &selectedDataStore); err != nil {
		return "", errors.Wrap(err, "failed UserInput")
	}

	return selectedDataStore, nil
}

func getNetwork(ctx context.Context, datacenter string, cluster string, finder Finder, client *vim25.Client) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Get a list of networks from the previously selected Datacenter and Cluster
	networks, err := GetClusterNetworks(ctx, finder, datacenter, cluster)
	if err != nil {
		return "", errors.Wrap(err, "unable to list networks")
	}

	// API returns an error when no results, but let's leave this in to be defensive.
	if len(networks) == 0 {
		return "", errors.New("did not find any networks")
	}
	if len(networks) == 1 {
		n, err := GetNetworkName(ctx, client, networks[0])
		if err != nil {
			return "", errors.Wrap(err, "unable to get network name")
		}
		logrus.Infof("Defaulting to only available network: %s", n)
		return n, nil
	}

	validNetworkTypes := sets.NewString(
		"DistributedVirtualPortgroup",
		"Network",
		"OpaqueNetwork",
	)

	var networkChoices []string
	for _, network := range networks {
		if validNetworkTypes.Has(network.Reference().Type) {
			// Below results in an API call. Can it be eliminated somehow?
			n, err := GetNetworkName(ctx, client, network)
			if err != nil {
				return "", errors.Wrap(err, "unable to get network name")
			}
			networkChoices = append(networkChoices, n)
		}
	}
	if len(networkChoices) == 0 {
		return "", errors.New("could not find any networks of the type DistributedVirtualPortgroup or Network")
	}
	sort.Strings(networkChoices)

	var selectednetwork string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Network",
				Options: networkChoices,
				Help:    "The network to be used for installation.",
			},
			Validate: survey.Required,
		},
	}, &selectednetwork); err != nil {
		return "", errors.Wrap(err, "failed UserInput")
	}

	return selectednetwork, nil
}

func getVIPs() (string, string, error) {
	var apiVIP, ingressVIP string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Virtual IP Address for API",
				Help:    "The VIP to be used for the OpenShift API.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.IP((ans).(string))
			}),
		},
	}, &apiVIP); err != nil {
		return "", "", errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Virtual IP Address for Ingress",
				Help:    "The VIP to be used for ingress to the cluster.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				if apiVIP == (ans.(string)) {
					return fmt.Errorf("%q should not be equal to the Virtual IP address for the API", ans.(string))
				}
				return validate.IP((ans).(string))
			}),
		},
	}, &ingressVIP); err != nil {
		return "", "", errors.Wrap(err, "failed UserInput")
	}

	return apiVIP, ingressVIP, nil
}

// formatPath is a helper function that appends "/..." to enable recursive
// find in a root object. For details, see the introduction at:
// https://godoc.org/github.com/vmware/govmomi/find
func formatPath(rootObject string) string {
	return fmt.Sprintf("%s/...", rootObject)
}
