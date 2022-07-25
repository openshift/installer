package powervs

import (
	"net/http"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const securityGroupTypeName = "security group"

// listSecurityGroups lists security groups in the vpc.
func (o *ClusterUninstaller) listSecurityGroups() (cloudResources, error) {
	o.Logger.Debugf("Listing security groups")

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listSecurityGroups: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListSecurityGroupsOptions()
	resources, _, err := o.vpcSvc.ListSecurityGroupsWithContext(ctx, options)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list security groups")
	}

	var foundOne = false

	result := []cloudResource{}
	for _, securityGroup := range resources.SecurityGroups {
		if strings.Contains(*securityGroup.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listSecurityGroups: FOUND: %s, %s", *securityGroup.ID, *securityGroup.Name)
			result = append(result, cloudResource{
				key:      *securityGroup.ID,
				name:     *securityGroup.Name,
				status:   "",
				typeName: securityGroupTypeName,
				id:       *securityGroup.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listSecurityGroups: NO matching security group against: %s", o.InfraID)
		for _, securityGroup := range resources.SecurityGroups {
			o.Logger.Debugf("listSecurityGroups: security group: %s", *securityGroup.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteSecurityGroup(item cloudResource) error {
	var getOptions *vpcv1.GetSecurityGroupOptions
	var response *core.DetailedResponse
	var err error

	getOptions = o.vpcSvc.NewGetSecurityGroupOptions(item.id)
	_, response, err = o.vpcSvc.GetSecurityGroup(getOptions)

	if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted security group %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
		o.Logger.Infof("deleteSecurityGroup: internal server error")
		return nil
	}

	o.Logger.Debugf("Deleting security group %q", item.name)

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteSecurityGroup: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	deleteOptions := o.vpcSvc.NewDeleteSecurityGroupOptions(item.id)
	_, err = o.vpcSvc.DeleteSecurityGroupWithContext(ctx, deleteOptions)

	if err != nil {
		return errors.Wrapf(err, "failed to delete security group %s", item.name)
	}

	return nil
}

// destroySecurityGroups removes all security group resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySecurityGroups() error {
	found, err := o.listSecurityGroups()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(securityGroupTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroySecurityGroups: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

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

		items = o.getPendingItems(securityGroupTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(securityGroupTypeName); len(items) > 0 {
		return errors.Errorf("destroySecurityGroups: %d undeleted items pending", len(items))
	}
	return nil
}
