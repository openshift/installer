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

package elb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/hash"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ResourceGroups are filtered by ARN identifier: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#arns-syntax
// this is the identifier for classic ELBs: https://docs.aws.amazon.com/IAM/latest/UserGuide/list_elasticloadbalancing.html#elasticloadbalancing-resources-for-iam-policies
const elbResourceType = "elasticloadbalancing:loadbalancer"

// maxELBsDescribeTagsRequest is the maximum number of loadbalancers for the DescribeTags API call
// see: https://docs.aws.amazon.com/elasticloadbalancing/2012-06-01/APIReference/API_DescribeTags.html
const maxELBsDescribeTagsRequest = 20

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers() error {
	s.scope.Debug("Reconciling load balancers")

	// do a switch and reconcile different load-balancer types
	switch s.scope.ControlPlaneLoadBalancer().LoadBalancerType {
	case infrav1.LoadBalancerTypeClassic:
		return s.reconcileClassicLoadBalancer()
	case infrav1.LoadBalancerTypeNLB, infrav1.LoadBalancerTypeALB, infrav1.LoadBalancerTypeELB:
		return s.reconcileV2LB()
	default:
		return fmt.Errorf("unknown or unsupported load balancer type: %s", s.scope.ControlPlaneLoadBalancer().LoadBalancerType)
	}
}

// reconcileV2LB creates a load balancer. It also takes care of generating unique names across
// namespaces by appending the namespace to the name.
func (s *Service) reconcileV2LB() error {
	name, err := LBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}

	// Get default api server spec.
	spec, err := s.getAPIServerLBSpec(name)
	if err != nil {
		return err
	}
	lb, err := s.describeLB(name)
	switch {
	case IsNotFound(err) && s.scope.ControlPlaneEndpoint().IsValid():
		// if elb is not found and owner cluster ControlPlaneEndpoint is already populated, then we should not recreate the elb.
		return errors.Wrapf(err, "no loadbalancer exists for the AWSCluster %s, the cluster has become unrecoverable and should be deleted manually", s.scope.InfraClusterName())
	case IsNotFound(err):
		lb, err = s.createLB(spec)
		if err != nil {
			s.scope.Error(err, "failed to create LB")
			return err
		}

		s.scope.Debug("Created new network load balancer for apiserver", "api-server-lb-name", lb.Name)
	case err != nil:
		// Failed to describe the classic ELB
		return err
	}

	// set up the type for later processing
	lb.LoadBalancerType = s.scope.ControlPlaneLoadBalancer().LoadBalancerType
	if lb.IsManaged(s.scope.Name()) {
		if !cmp.Equal(spec.ELBAttributes, lb.ELBAttributes) {
			if err := s.configureLBAttributes(lb.ARN, spec.ELBAttributes); err != nil {
				return err
			}
		}

		if err := s.reconcileV2LBTags(lb, spec.Tags); err != nil {
			return errors.Wrapf(err, "failed to reconcile tags for apiserver load balancer %q", lb.Name)
		}

		// Reconcile the subnets and availability zones from the spec
		// and the ones currently attached to the load balancer.
		if len(lb.SubnetIDs) != len(spec.SubnetIDs) {
			_, err := s.ELBV2Client.SetSubnets(&elbv2.SetSubnetsInput{
				LoadBalancerArn: &lb.ARN,
				Subnets:         aws.StringSlice(spec.SubnetIDs),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to set subnets for apiserver load balancer '%s'", lb.Name)
			}
		}
		if len(lb.AvailabilityZones) != len(spec.AvailabilityZones) {
			lb.AvailabilityZones = spec.AvailabilityZones
		}

		// Reconcile the security groups from the spec and the ones currently attached to the load balancer
		if s.scope.ControlPlaneLoadBalancer().LoadBalancerType != infrav1.LoadBalancerTypeNLB && !sets.NewString(lb.SecurityGroupIDs...).Equal(sets.NewString(spec.SecurityGroupIDs...)) {
			_, err := s.ELBV2Client.SetSecurityGroups(&elbv2.SetSecurityGroupsInput{
				LoadBalancerArn: &lb.ARN,
				SecurityGroups:  aws.StringSlice(spec.SecurityGroupIDs),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to apply security groups to load balancer %q", lb.Name)
			}
		}
	} else {
		s.scope.Trace("Unmanaged control plane load balancer, skipping load balancer configuration", "api-server-elb", lb)
	}
	lb.DeepCopyInto(&s.scope.Network().APIServerELB)
	return nil
}

func (s *Service) getAPIServerLBSpec(elbName string) (*infrav1.LoadBalancer, error) {
	var securityGroupIDs []string
	controlPlaneLoadBalancer := s.scope.ControlPlaneLoadBalancer()
	if controlPlaneLoadBalancer != nil && controlPlaneLoadBalancer.LoadBalancerType != infrav1.LoadBalancerTypeNLB {
		securityGroupIDs = append(securityGroupIDs, controlPlaneLoadBalancer.AdditionalSecurityGroups...)
		securityGroupIDs = append(securityGroupIDs, s.scope.SecurityGroups()[infrav1.SecurityGroupAPIServerLB].ID)
	}

	res := &infrav1.LoadBalancer{
		Name:          elbName,
		Scheme:        s.scope.ControlPlaneLoadBalancerScheme(),
		ELBAttributes: make(map[string]*string),
		ELBListeners: []infrav1.Listener{
			{
				Protocol: infrav1.ELBProtocolTCP,
				Port:     infrav1.DefaultAPIServerPort,
				TargetGroup: infrav1.TargetGroupSpec{
					Name:     fmt.Sprintf("apiserver-target-%d", time.Now().Unix()),
					Port:     infrav1.DefaultAPIServerPort,
					Protocol: infrav1.ELBProtocolTCP,
					VpcID:    s.scope.VPC().ID,
					HealthCheck: &infrav1.TargetGroupHealthCheck{
						Protocol: aws.String(string(infrav1.ELBProtocolTCP)),
						Port:     aws.String(infrav1.DefaultAPIServerPortString),
					},
				},
			},
		},
		SecurityGroupIDs: securityGroupIDs,
	}

	if s.scope.ControlPlaneLoadBalancer() != nil {
		for _, additionalListeners := range controlPlaneLoadBalancer.AdditionalListeners {
			res.ELBListeners = append(res.ELBListeners, infrav1.Listener{
				Protocol: additionalListeners.Protocol,
				Port:     additionalListeners.Port,
				TargetGroup: infrav1.TargetGroupSpec{
					Name:     fmt.Sprintf("additional-listener-%d", time.Now().Unix()),
					Port:     additionalListeners.Port,
					Protocol: additionalListeners.Protocol,
					VpcID:    s.scope.VPC().ID,
					HealthCheck: &infrav1.TargetGroupHealthCheck{
						Protocol: aws.String(string(additionalListeners.Protocol)),
						Port:     aws.String(fmt.Sprintf("%d", additionalListeners.Port)),
					},
				},
			})
		}
	}

	if s.scope.ControlPlaneLoadBalancer() != nil && s.scope.ControlPlaneLoadBalancer().LoadBalancerType != infrav1.LoadBalancerTypeNLB {
		res.ELBAttributes[infrav1.LoadBalancerAttributeIdleTimeTimeoutSeconds] = aws.String(infrav1.LoadBalancerAttributeIdleTimeDefaultTimeoutSecondsInSeconds)
	}

	if s.scope.ControlPlaneLoadBalancer() != nil {
		res.ELBAttributes[infrav1.LoadBalancerAttributeEnableLoadBalancingCrossZone] = aws.String(fmt.Sprintf("%t", s.scope.ControlPlaneLoadBalancer().CrossZoneLoadBalancing))
	}

	res.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(elbName),
		Role:        aws.String(infrav1.APIServerRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	})

	// If subnet IDs have been specified for this load balancer
	if s.scope.ControlPlaneLoadBalancer() != nil && len(s.scope.ControlPlaneLoadBalancer().Subnets) > 0 {
		// This set of subnets may not match the subnets specified on the Cluster, so we may not have already discovered them
		// We need to call out to AWS to describe them just in case
		input := &ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice(s.scope.ControlPlaneLoadBalancer().Subnets),
		}
		out, err := s.EC2Client.DescribeSubnetsWithContext(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, sn := range out.Subnets {
			res.AvailabilityZones = append(res.AvailabilityZones, *sn.AvailabilityZone)
			res.SubnetIDs = append(res.SubnetIDs, *sn.SubnetId)
		}
	} else {
		// The load balancer APIs require us to only attach one subnet for each AZ.
		subnets := s.scope.Subnets().FilterPrivate()

		if s.scope.ControlPlaneLoadBalancerScheme() == infrav1.ELBSchemeInternetFacing {
			subnets = s.scope.Subnets().FilterPublic()
		}

	subnetLoop:
		for _, sn := range subnets {
			for _, az := range res.AvailabilityZones {
				if sn.AvailabilityZone == az {
					// If we already attached another subnet in the same AZ, there is no need to
					// add this subnet to the list of the ELB's subnets.
					continue subnetLoop
				}
			}
			res.AvailabilityZones = append(res.AvailabilityZones, sn.AvailabilityZone)
			res.SubnetIDs = append(res.SubnetIDs, sn.GetResourceID())
		}
	}

	return res, nil
}

