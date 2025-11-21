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
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	rgapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	rgapitypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/utils/ptr"

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

// apiServerTargetGroupPrefix is the target group name prefix used when creating a target group for the API server
// listener.
const apiServerTargetGroupPrefix = "apiserver-target-"

// additionalTargetGroupPrefix is the target group name prefix used when creating target groups for additional
// listeners.
const additionalTargetGroupPrefix = "additional-listener-"

// cantAttachSGToNLBRegions is a set of regions that do not support Security Groups in NLBs.
var cantAttachSGToNLBRegions = sets.New("us-iso-east-1", "us-iso-west-1", "us-isob-east-1")

type lbReconciler func() error

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers(ctx context.Context) error {
	s.scope.Debug("Reconciling load balancers")

	var errs []error
	var lbReconcilers []lbReconciler

	// The following splits load balancer reconciliation into 2 phases:
	// 1. Get or create the load balancer
	// 2. Reconcile the load balancer
	// We ensure that we only wait for the load balancer to become available in
	// the reconcile phase. This is useful when creating multiple load
	// balancers, as they can take several minutes to become available.

	for _, lbSpec := range s.scope.ControlPlaneLoadBalancers() {
		if lbSpec == nil {
			continue
		}
		switch lbSpec.LoadBalancerType {
		case infrav1.LoadBalancerTypeClassic:
			reconciler, err := s.getOrCreateClassicLoadBalancer(ctx)
			if err != nil {
				errs = append(errs, err)
			} else {
				lbReconcilers = append(lbReconcilers, reconciler)
			}
		case infrav1.LoadBalancerTypeNLB, infrav1.LoadBalancerTypeALB, infrav1.LoadBalancerTypeELB:
			reconciler, err := s.getOrCreateV2LB(ctx, lbSpec)
			if err != nil {
				errs = append(errs, err)
			} else {
				lbReconcilers = append(lbReconcilers, reconciler)
			}
		default:
			errs = append(errs, fmt.Errorf("unknown or unsupported load balancer type on primary load balancer: %s", lbSpec.LoadBalancerType))
		}
	}

	// Reconcile all load balancers
	for _, reconciler := range lbReconcilers {
		if err := reconciler(); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

// getOrCreateV2LB gets an existing load balancer, or creates a new one if it does not exist.
// It also takes care of generating unique names across namespaces by appending the namespace to the name.
// It returns a function that reconciles the load balancer.
func (s *Service) getOrCreateV2LB(ctx context.Context, lbSpec *infrav1.AWSLoadBalancerSpec) (lbReconciler, error) {
	name, err := LBName(s.scope, lbSpec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get control plane load balancer name")
	}

	// Get default api server spec.
	desiredLB, err := s.getAPIServerLBSpec(ctx, name, lbSpec)
	if err != nil {
		return nil, err
	}
	lb, err := s.describeLB(ctx, name, lbSpec)
	switch {
	case IsNotFound(err) && s.scope.ControlPlaneEndpoint().IsValid():
		// if elb is not found and owner cluster ControlPlaneEndpoint is already populated, then we should not recreate the elb.
		return nil, errors.Wrapf(err, "no loadbalancer exists for the AWSCluster %s, the cluster has become unrecoverable and should be deleted manually", s.scope.InfraClusterName())
	case IsNotFound(err):
		lb, err = s.createLB(ctx, desiredLB, lbSpec)
		if err != nil {
			s.scope.Error(err, "failed to create LB")
			return nil, err
		}

		s.scope.Debug("Created new network load balancer for apiserver", "api-server-lb-name", lb.Name)
	case err != nil:
		// Failed to describe the classic ELB
		return nil, err
	}

	return func() error {
		return s.reconcileV2LB(ctx, lb, desiredLB, lbSpec)
	}, nil
}

func (s *Service) reconcileV2LB(ctx context.Context, lb *infrav1.LoadBalancer, desiredLB *infrav1.LoadBalancer, lbSpec *infrav1.AWSLoadBalancerSpec) error {
	wReq := &elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{lb.ARN},
	}
	s.scope.Debug("Waiting for LB to become active", "api-server-lb-name", lb.Name)
	waitStart := time.Now()
	if err := s.ELBV2Client.WaitUntilLoadBalancerAvailable(ctx, wReq, s.scope.MaxWaitDuration()); err != nil {
		s.scope.Error(err, "failed to wait for LB to become available", "time", time.Since(waitStart))
		return err
	}
	s.scope.Debug("LB reports active state", "api-server-lb-name", lb.Name, "time", time.Since(waitStart))

	// set up the type for later processing
	lb.LoadBalancerType = lbSpec.LoadBalancerType
	if lb.IsManaged(s.scope.Name()) {
		// Reconcile the target groups and listeners from the spec and the ones currently attached to the load balancer.
		// Pass in the ARN that AWS gave us, as well as the rest of the desired specification.
		_, _, err := s.reconcileTargetGroupsAndListeners(ctx, lb.ARN, desiredLB, lbSpec)
		if err != nil {
			return errors.Wrapf(err, "failed to create target groups/listeners for load balancer %q", lb.Name)
		}

		if !cmp.Equal(desiredLB.ELBAttributes, lb.ELBAttributes) {
			if err := s.configureLBAttributes(ctx, lb.ARN, desiredLB.ELBAttributes); err != nil {
				return err
			}
		}

		if err := s.reconcileV2LBTags(ctx, lb, desiredLB.Tags); err != nil {
			return errors.Wrapf(err, "failed to reconcile tags for apiserver load balancer %q", lb.Name)
		}

		// Reconcile the subnets and availability zones from the desiredLB
		// and the ones currently attached to the load balancer.
		if len(lb.SubnetIDs) != len(desiredLB.SubnetIDs) {
			_, err := s.ELBV2Client.SetSubnets(ctx, &elbv2.SetSubnetsInput{
				LoadBalancerArn: &lb.ARN,
				Subnets:         desiredLB.SubnetIDs,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to set subnets for apiserver load balancer '%s'", lb.Name)
			}
		}
		if len(lb.AvailabilityZones) != len(desiredLB.AvailabilityZones) {
			lb.AvailabilityZones = desiredLB.AvailabilityZones
		}

		// Reconcile the security groups from the desiredLB and the ones currently attached to the load balancer
		if shouldReconcileSGs(s.scope, lb, desiredLB.SecurityGroupIDs) {
			_, err := s.ELBV2Client.SetSecurityGroups(ctx, &elbv2.SetSecurityGroupsInput{
				LoadBalancerArn: &lb.ARN,
				SecurityGroups:  desiredLB.SecurityGroupIDs,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to apply security groups to load balancer %q", lb.Name)
			}
		}
	} else {
		s.scope.Trace("Unmanaged control plane load balancer, skipping load balancer configuration", "api-server-elb", lb)
	}

	if s.scope.ControlPlaneLoadBalancers()[1] != nil && lb.Name == *s.scope.ControlPlaneLoadBalancers()[1].Name {
		lb.DeepCopyInto(&s.scope.Network().SecondaryAPIServerELB)
	} else {
		lb.DeepCopyInto(&s.scope.Network().APIServerELB)
	}

	return nil
}

// getAPITargetGroupHealthCheck creates the health check for the Kube apiserver target group,
// limiting the customization for the health check probe counters (skipping standarized/reserved
// fields: Protocol, Port or Path). To customize the health check protocol, use HealthCheckProtocol instead.
func (s *Service) getAPITargetGroupHealthCheck(lbSpec *infrav1.AWSLoadBalancerSpec) *infrav1.TargetGroupHealthCheck {
	apiHealthCheckProtocol := infrav1.ELBProtocolTCP.String()
	if lbSpec != nil && lbSpec.HealthCheckProtocol != nil {
		s.scope.Trace("Found API health check protocol override in the Load Balancer spec, applying it to the API Target Group", "api-server-elb", lbSpec.HealthCheckProtocol.String())
		apiHealthCheckProtocol = lbSpec.HealthCheckProtocol.String()
	}
	apiHealthCheck := &infrav1.TargetGroupHealthCheck{
		Protocol:                aws.String(apiHealthCheckProtocol),
		Port:                    aws.String(infrav1.DefaultAPIServerPortString),
		Path:                    nil,
		IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
		TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
		ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
		UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
	}
	if apiHealthCheckProtocol == infrav1.ELBProtocolHTTP.String() || apiHealthCheckProtocol == infrav1.ELBProtocolHTTPS.String() {
		apiHealthCheck.Path = aws.String(infrav1.DefaultAPIServerHealthCheckPath)
	}

	if lbSpec != nil && lbSpec.HealthCheck != nil {
		s.scope.Trace("Found API health check override in the Load Balancer spec, applying it to the API Target Group", "api-server-elb", lbSpec.HealthCheck)
		if lbSpec.HealthCheck.IntervalSeconds != nil {
			apiHealthCheck.IntervalSeconds = lbSpec.HealthCheck.IntervalSeconds
		}
		if lbSpec.HealthCheck.TimeoutSeconds != nil {
			apiHealthCheck.TimeoutSeconds = lbSpec.HealthCheck.TimeoutSeconds
		}
		if lbSpec.HealthCheck.ThresholdCount != nil {
			apiHealthCheck.ThresholdCount = lbSpec.HealthCheck.ThresholdCount
		}
		if lbSpec.HealthCheck.UnhealthyThresholdCount != nil {
			apiHealthCheck.UnhealthyThresholdCount = lbSpec.HealthCheck.UnhealthyThresholdCount
		}
	}
	return apiHealthCheck
}

// getAdditionalTargetGroupHealthCheck creates the target group health check for additional listener.
// Additional listeners allows to set customized attributes for health check.
func (s *Service) getAdditionalTargetGroupHealthCheck(ln infrav1.AdditionalListenerSpec) *infrav1.TargetGroupHealthCheck {
	healthCheck := &infrav1.TargetGroupHealthCheck{
		Port:                    aws.String(fmt.Sprintf("%d", ln.Port)),
		Protocol:                aws.String(ln.Protocol.String()),
		Path:                    nil,
		IntervalSeconds:         aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
		TimeoutSeconds:          aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
		ThresholdCount:          aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
		UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
	}
	if ln.HealthCheck == nil {
		return healthCheck
	}
	if ln.HealthCheck.Protocol != nil {
		healthCheck.Protocol = aws.String(*ln.HealthCheck.Protocol)
	}
	if ln.HealthCheck.Port != nil {
		healthCheck.Port = aws.String(*ln.HealthCheck.Port)
	}
	if ln.HealthCheck.Path != nil {
		healthCheck.Path = aws.String(*ln.HealthCheck.Path)
	}
	if ln.HealthCheck.IntervalSeconds != nil {
		healthCheck.IntervalSeconds = aws.Int64(*ln.HealthCheck.IntervalSeconds)
	}
	if ln.HealthCheck.TimeoutSeconds != nil {
		healthCheck.TimeoutSeconds = aws.Int64(*ln.HealthCheck.TimeoutSeconds)
	}
	if ln.HealthCheck.ThresholdCount != nil {
		healthCheck.ThresholdCount = aws.Int64(*ln.HealthCheck.ThresholdCount)
	}
	if ln.HealthCheck.UnhealthyThresholdCount != nil {
		healthCheck.UnhealthyThresholdCount = aws.Int64(*ln.HealthCheck.UnhealthyThresholdCount)
	}

	return healthCheck
}

func (s *Service) getAPIServerLBSpec(ctx context.Context, elbName string, lbSpec *infrav1.AWSLoadBalancerSpec) (*infrav1.LoadBalancer, error) {
	var securityGroupIDs []string
	if lbSpec != nil {
		securityGroupIDs = append(securityGroupIDs, lbSpec.AdditionalSecurityGroups...)
		securityGroupIDs = append(securityGroupIDs, s.scope.SecurityGroups()[infrav1.SecurityGroupAPIServerLB].ID)
	}

	// Since we're no longer relying on s.scope.ControlPlaneLoadBalancerScheme to do the defaulting for us, do it here.
	scheme := infrav1.ELBSchemeInternetFacing
	if lbSpec != nil && lbSpec.Scheme != nil {
		scheme = *lbSpec.Scheme
	}

	// The default API health check is TCP, allowing customization to HTTP or HTTPS when HealthCheckProtocol is set.
	apiHealthCheck := s.getAPITargetGroupHealthCheck(lbSpec)
	res := &infrav1.LoadBalancer{
		Name:          elbName,
		Scheme:        scheme,
		ELBAttributes: make(map[string]*string),
		ELBListeners: []infrav1.Listener{
			{
				Protocol: infrav1.ELBProtocolTCP,
				Port:     infrav1.DefaultAPIServerPort,
				TargetGroup: infrav1.TargetGroupSpec{
					Name:        names.SimpleNameGenerator.GenerateName(apiServerTargetGroupPrefix),
					Port:        infrav1.DefaultAPIServerPort,
					Protocol:    infrav1.ELBProtocolTCP,
					VpcID:       s.scope.VPC().ID,
					HealthCheck: apiHealthCheck,
				},
			},
		},
		SecurityGroupIDs: securityGroupIDs,
	}

	if lbSpec != nil {
		for _, listener := range lbSpec.AdditionalListeners {
			lnHealthCheck := &infrav1.TargetGroupHealthCheck{
				Protocol: aws.String(string(listener.Protocol)),
				Port:     aws.String(strconv.FormatInt(listener.Port, 10)),
			}
			if listener.HealthCheck != nil {
				s.scope.Trace("Found health check override in the additional listener spec, applying it to the Target Group", listener.HealthCheck)
				lnHealthCheck = s.getAdditionalTargetGroupHealthCheck(listener)
			}
			res.ELBListeners = append(res.ELBListeners, infrav1.Listener{
				Protocol: listener.Protocol,
				Port:     listener.Port,
				TargetGroup: infrav1.TargetGroupSpec{
					Name:        names.SimpleNameGenerator.GenerateName(additionalTargetGroupPrefix),
					Port:        listener.Port,
					Protocol:    listener.Protocol,
					VpcID:       s.scope.VPC().ID,
					HealthCheck: lnHealthCheck,
				},
			})
		}
	}

	if lbSpec != nil && lbSpec.LoadBalancerType != infrav1.LoadBalancerTypeNLB {
		res.ELBAttributes[infrav1.LoadBalancerAttributeIdleTimeTimeoutSeconds] = aws.String(infrav1.LoadBalancerAttributeIdleTimeDefaultTimeoutSecondsInSeconds)
	}

	if lbSpec != nil {
		isCrossZoneLB := lbSpec.CrossZoneLoadBalancing
		res.ELBAttributes[infrav1.LoadBalancerAttributeEnableLoadBalancingCrossZone] = aws.String(strconv.FormatBool(isCrossZoneLB))
	}

	res.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(elbName),
		Role:        aws.String(infrav1.APIServerRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	})

	// If subnet IDs have been specified for this load balancer
	if lbSpec != nil && len(lbSpec.Subnets) > 0 {
		// This set of subnets may not match the subnets specified on the Cluster, so we may not have already discovered them
		// We need to call out to AWS to describe them just in case
		input := &ec2.DescribeSubnetsInput{
			SubnetIds: lbSpec.Subnets,
		}
		out, err := s.EC2Client.DescribeSubnets(ctx, input)
		if err != nil {
			return nil, err
		}
		for _, sn := range out.Subnets {
			res.AvailabilityZones = append(res.AvailabilityZones, *sn.AvailabilityZone)
			res.SubnetIDs = append(res.SubnetIDs, *sn.SubnetId)
		}
	} else {
		// The load balancer APIs require us to only attach one subnet for each AZ.
		subnets := s.scope.Subnets().FilterPrivate().FilterNonCni()

		// public-only setup has no private subnets
		if scheme == infrav1.ELBSchemeInternetFacing || len(subnets) == 0 {
			subnets = s.scope.Subnets().FilterPublic().FilterNonCni()
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

func (s *Service) createLB(ctx context.Context, spec *infrav1.LoadBalancer, lbSpec *infrav1.AWSLoadBalancerSpec) (*infrav1.LoadBalancer, error) {
	var t elbv2types.LoadBalancerTypeEnum

	switch lbSpec.LoadBalancerType {
	case infrav1.LoadBalancerTypeNLB:
		t = elbv2types.LoadBalancerTypeEnumNetwork
	case infrav1.LoadBalancerTypeALB:
		t = elbv2types.LoadBalancerTypeEnumApplication
	case infrav1.LoadBalancerTypeELB:
		t = elbv2types.LoadBalancerTypeEnumGateway
	}
	input := &elbv2.CreateLoadBalancerInput{
		Name:           aws.String(spec.Name),
		Subnets:        spec.SubnetIDs,
		Tags:           converters.MapToV2Tags(spec.Tags),
		Scheme:         SchemeToSDKScheme(spec.Scheme),
		SecurityGroups: spec.SecurityGroupIDs,
		Type:           t,
	}

	if s.scope.VPC().IsIPv6Enabled() {
		input.IpAddressType = elbv2types.IpAddressTypeDualstack
	}

	// TODO: remove when security groups on NLBs is supported in all regions.
	if cantAttachSGToNLBRegions.Has(s.scope.Region()) {
		input.SecurityGroups = nil
	}

	// Allocate custom addresses (Elastic IP) to internet-facing Load Balancers, when defined.
	// Custom, or BYO, Public IPv4 Pool need to be created prior install, and the Pool ID must be
	// set in the VpcSpec.ElasticIPPool.PublicIPv4Pool to allow Elastic IP be consumed from
	// public ip address of user-provided CIDR blocks.
	if spec.Scheme == infrav1.ELBSchemeInternetFacing {
		if err := s.allocatePublicIpv4AddressFromByoIPPool(input); err != nil {
			return nil, fmt.Errorf("failed to allocate addresses to load balancer: %w", err)
		}
	}

	// Subnets and SubnetMappings are mutually exclusive. SubnetMappings is set by users or when
	// BYO Public IPv4 Pool is set.
	if len(input.SubnetMappings) == 0 {
		input.Subnets = spec.SubnetIDs
	}

	out, err := s.ELBV2Client.CreateLoadBalancer(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create load balancer: %v", spec)
	}

	if len(out.LoadBalancers) == 0 {
		return nil, errors.New("no new network load balancer was created; the returned list is empty")
	}

	// Target Groups and listeners will be reconciled separately

	if out.LoadBalancers[0].DNSName == nil {
		return nil, fmt.Errorf("CreateLoadBalancer did not return a DNS name for %s", spec.Name)
	}
	dnsName := *out.LoadBalancers[0].DNSName
	if out.LoadBalancers[0].LoadBalancerArn == nil {
		return nil, fmt.Errorf("CreateLoadBalancer did not return an ARN for %s", spec.Name)
	}
	arn := *out.LoadBalancers[0].LoadBalancerArn

	s.scope.Info("Created network load balancer", "dns-name", dnsName)

	res := spec.DeepCopy()
	s.scope.Debug("applying load balancer DNS to result", "dns", dnsName)
	res.DNSName = dnsName
	res.ARN = arn
	return res, nil
}

func (s *Service) describeLB(ctx context.Context, name string, lbSpec *infrav1.AWSLoadBalancerSpec) (*infrav1.LoadBalancer, error) {
	input := &elbv2.DescribeLoadBalancersInput{
		Names: []string{name},
	}

	out, err := s.ELBV2Client.DescribeLoadBalancers(ctx, input)
	smithyErr := awserrors.ParseSmithyError(err)
	if smithyErr != nil {
		switch smithyErr.ErrorCode() {
		case (&elbtypes.AccessPointNotFoundException{}).ErrorCode():
			return nil, NewNotFound(fmt.Sprintf("no load balancer found with name: %q", name))
		case (&elbtypes.DependencyThrottleException{}).ErrorCode():
			return nil, errors.Wrap(err, "too many requests made to the ELB service")
		default:
			return nil, errors.Wrap(err, "unexpected aws error")
		}
	}

	if out != nil && len(out.LoadBalancers) == 0 {
		return nil, NewNotFound(fmt.Sprintf("no load balancer found with name %q", name))
	}

	// Direct usage of indices here is alright because the query to AWS is providing exactly one name,
	// and the name uniqueness constraints prevent us from getting more than one entry back.
	if s.scope.VPC().ID != "" && s.scope.VPC().ID != *out.LoadBalancers[0].VpcId {
		return nil, errors.Errorf(
			"Load balancer names must be unique within a region: %q load balancer already exists in this region in VPC %q",
			name, *out.LoadBalancers[0].VpcId)
	}

	if lbSpec != nil &&
		lbSpec.Scheme != nil &&
		SchemeToSDKScheme(*lbSpec.Scheme) != out.LoadBalancers[0].Scheme {
		return nil, errors.Errorf(
			"Load balancer names must be unique within a region: %q Load balancer already exists in this region with a different scheme %q",
			name, out.LoadBalancers[0].Scheme)
	}

	outAtt, err := s.ELBV2Client.DescribeLoadBalancerAttributes(ctx, &elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: out.LoadBalancers[0].LoadBalancerArn,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe load balancer %q attributes", name)
	}

	tags, err := s.describeLBTags(ctx, aws.ToString(out.LoadBalancers[0].LoadBalancerArn))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe load balancer tags")
	}

	return fromSDKTypeToLB(out.LoadBalancers[0], outAtt.Attributes, tags), nil
}

// getOrCreateClassicLoadBalancer gets an existing classic load balancer, or creates a new one if it does not exist.
// It also takes care of generating unique names across namespaces by appending the namespace to the name.
// It returns a function that reconciles the load balancer.
func (s *Service) getOrCreateClassicLoadBalancer(ctx context.Context) (lbReconciler, error) {
	// Generate a default control plane load balancer name. The load balancer name cannot be
	// generated by the defaulting webhook, because it is derived from the cluster name, and that
	// name is undefined at defaulting time when generateName is used.
	name, err := ELBName(s.scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get control plane load balancer name")
	}

	// Get default api server spec.
	spec, err := s.getAPIServerClassicELBSpec(ctx, name)
	if err != nil {
		return nil, err
	}

	apiELB, err := s.describeClassicELB(ctx, spec.Name)
	switch {
	case IsNotFound(err) && s.scope.ControlPlaneEndpoint().IsValid():
		// if elb is not found and owner cluster ControlPlaneEndpoint is already populated, then we should not recreate the elb.
		return nil, errors.Wrapf(err, "no loadbalancer exists for the AWSCluster %s, the cluster has become unrecoverable and should be deleted manually", s.scope.InfraClusterName())
	case IsNotFound(err):
		apiELB, err = s.createClassicELB(ctx, spec)
		if err != nil {
			return nil, err
		}
		s.scope.Debug("Created new classic load balancer for apiserver", "api-server-elb-name", apiELB.Name)
	case err != nil:
		// Failed to describe the classic ELB
		return nil, err
	}

	return func() error {
		return s.reconcileClassicLoadBalancer(ctx, apiELB, spec)
	}, nil
}

func (s *Service) reconcileClassicLoadBalancer(ctx context.Context, apiELB *infrav1.LoadBalancer, spec *infrav1.LoadBalancer) error {
	if apiELB.IsManaged(s.scope.Name()) {
		if !cmp.Equal(spec.ClassicElbAttributes, apiELB.ClassicElbAttributes) {
			err := s.configureAttributes(ctx, apiELB.Name, spec.ClassicElbAttributes)
			if err != nil {
				return err
			}
		}

		// BUG: note that describeClassicELB doesn't set HealthCheck in its output,
		// so we're configuring the health check on every reconcile whether it's
		// needed or not.
		if !cmp.Equal(spec.HealthCheck, apiELB.HealthCheck) {
			s.scope.Debug("Reconciling health check for apiserver load balancer", "health-check", spec.HealthCheck)
			err := s.configureHealthCheck(ctx, apiELB.Name, spec.HealthCheck)
			if err != nil {
				return err
			}
		}

		if err := s.reconcileELBTags(ctx, apiELB, spec.Tags); err != nil {
			return errors.Wrapf(err, "failed to reconcile tags for apiserver load balancer %q", apiELB.Name)
		}

		// Reconcile the subnets and availability zones from the spec
		// and the ones currently attached to the load balancer.
		if len(apiELB.SubnetIDs) != len(spec.SubnetIDs) {
			_, err := s.ELBClient.AttachLoadBalancerToSubnets(ctx, &elb.AttachLoadBalancerToSubnetsInput{
				LoadBalancerName: &apiELB.Name,
				Subnets:          spec.SubnetIDs,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to attach apiserver load balancer %q to subnets", apiELB.Name)
			}
		}

		// Reconcile the security groups from the spec and the ones currently attached to the load balancer
		if !sets.NewString(apiELB.SecurityGroupIDs...).Equal(sets.NewString(spec.SecurityGroupIDs...)) {
			_, err := s.ELBClient.ApplySecurityGroupsToLoadBalancer(ctx, &elb.ApplySecurityGroupsToLoadBalancerInput{
				LoadBalancerName: &apiELB.Name,
				SecurityGroups:   spec.SecurityGroupIDs,
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

func (s *Service) configureHealthCheck(ctx context.Context, name string, healthCheck *infrav1.ClassicELBHealthCheck) error {
	healthCheckInput := &elb.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String(name),
		HealthCheck: &elbtypes.HealthCheck{
			Target:             aws.String(healthCheck.Target),
			Interval:           aws.Int32(int32(healthCheck.Interval.Seconds())),
			Timeout:            aws.Int32(int32(healthCheck.Timeout.Seconds())),
			HealthyThreshold:   aws.Int32(int32(healthCheck.HealthyThreshold)),   //#nosec G115
			UnhealthyThreshold: aws.Int32(int32(healthCheck.UnhealthyThreshold)), //#nosec G115
		},
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.ELBClient.ConfigureHealthCheck(ctx, healthCheckInput); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.LoadBalancerNotFound); err != nil {
		return errors.Wrapf(err, "failed to configure health check for classic load balancer: %s", name)
	}
	return nil
}

func (s *Service) deleteAPIServerELB(ctx context.Context) error {
	s.scope.Debug("Deleting control plane load balancer")

	elbName, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}

	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	apiELB, err := s.describeClassicELB(ctx, elbName)
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
	if err := s.deleteClassicELB(ctx, elbName); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		_, err = s.describeClassicELB(ctx, elbName)
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
func (s *Service) deleteAWSCloudProviderELBs(ctx context.Context) error {
	s.scope.Debug("Deleting AWS cloud provider load balancers (created for LoadBalancer-type Services)")

	elbs, err := s.listAWSCloudProviderOwnedELBs(ctx)
	if err != nil {
		return err
	}

	for _, elb := range elbs {
		s.scope.Debug("Deleting AWS cloud provider load balancer", "arn", elb)
		if err := s.deleteClassicELB(ctx, elb); err != nil {
			return err
		}
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		elbs, err := s.listAWSCloudProviderOwnedELBs(ctx)
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
func (s *Service) DeleteLoadbalancers(ctx context.Context) error {
	s.scope.Debug("Deleting load balancers")

	if err := s.deleteAPIServerELB(ctx); err != nil {
		return errors.Wrap(err, "failed to delete control plane load balancer")
	}

	if err := s.deleteAWSCloudProviderELBs(ctx); err != nil {
		return errors.Wrap(err, "failed to delete AWS cloud provider load balancer(s)")
	}

	if err := s.deleteExistingNLBs(ctx); err != nil {
		return errors.Wrap(err, "failed to delete AWS cloud provider load balancer(s)")
	}

	return nil
}

func (s *Service) deleteExistingNLBs(ctx context.Context) error {
	errs := make([]error, 0)

	for _, lbSpec := range s.scope.ControlPlaneLoadBalancers() {
		if lbSpec == nil {
			continue
		}
		errs = append(errs, s.deleteExistingNLB(ctx, lbSpec))
	}

	return kerrors.NewAggregate(errs)
}

func (s *Service) deleteExistingNLB(ctx context.Context, lbSpec *infrav1.AWSLoadBalancerSpec) error {
	name, err := LBName(s.scope, lbSpec)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	lb, err := s.describeLB(ctx, name, lbSpec)
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
	if err := s.deleteLB(ctx, lb.ARN); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.LoadBalancerReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (done bool, err error) {
		_, err = s.describeLB(ctx, name, lbSpec)
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
func (s *Service) IsInstanceRegisteredWithAPIServerELB(ctx context.Context, i *infrav1.Instance) (bool, error) {
	name, err := ELBName(s.scope)
	if err != nil {
		return false, errors.Wrap(err, "failed to get control plane load balancer name")
	}

	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{name},
	}

	output, err := s.ELBClient.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return false, errors.Wrapf(err, "error describing ELB %q", name)
	}
	if len(output.LoadBalancerDescriptions) != 1 {
		return false, errors.Errorf("expected 1 ELB description for %q, got %d", name, len(output.LoadBalancerDescriptions))
	}

	for _, registeredInstance := range output.LoadBalancerDescriptions[0].Instances {
		if aws.ToString(registeredInstance.InstanceId) == i.ID {
			return true, nil
		}
	}

	return false, nil
}

// IsInstanceRegisteredWithAPIServerLB returns true if the instance is already registered with the APIServer LB.
func (s *Service) IsInstanceRegisteredWithAPIServerLB(ctx context.Context, i *infrav1.Instance, lb *infrav1.AWSLoadBalancerSpec) ([]string, bool, error) {
	name, err := LBName(s.scope, lb)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get control plane load balancer name")
	}

	input := &elbv2.DescribeLoadBalancersInput{
		Names: []string{name},
	}

	output, err := s.ELBV2Client.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error describing ELB %q", name)
	}
	if len(output.LoadBalancers) != 1 {
		return nil, false, errors.Errorf("expected 1 ELB description for %q, got %d", name, len(output.LoadBalancers))
	}

	describeTargetGroupInput := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: output.LoadBalancers[0].LoadBalancerArn,
	}

	targetGroups, err := s.ELBV2Client.DescribeTargetGroups(ctx, describeTargetGroupInput)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error describing ELB's target groups %q", name)
	}

	targetGroupARNs := []string{}
	for _, tg := range targetGroups.TargetGroups {
		healthInput := &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		}
		instanceHealth, err := s.ELBV2Client.DescribeTargetHealth(ctx, healthInput)
		if err != nil {
			return nil, false, errors.Wrapf(err, "error describing ELB's target groups health %q", name)
		}
		for _, id := range instanceHealth.TargetHealthDescriptions {
			if aws.ToString(id.Target.Id) == i.ID {
				targetGroupARNs = append(targetGroupARNs, aws.ToString(tg.TargetGroupArn))
			}
		}
	}
	if len(targetGroupARNs) > 0 {
		return targetGroupARNs, true, nil
	}

	return nil, false, nil
}

// RegisterInstanceWithAPIServerELB registers an instance with a classic ELB.
func (s *Service) RegisterInstanceWithAPIServerELB(ctx context.Context, i *infrav1.Instance) error {
	name, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}
	out, err := s.describeClassicELB(ctx, name)
	if err != nil {
		return err
	}

	// Validate that the subnets associated with the load balancer has the instance AZ.
	subnets := s.scope.Subnets()
	instanceSubnet := subnets.FindByID(i.SubnetID)
	if instanceSubnet == nil {
		return errors.Errorf("failed to attach load balancer subnets, could not find subnet %q description in AWSCluster", i.SubnetID)
	}
	instanceAZ := instanceSubnet.AvailabilityZone

	if s.scope.ControlPlaneLoadBalancer() != nil && len(s.scope.ControlPlaneLoadBalancer().Subnets) > 0 {
		subnets, err = s.getControlPlaneLoadBalancerSubnets(ctx)
		if err != nil {
			return err
		}
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
		Instances:        []elbtypes.Instance{{InstanceId: aws.String(i.ID)}},
		LoadBalancerName: aws.String(name),
	}

	_, err = s.ELBClient.RegisterInstancesWithLoadBalancer(ctx, input)
	return err
}

// RegisterInstanceWithAPIServerLB registers an instance with a LB.
func (s *Service) RegisterInstanceWithAPIServerLB(ctx context.Context, instance *infrav1.Instance, lbSpec *infrav1.AWSLoadBalancerSpec) error {
	name, err := LBName(s.scope, lbSpec)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}
	out, err := s.describeLB(ctx, name, lbSpec)
	if err != nil {
		return err
	}
	s.scope.Debug("found load balancer with name", "name", out.Name)
	describeTargetGroupInput := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(out.ARN),
	}

	targetGroups, err := s.ELBV2Client.DescribeTargetGroups(ctx, describeTargetGroupInput)
	if err != nil {
		return errors.Wrapf(err, "error describing ELB's target groups %q", name)
	}
	if len(targetGroups.TargetGroups) == 0 {
		return fmt.Errorf("no target groups found for load balancer with arn '%s'", out.ARN)
	}
	// Since TargetGroups and Listeners don't care, or are not aware, of subnets before registration, we ignore that check.
	// Also, registering with AZ is not supported using the an InstanceID.
	s.scope.Debug("found number of target groups", "target-groups", len(targetGroups.TargetGroups))
	for _, tg := range targetGroups.TargetGroups {
		input := &elbv2.RegisterTargetsInput{
			TargetGroupArn: tg.TargetGroupArn,
			Targets: []elbv2types.TargetDescription{
				{
					Id:   aws.String(instance.ID),
					Port: tg.Port,
				},
			},
		}
		if _, err = s.ELBV2Client.RegisterTargets(ctx, input); err != nil {
			return fmt.Errorf("failed to register instance with target group '%s': %w", *tg.TargetGroupName, err)
		}
	}

	return nil
}

// getControlPlaneLoadBalancerSubnets retrieves ControlPlaneLoadBalancer subnets information.
func (s *Service) getControlPlaneLoadBalancerSubnets(ctx context.Context) (infrav1.Subnets, error) {
	var subnets infrav1.Subnets

	input := &ec2.DescribeSubnetsInput{
		SubnetIds: s.scope.ControlPlaneLoadBalancer().Subnets,
	}
	res, err := s.EC2Client.DescribeSubnets(ctx, input)
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
func (s *Service) DeregisterInstanceFromAPIServerELB(ctx context.Context, i *infrav1.Instance) error {
	name, err := ELBName(s.scope)
	if err != nil {
		return errors.Wrap(err, "failed to get control plane load balancer name")
	}

	input := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances:        []elbtypes.Instance{{InstanceId: aws.String(i.ID)}},
		LoadBalancerName: aws.String(name),
	}

	_, err = s.ELBClient.DeregisterInstancesFromLoadBalancer(ctx, input)
	smithyErr := awserrors.ParseSmithyError(err)
	if smithyErr != nil {
		switch smithyErr.ErrorCode() {
		case (&elbtypes.AccessPointNotFoundException{}).ErrorCode(), (&elbtypes.InvalidEndPointException{}).ErrorCode():
			// Ignoring LoadBalancerNotFound and InvalidInstance when deregistering
			return nil
		default:
			return err
		}
	}
	return err
}

// DeregisterInstanceFromAPIServerLB de-registers an instance from a LB.
func (s *Service) DeregisterInstanceFromAPIServerLB(ctx context.Context, targetGroupArn string, i *infrav1.Instance) error {
	input := &elbv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String(targetGroupArn),
		Targets: []elbv2types.TargetDescription{
			{
				Id: aws.String(i.ID),
			},
		},
	}

	_, err := s.ELBV2Client.DeregisterTargets(ctx, input)
	smithyErr := awserrors.ParseSmithyError(err)
	if smithyErr != nil {
		switch smithyErr.ErrorCode() {
		case (&elbtypes.AccessPointNotFoundException{}).ErrorCode(), (&elbtypes.InvalidEndPointException{}).ErrorCode():
			// Ignoring LoadBalancerNotFound and InvalidInstance when deregistering
			return nil
		default:
			return err
		}
	}
	return err
}

// ELBName returns the user-defined API Server ELB name, or a generated default if the user has not defined the ELB
// name.
// This is only for the primary load balancer.
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

// LBName returns the user-defined API Server LB name, or a generated default if the user has not defined the LB
// name.
// This is used for both the primary and secondary load balancers.
func LBName(s scope.ELBScope, lbSpec *infrav1.AWSLoadBalancerSpec) (string, error) {
	if lbSpec != nil && lbSpec.Name != nil {
		return *lbSpec.Name, nil
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

func (s *Service) getAPIServerClassicELBSpec(ctx context.Context, elbName string) (*infrav1.LoadBalancer, error) {
	securityGroupIDs := []string{}
	controlPlaneLoadBalancer := s.scope.ControlPlaneLoadBalancer()
	if controlPlaneLoadBalancer != nil && len(controlPlaneLoadBalancer.AdditionalSecurityGroups) != 0 {
		securityGroupIDs = append(securityGroupIDs, controlPlaneLoadBalancer.AdditionalSecurityGroups...)
	}
	securityGroupIDs = append(securityGroupIDs, s.scope.SecurityGroups()[infrav1.SecurityGroupAPIServerLB].ID)

	scheme := infrav1.ELBSchemeInternetFacing
	if controlPlaneLoadBalancer != nil && controlPlaneLoadBalancer.Scheme != nil {
		scheme = *controlPlaneLoadBalancer.Scheme
	}

	res := &infrav1.LoadBalancer{
		Name:   elbName,
		Scheme: scheme,
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
			Interval:           infrav1.DefaultAPIServerHealthCheckIntervalSec * time.Second,
			Timeout:            infrav1.DefaultAPIServerHealthCheckTimeoutSec * time.Second,
			HealthyThreshold:   infrav1.DefaultAPIServerHealthThresholdCount,
			UnhealthyThreshold: infrav1.DefaultAPIServerUnhealthThresholdCount,
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
			SubnetIds: s.scope.ControlPlaneLoadBalancer().Subnets,
		}
		out, err := s.EC2Client.DescribeSubnets(ctx, input)
		if err != nil {
			return nil, err
		}
		for _, sn := range out.Subnets {
			res.AvailabilityZones = append(res.AvailabilityZones, *sn.AvailabilityZone)
			res.SubnetIDs = append(res.SubnetIDs, *sn.SubnetId)
		}
	} else {
		// The load balancer APIs require us to only attach one subnet for each AZ.
		subnets := s.scope.Subnets().FilterPrivate().FilterNonCni()

		// public-only setup has no private subnets
		if scheme == infrav1.ELBSchemeInternetFacing || len(subnets) == 0 {
			subnets = s.scope.Subnets().FilterPublic().FilterNonCni()
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

func (s *Service) createClassicELB(ctx context.Context, spec *infrav1.LoadBalancer) (*infrav1.LoadBalancer, error) {
	input := &elb.CreateLoadBalancerInput{
		LoadBalancerName: aws.String(spec.Name),
		Subnets:          spec.SubnetIDs,
		SecurityGroups:   spec.SecurityGroupIDs,
		Scheme:           aws.String(string(spec.Scheme)),
		Tags:             converters.MapToELBTags(spec.Tags),
	}

	for _, ln := range spec.ClassicELBListeners {
		input.Listeners = append(input.Listeners, elbtypes.Listener{
			Protocol:         aws.String(string(ln.Protocol)),
			LoadBalancerPort: int32(ln.Port), //#nosec G115
			InstanceProtocol: aws.String(string(ln.InstanceProtocol)),
			InstancePort:     aws.Int32(int32(ln.InstancePort)), //#nosec G115
		})
	}

	out, err := s.ELBClient.CreateLoadBalancer(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create classic load balancer: %v", spec)
	}

	s.scope.Info("Created classic load balancer", "dns-name", *out.DNSName)

	res := spec.DeepCopy()
	res.DNSName = *out.DNSName

	// We haven't configured any health check yet. Don't report it here so it
	// will be set later during reconciliation.
	res.HealthCheck = nil

	return res, nil
}

func (s *Service) configureAttributes(ctx context.Context, name string, attributes infrav1.ClassicELBAttributes) error {
	attrs := &elb.ModifyLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(name),
		LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
			CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
				Enabled: attributes.CrossZoneLoadBalancing,
			},
		},
	}

	if attributes.IdleTimeout > 0 {
		attrs.LoadBalancerAttributes.ConnectionSettings = &elbtypes.ConnectionSettings{
			IdleTimeout: aws.Int32(int32(attributes.IdleTimeout.Seconds())),
		}
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.ELBClient.ModifyLoadBalancerAttributes(ctx, attrs); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.LoadBalancerNotFound); err != nil {
		return errors.Wrapf(err, "failed to configure attributes for classic load balancer: %v", name)
	}

	return nil
}

func (s *Service) configureLBAttributes(ctx context.Context, arn string, attributes map[string]*string) error {
	attrs := make([]elbv2types.LoadBalancerAttribute, 0)
	for k, v := range attributes {
		attrs = append(attrs, elbv2types.LoadBalancerAttribute{
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
		if _, err := s.ELBV2Client.ModifyLoadBalancerAttributes(ctx, modifyInput); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.LoadBalancerNotFound); err != nil {
		return errors.Wrapf(err, "failed to configure attributes for load balancer: %v", arn)
	}
	return nil
}

func (s *Service) deleteClassicELB(ctx context.Context, name string) error {
	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	if _, err := s.ELBClient.DeleteLoadBalancer(ctx, input); err != nil {
		return err
	}

	s.scope.Info("Deleted AWS cloud provider load balancers")
	return nil
}

func (s *Service) deleteLB(ctx context.Context, arn string) error {
	// remove listeners and target groups
	// Order is important. ClassicELBListeners have to be deleted first.
	// However, we must first gather the groups because after the listeners are deleted the groups
	// are no longer associated with the LB, so we can't describe them afterwards.
	groups, err := s.ELBV2Client.DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(arn),
	})
	if err != nil {
		return fmt.Errorf("failed to gather target groups for LB: %w", err)
	}
	listeners, err := s.ELBV2Client.DescribeListeners(ctx, &elbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(arn),
	})
	if err != nil {
		return fmt.Errorf("failed to gather listeners: %w", err)
	}
	for _, listener := range listeners.Listeners {
		s.scope.Debug("deleting listener", "arn", aws.ToString(listener.ListenerArn))
		deleteListener := &elbv2.DeleteListenerInput{
			ListenerArn: listener.ListenerArn,
		}
		if _, err := s.ELBV2Client.DeleteListener(ctx, deleteListener); err != nil {
			return fmt.Errorf("failed to delete listener '%s': %w", aws.ToString(listener.ListenerArn), err)
		}
	}
	s.scope.Info("Successfully deleted all associated ClassicELBListeners")

	for _, group := range groups.TargetGroups {
		s.scope.Debug("deleting target group", "name", aws.ToString(group.TargetGroupName))
		deleteTargetGroup := &elbv2.DeleteTargetGroupInput{
			TargetGroupArn: group.TargetGroupArn,
		}
		if _, err := s.ELBV2Client.DeleteTargetGroup(ctx, deleteTargetGroup); err != nil {
			return fmt.Errorf("failed to delete target group '%s': %w", aws.ToString(group.TargetGroupName), err)
		}
	}

	s.scope.Info("Successfully deleted all associated Target Groups")

	deleteLoadBalancerInput := &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	}

	if _, err := s.ELBV2Client.DeleteLoadBalancer(ctx, deleteLoadBalancerInput); err != nil {
		return err
	}

	s.scope.Info("Deleted AWS cloud provider load balancers")
	return nil
}

func (s *Service) listByTag(ctx context.Context, tag string) ([]string, error) {
	input := rgapi.GetResourcesInput{
		ResourceTypeFilters: []string{elbResourceType},
		TagFilters: []rgapitypes.TagFilter{
			{
				Key:    aws.String(tag),
				Values: []string{string(infrav1.ResourceLifecycleOwned)},
			},
		},
	}

	names := []string{}

	err := s.ResourceTaggingClient.GetResourcesPages(ctx, &input, func(r *rgapi.GetResourcesOutput) {
		for _, tagmapping := range r.ResourceTagMappingList {
			if tagmapping.ResourceARN == nil {
				continue
			}
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
	})

	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedListELBsByTag", "Failed to list %s ELB by Tags: %v", s.scope.Name(), err)
		return nil, errors.Wrapf(err, "failed to list %s ELBs by tag group", s.scope.Name())
	}

	return names, nil
}

func (s *Service) filterByOwnedTag(ctx context.Context, tagKey string) ([]string, error) {
	var names []string
	err := s.ELBClient.DescribeLoadBalancersPages(ctx, &elb.DescribeLoadBalancersInput{}, func(r *elb.DescribeLoadBalancersOutput) {
		for _, lb := range r.LoadBalancerDescriptions {
			names = append(names, *lb.LoadBalancerName)
		}
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
		output, err := s.ELBClient.DescribeTags(ctx, &elb.DescribeTagsInput{
			LoadBalancerNames: chunk,
		})
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

func (s *Service) listAWSCloudProviderOwnedELBs(ctx context.Context) ([]string, error) {
	// k8s.io/cluster/<name>, created by k/k cloud provider
	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())
	arns, err := s.listByTag(ctx, serviceTag)
	if err != nil {
		// retry by listing all ELBs as listByTag will fail in air-gapped environments
		arns, err = s.filterByOwnedTag(ctx, serviceTag)
		if err != nil {
			return nil, err
		}
	}

	return arns, nil
}

func (s *Service) describeClassicELB(ctx context.Context, name string) (*infrav1.LoadBalancer, error) {
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{name},
	}

	out, err := s.ELBClient.DescribeLoadBalancers(ctx, input)
	smithyErr := awserrors.ParseSmithyError(err)
	if smithyErr != nil {
		switch smithyErr.ErrorCode() {
		case (&elbtypes.AccessPointNotFoundException{}).ErrorCode():
			return nil, NewNotFound(fmt.Sprintf("no classic load balancer found with name: %q", name))
		case (&elbtypes.DependencyThrottleException{}).ErrorCode():
			return nil, errors.Wrap(err, "too many requests made to the ELB service")
		default:
			return nil, errors.Wrap(err, "unexpected aws error")
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
		string(*s.scope.ControlPlaneLoadBalancer().Scheme) != aws.ToString(out.LoadBalancerDescriptions[0].Scheme) {
		return nil, errors.Errorf(
			"ELB names must be unique within a region: %q ELB already exists in this region with a different scheme %q",
			name, *out.LoadBalancerDescriptions[0].Scheme)
	}

	outAtt, err := s.ELBClient.DescribeLoadBalancerAttributes(ctx, &elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe classic load balancer %q attributes", name)
	}

	tags, err := s.describeClassicELBTags(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe classic load balancer tags")
	}

	return fromSDKTypeToClassicELB(&out.LoadBalancerDescriptions[0], outAtt.LoadBalancerAttributes, tags), nil
}

func (s *Service) describeClassicELBTags(ctx context.Context, name string) ([]elbtypes.Tag, error) {
	output, err := s.ELBClient.DescribeTags(ctx, &elb.DescribeTagsInput{
		LoadBalancerNames: []string{name},
	})
	if err != nil {
		return nil, err
	}

	if len(output.TagDescriptions) == 0 {
		return nil, errors.Errorf("no tag information returned for load balancer %q", name)
	}

	return output.TagDescriptions[0].Tags, nil
}

func (s *Service) describeLBTags(ctx context.Context, arn string) ([]elbv2types.Tag, error) {
	output, err := s.ELBV2Client.DescribeTags(ctx, &elbv2.DescribeTagsInput{
		ResourceArns: []string{arn},
	})
	if err != nil {
		return nil, err
	}

	if len(output.TagDescriptions) == 0 {
		return nil, errors.Errorf("no tag information returned for load balancer %q", arn)
	}

	return output.TagDescriptions[0].Tags, nil
}

func (s *Service) reconcileELBTags(ctx context.Context, lb *infrav1.LoadBalancer, desiredTags map[string]string) error {
	addTagsInput := &elb.AddTagsInput{
		LoadBalancerNames: []string{lb.Name},
	}

	removeTagsInput := &elb.RemoveTagsInput{
		LoadBalancerNames: []string{lb.Name},
	}

	currentTags := infrav1.Tags(lb.Tags)

	for k, v := range desiredTags {
		if val, ok := currentTags[k]; !ok || val != v {
			s.scope.Trace("adding tag to load balancer", "elb-name", lb.Name, "key", k, "value", v)
			addTagsInput.Tags = append(addTagsInput.Tags, elbtypes.Tag{Key: aws.String(k), Value: aws.String(v)})
		}
	}

	for k := range currentTags {
		if _, ok := desiredTags[k]; !ok {
			s.scope.Trace("removing tag from load balancer", "elb-name", lb.Name, "key", k)
			removeTagsInput.Tags = append(removeTagsInput.Tags, elbtypes.TagKeyOnly{Key: aws.String(k)})
		}
	}

	if len(addTagsInput.Tags) > 0 {
		if _, err := s.ELBClient.AddTags(ctx, addTagsInput); err != nil {
			return err
		}
	}

	if len(removeTagsInput.Tags) > 0 {
		if _, err := s.ELBClient.RemoveTags(ctx, removeTagsInput); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) reconcileV2LBTags(ctx context.Context, lb *infrav1.LoadBalancer, desiredTags map[string]string) error {
	addTagsInput := &elbv2.AddTagsInput{
		ResourceArns: []string{lb.ARN},
	}

	removeTagsInput := &elbv2.RemoveTagsInput{
		ResourceArns: []string{lb.ARN},
	}

	currentTags := infrav1.Tags(lb.Tags)

	for k, v := range desiredTags {
		if val, ok := currentTags[k]; !ok || val != v {
			s.scope.Trace("adding tag to load balancer", "elb-name", lb.Name, "key", k, "value", v)
			addTagsInput.Tags = append(addTagsInput.Tags, elbv2types.Tag{Key: aws.String(k), Value: aws.String(v)})
		}
	}

	for k := range currentTags {
		if _, ok := desiredTags[k]; !ok {
			s.scope.Trace("removing tag from load balancer", "elb-name", lb.Name, "key", k)
			removeTagsInput.TagKeys = append(removeTagsInput.TagKeys, k)
		}
	}

	if len(addTagsInput.Tags) > 0 {
		if _, err := s.ELBV2Client.AddTags(ctx, addTagsInput); err != nil {
			return err
		}
	}

	if len(removeTagsInput.TagKeys) > 0 {
		if _, err := s.ELBV2Client.RemoveTags(ctx, removeTagsInput); err != nil {
			return err
		}
	}

	return nil
}

// reconcileTargetGroupsAndListeners reconciles a Load Balancer's defined listeners with corresponding AWS Target Groups and Listeners.
// These are combined into a single function since they are tightly integrated.
func (s *Service) reconcileTargetGroupsAndListeners(ctx context.Context, lbARN string, spec *infrav1.LoadBalancer, lbSpec *infrav1.AWSLoadBalancerSpec) ([]*elbv2types.TargetGroup, []*elbv2types.Listener, error) {
	existingTargetGroups, err := s.ELBV2Client.DescribeTargetGroups(
		ctx,
		&elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: aws.String(lbARN),
		})
	if err != nil {
		s.scope.Error(err, "could not describe target groups for load balancer", "arn", lbARN)
		return nil, nil, err
	}

	existingListeners, err := s.ELBV2Client.DescribeListeners(
		ctx,
		&elbv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(lbARN),
		})
	if err != nil {
		s.scope.Error(err, "could not describe listeners for load balancer", "arn", lbARN)
	}

	createdTargetGroups := make([]*elbv2types.TargetGroup, 0, len(spec.ELBListeners))
	createdListeners := make([]*elbv2types.Listener, 0, len(spec.ELBListeners))

	// TODO(Skarlso): Add options to set up SSL.
	// https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/3899
	for _, ln := range spec.ELBListeners {
		var group *elbv2types.TargetGroup
		tgSpec := ln.TargetGroup
		for _, g := range existingTargetGroups.TargetGroups {
			if isSDKTargetGroupEqualToTargetGroup(&g, &tgSpec) {
				group = &g
				break
			}
		}
		// create the target group first
		if group == nil {
			group, err = s.createTargetGroup(ctx, ln, spec.Tags)
			if err != nil {
				return nil, nil, err
			}
			createdTargetGroups = append(createdTargetGroups, group)

			targetGroupAttributeInput := &elbv2.ModifyTargetGroupAttributesInput{TargetGroupArn: group.TargetGroupArn}

			if lbSpec.LoadBalancerType == infrav1.LoadBalancerTypeNLB {
				targetGroupAttributeInput.Attributes = append(targetGroupAttributeInput.Attributes,
					elbv2types.TargetGroupAttribute{
						Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
						Value: aws.String("false"),
					},
					elbv2types.TargetGroupAttribute{
						Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
						Value: aws.String("300"),
					},
				)
			}

			if !lbSpec.PreserveClientIP {
				targetGroupAttributeInput.Attributes = append(targetGroupAttributeInput.Attributes,
					elbv2types.TargetGroupAttribute{
						Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
						Value: aws.String("false"),
					},
				)
			}

			if len(targetGroupAttributeInput.Attributes) > 0 {
				s.scope.Debug("configuring target group attributes", "attributes", targetGroupAttributeInput)
				if _, err := s.ELBV2Client.ModifyTargetGroupAttributes(ctx, targetGroupAttributeInput); err != nil {
					return nil, nil, errors.Wrapf(err, "failed to modify target group attribute")
				}
			}
		}

		var listener *elbv2types.Listener
		for _, l := range existingListeners.Listeners {
			if len(l.DefaultActions) > 0 && *l.DefaultActions[0].TargetGroupArn == *group.TargetGroupArn {
				listener = &l
				break
			}
		}

		if listener == nil {
			listener, err = s.createListener(ctx, ln, group, lbARN, spec.Tags)
			if err != nil {
				return nil, nil, err
			}
			createdListeners = append(createdListeners, listener)
		}
	}

	return createdTargetGroups, createdListeners, nil
}

// createListener creates a single Listener.
func (s *Service) createListener(ctx context.Context, ln infrav1.Listener, group *elbv2types.TargetGroup, lbARN string, tags map[string]string) (*elbv2types.Listener, error) {
	listenerInput := &elbv2.CreateListenerInput{
		DefaultActions: []elbv2types.Action{
			{
				TargetGroupArn: group.TargetGroupArn,
				Type:           elbv2types.ActionTypeEnumForward,
			},
		},
		LoadBalancerArn: aws.String(lbARN),
		Port:            aws.Int32(int32(ln.Port)), //#nosec G115
		Protocol:        elbProtocolToSDKProtocol(ln.Protocol),
		Tags:            converters.MapToV2Tags(tags),
	}
	// Create ClassicELBListeners
	listener, err := s.ELBV2Client.CreateListener(ctx, listenerInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create listener")
	}
	if len(listener.Listeners) == 0 {
		return nil, errors.New("no listener was created; the returned list is empty")
	}
	if len(listener.Listeners) > 1 {
		return nil, errors.New("more than one listener created; expected only one")
	}
	return &listener.Listeners[0], nil
}

// createTargetGroup creates a single Target Group.
func (s *Service) createTargetGroup(ctx context.Context, ln infrav1.Listener, tags map[string]string) (*elbv2types.TargetGroup, error) {
	targetGroupInput := &elbv2.CreateTargetGroupInput{
		Name:                       aws.String(ln.TargetGroup.Name),
		Port:                       aws.Int32(int32(ln.TargetGroup.Port)), //#nosec G115
		Protocol:                   elbProtocolToSDKProtocol(ln.TargetGroup.Protocol),
		VpcId:                      aws.String(ln.TargetGroup.VpcID),
		Tags:                       converters.MapToV2Tags(tags),
		HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
		HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
		HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
		UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
	}
	if s.scope.VPC().IsIPv6Enabled() {
		targetGroupInput.IpAddressType = elbv2types.TargetGroupIpAddressTypeEnumIpv6
	}
	if ln.TargetGroup.HealthCheck != nil {
		targetGroupInput.HealthCheckEnabled = aws.Bool(true)

		if ln.TargetGroup.HealthCheck.Protocol != nil {
			targetGroupInput.HealthCheckProtocol = elbv2types.ProtocolEnum(strings.ToUpper(aws.ToString(ln.TargetGroup.HealthCheck.Protocol)))
		}
		targetGroupInput.HealthCheckPort = ln.TargetGroup.HealthCheck.Port
		if ln.TargetGroup.HealthCheck.Path != nil {
			targetGroupInput.HealthCheckPath = ln.TargetGroup.HealthCheck.Path
		}
		if ln.TargetGroup.HealthCheck.IntervalSeconds != nil {
			targetGroupInput.HealthCheckIntervalSeconds = aws.Int32(int32(*ln.TargetGroup.HealthCheck.IntervalSeconds)) //#nosec G115
		}
		if ln.TargetGroup.HealthCheck.TimeoutSeconds != nil {
			targetGroupInput.HealthCheckTimeoutSeconds = aws.Int32(int32(*ln.TargetGroup.HealthCheck.TimeoutSeconds)) //#nosec G115
		}
		if ln.TargetGroup.HealthCheck.ThresholdCount != nil {
			targetGroupInput.HealthyThresholdCount = aws.Int32(int32(*ln.TargetGroup.HealthCheck.ThresholdCount)) //#nosec G115
		}
		if ln.TargetGroup.HealthCheck.UnhealthyThresholdCount != nil {
			targetGroupInput.UnhealthyThresholdCount = aws.Int32(int32(*ln.TargetGroup.HealthCheck.UnhealthyThresholdCount)) //#nosec G115
		}
	}
	s.scope.Debug("creating target group", "group", targetGroupInput, "listener", ln)
	group, err := s.ELBV2Client.CreateTargetGroup(ctx, targetGroupInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create target group for load balancer")
	}
	if len(group.TargetGroups) == 0 {
		return nil, errors.New("no target group was created; the returned list is empty")
	}
	if len(group.TargetGroups) > 1 {
		return nil, errors.New("more than one target group created; expected only one")
	}
	return &group.TargetGroups[0], nil
}

func (s *Service) getHealthCheckTarget() string {
	controlPlaneELB := s.scope.ControlPlaneLoadBalancer()
	protocol := &infrav1.ELBProtocolTCP
	if controlPlaneELB != nil && controlPlaneELB.HealthCheckProtocol != nil {
		protocol = controlPlaneELB.HealthCheckProtocol
		if protocol.String() == infrav1.ELBProtocolHTTP.String() || protocol.String() == infrav1.ELBProtocolHTTPS.String() {
			return fmt.Sprintf("%v:%d%s", protocol, infrav1.DefaultAPIServerPort, infrav1.DefaultAPIServerHealthCheckPath)
		}
	}
	return fmt.Sprintf("%v:%d", protocol, infrav1.DefaultAPIServerPort)
}

func fromSDKTypeToClassicELB(v *elbtypes.LoadBalancerDescription, attrs *elbtypes.LoadBalancerAttributes, tags []elbtypes.Tag) *infrav1.LoadBalancer {
	res := &infrav1.LoadBalancer{
		Name:             aws.ToString(v.LoadBalancerName),
		Scheme:           infrav1.ELBScheme(*v.Scheme),
		SubnetIDs:        v.Subnets,
		SecurityGroupIDs: v.SecurityGroups,
		DNSName:          aws.ToString(v.DNSName),
		Tags:             converters.ELBTagsToMap(tags),
		LoadBalancerType: infrav1.LoadBalancerTypeClassic,
	}

	if attrs.ConnectionSettings != nil && attrs.ConnectionSettings.IdleTimeout != nil {
		res.ClassicElbAttributes.IdleTimeout = time.Duration(*attrs.ConnectionSettings.IdleTimeout) * time.Second
	}

	res.ClassicElbAttributes.CrossZoneLoadBalancing = attrs.CrossZoneLoadBalancing.Enabled

	return res
}

func fromSDKTypeToLB(v elbv2types.LoadBalancer, attrs []elbv2types.LoadBalancerAttribute, tags []elbv2types.Tag) *infrav1.LoadBalancer {
	subnetIDs := make([]string, len(v.AvailabilityZones))
	availabilityZones := make([]string, len(v.AvailabilityZones))
	for i, az := range v.AvailabilityZones {
		subnetIDs[i] = aws.ToString(az.SubnetId)
		availabilityZones[i] = aws.ToString(az.ZoneName)
	}
	res := &infrav1.LoadBalancer{
		ARN:               aws.ToString(v.LoadBalancerArn),
		Name:              aws.ToString(v.LoadBalancerName),
		Scheme:            infrav1.ELBScheme(v.Scheme),
		SubnetIDs:         subnetIDs,
		SecurityGroupIDs:  v.SecurityGroups,
		AvailabilityZones: availabilityZones,
		DNSName:           aws.ToString(v.DNSName),
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

func shouldReconcileSGs(scope scope.ELBScope, lb *infrav1.LoadBalancer, specSGs []string) bool {
	// Backwards compat: NetworkLoadBalancers were not always capable of having security groups attached.
	// Once created without a security group, the NLB can never have any added.
	// (https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-security-groups.html)
	if lb.LoadBalancerType == infrav1.LoadBalancerTypeNLB && len(lb.SecurityGroupIDs) == 0 {
		if cantAttachSGToNLBRegions.Has(scope.Region()) {
			scope.Info("Region doesn't support NLB security groups, cannot reconcile security groups.", "region", scope.Region(), "elb-name", lb.Name)
		} else {
			scope.Info("Pre-existing NLB without security groups, cannot reconcile security groups.", "elb-name", lb.Name)
		}
		return false
	}
	if !sets.NewString(lb.SecurityGroupIDs...).Equal(sets.NewString(specSGs...)) {
		return true
	}
	return true
}

// isSDKTargetGroupEqualToTargetGroup checks if a given AWS SDK target group matches a target group spec.
func isSDKTargetGroupEqualToTargetGroup(elbTG *elbv2types.TargetGroup, spec *infrav1.TargetGroupSpec) bool {
	// We can't check only the target group's name because it's randomly generated every time we get a spec
	// But CAPA-created target groups are guaranteed to have the "apiserver-target-" or "additional-listener-" prefix.
	switch {
	case strings.HasPrefix(*elbTG.TargetGroupName, apiServerTargetGroupPrefix):
		if !strings.HasPrefix(spec.Name, apiServerTargetGroupPrefix) {
			return false
		}
	case strings.HasPrefix(*elbTG.TargetGroupName, additionalTargetGroupPrefix):
		if !strings.HasPrefix(spec.Name, additionalTargetGroupPrefix) {
			return false
		}
	default:
		// Not created by CAPA
		return false
	}
	return int64(ptr.Deref(elbTG.Port, 0)) == spec.Port && strings.EqualFold(string(elbTG.Protocol), spec.Protocol.String())
}

// SchemeToSDKScheme converts infrav1.ELBScheme to elbv2types.LoadBalancerSchemeEnum.
func SchemeToSDKScheme(scheme infrav1.ELBScheme) elbv2types.LoadBalancerSchemeEnum {
	if scheme == infrav1.ELBSchemeInternetFacing {
		return elbv2types.LoadBalancerSchemeEnumInternetFacing
	}
	return elbv2types.LoadBalancerSchemeEnumInternal
}

func elbProtocolToSDKProtocol(protocol infrav1.ELBProtocol) elbv2types.ProtocolEnum {
	if protocol == infrav1.ELBProtocolSSL {
		return elbv2types.ProtocolEnumTls
	}
	return elbv2types.ProtocolEnum(protocol)
}

// WaitUntilLoadBalancerAvailable is a blocking function to wait until LoadBalancerV2 is Available.
func (c *ELBV2Client) WaitUntilLoadBalancerAvailable(ctx context.Context, input *elbv2.DescribeLoadBalancersInput, maxWait time.Duration) error {
	waiter := elbv2.NewLoadBalancerAvailableWaiter(c, func(o *elbv2.LoadBalancerAvailableWaiterOptions) {
		o.LogWaitAttempts = true
	})

	return waiter.Wait(ctx, input, maxWait)
}

// GetResourcesPages implementation of SDK V2.
func (c *ResourceGroupsTaggingAPIClient) GetResourcesPages(ctx context.Context, input *rgapi.GetResourcesInput, fn func(*rgapi.GetResourcesOutput)) error {
	paginator := rgapi.NewGetResourcesPaginator(c, input)
	for paginator.HasMorePages() {
		r, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		fn(r)
	}
	return nil
}

// DescribeLoadBalancersPages implementation of SDK V2.
func (c *ELBClient) DescribeLoadBalancersPages(ctx context.Context, input *elb.DescribeLoadBalancersInput, fn func(*elb.DescribeLoadBalancersOutput)) error {
	paginator := elb.NewDescribeLoadBalancersPaginator(c, input)
	for paginator.HasMorePages() {
		r, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		fn(r)
	}
	return nil
}

// DescribeLoadBalancersPages implementation of SDK V2.
func (c *ELBV2Client) DescribeLoadBalancersPages(ctx context.Context, input *elbv2.DescribeLoadBalancersInput, fn func(*elbv2.DescribeLoadBalancersOutput)) error {
	paginator := elbv2.NewDescribeLoadBalancersPaginator(c, input)
	for paginator.HasMorePages() {
		r, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		fn(r)
	}
	return nil
}
