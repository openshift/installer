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
		return parts[2], parts[4]
	}
	return "", ""
}

func (o *ClusterUninstaller) listInstanceGroupsWithFilter(filter string) ([]nameAndZone, error) {
	o.Logger.Debugf("Listing instance groups")
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.AggregatedList(o.ProjectID).Fields("items/*/instanceGroups(name,zone)").Filter(filter)
	err := req.Pages(ctx, func(list *compute.InstanceGroupAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, ig := range scopedList.InstanceGroups {
				zoneName := o.getZoneName(ig.Zone)
				result = append(result, nameAndZone{
					name: ig.Name,
					zone: zoneName,
				})
				o.Logger.Debugf("Found instance group %s in zone %s\n", ig.Name, zoneName)
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
	req := o.computeSvc.InstanceGroups.ListInstances(o.ProjectID, ig.zone, ig.name, &compute.InstanceGroupsListInstancesRequest{}).Fields("items(instance)")
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
	_, err := o.computeSvc.InstanceGroups.Delete(o.ProjectID, ig.zone, ig.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete instance group %s in zone %s", ig.name, ig.zone)
	}
	return nil
}

func (o *ClusterUninstaller) listComputeInstances() ([]nameAndZone, error) {
	o.Logger.Debugf("Listing compute instances")
	result := []nameAndZone{}
	filter := fmt.Sprintf("(status ne TERMINATED)(%s)", o.clusterIDFilter())
	req := o.computeSvc.Instances.AggregatedList(o.ProjectID).Filter(filter).Fields("items/*/instances(name,zone,status)")
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
	_, err := o.computeSvc.Instances.Delete(o.ProjectID, instance.zone, instance.name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete instance %s in zone %s", instance.name, instance.zone)
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
	for _, instance := range instances {
		err := o.deleteComputeInstance(instance)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listImages() ([]string, error) {
	o.Logger.Debugf("Listing images")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Images.List(o.ProjectID).Fields("items(name)").Filter(o.clusterIDFilter())
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
	_, err := o.computeSvc.Images.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete image %s", name)
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
	for _, image := range images {
		err := o.deleteImage(image)
		if err != nil {
			return err
		}
	}
	return nil
}