func (s *Service) createLB(spec *infrav1.LoadBalancer) (*infrav1.LoadBalancer, error) {
	var t *string
	switch s.scope.ControlPlaneLoadBalancer().LoadBalancerType {
	case infrav1.LoadBalancerTypeNLB:
		t = aws.String(elbv2.LoadBalancerTypeEnumNetwork)
	case infrav1.LoadBalancerTypeALB:
		t = aws.String(elbv2.LoadBalancerTypeEnumApplication)
	case infrav1.LoadBalancerTypeELB:
		t = aws.String(elbv2.LoadBalancerTypeEnumGateway)
	}
	input := &elbv2.CreateLoadBalancerInput{
		Name:    aws.String(spec.Name),
		Subnets: aws.StringSlice(spec.SubnetIDs),
		Tags:    converters.MapToV2Tags(spec.Tags),
		Scheme:  aws.String(string(spec.Scheme)),
		Type:    t,
	}
	if s.scope.ControlPlaneLoadBalancer().LoadBalancerType != infrav1.LoadBalancerTypeNLB {
		input.SecurityGroups = aws.StringSlice(spec.SecurityGroupIDs)
	}

	if s.scope.VPC().IsIPv6Enabled() {
		input.IpAddressType = aws.String("dualstack")
	}

	out, err := s.ELBV2Client.CreateLoadBalancer(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create load balancer: %v", spec)
	}

	if len(out.LoadBalancers) == 0 {
		return nil, errors.New("no new network load balancer was created; the returned list is empty")
	}

	// TODO(Skarlso): Add options to set up SSL.
	// https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/3899
	for _, ln := range spec.ELBListeners {
		// create the target group first
		targetGroupInput := &elbv2.CreateTargetGroupInput{
			Name:     aws.String(ln.TargetGroup.Name),
			Port:     aws.Int64(ln.TargetGroup.Port),
			Protocol: aws.String(ln.TargetGroup.Protocol.String()),
			VpcId:    aws.String(ln.TargetGroup.VpcID),
			Tags:     input.Tags,
		}
		if s.scope.VPC().IsIPv6Enabled() {
			targetGroupInput.IpAddressType = aws.String("ipv6")
		}
		if ln.TargetGroup.HealthCheck != nil {
			targetGroupInput.HealthCheckEnabled = aws.Bool(true)
			targetGroupInput.HealthCheckProtocol = ln.TargetGroup.HealthCheck.Protocol
			targetGroupInput.HealthCheckPort = ln.TargetGroup.HealthCheck.Port
		}
		s.scope.Debug("creating target group", "group", targetGroupInput, "listener", ln)
		group, err := s.ELBV2Client.CreateTargetGroup(targetGroupInput)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create target group for load balancer")
		}
		if len(group.TargetGroups) == 0 {
			return nil, errors.New("no target group was created; the returned list is empty")
		}

		if !s.scope.ControlPlaneLoadBalancer().PreserveClientIP {
			targetGroupAttributeInput := &elbv2.ModifyTargetGroupAttributesInput{
				TargetGroupArn: group.TargetGroups[0].TargetGroupArn,
				Attributes: []*elbv2.TargetGroupAttribute{
					{
						Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
						Value: aws.String("false"),
					},
				},
			}
			if _, err := s.ELBV2Client.ModifyTargetGroupAttributes(targetGroupAttributeInput); err != nil {
				return nil, errors.Wrapf(err, "failed to modify target group attribute")
			}
		}

		listenerInput := &elbv2.CreateListenerInput{
			DefaultActions: []*elbv2.Action{
				{
					TargetGroupArn: group.TargetGroups[0].TargetGroupArn,
					Type:           aws.String(elbv2.ActionTypeEnumForward),
				},
			},
			LoadBalancerArn: out.LoadBalancers[0].LoadBalancerArn,
			Port:            aws.Int64(ln.Port),
			Protocol:        aws.String(string(ln.Protocol)),
			Tags:            converters.MapToV2Tags(spec.Tags),
		}
		// Create ClassicELBListeners
		listener, err := s.ELBV2Client.CreateListener(listenerInput)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create listener")
		}
		if len(listener.Listeners) == 0 {
			return nil, errors.New("no listener was created; the returned list is empty")
		}
	}

	s.scope.Info("Created network load balancer", "dns-name", *out.LoadBalancers[0].DNSName)

	res := spec.DeepCopy()
	s.scope.Debug("applying load balancer DNS to result", "dns", *out.LoadBalancers[0].DNSName)
	res.DNSName = *out.LoadBalancers[0].DNSName
	return res, nil
}

