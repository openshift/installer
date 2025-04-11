// Package openstack contains OpenStack-specific Terraform-variable logic.
package openstack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"

	configv1 "github.com/openshift/api/config/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/rhcos"
	types_openstack "github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(
	ctx context.Context,
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

	conn, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, fmt.Errorf("failed to build an OpenStack service client: %w", err)
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

	// If baseImage is a URL, the corresponding image will be uploaded to
	// Glance in the PreTerraform hook of the Cluster asset.
	imageName, _ := rhcos.GenerateOpenStackImageName(baseImage, clusterID.InfraID)

	serviceCatalog, err := getServiceCatalog(ctx, cloud)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve service catalog: %w", err)
	}

	octaviaSupport, err := isOctaviaSupported(serviceCatalog)
	if err != nil {
		return nil, err
	}

	var rootVolumeSize int
	if rootVolume := masterSpecs[0].RootVolume; rootVolume != nil {
		rootVolumeSize = rootVolume.Size
	}

	masterServerGroupPolicy := GetServerGroupPolicy(mastermpool, defaultmpool)
	masterServerGroupName := masterSpecs[0].ServerGroupName
	if masterSpecs[0].ServerGroupID != "" {
		return nil, fmt.Errorf("the field ServerGroupID is not implemented in the Installer. Please use ServerGroupName for automatic creation of the Control Plane server group")
	}

	workerServerGroupPolicy := GetServerGroupPolicy(workermpool, defaultmpool)
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
		port, err := portTargetToTerraformPort(ctx, conn, *controlPlanePort)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve portTarget :%w", err)
		}
		defaultMachinesPort = &port
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
		OctaviaSupport:                    octaviaSupport,
		RootVolumeSize:                    rootVolumeSize,
		BootstrapShim:                     bootstrapIgn,
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

// getServiceCatalog fetches OpenStack service catalog with service endpoints
func getServiceCatalog(ctx context.Context, cloud string) (*tokens.ServiceCatalog, error) {
	conn, err := openstackdefaults.NewServiceClient(ctx, "identity", openstackdefaults.DefaultClientOpts(cloud))
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

// GetServerGroupPolicy returns the server group policy set in the given machine-pool, or in the default one, or falls back to soft-anti-affinity.
func GetServerGroupPolicy(machinePool, defaultMachinePool *types_openstack.MachinePool) types_openstack.ServerGroupPolicy {
	if machinePool != nil && machinePool.ServerGroupPolicy.IsSet() {
		return machinePool.ServerGroupPolicy
	}
	if defaultMachinePool != nil && defaultMachinePool.ServerGroupPolicy.IsSet() {
		return defaultMachinePool.ServerGroupPolicy
	}
	return types_openstack.SGPolicySoftAntiAffinity
}
