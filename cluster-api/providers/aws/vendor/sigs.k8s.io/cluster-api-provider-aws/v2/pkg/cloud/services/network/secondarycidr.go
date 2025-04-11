/*
Copyright 2018 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func isVPCPresent(vpcs *ec2.DescribeVpcsOutput) bool {
	return vpcs != nil && len(vpcs.Vpcs) > 0
}

func (s *Service) associateSecondaryCidrs() error {
	secondaryCidrBlocks := s.scope.AllSecondaryCidrBlocks()
	if len(secondaryCidrBlocks) == 0 {
		return nil
	}

	vpcs, err := s.EC2Client.DescribeVpcsWithContext(context.TODO(), &ec2.DescribeVpcsInput{
		VpcIds: []*string{&s.scope.VPC().ID},
	})
	if err != nil {
		return err
	}

	if !isVPCPresent(vpcs) {
		return errors.Errorf("failed to associateSecondaryCidr as there are no VPCs present")
	}

	// We currently only *add* associations. Here, we do not reconcile exactly against the provided list
	// (i.e. disassociate what isn't listed in the spec).
	existingAssociations := vpcs.Vpcs[0].CidrBlockAssociationSet
	for _, desiredCidrBlock := range secondaryCidrBlocks {
		found := false
		for _, existing := range existingAssociations {
			if *existing.CidrBlock == desiredCidrBlock.IPv4CidrBlock {
				found = true
				break
			}
		}
		if found {
			continue
		}

		out, err := s.EC2Client.AssociateVpcCidrBlockWithContext(context.TODO(), &ec2.AssociateVpcCidrBlockInput{
			VpcId:     &s.scope.VPC().ID,
			CidrBlock: &desiredCidrBlock.IPv4CidrBlock,
		})
		if err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedAssociateSecondaryCidr", "Failed associating secondary CIDR %q with VPC %v", desiredCidrBlock.IPv4CidrBlock, err)
			return err
		}

		// Once IPv6 is supported, we need to consider both `out.CidrBlockAssociation.AssociationId` and
		// `out.Ipv6CidrBlockAssociation.AssociationId`
		record.Eventf(s.scope.InfraCluster(), "SuccessfulAssociateSecondaryCidr", "Associated secondary CIDR %q with VPC %q", desiredCidrBlock.IPv4CidrBlock, *out.CidrBlockAssociation.AssociationId)
	}

	return nil
}

func (s *Service) disassociateSecondaryCidrs() error {
	// If the VPC is unmanaged or not yet populated, return early.
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) || s.scope.VPC().ID == "" {
		return nil
	}

	secondaryCidrBlocks := s.scope.AllSecondaryCidrBlocks()
	if len(secondaryCidrBlocks) == 0 {
		return nil
	}

	vpcs, err := s.EC2Client.DescribeVpcsWithContext(context.TODO(), &ec2.DescribeVpcsInput{
		VpcIds: []*string{&s.scope.VPC().ID},
	})
	if err != nil {
		return err
	}

	if !isVPCPresent(vpcs) {
		return errors.Errorf("failed to disassociateSecondaryCidr as there are no VPCs present")
	}

	existingAssociations := vpcs.Vpcs[0].CidrBlockAssociationSet
	for _, cidrBlockToDelete := range secondaryCidrBlocks {
		for _, existing := range existingAssociations {
			if *existing.CidrBlock == cidrBlockToDelete.IPv4CidrBlock {
				if _, err := s.EC2Client.DisassociateVpcCidrBlockWithContext(context.TODO(), &ec2.DisassociateVpcCidrBlockInput{
					AssociationId: existing.AssociationId,
				}); err != nil {
					record.Warnf(s.scope.InfraCluster(), "FailedDisassociateSecondaryCidr", "Failed disassociating secondary CIDR %q from VPC %v", cidrBlockToDelete.IPv4CidrBlock, err)
					return err
				}
				break
			}
		}
	}

	return nil
}
