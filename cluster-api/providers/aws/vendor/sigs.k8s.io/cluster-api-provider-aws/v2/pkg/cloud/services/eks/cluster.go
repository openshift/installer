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
	"net"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/blang/semver"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/internal/cidr"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/internal/cmp"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/internal/tristate"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func (s *Service) reconcileCluster(ctx context.Context) error {
	s.scope.Debug("Reconciling EKS cluster")

	eksClusterName := s.scope.KubernetesClusterName()

	cluster, err := s.describeEKSCluster(ctx, eksClusterName)
	if err != nil {
		return errors.Wrap(err, "failed to describe eks clusters")
	}

	if cluster == nil {
		cluster, err = s.createCluster(ctx, eksClusterName)
		if err != nil {
			return errors.Wrap(err, "failed to create cluster")
		}
	} else {
		tagKey := infrav1.ClusterAWSCloudProviderTagKey(eksClusterName)
		ownedTag := aws.String(cluster.Tags[tagKey])
		// Prior to https://github.com/kubernetes-sigs/cluster-api-provider-aws/pull/3573,
		// Clusters were tagged using s.scope.Name()
		// To support upgrading older clusters, check for both tags
		oldTagKey := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())
		oldOwnedTag := aws.String(cluster.Tags[oldTagKey])

		if ownedTag == nil && oldOwnedTag == nil {
			return fmt.Errorf("EKS cluster resource %q must have a tag with key %q or %q", eksClusterName, oldTagKey, tagKey)
		}

		s.scope.Debug("Found owned EKS cluster in AWS", "cluster", klog.KRef("", eksClusterName))
	}

	if err := s.setStatus(cluster); err != nil {
		return errors.Wrap(err, "failed to set status")
	}

	// Wait for our cluster to be ready if necessary
	switch cluster.Status {
	case ekstypes.ClusterStatusUpdating, ekstypes.ClusterStatusCreating:
		cluster, err = s.waitForClusterActive(ctx)
	default:
		break
	}
	if err != nil {
		return errors.Wrap(err, "failed to wait for cluster to be active")
	}

	if !s.scope.ControlPlane.Status.Ready {
		return nil
	}

	s.scope.Debug("EKS Control Plane active", "endpoint", *cluster.Endpoint)

	s.scope.ControlPlane.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
		Host: *cluster.Endpoint,
		Port: 443,
	}

	if err := s.reconcileSecurityGroups(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling security groups")
	}

	if err := s.reconcileKubeconfig(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling kubeconfig")
	}

	if err := s.reconcileAdditionalKubeconfigs(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling additional kubeconfigs")
	}

	if err := s.reconcileClusterVersion(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling cluster version")
	}

	if err := s.reconcileClusterConfig(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling cluster config")
	}

	if err := s.reconcileLogging(ctx, cluster.Logging); err != nil {
		return errors.Wrap(err, "failed reconciling logging")
	}

	if err := s.reconcileEKSEncryptionConfig(ctx, cluster.EncryptionConfig); err != nil {
		return errors.Wrap(err, "failed reconciling eks encryption config")
	}

	if err := s.reconcileTags(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed updating cluster tags")
	}

	if err := s.reconcileOIDCProvider(ctx, cluster); err != nil {
		return errors.Wrap(err, "failed reconciling OIDC provider for cluster")
	}

	return nil
}

