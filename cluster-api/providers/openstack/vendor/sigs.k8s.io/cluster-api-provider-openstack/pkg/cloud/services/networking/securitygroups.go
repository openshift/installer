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
	"errors"
	"fmt"
	"slices"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"k8s.io/utils/net"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/filterconvert"
)

const (
	secGroupPrefix     string = "k8s"
	controlPlaneSuffix string = "controlplane"
	workerSuffix       string = "worker"
	bastionSuffix      string = "bastion"
	allNodesSuffix     string = "allNodes"
	remoteGroupIDSelf  string = "self"
)

// ReconcileSecurityGroups reconcile the security groups.
func (s *Service) ReconcileSecurityGroups(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	s.scope.Logger().Info("Reconciling security groups")
	if openStackCluster.Spec.ManagedSecurityGroups == nil {
		s.scope.Logger().V(4).Info("No need to reconcile security groups")
		return nil
	}

	bastionEnabled := openStackCluster.Spec.Bastion.IsEnabled()

	secControlPlaneGroupName := getSecControlPlaneGroupName(clusterResourceName)
	secWorkerGroupName := getSecWorkerGroupName(clusterResourceName)
	suffixToNameMap := map[string]string{
		controlPlaneSuffix: secControlPlaneGroupName,
		workerSuffix:       secWorkerGroupName,
	}

	secBastionGroupName := getSecBastionGroupName(clusterResourceName)
	if bastionEnabled {
		suffixToNameMap[bastionSuffix] = secBastionGroupName
	} else {
		// We reconcile the security groups before the bastion, because the bastion
		// needs its security group to be created first when managed security groups are enabled.
		// When the bastion is disabled, we will try to delete the security group if it exists.
		// In the first attempt, the security group will still be in-use by the bastion instance
		// but then the bastion instance will be deleted in the next reconcile loop.
		// We do that here because we don't want to manage the bastion security group from
		// elsewhere, that could cause infinite loops between ReconCileSecurityGroups and ReconcileBastion.
		// Therefore we try to delete the bastion security group as a best effort here
		// and also when the cluster is deleted so we're sure it will be deleted at some point.
		// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues/2113
		if err := s.deleteSecurityGroup(openStackCluster, secBastionGroupName); err != nil {
			s.scope.Logger().Info("Non-fatal error when deleting the bastion security group", "name", secBastionGroupName, "error", err)
			return nil
		}
	}

	// create security groups first, because desired rules use group ids.
	observedSecGroupBySuffix := make(map[string]*groups.SecGroup)
	for suffix, secGroupName := range suffixToNameMap {
		group, err := s.getOrCreateSecurityGroup(openStackCluster, secGroupName)
		if err != nil {
			return err
		}
		observedSecGroupBySuffix[suffix] = group

		normaliseTags := func(tags []string) []string {
			tags = slices.Clone(tags)
			slices.Sort(tags)
			return slices.Compact(tags)
		}

		if !slices.Equal(normaliseTags(openStackCluster.Spec.Tags), normaliseTags(group.Tags)) {
			_, err = s.client.ReplaceAllAttributesTags("security-groups", group.ID, attributestags.ReplaceAllOpts{
				Tags: openStackCluster.Spec.Tags,
			})
			if err != nil {
				return err
			}
			s.scope.Logger().V(6).Info("Updated tags for security group", "name", group.Name, "id", group.ID)
		}
	}

	// create desired security groups
	desiredSecGroupsBySuffix, err := s.generateDesiredSecGroups(openStackCluster, suffixToNameMap, observedSecGroupBySuffix)
	if err != nil {
		return err
	}

	for suffix := range desiredSecGroupsBySuffix {
		desiredSecGroup := desiredSecGroupsBySuffix[suffix]
		observedSecGroup, ok := observedSecGroupBySuffix[suffix]
		if !ok {
			// This should never happen
			return fmt.Errorf("unable to reconcile security groups: security group %s not found", suffix)
		}

		err := s.reconcileGroupRules(&desiredSecGroup, observedSecGroup)
		if err != nil {
			return err
		}
		continue
	}

	openStackCluster.Status.ControlPlaneSecurityGroup = convertOSSecGroupToConfigSecGroup(observedSecGroupBySuffix[controlPlaneSuffix])
	openStackCluster.Status.WorkerSecurityGroup = convertOSSecGroupToConfigSecGroup(observedSecGroupBySuffix[workerSuffix])
	if bastionEnabled {
		openStackCluster.Status.BastionSecurityGroup = convertOSSecGroupToConfigSecGroup(observedSecGroupBySuffix[bastionSuffix])
	} else {
		openStackCluster.Status.BastionSecurityGroup = nil
	}

	return nil
}

