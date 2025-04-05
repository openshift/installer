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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const (
	defaultVPCCidr             = "10.0.0.0/16"
	defaultIpamV4NetmaskLength = 16
	defaultIpamV6NetmaskLength = 56
)

func (s *Service) reconcileVPC() error {
	s.scope.Debug("Reconciling VPC")

	// If the ID is not nil, VPC is either managed or unmanaged but should exist in the AWS.
	if s.scope.VPC().ID != "" {
		vpc, err := s.describeVPCByID()
		if err != nil {
			return errors.Wrap(err, ".spec.vpc.id is set but VPC resource is missing in AWS; failed to describe VPC resources. (might be in creation process)")
		}

		s.scope.VPC().CidrBlock = vpc.CidrBlock
		if s.scope.VPC().IsIPv6Enabled() {
			s.scope.VPC().IPv6 = vpc.IPv6
		}
		if s.scope.TagUnmanagedNetworkResources() {
			s.scope.VPC().Tags = vpc.Tags
		}

		// If VPC is unmanaged, return early.
		if vpc.IsUnmanaged(s.scope.Name()) {
			s.scope.Debug("Working on unmanaged VPC", "vpc-id", vpc.ID)
			if err := s.scope.PatchObject(); err != nil {
				return errors.Wrap(err, "failed to patch unmanaged VPC fields")
			}
			record.Eventf(s.scope.InfraCluster(), "SuccessfulSetVPCAttributes", "Set managed VPC attributes for %q", vpc.ID)
			return nil
		}

		if !s.scope.TagUnmanagedNetworkResources() {
			s.scope.VPC().Tags = vpc.Tags
		}

		// Make sure tags are up-to-date.
		// **Only** do this for managed VPCs. Make sure this logic is below the above `vpc.IsUnmanaged` check.
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			buildParams := s.getVPCTagParams(s.scope.VPC().ID)
			tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
			if err := tagsBuilder.Ensure(s.scope.VPC().Tags); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.VPCNotFound); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedTagVPC", "Failed ensure managed VPC %q: %v", s.scope.VPC().ID, err)
			return errors.Wrapf(err, "failed to ensure tags on vpc %q", s.scope.VPC().ID)
		}

		// if the VPC is managed, make managed sure attributes are configured.
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if err := s.ensureManagedVPCAttributes(vpc); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.VPCNotFound); err != nil {
			return errors.Wrapf(err, "failed to set vpc attributes for %q", vpc.ID)
		}

		return nil
	}

	// .spec.vpc.id is nil. This means no managed VPC exists or we failed to save its ID before. Check if a managed VPC
	// with the desired name exists, or if not, create a new managed VPC.

	vpc, err := s.describeVPCByName()
	if err == nil {
		// An VPC already exists with the desired name

		if !vpc.Tags.HasOwned(s.scope.Name()) {
			return errors.Errorf(
				"found VPC %q which cannot be managed by CAPA due to lack of tags (either tag the VPC manually with `%s=%s`, or provide the `vpc.id` field instead if you wish to bring your own VPC as shown in https://cluster-api-aws.sigs.k8s.io/topics/bring-your-own-aws-infrastructure)",
				vpc.ID,
				infrav1.ClusterTagKey(s.scope.Name()),
				infrav1.ResourceLifecycleOwned)
		}
	} else {
		if !awserrors.IsNotFound(err) {
			return errors.Wrap(err, "failed to describe VPC resources by name")
		}

		// VPC with that name does not exist yet. Create it.
		vpc, err = s.createVPC()
		if err != nil {
			return errors.Wrap(err, "failed to create new managed VPC")
		}
		s.scope.Info("Created VPC", "vpc-id", vpc.ID)
	}

	s.scope.VPC().CidrBlock = vpc.CidrBlock
	s.scope.VPC().IPv6 = vpc.IPv6
	s.scope.VPC().Tags = vpc.Tags
	s.scope.VPC().ID = vpc.ID

	if !conditions.Has(s.scope.InfraCluster(), infrav1.VpcReadyCondition) {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, infrav1.VpcCreationStartedReason, clusterv1.ConditionSeverityInfo, "")
		if err := s.scope.PatchObject(); err != nil {
			return errors.Wrap(err, "failed to patch conditions")
		}
	}

	// Make sure attributes are configured
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if err := s.ensureManagedVPCAttributes(vpc); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.VPCNotFound); err != nil {
		return errors.Wrapf(err, "failed to set vpc attributes for %q", vpc.ID)
	}

	return nil
}

