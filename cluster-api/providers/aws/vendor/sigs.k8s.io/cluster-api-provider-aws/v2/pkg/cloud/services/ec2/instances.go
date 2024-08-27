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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

// GetRunningInstanceByTags returns the existing instance or nothing if it doesn't exist.
func (s *Service) GetRunningInstanceByTags(scope *scope.MachineScope) (*infrav1.Instance, error) {
	s.scope.Debug("Looking for existing machine instance by tags")

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.ClusterOwned(s.scope.Name()),
			filter.EC2.Name(scope.Name()),
			filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
		},
	}

	out, err := s.EC2Client.DescribeInstancesWithContext(context.TODO(), input)
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
		InstanceIds: []*string{id},
	}

	out, err := s.EC2Client.DescribeInstancesWithContext(context.TODO(), input)
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
func (s *Service) CreateInstance(scope *scope.MachineScope, userData []byte, userDataFormat string) (*infrav1.Instance, error) {
	s.scope.Debug("Creating an instance for a machine")

	input := &infrav1.Instance{
		Type:              scope.AWSMachine.Spec.InstanceType,
		IAMProfile:        scope.AWSMachine.Spec.IAMInstanceProfile,
		RootVolume:        scope.AWSMachine.Spec.RootVolume.DeepCopy(),
		NonRootVolumes:    scope.AWSMachine.Spec.NonRootVolumes,
		NetworkInterfaces: scope.AWSMachine.Spec.NetworkInterfaces,
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

	imageArchitecture, err := s.pickArchitectureForInstanceType(input.Type)
	if err != nil {
		return nil, err
	}

	// Pick image from the machine configuration, or use a default one.
	if scope.AWSMachine.Spec.AMI.ID != nil { //nolint:nestif
		input.ImageID = *scope.AWSMachine.Spec.AMI.ID
	} else {
		if scope.Machine.Spec.Version == nil {
			err := errors.New("Either AWSMachine's spec.ami.id or Machine's spec.version must be defined")
			scope.SetFailureReason(capierrors.CreateMachineError)
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
			input.ImageID, err = s.eksAMILookup(*scope.Machine.Spec.Version, imageArchitecture, scope.AWSMachine.Spec.AMI.EKSOptimizedLookupType)
			if err != nil {
				return nil, err
			}
		} else {
			input.ImageID, err = s.defaultAMIIDLookup(imageLookupFormat, imageLookupOrg, imageLookupBaseOS, imageArchitecture, *scope.Machine.Spec.Version)
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
		criteria := []*ec2.Filter{
			filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
		}
		if scope.AWSMachine.Spec.Subnet.ID != nil {
			criteria = append(criteria, &ec2.Filter{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{*scope.AWSMachine.Spec.Subnet.ID})})
		}
		for _, f := range scope.AWSMachine.Spec.Subnet.Filters {
			criteria = append(criteria, &ec2.Filter{Name: aws.String(f.Name), Values: aws.StringSlice(f.Values)})
		}

		subnets, err := s.getFilteredSubnets(criteria...)
		if err != nil {
			return "", errors.Wrapf(err, "failed to filter subnets for criteria %q", criteria)
		}
		if len(subnets) == 0 {
			errMessage := fmt.Sprintf("failed to run machine %q, no subnets available matching criteria %q",
				scope.Name(), criteria)
			record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}

		var filtered []*ec2.Subnet
		var errMessage string
		for _, subnet := range subnets {
			if failureDomain != nil && *subnet.AvailabilityZone != *failureDomain {
				// we could have included the failure domain in the query criteria, but then we end up with EC2 error
				// messages that don't give a good hint about what is really wrong
				errMessage += fmt.Sprintf(" subnet %q availability zone %q does not match failure domain %q.",
					*subnet.SubnetId, *subnet.AvailabilityZone, *failureDomain)
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
		// prefer a subnet in the cluster VPC if multiple match
		clusterVPC := s.scope.VPC().ID
		sort.SliceStable(filtered, func(i, j int) bool {
			return strings.Compare(*filtered[i].VpcId, clusterVPC) > strings.Compare(*filtered[j].VpcId, clusterVPC)
		})
		if len(filtered) == 0 {
			errMessage = fmt.Sprintf("failed to run machine %q, found %d subnets matching criteria but post-filtering failed.",
				scope.Name(), len(subnets)) + errMessage
			record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
			return "", awserrors.NewFailedDependency(errMessage)
		}
		return *filtered[0].SubnetId, nil
	case failureDomain != nil:
		if scope.AWSMachine.Spec.PublicIP != nil && *scope.AWSMachine.Spec.PublicIP {
			subnets := s.scope.Subnets().FilterPublic().FilterNonCni().FilterByZone(*failureDomain)
			if len(subnets) == 0 {
				errMessage := fmt.Sprintf("failed to run machine %q with public IP, no public subnets available in availability zone %q",
					scope.Name(), *failureDomain)
				record.Warnf(scope.AWSMachine, "FailedCreate", errMessage)
				return "", awserrors.NewFailedDependency(errMessage)
			}
			return subnets[0].GetResourceID(), nil
		}

		subnets := s.scope.Subnets().FilterPrivate().FilterNonCni().FilterByZone(*failureDomain)
		if len(subnets) == 0 {
			errMessage := fmt.Sprintf("failed to run machine %q, no subnets available in availability zone %q",
				scope.Name(), *failureDomain)
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
func (s *Service) getFilteredSubnets(criteria ...*ec2.Filter) ([]*ec2.Subnet, error) {
	out, err := s.EC2Client.DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{Filters: criteria})
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
		InstanceIds: aws.StringSlice([]string{instanceID}),
	}

	if _, err := s.EC2Client.TerminateInstancesWithContext(context.TODO(), input); err != nil {
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
		InstanceIds: aws.StringSlice([]string{instanceID}),
	}

	if err := s.EC2Client.WaitUntilInstanceTerminatedWithContext(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to wait for instance %q termination", instanceID)
	}

	return nil
}

func (s *Service) runInstance(role string, i *infrav1.Instance) (*infrav1.Instance, error) {
	input := &ec2.RunInstancesInput{
		InstanceType: aws.String(i.Type),
		ImageId:      aws.String(i.ImageID),
		KeyName:      i.SSHKeyName,
		EbsOptimized: i.EBSOptimized,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		UserData:     i.UserData,
	}

	s.scope.Debug("userData size", "bytes", len(*i.UserData), "role", role)

	if len(i.NetworkInterfaces) > 0 {
		netInterfaces := make([]*ec2.InstanceNetworkInterfaceSpecification, 0, len(i.NetworkInterfaces))

		for index, id := range i.NetworkInterfaces {
			netInterfaces = append(netInterfaces, &ec2.InstanceNetworkInterfaceSpecification{
				NetworkInterfaceId: aws.String(id),
				DeviceIndex:        aws.Int64(int64(index)),
			})
		}
		netInterfaces[0].AssociatePublicIpAddress = i.PublicIPOnLaunch

		input.NetworkInterfaces = netInterfaces
	} else {
		input.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int64(0),
				SubnetId:                 aws.String(i.SubnetID),
				Groups:                   aws.StringSlice(i.SecurityGroupIDs),
				AssociatePublicIpAddress: i.PublicIPOnLaunch,
			},
		}
	}

	if i.IAMProfile != "" {
		input.IamInstanceProfile = &ec2.IamInstanceProfileSpecification{
			Name: aws.String(i.IAMProfile),
		}
	}

	blockdeviceMappings := []*ec2.BlockDeviceMapping{}

	if i.RootVolume != nil {
		rootDeviceName, err := s.checkRootVolume(i.RootVolume, i.ImageID)
		if err != nil {
			return nil, err
		}

		i.RootVolume.DeviceName = aws.StringValue(rootDeviceName)
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
		resources := []string{ec2.ResourceTypeInstance, ec2.ResourceTypeVolume, ec2.ResourceTypeNetworkInterface}
		for _, r := range resources {
			spec := &ec2.TagSpecification{ResourceType: aws.String(r)}

			// We need to sort keys for tests to work
			keys := make([]string, 0, len(i.Tags))
			for k := range i.Tags {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				spec.Tags = append(spec.Tags, &ec2.Tag{
					Key:   aws.String(key),
					Value: aws.String(i.Tags[key]),
				})
			}

			input.TagSpecifications = append(input.TagSpecifications, spec)
		}
	}

	input.InstanceMarketOptions = getInstanceMarketOptionsRequest(i.SpotMarketOptions)
	input.MetadataOptions = getInstanceMetadataOptionsRequest(i.InstanceMetadataOptions)
	input.PrivateDnsNameOptions = getPrivateDNSNameOptionsRequest(i.PrivateDNSName)
	input.CapacityReservationSpecification = getCapacityReservationSpecification(i.CapacityReservationID)

	if i.Tenancy != "" {
		input.Placement = &ec2.Placement{
			Tenancy: &i.Tenancy,
		}
	}

	if i.PlacementGroupName == "" && i.PlacementGroupPartition != 0 {
		return nil, errors.Errorf("placementGroupPartition is set but placementGroupName is empty")
	}

	if i.PlacementGroupName != "" {
		if input.Placement == nil {
			input.Placement = &ec2.Placement{}
		}
		input.Placement.GroupName = &i.PlacementGroupName
		if i.PlacementGroupPartition != 0 {
			input.Placement.PartitionNumber = &i.PlacementGroupPartition
		}
	}

	out, err := s.EC2Client.RunInstancesWithContext(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run instance")
	}

	if len(out.Instances) == 0 {
		return nil, errors.Errorf("no instance returned for reservation %v", out.GoString())
	}

	return s.SDKToInstance(out.Instances[0])
}

func volumeToBlockDeviceMapping(v *infrav1.Volume) *ec2.BlockDeviceMapping {
	ebsDevice := &ec2.EbsBlockDevice{
		DeleteOnTermination: aws.Bool(true),
		VolumeSize:          aws.Int64(v.Size),
		Encrypted:           v.Encrypted,
	}

	if v.Throughput != nil {
		ebsDevice.Throughput = v.Throughput
	}

	if v.IOPS != 0 {
		ebsDevice.Iops = aws.Int64(v.IOPS)
	}

	if v.EncryptionKey != "" {
		ebsDevice.Encrypted = aws.Bool(true)
		ebsDevice.KmsKeyId = aws.String(v.EncryptionKey)
	}

	if v.Type != "" {
		ebsDevice.VolumeType = aws.String(string(v.Type))
	}

	return &ec2.BlockDeviceMapping{
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
			groups = append(groups, aws.StringValue(group.GroupId))
		}
		out[aws.StringValue(eni.NetworkInterfaceId)] = groups
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
		if err := s.attachSecurityGroupsToNetworkInterface(ids, aws.StringValue(eni.NetworkInterfaceId)); err != nil {
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
	s.scope.Debug("Attempting to update tags on resource", "resource-id", *resourceID)

	// If we have anything to create or update
	if len(create) > 0 {
		s.scope.Debug("Attempting to create tags on resource", "resource-id", *resourceID)

		// Convert our create map into an array of *ec2.Tag
		createTagsInput := converters.MapToTags(create)

		// Create the CreateTags input.
		input := &ec2.CreateTagsInput{
			Resources: []*string{resourceID},
			Tags:      createTagsInput,
		}

		// Create/Update tags in AWS.
		if _, err := s.EC2Client.CreateTagsWithContext(context.TODO(), input); err != nil {
			return errors.Wrapf(err, "failed to create tags for resource %q: %+v", *resourceID, create)
		}
	}

	// If we have anything to remove
	if len(remove) > 0 {
		s.scope.Debug("Attempting to delete tags on resource", "resource-id", *resourceID)

		// Convert our remove map into an array of *ec2.Tag
		removeTagsInput := converters.MapToTags(remove)

		// Create the DeleteTags input
		input := &ec2.DeleteTagsInput{
			Resources: []*string{resourceID},
			Tags:      removeTagsInput,
		}

		// Delete tags in AWS.
		if _, err := s.EC2Client.DeleteTagsWithContext(context.TODO(), input); err != nil {
			return errors.Wrapf(err, "failed to delete tags for resource %q: %v", *resourceID, remove)
		}
	}

	return nil
}

func (s *Service) getInstanceENIs(instanceID string) ([]*ec2.NetworkInterface, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []*string{aws.String(instanceID)},
			},
		},
	}

	output, err := s.EC2Client.DescribeNetworkInterfacesWithContext(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return output.NetworkInterfaces, nil
}

func (s *Service) getImageRootDevice(imageID string) (*string, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(imageID)},
	}

	output, err := s.EC2Client.DescribeImagesWithContext(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(output.Images) == 0 {
		return nil, errors.Errorf("no images returned when looking up ID %q", imageID)
	}

	return output.Images[0].RootDeviceName, nil
}

