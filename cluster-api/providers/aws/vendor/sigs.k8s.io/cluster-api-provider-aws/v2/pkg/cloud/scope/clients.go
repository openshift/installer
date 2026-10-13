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
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	rgapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/endpoints"
	awslogs "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/logs"
	awsmetrics "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/metrics"
	stsservice "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

// NewASGClient creates a new ASG API client for a given session.
func NewASGClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *autoscaling.Client {
	cfg := session.Session()

	autoscalingOpts := []func(*autoscaling.Options){
		func(o *autoscaling.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
		},
		autoscaling.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return autoscaling.NewFromConfig(cfg, autoscalingOpts...)
}

// NewEC2Client creates a new EC2 API client for a given session.
func NewEC2Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *ec2.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	ec2EndpointResolver := &endpoints.EC2EndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	ec2opts := []func(*ec2.Options){
		func(o *ec2.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = ec2EndpointResolver
		},
		ec2.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
			throttle.WithServiceLimiterMiddleware(session.ServiceLimiter(ec2.ServiceID)),
		),
	}

	return ec2.NewFromConfig(cfg, ec2opts...)
}

// NewELBClient creates a new ELB API client for a given session.
func NewELBClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *elb.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	endpointResolver := &endpoints.ELBEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	opts := []func(*elb.Options){
		func(o *elb.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = endpointResolver
		},
		elb.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
			throttle.WithServiceLimiterMiddleware(session.ServiceLimiter(elb.ServiceID)),
		),
	}

	return elb.NewFromConfig(cfg, opts...)
}

// NewELBv2Client creates a new ELB v2 API client for a given session.
func NewELBv2Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *elbv2.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	endpointResolver := &endpoints.ELBV2EndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	opts := []func(*elbv2.Options){
		func(o *elbv2.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = endpointResolver
		},
		elbv2.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
			throttle.WithServiceLimiterMiddleware(session.ServiceLimiter(elbv2.ServiceID)),
		),
	}

	return elbv2.NewFromConfig(cfg, opts...)
}

// NewEventBridgeClient creates a new EventBridge API client for a given session.
func NewEventBridgeClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) *eventbridge.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	endpointResolver := &endpoints.EventBridgeEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	opts := []func(*eventbridge.Options){
		func(o *eventbridge.Options) {
			o.EndpointResolverV2 = endpointResolver
		},
		eventbridge.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return eventbridge.NewFromConfig(cfg, opts...)
}

// NewSQSClient creates a new SQS API client for a given session.
func NewSQSClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) *sqs.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	endpointResolver := &endpoints.SQSEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	opts := []func(*sqs.Options){
		func(o *sqs.Options) {
			o.EndpointResolverV2 = endpointResolver
		},
		sqs.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return sqs.NewFromConfig(cfg, opts...)
}

// NewGlobalSQSClient for creating a new SQS API client that isn't tied to a cluster.
func NewGlobalSQSClient(scopeUser cloud.ScopeUsage, session cloud.Session) *sqs.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	endpointResolver := &endpoints.SQSEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	opts := []func(*sqs.Options){
		func(o *sqs.Options) {
			o.EndpointResolverV2 = endpointResolver
		},
		sqs.WithAPIOptions(
			awsmetrics.WithRequestMetricContextMiddleware(),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return sqs.NewFromConfig(cfg, opts...)
}

// NewResourgeTaggingClient creates a new Resource Tagging API client for a given session.
func NewResourgeTaggingClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *rgapi.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	endpointResolver := &endpoints.RGAPIEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	opts := []func(*rgapi.Options){
		func(o *rgapi.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = endpointResolver
		},
		rgapi.WithAPIOptions(awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target), awsmetrics.WithCAPAUserAgentMiddleware()),
	}

	return rgapi.NewFromConfig(cfg, opts...)
}

// NewSecretsManagerClient creates a new Secrets API client for a given session..
func NewSecretsManagerClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *secretsmanager.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	secretsManagerEndpointResolver := &endpoints.SecretsManagerEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}
	secretsManagerOpts := []func(*secretsmanager.Options){
		func(o *secretsmanager.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = secretsManagerEndpointResolver
		},
		secretsmanager.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return secretsmanager.NewFromConfig(cfg, secretsManagerOpts...)
}

// NewEKSClient creates a new EKS API client for a given session.
func NewEKSClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *eks.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	eksEndpointResolver := &endpoints.EKSEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}
	s3Opts := []func(*eks.Options){
		func(o *eks.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = eksEndpointResolver
		},
		eks.WithAPIOptions(awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target), awsmetrics.WithCAPAUserAgentMiddleware()),
	}
	return eks.NewFromConfig(cfg, s3Opts...)
}

// NewIAMClient creates a new IAM API client for a given session.
func NewIAMClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *iam.Client {
	cfg := session.Session()

	iamOpts := []func(*iam.Options){
		func(o *iam.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
		},
		iam.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return iam.NewFromConfig(cfg, iamOpts...)
}

// NewSTSClient creates a new STS API client for a given session.
func NewSTSClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) stsservice.STSClient {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	stsEndpointResolver := &endpoints.STSEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}

	stsOpts := []func(*stsv2.Options){
		func(o *stsv2.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = stsEndpointResolver
		},
		stsv2.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return stsservice.NewClientWrapper(stsv2.NewFromConfig(cfg, stsOpts...))
}

// NewSSMClient creates a new Secrets API client for a given session.
func NewSSMClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *ssm.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	ssmEndpointResolver := &endpoints.SSMEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}
	ssmOpts := []func(*ssm.Options){
		func(o *ssm.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = ssmEndpointResolver
		},
		ssm.WithAPIOptions(
			awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetrics.WithCAPAUserAgentMiddleware(),
		),
	}

	return ssm.NewFromConfig(cfg, ssmOpts...)
}

// NewS3Client creates a new S3 API client for a given session.
func NewS3Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *s3.Client {
	cfg := session.Session()
	multiSvcEndpointResolver := endpoints.NewMultiServiceEndpointResolver()
	s3EndpointResolver := &endpoints.S3EndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}
	s3Opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevel(logger.GetLogger())
			o.EndpointResolverV2 = s3EndpointResolver
		},
		s3.WithAPIOptions(awsmetrics.WithMiddlewares(scopeUser.ControllerName(), target), awsmetrics.WithCAPAUserAgentMiddleware()),
	}
	return s3.NewFromConfig(cfg, s3Opts...)
}

// AWSClients contains all the aws clients used by the scopes.
type AWSClients struct {
	ELB             *elb.Client
	SecretsManager  *secretsmanager.Client
	ResourceTagging *rgapi.Client
	ASG             *autoscaling.Client
	EC2             *ec2.Client
	ELBV2           *elbv2.Client
}
