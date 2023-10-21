// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/smithy-go/middleware"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/awsconfig"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/constants"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/endpoints"
	"github.com/hashicorp/aws-sdk-go-base/v2/logging"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const loggerName string = "aws-base"

func configCommonLogging(ctx context.Context) context.Context {
	// Catch as last resort, but prefer the custom masking in the request-response logging
	return tflog.MaskAllFieldValuesRegexes(ctx, logging.UniqueIDRegex)
}

func GetAwsConfig(ctx context.Context, c *Config) (context.Context, aws.Config, error) {
	ctx = configCommonLogging(ctx)

	baseCtx, logger := logging.New(ctx, loggerName)
	baseCtx = logging.RegisterLogger(baseCtx, logger)

	logger.Trace(baseCtx, "Resolving AWS configuration")

	if metadataUrl := os.Getenv("AWS_METADATA_URL"); metadataUrl != "" {
		logger.Warn(baseCtx, `The environment variable "AWS_METADATA_URL" is deprecated. Use "AWS_EC2_METADATA_SERVICE_ENDPOINT" instead.`)
		if ec2MetadataServiceEndpoint := os.Getenv("AWS_EC2_METADATA_SERVICE_ENDPOINT"); ec2MetadataServiceEndpoint != "" {
			if ec2MetadataServiceEndpoint != metadataUrl {
				logger.Warn(baseCtx, fmt.Sprintf(`The environment variable "AWS_EC2_METADATA_SERVICE_ENDPOINT" is already set to %q. Ignoring "AWS_METADATA_URL".`, ec2MetadataServiceEndpoint))
			}
		} else {
			logger.Warn(baseCtx, fmt.Sprintf(`Setting "AWS_EC2_METADATA_SERVICE_ENDPOINT" to %q.`, metadataUrl))
			os.Setenv("AWS_EC2_METADATA_SERVICE_ENDPOINT", metadataUrl)
		}
	}

	logger.Debug(baseCtx, "Resolving credentials provider")
	credentialsProvider, initialSource, err := getCredentialsProvider(baseCtx, c)
	if err != nil {
		return ctx, aws.Config{}, err
	}
	creds, err := credentialsProvider.Retrieve(baseCtx)
	if err != nil {
		return ctx, aws.Config{}, fmt.Errorf("retrieving credentials: %w", err)
	}
	logger.Info(baseCtx, "Retrieved credentials", map[string]any{
		"tf_aws.credentials_source": creds.Source,
	})

	loadOptions, err := commonLoadOptions(baseCtx, c)
	if err != nil {
		return ctx, aws.Config{}, err
	}

	// The providers set `MaxRetries` to a very large value.
	// Add retries here so that authentication has a reasonable number of retries
	if c.MaxRetries != 0 {
		loadOptions = append(
			loadOptions,
			config.WithRetryMaxAttempts(c.MaxRetries),
		)
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

	logger.Debug(baseCtx, "Loading configuration")
	awsConfig, err := config.LoadDefaultConfig(baseCtx, loadOptions...)
	if err != nil {
		return ctx, aws.Config{}, fmt.Errorf("loading configuration: %w", err)
	}

	resolveRetryer(baseCtx, &awsConfig)

	if !c.SkipCredsValidation {
		if _, _, err := getAccountIDAndPartitionFromSTSGetCallerIdentity(baseCtx, stsClient(baseCtx, awsConfig, c)); err != nil {
			return ctx, awsConfig, fmt.Errorf("validating provider credentials: %w", err)
		}
	}

	return ctx, awsConfig, nil
}

// Adapted from the per-service-client `resolveRetryer()` functions in the AWS SDK for Go v2
// e.g. https://github.com/aws/aws-sdk-go-v2/blob/main/service/accessanalyzer/api_client.go
func resolveRetryer(ctx context.Context, awsConfig *aws.Config) {
	retryMode := awsConfig.RetryMode
	if len(retryMode) == 0 {
		defaultsMode := resolveDefaultsMode(ctx, awsConfig)
		modeConfig, err := defaults.GetModeConfiguration(defaultsMode)
		if err == nil {
			retryMode = modeConfig.RetryMode
		}
	}
	if len(retryMode) == 0 {
		retryMode = aws.RetryModeStandard
	}

	var standardOptions []func(*retry.StandardOptions)
	if v, found, _ := awsconfig.GetRetryMaxAttempts(ctx, awsConfig.ConfigSources); found && v != 0 {
		standardOptions = append(standardOptions, func(so *retry.StandardOptions) {
			so.MaxAttempts = v
		})
	}

	var retryer aws.RetryerV2
	switch retryMode {
	case aws.RetryModeAdaptive:
		var adaptiveOptions []func(*retry.AdaptiveModeOptions)
		if len(standardOptions) != 0 {
			adaptiveOptions = append(adaptiveOptions, func(ao *retry.AdaptiveModeOptions) {
				ao.StandardOptions = append(ao.StandardOptions, standardOptions...)
			})
		}
		retryer = retry.NewAdaptiveMode(adaptiveOptions...)

	default:
		retryer = retry.NewStandard(standardOptions...)
	}

	awsConfig.Retryer = func() aws.Retryer {
		return &networkErrorShortcutter{
			RetryerV2: retryer,
		}
	}
}

// Adapted from the per-service-client `setResolvedDefaultsMode()` functions in the AWS SDK for Go v2
// e.g. https://github.com/aws/aws-sdk-go-v2/blob/main/service/accessanalyzer/api_client.go
func resolveDefaultsMode(_ context.Context, awsConfig *aws.Config) aws.DefaultsMode {
	var mode aws.DefaultsMode
	mode.SetFromString(string(awsConfig.DefaultsMode))

	if mode == aws.DefaultsModeAuto {
		mode = defaults.ResolveDefaultsModeAuto(awsConfig.Region, awsConfig.RuntimeEnvironment)
	}

	return mode
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
				// TODO: figure out how to get correct logger here
				log.Printf("[WARN] Disabling retries after next request due to networking error: %s", err)
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
	ctx, logger := logging.New(ctx, loggerName)
	ctx = logging.RegisterLogger(ctx, logger)

	if !c.SkipCredsValidation {
		stsClient := stsClient(ctx, awsConfig, c)
		accountID, partition, err := getAccountIDAndPartitionFromSTSGetCallerIdentity(ctx, stsClient)
		if err != nil {
			return "", "", fmt.Errorf("validating provider credentials: %w", err)
		}

		return accountID, partition, nil
	}

	if !c.SkipRequestingAccountId {
		credentialsProviderName := ""
		if credentialsValue, err := awsConfig.Credentials.Retrieve(context.Background()); err == nil {
			credentialsProviderName = credentialsValue.Source
		}

		iamClient := iamClient(ctx, awsConfig, c)
		stsClient := stsClient(ctx, awsConfig, c)
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

func commonLoadOptions(ctx context.Context, c *Config) ([]func(*config.LoadOptions) error, error) {
	logger := logging.RetrieveLogger(ctx)

	var err error
	var httpClient config.HTTPClient

	if v := c.HTTPClient; v == nil {
		logger.Trace(ctx, "Building default HTTP client")
		httpClient, err = defaultHttpClient(c)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Debug(ctx, "Setting HTTP client", map[string]any{
			"tf_aws.http_client.source": configSourceProviderConfig,
		})
		httpClient = v
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
		logger.Debug(ctx, "Adding User-Agent info", map[string]any{
			"source": fmt.Sprintf("envvar(%q)", constants.AppendUserAgentEnvVar),
			"value":  v,
		})
		apiOptions = append(apiOptions, awsmiddleware.AddUserAgentKey(v))
	}

	if !c.SuppressDebugLog {
		apiOptions = append(apiOptions, func(stack *middleware.Stack) error {
			return stack.Deserialize.Add(&requestResponseLogger{}, middleware.After)
		})
	}

	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion(c.Region),
		config.WithHTTPClient(httpClient),
		config.WithAPIOptions(apiOptions),
		config.WithEC2IMDSClientEnableState(c.EC2MetadataServiceEnableState),
	}

	if !c.SuppressDebugLog {
		loadOptions = append(
			loadOptions,
			config.WithClientLogMode(aws.LogDeprecatedUsage|aws.LogRetries),
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

	if c.RetryMode != "" {
		loadOptions = append(loadOptions,
			config.WithRetryMode(c.RetryMode),
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
