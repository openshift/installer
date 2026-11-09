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
	"strings"

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
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

const (
	mainRouteTableInVPCKey = "main"
)

func (s *Service) reconcileRouteTables() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping routing tables reconcile in unmanaged mode")
		return nil
	}

	s.scope.Debug("Reconciling routing tables")

	subnetRouteMap, err := s.describeVpcRouteTablesBySubnet()
	if err != nil {
		return err
	}

	subnets := s.scope.Subnets()
	defer func() {
		s.scope.SetSubnets(subnets)
	}()

	for i := range subnets {
		sn := &subnets[i]
		// We need to compile the minimum routes for this subnet first, so we can compare it or create them.
		routes, err := s.getRoutesForSubnet(sn)
		if err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedRouteTableRoutes", "Failed to get routes for managed RouteTable for subnet %s: %v", sn.ID, err)
			return errors.Wrapf(err, "failed to discover routes on route table %s", sn.ID)
		}

		if rt, ok := subnetRouteMap[sn.GetResourceID()]; ok {
			s.scope.Debug("Subnet is already associated with route table", "subnet-id", sn.GetResourceID(), "route-table-id", aws.ToString(rt.RouteTableId))
			// TODO(vincepri): check that everything is in order, e.g. routes match the subnet type.

			// For managed environments we reconcile every desired route. A route whose
			// destination is missing gets created (e.g. for bring-your-own VPC that didn't have it),
			// and a route whose destination matches but points at the wrong target gets
			// replaced (e.g. a gateway that was deleted and re-created).
			for _, route := range routes {
				if err := s.reconcileRoute(route, rt); err != nil {
					return err
				}
			}

			// Make sure tags are up-to-date.
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				buildParams := s.getRouteTableTagParams(aws.ToString(rt.RouteTableId), sn.IsPublic, sn.AvailabilityZone)
				tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
				if err := tagsBuilder.Ensure(converters.TagsToMap(rt.Tags)); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.RouteTableNotFound); err != nil {
				record.Warnf(s.scope.InfraCluster(), "FailedTagRouteTable", "Failed to tag managed RouteTable %q: %v", aws.ToString(rt.RouteTableId), err)
				return errors.Wrapf(err, "failed to ensure tags on route table %q", aws.ToString(rt.RouteTableId))
			}

			// Not recording "SuccessfulTagRouteTable" here as we don't know if this was a no-op or an actual change
			continue
		}
		s.scope.Debug("Subnet isn't associated with route table", "subnet-id", sn.GetResourceID())

		// For each subnet that doesn't have a routing table associated with it,
		// create a new table with the appropriate default routes and associate it to the subnet.
		rt, err := s.createRouteTableWithRoutes(routes, sn.IsPublic, sn.AvailabilityZone)
		if err != nil {
			return err
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if err := s.associateRouteTable(rt, sn.GetResourceID()); err != nil {
				s.scope.Error(err, "trying to associate route table", "subnet_id", sn.GetResourceID())
				return false, err
			}
			return true, nil
		}, awserrors.RouteTableNotFound, awserrors.SubnetNotFound); err != nil {
			return err
		}

		s.scope.Debug("Subnet has been associated with route table", "subnet-id", sn.GetResourceID(), "route-table-id", rt.ID)
		sn.RouteTableID = aws.String(rt.ID)
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition)
	return nil
}

