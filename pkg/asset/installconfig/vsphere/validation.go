package vsphere

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	vapitags "github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

//go:generate mockgen -source=./validation.go -destination=./mock/tagmanager_generated.go -package=mock

// TagManager defines an interface to an implementation of the AuthorizationManager to facilitate mocking.
type TagManager interface {
	ListCategories(ctx context.Context) ([]string, error)
	GetCategories(ctx context.Context) ([]vapitags.Category, error)
	GetCategory(ctx context.Context, id string) (*vapitags.Category, error)
	GetTagsForCategory(ctx context.Context, id string) ([]vapitags.Tag, error)
	GetAttachedTags(ctx context.Context, ref mo.Reference) ([]vapitags.Tag, error)
	GetAttachedTagsOnObjects(ctx context.Context, objectID []mo.Reference) ([]vapitags.AttachedTags, error)
}

const (
	esxi7U2BuildNumber    int    = 17630552
	vcenter7U2BuildNumber int    = 17694817
	vcenter7U2Version     string = "7.0.2"
)

var localLogger = logrus.New()

type validationContext struct {
	AuthManager         AuthManager
	Finder              Finder
	Client              *vim25.Client
	TagManager          TagManager
	regionTagCategoryID string
	zoneTagCategoryID   string
	rhcosStream         *stream.Stream
}

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}
	return validation.ValidatePlatform(ic.Platform.VSphere, false, field.NewPath("platform").Child("vsphere"), ic).ToAggregate()
}

func getVCenterClient(failureDomain vsphere.FailureDomain, ic *types.InstallConfig) (*validationContext, ClientLogout, error) {
	server := failureDomain.Server
	ctx := context.TODO()
	for _, vcenter := range ic.VSphere.VCenters {
		if vcenter.Server == server {
			vim25Client, vim25RestClient, cleanup, err := CreateVSphereClients(ctx,
				vcenter.Server,
				vcenter.Username,
				vcenter.Password)

			if err != nil {
				return nil, nil, err
			}

			validationCtx := validationContext{
				TagManager:  vapitags.NewManager(vim25RestClient),
				AuthManager: newAuthManager(vim25Client),
				Finder:      find.NewFinder(vim25Client),
				Client:      vim25Client,
			}
			return &validationCtx, cleanup, err
		}
	}
	return nil, nil, fmt.Errorf("vcenter %s not defined in vcenters", server)
}

// ValidateForProvisioning performs platform validation specifically
// for multi-zone installer-provisioned infrastructure. In this case,
// self-hosted networking is a requirement when the installer creates
// infrastructure for vSphere clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	// If APIVIPs and IngressVIPs is equal to zero
	// then don't validate the VIPs.
	// Instead, ensure there is a configured
	// DNS record for api and test if the load
	// balancer is configured.

	// The VIP parameters within the Infrastructure status object
	// will be empty. This will cause MCO to not deploy
	// the static pods: haproxy, keepalived and coredns.
	// This will allow the use of an external load balancer
	// and RHCOS nodes to be on multiple L2 segments.
	if len(ic.Platform.VSphere.APIVIPs) == 0 && len(ic.Platform.VSphere.IngressVIPs) == 0 {
		allErrs = append(allErrs, ensureDNS(ic, field.NewPath("platform"), nil)...)
		ensureLoadBalancer(ic)
	}

	var clients = make(map[string]*validationContext, 0)

	checkTags := false
	if len(ic.VSphere.FailureDomains) > 1 {
		checkTags = true
	}

	for i, failureDomain := range ic.VSphere.FailureDomains {
		if _, exists := clients[failureDomain.Server]; !exists {
			validationCtx, cleanup, err := getVCenterClient(failureDomain, ic)
			if err != nil {
				return err
			}
			defer cleanup()

			err = getRhcosStream(validationCtx)
			if err != nil {
				return err
			}

			allErrs = append(allErrs, validateVCenterVersion(validationCtx, field.NewPath("platform").Child("vsphere").Child("vcenters"))...)
			clients[failureDomain.Server] = validationCtx
		}

		validationCtx := clients[failureDomain.Server]
		allErrs = append(allErrs, validateFailureDomain(validationCtx, &ic.VSphere.FailureDomains[i], checkTags)...)
	}
	return allErrs.ToAggregate()
}

