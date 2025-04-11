/**
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

package ocm

import (
	"crypto/x509"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	semver "github.com/hashicorp/go-version"
	common "github.com/openshift-online/ocm-common/pkg/ocm/validations"
	amsv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/output"
	"github.com/openshift/rosa/pkg/reporter"
)

const (
	ANY                 = "any"
	HibernateCapability = "capability.organization.hibernate_cluster"
	//Pendo Events
	Success             = "Success"
	Failure             = "Failure"
	Response            = "Response"
	ClusterID           = "ClusterID"
	OperatorRolesPrefix = "OperatorRolePrefix"
	Version             = "Version"
	Username            = "Username"
	URL                 = "URL"
	IsThrottle          = "IsThrottle"
	Account             = "Account"
	Organization        = "Organization"
	RoleName            = "RoleName"
	PolicyArn           = "PolicyArn"

	OCMRoleLabel  = "sts_ocm_role"
	USERRoleLabel = "sts_user_role"

	MaxClusterNameLength         = 54
	MaxClusterDomainPrefixLength = 15

	HcpProduct        = "hcp"
	HcpBillingAccount = "hcp-billing"
)

// Regular expression to used to make sure that the identifier or name given by the user is
// safe and that it there is no risk of SQL injection:
var clusterKeyRE = regexp.MustCompile(`^(\w|-)+$`)

var kubernetesLabelRE = regexp.MustCompile(`^[a-z0-9A-Z]+[-_.a-z0-9A-Z/]*$`)

// Cluster names must be valid DNS-1035 labels, so they must consist of at most 54 lower case alphanumeric
// characters or '-', start with an alphabetic character, and end with an alphanumeric character
var clusterNameRE = regexp.MustCompile(`^[a-z]([-a-z0-9]{0,52}[a-z0-9])?$`)

// Cluster domain prefix must be valid DNS-1035 labels, so they must consist of at most 15 lower case alphanumeric
// characters or '-', start with an alphabetic character, and end with an alphanumeric character
var clusterDomainPrefixRE = regexp.MustCompile(`^[a-z]([-a-z0-9]{0,13}[a-z0-9])?$`)

var badUsernameRE = regexp.MustCompile(`^(~|\.?\.|.*[:\/%].*)$`)

// math.MaxInt64 needs to be bound by int64() as that can otherwise overflow
// if the helpers are compiled for non-64bit architectures (e.g. armv7l)
// as per https://github.com/golang/go/issues/23086.
var maxInt = int64(math.MaxInt64)

func IsValidClusterKey(clusterKey string) bool {
	return clusterKeyRE.MatchString(clusterKey)
}

func IsValidClusterName(clusterName string) bool {
	return clusterNameRE.MatchString(clusterName)
}

func ClusterNameValidator(name interface{}) error {
	if str, ok := name.(string); ok {
		str := strings.Trim(str, " \t")
		if !IsValidClusterName(str) {
			return fmt.Errorf("Cluster name must consist of no more than %d lowercase "+
				"alphanumeric characters or '-', start with a letter, and end with an "+
				"alphanumeric character.", MaxClusterNameLength)
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got '%v'", name)
}

func IsValidClusterDomainPrefix(domainPrefix string) bool {
	return clusterDomainPrefixRE.MatchString(domainPrefix)
}

func ClusterDomainPrefixValidator(domainPrefix interface{}) error {
	if str, ok := domainPrefix.(string); ok {
		str := strings.Trim(str, " \t")
		if str != "" && !IsValidClusterDomainPrefix(str) {
			return fmt.Errorf("Cluster domain prefix must consist of no more than %d lowercase "+
				"alphanumeric characters or '-', start with a letter, and end with an "+
				"alphanumeric character.", MaxClusterDomainPrefixLength)
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got '%v'", domainPrefix)
}

func ValidateHTTPProxy(val interface{}) error {
	if httpProxy, ok := val.(string); ok {
		if httpProxy == "" {
			return nil
		}
		url, err := url.ParseRequestURI(httpProxy)
		if err != nil {
			return fmt.Errorf("Invalid http-proxy value '%s'", httpProxy)
		}
		if url.Scheme != "http" {
			return errors.Errorf("%s", "Expected http-proxy to have an http:// scheme")
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got '%v'", val)
}

func ValidateRegistryAdditionalCa(input map[string]string) error {
	var caRE = regexp.MustCompile(
		`(-----BEGIN CERTIFICATE-----\n)([^-----]*)(-----END CERTIFICATE-----)`)
	for _, ca := range input {
		if !caRE.MatchString(ca) {
			return fmt.Errorf("invalid PEM-encoded certificate.' ")
		}
	}
	return nil
}

func ValidateAllowedRegistriesForImport(input interface{}) error {
	var idRE = regexp.MustCompile(`^(.+):(true|false)$`)
	var registries []string
	inputType := reflect.TypeOf(input).Kind()
	switch inputType {
	case reflect.String:
		if input.(string) == "" {
			return nil
		}
		registries = strings.Split(input.(string), ",")
		for _, registry := range registries {
			if !idRE.MatchString(registry) {
				return fmt.Errorf("invalid identifier '%s' for 'allowed registries for import.' "+
					"Should be in a <registry>:<boolean> format. "+
					"The boolean indicates whether the registry is secure or not.", registry)
			}
		}
	case reflect.Slice:
		if reflect.TypeOf(input).Elem().Kind() != reflect.String {
			return fmt.Errorf("unable to verify allowed element, incompatible type, expected slice of string got: '%s'",
				inputType.String())
		}
	default:
		return fmt.Errorf("can only validate string types, got %v", inputType.String())
	}
	return nil
}

func ValidateAdditionalTrustBundle(val interface{}) error {
	if additionalTrustBundleFile, ok := val.(string); ok {
		if additionalTrustBundleFile == "" {
			return nil
		}
		cert, err := os.ReadFile(additionalTrustBundleFile)
		if err != nil {
			return err
		}
		additionalTrustBundle := string(cert)
		if additionalTrustBundle == "" {
			return errors.Errorf("%s", "Trust bundle file is empty")
		}
		additionalTrustBundleBytes := []byte(additionalTrustBundle)
		if !x509.NewCertPool().AppendCertsFromPEM(additionalTrustBundleBytes) {
			return errors.Errorf("%s", "Failed to parse additional trust bundle")
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got '%v'", val)
}

func IsValidUsername(username string) bool {
	return !badUsernameRE.MatchString(username)
}

func IsEmptyCIDR(cidr net.IPNet) bool {
	return cidr.String() == "<nil>"
}

// Determine whether a resources is compatible with ROSA clusters in general
func isCompatible(relatedResource *amsv1.RelatedResource) bool {
	product := strings.ToLower(relatedResource.Product())
	cloudProvider := strings.ToLower(relatedResource.CloudProvider())
	byoc := strings.ToLower(relatedResource.BYOC())

	// nolint:goconst
	return (product == ANY || product == "rosa" || product == "moa") &&
		(cloudProvider == ANY || cloudProvider == "aws") &&
		(byoc == ANY || byoc == "byoc")
}

func handleErr(res *ocmerrors.Error, err error) error {
	msg := err.Error()
	// Need to check nil as .Error will try to access internal values of the pointer
	if res != nil && res.Reason() != "" {
		msg = res.Error()
	}
	// This works because the following gets return zero value for their type
	// string->""
	// int->0
	// Hack to always display the correct terms and conditions message
	if res.Code() == "CLUSTERS-MGMT-451" {
		msg = "You must accept the Terms and Conditions in order to continue.\n" +
			"Go to https://www.redhat.com/wapps/tnc/ackrequired?site=ocm&event=register\n" +
			"Once you accept the terms, you will need to retry the action that was blocked."
	}
	// The error type set will be No Type though
	errType := errors.ErrorType(res.Status())
	return errType.Set(errors.Errorf("%s", msg))
}

func (c *Client) GetDefaultClusterFlavors(flavour string) (dMachinecidr *net.IPNet, dPodcidr *net.IPNet,
	dServicecidr *net.IPNet, dhostPrefix, defaultMachineRootVolumeSize int, computeInstanceType string) {
	flavourGetResponse, err := c.ocm.ClustersMgmt().V1().Flavours().Flavour(flavour).Get().Send()
	if err != nil {
		flavourGetResponse, _ = c.ocm.ClustersMgmt().V1().Flavours().Flavour("osd-4").Get().Send()
	}
	aws, ok := flavourGetResponse.Body().GetAWS()
	if !ok {
		return nil, nil, nil, 0, 0, ""
	}
	computeInstanceType = aws.ComputeInstanceType()
	network, ok := flavourGetResponse.Body().GetNetwork()
	if !ok {
		return nil, nil, nil, 0, 0, computeInstanceType
	}
	_, dMachinecidr, err = net.ParseCIDR(network.MachineCIDR())
	if err != nil {
		dMachinecidr = nil
	}
	_, dPodcidr, err = net.ParseCIDR(network.PodCIDR())
	if err != nil {
		dPodcidr = nil
	}
	_, dServicecidr, err = net.ParseCIDR(network.ServiceCIDR())
	if err != nil {
		dServicecidr = nil
	}
	dhostPrefix, _ = network.GetHostPrefix()

	// default machine root volume size
	defaultMachineRootVolumeSize = aws.WorkerVolume().Size()

	return dMachinecidr, dPodcidr, dServicecidr, dhostPrefix, defaultMachineRootVolumeSize, computeInstanceType
}

func (c *Client) LogEvent(key string, body map[string]string) {
	event, err := cmv1.NewEvent().Key(key).Body(body).Build()
	if err == nil {
		_, _ = c.ocm.ClustersMgmt().V1().
			Events().
			Add().
			Body(event).
			Send()
	}
}

func (c *Client) GetCurrentAccount() (*amsv1.Account, error) {
	response, err := c.ocm.AccountsMgmt().V1().
		CurrentAccount().
		Get().
		Send()
	if err != nil {
		if response.Status() == http.StatusNotFound {
			return nil, nil
		}
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) GetCurrentOrganization() (id string, externalID string, err error) {
	acctResponse, err := c.GetCurrentAccount()

	if err != nil {
		return
	}
	id = acctResponse.Organization().ID()
	externalID = acctResponse.Organization().ExternalID()

	return
}

func (c *Client) IsCapabilityEnabled(capability string) (enabled bool, err error) {
	organizationID, _, err := c.GetCurrentOrganization()
	if err != nil {
		return
	}
	isCapabilityEnable, err := c.isCapabilityEnabled(capability, organizationID)

	if err != nil {
		return
	}
	if !isCapabilityEnable {
		return false, nil
	}
	return true, nil
}

func (c *Client) isCapabilityEnabled(capabilityName string, orgID string) (bool, error) {
	capabilityResponse, err := c.ocm.AccountsMgmt().V1().Organizations().
		Organization(orgID).Get().Parameter("fetchCapabilities", true).Send()

	if err != nil {
		return false, handleErr(capabilityResponse.Error(), err)
	}
	if len(capabilityResponse.Body().Capabilities()) > 0 {
		for _, capability := range capabilityResponse.Body().Capabilities() {
			if capability.Name() == capabilityName {
				return capability.Value() == "true", nil
			}
		}
	}
	return false, nil
}

func (c *Client) UnlinkUserRoleFromAccount(accountID string, roleARN string) error {
	linkedRoles, err := c.GetAccountLinkedUserRoles(accountID)
	if err != nil {
		return err
	}

	if helper.Contains(linkedRoles, roleARN) {
		linkedRoles = helper.RemoveStrFromSlice(linkedRoles, roleARN)

		if len(linkedRoles) > 0 {
			newRoleARN := strings.Join(linkedRoles, ",")
			label, err := amsv1.NewLabel().Key(USERRoleLabel).Value(newRoleARN).Build()
			if err != nil {
				return err
			}

			resp, err := c.ocm.AccountsMgmt().V1().Accounts().Account(accountID).Labels().
				Labels(USERRoleLabel).Update().Body(label).Send()
			if err != nil {
				return handleErr(resp.Error(), err)
			}
		} else {
			resp, err := c.ocm.AccountsMgmt().V1().Accounts().Account(accountID).Labels().
				Labels(USERRoleLabel).Delete().Send()
			if err != nil {
				return handleErr(resp.Error(), err)
			}
		}

		return nil
	}

	return errors.UserErrorf("Role ARN '%s' is not linked with the current account '%s'", roleARN, accountID)
}

func (c *Client) LinkAccountRole(accountID string, roleARN string) error {
	resp, err := c.ocm.AccountsMgmt().V1().Accounts().Account(accountID).
		Labels().Labels("sts_user_role").Get().Send()
	if err != nil && resp.Status() != 404 {
		if resp.Status() == 403 {
			return errors.Forbidden.UserErrorf("%v", err)
		}
		return handleErr(resp.Error(), err)
	}
	existingARN := resp.Body().Value()
	exists := false
	if existingARN != "" {
		existingARNArr := strings.Split(existingARN, ",")
		if len(existingARNArr) > 0 {
			for _, value := range existingARNArr {
				if value == roleARN {
					exists = true
					break
				}
			}
		}
	}
	if exists {
		return nil
	}
	if existingARN != "" {
		roleARN = existingARN + "," + roleARN
	}
	labelBuilder, err := amsv1.NewLabel().Key("sts_user_role").Value(roleARN).Build()
	if err != nil {
		return err
	}
	_, err = c.ocm.AccountsMgmt().V1().Accounts().Account(accountID).
		Labels().Add().Body(labelBuilder).Send()
	if err != nil {
		return handleErr(resp.Error(), err)
	}
	return err
}

func (c *Client) UnlinkOCMRoleFromOrg(orgID string, roleARN string) error {
	linkedRoles, err := c.GetOrganizationLinkedOCMRoles(orgID)
	if err != nil {
		return err
	}

	if helper.Contains(linkedRoles, roleARN) {
		linkedRoles = helper.RemoveStrFromSlice(linkedRoles, roleARN)

		if len(linkedRoles) > 0 {
			newRoleARN := strings.Join(linkedRoles, ",")
			label, err := amsv1.NewLabel().Key(OCMRoleLabel).Value(newRoleARN).Build()
			if err != nil {
				return err
			}

			resp, err := c.ocm.AccountsMgmt().V1().Organizations().Organization(orgID).Labels().
				Labels(OCMRoleLabel).Update().Body(label).Send()
			if err != nil {
				return handleErr(resp.Error(), err)
			}
		} else {
			resp, err := c.ocm.AccountsMgmt().V1().Organizations().Organization(orgID).Labels().
				Labels(OCMRoleLabel).Delete().Send()
			if err != nil {
				return handleErr(resp.Error(), err)
			}
		}

		return nil
	}

	return errors.UserErrorf("Role-arn '%s' is not linked with the organization account '%s'", roleARN, orgID)
}

func (c *Client) LinkOrgToRole(orgID string, roleARN string) (bool, error) {
	parsedARN, err := arn.Parse(roleARN)
	if err != nil {
		return false, err
	}
	exists, existingARN, selectedARN, err := c.CheckIfAWSAccountExists(orgID, parsedARN.AccountID)
	if err != nil {
		return false, err
	}
	if exists {
		if selectedARN != roleARN {
			return false, errors.UserErrorf("User organization '%s' has role-arn '%s' associated. "+
				"Only one role can be linked per AWS account per organization", orgID, selectedARN)
		}
		return false, nil
	}
	if existingARN != "" {
		roleARN = existingARN + "," + roleARN
	}
	labelBuilder, err := amsv1.NewLabel().Key(OCMRoleLabel).Value(roleARN).Build()
	if err != nil {
		return false, err
	}

	resp, err := c.ocm.AccountsMgmt().V1().Organizations().Organization(orgID).
		Labels().Add().Body(labelBuilder).Send()
	if err != nil {
		return false, handleErr(resp.Error(), err)
	}
	return true, nil
}

func (c *Client) GetAccountLinkedUserRoles(accountID string) ([]string, error) {
	resp, err := c.ocm.AccountsMgmt().V1().Accounts().Account(accountID).
		Labels().Labels(USERRoleLabel).Get().Send()
	if err != nil && resp.Status() != http.StatusNotFound {
		return nil, handleErr(resp.Error(), err)
	}

	return strings.Split(resp.Body().Value(), ","), nil
}

func (c *Client) GetOrganizationLinkedOCMRoles(orgID string) ([]string, error) {
	resp, err := c.ocm.AccountsMgmt().V1().Organizations().Organization(orgID).
		Labels().Labels(OCMRoleLabel).Get().Send()
	if err != nil && resp.Status() != http.StatusNotFound {
		return nil, err
	}

	return strings.Split(resp.Body().Value(), ","), nil
}

func (c *Client) CheckIfAWSAccountExists(orgID string, awsAccountID string) (bool, string, string, error) {
	resp, err := c.ocm.AccountsMgmt().V1().Organizations().Organization(orgID).
		Labels().Labels(OCMRoleLabel).Get().Send()
	if err != nil && resp.Status() != 404 {
		if resp.Status() == 403 {
			return false, "", "", errors.Forbidden.UserErrorf("%v", err)
		}
		return false, "", "", handleErr(resp.Error(), err)
	}
	existingARN := resp.Body().Value()
	exists := false
	selectedARN := ""
	if existingARN != "" {
		existingARNArr := strings.Split(existingARN, ",")
		if len(existingARNArr) > 0 {
			for _, value := range existingARNArr {
				parsedARN, err := arn.Parse(value)
				if err != nil {
					return false, "", "", err
				}
				if parsedARN.AccountID == awsAccountID {
					exists = true
					selectedARN = value
					break
				}
			}
		}
	}
	return exists, existingARN, selectedARN, nil
}

/*
We should allow only one role per aws account per organization
If the user request same ocm role we should let them proceed to ensure they can add admin role
if not exists or attach policies or link etc
if the user request diff ocm role name we error out
*/
func (c *Client) CheckRoleExists(orgID string, roleName string, awsAccountID string) (bool, string, string, error) {
	exists, _, selectedARN, err := c.CheckIfAWSAccountExists(orgID, awsAccountID)
	if err != nil {
		return false, "", "", err
	}
	if !exists {
		return false, "", "", nil
	}
	existingRole := strings.SplitN(selectedARN, "/", 2)
	if len(existingRole) > 1 && existingRole[1] == roleName {
		return false, "", "", nil
	}
	return true, existingRole[1], selectedARN, nil
}

