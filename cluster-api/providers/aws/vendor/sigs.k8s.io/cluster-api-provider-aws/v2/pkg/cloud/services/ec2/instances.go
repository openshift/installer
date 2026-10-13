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

package ec2

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/common"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/utils"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// GetRunningInstanceByTags returns the existing instance or nothing if it doesn't exist.
func (s *Service) GetRunningInstanceByTags(scope *scope.MachineScope) (*infrav1.Instance, error) {
	s.scope.Debug("Looking for existing machine instance by tags")

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			filter.EC2.ClusterOwned(s.scope.Name()),
			filter.EC2.Name(scope.Name()),
			filter.EC2.InstanceStates(types.InstanceStateNamePending, types.InstanceStateNameRunning),
		},
	}

	out, err := s.EC2Client.DescribeInstances(context.TODO(), input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeInstances", "Failed to describe instances by tags: %v", err)
		return nil, errors.Wrap(err, "failed to describe instances by tags")
	}

	// TODO: currently just returns the first matched instance, need to
	// better rationalize how to find the right instance to return if multiple
	// match
	for _, res := range out.Reservations {
		for _, inst := range res.Instances {
			return s.SDKToInstance(inst)
		}
	}

	return nil, nil
}

// InstanceIfExists returns the existing instance by id and errors if it cannot find the instance(ErrInstanceNotFoundByID) or API call fails (ErrDescribeInstance).
// Returns empty instance with nil error, only when providerID is nil.
func (s *Service) InstanceIfExists(id *string) (*infrav1.Instance, error) {
	if id == nil {
		s.scope.Info("Instance does not have an instance id")
		return nil, nil
	}

	s.scope.Debug("Looking for instance by id", "instance-id", *id)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{aws.ToString(id)},
	}

	out, err := s.EC2Client.DescribeInstances(context.TODO(), input)
	switch {
	case awserrors.IsNotFound(err):
		record.Eventf(s.scope.InfraCluster(), "FailedFindInstances", "failed to find instance by providerId %q: %v", *id, err)
		return nil, ErrInstanceNotFoundByID
	case err != nil:
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeInstances", "failed to describe instance %q: %v", *id, err)
		return nil, ErrDescribeInstance
	}

	if len(out.Reservations) > 0 && len(out.Reservations[0].Instances) > 0 {
		return s.SDKToInstance(out.Reservations[0].Instances[0])
	}

	// Failed to find instance with provider id.
	record.Eventf(s.scope.InfraCluster(), "FailedFindInstances", "failed to find instance by providerId %q: %v", *id, err)
	return nil, ErrInstanceNotFoundByID
}

