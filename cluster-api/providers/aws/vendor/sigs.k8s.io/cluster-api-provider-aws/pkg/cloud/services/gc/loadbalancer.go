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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func (s *Service) deleteLoadBalancers(ctx context.Context, resources []*AWSResource) error {
	for _, resource := range resources {
		if !s.isELBResourceToDelete(resource, "loadbalancer") {
			s.scope.V(5).Info("Resource not a load balancer for deletion", "arn", resource.ARN.String())
			continue
		}

		switch {
		case strings.HasPrefix(resource.ARN.Resource, "loadbalancer/app/"):
			s.scope.V(5).Info("Deleting ALB for Service", "arn", resource.ARN.String())
			if err := s.deleteLoadBalancerV2(ctx, resource.ARN.String()); err != nil {
				return fmt.Errorf("deleting ALB: %w", err)
			}
		case strings.HasPrefix(resource.ARN.Resource, "loadbalancer/net/"):
			s.scope.V(5).Info("Deleting NLB for Service", "arn", resource.ARN.String())
			if err := s.deleteLoadBalancerV2(ctx, resource.ARN.String()); err != nil {
				return fmt.Errorf("deleting NLB: %w", err)
			}
		case strings.HasPrefix(resource.ARN.Resource, "loadbalancer/"):
			name := strings.ReplaceAll(resource.ARN.Resource, "loadbalancer/", "")
			s.scope.V(5).Info("Deleting classic ELB for Service", "arn", resource.ARN.String(), "name", name)
			if err := s.deleteLoadBalancer(ctx, name); err != nil {
				return fmt.Errorf("deleting classic ELB: %w", err)
			}
		default:
			s.scope.V(4).Info("Unexpected elasticloadbalancing resource, ignoring", "arn", resource.ARN.String())
		}
	}

	s.scope.V(2).Info("Finished processing tagged resources for load balancers")

	return nil
}

func (s *Service) deleteTargetGroups(ctx context.Context, resources []*AWSResource) error {
	for _, resource := range resources {
		if !s.isELBResourceToDelete(resource, "targetgroup") {
			s.scope.V(4).Info("Resource not a target group for deletion", "arn", resource.ARN.String())
			continue
		}

		name := strings.ReplaceAll(resource.ARN.Resource, "targetgroup/", "")
		if err := s.deleteTargetGroup(ctx, resource.ARN.String()); err != nil {
			return fmt.Errorf("deleting target group %s: %w", name, err)
		}
	}
	s.scope.V(2).Info("Finished processing resources for target group deletion")

	return nil
}

func (s *Service) isELBResourceToDelete(resource *AWSResource, resourceName string) bool {
	if !s.isMatchingResource(resource, elb.ServiceName, resourceName) {
		return false
	}

	if serviceName := resource.Tags[serviceNameTag]; serviceName == "" {
		s.scope.V(5).Info("Resource wasn't created for a Service via CCM", "arn", resource.ARN.String(), "resource_name", resourceName)
		return false
	}

	return true
}

func (s *Service) deleteLoadBalancerV2(ctx context.Context, lbARN string) error {
	input := elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(lbARN),
	}

	s.scope.V(2).Info("Deleting v2 load balancer", "arn", lbARN)
	if _, err := s.elbv2Client.DeleteLoadBalancerWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting v2 load balancer: %w", err)
	}

	return nil
}

func (s *Service) deleteLoadBalancer(ctx context.Context, name string) error {
	input := elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	s.scope.V(2).Info("Deleting classic load balancer", "name", name)
	if _, err := s.elbClient.DeleteLoadBalancerWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting classic load balancer: %w", err)
	}

	return nil
}

func (s *Service) deleteTargetGroup(ctx context.Context, targetGroupARN string) error {
	input := elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(targetGroupARN),
	}

	s.scope.V(2).Info("Deleting target group", "arn", targetGroupARN)
	if _, err := s.elbv2Client.DeleteTargetGroupWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting target group: %w", err)
	}

	return nil
}
