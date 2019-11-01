package gcp

import (
	"github.com/pkg/errors"

	storage "google.golang.org/api/storage/v1"
)

func (o *ClusterUninstaller) listBucketObjects(bucket string) ([]string, error) {
	o.Logger.Debugf("Listing objects for storage bucket %s", bucket)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []string{}
	req := o.storageSvc.Objects.List(bucket).Fields("items(name),nextPageToken")
	err := req.Pages(ctx, func(objects *storage.Objects) error {
		for _, object := range objects.Items {
			o.Logger.Debugf("Found storage object %s/%s", bucket, object.Name)
			result = append(result, object.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch objects for bucket %s", bucket)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteStorageObject(bucket, object string) error {
	o.Logger.Debugf("Deleting storate object %s/%s", bucket, object)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := o.storageSvc.Objects.Delete(bucket, object).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete bucket object %s/%s", bucket, object)
	}
	if err == nil {
		o.Logger.Infof("Deleted storage object %s/%s", bucket, object)
	}
	return nil
}
