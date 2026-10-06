/*
Copyright 2025 The Kubernetes Authors.

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

package network

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
)

// deleteOrphanedENIs deletes network interfaces in the cluster's VPC that are
// in "available" state (not attached to any instance or service) and reference
// one of the cluster's security groups. These are typically left behind by NLBs
// whose cross-AZ ENIs are cleaned up asynchronously by AWS after LB deletion.
func (s *Service) deleteOrphanedENIs(ctx context.Context) error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping orphaned ENI cleanup in unmanaged mode")
		return nil
	}

	sgIDs := s.clusterSecurityGroupIDs()
	if len(sgIDs) == 0 {
		s.scope.Debug("No cluster security groups found, skipping orphaned ENI cleanup")
		return nil
	}

	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{s.scope.VPC().ID},
			},
			{
				Name:   aws.String("group-id"),
				Values: sgIDs,
			},
			{
				Name:   aws.String("status"),
				Values: []string{"available"},
			},
		},
	}

	out, err := s.EC2Client.DescribeNetworkInterfaces(ctx, input)
	if err != nil {
		return fmt.Errorf("describing orphaned network interfaces: %w", err)
	}

	if len(out.NetworkInterfaces) == 0 {
		return nil
	}

	s.scope.Info("Deleting orphaned network interfaces", "count", len(out.NetworkInterfaces))

	var errs []error
	for _, eni := range out.NetworkInterfaces {
		eniID := aws.ToString(eni.NetworkInterfaceId)
		s.scope.Debug("Deleting orphaned network interface", "eni-id", eniID)

		if _, err := s.EC2Client.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: eni.NetworkInterfaceId,
		}); err != nil {
			errs = append(errs, fmt.Errorf("deleting network interface %q: %w", eniID, err))
			continue
		}

		s.scope.Info("Deleted orphaned network interface", "eni-id", eniID)
	}

	return kerrors.NewAggregate(errs)
}

func (s *Service) clusterSecurityGroupIDs() []string {
	sgs := s.scope.SecurityGroups()
	ids := make([]string, 0, len(sgs))
	for _, sg := range sgs {
		if sg.ID != "" {
			ids = append(ids, sg.ID)
		}
	}
	return ids
}