func GetVersionMinor(ver string) string {
	rawID := strings.Replace(ver, VersionPrefix, "", 1)
	version, err := semver.NewVersion(rawID)
	if err != nil {
		segments := strings.Split(rawID, ".")
		return fmt.Sprintf("%s.%s", segments[0], segments[1])
	}
	segments := version.Segments()
	return fmt.Sprintf("%d.%d", segments[0], segments[1])
}

func CheckSupportedVersion(clusterVersion string, operatorVersion string) (bool, error) {
	v1, err := semver.NewVersion(clusterVersion)
	if err != nil {
		return false, err
	}
	v2, err := semver.NewVersion(operatorVersion)
	if err != nil {
		return false, err
	}
	//Cluster version is greater than or equal to operator version
	return v1.GreaterThanOrEqual(v2), nil
}

func (c *Client) GetPolicies(policyType string) (map[string]*cmv1.AWSSTSPolicy, error) {

	query := fmt.Sprintf("policy_type = '%s'", policyType)

	m := make(map[string]*cmv1.AWSSTSPolicy)

	stmt := c.ocm.ClustersMgmt().V1().AWSInquiries().STSPolicies().List()
	if policyType != "" {
		stmt = stmt.Search(query)
	}
	accountRolePoliciesResponse, err := stmt.Send()
	if err != nil {
		return m, handleErr(accountRolePoliciesResponse.Error(), err)
	}
	accountRolePoliciesResponse.Items().Each(func(awsPolicy *cmv1.AWSSTSPolicy) bool {
		m[awsPolicy.ID()] = awsPolicy
		return true
	})
	return m, nil
}

