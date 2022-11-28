package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/destroy/aws/mock"
)

type (
	describeOutputClassicCB func(*elb.DescribeLoadBalancersOutput, bool) bool
	describeOutputv2CB      func(*elbv2.DescribeLoadBalancersOutput, bool) bool
)

func TestDeleteLBClassicByVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	elbClient := mock.NewMockELBAPI(mockCtrl)

	cases := []struct {
		name     string
		errorMsg string
	}{
		{
			name:     "Delete LB Classic by VPC succeeds",
			errorMsg: "",
		},
		{
			name:     "Delete LB Classic by VPC ignores LBs with no VPC",
			errorMsg: "",
		},
		{
			name:     "Delete LB Classic by VPC ignores LBs of wrong VPC",
			errorMsg: "",
		},
		{
			name:     "Delete LB Classic by VPC fails when delete fails",
			errorMsg: "deleting classic load balancer .*",
		},
		{
			name:     "Delete LB Classic by VPC fails when listing fails",
			errorMsg: "some aws elb error",
		},
	}

	gomock.InOrder(
		elbClient.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputClassicCB) error {
				results := &elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbClient.EXPECT().DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputClassicCB) error {
				results := &elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{LoadBalancerName: aws.String("lb-no-vpc"), VPCId: nil},
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbClient.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputClassicCB) error {
				results := &elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-other-vpc"), VPCId: aws.String("vpc-other")},
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-other-vpc"), VPCId: aws.String("vpc-other")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbClient.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputClassicCB) error {
				results := &elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-delete-fails"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-name"), VPCId: aws.String("vpc-name")},
						{LoadBalancerName: aws.String("lb-delete-fails"), VPCId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbClient.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			Return(errors.New("some aws elb error listing")),
	)

	elbClient.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("lb-name")).
		Return(nil).
		AnyTimes()
	elbClient.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("lb-no-vpc")).
		Times(0) // Should never be called
	elbClient.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("lb-other-vpc")).
		Times(0) // Should never be called
	elbClient.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("lb-delete-fails")).
		Return(errors.New("some aws elb error deleting")).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteElasticLoadBalancerClassicByVPC(context.TODO(), elbClient, "vpc-name", nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteLBV2ByVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	elbv2Client := mock.NewMockELBv2API(mockCtrl)

	cases := []struct {
		name     string
		errorMsg string
	}{
		{
			name:     "Delete LB v2 by VPC succeeds",
			errorMsg: "",
		},
		{
			name:     "Delete LB v2 by VPC ignores LBs with no VPC",
			errorMsg: "",
		},
		{
			name:     "Delete LB v2 by VPC ignores LBs of wrong VPC",
			errorMsg: "",
		},
		{
			name:     "Delete LB v2 by VPC fails when delete fails",
			errorMsg: "deleting .*",
		},
		{
			name:     "Delete LB v2 by VPC fails when parsing ARN fails",
			errorMsg: "parse ARN for load balancer",
		},
		{
			name:     "Delete LB v2 by VPC fails when listing fails",
			errorMsg: "some aws elbv2 error",
		},
	}

	gomock.InOrder(
		elbv2Client.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputv2CB) error {
				results := &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbv2Client.EXPECT().DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputv2CB) error {
				results := &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-no-vpc"), VpcId: nil},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbv2Client.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputv2CB) error {
				results := &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-other-vpc"), VpcId: aws.String("vpc-other")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-other-vpc"), VpcId: aws.String("vpc-other")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbv2Client.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputv2CB) error {
				results := &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-delete-fails"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-name"), VpcId: aws.String("vpc-name")},
						{LoadBalancerArn: aws.String("arn:aws:elbv2:::load-balancer/lb-delete-fails"), VpcId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbv2Client.EXPECT().DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn describeOutputv2CB) error {
				results := &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{LoadBalancerArn: aws.String("lb-no-vpc"), VpcId: nil},
						{LoadBalancerArn: aws.String("lb-name"), VpcId: aws.String("vpc-name")},
					},
				}
				fn(results, true)
				return nil
			}),
		elbv2Client.
			EXPECT().
			DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).
			Return(errors.New("some aws elbv2 error listing")),
	)

	elbv2Client.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("arn:aws:elbv2:::load-balancer/lb-name")).
		Return(nil).
		AnyTimes()
	elbv2Client.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("arn:aws:elbv2:::load-balancer/lb-no-vpc")).
		Times(0) // Should never be called
	elbv2Client.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("arn:aws:elbv2:::load-balancer/lb-other-vpc")).
		Times(0) // Should never be called
	elbv2Client.
		EXPECT().
		DeleteLoadBalancer(gomock.Any(), gomock.Eq("arn:aws:elbv2:::load-balancer/lb-delete-fails")).
		Return(errors.New("some aws elbv2 error deleting")).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteElasticLoadBalancerV2ByVPC(context.TODO(), elbv2Client, "vpc-name", nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
