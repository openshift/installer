package gcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	googleoauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	gcpsession "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/gather"
	"github.com/openshift/installer/pkg/gather/providers"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// Gather holds options for resources we want to gather.
type Gather struct {
	credentials     *googleoauth.Credentials
	clusterName     string
	clusterID       string
	infraID         string
	logger          logrus.FieldLogger
	serialLogBundle string
	bootstrap       string
	masters         []string
	directory       string
	endpoint        *gcptypes.PSCEndpoint
}

// New returns a GCP Gather from ClusterMetadata.
func New(logger logrus.FieldLogger, serialLogBundle string, bootstrap string, masters []string, metadata *types.ClusterMetadata) (providers.Gather, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	session, err := gcpsession.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	return &Gather{
		credentials:     session.Credentials,
		clusterName:     metadata.ClusterName,
		clusterID:       metadata.ClusterID,
		infraID:         metadata.InfraID,
		logger:          logger,
		serialLogBundle: serialLogBundle,
		bootstrap:       bootstrap,
		masters:         masters,
		directory:       filepath.Dir(serialLogBundle),
		endpoint:        metadata.GCP.Endpoint,
	}, nil
}

// Run is the entrypoint to start the gather process.
func (g *Gather) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	opts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(g.endpoint) {
		opts = append(opts, gcpsession.CreateEndpointOption(g.endpoint.Name, gcpsession.ServiceNameGCPCompute))
	}
	svc, err := gcpsession.GetComputeService(ctx, opts...)
	if err != nil {
		return err
	}
	isvc := compute.NewInstancesService(svc)

	var files []string
	var errs []error

	serialLogBundleDir := strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz")
	filePathDir := filepath.Join(g.directory, serialLogBundleDir)
	err = os.MkdirAll(filePathDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	req := svc.Instances.AggregatedList(g.credentials.ProjectID).Filter(fmt.Sprintf("name = %s-*", g.infraID))
	err = req.Pages(ctx, func(list *compute.InstanceAggregatedList) error {
		for _, aggListItem := range list.Items {
			for _, instance := range aggListItem.Instances {
				filename := filepath.Join(filePathDir, fmt.Sprintf("%s-serial.log", instance.Name))

				serialOutput, err := isvc.GetSerialPortOutput(g.credentials.ProjectID, filepath.Base(instance.Zone), instance.Name).Port(1).Do()
				if err != nil {
					errs = append(errs, err)
					continue
				}

				file, err := os.Create(filename)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				defer file.Close()

				_, err = file.Write([]byte(serialOutput.Contents))
				if err != nil {
					errs = append(errs, err)
					continue
				}

				files = append(files, filename)
			}
		}
		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	if len(files) > 0 {
		err := gather.CreateArchive(files, g.serialLogBundle)
		if err != nil {
			g.logger.Debugf("failed to create archive: %s", err.Error())
		}
	}

	err = gather.DeleteArchiveDirectory(filePathDir)
	if err != nil {
		g.logger.Debugf("failed to remove archive directory: %v", err)
	}

	return utilerrors.NewAggregate(errs)
}
