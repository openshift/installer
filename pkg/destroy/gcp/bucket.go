package gcp

import (
	"github.com/pkg/errors"

	storage "google.golang.org/api/storage/v1"
)

func (o *ClusterUninstaller) listStorageBuckets() ([]cloudResource, error) {
	o.Logger.Debug("Listing storage buckets")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.storageSvc.Buckets.List(o.ProjectID).Fields("items(name),nextPageToken").Prefix(o.ClusterID + "-")
	result := []cloudResource{}
	err := req.Pages(ctx, func(buckets *storage.Buckets) error {
		for _, bucket := range buckets.Items {
			o.Logger.Debugf("Found storage bucket %s", bucket.Name)
			result = append(result, cloudResource{
				key:      bucket.Name,
				name:     bucket.Name,
				typeName: "storgebucket",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object storage buckets")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteStorageBucket(item cloudResource) error {
	o.Logger.Debugf("Deleting storate bucket %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := o.storageSvc.Buckets.Delete(item.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete bucket %s", item.name)
	}
	if err == nil {
		o.Logger.Infof("Deleted storage bucket %s", item.name)
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
