package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// NewEC2Client creates a new EC2 API client.
func NewEC2Client(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*ec2.Options)) (*ec2.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	ec2Opts := []func(*ec2.Options){
		func(o *ec2.Options) {
			o.EndpointResolverV2 = &EC2EndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	ec2Opts = append(ec2Opts, optFns...)

	return ec2.NewFromConfig(cfg, ec2Opts...), nil
}

// NewELBClient creates a new ELB (classic) client.
func NewELBClient(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*elb.Options)) (*elb.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	elbOpts := []func(*elb.Options){
		func(o *elb.Options) {
			o.EndpointResolverV2 = &ELBEndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	elbOpts = append(elbOpts, optFns...)

	return elb.NewFromConfig(cfg, elbOpts...), nil
}

// NewELBV2Client creates a new ELBV2 client.
func NewELBV2Client(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*elbv2.Options)) (*elbv2.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	elbv2Opts := []func(*elbv2.Options){
		func(o *elbv2.Options) {
			o.EndpointResolverV2 = &ELBV2EndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	elbv2Opts = append(elbv2Opts, optFns...)

	return elbv2.NewFromConfig(cfg, elbv2Opts...), nil
}

// NewIAMClient creates a new IAM API client.
func NewIAMClient(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*iam.Options)) (*iam.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	iamOpts := []func(*iam.Options){
		func(o *iam.Options) {
			o.EndpointResolverV2 = &IAMEndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	iamOpts = append(iamOpts, optFns...)

	return iam.NewFromConfig(cfg, iamOpts...), nil
}

// NewRoute53Client creates a new Route 53 API client.
func NewRoute53Client(ctx context.Context, endpointOpts EndpointOptions, roleArn string, optFns ...func(*route53.Options)) (*route53.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	if len(roleArn) > 0 {
		stsClient, err := NewSTSClient(ctx, endpointOpts)
		if err != nil {
			return nil, err
		}
		creds := stscreds.NewAssumeRoleProvider(stsClient, roleArn)
		cfg.Credentials = aws.NewCredentialsCache(creds)
	}

	route53Opts := []func(*route53.Options){
		func(o *route53.Options) {
			o.EndpointResolverV2 = &Route53EndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	route53Opts = append(route53Opts, optFns...)

	return route53.NewFromConfig(cfg, route53Opts...), nil
}

// NewSTSClient creates a new STS API client.
func NewSTSClient(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*sts.Options)) (*sts.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	stsOpts := []func(*sts.Options){
		func(o *sts.Options) {
			o.EndpointResolverV2 = &STSEndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	stsOpts = append(stsOpts, optFns...)

	return sts.NewFromConfig(cfg, stsOpts...), nil
}

// NewServiceQuotasClient creates a new Service Quotas API client.
func NewServiceQuotasClient(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*servicequotas.Options)) (*servicequotas.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	sqOpts := []func(*servicequotas.Options){
		func(o *servicequotas.Options) {
			o.EndpointResolverV2 = &ServiceQuotasEndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	sqOpts = append(sqOpts, optFns...)

	return servicequotas.NewFromConfig(cfg, sqOpts...), nil
}

// NewS3Client creates a new S3 API client.
func NewS3Client(ctx context.Context, endpointOpts EndpointOptions, optFns ...func(*s3.Options)) (*s3.Client, error) {
	cfg, err := GetConfigWithOptions(ctx, config.WithRegion(endpointOpts.Region))
	if err != nil {
		return nil, err
	}

	s3Opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.EndpointResolverV2 = &S3EndpointResolver{
				ServiceEndpointResolver: NewServiceEndpointResolver(endpointOpts),
			}
		},
	}
	s3Opts = append(s3Opts, optFns...)

	return s3.NewFromConfig(cfg, s3Opts...), nil
}
