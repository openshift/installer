package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	awsInfra "github.com/openshift/installer/pkg/infrastructure/aws"
	"github.com/sirupsen/logrus"
)

const (
	readyzPath  = "/readyz"
	healthzPath = "/healthz"
	apiPort     = 6443
	servicePort = 22623
)

func createIntLB(client elbv2iface.ELBV2API, subnets []string, tags map[string]string, infraID, vpcID string) (*elbv2.LoadBalancer, *elbv2.TargetGroup, *elbv2.TargetGroup, error) {
	logger := logrus.StandardLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	lbName := fmt.Sprintf("%s-int", infraID)
	lb, err := awsInfra.EnsureLoadBalancer(ctx, logger, client, lbName, subnets, false, tags)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create api-int target group
	aintTGName := fmt.Sprintf("%s-aint", infraID)
	aintTG, err := awsInfra.EnsureTargetGroup(ctx, logger, client, aintTGName, vpcID, readyzPath, apiPort, tags)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create external target group: %w", err)
	}

	aintListenerName := fmt.Sprintf("%s-aint", infraID)
	aintListener, err := awsInfra.CreateListener(ctx, client, aintListenerName, lb.LoadBalancerArn, aintTG.TargetGroupArn, apiPort, tags)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create external listener: %w", err)
	}

	logger.WithField("arn", aws.StringValue(aintListener.ListenerArn)).Infoln("Created api-int listener")

	// Create machine-config server target group
	sintTGName := fmt.Sprintf("%s-sint", infraID)
	sintTG, err := awsInfra.EnsureTargetGroup(ctx, logger, client, sintTGName, vpcID, healthzPath, servicePort, tags)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create external target group: %w", err)
	}
	sintListenerName := fmt.Sprintf("%s-sint", infraID)
	sintListener, err := awsInfra.CreateListener(ctx, client, sintListenerName, lb.LoadBalancerArn, sintTG.TargetGroupArn, servicePort, tags)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create external listener: %w", err)
	}
	logger.WithField("arn", aws.StringValue(sintListener.ListenerArn)).Infoln("Created mcs listener")

	return lb, aintTG, sintTG, nil
}

func registerControlPlane(client elbv2iface.ELBV2API, ids []*string, tg *elbv2.TargetGroup) error {
	targets := []*elbv2.TargetDescription{}
	for _, id := range ids {
		targets = append(targets, &elbv2.TargetDescription{Id: id})
	}

	_, err := client.RegisterTargetsWithContext(context.TODO(), &elbv2.RegisterTargetsInput{
		TargetGroupArn: tg.TargetGroupArn,
		Targets:        targets,
	})
	if err != nil {
		return fmt.Errorf("failed to register target group (%s): %w", *tg.TargetGroupName, err)
	}
	return nil
}