// CreateInstance runs an ec2 instance.
//
//nolint:gocyclo // this function has multiple processes to perform
func (s *Service) CreateInstance(ctx context.Context, scope *scope.MachineScope, userData []byte, userDataFormat string) (*infrav1.Instance, error) {
	s.scope.Debug("Creating an instance for a machine")

	input := &infrav1.Instance{
		Type:                 scope.AWSMachine.Spec.InstanceType,
		IAMProfile:           scope.AWSMachine.Spec.IAMInstanceProfile,
		RootVolume:           scope.AWSMachine.Spec.RootVolume.DeepCopy(),
		NonRootVolumes:       scope.AWSMachine.Spec.NonRootVolumes,
		NetworkInterfaces:    scope.AWSMachine.Spec.NetworkInterfaces,
		AssignPrimaryIPv6:    scope.AWSMachine.Spec.AssignPrimaryIPv6,
		NetworkInterfaceType: scope.AWSMachine.Spec.NetworkInterfaceType,
	}

	// Make sure to use the MachineScope here to get the merger of AWSCluster and AWSMachine tags
	additionalTags := scope.AdditionalTags()
	input.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.KubernetesClusterName(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String(scope.Role()),
		Additional:  additionalTags,
	}.WithCloudProvider(s.scope.KubernetesClusterName()).WithMachineName(scope.Machine))

	var err error

	imageArchitecture, err := s.pickArchitectureForInstanceType(types.InstanceType(input.Type))
	if err != nil {
		return nil, err
	}

	// Pick image from the machine configuration, or use a default one.
	if scope.AWSMachine.Spec.AMI.ID != nil { //nolint:nestif
		input.ImageID = *scope.AWSMachine.Spec.AMI.ID
	} else {
		if scope.Machine.Spec.Version == "" {
			err := errors.New("Either AWSMachine's spec.ami.id or Machine's spec.version must be defined")
			scope.SetFailureReason("CreateError")
			scope.SetFailureMessage(err)
			return nil, err
		}

		imageLookupFormat := scope.AWSMachine.Spec.ImageLookupFormat
		if imageLookupFormat == "" {
			imageLookupFormat = scope.InfraCluster.ImageLookupFormat()
		}

		imageLookupOrg := scope.AWSMachine.Spec.ImageLookupOrg
		if imageLookupOrg == "" {
			imageLookupOrg = scope.InfraCluster.ImageLookupOrg()
		}

		imageLookupBaseOS := scope.AWSMachine.Spec.ImageLookupBaseOS
		if imageLookupBaseOS == "" {
			imageLookupBaseOS = scope.InfraCluster.ImageLookupBaseOS()
		}

		if scope.IsEKSManaged() && imageLookupFormat == "" && imageLookupOrg == "" && imageLookupBaseOS == "" {
			input.ImageID, err = s.eksAMILookup(ctx, scope.Machine.Spec.Version, imageArchitecture, scope.AWSMachine.Spec.AMI.EKSOptimizedLookupType)
			if err != nil {
				return nil, err
			}
		} else {
			input.ImageID, err = s.defaultAMIIDLookup(imageLookupFormat, imageLookupOrg, imageLookupBaseOS, imageArchitecture, scope.Machine.Spec.Version)
			if err != nil {
				return nil, err
			}
		}
	}

	subnetID, err := s.findSubnet(scope)
	if err != nil {
		return nil, err
	}
	input.SubnetID = subnetID

	// Preserve user-defined PublicIp option.
	input.PublicIPOnLaunch = scope.AWSMachine.Spec.PublicIP

	// Public address from BYO Public IPv4 Pools need to be associated after launch (main machine
	// reconciliate loop) preventing duplicated public IP. The map on launch is explicitly
	// disabled in instances with PublicIP defined to true.
	if scope.AWSMachine.Spec.ElasticIPPool != nil && scope.AWSMachine.Spec.ElasticIPPool.PublicIpv4Pool != nil {
		input.PublicIPOnLaunch = ptr.To(false)
	}

	if !scope.IsControlPlaneExternallyManaged() && !scope.IsExternallyManaged() && !scope.IsEKSManaged() && s.scope.Network().APIServerELB.DNSName == "" {
		record.Eventf(s.scope.InfraCluster(), "FailedCreateInstance", "Failed to run controlplane, APIServer ELB not available")
		return nil, awserrors.NewFailedDependency("failed to run controlplane, APIServer ELB not available")
	}

	if scope.CompressUserData(userDataFormat) {
		userData, err = userdata.GzipBytes(userData)
		if err != nil {
			return nil, errors.New("failed to gzip userdata")
		}
	}

	input.UserData = ptr.To[string](base64.StdEncoding.EncodeToString(userData))

	// Set security groups.
	ids, err := s.GetCoreSecurityGroups(scope)
	if err != nil {
		return nil, err
	}
	input.SecurityGroupIDs = append(input.SecurityGroupIDs, ids...)

	// If SSHKeyName WAS NOT provided in the AWSMachine Spec, fallback to the value provided in the AWSCluster Spec.
	// If a value was not provided in the AWSCluster Spec, then use the defaultSSHKeyName
	// Note that:
	// - a nil AWSMachine.Spec.SSHKeyName value means use the AWSCluster.Spec.SSHKeyName SSH key name value
	// - nil values for both AWSCluster.Spec.SSHKeyName and AWSMachine.Spec.SSHKeyName means use the default SSH key name value
	// - an empty string means do not set an SSH key name at all
	// - otherwise use the value specified in either AWSMachine or AWSCluster
	var prioritizedSSHKeyName string
	switch {
	case scope.AWSMachine.Spec.SSHKeyName != nil:
		// prefer AWSMachine.Spec.SSHKeyName if it is defined
		prioritizedSSHKeyName = *scope.AWSMachine.Spec.SSHKeyName
	case scope.InfraCluster.SSHKeyName() != nil:
		// fallback to AWSCluster.Spec.SSHKeyName if it is defined
		prioritizedSSHKeyName = *scope.InfraCluster.SSHKeyName()
	default:
		if !scope.IsExternallyManaged() {
			prioritizedSSHKeyName = defaultSSHKeyName
		}
	}

	// Only set input.SSHKeyName if the user did not explicitly request no ssh key be set (explicitly setting "" on either the Machine or related Cluster)
	if prioritizedSSHKeyName != "" {
		input.SSHKeyName = aws.String(prioritizedSSHKeyName)
	}

	input.SpotMarketOptions = scope.AWSMachine.Spec.SpotMarketOptions

	input.InstanceMetadataOptions = scope.AWSMachine.Spec.InstanceMetadataOptions

	input.Tenancy = scope.AWSMachine.Spec.Tenancy

	input.PlacementGroupName = scope.AWSMachine.Spec.PlacementGroupName

	input.PlacementGroupPartition = scope.AWSMachine.Spec.PlacementGroupPartition

	input.PrivateDNSName = scope.AWSMachine.Spec.PrivateDNSName

	input.CapacityReservationID = scope.AWSMachine.Spec.CapacityReservationID

	input.MarketType = scope.AWSMachine.Spec.MarketType

	// Handle dynamic host allocation if specified
	if scope.AWSMachine.Spec.DynamicHostAllocation != nil {
		hostID, err := s.ensureDedicatedHostAllocation(ctx, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to allocate dedicated host")
		}
		input.HostID = aws.String(hostID)
		input.HostAffinity = aws.String("host")

		if scope.AWSMachine.Status.DedicatedHost == nil {
			scope.AWSMachine.Status.DedicatedHost = &infrav1.DedicatedHostStatus{}
		}
		// Update machine status with allocated host ID
		scope.AWSMachine.Status.DedicatedHost.ID = &hostID
	} else {
		// Use static host allocation if specified
		input.HostID = scope.AWSMachine.Spec.HostID
		input.HostAffinity = scope.AWSMachine.Spec.HostAffinity
	}

	input.CapacityReservationPreference = scope.AWSMachine.Spec.CapacityReservationPreference

	input.CPUOptions = scope.AWSMachine.Spec.CPUOptions

	s.scope.Debug("Running instance", "machine-role", scope.Role())
	s.scope.Debug("Running instance with instance metadata options", "metadata options", input.InstanceMetadataOptions)
	out, err := s.runInstance(scope.Role(), input)
	if err != nil {
		// Only record the failure event if the error is not related to failed dependencies.
		// This is to avoid spamming failure events since the machine will be requeued by the actuator.
		if !awserrors.IsFailedDependency(errors.Cause(err)) {
			record.Warnf(scope.AWSMachine, "FailedCreate", "Failed to create instance: %v", err)
		}
		return nil, err
	}

	// Set the providerID and instanceID as soon as we create an instance so that we keep it in case of errors afterward
	scope.SetProviderID(out.ID, out.AvailabilityZone)
	scope.SetInstanceID(out.ID)

	if len(input.NetworkInterfaces) > 0 {
		for _, id := range input.NetworkInterfaces {
			s.scope.Debug("Attaching security groups to provided network interface", "groups", input.SecurityGroupIDs, "interface", id)
			if err := s.attachSecurityGroupsToNetworkInterface(input.SecurityGroupIDs, id); err != nil {
				return nil, err
			}
		}
	}

	s.scope.Debug("Adding tags on each network interface from resource", "resource-id", out.ID)

	// Fetching the network interfaces attached to the specific instance
	networkInterfaces, err := s.getInstanceENIs(out.ID)
	if err != nil {
		return nil, err
	}

	s.scope.Debug("Fetched the network interfaces")

	// Once all the network interfaces attached to the specific instance are found, the similar tags of instance are created for network interfaces too
	if len(networkInterfaces) > 0 {
		s.scope.Debug("Attempting to create tags from resource", "resource-id", out.ID)
		for _, networkInterface := range networkInterfaces {
			// Create/Update tags in AWS.
			if err := s.UpdateResourceTags(networkInterface.NetworkInterfaceId, out.Tags, nil); err != nil {
				return nil, errors.Wrapf(err, "failed to create tags for resource %q: ", *networkInterface.NetworkInterfaceId)
			}
		}
	}

	record.Eventf(scope.AWSMachine, "SuccessfulCreate", "Created new %s instance with id %q", scope.Role(), out.ID)
	return out, nil
}

