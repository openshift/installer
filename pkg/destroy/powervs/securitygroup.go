package powervs

import (
	"context"
	"fmt"
	"math"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const securityGroupTypeName = "security group"

// listSecurityGroups lists security groups in the vpc.
func (o *ClusterUninstaller) listSecurityGroups() (cloudResources, error) {
	var (
		sgIDs         []string
		sgID          string
		ctx           context.Context
		cancel        context.CancelFunc
		result        = make([]cloudResource, 0, 1)
		options       *vpcv1.GetSecurityGroupOptions
		securityGroup *vpcv1.SecurityGroup
		response      *core.DetailedResponse
		err           error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		sgIDs, err = o.listByTag(TagTypeSecurityGroup)
	} else {
		// Otherwise list will list by name matching.
		sgIDs, err = o.listSecurityGroupsByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, sgID = range sgIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listSecurityGroups: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetSecurityGroupOptions(sgID)

		securityGroup, response, err = o.vpcSvc.GetSecurityGroupWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The security group could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get security group (%s): err = %w, response = %v", sgID, err, response)
		}

		result = append(result, cloudResource{
			key:      *securityGroup.ID,
			name:     *securityGroup.Name,
			status:   "",
			typeName: securityGroupTypeName,
			id:       *securityGroup.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listSecurityGroupsByName lists security groups in the vpc.
func (o *ClusterUninstaller) listSecurityGroupsByName() ([]string, error) {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		options       *vpcv1.ListSecurityGroupsOptions
		sgCollection  *vpcv1.SecurityGroupCollection
		response      *core.DetailedResponse
		foundOne      = false
		result        = make([]string, 0, 1)
		securityGroup vpcv1.SecurityGroup
		err           error
	)

	o.Logger.Debugf("Listing security groups by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listSecurityGroupsByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options = o.vpcSvc.NewListSecurityGroupsOptions()
	options.SetResourceGroupID(o.resourceGroupID)

	sgCollection, response, err = o.vpcSvc.ListSecurityGroupsWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list security groups: response = %v, err = %w", response, err)
	}

	for _, securityGroup = range sgCollection.SecurityGroups {
		if strings.Contains(*securityGroup.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listSecurityGroupsByName: FOUND: %s, %s", *securityGroup.ID, *securityGroup.Name)
			result = append(result, *securityGroup.ID)
		}
	}
	if !foundOne {
		o.Logger.Debugf("listSecurityGroupsByName: NO matching security group against: %s", o.InfraID)
		for _, securityGroup = range sgCollection.SecurityGroups {
			o.Logger.Debugf("listSecurityGroupsByName: security group: %s", *securityGroup.Name)
		}
	}

	return result, nil
}

// deleteSecurityGroup deletes the specified security group.
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

	if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Security Group %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusInternalServerError {
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
