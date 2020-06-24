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

	errorTracker
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
func (o *ClusterUninstaller) Run(context.Context) error {
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
	stagedFuncs := [][]struct {
		name    string
		execute func() error
	}{{
		{name: "Stop instances", execute: o.stopInstances},
	}, {
		{name: "Cloud controller resources", execute: o.discoverCloudControllerResources},
	}, {
		{name: "Instances", execute: o.destroyInstances},
		{name: "Disks", execute: o.destroyDisks},
		{name: "Service accounts", execute: o.destroyServiceAccounts},
		{name: "Policy bindings", execute: o.destroyIAMPolicyBindings},
		{name: "Images", execute: o.destroyImages},
		{name: "DNS", execute: o.destroyDNS},
		{name: "Buckets", execute: o.destroyBuckets},
		{name: "Routes", execute: o.destroyRoutes},
		{name: "Firewalls", execute: o.destroyFirewalls},
		{name: "Addresses", execute: o.destroyAddresses},
		{name: "Target Pools", execute: o.destroyTargetPools},
		{name: "Instance groups", execute: o.destroyInstanceGroups},
		{name: "Forwarding rules", execute: o.destroyForwardingRules},
		{name: "Backend services", execute: o.destroyBackendServices},
		{name: "Health checks", execute: o.destroyHealthChecks},
		{name: "HTTP Health checks", execute: o.destroyHTTPHealthChecks},
		{name: "Routers", execute: o.destroyRouters},
		{name: "Subnetworks", execute: o.destroySubnetworks},
		{name: "Networks", execute: o.destroyNetworks},
	}}
	done := true
	for _, stage := range stagedFuncs {
		if done {
			for _, f := range stage {
				err := f.execute()
				if err != nil {
					o.Logger.Debugf("%s: %v", f.name, err)
					done = false
				}
			}
		}
	}
	return done, nil
}

// getZoneName extracts a zone name from a zone URL of the form:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a
// Splitting the URL with the delimiter `/projects`, leaves a string like: project-id/zones/us-central1-a
func (o *ClusterUninstaller) getZoneName(zoneURL string) string {
	parts := strings.Split(zoneURL, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

func (o *ClusterUninstaller) areAllClusterInstances(instances []cloudResource) bool {
	for _, instance := range instances {
		if !o.isClusterResource(instance.name) {
			return false
		}
	}
	return true
}

func (o *ClusterUninstaller) isClusterResource(name string) bool {
	return strings.HasPrefix(name, o.ClusterID+"-")
}

func (o *ClusterUninstaller) clusterIDFilter() string {
	return fmt.Sprintf("name eq \"%s-.*\"", o.ClusterID)
}

func (o *ClusterUninstaller) clusterLabelFilter() string {
	return fmt.Sprintf("labels.kubernetes-io-cluster-%s eq \"owned\"", o.ClusterID)
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
	pendingItems map[string]cloudResources
}

func newPendingItemTracker() pendingItemTracker {
	return pendingItemTracker{
		pendingItems: map[string]cloudResources{},
	}
}

// getPendingItems returns the list of resources to be deleted.
func (t pendingItemTracker) getPendingItems(itemType string) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	return lastFound.list()
}

// insertPendingItems adds to the list of resources to be deleted.
func (t pendingItemTracker) insertPendingItems(itemType string, items []cloudResource) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	lastFound = lastFound.insert(items...)
	t.pendingItems[itemType] = lastFound
	return lastFound.list()
}

// deletePendingItems removes from the list of resources to be deleted.
func (t pendingItemTracker) deletePendingItems(itemType string, items []cloudResource) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
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