// findSubnet attempts to retrieve a subnet ID in the following order:
// - subnetID specified in machine configuration,
// - subnet based on filters in machine configuration
// - subnet based on the availability zone specified,
// - default to the first private subnet available.
func (s *Service) findSubnet(scope *scope.MachineScope) (string, error) {
	// Check Machine.Spec.FailureDomain first as it's used by KubeadmControlPlane to spread machines across failure domains.
	failureDomain := scope.Machine.Spec.FailureDomain

	// We basically have 2 sources for subnets:
	//   1. If subnet.id or subnet.filters are specified, we directly query AWS
	//   2. All other cases use the subnets provided in the cluster network spec without ever calling AWS

	switch {
	case scope.AWSMachine.Spec.Subnet != nil && (scope.AWSMachine.Spec.Subnet.ID != nil || scope.AWSMachine.Spec.Subnet.Filters != nil):
		criteria := []types.Filter{
			filter.EC2.SubnetStates(types.SubnetStatePending, types.SubnetStateAvailable),
		}
		if scope.AWSMachine.Spec.Subnet.ID != nil {
			criteria = append(criteria, types.Filter{Name: aws.String("subnet-id"), Values: []string{*scope.AWSMachine.Spec.Subnet.ID}})
		}
		for _, f := range scope.AWSMachine.Spec.Subnet.Filters {
			criteria = append(criteria, types.Filter{Name: aws.String(f.Name), Values: f.Values})
		}

		subnets, err := s.getFilteredSubnets(criteria...)
		if err != nil {
			return "", errors.Wrapf(err, "failed to filter subnets for criteria %v", criteria)
		}
		if len(subnets) == 0 {
			errMessage := fmt.Sprintf("failed to run machine %q, no subnets available matching criteria %v",
				scope.Name(), criteria)
			record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}

		var filtered []types.Subnet
		var errMessage string
		for _, subnet := range subnets {
			if failureDomain != "" && *subnet.AvailabilityZone != failureDomain {
				// we could have included the failure domain in the query criteria, but then we end up with EC2 error
				// messages that don't give a good hint about what is really wrong
				errMessage += fmt.Sprintf(" subnet %q availability zone %q does not match failure domain %q.",
					*subnet.SubnetId, *subnet.AvailabilityZone, failureDomain)
				continue
			}

			if ptr.Deref(scope.AWSMachine.Spec.PublicIP, false) {
				matchingSubnet := s.scope.Subnets().FindByID(*subnet.SubnetId)
				if matchingSubnet == nil {
					errMessage += fmt.Sprintf(" unable to find subnet %q among the AWSCluster subnets.", *subnet.SubnetId)
					continue
				}
				if !matchingSubnet.IsPublic {
					errMessage += fmt.Sprintf(" subnet %q is a private subnet.", *subnet.SubnetId)
					continue
				}
			}

			tags := converters.TagsToMap(subnet.Tags)
			if tags[infrav1.NameAWSSubnetAssociation] == infrav1.SecondarySubnetTagValue {
				errMessage += fmt.Sprintf(" subnet %q belongs to a secondary CIDR block which won't be used to create instances.", *subnet.SubnetId)
				continue
			}

			filtered = append(filtered, subnet)
		}
		// keep AWS returned orderz stable, but prefer a subnet in the cluster VPC
		clusterVPC := s.scope.VPC().ID
		sort.SliceStable(filtered, func(i, j int) bool {
			return *filtered[i].VpcId == clusterVPC
		})
		if len(filtered) == 0 {
			errMessage = fmt.Sprintf("failed to run machine %q, found %d subnets matching criteria but post-filtering failed.",
				scope.Name(), len(subnets)) + errMessage
			record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}
		return *filtered[0].SubnetId, nil
	case failureDomain != "":
		if scope.AWSMachine.Spec.PublicIP != nil && *scope.AWSMachine.Spec.PublicIP {
			subnets := s.scope.Subnets().FilterPublic().FilterNonCni().FilterByZone(failureDomain)
			if len(subnets) == 0 {
				errMessage := fmt.Sprintf("failed to run machine %q with public IP, no public subnets available in availability zone %q",
					scope.Name(), failureDomain)
				record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
				return "", awserrors.NewFailedDependency(errMessage)
			}
			return subnets[0].GetResourceID(), nil
		}

		subnets := s.scope.Subnets().FilterPrivate().FilterNonCni().FilterByZone(failureDomain)
		if len(subnets) == 0 {
			errMessage := fmt.Sprintf("failed to run machine %q, no subnets available in availability zone %q",
				scope.Name(), failureDomain)
			record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}
		return subnets[0].GetResourceID(), nil
	case scope.AWSMachine.Spec.PublicIP != nil && *scope.AWSMachine.Spec.PublicIP:
		subnets := s.scope.Subnets().FilterPublic().FilterNonCni()
		if len(subnets) == 0 {
			errMessage := fmt.Sprintf("failed to run machine %q with public IP, no public subnets available", scope.Name())
			record.Eventf(scope.AWSMachine, "FailedCreate", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}
		return subnets[0].GetResourceID(), nil

		// TODO(vincepri): Define a tag that would allow to pick a preferred subnet in an AZ when working
		// with control plane machines.

	default:
		sns := s.scope.Subnets().FilterPrivate().FilterNonCni()
		if len(sns) == 0 {
			errMessage := fmt.Sprintf("failed to run machine %q, no subnets available", scope.Name())
			record.Eventf(s.scope.InfraCluster(), "FailedCreateInstance", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}
		return sns[0].GetResourceID(), nil
	}
}

// getFilteredSubnets fetches subnets filtered based on the criteria passed.
func (s *Service) getFilteredSubnets(criteria ...types.Filter) ([]types.Subnet, error) {
	out, err := s.EC2Client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{Filters: criteria})
	if err != nil {
		return nil, err
	}
	return out.Subnets, nil
}

// GetCoreSecurityGroups looks up the security group IDs managed by this actuator
// They are considered "core" to its proper functioning.
func (s *Service) GetCoreSecurityGroups(scope *scope.MachineScope) ([]string, error) {
	if scope.IsExternallyManaged() {
		ids := make([]string, 0)
		for _, sg := range scope.AWSMachine.Spec.AdditionalSecurityGroups {
			if sg.ID == nil {
				continue
			}
			ids = append(ids, *sg.ID)
		}
		return ids, nil
	}

	// These are common across both controlplane and node machines
	sgRoles := []infrav1.SecurityGroupRole{
		infrav1.SecurityGroupNode,
	}

	if !scope.IsEKSManaged() {
		sgRoles = append(sgRoles, infrav1.SecurityGroupLB)
	}

	switch scope.Role() {
	case "node":
		// Just the common security groups above
		if scope.IsEKSManaged() {
			sgRoles = append(sgRoles, infrav1.SecurityGroupEKSNodeAdditional)
		}
	case "control-plane":
		sgRoles = append(sgRoles, infrav1.SecurityGroupControlPlane)
	default:
		return nil, errors.Errorf("Unknown node role %q", scope.Role())
	}
	ids := make([]string, 0, len(sgRoles))
	for _, sg := range sgRoles {
		if _, ok := scope.AWSMachine.Spec.SecurityGroupOverrides[sg]; ok {
			ids = append(ids, scope.AWSMachine.Spec.SecurityGroupOverrides[sg])
			continue
		}
		if _, ok := s.scope.SecurityGroups()[sg]; ok {
			ids = append(ids, s.scope.SecurityGroups()[sg].ID)
			continue
		}
		return nil, awserrors.NewFailedDependency(fmt.Sprintf("%s security group not available", sg))
	}
	return ids, nil
}

// GetCoreNodeSecurityGroups looks up the security group IDs managed by this actuator
// They are considered "core" to its proper functioning.
func (s *Service) GetCoreNodeSecurityGroups(scope scope.LaunchTemplateScope) ([]string, error) {
	// These are common across both controlplane and node machines
	sgRoles := []infrav1.SecurityGroupRole{
		infrav1.SecurityGroupNode,
	}

	if !scope.IsEKSManaged() {
		sgRoles = append(sgRoles, infrav1.SecurityGroupLB)
	} else {
		sgRoles = append(sgRoles, infrav1.SecurityGroupEKSNodeAdditional)
	}

	ids := make([]string, 0, len(sgRoles))
	for _, sg := range sgRoles {
		if _, ok := s.scope.SecurityGroups()[sg]; !ok {
			return nil, awserrors.NewFailedDependency(
				fmt.Sprintf("%s security group not available", sg),
			)
		}
		ids = append(ids, s.scope.SecurityGroups()[sg].ID)
	}
	return ids, nil
}

// TerminateInstance terminates an EC2 instance.
// Returns nil on success, error in all other cases.
func (s *Service) TerminateInstance(instanceID string) error {
	s.scope.Debug("Attempting to terminate instance", "instance-id", instanceID)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	if _, err := s.EC2Client.TerminateInstances(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to terminate instance with id %q", instanceID)
	}

	s.scope.Debug("Terminated instance", "instance-id", instanceID)
	return nil
}

// TerminateInstanceAndWait terminates and waits
// for an EC2 instance to terminate.
func (s *Service) TerminateInstanceAndWait(instanceID string) error {
	if err := s.TerminateInstance(instanceID); err != nil {
		return err
	}

	s.scope.Debug("Waiting for EC2 instance to terminate", "instance-id", instanceID)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	if err := ec2.NewInstanceTerminatedWaiter(s.EC2Client).Wait(context.TODO(), input, time.Minute*2); err != nil {
		return errors.Wrapf(err, "failed to wait for instance %q termination", instanceID)
	}

	return nil
}

func (s *Service) runInstance(role string, i *infrav1.Instance) (*infrav1.Instance, error) {
	input := &ec2.RunInstancesInput{
		InstanceType: types.InstanceType(i.Type),
		ImageId:      aws.String(i.ImageID),
		KeyName:      i.SSHKeyName,
		EbsOptimized: i.EBSOptimized,
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
		UserData:     i.UserData,
	}

	s.scope.Debug("userData size", "bytes", len(*i.UserData), "role", role)

	if len(i.NetworkInterfaces) > 0 {
		netInterfaces := make([]types.InstanceNetworkInterfaceSpecification, 0, len(i.NetworkInterfaces))

		for index, id := range i.NetworkInterfaces {
			netInterfaces = append(netInterfaces, types.InstanceNetworkInterfaceSpecification{
				NetworkInterfaceId: aws.String(id),
				DeviceIndex:        aws.Int32(int32(index)),
			})
		}
		netInterfaces[0].AssociatePublicIpAddress = i.PublicIPOnLaunch

		input.NetworkInterfaces = netInterfaces
	} else {
		netInterface := types.InstanceNetworkInterfaceSpecification{
			DeviceIndex:              aws.Int32(0),
			SubnetId:                 aws.String(i.SubnetID),
			Groups:                   i.SecurityGroupIDs,
			AssociatePublicIpAddress: i.PublicIPOnLaunch,
		}

		// When registering targets by instance ID for an IPv6 target group, the targets must have an assigned primary IPv6 address.
		// Use case: registering controlplane nodes to the API LBs.
		enablePrimaryIpv6, err := s.shouldEnablePrimaryIpv6(i)
		if err != nil {
			return nil, fmt.Errorf("failed to determine whether to enable PrimaryIpv6 for instance: %w", err)
		}
		if enablePrimaryIpv6 {
			netInterface.PrimaryIpv6 = aws.Bool(true)
			netInterface.Ipv6AddressCount = aws.Int32(1)
		}

		input.NetworkInterfaces = []types.InstanceNetworkInterfaceSpecification{netInterface}
	}

	if i.NetworkInterfaceType != "" {
		input.NetworkInterfaces[0].InterfaceType = aws.String(string(i.NetworkInterfaceType))
	}

	if i.IAMProfile != "" {
		input.IamInstanceProfile = &types.IamInstanceProfileSpecification{
			Name: aws.String(i.IAMProfile),
		}
	}

	blockdeviceMappings := []types.BlockDeviceMapping{}

	if i.RootVolume != nil {
		rootDeviceName, err := s.checkRootVolume(i.RootVolume, i.ImageID)
		if err != nil {
			return nil, err
		}

		i.RootVolume.DeviceName = aws.ToString(rootDeviceName)
		blockDeviceMapping := volumeToBlockDeviceMapping(i.RootVolume)
		blockdeviceMappings = append(blockdeviceMappings, blockDeviceMapping)
	}

	for vi := range i.NonRootVolumes {
		nonRootVolume := i.NonRootVolumes[vi]

		if nonRootVolume.DeviceName == "" {
			return nil, errors.Errorf("non root volume should have device name specified")
		}

		blockDeviceMapping := volumeToBlockDeviceMapping(&nonRootVolume)
		blockdeviceMappings = append(blockdeviceMappings, blockDeviceMapping)
	}

	if len(blockdeviceMappings) != 0 {
		input.BlockDeviceMappings = blockdeviceMappings
	}

	if len(i.Tags) > 0 {
		resources := []types.ResourceType{types.ResourceTypeInstance, types.ResourceTypeVolume}

		if len(i.NetworkInterfaces) == 0 {
			resources = append(resources, types.ResourceTypeNetworkInterface)
		}

		for _, r := range resources {
			spec := types.TagSpecification{ResourceType: r}

			// We need to sort keys for tests to work
			keys := make([]string, 0, len(i.Tags))
			for k := range i.Tags {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				spec.Tags = append(spec.Tags, types.Tag{
					Key:   aws.String(key),
					Value: aws.String(i.Tags[key]),
				})
			}

			input.TagSpecifications = append(input.TagSpecifications, spec)
		}
	}
	marketOptions, err := getInstanceMarketOptionsRequest(i)
	if err != nil {
		return nil, err
	}
	if marketOptions != nil {
		input.InstanceMarketOptions = marketOptions
	}
	input.MetadataOptions = getInstanceMetadataOptionsRequest(i.InstanceMetadataOptions)
	input.PrivateDnsNameOptions = getPrivateDNSNameOptionsRequest(i.PrivateDNSName)
	input.CapacityReservationSpecification = getCapacityReservationSpecification(i.CapacityReservationID, i.CapacityReservationPreference)
	input.CpuOptions = getInstanceCPUOptionsRequest(i.CPUOptions)

	if i.Tenancy != "" {
		input.Placement = &types.Placement{
			Tenancy: types.Tenancy(i.Tenancy),
		}
	}

	if i.PlacementGroupName == "" && i.PlacementGroupPartition != 0 {
		return nil, errors.Errorf("placementGroupPartition is set but placementGroupName is empty")
	}

	if i.PlacementGroupName != "" {
		if input.Placement == nil {
			input.Placement = &types.Placement{}
		}
		input.Placement.GroupName = &i.PlacementGroupName
		if i.PlacementGroupPartition != 0 {
			input.Placement.PartitionNumber = utils.ToInt32Pointer(&i.PlacementGroupPartition)
		}
	}

	if i.HostID != nil {
		if i.HostAffinity == nil {
			// If HostAffinity is not specified, default to "default" Affinity (flexible affinity).
			i.HostAffinity = aws.String("default")
		}
		if len(i.Tenancy) == 0 {
			// If Tenancy is not specified with HostID set, default to "host" Tenancy.
			i.Tenancy = "host"
		}

		s.scope.Debug("Running instance with dedicated host placement",
			"hostId", i.HostID,
			"affinity", i.HostAffinity)
		if input.Placement != nil {
			s.scope.Warn("Placement already set for instance, overwriting with dedicated host placement",
				"hostId", i.HostID,
				"affinity", i.HostAffinity,
				"placement", input.Placement)
		}

		input.Placement = &types.Placement{
			Tenancy:  types.Tenancy(i.Tenancy),
			Affinity: i.HostAffinity,
			HostId:   i.HostID,
		}
	}

	out, err := s.EC2Client.RunInstances(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run instance")
	}

	if len(out.Instances) == 0 {
		return nil, errors.Errorf("no instance returned for reservation %v", out)
	}

	return s.SDKToInstance(out.Instances[0])
}

func volumeToBlockDeviceMapping(v *infrav1.Volume) types.BlockDeviceMapping {
	ebsDevice := &types.EbsBlockDevice{
		DeleteOnTermination: aws.Bool(true),
		VolumeSize:          utils.ToInt32Pointer(&v.Size),
		Encrypted:           v.Encrypted,
	}

	if v.Throughput != nil {
		ebsDevice.Throughput = utils.ToInt32Pointer(v.Throughput)
	}

	if v.IOPS != 0 {
		ebsDevice.Iops = utils.ToInt32Pointer(&v.IOPS)
	}

	if v.EncryptionKey != "" {
		ebsDevice.Encrypted = aws.Bool(true)
		ebsDevice.KmsKeyId = aws.String(v.EncryptionKey)
	}

	if v.Type != "" {
		ebsDevice.VolumeType = types.VolumeType(string(v.Type))
	}

	return types.BlockDeviceMapping{
		DeviceName: &v.DeviceName,
		Ebs:        ebsDevice,
	}
}

// GetInstanceSecurityGroups returns a map from ENI id to the security groups applied to that ENI
// While some security group operations take place at the "instance" level, these are in fact an API convenience for manipulating the first ("primary") ENI's properties.
func (s *Service) GetInstanceSecurityGroups(instanceID string) (map[string][]string, error) {
	enis, err := s.getInstanceENIs(instanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ENIs for instance %q", instanceID)
	}

	out := make(map[string][]string)
	for _, eni := range enis {
		var groups []string
		for _, group := range eni.Groups {
			groups = append(groups, aws.ToString(group.GroupId))
		}
		out[aws.ToString(eni.NetworkInterfaceId)] = groups
	}
	return out, nil
}

// UpdateInstanceSecurityGroups modifies the security groups of the given
// EC2 instance.
func (s *Service) UpdateInstanceSecurityGroups(instanceID string, ids []string) error {
	s.scope.Debug("Attempting to update security groups on instance", "instance-id", instanceID)

	enis, err := s.getInstanceENIs(instanceID)
	if err != nil {
		return errors.Wrapf(err, "failed to get ENIs for instance %q", instanceID)
	}

	s.scope.Debug("Found ENIs on instance", "number-of-enis", len(enis), "instance-id", instanceID)

	for _, eni := range enis {
		if err := s.attachSecurityGroupsToNetworkInterface(ids, aws.ToString(eni.NetworkInterfaceId)); err != nil {
			return errors.Wrapf(err, "failed to modify network interfaces on instance %q", instanceID)
		}
	}

	return nil
}

// UpdateResourceTags updates the tags for an instance.
// This will be called if there is anything to create (update) or delete.
// We may not always have to perform each action, so we check what we're
// receiving to avoid calling AWS if we don't need to.
func (s *Service) UpdateResourceTags(resourceID *string, create, remove map[string]string) error {
	s.scope.Debug("Attempting to update tags on resource", "resource-id", aws.ToString(resourceID))

	// If we have anything to create or update
	if len(create) > 0 {
		s.scope.Debug("Attempting to create tags on resource", "resource-id", aws.ToString(resourceID))

		// Convert our create map into an array of *ec2.Tag
		createTagsInput := converters.MapToTags(create)

		// Create the CreateTags input.
		input := &ec2.CreateTagsInput{
			Resources: []string{aws.ToString(resourceID)},
			Tags:      createTagsInput,
		}

		// Create/Update tags in AWS.
		if _, err := s.EC2Client.CreateTags(context.TODO(), input); err != nil {
			return errors.Wrapf(err, "failed to create tags for resource %q: %+v", aws.ToString(resourceID), create)
		}
	}

	// If we have anything to remove
	if len(remove) > 0 {
		s.scope.Debug("Attempting to delete tags on resource", "resource-id", aws.ToString(resourceID))

		// Convert our remove map into an array of *ec2.Tag
		removeTagsInput := converters.MapToTags(remove)

		// Create the DeleteTags input
		input := &ec2.DeleteTagsInput{
			Resources: []string{aws.ToString(resourceID)},
			Tags:      removeTagsInput,
		}

		// Delete tags in AWS.
		if _, err := s.EC2Client.DeleteTags(context.TODO(), input); err != nil {
			return errors.Wrapf(err, "failed to delete tags for resource %q: %v", aws.ToString(resourceID), remove)
		}
	}

	return nil
}

func (s *Service) getInstanceENIs(instanceID string) ([]types.NetworkInterface, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []string{instanceID},
			},
		},
	}

	output, err := s.EC2Client.DescribeNetworkInterfaces(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return output.NetworkInterfaces, nil
}

