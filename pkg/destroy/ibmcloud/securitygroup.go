package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const securityGroupTypeName = "security group"

// listSecurityGroups lists security groups in the vpc
func (o *ClusterUninstaller) listSecurityGroups() (cloudResources, error) {
	o.Logger.Debugf("Listing security groups")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListSecurityGroupsOptions()
	resources, _, err := o.vpcSvc.ListSecurityGroupsWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list security groups")
	}

	result := []cloudResource{}
	for _, securityGroup := range resources.SecurityGroups {
		if strings.Contains(*securityGroup.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *securityGroup.ID,
				name:     *securityGroup.Name,
				status:   "",
				typeName: securityGroupTypeName,
				id:       *securityGroup.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteSecurityGroup(item cloudResource) error {
	o.Logger.Debugf("Deleting security group %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteSecurityGroupOptions(item.id)
	details, err := o.vpcSvc.DeleteSecurityGroupWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted security group %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete security group %s", item.name)
	}

	return nil
}

// destroySecurityGroups removes all security group resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySecurityGroups() error {
	if o.UserProvidedVPC == "" {
		o.Logger.Info("Skipping deletion of security groups with generated VPC")
		return nil
	}

	found, err := o.listSecurityGroups()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(securityGroupTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted security group %q", item.name)
			continue
		}
		err = o.deleteSecurityGroup(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(securityGroupTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
