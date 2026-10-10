package powervc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	od "github.com/openshift/installer/pkg/destroy/openstack"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

const (
	// serverDeletePollInterval and serverDeleteTimeout bound how long destroy
	// waits for Nova to finish deleting the cluster's servers.
	serverDeletePollInterval = 10 * time.Second
	serverDeleteTimeout      = 15 * time.Minute

	// serverStatusError is the Nova status of a server whose delete failed.
	serverStatusError = "ERROR"

	// containerDeletePollInterval and containerDeleteTimeout bound how long
	// destroy retries deleting a container that is not empty yet.
	containerDeletePollInterval = 5 * time.Second
	containerDeleteTimeout      = 2 * time.Minute
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	Metadata *types.ClusterMetadata
	Logger   logrus.FieldLogger
}

// New returns an PowerVC destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Metadata: metadata,
		Logger:   logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	openstackMetadata := &openstack.Metadata{
		Cloud:      o.Metadata.ClusterPlatformMetadata.PowerVC.Cloud,
		Identifier: o.Metadata.ClusterPlatformMetadata.PowerVC.Identifier,
	}
	o.Metadata.ClusterPlatformMetadata.OpenStack = openstackMetadata

	openstackDestroyer, err := od.New(o.Logger, o.Metadata)
	if err != nil {
		return nil, errors.New("destroy PowerVC cannot call New OpenStack")
	}

	// The OpenStack destroyer deletes the objects of the cluster's Swift
	// containers (for example the image registry's) with Swift bulk delete.
	// PowerVC's Swift has no bulk delete, so that fails and the OpenStack
	// destroyer exits fatally. Delete those containers first, one object at a
	// time, so the OpenStack destroyer finds no containers left to delete.
	if err := deleteClusterContainers(context.TODO(), openstackMetadata.Cloud, openstackMetadata.Identifier, o.Logger); err != nil {
		return nil, err
	}

	quota, err := openstackDestroyer.Run()
	if err != nil {
		return quota, err
	}

	// On PowerVC the bootstrap Ignition is stored in the Swift container
	// "<infraID>-ignition", which has no openshiftClusterID metadata, so the
	// OpenStack destroyer's container cleanup does not find it. Remove it by
	// name, deleting its objects one at a time: PowerVC's Swift has no bulk
	// delete, which the OpenStack destroyer relies on. It does not depend on
	// the servers, so do it before waiting for them.
	if err := deleteIgnitionContainer(context.TODO(), openstackMetadata.Cloud, o.Metadata.InfraID, o.Logger); err != nil {
		return quota, err
	}

	// The OpenStack destroyer treats a server as deleted once Nova accepts
	// the DELETE request. On PowerVC the delete can then fail inside the
	// compute service and leave the server in ERROR (for example, PowerVC
	// rejects a non-admin user's delete of a volume-backed server). Wait for
	// the servers to be gone so destroy does not report success while they
	// still exist.
	if err := waitForServersDeleted(context.TODO(), openstackMetadata.Cloud, openstackMetadata.Identifier, o.Logger); err != nil {
		return quota, err
	}

	// PowerVC does not allow Neutron port tagging, so the machine ports
	// created by CAPO are untagged and are not found by the OpenStack
	// destroyer's tag-based port cleanup. Remove them by name instead. This
	// runs after the servers are gone so that the ports are no longer
	// attached.
	if err := deleteUntaggedPorts(context.TODO(), openstackMetadata.Cloud, o.Metadata.InfraID, o.Logger); err != nil {
		return quota, err
	}

	return quota, nil
}

// ignitionContainerName returns the name of the Swift container holding the
// bootstrap Ignition, or "" when infraID is empty.
func ignitionContainerName(infraID string) string {
	if infraID == "" {
		return ""
	}
	return infraID + "-ignition"
}

// deleteIgnitionContainer deletes the bootstrap Ignition container and its
// objects. A missing container or Swift endpoint is not an error.
func deleteIgnitionContainer(ctx context.Context, cloud, infraID string, logger logrus.FieldLogger) error {
	name := ignitionContainerName(infraID)
	if name == "" {
		return nil
	}
	logger.Debugf("Deleting PowerVC bootstrap Ignition container %q", name)
	defer logger.Debug("Exiting deleting PowerVC bootstrap Ignition container")

	conn, err := newObjectStoreClient(ctx, cloud)
	if err != nil || conn == nil {
		return err
	}
	return deleteContainer(ctx, conn, name, logger)
}