type securityGroupSpec struct {
	Name  string
	Rules []resolvedSecurityGroupRuleSpec
}

type resolvedSecurityGroupRuleSpec struct {
	Description    string `json:"description,omitempty"`
	Direction      string `json:"direction,omitempty"`
	EtherType      string `json:"etherType,omitempty"`
	PortRangeMin   int    `json:"portRangeMin,omitempty"`
	PortRangeMax   int    `json:"portRangeMax,omitempty"`
	Protocol       string `json:"protocol,omitempty"`
	RemoteGroupID  string `json:"remoteGroupID,omitempty"`
	RemoteIPPrefix string `json:"remoteIPPrefix,omitempty"`
}

func (r resolvedSecurityGroupRuleSpec) Matches(other rules.SecGroupRule) bool {
	return r.Description == other.Description &&
		r.Direction == other.Direction &&
		r.EtherType == other.EtherType &&
		r.PortRangeMin == other.PortRangeMin &&
		r.PortRangeMax == other.PortRangeMax &&
		r.Protocol == other.Protocol &&
		r.RemoteGroupID == other.RemoteGroupID &&
		r.RemoteIPPrefix == other.RemoteIPPrefix
}

func (s *Service) generateDesiredSecGroups(openStackCluster *infrav1.OpenStackCluster, suffixToNameMap map[string]string, observedSecGroupsBySuffix map[string]*groups.SecGroup) (map[string]securityGroupSpec, error) {
	if openStackCluster.Spec.ManagedSecurityGroups == nil {
		return nil, nil
	}

	var secControlPlaneGroupID string
	var secWorkerGroupID string
	var secBastionGroupID string

	// remoteManagedGroups is a map of suffix to security group ID.
	// It will be used to fill in the RemoteGroupID field of the security group rules
	// that reference a managed security group.
	// For now, we only reference the control plane and worker security groups.
	remoteManagedGroups := make(map[string]string)

	for suffix := range suffixToNameMap {
		secGroup, ok := observedSecGroupsBySuffix[suffix]
		if !ok {
			// This should never happen, as we should have created the security group earlier in this reconcile if it does not exist.
			return nil, fmt.Errorf("unable to generate desired security group rules: security group for %s not found", suffix)
		}
		switch suffix {
		case controlPlaneSuffix:
			secControlPlaneGroupID = secGroup.ID
			remoteManagedGroups[controlPlaneSuffix] = secControlPlaneGroupID
		case workerSuffix:
			secWorkerGroupID = secGroup.ID
			remoteManagedGroups[workerSuffix] = secWorkerGroupID
		case bastionSuffix:
			secBastionGroupID = secGroup.ID
			remoteManagedGroups[bastionSuffix] = secBastionGroupID
		}
	}

	// Start with the default rules
	controlPlaneRules := append([]resolvedSecurityGroupRuleSpec{}, defaultRules...)
	workerRules := append([]resolvedSecurityGroupRuleSpec{}, defaultRules...)

	controlPlaneRules = append(controlPlaneRules, getSGControlPlaneHTTPS()...)

	// Fetch subnet to use for worker node port rules
	// In the future IPv6 support need to be added here
	if openStackCluster.Status.Network != nil {
		for _, subnet := range openStackCluster.Status.Network.Subnets {
			if net.IsIPv4CIDRString(subnet.CIDR) {
				workerRules = append(workerRules, getSGWorkerNodePortCIDR(subnet.CIDR)...)
			}
		}
	}

	// Add rules allowing nodepors from all cluster nodes, this will take effect even if no subnetCIDR is found to ensure all nodes can commincate over nodeports at all time
	workerRules = append(workerRules, getSGWorkerNodePort(secWorkerGroupID, secControlPlaneGroupID)...)

	// If we set additional ports to LB, we need create secgroup rules those ports, this apply to controlPlaneRules only
	if openStackCluster.Spec.APIServerLoadBalancer.IsEnabled() {
		controlPlaneRules = append(controlPlaneRules, getSGControlPlaneAdditionalPorts(openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts)...)
	}

	if openStackCluster.Spec.ManagedSecurityGroups != nil && openStackCluster.Spec.ManagedSecurityGroups.AllowAllInClusterTraffic {
		// Permit all ingress from the cluster security groups
		controlPlaneRules = append(controlPlaneRules, getSGControlPlaneAllowAll(remoteGroupIDSelf, secWorkerGroupID)...)
		workerRules = append(workerRules, getSGWorkerAllowAll(remoteGroupIDSelf, secControlPlaneGroupID)...)
	} else {
		controlPlaneRules = append(controlPlaneRules, getSGControlPlaneGeneral(remoteGroupIDSelf, secWorkerGroupID)...)
		workerRules = append(workerRules, getSGWorkerGeneral(remoteGroupIDSelf, secControlPlaneGroupID)...)
	}

	// Append any additional rules for control plane and worker nodes
	controlPlaneExtraRules, err := getRulesFromSpecs(remoteManagedGroups, openStackCluster.Spec.ManagedSecurityGroups.ControlPlaneNodesSecurityGroupRules)
	if err != nil {
		return nil, err
	}
	controlPlaneRules = append(controlPlaneRules, controlPlaneExtraRules...)
	workersExtraRules, err := getRulesFromSpecs(remoteManagedGroups, openStackCluster.Spec.ManagedSecurityGroups.WorkerNodesSecurityGroupRules)
	if err != nil {
		return nil, err
	}
	workerRules = append(workerRules, workersExtraRules...)

	// For now, we do not create a separate security group for allNodes.
	// Instead, we append the rules for allNodes to the control plane and worker security groups.
	allNodesRules, err := getRulesFromSpecs(remoteManagedGroups, openStackCluster.Spec.ManagedSecurityGroups.AllNodesSecurityGroupRules)
	if err != nil {
		return nil, err
	}
	controlPlaneRules = append(controlPlaneRules, allNodesRules...)
	workerRules = append(workerRules, allNodesRules...)

	desiredSecGroupsBySuffix := make(map[string]securityGroupSpec)

	if openStackCluster.Spec.Bastion.IsEnabled() {
		controlPlaneRules = append(controlPlaneRules, getSGControlPlaneSSH(secBastionGroupID)...)
		workerRules = append(workerRules, getSGWorkerSSH(secBastionGroupID)...)

		desiredSecGroupsBySuffix[bastionSuffix] = securityGroupSpec{
			Name: suffixToNameMap[bastionSuffix],
			Rules: append(
				[]resolvedSecurityGroupRuleSpec{
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

	desiredSecGroupsBySuffix[controlPlaneSuffix] = securityGroupSpec{
		Name:  suffixToNameMap[controlPlaneSuffix],
		Rules: controlPlaneRules,
	}

	desiredSecGroupsBySuffix[workerSuffix] = securityGroupSpec{
		Name:  suffixToNameMap[workerSuffix],
		Rules: workerRules,
	}
	return desiredSecGroupsBySuffix, nil
}

// getAllNodesRules returns the rules for the allNodes security group that should be created.
func getRulesFromSpecs(remoteManagedGroups map[string]string, securityGroupRules []infrav1.SecurityGroupRuleSpec) ([]resolvedSecurityGroupRuleSpec, error) {
	rules := make([]resolvedSecurityGroupRuleSpec, 0, len(securityGroupRules))
	for _, rule := range securityGroupRules {
		if err := validateRemoteManagedGroups(remoteManagedGroups, rule.RemoteManagedGroups); err != nil {
			return nil, err
		}
		r := resolvedSecurityGroupRuleSpec{
			Direction: rule.Direction,
		}
		if rule.Description != nil {
			r.Description = *rule.Description
		}
		if rule.EtherType != nil {
			r.EtherType = *rule.EtherType
		}
		if rule.PortRangeMin != nil {
			r.PortRangeMin = *rule.PortRangeMin
		}
		if rule.PortRangeMax != nil {
			r.PortRangeMax = *rule.PortRangeMax
		}
		if rule.Protocol != nil {
			r.Protocol = *rule.Protocol
		}
		if rule.RemoteGroupID != nil {
			r.RemoteGroupID = *rule.RemoteGroupID
		}
		if rule.RemoteIPPrefix != nil {
			r.RemoteIPPrefix = *rule.RemoteIPPrefix
		}

		if len(rule.RemoteManagedGroups) > 0 {
			if rule.RemoteGroupID != nil {
				return nil, fmt.Errorf("remoteGroupID must not be set when remoteManagedGroups is set")
			}

			for _, rg := range rule.RemoteManagedGroups {
				rc := r
				rc.RemoteGroupID = remoteManagedGroups[rg.String()]
				rules = append(rules, rc)
			}
		} else {
			rules = append(rules, r)
		}
	}
	return rules, nil
}

// validateRemoteManagedGroups validates that the remoteManagedGroups target existing managed security groups.
func validateRemoteManagedGroups(remoteManagedGroups map[string]string, ruleRemoteManagedGroups []infrav1.ManagedSecurityGroupName) error {
	for _, group := range ruleRemoteManagedGroups {
		if _, ok := remoteManagedGroups[group.String()]; !ok {
			return fmt.Errorf("remoteManagedGroups: %s is not a valid remote managed security group", group)
		}
	}
	return nil
}

func (s *Service) GetSecurityGroups(securityGroupParams []infrav1.SecurityGroupParam) ([]string, error) {
	var sgIDs []string
	for i := range securityGroupParams {
		sg := &securityGroupParams[i]

		// Don't validate an explicit UUID if we were given one
		if sg.ID != nil {
			if isDuplicate(sgIDs, *sg.ID) {
				continue
			}
			sgIDs = append(sgIDs, *sg.ID)
			continue
		}

		if sg.Filter == nil {
			// Should have been caught by validation
			return nil, errors.New("security group param must have id or filter")
		}

		listOpts := filterconvert.SecurityGroupFilterToListOpts(sg.Filter)
		if listOpts.ProjectID == "" {
			listOpts.ProjectID = s.scope.ProjectID()
		}
		SGList, err := s.client.ListSecGroup(listOpts)
		if err != nil {
			return nil, err
		}

		if len(SGList) == 0 {
			return nil, fmt.Errorf("security group %d not found", i)
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

func (s *Service) DeleteSecurityGroups(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	secGroupNames := []string{
		getSecControlPlaneGroupName(clusterResourceName),
		getSecWorkerGroupName(clusterResourceName),
		// Even if the bastion might be disabled, we still try to delete the security group in case
		// we had a bastion before and for some reason we didn't delete its security group.
		// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues/2113
		getSecBastionGroupName(clusterResourceName),
	}

	for _, secGroupName := range secGroupNames {
		if err := s.deleteSecurityGroup(openStackCluster, secGroupName); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) deleteSecurityGroup(openStackCluster *infrav1.OpenStackCluster, name string) error {
	group, err := s.getSecurityGroupByName(name)
	if err != nil {
		return err
	}
	if group == nil {
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
func (s *Service) reconcileGroupRules(desired *securityGroupSpec, observed *groups.SecGroup) error {
	var rulesToDelete []string
	// fills rulesToDelete by calculating observed - desired
	for _, observedRule := range observed.Rules {
		deleteRule := true
		for _, desiredRule := range desired.Rules {
			r := desiredRule
			if r.RemoteGroupID == remoteGroupIDSelf {
				r.RemoteGroupID = observed.ID
			}
			if r.Matches(observedRule) {
				deleteRule = false
				break
			}
		}
		if deleteRule {
			rulesToDelete = append(rulesToDelete, observedRule.ID)
		}
	}

	rulesToCreate := []resolvedSecurityGroupRuleSpec{}
	// fills rulesToCreate by calculating desired - observed
	// also adds rules which are in observed and desired to reconcileGroupRules.
	for _, desiredRule := range desired.Rules {
		r := desiredRule
		if r.RemoteGroupID == remoteGroupIDSelf {
			r.RemoteGroupID = observed.ID
		}
		createRule := true
		for _, observedRule := range observed.Rules {
			if r.Matches(observedRule) {
				// add already existing rules to reconciledRules because we won't touch them anymore
				createRule = false
				break
			}
		}
		if createRule {
			rulesToCreate = append(rulesToCreate, desiredRule)
		}
	}

	if len(rulesToDelete) > 0 {
		s.scope.Logger().V(4).Info("Deleting rules not needed anymore for group", "name", observed.Name, "amount", len(rulesToDelete))
		for _, rule := range rulesToDelete {
			s.scope.Logger().V(6).Info("Deleting rule", "ID", rule, "name", observed.Name)
			err := s.client.DeleteSecGroupRule(rule)
			if err != nil {
				return err
			}
		}
	}

	if len(rulesToCreate) > 0 {
		s.scope.Logger().V(4).Info("Creating new rules needed for group", "name", observed.Name, "amount", len(rulesToCreate))
		for _, rule := range rulesToCreate {
			r := rule
			if r.RemoteGroupID == remoteGroupIDSelf {
				r.RemoteGroupID = observed.ID
			}
			err := s.createRule(observed.ID, r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) getOrCreateSecurityGroup(openStackCluster *infrav1.OpenStackCluster, groupName string) (*groups.SecGroup, error) {
	secGroup, err := s.getSecurityGroupByName(groupName)
	if err != nil {
		return nil, err
	}
	if secGroup != nil {
		s.scope.Logger().V(6).Info("Reusing existing SecurityGroup", "name", groupName, "id", secGroup.ID)
		return secGroup, nil
	}

	s.scope.Logger().V(6).Info("Group doesn't exist, creating it", "name", groupName)

	createOpts := groups.CreateOpts{
		Name:        groupName,
		Description: "Cluster API managed group",
	}
	s.scope.Logger().V(6).Info("Creating group", "name", groupName)

	group, err := s.client.CreateSecGroup(createOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateSecurityGroup", "Failed to create security group %s: %v", groupName, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulCreateSecurityGroup", "Created security group %s with id %s", groupName, group.ID)
	return group, nil
}

func (s *Service) getSecurityGroupByName(name string) (*groups.SecGroup, error) {
	opts := groups.ListOpts{
		Name: name,
	}

	s.scope.Logger().V(6).Info("Attempting to fetch security group with", "name", name)
	allGroups, err := s.client.ListSecGroup(opts)
	if err != nil {
		return nil, err
	}

	switch len(allGroups) {
	case 0:
		return nil, nil
	case 1:
		return &allGroups[0], nil
	}

	return nil, fmt.Errorf("more than one security group found named: %s", name)
}

func (s *Service) createRule(securityGroupID string, r resolvedSecurityGroupRuleSpec) error {
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
		SecGroupID:     securityGroupID,
	}
	s.scope.Logger().V(6).Info("Creating rule", "description", r.Description, "direction", dir, "portRangeMin", r.PortRangeMin, "portRangeMax", r.PortRangeMax, "proto", proto, "etherType", etherType, "remoteGroupID", r.RemoteGroupID, "remoteIPPrefix", r.RemoteIPPrefix, "securityGroupID", securityGroupID)
	_, err := s.client.CreateSecGroupRule(createOpts)
	if err != nil {
		return err
	}
	return nil
}

func getSecControlPlaneGroupName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s-secgroup-%s", secGroupPrefix, clusterResourceName, controlPlaneSuffix)
}

func getSecWorkerGroupName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s-secgroup-%s", secGroupPrefix, clusterResourceName, workerSuffix)
}

func getSecBastionGroupName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s-secgroup-%s", secGroupPrefix, clusterResourceName, bastionSuffix)
}

func convertOSSecGroupToConfigSecGroup(osSecGroup *groups.SecGroup) *infrav1.SecurityGroupStatus {
	return &infrav1.SecurityGroupStatus{
		ID:   osSecGroup.ID,
		Name: osSecGroup.Name,
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
