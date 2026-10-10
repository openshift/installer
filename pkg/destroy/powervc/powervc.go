package powervc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	od "github.com/openshift/installer/pkg/destroy/openstack"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
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

	// PowerVC does not allow Neutron port tagging, so the machine ports
	// created by CAPO are untagged and are not found by the OpenStack
	// destroyer's tag-based port cleanup. Remove them by name instead. This
	// runs after the OpenStack destroyer so that the servers are gone and the
	// ports are no longer attached.
	if err := deleteUntaggedPorts(context.TODO(), openstackMetadata.Cloud, o.Metadata.InfraID, o.Logger); err != nil {
		return quota, err
	}

	return quota, nil
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
