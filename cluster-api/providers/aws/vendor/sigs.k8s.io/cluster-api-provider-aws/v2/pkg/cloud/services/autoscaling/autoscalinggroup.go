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

// Package asg provides a service for managing AWS AutoScalingGroups.
package asg

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/utils"
	"sigs.k8s.io/cluster-api/util/annotations"
)

// SDKToAutoScalingGroup converts an AWS EC2 SDK AutoScalingGroup to the CAPA AutoScalingGroup type.
func (s *Service) SDKToAutoScalingGroup(v *autoscalingtypes.AutoScalingGroup) (*expinfrav1.AutoScalingGroup, error) {
	i := &expinfrav1.AutoScalingGroup{
		ID:   aws.StringValue(v.AutoScalingGroupARN),
		Name: aws.StringValue(v.AutoScalingGroupName),
		// TODO(rudoi): this is just terrible
		DesiredCapacity:   v.DesiredCapacity,
		MaxSize:           aws.Int32Value(v.MaxSize), //#nosec G115
		MinSize:           aws.Int32Value(v.MinSize), //#nosec G115
		CapacityRebalance: aws.BoolValue(v.CapacityRebalance),
		// TODO: determine what additional values go here and what else should be in the struct
	}

	if v.VPCZoneIdentifier != nil {
		i.Subnets = strings.Split(*v.VPCZoneIdentifier, ",")
	}

	if v.MixedInstancesPolicy != nil {
		i.MixedInstancesPolicy = &expinfrav1.MixedInstancesPolicy{
			InstancesDistribution: &expinfrav1.InstancesDistribution{
				OnDemandBaseCapacity:                utils.ToInt64Pointer(v.MixedInstancesPolicy.InstancesDistribution.OnDemandBaseCapacity),
				OnDemandPercentageAboveBaseCapacity: utils.ToInt64Pointer(v.MixedInstancesPolicy.InstancesDistribution.OnDemandPercentageAboveBaseCapacity),
			},
		}

		for _, override := range v.MixedInstancesPolicy.LaunchTemplate.Overrides {
			i.MixedInstancesPolicy.Overrides = append(i.MixedInstancesPolicy.Overrides, expinfrav1.Overrides{InstanceType: aws.StringValue(override.InstanceType)})
		}

		onDemandAllocationStrategy := aws.StringValue(v.MixedInstancesPolicy.InstancesDistribution.OnDemandAllocationStrategy)
		switch onDemandAllocationStrategy {
		case string(expinfrav1.OnDemandAllocationStrategyPrioritized):
			i.MixedInstancesPolicy.InstancesDistribution.OnDemandAllocationStrategy = expinfrav1.OnDemandAllocationStrategyPrioritized
		case string(expinfrav1.OnDemandAllocationStrategyLowestPrice):
			i.MixedInstancesPolicy.InstancesDistribution.OnDemandAllocationStrategy = expinfrav1.OnDemandAllocationStrategyLowestPrice
		default:
			return nil, fmt.Errorf("unsupported on-demand allocation strategy: %s", onDemandAllocationStrategy)
		}

		spotAllocationStrategy := aws.StringValue(v.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy)
		switch spotAllocationStrategy {
		case string(expinfrav1.SpotAllocationStrategyLowestPrice):
			i.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy = expinfrav1.SpotAllocationStrategyLowestPrice
		case string(expinfrav1.SpotAllocationStrategyCapacityOptimized):
			i.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy = expinfrav1.SpotAllocationStrategyCapacityOptimized
		case string(expinfrav1.SpotAllocationStrategyCapacityOptimizedPrioritized):
			i.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy = expinfrav1.SpotAllocationStrategyCapacityOptimizedPrioritized
		case string(expinfrav1.SpotAllocationStrategyPriceCapacityOptimized):
			i.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy = expinfrav1.SpotAllocationStrategyPriceCapacityOptimized
		default:
			return nil, fmt.Errorf("unsupported spot allocation strategy: %s", spotAllocationStrategy)
		}
	}

	if v.Status != nil {
		i.Status = expinfrav1.ASGStatus(*v.Status)
	}

	if len(v.Tags) > 0 {
		i.Tags = converters.ASGTagsToMap(v.Tags)
	}

	if len(v.Instances) > 0 {
		for _, autoscalingInstance := range v.Instances {
			tmp := &infrav1.Instance{
				ID:               aws.StringValue(autoscalingInstance.InstanceId),
				State:            infrav1.InstanceState(autoscalingInstance.LifecycleState),
				AvailabilityZone: *autoscalingInstance.AvailabilityZone,
			}
			i.Instances = append(i.Instances, *tmp)
		}
	}

	if len(v.SuspendedProcesses) > 0 {
		currentlySuspendedProcesses := make([]string, len(v.SuspendedProcesses))
		for i, service := range v.SuspendedProcesses {
			currentlySuspendedProcesses[i] = aws.StringValue(service.ProcessName)
		}
		i.CurrentlySuspendProcesses = currentlySuspendedProcesses
	}

	return i, nil
}

