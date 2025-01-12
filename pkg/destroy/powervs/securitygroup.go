package powervs

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const securityGroupTypeName = "security group"

// listSecurityGroups lists security groups in the vpc.
func (o *ClusterUninstaller) listSecurityGroups() (cloudResources, error) {
	o.Logger.Debugf("Listing security groups")

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listSecurityGroups: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListSecurityGroupsOptions()
	options.SetResourceGroupID(o.resourceGroupID)

	resources, _, err := o.vpcSvc.ListSecurityGroupsWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list security groups: %w", err)
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

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteSecurityGroup: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getOptions = o.vpcSvc.NewGetSecurityGroupOptions(item.id)
	_, response, err = o.vpcSvc.GetSecurityGroup(getOptions)

	if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Security Group %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
		o.Logger.Infof("deleteSecurityGroup: internal server error")
		return nil
	}

	deleteOptions := o.vpcSvc.NewDeleteSecurityGroupOptions(item.id)

	_, err = o.vpcSvc.DeleteSecurityGroupWithContext(ctx, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete security group %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Security Group %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroySecurityGroups removes all security group resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySecurityGroups() error {
	firstPassList, err := o.listSecurityGroups()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(securityGroupTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroySecurityGroups: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteSecurityGroup(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroySecurityGroups: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(securityGroupTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyServiceInstances: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroySecurityGroups: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listSecurityGroups()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroySecurityGroups: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroySecurityGroups: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
