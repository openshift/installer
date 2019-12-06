package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (o *ClusterUninstaller) listImages() ([]cloudResource, error) {
	return o.listImagesWithFilter("items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listImagesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listImagesWithFilter(fields string, filter string, filterFunc func(*compute.Image) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing images")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Images.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.ImageList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found image: %s\n", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "image",
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

func (o *ClusterUninstaller) deleteImage(item cloudResource) error {
	o.Logger.Debugf("Deleting image %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Images.Delete(o.ProjectID, item.name).Context(ctx).RequestId(o.requestID(item.typeName, item.name)).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete image %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete image %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted image %s", item.name)
	}
	return nil
}

// destroyImages removes all image resources with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyImages() error {
	found, err := o.listImages()
	if err != nil {
		return err
	}
	items := o.insertPendingItems("image", found)
	errs := []error{}
	for _, item := range items {
		err := o.deleteImage(item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	items = o.getPendingItems("image")
	return aggregateError(errs, len(items))
}