// ASGIfExists returns the existing autoscaling group or nothing if it doesn't exist.
func (s *Service) ASGIfExists(name *string) (*expinfrav1.AutoScalingGroup, error) {
	if name == nil {
		s.scope.Info("Autoscaling Group does not have a name")
		return nil, nil
	}

	s.scope.Info("Looking for asg by name", "name", *name)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{*name},
	}

	out, err := s.ASGClient.DescribeAutoScalingGroups(context.TODO(), input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeAutoScalingGroups", "failed to describe ASG %q: %v", *name, err)
		return nil, errors.Wrapf(err, "failed to describe AutoScaling Group: %q", *name)
	case len(out.AutoScalingGroups) == 0:
		record.Eventf(s.scope.InfraCluster(), corev1.EventTypeNormal, expinfrav1.ASGNotFoundReason, "Unable to find ASG matching %q", *name)
		return nil, nil
	}
	return s.SDKToAutoScalingGroup(&out.AutoScalingGroups[0])
}

// GetASGByName returns the existing ASG or nothing if it doesn't exist.
func (s *Service) GetASGByName(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
	name := scope.Name()
	return s.ASGIfExists(&name)
}

// CreateASG runs an autoscaling group.
func (s *Service) CreateASG(machinePoolScope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
	subnets, err := s.SubnetIDs(machinePoolScope)
	if err != nil {
		return nil, fmt.Errorf("getting subnets for ASG: %w", err)
	}

	name := machinePoolScope.Name()
	s.scope.Info("Creating ASG", "name", name)

	// Default value of MachinePool replicas set by CAPI is 1.
	mpReplicas := *machinePoolScope.MachinePool.Spec.Replicas
	var desiredCapacity *int32

	// Check that MachinePool replicas number is between the minimum and maximum size of the AWSMachinePool.
	// Ignore the problem for externally managed clusters because MachinePool replicas will be updated to the right value automatically.
	if mpReplicas >= machinePoolScope.AWSMachinePool.Spec.MinSize && mpReplicas <= machinePoolScope.AWSMachinePool.Spec.MaxSize {
		desiredCapacity = &mpReplicas
	} else if !annotations.ReplicasManagedByExternalAutoscaler(machinePoolScope.MachinePool) {
		return nil, fmt.Errorf("incorrect number of replicas %d in MachinePool %v", mpReplicas, machinePoolScope.MachinePool.Name)
	}

	if machinePoolScope.AWSMachinePool.Status.LaunchTemplateID == "" {
		return nil, errors.New("AWSMachinePool has no LaunchTemplateID for some reason")
	}

	// Make sure to use the MachinePoolScope here to get the merger of AWSCluster and AWSMachinePool tags
	additionalTags := machinePoolScope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())] = string(infrav1.ResourceLifecycleOwned)

	input := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName:           aws.String(name),
		MaxSize:                        aws.Int32(machinePoolScope.AWSMachinePool.Spec.MaxSize),
		MinSize:                        aws.Int32(machinePoolScope.AWSMachinePool.Spec.MinSize),
		VPCZoneIdentifier:              aws.String(strings.Join(subnets, ", ")),
		DefaultCooldown:                aws.Int32(int32(machinePoolScope.AWSMachinePool.Spec.DefaultCoolDown.Duration.Seconds())),
		DefaultInstanceWarmup:          aws.Int32(int32(machinePoolScope.AWSMachinePool.Spec.DefaultInstanceWarmup.Duration.Seconds())),
		CapacityRebalance:              aws.Bool(machinePoolScope.AWSMachinePool.Spec.CapacityRebalance),
		LifecycleHookSpecificationList: getLifecycleHookSpecificationList(machinePoolScope.GetLifecycleHooks()),
	}

	if desiredCapacity != nil {
		input.DesiredCapacity = aws.Int32(*desiredCapacity)
	}

	if machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy != nil {
		input.MixedInstancesPolicy = createSDKMixedInstancesPolicy(name, machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy)
	} else {
		input.LaunchTemplate = &autoscalingtypes.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(machinePoolScope.AWSMachinePool.Status.LaunchTemplateID),
			Version:          aws.String(expinfrav1.LaunchTemplateLatestVersion),
		}
	}

	input.Tags = BuildTagsFromMap(name, infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.KubernetesClusterName(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	}))

	if _, err := s.ASGClient.CreateAutoScalingGroup(context.TODO(), input); err != nil {
		s.scope.Error(err, "unable to create AutoScalingGroup")
		return nil, errors.Wrap(err, "failed to create autoscaling group")
	}

	record.Eventf(machinePoolScope.AWSMachinePool, "SuccessfulCreate", "Created new ASG: %s", machinePoolScope.Name())

	return nil, nil
}

