// Package openstack contains OpenStack-specific Terraform-variable logic.
package openstack

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/utils/openstack/clientconfig"

	configv1 "github.com/openshift/api/config/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	installconfig_openstack "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	types_openstack "github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(
	installConfig *installconfig.InstallConfig,
	mastersAsset *machines.Master,
	workersAsset *machines.Worker,
	baseImage string,
	clusterID *installconfig.ClusterID,
	bootstrapIgn string,
) ([]byte, error) {
	var (
		cloud        = installConfig.Config.Platform.OpenStack.Cloud
		mastermpool  = installConfig.Config.ControlPlane.Platform.OpenStack
		defaultmpool = installConfig.Config.OpenStack.DefaultMachinePlatform
	)

	conn, err := openstackdefaults.NewServiceClient("network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, fmt.Errorf("failed to build an OpenStack service client: %w", err)
	}

	var userCA string
	{
		cloud, err := installconfig_openstack.GetSession(installConfig.Config.Platform.OpenStack.Cloud)
		if err != nil {
			return nil, fmt.Errorf("failed to get cloud config for openstack: %w", err)
		}
		// Get the ca-cert-bundle key if there is a value for cacert in clouds.yaml
		if caPath := cloud.CloudConfig.CACertFile; caPath != "" {
			caFile, err := os.ReadFile(caPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read clouds.yaml ca-cert from disk: %w", err)
			}
			userCA = string(caFile)
		}
	}

	var masterSpecs []*machinev1alpha1.OpenstackProviderSpec
	{
		masters, err := mastersAsset.Machines()
		if err != nil {
			return nil, err
		}

		for _, master := range masters {
			masterSpecs = append(masterSpecs, master.Spec.ProviderSpec.Value.Object.(*machinev1alpha1.OpenstackProviderSpec))
		}
	}

	var workerSpecs []*machinev1alpha1.OpenstackProviderSpec
	{
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return nil, err
		}

		for _, worker := range workers {
			workerSpecs = append(workerSpecs, worker.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1alpha1.OpenstackProviderSpec))
		}
	}

	var workermpool *types_openstack.MachinePool
	if len(installConfig.Config.Compute) > 0 {
		// Only considering the first Compute machinepool here, because
		// the current Installer implementation allows for one only.
		//
		// This validation code[1] errors if the pool is not named
		// "worker", and also errors in case of duplicate names,
		// factually rendering impossible to have two machinepools in
		// the install-config YAML array.
		//
		// [1]: https://github.com/openshift/installer/blob/252facf5e6e1238ee60b5f78607214e8691a3eab/pkg/types/validation/installconfig.go#L404-L410
		if len(installConfig.Config.Compute) > 1 {
			panic("Multiple machine-pools are currently not supported by the OpenShift installer on OpenStack platform")
		}
		workermpool = installConfig.Config.Compute[0].Platform.OpenStack
	}

	var userManagedLoadBalancer bool
	if lb := installConfig.Config.Platform.OpenStack.LoadBalancer; lb != nil && lb.Type == configv1.LoadBalancerTypeUserManaged {
		userManagedLoadBalancer = true
	}

	// computeAvailabilityZones is a slice where each index targets a master.
	computeAvailabilityZones := make([]string, len(masterSpecs))
	for i := range computeAvailabilityZones {
		computeAvailabilityZones[i] = masterSpecs[i].AvailabilityZone
	}

	// storageAvailabilityZones is a slice where each index targets a master.
	storageAvailabilityZones := make([]string, len(masterSpecs))
	for i := range storageAvailabilityZones {
		if masterSpecs[i].RootVolume != nil {
			storageAvailabilityZones[i] = masterSpecs[i].RootVolume.Zone
		}
	}

	// storageVolumeTypes is a slice where each index targets a master.
	storageVolumeTypes := make([]string, len(masterSpecs))
	for i := range storageVolumeTypes {
		if masterSpecs[i].RootVolume != nil {
			storageVolumeTypes[i] = masterSpecs[i].RootVolume.VolumeType
		}
	}

	// Normally baseImage contains a URL that we will use to create a new Glance image, but for testing
	// purposes we also allow to set a custom Glance image name to skip the uploading. Here we check
	// whether baseImage is a URL or not. If this is the first case, it means that the image should be
	// created by the installer from the URL. Otherwise, it means that we are given the name of the pre-created
	// Glance image, which we should use for instances.
	imageName, isURL := rhcos.GenerateOpenStackImageName(baseImage, clusterID.InfraID)
	if isURL {
		// Valid URL -> use baseImage as a URL that will be used to create new Glance image with name "<infraID>-rhcos".
		if err := uploadBaseImage(cloud, baseImage, imageName, clusterID.InfraID, installConfig.Config.Platform.OpenStack.ClusterOSImageProperties); err != nil {
			return nil, err
		}
	}

	serviceCatalog, err := getServiceCatalog(cloud)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service catalog: %w", err)
	}

	bootstrapShim, err := getBootstrapShim(cloud, clusterID.InfraID, serviceCatalog, installConfig.Config.Proxy, bootstrapIgn, userCA)
	if err != nil {
		return nil, err
	}

	octaviaSupport, err := isOctaviaSupported(serviceCatalog)
	if err != nil {
		return nil, err
	}

	var rootVolumeSize int
	if rootVolume := masterSpecs[0].RootVolume; rootVolume != nil {
		rootVolumeSize = rootVolume.Size
	}

	masterServerGroupPolicy := getServerGroupPolicy(mastermpool, defaultmpool, types_openstack.SGPolicySoftAntiAffinity)
	masterServerGroupName := masterSpecs[0].ServerGroupName
	if masterSpecs[0].ServerGroupID != "" {
		return nil, fmt.Errorf("the field ServerGroupID is not implemented in the Installer. Please use ServerGroupName for automatic creation of the Control Plane server group")
	}

	workerServerGroupPolicy := getServerGroupPolicy(workermpool, defaultmpool, types_openstack.SGPolicySoftAntiAffinity)
	var workerServerGroupNames []string
	{
		for _, workerConfig := range workerSpecs {
			workerServerGroupNames = append(workerServerGroupNames, workerConfig.ServerGroupName)
			if workerConfig.ServerGroupID != "" {
				return nil, fmt.Errorf("the field ServerGroupID is not implemented in the Installer. Please use ServerGroupName for automatic creation of the Compute server group")
			}
		}
	}

	var additionalNetworkIDs []string
	if mastermpool != nil {
		additionalNetworkIDs = mastermpool.AdditionalNetworkIDs
	}

	// defaultMachinesPort carries the machine subnets and the network.
	var defaultMachinesPort *terraformPort
	if controlPlanePort := installConfig.Config.Platform.OpenStack.ControlPlanePort; controlPlanePort != nil {
		port, err := portTargetToTerraformPort(conn, *controlPlanePort)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve portTarget :%w", err)
		}
		defaultMachinesPort = &port
		// tagging the API port if pre-created by the user.
		if apiVIPS := installConfig.Config.OpenStack.APIVIPs; len(apiVIPS) > 0 {
			// Assuming the API VIPs addresses are on the same Port
			err = tagVIPsPort(cloud, clusterID.InfraID, apiVIPS[0], port.NetworkID)
			if err != nil {
				return nil, err
			}
		}
		// tagging the Ingress port if pre-created by the user.
		if ingressVIPS := installConfig.Config.OpenStack.IngressVIPs; len(ingressVIPS) > 0 {
			// Assuming the Ingress VIPs addresses are on the same Port
			err = tagVIPsPort(cloud, clusterID.InfraID, ingressVIPS[0], port.NetworkID)
			if err != nil {
				return nil, err
			}
		}

	}

	// machinesPorts defines the primary port for a master. A nil value
	// signals Terraform to fill in the blank with the network it creates.
	// Each slice index targets a master.
	machinesPorts := make([]*terraformPort, len(masterSpecs))

	// additionalPorts translates non-control-plane
	// `failureDomain.portTarget` information in Terraform-understandable
	// syntax. Each slice index targets a master.
	additionalPorts := make([][]terraformPort, len(masterSpecs))
	for i := range masterSpecs {
		// Assign a slice to each master's index, no matter what.
		// Terraform expects each master to get an array, empty or otherwise.
		additionalPorts[i] = []terraformPort{}
		machinesPorts[i] = defaultMachinesPort
	}

	var additionalSecurityGroupIDs []string
	if mastermpool != nil {
		additionalSecurityGroupIDs = mastermpool.AdditionalSecurityGroupIDs
	}

	return json.MarshalIndent(struct {
		BaseImageName                     string                            `json:"openstack_base_image_name,omitempty"`
		ExternalNetwork                   string                            `json:"openstack_external_network,omitempty"`
		Cloud                             string                            `json:"openstack_credentials_cloud,omitempty"`
		FlavorName                        string                            `json:"openstack_master_flavor_name,omitempty"`
		APIFloatingIP                     string                            `json:"openstack_api_floating_ip,omitempty"`
		IngressFloatingIP                 string                            `json:"openstack_ingress_floating_ip,omitempty"`
		APIVIPs                           []string                          `json:"openstack_api_int_ips,omitempty"`
		IngressVIPs                       []string                          `json:"openstack_ingress_ips,omitempty"`
		TrunkSupport                      bool                              `json:"openstack_trunk_support,omitempty"`
		OctaviaSupport                    bool                              `json:"openstack_octavia_support,omitempty"`
		RootVolumeSize                    int                               `json:"openstack_master_root_volume_size,omitempty"`
		BootstrapShim                     string                            `json:"openstack_bootstrap_shim_ignition,omitempty"`
		ExternalDNS                       []string                          `json:"openstack_external_dns,omitempty"`
		MasterServerGroupName             string                            `json:"openstack_master_server_group_name,omitempty"`
		MasterServerGroupPolicy           types_openstack.ServerGroupPolicy `json:"openstack_master_server_group_policy"`
		WorkerServerGroupNames            []string                          `json:"openstack_worker_server_group_names,omitempty"`
		WorkerServerGroupPolicy           types_openstack.ServerGroupPolicy `json:"openstack_worker_server_group_policy"`
		AdditionalNetworkIDs              []string                          `json:"openstack_additional_network_ids,omitempty"`
		AdditionalPorts                   [][]terraformPort                 `json:"openstack_additional_ports"`
		AdditionalSecurityGroupIDs        []string                          `json:"openstack_master_extra_sg_ids,omitempty"`
		DefaultMachinesPort               *terraformPort                    `json:"openstack_default_machines_port,omitempty"`
		MachinesPorts                     []*terraformPort                  `json:"openstack_machines_ports"`
		MasterAvailabilityZones           []string                          `json:"openstack_master_availability_zones,omitempty"`
		MasterRootVolumeAvailabilityZones []string                          `json:"openstack_master_root_volume_availability_zones,omitempty"`
		MasterRootVolumeTypes             []string                          `json:"openstack_master_root_volume_types,omitempty"`
		UserManagedLoadBalancer           bool                              `json:"openstack_user_managed_load_balancer"`
	}{
		BaseImageName:                     imageName,
		ExternalNetwork:                   installConfig.Config.Platform.OpenStack.ExternalNetwork,
		Cloud:                             cloud,
		FlavorName:                        masterSpecs[0].Flavor,
		APIFloatingIP:                     installConfig.Config.Platform.OpenStack.APIFloatingIP,
		IngressFloatingIP:                 installConfig.Config.Platform.OpenStack.IngressFloatingIP,
		APIVIPs:                           installConfig.Config.Platform.OpenStack.APIVIPs,
		IngressVIPs:                       installConfig.Config.Platform.OpenStack.IngressVIPs,
		TrunkSupport:                      masterSpecs[0].Trunk,
		OctaviaSupport:                    octaviaSupport,
		RootVolumeSize:                    rootVolumeSize,
		BootstrapShim:                     bootstrapShim,
		ExternalDNS:                       installConfig.Config.Platform.OpenStack.ExternalDNS,
		MasterServerGroupName:             masterServerGroupName,
		MasterServerGroupPolicy:           masterServerGroupPolicy,
		WorkerServerGroupNames:            workerServerGroupNames,
		WorkerServerGroupPolicy:           workerServerGroupPolicy,
		AdditionalNetworkIDs:              additionalNetworkIDs,
		AdditionalPorts:                   additionalPorts,
		AdditionalSecurityGroupIDs:        additionalSecurityGroupIDs,
		DefaultMachinesPort:               defaultMachinesPort,
		MachinesPorts:                     machinesPorts,
		MasterAvailabilityZones:           computeAvailabilityZones,
		MasterRootVolumeAvailabilityZones: storageAvailabilityZones,
		MasterRootVolumeTypes:             storageVolumeTypes,
		UserManagedLoadBalancer:           userManagedLoadBalancer,
	}, "", "  ")
}