func (s *Service) getImageRootDevice(imageID string) (*string, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{imageID},
	}

	output, err := s.EC2Client.DescribeImages(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(output.Images) == 0 {
		return nil, errors.Errorf("no images returned when looking up ID %q", imageID)
	}

	return output.Images[0].RootDeviceName, nil
}

func (s *Service) getImageSnapshotSize(imageID string) (*int32, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{imageID},
	}

	output, err := s.EC2Client.DescribeImages(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(output.Images) == 0 {
		return nil, errors.Errorf("no images returned when looking up ID %q", imageID)
	}

	if len(output.Images[0].BlockDeviceMappings) == 0 {
		return nil, errors.Errorf("no block device mappings returned when looking up ID %q", imageID)
	}

	if output.Images[0].BlockDeviceMappings[0].Ebs == nil {
		return nil, errors.Errorf("no EBS returned when looking up ID %q", imageID)
	}

	if output.Images[0].BlockDeviceMappings[0].Ebs.VolumeSize == nil {
		return nil, errors.Errorf("no EBS volume size returned when looking up ID %q", imageID)
	}

	return output.Images[0].BlockDeviceMappings[0].Ebs.VolumeSize, nil
}

// SDKToInstance converts an AWS EC2 SDK instance to the CAPA instance type.
// SDKToInstance populates all instance fields except for rootVolumeSize,
// because EC2.DescribeInstances does not return the size of storage devices. An
// additional call to EC2 is required to get this value.
func (s *Service) SDKToInstance(v types.Instance) (*infrav1.Instance, error) {
	i := &infrav1.Instance{
		ID:           aws.ToString(v.InstanceId),
		State:        infrav1.InstanceState(string(v.State.Name)),
		Type:         string(v.InstanceType),
		SubnetID:     aws.ToString(v.SubnetId),
		ImageID:      aws.ToString(v.ImageId),
		SSHKeyName:   v.KeyName,
		PrivateIP:    v.PrivateIpAddress,
		IPv6Address:  v.Ipv6Address,
		PublicIP:     v.PublicIpAddress,
		ENASupport:   v.EnaSupport,
		EBSOptimized: v.EbsOptimized,
	}

	// Extract IAM Instance Profile name from ARN
	// TODO: Handle this comparison more safely, perhaps by querying IAM for the
	// instance profile ARN and comparing to the ARN returned by EC2
	if v.IamInstanceProfile != nil && v.IamInstanceProfile.Arn != nil {
		split := strings.Split(aws.ToString(v.IamInstanceProfile.Arn), "instance-profile/")
		if len(split) > 1 && split[1] != "" {
			i.IAMProfile = split[1]
		}
	}

	for _, sg := range v.SecurityGroups {
		i.SecurityGroupIDs = append(i.SecurityGroupIDs, *sg.GroupId)
	}

	if len(v.Tags) > 0 {
		i.Tags = converters.TagsToMap(v.Tags)
	}

	i.Addresses = s.getInstanceAddresses(v)

	// Extract whether the instance has a primary IPv6 assigned
	for _, eni := range v.NetworkInterfaces {
		for _, addr := range eni.Ipv6Addresses {
			if aws.ToBool(addr.IsPrimaryIpv6) {
				enabled := infrav1.PrimaryIPv6AssignmentStateEnabled
				i.AssignPrimaryIPv6 = &enabled
				break
			}
		}
	}

	i.AvailabilityZone = aws.ToString(v.Placement.AvailabilityZone)

	for _, volume := range v.BlockDeviceMappings {
		i.VolumeIDs = append(i.VolumeIDs, *volume.Ebs.VolumeId)
	}

	if v.MetadataOptions != nil {
		metadataOptions := &infrav1.InstanceMetadataOptions{}
		metadataOptions.HTTPEndpoint = infrav1.InstanceMetadataState(string(v.MetadataOptions.HttpEndpoint))
		metadataOptions.HTTPTokens = infrav1.HTTPTokensState(string(v.MetadataOptions.HttpTokens))
		metadataOptions.InstanceMetadataTags = infrav1.InstanceMetadataState(string(v.MetadataOptions.InstanceMetadataTags))
		metadataOptions.HTTPProtocolIPv6 = infrav1.InstanceMetadataState(v.MetadataOptions.HttpProtocolIpv6)
		if v.MetadataOptions.HttpPutResponseHopLimit != nil {
			metadataOptions.HTTPPutResponseHopLimit = int64(*v.MetadataOptions.HttpPutResponseHopLimit)
		}

		i.InstanceMetadataOptions = metadataOptions
	}

	if v.PrivateDnsNameOptions != nil {
		i.PrivateDNSName = &infrav1.PrivateDNSName{
			EnableResourceNameDNSAAAARecord: v.PrivateDnsNameOptions.EnableResourceNameDnsAAAARecord,
			EnableResourceNameDNSARecord:    v.PrivateDnsNameOptions.EnableResourceNameDnsARecord,
			HostnameType:                    aws.String(string(v.PrivateDnsNameOptions.HostnameType)),
		}
	}

	return i, nil
}