// DeleteASGAndWait will delete an ASG and wait until it is deleted.
func (s *Service) DeleteASGAndWait(name string) error {
	if err := s.DeleteASG(name); err != nil {
		return err
	}

	s.scope.Debug("Waiting for ASG to be deleted", "name", name)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}

	waiter := autoscaling.NewGroupNotExistsWaiter(s.ASGClient)

	if err := waiter.Wait(context.TODO(), input, 15*time.Minute); err != nil {
		return errors.Wrapf(err, "failed to wait for ASG %q deletion", name)
	}

	return nil
}

// DeleteASG will delete the ASG of a service.
func (s *Service) DeleteASG(name string) error {
	s.scope.Debug("Attempting to delete ASG", "name", name)

	input := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		ForceDelete:          aws.Bool(true),
	}

	if _, err := s.ASGClient.DeleteAutoScalingGroup(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to delete ASG %q", name)
	}

	s.scope.Debug("Deleted ASG", "name", name)
	return nil
}

// UpdateASG will update the ASG of a service.
func (s *Service) UpdateASG(machinePoolScope *scope.MachinePoolScope) error {
	subnetIDs, err := s.SubnetIDs(machinePoolScope)
	if err != nil {
		return fmt.Errorf("getting subnets for ASG: %w", err)
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(machinePoolScope.Name()), // TODO: define dynamically - borrow logic from ec2
		MaxSize:              aws.Int32(machinePoolScope.AWSMachinePool.Spec.MaxSize),
		MinSize:              aws.Int32(machinePoolScope.AWSMachinePool.Spec.MinSize),
		VPCZoneIdentifier:    aws.String(strings.Join(subnetIDs, ",")),
		CapacityRebalance:    aws.Bool(machinePoolScope.AWSMachinePool.Spec.CapacityRebalance),
	}

	if machinePoolScope.MachinePool.Spec.Replicas != nil && !annotations.ReplicasManagedByExternalAutoscaler(machinePoolScope.MachinePool) {
		input.DesiredCapacity = aws.Int32(*machinePoolScope.MachinePool.Spec.Replicas)
	}

	if machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy != nil {
		input.MixedInstancesPolicy = createSDKMixedInstancesPolicy(machinePoolScope.Name(), machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy)
	} else {
		input.LaunchTemplate = &autoscalingtypes.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(machinePoolScope.AWSMachinePool.Status.LaunchTemplateID),
			Version:          aws.String(expinfrav1.LaunchTemplateLatestVersion),
		}
	}

	if _, err := s.ASGClient.UpdateAutoScalingGroup(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to update ASG %q", machinePoolScope.Name())
	}

	return nil
}

// CanStartASGInstanceRefresh will start an ASG instance with refresh.
func (s *Service) CanStartASGInstanceRefresh(scope *scope.MachinePoolScope) (bool, error) {
	describeInput := &autoscaling.DescribeInstanceRefreshesInput{AutoScalingGroupName: aws.String(scope.Name())}
	refreshes, err := s.ASGClient.DescribeInstanceRefreshes(context.TODO(), describeInput)
	if err != nil {
		return false, err
	}
	hasUnfinishedRefresh := false
	if len(refreshes.InstanceRefreshes) != 0 {
		for i := range refreshes.InstanceRefreshes {
			if refreshes.InstanceRefreshes[i].Status == autoscalingtypes.InstanceRefreshStatusInProgress ||
				refreshes.InstanceRefreshes[i].Status == autoscalingtypes.InstanceRefreshStatusPending ||
				refreshes.InstanceRefreshes[i].Status == autoscalingtypes.InstanceRefreshStatusCancelling {
				hasUnfinishedRefresh = true
			}
		}
	}
	if hasUnfinishedRefresh {
		return false, nil
	}
	return true, nil
}

