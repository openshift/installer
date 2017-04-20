package cloudforms

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ec2Service interface {
	CreateVolume(*ec2.CreateVolumeInput) (*ec2.Volume, error)
	DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeKeyPairs(*ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error)
}

type Cluster struct {
	ClusterName      string
	ControllerDomain string
	Region           string
	StackBody        string
	StackURL         string
}

func (cb *Cluster) Deploy(sess *session.Session, tags []Tag) (*cloudformation.CreateStackOutput, error) {
	cfSvc := cloudformation.New(sess)

	// Add the specified tags to the stack.
	cfTags := []*cloudformation.Tag{
		{Key: aws.String("KubernetesCluster"), Value: aws.String(cb.ClusterName)},
	}
	for _, tag := range tags {
		t := &cloudformation.Tag{
			Key:   aws.String(tag.Key),
			Value: aws.String(tag.Value),
		}
		cfTags = append(cfTags, t)
	}

	// Create the stack.
	creq := &cloudformation.CreateStackInput{
		StackName:    aws.String(cb.ClusterName),
		OnFailure:    aws.String(cloudformation.OnFailureDoNothing),
		Capabilities: []*string{aws.String(cloudformation.CapabilityCapabilityIam)},
		TemplateURL:  aws.String(cb.StackURL),
		Tags:         cfTags,
	}
	cf, err := cfSvc.CreateStack(creq)

	return cf, maybeAwsErr(err)
}

type Status struct {
	Name         string
	ID           string
	StatusString string
	Events       []string
	Resources    []*cloudformation.StackResourceSummary
	Ready        bool
	Error        bool
}

func maybeAwsErr(e error) error {
	if e == nil {
		return nil
	}
	if awsErr, ok := e.(awserr.Error); ok {
		return fmt.Errorf("%s", awsErr.Message())
	}
	return e
}

// returns error if, for some reason, we can't query the stack status
func (cb *Cluster) Status(sess *session.Session) (*Status, error) {
	cfSvc := cloudformation.New(sess)

	req := cloudformation.DescribeStacksInput{
		StackName: aws.String(cb.ClusterName),
	}
	resp, err := cfSvc.DescribeStacks(&req)
	if err != nil {
		return nil, maybeAwsErr(err)
	}
	if len(resp.Stacks) == 0 {
		return nil, fmt.Errorf("Stack %s not found", cb.ClusterName)
	}

	stack := resp.Stacks[0]

	describeEventsInput := &cloudformation.DescribeStackEventsInput{
		StackName: stack.StackName,
	}
	stackEventsOutput, err := cfSvc.DescribeStackEvents(describeEventsInput)
	if err != nil {
		return nil, maybeAwsErr(err)
	}

	listResourcesInput := &cloudformation.ListStackResourcesInput{
		StackName: stack.StackName,
	}
	stackResourcesOutput, err := cfSvc.ListStackResources(listResourcesInput)
	if err != nil {
		return nil, maybeAwsErr(err)
	}

	stackStatus := &Status{
		ID:           aws.StringValue(stack.StackId),
		Name:         cb.ClusterName,
		StatusString: aws.StringValue(stack.StackStatus),
		Events:       stackEventErrMsgs(stackEventsOutput.StackEvents),
		Resources:    stackResourcesOutput.StackResourceSummaries,
	}

	if stackStatus.StatusString == cloudformation.ResourceStatusCreateComplete {
		stackStatus.Ready = true
	} else if stackStatus.StatusString == cloudformation.ResourceStatusCreateFailed {
		stackStatus.Error = true
	}

	return stackStatus, nil
}

func (cb *Cluster) Destroy(sess *session.Session) error {
	cfSvc := cloudformation.New(sess)
	dreq := &cloudformation.DeleteStackInput{
		StackName: aws.String(cb.ClusterName),
	}
	_, err := cfSvc.DeleteStack(dreq)
	return maybeAwsErr(err)
}

func stackEventErrMsgs(events []*cloudformation.StackEvent) []string {
	var errMsgs []string

	for _, event := range events {
		if aws.StringValue(event.ResourceStatus) == cloudformation.ResourceStatusCreateFailed {
			// Only show actual failures, not cancelled dependent resources.
			if aws.StringValue(event.ResourceStatusReason) != "Resource creation cancelled" {
				errMsgs = append(errMsgs,
					strings.TrimSpace(
						strings.Join([]string{
							aws.StringValue(event.ResourceStatus),
							aws.StringValue(event.ResourceType),
							aws.StringValue(event.LogicalResourceId),
							aws.StringValue(event.ResourceStatusReason),
						}, " ")))
			}
		}
	}

	return errMsgs
}
