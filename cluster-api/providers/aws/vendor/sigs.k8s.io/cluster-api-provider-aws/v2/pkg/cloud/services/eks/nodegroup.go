/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/version"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func (s *NodegroupService) describeNodegroup() (*eks.Nodegroup, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	s.scope.Debug("describing eks node group", "cluster", eksClusterName, "nodegroup", nodegroupName)
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
	}

	out, err := s.EKSClient.DescribeNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				return nil, nil
			default:
				return nil, errors.Wrap(err, "failed to describe nodegroup")
			}
		} else {
			return nil, errors.Wrap(err, "failed to describe nodegroup")
		}
	}

	return out.Nodegroup, nil
}

func (s *NodegroupService) describeASGs(ng *eks.Nodegroup) (*autoscaling.Group, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	s.scope.Debug("describing node group ASG", "cluster", eksClusterName, "nodegroup", nodegroupName)

	if len(ng.Resources.AutoScalingGroups) == 0 {
		return nil, nil
	}

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			ng.Resources.AutoScalingGroups[0].Name,
		},
	}

	out, err := s.AutoscalingClient.DescribeAutoScalingGroups(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "failed to describe ASGs")
	case len(out.AutoScalingGroups) == 0:
		return nil, errors.Wrap(err, "no ASG found")
	}

	return out.AutoScalingGroups[0], nil
}

func (s *NodegroupService) scalingConfig() *eks.NodegroupScalingConfig {
	var replicas int32 = 1
	if s.scope.MachinePool.Spec.Replicas != nil {
		replicas = *s.scope.MachinePool.Spec.Replicas
	}
	cfg := eks.NodegroupScalingConfig{
		DesiredSize: aws.Int64(int64(replicas)),
	}
	scaling := s.scope.ManagedMachinePool.Spec.Scaling
	if scaling == nil {
		return &cfg
	}
	if scaling.MaxSize != nil {
		cfg.MaxSize = aws.Int64(int64(*scaling.MaxSize))
	}
	if scaling.MaxSize != nil {
		cfg.MinSize = aws.Int64(int64(*scaling.MinSize))
	}
	return &cfg
}

func (s *NodegroupService) updateConfig() *eks.NodegroupUpdateConfig {
	updateConfig := s.scope.ManagedMachinePool.Spec.UpdateConfig

	return converters.NodegroupUpdateconfigToSDK(updateConfig)
}

func (s *NodegroupService) roleArn() (*string, error) {
	var role *iam.Role
	if s.scope.RoleName() != "" {
		var err error
		role, err = s.GetIAMRole(s.scope.RoleName())
		if err != nil {
			return nil, errors.Wrapf(err, "error getting node group IAM role: %s", s.scope.RoleName())
		}
	}
	return role.Arn, nil
}

func ngTags(key string, additionalTags infrav1.Tags) map[string]string {
	tags := additionalTags.DeepCopy()
	tags[infrav1.ClusterAWSCloudProviderTagKey(key)] = string(infrav1.ResourceLifecycleOwned)
	return tags
}

