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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/google/uuid"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/sirupsen/logrus"
	"github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/aws/profile"
	regionflag "github.com/openshift/rosa/pkg/aws/region"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/fedramp"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/info"
	"github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/reporter"
)

// Name of the AWS user that will be used to create all the resources of the cluster:
//
//go:generate mockgen -source=client.go -package=aws -destination=mock_client.go
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
)

// addROSAVersionToUserAgent is a named handler that will add ROSA CLI
// version information to requests made by the AWS SDK.
var addROSAVersionToUserAgent = request.NamedHandler{
	Name: "rosa.ROSAVersionUserAgentHandler",
	Fn:   request.MakeAddToUserAgentHandler(info.UserAgent, info.Version),
}

// Client defines a client interface
type Client interface {
	CheckAdminUserNotExisting(userName string) (err error)
	CheckAdminUserExists(userName string) (err error)
	CheckStackReadyOrNotExisting(stackName string) (stackReady bool, stackStatus *string, err error)
	CheckRoleExists(roleName string) (bool, string, error)
	ValidateRoleARNAccountIDMatchCallerAccountID(roleARN string) error
	GetIAMCredentials() (credentials.Value, error)
	GetRegion() string
	ValidateCredentials() (isValid bool, err error)
	EnsureOsdCcsAdminUser(stackName string, adminUserName string, awsRegion string) (bool, error)
	DeleteOsdCcsAdminUser(stackName string) error
	AccessKeyGetter
	GetCreator() (*Creator, error)
	ValidateSCP(*string, map[string]*cmv1.AWSSTSPolicy) (bool, error)
	ListSubnets(subnetIds ...string) ([]*ec2.Subnet, error)
	GetSubnetAvailabilityZone(subnetID string) (string, error)
	GetAvailabilityZoneType(availabilityZoneName string) (string, error)
	GetVPCSubnets(subnetID string) ([]*ec2.Subnet, error)
	GetVPCPrivateSubnets(subnetID string) ([]*ec2.Subnet, error)
	FilterVPCsPrivateSubnets(subnets []*ec2.Subnet) ([]*ec2.Subnet, error)
	ValidateQuota() (bool, error)
	TagUserRegion(username string, region string) error
	GetClusterRegionTagForUser(username string) (string, error)
	EnsureRole(name string, policy string, permissionsBoundary string,
		version string, tagList map[string]string, path string, managedPolicies bool) (string, error)
	ValidateRoleNameAvailable(name string) (err error)
	PutRolePolicy(roleName string, policyName string, policy string) error
	ForceEnsurePolicy(policyArn string, document string, version string, tagList map[string]string,
		path string) (string, error)
	EnsurePolicy(policyArn string, document string, version string, tagList map[string]string,
		path string) (string, error)
	AttachRolePolicy(roleName string, policyARN string) error
	CreateOpenIDConnectProvider(issuerURL string, thumbprint string, clusterID string) (string, error)
	DeleteOpenIDConnectProvider(providerURL string) error
	HasOpenIDConnectProvider(issuerURL string, accountID string) (bool, error)
	FindRoleARNs(roleType string, version string) ([]string, error)
	FindRoleARNsClassic(roleType string, version string) ([]string, error)
	FindRoleARNsHostedCp(roleType string, version string) ([]string, error)
	FindPolicyARN(operator Operator, version string) (string, error)
	ListUserRoles() ([]Role, error)
	ListOCMRoles() ([]Role, error)
	ListAccountRoles(version string) ([]Role, error)
	ListOperatorRoles(version string, clusterID string) (map[string][]OperatorRoleDetail, error)
	ListOidcProviders(targetClusterId string, config *cmv1.OidcConfig) ([]OidcProviderOutput, error)
	GetRoleByARN(roleARN string) (*iam.Role, error)
	DeleteOperatorRole(roles string, managedPolicies bool) error
	GetOperatorRolesFromAccountByClusterID(
		clusterID string,
		credRequests map[string]*cmv1.STSOperator,
	) ([]string, error)
	GetOperatorRolesFromAccountByPrefix(prefix string, credRequest map[string]*cmv1.STSOperator) ([]string, error)
	GetPolicies(roles []string) (map[string][]string, error)
	GetAccountRolesForCurrentEnv(env string, accountID string) ([]Role, error)
	GetAccountRoleForCurrentEnv(env string, roleName string) (Role, error)
	GetAccountRoleForCurrentEnvWithPrefix(env string, rolePrefix string,
		accountRolesMap map[string]AccountRole) ([]Role, error)
	DeleteAccountRole(roleName string, managedPolicies bool) error
	DeleteOCMRole(roleARN string, managedPolicies bool) error
	DeleteUserRole(roleName string) error
	GetAccountRolePolicies(roles []string) (map[string][]PolicyDetail, error)
	GetAttachedPolicy(role *string) ([]PolicyDetail, error)
	HasPermissionsBoundary(roleName string) (bool, error)
	GetOpenIDConnectProviderByClusterIdTag(clusterID string) (string, error)
	GetOpenIDConnectProviderByOidcEndpointUrl(oidcEndpointUrl string) (string, error)
	GetInstanceProfilesForRole(role string) ([]string, error)
	IsUpgradedNeededForAccountRolePolicies(rolePrefix string, version string) (bool, error)
	IsUpgradedNeededForAccountRolePoliciesUsingCluster(clusterID *cmv1.Cluster, version string) (bool, error)
	IsUpgradedNeededForOperatorRolePoliciesUsingCluster(
		cluster *cmv1.Cluster,
		accountID string,
		version string,
		credRequests map[string]*cmv1.STSOperator,
		operatorRolePolicyPrefix string,
	) (bool, error)
	IsUpgradedNeededForOperatorRolePoliciesUsingPrefix(
		rolePrefix string,
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
	GetAccountRoleByArn(roleArn string) (*Role, error)
	GetSecurityGroupIds(vpcId string) ([]*ec2.SecurityGroup, error)
	FetchPublicSubnetMap(subnets []*ec2.Subnet) (map[string]bool, error)
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
	logger              *logrus.Logger
	iamClient           iamiface.IAMAPI
	ec2Client           ec2iface.EC2API
	orgClient           organizationsiface.OrganizationsAPI
	s3Client            s3iface.S3API
	smClient            secretsmanageriface.SecretsManagerAPI
	stsClient           stsiface.STSAPI
	cfClient            cloudformationiface.CloudFormationAPI
	servicequotasClient servicequotasiface.ServiceQuotasAPI
	awsSession          *session.Session
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
	logger *logrus.Logger,
	iamClient iamiface.IAMAPI,
	ec2Client ec2iface.EC2API,
	orgClient organizationsiface.OrganizationsAPI,
	s3Client s3iface.S3API,
	smClient secretsmanageriface.SecretsManagerAPI,
	stsClient stsiface.STSAPI,
	cfClient cloudformationiface.CloudFormationAPI,
	servicequotasClient servicequotasiface.ServiceQuotasAPI,
	awsSession *session.Session,
	awsAccessKeys *AccessKey,
	useLocalCredentials bool,

) Client {
	return &awsClient{
		logger,
		iamClient,
		ec2Client,
		orgClient,
		s3Client,
		smClient,
		stsClient,
		cfClient,
		servicequotasClient,
		awsSession,
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
func (b *ClientBuilder) BuildSessionWithOptionsCredentials(value *AccessKey) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        b.region,
			Credentials: credentials.NewStaticCredentials(
				value.AccessKeyID,
				value.SecretAccessKey,
				"",
			),
		},
	})
}