func (s *Service) describeVPCEndpoints(filters ...*ec2.Filter) ([]*ec2.VpcEndpoint, error) {
	vpc := s.scope.VPC()
	if vpc == nil || vpc.ID == "" {
		return nil, errors.New("vpc is nil or vpc id is not set")
	}
	input := &ec2.DescribeVpcEndpointsInput{
		Filters: append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: []*string{&vpc.ID},
		}),
	}
	endpoints := []*ec2.VpcEndpoint{}
	if err := s.EC2Client.DescribeVpcEndpointsPages(input, func(dveo *ec2.DescribeVpcEndpointsOutput, lastPage bool) bool {
		endpoints = append(endpoints, dveo.VpcEndpoints...)
		return true
	}); err != nil {
		return nil, errors.Wrap(err, "failed to describe vpc endpoints")
	}
	return endpoints, nil
}

// reconcileVPCEndpoints registers the AWS endpoints for the services that need to be enabled
// in the VPC routing tables. If the VPC is unmanaged, this is a no-op.
// For more information, see: https://docs.aws.amazon.com/vpc/latest/privatelink/gateway-endpoints.html
func (s *Service) reconcileVPCEndpoints() error {
	// If the VPC is unmanaged or not yet populated, return early.
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) || s.scope.VPC().ID == "" {
		return nil
	}

	// Gather all services that need to be enabled.
	services := sets.New[string]()
	if s.scope.Bucket() != nil {
		services.Insert(fmt.Sprintf("com.amazonaws.%s.s3", s.scope.Region()))
	}
	if services.Len() == 0 {
		return nil
	}

	// Gather the current routes.
	routeTables := sets.New[string]()
	for _, rt := range s.scope.Subnets() {
		if rt.RouteTableID != nil && *rt.RouteTableID != "" {
			routeTables.Insert(*rt.RouteTableID)
		}
	}
	if routeTables.Len() == 0 {
		return nil
	}

	// Build the filters based on all the services we need to enable.
	// A single filter with multiple values functions as an OR.
	filters := []*ec2.Filter{
		{
			Name:   aws.String("service-name"),
			Values: aws.StringSlice(services.UnsortedList()),
		},
	}

	// Get all existing endpoints.
	endpoints, err := s.describeVPCEndpoints(filters...)
	if err != nil {
		return errors.Wrap(err, "failed to describe vpc endpoints")
	}

	// Iterate over all services and create missing endpoints.
	for _, service := range services.UnsortedList() {
		var existing *ec2.VpcEndpoint
		for _, ep := range endpoints {
			if aws.StringValue(ep.ServiceName) == service {
				existing = ep
				break
			}
		}

		// Handle the case where the endpoint already exists.
		// If the route tables are different, modify the endpoint.
		if existing != nil {
			existingRouteTables := sets.New(aws.StringValueSlice(existing.RouteTableIds)...)
			existingRouteTables.Delete("")
			additions := routeTables.Difference(existingRouteTables)
			removals := existingRouteTables.Difference(routeTables)
			if additions.Len() > 0 || removals.Len() > 0 {
				modify := &ec2.ModifyVpcEndpointInput{
					VpcEndpointId: existing.VpcEndpointId,
				}
				if additions.Len() > 0 {
					modify.AddRouteTableIds = aws.StringSlice(additions.UnsortedList())
				}
				if removals.Len() > 0 {
					modify.RemoveRouteTableIds = aws.StringSlice(removals.UnsortedList())
				}
				if _, err := s.EC2Client.ModifyVpcEndpoint(modify); err != nil {
					return errors.Wrapf(err, "failed to modify vpc endpoint for service %q", service)
				}
			}
			continue
		}

		// Create the endpoint.
		if _, err := s.EC2Client.CreateVpcEndpoint(&ec2.CreateVpcEndpointInput{
			VpcId:         aws.String(s.scope.VPC().ID),
			ServiceName:   aws.String(service),
			RouteTableIds: aws.StringSlice(routeTables.UnsortedList()),
			TagSpecifications: []*ec2.TagSpecification{
				tags.BuildParamsToTagSpecification(ec2.ResourceTypeVpcEndpoint, s.getVPCEndpointTagParams()),
			},
		}); err != nil {
			return errors.Wrapf(err, "failed to create vpc endpoint for service %q", service)
		}
	}

	return nil
}

