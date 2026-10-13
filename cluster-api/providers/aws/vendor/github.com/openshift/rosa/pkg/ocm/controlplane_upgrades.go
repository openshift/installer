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

import cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

func (c *Client) CancelControlPlaneUpgrade(clusterID, upgradeID string) (bool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).ControlPlane().UpgradePolicies().
		ControlPlaneUpgradePolicy(upgradeID).Delete().Send()
	if err != nil {
		return false, handleErr(response.Error(), err)
	}
	return true, nil
}

func (c *Client) GetControlPlaneScheduledUpgrade(clusterID string) (*cmv1.ControlPlaneUpgradePolicy, error) {
	upgradePolicies, err := c.GetControlPlaneUpgradePolicies(clusterID)
	if err != nil {
		return nil, err
	}
	for _, upgradePolicy := range upgradePolicies {
		if upgradePolicy.UpgradeType() == cmv1.UpgradeTypeControlPlane {
			return upgradePolicy, nil
		}
	}

	return nil, nil
}

func (c *Client) GetControlPlaneUpgradePolicies(clusterID string) (
	controlPlaneUpgradePolicies []*cmv1.ControlPlaneUpgradePolicy,
	err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ControlPlane().
		UpgradePolicies()
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
		controlPlaneUpgradePolicies = append(controlPlaneUpgradePolicies, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
}

func (c *Client) ScheduleHypershiftControlPlaneUpgrade(clusterID string,
	upgradePolicy *cmv1.ControlPlaneUpgradePolicy) (*cmv1.ControlPlaneUpgradePolicy, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).ControlPlane().
		UpgradePolicies().
		Add().Body(upgradePolicy).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}
