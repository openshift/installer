package powervc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
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

	quota, err := openstackDestroyer.Run()
	if err != nil {
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
