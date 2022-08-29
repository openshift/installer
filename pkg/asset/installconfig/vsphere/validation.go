package vsphere

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

type validationContext struct {
	User        string
	AuthManager AuthManager
	Finder      Finder
	Client      *vim25.Client
}

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}
	return validation.ValidatePlatform(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere")).ToAggregate()
}

func getVCenterClient(deploymentZone vsphere.DeploymentZone, ic *types.InstallConfig) (*validationContext, ClientLogout, error) {
	server := deploymentZone.Server
	ctx := context.TODO()
	for _, vcenter := range ic.VSphere.VCenters {
		if vcenter.Server == server {
			vim25Client, _, cleanup, err := CreateVSphereClients(ctx,
				vcenter.Server,
				vcenter.Username,
				vcenter.Password)

			validationCtx := validationContext{
				User:        vcenter.Username,
				AuthManager: newAuthManager(vim25Client),
				Finder:      find.NewFinder(vim25Client),
				Client:      vim25Client,
			}
			return &validationCtx, cleanup, err
		}
	}
	return nil, nil, fmt.Errorf("vcenter %s not defined in vcenters", server)
}

func getAssociatedFailureDomain(deploymentZone vsphere.DeploymentZone, ic *types.InstallConfig) (*vsphere.FailureDomain, error) {
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
	var clients = make(map[string]*validationContext, 0)
	for _, deploymentZone := range ic.VSphere.DeploymentZones {
		failureDomain, err := getAssociatedFailureDomain(deploymentZone, ic)
		if err != nil {
			return err
		}
		if _, exists := clients[deploymentZone.Server]; !exists {
			validationCtx, cleanup, err := getVCenterClient(deploymentZone, ic)
			if err != nil {
				return err
			}
			defer cleanup()
			clients[deploymentZone.Server] = validationCtx
		}

		validationCtx := clients[deploymentZone.Server]
		allErrs = append(allErrs, validateMultiZoneProvisioning(validationCtx, failureDomain, &deploymentZone)...)
	}
	return allErrs.ToAggregate()
}

