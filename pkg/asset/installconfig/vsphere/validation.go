package vsphere

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	return validation.ValidatePlatform(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere")).ToAggregate()
}

func getVCenterClient(deploymentZone vsphere.DeploymentZoneSpec, ic *types.InstallConfig) (*vim25.Client, ClientLogout, error) {
	server := deploymentZone.Server
	for _, vcenter := range ic.VSphere.VCenters {
		if vcenter.Server == server {
			vim25Client, _, cleanup, err := CreateVSphereClients(context.TODO(),
				vcenter.Server,
				vcenter.Username,
				vcenter.Password)
			return vim25Client, cleanup, err
		}
	}
	return nil, nil, fmt.Errorf("vcenter %s not defined in vcenters", server)
}

func getAssociatedFailureDomain(deploymentZone vsphere.DeploymentZoneSpec, ic *types.InstallConfig) (*vsphere.FailureDomainSpec, error) {
	failureDomainName := deploymentZone.FailureDomain
	for _, failureDomain := range ic.VSphere.FailureDomains {
		if failureDomainName == failureDomain.Name {
			return &failureDomain, nil
		}
	}
	return nil, fmt.Errorf("failure domain %s not defined in failureDomains", failureDomainName)
}

// ValidateMultiZoneForProvisioning performs platform validation specifically
// for multi-zone installer-provisioned infrastructure. In this case,
// self-hosted networking is a requirement when the installer creates
// infrastructure for vSphere clusters.
func ValidateMultiZoneForProvisioning(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	var clients = make(map[string]*vim25.Client, 0)
	for _, deploymentZone := range ic.VSphere.DeploymentZones {
		var client *vim25.Client
		failureDomain, err := getAssociatedFailureDomain(deploymentZone, ic)
		if err != nil {
			return err
		}
		if _, exists := clients[deploymentZone.Server]; !exists {
			client, cleanup, err := getVCenterClient(deploymentZone, ic)
			if err != nil {
				return err
			}
			defer cleanup()
			clients[deploymentZone.Server] = client
		}

		client = clients[deploymentZone.Server]

		finder := NewFinder(client)
		allErrs = append(allErrs, validateMultiZoneProvisioning(client, finder, failureDomain, &deploymentZone)...)
	}
	return allErrs.ToAggregate()
}

func validateMultiZoneProvisioning(client *vim25.Client, finder Finder, failureDomain *vsphere.FailureDomainSpec, deploymentZone *vsphere.DeploymentZoneSpec) field.ErrorList {
	allErrs := field.ErrorList{}
	resourcePool := fmt.Sprintf("%s/Resources", failureDomain.Topology.ComputeCluster)
	if len(deploymentZone.PlacementConstraint.ResourcePool) != 0 {
		resourcePool = deploymentZone.PlacementConstraint.ResourcePool
	}

	vsphereField := field.NewPath("platform").Child("vsphere")
	topologyField := vsphereField.Child("failureDomains").Child("topology")
	placementConstraintField := vsphereField.Child("deploymentZones").Child("placementConstraint")

	allErrs = append(allErrs, resourcePoolExists(finder, resourcePool, placementConstraintField.Child("resourcePool"))...)

	if len(deploymentZone.PlacementConstraint.Folder) > 0 {
		allErrs = append(allErrs, folderExists(finder, deploymentZone.PlacementConstraint.Folder, placementConstraintField.Child("folder"))...)
	}

	computeCluster := failureDomain.Topology.ComputeCluster
	clusterPathRegexp := regexp.MustCompile("^\\/(.*?)\\/host\\/(.*?)$")
	clusterPathParts := clusterPathRegexp.FindStringSubmatch(computeCluster)
	if len(clusterPathParts) < 3 {
		return append(allErrs, field.Invalid(topologyField.Child("computeCluster"), computeCluster, "full path of cluster is required"))
	}
	computeClusterName := clusterPathParts[2]

	err := computeClusterExists(finder, computeCluster, topologyField.Child("computeCluster"))
	if len(err) > 0 {
		return append(allErrs, err...)
	}

	for _, network := range failureDomain.Topology.Networks {
		allErrs = append(allErrs, validateNetwork(client, finder, failureDomain.Topology.Datacenter, computeClusterName, network, topologyField)...)
	}

	return allErrs
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for vSphere clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	p := ic.Platform.VSphere
	vim25Client, _, cleanup, err := CreateVSphereClients(context.TODO(),
		p.VCenter,
		p.Username,
		p.Password)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "vsphere"), errors.Wrapf(err, "unable to connect to vCenter %s.", p.VCenter)).Error())
	}
	defer cleanup()

	finder := NewFinder(vim25Client)
	return validateProvisioning(vim25Client, finder, ic)
}

func validateProvisioning(client *vim25.Client, finder Finder, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, folderExists(finder, ic.VSphere.Folder, field.NewPath("platform").Child("vsphere").Child("folder"))...)
	allErrs = append(allErrs, resourcePoolExists(finder, ic.VSphere.ResourcePool, field.NewPath("platform").Child("vsphere").Child("resourcePool"))...)
	if p := ic.Platform.VSphere; p.Network != "" {
		allErrs = append(allErrs, validateNetwork(client, finder, p.Datacenter, p.Cluster, p.Network, field.NewPath("platform").Child("vsphere").Child("network"))...)
	}

	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(finder Finder, folderPath string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// If no folder is specified, skip this check as the folder will be created.
	if folderPath == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err := finder.Folder(ctx, folderPath); err != nil {
		return append(allErrs, field.Invalid(fldPath, folderPath, err.Error()))
	}
	return nil
}

func validateNetwork(client *vim25.Client, finder Finder, datacenterName string, clusterName string, networkName string, fldPath *field.Path) field.ErrorList {
	// It's not possible to validate a networkName if datacenterName or clusterName are empty strings
	if datacenterName == "" || clusterName == "" {
		return nil
	}
	datacenterPath := datacenterName
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if !strings.HasPrefix(datacenterName, "/") && !strings.HasPrefix(datacenterName, "./") {
		datacenterPath = "./" + datacenterName
	}

	dataCenter, err := finder.Datacenter(ctx, datacenterPath)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, datacenterName, err.Error())}
	}
	// Remove any trailing backslash before getting networkMoID
	trimmedPath := strings.TrimPrefix(dataCenter.InventoryPath, "/")
	_, err = GetNetworkMoID(ctx, client, finder, trimmedPath, clusterName, networkName)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, networkName, err.Error())}
	}
	return nil
}

// resourcePoolExists returns an error if a resourcePool is specified in the vSphere platform but a resourcePool with that name is not found in the datacenter.
func computeClusterExists(finder Finder, computeCluster string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	//cfg := ic.VSphere

	// If no resourcePool is specified, skip this check as the root resourcePool will be used.
	if computeCluster == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err := finder.ClusterComputeResource(ctx, computeCluster); err != nil {
		return append(allErrs, field.Invalid(fldPath, computeCluster, err.Error()))
	}
	return nil
}

// resourcePoolExists returns an error if a resourcePool is specified in the vSphere platform but a resourcePool with that name is not found in the datacenter.
func resourcePoolExists(finder Finder, resourcePool string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	//cfg := ic.VSphere

	// If no resourcePool is specified, skip this check as the root resourcePool will be used.
	if resourcePool == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err := finder.ResourcePool(ctx, resourcePool); err != nil {
		return append(allErrs, field.Invalid(fldPath, resourcePool, err.Error()))
	}
	return nil
}