func (s *NodegroupService) remoteAccess() (*eks.RemoteAccessConfig, error) {
	pool := s.scope.ManagedMachinePool.Spec
	if pool.RemoteAccess == nil {
		return nil, nil
	}

	controlPlane := s.scope.ControlPlane

	// SourceSecurityGroups is validated to be empty if PublicAccess is true
	// but just in case we use an empty list to take advantage of the documented
	// API behavior
	var sSGs = []string{}

	if !pool.RemoteAccess.Public {
		sSGs = pool.RemoteAccess.SourceSecurityGroups
		// We add the EKS created cluster security group to the allowed security
		// groups by default to prevent the API default of 0.0.0.0/0 from taking effect
		// in case SourceSecurityGroups is empty
		clusterSG, ok := controlPlane.Status.Network.SecurityGroups[ekscontrolplanev1.SecurityGroupCluster]
		if !ok {
			return nil, errors.Errorf("%s security group not found on control plane", ekscontrolplanev1.SecurityGroupCluster)
		}
		sSGs = append(sSGs, clusterSG.ID)

		if controlPlane.Spec.Bastion.Enabled {
			bastionSG, ok := controlPlane.Status.Network.SecurityGroups[infrav1.SecurityGroupBastion]
			if !ok {
				return nil, errors.Errorf("%s security group not found on control plane", infrav1.SecurityGroupBastion)
			}
			sSGs = append(
				sSGs,
				bastionSG.ID,
			)
		}
	}

	sshKeyName := pool.RemoteAccess.SSHKeyName
	if sshKeyName == nil {
		sshKeyName = controlPlane.Spec.SSHKeyName
	}

	return &eks.RemoteAccessConfig{
		SourceSecurityGroups: aws.StringSlice(sSGs),
		Ec2SshKey:            sshKeyName,
	}, nil
}

func (s *NodegroupService) createNodegroup() (*eks.Nodegroup, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	additionalTags := s.scope.AdditionalTags()
	roleArn, err := s.roleArn()
	if err != nil {
		return nil, err
	}
	managedPool := s.scope.ManagedMachinePool.Spec
	tags := ngTags(s.scope.ClusterName(), additionalTags)

	remoteAccess, err := s.remoteAccess()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create remote access configuration")
	}

	subnets, err := s.scope.SubnetIDs()
	if err != nil {
		return nil, fmt.Errorf("failed getting nodegroup subnets: %w", err)
	}

	input := &eks.CreateNodegroupInput{
		ScalingConfig: s.scalingConfig(),
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
		Subnets:       aws.StringSlice(subnets),
		NodeRole:      roleArn,
		Labels:        aws.StringMap(managedPool.Labels),
		Tags:          aws.StringMap(tags),
		RemoteAccess:  remoteAccess,
		UpdateConfig:  s.updateConfig(),
	}
	if managedPool.AMIType != nil && (managedPool.AWSLaunchTemplate == nil || managedPool.AWSLaunchTemplate.AMI.ID == nil) {
		input.AmiType = aws.String(string(*managedPool.AMIType))
	}
	if managedPool.DiskSize != nil {
		input.DiskSize = aws.Int64(int64(*managedPool.DiskSize))
	}
	if managedPool.InstanceType != nil {
		input.InstanceTypes = []*string{managedPool.InstanceType}
	}
	if len(managedPool.Taints) > 0 {
		s.Info("adding taints to nodegroup", "nodegroup", nodegroupName)
		taints, err := converters.TaintsToSDK(managedPool.Taints)
		if err != nil {
			return nil, fmt.Errorf("converting taints: %w", err)
		}
		input.Taints = taints
	}
	if managedPool.CapacityType != nil {
		capacityType, err := converters.CapacityTypeToSDK(*managedPool.CapacityType)
		if err != nil {
			return nil, fmt.Errorf("converting capacity type: %w", err)
		}
		input.CapacityType = aws.String(capacityType)
	}
	if managedPool.AWSLaunchTemplate != nil {
		input.LaunchTemplate = &eks.LaunchTemplateSpecification{
			Id:      s.scope.ManagedMachinePool.Status.LaunchTemplateID,
			Version: s.scope.ManagedMachinePool.Status.LaunchTemplateVersion,
		}
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "created invalid CreateNodegroupInput")
	}

	out, err := s.EKSClient.CreateNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// TODO
			case eks.ErrCodeResourceNotFoundException:
				return nil, nil
			default:
				return nil, errors.Wrap(err, "failed to create nodegroup")
			}
		} else {
			return nil, errors.Wrap(err, "failed to create nodegroup")
		}
	}

	return out.Nodegroup, nil
}

