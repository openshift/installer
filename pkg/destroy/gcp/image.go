package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listImages() ([]cloudResource, error) {
	o.Logger.Debugf("Listing images")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Images.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(imageList *compute.ImageList) error {
		for _, image := range imageList.Items {
			result = append(result, cloudResource{
				key:      image.Name,
				name:     image.Name,
				typeName: "image",
			})
			o.Logger.Debugf("Found image %s\n", image.Name)
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
	return nil
}

// destroyImages removes all image resources with a name prefixed by the
// cluster's infra ID
func (o *ClusterUninstaller) destroyImages() error {
	images, err := o.listImages()
	if err != nil {
		return err
	}
	errs := []error{}
	found := cloudResources{}
	for _, image := range images {
		found.insert(image)
		err := o.deleteImage(image)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedImages := o.setPendingItems("image", found)
	for _, item := range deletedImages {
		o.Logger.Infof("Deleted image %s", item.name)
	}
	return aggregateError(errs, len(found))
}
