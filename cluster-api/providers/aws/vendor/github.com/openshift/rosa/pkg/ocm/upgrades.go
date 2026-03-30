/*
Copyright (c) 2020 Red Hat, Inc.

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
	"encoding/json"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"
)

func (c *Client) GetUpgradePolicies(clusterID string) (upgradePolicies []*cmv1.UpgradePolicy, err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
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
		upgradePolicies = append(upgradePolicies, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
}

func (c *Client) GetScheduledUpgrade(clusterID string) (*cmv1.UpgradePolicy, *cmv1.UpgradePolicyState, error) {
	upgradePolicies, err := c.GetUpgradePolicies(clusterID)
	if err != nil {
		return nil, nil, err
	}
	for _, upgradePolicy := range upgradePolicies {
		if upgradePolicy.UpgradeType() == cmv1.UpgradeTypeOSD {
			state, err := c.ocm.ClustersMgmt().V1().
				Clusters().Cluster(clusterID).
				UpgradePolicies().UpgradePolicy(upgradePolicy.ID()).
				State().
				Get().
				Send()
			if err != nil {
				return nil, nil, err
			}

			return upgradePolicy, state.Body(), nil
		}
	}

	return nil, nil, nil
}

func (c *Client) ScheduleUpgrade(clusterID string, upgradePolicy *cmv1.UpgradePolicy) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		UpgradePolicies().
		Add().Body(upgradePolicy).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func (c *Client) CancelUpgrade(clusterID string) (bool, error) {
	scheduledUpgrade, _, err := c.GetScheduledUpgrade(clusterID)
	if err != nil || scheduledUpgrade == nil {
		return false, err
	}
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		UpgradePolicies().UpgradePolicy(scheduledUpgrade.ID()).
		Delete().
		Send()
	if err != nil {
		return false, handleErr(response.Error(), err)
	}
	return true, nil
}

func (c *Client) GetMissingGateAgreementsHypershift(
	clusterID string,
	upgradePolicy *cmv1.ControlPlaneUpgradePolicy) ([]*cmv1.VersionGate, error) {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).ControlPlane().UpgradePolicies().Add().Parameter("dryRun", true).Body(upgradePolicy).Send()

	if err != nil {
		if response.Error() != nil {
			// parse gates list
			errorDetails, ok := response.Error().GetDetails()
			if !ok {
				return []*cmv1.VersionGate{}, handleErr(response.Error(), err)
			}
			data, err := json.Marshal(errorDetails)
			if err != nil {
				return []*cmv1.VersionGate{}, handleErr(response.Error(), err)
			}
			gates, err := cmv1.UnmarshalVersionGateList(data)
			if err != nil {
				return []*cmv1.VersionGate{}, handleErr(response.Error(), err)
			}
			// return original error if invaild version gate detected
			if len(gates) > 0 && gates[0].ID() == "" {
				errType := errors.ErrorType(response.Error().Status())
				return []*cmv1.VersionGate{}, errType.Set(errors.Errorf("%v", response.Error().Reason()))
			}
			return gates, nil
		}
	}
	return []*cmv1.VersionGate{}, nil
}

func (c *Client) GetMissingGateAgreementsClassic(
	clusterID string,
	upgradePolicy *cmv1.UpgradePolicy) ([]*cmv1.VersionGate, error) {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).UpgradePolicies().Add().Parameter("dryRun", true).Body(upgradePolicy).Send()

	if err != nil {
		if response.Error() != nil {
			// parse gates list
			errorDetails, ok := response.Error().GetDetails()
			if !ok {
				return []*cmv1.VersionGate{}, handleErr(response.Error(), err)
			}
			data, err := json.Marshal(errorDetails)
			if err != nil {
				return []*cmv1.VersionGate{}, handleErr(response.Error(), err)
			}
			gates, err := cmv1.UnmarshalVersionGateList(data)
			if err != nil {
				return []*cmv1.VersionGate{}, handleErr(response.Error(), err)
			}
			// return original error if invaild version gate detected
			if len(gates) > 0 && gates[0].ID() == "" {
				errType := errors.ErrorType(response.Error().Status())
				return []*cmv1.VersionGate{}, errType.Set(errors.Errorf("%v", response.Error().Reason()))
			}
			return gates, nil
		}
	}
	return []*cmv1.VersionGate{}, nil
}

func (c *Client) AckVersionGate(
	clusterID string,
	gateID string) error {
	agreement, err := cmv1.NewVersionGateAgreement().
		VersionGate(cmv1.NewVersionGate().ID(gateID)).
		Build()
	if err != nil {
		return err
	}
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		GateAgreements().
		Add().
		Body(agreement).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}