func (s *Service) getInstanceAddresses(instance types.Instance) []clusterv1beta1.MachineAddress {
	addresses := []clusterv1beta1.MachineAddress{}
	// Check if the DHCP Option Set has domain name set
	domainName := s.GetDHCPOptionSetDomainName(s.EC2Client, instance.VpcId)
	for _, eni := range instance.NetworkInterfaces {
		if addr := aws.ToString(eni.PrivateDnsName); addr != "" {
			privateDNSAddress := clusterv1beta1.MachineAddress{
				Type:    clusterv1beta1.MachineInternalDNS,
				Address: addr,
			}
			addresses = append(addresses, privateDNSAddress)

			if domainName != nil {
				// Add secondary private DNS Name with domain name set in DHCP Option Set
				additionalPrivateDNSAddress := clusterv1beta1.MachineAddress{
					Type:    clusterv1beta1.MachineInternalDNS,
					Address: fmt.Sprintf("%s.%s", strings.Split(privateDNSAddress.Address, ".")[0], *domainName),
				}
				addresses = append(addresses, additionalPrivateDNSAddress)
			}
		}

		if addr := aws.ToString(eni.PrivateIpAddress); addr != "" {
			privateIPAddress := clusterv1beta1.MachineAddress{
				Type:    clusterv1beta1.MachineInternalIP,
				Address: addr,
			}
			addresses = append(addresses, privateIPAddress)
		}

		// An elastic IP is attached if association is non nil pointer
		if eni.Association != nil {
			if addr := aws.ToString(eni.Association.PublicDnsName); addr != "" {
				publicDNSAddress := clusterv1beta1.MachineAddress{
					Type:    clusterv1beta1.MachineExternalDNS,
					Address: addr,
				}
				addresses = append(addresses, publicDNSAddress)
			}

			if addr := aws.ToString(eni.Association.PublicIp); addr != "" {
				publicIPAddress := clusterv1beta1.MachineAddress{
					Type:    clusterv1beta1.MachineExternalIP,
					Address: addr,
				}
				addresses = append(addresses, publicIPAddress)
			}
		}
	}

	return addresses
}

