package awsbase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/smithy-go/middleware"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/awsconfig"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/constants"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/endpoints"
)

func GetAwsConfig(ctx context.Context, c *Config) (aws.Config, error) {
	if metadataUrl := os.Getenv("AWS_METADATA_URL"); metadataUrl != "" {
		log.Println(`[WARN] The environment variable "AWS_METADATA_URL" is deprecated. Use "AWS_EC2_METADATA_SERVICE_ENDPOINT" instead.`)
		if ec2MetadataServiceEndpoint := os.Getenv("AWS_EC2_METADATA_SERVICE_ENDPOINT"); ec2MetadataServiceEndpoint != "" {
			if ec2MetadataServiceEndpoint != metadataUrl {
				log.Printf(`[WARN] The environment variable "AWS_EC2_METADATA_SERVICE_ENDPOINT" is already set to %q. Ignoring "AWS_METADATA_URL".`, ec2MetadataServiceEndpoint)
			}
		} else {
			log.Printf(`[WARN] Setting "AWS_EC2_METADATA_SERVICE_ENDPOINT" to %q.`, metadataUrl)
			os.Setenv("AWS_EC2_METADATA_SERVICE_ENDPOINT", metadataUrl)
		}
	}

	credentialsProvider, initialSource, err := getCredentialsProvider(ctx, c)
	if err != nil {
		return aws.Config{}, err
	}
	creds, _ := credentialsProvider.Retrieve(ctx)
	log.Printf("[INFO] Retrieved credentials from %q", creds.Source)

	loadOptions, err := commonLoadOptions(c)
	if err != nil {
		return aws.Config{}, err
	}
	loadOptions = append(
		loadOptions,
		config.WithCredentialsProvider(credentialsProvider),
	)
	if initialSource == ec2rolecreds.ProviderName {
		loadOptions = append(
			loadOptions,
			config.WithEC2IMDSRegion(),
		)
	}
	awsConfig, err := config.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return awsConfig, fmt.Errorf("loading configuration: %w", err)
	}

	resolveRetryer(ctx, &awsConfig)

	if !c.SkipCredsValidation {
		if _, _, err := getAccountIDAndPartitionFromSTSGetCallerIdentity(ctx, stsClient(awsConfig, c)); err != nil {
			return awsConfig, fmt.Errorf("error validating provider credentials: %w", err)
		}
	}

	return awsConfig, nil
}

// Adapted from the per-service-client `resolveRetryer()` functions in the AWS SDK for Go v2
// e.g. https://github.com/aws/aws-sdk-go-v2/blob/main/service/accessanalyzer/api_client.go
// Currently only supports "standard" retry mode
func resolveRetryer(ctx context.Context, awsConfig *aws.Config) {
	var standardOptions []func(*retry.StandardOptions)

	if v, found, _ := awsconfig.GetRetryMaxAttempts(ctx, awsConfig.ConfigSources); found && v != 0 {
		standardOptions = append(standardOptions, func(so *retry.StandardOptions) {
			so.MaxAttempts = v
		})
	}

	awsConfig.Retryer = func() aws.Retryer {
		return &networkErrorShortcutter{
			RetryerV2: retry.NewStandard(standardOptions...),
		}
	}
}

// networkErrorShortcutter is used to enable networking error shortcutting
type networkErrorShortcutter struct {
	aws.RetryerV2
}

// We're misusing RetryDelay here, since this is the only function that takes the attempt count
func (r *networkErrorShortcutter) RetryDelay(attempt int, err error) (time.Duration, error) {
	if attempt >= constants.MaxNetworkRetryCount {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) {
			// It's disappointing that we have to do string matching here, rather than being able to using `errors.Is()` or even strings exported by the Go `net` package
			if strings.Contains(netOpErr.Error(), "no such host") || strings.Contains(netOpErr.Error(), "connection refused") {
				log.Printf("[WARN] Disabling retries after next request due to networking issue: %s", err)
				return 0, &retry.MaxAttemptsError{
					Attempt: attempt,
					Err:     err,
				}
			}
		}
	}

	return r.RetryerV2.RetryDelay(attempt, err)
}