// The actual values might differ from classic to hcp
// prefer using GetCredRequests(isHypershift bool) when there is prior knowledge of the topology
func (c *Client) GetAllCredRequests() (map[string]*cmv1.STSOperator, error) {
	result := make(map[string]*cmv1.STSOperator)
	classic, err := c.GetCredRequests(false)
	if err != nil {
		return result, err
	}
	hcp, err := c.GetCredRequests(true)
	if err != nil {
		return result, err
	}
	for key, value := range classic {
		result[key] = value
	}
	for key, value := range hcp {
		result[key] = value
	}
	return result, nil
}

func (c *Client) GetCredRequests(isHypershift bool) (map[string]*cmv1.STSOperator, error) {
	m := make(map[string]*cmv1.STSOperator)
	stsCredentialResponse, err := c.ocm.ClustersMgmt().
		V1().
		AWSInquiries().
		STSCredentialRequests().
		List().
		Parameter("is_hypershift", isHypershift).
		Send()
	if err != nil {
		return m, handleErr(stsCredentialResponse.Error(), err)
	}

	stsCredentialResponse.Items().Each(func(stsCredentialRequest *cmv1.STSCredentialRequest) bool {
		m[stsCredentialRequest.Name()] = stsCredentialRequest.Operator()
		return true
	})
	return m, nil
}

