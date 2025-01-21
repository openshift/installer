package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// GetConfig returns the default config AWS configuration.
func GetConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

// GetConfigWithOptions returns the AWS Config with options set by the user information.
func GetConfigWithOptions(ctx context.Context, region, service string, endpoints []awstypes.ServiceEndpoint) (aws.Config, error) {
	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	return config.LoadDefaultConfig(ctx, loadOptions...)
}