func validateFailureDomain(validationCtx *validationContext, failureDomain *vsphere.FailureDomain, checkTags bool) field.ErrorList {
	allErrs := field.ErrorList{}
	checkDatacenterPrivileges := true
	checkComputeClusterPrivileges := true

	resourcePool := fmt.Sprintf("%s/Resources", failureDomain.Topology.ComputeCluster)
	if len(failureDomain.Topology.ResourcePool) != 0 {
		resourcePool = failureDomain.Topology.ResourcePool
		checkComputeClusterPrivileges = false
	}

	vsphereField := field.NewPath("platform").Child("vsphere")
	topologyField := vsphereField.Child("failureDomains").Child("topology")

	if checkTags {
		regionTagCategoryID, zoneTagCategoryID, err := validateTagCategories(validationCtx)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(vsphereField, err))
		}
		validationCtx.regionTagCategoryID = regionTagCategoryID
		validationCtx.zoneTagCategoryID = zoneTagCategoryID
	}

	allErrs = append(allErrs, resourcePoolExists(validationCtx, resourcePool, topologyField.Child("resourcePool"))...)

	if len(failureDomain.Topology.Folder) > 0 {
		allErrs = append(allErrs, folderExists(validationCtx, failureDomain.Topology.Folder, topologyField.Child("folder"))...)
		checkDatacenterPrivileges = false
	}

	allErrs = append(allErrs, validateESXiVersion(validationCtx, failureDomain.Topology.ComputeCluster, vsphereField, topologyField.Child("computeCluster"))...)
	allErrs = append(allErrs, validateVcenterPrivileges(validationCtx, topologyField.Child("server"))...)
	allErrs = append(allErrs, computeClusterExists(validationCtx, failureDomain.Topology.ComputeCluster, topologyField.Child("computeCluster"), checkComputeClusterPrivileges, checkTags)...)
	allErrs = append(allErrs, datacenterExists(validationCtx, failureDomain.Topology.Datacenter, topologyField.Child("datacenter"), checkDatacenterPrivileges)...)
	allErrs = append(allErrs, datastoreExists(validationCtx, failureDomain.Topology.Datacenter, failureDomain.Topology.Datastore, topologyField.Child("datastore"))...)

	if failureDomain.Topology.Template != "" {
		allErrs = append(allErrs, validateTemplate(validationCtx, failureDomain.Topology.Template, topologyField.Child("template"))...)
	}

	for _, network := range failureDomain.Topology.Networks {
		allErrs = append(allErrs, validateNetwork(validationCtx, failureDomain.Topology.Datacenter, failureDomain.Topology.ComputeCluster, network, topologyField)...)
	}

	return allErrs
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

func validateVCenterVersion(validationCtx *validationContext, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	constraints, err := version.NewConstraint(fmt.Sprintf("< %s", vcenter7U2Version))
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
	}

	vCenterVersion, err := version.NewVersion(validationCtx.Client.ServiceContent.About.Version)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
	}
	build, err := strconv.Atoi(validationCtx.Client.ServiceContent.About.Build)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
	}

	detail := fmt.Sprintf("The vSphere storage driver requires a minimum of vSphere 7 Update 2. Current vCenter version: %s, build: %s",
		validationCtx.Client.ServiceContent.About.Version, validationCtx.Client.ServiceContent.About.Build)

	if constraints.Check(vCenterVersion) {
		allErrs = append(allErrs, field.Required(fldPath, detail))
	} else if build < vcenter7U2BuildNumber {
		allErrs = append(allErrs, field.Required(fldPath, detail))
	}

	return allErrs
}