func (b *ClientBuilder) BuildSessionWithOptions() (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile.Profile(),
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        b.region,
		},
	})
}

// Build uses the information stored in the builder to build a new AWS client.
func (b *ClientBuilder) Build() (Client, error) {
	// Check parameters:
	if b.logger == nil {
		return nil, fmt.Errorf("Logger is mandatory")
	}

	// Create the AWS logger:
	logger, err := logging.NewAWSLogger().
		Logger(b.logger).
		Build()
	if err != nil {
		return nil, err
	}

	var sess *session.Session

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
		return nil, fmt.Errorf("Failed to connect to AWS. Use a GovCloud region in your profile")
	}

	// Create the AWS session:
	if b.credentials != nil {
		sess, err = b.BuildSessionWithOptionsCredentials(b.credentials)
	} else {
		sess, err = b.BuildSessionWithOptions()
	}
	if err != nil {
		return nil, err
	}

	// Add ROSACLI as user-agent
	sess.Handlers.Build.PushFrontNamed(addROSAVersionToUserAgent)

	if profile.Profile() != "" {
		b.logger.Debugf("Using AWS profile: %s", profile.Profile())
	}

	// Check that the AWS credentials are available:
	// TODO: No need to do this twice, we're essentially doing the
	// same thing in getClientDetails()
	// We should implement getClientDetails() here or a new validation func
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		b.logger.Debugf("Failed to find credentials: %v", err)
		return nil, fmt.Errorf("Failed to find credentials. Check your AWS configuration and try again")
	}

	// Check that the region is set:
	region := aws.StringValue(sess.Config.Region)
	if region == "" {
		return nil, fmt.Errorf("Region is not set. Use --region to set the region")
	}

	// Update session config
	sess = sess.Copy(&aws.Config{
		Retryer: buildCustomRetryer(),
		Logger:  logger,
		HTTPClient: &http.Client{
			Transport: http.DefaultTransport,
		},
	})

	if b.logger.IsLevelEnabled(logrus.DebugLevel) {
		var dumper http.RoundTripper
		dumper, err = logging.NewRoundTripper().
			Logger(b.logger).
			Next(sess.Config.HTTPClient.Transport).
			Build()
		if err != nil {
			return nil, err
		}
		sess.Config.HTTPClient.Transport = dumper
	}

	// Create and populate the object:
	c := &awsClient{
		logger:              b.logger,
		iamClient:           iam.New(sess),
		ec2Client:           ec2.New(sess),
		orgClient:           organizations.New(sess),
		s3Client:            s3.New(sess),
		smClient:            secretsmanager.New(sess),
		stsClient:           sts.New(sess),
		cfClient:            cloudformation.New(sess),
		servicequotasClient: servicequotas.New(sess),
		awsSession:          sess,
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

func (c *awsClient) GetIAMCredentials() (credentials.Value, error) {
	return c.awsSession.Config.Credentials.Get()
}

func (c *awsClient) GetRegion() string {
	return aws.StringValue(c.awsSession.Config.Region)
}

func (c *awsClient) FetchPublicSubnetMap(subnets []*ec2.Subnet) (map[string]bool, error) {
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
		routeTablesResp, err := c.ec2Client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("association.subnet-id"),
					Values: subnetIds,
				},
			},
		})
		if err != nil {
			return mapSubnetIdToPublic, err
		}
		if routeTablesResp == nil {
			return mapSubnetIdToPublic, fmt.Errorf(
				"No route table found for associated subnets '%s'",
				helper.SliceToSortedString(aws.StringValueSlice(subnetIds)),
			)
		}
		for _, routes := range routeTablesResp.RouteTables {
			for _, association := range routes.Associations {
				subnetAssociation := aws.StringValue(association.SubnetId)
				mapSubnetIdToPublic[subnetAssociation] = false
				for _, route := range routes.Routes {
					if strings.HasPrefix(aws.StringValue(route.GatewayId), "igw") {
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

func (c *awsClient) ListSubnets(subnetIds ...string) ([]*ec2.Subnet, error) {

	if len(subnetIds) == 0 {
		return c.getSubnetIDs(&ec2.DescribeSubnetsInput{})
	}

	var ids []*string

	for i := range subnetIds {
		ids = append(ids, &subnetIds[i])
	}

	return c.getSubnetIDs(&ec2.DescribeSubnetsInput{
		SubnetIds: ids,
	})
}
func (c *awsClient) GetSubnetAvailabilityZone(subnetID string) (string, error) {
	res, err := c.ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{SubnetIds: []*string{aws.String(subnetID)}})
	if err != nil {
		return "", err
	}
	if len(res.Subnets) < 1 {
		return "", fmt.Errorf("Failed to get subnet with ID '%s'", subnetID)
	}

	return *res.Subnets[0].AvailabilityZone, nil
}

func (c *awsClient) GetVPCPrivateSubnets(subnetID string) ([]*ec2.Subnet, error) {
	subnets, err := c.GetVPCSubnets(subnetID)
	if err != nil {
		return nil, err
	}

	return c.FilterVPCsPrivateSubnets(subnets)
}

// getVPCSubnets gets a subnet ID and fetches all the subnets that belong to the same VPC as the provided subnet.
func (c *awsClient) GetVPCSubnets(subnetID string) ([]*ec2.Subnet, error) {
	// Fetch the subnet details
	subnets, err := c.getSubnetIDs(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: []*string{aws.String(subnetID)},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(subnets) < 1 {
		return nil, fmt.Errorf("Failed to get subnet with ID '%s'", subnetID)
	}

	// Fetch VPC's subnets
	vpcID := subnets[0].VpcId
	subnets, err = c.getSubnetIDs(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{vpcID},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(subnets) < 1 {
		return nil, fmt.Errorf("Failed to get the subnets of VPC with ID '%s'", *vpcID)
	}

	return subnets, nil
}

// FilterPrivateSubnets gets a slice of subnets that belongs to the same VPC and filters the private subnets.
// Assumption: subnets - non-empty slice.
func (c *awsClient) FilterVPCsPrivateSubnets(subnets []*ec2.Subnet) ([]*ec2.Subnet, error) {
	// Fetch VPC route tables
	vpcID := subnets[0].VpcId
	describeRouteTablesOutput, err := c.ec2Client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{vpcID},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(describeRouteTablesOutput.RouteTables) < 1 {
		return nil, fmt.Errorf("Failed to find VPC '%s' route table", *vpcID)
	}

	var privateSubnets []*ec2.Subnet
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
		return nil, fmt.Errorf("Failed to find private subnets associated with VPC '%s'", *subnets[0].VpcId)
	}

	return privateSubnets, nil
}

// isPublicSubnet a public subnet is a subnet that's associated with a route table that has a route to an
// internet gateway
func (c *awsClient) isPublicSubnet(subnetID *string, routeTables []*ec2.RouteTable) (bool, error) {
	subnetRouteTable, err := c.getSubnetRouteTable(subnetID, routeTables)
	if err != nil {
		return false, err
	}

	for _, route := range subnetRouteTable.Routes {
		if strings.Contains(aws.StringValue(route.GatewayId), "igw") {
			return true, nil
		}
	}

	return false, nil
}

func (c *awsClient) getSubnetRouteTable(subnetID *string, routeTables []*ec2.RouteTable) (*ec2.RouteTable, error) {
	// Subnet route table — A route table that's associated with a subnet
	for _, routeTable := range routeTables {
		for _, association := range routeTable.Associations {
			if aws.StringValue(association.SubnetId) == aws.StringValue(subnetID) {
				return routeTable, nil
			}
		}
	}

	// A subnet can be explicitly associated with custom route table, or implicitly or explicitly associated with the
	// main route table.
	for _, routeTable := range routeTables {
		for _, association := range routeTable.Associations {
			if aws.BoolValue(association.Main) {
				return routeTable, nil
			}
		}
	}

	// Each subnet in the VPC must be associated with a route table
	return nil, fmt.Errorf("Failed to find subnet '%s' route table", *subnetID)
}

// getSubnetIDs will return the list of subnetsIDs supported for the region picked.
// It is possible to pass non-empty `describeSubnetsInput` to filter results.
func (c *awsClient) getSubnetIDs(describeSubnetsInput *ec2.DescribeSubnetsInput) ([]*ec2.Subnet, error) {
	res, err := c.ec2Client.DescribeSubnets(describeSubnetsInput)
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
}

func (c *awsClient) GetCreator() (*Creator, error) {
	getCallerIdentityOutput, err := c.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return CreatorForCallerIdentity(getCallerIdentityOutput)
}

// CreatorForCallerIdentity adapts an STS CallerIdentity to the ROSA *Creator
func CreatorForCallerIdentity(identity *sts.GetCallerIdentityOutput) (*Creator, error) {
	creatorARN := aws.StringValue(identity.Arn)

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
	}, nil
}

// Checks if given credentials are valid.
func (c *awsClient) ValidateCredentials() (bool, error) {
	// Validate the AWS credentials by calling STS GetCallerIdentity
	// This will fail if the AWS access key and secret key are invalid. This
	// will also work for STS credentials with access key, secret key and session
	// token
	_, err := c.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "InvalidClientTokenId") {
			awsErr := awserr.New("InvalidClientTokenId",
				"Invalid AWS Credentials. For help configuring your credentials, see "+
					"https://docs.openshift.com/rosa/rosa_install_access_delete_clusters/rosa_getting_started_iam/"+
					"rosa-config-aws-account.html#rosa-configuring-aws-account_rosa-config-aws-account",
				err)
			return false, awsErr

		}
		return false, err
	}

	return true, nil
}

func (c *awsClient) CheckAdminUserNotExisting(userName string) (err error) {
	userList, err := c.iamClient.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return err
	}
	for _, user := range userList.Users {
		if *user.UserName == userName {
			return fmt.Errorf("Error creating user: IAM user '%s' already exists.\n"+
				"Ensure user '%s' IAM user does not exist, then retry with\n"+
				"rosa init",
				*user.UserName, *user.UserName)
		}
	}
	return nil
}

