package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

type nameAndZone struct {
	name string
	zone string
}

func (n nameAndZone) String() string {
	return fmt.Sprintf("Name: %s, Zone: %s", n.name, n.zone)
}

// getZoneName extracts a zone name from a zone URL of the form:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a
// where the compute service's basepath is:
// https://www.googleapis.com/compute/v1/projects/
// Trimming the URL, leaves a string like: project-id/zones/us-central1-a
func (o *ClusterUninstaller) getZoneName(zoneURL string) string {
	path := strings.TrimLeft(zoneURL, o.computeSvc.BasePath)
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

// getInstanceNameAndZone extracts an instance and zone name from an instance URL in the form:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a/instances/instance-name
// After trimming the service's base path, you get:
// project-id/zones/us-central1-a/instances/instance-name
func (o *ClusterUninstaller) getInstanceNameAndZone(instanceURL string) (string, string) {
	path := strings.TrimLeft(instanceURL, o.computeSvc.BasePath)
	parts := strings.Split(path, "/")
	if len(parts) >= 5 {
		return parts[4], parts[2]
	}
	return "", ""
}

func (o *ClusterUninstaller) listInstanceGroupsWithFilter(filter string) ([]nameAndZone, error) {
	o.Logger.Debugf("Listing instance groups")
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.AggregatedList(o.ProjectID).Fields("items/*/instanceGroups(name,zone),nextPageToken").Filter(filter)
	err := req.Pages(ctx, func(list *compute.InstanceGroupAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, ig := range scopedList.InstanceGroups {
				zoneName := o.getZoneName(ig.Zone)
				result = append(result, nameAndZone{
					name: ig.Name,
					zone: zoneName,
				})
				o.Logger.Debugf("Found instance group %s in zone %s", ig.Name, zoneName)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch compute instance groups")
	}
	return result, nil
}

func (o *ClusterUninstaller) listInstanceGroupInstances(ig nameAndZone) ([]nameAndZone, error) {
	o.Logger.Debugf("Listing instance group instances for %v", ig)
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.ListInstances(o.ProjectID, ig.zone, ig.name, &compute.InstanceGroupsListInstancesRequest{}).Fields("items(instance),nextPageToken")
	err := req.Pages(ctx, func(list *compute.InstanceGroupsListInstances) error {
		for _, item := range list.Items {
			name, zone := o.getInstanceNameAndZone(item.Instance)
			if len(name) == 0 {
				continue
			}
			result = append(result, nameAndZone{
				name: name,
				zone: zone,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch instance group instances")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteInstanceGroup(ig nameAndZone) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting instance group %s in zone %s", ig.name, ig.zone)
	op, err := o.computeSvc.InstanceGroups.Delete(o.ProjectID, ig.zone, ig.name).RequestId(o.requestID("instancegroup", ig.zone, ig.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("instancegroup", ig.zone, ig.name)
		return errors.Wrapf(err, "failed to delete instance group %s in zone %s", ig.name, ig.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("instancegroup", ig.zone, ig.name)
		return errors.Errorf("failed to delete instance group %s in zone %s with error: %s", ig.name, ig.zone, operationErrorMessage(op))
	}
	return nil
}

// destroyInstanceGroups searches for instance groups across all zones that have a name that starts with
// the infra ID prefix, and then deletes each instance found.
func (o *ClusterUninstaller) destroyInstanceGroups() error {
	groups, err := o.listInstanceGroupsWithFilter(o.clusterIDFilter())
	if err != nil {
		return err
	}
	errs := []error{}
	found := make([]string, 0, len(groups))
	for _, group := range groups {
		found = append(found, fmt.Sprintf("%s/%s", group.zone, group.name))
		err := o.deleteInstanceGroup(group)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("computeinstancegroup", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted instance group %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listComputeInstances() ([]nameAndZone, error) {
	o.Logger.Debugf("Listing compute instances")
	result := []nameAndZone{}
	req := o.computeSvc.Instances.AggregatedList(o.ProjectID).Filter(o.clusterIDFilter()).Fields("items/*/instances(name,zone,status),nextPageToken")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := req.Pages(ctx, func(list *compute.InstanceAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, instance := range scopedList.Instances {
				zoneName := o.getZoneName(instance.Zone)
				result = append(result, nameAndZone{
					name: instance.Name,
					zone: zoneName,
				})
				o.Logger.Debugf("Found instance %s in zone %s, status %s", instance.Name, zoneName, instance.Status)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch compute instances")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteComputeInstance(instance nameAndZone) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting compute instance %s in zone %s", instance.name, instance.zone)
	op, err := o.computeSvc.Instances.Delete(o.ProjectID, instance.zone, instance.name).RequestId(o.requestID("instance", instance.zone, instance.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("instance", instance.zone, instance.name)
		return errors.Wrapf(err, "failed to delete instance %s in zone %s", instance.name, instance.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("instance", instance.zone, instance.name)
		return errors.Errorf("failed to delete instance %s in zone %s with error: %s", instance.name, instance.zone, operationErrorMessage(op))
	}
	return nil
}

// destroyComputeInstances searches for instances across all zones that have a name that starts with
// the infra ID prefix and are not yet terminated. It then deletes each instance found.
func (o *ClusterUninstaller) destroyComputeInstances() error {
	instances, err := o.listComputeInstances()
	if err != nil {
		return err
	}
	errs := []error{}
	found := make([]string, 0, len(instances))
	for _, instance := range instances {
		found = append(found, fmt.Sprintf("%s/%s", instance.zone, instance.name))
		err := o.deleteComputeInstance(instance)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("computeinstance", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted instance %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listImages() ([]string, error) {
	o.Logger.Debugf("Listing images")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Images.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(imageList *compute.ImageList) error {
		for _, image := range imageList.Items {
			result = append(result, image.Name)
			o.Logger.Debugf("Found image %s\n", image.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch images")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteImage(name string) error {
	o.Logger.Debugf("Deleting image %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Images.Delete(o.ProjectID, name).Context(ctx).RequestId(o.requestID("image", name)).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("image", name)
		return errors.Wrapf(err, "failed to delete image %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("image", name)
		return errors.Errorf("failed to delete image %s with error: %s", name, operationErrorMessage(op))
	}
	return nil
}

// destroyImages removes all image resources with a name prefixed by the
// cluster's infra ID
func (o *ClusterUninstaller) destroyImages() error {
	images, err := o.listImages()
	if err != nil {
		return err
	}
	errs := []error{}
	found := make([]string, 0, len(images))
	for _, image := range images {
		found = append(found, image)
		err := o.deleteImage(image)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedImages := o.setPendingItems("image", found)
	for _, item := range deletedImages {
		o.Logger.Infof("Deleted image %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listDisks() ([]nameAndZone, error) {
	o.Logger.Debug("Listing disks")
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Disks.AggregatedList(o.ProjectID).Fields("items/*/disks(name,zone),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(aggregatedList *compute.DiskAggregatedList) error {
		for _, scopedList := range aggregatedList.Items {
			for _, disk := range scopedList.Disks {
				zone := o.getZoneName(disk.Zone)
				result = append(result, nameAndZone{name: disk.Name, zone: zone})
				o.Logger.Debugf("Found disk %s in zone %s", disk.Name, zone)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch disks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteDisk(disk nameAndZone) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting disk %s in zone %s", disk.name, disk.zone)
	op, err := o.computeSvc.Disks.Delete(o.ProjectID, disk.zone, disk.name).RequestId(o.requestID("disk", disk.zone, disk.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("disk", disk.zone, disk.name)
		return errors.Wrapf(err, "failed to delete disk %s in zone %s", disk.name, disk.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("disk", disk.zone, disk.name)
		return errors.Errorf("failed to delete disk %s in zone %s with error: %s", disk.name, disk.zone, operationErrorMessage(op))
	}
	return nil
}

// destroyDisks searches for disks across all zones that have a name that starts with
// the infra ID prefix. It then deletes each disk found.
func (o *ClusterUninstaller) destroyDisks() error {
	disks, err := o.listDisks()
	if err != nil {
		return err
	}
	errs := []error{}
	found := make([]string, 0, len(disks))
	for _, disk := range disks {
		found = append(found, fmt.Sprintf("%s/%s", disk.zone, disk.name))
		err := o.deleteDisk(disk)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("disk", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted disk %s", item)
	}
	return aggregateError(errs, len(found))
}