func (s *NodegroupService) deleteNodegroupAndWait() (reterr error) {
	eksClusterName := s.scope.KubernetesClusterName()
	nodegroupName := s.scope.NodegroupName()
	if err := s.scope.NodegroupReadyFalse(clusterv1.DeletingReason, ""); err != nil {
		return err
	}
	defer func() {
		if reterr != nil {
			record.Warnf(
				s.scope.ManagedMachinePool, "FailedDeleteEKSNodegroup", "Failed to delete EKS nodegroup %s: %v", s.scope.NodegroupName(), reterr,
			)
			if err := s.scope.NodegroupReadyFalse("DeletingFailed", reterr.Error()); err != nil {
				reterr = err
			}
		} else if err := s.scope.NodegroupReadyFalse(clusterv1.DeletedReason, ""); err != nil {
			reterr = err
		}
	}()
	input := &eks.DeleteNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
	}
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "created invalid DeleteNodegroupInput")
	}

	_, err := s.EKSClient.DeleteNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// TODO
			case eks.ErrCodeResourceNotFoundException:
				return nil
			default:
				return errors.Wrap(err, "failed to delete nodegroup")
			}
		} else {
			return errors.Wrap(err, "failed to delete nodegroup")
		}
	}

	waitInput := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodegroupName),
	}
	err = s.EKSClient.WaitUntilNodegroupDeleted(waitInput)
	if err != nil {
		return errors.Wrapf(err, "failed waiting for EKS nodegroup %s to delete", nodegroupName)
	}

	return nil
}

func (s *NodegroupService) reconcileNodegroupVersion(ng *eks.Nodegroup) error {
	var specVersion *version.Version
	if s.scope.Version() != nil {
		specVersion = parseEKSVersion(*s.scope.Version())
	}
	ngVersion := version.MustParseGeneric(*ng.Version)
	specAMI := s.scope.ManagedMachinePool.Spec.AMIVersion
	ngAMI := *ng.ReleaseVersion
	statusLaunchTemplateVersion := s.scope.ManagedMachinePool.Status.LaunchTemplateVersion
	var ngLaunchTemplateVersion *string
	if ng.LaunchTemplate != nil {
		ngLaunchTemplateVersion = ng.LaunchTemplate.Version
	}

	eksClusterName := s.scope.KubernetesClusterName()
	if (specVersion != nil && ngVersion.LessThan(specVersion)) || (specAMI != nil && *specAMI != ngAMI) || (statusLaunchTemplateVersion != nil && *statusLaunchTemplateVersion != *ngLaunchTemplateVersion) {
		input := &eks.UpdateNodegroupVersionInput{
			ClusterName:   aws.String(eksClusterName),
			NodegroupName: aws.String(s.scope.NodegroupName()),
		}

		var updateMsg string
		// Either update k8s version or AMI version
		switch {
		case specVersion != nil && ngVersion.LessThan(specVersion):
			// NOTE: you can only upgrade increments of minor versions. If you want to upgrade 1.14 to 1.16 we
			// need to go 1.14-> 1.15 and then 1.15 -> 1.16.
			input.Version = aws.String(versionToEKS(ngVersion.WithMinor(ngVersion.Minor() + 1)))
			updateMsg = fmt.Sprintf("to version %s", *input.Version)
		case specAMI != nil && *specAMI != ngAMI:
			input.ReleaseVersion = specAMI
			updateMsg = fmt.Sprintf("to AMI version %s", *input.ReleaseVersion)
		case statusLaunchTemplateVersion != nil && *statusLaunchTemplateVersion != *ngLaunchTemplateVersion:
			input.LaunchTemplate = &eks.LaunchTemplateSpecification{
				Id:      s.scope.ManagedMachinePool.Status.LaunchTemplateID,
				Version: statusLaunchTemplateVersion,
			}
			updateMsg = fmt.Sprintf("to launch template version %s", *statusLaunchTemplateVersion)
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateNodegroupVersion(input); err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					return false, aerr
				}
				return false, err
			}
			record.Eventf(s.scope.ManagedMachinePool, "SuccessfulUpdateEKSNodegroup", "Updated EKS nodegroup %s %s", eksClusterName, updateMsg)
			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ManagedMachinePool, "FailedUpdateEKSNodegroup", "failed to update the EKS nodegroup %s %s: %v", eksClusterName, updateMsg, err)
			return errors.Wrapf(err, "failed to update EKS nodegroup")
		}
	}
	return nil
}

