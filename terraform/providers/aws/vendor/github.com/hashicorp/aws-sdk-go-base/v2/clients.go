package awsbase

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func iamClient(awsConfig aws.Config, c *Config) *iam.Client {
	return iam.NewFromConfig(awsConfig, func(opts *iam.Options) {
		if c.IamEndpoint != "" {
			log.Printf("[INFO] IAM client: setting custom endpoint: %s", c.IamEndpoint)
			opts.EndpointResolver = iam.EndpointResolverFromURL(c.IamEndpoint)
		}
	})
}

func stsClient(awsConfig aws.Config, c *Config) *sts.Client {
	return sts.NewFromConfig(awsConfig, func(opts *sts.Options) {
		if c.StsRegion != "" {
			log.Printf("[INFO] STS client: setting region: %s", c.StsRegion)
			opts.Region = c.StsRegion
		}
		if c.StsEndpoint != "" {
			log.Printf("[INFO] STS client: setting custom endpoint: %s", c.StsEndpoint)
			opts.EndpointResolver = sts.EndpointResolverFromURL(c.StsEndpoint)
		}
	})
}