func (s *Service) getNetworkInterfaceSecurityGroups(interfaceID string) ([]string, error) {
	input := &ec2.DescribeNetworkInterfaceAttributeInput{
		Attribute:          types.NetworkInterfaceAttributeGroupSet,
		NetworkInterfaceId: aws.String(interfaceID),
	}

	output, err := s.EC2Client.DescribeNetworkInterfaceAttribute(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	groups := make([]string, len(output.Groups))
	for i := range output.Groups {
		groups[i] = aws.ToString(output.Groups[i].GroupId)
	}

	return groups, nil
}

func (s *Service) attachSecurityGroupsToNetworkInterface(groups []string, interfaceID string) error {
	s.scope.Info("Updating security groups", "groups", groups)

	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(interfaceID),
		Groups:             groups,
	}

	if _, err := s.EC2Client.ModifyNetworkInterfaceAttribute(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to modify interface %q to have security groups %v", interfaceID, groups)
	}
	return nil
}

// DetachSecurityGroupsFromNetworkInterface looks up an ENI by interfaceID and
// detaches a list of Security Groups from that ENI.
func (s *Service) DetachSecurityGroupsFromNetworkInterface(groups []string, interfaceID string) error {
	existingGroups, err := s.getNetworkInterfaceSecurityGroups(interfaceID)
	if err != nil {
		return errors.Wrapf(err, "failed to look up network interface security groups")
	}

	remainingGroups := existingGroups
	for _, group := range groups {
		remainingGroups = filterGroups(remainingGroups, group)
	}

	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(interfaceID),
		Groups:             remainingGroups,
	}

	if _, err := s.EC2Client.ModifyNetworkInterfaceAttribute(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to modify interface %q", interfaceID)
	}
	return nil
}

// checkRootVolume checks the input root volume options against the requested AMI's defaults
// and returns the AMI's root device name.
func (s *Service) checkRootVolume(rootVolume *infrav1.Volume, imageID string) (*string, error) {
	rootDeviceName, err := s.getImageRootDevice(imageID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get root volume from image %q", imageID)
	}

	snapshotSize, err := s.getImageSnapshotSize(imageID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get root volume from image %q", imageID)
	}

	if rootVolume.Size < int64(*snapshotSize) {
		return nil, errors.Errorf("root volume size (%d) must be greater than or equal to snapshot size (%d)", rootVolume.Size, *snapshotSize)
	}

	return rootDeviceName, nil
}

