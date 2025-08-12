package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/route53"
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