// deleteClusterContainers deletes the Swift containers whose metadata matches
// the cluster filter, and their objects. A missing Swift endpoint, or a user
// who may not list containers, is not an error.
func deleteClusterContainers(ctx context.Context, cloud string, filter map[string]string, logger logrus.FieldLogger) error {
	if len(filter) == 0 {
		return nil
	}
	logger.Debug("Deleting PowerVC cluster containers")
	defer logger.Debug("Exiting deleting PowerVC cluster containers")

	conn, err := newObjectStoreClient(ctx, cloud)
	if err != nil || conn == nil {
		return err
	}

	allPages, err := containers.List(conn, nil).AllPages(ctx)
	if err != nil {
		// Same as the OpenStack destroyer: without a Swift operator role,
		// Swift returns 403 (Keystone) or 401 (Swauth).
		if gophercloud.ResponseCodeIs(err, http.StatusForbidden) || gophercloud.ResponseCodeIs(err, http.StatusUnauthorized) {
			logger.Debug("Skip container deletion because the user may not list containers")
			return nil
		}
		return fmt.Errorf("failed to list containers: %w", err)
	}
	names, err := containers.ExtractNames(allPages)
	if err != nil {
		return fmt.Errorf("failed to extract containers: %w", err)
	}
	for _, name := range names {
		metadata, err := containers.Get(ctx, conn, name, nil).ExtractMetadata()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				continue
			}
			return fmt.Errorf("failed to get metadata of container %q: %w", name, err)
		}
		if !containerMatchesFilter(metadata, filter) {
			continue
		}
		if err := deleteContainer(ctx, conn, name, logger); err != nil {
			return err
		}
	}
	return nil
}

// containerMatchesFilter reports whether the container metadata contains
// every key/value pair of the filter. Swift changes the case of metadata keys
// (openshiftClusterID is returned as Openshiftclusterid), so keys are compared
// case-insensitively. An empty filter matches nothing.
func containerMatchesFilter(metadata, filter map[string]string) bool {
	if len(filter) == 0 {
		return false
	}
	for key, val := range filter {
		found := false
		for k, v := range metadata {
			if strings.EqualFold(k, key) && v == val {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// newObjectStoreClient returns an object-store client, or nil and no error when
// the cloud has no Swift endpoint.
func newObjectStoreClient(ctx context.Context, cloud string) (*gophercloud.ServiceClient, error) {
	conn, err := openstackdefaults.NewServiceClient(ctx, "object-store", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		var endpointErr *gophercloud.ErrEndpointNotFound
		if errors.As(err, &endpointErr) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to create object-store client: %w", err)
	}
	return conn, nil
}

// objectURL returns the URL of an object. Each "/"-separated segment of the
// object name is escaped on its own, so "/" stays a path separator.
// gophercloud's objects package escapes "/" as "%2F", which PowerVC's Swift
// does not decode: it answers 404 for objects whose names contain "/", such
// as the image registry's.
func objectURL(conn *gophercloud.ServiceClient, container, object string) string {
	segments := strings.Split(object, "/")
	for i := range segments {
		segments[i] = url.PathEscape(segments[i])
	}
	return conn.ServiceURL(url.PathEscape(container), strings.Join(segments, "/"))
}

// deleteObject deletes one object. A missing object is not an error.
func deleteObject(ctx context.Context, conn *gophercloud.ServiceClient, container, object string) error {
	if _, err := conn.Delete(ctx, objectURL(conn, container, object), nil); err != nil && !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
		return err
	}
	return nil
}

// deleteContainer deletes a container's objects one at a time, then the
// container. Swift bulk delete is not used because PowerVC's Swift does not
// support it. If objects are added while it runs (Swift returns 409 Conflict
// for the container delete), it retries until containerDeleteTimeout. A missing
// container is not an error.
func deleteContainer(ctx context.Context, conn *gophercloud.ServiceClient, name string, logger logrus.FieldLogger) error {
	var lastErr error
	err := wait.PollUntilContextTimeout(ctx, containerDeletePollInterval, containerDeleteTimeout, true, func(ctx context.Context) (bool, error) {
		allPages, err := objects.List(conn, name, nil).AllPages(ctx)
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				logger.Debugf("Container %q not found, nothing to delete", name)
				return true, nil
			}
			return false, fmt.Errorf("failed to list objects in container %q: %w", name, err)
		}
		objectNames, err := objects.ExtractNames(allPages)
		if err != nil {
			return false, fmt.Errorf("failed to extract objects in container %q: %w", name, err)
		}
		logger.Debugf("Deleting %d object(s) in container %q", len(objectNames), name)
		for _, object := range objectNames {
			if err := deleteObject(ctx, conn, name, object); err != nil {
				return false, fmt.Errorf("failed to delete object %q in container %q: %w", object, name, err)
			}
		}
		if _, err := containers.Delete(ctx, conn, name).Extract(); err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return true, nil
			}
			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				lastErr = err
				logger.Debugf("Container %q is not empty yet, retrying", name)
				return false, nil
			}
			return false, fmt.Errorf("failed to delete container %q: %w", name, err)
		}
		logger.Infof("Deleted container %q", name)
		return true, nil
	})
	if err != nil && lastErr != nil && wait.Interrupted(err) {
		return fmt.Errorf("failed to delete container %q: %w", name, lastErr)
	}
	return err
}