// reconcileRoute ensures a single desired route exists in the route table pointing
// at the desired target.
func (s *Service) reconcileRoute(specRoute *ec2.CreateRouteInput, rt types.RouteTable) error {
	var current *types.Route
	for i := range rt.Routes {
		route := &rt.Routes[i]
		if (specRoute.DestinationCidrBlock != nil && aws.ToString(route.DestinationCidrBlock) == aws.ToString(specRoute.DestinationCidrBlock)) ||
			(specRoute.DestinationIpv6CidrBlock != nil && aws.ToString(route.DestinationIpv6CidrBlock) == aws.ToString(specRoute.DestinationIpv6CidrBlock)) {
			current = route
			break
		}
	}

	if current != nil &&
		aws.ToString(current.GatewayId) == aws.ToString(specRoute.GatewayId) &&
		aws.ToString(current.NatGatewayId) == aws.ToString(specRoute.NatGatewayId) &&
		aws.ToString(current.EgressOnlyInternetGatewayId) == aws.ToString(specRoute.EgressOnlyInternetGatewayId) &&
		aws.ToString(current.CarrierGatewayId) == aws.ToString(specRoute.CarrierGatewayId) {
		// The route already exists with the desired target, nothing to do
		return nil
	}

	routeDesc := describeRoute(specRoute)
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if current == nil {
			specRoute.RouteTableId = rt.RouteTableId
			_, err := s.EC2Client.CreateRoute(context.TODO(), specRoute)
			return err == nil, err
		}
		_, err := s.EC2Client.ReplaceRoute(context.TODO(), &ec2.ReplaceRouteInput{
			RouteTableId:                rt.RouteTableId,
			DestinationCidrBlock:        specRoute.DestinationCidrBlock,
			DestinationIpv6CidrBlock:    specRoute.DestinationIpv6CidrBlock,
			DestinationPrefixListId:     specRoute.DestinationPrefixListId,
			GatewayId:                   specRoute.GatewayId,
			NatGatewayId:                specRoute.NatGatewayId,
			EgressOnlyInternetGatewayId: specRoute.EgressOnlyInternetGatewayId,
			CarrierGatewayId:            specRoute.CarrierGatewayId,
		})
		return err == nil, err
	}, awserrors.RouteTableNotFound, awserrors.NATGatewayNotFound, awserrors.GatewayNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedReconcileRoute", "Failed to reconcile route %s on managed RouteTable %q: %v", routeDesc, aws.ToString(rt.RouteTableId), err)
		return errors.Wrapf(err, "failed to reconcile route %s on route table %q", routeDesc, aws.ToString(rt.RouteTableId))
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulReconcileRoute", "Reconciled route %s on RouteTable %q", routeDesc, aws.ToString(rt.RouteTableId))
	return nil
}

func (s *Service) describeVpcRouteTablesBySubnet() (map[string]types.RouteTable, error) {
	rts, err := s.describeVpcRouteTables()
	if err != nil {
		return nil, err
	}

	// Amazon allows a subnet to be associated only with a single routing table
	// https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html.
	res := make(map[string]types.RouteTable)
	for _, rt := range rts {
		for _, as := range rt.Associations {
			if as.Main != nil && *as.Main {
				res[mainRouteTableInVPCKey] = rt
			}
			if as.SubnetId == nil {
				continue
			}

			res[aws.ToString(as.SubnetId)] = rt
		}
	}

	return res, nil
}

func (s *Service) deleteRouteTable(rt types.RouteTable) error {
	for _, as := range rt.Associations {
		if as.SubnetId == nil {
			continue
		}

		if _, err := s.EC2Client.DisassociateRouteTable(context.TODO(), &ec2.DisassociateRouteTableInput{AssociationId: as.RouteTableAssociationId}); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedDisassociateRouteTable", "Failed to disassociate managed RouteTable %q from Subnet %q: %v", aws.ToString(rt.RouteTableId), aws.ToString(as.SubnetId), err)
			return errors.Wrapf(err, "failed to disassociate route table %q from subnet %q", aws.ToString(rt.RouteTableId), aws.ToString(as.SubnetId))
		}

		record.Eventf(s.scope.InfraCluster(), "SuccessfulDisassociateRouteTable", "Disassociated managed RouteTable %q from subnet %q", aws.ToString(rt.RouteTableId), aws.ToString(as.SubnetId))
		s.scope.Debug("Deleted association between route table and subnet", "route-table-id", aws.ToString(rt.RouteTableId), "subnet-id", aws.ToString(as.SubnetId))
	}

	if _, err := s.EC2Client.DeleteRouteTable(context.TODO(), &ec2.DeleteRouteTableInput{RouteTableId: rt.RouteTableId}); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedDeleteRouteTable", "Failed to delete managed RouteTable %q: %v", aws.ToString(rt.RouteTableId), err)
		return errors.Wrapf(err, "failed to delete route table %q", aws.ToString(rt.RouteTableId))
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteRouteTable", "Deleted managed RouteTable %q", aws.ToString(rt.RouteTableId))
	s.scope.Info("Deleted route table", "route-table-id", aws.ToString(rt.RouteTableId))

	return nil
}

func (s *Service) deleteRouteTables() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.Trace("Skipping routing tables deletion in unmanaged mode")
		return nil
	}

	rts, err := s.describeVpcRouteTables()
	if err != nil {
		return errors.Wrapf(err, "failed to describe route tables in vpc %q", s.scope.VPC().ID)
	}

	for _, rt := range rts {
		err := s.deleteRouteTable(rt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) describeVpcRouteTables() ([]types.RouteTable, error) {
	filters := []types.Filter{
		filter.EC2.VPC(s.scope.VPC().ID),
	}

	if !s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		filters = append(filters, filter.EC2.Cluster(s.scope.Name()))
	}

	out, err := s.EC2Client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{
		Filters: filters,
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeVPCRouteTable", "Failed to describe route tables in vpc %q: %v", s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to describe route tables in vpc %q", s.scope.VPC().ID)
	}

	return out.RouteTables, nil
}

