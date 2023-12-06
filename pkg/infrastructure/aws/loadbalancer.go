package aws

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	readyzPath  = "/readyz"
	healthzPath = "/healthz"
	apiPort     = 6443
	servicePort = 22623
)

type lbInputOptions struct {
	infraID          string
	vpcID            string
	isPrivateCluster bool
	tags             map[string]string
	privateSubnetIDs []string
	publicSubnetIDs  []string
}

type lbState struct {
	input           *lbInputOptions
	targetGroupArns sets.Set[string]
}

type lbOutput struct {
	internal struct {
		zoneID  string
		dnsName string
	}
	external struct {
		zoneID  string
		dnsName string
	}
	targetGroupArns []string
}

func createLoadBalancers(ctx context.Context, logger logrus.FieldLogger, elbClient elbv2iface.ELBV2API, input *lbInputOptions) (*lbOutput, error) {
	state := lbState{
		input:           input,
		targetGroupArns: sets.New[string](),
	}
	output := &lbOutput{}

	internalLB, err := state.ensureInternalLoadBalancer(ctx, logger, elbClient, input.privateSubnetIDs, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create internal load balancer: %w", err)
	}
	output.internal.zoneID = aws.StringValue(internalLB.CanonicalHostedZoneId)
	output.internal.dnsName = aws.StringValue(internalLB.DNSName)

	if input.isPrivateCluster {
		logger.Debugln("Skipping creation of public LB because of private cluster")
		output.targetGroupArns = sets.List(state.targetGroupArns)
		return output, nil
	}

	externalLB, err := state.ensureExternalLoadBalancer(ctx, logger, elbClient, input.publicSubnetIDs, input.tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create external load balancer: %w", err)
	}
	output.external.zoneID = aws.StringValue(externalLB.CanonicalHostedZoneId)
	output.external.dnsName = aws.StringValue(externalLB.DNSName)

	output.targetGroupArns = sets.List(state.targetGroupArns)

	return output, nil
}

func (o *lbState) ensureInternalLoadBalancer(ctx context.Context, logger logrus.FieldLogger, client elbv2iface.ELBV2API, subnets []string, tags map[string]string) (*elbv2.LoadBalancer, error) {
	lbName := fmt.Sprintf("%s-int", o.input.infraID)
	lb, err := ensureLoadBalancer(ctx, logger, client, lbName, subnets, false, tags)
	if err != nil {
		return nil, err
	}

	// Create internalA target group
	aTGName := fmt.Sprintf("%s-aint", o.input.infraID)
	aTG, err := ensureTargetGroup(ctx, logger, client, aTGName, o.input.vpcID, readyzPath, apiPort, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create internalA target group: %w", err)
	}
	o.targetGroupArns.Insert(aws.StringValue(aTG.TargetGroupArn))

	// Create internalA listener
	aListenerName := fmt.Sprintf("%s-aint", o.input.infraID)
	aListener, err := createListener(ctx, client, aListenerName, lb.LoadBalancerArn, aTG.TargetGroupArn, 6443, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create internalA listener: %w", err)
	}
	logger.WithField("arn", aws.StringValue(aListener.ListenerArn)).Infoln("Created internalA listener")

	// Create internalS target group
	sTGName := fmt.Sprintf("%s-sint", o.input.infraID)
	sTG, err := ensureTargetGroup(ctx, logger, client, sTGName, o.input.vpcID, healthzPath, servicePort, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create internalS target group: %w", err)
	}
	o.targetGroupArns.Insert(aws.StringValue(sTG.TargetGroupArn))

	// Create internalS listener
	sListenerName := fmt.Sprintf("%s-sint", o.input.infraID)
	sListener, err := createListener(ctx, client, sListenerName, lb.LoadBalancerArn, sTG.TargetGroupArn, servicePort, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create internalS listener: %w", err)
	}
	logger.WithField("arn", aws.StringValue(sListener.ListenerArn)).Infoln("Created internalS listener")

	return lb, nil
}

func (o *lbState) ensureExternalLoadBalancer(ctx context.Context, logger logrus.FieldLogger, client elbv2iface.ELBV2API, subnets []string, tags map[string]string) (*elbv2.LoadBalancer, error) {
	lbName := fmt.Sprintf("%s-ext", o.input.infraID)
	lb, err := ensureLoadBalancer(ctx, logger, client, lbName, subnets, true, tags)
	if err != nil {
		return nil, err
	}

	// Create target group
	tgName := fmt.Sprintf("%s-aext", o.input.infraID)
	tg, err := ensureTargetGroup(ctx, logger, client, tgName, o.input.vpcID, readyzPath, apiPort, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create external target group: %w", err)
	}
	o.targetGroupArns.Insert(aws.StringValue(tg.TargetGroupArn))

	listenerName := fmt.Sprintf("%s-aext", o.input.infraID)
	listener, err := createListener(ctx, client, listenerName, lb.LoadBalancerArn, tg.TargetGroupArn, apiPort, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create external listener: %w", err)
	}
	logger.WithField("arn", aws.StringValue(listener.ListenerArn)).Infoln("Created external listener")

	return lb, nil
}