// tagVIPsPort tags a provided Port to allow detach of security group when destroying the cluster.
func tagVIPsPort(cloud, infraID, portIP, networkID string) error {
	conn, err := openstackdefaults.NewServiceClient("network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return err
	}

	listOpts := ports.ListOpts{
		NetworkID: networkID,
		FixedIPs:  []ports.FixedIPOpts{{IPAddress: portIP}},
	}

	allPagesPort, err := ports.List(conn, listOpts).AllPages()
	if err != nil {
		return err
	}

	allPorts, err := ports.ExtractPorts(allPagesPort)
	if err != nil {
		return err
	}

	if len(allPorts) > 1 {
		return fmt.Errorf("found multiple ports matching network ID %s and Port IP %s, cannot proceed", networkID, portIP)
	} else if len(allPorts) == 1 {
		tag := infraID + openstackdefaults.DualStackVIPsPortTag
		err = attributestags.Add(conn, "ports", allPorts[0].ID, tag).ExtractErr()
		if err != nil {
			return err
		}
	}
	return nil
}

// getServiceCatalog fetches OpenStack service catalog with service endpoints
func getServiceCatalog(cloud string) (*tokens.ServiceCatalog, error) {
	conn, err := openstackdefaults.NewServiceClient("identity", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, err
	}

	authResult := conn.GetAuthResult()
	auth, ok := authResult.(tokens.CreateResult)
	if !ok {
		return nil, fmt.Errorf("unable to extract service catalog")
	}

	serviceCatalog, err := auth.ExtractServiceCatalog()
	if err != nil {
		return nil, err
	}

	return serviceCatalog, nil
}

