package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/file/v1"

	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
)

func (o *ClusterUninstaller) filestoreParentPath() string {
	return fmt.Sprintf("projects/%s/locations/-", o.ProjectID)
}

func (o *ClusterUninstaller) clusterFilestoreLabelFilter() string {
	return fmt.Sprintf("labels.%s = \"owned\"", fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, o.ClusterID))
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
			return nil, fmt.Errorf("error retrieving filestore instances: %w", err)
		}

		for _, activeInstance := range instances.Instances {
			o.Logger.Debugf("Found filestore %s", activeInstance.Name)
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
	o.Logger.Debugf("Deleting filestore %s", item.name)

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	instSvc := file.NewProjectsLocationsInstancesService(o.fileSvc)
	_, err := instSvc.Delete(item.name).Context(ctx).Do()
	if err != nil && isForbidden(err) {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		return fmt.Errorf("insufficient permissions to delete filestore %s", item.name)
	}
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete filestore %s", item.name)
	}
	if err != nil && isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted filestore %s", item.name)
	}

	return nil
}

func (o *ClusterUninstaller) destroyFilestores(ctx context.Context) error {
	found, err := o.listFilestores(ctx)
	if err != nil {
		if isForbidden(err) {
			o.Logger.Warning("Skipping deletion of filestores: insufficient Filestore API permissions or API disabled")
			return nil
		}
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
		return fmt.Errorf("%d items pending", len(items))
	}

	return nil
}
