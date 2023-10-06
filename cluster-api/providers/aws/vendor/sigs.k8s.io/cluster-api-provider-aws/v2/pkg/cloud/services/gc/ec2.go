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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
)

func (s *Service) deleteSecurityGroups(ctx context.Context, resources []*AWSResource) error {
	for _, resource := range resources {
		if !s.isSecurityGroupToDelete(resource) {
			s.scope.Debug("Resource not a security group for deletion", "arn", resource.ARN.String())
			continue
		}

		groupID := strings.ReplaceAll(resource.ARN.Resource, "security-group/", "")
		if err := s.deleteSecurityGroup(ctx, groupID); err != nil {
			return fmt.Errorf("deleting security group %s: %w", groupID, err)
		}
	}
	s.scope.Debug("Finished processing resources for security group deletion")

	return nil
}

func (s *Service) isSecurityGroupToDelete(resource *AWSResource) bool {
	if !s.isMatchingResource(resource, ec2.ServiceName, "security-group") {
		return false
	}
	if eksClusterName := resource.Tags[eksClusterNameTag]; eksClusterName != "" {
		s.scope.Debug("Security group was created by EKS directly", "arn", resource.ARN.String(), "check", "securitygroup", "cluster_name", eksClusterName)
		return false
	}
	s.scope.Debug("Resource is a security group to delete", "arn", resource.ARN.String(), "check", "securitygroup")

	return true
}

func (s *Service) deleteSecurityGroup(ctx context.Context, securityGroupID string) error {
	input := ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(securityGroupID),
	}

	s.scope.Debug("Deleting security group", "group_id", securityGroupID)
	if _, err := s.ec2Client.DeleteSecurityGroupWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting security group: %w", err)
	}

	return nil
}

// getProviderOwnedSecurityGroups gets cloud provider created security groups of ELBs for this cluster, filtering by tag: kubernetes.io/cluster/<cluster-name>:owned and VPC Id.
func (s *Service) getProviderOwnedSecurityGroups(ctx context.Context) ([]*AWSResource, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			filter.EC2.ProviderOwned(s.scope.KubernetesClusterName()),
		},
	}

	var resources []*AWSResource
	err := s.ec2Client.DescribeSecurityGroupsPages(input, func(out *ec2.DescribeSecurityGroupsOutput, last bool) bool {
		for _, group := range out.SecurityGroups {
			arn := composeFakeArn(sgService, sgResourcePrefix+*group.GroupId)
			resource, err := composeAWSResource(arn, converters.TagsToMap(group.Tags))
			if err != nil {
				s.scope.Error(err, "error compose aws security group resource: %v", "name", arn)
				continue
			}
			resources = append(resources, resource)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("describe security groups error: %w", err)
	}

	return resources, nil
}