func createLabelUpdate(specLabels map[string]string, ng *eks.Nodegroup) *eks.UpdateLabelsPayload {
	current := ng.Labels
	payload := eks.UpdateLabelsPayload{
		AddOrUpdateLabels: map[string]*string{},
	}
	for k, v := range specLabels {
		if currentV, ok := current[k]; !ok || currentV == nil || v != *currentV {
			payload.AddOrUpdateLabels[k] = aws.String(v)
		}
	}
	for k := range current {
		if _, ok := specLabels[k]; !ok {
			payload.RemoveLabels = append(payload.RemoveLabels, aws.String(k))
		}
	}
	if len(payload.AddOrUpdateLabels) > 0 || len(payload.RemoveLabels) > 0 {
		return &payload
	}
	return nil
}

func (s *NodegroupService) createTaintsUpdate(specTaints expinfrav1.Taints, ng *eks.Nodegroup) (*eks.UpdateTaintsPayload, error) {
	s.Debug("Creating taints update for node group", "name", *ng.NodegroupName, "num_current", len(ng.Taints), "num_required", len(specTaints))
	current, err := converters.TaintsFromSDK(ng.Taints)
	if err != nil {
		return nil, fmt.Errorf("converting taints: %w", err)
	}
	payload := eks.UpdateTaintsPayload{}
	for _, specTaint := range specTaints {
		st := specTaint.DeepCopy()
		if !current.Contains(st) {
			sdkTaint, err := converters.TaintToSDK(*st)
			if err != nil {
				return nil, fmt.Errorf("converting taint to sdk: %w", err)
			}
			payload.AddOrUpdateTaints = append(payload.AddOrUpdateTaints, sdkTaint)
		}
	}
	for _, currentTaint := range current {
		ct := currentTaint.DeepCopy()
		if !specTaints.Contains(ct) {
			sdkTaint, err := converters.TaintToSDK(*ct)
			if err != nil {
				return nil, fmt.Errorf("converting taint to sdk: %w", err)
			}
			payload.RemoveTaints = append(payload.RemoveTaints, sdkTaint)
		}
	}
	if len(payload.AddOrUpdateTaints) > 0 || len(payload.RemoveTaints) > 0 {
		s.Debug("Node group taints update required", "name", *ng.NodegroupName, "addupdate", len(payload.AddOrUpdateTaints), "remove", len(payload.RemoveTaints))
		return &payload, nil
	}

	s.Debug("No updates required for node group taints", "name", *ng.NodegroupName)
	return nil, nil
}