func (s *Service) describeLB(name string) (*infrav1.LoadBalancer, error) {
	input := &elbv2.DescribeLoadBalancersInput{
		Names: aws.StringSlice([]string{name}),
	}

	out, err := s.ELBV2Client.DescribeLoadBalancers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				return nil, NewNotFound(fmt.Sprintf("no load balancer found with name: %q", name))
			case elb.ErrCodeDependencyThrottleException:
				return nil, errors.Wrap(err, "too many requests made to the ELB service")
			default:
				return nil, errors.Wrap(err, "unexpected aws error")
			}
		} else {
			return nil, errors.Wrapf(err, "failed to describe load balancer: %s", name)
		}
	}

	if out != nil && len(out.LoadBalancers) == 0 {
		return nil, NewNotFound(fmt.Sprintf("no load balancer found with name %q", name))
	}

	if s.scope.VPC().ID != "" && s.scope.VPC().ID != *out.LoadBalancers[0].VpcId {
		return nil, errors.Errorf(
			"Load balancer names must be unique within a region: %q load balancer already exists in this region in VPC %q",
			name, *out.LoadBalancers[0].VpcId)
	}

	if s.scope.ControlPlaneLoadBalancer() != nil &&
		s.scope.ControlPlaneLoadBalancer().Scheme != nil &&
		string(*s.scope.ControlPlaneLoadBalancer().Scheme) != aws.StringValue(out.LoadBalancers[0].Scheme) {
		return nil, errors.Errorf(
			"Load balancer names must be unique within a region: %q Load balancer already exists in this region with a different scheme %q",
			name, *out.LoadBalancers[0].Scheme)
	}

	outAtt, err := s.ELBV2Client.DescribeLoadBalancerAttributes(&elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: out.LoadBalancers[0].LoadBalancerArn,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe load balancer %q attributes", name)
	}

	tags, err := s.describeLBTags(aws.StringValue(out.LoadBalancers[0].LoadBalancerArn))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe load balancer tags")
	}

	return fromSDKTypeToLB(out.LoadBalancers[0], outAtt.Attributes, tags), nil
}

func (s *Service) reconcileClassicLoadBalancer() error {
	// Generate a default control plane load balancer name. The load balancer name cannot be
	// generated by the defaulting webhook, because it is derived from the cluster name, and that
	// name is undefined at defaulting time when generateName is used.
	name, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}

	// Get default api server spec.
	spec, err := s.getAPIServerClassicELBSpec(name)
	if err != nil {
		return err
	}

	apiELB, err := s.describeClassicELB(spec.Name)
	switch {
	case IsNotFound(err) && s.scope.ControlPlaneEndpoint().IsValid():
		// if elb is not found and owner cluster ControlPlaneEndpoint is already populated, then we should not recreate the elb.
		return errors.Wrapf(err, "no loadbalancer exists for the AWSCluster %s, the cluster has become unrecoverable and should be deleted manually", s.scope.InfraClusterName())
	case IsNotFound(err):
		apiELB, err = s.createClassicELB(spec)
		if err != nil {
			return err
		}
		s.scope.Debug("Created new classic load balancer for apiserver", "api-server-elb-name", apiELB.Name)
	case err != nil:
		// Failed to describe the classic ELB
		return err
	}

	if apiELB.IsManaged(s.scope.Name()) {
		if !cmp.Equal(spec.ClassicElbAttributes, apiELB.ClassicElbAttributes) {
			err := s.configureAttributes(apiELB.Name, spec.ClassicElbAttributes)
			if err != nil {
				return err
			}
		}

		if err := s.reconcileELBTags(apiELB, spec.Tags); err != nil {
			return errors.Wrapf(err, "failed to reconcile tags for apiserver load balancer %q", apiELB.Name)
		}

		// Reconcile the subnets and availability zones from the spec
		// and the ones currently attached to the load balancer.
		if len(apiELB.SubnetIDs) != len(spec.SubnetIDs) {
			_, err := s.ELBClient.AttachLoadBalancerToSubnets(&elb.AttachLoadBalancerToSubnetsInput{
				LoadBalancerName: &apiELB.Name,
				Subnets:          aws.StringSlice(spec.SubnetIDs),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to attach apiserver load balancer %q to subnets", apiELB.Name)
			}
		}

		// Reconcile the security groups from the spec and the ones currently attached to the load balancer
		if !sets.NewString(apiELB.SecurityGroupIDs...).Equal(sets.NewString(spec.SecurityGroupIDs...)) {
			_, err := s.ELBClient.ApplySecurityGroupsToLoadBalancer(&elb.ApplySecurityGroupsToLoadBalancerInput{
				LoadBalancerName: &apiELB.Name,
				SecurityGroups:   aws.StringSlice(spec.SecurityGroupIDs),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to apply security groups to load balancer %q", apiELB.Name)
			}
		}
	} else {
		s.scope.Trace("Unmanaged control plane load balancer, skipping load balancer configuration", "api-server-elb", apiELB)
	}

	if len(apiELB.AvailabilityZones) != len(spec.AvailabilityZones) {
		apiELB.AvailabilityZones = spec.AvailabilityZones
	}

	// TODO(vincepri): check if anything has changed and reconcile as necessary.
	apiELB.DeepCopyInto(&s.scope.Network().APIServerELB)
	s.scope.Trace("Control plane load balancer", "api-server-elb", apiELB)

	s.scope.Debug("Reconcile load balancers completed successfully")
	return nil
}

func (s *Service) deleteAPIServerELB() error {
	s.scope.Debug("Deleting control plane load balancer")

	elbName, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}

	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	apiELB, err := s.describeClassicELB(elbName)
	if IsNotFound(err) {
		s.scope.Debug("Control plane load balancer not found, skipping deletion")
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
		return nil
	}
	if err != nil {
		return err
	}

	if apiELB.IsUnmanaged(s.scope.Name()) {
		s.scope.Debug("Found unmanaged classic load balancer for apiserver, skipping deletion", "api-server-elb-name", apiELB.Name)
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
		return nil
	}

	s.scope.Debug("deleting load balancer", "name", elbName)
	if err := s.deleteClassicELB(elbName); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		_, err = s.describeClassicELB(elbName)
		done = IsNotFound(err)
		return done, nil
	}); err != nil {
		return errors.Wrapf(err, "failed to wait for %q load balancer deletion", s.scope.Name())
	}

	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
	s.scope.Info("Deleted control plane load balancer", "name", elbName)
	return nil
}