// computeCurrentStatusVersion returns the computed current EKS cluster kubernetes version.
// The computation has awareness of the fact that EKS clusters only return a major.minor kubernetes version,
// and returns a compatible version for te status according to the one the user specified in the spec.
func computeCurrentStatusVersion(specV *string, clusterV *string) *string {
	specVersion := ""
	if specV != nil {
		specVersion = *specV
	}

	clusterVersion := ""
	if clusterV != nil {
		clusterVersion = *clusterV
	}

	// Ignore parsing errors as these are already validated by the kubebuilder validation and the AWS API.
	// Also specVersion might not be specified in the spec.Version for AWSManagedControlPlane, this results in a "0.0.0" version.
	// Also clusterVersion might not yet be returned by the AWS EKS API, as the cluster might still be initializing, this results in a "0.0.0" version.
	specSemverVersion, _ := semver.ParseTolerant(specVersion)
	currentSemverVersion, _ := semver.ParseTolerant(clusterVersion)

	// If AWS EKS API is not returning a version, set the status.Version to empty string.
	if currentSemverVersion.String() == "0.0.0" {
		return ptr.To("")
	}

	if currentSemverVersion.Major == specSemverVersion.Major &&
		currentSemverVersion.Minor == specSemverVersion.Minor &&
		specSemverVersion.Patch != 0 {
		// Treat this case differently as we want it to exactly match the spec.Version,
		// including its Patch, in the status.Version.
		currentSemverVersion.Patch = specSemverVersion.Patch

		return ptr.To(currentSemverVersion.String())
	}

	// For all the other cases it doesn't matter to have a patch version, as EKS ignores it internally.
	// So set the current cluster.Version (this always is a major.minor version format (e.g. "1.31")) in the status.Version.
	// Even in the event where in the spec.Version a zero patch version is specified (e.g. "1.31.0"),
	// the call to semver.ParseTolerant on the consumer side
	// will make sure the version with and without the trailing zero actually result in a match.
	return clusterV
}

func (s *Service) setStatus(cluster *ekstypes.Cluster) error {
	// Set the current Kubernetes control plane version in the status.
	s.scope.ControlPlane.Status.Version = computeCurrentStatusVersion(s.scope.ControlPlane.Spec.Version, cluster.Version)

	// Set the current cluster status in the control plane status.
	switch cluster.Status {
	case ekstypes.ClusterStatusDeleting:
		s.scope.ControlPlane.Status.Ready = false
	case ekstypes.ClusterStatusFailed:
		s.scope.ControlPlane.Status.Ready = false
		// TODO FailureReason
		failureMsg := fmt.Sprintf("EKS cluster in unexpected %s state", cluster.Status)
		s.scope.ControlPlane.Status.FailureMessage = &failureMsg
	case ekstypes.ClusterStatusActive:
		s.scope.ControlPlane.Status.Ready = true
		s.scope.ControlPlane.Status.FailureMessage = nil
		if conditions.IsTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneCreatingCondition) {
			record.Eventf(s.scope.ControlPlane, "SuccessfulCreateEKSControlPlane", "Created new EKS control plane %s", s.scope.KubernetesClusterName())
			conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneCreatingCondition, "created", clusterv1.ConditionSeverityInfo, "")
		}
		if conditions.IsTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneUpdatingCondition) {
			conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneUpdatingCondition, "updated", clusterv1.ConditionSeverityInfo, "")
			record.Eventf(s.scope.ControlPlane, "SuccessfulUpdateEKSControlPlane", "Updated EKS control plane %s", s.scope.KubernetesClusterName())
		}
		// TODO FailureReason
	case ekstypes.ClusterStatusCreating:
		s.scope.ControlPlane.Status.Ready = false
	case ekstypes.ClusterStatusUpdating:
		s.scope.ControlPlane.Status.Ready = true
	default:
		return errors.Errorf("unexpected EKS cluster status %s", cluster.Status)
	}
	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to update control plane")
	}
	return nil
}

// deleteCluster deletes an EKS cluster.
func (s *Service) deleteCluster(ctx context.Context) error {
	eksClusterName := s.scope.KubernetesClusterName()

	if eksClusterName == "" {
		s.scope.Debug("no EKS cluster name, skipping EKS cluster deletion")
		return nil
	}

	cluster, err := s.describeEKSCluster(ctx, eksClusterName)
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.Trace("eks cluster does not exist")
			return nil
		}
		return errors.Wrap(err, "unable to describe eks cluster")
	}
	if cluster == nil {
		return nil
	}

	err = s.deleteClusterAndWait(ctx, cluster)
	if err != nil {
		record.Warnf(s.scope.ControlPlane, "FailedDeleteEKSCluster", "Failed to delete EKS cluster %s: %v", s.scope.KubernetesClusterName(), err)
		return errors.Wrap(err, "unable to delete EKS cluster")
	}
	record.Eventf(s.scope.ControlPlane, "SuccessfulDeleteEKSCluster", "Deleted EKS Cluster %s", s.scope.KubernetesClusterName())

	return nil
}