// StartASGInstanceRefresh will start an ASG instance with refresh.
func (s *Service) StartASGInstanceRefresh(scope *scope.MachinePoolScope) error {
	strategy := ptr.To(autoscalingtypes.RefreshStrategyRolling)
	var minHealthyPercentage, maxHealthyPercentage, instanceWarmup *int32
	if scope.AWSMachinePool.Spec.RefreshPreferences != nil {
		if scope.AWSMachinePool.Spec.RefreshPreferences.Strategy != nil {
			strategy = ptr.To(autoscalingtypes.RefreshStrategy(*scope.AWSMachinePool.Spec.RefreshPreferences.Strategy))
		}
		if scope.AWSMachinePool.Spec.RefreshPreferences.InstanceWarmup != nil {
			instanceWarmup = utils.ToInt32Pointer(scope.AWSMachinePool.Spec.RefreshPreferences.InstanceWarmup)
		}
		if scope.AWSMachinePool.Spec.RefreshPreferences.MinHealthyPercentage != nil {
			minHealthyPercentage = utils.ToInt32Pointer(scope.AWSMachinePool.Spec.RefreshPreferences.MinHealthyPercentage)
		}
		if scope.AWSMachinePool.Spec.RefreshPreferences.MaxHealthyPercentage != nil {
			maxHealthyPercentage = utils.ToInt32Pointer(scope.AWSMachinePool.Spec.RefreshPreferences.MaxHealthyPercentage)
		}
	}

	input := &autoscaling.StartInstanceRefreshInput{
		AutoScalingGroupName: aws.String(scope.Name()),
		Strategy:             *strategy,
		Preferences: &autoscalingtypes.RefreshPreferences{
			InstanceWarmup:       instanceWarmup,
			MinHealthyPercentage: minHealthyPercentage,
			MaxHealthyPercentage: maxHealthyPercentage,
		},
	}

	if _, err := s.ASGClient.StartInstanceRefresh(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to start ASG instance refresh %q", scope.Name())
	}

	return nil
}

func createSDKMixedInstancesPolicy(name string, i *expinfrav1.MixedInstancesPolicy) *autoscalingtypes.MixedInstancesPolicy {
	mixedInstancesPolicy := &autoscalingtypes.MixedInstancesPolicy{
		LaunchTemplate: &autoscalingtypes.LaunchTemplate{
			LaunchTemplateSpecification: &autoscalingtypes.LaunchTemplateSpecification{
				LaunchTemplateName: aws.String(name),
				Version:            aws.String(expinfrav1.LaunchTemplateLatestVersion),
			},
		},
	}

	if i.InstancesDistribution != nil {
		mixedInstancesPolicy.InstancesDistribution = &autoscalingtypes.InstancesDistribution{
			OnDemandAllocationStrategy:          aws.String(string(i.InstancesDistribution.OnDemandAllocationStrategy)),
			OnDemandBaseCapacity:                utils.ToInt32Pointer(i.InstancesDistribution.OnDemandBaseCapacity),
			OnDemandPercentageAboveBaseCapacity: utils.ToInt32Pointer(i.InstancesDistribution.OnDemandPercentageAboveBaseCapacity),
			SpotAllocationStrategy:              aws.String(string(i.InstancesDistribution.SpotAllocationStrategy)),
		}
	}

	for _, override := range i.Overrides {
		mixedInstancesPolicy.LaunchTemplate.Overrides = append(mixedInstancesPolicy.LaunchTemplate.Overrides, autoscalingtypes.LaunchTemplateOverrides{
			InstanceType: aws.String(override.InstanceType),
		})
	}

	return mixedInstancesPolicy
}