// deleteAWSCloudProviderELBs deletes ELBs owned by the AWS Cloud Provider. For every
// LoadBalancer-type Service on the cluster, there is one ELB. If the Service is deleted before the
// cluster is deleted, its ELB is deleted; the ELBs found in this function will typically be for
// Services that were not deleted before the cluster was deleted.
func (s *Service) deleteAWSCloudProviderELBs() error {
	s.scope.Debug("Deleting AWS cloud provider load balancers (created for LoadBalancer-type Services)")

	elbs, err := s.listAWSCloudProviderOwnedELBs()
	if err != nil {
		return err
	}

	for _, elb := range elbs {
		s.scope.Debug("Deleting AWS cloud provider load balancer", "arn", elb)
		if err := s.deleteClassicELB(elb); err != nil {
			return err
		}
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		elbs, err := s.listAWSCloudProviderOwnedELBs()
		if err != nil {
			return false, err
		}
		done = len(elbs) == 0
		return done, nil
	}); err != nil {
		return errors.Wrapf(err, "failed to wait for %q load balancer deletions", s.scope.Name())
	}

	return nil
}

// DeleteLoadbalancers deletes the load balancers for the given cluster.
func (s *Service) DeleteLoadbalancers() error {
	s.scope.Debug("Deleting load balancers")

	if err := s.deleteAPIServerELB(); err != nil {
		return errors.Wrap(err, "failed to delete control plane load balancer")
	}

	if err := s.deleteAWSCloudProviderELBs(); err != nil {
		return errors.Wrap(err, "failed to delete AWS cloud provider load balancer(s)")
	}

	if err := s.deleteExistingNLBs(); err != nil {
		return errors.Wrap(err, "failed to delete AWS cloud provider load balancer(s)")
	}

	return nil
}

func (s *Service) deleteExistingNLBs() error {
	name, err := LBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	lb, err := s.describeLB(name)
	if IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if lb.IsUnmanaged(s.scope.Name()) {
		s.scope.Debug("Found unmanaged load balancer for apiserver, skipping deletion", "api-server-elb-name", lb.Name)
		return nil
	}
	s.scope.Debug("deleting load balancer", "name", name)
	if err := s.deleteLB(lb.ARN); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		_, err = s.describeLB(name)
		done = IsNotFound(err)
		return done, nil
	}); err != nil {
		return errors.Wrapf(err, "failed to wait for %q load balancer deletion", s.scope.Name())
	}

	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
	s.scope.Info("Deleted control plane load balancer", "name", name)

	return nil
}

// IsInstanceRegisteredWithAPIServerELB returns true if the instance is already registered with the APIServer ELB.
func (s *Service) IsInstanceRegisteredWithAPIServerELB(i *infrav1.Instance) (bool, error) {
	name, err := ELBName(s.scope)
	if err != nil {
		return false, errors.Wrap(err, "failed to get control plane load balancer name")
	}

	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(name)},
	}

	output, err := s.ELBClient.DescribeLoadBalancers(input)
	if err != nil {
		return false, errors.Wrapf(err, "error describing ELB %q", name)
	}
	if len(output.LoadBalancerDescriptions) != 1 {
		return false, errors.Errorf("expected 1 ELB description for %q, got %d", name, len(output.LoadBalancerDescriptions))
	}

	for _, registeredInstance := range output.LoadBalancerDescriptions[0].Instances {
		if aws.StringValue(registeredInstance.InstanceId) == i.ID {
			return true, nil
		}
	}

	return false, nil
}

// IsInstanceRegisteredWithAPIServerLB returns true if the instance is already registered with the APIServer LB.
func (s *Service) IsInstanceRegisteredWithAPIServerLB(i *infrav1.Instance) ([]string, bool, error) {
	name, err := LBName(s.scope)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get control plane load balancer name")
	}

	input := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{aws.String(name)},
	}

	output, err := s.ELBV2Client.DescribeLoadBalancers(input)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error describing ELB %q", name)
	}
	if len(output.LoadBalancers) != 1 {
		return nil, false, errors.Errorf("expected 1 ELB description for %q, got %d", name, len(output.LoadBalancers))
	}

	describeTargetGroupInput := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: output.LoadBalancers[0].LoadBalancerArn,
	}

	targetGroups, err := s.ELBV2Client.DescribeTargetGroups(describeTargetGroupInput)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error describing ELB's target groups %q", name)
	}

	targetGroupARNs := []string{}
	for _, tg := range targetGroups.TargetGroups {
		healthInput := &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		}
		instanceHealth, err := s.ELBV2Client.DescribeTargetHealth(healthInput)
		if err != nil {
			return nil, false, errors.Wrapf(err, "error describing ELB's target groups health %q", name)
		}
		for _, id := range instanceHealth.TargetHealthDescriptions {
			if aws.StringValue(id.Target.Id) == i.ID {
				targetGroupARNs = append(targetGroupARNs, aws.StringValue(tg.TargetGroupArn))
			}
		}
	}
	if len(targetGroupARNs) > 0 {
		return targetGroupARNs, true, nil
	}

	return nil, false, nil
}

// RegisterInstanceWithAPIServerELB registers an instance with a classic ELB.
func (s *Service) RegisterInstanceWithAPIServerELB(i *infrav1.Instance) error {
	name, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}
	out, err := s.describeClassicELB(name)
	if err != nil {
		return err
	}

	// Validate that the subnets associated with the load balancer has the instance AZ.
	subnet := s.scope.Subnets().FindByID(i.SubnetID)
	if subnet == nil {
		return errors.Errorf("failed to attach load balancer subnets, could not find subnet %q description in AWSCluster", i.SubnetID)
	}
	instanceAZ := subnet.AvailabilityZone

	var subnets infrav1.Subnets
	if s.scope.ControlPlaneLoadBalancer() != nil && len(s.scope.ControlPlaneLoadBalancer().Subnets) > 0 {
		subnets, err = s.getControlPlaneLoadBalancerSubnets()
		if err != nil {
			return err
		}
	} else {
		subnets = s.scope.Subnets()
	}

	found := false
	for _, subnetID := range out.SubnetIDs {
		if subnet := subnets.FindByID(subnetID); subnet != nil && instanceAZ == subnet.AvailabilityZone {
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("failed to register instance with APIServer ELB %q: instance is in availability zone %q, no public subnets attached to the ELB in the same zone", name, instanceAZ)
	}

	input := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        []*elb.Instance{{InstanceId: aws.String(i.ID)}},
		LoadBalancerName: aws.String(name),
	}

	_, err = s.ELBClient.RegisterInstancesWithLoadBalancer(input)
	return err
}