func (c *awsClient) CheckAdminUserExists(userName string) (err error) {
	_, err = c.iamClient.GetUser(&iam.GetUserInput{UserName: aws.String(userName)})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) GetClusterRegionTagForUser(username string) (string, error) {
	user, err := c.iamClient.GetUser(&iam.GetUserInput{UserName: aws.String(username)})
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
	_, err := c.iamClient.TagUser(&iam.TagUserInput{
		UserName: aws.String(username),
		Tags: []*iam.Tag{
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
	creds, err := c.awsSession.Config.Credentials.Get()
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
			switch typed := err.(type) {
			case awserr.Error:
				// Waiter reached maximum attempts waiting for the resource to be ready
				if typed.Code() == "InvalidClientTokenId" {
					wait := time.Duration((i * 200)) * time.Millisecond
					waited := time.Since(start)
					logger.Debug(fmt.Sprintf("InvalidClientTokenId, waited %.2f\n", waited.Seconds()))
					time.Sleep(wait)
				}
				if typed.Code() == "AccessDenied" {
					wait := time.Duration((i * 200)) * time.Millisecond
					waited := time.Since(start)
					logger.Debug(fmt.Printf("AccessDenied, waited %.2f\n", waited.Seconds()))
					time.Sleep(wait)
				}
			}

			// If we've still got an error on the last attempt return it
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
	createIAMUserAccessKeyOutput, err := c.iamClient.CreateAccessKey(
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
	listAccessKeysOutput, err := c.iamClient.ListAccessKeys(
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
		_, err = c.iamClient.DeleteAccessKey(
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
	role, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return false, "", nil
			default:
				return false, "", err
			}
		}
	}

	return true, aws.StringValue(role.Role.Arn), nil
}

func (c *awsClient) GetRoleByARN(roleARN string) (*iam.Role, error) {
	// validate arn
	parsedARN, err := arn.Parse(roleARN)
	if err != nil {
		return nil, fmt.Errorf("expected '%s' to be a valid IAM role ARN: %s", roleARN, err)
	}

	// validate arn is for a role resource
	resource := parsedARN.Resource
	isRole := strings.Contains(resource, "role/")
	if !isRole {
		return nil, fmt.Errorf("expected ARN '%s' to be IAM role resource", roleARN)
	}

	// get resource name

	m := strings.LastIndex(resource, "/")
	roleName := resource[m+1:]

	roleOutput, err := c.iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return nil, err
	}
	return roleOutput.Role, nil
}

// DescribeAvailabilityZones fetches the region's availability zones with type `availability-zone`
func (c *awsClient) DescribeAvailabilityZones() ([]string, error) {
	describeAvailabilityZonesOutput, err := c.ec2Client.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("zone-type"),
				Values: []*string{aws.String("availability-zone")},
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
	availabilityZones, err := c.ec2Client.DescribeAvailabilityZones(
		&ec2.DescribeAvailabilityZonesInput{ZoneNames: []*string{aws.String(availabilityZoneName)}})
	if err != nil {
		return false, err
	}
	if len(availabilityZones.AvailabilityZones) < 1 {
		return false, fmt.Errorf("Failed to find availability zone '%s'", availabilityZoneName)
	}

	return aws.StringValue(availabilityZones.AvailabilityZones[0].ZoneType) == LocalZone, nil
}

func (c *awsClient) GetAvailabilityZoneType(availabilityZoneName string) (string, error) {
	availabilityZones, err := c.ec2Client.DescribeAvailabilityZones(
		&ec2.DescribeAvailabilityZonesInput{ZoneNames: []*string{aws.String(availabilityZoneName)}})
	if err != nil {
		return "", err
	}
	if len(availabilityZones.AvailabilityZones) < 1 {
		return "", fmt.Errorf("Failed to find availability zone '%s'", availabilityZoneName)
	}
	return aws.StringValue(availabilityZones.AvailabilityZones[0].ZoneType), nil
}

func (c *awsClient) DetachRolePolicies(roleName string) error {
	attachedPolicies := make([]*iam.AttachedPolicy, 0)
	isTruncated := true
	var marker *string
	for isTruncated {
		resp, err := c.iamClient.ListAttachedRolePolicies(
			&iam.ListAttachedRolePoliciesInput{
				Marker:   marker,
				RoleName: &roleName,
			},
		)
		if err != nil {
			return err
		}
		isTruncated = *resp.IsTruncated
		marker = resp.Marker
		attachedPolicies = append(attachedPolicies, resp.AttachedPolicies...)
	}
	for _, attachedPolicy := range attachedPolicies {
		err := c.detachRolePolicy(*attachedPolicy.PolicyArn, roleName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *awsClient) detachRolePolicy(policyArn string, roleName string) error {
	_, err := c.iamClient.DetachRolePolicy(&iam.DetachRolePolicyInput{PolicyArn: &policyArn, RoleName: &roleName})
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
	_, err := c.s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return weberr.Errorf("Bucket '%s' already exists.", bucketName)
	}
	bucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if region != DefaultRegion {
		bucketInput.SetCreateBucketConfiguration(&s3.CreateBucketConfiguration{
			LocationConstraint: &region,
		})
	}
	_, err = c.s3Client.CreateBucket(bucketInput)
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutPublicAccessBlock(&s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(true),
			IgnorePublicAcls:      aws.Bool(true),
			BlockPublicPolicy:     aws.Bool(false),
			RestrictPublicBuckets: aws.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(fmt.Sprintf(ReadOnlyAnonUserPolicyTemplate, bucketName)),
	})
	if err != nil {
		return err
	}

	_, err = c.s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
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
	_, err := c.s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NotFound" {
			return nil
		}
		return err
	}
	err = c.emptyS3Bucket(bucketName)
	if err != nil {
		return err
	}
	_, err = c.s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) emptyS3Bucket(bucketName string) error {
	objects, err := c.s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}
	for _, object := range (*objects).Contents {
		_, err = c.s3Client.DeleteObject(&s3.DeleteObjectInput{
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
	_, err := c.s3Client.PutObject(&s3.PutObjectInput{
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
	createSecretResponse, err := c.smClient.CreateSecret(
		&secretsmanager.CreateSecretInput{
			Description:  aws.String(fmt.Sprintf("Secret for %s", name)),
			Name:         aws.String(name),
			SecretString: aws.String(secret),
			Tags: []*secretsmanager.Tag{{
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
	_, err := c.smClient.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretArn),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
			return nil
		}
	}
	_, err = c.smClient.DeleteSecret(
		&secretsmanager.DeleteSecretInput{
			ForceDeleteWithoutRecovery: aws.Bool(true),
			SecretId:                   aws.String(secretArn),
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *awsClient) GetSecurityGroupIds(vpcId string) ([]*ec2.SecurityGroup, error) {
	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{vpcId}),
			},
		},
	}
	securityGroups := []*ec2.SecurityGroup{}
	err := c.ec2Client.DescribeSecurityGroupsPages(describeSecurityGroupsInput,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			for _, sg := range page.SecurityGroups {
				if tags.Ec2ResourceHasTag(sg.Tags, tags.RedHatManaged, strconv.FormatBool(true)) {
					continue
				}
				if aws.StringValue(sg.GroupName) == "default" {
					continue
				}
				securityGroups = append(securityGroups, sg)
			}
			return page.NextToken != nil
		})
	if err != nil {
		return []*ec2.SecurityGroup{}, err
	}
	return securityGroups, nil
}

// CustomRetryer wraps the aws SDK's built in DefaultRetryer allowing for
// additional custom features
type CustomRetryer struct {
	client.DefaultRetryer
}

// ShouldRetry overrides the SDK's built in DefaultRetryer adding customization
// to not retry 5xx status codes.
func (r CustomRetryer) ShouldRetry(req *request.Request) bool {
	if req.HTTPResponse.StatusCode >= 500 {
		return false
	}
	logger := logging.NewLogger()

	if req.HTTPRequest.Header.Get("ROSA-Request-Id") == "" {
		req.HTTPRequest.Header.Add("ROSA-Request-Id", uuid.New().String())
	}

	if strings.Contains(req.Error.Error(), "Throttling") {
		logger.Warn(fmt.Sprintf(
			"Throttling Rate limit exceeded. Retrying [ROSA-Request-Id: %s / %s %s]: %v/%v",
			req.HTTPRequest.Header.Get("ROSA-Request-Id"),
			req.HTTPRequest.Method,
			req.HTTPRequest.URL.Host,
			req.RetryCount+1,
			r.MaxRetries()))
	}

	return r.DefaultRetryer.ShouldRetry(req)
}

func buildCustomRetryer() CustomRetryer {
	return CustomRetryer{
		DefaultRetryer: client.DefaultRetryer{
			NumMaxRetries:    numMaxRetries,
			MinRetryDelay:    minRetryDelay,
			MinThrottleDelay: minThrottleDelay,
			MaxThrottleDelay: maxThrottleDelay,
		},
	}
}
