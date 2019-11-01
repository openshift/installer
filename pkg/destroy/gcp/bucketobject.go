package gcp

import (
	"github.com/pkg/errors"

	storage "google.golang.org/api/storage/v1"
)

func (o *ClusterUninstaller) listBucketObjects(bucket cloudResource) ([]cloudResource, error) {
	o.Logger.Debugf("Listing objects for storage bucket %s", bucket.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []cloudResource{}
	req := o.storageSvc.Objects.List(bucket.name).Fields("items(name),nextPageToken")
	err := req.Pages(ctx, func(objects *storage.Objects) error {
		for _, object := range objects.Items {
			o.Logger.Debugf("Found storage object %s/%s", bucket.name, object.Name)
			result = append(result, cloudResource{
				key:      object.Name,
				name:     object.Name,
				typeName: "storageobject",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch objects for bucket %s", bucket.name)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteStorageObject(bucket cloudResource, item cloudResource) error {
	o.Logger.Debugf("Deleting storate object %s/%s", bucket.name, item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := o.storageSvc.Objects.Delete(bucket.name, item.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete bucket object %s/%s", bucket.name, item.name)
	}
	if err == nil {
		o.Logger.Infof("Deleted storage object %s/%s", bucket.name, item.name)
	}
	return nil
}