func (s *Service) deleteVPCEndpoints() error {
	// If the VPC is unmanaged or not yet populated, return early.
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) || s.scope.VPC().ID == "" {
		return nil
	}

	// Get all existing endpoints.
	endpoints, err := s.describeVPCEndpoints(filter.EC2.ClusterOwned(s.scope.Name()))
	if err != nil {
		return errors.Wrap(err, "failed to describe vpc endpoints")
	}

	// Gather all endpoint IDs.
	ids := []*string{}
	for _, ep := range endpoints {
		if ep.VpcEndpointId == nil || *ep.VpcEndpointId == "" {
			continue
		}
		ids = append(ids, ep.VpcEndpointId)
	}

	if len(ids) == 0 {
		return nil
	}

	// Iterate over all services and delete endpoints.
	if _, err := s.EC2Client.DeleteVpcEndpoints(&ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: ids,
	}); err != nil {
		return errors.Wrapf(err, "failed to delete vpc endpoints %+v", ids)
	}
	return nil
}

func (s *Service) ensureManagedVPCAttributes(vpc *infrav1.VPCSpec) error {
	var (
		errs    []error
		updated bool
	)

	// Cannot get or set both attributes at the same time.
	descAttrInput := &ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String(vpc.ID),
		Attribute: aws.String("enableDnsHostnames"),
	}
	vpcAttr, err := s.EC2Client.DescribeVpcAttributeWithContext(context.TODO(), descAttrInput)
	if err != nil {
		// If the returned error is a 'NotFound' error it should trigger retry
		if code, ok := awserrors.Code(errors.Cause(err)); ok && code == awserrors.VPCNotFound {
			return err
		}
		errs = append(errs, errors.Wrap(err, "failed to describe enableDnsHostnames vpc attribute"))
	} else if !aws.BoolValue(vpcAttr.EnableDnsHostnames.Value) {
		attrInput := &ec2.ModifyVpcAttributeInput{
			VpcId:              aws.String(vpc.ID),
			EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}
		if _, err := s.EC2Client.ModifyVpcAttributeWithContext(context.TODO(), attrInput); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to set enableDnsHostnames vpc attribute"))
		} else {
			updated = true
		}
	}

	descAttrInput = &ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String(vpc.ID),
		Attribute: aws.String("enableDnsSupport"),
	}
	vpcAttr, err = s.EC2Client.DescribeVpcAttributeWithContext(context.TODO(), descAttrInput)
	if err != nil {
		// If the returned error is a 'NotFound' error it should trigger retry
		if code, ok := awserrors.Code(errors.Cause(err)); ok && code == awserrors.VPCNotFound {
			return err
		}
		errs = append(errs, errors.Wrap(err, "failed to describe enableDnsSupport vpc attribute"))
	} else if !aws.BoolValue(vpcAttr.EnableDnsSupport.Value) {
		attrInput := &ec2.ModifyVpcAttributeInput{
			VpcId:            aws.String(vpc.ID),
			EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}
		if _, err := s.EC2Client.ModifyVpcAttributeWithContext(context.TODO(), attrInput); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to set enableDnsSupport vpc attribute"))
		} else {
			updated = true
		}
	}

	if len(errs) > 0 {
		record.Warnf(s.scope.InfraCluster(), "FailedSetVPCAttributes", "Failed to set managed VPC attributes for %q: %v", vpc.ID, err)
		return kerrors.NewAggregate(errs)
	}

	if updated {
		record.Eventf(s.scope.InfraCluster(), "SuccessfulSetVPCAttributes", "Set managed VPC attributes for %q", vpc.ID)
	}

	return nil
}

