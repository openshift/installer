package gcp

import (
	"github.com/pkg/errors"

	storage "google.golang.org/api/storage/v1"
)

func (o *ClusterUninstaller) listStorageBuckets() ([]string, error) {
	o.Logger.Debug("Listing storage buckets")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.storageSvc.Buckets.List(o.ProjectID).Fields("items(name),nextPageToken").Prefix(o.ClusterID + "-")
	result := []string{}
	err := req.Pages(ctx, func(buckets *storage.Buckets) error {
		for _, bucket := range buckets.Items {
			o.Logger.Debugf("Found storage bucket %s", bucket.Name)
			result = append(result, bucket.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object storage buckets")
	}
	return result, nil
}

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

func (o *ClusterUninstaller) deleteStorageBucket(bucket string) error {
	o.Logger.Debugf("Deleting storate bucket %s", bucket)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := o.storageSvc.Buckets.Delete(bucket).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete bucket %s", bucket)
	}
	if err == nil {
		o.Logger.Infof("Deleted storage bucket %s", bucket)
	}
	return nil
}

// destroyObjectStorage finds all gcs buckets that have a name prefixed with
// the cluster's infra ID. It then removes all objects in each bucket and deletes it.
func (o *ClusterUninstaller) destroyObjectStorage() error {
	buckets, err := o.listStorageBuckets()
	if err != nil {
		return err
	}
	errs := []error{}
	for _, bucket := range buckets {
		objects, err := o.listBucketObjects(bucket)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, object := range objects {
			err = o.deleteStorageObject(bucket, object)
			if err != nil {
				errs = append(errs, err)
			}
		}
		err = o.deleteStorageBucket(bucket)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return aggregateError(errs, len(buckets))
}
