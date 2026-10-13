package aws

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsCommonUtils "github.com/openshift-online/ocm-common/pkg/aws/utils"
	awsCommonValidations "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/sirupsen/logrus"
	"github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/arguments"
	client "github.com/openshift/rosa/pkg/aws/api_interface"
	awscb "github.com/openshift/rosa/pkg/aws/commandbuilder"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/constants"
	"github.com/openshift/rosa/pkg/fedramp"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/reporter"
	rprtr "github.com/openshift/rosa/pkg/reporter"
)

// AWS accepted role name: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html
var RoleNameRE = regexp.MustCompile(`^[\w+=,.@-]+$`)

// AWS accepted arn format: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
var RoleArnRE = regexp.MustCompile(
	`^arn:aws[\w-]*:iam::\d{12}:role(?:\/+[\w+=,.@-]+)+$`,
)
var PolicyArnRE = regexp.MustCompile(
	`^arn:aws[\w-]*:iam::(\d{12}|aws):policy(?:\/+[\w+=,.@-]+)+$`,
)

// UserTagKeyRE , UserTagValueRE - https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html#tag-conventions
var UserTagKeyRE = regexp.MustCompile(`^[\pL\pZ\pN_.:/=+\-@]{1,128}$`)
var UserTagValueRE = regexp.MustCompile(`^[\pL\pZ\pN_.:/=+\-@]{0,256}$`)

// the following regex defines five different patterns:
// first pattern is to validate IPv4 address
// second,is for IPv4 CIDR range validation
// third pattern is to validate domains
// and the fifth petterrn is to be able to remove the existing no-proxy value by typing empty string ("").
// nolint
var UserNoProxyRE = regexp.MustCompile(
	`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$|^(.?[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$|^""$`,
)

const (
	SecretsManager = "secretsmanager"
	awsDefaultId   = "aws"
	maxWaitDur     = 5 * time.Minute
)

func GetJumpAccount(env string) string {
	jumpAccounts := JumpAccounts
	if fedramp.Enabled() {
		jumpAccounts = fedramp.JumpAccounts
	}
	return jumpAccounts[env]
}

// JumpAccounts are the various of AWS accounts used for the installer jump role in the various OCM environments
var JumpAccounts = map[string]string{
	"production":  "710019948333",
	"staging":     "644306948063",
	"integration": "896164604406",
	"local":       "765374464689",
	"local-proxy": "765374464689",
	"crc":         "765374464689",
}

var ARNPath = regexp.MustCompile(`^\/[a-zA-Z0-9\/]*\/$`)