func (s *Service) getIPAMPoolID() (*string, error) {
	input := &ec2.DescribeIpamPoolsInput{}

	if s.scope.VPC().IPAMPool.ID != "" {
		input.Filters = append(input.Filters, filter.EC2.IPAM(s.scope.VPC().IPAMPool.ID))
	}

	if s.scope.VPC().IPAMPool.Name != "" {
		input.Filters = append(input.Filters, filter.EC2.Name(s.scope.VPC().IPAMPool.Name))
	}

	output, err := s.EC2Client.DescribeIpamPools(input)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateVPC", "Failed to describe IPAM Pools: %v", err)
		return nil, errors.Wrap(err, "failed to describe IPAM Pools")
	}

	switch len(output.IpamPools) {
	case 0:
		record.Warnf(s.scope.InfraCluster(), "FailedCreateVPC", "IPAM not found")
		return nil, fmt.Errorf("IPAM not found")
	case 1:
		return output.IpamPools[0].IpamPoolId, nil
	default:
		record.Warnf(s.scope.InfraCluster(), "FailedCreateVPC", "multiple IPAMs found")
		return nil, fmt.Errorf("multiple IPAMs found")
	}
}

func (s *Service) createVPC() (*infrav1.VPCSpec, error) {
	input := &ec2.CreateVpcInput{
		TagSpecifications: []*ec2.TagSpecification{
			tags.BuildParamsToTagSpecification(ec2.ResourceTypeVpc, s.getVPCTagParams(services.TemporaryResourceID)),
		},
	}

	// IPv6-specific configuration
	if s.scope.VPC().IsIPv6Enabled() {
		switch {
		case s.scope.VPC().IPv6.CidrBlock != "":
			input.Ipv6CidrBlock = aws.String(s.scope.VPC().IPv6.CidrBlock)
			input.Ipv6Pool = aws.String(s.scope.VPC().IPv6.PoolID)
			input.AmazonProvidedIpv6CidrBlock = aws.Bool(false)
		case s.scope.VPC().IPv6.IPAMPool != nil:
			ipamPoolID, err := s.getIPAMPoolID()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get IPAM Pool ID")
			}

			if s.scope.VPC().IPv6.IPAMPool.NetmaskLength == 0 {
				s.scope.VPC().IPv6.IPAMPool.NetmaskLength = defaultIpamV6NetmaskLength
			}

			input.Ipv6IpamPoolId = ipamPoolID
			input.Ipv6NetmaskLength = aws.Int64(s.scope.VPC().IPv6.IPAMPool.NetmaskLength)
		default:
			input.AmazonProvidedIpv6CidrBlock = aws.Bool(s.scope.VPC().IsIPv6Enabled())
		}
	}

	// IPv4-specific configuration
	if s.scope.VPC().IPAMPool != nil {
		ipamPoolID, err := s.getIPAMPoolID()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get IPAM Pool ID")
		}

		if s.scope.VPC().IPAMPool.NetmaskLength == 0 {
			s.scope.VPC().IPAMPool.NetmaskLength = defaultIpamV4NetmaskLength
		}

		input.Ipv4IpamPoolId = ipamPoolID
		input.Ipv4NetmaskLength = aws.Int64(s.scope.VPC().IPAMPool.NetmaskLength)
	} else {
		if s.scope.VPC().CidrBlock == "" {
			s.scope.VPC().CidrBlock = defaultVPCCidr
		}

		input.CidrBlock = &s.scope.VPC().CidrBlock
	}

	out, err := s.EC2Client.CreateVpcWithContext(context.TODO(), input)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateVPC", "Failed to create new managed VPC: %v", err)
		return nil, errors.Wrap(err, "failed to create vpc")
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateVPC", "Created new managed VPC %q", *out.Vpc.VpcId)
	s.scope.Debug("Created new VPC with cidr", "vpc-id", *out.Vpc.VpcId, "cidr-block", *out.Vpc.CidrBlock)

	if !s.scope.VPC().IsIPv6Enabled() {
		return &infrav1.VPCSpec{
			ID:        *out.Vpc.VpcId,
			CidrBlock: *out.Vpc.CidrBlock,
			Tags:      converters.TagsToMap(out.Vpc.Tags),
		}, nil
	}

	// BYOIP was defined, no need to look up the VPC.
	if s.scope.VPC().IsIPv6Enabled() && s.scope.VPC().IPv6.CidrBlock != "" {
		return &infrav1.VPCSpec{
			ID:        *out.Vpc.VpcId,
			CidrBlock: *out.Vpc.CidrBlock,
			IPv6: &infrav1.IPv6{
				CidrBlock: s.scope.VPC().IPv6.CidrBlock,
				PoolID:    s.scope.VPC().IPv6.PoolID,
			},
			Tags: converters.TagsToMap(out.Vpc.Tags),
		}, nil
	}

	// We have to describe the VPC again because the `create` output will **NOT** contain the associated IPv6 address.
	vpc, err := s.EC2Client.DescribeVpcsWithContext(context.TODO(), &ec2.DescribeVpcsInput{
		VpcIds: aws.StringSlice([]string{aws.StringValue(out.Vpc.VpcId)}),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "DescribeVpcs", "Failed to describe the new ipv6 vpc: %v", err)
		return nil, errors.Wrap(err, "failed to describe new ipv6 vpc")
	}
	if len(vpc.Vpcs) == 0 {
		record.Warnf(s.scope.InfraCluster(), "DescribeVpcs", "Failed to find the new ipv6 vpc, returned list was empty.")
		return nil, errors.New("failed to find new ipv6 vpc; returned list was empty")
	}
	for _, set := range vpc.Vpcs[0].Ipv6CidrBlockAssociationSet {
		if *set.Ipv6CidrBlockState.State == ec2.SubnetCidrBlockStateCodeAssociated {
			return &infrav1.VPCSpec{
				IPv6: &infrav1.IPv6{
					CidrBlock: aws.StringValue(set.Ipv6CidrBlock),
					PoolID:    aws.StringValue(set.Ipv6Pool),
				},
				ID:        *vpc.Vpcs[0].VpcId,
				CidrBlock: *out.Vpc.CidrBlock,
				Tags:      converters.TagsToMap(vpc.Vpcs[0].Tags),
			}, nil
		}
	}

	return nil, fmt.Errorf("no IPv6 associated CIDR block sets found for IPv6 enabled cluster with vpc id %s", *out.Vpc.VpcId)
}

