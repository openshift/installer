// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsbase

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/logging"
)

const (
	configSourceProviderConfig      = "provider"
	configSourceEnvironmentVariable = "envvar"
)

func getCredentialsProvider(ctx context.Context, c *Config) (aws.CredentialsProvider, string, error) {
	logger := logging.RetrieveLogger(ctx)

	loadOptions, err := commonLoadOptions(ctx, c)
	if err != nil {
		return nil, "", err
	}
	loadOptions = append(
		loadOptions,
		// The endpoint resolver is added here instead of in commonLoadOptions() so that it
		// is not included in the aws.Config returned to the caller
		config.WithEndpointResolverWithOptions(credentialsEndpointResolver(ctx, c)),
		config.WithLogConfigurationWarnings(true),
	)

	envConfig, err := config.NewEnvConfig()
	if err != nil {
		return nil, "", err
	}

	if c.Profile != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		logger.Warn(ctx, `A Profile was specified along with the environment variables "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY". `+
			"The Profile is now used instead of the environment variable credentials. This may lead to unexpected behavior.")
	}

	// The default AWS SDK authentication flow silently ignores invalid Profiles. Pre-validate that the Profile exists
	// https://github.com/aws/aws-sdk-go-v2/issues/1591
	profile := c.Profile
	var profileSource string
	if profile != "" {
		profileSource = configSourceProviderConfig
	} else {
		profile = envConfig.SharedConfigProfile
		profileSource = configSourceEnvironmentVariable
	}

	if profile != "" {
		logger.Debug(ctx, "Using profile", map[string]any{
			"tf_aws.profile":        profile,
			"tf_aws.profile.source": profileSource,
		})
		sharedCredentialsFiles, err := c.ResolveSharedCredentialsFiles()
		if err != nil {
			return nil, "", err
		}
		if len(sharedCredentialsFiles) != 0 {
			f := make([]string, len(sharedCredentialsFiles))
			for i, v := range sharedCredentialsFiles {
				f[i] = fmt.Sprintf(`"%s"`, v)
			}
			logger.Debug(ctx, "Using shared credentials files", map[string]any{
				"tf_aws.shared_credentials_files":        f,
				"tf_aws.shared_credentials_files.source": configSourceProviderConfig,
			})
		} else {
			if envConfig.SharedCredentialsFile != "" {
				sharedCredentialsFiles = []string{envConfig.SharedCredentialsFile}
				logger.Debug(ctx, "Using shared credentials files", map[string]any{
					"tf_aws.shared_credentials_files":        sharedCredentialsFiles,
					"tf_aws.shared_credentials_files.source": configSourceEnvironmentVariable,
				})
			}
		}

		sharedConfigFiles, err := c.ResolveSharedConfigFiles()
		if err != nil {
			return nil, "", err
		}
		if len(sharedConfigFiles) != 0 {
			f := make([]string, len(sharedConfigFiles))
			for i, v := range sharedConfigFiles {
				f[i] = fmt.Sprintf(`"%s"`, v)
			}
			logger.Debug(ctx, "Using shared configuration files", map[string]any{
				"tf_aws.shared_config_files":        f,
				"tf_aws.shared_config_files.source": configSourceProviderConfig,
			})
		} else {
			if envConfig.SharedConfigFile != "" {
				sharedConfigFiles = []string{envConfig.SharedConfigFile}
				logger.Debug(ctx, "Using shared configuration files", map[string]any{
					"tf_aws.shared_config_files":        sharedConfigFiles,
					"tf_aws.shared_config_files.source": configSourceEnvironmentVariable,
				})
			}
		}

		logger.Debug(ctx, "Loading profile", map[string]any{
			"tf_aws.profile": profile,
		})
		_, err = config.LoadSharedConfigProfile(ctx, profile, func(opts *config.LoadSharedConfigOptions) {
			if len(sharedCredentialsFiles) != 0 {
				opts.CredentialsFiles = sharedCredentialsFiles
			}
			if len(sharedConfigFiles) != 0 {
				opts.ConfigFiles = sharedConfigFiles
			}
		})
		if err != nil {
			return nil, "", err
		}
	}
	// We need to validate both the configured and envvar named profiles for validity,
	// but to use proper precedence, we only set the configured named profile
	if c.Profile != "" {
		logger.Debug(ctx, "Setting profile", map[string]any{
			"tf_aws.profile":        profile,
			"tf_aws.profile.source": configSourceProviderConfig,
		})
		loadOptions = append(
			loadOptions,
			config.WithSharedConfigProfile(c.Profile),
		)
	}

	if c.AccessKey != "" || c.SecretKey != "" || c.Token != "" {
		params := make([]string, 0, 3) //nolint:gomnd
		if c.AccessKey != "" {
			params = append(params, "access key")
		}
		if c.SecretKey != "" {
			params = append(params, "secret key")
		}
		if c.Token != "" {
			params = append(params, "token")
		}
		logger.Debug(ctx, "Using authentication parameters", map[string]any{
			"tf_aws.auth_fields":        params,
			"tf_aws.auth_fields.source": configSourceProviderConfig,
		})
		loadOptions = append(
			loadOptions,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					c.AccessKey,
					c.SecretKey,
					c.Token,
				),
			),
		)
	}

	logger.Debug(ctx, "Loading configuration")
	cfg, err := config.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, "", fmt.Errorf("loading configuration: %w", err)
	}

	// This can probably be configured directly in commonLoadOptions() once
	// https://github.com/aws/aws-sdk-go-v2/pull/1682 is merged
	if c.AssumeRoleWithWebIdentity != nil {
		if c.AssumeRoleWithWebIdentity.RoleARN == "" {
			return nil, "", errors.New("Assume Role With Web Identity: role ARN not set")
		}
		if c.AssumeRoleWithWebIdentity.WebIdentityToken == "" && c.AssumeRoleWithWebIdentity.WebIdentityTokenFile == "" {
			return nil, "", errors.New("Assume Role With Web Identity: one of WebIdentityToken, WebIdentityTokenFile must be set")
		}
		provider, err := webIdentityCredentialsProvider(ctx, cfg, c)
		if err != nil {
			return nil, "", err
		}
		cfg.Credentials = provider
	}

	logger.Debug(ctx, "Retrieving credentials")
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		if c.Profile != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
			err = fmt.Errorf(`A Profile was specified along with the environment variables "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY". The Profile is now used instead of the environment variable credentials.

AWS Error: %w`, err)
		}
		return nil, "", c.NewNoValidCredentialSourcesError(err)
	}

	if c.AssumeRole == nil {
		return cfg.Credentials, creds.Source, nil
	}

	logger.Info(ctx, "Retrieved initial credentials", map[string]any{
		"tf_aws.credentials_source": creds.Source,
	})
	provider, err := assumeRoleCredentialsProvider(ctx, cfg, c)

	return provider, creds.Source, err
}