// RegisterInstanceWithAPIServerLB registers an instance with a LB.
func (s *Service) RegisterInstanceWithAPIServerLB(instance *infrav1.Instance) error {
	name, err := LBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}
	out, err := s.describeLB(name)
	if err != nil {
		return err
	}
	s.scope.Debug("found load balancer with name", "name", out.Name)
	describeTargetGroupInput := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(out.ARN),
	}

	targetGroups, err := s.ELBV2Client.DescribeTargetGroups(describeTargetGroupInput)
	if err != nil {
		return errors.Wrapf(err, "error describing ELB's target groups %q", name)
	}
	if len(targetGroups.TargetGroups) == 0 {
		return errors.New(fmt.Sprintf("no target groups found for load balancer with arn '%s'", out.ARN))
	}
	// Since TargetGroups and Listeners don't care, or are not aware, of subnets before registration, we ignore that check.
	// Also, registering with AZ is not supported using the an InstanceID.
	s.scope.Debug("found number of target groups", "target-groups", len(targetGroups.TargetGroups))
	for _, tg := range targetGroups.TargetGroups {
		input := &elbv2.RegisterTargetsInput{
			TargetGroupArn: tg.TargetGroupArn,
			Targets: []*elbv2.TargetDescription{
				{
					Id:   aws.String(instance.ID),
					Port: tg.Port,
				},
			},
		}
		if _, err = s.ELBV2Client.RegisterTargets(input); err != nil {
			return fmt.Errorf("failed to register instance with target group '%s': %w", aws.StringValue(tg.TargetGroupName), err)
		}
	}

	return nil
}

// getControlPlaneLoadBalancerSubnets retrieves ControlPlaneLoadBalancer subnets information.
func (s *Service) getControlPlaneLoadBalancerSubnets() (infrav1.Subnets, error) {
	var subnets infrav1.Subnets

	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(s.scope.ControlPlaneLoadBalancer().Subnets),
	}
	res, err := s.EC2Client.DescribeSubnetsWithContext(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	for _, sn := range res.Subnets {
		lbSn := infrav1.SubnetSpec{
			AvailabilityZone: *sn.AvailabilityZone,
			ID:               *sn.SubnetId,
			ResourceID:       *sn.SubnetId,
		}
		subnets = append(subnets, lbSn)
	}

	return subnets, nil
}

// DeregisterInstanceFromAPIServerELB de-registers an instance from a classic ELB.
func (s *Service) DeregisterInstanceFromAPIServerELB(i *infrav1.Instance) error {
	name, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}

	input := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances:        []*elb.Instance{{InstanceId: aws.String(i.ID)}},
		LoadBalancerName: aws.String(name),
	}

	_, err = s.ELBClient.DeregisterInstancesFromLoadBalancer(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException, elb.ErrCodeInvalidEndPointException:
				// Ignoring LoadBalancerNotFound and InvalidInstance when deregistering
				return nil
			default:
				return err
			}
		}
	}
	return err
}

// DeregisterInstanceFromAPIServerLB de-registers an instance from a LB.
func (s *Service) DeregisterInstanceFromAPIServerLB(targetGroupArn string, i *infrav1.Instance) error {
	input := &elbv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String(targetGroupArn),
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String(i.ID),
			},
		},
	}

	_, err := s.ELBV2Client.DeregisterTargets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException, elb.ErrCodeInvalidEndPointException:
				// Ignoring LoadBalancerNotFound and InvalidInstance when deregistering
				return nil
			default:
				return err
			}
		}
	}
	return err
}

// ELBName returns the user-defined API Server ELB name, or a generated default if the user has not defined the ELB
// name.
func ELBName(s scope.ELBScope) (string, error) {
	if userDefinedName := s.ControlPlaneLoadBalancerName(); userDefinedName != nil {
		return *userDefinedName, nil
	}
	name, err := GenerateELBName(s.Name())
	if err != nil {
		return "", fmt.Errorf("failed to generate name: %w", err)
	}
	return name, nil
}

// LBName returns the user-defined API Server ELB name, or a generated default if the user has not defined the ELB
// name.
func LBName(s scope.ELBScope) (string, error) {
	if userDefinedName := s.ControlPlaneLoadBalancerName(); userDefinedName != nil {
		return *userDefinedName, nil
	}
	name, err := GenerateELBName(fmt.Sprintf("%s-%s", s.Namespace(), s.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to generate name: %w", err)
	}
	return name, nil
}

// GenerateELBName generates a formatted ELB name via either
// concatenating the cluster name to the "-apiserver" suffix
// or computing a hash for clusters with names above 32 characters.
//
// WARNING If this function's output is changed, a controller using the
// new function will fail to generate the load balancer of an existing
// cluster whose load balancer name was generated using the old
// function.
func GenerateELBName(clusterName string) (string, error) {
	standardELBName := generateStandardELBName(clusterName)
	if len(standardELBName) <= 32 {
		return standardELBName, nil
	}

	elbName, err := generateHashedELBName(clusterName)
	if err != nil {
		return "", err
	}

	return elbName, nil
}

// generateStandardELBName generates a formatted ELB name based on cluster
// and ELB name.
func generateStandardELBName(clusterName string) string {
	elbCompatibleClusterName := strings.ReplaceAll(clusterName, ".", "-")
	return fmt.Sprintf("%s-%s", elbCompatibleClusterName, infrav1.APIServerRoleTagValue)
}

// generateHashedELBName generates a 32-character hashed name based on cluster
// and ELB name.
func generateHashedELBName(clusterName string) (string, error) {
	// hashSize = 32 - length of "k8s" - length of "-" = 28
	shortName, err := hash.Base36TruncatedHash(clusterName, 28)
	if err != nil {
		return "", errors.Wrap(err, "unable to create ELB name")
	}

	return fmt.Sprintf("%s-%s", shortName, "k8s"), nil
}

func (s *Service) getAPIServerClassicELBSpec(elbName string) (*infrav1.LoadBalancer, error) {
	securityGroupIDs := []string{}
	controlPlaneLoadBalancer := s.scope.ControlPlaneLoadBalancer()
	if controlPlaneLoadBalancer != nil && len(controlPlaneLoadBalancer.AdditionalSecurityGroups) != 0 {
		securityGroupIDs = append(securityGroupIDs, controlPlaneLoadBalancer.AdditionalSecurityGroups...)
	}
	securityGroupIDs = append(securityGroupIDs, s.scope.SecurityGroups()[infrav1.SecurityGroupAPIServerLB].ID)

	res := &infrav1.LoadBalancer{
		Name:   elbName,
		Scheme: s.scope.ControlPlaneLoadBalancerScheme(),
		ClassicELBListeners: []infrav1.ClassicELBListener{
			{
				Protocol:         infrav1.ELBProtocolTCP,
				Port:             int64(s.scope.APIServerPort()),
				InstanceProtocol: infrav1.ELBProtocolTCP,
				InstancePort:     infrav1.DefaultAPIServerPort,
			},
		},
		HealthCheck: &infrav1.ClassicELBHealthCheck{
			Target:             s.getHealthCheckTarget(),
			Interval:           10 * time.Second,
			Timeout:            5 * time.Second,
			HealthyThreshold:   5,
			UnhealthyThreshold: 3,
		},
		SecurityGroupIDs: securityGroupIDs,
		ClassicElbAttributes: infrav1.ClassicELBAttributes{
			IdleTimeout: 10 * time.Minute,
		},
	}

	if s.scope.ControlPlaneLoadBalancer() != nil {
		res.ClassicElbAttributes.CrossZoneLoadBalancing = s.scope.ControlPlaneLoadBalancer().CrossZoneLoadBalancing
	}

	res.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(elbName),
		Role:        aws.String(infrav1.APIServerRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	})

	// If subnet IDs have been specified for this load balancer
	if s.scope.ControlPlaneLoadBalancer() != nil && len(s.scope.ControlPlaneLoadBalancer().Subnets) > 0 {
		// This set of subnets may not match the subnets specified on the Cluster, so we may not have already discovered them
		// We need to call out to AWS to describe them just in case
		input := &ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice(s.scope.ControlPlaneLoadBalancer().Subnets),
		}
		out, err := s.EC2Client.DescribeSubnetsWithContext(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, sn := range out.Subnets {
			res.AvailabilityZones = append(res.AvailabilityZones, *sn.AvailabilityZone)
			res.SubnetIDs = append(res.SubnetIDs, *sn.SubnetId)
		}
	} else {
		// The load balancer APIs require us to only attach one subnet for each AZ.
		subnets := s.scope.Subnets().FilterPrivate()

		if s.scope.ControlPlaneLoadBalancerScheme() == infrav1.ELBSchemeInternetFacing {
			subnets = s.scope.Subnets().FilterPublic()
		}

	subnetLoop:
		for _, sn := range subnets {
			for _, az := range res.AvailabilityZones {
				if sn.AvailabilityZone == az {
					// If we already attached another subnet in the same AZ, there is no need to
					// add this subnet to the list of the ELB's subnets.
					continue subnetLoop
				}
			}
			res.AvailabilityZones = append(res.AvailabilityZones, sn.AvailabilityZone)
			res.SubnetIDs = append(res.SubnetIDs, sn.GetResourceID())
		}
	}

	return res, nil
}

