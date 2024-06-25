package gcp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	resourcemanager "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/version"
)

var (
	defaultTimeout = 2 * time.Minute
	longTimeout    = 10 * time.Minute
)

type resourceScope string

const (
	// capgProviderOwnedLabelFmt is the format string for the label
	// used for resources created by the Cluster API GCP provider.
	capgProviderOwnedLabelFmt = "capg-cluster-%s"

	// gcpGlobalResource is an identifier to indicate that the resource(s)
	// that are being deleted are globally scoped.
	gcpGlobalResource resourceScope = "global"

	// gcpRegionalResource is an identifier to indicate that the resource(s)
	// that are being deleted are regionally scoped.
	gcpRegionalResource resourceScope = "regional"

	addressResource = iota + 1
	backendServiceResource
	bucketResource
	bucketObjectResource
	cloudControllerResource
	diskResource
	dnsResource
	firewallResource
	forwardingRuleResource
	healthCheckResource
	httpHealthCheckResource
	imageResource
	instanceResource
	instanceGroupResource
	networkResource
	policyBindingResource
	routeResource
	routerResource
	serviceAccountResource
	subnetworkResource
	targetPoolResource
	targetTCPProxiesResource
)

// ClusterUninstaller holds the various options for the cluster we want to delete
type ClusterUninstaller struct {
	Logger            logrus.FieldLogger
	Region            string
	ProjectID         string
	NetworkProjectID  string
	PrivateZoneDomain string
	ClusterID         string

	computeSvc *compute.Service
	iamSvc     *iam.Service
	dnsSvc     *dns.Service
	storageSvc *storage.Service
	rmSvc      *resourcemanager.Service

	// cpusByMachineType caches the number of CPUs per machine type, used in quota
	// calculations on deletion
	cpusByMachineType map[string]int64

	// cloudControllerUID is the cluster ID used by the cluster's cloud controller
	// to generate load balancer related resources. It can be obtained either
	// from metadata or by inferring it from existing cluster resources.
	cloudControllerUID string

	errorTracker
	requestIDTracker
	pendingItemTracker
}

