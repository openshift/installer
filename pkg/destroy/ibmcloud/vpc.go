package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const vpcTypeName = "vpc"

// listVPCs lists subnets in the vpc
func (o *ClusterUninstaller) listVPCs() (cloudResources, error) {
	o.Logger.Debugf("Listing VPCs")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListVpcsOptions()
	resources, _, err := o.vpcSvc.ListVpcsWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list vpcs")
	}

	result := []cloudResource{}
	for _, vpc := range resources.Vpcs {
		if strings.Contains(*vpc.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *vpc.ID,
				name:     *vpc.Name,
				status:   *vpc.Status,
				typeName: vpcTypeName,
				id:       *vpc.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteVPC(item cloudResource) error {
	if item.status == vpcv1.VPCStatusDeletingConst {
		o.Logger.Debugf("Waiting for VPC %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting VPC %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteVPCOptions(item.id)
	details, err := o.vpcSvc.DeleteVPCWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted VPC %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete VOC %s", item.name)
	}

	return nil
}

// listVPCs removes all VPC resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyVPCs() error {
	if o.UserProvidedVPC != "" {
		o.Logger.Infof("Skipping deletion of user-provided VPC %q", o.UserProvidedVPC)
		return nil
	}

	found, err := o.listVPCs()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(vpcTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted VPC %q", item.name)
			continue
		}
		err := o.deleteVPC(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(vpcTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