func (s *Service) createClassicELB(spec *infrav1.LoadBalancer) (*infrav1.LoadBalancer, error) {
	input := &elb.CreateLoadBalancerInput{
		LoadBalancerName: aws.String(spec.Name),
		Subnets:          aws.StringSlice(spec.SubnetIDs),
		SecurityGroups:   aws.StringSlice(spec.SecurityGroupIDs),
		Scheme:           aws.String(string(spec.Scheme)),
		Tags:             converters.MapToELBTags(spec.Tags),
	}

	for _, ln := range spec.ClassicELBListeners {
		input.Listeners = append(input.Listeners, &elb.Listener{
			Protocol:         aws.String(string(ln.Protocol)),
			LoadBalancerPort: aws.Int64(ln.Port),
			InstanceProtocol: aws.String(string(ln.InstanceProtocol)),
			InstancePort:     aws.Int64(ln.InstancePort),
		})
	}

	out, err := s.ELBClient.CreateLoadBalancer(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create classic load balancer: %v", spec)
	}

	if spec.HealthCheck != nil {
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.ELBClient.ConfigureHealthCheck(&elb.ConfigureHealthCheckInput{
				LoadBalancerName: aws.String(spec.Name),
				HealthCheck: &elb.HealthCheck{
					Target:             aws.String(spec.HealthCheck.Target),
					Interval:           aws.Int64(int64(spec.HealthCheck.Interval.Seconds())),
					Timeout:            aws.Int64(int64(spec.HealthCheck.Timeout.Seconds())),
					HealthyThreshold:   aws.Int64(spec.HealthCheck.HealthyThreshold),
					UnhealthyThreshold: aws.Int64(spec.HealthCheck.UnhealthyThreshold),
				},
			}); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.LoadBalancerNotFound); err != nil {
			return nil, errors.Wrapf(err, "failed to configure health check for classic load balancer: %v", spec)
		}
	}

	s.scope.Info("Created classic load balancer", "dns-name", *out.DNSName)

	res := spec.DeepCopy()
	res.DNSName = *out.DNSName
	return res, nil
}

func (s *Service) configureAttributes(name string, attributes infrav1.ClassicELBAttributes) error {
	attrs := &elb.ModifyLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(name),
		LoadBalancerAttributes: &elb.LoadBalancerAttributes{
			CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
				Enabled: aws.Bool(attributes.CrossZoneLoadBalancing),
			},
		},
	}

	if attributes.IdleTimeout > 0 {
		attrs.LoadBalancerAttributes.ConnectionSettings = &elb.ConnectionSettings{
			IdleTimeout: aws.Int64(int64(attributes.IdleTimeout.Seconds())),
		}
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.ELBClient.ModifyLoadBalancerAttributes(attrs); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.LoadBalancerNotFound); err != nil {
		return errors.Wrapf(err, "failed to configure attributes for classic load balancer: %v", name)
	}

	return nil
}

func (s *Service) configureLBAttributes(arn string, attributes map[string]*string) error {
	attrs := make([]*elbv2.LoadBalancerAttribute, 0)
	for k, v := range attributes {
		attrs = append(attrs, &elbv2.LoadBalancerAttribute{
			Key:   aws.String(k),
			Value: v,
		})
	}
	s.scope.Debug("adding attributes to load balancer", "attrs", attrs)
	modifyInput := &elbv2.ModifyLoadBalancerAttributesInput{
		Attributes:      attrs,
		LoadBalancerArn: aws.String(arn),
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.ELBV2Client.ModifyLoadBalancerAttributes(modifyInput); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.LoadBalancerNotFound); err != nil {
		return errors.Wrapf(err, "failed to configure attributes for load balancer: %v", arn)
	}
	return nil
}

func (s *Service) deleteClassicELB(name string) error {
	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	if _, err := s.ELBClient.DeleteLoadBalancer(input); err != nil {
		return err
	}

	s.scope.Info("Deleted AWS cloud provider load balancers")
	return nil
}

