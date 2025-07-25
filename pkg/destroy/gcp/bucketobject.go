package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

const (
	bucketObjectResourceName = "bucketobject"
)

func (o *ClusterUninstaller) listBucketObjects(ctx context.Context, bucket cloudResource) ([]cloudResource, error) {
	o.Logger.Debugf("Listing objects for storage bucket %s", bucket.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}

	it := o.storageSvc.Bucket(bucket.name).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list bucket objects: %w", err)
		}
		o.Logger.Debugf("Found storage object %s/%s", bucket.name, attrs.Name)
		result = append(result, cloudResource{
			key:      attrs.Name,
			name:     attrs.Name,
			typeName: bucketObjectResourceName,
		})
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBucketObject(ctx context.Context, bucket cloudResource, item cloudResource) error {
	o.Logger.Debugf("Deleting storage object %s/%s", bucket.name, item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	err := o.storageSvc.Bucket(bucket.name).Object(item.name).Delete(ctx)
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete bucket object %s/%s", bucket.name, item.name)
	}
	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted bucket object %s", item.name)
	return nil
}
