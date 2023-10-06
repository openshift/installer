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

func (s *Service) reconcileInternetGateways() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping internet gateways reconcile in unmanaged mode")
		return nil
	}

	s.scope.Debug("Reconciling internet gateways")

	igs, err := s.describeVpcInternetGateways()
	if awserrors.IsNotFound(err) {
		if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
			return errors.Errorf("failed to validate network: no internet gateways found in VPC %q", s.scope.VPC().ID)
		}

		ig, err := s.createInternetGateway()
		if err != nil {
			return err
		}
		igs = []*ec2.InternetGateway{ig}
	} else if err != nil {
		return err
	}

	gateway := igs[0]
	s.scope.VPC().InternetGatewayID = gateway.InternetGatewayId

	// Make sure tags are up-to-date.
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		buildParams := s.getGatewayTagParams(*gateway.InternetGatewayId)
		tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
		if err := tagsBuilder.Ensure(converters.TagsToMap(gateway.Tags)); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.InternetGatewayNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedTagInternetGateway", "Failed to tag managed Internet Gateway %q: %v", gateway.InternetGatewayId, err)
		return errors.Wrapf(err, "failed to tag internet gateway %q", *gateway.InternetGatewayId)
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition)
	return nil
}

func (s *Service) deleteInternetGateways() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping internet gateway deletion in unmanaged mode")
		return nil
	}

	igs, err := s.describeVpcInternetGateways()
	if awserrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, ig := range igs {
		detachReq := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
			VpcId:             aws.String(s.scope.VPC().ID),
		}

		if _, err := s.EC2Client.DetachInternetGateway(detachReq); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedDetachInternetGateway", "Failed to detach Internet Gateway %q from VPC %q: %v", *ig.InternetGatewayId, s.scope.VPC().ID, err)
			return errors.Wrapf(err, "failed to detach internet gateway %q", *ig.InternetGatewayId)
		}

		record.Eventf(s.scope.InfraCluster(), "SuccessfulDetachInternetGateway", "Detached Internet Gateway %q from VPC %q", *ig.InternetGatewayId, s.scope.VPC().ID)
		s.scope.Debug("Detached internet gateway from VPC", "internet-gateway-id", *ig.InternetGatewayId, "vpc-id", s.scope.VPC().ID)

		deleteReq := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
		}

		if _, err = s.EC2Client.DeleteInternetGateway(deleteReq); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedDeleteInternetGateway", "Failed to delete Internet Gateway %q previously attached to VPC %q: %v", *ig.InternetGatewayId, s.scope.VPC().ID, err)
			return errors.Wrapf(err, "failed to delete internet gateway %q", *ig.InternetGatewayId)
		}

		record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteInternetGateway", "Deleted Internet Gateway %q previously attached to VPC %q", *ig.InternetGatewayId, s.scope.VPC().ID)
		s.scope.Info("Deleted Internet gateway in VPC", "internet-gateway-id", *ig.InternetGatewayId, "vpc-id", s.scope.VPC().ID)
	}

	return nil
}

func (s *Service) createInternetGateway() (*ec2.InternetGateway, error) {
	ig, err := s.EC2Client.CreateInternetGateway(&ec2.CreateInternetGatewayInput{
		TagSpecifications: []*ec2.TagSpecification{
			tags.BuildParamsToTagSpecification(ec2.ResourceTypeInternetGateway, s.getGatewayTagParams(services.TemporaryResourceID)),
		},
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateInternetGateway", "Failed to create new managed Internet Gateway: %v", err)
		return nil, errors.Wrap(err, "failed to create internet gateway")
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateInternetGateway", "Created new managed Internet Gateway %q", *ig.InternetGateway.InternetGatewayId)
	s.scope.Info("Created Internet gateway for VPC", "internet-gateway-id", *ig.InternetGateway.InternetGatewayId, "vpc-id", s.scope.VPC().ID)

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.EC2Client.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
			InternetGatewayId: ig.InternetGateway.InternetGatewayId,
			VpcId:             aws.String(s.scope.VPC().ID),
		}); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.InternetGatewayNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAttachInternetGateway", "Failed to attach managed Internet Gateway %q to vpc %q: %v", *ig.InternetGateway.InternetGatewayId, s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to attach internet gateway %q to vpc %q", *ig.InternetGateway.InternetGatewayId, s.scope.VPC().ID)
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulAttachInternetGateway", "Internet Gateway %q attached to VPC %q", *ig.InternetGateway.InternetGatewayId, s.scope.VPC().ID)
	s.scope.Debug("attached internet gateway to VPC", "internet-gateway-id", *ig.InternetGateway.InternetGatewayId, "vpc-id", s.scope.VPC().ID)

	return ig.InternetGateway, nil
}

func (s *Service) describeVpcInternetGateways() ([]*ec2.InternetGateway, error) {
	out, err := s.EC2Client.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPCAttachment(s.scope.VPC().ID),
		},
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeInternetGateway", "Failed to describe internet gateways in vpc %q: %v", s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to describe internet gateways in vpc %q", s.scope.VPC().ID)
	}

	if len(out.InternetGateways) == 0 {
		return nil, awserrors.NewNotFound(fmt.Sprintf("no internet gateways found in vpc %q", s.scope.VPC().ID))
	}

	return out.InternetGateways, nil
}

func (s *Service) getGatewayTagParams(id string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-igw", s.scope.Name())

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}
