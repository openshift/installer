package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
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
			err := NewELBClient(session).DeleteLoadBalancer(ctx, id)
			if err == nil {
				logger.Info("Deleted")
			}
			return err
		} else if len(segments) != 2 {
			return errors.Errorf("cannot parse subresource %q into {subtype}/{id}", id)
		}
		subtype := segments[0]
		switch subtype {
		case "net":
			err := NewELBv2Client(session).DeleteLoadBalancer(ctx, arn.String())
			if err == nil {
				logger.Info("Deleted")
			}
			return err
		default:
			return errors.Errorf("unrecognized elastic load balancing resource subtype %s", subtype)
		}
	case "targetgroup":
		err := NewELBv2Client(session).DeleteTargetGroup(ctx, arn.String())
		if err == nil {
			logger.Info("Deleted")
		}
		return err
	default:
		return errors.Errorf("unrecognized elastic load balancing resource type %s", resourceType)
	}
}

func deleteElasticLoadBalancerClassicByVPC(ctx context.Context, client ELBAPI, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPages(ctx,
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

				err := client.DeleteLoadBalancer(ctx, *lb.LoadBalancerName)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting classic load balancer %s", *lb.LoadBalancerName)
				} else {
					lbLogger.Info("Deleted")
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

func deleteElasticLoadBalancerV2ByVPC(ctx context.Context, client ELBv2API, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPages(ctx,
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

				err = client.DeleteLoadBalancer(ctx, parsed.String())
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting %s", parsed.String())
				} else {
					logger.WithField("load balancer", parsed.Resource).Info("Deleted")
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

func deleteElasticLoadBalancingByVPC(ctx context.Context, session *session.Session, vpc string, logger logrus.FieldLogger) error {
	elbClient := NewELBClient(session)
	v1lbError := deleteElasticLoadBalancerClassicByVPC(ctx, elbClient, vpc, logger)

	elbv2Client := NewELBv2Client(session)
	v2lbError := deleteElasticLoadBalancerV2ByVPC(ctx, elbv2Client, vpc, logger)

	return utilerrors.NewAggregate([]error{v1lbError, v2lbError})
}
