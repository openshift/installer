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
	"net/http"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

type AutoscalerConfig struct {
	BalanceSimilarNodeGroups    bool
	SkipNodesWithLocalStorage   bool
	LogVerbosity                int
	MaxPodGracePeriod           int
	PodPriorityThreshold        int
	IgnoreDaemonsetsUtilization bool
	MaxNodeProvisionTime        string
	BalancingIgnoredLabels      []string
	ResourceLimits              ResourceLimits
	ScaleDown                   ScaleDownConfig
}

type ResourceLimits struct {
	MaxNodesTotal int
	Cores         ResourceRange
	Memory        ResourceRange
	GPULimits     []GPULimit
}

type ResourceRange struct {
	Min int
	Max int
}

type GPULimit struct {
	Type  string
	Range ResourceRange
}

type ScaleDownConfig struct {
	Enabled              bool
	UnneededTime         string
	UtilizationThreshold float64
	DelayAfterAdd        string
	DelayAfterDelete     string
	DelayAfterFailure    string
}

func BuildClusterAutoscalerHostedCp(config *AutoscalerConfig) *cmv1.ClusterAutoscalerBuilder {
	return cmv1.NewClusterAutoscaler().
		MaxPodGracePeriod(config.MaxPodGracePeriod).
		PodPriorityThreshold(config.PodPriorityThreshold).
		MaxNodeProvisionTime(config.MaxNodeProvisionTime).
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().
			MaxNodesTotal(config.ResourceLimits.MaxNodesTotal))
}

func BuildClusterAutoscaler(config *AutoscalerConfig) *cmv1.ClusterAutoscalerBuilder {
	if config == nil {
		return nil
	}

	gpuLimits := []*cmv1.AutoscalerResourceLimitsGPULimitBuilder{}
	for _, gpuLimit := range config.ResourceLimits.GPULimits {
		gpuLimits = append(
			gpuLimits,
			cmv1.NewAutoscalerResourceLimitsGPULimit().
				Type(gpuLimit.Type).
				Range(cmv1.NewResourceRange().
					Min(gpuLimit.Range.Min).
					Max(gpuLimit.Range.Max)),
		)
	}

	return cmv1.NewClusterAutoscaler().
		BalanceSimilarNodeGroups(config.BalanceSimilarNodeGroups).
		SkipNodesWithLocalStorage(config.SkipNodesWithLocalStorage).
		LogVerbosity(config.LogVerbosity).
		MaxPodGracePeriod(config.MaxPodGracePeriod).
		PodPriorityThreshold(config.PodPriorityThreshold).
		IgnoreDaemonsetsUtilization(config.IgnoreDaemonsetsUtilization).
		MaxNodeProvisionTime(config.MaxNodeProvisionTime).
		BalancingIgnoredLabels(config.BalancingIgnoredLabels...).
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().
			MaxNodesTotal(config.ResourceLimits.MaxNodesTotal).
			Cores(cmv1.NewResourceRange().
				Min(config.ResourceLimits.Cores.Min).
				Max(config.ResourceLimits.Cores.Max)).
			Memory(cmv1.NewResourceRange().
				Min(config.ResourceLimits.Memory.Min).
				Max(config.ResourceLimits.Memory.Max)).
			GPUS(gpuLimits...)).
		ScaleDown(cmv1.NewAutoscalerScaleDownConfig().
			Enabled(config.ScaleDown.Enabled).
			UnneededTime(config.ScaleDown.UnneededTime).
			UtilizationThreshold(fmt.Sprintf("%f", config.ScaleDown.UtilizationThreshold)).
			DelayAfterAdd(config.ScaleDown.DelayAfterAdd).
			DelayAfterDelete(config.ScaleDown.DelayAfterDelete).
			DelayAfterFailure(config.ScaleDown.DelayAfterFailure))
}

func (c *Client) GetClusterAutoscaler(clusterID string) (*cmv1.ClusterAutoscaler, error) {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).Autoscaler().Get().Send()

	if response.Status() == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *Client) CreateClusterAutoscaler(clusterId string, config *AutoscalerConfig) (*cmv1.ClusterAutoscaler, error) {
	object, err := BuildClusterAutoscaler(config).Build()

	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().Post().Request(object).Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) UpdateClusterAutoscaler(clusterId string, isHostedCp bool,
	config *AutoscalerConfig) (*cmv1.ClusterAutoscaler, error) {

	var object *cmv1.ClusterAutoscaler
	var err error

	if isHostedCp {
		object, err = BuildClusterAutoscalerHostedCp(config).Build()
	} else {
		object, err = BuildClusterAutoscaler(config).Build()
	}

	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().Update().Body(object).Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) DeleteClusterAutoscaler(clusterID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Autoscaler().
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}
