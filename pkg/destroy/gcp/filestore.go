package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/file/v1"
)

func (o *ClusterUninstaller) filestoreParentPath() string {
	return fmt.Sprintf("projects/%s/locations/-", o.ProjectID)
}

func (o *ClusterUninstaller) clusterFilestoreLabelFilter() string {
	return fmt.Sprintf("labels.kubernetes-io-cluster-%s = \"owned\"", o.ClusterID)
}

func (o *ClusterUninstaller) listFilestores(ctx context.Context) ([]cloudResource, error) {
	o.Logger.Debug("Listing filestores")
	result := []cloudResource{}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	instSvc := file.NewProjectsLocationsInstancesService(o.fileSvc)
	lCall := instSvc.List(o.filestoreParentPath()).Context(ctx).Filter(o.clusterFilestoreLabelFilter())

	nextPageToken := "pageToken"
	for nextPageToken != "" {
		instances, err := lCall.Do()
		if err != nil {
			return nil, errors.Wrapf(err, "error retrieving filestore instances")
		}

		for _, activeInstance := range instances.Instances {
			if len(activeInstance.FileShares) == 0 {
				// skip multi-share instances
				o.Logger.Debug("Skipping multi-share filestore %s", activeInstance.Name)
				continue
			}
			o.Logger.Debug("Found filestore %s", activeInstance.Name)
			result = append(result, cloudResource{
				name:     activeInstance.Name,
				typeName: "filestore",
			})
		}

		nextPageToken = instances.NextPageToken
		lCall.PageToken(nextPageToken)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteFilestore(ctx context.Context, item cloudResource) error {
	o.Logger.Debug("Deleting filestore %s", item.name)

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	instSvc := file.NewProjectsLocationsInstancesService(o.fileSvc)
	_, err := instSvc.Delete(item.name).Context(ctx).Do()

	return err
}

func (o *ClusterUninstaller) destroyFilestores(ctx context.Context) error {
	found, err := o.listFilestores(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("filestore", found)
	for _, item := range items {
		err := o.deleteFilestore(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("filestore"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}

	return nil
}
