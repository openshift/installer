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

package scope

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	awslogs "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/logs"
	awsmetrics "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/metrics"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

// NewASGClient creates a new ASG API client for a given session.
func NewASGClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) autoscalingiface.AutoScalingAPI {
	asgClient := autoscaling.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	asgClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	asgClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	asgClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return asgClient
}

// NewEC2Client creates a new EC2 API client for a given session.
func NewEC2Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ec2iface.EC2API {
	ec2Client := ec2.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	ec2Client.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	if session.ServiceLimiter(ec2.ServiceID) != nil {
		ec2Client.Handlers.Sign.PushFront(session.ServiceLimiter(ec2.ServiceID).LimitRequest)
	}
	ec2Client.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	if session.ServiceLimiter(ec2.ServiceID) != nil {
		ec2Client.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(ec2.ServiceID).ReviewResponse)
	}
	ec2Client.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return ec2Client
}

// NewELBClient creates a new ELB API client for a given session.
func NewELBClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) elbiface.ELBAPI {
	elbClient := elb.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	elbClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	elbClient.Handlers.Sign.PushFront(session.ServiceLimiter(elb.ServiceID).LimitRequest)
	elbClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	elbClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(elb.ServiceID).ReviewResponse)
	elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return elbClient
}

// NewELBv2Client creates a new ELB v2 API client for a given session.
func NewELBv2Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) elbv2iface.ELBV2API {
	elbClient := elbv2.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	elbClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	elbClient.Handlers.Sign.PushFront(session.ServiceLimiter(elbv2.ServiceID).LimitRequest)
	elbClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	elbClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(elbv2.ServiceID).ReviewResponse)
	elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return elbClient
}

// NewEventBridgeClient creates a new EventBridge API client for a given session.
func NewEventBridgeClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) eventbridgeiface.EventBridgeAPI {
	eventBridgeClient := eventbridge.New(session.Session())
	eventBridgeClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	eventBridgeClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	eventBridgeClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return eventBridgeClient
}

// NewSQSClient creates a new SQS API client for a given session.
func NewSQSClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) sqsiface.SQSAPI {
	SQSClient := sqs.New(session.Session())
	SQSClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	SQSClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	SQSClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return SQSClient
}

// NewGlobalSQSClient for creating a new SQS API client that isn't tied to a cluster.
func NewGlobalSQSClient(scopeUser cloud.ScopeUsage, session cloud.Session) sqsiface.SQSAPI {
	SQSClient := sqs.New(session.Session())
	SQSClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	SQSClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))

	return SQSClient
}

// NewResourgeTaggingClient creates a new Resource Tagging API client for a given session.
func NewResourgeTaggingClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI {
	resourceTagging := resourcegroupstaggingapi.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	resourceTagging.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	resourceTagging.Handlers.Sign.PushFront(session.ServiceLimiter(resourceTagging.ServiceID).LimitRequest)
	resourceTagging.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	resourceTagging.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(resourceTagging.ServiceID).ReviewResponse)
	resourceTagging.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return resourceTagging
}

// NewSecretsManagerClient creates a new Secrets API client for a given session..
func NewSecretsManagerClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) secretsmanageriface.SecretsManagerAPI {
	secretsClient := secretsmanager.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	secretsClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	secretsClient.Handlers.Sign.PushFront(session.ServiceLimiter(secretsClient.ServiceID).LimitRequest)
	secretsClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	secretsClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(secretsClient.ServiceID).ReviewResponse)
	secretsClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return secretsClient
}

// NewEKSClient creates a new EKS API client for a given session.
func NewEKSClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) eksiface.EKSAPI {
	eksClient := eks.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	eksClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	eksClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	eksClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return eksClient
}

// NewIAMClient creates a new IAM API client for a given session.
func NewIAMClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) iamiface.IAMAPI {
	iamClient := iam.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	iamClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	iamClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	iamClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return iamClient
}

// NewSTSClient creates a new STS API client for a given session.
func NewSTSClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) stsiface.STSAPI {
	stsClient := sts.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	stsClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	stsClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	stsClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return stsClient
}

// NewSSMClient creates a new Secrets API client for a given session.
func NewSSMClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ssmiface.SSMAPI {
	ssmClient := ssm.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	ssmClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	ssmClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	ssmClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return ssmClient
}

// NewS3Client creates a new S3 API client for a given session.
func NewS3Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) s3iface.S3API {
	s3Client := s3.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	s3Client.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	s3Client.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	s3Client.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return s3Client
}

func recordAWSPermissionsIssue(target runtime.Object) func(r *request.Request) {
	return func(r *request.Request) {
		if awsErr, ok := r.Error.(awserr.Error); ok {
			switch awsErr.Code() {
			case "AuthFailure", "UnauthorizedOperation", "NoCredentialProviders":
				record.Warnf(target, awsErr.Code(), "Operation %s failed with a credentials or permission issue", r.Operation.Name)
			}
		}
	}
}

func getUserAgentHandler() request.NamedHandler {
	return request.NamedHandler{
		Name: "capa/user-agent",
		Fn:   request.MakeAddToUserAgentHandler("aws.cluster.x-k8s.io", version.Get().String()),
	}
}

// AWSClients contains all the aws clients used by the scopes.
type AWSClients struct {
	ASG             autoscalingiface.AutoScalingAPI
	EC2             ec2iface.EC2API
	ELB             elbiface.ELBAPI
	SecretsManager  secretsmanageriface.SecretsManagerAPI
	ResourceTagging resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}