// waitForServersDeleted waits until no server matching the cluster filter is
// left. It returns an error naming the servers that ended up in ERROR, or that
// were still present when the timeout expired.
func waitForServersDeleted(ctx context.Context, cloud string, filter map[string]string, logger logrus.FieldLogger) error {
	logger.Debug("Waiting for PowerVC servers to be deleted")
	defer logger.Debug("Exiting waiting for PowerVC servers to be deleted")

	conn, err := openstackdefaults.NewServiceClient(ctx, "compute", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return fmt.Errorf("failed to create compute client: %w", err)
	}

	var pending, failed []servers.Server
	err = wait.PollUntilContextTimeout(ctx, serverDeletePollInterval, serverDeleteTimeout, true, func(ctx context.Context) (bool, error) {
		allPages, err := servers.List(conn, servers.ListOpts{}).AllPages(ctx)
		if err != nil {
			logger.Debugf("Failed to list servers, retrying: %v", err)
			return false, nil
		}
		allServers, err := servers.ExtractServers(allPages)
		if err != nil {
			logger.Debugf("Failed to extract servers, retrying: %v", err)
			return false, nil
		}
		pending, failed = remainingServers(allServers, filter)
		if len(pending) > 0 {
			logger.Debugf("Waiting for %d server(s) to be deleted", len(pending))
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("timed out waiting for servers to be deleted: %s", describeServers(pending, false))
	}
	if len(failed) > 0 {
		return fmt.Errorf("failed to delete servers, they are in %s state: %s. PowerVC does not let non-admin users delete volume-backed servers; delete them as an admin user and run destroy again",
			serverStatusError, describeServers(failed, true))
	}
	return nil
}

// remainingServers returns the cluster's servers that are still being deleted
// (pending) and those whose delete failed (failed). A server belongs to the
// cluster when its metadata contains every key/value pair of the filter. An
// empty filter matches nothing. A server counts as failed only when it is in
// ERROR with no task in progress: a server left in ERROR by an earlier destroy
// keeps that status while Nova processes the new delete request.
func remainingServers(allServers []servers.Server, filter map[string]string) (pending, failed []servers.Server) {
	if len(filter) == 0 {
		return nil, nil
	}
	for _, server := range allServers {
		if !matchesFilter(server.Metadata, filter) {
			continue
		}
		if server.Status == serverStatusError && server.TaskState == "" {
			failed = append(failed, server)
		} else {
			pending = append(pending, server)
		}
	}
	return pending, failed
}

func matchesFilter(metadata, filter map[string]string) bool {
	for key, val := range filter {
		if v, ok := metadata[key]; !ok || v != val {
			return false
		}
	}
	return true
}

// describeServers formats servers for an error message.
func describeServers(list []servers.Server, withFault bool) string {
	descriptions := make([]string, 0, len(list))
	for _, server := range list {
		d := fmt.Sprintf("%q (%s)", server.Name, server.ID)
		if withFault && server.Fault.Message != "" {
			d += ": " + server.Fault.Message
		}
		descriptions = append(descriptions, d)
	}
	return strings.Join(descriptions, ", ")
}

// deleteUntaggedPorts deletes the ports whose names begin with the infraID.
func deleteUntaggedPorts(ctx context.Context, cloud, infraID string, logger logrus.FieldLogger) error {
	logger.Debug("Deleting PowerVC cluster ports by name")
	defer logger.Debug("Exiting deleting PowerVC cluster ports by name")

	conn, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return fmt.Errorf("failed to create network client: %w", err)
	}

	allPages, err := ports.List(conn, ports.ListOpts{}).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list ports: %w", err)
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return fmt.Errorf("failed to extract ports: %w", err)
	}

	var errs []string
	for _, port := range clusterPorts(allPorts, infraID) {
		logger.Debugf("Deleting port %q (%s)", port.Name, port.ID)
		if err := ports.Delete(ctx, conn, port.ID).ExtractErr(); err != nil && !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			errs = append(errs, fmt.Sprintf("port %q (%s): %v", port.Name, port.ID, err))
			continue
		}
		logger.Infof("Deleted port %q", port.Name)
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to delete ports: %s", strings.Join(errs, "; "))
	}
	return nil
}

// clusterPorts returns the ports that belong to the cluster, identified by
// the CAPO port naming scheme "<infraID>-<machine>-<index>".
func clusterPorts(allPorts []ports.Port, infraID string) []ports.Port {
	if infraID == "" {
		return nil
	}
	prefix := infraID + "-"
	var result []ports.Port
	for _, port := range allPorts {
		if strings.HasPrefix(port.Name, prefix) {
			result = append(result, port)
		}
	}
	return result
}