func (s *Service) deleteLB(arn string) error {
	// remove listeners and target groups
	// Order is important. ClassicELBListeners have to be deleted first.
	// However, we must first gather the groups because after the listeners are deleted the groups
	// are no longer associated with the LB, so we can't describe them afterwards.
	groups, err := s.ELBV2Client.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(arn),
	})
	if err != nil {
		return fmt.Errorf("failed to gather target groups for LB: %w", err)
	}
	listeners, err := s.ELBV2Client.DescribeListeners(&elbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(arn),
	})
	if err != nil {
		return fmt.Errorf("failed to gather listeners: %w", err)
	}
	for _, listener := range listeners.Listeners {
		s.scope.Debug("deleting listener", "arn", aws.StringValue(listener.ListenerArn))
		deleteListener := &elbv2.DeleteListenerInput{
			ListenerArn: listener.ListenerArn,
		}
		if _, err := s.ELBV2Client.DeleteListener(deleteListener); err != nil {
			return fmt.Errorf("failed to delete listener '%s': %w", aws.StringValue(listener.ListenerArn), err)
		}
	}
	s.scope.Info("Successfully deleted all associated ClassicELBListeners")

	for _, group := range groups.TargetGroups {
		s.scope.Debug("deleting target group", "name", aws.StringValue(group.TargetGroupName))
		deleteTargetGroup := &elbv2.DeleteTargetGroupInput{
			TargetGroupArn: group.TargetGroupArn,
		}
		if _, err := s.ELBV2Client.DeleteTargetGroup(deleteTargetGroup); err != nil {
			return fmt.Errorf("failed to delete target group '%s': %w", aws.StringValue(group.TargetGroupName), err)
		}
	}

	s.scope.Info("Successfully deleted all associated Target Groups")

	deleteLoadBalancerInput := &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	}

	if _, err := s.ELBV2Client.DeleteLoadBalancer(deleteLoadBalancerInput); err != nil {
		return err
	}

	s.scope.Info("Deleted AWS cloud provider load balancers")
	return nil
}

func (s *Service) listByTag(tag string) ([]string, error) {
	input := rgapi.GetResourcesInput{
		ResourceTypeFilters: aws.StringSlice([]string{elbResourceType}),
		TagFilters: []*rgapi.TagFilter{
			{
				Key:    aws.String(tag),
				Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
			},
		},
	}

	names := []string{}

	err := s.ResourceTaggingClient.GetResourcesPages(&input, func(r *rgapi.GetResourcesOutput, last bool) bool {
		for _, tagmapping := range r.ResourceTagMappingList {
			if tagmapping.ResourceARN != nil {
				parsedARN, err := arn.Parse(*tagmapping.ResourceARN)
				if err != nil {
					s.scope.Info("failed to parse ARN", "arn", *tagmapping.ResourceARN, "tag", tag)
					continue
				}
				if strings.Contains(parsedARN.Resource, "loadbalancer/net/") {
					s.scope.Info("ignoring nlb created by service, consider enabling garbage collection", "arn", *tagmapping.ResourceARN, "tag", tag)
					continue
				}
				if strings.Contains(parsedARN.Resource, "loadbalancer/app/") {
					s.scope.Info("ignoring alb created by service, consider enabling garbage collection", "arn", *tagmapping.ResourceARN, "tag", tag)
					continue
				}
				name := strings.ReplaceAll(parsedARN.Resource, "loadbalancer/", "")
				if name == "" {
					s.scope.Info("failed to parse ARN", "arn", *tagmapping.ResourceARN, "tag", tag)
					continue
				}
				names = append(names, name)
			}
		}
		return true
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedListELBsByTag", "Failed to list %s ELB by Tags: %v", s.scope.Name(), err)
		return nil, errors.Wrapf(err, "failed to list %s ELBs by tag group", s.scope.Name())
	}

	return names, nil
}

func (s *Service) filterByOwnedTag(tagKey string) ([]string, error) {
	var names []string
	err := s.ELBClient.DescribeLoadBalancersPages(&elb.DescribeLoadBalancersInput{}, func(r *elb.DescribeLoadBalancersOutput, last bool) bool {
		for _, lb := range r.LoadBalancerDescriptions {
			names = append(names, *lb.LoadBalancerName)
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return nil, nil
	}

	var ownedElbs []string
	lbChunks := chunkELBs(names)
	for _, chunk := range lbChunks {
		output, err := s.ELBClient.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: aws.StringSlice(chunk)})
		if err != nil {
			return nil, err
		}
		for _, tagDesc := range output.TagDescriptions {
			for _, tag := range tagDesc.Tags {
				if *tag.Key == tagKey && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
					ownedElbs = append(ownedElbs, *tagDesc.LoadBalancerName)
				}
			}
		}
	}

	return ownedElbs, nil
}

func (s *Service) listAWSCloudProviderOwnedELBs() ([]string, error) {
	// k8s.io/cluster/<name>, created by k/k cloud provider
	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())
	arns, err := s.listByTag(serviceTag)
	if err != nil {
		// retry by listing all ELBs as listByTag will fail in air-gapped environments
		arns, err = s.filterByOwnedTag(serviceTag)
		if err != nil {
			return nil, err
		}
	}

	return arns, nil
}

func (s *Service) describeClassicELB(name string) (*infrav1.LoadBalancer, error) {
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{name}),
	}

	out, err := s.ELBClient.DescribeLoadBalancers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				return nil, NewNotFound(fmt.Sprintf("no classic load balancer found with name: %q", name))
			case elb.ErrCodeDependencyThrottleException:
				return nil, errors.Wrap(err, "too many requests made to the ELB service")
			default:
				return nil, errors.Wrap(err, "unexpected aws error")
			}
		} else {
			return nil, errors.Wrapf(err, "failed to describe classic load balancer: %s", name)
		}
	}

	if out != nil && len(out.LoadBalancerDescriptions) == 0 {
		return nil, NewNotFound(fmt.Sprintf("no classic load balancer found with name %q", name))
	}

	if s.scope.VPC().ID != "" && s.scope.VPC().ID != *out.LoadBalancerDescriptions[0].VPCId {
		return nil, errors.Errorf(
			"ELB names must be unique within a region: %q ELB already exists in this region in VPC %q",
			name, *out.LoadBalancerDescriptions[0].VPCId)
	}

	if s.scope.ControlPlaneLoadBalancer() != nil &&
		s.scope.ControlPlaneLoadBalancer().Scheme != nil &&
		string(*s.scope.ControlPlaneLoadBalancer().Scheme) != aws.StringValue(out.LoadBalancerDescriptions[0].Scheme) {
		return nil, errors.Errorf(
			"ELB names must be unique within a region: %q ELB already exists in this region with a different scheme %q",
			name, *out.LoadBalancerDescriptions[0].Scheme)
	}

	outAtt, err := s.ELBClient.DescribeLoadBalancerAttributes(&elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe classic load balancer %q attributes", name)
	}

	tags, err := s.describeClassicELBTags(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe classic load balancer tags")
	}

	return fromSDKTypeToClassicELB(out.LoadBalancerDescriptions[0], outAtt.LoadBalancerAttributes, tags), nil
}

func (s *Service) describeClassicELBTags(name string) ([]*elb.Tag, error) {
	output, err := s.ELBClient.DescribeTags(&elb.DescribeTagsInput{
		LoadBalancerNames: []*string{aws.String(name)},
	})
	if err != nil {
		return nil, err
	}

	if len(output.TagDescriptions) == 0 {
		return nil, errors.Errorf("no tag information returned for load balancer %q", name)
	}

	return output.TagDescriptions[0].Tags, nil
}

