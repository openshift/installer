/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
)

func (s *Service) reconcileSecurityGroups(ctx context.Context, cluster *ekstypes.Cluster) error {
	s.scope.Info("Reconciling EKS security groups", "cluster-name", ptr.Deref(cluster.Name, ""))

	if s.scope.Network().SecurityGroups == nil {
		s.scope.Network().SecurityGroups = make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
	}

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:aws:eks:cluster-name"),
				Values: []*string{cluster.Name},
			},
		},
	}

	output, err := s.EC2Client.DescribeSecurityGroupsWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("describing security groups: %w", err)
	}

	if len(output.SecurityGroups) == 0 {
		return ErrNoSecurityGroup
	}

	sg := infrav1.SecurityGroup{
		ID:   *output.SecurityGroups[0].GroupId,
		Name: *output.SecurityGroups[0].GroupName,
		Tags: converters.TagsToMap(output.SecurityGroups[0].Tags),
	}
	s.scope.ControlPlane.Status.Network.SecurityGroups[infrav1.SecurityGroupNode] = sg

	input = &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			cluster.ResourcesVpcConfig.ClusterSecurityGroupId,
		},
	}

	output, err = s.EC2Client.DescribeSecurityGroupsWithContext(ctx, input)
	if err != nil || len(output.SecurityGroups) == 0 {
		return fmt.Errorf("describing EKS cluster security group: %w", err)
	}

	s.scope.ControlPlane.Status.Network.SecurityGroups[ekscontrolplanev1.SecurityGroupCluster] = infrav1.SecurityGroup{
		ID:   aws.ToString(cluster.ResourcesVpcConfig.ClusterSecurityGroupId),
		Name: *output.SecurityGroups[0].GroupName,
		Tags: converters.TagsToMap(output.SecurityGroups[0].Tags),
	}

	return nil
}