func (c *Client) FindMissingOperatorRolesForUpgrade(cluster *cmv1.Cluster,
	newMinorVersion string) (map[string]*cmv1.STSOperator, error) {
	missingRoles := make(map[string]*cmv1.STSOperator)

	credRequests, err := c.GetCredRequests(cluster.Hypershift().Enabled())
	if err != nil {
		return nil, errors.Errorf("Error getting operator credential request from OCM %s", err)
	}
	for credRequest, operator := range credRequests {
		if operator.MinVersion() != "" {
			clusterUpgradeVersion, err := semver.NewVersion(newMinorVersion)
			if err != nil {
				return nil, err
			}
			operatorMinVersion, err := semver.NewVersion(operator.MinVersion())
			if err != nil {
				return nil, err
			}

			if clusterUpgradeVersion.GreaterThanOrEqual(operatorMinVersion) {
				if !isOperatorRoleAlreadyExist(cluster, operator) {
					missingRoles[credRequest] = operator
				}
			}
		}
	}
	return missingRoles, nil
}

func (c *Client) createCloudProviderDataBuilder(roleARN string, awsClient aws.Client,
	externalID string) (*cmv1.CloudProviderDataBuilder, error) {
	var awsBuilder *cmv1.AWSBuilder
	if roleARN != "" {
		stsBuilder := cmv1.NewSTS().RoleARN(roleARN)
		if externalID != "" {
			stsBuilder = stsBuilder.ExternalID(externalID)
		}
		awsBuilder = cmv1.NewAWS().STS(stsBuilder)
	} else {
		accessKeys, err := awsClient.GetAWSAccessKeys()
		if err != nil {
			return &cmv1.CloudProviderDataBuilder{}, err
		}
		awsBuilder = cmv1.NewAWS().AccessKeyID(accessKeys.AccessKeyID).SecretAccessKey(accessKeys.SecretAccessKey)
	}

	return cmv1.NewCloudProviderData().AWS(awsBuilder), nil
}