// BuildTagsFromMap takes a map of keys and values and returns them as autoscaling group tags.
func BuildTagsFromMap(asgName string, inTags map[string]string) []autoscalingtypes.Tag {
	if inTags == nil {
		return nil
	}
	tags := make([]autoscalingtypes.Tag, 0)
	for k, v := range inTags {
		tags = append(tags, autoscalingtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
			// We set the instance tags in the LaunchTemplate, disabling propagation to prevent the two
			// resources from clobbering the tags set in the LaunchTemplate
			PropagateAtLaunch: aws.Bool(false),
			ResourceId:        aws.String(asgName),
			ResourceType:      aws.String("auto-scaling-group"),
		})
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// UpdateResourceTags updates the tags for an autoscaling group.
// This will be called if there is anything to create (update) or delete.
// We may not always have to perform each action, so we check what we're
// receiving to avoid calling AWS if we don't need to.
func (s *Service) UpdateResourceTags(resourceID *string, create, remove map[string]string) error {
	s.scope.Debug("Attempting to update tags on resource", "resource-id", *resourceID)
	s.scope.Info("updating tags on resource", "resource-id", *resourceID, "create", create, "remove", remove)

	// If we have anything to create or update
	if len(create) > 0 {
		s.scope.Debug("Attempting to create tags on resource", "resource-id", *resourceID)

		createOrUpdateTagsInput := &autoscaling.CreateOrUpdateTagsInput{}

		createOrUpdateTagsInput.Tags = mapToTags(create, resourceID)

		if _, err := s.ASGClient.CreateOrUpdateTags(context.TODO(), createOrUpdateTagsInput); err != nil {
			return errors.Wrapf(err, "failed to update tags on AutoScalingGroup %q", *resourceID)
		}
	}

	// If we have anything to remove
	if len(remove) > 0 {
		s.scope.Debug("Attempting to delete tags on resource", "resource-id", *resourceID)

		// Convert our remove map into an array of *ec2.Tag
		removeTagsInput := mapToTags(remove, resourceID)

		// Create the DeleteTags input
		input := &autoscaling.DeleteTagsInput{
			Tags: removeTagsInput,
		}

		// Delete tags in AWS.
		if _, err := s.ASGClient.DeleteTags(context.TODO(), input); err != nil {
			return errors.Wrapf(err, "failed to delete tags on AutoScalingGroup %q: %v", *resourceID, remove)
		}
	}

	return nil
}

// SuspendProcesses suspends the processes for an autoscaling group.
func (s *Service) SuspendProcesses(name string, processes []string) error {
	input := autoscaling.SuspendProcessesInput{
		AutoScalingGroupName: aws.String(name),
		ScalingProcesses:     processes,
	}
	if _, err := s.ASGClient.SuspendProcesses(context.TODO(), &input); err != nil {
		return errors.Wrapf(err, "failed to suspend processes for AutoScalingGroup: %q", name)
	}
	return nil
}

// ResumeProcesses resumes the processes for an autoscaling group.
func (s *Service) ResumeProcesses(name string, processes []string) error {
	input := autoscaling.ResumeProcessesInput{
		AutoScalingGroupName: aws.String(name),
		ScalingProcesses:     processes,
	}
	if _, err := s.ASGClient.ResumeProcesses(context.TODO(), &input); err != nil {
		return errors.Wrapf(err, "failed to resume processes for AutoScalingGroup: %q", name)
	}
	return nil
}

func mapToTags(input map[string]string, resourceID *string) []autoscalingtypes.Tag {
	tags := make([]autoscalingtypes.Tag, 0)
	for k, v := range input {
		tags = append(tags, autoscalingtypes.Tag{
			Key:               aws.String(k),
			PropagateAtLaunch: aws.Bool(false),
			ResourceId:        resourceID,
			ResourceType:      aws.String("auto-scaling-group"),
			Value:             aws.String(v),
		})
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// SubnetIDs return subnet IDs of a AWSMachinePool based on given subnetIDs and filters.
func (s *Service) SubnetIDs(scope *scope.MachinePoolScope) ([]string, error) {
	subnetIDs := make([]string, 0)
	var inputFilters = make([]*ec2.Filter, 0)

	for _, subnet := range scope.AWSMachinePool.Spec.Subnets {
		switch {
		case subnet.ID != nil:
			subnetIDs = append(subnetIDs, aws.StringValue(subnet.ID))
		case subnet.Filters != nil:
			for _, eachFilter := range subnet.Filters {
				inputFilters = append(inputFilters, &ec2.Filter{
					Name:   aws.String(eachFilter.Name),
					Values: aws.StringSlice(eachFilter.Values),
				})
			}
		}
	}

	if len(inputFilters) > 0 {
		out, err := s.EC2Client.DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
			Filters: inputFilters,
		})
		if err != nil {
			return nil, err
		}

		for _, subnet := range out.Subnets {
			tags := converters.TagsToMap(subnet.Tags)
			if tags[infrav1.NameAWSSubnetAssociation] == infrav1.SecondarySubnetTagValue {
				// Subnet belongs to a secondary CIDR block which won't be used to create instances
				continue
			}

			subnetIDs = append(subnetIDs, *subnet.SubnetId)
		}

		if len(subnetIDs) == 0 {
			errMessage := fmt.Sprintf("failed to create ASG %q, no subnets available matching criteria %q", scope.Name(), inputFilters)
			record.Warnf(scope.AWSMachinePool, "FailedCreate", errMessage)
			return subnetIDs, awserrors.NewFailedDependency(errMessage)
		}
	}

	return scope.SubnetIDs(subnetIDs)
}