// ModifyInstanceMetadataOptions modifies the metadata options of the given EC2 instance.
func (s *Service) ModifyInstanceMetadataOptions(instanceID string, options *infrav1.InstanceMetadataOptions) error {
	input := &ec2.ModifyInstanceMetadataOptionsInput{
		HttpEndpoint:            types.InstanceMetadataEndpointState(string(options.HTTPEndpoint)),
		HttpPutResponseHopLimit: utils.ToInt32Pointer(&options.HTTPPutResponseHopLimit),
		HttpTokens:              types.HttpTokensState(string(options.HTTPTokens)),
		InstanceMetadataTags:    types.InstanceMetadataTagsState(string(options.InstanceMetadataTags)),
		HttpProtocolIpv6:        types.InstanceMetadataProtocolState(string(options.HTTPProtocolIPv6)),
		InstanceId:              aws.String(instanceID),
	}

	s.scope.Info("Updating instance metadata options", "instance id", instanceID, "options", input)
	if _, err := s.EC2Client.ModifyInstanceMetadataOptions(context.TODO(), input); err != nil {
		return err
	}

	return nil
}

// GetDHCPOptionSetDomainName returns the domain DNS name for the VPC from the DHCP Options.
func (s *Service) GetDHCPOptionSetDomainName(ec2client common.EC2API, vpcID *string) *string {
	log := s.scope.GetLogger()

	if vpcID == nil {
		log.V(4).Info("vpcID is nil, skipping DHCP Option Set discovery")
		return nil
	}

	vpcInput := &ec2.DescribeVpcsInput{
		VpcIds: []string{aws.ToString(vpcID)},
	}

	vpcResult, err := ec2client.DescribeVpcs(context.TODO(), vpcInput)
	if err != nil {
		log.Info("failed to describe VPC, skipping DHCP Option Set discovery", "vpcID", aws.ToString(vpcID), "Error", err.Error())
		return nil
	}

	dhcpInput := &ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []string{aws.ToString(vpcResult.Vpcs[0].DhcpOptionsId)},
	}

	dhcpResult, err := ec2client.DescribeDhcpOptions(context.TODO(), dhcpInput)
	if err != nil {
		log.Error(err, "failed to describe DHCP Options Set", "input", *dhcpInput)
		return nil
	}

	for _, dhcpConfig := range dhcpResult.DhcpOptions[0].DhcpConfigurations {
		if aws.ToString(dhcpConfig.Key) == "domain-name" {
			if len(dhcpConfig.Values) == 0 {
				return nil
			}
			domainName := dhcpConfig.Values[0].Value
			// default domainName is 'ec2.internal' in us-east-1 and 'region.compute.internal' in the other regions.
			if (s.scope.Region() == "us-east-1" && aws.ToString(domainName) == "ec2.internal") ||
				(s.scope.Region() != "us-east-1" && aws.ToString(domainName) == fmt.Sprintf("%s.compute.internal", s.scope.Region())) {
				return nil
			}

			return domainName
		}
	}

	return nil
}

// filterGroups filters a list for a string.
func filterGroups(list []string, strToFilter string) (newList []string) {
	for _, item := range list {
		if item != strToFilter {
			newList = append(newList, item)
		}
	}
	return
}

func getCapacityReservationSpecification(capacityReservationID *string, capacityReservationPreference infrav1.CapacityReservationPreference) *types.CapacityReservationSpecification {
	if capacityReservationID == nil && capacityReservationPreference == "" {
		return nil
	}
	var spec types.CapacityReservationSpecification
	if capacityReservationID != nil {
		spec.CapacityReservationTarget = &types.CapacityReservationTarget{
			CapacityReservationId: capacityReservationID,
		}
	}
	spec.CapacityReservationPreference = CapacityReservationPreferenceToSDK(capacityReservationPreference)
	return &spec
}

func getInstanceMarketOptionsRequest(i *infrav1.Instance) (*types.InstanceMarketOptionsRequest, error) {
	if i.MarketType != "" && i.MarketType == infrav1.MarketTypeCapacityBlock && i.SpotMarketOptions != nil {
		return nil, errors.New("can't create spot capacity-blocks, remove spot market request")
	}

	if (i.MarketType == infrav1.MarketTypeSpot || i.SpotMarketOptions != nil) && i.CapacityReservationID != nil {
		return nil, errors.New("unable to generate marketOptions for spot instance, capacityReservationID is incompatible with marketType spot and spotMarketOptions")
	}

	// Infer MarketType if not explicitly set
	if i.SpotMarketOptions != nil && i.MarketType == "" {
		i.MarketType = infrav1.MarketTypeSpot
	}

	if i.MarketType == "" {
		i.MarketType = infrav1.MarketTypeOnDemand
	}

	if i.MarketType == infrav1.MarketTypeSpot && i.SpotMarketOptions == nil {
		i.SpotMarketOptions = &infrav1.SpotMarketOptions{}
	}

	switch i.MarketType {
	case infrav1.MarketTypeCapacityBlock:
		if i.CapacityReservationID == nil {
			return nil, errors.Errorf("capacityReservationID is required when CapacityBlock is enabled")
		}
		return &types.InstanceMarketOptionsRequest{
			MarketType: types.MarketTypeCapacityBlock,
		}, nil

	case infrav1.MarketTypeSpot:
		// Set required values for Spot instances
		spotOpts := &types.SpotMarketOptions{
			// The following two options ensure that:
			// - If an instance is interrupted, it is terminated rather than hibernating or stopping
			// - No replacement instance will be created if the instance is interrupted
			// - If the spot request cannot immediately be fulfilled, it will not be created
			// This behaviour should satisfy the 1:1 mapping of Machines to Instances as
			// assumed by the Cluster API.
			InstanceInterruptionBehavior: types.InstanceInterruptionBehaviorTerminate,
			SpotInstanceType:             types.SpotInstanceTypeOneTime,
		}

		if maxPrice := aws.ToString(i.SpotMarketOptions.MaxPrice); maxPrice != "" {
			spotOpts.MaxPrice = aws.String(maxPrice)
		}

		return &types.InstanceMarketOptionsRequest{
			MarketType:  types.MarketTypeSpot,
			SpotOptions: spotOpts,
		}, nil
	case infrav1.MarketTypeOnDemand:
		// Instance is on-demand or empty
		return nil, nil
	default:
		// Invalid MarketType provided
		return nil, errors.Errorf("invalid MarketType %q", i.MarketType)
	}
}