func (s *Service) deleteVPC() error {
	vpc := s.scope.VPC()

	if vpc.IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping VPC deletion in unmanaged mode")
		return nil
	}

	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpc.ID),
	}

	if _, err := s.EC2Client.DeleteVpcWithContext(context.TODO(), input); err != nil {
		// Ignore if it's already deleted
		if code, ok := awserrors.Code(err); ok && code == awserrors.VPCNotFound {
			s.scope.Trace("Skipping VPC deletion, VPC not found")
			return nil
		}

		// Ignore if VPC ID is not present,
		if code, ok := awserrors.Code(err); ok && code == awserrors.VPCMissingParameter {
			s.scope.Trace("Skipping VPC deletion, VPC ID not present")
			return nil
		}

		record.Warnf(s.scope.InfraCluster(), "FailedDeleteVPC", "Failed to delete managed VPC %q: %v", vpc.ID, err)
		return errors.Wrapf(err, "failed to delete vpc %q", vpc.ID)
	}

	s.scope.Info("Deleted VPC", "vpc-id", vpc.ID)
	record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteVPC", "Deleted managed VPC %q", vpc.ID)
	return nil
}

func (s *Service) describeVPCByID() (*infrav1.VPCSpec, error) {
	if s.scope.VPC().ID == "" {
		return nil, errors.New("VPC ID is not set, failed to describe VPCs by ID")
	}

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPCStates(ec2.VpcStatePending, ec2.VpcStateAvailable),
		},
	}

	input.VpcIds = []*string{aws.String(s.scope.VPC().ID)}

	out, err := s.EC2Client.DescribeVpcsWithContext(context.TODO(), input)
	if err != nil {
		if awserrors.IsNotFound(err) {
			return nil, err
		}

		return nil, errors.Wrap(err, "failed to query ec2 for VPCs")
	}

	if len(out.Vpcs) == 0 {
		return nil, awserrors.NewNotFound(fmt.Sprintf("could not find vpc %q", s.scope.VPC().ID))
	} else if len(out.Vpcs) > 1 {
		return nil, awserrors.NewConflict(fmt.Sprintf("found %v VPCs with matching tags for %v. Only one VPC per cluster name is supported. Ensure duplicate VPCs are deleted for this AWS account and there are no conflicting instances of Cluster API Provider AWS. filtered VPCs: %v", len(out.Vpcs), s.scope.Name(), out.GoString()))
	}

	switch *out.Vpcs[0].State {
	case ec2.VpcStateAvailable, ec2.VpcStatePending:
	default:
		return nil, awserrors.NewNotFound("could not find available or pending vpc")
	}

	vpc := &infrav1.VPCSpec{
		ID:        *out.Vpcs[0].VpcId,
		CidrBlock: *out.Vpcs[0].CidrBlock,
		Tags:      converters.TagsToMap(out.Vpcs[0].Tags),
	}
	for _, set := range out.Vpcs[0].Ipv6CidrBlockAssociationSet {
		if *set.Ipv6CidrBlockState.State == ec2.SubnetCidrBlockStateCodeAssociated {
			vpc.IPv6 = &infrav1.IPv6{
				CidrBlock: aws.StringValue(set.Ipv6CidrBlock),
				PoolID:    aws.StringValue(set.Ipv6Pool),
			}
			break
		}
	}
	return vpc, nil
}