func webIdentityCredentialsProvider(ctx context.Context, awsConfig aws.Config, c *Config) (aws.CredentialsProvider, error) {
	logger := logging.RetrieveLogger(ctx)

	ar := c.AssumeRoleWithWebIdentity

	logger.Info(ctx, "Assuming IAM Role With Web Identity", map[string]any{
		"tf_aws.assume_role_with_web_identity.role_arn":     ar.RoleARN,
		"tf_aws.assume_role_with_web_identity.session_name": ar.SessionName,
	})

	client := stsClient(ctx, awsConfig, c)

	appCreds := stscreds.NewWebIdentityRoleProvider(client, ar.RoleARN, ar, func(opts *stscreds.WebIdentityRoleOptions) {
		opts.RoleSessionName = ar.SessionName
		opts.Duration = ar.Duration

		if ar.Policy != "" {
			opts.Policy = aws.String(ar.Policy)
		}

		if len(ar.PolicyARNs) > 0 {
			opts.PolicyARNs = getPolicyDescriptorTypes(ar.PolicyARNs)
		}
	})

	_, err := appCreds.Retrieve(ctx)
	if err != nil {
		return nil, c.NewCannotAssumeRoleWithWebIdentityError(err)
	}
	return aws.NewCredentialsCache(appCreds), nil
}

func assumeRoleCredentialsProvider(ctx context.Context, awsConfig aws.Config, c *Config) (aws.CredentialsProvider, error) {
	logger := logging.RetrieveLogger(ctx)

	ar := c.AssumeRole

	if ar.RoleARN == "" {
		return nil, errors.New("Assume Role: role ARN not set")
	}

	logger.Info(ctx, "Assuming IAM Role", map[string]any{
		"tf_aws.assume_role.role_arn":        ar.RoleARN,
		"tf_aws.assume_role.session_name":    ar.SessionName,
		"tf_aws.assume_role.external_id":     ar.ExternalID,
		"tf_aws.assume_role.source_identity": ar.SourceIdentity,
	})

	// When assuming a role, we need to first authenticate the base credentials above, then assume the desired role
	client := stsClient(ctx, awsConfig, c)

	appCreds := stscreds.NewAssumeRoleProvider(client, ar.RoleARN, func(opts *stscreds.AssumeRoleOptions) {
		opts.RoleSessionName = ar.SessionName
		opts.Duration = ar.Duration

		if ar.ExternalID != "" {
			opts.ExternalID = aws.String(ar.ExternalID)
		}

		if ar.Policy != "" {
			opts.Policy = aws.String(ar.Policy)
		}

		if len(ar.PolicyARNs) > 0 {
			opts.PolicyARNs = getPolicyDescriptorTypes(ar.PolicyARNs)
		}

		if len(ar.Tags) > 0 {
			var tags []types.Tag
			for k, v := range ar.Tags {
				tag := types.Tag{
					Key:   aws.String(k),
					Value: aws.String(v),
				}
				tags = append(tags, tag)
			}

			opts.Tags = tags
		}

		if len(ar.TransitiveTagKeys) > 0 {
			opts.TransitiveTagKeys = ar.TransitiveTagKeys
		}

		if ar.SourceIdentity != "" {
			opts.SourceIdentity = aws.String(ar.SourceIdentity)
		}
	})
	_, err := appCreds.Retrieve(ctx)
	if err != nil {
		return nil, c.NewCannotAssumeRoleError(err)
	}
	return aws.NewCredentialsCache(appCreds), nil
}

func getPolicyDescriptorTypes(policyARNs []string) []types.PolicyDescriptorType {
	var policyDescriptorTypes []types.PolicyDescriptorType

	for _, policyARN := range policyARNs {
		policyDescriptorType := types.PolicyDescriptorType{
			Arn: aws.String(policyARN),
		}
		policyDescriptorTypes = append(policyDescriptorTypes, policyDescriptorType)
	}
	return policyDescriptorTypes
}