func (s *Service) deleteClusterAndWait(ctx context.Context, cluster *ekstypes.Cluster) error {
	s.scope.Info("Deleting EKS cluster", "cluster", klog.KRef("", s.scope.KubernetesClusterName()))

	input := &eks.DeleteClusterInput{
		Name: cluster.Name,
	}
	_, err := s.EKSClient.DeleteCluster(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "failed to request delete of eks cluster %s", *cluster.Name)
	}

	waitInput := &eks.DescribeClusterInput{
		Name: cluster.Name,
	}

	err = s.EKSClient.WaitUntilClusterDeleted(ctx, waitInput, s.scope.MaxWaitActiveUpdateDelete)
	if err != nil {
		return errors.Wrapf(err, "failed waiting for eks cluster %s to delete", *cluster.Name)
	}

	return nil
}

func makeEksEncryptionConfigs(encryptionConfig *ekscontrolplanev1.EncryptionConfig) []ekstypes.EncryptionConfig {
	cfg := []ekstypes.EncryptionConfig{}

	if encryptionConfig == nil {
		return cfg
	}
	// TODO: change EncryptionConfig so that provider and resources are required  if encruptionConfig is specified
	if encryptionConfig.Provider == nil || len(*encryptionConfig.Provider) == 0 {
		return cfg
	}
	if len(encryptionConfig.Resources) == 0 {
		return cfg
	}

	return append(cfg, ekstypes.EncryptionConfig{
		Provider: &ekstypes.Provider{
			KeyArn: encryptionConfig.Provider,
		},
		Resources: aws.ToStringSlice(encryptionConfig.Resources),
	})
}

func makeKubernetesNetworkConfig(serviceCidrs *clusterv1.NetworkRanges) (*ekstypes.KubernetesNetworkConfigRequest, error) {
	if serviceCidrs == nil || len(serviceCidrs.CIDRBlocks) == 0 {
		return nil, nil
	}

	ipv4cidrs, err := cidr.GetIPv4Cidrs(serviceCidrs.CIDRBlocks)
	if err != nil {
		return nil, fmt.Errorf("filtering service cidr blocks to IPv4: %w", err)
	}

	if len(ipv4cidrs) == 0 {
		return nil, nil
	}

	return &ekstypes.KubernetesNetworkConfigRequest{
		ServiceIpv4Cidr: &ipv4cidrs[0],
	}, nil
}

func makeVpcConfig(subnets infrav1.Subnets, endpointAccess ekscontrolplanev1.EndpointAccess, securityGroups map[infrav1.SecurityGroupRole]infrav1.SecurityGroup) (*ekstypes.VpcConfigRequest, error) {
	// TODO: Do we need to just add the private subnets?
	if len(subnets) < 2 {
		return nil, awserrors.NewFailedDependency("at least 2 subnets is required")
	}

	if zones := subnets.GetUniqueZones(); len(zones) < 2 {
		return nil, awserrors.NewFailedDependency("subnets in at least 2 different az's are required")
	}

	subnetIDs := make([]string, 0)
	for i := range subnets {
		subnet := subnets[i]
		subnetID := subnet.GetResourceID()
		subnetIDs = append(subnetIDs, subnetID)
	}

	cidrs := make([]string, 0)
	for _, cidr := range endpointAccess.PublicCIDRs {
		_, ipNet, err := net.ParseCIDR(*cidr)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't parse PublicCIDRs")
		}
		parsedCIDR := ipNet.String()
		cidrs = append(cidrs, parsedCIDR)
	}

	vpcConfig := &ekstypes.VpcConfigRequest{
		EndpointPublicAccess:  endpointAccess.Public,
		EndpointPrivateAccess: endpointAccess.Private,
		SubnetIds:             subnetIDs,
	}

	if len(cidrs) > 0 {
		vpcConfig.PublicAccessCidrs = cidrs
	}
	sg, ok := securityGroups[infrav1.SecurityGroupEKSNodeAdditional]
	if ok {
		vpcConfig.SecurityGroupIds = append(vpcConfig.SecurityGroupIds, sg.ID)
	}
	return vpcConfig, nil
}

