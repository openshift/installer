package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func deleteElasticLoadBalancing(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "loadbalancer":
		segments := strings.SplitN(id, "/", 2)
		if len(segments) == 1 {
			return deleteElasticLoadBalancerClassic(ctx, elb.New(session), id, logger)
		} else if len(segments) != 2 {
			return errors.Errorf("cannot parse subresource %q into {subtype}/{id}", id)
		}
		subtype := segments[0]
		id = segments[1]
		switch subtype {
		case "net":
			return deleteElasticLoadBalancerV2(ctx, elbv2.New(session), arn, logger)
		default:
			return errors.Errorf("unrecognized elastic load balancing resource subtype %s", subtype)
		}
	case "targetgroup":
		return deleteElasticLoadBalancerTargetGroup(ctx, elbv2.New(session), arn, logger)
	case "listener":
		return deleteElasticLoadBalancerListener(ctx, elbv2.New(session), arn, logger)
	default:
		return errors.Errorf("unrecognized elastic load balancing resource type %s", resourceType)
	}
}

func deleteElasticLoadBalancerClassic(ctx context.Context, client *elb.ELB, name string, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancerWithContext(ctx, &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerClassicByVPC(ctx context.Context, client *elb.ELB, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPagesWithContext(
		ctx,
		&elb.DescribeLoadBalancersInput{},
		func(results *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			logger.Debugf("iterating over a page of %d v1 load balancers", len(results.LoadBalancerDescriptions))
			for _, lb := range results.LoadBalancerDescriptions {
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
					lastError = errors.Wrapf(err, "deleting classic load balancer %s", *lb.LoadBalancerName)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteElasticLoadBalancerTargetGroup(ctx context.Context, client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteTargetGroupWithContext(ctx, &elbv2.DeleteTargetGroupInput{
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
		if aerr.Code() == elbv2.ErrCodeListenerNotFoundException {
			return true
		}
	}
	return false
}

func deleteElasticLoadBalancerListener(ctx context.Context, client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteListenerWithContext(ctx, &elbv2.DeleteListenerInput{
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

func deleteElasticLoadBalancerV2(ctx context.Context, client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancerWithContext(ctx, &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerV2ByVPC(ctx context.Context, client *elbv2.ELBV2, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPagesWithContext(
		ctx,
		&elbv2.DescribeLoadBalancersInput{},
		func(results *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			logger.Debugf("iterating over a page of %d v2 load balancers", len(results.LoadBalancers))
			for _, lb := range results.LoadBalancers {
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
					lastError = errors.Wrap(err, "parse ARN for load balancer")
					continue
				}

				err = deleteElasticLoadBalancerV2(ctx, client, parsed, logger.WithField("load balancer", parsed.Resource))
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting %s", parsed.String())
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}
