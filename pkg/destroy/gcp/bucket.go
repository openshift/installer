package gcp

import (
	"regexp"

	"github.com/pkg/errors"

	"google.golang.org/api/googleapi"
	storage "google.golang.org/api/storage/v1"
)

var (
	// multiDashes is a regexp matching multiple dashes in a sequence.
	multiDashes = regexp.MustCompile(`-{2,}`)
)

func (o *ClusterUninstaller) listBuckets() ([]cloudResource, error) {
	return o.listBucketsWithFilter("items(name),nextPageToken", o.ClusterID+"-", nil)
}

// listBucketsWithFilter lists buckets in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a prefix string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listBucketsWithFilter(fields string, prefix string, filterFunc func(*storage.Bucket) bool) ([]cloudResource, error) {
	o.Logger.Debug("Listing storage buckets")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []cloudResource{}
	req := o.storageSvc.Buckets.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(prefix) > 0 {
		prefix = multiDashes.ReplaceAllString(prefix, "-")
		req = req.Prefix(prefix)
	}
	err := req.Pages(ctx, func(list *storage.Buckets) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found bucket: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "bucket",
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object storage buckets")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBucket(item cloudResource) error {
	o.Logger.Debugf("Deleting storate bucket %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := o.storageSvc.Buckets.Delete(item.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete bucket %s", item.name)
	}
	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted bucket %s", item.name)
	return nil
}

// destroyBuckets finds all gcs buckets that have a name prefixed
// with the cluster's infra ID. It then removes all the objects in each bucket and deletes it.
func (o *ClusterUninstaller) destroyBuckets() error {
	found, err := o.listBuckets()
	if err != nil {
		return err
	}
	items := o.insertPendingItems("bucket", found)
	for _, item := range items {
		foundObjects, err := o.listBucketObjects(item)
		if err != nil {
			return err
		}
		objects := o.insertPendingItems("bucketobject", foundObjects)
		for _, object := range objects {
			err = o.deleteBucketObject(item, object)
			if err != nil {
				o.errorTracker.suppressWarning(object.key, err, o.Logger)
			}
		}
		err = o.deleteBucket(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("bucket"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
