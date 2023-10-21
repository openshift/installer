package conns

import (
	"context"
	"log"
	"strings"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/aws/aws-sdk-go/service/chime"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudhsmv2"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/fms"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53recoverycontrolconfig"
	"github.com/aws/aws-sdk-go/service/route53recoveryreadiness"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/aws/aws-sdk-go/service/shield"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/aws/aws-sdk-go/service/storagegateway"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/wafv2"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	awsbasev1 "github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	tenBackOff = 10
)

type Config struct {
	AccessKey                      string
	AllowedAccountIds              []string
	AssumeRole                     *awsbase.AssumeRole
	AssumeRoleWithWebIdentity      *awsbase.AssumeRoleWithWebIdentity
	CustomCABundle                 string
	DefaultTagsConfig              *tftags.DefaultConfig
	EC2MetadataServiceEnableState  imds.ClientEnableState
	EC2MetadataServiceEndpoint     string
	EC2MetadataServiceEndpointMode string
	Endpoints                      map[string]string
	ForbiddenAccountIds            []string
	HTTPProxy                      string
	IgnoreTagsConfig               *tftags.IgnoreConfig
	Insecure                       bool
	MaxRetries                     int
	Profile                        string
	Region                         string
	RetryMode                      awsv2.RetryMode
	S3UsePathStyle                 bool
	SecretKey                      string
	SharedConfigFiles              []string
	SharedCredentialsFiles         []string
	SkipCredsValidation            bool
	SkipRegionValidation           bool
	SkipRequestingAccountId        bool
	STSRegion                      string
	SuppressDebugLog               bool
	TerraformVersion               string
	Token                          string
	UseDualStackEndpoint           bool
	UseFIPSEndpoint                bool
}

