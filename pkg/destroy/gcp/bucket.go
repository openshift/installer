package gcp

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
	storage "google.golang.org/api/storage/v1"
)

var (
	// multiDashes is a regexp matching multiple dashes in a sequence.
	multiDashes = regexp.MustCompile(`-{2,}`)
)

const (
	bucketResourceName = "bucket"
)

func (o *ClusterUninstaller) listBuckets(ctx context.Context) ([]cloudResource, error) {
	return o.listBucketsWithFilter(ctx, "items(name),nextPageToken",
		func(itemName string) bool {
			prefix := multiDashes.ReplaceAllString(o.ClusterID+"-", "-")
			return strings.HasPrefix(itemName, prefix)
		})
}

// listBucketsWithFilter lists buckets in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a prefix string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listBucketsWithFilter(ctx context.Context, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debug("Listing storage buckets")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.storageSvc.Buckets.List(o.ProjectID).Fields(googleapi.Field(fields))

	err := req.Pages(ctx, func(list *storage.Buckets) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
				o.Logger.Debugf("Found bucket: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: bucketResourceName,
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch object storage buckets: %w", err)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBucket(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting storage bucket %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
func (o *ClusterUninstaller) destroyBuckets(ctx context.Context) error {
	found, err := o.listBuckets(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(bucketResourceName, found)
	for _, item := range items {
		foundObjects, err := o.listBucketObjects(ctx, item)
		if err != nil {
			return err
		}
		objects := o.insertPendingItems(bucketObjectResourceName, foundObjects)
		for _, object := range objects {
			err = o.deleteBucketObject(ctx, item, object)
			if err != nil {
				o.errorTracker.suppressWarning(object.key, err, o.Logger)
			}
		}
		err = o.deleteBucket(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(bucketResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
