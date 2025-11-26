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

package network

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func (s *Service) reconcileEgressOnlyInternetGateways() error {
	if !s.scope.VPC().IsIPv6Enabled() {
		s.scope.Trace("Skipping egress only internet gateways reconcile in not ipv6 mode")
		return nil
	}

	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping egress only internet gateway reconcile in unmanaged mode")
		return nil
	}

	s.scope.Debug("Reconciling egress only internet gateways")

	eigws, err := s.describeEgressOnlyVpcInternetGateways()
	if awserrors.IsNotFound(err) {
		if !s.scope.VPC().IsIPv6Enabled() {
			return errors.Errorf("failed to validate network: no egress only internet gateways found in VPC %q", s.scope.VPC().ID)
		}

		ig, err := s.createEgressOnlyInternetGateway()
		if err != nil {
			return err
		}
		eigws = []types.EgressOnlyInternetGateway{*ig}
	} else if err != nil {
		return err
	}

	gateway := eigws[0]
	s.scope.VPC().IPv6.EgressOnlyInternetGatewayID = gateway.EgressOnlyInternetGatewayId

	// Make sure tags are up to date.
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		buildParams := s.getEgressOnlyGatewayTagParams(*gateway.EgressOnlyInternetGatewayId)
		tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
		if err := tagsBuilder.Ensure(converters.TagsToMap(gateway.Tags)); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.EgressOnlyInternetGatewayNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedTagEgressOnlyInternetGateway", "Failed to tag managed Egress Only Internet Gateway %q: %v", gateway.EgressOnlyInternetGatewayId, err)
		return errors.Wrapf(err, "failed to tag egress only internet gateway %q", *gateway.EgressOnlyInternetGatewayId)
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition)
	return nil
}

func (s *Service) deleteEgressOnlyInternetGateways() error {
	if !s.scope.VPC().IsIPv6Enabled() {
		s.scope.Trace("Skipping egress only internet gateway deletion in none ipv6 mode")
		return nil
	}

	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping egress only internet gateway deletion in unmanaged mode")
		return nil
	}

	eigws, err := s.describeEgressOnlyVpcInternetGateways()
	if awserrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, ig := range eigws {
		deleteReq := &ec2.DeleteEgressOnlyInternetGatewayInput{
			EgressOnlyInternetGatewayId: ig.EgressOnlyInternetGatewayId,
		}

		if _, err = s.EC2Client.DeleteEgressOnlyInternetGateway(context.TODO(), deleteReq); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedDeleteEgressOnlyInternetGateway", "Failed to delete Egress Only Internet Gateway %q previously attached to VPC %q: %v", *ig.EgressOnlyInternetGatewayId, s.scope.VPC().ID, err)
			return errors.Wrapf(err, "failed to delete egress only internet gateway %q", *ig.EgressOnlyInternetGatewayId)
		}

		record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteEgressOnlyInternetGateway", "Deleted Egress Only Internet Gateway %q previously attached to VPC %q", *ig.EgressOnlyInternetGatewayId, s.scope.VPC().ID)
		s.scope.Info("Deleted Egress Only Internet gateway in VPC", "egress-only-internet-gateway-id", *ig.EgressOnlyInternetGatewayId, "vpc-id", s.scope.VPC().ID)
	}

	return nil
}

func (s *Service) createEgressOnlyInternetGateway() (*types.EgressOnlyInternetGateway, error) {
	ig, err := s.EC2Client.CreateEgressOnlyInternetGateway(context.TODO(), &ec2.CreateEgressOnlyInternetGatewayInput{
		TagSpecifications: []types.TagSpecification{
			tags.BuildParamsToTagSpecification(types.ResourceTypeEgressOnlyInternetGateway, s.getEgressOnlyGatewayTagParams(services.TemporaryResourceID)),
		},
		VpcId: aws.String(s.scope.VPC().ID),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateEgressOnlyInternetGateway", "Failed to create new managed Egress Only Internet Gateway: %v", err)
		return nil, errors.Wrap(err, "failed to create egress only internet gateway")
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateEgressOnlyInternetGateway", "Created new managed Egress Only Internet Gateway %q", *ig.EgressOnlyInternetGateway.EgressOnlyInternetGatewayId)
	s.scope.Info("Created Egress Only Internet gateway", "egress-only-internet-gateway-id", *ig.EgressOnlyInternetGateway.EgressOnlyInternetGatewayId, "vpc-id", s.scope.VPC().ID)

	return ig.EgressOnlyInternetGateway, nil
}

func (s *Service) describeEgressOnlyVpcInternetGateways() ([]types.EgressOnlyInternetGateway, error) {
	// The API for DescribeEgressOnlyInternetGateways does not support filtering by VPC ID attachment.
	// More details: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeEgressOnlyInternetGateways.html
	// Since the eigw is managed by CAPA, we can filter by the kubernetes cluster tag.
	out, err := s.EC2Client.DescribeEgressOnlyInternetGateways(context.TODO(), &ec2.DescribeEgressOnlyInternetGatewaysInput{
		Filters: []types.Filter{
			filter.EC2.Cluster(s.scope.Name()),
		},
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeEgressOnlyInternetGateway", "Failed to describe egress only internet gateway in vpc %q: %v", s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to describe egress only internet gateways in vpc %q", s.scope.VPC().ID)
	}

	// For safeguarding, we collect only egress-only internet gateways
	// that are attached to the VPC.
	eigws := make([]types.EgressOnlyInternetGateway, 0)
	for _, eigw := range out.EgressOnlyInternetGateways {
		for _, attachment := range eigw.Attachments {
			if aws.ToString(attachment.VpcId) == s.scope.VPC().ID {
				eigws = append(eigws, eigw)
			}
		}
	}

	if len(eigws) == 0 {
		return nil, awserrors.NewNotFound(fmt.Sprintf("no egress only internet gateways found in vpc %q", s.scope.VPC().ID))
	} else if len(eigws) > 1 {
		eigwIDs := make([]string, len(eigws))
		for i, eigw := range eigws {
			eigwIDs[i] = aws.ToString(eigw.EgressOnlyInternetGatewayId)
		}
		return nil, awserrors.NewConflict(fmt.Sprintf("expected 1 egress only internet gateway in vpc %q, but found %v: %v", s.scope.VPC().ID, len(eigws), eigwIDs))
	}

	return eigws, nil
}

func (s *Service) getEgressOnlyGatewayTagParams(id string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-eigw", s.scope.Name())

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}
