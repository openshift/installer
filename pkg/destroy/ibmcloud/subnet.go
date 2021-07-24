package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const subnetTypeName = "subnet"

// listSubnets lists subnets in the vpc
func (o *ClusterUninstaller) listSubnets() (cloudResources, error) {
	o.Logger.Debugf("Listing subnets")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListSubnetsOptions()
	resources, _, err := o.vpcSvc.ListSubnetsWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list subnets")
	}

	result := []cloudResource{}
	for _, subnet := range resources.Subnets {
		if strings.Contains(*subnet.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *subnet.ID,
				name:     *subnet.Name,
				status:   *subnet.Status,
				typeName: subnetTypeName,
				id:       *subnet.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteSubnet(item cloudResource) error {
	if item.status == vpcv1.SubnetStatusDeletingConst {
		o.Logger.Debugf("Waiting for subnet %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting subnet %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteSubnetOptions(item.id)
	details, err := o.vpcSvc.DeleteSubnetWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted subnet %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete subnet %s", item.name)
	}

	return nil
}

// destroySubnets removes all subnet resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySubnets() error {
	if len(o.UserProvidedSubnets) > 0 {
		o.Logger.Infof("Skipping deletion of user-provided subnets %v", o.UserProvidedSubnets)
		return nil
	}

	found, err := o.listSubnets()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(subnetTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted subnet %q", item.name)
			continue
		}
		err = o.deleteSubnet(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(subnetTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