func (s *Service) getImageSnapshotSize(imageID string) (*int64, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(imageID)},
	}

	output, err := s.EC2Client.DescribeImagesWithContext(context.TODO(), input)
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
func (s *Service) SDKToInstance(v *ec2.Instance) (*infrav1.Instance, error) {
	i := &infrav1.Instance{
		ID:           aws.StringValue(v.InstanceId),
		State:        infrav1.InstanceState(*v.State.Name),
		Type:         aws.StringValue(v.InstanceType),
		SubnetID:     aws.StringValue(v.SubnetId),
		ImageID:      aws.StringValue(v.ImageId),
		SSHKeyName:   v.KeyName,
		PrivateIP:    v.PrivateIpAddress,
		PublicIP:     v.PublicIpAddress,
		ENASupport:   v.EnaSupport,
		EBSOptimized: v.EbsOptimized,
	}

	// Extract IAM Instance Profile name from ARN
	// TODO: Handle this comparison more safely, perhaps by querying IAM for the
	// instance profile ARN and comparing to the ARN returned by EC2
	if v.IamInstanceProfile != nil && v.IamInstanceProfile.Arn != nil {
		split := strings.Split(aws.StringValue(v.IamInstanceProfile.Arn), "instance-profile/")
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

	i.AvailabilityZone = aws.StringValue(v.Placement.AvailabilityZone)

	for _, volume := range v.BlockDeviceMappings {
		i.VolumeIDs = append(i.VolumeIDs, *volume.Ebs.VolumeId)
	}

	if v.MetadataOptions != nil {
		metadataOptions := &infrav1.InstanceMetadataOptions{}
		if v.MetadataOptions.HttpEndpoint != nil {
			metadataOptions.HTTPEndpoint = infrav1.InstanceMetadataState(*v.MetadataOptions.HttpEndpoint)
		}
		if v.MetadataOptions.HttpPutResponseHopLimit != nil {
			metadataOptions.HTTPPutResponseHopLimit = *v.MetadataOptions.HttpPutResponseHopLimit
		}
		if v.MetadataOptions.HttpTokens != nil {
			metadataOptions.HTTPTokens = infrav1.HTTPTokensState(*v.MetadataOptions.HttpTokens)
		}
		if v.MetadataOptions.InstanceMetadataTags != nil {
			metadataOptions.InstanceMetadataTags = infrav1.InstanceMetadataState(*v.MetadataOptions.InstanceMetadataTags)
		}

		i.InstanceMetadataOptions = metadataOptions
	}

	if v.PrivateDnsNameOptions != nil {
		i.PrivateDNSName = &infrav1.PrivateDNSName{
			EnableResourceNameDNSAAAARecord: v.PrivateDnsNameOptions.EnableResourceNameDnsAAAARecord,
			EnableResourceNameDNSARecord:    v.PrivateDnsNameOptions.EnableResourceNameDnsARecord,
			HostnameType:                    v.PrivateDnsNameOptions.HostnameType,
		}
	}

	return i, nil
}

func (s *Service) getInstanceAddresses(instance *ec2.Instance) []clusterv1.MachineAddress {
	addresses := []clusterv1.MachineAddress{}
	// Check if the DHCP Option Set has domain name set
	domainName := s.GetDHCPOptionSetDomainName(s.EC2Client, instance.VpcId)
	for _, eni := range instance.NetworkInterfaces {
		privateDNSAddress := clusterv1.MachineAddress{
			Type:    clusterv1.MachineInternalDNS,
			Address: aws.StringValue(eni.PrivateDnsName),
		}
		privateIPAddress := clusterv1.MachineAddress{
			Type:    clusterv1.MachineInternalIP,
			Address: aws.StringValue(eni.PrivateIpAddress),
		}

		addresses = append(addresses, privateDNSAddress, privateIPAddress)

		if domainName != nil {
			// Add secondary private DNS Name with domain name set in DHCP Option Set
			additionalPrivateDNSAddress := clusterv1.MachineAddress{
				Type:    clusterv1.MachineInternalDNS,
				Address: fmt.Sprintf("%s.%s", strings.Split(privateDNSAddress.Address, ".")[0], *domainName),
			}
			addresses = append(addresses, additionalPrivateDNSAddress)
		}

		// An elastic IP is attached if association is non nil pointer
		if eni.Association != nil {
			publicDNSAddress := clusterv1.MachineAddress{
				Type:    clusterv1.MachineExternalDNS,
				Address: aws.StringValue(eni.Association.PublicDnsName),
			}
			publicIPAddress := clusterv1.MachineAddress{
				Type:    clusterv1.MachineExternalIP,
				Address: aws.StringValue(eni.Association.PublicIp),
			}
			addresses = append(addresses, publicDNSAddress, publicIPAddress)
		}
	}

	return addresses
}

func (s *Service) getNetworkInterfaceSecurityGroups(interfaceID string) ([]string, error) {
	input := &ec2.DescribeNetworkInterfaceAttributeInput{
		Attribute:          aws.String("groupSet"),
		NetworkInterfaceId: aws.String(interfaceID),
	}

	output, err := s.EC2Client.DescribeNetworkInterfaceAttributeWithContext(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	groups := make([]string, len(output.Groups))
	for i := range output.Groups {
		groups[i] = aws.StringValue(output.Groups[i].GroupId)
	}

	return groups, nil
}

func (s *Service) attachSecurityGroupsToNetworkInterface(groups []string, interfaceID string) error {
	s.scope.Info("Updating security groups", "groups", groups)

	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(interfaceID),
		Groups:             aws.StringSlice(groups),
	}

	if _, err := s.EC2Client.ModifyNetworkInterfaceAttributeWithContext(context.TODO(), input); err != nil {
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
		Groups:             aws.StringSlice(remainingGroups),
	}

	if _, err := s.EC2Client.ModifyNetworkInterfaceAttributeWithContext(context.TODO(), input); err != nil {
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

	if rootVolume.Size < *snapshotSize {
		return nil, errors.Errorf("root volume size (%d) must be greater than or equal to snapshot size (%d)", rootVolume.Size, *snapshotSize)
	}

	return rootDeviceName, nil
}

// ModifyInstanceMetadataOptions modifies the metadata options of the given EC2 instance.
func (s *Service) ModifyInstanceMetadataOptions(instanceID string, options *infrav1.InstanceMetadataOptions) error {
	input := &ec2.ModifyInstanceMetadataOptionsInput{
		HttpEndpoint:            aws.String(string(options.HTTPEndpoint)),
		HttpPutResponseHopLimit: aws.Int64(options.HTTPPutResponseHopLimit),
		HttpTokens:              aws.String(string(options.HTTPTokens)),
		InstanceMetadataTags:    aws.String(string(options.InstanceMetadataTags)),
		InstanceId:              aws.String(instanceID),
	}

	s.scope.Info("Updating instance metadata options", "instance id", instanceID, "options", input)
	if _, err := s.EC2Client.ModifyInstanceMetadataOptionsWithContext(context.TODO(), input); err != nil {
		return err
	}

	return nil
}

// GetDHCPOptionSetDomainName returns the domain DNS name for the VPC from the DHCP Options.
func (s *Service) GetDHCPOptionSetDomainName(ec2client ec2iface.EC2API, vpcID *string) *string {
	log := s.scope.GetLogger()

	if vpcID == nil {
		log.Info("vpcID is nil, skipping DHCP Option Set discovery")
		return nil
	}

	vpcInput := &ec2.DescribeVpcsInput{
		VpcIds: []*string{vpcID},
	}

	vpcResult, err := ec2client.DescribeVpcs(vpcInput)
	if err != nil {
		log.Info("failed to describe VPC, skipping DHCP Option Set discovery", "vpcID", *vpcID, "Error", err.Error())
		return nil
	}

	dhcpInput := &ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []*string{vpcResult.Vpcs[0].DhcpOptionsId},
	}

	dhcpResult, err := ec2client.DescribeDhcpOptions(dhcpInput)
	if err != nil {
		log.Error(err, "failed to describe DHCP Options Set", "DhcpOptionsSet", *dhcpResult)
		return nil
	}

	for _, dhcpConfig := range dhcpResult.DhcpOptions[0].DhcpConfigurations {
		if *dhcpConfig.Key == "domain-name" {
			if len(dhcpConfig.Values) == 0 {
				return nil
			}
			domainName := dhcpConfig.Values[0].Value
			// default domainName is 'ec2.internal' in us-east-1 and 'region.compute.internal' in the other regions.
			if (s.scope.Region() == "us-east-1" && *domainName == "ec2.internal") ||
				(s.scope.Region() != "us-east-1" && *domainName == fmt.Sprintf("%s.compute.internal", s.scope.Region())) {
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

func getCapacityReservationSpecification(capacityReservationID *string) *ec2.CapacityReservationSpecification {
	if capacityReservationID == nil {
		//  Not targeting any specific Capacity Reservation
		return nil
	}

	return &ec2.CapacityReservationSpecification{
		CapacityReservationTarget: &ec2.CapacityReservationTarget{
			CapacityReservationId: capacityReservationID,
		},
	}
}

func getInstanceMarketOptionsRequest(spotMarketOptions *infrav1.SpotMarketOptions) *ec2.InstanceMarketOptionsRequest {
	if spotMarketOptions == nil {
		// Instance is not a Spot instance
		return nil
	}

	// Set required values for Spot instances
	spotOptions := &ec2.SpotMarketOptions{}

	// The following two options ensure that:
	// - If an instance is interrupted, it is terminated rather than hibernating or stopping
	// - No replacement instance will be created if the instance is interrupted
	// - If the spot request cannot immediately be fulfilled, it will not be created
	// This behaviour should satisfy the 1:1 mapping of Machines to Instances as
	// assumed by the Cluster API.
	spotOptions.SetInstanceInterruptionBehavior(ec2.InstanceInterruptionBehaviorTerminate)
	spotOptions.SetSpotInstanceType(ec2.SpotInstanceTypeOneTime)

	maxPrice := spotMarketOptions.MaxPrice
	if maxPrice != nil && *maxPrice != "" {
		spotOptions.SetMaxPrice(*maxPrice)
	}

	instanceMarketOptionsRequest := &ec2.InstanceMarketOptionsRequest{}
	instanceMarketOptionsRequest.SetMarketType(ec2.MarketTypeSpot)
	instanceMarketOptionsRequest.SetSpotOptions(spotOptions)

	return instanceMarketOptionsRequest
}

func getInstanceMetadataOptionsRequest(metadataOptions *infrav1.InstanceMetadataOptions) *ec2.InstanceMetadataOptionsRequest {
	if metadataOptions == nil {
		return nil
	}

	request := &ec2.InstanceMetadataOptionsRequest{}
	if metadataOptions.HTTPEndpoint != "" {
		request.SetHttpEndpoint(string(metadataOptions.HTTPEndpoint))
	}
	if metadataOptions.HTTPPutResponseHopLimit != 0 {
		request.SetHttpPutResponseHopLimit(metadataOptions.HTTPPutResponseHopLimit)
	}
	if metadataOptions.HTTPTokens != "" {
		request.SetHttpTokens(string(metadataOptions.HTTPTokens))
	}
	if metadataOptions.InstanceMetadataTags != "" {
		request.SetInstanceMetadataTags(string(metadataOptions.InstanceMetadataTags))
	}

	return request
}

func getPrivateDNSNameOptionsRequest(privateDNSName *infrav1.PrivateDNSName) *ec2.PrivateDnsNameOptionsRequest {
	if privateDNSName == nil {
		return nil
	}

	return &ec2.PrivateDnsNameOptionsRequest{
		EnableResourceNameDnsAAAARecord: privateDNSName.EnableResourceNameDNSAAAARecord,
		EnableResourceNameDnsARecord:    privateDNSName.EnableResourceNameDNSARecord,
		HostnameType:                    privateDNSName.HostnameType,
	}
}
