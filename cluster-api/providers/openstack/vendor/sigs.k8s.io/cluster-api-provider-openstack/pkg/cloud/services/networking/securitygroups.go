/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package networking

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
)

const (
	secGroupPrefix     string = "k8s"
	controlPlaneSuffix string = "controlplane"
	workerSuffix       string = "worker"
	bastionSuffix      string = "bastion"
	remoteGroupIDSelf  string = "self"
)

// ReconcileSecurityGroups reconcile the security groups.
func (s *Service) ReconcileSecurityGroups(openStackCluster *infrav1.OpenStackCluster, clusterName string) error {
	s.scope.Logger().Info("Reconciling security groups")
	if !openStackCluster.Spec.ManagedSecurityGroups {
		s.scope.Logger().V(4).Info("No need to reconcile security groups")
		return nil
	}

	secControlPlaneGroupName := getSecControlPlaneGroupName(clusterName)
	secWorkerGroupName := getSecWorkerGroupName(clusterName)
	secGroupNames := map[string]string{
		controlPlaneSuffix: secControlPlaneGroupName,
		workerSuffix:       secWorkerGroupName,
	}

	if openStackCluster.Spec.Bastion != nil && openStackCluster.Spec.Bastion.Enabled {
		secBastionGroupName := getSecBastionGroupName(clusterName)
		secGroupNames[bastionSuffix] = secBastionGroupName
	}

	// create security groups first, because desired rules use group ids.
	for _, v := range secGroupNames {
		if err := s.createSecurityGroupIfNotExists(openStackCluster, v); err != nil {
			return err
		}
	}
	// create desired security groups
	desiredSecGroups, err := s.generateDesiredSecGroups(openStackCluster, secGroupNames)
	if err != nil {
		return err
	}

	observedSecGroups := make(map[string]*infrav1.SecurityGroup)
	for k, desiredSecGroup := range desiredSecGroups {
		var err error
		observedSecGroups[k], err = s.getSecurityGroupByName(desiredSecGroup.Name)

		if err != nil {
			return err
		}

		if observedSecGroups[k].ID != "" {
			observedSecGroup, err := s.reconcileGroupRules(desiredSecGroup, *observedSecGroups[k])
			if err != nil {
				return err
			}
			observedSecGroups[k] = &observedSecGroup
			continue
		}
	}

	openStackCluster.Status.ControlPlaneSecurityGroup = observedSecGroups[controlPlaneSuffix]
	openStackCluster.Status.WorkerSecurityGroup = observedSecGroups[workerSuffix]
	openStackCluster.Status.BastionSecurityGroup = observedSecGroups[bastionSuffix]

	return nil
}

