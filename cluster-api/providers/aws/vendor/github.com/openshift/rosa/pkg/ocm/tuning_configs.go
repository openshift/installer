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

func (c *Client) GetTuningConfigs(clusterID string) ([]*cmv1.TuningConfig, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		TuningConfigs().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Items().Slice(), nil
}

func (c *Client) CreateTuningConfig(clusterID string, tuningConfig *cmv1.TuningConfig) (*cmv1.TuningConfig, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		TuningConfigs().
		Add().Body(tuningConfig).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) UpdateTuningConfig(clusterID string, tuningConfig *cmv1.TuningConfig) (*cmv1.TuningConfig, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		TuningConfigs().TuningConfig(tuningConfig.ID()).
		Update().Body(tuningConfig).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) DeleteTuningConfig(clusterID string, tuningConfigID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		TuningConfigs().TuningConfig(tuningConfigID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func (c *Client) FindTuningConfigByName(clusterID string, tuningConfigName string) (*cmv1.TuningConfig, error) {
	// TODO search directly by name instead of post processing. Currently, search is not exposed by the backend.
	tuningConfigs, err := c.GetTuningConfigs(clusterID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tuning configs for cluster '%s': %v", clusterID, err)
	}

	var tuningConfig *cmv1.TuningConfig
	for _, item := range tuningConfigs {
		if item != nil && item.Name() == tuningConfigName {
			tuningConfig = item
			break
		}
	}
	if tuningConfig == nil {
		return nil, fmt.Errorf("Tuning config '%s' does not exist on cluster '%s'", tuningConfigName, clusterID)
	}
	return tuningConfig, nil
}

func (c *Client) GetTuningConfigsName(clusterID string) ([]string, error) {
	var tuningConfigsNames []string
	tuningConfigs, err := c.GetTuningConfigs(clusterID)
	if err != nil {
		return tuningConfigsNames, err
	}
	for i := range tuningConfigs {
		if tuningConfigs[i] != nil {
			tuningConfigsNames = append(tuningConfigsNames, tuningConfigs[i].Name())
		}
	}
	return tuningConfigsNames, nil
}