// ConfigureProvider configures the provided provider Meta (instance data).
func (c *Config) ConfigureProvider(ctx context.Context, client *AWSClient) (*AWSClient, diag.Diagnostics) {
	awsbaseConfig := awsbase.Config{
		AccessKey:                     c.AccessKey,
		APNInfo:                       StdUserAgentProducts(c.TerraformVersion),
		AssumeRoleWithWebIdentity:     c.AssumeRoleWithWebIdentity,
		CallerDocumentationURL:        "https://registry.terraform.io/providers/hashicorp/aws",
		CallerName:                    "Terraform AWS Provider",
		EC2MetadataServiceEnableState: c.EC2MetadataServiceEnableState,
		IamEndpoint:                   c.Endpoints[names.IAM],
		Insecure:                      c.Insecure,
		HTTPClient:                    client.HTTPClient(),
		HTTPProxy:                     c.HTTPProxy,
		MaxRetries:                    c.MaxRetries,
		Profile:                       c.Profile,
		Region:                        c.Region,
		RetryMode:                     c.RetryMode,
		SecretKey:                     c.SecretKey,
		SkipCredsValidation:           c.SkipCredsValidation,
		SkipRequestingAccountId:       c.SkipRequestingAccountId,
		StsEndpoint:                   c.Endpoints[names.STS],
		SuppressDebugLog:              c.SuppressDebugLog,
		Token:                         c.Token,
		UseDualStackEndpoint:          c.UseDualStackEndpoint,
		UseFIPSEndpoint:               c.UseFIPSEndpoint,
	}

	if c.AssumeRole != nil && c.AssumeRole.RoleARN != "" {
		awsbaseConfig.AssumeRole = c.AssumeRole
	}

	if c.CustomCABundle != "" {
		awsbaseConfig.CustomCABundle = c.CustomCABundle
	}

	if c.EC2MetadataServiceEndpoint != "" {
		awsbaseConfig.EC2MetadataServiceEndpoint = c.EC2MetadataServiceEndpoint
		awsbaseConfig.EC2MetadataServiceEndpointMode = c.EC2MetadataServiceEndpointMode
	}

	if len(c.SharedConfigFiles) != 0 {
		awsbaseConfig.SharedConfigFiles = c.SharedConfigFiles
	}

	if len(c.SharedCredentialsFiles) != 0 {
		awsbaseConfig.SharedCredentialsFiles = c.SharedCredentialsFiles
	}

	if c.STSRegion != "" {
		awsbaseConfig.StsRegion = c.STSRegion
	}

	tflog.Debug(ctx, "Configuring Terraform AWS Provider")
	ctx, cfg, err := awsbase.GetAwsConfig(ctx, &awsbaseConfig)
	if err != nil {
		return nil, diag.Errorf("configuring Terraform AWS Provider: %s", err)
	}

	if !c.SkipRegionValidation {
		if err := awsbase.ValidateRegion(cfg.Region); err != nil {
			return nil, diag.FromErr(err)
		}
	}
	c.Region = cfg.Region

	tflog.Debug(ctx, "Creating AWS SDK v1 session")
	sess, err := awsbasev1.GetSession(ctx, &cfg, &awsbaseConfig)
	if err != nil {
		return nil, diag.Errorf("creating AWS SDK v1 session: %s", err)
	}

	tflog.Debug(ctx, "Retrieving AWS account details")
	accountID, partition, err := awsbase.GetAwsAccountIDAndPartition(ctx, cfg, &awsbaseConfig)
	if err != nil {
		return nil, diag.Errorf("retrieving AWS account details: %s", err)
	}

	if accountID == "" {
		// TODO: Make this a Warning Diagnostic
		log.Println("[WARN] AWS account ID not found for provider. See https://www.terraform.io/docs/providers/aws/index.html#skip_requesting_account_id for implications.")
	}

	if len(c.ForbiddenAccountIds) > 0 {
		for _, forbiddenAccountID := range c.AllowedAccountIds {
			if accountID == forbiddenAccountID {
				return nil, diag.Errorf("AWS account ID not allowed: %s", accountID)
			}
		}
	}
	if len(c.AllowedAccountIds) > 0 {
		found := false
		for _, allowedAccountID := range c.AllowedAccountIds {
			if accountID == allowedAccountID {
				found = true
				break
			}
		}
		if !found {
			return nil, diag.Errorf("AWS account ID not allowed: %s", accountID)
		}
	}

	DNSSuffix := "amazonaws.com"
	if p, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), c.Region); ok {
		DNSSuffix = p.DNSSuffix()
	}

	client.AccountID = accountID
	client.DefaultTagsConfig = c.DefaultTagsConfig
	client.DNSSuffix = DNSSuffix
	client.IgnoreTagsConfig = c.IgnoreTagsConfig
	client.Partition = partition
	client.Region = c.Region
	client.ReverseDNSPrefix = ReverseDNS(DNSSuffix)
	client.SetHTTPClient(sess.Config.HTTPClient) // Must be called while client.Session is nil.
	client.Session = sess
	client.TerraformVersion = c.TerraformVersion

	// API clients (generated).
	c.sdkv1Conns(client, sess)
	c.sdkv2Conns(client, cfg)
	c.sdkv2LazyConns(client, cfg)

	// AWS SDK for Go v1 custom API clients.

	// STS.
	stsConfig := &aws.Config{
		Endpoint: aws.String(c.Endpoints[names.STS]),
	}
	if c.STSRegion != "" {
		stsConfig.Region = aws.String(c.STSRegion)
	}
	client.stsConn = sts.New(sess.Copy(stsConfig))

	// Services that require multiple client configurations.
	s3Config := &aws.Config{
		Endpoint:         aws.String(c.Endpoints[names.S3]),
		S3ForcePathStyle: aws.Bool(c.S3UsePathStyle),
	}
	client.s3Conn = s3.New(sess.Copy(s3Config))

	// "Global" services that require customizations.
	globalAcceleratorConfig := &aws.Config{
		Endpoint: aws.String(c.Endpoints[names.GlobalAccelerator]),
	}
	route53Config := &aws.Config{
		Endpoint: aws.String(c.Endpoints[names.Route53]),
	}
	route53RecoveryControlConfigConfig := &aws.Config{
		Endpoint: aws.String(c.Endpoints[names.Route53RecoveryControlConfig]),
	}
	route53RecoveryReadinessConfig := &aws.Config{
		Endpoint: aws.String(c.Endpoints[names.Route53RecoveryReadiness]),
	}
	shieldConfig := &aws.Config{
		Endpoint: aws.String(c.Endpoints[names.Shield]),
	}

	// Force "global" services to correct Regions.
	switch partition {
	case endpoints.AwsPartitionID:
		globalAcceleratorConfig.Region = aws.String(endpoints.UsWest2RegionID)
		route53Config.Region = aws.String(endpoints.UsEast1RegionID)
		route53RecoveryControlConfigConfig.Region = aws.String(endpoints.UsWest2RegionID)
		route53RecoveryReadinessConfig.Region = aws.String(endpoints.UsWest2RegionID)
		shieldConfig.Region = aws.String(endpoints.UsEast1RegionID)
	case endpoints.AwsCnPartitionID:
		// The AWS Go SDK is missing endpoint information for Route 53 in the AWS China partition.
		// This can likely be removed in the future.
		if aws.StringValue(route53Config.Endpoint) == "" {
			route53Config.Endpoint = aws.String("https://api.route53.cn")
		}
		route53Config.Region = aws.String(endpoints.CnNorthwest1RegionID)
	case endpoints.AwsUsGovPartitionID:
		route53Config.Region = aws.String(endpoints.UsGovWest1RegionID)
	}

	client.globalacceleratorConn = globalaccelerator.New(sess.Copy(globalAcceleratorConfig))
	client.route53Conn = route53.New(sess.Copy(route53Config))
	client.route53recoverycontrolconfigConn = route53recoverycontrolconfig.New(sess.Copy(route53RecoveryControlConfigConfig))
	client.route53recoveryreadinessConn = route53recoveryreadiness.New(sess.Copy(route53RecoveryReadinessConfig))
	client.shieldConn = shield.New(sess.Copy(shieldConfig))

	client.apigatewayConn.Handlers.Retry.PushBack(func(r *request.Request) {
		// Many operations can return an error such as:
		//   ConflictException: Unable to complete operation due to concurrent modification. Please try again later.
		// Handle them all globally for the service client.
		if tfawserr.ErrMessageContains(r.Error, apigateway.ErrCodeConflictException, "try again later") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.apigatewayv2Conn.Handlers.Retry.PushBack(func(r *request.Request) {
		// Many operations can return an error such as:
		//   ConflictException: Unable to complete operation due to concurrent modification. Please try again later.
		// Handle them all globally for the service client.
		if tfawserr.ErrMessageContains(r.Error, apigatewayv2.ErrCodeConflictException, "try again later") {
			r.Retryable = aws.Bool(true)
		}
	})

	// Workaround for https://github.com/aws/aws-sdk-go/issues/1472
	client.applicationautoscalingConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if !strings.HasPrefix(r.Operation.Name, "Describe") && !strings.HasPrefix(r.Operation.Name, "List") {
			return
		}
		if tfawserr.ErrCodeEquals(r.Error, applicationautoscaling.ErrCodeFailedResourceAccessException) {
			r.Retryable = aws.Bool(true)
		}
	})

	// StartDeployment operations can return a ConflictException
	// if ongoing deployments are in-progress, thus we handle them
	// here for the service client.
	client.appconfigConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "StartDeployment" {
			if tfawserr.ErrCodeEquals(r.Error, appconfig.ErrCodeConflictException) {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.appsyncConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "CreateGraphqlApi" {
			if tfawserr.ErrMessageContains(r.Error, appsync.ErrCodeConcurrentModificationException, "a GraphQL API creation is already in progress") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.chimeConn.Handlers.Retry.PushBack(func(r *request.Request) {
		// When calling CreateVoiceConnector across multiple resources,
		// the API can randomly return a BadRequestException without explanation
		if r.Operation.Name == "CreateVoiceConnector" {
			if tfawserr.ErrMessageContains(r.Error, chime.ErrCodeBadRequestException, "Service received a bad request") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.cloudhsmv2Conn.Handlers.Retry.PushBack(func(r *request.Request) {
		if tfawserr.ErrMessageContains(r.Error, cloudhsmv2.ErrCodeCloudHsmInternalFailureException, "request was rejected because of an AWS CloudHSM internal failure") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.configserviceConn.Handlers.Retry.PushBack(func(r *request.Request) {
		// When calling Config Organization Rules API actions immediately
		// after Organization creation, the API can randomly return the
		// OrganizationAccessDeniedException error for a few minutes, even
		// after succeeding a few requests.
		switch r.Operation.Name {
		case "DeleteOrganizationConfigRule", "DescribeOrganizationConfigRules", "DescribeOrganizationConfigRuleStatuses", "PutOrganizationConfigRule":
			if !tfawserr.ErrMessageContains(r.Error, configservice.ErrCodeOrganizationAccessDeniedException, "This action can be only made by AWS Organization's master account.") {
				return
			}

			// We only want to retry briefly as the default max retry count would
			// excessively retry when the error could be legitimate.
			// We currently depend on the DefaultRetryer exponential backoff here.
			// ~10 retries gives a fair backoff of a few seconds.
			if r.RetryCount < 9 {
				r.Retryable = aws.Bool(true)
			} else {
				r.Retryable = aws.Bool(false)
			}
		case "DeleteOrganizationConformancePack", "DescribeOrganizationConformancePacks", "DescribeOrganizationConformancePackStatuses", "PutOrganizationConformancePack":
			if !tfawserr.ErrCodeEquals(r.Error, configservice.ErrCodeOrganizationAccessDeniedException) {
				if r.Operation.Name == "DeleteOrganizationConformancePack" && tfawserr.ErrCodeEquals(err, configservice.ErrCodeResourceInUseException) {
					r.Retryable = aws.Bool(true)
				}
				return
			}

			// We only want to retry briefly as the default max retry count would
			// excessively retry when the error could be legitimate.
			// We currently depend on the DefaultRetryer exponential backoff here.
			// ~10 retries gives a fair backoff of a few seconds.
			if r.RetryCount < 9 {
				r.Retryable = aws.Bool(true)
			} else {
				r.Retryable = aws.Bool(false)
			}
		}
	})

	client.cloudformationConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if tfawserr.ErrMessageContains(r.Error, cloudformation.ErrCodeOperationInProgressException, "Another Operation on StackSet") {
			r.Retryable = aws.Bool(true)
		}
	})

	// See https://github.com/aws/aws-sdk-go/pull/1276
	client.dynamodbConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name != "PutItem" && r.Operation.Name != "UpdateItem" && r.Operation.Name != "DeleteItem" {
			return
		}
		if tfawserr.ErrMessageContains(r.Error, dynamodb.ErrCodeLimitExceededException, "Subscriber limit exceeded:") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.ec2Conn.Handlers.Retry.PushBack(func(r *request.Request) {
		switch err := r.Error; r.Operation.Name {
		case "AttachVpnGateway", "DetachVpnGateway":
			if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "This call cannot be completed because there are pending VPNs or Virtual Interfaces") {
				r.Retryable = aws.Bool(true)
			}

		case "CreateClientVpnEndpoint":
			if tfawserr.ErrMessageContains(err, "OperationNotPermitted", "Endpoint cannot be created while another endpoint is being created") {
				r.Retryable = aws.Bool(true)
			}

		case "CreateClientVpnRoute", "DeleteClientVpnRoute":
			if tfawserr.ErrMessageContains(err, "ConcurrentMutationLimitExceeded", "Cannot initiate another change for this endpoint at this time") {
				r.Retryable = aws.Bool(true)
			}

		case "CreateVpnConnection":
			if tfawserr.ErrMessageContains(err, "VpnConnectionLimitExceeded", "maximum number of mutating objects has been reached") {
				r.Retryable = aws.Bool(true)
			}

		case "CreateVpnGateway":
			if tfawserr.ErrMessageContains(err, "VpnGatewayLimitExceeded", "maximum number of mutating objects has been reached") {
				r.Retryable = aws.Bool(true)
			}

		case "RunInstances":
			// `InsufficientInstanceCapacity` error has status code 500 and AWS SDK try retry this error by default.
			if tfawserr.ErrCodeEquals(err, "InsufficientInstanceCapacity") {
				r.Retryable = aws.Bool(false)
			}
		}
	})

	client.fmsConn.Handlers.Retry.PushBack(func(r *request.Request) {
		// Acceptance testing creates and deletes resources in quick succession.
		// The FMS onboarding process into Organizations is opaque to consumers.
		// Since we cannot reasonably check this status before receiving the error,
		// set the operation as retryable.
		switch r.Operation.Name {
		case "AssociateAdminAccount":
			if tfawserr.ErrMessageContains(r.Error, fms.ErrCodeInvalidOperationException, "Your AWS Organization is currently offboarding with AWS Firewall Manager. Please submit onboard request after offboarded.") {
				r.Retryable = aws.Bool(true)
			}
		case "DisassociateAdminAccount":
			if tfawserr.ErrMessageContains(r.Error, fms.ErrCodeInvalidOperationException, "Your AWS Organization is currently onboarding with AWS Firewall Manager and cannot be offboarded.") {
				r.Retryable = aws.Bool(true)
			}
		// System problems can arise during FMS policy updates (maybe also creation),
		// so we set the following operation as retryable.
		// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/23946
		case "PutPolicy":
			if tfawserr.ErrCodeEquals(r.Error, fms.ErrCodeInternalErrorException) {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.kafkaConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if tfawserr.ErrMessageContains(r.Error, kafka.ErrCodeTooManyRequestsException, "Too Many Requests") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.kinesisConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "CreateStream" {
			if tfawserr.ErrMessageContains(r.Error, kinesis.ErrCodeLimitExceededException, "simultaneously be in CREATING or DELETING") {
				r.Retryable = aws.Bool(true)
			}
		}
		if r.Operation.Name == "CreateStream" || r.Operation.Name == "DeleteStream" {
			if tfawserr.ErrMessageContains(r.Error, kinesis.ErrCodeLimitExceededException, "Rate exceeded for stream") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.organizationsConn.Handlers.Retry.PushBack(func(r *request.Request) {
		// Retry on the following error:
		// ConcurrentModificationException: AWS Organizations can't complete your request because it conflicts with another attempt to modify the same entity. Try again later.
		if tfawserr.ErrMessageContains(r.Error, organizations.ErrCodeConcurrentModificationException, "Try again later") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.s3Conn.Handlers.Retry.PushBack(func(r *request.Request) {
		if tfawserr.ErrMessageContains(r.Error, "OperationAborted", "A conflicting conditional operation is currently in progress against this resource. Please try again.") {
			r.Retryable = aws.Bool(true)
		}
	})

	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/17996
	client.securityhubConn.Handlers.Retry.PushBack(func(r *request.Request) {
		switch r.Operation.Name {
		case "EnableOrganizationAdminAccount":
			if tfawserr.ErrCodeEquals(r.Error, securityhub.ErrCodeResourceConflictException) {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/19215
	client.ssoadminConn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "AttachManagedPolicyToPermissionSet" || r.Operation.Name == "DetachManagedPolicyFromPermissionSet" {
			if tfawserr.ErrCodeEquals(r.Error, ssoadmin.ErrCodeConflictException) {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.storagegatewayConn.Handlers.Retry.PushBack(func(r *request.Request) {
		// InvalidGatewayRequestException: The specified gateway proxy network connection is busy.
		if tfawserr.ErrMessageContains(r.Error, storagegateway.ErrCodeInvalidGatewayRequestException, "The specified gateway proxy network connection is busy") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.wafv2Conn.Handlers.Retry.PushBack(func(r *request.Request) {
		if tfawserr.ErrMessageContains(r.Error, wafv2.ErrCodeWAFInternalErrorException, "Retry your request") {
			r.Retryable = aws.Bool(true)
		}

		if tfawserr.ErrMessageContains(r.Error, wafv2.ErrCodeWAFServiceLinkedRoleErrorException, "Retry") {
			r.Retryable = aws.Bool(true)
		}

		if r.Operation.Name == "CreateIPSet" || r.Operation.Name == "CreateRegexPatternSet" ||
			r.Operation.Name == "CreateRuleGroup" || r.Operation.Name == "CreateWebACL" {
			// WAFv2 supports tag on create which can result in the below error codes according to the documentation
			if tfawserr.ErrMessageContains(r.Error, wafv2.ErrCodeWAFTagOperationException, "Retry your request") {
				r.Retryable = aws.Bool(true)
			}
			if tfawserr.ErrMessageContains(err, wafv2.ErrCodeWAFTagOperationInternalErrorException, "Retry your request") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	// AWS SDK for Go v2 custom API clients.

	client.route53domainsClient = route53domains.NewFromConfig(cfg, func(o *route53domains.Options) {
		if endpoint := c.Endpoints[names.Route53Domains]; endpoint != "" {
			o.EndpointResolver = route53domains.EndpointResolverFromURL(endpoint)
		} else if partition == endpoints.AwsPartitionID {
			// Route 53 Domains is only available in AWS Commercial us-east-1 Region.
			o.Region = endpoints.UsEast1RegionID
		}
	})

	client.lightsailClient = lightsail.NewFromConfig(cfg, func(o *lightsail.Options) {
		lightsailRetryable := retry.IsErrorRetryableFunc(func(e error) awsv2.Ternary {
			if strings.Contains(e.Error(), "Please try again in a few minutes") || strings.Contains(e.Error(), "Please wait for it to complete before trying again") {
				return awsv2.TrueTernary
			}
			return awsv2.UnknownTernary
		})

		if endpoint := c.Endpoints[names.Lightsail]; endpoint != "" {
			o.EndpointResolver = lightsail.EndpointResolverFromURL(endpoint)
		}

		o.Retryer = retry.NewStandard(func(options *retry.StandardOptions) {
			options.Retryables = append(options.Retryables, lightsailRetryable)
			options.MaxAttempts = 18
			options.Backoff = retry.NewExponentialJitterBackoff(time.Second * tenBackOff)
		})
	})

	return client, nil
}