func isOperatorRoleAlreadyExist(cluster *cmv1.Cluster, operator *cmv1.STSOperator) bool {
	for _, role := range cluster.AWS().STS().OperatorIAMRoles() {
		//FIXME: Check it does not exist on AWS itself too
		// the iam roles will only return up to the version of the cluster
		if role.Namespace() == operator.Namespace() && role.Name() == operator.Name() {
			return true
		}
	}
	return false
}

const (
	BYOVPCSingleAZSubnetsCount      = 2
	BYOVPCMultiAZSubnetsCount       = 6
	privateLinkSingleAZSubnetsCount = 1
	privateLinkMultiAZSubnetsCount  = 3
)

func expectedSubnetsCount(multiAZ, privateLink bool) int {
	return map[bool]map[bool]int{
		true: {
			true:  privateLinkMultiAZSubnetsCount,
			false: BYOVPCMultiAZSubnetsCount,
		},
		false: {
			true:  privateLinkSingleAZSubnetsCount,
			false: BYOVPCSingleAZSubnetsCount,
		},
	}[multiAZ][privateLink]
}

func ValidateSubnetsCount(multiAZ bool, privateLink bool, subnetsInputCount int) error {
	expected := expectedSubnetsCount(multiAZ, privateLink)
	if subnetsInputCount != expected {
		clusterPrefix := "cluster"
		if privateLink {
			clusterPrefix = "private link cluster"
		}
		azPrefix := "single AZ"
		if multiAZ {
			azPrefix = "multi-AZ"
		}
		return fmt.Errorf("The number of subnets for a '%s' '%s' should be '%d', "+
			"instead received: '%d'", azPrefix, clusterPrefix, expected, subnetsInputCount)
	}
	return nil
}

