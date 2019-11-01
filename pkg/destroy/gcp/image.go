package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listImages() ([]string, error) {
	o.Logger.Debugf("Listing images")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Images.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(imageList *compute.ImageList) error {
		for _, image := range imageList.Items {
			result = append(result, image.Name)
			o.Logger.Debugf("Found image %s\n", image.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch images")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteImage(name string) error {
	o.Logger.Debugf("Deleting image %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Images.Delete(o.ProjectID, name).Context(ctx).RequestId(o.requestID("image", name)).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("image", name)
		return errors.Wrapf(err, "failed to delete image %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("image", name)
		return errors.Errorf("failed to delete image %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(images))
	for _, image := range images {
		found = append(found, image)
		err := o.deleteImage(image)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedImages := o.setPendingItems("image", found)
	for _, item := range deletedImages {
		o.Logger.Infof("Deleted image %s", item)
	}
	return aggregateError(errs, len(found))
}
