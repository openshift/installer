/*
Copyright (c) 2023 Red Hat, Inc.

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

package ocm

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func (c *Client) ScheduleNodePoolUpgrade(clusterID string, nodePoolId string,
	upgradePolicy *cmv1.NodePoolUpgradePolicy) (*cmv1.NodePoolUpgradePolicy, error) {
	if upgradePolicy == nil {
		return nil, fmt.Errorf("upgrade policy is nil")
	}
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).NodePools().NodePool(nodePoolId).
		UpgradePolicies().
		Add().Body(upgradePolicy).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) BuildNodeUpgradePolicy(version string, machinePoolID string,
	scheduling UpgradeScheduling) (*cmv1.NodePoolUpgradePolicy, error) {
	upgradePolicyBuilder := cmv1.NewNodePoolUpgradePolicy().UpgradeType(cmv1.UpgradeTypeNodePool).NodePoolID(machinePoolID)
	if scheduling.AutomaticUpgrades {
		upgradePolicyBuilder.ScheduleType(cmv1.ScheduleTypeAutomatic).
			EnableMinorVersionUpgrades(scheduling.AllowMinorVersionUpdates).Schedule(scheduling.Schedule)
	} else {
		upgradePolicyBuilder.ScheduleType(cmv1.ScheduleTypeManual).Version(version).NextRun(scheduling.NextRun)
	}
	return upgradePolicyBuilder.Build()
}

func (c *Client) CancelNodePoolUpgrade(clusterID, nodePoolID string, upgradeID string) (bool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).NodePools().NodePool(nodePoolID).UpgradePolicies().
		NodePoolUpgradePolicy(upgradeID).Delete().Send()
	if err != nil {
		return false, handleErr(response.Error(), err)
	}
	return true, nil
}

func (c *Client) getNodePoolUpgradePolicies(clusterID string, nodePoolID string) (
	nodePoolUpgradePolicies []*cmv1.NodePoolUpgradePolicy,
	err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).NodePools().NodePool(nodePoolID).UpgradePolicies()
	page := 1
	size := 100
	for {
		response, err := collection.List().
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return nil, handleErr(response.Error(), err)
		}
		nodePoolUpgradePolicies = append(nodePoolUpgradePolicies, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
}

func (c *Client) GetHypershiftNodePoolUpgrades(clusterID, clusterKey,
	nodePoolID string) (*cmv1.NodePool, []*cmv1.NodePoolUpgradePolicy, error) {
	nodePool, exists, err := c.GetNodePool(clusterID, nodePoolID)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get machine pools for hosted cluster '%s': %v", clusterKey, err)
	}
	if !exists {
		return nil, nil, fmt.Errorf("Machine pool '%s' does not exist for hosted cluster '%s'", nodePoolID, clusterKey)
	}

	scheduledUpgrades, err := c.getNodePoolUpgradePolicies(clusterID, nodePoolID)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get scheduled upgrades for machine pool '%s': %v", nodePoolID, err)
	}

	return nodePool, scheduledUpgrades, nil
}

func (c *Client) GetHypershiftNodePoolUpgrade(clusterID, clusterKey,
	nodePoolID string) (*cmv1.NodePool, *cmv1.NodePoolUpgradePolicy, error) {
	nodePool, upgradePolicies, err := c.GetHypershiftNodePoolUpgrades(clusterID, clusterKey, nodePoolID)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get scheduled upgrades for machine pool '%s': %v", nodePoolID, err)
	}
	for _, upgradePolicy := range upgradePolicies {
		if upgradePolicy.UpgradeType() == cmv1.UpgradeTypeNodePool {
			return nodePool, upgradePolicy, nil
		}
	}

	return nodePool, nil, nil
}