func ValidateHostedClusterSubnets(awsClient aws.Client, isPrivate bool, subnetIDs []string) (int, error) {
	if isPrivate && len(subnetIDs) < 1 {
		return 0, fmt.Errorf("The number of subnets for a private hosted cluster should be at least one")
	}
	if !isPrivate && len(subnetIDs) < 2 {
		return 0, fmt.Errorf("The number of subnets for a public hosted cluster should be at least two")
	}
	vpcSubnets, vpcSubnetsErr := awsClient.GetVPCSubnets(subnetIDs[0])
	if vpcSubnetsErr != nil {
		return 0, vpcSubnetsErr
	}

	var subnets []ec2types.Subnet
	for _, subnet := range vpcSubnets {
		for _, subnetId := range subnetIDs {
			if awssdk.ToString(subnet.SubnetId) == subnetId {
				subnets = append(subnets, subnet)
				break
			}
		}
	}

	privateSubnets, privateSubnetsErr := awsClient.FilterVPCsPrivateSubnets(subnets)
	if privateSubnetsErr != nil {
		return 0, privateSubnetsErr
	}

	privateSubnetCount := len(privateSubnets)
	publicSubnetsCount := len(subnets) - privateSubnetCount

	if isPrivate {
		if publicSubnetsCount > 0 {
			return 0, fmt.Errorf("The number of public subnets for a private hosted cluster should be zero")
		}
	} else {
		if publicSubnetsCount == 0 {
			return 0, fmt.Errorf("The number of public subnets for a public hosted " +
				"cluster should be at least one")
		}
	}
	return privateSubnetCount, nil
}

func (c *Client) CheckUpgradeClusterVersion(
	availableUpgrades []string,
	clusterUpgradeVersion string,
	cluster *cmv1.Cluster,
) (err error) {
	clusterVersion := cluster.OpenshiftVersion()
	if clusterVersion == "" {
		clusterVersion = cluster.Version().RawID()
	}
	validVersion := false
	for _, v := range availableUpgrades {
		isValidVersion, err := IsValidVersion(clusterUpgradeVersion, v, clusterVersion)
		if err != nil {
			return err
		}
		if isValidVersion {
			validVersion = true
			break
		}
	}
	if !validVersion {
		return errors.Errorf(
			"Expected a valid version to upgrade cluster to.\nValid versions: %s",
			helper.SliceToSortedString(availableUpgrades),
		)
	}
	return nil
}

func (c *Client) GetPolicyVersion(userRequestedVersion string, channelGroup string) (string, error) {
	if userRequestedVersion == "" {
		version, err := c.GetLatestVersion(channelGroup)
		if err != nil {
			return userRequestedVersion, err
		}
		return version, nil
	}

	versionList, err := c.GetVersionsList(channelGroup, false)
	if err != nil {
		err := fmt.Errorf("%v", err)
		return userRequestedVersion, err
	}

	hasVersion := false
	for _, vs := range versionList {
		if vs == userRequestedVersion {
			hasVersion = true
			break
		}
	}

	if !hasVersion {
		versionSet := helper.SliceToMap(versionList)
		err := errors.Errorf(
			"A valid policy version number must be specified\nValid versions: %v",
			helper.MapKeysToString(versionSet),
		)
		return userRequestedVersion, err
	}

	return userRequestedVersion, nil
}

