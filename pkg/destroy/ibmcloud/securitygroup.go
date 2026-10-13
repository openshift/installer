package ibmcloud

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const (
	securityGroupTypeName     = "security group"
	securityGroupRuleTypeName = "security group rule"
)

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

func (o *ClusterUninstaller) listSecurityGroupRules(securityGroupID string) (cloudResources, error) {
	o.Logger.Debugf("Listing security group rules for %q", securityGroupID)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListSecurityGroupRulesOptions(securityGroupID)
	resources, _, err := o.vpcSvc.ListSecurityGroupRulesWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list security group rules for %q", securityGroupID)
	}

	result := []cloudResource{}
	for _, securityGroupRule := range resources.Rules {
		switch reflect.TypeOf(securityGroupRule).String() {

		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
			{
				rule := securityGroupRule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
				result = append(result, cloudResource{
					key:      *rule.ID,
					name:     *rule.ID,
					status:   "",
					typeName: securityGroupRuleTypeName,
					id:       *rule.ID,
				})
			}
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
			{
				rule := securityGroupRule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
				result = append(result, cloudResource{
					key:      *rule.ID,
					name:     *rule.ID,
					status:   "",
					typeName: securityGroupRuleTypeName,
					id:       *rule.ID,
				})
			}
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
			{
				rule := securityGroupRule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
				result = append(result, cloudResource{
					key:      *rule.ID,
					name:     *rule.ID,
					status:   "",
					typeName: securityGroupRuleTypeName,
					id:       *rule.ID,
				})
			}
		default:
			{
				o.Logger.Debugf("Unknown rule: %q", securityGroupRule)
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteSecurityGroup(item cloudResource) error {
	o.Logger.Debugf("Deleting security group %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	found, err := o.listSecurityGroupRules(item.id)
	if err != nil {
		return err
	}

	rules := o.insertPendingItems(securityGroupRuleTypeName, found.list())

	for _, rule := range rules {
		if _, ok := found[rule.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(rule.typeName, []cloudResource{rule})
			o.Logger.Infof("Deleted security group rule %q", rule.name)
			continue
		}
		err = o.deleteSecurityGroupRule(rule, item.id)
		if err != nil {
			o.errorTracker.suppressWarning(rule.key, err, o.Logger)
		}
	}

	if rules = o.getPendingItems(securityGroupRuleTypeName); len(rules) > 0 {
		return errors.Errorf("%d items pending", len(rules))
	}

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

func (o *ClusterUninstaller) deleteSecurityGroupRule(item cloudResource, securityGroupID string) error {
	o.Logger.Debugf("Deleting security group rule %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteSecurityGroupRuleOptions(securityGroupID, item.id)
	details, err := o.vpcSvc.DeleteSecurityGroupRuleWithContext(ctx, options)
	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted security group rule %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete security group rule %s", item.name)
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
