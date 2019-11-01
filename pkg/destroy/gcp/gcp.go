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
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	storage "google.golang.org/api/storage/v1"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/version"
)

var (
	defaultTimeout = 2 * time.Minute
)

// ClusterUninstaller holds the various options for the cluster we want to delete
type ClusterUninstaller struct {
	Logger    logrus.FieldLogger
	Region    string
	ProjectID string
	ClusterID string
	Context   context.Context

	computeSvc *compute.Service
	iamSvc     *iam.Service
	dnsSvc     *dns.Service
	storageSvc *storage.Service
	rmSvc      *resourcemanager.Service

	// cloudControllerUID is the cluster ID used by the cluster's cloud controller
	// to generate load balancer related resources. It can be obtained either
	// from metadata or by inferring it from existing cluster resources.
	cloudControllerUID string

	requestIDTracker
	pendingItemTracker
}

// New returns an AWS destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Logger:             logger,
		Region:             metadata.ClusterPlatformMetadata.GCP.Region,
		ProjectID:          metadata.ClusterPlatformMetadata.GCP.ProjectID,
		ClusterID:          metadata.InfraID,
		Context:            context.Background(),
		cloudControllerUID: gcptypes.CloudControllerUID(metadata.InfraID),
		requestIDTracker:   newRequestIDTracker(),
		pendingItemTracker: newPendingItemTracker(),
	}, nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	ssn, err := gcpconfig.GetSession(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get session")
	}

	options := []option.ClientOption{
		option.WithCredentials(ssn.Credentials),
		option.WithUserAgent(fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw)),
	}

	o.computeSvc, err = compute.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create compute service")
	}

	o.iamSvc, err = iam.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create iam service")
	}

	o.dnsSvc, err = dns.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create dns service")
	}

	o.storageSvc, err = storage.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create storage service")
	}

	o.rmSvc, err = resourcemanager.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create resourcemanager service")
	}

	err = wait.PollImmediateInfinite(
		time.Second*10,
		o.destroyCluster,
	)
	return nil

}

func (o *ClusterUninstaller) destroyCluster() (bool, error) {
	destroyFuncs := []struct {
		name    string
		destroy func() error
	}{
		{name: "Compute instances", destroy: o.destroyComputeInstances},
		{name: "Disks", destroy: o.destroyDisks},
		{name: "Service accounts", destroy: o.destroyServiceAccounts},
		{name: "Policy bindings", destroy: o.destroyIAMPolicyBindings},
		{name: "Images", destroy: o.destroyImages},
		{name: "DNS", destroy: o.destroyDNS},
		{name: "Object storage", destroy: o.destroyObjectStorage},
		{name: "Routes", destroy: o.destroyRoutes},
		{name: "Firewalls", destroy: o.destroyFirewalls},
		{name: "Addresses", destroy: o.destroyAddresses},
		{name: "Target Pools", destroy: o.destroyTargetPools},
		{name: "Compute instance groups", destroy: o.destroyInstanceGroups},
		{name: "Forwarding rules", destroy: o.destroyForwardingRules},
		{name: "Backend services", destroy: o.destroyBackendServices},
		{name: "Health checks", destroy: o.destroyHealthChecks},
		{name: "HTTP Health checks", destroy: o.destroyHTTPHealthChecks},
		{name: "Cloud controller internal LBs", destroy: o.destroyCloudControllerInternalLBs},
		{name: "Cloud controller external LBs", destroy: o.destroyCloudControllerExternalLBs},
		{name: "Cloud routers", destroy: o.destroyRouters},
		{name: "Subnetworks", destroy: o.destroySubNetworks},
		{name: "Networks", destroy: o.destroyNetworks},
	}
	done := true
	for _, f := range destroyFuncs {
		err := f.destroy()
		if err != nil {
			o.Logger.Debugf("%s: %v", f.name, err)
			done = false
		}
	}
	return done, nil
}

type nameAndURL struct {
	name string
	url  string
}

func (n nameAndURL) String() string {
	return fmt.Sprintf("Name: %s, URL: %s\n", n.name, n.url)
}

type nameAndZone struct {
	name   string
	zone   string
	status string
}

func (n nameAndZone) String() string {
	return fmt.Sprintf("Name: %s, Zone: %s", n.name, n.zone)
}

// getZoneName extracts a zone name from a zone URL of the form:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a
// where the compute service's basepath is:
// https://www.googleapis.com/compute/v1/projects/
// Trimming the URL, leaves a string like: project-id/zones/us-central1-a
func (o *ClusterUninstaller) getZoneName(zoneURL string) string {
	path := strings.TrimLeft(zoneURL, o.computeSvc.BasePath)
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func (o *ClusterUninstaller) areAllClusterInstances(instances []nameAndZone) bool {
	for _, instance := range instances {
		if !o.isClusterResource(instance.name) {
			return false
		}
	}
	return true
}

func (o *ClusterUninstaller) getInstanceGroupURL(ig nameAndZone) string {
	return fmt.Sprintf("%s%s/zones/%s/instanceGroups/%s", o.computeSvc.BasePath, o.ProjectID, ig.zone, ig.name)
}

func (o *ClusterUninstaller) isClusterResource(name string) bool {
	return strings.HasPrefix(name, o.ClusterID+"-")
}

func (o *ClusterUninstaller) clusterIDFilter() string {
	return fmt.Sprintf("name eq \"%s-.*\"", o.ClusterID)
}

func isNoOp(err error) bool {
	if err == nil {
		return false
	}
	ae, ok := err.(*googleapi.Error)
	return ok && (ae.Code == http.StatusNotFound || ae.Code == http.StatusNotModified)
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	ae, ok := err.(*googleapi.Error)
	return ok && ae.Code == http.StatusNotFound
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

func (o *ClusterUninstaller) contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(o.Context, defaultTimeout)
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
	if _, exists := t.requestIDs[key]; exists {
		delete(t.requestIDs, key)
	}
}

// pendingItemTracker tracks a set of pending item names for a given type of resource
type pendingItemTracker struct {
	pendingItems map[string]sets.String
}

func newPendingItemTracker() pendingItemTracker {
	return pendingItemTracker{
		pendingItems: map[string]sets.String{},
	}
}

// setPendingItems sets the list of items pending deletion for a particular item type.
// It returns items that were previously pending that are no longer in the list
// of pending items. These are items that have been deleted.
func (t pendingItemTracker) setPendingItems(itemType string, items []string) []string {
	found := sets.NewString(items...)
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = sets.NewString()
	}
	deletedItems := lastFound.Difference(found)
	t.pendingItems[itemType] = found
	return deletedItems.List()
}

func isErrorStatus(code int64) bool {
	return code < 200 || code >= 300
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