func (s *Service) generateDesiredSecGroups(openStackCluster *infrav1.OpenStackCluster, secGroupNames map[string]string) (map[string]infrav1.SecurityGroup, error) {
	desiredSecGroups := make(map[string]infrav1.SecurityGroup)

	var secControlPlaneGroupID string
	var secWorkerGroupID string
	var secBastionGroupID string
	for i, v := range secGroupNames {
		secGroup, err := s.getSecurityGroupByName(v)
		if err != nil {
			return desiredSecGroups, err
		}
		switch i {
		case controlPlaneSuffix:
			secControlPlaneGroupID = secGroup.ID
		case workerSuffix:
			secWorkerGroupID = secGroup.ID
		case bastionSuffix:
			secBastionGroupID = secGroup.ID
		}
	}

	// Start with the default rules
	controlPlaneRules := append([]infrav1.SecurityGroupRule{}, defaultRules...)
	workerRules := append([]infrav1.SecurityGroupRule{}, defaultRules...)

	controlPlaneRules = append(controlPlaneRules, GetSGControlPlaneHTTPS()...)
	workerRules = append(workerRules, GetSGWorkerNodePort()...)

	// If we set additional ports to LB, we need create secgroup rules those ports, this apply to controlPlaneRules only
	if openStackCluster.Spec.APIServerLoadBalancer.Enabled {
		controlPlaneRules = append(controlPlaneRules, GetSGControlPlaneAdditionalPorts(openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts)...)
	}

	if openStackCluster.Spec.AllowAllInClusterTraffic {
		// Permit all ingress from the cluster security groups
		controlPlaneRules = append(controlPlaneRules, GetSGControlPlaneAllowAll(remoteGroupIDSelf, secWorkerGroupID)...)
		workerRules = append(workerRules, GetSGWorkerAllowAll(remoteGroupIDSelf, secControlPlaneGroupID)...)
	} else {
		controlPlaneRules = append(controlPlaneRules, GetSGControlPlaneGeneral(remoteGroupIDSelf, secWorkerGroupID)...)
		workerRules = append(workerRules, GetSGWorkerGeneral(remoteGroupIDSelf, secControlPlaneGroupID)...)
	}

	if openStackCluster.Spec.Bastion != nil && openStackCluster.Spec.Bastion.Enabled {
		controlPlaneRules = append(controlPlaneRules, GetSGControlPlaneSSH(secBastionGroupID)...)
		workerRules = append(workerRules, GetSGWorkerSSH(secBastionGroupID)...)

		desiredSecGroups[bastionSuffix] = infrav1.SecurityGroup{
			Name: secGroupNames[bastionSuffix],
			Rules: append(
				[]infrav1.SecurityGroupRule{
					{
						Description:  "SSH",
						Direction:    "ingress",
						EtherType:    "IPv4",
						PortRangeMin: 22,
						PortRangeMax: 22,
						Protocol:     "tcp",
					},
				},
				defaultRules...,
			),
		}
	}

	desiredSecGroups[controlPlaneSuffix] = infrav1.SecurityGroup{
		Name:  secGroupNames[controlPlaneSuffix],
		Rules: controlPlaneRules,
	}

	desiredSecGroups[workerSuffix] = infrav1.SecurityGroup{
		Name:  secGroupNames[workerSuffix],
		Rules: workerRules,
	}

	return desiredSecGroups, nil
}

func (s *Service) GetSecurityGroups(securityGroupParams []infrav1.SecurityGroupFilter) ([]string, error) {
	var sgIDs []string
	for _, sg := range securityGroupParams {
		// Don't validate an explicit UUID if we were given one
		if sg.ID != "" {
			if isDuplicate(sgIDs, sg.ID) {
				continue
			}
			sgIDs = append(sgIDs, sg.ID)
			continue
		}

		listOpts := sg.ToListOpt()
		if listOpts.ProjectID == "" {
			listOpts.ProjectID = s.scope.ProjectID()
		}
		SGList, err := s.client.ListSecGroup(listOpts)
		if err != nil {
			return nil, err
		}

		if len(SGList) == 0 {
			return nil, fmt.Errorf("security group %s not found", sg.Name)
		}

		for _, group := range SGList {
			if isDuplicate(sgIDs, group.ID) {
				continue
			}
			sgIDs = append(sgIDs, group.ID)
		}
	}
	return sgIDs, nil
}