func validateMultiZoneProvisioning(validationCtx *validationContext, failureDomain *vsphere.FailureDomain, deploymentZone *vsphere.DeploymentZone) field.ErrorList {
	allErrs := field.ErrorList{}
	resourcePool := fmt.Sprintf("%s/Resources", failureDomain.Topology.ComputeCluster)
	if len(deploymentZone.PlacementConstraint.ResourcePool) != 0 {
		resourcePool = deploymentZone.PlacementConstraint.ResourcePool
	}

	vsphereField := field.NewPath("platform").Child("vsphere")
	topologyField := vsphereField.Child("failureDomains").Child("topology")
	placementConstraintField := vsphereField.Child("deploymentZones").Child("placementConstraint")

	allErrs = append(allErrs, resourcePoolExists(validationCtx, resourcePool, placementConstraintField.Child("resourcePool"))...)

	if len(deploymentZone.PlacementConstraint.Folder) > 0 {
		allErrs = append(allErrs, folderExists(validationCtx, deploymentZone.PlacementConstraint.Folder, placementConstraintField.Child("folder"))...)
	}

	computeCluster := failureDomain.Topology.ComputeCluster
	clusterPathRegexp := regexp.MustCompile("^\\/(.*?)\\/host\\/(.*?)$")
	clusterPathParts := clusterPathRegexp.FindStringSubmatch(computeCluster)
	if len(clusterPathParts) < 3 {
		return append(allErrs, field.Invalid(topologyField.Child("computeCluster"), computeCluster, "full path of cluster is required"))
	}
	computeClusterName := clusterPathParts[2]
	errs := validateVcenterPrivileges(validationCtx, topologyField.Child("server"))
	if len(errs) > 0 {
		return append(allErrs, errs...)
	}
	errs = computeClusterExists(validationCtx, computeCluster, topologyField.Child("computeCluster"))
	if len(errs) > 0 {
		return append(allErrs, errs...)
	}
	errs = datacenterExists(validationCtx, failureDomain.Topology.Datacenter, topologyField.Child("datacenter"))
	if len(errs) > 0 {
		return append(allErrs, errs...)
	}
	errs = datastoreExists(validationCtx, failureDomain.Topology.Datacenter, failureDomain.Topology.Datastore, topologyField.Child("datastore"))
	if len(errs) > 0 {
		return append(allErrs, errs...)
	}

	for _, network := range failureDomain.Topology.Networks {
		allErrs = append(allErrs, validateNetwork(validationCtx, failureDomain.Topology.Datacenter, computeClusterName, network, topologyField)...)
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
	validationCtx := &validationContext{
		User:        ic.VSphere.Username,
		AuthManager: object.NewAuthorizationManager(vim25Client),
		Finder:      finder,
		Client:      vim25Client,
	}

	return validateProvisioning(validationCtx, ic)
}

func validateProvisioning(validationCtx *validationContext, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	platform := ic.Platform.VSphere
	vsphereField := field.NewPath("platform").Child("vsphere")
	allErrs = append(allErrs, validation.ValidateForProvisioning(platform, vsphereField)...)
	allErrs = append(allErrs, validateVcenterPrivileges(validationCtx, vsphereField.Child("vcenter"))...)
	allErrs = append(allErrs, folderExists(validationCtx, ic.VSphere.Folder, vsphereField.Child("folder"))...)
	allErrs = append(allErrs, resourcePoolExists(validationCtx, ic.VSphere.ResourcePool, vsphereField.Child("resourcePool"))...)

	// if the datacenter or cluster fail to be found or is missing privileges, this will cascade through the balance
	// of checks.  exit if they fail to limit multiple errors from being thrown.
	errs := datacenterExists(validationCtx, platform.Datacenter, vsphereField.Child("datacenter"))
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
		return allErrs.ToAggregate()
	}
	computeCluster := platform.Cluster
	if computeCluster == "" {
		return field.Required(vsphereField.Child("cluster"), "must specify the cluster")
	}
	clusterPathRegexp := regexp.MustCompile("^\\/(.*?)\\/host\\/(.*?)$")
	clusterPathParts := clusterPathRegexp.FindStringSubmatch(computeCluster)
	if len(clusterPathParts) < 3 {
		computeCluster = fmt.Sprintf("/%s/host/%s", platform.Datacenter, computeCluster)
	}
	errs = computeClusterExists(validationCtx, computeCluster, vsphereField.Child("cluster"))
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
		return allErrs.ToAggregate()
	}

	errs = validateNetwork(validationCtx, platform.Datacenter, platform.Cluster, platform.Network, vsphereField.Child("network"))
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
		return allErrs.ToAggregate()
	}

	allErrs = append(allErrs, datastoreExists(validationCtx, platform.Datacenter, platform.DefaultDatastore, vsphereField.Child("defaultDatastore"))...)
	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(validationCtx *validationContext, folderPath string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	finder := validationCtx.Finder
	// If no folder is specified, skip this check as the folder will be created.
	if folderPath == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	folder, err := finder.Folder(ctx, folderPath)
	if err != nil {
		return append(allErrs, field.Invalid(fldPath, folderPath, err.Error()))
	}
	permissionGroup := permissions[permissionFolder]

	err = comparePrivileges(ctx, validationCtx, folder.Reference(), permissionGroup)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}
	return allErrs
}

func validateNetwork(validationCtx *validationContext, datacenterName string, clusterName string, networkName string, fldPath *field.Path) field.ErrorList {
	finder := validationCtx.Finder
	client := validationCtx.Client

	// It's not possible to validate a networkName if datacenterName or clusterName are empty strings
	if datacenterName == "" || clusterName == "" || networkName == "" {
		return field.ErrorList{}
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

	network, err := GetNetworkMo(ctx, client, finder, trimmedPath, clusterName, networkName)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, networkName, err.Error())}
	}
	permissionGroup := permissions[permissionPortgroup]
	err = comparePrivileges(ctx, validationCtx, network.Reference(), permissionGroup)
	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}
	return field.ErrorList{}
}

