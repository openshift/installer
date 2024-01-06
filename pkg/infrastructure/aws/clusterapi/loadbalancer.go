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

func createExtLB(client elbv2iface.ELBV2API, subnets []string, tags map[string]string, infraID, vpcID string) (*elbv2.LoadBalancer, *elbv2.TargetGroup, error) {
	logger := logrus.StandardLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	lbName := fmt.Sprintf("%s-ext", infraID)
	lb, err := awsInfra.EnsureLoadBalancer(ctx, logger, client, lbName, subnets, true, tags)
	if err != nil {
		return nil, nil, err
	}

	// Create target group
	tgName := fmt.Sprintf("%s-aext", infraID)
	tg, err := awsInfra.EnsureTargetGroup(ctx, logger, client, tgName, vpcID, readyzPath, apiPort, tags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create external target group: %w", err)
	}

	listenerName := fmt.Sprintf("%s-aext", infraID)
	listener, err := awsInfra.CreateListener(ctx, client, listenerName, lb.LoadBalancerArn, tg.TargetGroupArn, apiPort, tags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create external listener: %w", err)
	}
	logger.WithField("arn", aws.StringValue(listener.ListenerArn)).Infoln("Created external listener")

	return lb, tg, nil
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
