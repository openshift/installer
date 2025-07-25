package gcp

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
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

	it := o.storageSvc.Buckets(ctx, o.ProjectID)
	for {
		bucketAttrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to fetch storage bucket: %w", err)
		}
		if filterFunc(bucketAttrs.Name) {
			o.Logger.Debugf("Found bucket: %s", bucketAttrs.Name)
			result = append(result, cloudResource{
				key:      bucketAttrs.Name,
				name:     bucketAttrs.Name,
				typeName: bucketResourceName,
			})
		}
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBucket(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting storage bucket %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	err := o.storageSvc.Bucket(item.name).Delete(ctx)
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