// New returns a GCP destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Logger:             logger,
		Region:             metadata.ClusterPlatformMetadata.GCP.Region,
		ProjectID:          metadata.ClusterPlatformMetadata.GCP.ProjectID,
		NetworkProjectID:   metadata.ClusterPlatformMetadata.GCP.NetworkProjectID,
		PrivateZoneDomain:  metadata.ClusterPlatformMetadata.GCP.PrivateZoneDomain,
		ClusterID:          metadata.InfraID,
		cloudControllerUID: gcptypes.CloudControllerUID(metadata.InfraID),
		requestIDTracker:   newRequestIDTracker(),
		pendingItemTracker: newPendingItemTracker(),
	}, nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	ctx := context.Background()
	ssn, err := gcpconfig.GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	options := []option.ClientOption{
		option.WithCredentials(ssn.Credentials),
		option.WithUserAgent(fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw)),
	}

	o.computeSvc, err = compute.NewService(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}

	cctx, cancel := context.WithTimeout(ctx, longTimeout)
	defer cancel()

	o.cpusByMachineType = map[string]int64{}
	req := o.computeSvc.MachineTypes.AggregatedList(o.ProjectID).Fields("items/*/machineTypes(name,guestCpus),nextPageToken")
	if err := req.Pages(cctx, func(list *compute.MachineTypeAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, item := range scopedList.MachineTypes {
				o.cpusByMachineType[item.Name] = item.GuestCpus
			}
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to cache machine types")
	}

	o.iamSvc, err = iam.NewService(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create iam service")
	}

	o.dnsSvc, err = dns.NewService(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dns service")
	}

	o.storageSvc, err = storage.NewService(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create storage service")
	}

	o.rmSvc, err = resourcemanager.NewService(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create resourcemanager service")
	}

	err = wait.PollImmediateInfinite(
		time.Second*10,
		o.destroyCluster,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	quota := gcptypes.Quota(o.pendingItemTracker.removedQuota)
	return &types.ClusterQuota{GCP: &quota}, nil
}

func (o *ClusterUninstaller) destroyCluster() (bool, error) {
	stagedFuncs := [][]struct {
		name    string
		execute func(ctx context.Context) error
	}{{
		{name: "Stop instances", execute: o.stopInstances},
	}, {
		{name: "Cloud controller resources", execute: o.discoverCloudControllerResources},
	}, {
		{name: "Instances", execute: o.destroyInstances},
		{name: "Disks", execute: o.destroyDisks},
		{name: "Service accounts", execute: o.destroyServiceAccounts},
		{name: "Images", execute: o.destroyImages},
		{name: "DNS", execute: o.destroyDNS},
		{name: "Buckets", execute: o.destroyBuckets},
		{name: "Routes", execute: o.destroyRoutes},
		{name: "Firewalls", execute: o.destroyFirewalls},
		{name: "Addresses", execute: o.destroyResources},
		{name: "Forwarding rules", execute: o.destroyForwardingRules},
		{name: "Target Pools", execute: o.destroyTargetPools},
		{name: "Instance groups", execute: o.destroyInstanceGroups},
		{name: "Target TCP Proxies", execute: o.destroyTargetTCPProxies},
		{name: "Backend services", execute: o.destroyBackendServices},
		{name: "Health checks", execute: o.destroyHealthChecks},
		{name: "HTTP Health checks", execute: o.destroyHTTPHealthChecks},
		{name: "Routers", execute: o.destroyRouters},
		{name: "Subnetworks", execute: o.destroySubnetworks},
		{name: "Networks", execute: o.destroyNetworks},
	}}

	// create the main Context, so all stages can accept and make context children
	ctx := context.Background()

	done := true
	for _, stage := range stagedFuncs {
		if done {
			for _, f := range stage {
				err := f.execute(ctx)
				if err != nil {
					o.Logger.Debugf("%s: %v", f.name, err)
					done = false
				}
			}
		}
	}
	return done, nil
}

// getZoneName extracts a zone name from a zone URL
func (o *ClusterUninstaller) getZoneName(zoneURL string) string {
	return getNameFromURL("zones", zoneURL)
}

// getNameFromURL gets the item name from the full URL, ex:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a -> us-central1-a
// https://www.googleapis.com/compute/v1/projects/project-id/global/networks/something-network -> something-network
func getNameFromURL(item, url string) string {
	items := strings.Split(url, item+"/")
	if len(items) < 2 {
		return ""
	}
	return items[len(items)-1]
}

// getRegionFromZone extracts a region name from a zone name of the form: us-central1-a
// Splitting the name with the last delimiter `-`, leaves a string like: us-central1
func getRegionFromZone(zoneName string) string {
	return zoneName[:strings.LastIndex(zoneName, "-")]
}

// getDiskLimit determines the name of the quota Limit that applies to the disk type, ex:
// projects/project/zones/zone/diskTypes/pd-standard -> "ssd_total_storage"
func getDiskLimit(typeURL string) string {
	switch getNameFromURL("diskTypes", typeURL) {
	case "pd-balanced", "pd-ssd":
		return "ssd_total_storage"
	case "pd-standard":
		return "disks_total_storage"
	default:
		return "unknown"
	}
}

func (o *ClusterUninstaller) isClusterResource(name string) bool {
	return strings.HasPrefix(name, o.ClusterID+"-")
}

func (o *ClusterUninstaller) clusterIDFilter() string {
	return fmt.Sprintf("name : \"%s-*\"", o.ClusterID)
}

func (o *ClusterUninstaller) clusterLabelFilter() string {
	return fmt.Sprintf("(labels.%s = \"owned\") OR (labels.%s = \"owned\")",
		fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, o.ClusterID), fmt.Sprintf(capgProviderOwnedLabelFmt, o.ClusterID))
}

func (o *ClusterUninstaller) clusterLabelOrClusterIDFilter() string {
	return fmt.Sprintf("(%s) OR (%s)", o.clusterIDFilter(), o.clusterLabelFilter())
}

func isNoOp(err error) bool {
	if err == nil {
		return false
	}
	ae, ok := err.(*googleapi.Error)
	return ok && (ae.Code == http.StatusNotFound || ae.Code == http.StatusNotModified)
}

// aggregateError is a utility function that takes a slice of errors and an
// optional pending argument, and returns an error or nil
func aggregateError(errs []error, pending ...int) error {
	err := utilerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	if len(pending) > 0 && pending[0] > 0 {
		return errors.Errorf("%d items pending", pending[0])
	}
	return nil
}

// requestIDTracker keeps track of a set of request IDs mapped to a unique resource
// identifier
type requestIDTracker struct {
	requestIDs map[string]string
}

func newRequestIDTracker() requestIDTracker {
	return requestIDTracker{
		requestIDs: map[string]string{},
	}
}

// requestID returns a UID for a given item identifier. Unless the ID is reset, the
// same requestID will be returned every time for a given item.
func (t requestIDTracker) requestID(identifier ...string) string {
	key := strings.Join(identifier, "/")
	id, exists := t.requestIDs[key]
	if !exists {
		id = uuid.New()
		t.requestIDs[key] = id
	}
	return id
}

// resetRequestID resets the request ID used for a particular item. This
// should be called whenever a request fails, and a brand new request should be
// sent.
func (t requestIDTracker) resetRequestID(identifier ...string) {
	key := strings.Join(identifier, "/")
	delete(t.requestIDs, key)
}

// pendingItemTracker tracks a set of pending item names for a given type of resource
type pendingItemTracker struct {
	pendingItems map[string]cloudResources
	removedQuota []gcptypes.QuotaUsage
}

