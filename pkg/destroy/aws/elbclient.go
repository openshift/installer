package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

//go:generate mockgen -source=./elbclient.go -destination=mock/elbclient_generated.go -package=mock

// ELBAPI represents the calls made to the AWS ELB API
type ELBAPI interface {
	DeleteLoadBalancer(ctx context.Context, name string) error
	DescribeLoadBalancersPages(ctx context.Context, fn func(*elb.DescribeLoadBalancersOutput, bool) bool) error
}

// ELBv2API represents the calls made to the AWS ELBv2 API
type ELBv2API interface {
	DeleteLoadBalancer(ctx context.Context, name string) error
	DeleteTargetGroup(ctx context.Context, arn string) error
	DescribeLoadBalancersPages(ctx context.Context, fn func(*elbv2.DescribeLoadBalancersOutput, bool) bool) error
}

// ELBClient makes calls to the AWS ELB API
type ELBClient struct {
	client *elb.ELB
}

// ELBv2Client makes calls to the AWS ELBv2 API
type ELBv2Client struct {
	client *elbv2.ELBV2
}

// NewELBClient initializes a elbv1 client
func NewELBClient(awsSession *session.Session) *ELBClient {
	return &ELBClient{client: elb.New(awsSession)}
}

// NewELBv2Client initializes a elbv2 client
func NewELBv2Client(awsSession *session.Session) *ELBv2Client {
	return &ELBv2Client{client: elbv2.New(awsSession)}
}

// DeleteLoadBalancer deletes a LB named `name`
func (c *ELBClient) DeleteLoadBalancer(ctx context.Context, name string) error {
	_, err := c.client.DeleteLoadBalancerWithContext(ctx, &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// DescribeLoadBalancersPages runs `fn` for each LB page
func (c *ELBClient) DescribeLoadBalancersPages(ctx context.Context, fn func(*elb.DescribeLoadBalancersOutput, bool) bool) error {
	return c.client.DescribeLoadBalancersPagesWithContext(ctx, &elb.DescribeLoadBalancersInput{}, fn)
}

// DeleteLoadBalancer deletes a LB named `name`
func (c *ELBv2Client) DeleteLoadBalancer(ctx context.Context, arn string) error {
	_, err := c.client.DeleteLoadBalancerWithContext(ctx, &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteTargetGroup deletes a Target Group with Arn `arn`
func (c *ELBv2Client) DeleteTargetGroup(ctx context.Context, arn string) error {
	_, err := c.client.DeleteTargetGroupWithContext(ctx, &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	})
	if err != nil {
		return err
	}

	return nil
}

// DescribeLoadBalancersPages runs `fn` for each LB page
func (c *ELBv2Client) DescribeLoadBalancersPages(ctx context.Context, fn func(*elbv2.DescribeLoadBalancersOutput, bool) bool) error {
	return c.client.DescribeLoadBalancersPagesWithContext(ctx, &elbv2.DescribeLoadBalancersInput{}, fn)
}