func makeEksLogging(loggingSpec *ekscontrolplanev1.ControlPlaneLoggingSpec) *ekstypes.Logging {
	if loggingSpec == nil {
		return nil
	}
	on := true
	off := false
	var enabledTypes []ekstypes.LogType
	var disabledTypes []ekstypes.LogType

	appendToTypes := func(logType ekstypes.LogType, field bool) {
		if field {
			enabledTypes = append(enabledTypes, logType)
		} else {
			disabledTypes = append(disabledTypes, logType)
		}
	}

	appendToTypes(ekstypes.LogTypeApi, loggingSpec.APIServer)
	appendToTypes(ekstypes.LogTypeAudit, loggingSpec.Audit)
	appendToTypes(ekstypes.LogTypeAuthenticator, loggingSpec.Authenticator)
	appendToTypes(ekstypes.LogTypeControllerManager, loggingSpec.ControllerManager)
	appendToTypes(ekstypes.LogTypeScheduler, loggingSpec.Scheduler)

	var clusterLogging []ekstypes.LogSetup
	if len(enabledTypes) > 0 {
		enabled := ekstypes.LogSetup{
			Enabled: &on,
			Types:   enabledTypes,
		}
		clusterLogging = append(clusterLogging, enabled)
	}
	if len(disabledTypes) > 0 {
		disabled := ekstypes.LogSetup{
			Enabled: &off,
			Types:   disabledTypes,
		}
		clusterLogging = append(clusterLogging, disabled)
	}
	if len(clusterLogging) > 0 {
		return &ekstypes.Logging{
			ClusterLogging: clusterLogging,
		}
	}
	return nil
}