func newPendingItemTracker() pendingItemTracker {
	return pendingItemTracker{
		pendingItems: map[string]cloudResources{},
	}
}

// GetAllPendintItems returns a slice of all of the pending items across all types.
func (t *pendingItemTracker) GetAllPendingItems() []cloudResource {
	var items []cloudResource
	for _, is := range t.pendingItems {
		for _, i := range is {
			items = append(items, i)
		}
	}
	return items
}

// getPendingItems returns the list of resources to be deleted.
func (t *pendingItemTracker) getPendingItems(itemType string) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	return lastFound.list()
}

// insertPendingItems adds to the list of resources to be deleted.
func (t *pendingItemTracker) insertPendingItems(itemType string, items []cloudResource) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	lastFound = lastFound.insert(items...)
	t.pendingItems[itemType] = lastFound
	return lastFound.list()
}

// deletePendingItems removes from the list of resources to be deleted.
func (t *pendingItemTracker) deletePendingItems(itemType string, items []cloudResource) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	for _, item := range items {
		t.removedQuota = mergeAllUsage(t.removedQuota, item.quota)
	}
	lastFound = lastFound.delete(items...)
	t.pendingItems[itemType] = lastFound
	return lastFound.list()
}

func isErrorStatus(code int64) bool {
	return code != 0 && (code < 200 || code >= 300)
}

func operationErrorMessage(op *compute.Operation) string {
	errs := []string{}
	if op.Error != nil {
		for _, e := range op.Error.Errors {
			errs = append(errs, fmt.Sprintf("%s: %s", e.Code, e.Message))
		}
	}
	if len(errs) == 0 {
		return op.HttpErrorMessage
	}
	return strings.Join(errs, ", ")
}

type destroyFunc func(ctx context.Context, item cloudResource) (*compute.Operation, error)
type destroyCloudResource func(ctx context.Context, item cloudResource, deleteFunc destroyFunc) error

type uninstallerResource struct {
	resourceType int
	destroyFunc  destroyFunc

	addressListFunc addressListFunc
	// TODO: add other list functions below as they are implemented
}

type destroyer struct {
	name         string
	itemTypeName string
	resource     uninstallerResource
}

func (o *ClusterUninstaller) listResources(ctx context.Context, resourceName, typeName string, resource uninstallerResource) ([]cloudResource, error) {
	switch resource.resourceType {
	case addressResource:
		return o.listAddressesWithFilter(ctx, resourceName, typeName, "items(name,region,addressType),nextPageToken", o.clusterIDFilter(), resource.addressListFunc)
	case backendServiceResource:
	case bucketResource:
	case bucketObjectResource:
	case cloudControllerResource:
	case diskResource:
	case dnsResource:
	case firewallResource:
	case forwardingRuleResource:
	case healthCheckResource:
	case httpHealthCheckResource:
	case imageResource:
	case instanceResource:
	case instanceGroupResource:
	case networkResource:
	case policyBindingResource:
	case routeResource:
	case routerResource:
	case serviceAccountResource:
	case subnetworkResource:
	case targetPoolResource:
	case targetTCPProxiesResource:
	}

	return nil, fmt.Errorf("failed to find matching resource type to destroy: %d", resource.resourceType)
}

func (o *ClusterUninstaller) destroyResources(ctx context.Context) error {
	for _, resourceInfo := range []destroyer{
		{
			name:         "global address",
			itemTypeName: "internal_address",
			resource: uninstallerResource{
				resourceType:    addressResource,
				destroyFunc:     o.addressDelete,
				addressListFunc: o.addressList,
			},
		},
		{
			name:         "regional address",
			itemTypeName: "internal_address",
			resource: uninstallerResource{
				resourceType:    addressResource,
				destroyFunc:     o.regionAddressDelete,
				addressListFunc: o.regionAddressList,
			},
		},
	} {
		found, err := o.listResources(ctx, resourceInfo.name, resourceInfo.itemTypeName, resourceInfo.resource)
		if err != nil {
			return err
		}
		items := o.insertPendingItems(resourceInfo.itemTypeName, found)
		for _, item := range items {
			err := o.deleteResource(ctx, item, resourceInfo.resource.destroyFunc)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}
		if items = o.getPendingItems(resourceInfo.itemTypeName); len(items) > 0 {
			return fmt.Errorf("%s: %d items pending", resourceInfo.name, len(items))
		}
	}
	return nil
}

func (o *ClusterUninstaller) deleteResource(ctx context.Context, item cloudResource, deleteFunc destroyFunc) error {
	// TODO: see if we can get the resource type here
	o.Logger.Debugf("Deleting resource %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := deleteFunc(ctx, item)
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		// TODO: see if we can get the resource type here
		return fmt.Errorf("failed to delete resource %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		// TODO: see if we can get the resource type here
		return fmt.Errorf("failed to delete resource %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted address %s", item.name)
	}
	return nil
}