func (s *NodegroupService) reconcileNodegroupConfig(ng *eks.Nodegroup) error {
	eksClusterName := s.scope.KubernetesClusterName()
	s.Debug("reconciling node group config", "cluster", eksClusterName, "name", *ng.NodegroupName)

	managedPool := s.scope.ManagedMachinePool.Spec
	input := &eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(managedPool.EKSNodegroupName),
	}
	var needsUpdate bool
	if labelPayload := createLabelUpdate(managedPool.Labels, ng); labelPayload != nil {
		s.Debug("Nodegroup labels need an update", "nodegroup", ng.NodegroupName)
		input.Labels = labelPayload
		needsUpdate = true
	}
	taintsPayload, err := s.createTaintsUpdate(managedPool.Taints, ng)
	if err != nil {
		return fmt.Errorf("creating taints update payload: %w", err)
	}
	if taintsPayload != nil {
		s.Debug("nodegroup taints need updating")
		input.Taints = taintsPayload
		needsUpdate = true
	}
	if machinePool := s.scope.MachinePool.Spec; machinePool.Replicas == nil {
		if ng.ScalingConfig.DesiredSize != nil && *ng.ScalingConfig.DesiredSize != 1 {
			s.Debug("Nodegroup desired size differs from spec, updating scaling configuration", "nodegroup", ng.NodegroupName)
			input.ScalingConfig = s.scalingConfig()
			needsUpdate = true
		}
	} else if ng.ScalingConfig.DesiredSize == nil || int64(*machinePool.Replicas) != *ng.ScalingConfig.DesiredSize {
		s.Debug("Nodegroup has no desired size or differs from replicas, updating scaling configuration", "nodegroup", ng.NodegroupName)
		input.ScalingConfig = s.scalingConfig()
		needsUpdate = true
	}
	if managedPool.Scaling != nil && ((aws.Int64Value(ng.ScalingConfig.MaxSize) != int64(aws.Int32Value(managedPool.Scaling.MaxSize))) ||
		(aws.Int64Value(ng.ScalingConfig.MinSize) != int64(aws.Int32Value(managedPool.Scaling.MinSize)))) {
		s.Debug("Nodegroup min/max differ from spec, updating scaling configuration", "nodegroup", ng.NodegroupName)
		input.ScalingConfig = s.scalingConfig()
		needsUpdate = true
	}
	currentUpdateConfig := converters.NodegroupUpdateconfigFromSDK(ng.UpdateConfig)
	if !cmp.Equal(managedPool.UpdateConfig, currentUpdateConfig) {
		s.Debug("Nodegroup update configuration differs from spec, updating the nodegroup update config", "nodegroup", ng.NodegroupName)
		input.UpdateConfig = s.updateConfig()
		needsUpdate = true
	}
	if !needsUpdate {
		s.Debug("node group config update not needed", "cluster", eksClusterName, "name", *ng.NodegroupName)
		return nil
	}
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "created invalid UpdateNodegroupConfigInput")
	}

	_, err = s.EKSClient.UpdateNodegroupConfig(input)
	if err != nil {
		return errors.Wrap(err, "failed to update nodegroup config")
	}

	return nil
}

func (s *NodegroupService) reconcileNodegroup(ctx context.Context) error {
	ng, err := s.describeNodegroup()
	if err != nil {
		return errors.Wrap(err, "failed to describe nodegroup")
	}

	if eksClusterName, eksNodegroupName := s.scope.KubernetesClusterName(), s.scope.NodegroupName(); ng == nil {
		ng, err = s.createNodegroup()
		if err != nil {
			return errors.Wrap(err, "failed to create nodegroup")
		}
		s.scope.Info("Created EKS nodegroup in AWS", "cluster-name", eksClusterName, "nodegroup-name", eksNodegroupName)
	} else {
		tagKey := infrav1.ClusterAWSCloudProviderTagKey(s.scope.ClusterName())
		ownedTag := ng.Tags[tagKey]
		if ownedTag == nil {
			return errors.Errorf("owner of %s mismatch: %s", eksNodegroupName, s.scope.ClusterName())
		}
		s.scope.Debug("Found owned EKS nodegroup in AWS", "cluster-name", eksClusterName, "nodegroup-name", eksNodegroupName)
	}

	if err := s.setStatus(ng); err != nil {
		return errors.Wrap(err, "failed to set status")
	}

	switch *ng.Status {
	case eks.NodegroupStatusCreating, eks.NodegroupStatusUpdating:
		ng, err = s.waitForNodegroupActive()
	default:
		break
	}

	if scope.ReplicasExternallyManaged(s.scope.MachinePool) {
		// Set MachinePool replicas to the node group DesiredCapacity
		ngDesiredCapacity := int32(aws.Int64Value(ng.ScalingConfig.DesiredSize))
		if *s.scope.MachinePool.Spec.Replicas != ngDesiredCapacity {
			s.scope.Info("Setting MachinePool replicas to node group DesiredCapacity",
				"local", *s.scope.MachinePool.Spec.Replicas,
				"external", ngDesiredCapacity)
			s.scope.MachinePool.Spec.Replicas = &ngDesiredCapacity
			if err := s.scope.PatchCAPIMachinePoolObject(ctx); err != nil {
				return err
			}
		}
	}

	if err != nil {
		return errors.Wrap(err, "failed to wait for nodegroup to be active")
	}

	if err := s.reconcileNodegroupVersion(ng); err != nil {
		return errors.Wrap(err, "failed to reconcile nodegroup version")
	}

	if err := s.reconcileNodegroupConfig(ng); err != nil {
		return errors.Wrap(err, "failed to reconcile nodegroup config")
	}

	if err := s.reconcileTags(ng); err != nil {
		return errors.Wrapf(err, "failed to reconcile nodegroup tags")
	}

	if err := s.reconcileASGTags(ng); err != nil {
		return errors.Wrapf(err, "failed to reconcile asg tags")
	}

	return nil
}