func getInstanceMetadataOptionsRequest(metadataOptions *infrav1.InstanceMetadataOptions) *types.InstanceMetadataOptionsRequest {
	if metadataOptions == nil {
		return nil
	}

	request := &types.InstanceMetadataOptionsRequest{}
	if metadataOptions.HTTPEndpoint != "" {
		request.HttpEndpoint = types.InstanceMetadataEndpointState(string(metadataOptions.HTTPEndpoint))
	}
	if metadataOptions.HTTPProtocolIPv6 != "" {
		request.HttpProtocolIpv6 = types.InstanceMetadataProtocolState(string(metadataOptions.HTTPProtocolIPv6))
	}
	if metadataOptions.HTTPPutResponseHopLimit != 0 {
		request.HttpPutResponseHopLimit = utils.ToInt32Pointer(&metadataOptions.HTTPPutResponseHopLimit)
	}
	if metadataOptions.HTTPTokens != "" {
		request.HttpTokens = types.HttpTokensState(string(metadataOptions.HTTPTokens))
	}
	if metadataOptions.InstanceMetadataTags != "" {
		request.InstanceMetadataTags = types.InstanceMetadataTagsState(string(metadataOptions.InstanceMetadataTags))
	}

	return request
}

// ensureDedicatedHostAllocation ensures a dedicated host is allocated for the machine.
func (s *Service) ensureDedicatedHostAllocation(ctx context.Context, scope *scope.MachineScope) (string, error) {
	spec := scope.AWSMachine.Spec.DynamicHostAllocation
	if spec == nil {
		return "", errors.New("dynamic host allocation spec is nil")
	}

	// Check if a host is already allocated for this machine
	// Each machine gets its own dedicated host for complete isolation and resource dedication
	if scope.AWSMachine.Status.DedicatedHost != nil && scope.AWSMachine.Status.DedicatedHost.ID != nil {
		existingHostID := aws.ToString(scope.AWSMachine.Status.DedicatedHost.ID)
		s.scope.Info("Found existing allocated host for machine", "hostID", existingHostID, "machine", scope.Name())
		return existingHostID, nil
	}

	// Determine the availability zone for the host
	var availabilityZone *string

	// Get AZ from the machine's subnet
	if scope.AWSMachine.Spec.Subnet != nil {
		subnetID, err := s.findSubnet(scope)
		if err != nil {
			return "", errors.Wrap(err, "failed to find subnet for host allocation")
		}

		// Get the full subnet object to extract availability zone
		subnets, err := s.getFilteredSubnets(types.Filter{
			Name:   aws.String("subnet-id"),
			Values: []string{subnetID},
		})
		if err != nil {
			return "", errors.Wrap(err, "failed to get subnet details for host allocation")
		}

		if len(subnets) > 0 && subnets[0].AvailabilityZone != nil {
			availabilityZone = subnets[0].AvailabilityZone
		}
	}

	instanceType := scope.AWSMachine.Spec.InstanceType

	if availabilityZone == nil {
		return "", errors.New("availability zone could not be determined, please specify a subnet ID or subnet filters")
	}

	// Allocate the dedicated host
	hostID, err := s.AllocateDedicatedHost(ctx, spec, instanceType, *availabilityZone, scope)
	if err != nil {
		return "", errors.Wrap(err, "failed to allocate dedicated host")
	}

	s.scope.Info("Successfully allocated dedicated host for machine", "hostID", hostID, "machine", scope.Name())
	return hostID, nil
}

func getPrivateDNSNameOptionsRequest(privateDNSName *infrav1.PrivateDNSName) *types.PrivateDnsNameOptionsRequest {
	if privateDNSName == nil {
		return nil
	}

	return &types.PrivateDnsNameOptionsRequest{
		EnableResourceNameDnsAAAARecord: privateDNSName.EnableResourceNameDNSAAAARecord,
		EnableResourceNameDnsARecord:    privateDNSName.EnableResourceNameDNSARecord,
		HostnameType:                    types.HostnameType(aws.ToString(privateDNSName.HostnameType)),
	}
}

func getInstanceCPUOptionsRequest(cpuOptions infrav1.CPUOptions) *types.CpuOptionsRequest {
	request := &types.CpuOptionsRequest{}
	switch cpuOptions.ConfidentialCompute {
	case infrav1.AWSConfidentialComputePolicySEVSNP:
		request.AmdSevSnp = types.AmdSevSnpSpecificationEnabled
	case infrav1.AWSConfidentialComputePolicyDisabled:
		request.AmdSevSnp = types.AmdSevSnpSpecificationDisabled
	default:
	}

	if *request == (types.CpuOptionsRequest{}) {
		return nil
	}

	return request
}

// shouldEnablePrimaryIpv6 determines whether to enable a primary IPv6 address for an instance.
// This is required when registering instances by ID to IPv6 target groups.
func (s *Service) shouldEnablePrimaryIpv6(i *infrav1.Instance) (bool, error) {
	// We ignore IPv6-related fields when the users do not explicitly enable IPv6 capabilities.
	if !s.scope.VPC().IsIPv6Enabled() {
		// If explicitly set to enabled but VPC doesn't have IPv6 enabled, return error.
		if i.AssignPrimaryIPv6 != nil && *i.AssignPrimaryIPv6 == infrav1.PrimaryIPv6AssignmentStateEnabled {
			return false, fmt.Errorf("cannot enable PrimaryIPv6: VPC does not have IPv6 enabled")
		}
		return false, nil
	}

	// If explicitly set to disabled, return early without checking subnet capabilities.
	if i.AssignPrimaryIPv6 != nil && *i.AssignPrimaryIPv6 == infrav1.PrimaryIPv6AssignmentStateDisabled {
		return false, nil
	}

	// We need to know whether the subnet has IPv6 enabled (i.e. IPv6 only or dual-stack subnet)
	var hasIPv6CIDR bool
	if sn := s.scope.Subnets().FindByID(i.SubnetID); sn != nil {
		hasIPv6CIDR = sn.IsIPv6
	} else {
		// Subnet not in cluster VPC, query AWS API
		sns, err := s.getFilteredSubnets(types.Filter{Name: aws.String("subnet-id"), Values: []string{i.SubnetID}})
		if err != nil {
			return false, fmt.Errorf("failed to find subnet info with id %q for instance: %w", i.SubnetID, err)
		}
		if len(sns) == 0 {
			return false, fmt.Errorf("expected subnet %q for instance to exist, but found none", i.SubnetID)
		}
		if len(sns) > 1 {
			subnetIDs := make([]string, len(sns))
			for i, sn := range sns {
				subnetIDs[i] = aws.ToString(sn.SubnetId)
			}
			return false, fmt.Errorf("expected 1 subnet with id %q, but found %v: %v", i.SubnetID, len(sns), subnetIDs)
		}

		for _, set := range sns[0].Ipv6CidrBlockAssociationSet {
			if set.Ipv6CidrBlockState.State == types.SubnetCidrBlockStateCodeAssociated {
				hasIPv6CIDR = true
				break
			}
		}
	}

	// We should use the value provided by the users if any.
	if i.AssignPrimaryIPv6 != nil {
		// If explicitly set to enabled, validate subnet has IPv6.
		enablePrimaryIPv6 := *i.AssignPrimaryIPv6 == infrav1.PrimaryIPv6AssignmentStateEnabled
		if enablePrimaryIPv6 && !hasIPv6CIDR {
			return false, fmt.Errorf("cannot enable PrimaryIPv6: subnet %s does not have IPv6 CIDR block", i.SubnetID)
		}
		return enablePrimaryIPv6, nil
	}

	// Otherwise, we define the default behavior as follows:
	// - disabled if subnet is ipv4 only
	// - enabled if subnet is ipv6 only or dual-stack
	return hasIPv6CIDR, nil
}
