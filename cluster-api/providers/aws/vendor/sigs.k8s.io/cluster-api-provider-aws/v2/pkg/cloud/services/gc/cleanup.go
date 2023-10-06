/*
Copyright 2022 The Kubernetes Authors.

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

package gc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/annotations"
)

const (
	serviceNameTag    = "kubernetes.io/service-name"
	eksClusterNameTag = "aws:eks:cluster-name"
)

// ReconcileDelete is responsible for determining if the infra cluster needs to be garbage collected. If
// does then it will perform garbage collection. For example, it will delete the ELB/NLBs that where created
// as a result of Services of type load balancer.
func (s *Service) ReconcileDelete(ctx context.Context) error {
	s.scope.Info("reconciling deletion for garbage collection", "cluster", s.scope.InfraClusterName())

	val, found := annotations.Get(s.scope.InfraCluster(), expinfrav1.ExternalResourceGCAnnotation)
	if !found {
		val = "true"
	}

	shouldGC, err := strconv.ParseBool(val)
	if err != nil {
		return fmt.Errorf("converting value %s of annotation %s to bool: %w", val, expinfrav1.ExternalResourceGCAnnotation, err)
	}

	if !shouldGC {
		s.scope.Info("cluster opted-out of garbage collection")

		return nil
	}

	return s.deleteResources(ctx)
}

func (s *Service) deleteResources(ctx context.Context) error {
	s.scope.Info("deleting aws resources created by tenant cluster", "cluster", s.scope.InfraClusterName())

	resources, err := s.collectFuncs.Execute(ctx)
	if err != nil {
		return fmt.Errorf("collecting resources: %w", err)
	}

	if deleteErr := s.cleanupFuncs.Execute(ctx, resources); deleteErr != nil {
		return fmt.Errorf("deleting resources: %w", deleteErr)
	}

	return nil
}

func (s *Service) defaultGetResources(ctx context.Context) ([]*AWSResource, error) {
	s.scope.Info("get aws resources created by tenant cluster with resource group tagging API", "cluster", s.scope.InfraClusterName())

	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())

	awsInput := rgapi.GetResourcesInput{
		ResourceTypeFilters: nil,
		TagFilters: []*rgapi.TagFilter{
			{
				Key:    aws.String(serviceTag),
				Values: []*string{aws.String(string(infrav1.ResourceLifecycleOwned))},
			},
		},
	}

	awsOutput, err := s.resourceTaggingClient.GetResourcesWithContext(ctx, &awsInput)
	if err != nil {
		return nil, fmt.Errorf("getting tagged resources: %w", err)
	}

	resources := []*AWSResource{}

	for i := range awsOutput.ResourceTagMappingList {
		mapping := awsOutput.ResourceTagMappingList[i]
		parsedArn, err := arn.Parse(*mapping.ResourceARN)
		if err != nil {
			return nil, fmt.Errorf("parsing resource arn %s: %w", *mapping.ResourceARN, err)
		}

		tags := map[string]string{}
		for _, rgTag := range mapping.Tags {
			tags[*rgTag.Key] = *rgTag.Value
		}

		resources = append(resources, &AWSResource{
			ARN:  &parsedArn,
			Tags: tags,
		})
	}

	return resources, nil
}

func (s *Service) isMatchingResource(resource *AWSResource, serviceName, resourceName string) bool {
	if resource.ARN.Service != serviceName {
		s.scope.Debug("Resource not for service", "arn", resource.ARN.String(), "service_name", serviceName, "resource_name", resourceName)
		return false
	}
	if !strings.HasPrefix(resource.ARN.Resource, resourceName+"/") {
		s.scope.Debug("Resource type does not match", "arn", resource.ARN.String(), "service_name", serviceName, "resource_name", resourceName)
		return false
	}

	return true
}