func (s *NodegroupService) setStatus(ng *eks.Nodegroup) error {
	managedPool := s.scope.ManagedMachinePool
	switch *ng.Status {
	case eks.NodegroupStatusDeleting:
		managedPool.Status.Ready = false
	case eks.NodegroupStatusCreateFailed, eks.NodegroupStatusDeleteFailed:
		managedPool.Status.Ready = false
		// TODO FailureReason
		failureMsg := fmt.Sprintf("EKS nodegroup in failed %s status", *ng.Status)
		managedPool.Status.FailureMessage = &failureMsg
	case eks.NodegroupStatusActive:
		managedPool.Status.Ready = true
		managedPool.Status.FailureMessage = nil
		// TODO FailureReason
	case eks.NodegroupStatusCreating:
		managedPool.Status.Ready = false
	case eks.NodegroupStatusUpdating:
		managedPool.Status.Ready = true
	default:
		return errors.Errorf("unexpected EKS nodegroup status %s", *ng.Status)
	}
	if managedPool.Status.Ready && ng.Resources != nil && len(ng.Resources.AutoScalingGroups) > 0 {
		req := autoscaling.DescribeAutoScalingGroupsInput{}
		for _, asg := range ng.Resources.AutoScalingGroups {
			req.AutoScalingGroupNames = append(req.AutoScalingGroupNames, asg.Name)
		}
		groups, err := s.AutoscalingClient.DescribeAutoScalingGroups(&req)
		if err != nil {
			return errors.Wrap(err, "failed to describe AutoScalingGroup for nodegroup")
		}

		var replicas int32
		var providerIDList []string
		for _, group := range groups.AutoScalingGroups {
			replicas += int32(len(group.Instances))
			for _, instance := range group.Instances {
				providerIDList = append(providerIDList, fmt.Sprintf("aws:///%s/%s", *instance.AvailabilityZone, *instance.InstanceId))
			}
		}
		managedPool.Spec.ProviderIDList = providerIDList
		managedPool.Status.Replicas = replicas
	}
	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to update nodegroup")
	}
	return nil
}

func (s *NodegroupService) waitForNodegroupActive() (*eks.Nodegroup, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	eksNodegroupName := s.scope.NodegroupName()
	req := eks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(eksNodegroupName),
	}
	if err := s.EKSClient.WaitUntilNodegroupActive(&req); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for EKS nodegroup %q", *req.NodegroupName)
	}

	s.scope.Info("EKS nodegroup is now available", "nodegroup-name", eksNodegroupName)

	ng, err := s.describeNodegroup()
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe EKS nodegroup")
	}
	if err := s.setStatus(ng); err != nil {
		return nil, errors.Wrap(err, "failed to set status")
	}

	return ng, nil
}