func ARNValidator(input interface{}) error {
	if str, ok := input.(string); ok {
		if str == "" {
			return nil
		}
		_, err := arn.Parse(str)
		if err != nil {
			return fmt.Errorf("Invalid ARN: %s", err)
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got %v", input)
}

func SecretManagerArnValidator(input interface{}) error {
	str, ok := input.(string)
	if !ok {
		return fmt.Errorf("can only validate strings, got %v", input)
	}
	if str == "" {
		return nil
	}
	if !arn.IsARN(str) {
		return fmt.Errorf("Secret ARN '%s' is not a valid ARN", str)
	}
	parsedSecretArn, err := arn.Parse(str)
	if err != nil {
		return fmt.Errorf("Invalid ARN: %s", err)
	}
	if parsedSecretArn.Service != constants.SecretsManagerService {
		return fmt.Errorf("Secret ARN '%s' is not a valid secrets manager ARN", str)
	}
	return nil
}

func ARNPathValidator(input interface{}) error {
	if str, ok := input.(string); ok {
		if str == "" {
			return nil
		}
		if !ARNPath.MatchString(str) {
			return fmt.Errorf("invalid ARN Path. It must begin and end with / and " +
				"contain only alphanumeric characters")
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got %v", input)
}

// GetRegion will return a region selected by the user or given as a default to the AWS client.
// If the region given is empty, it will first attempt to use the default, and, failing that, will
// prompt for user input.
func GetRegion(region string) (string, error) {
	if region == "" {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return "", fmt.Errorf("Error loading default AWS configuration: %v", err)
		}

		region = cfg.Region
	}
	return region, nil
}

// getClientDetails will return the *iam.User associated with the provided client's credentials,
// a boolean indicating whether the user is the 'root' account, and any error encountered
// while trying to gather the info.
func getClientDetails(awsClient *awsClient) (*sts.GetCallerIdentityOutput, bool, error) {
	rootUser := false

	_, err := awsClient.ValidateCredentials()
	if err != nil {
		return nil, rootUser, err
	}

	user, err := awsClient.stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, rootUser, err
	}

	// Detect whether the AWS account's root user is being used
	parsed, err := arn.Parse(*user.Arn)
	if err != nil {
		return nil, rootUser, err
	}
	if parsed.AccountID == *user.UserId {
		rootUser = true
	}

	return user, rootUser, nil
}

// Currently user can rosa init using the region from their config or using --region
// When checking for cloud formation we need to check in the region used by the user
func GetAWSClientForUserRegion(reporter reporter.Logger, logger *logrus.Logger,
	supportedRegions []string, useLocalCreds bool) Client {
	// Get AWS region from env
	awsRegionInUserConfig, err := GetRegion(arguments.GetRegion())
	if err != nil {
		reporter.Errorf("Error getting region: %v", err)
		os.Exit(1)
	}
	if awsRegionInUserConfig == "" {
		reporter.Errorf("AWS Region not set")
		os.Exit(1)
	}
	if !helper.Contains(supportedRegions, awsRegionInUserConfig) {
		reporter.Errorf("Unsupported region '%s', available regions: %s",
			awsRegionInUserConfig, helper.SliceToSortedString(supportedRegions))
		os.Exit(1)
	}

	// Create the AWS client:
	client, err := NewClient().
		Logger(logger).
		Region(awsRegionInUserConfig).
		UseLocalCredentials(useLocalCreds).
		Build()
	if err != nil {
		reporter.Errorf("Error creating aws client for stack validation: %v", err)
		os.Exit(1)
	}
	regionUsedForInit, err := client.GetClusterRegionTagForUser(AdminUserName)
	if err != nil || regionUsedForInit == "" {
		return client
	}

	if regionUsedForInit != awsRegionInUserConfig {
		if !helper.Contains(supportedRegions, regionUsedForInit) {
			reporter.Errorf("Unsupported region '%s', available regions: %s",
				regionUsedForInit, helper.SliceToSortedString(supportedRegions))
			os.Exit(1)
		}
		// Create the AWS client with the region used in the init
		//So we can check for the stack in that region
		awsClient, err := NewClient().
			Logger(logger).
			Region(regionUsedForInit).
			UseLocalCredentials(useLocalCreds).
			Build()
		if err != nil {
			reporter.Errorf("Error creating aws client for stack validation: %v", err)
			os.Exit(1)
		}
		return awsClient
	}
	return client
}

func isSTS(ARN arn.ARN) bool {
	// If the client is using STS credentials we'll attempt to find the role
	// assumed by the user and validate that using PolicySimulator
	resource := strings.Split(ARN.Resource, "/")
	resourceType := 0
	// Example STS role ARN "arn:aws:sts::123456789123:assumed-role/OrganizationAccountAccessRole/UserAccess"
	// if the "service" is STS and the "resource-id" sectino of the ARN contains 3 sections delimited by
	// "/" we can validate its an assumed-role and assume the role name is the "parent-resource" and construct
	// a role ARN
	// https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
	if ARN.Service == "sts" &&
		resource[resourceType] == "assumed-role" {
		return true
	}
	return false
}

func resolveSTSRole(ARN arn.ARN) (*string, error) {
	// If the client is using STS credentials we'll attempt to find the role
	// assumed by the user and validate that using PolicySimulator
	resource := strings.Split(ARN.Resource, "/")
	parentResource := 1
	// Example STS role ARN "arn:aws:sts::123456789123:assumed-role/OrganizationAccountAccessRole/UserAccess"
	// if the "service" is STS and the "resource-id" sectino of the ARN contains 3 sections delimited by
	// "/" we can validate its an assumed-role and assume the role name is the "parent-resource" and construct
	// a role ARN
	// https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
	if isSTS(ARN) && len(resource) == 3 {
		// Construct IAM role ARN
		roleARNString := fmt.Sprintf(
			"arn:%s:iam::%s:role/%s", ARN.Partition, ARN.AccountID, resource[parentResource])
		// Parse it to validate its ok
		err := ARNValidator(roleARNString)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse role ARN %s created from sts role: %v", roleARNString, err)
		}
		return &roleARNString, nil
	}

	return nil, fmt.Errorf("ARN %s doesn't appear to have a a resource-id that confirms to an STS user", ARN.String())
}

func UserTagValidator(input interface{}) error {
	var inputTags []string
	inputType := reflect.TypeOf(input).Kind()
	switch inputType {
	case reflect.String:
		if input.(string) == "" {
			return nil
		}
		inputTags = strings.Split(input.(string), ",")
	case reflect.Slice:
		if reflect.TypeOf(input).Elem().Kind() != reflect.String {
			return fmt.Errorf("unable to verify tags, incompatible type, expected slice of string got: '%s'",
				inputType.String())
		}
		inputTags = input.([]string)
	default:
		return fmt.Errorf("can only validate string types, got %v", inputType.String())
	}

	if duplicate, hasDupe := hasDuplicateTagKey(inputTags); hasDupe {
		return fmt.Errorf("invalid tags, user tag keys must be unique, duplicate key '%s' found", duplicate)
	}

	delimiter := GetTagsDelimiter(inputTags)
	for _, t := range inputTags {
		tag := strings.Split(t, delimiter)
		if len(tag) != 2 {
			return fmt.Errorf("invalid tag format for tag '%s'. Expected tag format: 'key value'", tag)
		}

		if tag[0] == "" || tag[1] == "" {
			return fmt.Errorf("invalid tag format, tag key or tag value can not be empty")
		}

		if !UserTagKeyRE.MatchString(tag[0]) {
			return fmt.Errorf("expected a valid user tag key '%s' matching %s", tag[0], UserTagKeyRE.String())
		}

		if !UserTagValueRE.MatchString(tag[1]) {
			return fmt.Errorf("expected a valid user tag value '%s' matching %s", tag[1], UserTagValueRE.String())
		}
	}
	return nil
}

func GetTagsDelimiter(tags []string) string {
	sanitizeTags(tags)
	tagsString := strings.Join(tags, "")
	if strings.Contains(tagsString, " ") {
		return " "
	}
	return ":"
}

func UserTagDuplicateValidator(input interface{}) error {
	if str, ok := input.(string); ok {
		if str == "" {
			return nil
		}
		splitTags := strings.Split(str, ",")
		sanitizeTags(splitTags)
		duplicate, found := hasDuplicateTagKey(splitTags)
		if found {
			return fmt.Errorf("user tag keys must be unique, duplicate key '%s' found", duplicate)
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got %v", input)
}

func hasDuplicateTagKey(tags []string) (string, bool) {
	delimiter := GetTagsDelimiter(tags)
	visited := make(map[string]bool)
	for _, t := range tags {
		tag := strings.Split(t, delimiter)
		if visited[tag[0]] {
			return tag[0], true
		}
		visited[tag[0]] = true
	}
	return "", false
}

func sanitizeTags(tags []string) {
	for i, str := range tags {
		tags[i] = strings.TrimSpace(str)
	}
}

func UserNoProxyValidator(input interface{}) error {
	if str, ok := input.(string); ok {
		if str == "" {
			return nil
		}
		noProxyValues := strings.Split(str, ",")

		for _, v := range noProxyValues {
			if !UserNoProxyRE.MatchString(v) {
				return fmt.Errorf("expected a valid user no-proxy value: '%s' matching %s", v, UserNoProxyRE.String())
			}
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got %v", input)
}

func UserNoProxyDuplicateValidator(input interface{}) error {
	if str, ok := input.(string); ok {
		if str == "" {
			return nil
		}
		values := strings.Split(str, ",")
		duplicate, found := HasDuplicates(values)
		if found {
			return fmt.Errorf("no-proxy values must be unique, duplicate key '%s' found", duplicate)
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got %v", input)
}

func HasDuplicates(valSlice []string) (string, bool) {
	visited := make(map[string]bool)
	for _, v := range valSlice {
		if visited[v] {
			return v, true
		}
		visited[v] = true
	}
	return "", false
}

func GetTagValues(tagsValue []iamtypes.Tag) (roleType string, version string) {
	for _, tag := range tagsValue {
		switch aws.ToString(tag.Key) {
		case tags.RoleType:
			roleType = aws.ToString(tag.Value)
		case awsCommonValidations.OpenShiftVersion:
			version = aws.ToString(tag.Value)
		}
	}
	return
}

func GetOCMRoleName(prefix string, role string, postfix string) string {
	name := fmt.Sprintf("%s-%s-Role-%s", prefix, role, postfix)
	return awsCommonUtils.TruncateRoleName(name)
}

func GetUserRoleName(prefix string, role string, userName string) string {
	name := fmt.Sprintf("%s-%s-%s-Role", prefix, role, userName)
	return awsCommonUtils.TruncateRoleName(name)
}

func GetOperatorPolicyName(prefix string, namespace string, name string) string {
	policy := fmt.Sprintf("%s-%s-%s", prefix, namespace, name)
	return awsCommonUtils.TruncateRoleName(policy)
}

func GetAdminPolicyName(name string) string {
	return fmt.Sprintf("%s-Admin-Policy", name)
}

func GetPolicyName(name string) string {
	return fmt.Sprintf("%s-Policy", name)
}

func GetOperatorPolicyARN(partition string, accountID string,
	prefix string, namespace string, name string, path string) string {
	return getPolicyARN(partition, accountID, GetOperatorPolicyName(prefix, namespace, name), path)
}

func GetAdminPolicyARN(partition string, accountID string, name string, path string) string {
	return getPolicyARN(partition, accountID, GetAdminPolicyName(name), path)
}

func GetPolicyArnWithSuffix(partition string, accountID string, name string, path string) string {
	return getPolicyARN(partition, accountID, GetPolicyName(name), path)
}

func GetPolicyArn(partition string, accountID string, name string, path string) string {
	return getPolicyARN(partition, accountID, name, path)
}

func getPolicyARN(partition string, accountID string, name string, path string) string {
	str := fmt.Sprintf("arn:%s:iam::%s:policy", partition, accountID)
	if path != "" {
		str = fmt.Sprintf("%s%s", str, path)
		return fmt.Sprintf("%s%s", str, name)
	}
	return fmt.Sprintf("%s/%s", str, name)
}

func GetPathFromAccountRole(cluster *cmv1.Cluster, roleNameSuffix string) (string, error) {
	accRoles := GetAccountRolesArnsMap(cluster)
	if accRoles[roleNameSuffix] == "" {
		return "", nil
	}
	return GetPathFromARN(accRoles[roleNameSuffix])
}

func GetPathFromARN(arnStr string) (string, error) {
	parse, err := arn.Parse(arnStr)
	if err != nil {
		return "", err
	}
	resource := parse.Resource
	firstIndex := strings.Index(resource, "/")
	lastIndex := strings.LastIndex(resource, "/")
	if firstIndex == lastIndex {
		return "", nil
	}
	path := resource[firstIndex : lastIndex+1]
	return path, nil
}

func IsArnAssumedRole(arnStr string) (bool, error) {
	parse, err := arn.Parse(arnStr)
	if err != nil {
		return false, err
	}
	resource := strings.SplitN(parse.Resource, "/", 2)[0]
	if resource == AssumedRoleRolePrefix {
		return true, nil
	}
	return false, nil
}

func GetRoleARN(accountID string, name string, path string, partition string) string {
	if path == "" {
		path = "/"
	}
	return fmt.Sprintf("arn:%s:iam::%s:role%s%s", partition, accountID, path, name)
}

func GetOIDCProviderARN(partition string, accountID string, providerURL string) string {
	return fmt.Sprintf("arn:%s:iam::%s:oidc-provider/%s", partition, accountID, providerURL)
}

func GetPrefixFromAccountRole(cluster *cmv1.Cluster, roleNameSuffix string) (string, error) {
	// role name here is resource id (ex: dle-test-Worker-Role)
	roleName, err := GetAccountRoleName(cluster, roleNameSuffix)
	if err != nil {
		return "", err
	}

	var suffix string
	if IsHostedCPManagedPolicies(cluster) {
		suffix = fmt.Sprintf("-%s-%s-Role", HCPSuffixPattern, roleNameSuffix)
	} else {
		suffix = fmt.Sprintf("-%s-Role", roleNameSuffix)
	}

	rolePrefix := TrimRoleSuffix(roleName, suffix)
	return rolePrefix, nil
}

func GetPrefixFromInstallerAccountRole(cluster *cmv1.Cluster) (string, error) {
	return GetPrefixFromAccountRole(cluster, AccountRoles[InstallerAccountRole].Name)
}

// Role names can be truncated if they are over 64 chars, so we need to make sure we aren't missing a truncated suffix
func TrimRoleSuffix(orig, sufix string) string {
	for i := len(sufix); i >= 0; i-- {
		if strings.HasSuffix(orig, sufix[:i]) {
			return orig[:len(orig)-i]
		}
	}
	return orig
}

func GetPrefixFromOperatorRole(cluster *cmv1.Cluster) string {
	operator := cluster.AWS().STS().OperatorIAMRoles()[0]
	roleName, _ := GetResourceIdFromARN(operator.RoleARN())
	rolePrefix := TrimRoleSuffix(roleName, fmt.Sprintf("-%s-%s", operator.Namespace(), operator.Name()))
	return rolePrefix
}

func GetOperatorRolePolicyPrefixFromCluster(cluster *cmv1.Cluster, awsClient Client) (string, error) {
	installerRolePrefix, err := GetPrefixFromInstallerAccountRole(cluster)
	if err != nil {
		return "", err
	}
	// Check if installer role prefix follows standard
	installerRoleName, err := GetResourceIdFromARN(cluster.AWS().STS().RoleARN())
	if err != nil {
		return "", err
	}
	hasStandardAccRole := installerRolePrefix != installerRoleName
	if hasStandardAccRole {
		return installerRolePrefix, nil
	}

	// If is non standard try to find most common operator policy prefix from current attached policies
	operatorRoles := cluster.AWS().STS().OperatorIAMRoles()
	policyPrefixCountMap := make(map[string]int)
	policyNames := []string{}
	for _, operatorRole := range operatorRoles {
		roleName, _ := GetResourceIdFromARN(operatorRole.RoleARN())
		policiesDetails, err := awsClient.GetAttachedPolicy(&roleName)
		if err != nil {
			return "", err
		}
		attachedPoliciesDetail := FindAllAttachedPolicyDetails(policiesDetails)
		for _, attachedPolicyDetail := range attachedPoliciesDetail {
			policyNames = append(policyNames, attachedPolicyDetail.PolicyName)
			index := strings.LastIndex(attachedPolicyDetail.PolicyName, "-openshift")
			if index != -1 {
				policyPrefix := attachedPolicyDetail.PolicyName[0:index]
				if _, ok := policyPrefixCountMap[policyPrefix]; !ok {
					policyPrefixCountMap[policyPrefix] = 0
				}
				policyPrefixCountMap[policyPrefix]++
			}
		}
	}
	rankedPolicyPrefix := helper.RankMapStringInt(policyPrefixCountMap)
	if len(rankedPolicyPrefix) != 0 {
		return rankedPolicyPrefix[0], nil
	}

	//If no standard prefix is found, tries to look for the longest common prefix
	// TODO: check if it makes sense to only use this and remove "-openshift" later
	policyPrefix := helper.LongestCommonPrefixBySorting(policyNames)
	if policyPrefix != "" {
		return strings.TrimRight(policyPrefix, "-"), nil
	}

	// If nothing works uses operator role prefix as to not interrupt flow
	return cluster.AWS().STS().OperatorRolePrefix(), nil
}

func GetAccountRolesArnsMap(cluster *cmv1.Cluster) map[string]string {
	return map[string]string{
		AccountRoles[InstallerAccountRole].Name:    cluster.AWS().STS().RoleARN(),
		AccountRoles[SupportAccountRole].Name:      cluster.AWS().STS().SupportRoleARN(),
		AccountRoles[ControlPlaneAccountRole].Name: cluster.AWS().STS().InstanceIAMRoles().MasterRoleARN(),
		AccountRoles[WorkerAccountRole].Name:       cluster.AWS().STS().InstanceIAMRoles().WorkerRoleARN(),
	}
}

func GetAccountRoleName(cluster *cmv1.Cluster, accountRole string) (string, error) {
	accRoles := GetAccountRolesArnsMap(cluster)

	// checks to see if role arn exists from map
	if accRoles[accountRole] == "" {
		return "", nil
	}

	// outputs resource id (ex: "dle-test-Worker-Role") from role arn
	return GetResourceIdFromARN(accRoles[accountRole])
}

func GetInstallerAccountRoleName(cluster *cmv1.Cluster) (string, error) {
	return GetAccountRoleName(cluster, AccountRoles[InstallerAccountRole].Name)
}

func GenerateOperatorRolePolicyFiles(reporter reporter.Logger, policies map[string]*cmv1.AWSSTSPolicy,
	credRequests map[string]*cmv1.STSOperator, sharedVpcRoleArn string, partition string) error {
	isSharedVpc := sharedVpcRoleArn != ""
	for credrequest := range credRequests {
		filename := GetOperatorPolicyKey(credrequest, false, isSharedVpc)
		policyDetail := GetPolicyDetails(policies, filename)
		if isSharedVpc {
			policyDetail = InterpolatePolicyDocument(partition, policyDetail, map[string]string{
				"shared_vpc_role_arn": sharedVpcRoleArn,
			})
		}
		//In case any missing policy we don't want to block the user.This might not happen
		if policyDetail == "" {
			continue
		}
		reporter.Debugf("Saving '%s' to the current directory", filename)
		filename = GetFormattedFileName(filename)
		err := helper.SaveDocument(policyDetail, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateAccountRolePolicyFiles(reporter rprtr.Logger, env string, policies map[string]*cmv1.AWSSTSPolicy,
	skipPermissionFiles bool, accountRoles map[string]AccountRole, partition string) error {
	for file := range accountRoles {
		//Get trust policy
		filename := fmt.Sprintf("sts_%s_trust_policy", file)
		policyDetail := GetPolicyDetails(policies, filename)
		policy := InterpolatePolicyDocument(partition, policyDetail, map[string]string{
			"partition":      partition,
			"aws_account_id": GetJumpAccount(env),
		})
		filename = GetFormattedFileName(filename)
		reporter.Debugf("Saving '%s' to the current directory", filename)
		err := helper.SaveDocument(policy, filename)
		if err != nil {
			return err
		}

		//Get the permission policy
		if !skipPermissionFiles {
			err = generatePermissionPolicyFile(reporter, file, policies)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generatePermissionPolicyFile(reporter rprtr.Logger, file string, policies map[string]*cmv1.AWSSTSPolicy) error {
	filename := fmt.Sprintf("sts_%s_permission_policy", file)
	policyDetail := GetPolicyDetails(policies, filename)
	if policyDetail == "" {
		return nil
	}
	//Check and save it as json file
	filename = GetFormattedFileName(filename)
	reporter.Debugf("Saving '%s' to the current directory", filename)

	return helper.SaveDocument(policyDetail, filename)
}

func GetFormattedFileName(filename string) string {
	//Check and save it as json file
	ext := filepath.Ext(filename)
	if ext != ".json" {
		filename = fmt.Sprintf("%s.json", filename)
	}
	return filename
}

func BuildOperatorRolePolicies(prefix string, accountID string, partition string, awsClient Client, commands []string,
	defaultPolicyVersion string, credRequests map[string]*cmv1.STSOperator, path string) []string {
	for credrequest, operator := range credRequests {
		policyARN := GetOperatorPolicyARN(partition, accountID, prefix, operator.Namespace(), operator.Name(), path)
		_, err := awsClient.IsPolicyExists(policyARN)
		if err != nil {
			name := GetOperatorPolicyName(prefix, operator.Namespace(), operator.Name())
			iamTags := map[string]string{
				awsCommonValidations.OpenShiftVersion: defaultPolicyVersion,
				tags.RolePrefix:                       prefix,
				tags.OperatorNamespace:                operator.Namespace(),
				tags.OperatorName:                     operator.Name(),
				tags.RedHatManaged:                    "true",
			}
			createPolicy := awscb.NewIAMCommandBuilder().
				SetCommand(awscb.CreatePolicy).
				AddParam(awscb.PolicyName, name).
				AddParam(awscb.PolicyDocument, fmt.Sprintf("file://openshift_%s_policy.json", credrequest)).
				AddTags(iamTags).
				Build()
			commands = append(commands, createPolicy)
		} else {
			policyTags := map[string]string{
				awsCommonValidations.OpenShiftVersion: defaultPolicyVersion,
			}

			createPolicy := awscb.NewIAMCommandBuilder().
				SetCommand(awscb.CreatePolicyVersion).
				AddParam(awscb.PolicyArn, policyARN).
				AddParam(awscb.PolicyDocument, fmt.Sprintf("file://openshift_%s_policy.json", credrequest)).
				AddParamNoValue(awscb.SetAsDefault).
				Build()

			tagPolicy := awscb.NewIAMCommandBuilder().
				SetCommand(awscb.TagPolicy).
				AddTags(policyTags).
				AddParam(awscb.PolicyArn, policyARN).
				Build()
			commands = append(commands, createPolicy, tagPolicy)
		}
	}
	return commands
}

func FindAllAttachedPolicyDetails(policiesDetails []PolicyDetail) []PolicyDetail {
	attachedPolicies := make([]PolicyDetail, 0)
	for _, policy := range policiesDetails {
		if policy.PolicyType == Attached {
			attachedPolicies = append(attachedPolicies, policy)
		}
	}
	return attachedPolicies
}

func FindFirstAttachedPolicy(policiesDetails []PolicyDetail) PolicyDetail {
	for _, policy := range policiesDetails {
		if policy.PolicyType == Attached {
			return policy
		}
	}
	return PolicyDetail{}
}

func UpgradeOperatorRolePolicies(
	reporter reporter.Logger,
	awsClient Client,
	partition string,
	accountID string,
	prefix string,
	policies map[string]*cmv1.AWSSTSPolicy,
	defaultPolicyVersion string,
	credRequests map[string]*cmv1.STSOperator,
	path string,
	cluster *cmv1.Cluster,
) error {
	isSharedVpc := cluster.AWS().PrivateHostedZoneRoleARN() != ""
	for credrequest, operator := range credRequests {
		policyARN := GetOperatorPolicyARN(partition, accountID, prefix, operator.Namespace(), operator.Name(), path)
		filename := GetOperatorPolicyKey(credrequest, cluster.Hypershift().Enabled(), isSharedVpc)
		policyDetails := GetPolicyDetails(policies, filename)
		if isSharedVpc {
			policyDetails = InterpolatePolicyDocument(partition, policyDetails, map[string]string{
				"shared_vpc_role_arn": cluster.AWS().PrivateHostedZoneRoleARN(),
			})
		}
		policyARN, err := awsClient.EnsurePolicy(policyARN, policyDetails,
			defaultPolicyVersion, map[string]string{
				awsCommonValidations.OpenShiftVersion: defaultPolicyVersion,
				tags.RolePrefix:                       prefix,
				tags.OperatorNamespace:                operator.Namespace(),
				tags.OperatorName:                     operator.Name(),
			}, path)
		if err != nil {
			return err
		}
		reporter.Infof("Upgraded policy with ARN '%s' to version '%s'", policyARN, defaultPolicyVersion)
	}
	return nil
}

type Subnet struct {
	AvailabilityZone string
	OwnerID          string
}

const (
	subnetTemplate        = "%s ('%s','%s','%s', Owner ID: '%s')"
	securityGroupTemplate = "%s ('%s')"
)

// SetSubnetOption Creates a subnet option using a predefined template.
func SetSubnetOption(subnet ec2types.Subnet) string {
	subnetName := ""
	for _, tag := range subnet.Tags {
		switch *tag.Key {
		case "Name":
			subnetName = *tag.Value
		}
		if subnetName != "" {
			break
		}
	}
	return fmt.Sprintf(subnetTemplate, aws.ToString(subnet.SubnetId),
		subnetName, aws.ToString(subnet.VpcId), aws.ToString(subnet.AvailabilityZone),
		aws.ToString(subnet.OwnerId))
}

// SetSecurityGroupOption Creates a security group option using a predefined template.
func SetSecurityGroupOption(securityGroup ec2types.SecurityGroup) string {
	return fmt.Sprintf(securityGroupTemplate,
		aws.ToString(securityGroup.GroupId), aws.ToString(securityGroup.GroupName))
}

// Parse option expects the actual option as the first token followed by a space
func ParseOption(option string) string {
	return strings.Split(option, " ")[0]
}

// GetResourceIdFromARN
// function takes a full AWS ARN, parses it and extracts the last part of the resource field
// e.g. arn:partition:service:region:account-id:resource-type/<some-path>/resource-id
// an assumption is made that there is always a resource-type
// if resource-id is empty then error is returned
func GetResourceIdFromARN(stringARN string) (string, error) {
	parsedARN, err := arn.Parse(stringARN)

	if err != nil {
		return "", fmt.Errorf("couldn't parse arn '%s': %v", stringARN, err)
	}

	index := strings.LastIndex(parsedARN.Resource, "/")
	if index == -1 {
		return "", fmt.Errorf("can't find resource-id in ARN '%s'", stringARN)
	}

	// If the customer has created the provider using a / at the end of the URL for some reason
	if index == len(parsedARN.Resource)-1 {
		return GetResourceIdFromARN(stringARN[:index])
	}

	return parsedARN.Resource[index+1:], nil
}

func GetResourceIdFromOidcProviderARN(stringARN string) (string, error) {
	parsedARN, err := arn.Parse(stringARN)

	if err != nil {
		return "", fmt.Errorf("couldn't parse arn '%s': %v", stringARN, err)
	}

	index := strings.Index(parsedARN.Resource, "/")
	if index == -1 {
		return "", fmt.Errorf("can't find resource-id in ARN '%s'", stringARN)
	}

	// If the customer has created the provider using a / at the end of the URL for some reason
	if index == len(parsedARN.Resource)-1 {
		return GetResourceIdFromOidcProviderARN(stringARN[:index])
	}

	return parsedARN.Resource[index+1:], nil
}

func GetResourceIdFromSecretArn(secretArn string) (string, error) {
	parsedARN, err := arn.Parse(secretArn)

	if err != nil {
		return "", err
	}

	index := strings.LastIndex(parsedARN.Resource, ":")
	if index == -1 || index == len(parsedARN.Resource)-1 {
		return "", weberr.Errorf("can't find resource-id in ARN '%s'", secretArn)
	}

	return parsedARN.Resource[index+1:], nil
}

func FindOperatorRoleNameBySTSOperator(cluster *cmv1.Cluster, operator *cmv1.STSOperator) (string, bool) {
	for _, role := range cluster.AWS().STS().OperatorIAMRoles() {
		if role.Namespace() == operator.Namespace() && role.Name() == operator.Name() {
			name, _ := GetResourceIdFromARN(role.RoleARN())
			return name, true
		}
	}

	return "", false
}

func FindOperatorRoleBySTSOperator(operatorRoles []*cmv1.OperatorIAMRole, operator *cmv1.STSOperator) string {
	for _, operatorRole := range operatorRoles {
		if operatorRole.Name() == operator.Name() && operatorRole.Namespace() == operator.Namespace() {
			return operatorRole.RoleARN()
		}
	}
	return ""
}

// GetPolicyDetails retrieves from the map the policy details for unmanaged and managed policies.
func GetPolicyDetails(policies map[string]*cmv1.AWSSTSPolicy, key string) string {
	policy, ok := policies[key]
	if ok {
		return policy.Details()
	}

	return ""
}

func GetManagedPolicyARN(policies map[string]*cmv1.AWSSTSPolicy, key string) (string, error) {
	policy, ok := policies[key]
	if !ok {
		return "", fmt.Errorf("failed to find policy ARN for '%s'", key)
	}
	if policy.ARN() == "" {
		return "", fmt.Errorf("failed to find policy ARN for '%s'", key)
	}

	return policy.ARN(), nil
}

func GetOperatorPolicyKey(roleType string, hostedCP bool, sharedVpc bool) string {
	if hostedCP {
		return fmt.Sprintf("openshift_hcp_%s_policy", roleType)
	}

	if sharedVpc && roleType == IngressOperatorCloudCredentialsRoleType {
		return fmt.Sprintf("shared_vpc_openshift_%s_policy", roleType)
	}

	return fmt.Sprintf("openshift_%s_policy", roleType)
}

// GetAccountRolePolicyKeys returns the policy key for fetching the managed policy ARN
func GetAccountRolePolicyKeys(roleType string) []string {
	if roleType == InstallerAccountRole {
		return []string{
			InstallerCoreKey,
			InstallerVPCKey,
			InstallerPrivateLinkKey,
		}
	}

	return []string{fmt.Sprintf("sts_%s_permission_policy", roleType)}
}

// GetAccountRolePolicyKeys returns the policy key for fetching the managed policy ARN
func GetHcpAccountRolePolicyKeys(roleType string) []string {
	policyKeys := []string{fmt.Sprintf("sts_hcp_%s_permission_policy", roleType)}

	return policyKeys
}

func ComputeOperatorRoleArn(prefix string, operator *cmv1.STSOperator, creator *Creator, path string) string {
	role := fmt.Sprintf("%s-%s-%s", prefix, operator.Namespace(), operator.Name())
	role = awsCommonUtils.TruncateRoleName(role)
	str := fmt.Sprintf("arn:%s:iam::%s:role", creator.Partition, creator.AccountID)
	if path != "" {
		str = fmt.Sprintf("%s%s", str, path)
		return fmt.Sprintf("%s%s", str, role)
	}
	return fmt.Sprintf("%s/%s", str, role)
}

func IsStandardNamedAccountRole(accountRoleName, roleSuffix string) (bool, string) {
	accountRolePrefix := TrimRoleSuffix(accountRoleName, fmt.Sprintf("-%s-Role", roleSuffix))
	return accountRolePrefix != accountRoleName, accountRolePrefix
}

func IsHostedCPManagedPolicies(cluster *cmv1.Cluster) bool {
	return cluster.Hypershift().Enabled() && cluster.AWS().STS().ManagedPolicies()
}

func IsHostedCP(cluster *cmv1.Cluster) bool {
	return cluster.Hypershift().Enabled()
}

func buildDescribeStacksInput(stackName string) *cloudformation.DescribeStacksInput {
	return &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}
}

func waitForStackCreateComplete(ctx context.Context, cfClient client.CloudFormationApiClient, stackName string) error {
	waiter := cloudformation.NewStackCreateCompleteWaiter(cfClient)

	params := buildDescribeStacksInput(stackName)

	// You can also use WaitForOutput if you need the output
	return waiter.Wait(ctx, params, maxWaitDur, func(o *cloudformation.StackCreateCompleteWaiterOptions) {
		// Optionally set MinDelay, MaxDelay, and other options here
	})
}

func waitForStackUpdateComplete(ctx context.Context, cfClient client.CloudFormationApiClient, stackName string) error {
	waiter := cloudformation.NewStackUpdateCompleteWaiter(cfClient)

	params := buildDescribeStacksInput(stackName)

	// You can also use WaitForOutput if you need the output
	return waiter.Wait(ctx, params, maxWaitDur, func(o *cloudformation.StackUpdateCompleteWaiterOptions) {
		// Optionally set MinDelay, MaxDelay, and other options here
	})

}

func waitForStackDeleteComplete(ctx context.Context, cfClient client.CloudFormationApiClient, stackName string) error {
	waiter := cloudformation.NewStackDeleteCompleteWaiter(cfClient)

	params := buildDescribeStacksInput(stackName)

	// You can also use WaitForOutput if you need the output
	return waiter.Wait(ctx, params, maxWaitDur, func(o *cloudformation.StackDeleteCompleteWaiterOptions) {
		// Optionally set MinDelay, MaxDelay, and other options here
	})
}

func sortAccountRolesByHCPSuffix(roles []iamtypes.Role) {
	sort.Slice(roles, func(i, j int) bool {
		prefixI := strings.SplitAfter(*roles[i].RoleName, "-")[0]
		prefixJ := strings.SplitAfter(*roles[j].RoleName, "-")[0]

		roleGroupI, errI := strconv.Atoi(prefixI)
		roleGroupJ, errJ := strconv.Atoi(prefixJ)

		if errI == nil && errJ == nil {
			if roleGroupI != roleGroupJ {
				return roleGroupI < roleGroupJ
			}
		} else {
			if prefixI != prefixJ {
				return prefixI < prefixJ
			}
		}

		hcpI := strings.Contains(*roles[i].RoleName, HCPSuffixPattern)
		hcpJ := strings.Contains(*roles[j].RoleName, HCPSuffixPattern)
		if hcpI != hcpJ {
			return hcpI
		}
		return roleGroupI < roleGroupJ
	})
}