func ParseVersion(version string) (string, error) {
	parsedVersion, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	versionSplit := parsedVersion.Segments64()
	return fmt.Sprintf("%d.%d", versionSplit[0], versionSplit[1]), nil
}

func (c *Client) GetVersionsList(channelGroup string, defaultFirst bool) ([]string, error) {
	response, err := c.GetVersions(channelGroup, defaultFirst)
	if err != nil {
		err := fmt.Errorf("error getting versions: %s", err)
		return make([]string, 0), err
	}
	versionList := make([]string, 0)
	for _, v := range response {
		if !HasSTSSupport(v.RawID(), v.ChannelGroup()) {
			continue
		}
		parsedVersion, err := ParseVersion(v.RawID())
		if err != nil {
			err = fmt.Errorf("error parsing version")
			return versionList, err
		}
		versionList = append(versionList, parsedVersion)
	}

	if len(versionList) == 0 {
		err = fmt.Errorf("could not find versions for the provided channel-group: '%s'", channelGroup)
		return versionList, err
	}
	return versionList, nil
}

func ValidateOperatorRolesMatchOidcProvider(reporter *reporter.Object, awsClient aws.Client,
	operatorIAMRoleList []OperatorIAMRole, oidcEndpointUrl string,
	clusterVersion string, expectedOperatorRolePath string,
	accountRolesHasManagedPolicies bool, logOperatorRoles bool) error {
	operatorIAMRoles := operatorIAMRoleList
	parsedUrl, err := url.Parse(oidcEndpointUrl)
	if err != nil {
		return err
	}
	printInfo := reporter.IsTerminal() && !output.HasFlag() && logOperatorRoles
	if printInfo {
		reporter.Infof("Reusable OIDC Configuration detected. Validating trusted relationships to operator roles: ")
	}
	for _, operatorIAMRole := range operatorIAMRoles {
		roleObject, err := awsClient.GetRoleByARN(operatorIAMRole.RoleARN)
		if err != nil {
			return err
		}
		roleARN := *roleObject.Arn
		pathFromArn, err := aws.GetPathFromARN(roleARN)
		if err != nil {
			return err
		}
		if pathFromArn != expectedOperatorRolePath {
			return errors.Errorf("Operator Role '%s' does not match the path from Installer Role, "+
				"please choose correct Installer Role and try again.", roleARN)
		}
		if roleARN != operatorIAMRole.RoleARN {
			return errors.Errorf("Computed Operator Role '%s' does not match role ARN found in AWS '%s', "+
				"please check if the correct parameters have been supplied.", operatorIAMRole.RoleARN, roleARN)
		}
		err = common.ValidateIssuerUrlMatchesAssumePolicyDocument(
			roleARN, parsedUrl, *roleObject.AssumeRolePolicyDocument)
		if err != nil {
			return err
		}
		hasManagedPolicies, err := awsClient.HasManagedPolicies(roleARN)
		if err != nil {
			return err
		}
		if accountRolesHasManagedPolicies && !hasManagedPolicies {
			return errors.Errorf("Operator role '%s' has unmanaged policies and is not compatible with the account "+
				"role's managed policies.", roleARN)
		}
		if printInfo {
			reporter.Infof("Using '%s'", roleARN)
		}
		if hasManagedPolicies {
			// Managed policies should be compatible with all versions
			continue
		}
		policiesDetails, err := awsClient.GetAttachedPolicy(roleObject.RoleName)
		if err != nil {
			return err
		}
		for _, policyDetails := range policiesDetails {
			if policyDetails.PolicyType == aws.Inline {
				continue
			}
			isCompatible, err := awsClient.IsPolicyCompatible(policyDetails.PolicyArn, clusterVersion)
			if err != nil {
				return err
			}
			if !isCompatible {
				return errors.Errorf(
					"Operator role '%s' is not compatible with cluster version '%s'",
					roleARN,
					clusterVersion,
				)
			}
		}
	}
	return nil
}

func ValidateHttpTokensValue(val interface{}) error {
	if httpTokens, ok := val.(string); ok {
		if httpTokens == "" {
			return nil
		}
		switch cmv1.Ec2MetadataHttpTokens(httpTokens) {
		case cmv1.Ec2MetadataHttpTokensRequired, cmv1.Ec2MetadataHttpTokensOptional:
			return nil
		default:
			return errors.Errorf("ec2-metadata-http-tokens value should be one of '%s', '%s'",
				cmv1.Ec2MetadataHttpTokensRequired, cmv1.Ec2MetadataHttpTokensOptional)
		}
	}

	return fmt.Errorf("can only validate strings, got '%v'", val)
}