func validateESXiVersion(validationCtx *validationContext, clusterPath string, vSphereFldPath, computeClusterFldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	finder := validationCtx.Finder

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	clusters, err := finder.ClusterComputeResourceList(ctx, clusterPath)

	if err != nil {
		var notFoundError *find.NotFoundError
		var defaultNotFoundError *find.DefaultNotFoundError

		/* These error types also exist, but it seems less likely to occur.
		var *find.MultipleFoundError
		var *find.DefaultMultipleFoundError
		*/
		switch {
		case errors.As(err, &notFoundError):
			return field.ErrorList{field.Invalid(computeClusterFldPath, clusterPath, notFoundError.Error())}
		case errors.As(err, &defaultNotFoundError):
			return field.ErrorList{field.Invalid(computeClusterFldPath, clusterPath, defaultNotFoundError.Error())}
		default:
			return append(allErrs, field.InternalError(vSphereFldPath, err))
		}
	}

	v7, err := version.NewVersion("7.0")
	if err != nil {
		return append(allErrs, field.InternalError(vSphereFldPath, err))
	}

	hosts, err := clusters[0].Hosts(context.TODO())
	if err != nil {
		err = errors.Wrapf(err, "unable to find hosts from cluster on path: %s", clusterPath)
		return append(allErrs, field.InternalError(vSphereFldPath, err))
	}

	for _, h := range hosts {
		var mh mo.HostSystem
		err := h.Properties(context.TODO(), h.Reference(), []string{"config.product"}, &mh)
		if err != nil {
			return append(allErrs, field.InternalError(vSphereFldPath, err))
		}

		esxiHostVersion, err := version.NewVersion(mh.Config.Product.Version)
		if err != nil {
			return append(allErrs, field.InternalError(vSphereFldPath, err))
		}

		detail := fmt.Sprintf("The vSphere storage driver requires a minimum of vSphere 7 Update 2. The ESXi host: %s is version: %s and build: %s",
			h.Name(), mh.Config.Product.Version, mh.Config.Product.Build)

		if esxiHostVersion.LessThan(v7) {
			allErrs = append(allErrs, field.Required(computeClusterFldPath, detail))
		} else {
			build, err := strconv.Atoi(mh.Config.Product.Build)
			if err != nil {
				return append(allErrs, field.InternalError(vSphereFldPath, err))
			}
			if build < esxi7U2BuildNumber {
				allErrs = append(allErrs, field.Required(computeClusterFldPath, detail))
			}
		}
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
func computeClusterExists(validationCtx *validationContext, computeCluster string, fldPath *field.Path, checkPrivileges, checkTagAttachment bool) field.ErrorList {
	if computeCluster == "" {
		return field.ErrorList{field.Required(fldPath, "must specify the cluster")}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	computeClusterMo, err := validationCtx.Finder.ClusterComputeResource(ctx, computeCluster)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, computeCluster, err.Error())}
	}

	if checkPrivileges {
		permissionGroup := permissions[permissionCluster]
		err = comparePrivileges(ctx, validationCtx, computeClusterMo.Reference(), permissionGroup)

		if err != nil {
			return field.ErrorList{field.InternalError(fldPath, err)}
		}
	}

	if checkTagAttachment {
		err = validateTagAttachment(validationCtx, computeClusterMo.Reference())
		if err != nil {
			return field.ErrorList{field.InternalError(fldPath, err)}
		}
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
func datacenterExists(validationCtx *validationContext, datacenterName string, fldPath *field.Path, checkPrivileges bool) field.ErrorList {
	finder := validationCtx.Finder

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	dataCenter, err := finder.Datacenter(ctx, datacenterName)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, datacenterName, err.Error())}
	}
	if checkPrivileges {
		permissionGroup := permissions[permissionDatacenter]
		err = comparePrivileges(ctx, validationCtx, dataCenter.Reference(), permissionGroup)
		if err != nil {
			return field.ErrorList{field.InternalError(fldPath, err)}
		}
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
		if datastore.InventoryPath == datastoreName || datastore.Name() == datastoreName {
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

func ensureDNS(installConfig *types.InstallConfig, fldPath *field.Path, resolver *net.Resolver) field.ErrorList {
	var uris []string
	errList := field.ErrorList{}
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	uris = append(uris, fmt.Sprintf("api.%s", installConfig.ClusterDomain()))
	uris = append(uris, fmt.Sprintf("api-int.%s", installConfig.ClusterDomain()))

	if resolver == nil {
		resolver = &net.Resolver{
			PreferGo: true,
		}
	}

	// DNS lookup uri
	for _, u := range uris {
		logrus.Debugf("Performing DNS Lookup: %s", u)
		_, err := resolver.LookupHost(ctx, u)
		// Append error if DNS entry does not exist
		if err != nil {
			errList = append(errList, field.Invalid(fldPath, u, err.Error()))
		}
	}

	return errList
}

func ensureLoadBalancer(installConfig *types.InstallConfig) {
	var lastErr error
	dialTimeout := time.Second
	tcpTimeout := time.Second * 10
	errorCount := 0
	apiURIPort := fmt.Sprintf("api.%s:%s", installConfig.ClusterDomain(), "6443")
	tcpContext, cancel := context.WithTimeout(context.TODO(), tcpTimeout)
	defer cancel()

	// If the load balancer is configured properly even
	// without members we should be available to make
	// a connection to port 6443. Check for 10 seconds
	// emit debug message every 2 failures. If unavailable
	// after timeout emit warning only.
	wait.Until(func() {
		conn, err := net.DialTimeout("tcp", apiURIPort, dialTimeout)
		if err == nil {
			conn.Close()
			cancel()
		} else {
			lastErr = err
			if errorCount == 2 {
				logrus.Debug("Still waiting for load balancer...")
				errorCount = 0
			} else {
				errorCount++
			}
		}
	}, 2*time.Second, tcpContext.Done())

	err := tcpContext.Err()
	if err != nil && !errors.Is(err, context.Canceled) {
		if lastErr != nil {
			localLogger.Warnf("Installation may fail, load balancer not available: %v", lastErr)
		}
	}
}

func validateTagCategories(validationCtx *validationContext) (string, string, error) {
	if validationCtx.TagManager == nil {
		return "", "", nil
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	categories, err := validationCtx.TagManager.GetCategories(ctx)
	if err != nil {
		return "", "", err
	}

	regionTagCategoryID := ""
	zoneTagCategoryID := ""
	for _, category := range categories {
		switch category.Name {
		case vsphere.TagCategoryRegion:
			regionTagCategoryID = category.ID
		case vsphere.TagCategoryZone:
			zoneTagCategoryID = category.ID
		}
		if len(zoneTagCategoryID) > 0 && len(regionTagCategoryID) > 0 {
			break
		}
	}
	if len(zoneTagCategoryID) == 0 || len(regionTagCategoryID) == 0 {
		return "", "", errors.New("tag categories openshift-zone and openshift-region must be created")
	}
	return regionTagCategoryID, zoneTagCategoryID, nil
}

func validateTagAttachment(validationCtx *validationContext, reference vim25types.ManagedObjectReference) error {
	if validationCtx.TagManager == nil {
		return nil
	}
	client := validationCtx.Client
	tagManager := validationCtx.TagManager
	regionTagCategoryID := validationCtx.regionTagCategoryID
	zoneTagCategoryID := validationCtx.zoneTagCategoryID
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	referencesToCheck := []mo.Reference{reference}
	ancestors, err := mo.Ancestors(ctx,
		client.RoundTripper,
		client.ServiceContent.PropertyCollector,
		reference)
	if err != nil {
		return err
	}
	for _, ancestor := range ancestors {
		referencesToCheck = append(referencesToCheck, ancestor.Reference())
	}
	attachedTags, err := tagManager.GetAttachedTagsOnObjects(ctx, referencesToCheck)
	if err != nil {
		return err
	}
	regionTagAttached := false
	zoneTagAttached := false
	for _, attachedTag := range attachedTags {
		for _, tag := range attachedTag.Tags {
			if !regionTagAttached {
				if tag.CategoryID == regionTagCategoryID {
					regionTagAttached = true
				}
			}
			if !zoneTagAttached {
				if tag.CategoryID == zoneTagCategoryID {
					zoneTagAttached = true
				}
			}
			if regionTagAttached && zoneTagAttached {
				return nil
			}
		}
	}
	var errs []string
	if !regionTagAttached {
		errs = append(errs, fmt.Sprintf("tag associated with tag category %s not attached to this resource or ancestor", vsphere.TagCategoryRegion))
	}
	if !zoneTagAttached {
		errs = append(errs, fmt.Sprintf("tag associated with tag category %s not attached to this resource or ancestor", vsphere.TagCategoryZone))
	}
	return errors.New(strings.Join(errs, ","))
}

func validateTemplate(validationCtx *validationContext, template string, fldPath *field.Path) field.ErrorList {
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	var vmMo mo.VirtualMachine
	var arch stream.Arch
	var platformArtifacts stream.PlatformArtifacts
	var ok bool

	// Not using GetArchitectures() here to make it easier to test
	if arch, ok = validationCtx.rhcosStream.Architectures["x86_64"]; !ok {
		return field.ErrorList{field.InternalError(fldPath, errors.New("unable to find vmware rhcos artifacts"))}
	}

	if platformArtifacts, ok = arch.Artifacts["vmware"]; !ok {
		return field.ErrorList{field.InternalError(fldPath, errors.New("unable to find vmware rhcos artifacts"))}
	}
	rhcosReleaseVersion := platformArtifacts.Release

	vm, err := validationCtx.Finder.VirtualMachine(ctx, template)

	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, template, errors.Wrapf(err, "unable to find template %s", template).Error())}
	}
	err = vm.Properties(ctx, vm.Reference(), nil, &vmMo)
	if err != nil {
		return field.ErrorList{field.InternalError(fldPath, err)}
	}

	if vmMo.Summary.Config.Product != nil {
		templateProductVersion := vmMo.Summary.Config.Product.Version
		if templateProductVersion == "" {
			localLogger.Warnf("unable to determine RHCOS version of virtual machine: %s, installation may fail.", template)
			return nil
		}

		err := compareCurrentToTemplate(templateProductVersion, rhcosReleaseVersion)
		if err != nil {
			return field.ErrorList{field.InternalError(fldPath, fmt.Errorf("current template: %s %w", template, err))}
		}
	} else {
		localLogger.Warnf("unable to determine RHCOS version of virtual machine: %s, installation may fail.", template)
	}

	return nil
}

func compareCurrentToTemplate(templateProductVersion, rhcosStreamVersion string) error {
	if templateProductVersion != rhcosStreamVersion {
		templateVersion, err := strconv.Atoi(strings.Split(templateProductVersion, ".")[0])
		if err != nil {
			return err
		}
		currentRhcosVersion, err := strconv.Atoi(strings.Split(rhcosStreamVersion, ".")[0])
		if err != nil {
			return err
		}

		switch versionDiff := currentRhcosVersion - templateVersion; {
		case versionDiff < 0:
			return fmt.Errorf("rhcos version: %s is too many revisions ahead current version: %s", templateProductVersion, rhcosStreamVersion)
		case versionDiff >= 2:
			return fmt.Errorf("rhcos version: %s is too many revisions behind current version: %s", templateProductVersion, rhcosStreamVersion)
		case versionDiff == 1:
			localLogger.Warnf("rhcos version: %s is behind current version: %s, installation may fail", templateProductVersion, rhcosStreamVersion)
		}
	}
	return nil
}

func getRhcosStream(validationCtx *validationContext) error {
	var err error
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	validationCtx.rhcosStream, err = rhcos.FetchCoreOSBuild(ctx)

	if err != nil {
		return err
	}
	return nil
}
