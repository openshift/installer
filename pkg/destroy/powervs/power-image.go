package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"k8s.io/apimachinery/pkg/util/wait"
)

const imageTypeName = "image"

// listImages lists images in the vpc.
func (o *ClusterUninstaller) listImages() (cloudResources, error) {
	o.Logger.Debugf("Listing images")

	if o.imageClient == nil {
		o.Logger.Infof("Skipping deleting images because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listImages: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	images, err := o.imageClient.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, image := range images.Images {
		if strings.Contains(*image.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listImages: FOUND: %s, %s, %s", *image.ImageID, *image.Name, *image.State)
			result = append(result, cloudResource{
				key:      *image.ImageID,
				name:     *image.Name,
				status:   *image.State,
				typeName: imageTypeName,
				id:       *image.ImageID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listImages: NO matching image against: %s", o.InfraID)
		for _, image := range images.Images {
			o.Logger.Debugf("listImages: image: %s", *image.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteImage(item cloudResource) error {
	var img *models.Image
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteImage: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	img, err = o.imageClient.Get(item.id)
	if err != nil {
		o.Logger.Debugf("listImages: deleteImage: image %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Image %q", item.name)
		return nil
	}

	if !strings.EqualFold(img.State, "active") {
		o.Logger.Debugf("Waiting for image %q to delete", item.name)
		return nil
	}

	err = o.imageClient.Delete(item.id)
	if err != nil {
		return fmt.Errorf("failed to delete image %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Image %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyImages removes all image resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyImages() error {
	firstPassList, err := o.listImages()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(imageTypeName, firstPassList.list())
	for _, item := range items {
		o.Logger.Debugf("destroyImages: firstPassList: %v / %v", item.name, item.id)
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyImages: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteImage(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyImages: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(imageTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyImages: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyImages: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listImages()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyImages: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyImages: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