func ensureLoadBalancer(ctx context.Context, logger logrus.FieldLogger, client elbv2iface.ELBV2API, lbName string, subnets []string, isPublic bool, tags map[string]string) (*elbv2.LoadBalancer, error) {
	l := logger.WithField("name", lbName)
	createdOrFoundMsg := "Found existing load balancer"
	lb, err := existingLoadBalancer(ctx, client, lbName)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		createdOrFoundMsg = "Created load balancer"
		lb, err = createLoadBalancer(ctx, client, lbName, subnets, isPublic, tags)
		if err != nil {
			return nil, err
		}
	}
	l.Infoln(createdOrFoundMsg)

	// enable cross zone
	attrIn := &elbv2.ModifyLoadBalancerAttributesInput{
		LoadBalancerArn: lb.LoadBalancerArn,
		Attributes: []*elbv2.LoadBalancerAttribute{
			{
				Key:   aws.String("load_balancing.cross_zone.enabled"),
				Value: aws.String("true"),
			},
		},
	}
	_, err = client.ModifyLoadBalancerAttributesWithContext(ctx, attrIn)
	if err != nil {
		return nil, fmt.Errorf("failed to enable cross_zone attribute: %w", err)
	}
	l.Infoln("Enabled load balancer cross zone attribute")

	return lb, nil
}

func createLoadBalancer(ctx context.Context, client elbv2iface.ELBV2API, lbName string, subnets []string, isPublic bool, tags map[string]string) (*elbv2.LoadBalancer, error) {
	scheme := "internal"
	if isPublic {
		scheme = "internet-facing"
	}

	res, err := client.CreateLoadBalancerWithContext(ctx, &elbv2.CreateLoadBalancerInput{
		CustomerOwnedIpv4Pool: nil,
		IpAddressType:         nil,
		Name:                  aws.String(lbName),
		Scheme:                aws.String(scheme),
		SecurityGroups:        nil,
		SubnetMappings:        nil,
		Subnets:               aws.StringSlice(subnets),
		Tags:                  elbTags(tags),
		Type:                  aws.String(elbv2.LoadBalancerTypeEnumNetwork),
	})
	if err != nil {
		return nil, err
	}

	return res.LoadBalancers[0], nil
}

func existingLoadBalancer(ctx context.Context, client elbv2iface.ELBV2API, lbName string) (*elbv2.LoadBalancer, error) {
	lbInput := &elbv2.DescribeLoadBalancersInput{
		Names: aws.StringSlice([]string{lbName}),
	}
	res, err := client.DescribeLoadBalancersWithContext(ctx, lbInput)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == elbv2.ErrCodeLoadBalancerNotFoundException {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("failed to list load balancers: %w", err)
	}
	for _, lb := range res.LoadBalancers {
		return lb, nil
	}

	return nil, errNotFound
}

func ensureTargetGroup(ctx context.Context, logger logrus.FieldLogger, client elbv2iface.ELBV2API, targetName string, vpcID string, healthCheckPath string, port int64, tags map[string]string) (*elbv2.TargetGroup, error) {
	l := logger.WithField("name", targetName)
	createdOrFoundMsg := "Found existing Target Group"
	tg, err := existingTargetGroup(ctx, client, targetName)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		createdOrFoundMsg = "Created Target Group"
		tg, err = createTargetGroup(ctx, client, targetName, vpcID, healthCheckPath, port, tags)
		if err != nil {
			return nil, err
		}
	}
	l.Infoln(createdOrFoundMsg)

	return tg, nil
}

func existingTargetGroup(ctx context.Context, client elbv2iface.ELBV2API, targetName string) (*elbv2.TargetGroup, error) {
	input := &elbv2.DescribeTargetGroupsInput{
		Names: aws.StringSlice([]string{targetName}),
	}
	res, err := client.DescribeTargetGroupsWithContext(ctx, input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == elbv2.ErrCodeTargetGroupNotFoundException {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("failed to list target groups: %w", err)
	}
	for _, tg := range res.TargetGroups {
		return tg, nil
	}

	return nil, errNotFound
}

func createTargetGroup(ctx context.Context, client elbv2iface.ELBV2API, targetName string, vpcID string, healthCheckPath string, port int64, tags map[string]string) (*elbv2.TargetGroup, error) {
	ttags := mergeTags(tags, map[string]string{
		"Name": targetName,
	})
	input := &elbv2.CreateTargetGroupInput{
		HealthCheckEnabled:         aws.Bool(true),
		HealthCheckPath:            aws.String(healthCheckPath),
		HealthCheckPort:            aws.String(strconv.FormatInt(port, 10)),
		HealthCheckProtocol:        aws.String("HTTPS"),
		HealthCheckIntervalSeconds: aws.Int64(10),
		HealthyThresholdCount:      aws.Int64(2),
		UnhealthyThresholdCount:    aws.Int64(2),
		Name:                       aws.String(targetName),
		Port:                       aws.Int64(port),
		Protocol:                   aws.String("TCP"),
		Tags:                       elbTags(ttags),
		TargetType:                 aws.String("ip"),
		VpcId:                      aws.String(vpcID),
	}
	res, err := client.CreateTargetGroupWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return res.TargetGroups[0], nil
}

func createListener(ctx context.Context, client elbv2iface.ELBV2API, listenerName string, lbARN *string, tgARN *string, port int64, tags map[string]string) (*elbv2.Listener, error) {
	ltags := mergeTags(tags, map[string]string{
		"Name": listenerName,
	})
	input := &elbv2.CreateListenerInput{
		LoadBalancerArn: lbARN,
		Protocol:        aws.String("TCP"),
		Port:            aws.Int64(port),
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: tgARN,
						},
					},
				},
			},
		},
		Tags: elbTags(ltags),
	}
	// This operation is idempotent, which means that it completes at most one
	// time. If you attempt to create multiple listeners with the same
	// settings, each call succeeds.
	// https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/#ELBV2.CreateListener
	res, err := client.CreateListenerWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return res.Listeners[0], nil
}

func elbTags(tags map[string]string) []*elbv2.Tag {
	etags := make([]*elbv2.Tag, 0, len(tags))
	for k, v := range tags {
		k, v := k, v
		etags = append(etags, &elbv2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return etags
}