// resourcePoolExists returns an error if a resourcePool is specified in the vSphere platform but a resourcePool with that name is not found in the datacenter.
func computeClusterExists(validationCtx *validationContext, computeCluster string, fldPath *field.Path) field.ErrorList {
	if computeCluster == "" {
		return field.ErrorList{field.Required(fldPath, "must specify the cluster")}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	computeClusterMo, err := validationCtx.Finder.ClusterComputeResource(ctx, computeCluster)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, computeCluster, err.Error())}
	}
	permissionGroup := permissions[permissionCluster]
	err = comparePrivileges(ctx, validationCtx, computeClusterMo.Reference(), permissionGroup)

	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}

	return field.ErrorList{}
}

// resourcePoolExists returns an error if a resourcePool is specified in the vSphere platform but a resourcePool with that name is not found in the datacenter.
func resourcePoolExists(validationCtx *validationContext, resourcePool string, fldPath *field.Path) field.ErrorList {
	finder := validationCtx.Finder

	// If no resourcePool is specified, skip this check as the root resourcePool will be used.
	if resourcePool == "" {
		return field.ErrorList{}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	resourcePoolMo, err := finder.ResourcePool(ctx, resourcePool)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, resourcePool, err.Error())}
	}
	permissionGroup := permissions[permissionResourcePool]
	err = comparePrivileges(ctx, validationCtx, resourcePoolMo.Reference(), permissionGroup)
	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}
	return field.ErrorList{}
}

// datacenterExists returns an error if a datacenter is specified in the vSphere platform but a datacenter with that
// name is not found in the datacenter or the user does not hold adequate privileges for the datacenter.
func datacenterExists(validationCtx *validationContext, datacenterName string, fldPath *field.Path) field.ErrorList {
	finder := validationCtx.Finder

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	dataCenter, err := finder.Datacenter(ctx, datacenterName)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, datacenterName, err.Error())}
	}
	permissionGroup := permissions[permissionDatacenter]
	err = comparePrivileges(ctx, validationCtx, dataCenter.Reference(), permissionGroup)
	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}
	return field.ErrorList{}
}

// datastoreExists returns an error if a datastore is specified in the vSphere platform but a datastore with that
// name is not found in the datacenter or the user does not hold adequate privileges for the datastore.
func datastoreExists(validationCtx *validationContext, datacenterName string, datastoreName string, fldPath *field.Path) field.ErrorList {
	finder := validationCtx.Finder

	if datastoreName == "" {
		return field.ErrorList{field.Required(fldPath, "must specify the datastore")}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	dataCenter, err := finder.Datacenter(ctx, datacenterName)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, datacenterName, errors.Wrapf(err, "unable to find datacenter %s", datacenterName).Error())}
	}

	datastorePath := fmt.Sprintf("%s/datastore/...", dataCenter.InventoryPath)
	datastores, err := finder.DatastoreList(ctx, datastorePath)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, datastoreName, err.Error())}
	}

	var datastoreMo *vim25types.ManagedObjectReference
	for _, datastore := range datastores {
		if datastore.Name() == datastoreName {
			mo := datastore.Reference()
			datastoreMo = &mo
		}
	}

	if datastoreMo == nil {
		return field.ErrorList{field.Invalid(fldPath, datastoreName, fmt.Sprintf("could not find datastore %s", datastoreName))}
	}
	permissionGroup := permissions[permissionDatastore]
	err = comparePrivileges(ctx, validationCtx, datastoreMo.Reference(), permissionGroup)

	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}
	return field.ErrorList{}
}

// validateVcenterPrivileges verifies the privileges associated with
func validateVcenterPrivileges(validationCtx *validationContext, fldPath *field.Path) field.ErrorList {
	finder := validationCtx.Finder
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	rootFolder, err := finder.Folder(ctx, "/")
	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}
	permissionGroup := permissions[permissionVcenter]
	err = comparePrivileges(ctx, validationCtx, rootFolder.Reference(), permissionGroup)
	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}
	return field.ErrorList{}
}