func (s *Service) DeleteSecurityGroups(openStackCluster *infrav1.OpenStackCluster, clusterName string) error {
	secGroupNames := []string{
		getSecControlPlaneGroupName(clusterName),
		getSecWorkerGroupName(clusterName),
	}
	for _, secGroupName := range secGroupNames {
		if err := s.deleteSecurityGroup(openStackCluster, secGroupName); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteBastionSecurityGroup(openStackCluster *infrav1.OpenStackCluster, clusterName string) error {
	secBastionGroupName := getSecBastionGroupName(clusterName)
	return s.deleteSecurityGroup(openStackCluster, secBastionGroupName)
}

func (s *Service) deleteSecurityGroup(openStackCluster *infrav1.OpenStackCluster, name string) error {
	group, err := s.getSecurityGroupByName(name)
	if err != nil {
		return err
	}
	if group.ID == "" {
		// nothing to do
		return nil
	}
	err = s.client.DeleteSecGroup(group.ID)
	if err != nil {
		record.Warnf(openStackCluster, "FailedDeleteSecurityGroup", "Failed to delete security group %s with id %s: %v", group.Name, group.ID, err)
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulDeleteSecurityGroup", "Deleted security group %s with id %s", group.Name, group.ID)
	return nil
}

// reconcileGroupRules reconciles an already existing observed group by deleting rules not needed anymore and
// creating rules that are missing.
func (s *Service) reconcileGroupRules(desired, observed infrav1.SecurityGroup) (infrav1.SecurityGroup, error) {
	rulesToDelete := []infrav1.SecurityGroupRule{}
	// fills rulesToDelete by calculating observed - desired
	for _, observedRule := range observed.Rules {
		deleteRule := true
		for _, desiredRule := range desired.Rules {
			r := desiredRule
			if r.RemoteGroupID == remoteGroupIDSelf {
				r.RemoteGroupID = observed.ID
			}
			if r.Equal(observedRule) {
				deleteRule = false
				break
			}
		}
		if deleteRule {
			rulesToDelete = append(rulesToDelete, observedRule)
		}
	}

	rulesToCreate := []infrav1.SecurityGroupRule{}
	reconciledRules := make([]infrav1.SecurityGroupRule, 0, len(desired.Rules))
	// fills rulesToCreate by calculating desired - observed
	// also adds rules which are in observed and desired to reconciledRules.
	for _, desiredRule := range desired.Rules {
		r := desiredRule
		if r.RemoteGroupID == remoteGroupIDSelf {
			r.RemoteGroupID = observed.ID
		}
		createRule := true
		for _, observedRule := range observed.Rules {
			if r.Equal(observedRule) {
				// add already existing rules to reconciledRules because we won't touch them anymore
				reconciledRules = append(reconciledRules, observedRule)
				createRule = false
				break
			}
		}
		if createRule {
			rulesToCreate = append(rulesToCreate, desiredRule)
		}
	}

	s.scope.Logger().V(4).Info("Deleting rules not needed anymore for group", "name", observed.Name, "amount", len(rulesToDelete))
	for _, rule := range rulesToDelete {
		s.scope.Logger().V(6).Info("Deleting rule", "ID", rule.ID, "name", observed.Name)
		err := s.client.DeleteSecGroupRule(rule.ID)
		if err != nil {
			return infrav1.SecurityGroup{}, err
		}
	}

	s.scope.Logger().V(4).Info("Creating new rules needed for group", "name", observed.Name, "amount", len(rulesToCreate))
	for _, rule := range rulesToCreate {
		r := rule
		r.SecurityGroupID = observed.ID
		if r.RemoteGroupID == remoteGroupIDSelf {
			r.RemoteGroupID = observed.ID
		}
		newRule, err := s.createRule(r)
		if err != nil {
			return infrav1.SecurityGroup{}, err
		}
		reconciledRules = append(reconciledRules, newRule)
	}
	observed.Rules = reconciledRules

	return observed, nil
}

func (s *Service) createSecurityGroupIfNotExists(openStackCluster *infrav1.OpenStackCluster, groupName string) error {
	secGroup, err := s.getSecurityGroupByName(groupName)
	if err != nil {
		return err
	}
	if secGroup == nil || secGroup.ID == "" {
		s.scope.Logger().V(6).Info("Group doesn't exist, creating it", "name", groupName)

		createOpts := groups.CreateOpts{
			Name:        groupName,
			Description: "Cluster API managed group",
		}
		s.scope.Logger().V(6).Info("Creating group", "name", groupName)

		group, err := s.client.CreateSecGroup(createOpts)
		if err != nil {
			record.Warnf(openStackCluster, "FailedCreateSecurityGroup", "Failed to create security group %s: %v", groupName, err)
			return err
		}

		if len(openStackCluster.Spec.Tags) > 0 {
			_, err = s.client.ReplaceAllAttributesTags("security-groups", group.ID, attributestags.ReplaceAllOpts{
				Tags: openStackCluster.Spec.Tags,
			})
			if err != nil {
				return err
			}
		}

		record.Eventf(openStackCluster, "SuccessfulCreateSecurityGroup", "Created security group %s with id %s", groupName, group.ID)
		return nil
	}

	sInfo := fmt.Sprintf("Reuse Existing SecurityGroup %s with %s", groupName, secGroup.ID)
	s.scope.Logger().V(6).Info(sInfo)

	return nil
}

func (s *Service) getSecurityGroupByName(name string) (*infrav1.SecurityGroup, error) {
	opts := groups.ListOpts{
		Name: name,
	}

	s.scope.Logger().V(6).Info("Attempting to fetch security group with", "name", name)
	allGroups, err := s.client.ListSecGroup(opts)
	if err != nil {
		return &infrav1.SecurityGroup{}, err
	}

	switch len(allGroups) {
	case 0:
		return &infrav1.SecurityGroup{}, nil
	case 1:
		return convertOSSecGroupToConfigSecGroup(allGroups[0]), nil
	}

	return &infrav1.SecurityGroup{}, fmt.Errorf("more than one security group found named: %s", name)
}

func (s *Service) createRule(r infrav1.SecurityGroupRule) (infrav1.SecurityGroupRule, error) {
	dir := rules.RuleDirection(r.Direction)
	proto := rules.RuleProtocol(r.Protocol)
	etherType := rules.RuleEtherType(r.EtherType)

	createOpts := rules.CreateOpts{
		Description:    r.Description,
		Direction:      dir,
		PortRangeMin:   r.PortRangeMin,
		PortRangeMax:   r.PortRangeMax,
		Protocol:       proto,
		EtherType:      etherType,
		RemoteGroupID:  r.RemoteGroupID,
		RemoteIPPrefix: r.RemoteIPPrefix,
		SecGroupID:     r.SecurityGroupID,
	}
	s.scope.Logger().V(6).Info("Creating rule", "description", r.Description, "direction", dir, "portRangeMin", r.PortRangeMin, "portRangeMax", r.PortRangeMax, "proto", proto, "etherType", etherType, "remoteGroupID", r.RemoteGroupID, "remoteIPPrefix", r.RemoteIPPrefix, "securityGroupID", r.SecurityGroupID)
	rule, err := s.client.CreateSecGroupRule(createOpts)
	if err != nil {
		return infrav1.SecurityGroupRule{}, err
	}
	return convertOSSecGroupRuleToConfigSecGroupRule(*rule), nil
}

func getSecControlPlaneGroupName(clusterName string) string {
	return fmt.Sprintf("%s-cluster-%s-secgroup-%s", secGroupPrefix, clusterName, controlPlaneSuffix)
}

func getSecWorkerGroupName(clusterName string) string {
	return fmt.Sprintf("%s-cluster-%s-secgroup-%s", secGroupPrefix, clusterName, workerSuffix)
}

func getSecBastionGroupName(clusterName string) string {
	return fmt.Sprintf("%s-cluster-%s-secgroup-%s", secGroupPrefix, clusterName, bastionSuffix)
}

func convertOSSecGroupToConfigSecGroup(osSecGroup groups.SecGroup) *infrav1.SecurityGroup {
	securityGroupRules := make([]infrav1.SecurityGroupRule, len(osSecGroup.Rules))
	for i, rule := range osSecGroup.Rules {
		securityGroupRules[i] = convertOSSecGroupRuleToConfigSecGroupRule(rule)
	}
	return &infrav1.SecurityGroup{
		ID:    osSecGroup.ID,
		Name:  osSecGroup.Name,
		Rules: securityGroupRules,
	}
}

func convertOSSecGroupRuleToConfigSecGroupRule(osSecGroupRule rules.SecGroupRule) infrav1.SecurityGroupRule {
	return infrav1.SecurityGroupRule{
		ID:              osSecGroupRule.ID,
		Direction:       osSecGroupRule.Direction,
		Description:     osSecGroupRule.Description,
		EtherType:       osSecGroupRule.EtherType,
		SecurityGroupID: osSecGroupRule.SecGroupID,
		PortRangeMin:    osSecGroupRule.PortRangeMin,
		PortRangeMax:    osSecGroupRule.PortRangeMax,
		Protocol:        osSecGroupRule.Protocol,
		RemoteGroupID:   osSecGroupRule.RemoteGroupID,
		RemoteIPPrefix:  osSecGroupRule.RemoteIPPrefix,
	}
}

func isDuplicate(list []string, name string) bool {
	if len(list) == 0 {
		return false
	}
	for _, element := range list {
		if element == name {
			return true
		}
	}
	return false
}
