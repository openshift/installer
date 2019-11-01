package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listRouters() ([]cloudResource, error) {
	o.Logger.Debug("Listing routers")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Routers.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.RouterList) error {
		for _, router := range list.Items {
			o.Logger.Debugf("Found router: %s", router.Name)
			result = append(result, cloudResource{
				key:      router.Name,
				name:     router.Name,
				typeName: "router",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list routers")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRouter(item cloudResource) error {
	o.Logger.Debugf("Deleting router %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Routers.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete router %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete router %s with error: %s", item.name, operationErrorMessage(op))
	}
	return nil
}

// destroyRouters removes all router resources that have a name prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyRouters() error {
	routers, err := o.listRouters()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, router := range routers {
		found.insert(router)
		err := o.deleteRouter(router)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("router", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted router %s", item.name)
	}
	return aggregateError(errs, len(found))
}