func (s *Service) createCluster(ctx context.Context, eksClusterName string) (*ekstypes.Cluster, error) {
	var (
		vpcConfig *ekstypes.VpcConfigRequest
		err       error
	)
	logging := makeEksLogging(s.scope.ControlPlane.Spec.Logging)
	encryptionConfigs := makeEksEncryptionConfigs(s.scope.ControlPlane.Spec.EncryptionConfig)
	if s.scope.ControlPlane.Spec.RestrictPrivateSubnets {
		s.scope.Info("Filtering private subnets")
		vpcConfig, err = makeVpcConfig(s.scope.Subnets().FilterPrivate(), s.scope.ControlPlane.Spec.EndpointAccess, s.scope.SecurityGroups())
	} else {
		vpcConfig, err = makeVpcConfig(s.scope.Subnets(), s.scope.ControlPlane.Spec.EndpointAccess, s.scope.SecurityGroups())
	}
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create vpc config for cluster")
	}

	var netConfig *ekstypes.KubernetesNetworkConfigRequest
	if s.scope.VPC().IsIPv6Enabled() {
		netConfig = &ekstypes.KubernetesNetworkConfigRequest{
			IpFamily: ekstypes.IpFamilyIpv6,
		}
	} else {
		netConfig, err = makeKubernetesNetworkConfig(s.scope.ServiceCidrs())
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create Kubernetes network config for cluster")
		}
	}

	// Make sure to use the MachineScope here to get the merger of AWSCluster and AWSMachine tags
	additionalTags := s.scope.AdditionalTags()

	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(eksClusterName)] = string(infrav1.ResourceLifecycleOwned)
	tags := make(map[string]string)
	for k, v := range additionalTags {
		tagValue := v
		tags[k] = tagValue
	}

	role, err := s.GetIAMRole(ctx, *s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting control plane iam role: %s", *s.scope.ControlPlane.Spec.RoleName)
	}

	var eksVersion *string
	if s.scope.ControlPlane.Spec.Version != nil {
		specVersion, err := parseEKSVersion(*s.scope.ControlPlane.Spec.Version)
		if err != nil {
			return nil, fmt.Errorf("parsing EKS version from spec: %w", err)
		}
		v := versionToEKS(specVersion)
		eksVersion = &v
	}

	bootstrapAddon := s.scope.BootstrapSelfManagedAddons()
	input := &eks.CreateClusterInput{
		Name:                       aws.String(eksClusterName),
		Version:                    eksVersion,
		Logging:                    logging,
		EncryptionConfig:           encryptionConfigs,
		ResourcesVpcConfig:         vpcConfig,
		RoleArn:                    role.Arn,
		Tags:                       tags,
		KubernetesNetworkConfig:    netConfig,
		BootstrapSelfManagedAddons: bootstrapAddon,
	}

	var out *eks.CreateClusterOutput
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if out, err = s.EKSClient.CreateCluster(ctx, input); err != nil {
			return false, err
		}
		conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneCreatingCondition)
		record.Eventf(s.scope.ControlPlane, "InitiatedCreateEKSControlPlane", "Initiated creation of a new EKS control plane %s", s.scope.KubernetesClusterName())
		return true, nil
	}, awserrors.ResourceNotFound); err != nil { // TODO: change the error that can be retried
		record.Warnf(s.scope.ControlPlane, "FailedCreateEKSControlPlane", "Failed to initiate creation of a new EKS control plane: %v", err)
		return nil, errors.Wrapf(err, "failed to create EKS cluster")
	}

	s.scope.Info("Created EKS cluster in AWS", "cluster", klog.KRef("", eksClusterName))
	return out.Cluster, nil
}

func (s *Service) waitForClusterActive(ctx context.Context) (*ekstypes.Cluster, error) {
	eksClusterName := s.scope.KubernetesClusterName()
	req := eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
	}
	if err := s.EKSClient.WaitUntilClusterActive(ctx, &req, s.scope.MaxWaitActiveUpdateDelete); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for eks control plane %q", *req.Name)
	}

	s.scope.Info("EKS control plane is now active", "cluster", klog.KRef("", eksClusterName))

	cluster, err := s.describeEKSCluster(ctx, eksClusterName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe eks clusters")
	}

	if err := s.setStatus(cluster); err != nil {
		return nil, errors.Wrap(err, "failed to set status")
	}

	return cluster, nil
}

func (s *Service) reconcileClusterConfig(ctx context.Context, cluster *ekstypes.Cluster) error {
	var needsUpdate bool
	input := &eks.UpdateClusterConfigInput{Name: aws.String(s.scope.KubernetesClusterName())}

	updateVpcConfig, err := s.reconcileVpcConfig(cluster.ResourcesVpcConfig)
	if err != nil {
		return errors.Wrap(err, "couldn't create vpc config for cluster")
	}
	if updateVpcConfig != nil {
		needsUpdate = true
		input.ResourcesVpcConfig = updateVpcConfig
	}

	if needsUpdate {
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateClusterConfig(ctx, input); err != nil {
				return false, err
			}
			conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneUpdatingCondition)
			record.Eventf(s.scope.ControlPlane, "InitiatedUpdateEKSControlPlane", "Initiated update of a new EKS control plane %s", s.scope.KubernetesClusterName())
			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "Failed to update the EKS control plane: %v", err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}
	}
	return nil
}

