package aws

import (
	"context"
	"fmt"
	"strings"

	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	elbcodes "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (o *ClusterUninstaller) deleteElasticLoadBalancing(ctx context.Context, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "loadbalancer":
		segments := strings.SplitN(id, "/", 2)
		if len(segments) == 1 {
			return deleteElasticLoadBalancerClassic(ctx, o.ElbBaseClient, id, logger)
		} else if len(segments) != 2 {
			return errors.Errorf("cannot parse subresource %q into {subtype}/{id}", id)
		}
		subtype := segments[0]
		id = segments[1]
		switch subtype {
		case "net":
			return deleteElasticLoadBalancerV2(ctx, o.Elbv2Client, arn, logger)
		default:
			return errors.Errorf("unrecognized elastic load balancing resource subtype %s", subtype)
		}
	case "targetgroup":
		return deleteElasticLoadBalancerTargetGroup(ctx, o.Elbv2Client, arn, logger)
	case "listener":
		return deleteElasticLoadBalancerListener(ctx, o.Elbv2Client, arn, logger)
	default:
		return errors.Errorf("unrecognized elastic load balancing resource type %s", resourceType)
	}
}

func deleteElasticLoadBalancerClassic(ctx context.Context, client *elb.Client, name string, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancer(ctx, &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerClassicByVPC(ctx context.Context, client *elb.Client, vpc string, logger logrus.FieldLogger) error {
	var lastError error

	paginator := elb.NewDescribeLoadBalancersPaginator(client, &elb.DescribeLoadBalancersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list load balancers: %w", err)
		}

		logger.Debugf("iterating over a page of %d v1 load balancers", len(page.LoadBalancerDescriptions))
		for _, lb := range page.LoadBalancerDescriptions {
			lbLogger := logger.WithField("classic load balancer", *lb.LoadBalancerName)

			if lb.VPCId == nil {
				lbLogger.Warn("classic load balancer does not have a VPC ID so could not determine whether it should be deleted")
				continue
			}

			if *lb.VPCId != vpc {
				continue
			}

			err := deleteElasticLoadBalancerClassic(ctx, client, *lb.LoadBalancerName, lbLogger)
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting classic load balancer %s: %w", *lb.LoadBalancerName, err)
			}
		}
	}

	return lastError
}

func deleteElasticLoadBalancerTargetGroup(ctx context.Context, client *elbv2.Client, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteTargetGroup(ctx, &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func isListenerNotFound(err interface{}) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == elbcodes.ErrCodeListenerNotFoundException {
			return true
		}
	}
	return false
}

func deleteElasticLoadBalancerListener(ctx context.Context, client *elbv2.Client, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteListener(ctx, &elbv2.DeleteListenerInput{
		ListenerArn: aws.String(arn.String()),
	})
	if err != nil {
		if isListenerNotFound(err) {
			logger.Info("Not found or already deleted")
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerV2(ctx context.Context, client *elbv2.Client, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancer(ctx, &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerV2ByVPC(ctx context.Context, client *elbv2.Client, vpc string, logger logrus.FieldLogger) error {
	var lastError error

	paginator := elbv2.NewDescribeLoadBalancersPaginator(client, &elbv2.DescribeLoadBalancersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list load balancers: %w", err)
		}

		logger.Debugf("iterating over a page of %d v2 load balancers", len(page.LoadBalancers))
		for _, lb := range page.LoadBalancers {
			if lb.VpcId == nil {
				logger.WithField("load balancer", *lb.LoadBalancerArn).Warn("load balancer does not have a VPC ID so could not determine whether it should be deleted")
				continue
			}

			if *lb.VpcId != vpc {
				continue
			}

			parsed, err := arn.Parse(*lb.LoadBalancerArn)
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("parse ARN for load balancer: %w", err)
				continue
			}

			err = deleteElasticLoadBalancerV2(ctx, client, parsed, logger.WithField("load balancer", parsed.Resource))
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting %s: %w", parsed.String(), err)
			}
		}
	}

	return lastError
}