func GetAwsAccountIDAndPartition(ctx context.Context, awsConfig aws.Config, c *Config) (string, string, error) {
	if !c.SkipCredsValidation {
		stsClient := stsClient(awsConfig, c)
		accountID, partition, err := getAccountIDAndPartitionFromSTSGetCallerIdentity(ctx, stsClient)
		if err != nil {
			return "", "", fmt.Errorf("error validating provider credentials: %w", err)
		}

		return accountID, partition, nil
	}

	if !c.SkipRequestingAccountId {
		credentialsProviderName := ""
		if credentialsValue, err := awsConfig.Credentials.Retrieve(context.Background()); err == nil {
			credentialsProviderName = credentialsValue.Source
		}

		iamClient := iamClient(awsConfig, c)
		stsClient := stsClient(awsConfig, c)
		accountID, partition, err := getAccountIDAndPartition(ctx, iamClient, stsClient, credentialsProviderName)

		if err == nil {
			return accountID, partition, nil
		}

		return "", "", fmt.Errorf(
			"AWS account ID not previously found and failed retrieving via all available methods. "+
				"See https://www.terraform.io/docs/providers/aws/index.html#skip_requesting_account_id for workaround and implications. "+
				"Errors: %w", err)
	}

	return "", endpoints.PartitionForRegion(awsConfig.Region), nil
}

func commonLoadOptions(c *Config) ([]func(*config.LoadOptions) error, error) {
	httpClient, err := defaultHttpClient(c)
	if err != nil {
		return nil, err
	}

	apiOptions := make([]func(*middleware.Stack) error, 0)
	if c.APNInfo != nil {
		apiOptions = append(apiOptions, func(stack *middleware.Stack) error {
			// Because the default User-Agent middleware prepends itself to the contents of the User-Agent header,
			// we have to run after it and also prepend our custom User-Agent
			return stack.Build.Add(apnUserAgentMiddleware(*c.APNInfo), middleware.After)
		})
	}

	if len(c.UserAgent) > 0 {
		apiOptions = append(apiOptions, awsmiddleware.AddUserAgentKey(c.UserAgent.BuildUserAgentString()))
	}

	apiOptions = append(apiOptions, func(stack *middleware.Stack) error {
		return stack.Build.Add(userAgentFromContextMiddleware(), middleware.After)
	})

	if v := os.Getenv(constants.AppendUserAgentEnvVar); v != "" {
		log.Printf("[DEBUG] Using additional User-Agent Info: %s", v)
		apiOptions = append(apiOptions, awsmiddleware.AddUserAgentKey(v))
	}

	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion(c.Region),
		config.WithHTTPClient(httpClient),
		config.WithAPIOptions(apiOptions),
		config.WithEC2IMDSClientEnableState(c.EC2MetadataServiceEnableState),
	}

	if c.MaxRetries != 0 {
		loadOptions = append(
			loadOptions,
			config.WithRetryMaxAttempts(c.MaxRetries),
		)
	}

	if !c.SuppressDebugLog {
		loadOptions = append(
			loadOptions,
			config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody|aws.LogRetries),
			config.WithLogger(debugLogger{}),
		)
	}

	sharedCredentialsFiles, err := c.ResolveSharedCredentialsFiles()
	if err != nil {
		return nil, err
	}
	if len(sharedCredentialsFiles) > 0 {
		loadOptions = append(
			loadOptions,
			config.WithSharedCredentialsFiles(sharedCredentialsFiles),
		)
	}

	sharedConfigFiles, err := c.ResolveSharedConfigFiles()
	if err != nil {
		return nil, err
	}
	if len(sharedConfigFiles) > 0 {
		loadOptions = append(
			loadOptions,
			config.WithSharedConfigFiles(sharedConfigFiles),
		)
	}

	if c.CustomCABundle != "" {
		reader, err := c.CustomCABundleReader()
		if err != nil {
			return nil, err
		}
		loadOptions = append(loadOptions,
			config.WithCustomCABundle(reader),
		)
	}

	if c.EC2MetadataServiceEndpoint != "" {
		loadOptions = append(loadOptions,
			config.WithEC2IMDSEndpoint(c.EC2MetadataServiceEndpoint),
		)
	}

	if c.EC2MetadataServiceEndpointMode != "" {
		var endpointMode imds.EndpointModeState
		err := endpointMode.SetFromString(c.EC2MetadataServiceEndpointMode)
		if err != nil {
			return nil, err
		}
		loadOptions = append(loadOptions,
			config.WithEC2IMDSEndpointMode(endpointMode),
		)
	}

	// This should not be needed, but https://github.com/aws/aws-sdk-go-v2/issues/1398
	if c.EC2MetadataServiceEnableState == imds.ClientEnabled {
		os.Setenv("AWS_EC2_METADATA_DISABLED", "false")
	} else if c.EC2MetadataServiceEnableState == imds.ClientDisabled {
		os.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	}

	if c.UseDualStackEndpoint {
		loadOptions = append(loadOptions,
			config.WithUseDualStackEndpoint(aws.DualStackEndpointStateEnabled),
		)
	}

	if c.UseFIPSEndpoint {
		loadOptions = append(loadOptions,
			config.WithUseFIPSEndpoint(aws.FIPSEndpointStateEnabled),
		)
	}

	return loadOptions, nil
}