// describeVPCByName finds the VPC by `Name` tag. Use this if the ID is not available yet, either because no
// VPC was created until now or if storing the ID could have failed.
func (s *Service) describeVPCByName() (*infrav1.VPCSpec, error) {
	vpcName := *s.getVPCTagParams(services.TemporaryResourceID).Name

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: aws.StringSlice([]string{vpcName}),
			},
		},
	}

	out, err := s.EC2Client.DescribeVpcsWithContext(context.TODO(), input)
	if (err != nil && awserrors.IsNotFound(err)) || (out != nil && len(out.Vpcs) == 0) {
		return nil, awserrors.NewNotFound(fmt.Sprintf("could not find VPC by name %q", vpcName))
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query ec2 for VPCs by name %q", vpcName)
	}
	if len(out.Vpcs) > 1 {
		return nil, awserrors.NewConflict(fmt.Sprintf("found %v VPCs with name %q. Only one VPC per cluster name is supported. Ensure duplicate VPCs are deleted for this AWS account and there are no conflicting instances of Cluster API Provider AWS. Filtered VPCs: %v", len(out.Vpcs), vpcName, out.GoString()))
	}

	switch *out.Vpcs[0].State {
	case ec2.VpcStateAvailable, ec2.VpcStatePending:
	default:
		return nil, awserrors.NewNotFound(fmt.Sprintf("could not find available or pending VPC by name %q", vpcName))
	}

	vpc := &infrav1.VPCSpec{
		ID:        *out.Vpcs[0].VpcId,
		CidrBlock: *out.Vpcs[0].CidrBlock,
		Tags:      converters.TagsToMap(out.Vpcs[0].Tags),
	}
	for _, set := range out.Vpcs[0].Ipv6CidrBlockAssociationSet {
		if *set.Ipv6CidrBlockState.State == ec2.SubnetCidrBlockStateCodeAssociated {
			vpc.IPv6 = &infrav1.IPv6{
				CidrBlock: aws.StringValue(set.Ipv6CidrBlock),
				PoolID:    aws.StringValue(set.Ipv6Pool),
			}
			break
		}
	}
	return vpc, nil
}

func (s *Service) getVPCTagParams(id string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-vpc", s.scope.Name())

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}

func (s *Service) getVPCEndpointTagParams() infrav1.BuildParams {
	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}