func ParseDiskSizeToGigibyte(size string) (int, error) {
	// Empty string is valid, a default will be set later
	if size == "" {
		return 0, nil
	}

	suffixErrorString := "accepted units are Giga or Tera in the form of g, G, GB, GiB, Gi, t, T, TB, TiB, Ti"

	// Remove spaces if the value is '100 GB'
	size = strings.ReplaceAll(size, " ", "")

	// Do a bit of cleaning to avoid errors from the resources quantity parser
	// Convert units using regular expressions, units on the left will be converted to the unit on
	// the right if matched
	unitToConvert := map[string]string{
		"GB|gb|Gb|g":      "G",
		"gib|GIB|GiB|Gib": "Gi",
		"TB|tb|Tb|t":      "T",
		"tib|TIB|TiB|Tib": "Ti",
	}

	for badUnit, rightUnit := range unitToConvert {
		re, err := regexp.Compile(badUnit)
		if err != nil {
			return 0, err
		}
		size = re.ReplaceAllString(size, rightUnit)
	}

	qty, err := resource.ParseQuantity(size)
	if err != nil {
		if err == resource.ErrFormatWrong {
			return 0, fmt.Errorf("invalid disk size format: '%s'. %s", size, suffixErrorString)
		}
		return 0, fmt.Errorf("invalid disk size: '%s'. %w", size, err)
	}

	if qty.Value() < 0 {
		return 0, fmt.Errorf("invalid disk size: '%s'. positive size required", size)
	}

	// If the value is 0, this could mean the user forgot to specify the unit suffix
	if qty.IsZero() {
		return 0, nil
	}

	// resource.ParseQuantity() will not error when a value exceeds the max int64.
	if qty.Value() == maxInt {
		return 0, fmt.Errorf("invalid disk size: '%s'. maximum size exceeded", size)
	}

	// Check the suffix is correct
	diskSize := qty.String()

	// If we can convert the string to a number, it means there is no suffix, thus the user passes
	// bytes and forgot to include a suffix
	if diskSizeInt, err := strconv.Atoi(diskSize); err == nil {
		// If the conversion hit the maximum value of a 64 bit integer, it means the user passed a very
		// large value which is ok, we can still proceed and the backend will return an error since
		// the size is too large
		if int64(diskSizeInt) != maxInt {
			return 0, fmt.Errorf("missing unit suffix: '%s'. %s", diskSize, suffixErrorString)
		}
	}

	// If the suffix is too small, it means the user passed a very small value
	if strings.HasSuffix(diskSize, "M") ||
		strings.HasSuffix(diskSize, "Mi") ||
		strings.HasSuffix(diskSize, "K") ||
		strings.HasSuffix(diskSize, "Ki") {
		return 0, fmt.Errorf("invalid disk size format: '%s'. %s", diskSize, suffixErrorString)
	}

	// Return gibibytes since the AWS expects that format
	// qty.Value() returns the value in bytes
	return int(qty.Value() / 1024 / 1024 / 1024), nil
}

func ValidateBalancingIgnoredLabels(val interface{}) error {
	labelsInput, ok := val.(string)

	if !ok {
		return fmt.Errorf("can only validate strings, got '%v'", val)
	}

	labels := strings.Split(labelsInput, ",")

	for _, label := range labels {
		if label == "" {
			continue
		}

		if len(label) > 63 {
			return fmt.Errorf("label key '%v' has has exceeded allowed label length of 63 characters", label)
		}

		if !kubernetesLabelRE.MatchString(label) {
			return fmt.Errorf("label '%v' is not a valid Kubernetes label key. "+
				"It must start with an alphanumeric character and may additionally contain only forward-slashes, "+
				"dashes, underscores, and dots",
				label)
		}
	}

	return nil
}

func ValidateClaimValidationRules(input interface{}) error {
	var idRE = regexp.MustCompile(`^[0-9a-z]+([:][0-9a-z]+)`)
	var inputRules []string
	inputType := reflect.TypeOf(input).Kind()
	switch inputType {
	case reflect.String:
		if input.(string) == "" {
			return nil
		}
		inputRules = strings.Split(input.(string), ",")
		for _, inputRule := range inputRules {
			if !idRE.MatchString(inputRule) {
				return fmt.Errorf("invalid identifier '%s' for 'claim validation rule. '"+
					"Should be in a <claim>:<required_value> format.", inputRule)
			}
		}
	case reflect.Slice:
		if reflect.TypeOf(input).Elem().Kind() != reflect.String {
			return fmt.Errorf(
				"unable to verify claim validation rules, incompatible type, expected slice of string got: '%s'",
				inputType.String(),
			)
		}
	default:
		return fmt.Errorf("can only validate string types, got %v", inputType.String())
	}
	return nil
}