func (s *Service) createRouteTableWithRoutes(routes []*ec2.CreateRouteInput, isPublic bool, zone string) (*infrav1.RouteTable, error) {
	out, err := s.EC2Client.CreateRouteTable(context.TODO(), &ec2.CreateRouteTableInput{
		VpcId: aws.String(s.scope.VPC().ID),
		TagSpecifications: []types.TagSpecification{
			tags.BuildParamsToTagSpecification(types.ResourceTypeRouteTable, s.getRouteTableTagParams(services.TemporaryResourceID, isPublic, zone)),
		},
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateRouteTable", "Failed to create managed RouteTable: %v", err)
		return nil, errors.Wrapf(err, "failed to create route table in vpc %q", s.scope.VPC().ID)
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateRouteTable", "Created managed RouteTable %q", aws.ToString(out.RouteTable.RouteTableId))
	s.scope.Info("Created route table", "route-table-id", *out.RouteTable.RouteTableId)

	for i := range routes {
		route := routes[i]
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			route.RouteTableId = out.RouteTable.RouteTableId
			if _, err := s.EC2Client.CreateRoute(context.TODO(), route); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.RouteTableNotFound, awserrors.NATGatewayNotFound, awserrors.GatewayNotFound); err != nil {
			routeDesc := describeRoute(route)
			record.Warnf(s.scope.InfraCluster(), "FailedCreateRoute", "Failed to create route %s for RouteTable %q: %v", routeDesc, aws.ToString(out.RouteTable.RouteTableId), err)
			errDel := s.deleteRouteTable(*out.RouteTable)
			if errDel != nil {
				record.Warnf(s.scope.InfraCluster(), "FailedDeleteRouteTable", "Failed to delete managed RouteTable %q: %v", aws.ToString(out.RouteTable.RouteTableId), errDel)
			}
			return nil, errors.Wrapf(err, "failed to create route in route table %q: %s", aws.ToString(out.RouteTable.RouteTableId), routeDesc)
		}
		record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateRoute", "Created route %s for RouteTable %q", describeRoute(route), aws.ToString(out.RouteTable.RouteTableId))
	}

	return &infrav1.RouteTable{
		ID: *out.RouteTable.RouteTableId,
	}, nil
}

func (s *Service) associateRouteTable(rt *infrav1.RouteTable, subnetID string) error {
	_, err := s.EC2Client.AssociateRouteTable(context.TODO(), &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rt.ID),
		SubnetId:     aws.String(subnetID),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAssociateRouteTable", "Failed to associate managed RouteTable %q with Subnet %q: %v", rt.ID, subnetID, err)
		return errors.Wrapf(err, "failed to associate route table %q to subnet %q", rt.ID, subnetID)
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulAssociateRouteTable", "Associated managed RouteTable %q with subnet %q", rt.ID, subnetID)

	return nil
}

func (s *Service) getNatGatewayPrivateRoute(natGatewayID string) *ec2.CreateRouteInput {
	return &ec2.CreateRouteInput{
		NatGatewayId:         aws.String(natGatewayID),
		DestinationCidrBlock: aws.String(services.AnyIPv4CidrBlock),
	}
}

func (s *Service) getNat64PrivateRoute(natGatewayID string) *ec2.CreateRouteInput {
	return &ec2.CreateRouteInput{
		NatGatewayId:             aws.String(natGatewayID),
		DestinationIpv6CidrBlock: aws.String(services.NAT64CidrBlock),
	}
}

func (s *Service) getEgressOnlyInternetGateway() *ec2.CreateRouteInput {
	return &ec2.CreateRouteInput{
		DestinationIpv6CidrBlock:    aws.String(services.AnyIPv6CidrBlock),
		EgressOnlyInternetGatewayId: s.scope.VPC().IPv6.EgressOnlyInternetGatewayID,
	}
}

func (s *Service) getGatewayPublicRoute() *ec2.CreateRouteInput {
	return &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String(services.AnyIPv4CidrBlock),
		GatewayId:            aws.String(*s.scope.VPC().InternetGatewayID),
	}
}

func (s *Service) getGatewayPublicIPv6Route() *ec2.CreateRouteInput {
	return &ec2.CreateRouteInput{
		DestinationIpv6CidrBlock: aws.String(services.AnyIPv6CidrBlock),
		GatewayId:                aws.String(*s.scope.VPC().InternetGatewayID),
	}
}

