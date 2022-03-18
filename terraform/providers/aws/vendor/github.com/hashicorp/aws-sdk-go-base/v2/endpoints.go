package awsbase

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// This endpoint resolver is needed when authenticating because the AWS SDK makes internal
// calls to STS. The resolver should not be attached to the aws.Config returned to the
// client, since it should configure its own overrides
func credentialsEndpointResolver(c *Config) aws.EndpointResolverWithOptions {
	resolver := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		switch service {
		case iam.ServiceID:
			if endpoint := c.IamEndpoint; endpoint != "" {
				log.Printf("[INFO] Credentials resolution: setting custom IAM endpoint: %s", endpoint)
				return aws.Endpoint{
					URL:           endpoint,
					Source:        aws.EndpointSourceCustom,
					SigningRegion: region,
				}, nil
			}
		case sts.ServiceID:
			if endpoint := c.StsEndpoint; endpoint != "" {
				if c.StsRegion != "" {
					log.Printf("[INFO] Credentials resolution: setting custom STS endpoint: %s (with signing region %s)", endpoint, c.StsRegion)
					region = c.StsRegion
				} else {
					log.Printf("[INFO] Credentials resolution: setting custom STS endpoint: %s", endpoint)
				}
				return aws.Endpoint{
					URL:           endpoint,
					Source:        aws.EndpointSourceCustom,
					SigningRegion: region,
				}, nil
			}
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	return aws.EndpointResolverWithOptionsFunc(resolver)
}