func (s *Service) reconcileLogging(ctx context.Context, logging *ekstypes.Logging) error {
	input := &eks.UpdateClusterConfigInput{Name: aws.String(s.scope.KubernetesClusterName())}

	for _, logSetup := range logging.ClusterLogging {
		for _, l := range logSetup.Types {
			enabled := s.scope.ControlPlane.Spec.Logging.IsLogEnabled(string(l))
			if enabled != *logSetup.Enabled {
				input.Logging = makeEksLogging(s.scope.ControlPlane.Spec.Logging)
			}
		}
	}

	if input.Logging != nil {
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateClusterConfig(ctx, input); err != nil {
				return false, err
			}
			conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneUpdatingCondition)
			record.Eventf(s.scope.ControlPlane, "InitiatedUpdateEKSControlPlane", "Initiated logging update for EKS control plane %s", s.scope.KubernetesClusterName())
			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "Failed to update EKS control plane logging: %v", err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}
	}

	return nil
}

func publicAccessCIDRsEqual(as []string, bs []string) bool {
	allV4 := "0.0.0.0/0"
	allV6 := "::/0"
	asDefault := false
	bsDefault := false

	// For IPv6 only clusters
	if len(as) == 0 {
		as = []string{allV4, allV6}
		asDefault = true
	}
	if len(bs) == 0 {
		bs = []string{allV4, allV6}
		bsDefault = true
	}
	if sets.NewString(as...).Equal(sets.NewString(bs...)) {
		return true
	}

	// For IPv4 only clusters
	if asDefault {
		as = []string{allV4}
	}
	if bsDefault {
		bs = []string{allV4}
	}
	return sets.NewString(as...).Equal(
		sets.NewString(bs...),
	)
}

func (s *Service) reconcileVpcConfig(vpcConfig *ekstypes.VpcConfigResponse) (*ekstypes.VpcConfigRequest, error) {
	var (
		updatedVpcConfig *ekstypes.VpcConfigRequest
		err              error
	)
	endpointAccess := s.scope.ControlPlane.Spec.EndpointAccess
	if s.scope.ControlPlane.Spec.RestrictPrivateSubnets {
		updatedVpcConfig, err = makeVpcConfig(s.scope.Subnets().FilterPrivate(), endpointAccess, s.scope.SecurityGroups())
	} else {
		updatedVpcConfig, err = makeVpcConfig(s.scope.Subnets(), endpointAccess, s.scope.SecurityGroups())
	}
	if err != nil {
		return nil, err
	}
	needsUpdate := !tristate.EqualWithDefault(false, &vpcConfig.EndpointPrivateAccess, updatedVpcConfig.EndpointPrivateAccess) ||
		!tristate.EqualWithDefault(true, &vpcConfig.EndpointPublicAccess, updatedVpcConfig.EndpointPublicAccess) ||
		!publicAccessCIDRsEqual(vpcConfig.PublicAccessCidrs, updatedVpcConfig.PublicAccessCidrs)
	if needsUpdate {
		return &ekstypes.VpcConfigRequest{
			EndpointPublicAccess:  updatedVpcConfig.EndpointPublicAccess,
			EndpointPrivateAccess: updatedVpcConfig.EndpointPrivateAccess,
			PublicAccessCidrs:     updatedVpcConfig.PublicAccessCidrs,
		}, nil
	}
	return nil, nil
}

func (s *Service) reconcileEKSEncryptionConfig(ctx context.Context, currentClusterConfig []ekstypes.EncryptionConfig) error {
	s.Info("reconciling encryption configuration")
	if currentClusterConfig == nil {
		currentClusterConfig = []ekstypes.EncryptionConfig{}
	}

	encryptionConfigs := s.scope.ControlPlane.Spec.EncryptionConfig
	updatedEncryptionConfigs := makeEksEncryptionConfigs(encryptionConfigs)

	if compareEncryptionConfig(currentClusterConfig, updatedEncryptionConfigs) {
		s.Debug("encryption configuration unchanged, no action")
		return nil
	}

	if len(currentClusterConfig) == 0 && len(updatedEncryptionConfigs) > 0 {
		s.Debug("enabling encryption for eks cluster", "cluster", s.scope.KubernetesClusterName())
		if err := s.updateEncryptionConfig(ctx, updatedEncryptionConfigs); err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "failed to update the EKS control plane encryption configuration: %v", err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}

		return nil
	}

	record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "failed to update the EKS control plane: disabling EKS encryption is not allowed after it has been enabled")
	return errors.Errorf("failed to update the EKS control plane: disabling EKS encryption is not allowed after it has been enabled")
}

