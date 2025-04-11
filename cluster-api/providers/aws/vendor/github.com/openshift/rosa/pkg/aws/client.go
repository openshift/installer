/*
Copyright (c) 2020 Red Hat, Inc.

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

package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	awserr "github.com/openshift-online/ocm-common/pkg/aws/errors"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/sirupsen/logrus"
	"github.com/zgalor/weberr"

	client "github.com/openshift/rosa/pkg/aws/api_interface"
	"github.com/openshift/rosa/pkg/aws/profile"
	regionflag "github.com/openshift/rosa/pkg/aws/region"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/fedramp"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/info"
	"github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/reporter"
)

var (
	allErrorCodes      []string
	throttleErrorCodes []string
)

// Name of the AWS user that will be used to create all the resources of the cluster:
//
//go:generate mockgen -source=client.go -package=aws -destination=client_mock.go
const (
	AdminUserName        = "osdCcsAdmin"
	OsdCcsAdminStackName = "osdCcsAdminIAMUser"

	// Since CloudFormation stacks are region-dependent, we hard-code OCM's default region and
	// then use it to ensure that the user always gets the stack from the same region.
	DefaultRegion = "us-east-1"
	Inline        = "inline"
	Attached      = "attached"

	LocalZone      = "local-zone"
	WavelengthZone = "wavelength-zone"

	govPartition = "aws-us-gov"

	awsMaxFilterLength = 200

	numMaxRetries    = 12
	minRetryDelay    = 1 * time.Second
	minThrottleDelay = 5 * time.Second
	maxThrottleDelay = 5 * time.Second

	IAMServiceRegion = "us-east-1"
)

// Client defines a client interface
type Client interface {
	CheckAdminUserNotExisting(userName string) (err error)
	CheckAdminUserExists(userName string) (err error)
	CheckStackReadyOrNotExisting(stackName string) (stackReady bool, stackStatus *string, err error)
	CheckRoleExists(roleName string) (bool, string, error)
	ValidateRoleARNAccountIDMatchCallerAccountID(roleARN string) error
	GetIAMCredentials() (aws.Credentials, error)
	GetRegion() string
	ValidateCredentials() (isValid bool, err error)
	EnsureOsdCcsAdminUser(stackName string, adminUserName string, awsRegion string) (bool, error)
	DeleteOsdCcsAdminUser(stackName string) error
	AccessKeyGetter
	GetCreator() (*Creator, error)
	ValidateSCP(*string, map[string]*cmv1.AWSSTSPolicy) (bool, error)
	ListSubnets(subnetIds ...string) ([]ec2types.Subnet, error)
	GetSubnetAvailabilityZone(subnetID string) (string, error)
	GetAvailabilityZoneType(availabilityZoneName string) (string, error)
	GetVPCSubnets(subnetID string) ([]ec2types.Subnet, error)
	GetVPCPrivateSubnets(subnetID string) ([]ec2types.Subnet, error)
	FilterVPCsPrivateSubnets(subnets []ec2types.Subnet) ([]ec2types.Subnet, error)
	ValidateQuota() (bool, error)
	TagUserRegion(username string, region string) error
	GetClusterRegionTagForUser(username string) (string, error)
	EnsureRole(reporter *reporter.Object, name string, policy string, permissionsBoundary string,
		version string, tagList map[string]string, path string, managedPolicies bool) (string, error)
	ValidateRoleNameAvailable(name string) (err error)
	PutRolePolicy(roleName string, policyName string, policy string) error
	ForceEnsurePolicy(policyArn string, document string, version string, tagList map[string]string,
		path string) (string, error)
	EnsurePolicy(policyArn string, document string, version string, tagList map[string]string,
		path string) (string, error)
	AttachRolePolicy(reporter *reporter.Object, roleName string, policyARN string) error
	CreateOpenIDConnectProvider(issuerURL string, thumbprint string, clusterID string) (string, error)
	DeleteOpenIDConnectProvider(providerURL string) error
	HasOpenIDConnectProvider(issuerURL string, partition string, accountID string) (bool, error)
	FindRoleARNs(roleType string, version string) ([]string, error)
	FindRoleARNsClassic(roleType string, version string) ([]string, error)
	FindRoleARNsHostedCp(roleType string, version string) ([]string, error)
	FindPolicyARN(operator Operator, version string) (string, error)
	ListUserRoles() ([]Role, error)
	ListOCMRoles() ([]Role, error)
	ListAccountRoles(version string) ([]Role, error)
	ListOperatorRoles(version string, clusterID string, prefix string) (map[string][]OperatorRoleDetail, error)
	ListAttachedRolePolicies(roleName string) ([]string, error)
	ListOidcProviders(targetClusterId string, config *cmv1.OidcConfig) ([]OidcProviderOutput, error)
	GetRoleByARN(roleARN string) (iamtypes.Role, error)
	GetRoleByName(roleName string) (iamtypes.Role, error)
	DeleteOperatorRole(roles string, managedPolicies bool) error
	GetOperatorRolesFromAccountByClusterID(
		clusterID string,
		credRequests map[string]*cmv1.STSOperator,
	) ([]string, error)
	GetOperatorRolesFromAccountByPrefix(prefix string, credRequest map[string]*cmv1.STSOperator) ([]string, error)
	GetOperatorRolePolicies(roles []string) (map[string][]string, map[string][]string, error)
	GetAccountRolesForCurrentEnv(env string, accountID string) ([]Role, error)
	GetAccountRoleForCurrentEnv(env string, roleName string) (Role, error)
	GetAccountRoleForCurrentEnvWithPrefix(env string, rolePrefix string,
		accountRolesMap map[string]AccountRole) ([]Role, error)
	DeleteAccountRole(roleName string, prefix string, managedPolicies bool) error
	DeleteOCMRole(roleARN string, managedPolicies bool) error
	DeleteUserRole(roleName string) error
	GetAccountRolePolicies(roles []string, prefix string) (map[string][]PolicyDetail, map[string][]PolicyDetail, error)
	GetAttachedPolicy(role *string) ([]PolicyDetail, error)
	HasPermissionsBoundary(roleName string) (bool, error)
	GetOpenIDConnectProviderByClusterIdTag(clusterID string) (string, error)
	GetOpenIDConnectProviderByOidcEndpointUrl(oidcEndpointUrl string) (string, error)
	GetInstanceProfilesForRole(role string) ([]string, error)
	IsUpgradedNeededForAccountRolePolicies(rolePrefix string, version string) (bool, error)
	IsUpgradedNeededForAccountRolePoliciesUsingCluster(clusterID *cmv1.Cluster, version string) (bool, error)
	IsUpgradedNeededForOperatorRolePoliciesUsingCluster(
		cluster *cmv1.Cluster,
		partition string,
		accountID string,
		version string,
		credRequests map[string]*cmv1.STSOperator,
		operatorRolePolicyPrefix string,
	) (bool, error)
	IsUpgradedNeededForOperatorRolePoliciesUsingPrefix(
		rolePrefix string,
		partition string,
		accountID string,
		version string,
		credRequests map[string]*cmv1.STSOperator,
		path string,
	) (bool, error)
	UpdateTag(roleName string, defaultPolicyVersion string) error
	AddRoleTag(roleName string, key string, value string) error
	IsPolicyCompatible(policyArn string, version string) (bool, error)
	GetAccountRoleVersion(roleName string) (string, error)
	IsPolicyExists(policyARN string) (*iam.GetPolicyOutput, error)
	IsRolePolicyExists(roleName string, policyName string) (*iam.GetRolePolicyOutput, error)
	IsAdminRole(roleName string) (bool, error)
	DeleteInlineRolePolicies(roleName string) error
	IsUserRole(roleName *string) (bool, error)
	GetRoleARNPath(prefix string) (string, error)
	DescribeAvailabilityZones() ([]string, error)
	IsLocalAvailabilityZone(availabilityZoneName string) (bool, error)
	DetachRolePolicies(roleName string) error
	DetachRolePolicy(policyArn string, roleName string) error
	HasManagedPolicies(roleARN string) (bool, error)
	HasHostedCPPolicies(roleARN string) (bool, error)
	GetAccountRoleARN(prefix string, roleType string) (string, error)
	ValidateAccountRolesManagedPolicies(prefix string, policies map[string]*cmv1.AWSSTSPolicy) error
	ValidateHCPAccountRolesManagedPolicies(prefix string, policies map[string]*cmv1.AWSSTSPolicy) error
	ValidateOperatorRolesManagedPolicies(cluster *cmv1.Cluster, operatorRoles map[string]*cmv1.STSOperator,
		policies map[string]*cmv1.AWSSTSPolicy, hostedCPPolicies bool) error
	CreateS3Bucket(bucketName string, region string) error
	DeleteS3Bucket(bucketName string) error
	PutPublicReadObjectInS3Bucket(bucketName string, body io.ReadSeeker, key string) error
	CreateSecretInSecretsManager(name string, secret string) (string, error)
	DeleteSecretInSecretsManager(secretArn string) error
	ValidateAccountRoleVersionCompatibility(roleName string, roleType string, minVersion string) (bool, error)
	GetDefaultPolicyDocument(policyArn string) (string, error)
	GetAccountRoleByArn(roleArn string) (Role, error)
	GetSecurityGroupIds(vpcId string) ([]ec2types.SecurityGroup, error)
	FetchPublicSubnetMap(subnets []ec2types.Subnet) (map[string]bool, error)
	GetIAMServiceQuota(quotaCode string) (*servicequotas.GetServiceQuotaOutput, error)
	GetAccountRoleDefaultPolicy(roleName string, prefix string) (string, error)
	GetOperatorRoleDefaultPolicy(roleName string) (string, error)
}

type AccessKeyGetter interface {
	GetAWSAccessKeys() (*AccessKey, error)
	GetLocalAWSAccessKeys() (*AccessKey, error)
}

// ClientBuilder contains the information and logic needed to build a new AWS client.
type ClientBuilder struct {
	logger              *logrus.Logger
	region              *string
	credentials         *AccessKey
	useLocalCredentials bool
}

type awsClient struct {
	cfg                 aws.Config
	logger              *logrus.Logger
	iamClient           client.IamApiClient
	ec2Client           client.Ec2ApiClient
	orgClient           client.OrganizationsApiClient
	s3Client            client.S3ApiClient
	smClient            client.SecretsManagerApiClient
	stsClient           client.StsApiClient
	cfClient            client.CloudFormationApiClient
	serviceQuotasClient client.ServiceQuotasApiClient
	iamQuotaClient      client.ServiceQuotasApiClient
	awsAccessKeys       *AccessKey
	useLocalCredentials bool
}

func CreateNewClientOrExit(logger *logrus.Logger, reporter *reporter.Object) Client {
	awsClient, err := NewClient().
		Logger(logger).
		Build()
	if err != nil {
		reporter.Errorf("Failed to create AWS client: %v", err)
		os.Exit(1)
	}

	return awsClient
}

// NewClient creates a builder that can then be used to configure and build a new AWS client.
func NewClient() *ClientBuilder {
	return &ClientBuilder{}
}

func New(
	cfg aws.Config,
	logger *logrus.Logger,
	iamClient client.IamApiClient,
	ec2Client client.Ec2ApiClient,
	orgClient client.OrganizationsApiClient,
	s3Client client.S3ApiClient,
	smClient client.SecretsManagerApiClient,
	stsClient client.StsApiClient,
	cfClient client.CloudFormationApiClient,
	serviceQuotasClient client.ServiceQuotasApiClient,
	iamQuotaClient client.ServiceQuotasApiClient,
	awsAccessKeys *AccessKey,
	useLocalCredentials bool,

) Client {
	return &awsClient{
		cfg,
		logger,
		iamClient,
		ec2Client,
		orgClient,
		s3Client,
		smClient,
		stsClient,
		cfClient,
		serviceQuotasClient,
		iamQuotaClient,
		awsAccessKeys,
		useLocalCredentials,
	}
}

// Logger sets the logger that the AWS client will use to send messages to the log.
func (b *ClientBuilder) Logger(value *logrus.Logger) *ClientBuilder {
	b.logger = value
	return b
}

func (b *ClientBuilder) Region(value string) *ClientBuilder {
	b.region = aws.String(value)
	return b
}

func (b *ClientBuilder) AccessKeys(value *AccessKey) *ClientBuilder {
	b.credentials = value
	return b
}

func (b *ClientBuilder) UseLocalCredentials(value bool) *ClientBuilder {
	b.useLocalCredentials = value
	return b
}

// Create AWS session with a specific set of credentials
func (b *ClientBuilder) BuildSessionWithOptionsCredentials(value *AccessKey,
	logLevel aws.ClientLogMode) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(value.AccessKeyID,
			value.SecretAccessKey, "")),
		config.WithRegion(*b.region),
		config.WithHTTPClient(&http.Client{
			Transport: http.DefaultTransport,
		}),
		config.WithClientLogMode(logLevel),
		config.WithAPIOptions([]func(stack *middleware.Stack) error{
			smithyhttp.AddHeaderValue("User-Agent",
				strings.Join([]string{info.DefaultUserAgent, info.DefaultVersion}, ";")),
		}),
		config.WithRetryer(func() aws.Retryer {
			retryer := retry.AddWithMaxAttempts(retry.NewStandard(), numMaxRetries)
			retryer = retry.AddWithMaxBackoffDelay(retryer, time.Second)

			for _, code := range allErrorCodes {
				retryer = retry.AddWithErrorCodes(retryer, code)
			}

			return retryer
		}),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func (b *ClientBuilder) BuildSessionWithOptions(logLevel aws.ClientLogMode) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile.Profile()),
		config.WithRegion(*b.region),
		config.WithHTTPClient(awshttp.NewBuildableClient().WithTransportOptions()),
		config.WithClientLogMode(logLevel),
		config.WithAPIOptions([]func(stack *middleware.Stack) error{
			smithyhttp.AddHeaderValue("User-Agent",
				strings.Join([]string{info.DefaultUserAgent, info.DefaultVersion}, ";")),
		}),
		config.WithRetryer(func() aws.Retryer {
			retryer := retry.AddWithMaxAttempts(retry.NewStandard(), numMaxRetries)
			retryer = retry.AddWithMaxBackoffDelay(retryer, time.Second)

			for _, code := range allErrorCodes {
				retryer = retry.AddWithErrorCodes(retryer, code)
			}

			return retryer
		}),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func (b *ClientBuilder) BuildSession() (aws.Config, error) {
	var logLevel aws.ClientLogMode
	logLevel = 0
	if b.logger.Level == logrus.DebugLevel {
		logLevel = aws.LogRequestWithBody | aws.LogResponseWithBody
	}
	// Convert the map to a slice of strings.
	for code := range retry.DefaultThrottleErrorCodes {
		throttleErrorCodes = append(throttleErrorCodes, code)
	}
	allErrorCodes = append(throttleErrorCodes, awserr.InvalidClientTokenID)

	if b.credentials != nil {
		return b.BuildSessionWithOptionsCredentials(b.credentials, logLevel)
	}

	return b.BuildSessionWithOptions(logLevel)
}

// Build uses the information stored in the builder to build a new AWS client.
func (b *ClientBuilder) Build() (Client, error) {
	// Check parameters:
	if b.logger == nil {
		return nil, fmt.Errorf("logger is mandatory")
	}

	if b.region == nil || *b.region == "" {
		region, err := GetRegion(regionflag.Region())
		if err != nil {
			return nil, err
		}
		b.region = aws.String(region)
	}

	if fedramp.IsGovRegion(*b.region) {
		fedramp.Enable()
	} else if fedramp.Enabled() {
		return nil, fmt.Errorf("failed to connect to AWS. Use a GovCloud region in your profile")
	}

	// Create the AWS session:
	cfg, err := b.BuildSession()
	if err != nil {
		return nil, err
	}

	if cfg.Region == "" {
		return nil, fmt.Errorf("region is not set. Use --region to set the region")
	}

	if profile.Profile() != "" {
		b.logger.Debugf("Using AWS profile: %s", profile.Profile())
	}

	// IAM Service is only available in "us-east-1", need to create specific config for it
	iamCfg, err := b.BuildSession()
	if err != nil {
		return nil, err
	}
	iamCfg.Region = IAMServiceRegion

	// Create and populate the object:
	c := &awsClient{
		cfg:                 cfg,
		logger:              b.logger,
		iamClient:           iam.NewFromConfig(cfg),
		ec2Client:           ec2.NewFromConfig(cfg),
		orgClient:           organizations.NewFromConfig(cfg),
		s3Client:            s3.NewFromConfig(cfg),
		smClient:            secretsmanager.NewFromConfig(cfg),
		stsClient:           sts.NewFromConfig(cfg),
		cfClient:            cloudformation.NewFromConfig(cfg),
		serviceQuotasClient: servicequotas.NewFromConfig(cfg),
		iamQuotaClient:      servicequotas.NewFromConfig(iamCfg),
		useLocalCredentials: b.useLocalCredentials,
	}

	_, root, err := getClientDetails(c)
	if err != nil {
		return nil, err
	}

	if root {
		return nil, errors.New("using a root account is not supported, please use an IAM user instead")
	}

	return c, err
}

func (c *awsClient) GetIAMCredentials() (aws.Credentials, error) {
	return c.cfg.Credentials.Retrieve(context.TODO())
}

func (c *awsClient) GetRegion() string {
	return c.cfg.Region
}

func (c *awsClient) FetchPublicSubnetMap(subnets []ec2types.Subnet) (map[string]bool, error) {
	mapSubnetIdToPublic := map[string]bool{}
	if len(subnets) == 0 {
		return mapSubnetIdToPublic, nil
	}
	// AWS has a limit of 200 filters per query so it needs to be broken up in chunks
	chunks := helper.ChunkSlice(subnets, awsMaxFilterLength)

	for _, curChunk := range chunks {
		subnetIds := []*string{}
		for _, subnet := range curChunk {
			subnetIds = append(subnetIds, subnet.SubnetId)
		}
		routeTablesResp, err := c.ec2Client.DescribeRouteTables(context.Background(), &ec2.DescribeRouteTablesInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("association.subnet-id"),
					Values: aws.ToStringSlice(subnetIds),
				},
			},
		})
		if err != nil {
			return mapSubnetIdToPublic, err
		}
		if routeTablesResp == nil {
			return mapSubnetIdToPublic, fmt.Errorf(
				"No route table found for associated subnets '%s'",
				helper.SliceToSortedString(aws.ToStringSlice(subnetIds)),
			)
		}
		for _, routes := range routeTablesResp.RouteTables {
			for _, association := range routes.Associations {
				subnetAssociation := aws.ToString(association.SubnetId)
				mapSubnetIdToPublic[subnetAssociation] = false
				for _, route := range routes.Routes {
					if strings.HasPrefix(aws.ToString(route.GatewayId), "igw") {
						// There is no direct way in the AWS API to determine if a subnet is public or private.
						// A public subnet is one which has an internet gateway route
						// we look for the gatewayId and make sure it has the prefix of igw to differentiate
						// from the default in-subnet route which is called "local"
						// or other virtual gateway (starting with vgv)
						// or vpc peering connections (starting with pcx).
						mapSubnetIdToPublic[subnetAssociation] = true
					}
				}
			}
		}
	}
	return mapSubnetIdToPublic, nil
}

func (c *awsClient) ListSubnets(subnetIds ...string) ([]ec2types.Subnet, error) {

	if len(subnetIds) == 0 {
		return c.getSubnetIDs(&ec2.DescribeSubnetsInput{})
	}

	var ids []string

	ids = append(ids, subnetIds...)

	return c.getSubnetIDs(&ec2.DescribeSubnetsInput{
		SubnetIds: ids,
	})
}

func (c *awsClient) GetSubnetAvailabilityZone(subnetID string) (string, error) {
	res, err := c.ec2Client.DescribeSubnets(
		context.Background(),
		&ec2.DescribeSubnetsInput{SubnetIds: []string{subnetID}},
	)
	if err != nil {
		return "", err
	}
	if len(res.Subnets) < 1 {
		return "", fmt.Errorf("failed to get subnet with ID '%s'", subnetID)
	}

	return *res.Subnets[0].AvailabilityZone, nil
}

func (c *awsClient) GetVPCPrivateSubnets(subnetID string) ([]ec2types.Subnet, error) {
	subnets, err := c.GetVPCSubnets(subnetID)
	if err != nil {
		return nil, err
	}

	return c.FilterVPCsPrivateSubnets(subnets)
}

// getVPCSubnets gets a subnet ID and fetches all the subnets that belong to the same VPC as the provided subnet.
func (c *awsClient) GetVPCSubnets(subnetID string) ([]ec2types.Subnet, error) {
	// Fetch the subnet details
	subnets, err := c.getSubnetIDs(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: []string{subnetID},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(subnets) < 1 {
		return nil, fmt.Errorf("failed to get subnet with ID '%s'", subnetID)
	}

	// Fetch VPC's subnets
	vpcID := subnets[0].VpcId
	subnets, err = c.getSubnetIDs(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{*vpcID},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(subnets) < 1 {
		return nil, fmt.Errorf("failed to get the subnets of VPC with ID '%s'", *vpcID)
	}

	return subnets, nil
}

// FilterPrivateSubnets gets a slice of subnets that belongs to the same VPC and filters the private subnets.
// Assumption: subnets - non-empty slice.
func (c *awsClient) FilterVPCsPrivateSubnets(subnets []ec2types.Subnet) ([]ec2types.Subnet, error) {
	// Fetch VPC route tables
	vpcID := subnets[0].VpcId
	describeRouteTablesOutput, err := c.ec2Client.DescribeRouteTables(
		context.Background(),
		&ec2.DescribeRouteTablesInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{*vpcID},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(describeRouteTablesOutput.RouteTables) < 1 {
		return nil, fmt.Errorf("failed to find VPC '%s' route table", *vpcID)
	}

	var privateSubnets []ec2types.Subnet
	for _, subnet := range subnets {
		isPublic, err := c.isPublicSubnet(subnet.SubnetId, describeRouteTablesOutput.RouteTables)
		if err != nil {
			return nil, err
		}
		if !isPublic {
			privateSubnets = append(privateSubnets, subnet)
		}
	}

	if len(privateSubnets) < 1 {
		return nil, fmt.Errorf("failed to find private subnets associated with VPC '%s'", *subnets[0].VpcId)
	}

	return privateSubnets, nil
}

// isPublicSubnet a public subnet is a subnet that's associated with a route table that has a route to an
// internet gateway
func (c *awsClient) isPublicSubnet(subnetID *string, routeTables []ec2types.RouteTable) (bool, error) {
	subnetRouteTable, err := c.getSubnetRouteTable(subnetID, routeTables)
	if err != nil {
		return false, err
	}

	for _, route := range subnetRouteTable.Routes {
		if strings.Contains(aws.ToString(route.GatewayId), "igw") {
			return true, nil
		}
	}

	return false, nil
}

func (c *awsClient) getSubnetRouteTable(subnetID *string,
	routeTables []ec2types.RouteTable) (*ec2types.RouteTable, error) {
	// Subnet route table — A route table that's associated with a subnet
	for _, routeTable := range routeTables {
		for _, association := range routeTable.Associations {
			if aws.ToString(association.SubnetId) == aws.ToString(subnetID) {
				return &routeTable, nil
			}
		}
	}

	// A subnet can be explicitly associated with custom route table, or implicitly or explicitly associated with the
	// main route table.
	for _, routeTable := range routeTables {
		for _, association := range routeTable.Associations {
			if aws.ToBool(association.Main) {
				return &routeTable, nil
			}
		}
	}

	// Each subnet in the VPC must be associated with a route table
	return nil, fmt.Errorf("failed to find subnet '%s' route table", *subnetID)
}

// getSubnetIDs will return the list of subnetsIDs supported for the region picked.
// It is possible to pass non-empty `describeSubnetsInput` to filter results.
func (c *awsClient) getSubnetIDs(describeSubnetsInput *ec2.DescribeSubnetsInput) ([]ec2types.Subnet, error) {
	res, err := c.ec2Client.DescribeSubnets(context.Background(), describeSubnetsInput)
	if err != nil {
		return nil, err
	}

	return res.Subnets, nil
}

type Creator struct {
	ARN        string
	AccountID  string
	IsSTS      bool
	IsGovcloud bool
	Partition  string
}

func (c *awsClient) GetCreator() (*Creator, error) {
	getCallerIdentityOutput, err := c.stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return CreatorForCallerIdentity(getCallerIdentityOutput)
}

// CreatorForCallerIdentity adapts an STS CallerIdentity to the ROSA *Creator
func CreatorForCallerIdentity(identity *sts.GetCallerIdentityOutput) (*Creator, error) {
	creatorARN := aws.ToString(identity.Arn)

	// Extract the account identifier from the ARN of the user:
	creatorParsedARN, err := arn.Parse(creatorARN)
	if err != nil {
		return nil, err
	}

	// If the user is STS resolve the Role the user has assumed
	var stsRole *string
	if isSTS(creatorParsedARN) {
		stsRole, err = resolveSTSRole(creatorParsedARN)
		if err != nil {
			return nil, err
		}

		// resolveSTSRole ensures a parsed valid ARN before
		// returning it so we don't need to parse it again
		creatorARN = *stsRole
	}

	isGovcloud := creatorParsedARN.Partition == govPartition

	return &Creator{
		ARN:        creatorARN,
		AccountID:  creatorParsedARN.AccountID,
		IsSTS:      isSTS(creatorParsedARN),
		IsGovcloud: isGovcloud,
		Partition:  creatorParsedARN.Partition,
	}, nil
}

// Checks if given credentials are valid.
func (c *awsClient) ValidateCredentials() (bool, error) {
	// Validate the AWS credentials by calling STS GetCallerIdentity
	// This will fail if the AWS access key and secret key are invalid. This
	// will also work for STS credentials with access key, secret key and session
	// token
	_, err := c.stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "InvalidClientTokenId") {
			awsErr := fmt.Errorf("Invalid AWS Credentials: %s.\n For help configuring your credentials, see %s",
				err,
				"https://docs.openshift.com/rosa/rosa_install_access_delete_clusters/rosa_getting_started_iam/"+
					"rosa-config-aws-account.html#rosa-configuring-aws-account_rosa-config-aws-account")
			return false, awsErr
		}
		return false, err
	}

	return true, nil
}

func (c *awsClient) CheckAdminUserNotExisting(userName string) (err error) {
	userList, err := c.iamClient.ListUsers(context.Background(), &iam.ListUsersInput{})
	if err != nil {
		return err
	}
	for _, user := range userList.Users {
		if *user.UserName == userName {
			return fmt.Errorf("error creating user: IAM user '%s' already exists.\n"+
				"Ensure user '%s' IAM user does not exist, then retry with\n"+
				"rosa init",
				*user.UserName, *user.UserName)
		}
	}
	return nil
}

func (c *awsClient) CheckAdminUserExists(userName string) (err error) {
	_, err = c.iamClient.GetUser(context.Background(), &iam.GetUserInput{UserName: aws.String(userName)})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) GetClusterRegionTagForUser(username string) (string, error) {
	user, err := c.iamClient.GetUser(context.Background(), &iam.GetUserInput{UserName: aws.String(username)})
	if err != nil {
		return "", err
	}
	for _, tag := range user.User.Tags {
		if *tag.Key == tags.ClusterRegion {
			return *tag.Value, nil
		}
	}
	return "", nil
}

func (c *awsClient) TagUserRegion(username string, region string) error {
	_, err := c.iamClient.TagUser(context.Background(), &iam.TagUserInput{
		UserName: aws.String(username),
		Tags: []iamtypes.Tag{
			{
				Key:   aws.String(tags.ClusterRegion),
				Value: aws.String(region),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

type AccessKey struct {
	AccessKeyID     string
	SecretAccessKey string
}

// GetAWSAccessKeys uses UpsertAccessKey to delete and create new access keys
// for `osdCcsAdmin` each time we use the client to create a cluster.
// There is no need to permanently store these credentials since they are only used
// on create, the cluster uses a completely different set of IAM credentials
// provisioned by this user.
func (c *awsClient) GetAWSAccessKeys() (*AccessKey, error) {
	if c.awsAccessKeys != nil {
		return c.awsAccessKeys, nil
	}

	if c.useLocalCredentials {
		return c.GetLocalAWSAccessKeys()
	}

	accessKey, err := c.UpsertAccessKey(AdminUserName)
	if err != nil {
		return nil, err
	}

	err = c.ValidateAccessKeys(accessKey)
	if err != nil {
		return nil, err
	}

	c.awsAccessKeys = accessKey

	return c.awsAccessKeys, nil
}

func (c *awsClient) GetLocalAWSAccessKeys() (*AccessKey, error) {
	creds, err := c.cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		return nil, err
	}
	c.awsAccessKeys = &AccessKey{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
	}
	return c.awsAccessKeys, nil
}

// ValidateAccessKeys deals with AWS' eventual consistency, its attempts to call
// GetCallerIdentity and will try again if the error is access denied.
func (c *awsClient) ValidateAccessKeys(AccessKey *AccessKey) error {
	logger := logging.NewLogger()

	start := time.Now()
	maxAttempts := 15

	// Wait for credentials
	// 15 attempts should be enough, it takes generally around 10 seconds to ready
	// credentials
	for i := 0; i < maxAttempts; i++ {
		// Create the AWS client
		_, err := NewClient().
			Logger(logger).
			AccessKeys(AccessKey).
			Build()

		if err != nil {
			logger.Debug(fmt.Sprintf("%+v\n", err))
			if awserr.IsInvalidTokenException(err) {
				wait := time.Duration((i * 200)) * time.Millisecond
				waited := time.Since(start)
				logger.Debug(fmt.Printf("InvalidClientTokenId, waited %.2f\n", waited.Seconds()))
				time.Sleep(wait)
			}

			if awserr.IsAccessDeniedException(err) {
				wait := time.Duration((i * 200)) * time.Millisecond
				waited := time.Since(start)
				logger.Debug(fmt.Printf("AccessDenied, waited %.2f\n", waited.Seconds()))
				time.Sleep(wait)
			}

			if i == maxAttempts {
				logger.Error("Error waiting for IAM credentials to become ready")
				return err
			}
		} else {
			waited := time.Since(start)
			logger.Debug(fmt.Sprintf("\nCredentials ready in %.2fs\n", waited.Seconds()))
			break
		}
	}
	return nil
}

// UpsertAccessKey first deletes all access keys attached to `username` and then creates a
// new access key. DeleteAccessKey ensures we own the user before proceeding to delete
// access keys
func (c *awsClient) UpsertAccessKey(username string) (*AccessKey, error) {
	err := c.DeleteAccessKeys(username)
	if err != nil {
		return nil, err
	}

	createAccessKeyOutput, err := c.CreateAccessKey(username)
	if err != nil {
		return nil, err
	}

	return &AccessKey{
		AccessKeyID:     *createAccessKeyOutput.AccessKey.AccessKeyId,
		SecretAccessKey: *createAccessKeyOutput.AccessKey.SecretAccessKey,
	}, nil
}

// CreateAccessKey creates an IAM access key for `username`
func (c *awsClient) CreateAccessKey(username string) (*iam.CreateAccessKeyOutput, error) {
	// Create access key for IAM user
	createIAMUserAccessKeyOutput, err := c.iamClient.CreateAccessKey(context.Background(),
		&iam.CreateAccessKeyInput{
			UserName: aws.String(username),
		},
	)
	if err != nil {
		return nil, err
	}

	return createIAMUserAccessKeyOutput, nil
}

// DeleteAccessKeys deletes all access keys from `username`. We ensure
// that we own the user before deleting access keys by search for IAM Tags
func (c *awsClient) DeleteAccessKeys(username string) error {
	// List all access keys for user. Result wont be truncated since IAM users
	// can only have 2 access keys
	listAccessKeysOutput, err := c.iamClient.ListAccessKeys(context.Background(),
		&iam.ListAccessKeysInput{
			UserName: aws.String(username),
		},
	)
	if err != nil {
		return err
	}

	// Delete all access keys. Moactl owns this user since the CloudFormation stack
	// at this point is complete and the user is tagged by use on creation
	for _, key := range listAccessKeysOutput.AccessKeyMetadata {
		_, err = c.iamClient.DeleteAccessKey(context.Background(),
			&iam.DeleteAccessKeyInput{
				UserName:    aws.String(username),
				AccessKeyId: key.AccessKeyId,
			},
		)
		if err != nil {
			return err
		}
	}

	// Complete, deleted all accesskeys for `username`
	return nil
}

// CheckRoleExists checks to see if an IAM role with the same name
// already exists
func (c *awsClient) CheckRoleExists(roleName string) (bool, string, error) {
	role, err := c.iamClient.GetRole(context.Background(),
		&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
	if err != nil {
		if awserr.IsNoSuchEntityException(err) {
			return false, "", nil
		}
		return false, "", err
	}

	return true, aws.ToString(role.Role.Arn), nil
}

func (c *awsClient) GetRoleByARN(roleARN string) (iamtypes.Role, error) {
	// validate arn
	parsedARN, err := arn.Parse(roleARN)
	if err != nil {
		return iamtypes.Role{}, fmt.Errorf("expected '%s' to be a valid IAM role ARN: %s", roleARN, err)
	}

	// validate arn is for a role resource
	resource := parsedARN.Resource
	isRole := strings.Contains(resource, "role/")
	if !isRole {
		return iamtypes.Role{}, fmt.Errorf("expected ARN '%s' to be IAM role resource", roleARN)
	}

	// get resource name

	m := strings.LastIndex(resource, "/")
	roleName := resource[m+1:]

	return c.GetRoleByName(roleName)
}

func (c *awsClient) GetRoleByName(roleName string) (iamtypes.Role, error) {
	roleOutput, err := c.iamClient.GetRole(context.Background(),
		&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
	if err != nil {
		return iamtypes.Role{}, err
	}
	return *roleOutput.Role, nil
}

// DescribeAvailabilityZones fetches the region's availability zones with type `availability-zone`
func (c *awsClient) DescribeAvailabilityZones() ([]string, error) {
	describeAvailabilityZonesOutput, err := c.ec2Client.DescribeAvailabilityZones(context.Background(),
		&ec2.DescribeAvailabilityZonesInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("zone-type"),
					Values: []string{"availability-zone"},
				},
			},
		})
	if err != nil {
		return nil, err
	}

	var availabilityZones []string
	for _, az := range describeAvailabilityZonesOutput.AvailabilityZones {
		availabilityZones = append(availabilityZones, *az.ZoneName)
	}

	return availabilityZones, nil
}

func (c *awsClient) IsLocalAvailabilityZone(availabilityZoneName string) (bool, error) {
	availabilityZones, err := c.ec2Client.DescribeAvailabilityZones(context.Background(),
		&ec2.DescribeAvailabilityZonesInput{ZoneNames: []string{availabilityZoneName}})
	if err != nil {
		return false, err
	}
	if len(availabilityZones.AvailabilityZones) < 1 {
		return false, fmt.Errorf("failed to find availability zone '%s'", availabilityZoneName)
	}

	return aws.ToString(availabilityZones.AvailabilityZones[0].ZoneType) == LocalZone, nil
}

func (c *awsClient) GetAvailabilityZoneType(availabilityZoneName string) (string, error) {
	availabilityZones, err := c.ec2Client.DescribeAvailabilityZones(context.Background(),
		&ec2.DescribeAvailabilityZonesInput{ZoneNames: []string{availabilityZoneName}})
	if err != nil {
		return "", err
	}
	if len(availabilityZones.AvailabilityZones) < 1 {
		return "", fmt.Errorf("Failed to find availability zone '%s'", availabilityZoneName)
	}
	return aws.ToString(availabilityZones.AvailabilityZones[0].ZoneType), nil
}

func (c *awsClient) DetachRolePolicies(roleName string) error {
	attachedPolicies := make([]iamtypes.AttachedPolicy, 0)
	isTruncated := true
	var marker *string
	for isTruncated {
		resp, err := c.iamClient.ListAttachedRolePolicies(context.Background(),
			&iam.ListAttachedRolePoliciesInput{
				Marker:   marker,
				RoleName: &roleName,
			},
		)
		if err != nil {
			return err
		}
		isTruncated = resp.IsTruncated
		marker = resp.Marker
	}
	for _, attachedPolicy := range attachedPolicies {
		err := c.DetachRolePolicy(*attachedPolicy.PolicyArn, roleName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *awsClient) DetachRolePolicy(policyArn string, roleName string) error {
	_, err := c.iamClient.DetachRolePolicy(context.Background(),
		&iam.DetachRolePolicyInput{PolicyArn: &policyArn, RoleName: &roleName})
	if err != nil {
		return err
	}
	return nil
}

const ReadOnlyAnonUserPolicyTemplate = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "AllowReadPublicAccess",
			"Principal": "*",
			"Effect": "Allow",
			"Action": [
				"s3:GetObject"
			],
			"Resource": [
				"arn:aws:s3:::%s/*"
			]
		}
	]
}`

func (c *awsClient) CreateS3Bucket(bucketName string, region string) error {
	_, err := c.s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return weberr.Errorf("Bucket '%s' already exists.", bucketName)
	}

	bucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if region != DefaultRegion {
		bucketInput.CreateBucketConfiguration = &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(region),
		}
	}
	_, err = c.s3Client.CreateBucket(context.TODO(), bucketInput)
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutPublicAccessBlock(context.TODO(), &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(true),
			IgnorePublicAcls:      aws.Bool(true),
			BlockPublicPolicy:     aws.Bool(false),
			RestrictPublicBuckets: aws.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutBucketPolicy(context.TODO(), &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(fmt.Sprintf(ReadOnlyAnonUserPolicyTemplate, bucketName)),
	})
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutBucketTagging(context.TODO(), &s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3types.Tagging{
			TagSet: []s3types.Tag{
				{
					Key:   aws.String(tags.RedHatManaged),
					Value: aws.String(tags.True),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *awsClient) DeleteS3Bucket(bucketName string) error {
	_, err := c.s3Client.HeadBucket(context.Background(),
		&s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
	if err != nil {
		var notFound *s3types.NotFound
		if errors.As(err, &notFound) {
			return nil
		}
		return err
	}
	err = c.emptyS3Bucket(bucketName)
	if err != nil {
		return err
	}
	_, err = c.s3Client.DeleteBucket(context.Background(),
		&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) emptyS3Bucket(bucketName string) error {
	objects, err := c.s3Client.ListObjects(context.Background(),
		&s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
	if err != nil {
		return err
	}
	for _, object := range (*objects).Contents {
		_, err = c.s3Client.DeleteObject(context.Background(),
			&s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    object.Key,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *awsClient) PutPublicReadObjectInS3Bucket(bucketName string, body io.ReadSeeker, key string) error {
	_, err := c.s3Client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Body:    body,
			Bucket:  aws.String(bucketName),
			Key:     aws.String(key),
			Tagging: aws.String(fmt.Sprintf("%s=%s", tags.RedHatManaged, tags.True)),
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) CreateSecretInSecretsManager(name string, secret string) (string, error) {
	createSecretResponse, err := c.smClient.CreateSecret(context.Background(),
		&secretsmanager.CreateSecretInput{
			Description:  aws.String(fmt.Sprintf("Secret for %s", name)),
			Name:         aws.String(name),
			SecretString: aws.String(secret),
			Tags: []secretsmanagertypes.Tag{{
				Key:   aws.String(tags.RedHatManaged),
				Value: aws.String("true"),
			}},
		})
	if err != nil {
		return "", err
	}
	return *createSecretResponse.ARN, nil
}

func (c *awsClient) DeleteSecretInSecretsManager(secretArn string) error {
	_, err := c.smClient.DescribeSecret(context.Background(),
		&secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secretArn),
		})
	if err != nil {
		var resourceNotFound *secretsmanagertypes.ResourceNotFoundException
		if errors.As(err, &resourceNotFound) {
			return nil
		}
	}
	_, err = c.smClient.DeleteSecret(context.Background(),
		&secretsmanager.DeleteSecretInput{
			ForceDeleteWithoutRecovery: aws.Bool(true),
			SecretId:                   aws.String(secretArn),
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) GetSecurityGroupIds(vpcId string) ([]ec2types.SecurityGroup, error) {
	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcId},
			},
		},
	}
	resp, err := c.ec2Client.DescribeSecurityGroups(context.Background(), describeSecurityGroupsInput)
	if err != nil {
		return []ec2types.SecurityGroup{}, err
	}

	return resp.SecurityGroups, nil
}

func Ec2ResourceHasTag(tags []ec2types.Tag, tagName, tagValue string) bool {
	for _, tag := range tags {
		if aws.ToString(tag.Key) == tagName && aws.ToString(tag.Value) == tagValue {
			return true
		}
	}
	return false
}
