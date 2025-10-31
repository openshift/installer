package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	imageResourceName = "image"
)

func (o *ClusterUninstaller) listImages(ctx context.Context) ([]cloudResource, error) {
	return o.listImagesWithFilter(ctx, "items(name,labels),nextPageToken", func(item *compute.Image) bool {
		return o.isClusterResource(item.Name) && !o.isSharedResource(item.Labels)
	})
}

// listImagesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listImagesWithFilter(ctx context.Context, fields string, filterFunc func(item *compute.Image) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing images")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Images.List(o.ProjectID).Fields(googleapi.Field(fields))
	err := req.Pages(ctx, func(list *compute.ImageList) error {
		for _, item := range list.Items {
			if filterFunc(item) {
				o.Logger.Debugf("Found image: %s\n", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: imageResourceName,
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "images",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch images")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteImage(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting image %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Images.Delete(o.ProjectID, item.name).Context(ctx).RequestId(o.requestID(item.typeName, item.name)).Do()
	item.scope = global
	return o.handleOperation(ctx, op, err, item, "image")
}

// destroyImages removes all image resources with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyImages(ctx context.Context) error {
	found, err := o.listImages(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(imageResourceName, found)
	for _, item := range items {
		err := o.deleteImage(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(imageResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