func isOctaviaSupported(serviceCatalog *tokens.ServiceCatalog) (bool, error) {
	_, err := openstack.V3EndpointURL(serviceCatalog, gophercloud.EndpointOpts{
		Type:         "load-balancer",
		Name:         "octavia",
		Availability: gophercloud.AvailabilityPublic,
	})
	if err != nil {
		if _, ok := err.(*gophercloud.ErrEndpointNotFound); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// We need to obtain Glance public endpoint that will be used by Ignition to download bootstrap ignition files.
// By design this should be done by using https://www.terraform.io/docs/providers/openstack/d/identity_endpoint_v3.html
// but OpenStack default policies forbid to use this API for regular users.
// On the other hand when a user authenticates in OpenStack (i.e. gets a token), it includes the whole service
// catalog in the output json. So we are able to parse the data and get the endpoint from there
// https://docs.openstack.org/api-ref/identity/v3/?expanded=token-authentication-with-scoped-authorization-detail#token-authentication-with-scoped-authorization
// Unfortunately this feature is not currently supported by Terraform, so we had to implement it here.
// We do next:
//  1. In "getServiceCatalog" we authenticate in OpenStack (tokens.Create(..)),
//     parse the token and extract the service catalog: (ExtractServiceCatalog())
//  2. In getGlancePublicURL we iterate through the catalog and find "public" endpoint for "image".
func getBootstrapShim(cloud string, infraID string, serviceCatalog *tokens.ServiceCatalog, proxy *types.Proxy, bootstrapIgn string, userCA string) (string, error) {
	clientConfigCloud, err := clientconfig.GetCloudFromYAML(openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return "", err
	}
	regionName := clientConfigCloud.RegionName
	glancePublicURL, err := openstack.V3EndpointURL(serviceCatalog, gophercloud.EndpointOpts{
		Type:         "image",
		Availability: gophercloud.AvailabilityPublic,
		Region:       regionName,
	})
	if err != nil {
		return "", fmt.Errorf("cannot retrieve Glance URL from the service catalog: %w", err)
	}

	configLocation, err := uploadBootstrapConfig(cloud, bootstrapIgn, infraID)
	if err != nil {
		return "", err
	}

	tokenID, err := getAuthToken(cloud)
	if err != nil {
		return "", err
	}

	bootstrapConfigURL := fmt.Sprintf("%s%s", glancePublicURL, configLocation)

	return generateIgnitionShim(userCA, infraID, bootstrapConfigURL, tokenID, proxy)
}

func getServerGroupPolicy(machinePool, defaultMachinePool *types_openstack.MachinePool, defaultPolicy types_openstack.ServerGroupPolicy) types_openstack.ServerGroupPolicy {
	if machinePool != nil && machinePool.ServerGroupPolicy.IsSet() {
		return machinePool.ServerGroupPolicy
	}
	if defaultMachinePool != nil && defaultMachinePool.ServerGroupPolicy.IsSet() {
		return defaultMachinePool.ServerGroupPolicy
	}
	return defaultPolicy
}