func parseEKSVersion(raw string) (*version.Version, error) {
	v, err := version.ParseGeneric(raw)
	if err != nil {
		return nil, err
	}
	return version.MustParseGeneric(fmt.Sprintf("%d.%d", v.Major(), v.Minor())), nil
}

func versionToEKS(v *version.Version) string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

func (s *Service) reconcileClusterVersion(ctx context.Context, cluster *ekstypes.Cluster) error {
	var specVersion *version.Version
	if s.scope.ControlPlane.Spec.Version != nil {
		var err error
		specVersion, err = parseEKSVersion(*s.scope.ControlPlane.Spec.Version)
		if err != nil {
			return fmt.Errorf("parsing EKS version from spec: %w", err)
		}
	}

	clusterVersion := version.MustParseGeneric(*cluster.Version)

	if specVersion != nil && clusterVersion.LessThan(specVersion) {
		// NOTE: you can only upgrade increments of minor versions. If you want to upgrade 1.14 to 1.16 we
		// need to go 1.14-> 1.15 and then 1.15 -> 1.16.
		nextVersionString := versionToEKS(clusterVersion.WithMinor(clusterVersion.Minor() + 1))

		input := &eks.UpdateClusterVersionInput{
			Name:    aws.String(s.scope.KubernetesClusterName()),
			Version: &nextVersionString,
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EKSClient.UpdateClusterVersion(ctx, input); err != nil {
				return false, err
			}

			// Wait until status transitions to UPDATING because there's a short
			// window after UpdateClusterVersion returns where the cluster
			// status is ACTIVE and the update would be tried again
			if err := s.EKSClient.WaitUntilClusterUpdating(
				ctx,
				&eks.DescribeClusterInput{Name: aws.String(s.scope.KubernetesClusterName())},
				s.scope.MaxWaitActiveUpdateDelete,
			); err != nil {
				return false, err
			}

			conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneUpdatingCondition)
			record.Eventf(s.scope.ControlPlane, "InitiatedUpdateEKSControlPlane", "Initiated update of EKS control plane %s to version %s", s.scope.KubernetesClusterName(), nextVersionString)

			return true, nil
		}); err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedUpdateEKSControlPlane", "failed to update the EKS control plane: %v", err)
			return errors.Wrapf(err, "failed to update EKS cluster")
		}
	}
	return nil
}

func (s *Service) describeEKSCluster(ctx context.Context, eksClusterName string) (*ekstypes.Cluster, error) {
	input := &eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
	}

	out, err := s.EKSClient.DescribeCluster(ctx, input)
	if err != nil {
		smithyErr := awserrors.ParseSmithyError(err)
		notFoundErr := &ekstypes.ResourceNotFoundException{}
		if smithyErr.ErrorCode() == notFoundErr.ErrorCode() {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to describe cluster")
	}

	return out.Cluster, nil
}

func (s *Service) updateEncryptionConfig(ctx context.Context, updatedEncryptionConfigs []ekstypes.EncryptionConfig) error {
	input := &eks.AssociateEncryptionConfigInput{
		ClusterName:      aws.String(s.scope.KubernetesClusterName()),
		EncryptionConfig: updatedEncryptionConfigs,
	}
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.EKSClient.AssociateEncryptionConfig(ctx, input); err != nil {
			return false, err
		}

		// Wait until status transitions to UPDATING because there's a short
		// window after UpdateClusterVersion returns where the cluster
		// status is ACTIVE and the update would be tried again
		if err := s.EKSClient.WaitUntilClusterUpdating(
			ctx,
			&eks.DescribeClusterInput{Name: aws.String(s.scope.KubernetesClusterName())},
			s.scope.MaxWaitActiveUpdateDelete,
		); err != nil {
			return false, err
		}

		conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneUpdatingCondition)
		record.Eventf(s.scope.ControlPlane, "InitiatedUpdateEncryptionConfig", "Initiated update of encryption config in EKS control plane %s", s.scope.KubernetesClusterName())

		return true, nil
	}); err != nil {
		return err
	}
	return nil
}

