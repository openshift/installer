package powervs

import (
	"strings"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
)

const imageTypeName = "image"

// listImages lists images in the vpc.
func (o *ClusterUninstaller) listImages() (cloudResources, error) {
	o.Logger.Debugf("Listing images")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listImages: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	images, err := o.imageClient.GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list images")
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

	img, err = o.imageClient.Get(item.id)
	if err != nil {
		o.Logger.Debugf("listImages: deleteImage: image %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted image %q", item.name)
		return nil
	}

	if !strings.EqualFold(img.State, "active") {
		o.Logger.Debugf("Waiting for image %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting image %q", item.name)

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteImage: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	err = o.imageClient.Delete(item.id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete image %s", item.name)
	}

	return nil
}

// destroyImages removes all image resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyImages() error {
	found, err := o.listImages()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(imageTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyImages: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted image %q", item.name)
				continue
			}
			err := o.deleteImage(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(imageTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(imageTypeName); len(items) > 0 {
		return errors.Errorf("destroyImages: %d undeleted items pending", len(items))
	}
	return nil
}