func (s *Service) getCarrierGatewayPublicIPv4Route() *ec2.CreateRouteInput {
	return &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String(services.AnyIPv4CidrBlock),
		CarrierGatewayId:     aws.String(*s.scope.VPC().CarrierGatewayID),
	}
}

func (s *Service) getRouteTableTagParams(id string, public bool, zone string) infrav1.BuildParams {
	var name strings.Builder

	name.WriteString(s.scope.Name())
	name.WriteString("-rt-")
	if public {
		name.WriteString("public")
	} else {
		name.WriteString("private")
	}
	name.WriteString("-")
	name.WriteString(zone)

	additionalTags := s.scope.AdditionalTags()
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())] = string(infrav1.ResourceLifecycleOwned)

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name.String()),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  additionalTags,
	}
}

func (s *Service) getRoutesToPublicSubnet(sn *infrav1.SubnetSpec) ([]*ec2.CreateRouteInput, error) {
	var routes []*ec2.CreateRouteInput

	if sn.IsEdge() && sn.IsIPv6 {
		return nil, errors.Errorf("can't determine routes for unsupported ipv6 subnet in zone type %q", sn.ZoneType)
	}

	if sn.IsEdgeWavelength() {
		if s.scope.VPC().CarrierGatewayID == nil {
			return routes, errors.Errorf("failed to create carrier routing table: carrier gateway for VPC %q is not present", s.scope.VPC().ID)
		}
		routes = append(routes, s.getCarrierGatewayPublicIPv4Route())
		return routes, nil
	}

	if s.scope.VPC().InternetGatewayID == nil {
		return routes, errors.Errorf("failed to create routing tables: internet gateway for VPC %q is not present", s.scope.VPC().ID)
	}
	routes = append(routes, s.getGatewayPublicRoute())
	if sn.IsIPv6 {
		routes = append(routes, s.getGatewayPublicIPv6Route())
	}

	return routes, nil
}

func (s *Service) getRoutesToPrivateSubnet(sn *infrav1.SubnetSpec) (routes []*ec2.CreateRouteInput, err error) {
	var natGatewayID string

	if sn.IsEdge() && sn.IsIPv6 {
		return nil, errors.Errorf("can't determine routes for unsupported ipv6 subnet in zone type %q", sn.ZoneType)
	}

	natGatewayID, err = s.getNatGatewayForSubnet(sn)
	if err != nil {
		return routes, err
	}

	routes = append(routes, s.getNatGatewayPrivateRoute(natGatewayID))
	if sn.IsIPv6 {
		// We add the NAT64 route only if DNS64 is enabled for the subnet
		// That is when the subnet is private and IPv6-only.
		if sn.CidrBlock == "" {
			routes = append(routes, s.getNat64PrivateRoute(natGatewayID))
		}
		if !s.scope.VPC().IsIPv6Enabled() {
			// Safety net because EgressOnlyInternetGateway needs the ID from the ipv6 block.
			// if, for whatever reason by this point that is not available, we don't want to
			// panic because of a nil pointer access. This should never occur. Famous last words though.
			return routes, errors.Errorf("ipv6 block missing for ipv6 enabled subnet, can't create route for egress only internet gateway")
		}
		routes = append(routes, s.getEgressOnlyInternetGateway())
	}

	return routes, nil
}

func (s *Service) getRoutesForSubnet(sn *infrav1.SubnetSpec) ([]*ec2.CreateRouteInput, error) {
	if sn.IsPublic {
		return s.getRoutesToPublicSubnet(sn)
	}
	return s.getRoutesToPrivateSubnet(sn)
}

// describeRoute renders a CreateRouteInput as a human-readable string for logs and events.
func describeRoute(route *ec2.CreateRouteInput) string {
	if route == nil {
		return "<nil>"
	}
	var parts []string
	addPart := func(name string, value *string) {
		if v := aws.ToString(value); v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", name, v))
		}
	}
	addPart("destinationCidrBlock", route.DestinationCidrBlock)
	addPart("destinationIpv6CidrBlock", route.DestinationIpv6CidrBlock)
	addPart("destinationPrefixListId", route.DestinationPrefixListId)
	addPart("gatewayId", route.GatewayId)
	addPart("natGatewayId", route.NatGatewayId)
	addPart("egressOnlyInternetGatewayId", route.EgressOnlyInternetGatewayId)
	addPart("carrierGatewayId", route.CarrierGatewayId)
	return "{" + strings.Join(parts, ", ") + "}"
}