func compareEncryptionConfig(updatedEncryptionConfig, existingEncryptionConfig []ekstypes.EncryptionConfig) bool {
	if len(updatedEncryptionConfig) != len(existingEncryptionConfig) {
		return false
	}
	for index := range updatedEncryptionConfig {
		encryptionConfig := updatedEncryptionConfig[index]

		if getKeyArn(encryptionConfig) != getKeyArn(existingEncryptionConfig[index]) {
			return false
		}
		if !cmp.Equals(aws.StringSlice(encryptionConfig.Resources), aws.StringSlice(existingEncryptionConfig[index].Resources)) {
			return false
		}
	}
	return true
}

func getKeyArn(encryptionConfig ekstypes.EncryptionConfig) string {
	if encryptionConfig.Provider != nil {
		return aws.ToString(encryptionConfig.Provider.KeyArn)
	}
	return ""
}

// WaitUntilClusterActive is blocking function to wait until EKS Cluster is Active.
func (k *EKSClient) WaitUntilClusterActive(ctx context.Context, input *eks.DescribeClusterInput, maxWait time.Duration) error {
	waiter := eks.NewClusterActiveWaiter(k, func(o *eks.ClusterActiveWaiterOptions) {
		o.LogWaitAttempts = true
	})

	return waiter.Wait(ctx, input, maxWait)
}

// WaitUntilClusterDeleted is blocking function to wait until EKS Cluster is Deleted.
func (k *EKSClient) WaitUntilClusterDeleted(ctx context.Context, input *eks.DescribeClusterInput, maxWait time.Duration) error {
	waiter := eks.NewClusterDeletedWaiter(k)

	return waiter.Wait(ctx, input, maxWait)
}

// WaitUntilClusterUpdating is blocking function to wait until EKS Cluster is Updating.
func (k *EKSClient) WaitUntilClusterUpdating(ctx context.Context, input *eks.DescribeClusterInput, maxWait time.Duration) error {
	waiter := eks.NewClusterActiveWaiter(k, func(o *eks.ClusterActiveWaiterOptions) {
		o.LogWaitAttempts = true
		o.Retryable = clusterUpdatingStateRetryable
	})

	return waiter.Wait(ctx, input, maxWait)
}

// clusterUpdatingStateRetryable is adapted from aws-sdk-go-v2/service/eks/api_op_DescribeCluster.go.
func clusterUpdatingStateRetryable(ctx context.Context, input *eks.DescribeClusterInput, output *eks.DescribeClusterOutput, err error) (bool, error) {
	if err == nil {
		v1 := output.Cluster
		var v2 ekstypes.ClusterStatus
		if v1 != nil {
			v3 := v1.Status
			v2 = v3
		}
		expectedValue := "DELETING"
		pathValue := string(v2)
		if pathValue == expectedValue {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	if err == nil {
		v1 := output.Cluster
		var v2 ekstypes.ClusterStatus
		if v1 != nil {
			v3 := v1.Status
			v2 = v3
		}
		expectedValue := "FAILED"
		pathValue := string(v2)
		if pathValue == expectedValue {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	if err == nil {
		v1 := output.Cluster
		var v2 ekstypes.ClusterStatus
		if v1 != nil {
			v3 := v1.Status
			v2 = v3
		}
		expectedValue := "UPDATING"
		pathValue := string(v2)
		if pathValue == expectedValue {
			return false, nil
		}
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
