package gcp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	storage "google.golang.org/api/storage/v1"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
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
}

// New returns an AWS destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Logger:    logger,
		Region:    metadata.ClusterPlatformMetadata.GCP.Region,
		ProjectID: metadata.ClusterPlatformMetadata.GCP.ProjectID,
		ClusterID: metadata.InfraID,
		Context:   context.Background(),
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
		return errors.Wrap(err, "failed to create compute service")
	}

	o.dnsSvc, err = dns.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create dns service")
	}

	o.storageSvc, err = storage.NewService(ctx, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create storage service")
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
		{name: "Instance groups", destroy: o.destroyInstanceGroups},
		{name: "Service accounts", destroy: o.destroyServiceAccounts},
		{name: "Images", destroy: o.destroyImages},
		{name: "DNS", destroy: o.destroyDNS},
		{name: "Object storage", destroy: o.destroyObjectStorage},
		{name: "Routes", destroy: o.destroyRoutes},
		{name: "Firewalls", destroy: o.destroyFirewalls},
		{name: "Addresses", destroy: o.destroyAddresses},
		{name: "Target Pools", destroy: o.destroyTargetPools},
		{name: "Forwarding rules", destroy: o.destroyForwardingRules},
		{name: "Backend services", destroy: o.destroyBackendServices},
		{name: "Health checks", destroy: o.destroyHealthChecks},
		{name: "Cloud controller internal LBs", destroy: o.destroyCloudControllerInternalLBs},
		{name: "Cloud controller external LBs", destroy: o.destroyCloudControllerExternalLBs},
		{name: "Cloud routers", destroy: o.destroyRouters},
		{name: "Subnetworks", destroy: o.destroySubNetworks},
		{name: "Networks", destroy: o.destroyNetworks},
	}
	hasErr := false
	for _, f := range destroyFuncs {
		err := f.destroy()
		if err != nil {
			hasErr = true
			o.Logger.Errorf("%s: %v", f.name, err)
		}
	}
	return !hasErr, nil
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

func (o *ClusterUninstaller) contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(o.Context, defaultTimeout)
}
