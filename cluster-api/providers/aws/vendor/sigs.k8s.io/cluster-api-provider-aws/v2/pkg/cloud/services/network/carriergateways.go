/*
Copyright 2024 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

func (s *Service) reconcileCarrierGateway() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping carrier gateway reconcile in unmanaged mode")
		return nil
	}

	if !s.scope.Subnets().HasPublicSubnetWavelength() {
		s.scope.Trace("Skipping carrier gateway reconcile in VPC without subnets in zone type wavelength-zone")
		return nil
	}

	s.scope.Debug("Reconciling carrier gateway")

	cagw, err := s.describeVpcCarrierGateway()
	if awserrors.IsNotFound(err) {
		if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
			return errors.Errorf("failed to validate network: no carrier gateway found in VPC %q", s.scope.VPC().ID)
		}

		cg, err := s.createCarrierGateway()
		if err != nil {
			return err
		}
		cagw = cg
	} else if err != nil {
		return err
	}

	s.scope.VPC().CarrierGatewayID = cagw.CarrierGatewayId

	// Make sure tags are up-to-date.
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		buildParams := s.getGatewayTagParams(*cagw.CarrierGatewayId)
		tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
		if err := tagsBuilder.Ensure(converters.TagsToMap(cagw.Tags)); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.InvalidCarrierGatewayNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedTagCarrierGateway", "Failed to tag managed Carrier Gateway %q: %v", cagw.CarrierGatewayId, err)
		return errors.Wrapf(err, "failed to tag carrier gateway %q", *cagw.CarrierGatewayId)
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition)
	return nil
}

func (s *Service) deleteCarrierGateway() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping carrier gateway deletion in unmanaged mode")
		return nil
	}

	cagw, err := s.describeVpcCarrierGateway()
	if awserrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	deleteReq := &ec2.DeleteCarrierGatewayInput{
		CarrierGatewayId: cagw.CarrierGatewayId,
	}

	if _, err = s.EC2Client.DeleteCarrierGatewayWithContext(context.TODO(), deleteReq); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedDeleteCarrierGateway", "Failed to delete Carrier Gateway %q previously attached to VPC %q: %v", *cagw.CarrierGatewayId, s.scope.VPC().ID, err)
		return errors.Wrapf(err, "failed to delete carrier gateway %q", *cagw.CarrierGatewayId)
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteCarrierGateway", "Deleted Carrier Gateway %q previously attached to VPC %q", *cagw.CarrierGatewayId, s.scope.VPC().ID)
	s.scope.Info("Deleted Carrier Gateway in VPC", "carrier-gateway-id", *cagw.CarrierGatewayId, "vpc-id", s.scope.VPC().ID)

	return nil
}

func (s *Service) createCarrierGateway() (*ec2.CarrierGateway, error) {
	ig, err := s.EC2Client.CreateCarrierGatewayWithContext(context.TODO(), &ec2.CreateCarrierGatewayInput{
		VpcId: aws.String(s.scope.VPC().ID),
		TagSpecifications: []*ec2.TagSpecification{
			tags.BuildParamsToTagSpecification(ec2.ResourceTypeCarrierGateway, s.getGatewayTagParams(services.TemporaryResourceID)),
		},
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateCarrierGateway", "Failed to create new managed Internet Gateway: %v", err)
		return nil, errors.Wrap(err, "failed to create carrier gateway")
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateCarrierGateway", "Created new managed Internet Gateway %q", *ig.CarrierGateway.CarrierGatewayId)
	s.scope.Info("Created Internet gateway for VPC", "internet-gateway-id", *ig.CarrierGateway.CarrierGatewayId, "vpc-id", s.scope.VPC().ID)

	return ig.CarrierGateway, nil
}

func (s *Service) describeVpcCarrierGateway() (*ec2.CarrierGateway, error) {
	out, err := s.EC2Client.DescribeCarrierGatewaysWithContext(context.TODO(), &ec2.DescribeCarrierGatewaysInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(s.scope.VPC().ID),
		},
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeCarrierGateway", "Failed to describe carrier gateways in vpc %q: %v", s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to describe carrier gateways in vpc %q", s.scope.VPC().ID)
	}

	if len(out.CarrierGateways) == 0 {
		return nil, awserrors.NewNotFound(fmt.Sprintf("no carrier gateways found in vpc %q", s.scope.VPC().ID))
	}

	return out.CarrierGateways[0], nil
}