func (s *Service) describeLBTags(arn string) ([]*elbv2.Tag, error) {
	output, err := s.ELBV2Client.DescribeTags(&elbv2.DescribeTagsInput{
		ResourceArns: []*string{aws.String(arn)},
	})
	if err != nil {
		return nil, err
	}

	if len(output.TagDescriptions) == 0 {
		return nil, errors.Errorf("no tag information returned for load balancer %q", arn)
	}

	return output.TagDescriptions[0].Tags, nil
}

func (s *Service) reconcileELBTags(lb *infrav1.LoadBalancer, desiredTags map[string]string) error {
	addTagsInput := &elb.AddTagsInput{
		LoadBalancerNames: []*string{aws.String(lb.Name)},
	}

	removeTagsInput := &elb.RemoveTagsInput{
		LoadBalancerNames: []*string{aws.String(lb.Name)},
	}

	currentTags := infrav1.Tags(lb.Tags)

	for k, v := range desiredTags {
		if val, ok := currentTags[k]; !ok || val != v {
			s.scope.Trace("adding tag to load balancer", "elb-name", lb.Name, "key", k, "value", v)
			addTagsInput.Tags = append(addTagsInput.Tags, &elb.Tag{Key: aws.String(k), Value: aws.String(v)})
		}
	}

	for k := range currentTags {
		if _, ok := desiredTags[k]; !ok {
			s.scope.Trace("removing tag from load balancer", "elb-name", lb.Name, "key", k)
			removeTagsInput.Tags = append(removeTagsInput.Tags, &elb.TagKeyOnly{Key: aws.String(k)})
		}
	}

	if len(addTagsInput.Tags) > 0 {
		if _, err := s.ELBClient.AddTags(addTagsInput); err != nil {
			return err
		}
	}

	if len(removeTagsInput.Tags) > 0 {
		if _, err := s.ELBClient.RemoveTags(removeTagsInput); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) reconcileV2LBTags(lb *infrav1.LoadBalancer, desiredTags map[string]string) error {
	addTagsInput := &elbv2.AddTagsInput{
		ResourceArns: []*string{aws.String(lb.ARN)},
	}

	removeTagsInput := &elbv2.RemoveTagsInput{
		ResourceArns: []*string{aws.String(lb.ARN)},
	}

	currentTags := infrav1.Tags(lb.Tags)

	for k, v := range desiredTags {
		if val, ok := currentTags[k]; !ok || val != v {
			s.scope.Trace("adding tag to load balancer", "elb-name", lb.Name, "key", k, "value", v)
			addTagsInput.Tags = append(addTagsInput.Tags, &elbv2.Tag{Key: aws.String(k), Value: aws.String(v)})
		}
	}

	for k := range currentTags {
		if _, ok := desiredTags[k]; !ok {
			s.scope.Trace("removing tag from load balancer", "elb-name", lb.Name, "key", k)
			removeTagsInput.TagKeys = append(removeTagsInput.TagKeys, aws.String(k))
		}
	}

	if len(addTagsInput.Tags) > 0 {
		if _, err := s.ELBV2Client.AddTags(addTagsInput); err != nil {
			return err
		}
	}

	if len(removeTagsInput.TagKeys) > 0 {
		if _, err := s.ELBV2Client.RemoveTags(removeTagsInput); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) getHealthCheckTarget() string {
	controlPlaneELB := s.scope.ControlPlaneLoadBalancer()
	protocol := &infrav1.ELBProtocolSSL
	if controlPlaneELB != nil && controlPlaneELB.HealthCheckProtocol != nil {
		protocol = controlPlaneELB.HealthCheckProtocol
		if protocol.String() == infrav1.ELBProtocolHTTP.String() || protocol.String() == infrav1.ELBProtocolHTTPS.String() {
			return fmt.Sprintf("%v:%d/readyz", protocol, infrav1.DefaultAPIServerPort)
		}
	}
	return fmt.Sprintf("%v:%d", protocol, infrav1.DefaultAPIServerPort)
}

func fromSDKTypeToClassicELB(v *elb.LoadBalancerDescription, attrs *elb.LoadBalancerAttributes, tags []*elb.Tag) *infrav1.LoadBalancer {
	res := &infrav1.LoadBalancer{
		Name:             aws.StringValue(v.LoadBalancerName),
		Scheme:           infrav1.ELBScheme(*v.Scheme),
		SubnetIDs:        aws.StringValueSlice(v.Subnets),
		SecurityGroupIDs: aws.StringValueSlice(v.SecurityGroups),
		DNSName:          aws.StringValue(v.DNSName),
		Tags:             converters.ELBTagsToMap(tags),
		LoadBalancerType: infrav1.LoadBalancerTypeClassic,
	}

	if attrs.ConnectionSettings != nil && attrs.ConnectionSettings.IdleTimeout != nil {
		res.ClassicElbAttributes.IdleTimeout = time.Duration(*attrs.ConnectionSettings.IdleTimeout) * time.Second
	}

	res.ClassicElbAttributes.CrossZoneLoadBalancing = aws.BoolValue(attrs.CrossZoneLoadBalancing.Enabled)

	return res
}

func fromSDKTypeToLB(v *elbv2.LoadBalancer, attrs []*elbv2.LoadBalancerAttribute, tags []*elbv2.Tag) *infrav1.LoadBalancer {
	subnetIds := make([]*string, len(v.AvailabilityZones))
	availabilityZones := make([]*string, len(v.AvailabilityZones))
	for i, az := range v.AvailabilityZones {
		subnetIds[i] = az.SubnetId
		availabilityZones[i] = az.ZoneName
	}
	res := &infrav1.LoadBalancer{
		ARN:       aws.StringValue(v.LoadBalancerArn),
		Name:      aws.StringValue(v.LoadBalancerName),
		Scheme:    infrav1.ELBScheme(aws.StringValue(v.Scheme)),
		SubnetIDs: aws.StringValueSlice(subnetIds),
		// SecurityGroupIDs: aws.StringValueSlice(v.SecurityGroups),
		AvailabilityZones: aws.StringValueSlice(availabilityZones),
		DNSName:           aws.StringValue(v.DNSName),
		Tags:              converters.V2TagsToMap(tags),
	}

	infraAttrs := make(map[string]*string, len(attrs))
	for _, a := range attrs {
		infraAttrs[*a.Key] = a.Value
	}
	res.ELBAttributes = infraAttrs

	return res
}

// chunkELBs is similar to chunkResources in package pkg/cloud/services/gc.
func chunkELBs(names []string) [][]string {
	var chunked [][]string
	for i := 0; i < len(names); i += maxELBsDescribeTagsRequest {
		end := i + maxELBsDescribeTagsRequest
		if end > len(names) {
			end = len(names)
		}
		chunked = append(chunked, names[i:end])
	}
	return chunked
}
