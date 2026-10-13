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
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	idputils "github.com/openshift-online/ocm-common/pkg/idp/utils"
	ocmConsts "github.com/openshift-online/ocm-common/pkg/ocm/consts"
	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/fedramp"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/info"
	"github.com/openshift/rosa/pkg/interactive/consts"
	"github.com/openshift/rosa/pkg/logforwarding"
	"github.com/openshift/rosa/pkg/properties"
	rprtr "github.com/openshift/rosa/pkg/reporter"
)

const (
	legacyIngressSupportLabel = "ext-managed.openshift.io/legacy-ingress-support"
)

var NetworkTypes = []string{"OpenShiftSDN", "OVNKubernetes"}

type DefaultIngressSpec struct {
	RouteSelectors           map[string]string
	ExcludedNamespaces       []string
	WildcardPolicy           string
	NamespaceOwnershipPolicy string
}

func NewDefaultIngressSpec() DefaultIngressSpec {
	defaultIngressSpec := DefaultIngressSpec{}
	defaultIngressSpec.RouteSelectors = map[string]string{}
	defaultIngressSpec.ExcludedNamespaces = []string{}
	return defaultIngressSpec
}

// Spec is the configuration for a cluster spec.
type Spec struct {
	// Basic configs
	Name                      string
	DomainPrefix              string
	Region                    string
	MultiAZ                   bool
	Version                   string
	ChannelGroup              string
	Expiration                time.Time
	Flavour                   string
	DisableWorkloadMonitoring *bool

	//Encryption
	FIPS                 bool
	EtcdEncryption       bool
	KMSKeyArn            string
	EtcdEncryptionKMSArn string
	// Scaling config
	ComputeMachineType string
	ComputeNodes       int
	Autoscaling        bool
	AutoscalerConfig   *AutoscalerConfig
	MinReplicas        int
	MaxReplicas        int
	ComputeLabels      map[string]string

	// SubnetIDs
	SubnetIds []string

	// AvailabilityZones
	AvailabilityZones []string

	// Network config
	NetworkType                    string
	SubnetConfiguration            string
	OvnInternalSubnetConfiguration map[string]string
	MachineCIDR                    net.IPNet
	ServiceCIDR                    net.IPNet
	PodCIDR                        net.IPNet
	HostPrefix                     int
	Private                        *bool
	PrivateLink                    *bool
	PrivateIngress                 *bool

	// Properties
	CustomProperties map[string]string

	// User-defined tags for AWS resources
	Tags map[string]string

	// Simulate creating a cluster but don't actually create it
	DryRun *bool

	// Disable SCP checks in the installer by setting credentials mode as mint
	DisableSCPChecks *bool

	// Non-STS
	AWSAccessKey *aws.AccessKey
	AWSCreator   *aws.Creator

	// STS
	IsSTS               bool
	RoleARN             string
	ExternalID          string
	SupportRoleARN      string
	OperatorIAMRoles    []OperatorIAMRole
	ControlPlaneRoleARN string
	WorkerRoleARN       string
	OidcConfigId        string
	Mode                string

	// External authentication configuration
	ExternalAuthProvidersEnabled bool

	NodeDrainGracePeriodInMinutes float64

	EnableProxy               bool
	HTTPProxy                 *string
	HTTPSProxy                *string
	NoProxy                   *string
	AdditionalTrustBundleFile *string
	AdditionalTrustBundle     *string

	// HyperShift options:
	Hypershift                  Hypershift
	BillingAccount              string
	NoCni                       bool
	AdditionalAllowedPrincipals []string

	// Audit Log Forwarding
	AuditLogRoleARN *string

	// AutoNode configuration
	AutoNodeMode    string
	AutoNodeRoleARN string

	Ec2MetadataHttpTokens cmv1.Ec2MetadataHttpTokens

	// Cluster Admin
	ClusterAdminUser     string
	ClusterAdminPassword string

	// Default Ingress Attributes
	DefaultIngress DefaultIngressSpec

	// Machine pool's storage
	MachinePoolRootDisk *Volume

	// Shared VPC
	PrivateHostedZoneID string
	SharedVPCRoleArn    string
	BaseDomain          string

	// HCP Shared VPC
	VpcEndpointRoleArn                string
	InternalCommunicationHostedZoneId string

	// Worker Machine Pool attributes
	AdditionalComputeSecurityGroupIds []string

	// Infra Machine Pool attributes
	AdditionalInfraSecurityGroupIds []string

	// Control Plane Machine Pool attributes
	AdditionalControlPlaneSecurityGroupIds []string

	// Registry Config
	AllowedRegistries          []string
	BlockedRegistries          []string
	InsecureRegistries         []string
	AllowedRegistriesForImport string
	PlatformAllowlist          string
	AdditionalTrustedCaFile    string
	AdditionalTrustedCa        map[string]string

	// Master/Infra Machine Config
	MasterMachineType string
	InfraMachineType  string

	// LogForward
	S3LogForwarder         *logforwarding.S3LogForwarderConfig
	CloudWatchLogForwarder *logforwarding.CloudWatchLogForwarderConfig
}

// Volume represents a volume property for a disk
type Volume struct {
	Size int
}

type OperatorIAMRole struct {
	Name      string
	Namespace string
	RoleARN   string
	Path      string
}

func NewOperatorIamRoleFromCmv1(operatorIAMRole *cmv1.OperatorIAMRole) (*OperatorIAMRole, error) {
	path, err := aws.GetPathFromARN(operatorIAMRole.RoleARN())
	if err != nil {
		return nil, err
	}
	return &OperatorIAMRole{
		Name:      operatorIAMRole.Name(),
		Namespace: operatorIAMRole.Namespace(),
		RoleARN:   operatorIAMRole.RoleARN(),
		Path:      path,
	}, nil
}

type Hypershift struct {
	Enabled bool
}

// Generate a query that filters clusters running on the current AWS session account
func getClusterFilter(creator *aws.Creator) string {
	filter := "product.id = 'rosa'"
	if creator != nil {
		filter = fmt.Sprintf("%s AND (properties.%s LIKE '%%:%s:%%' OR aws.sts.role_arn LIKE '%%:%s:%%')",
			filter,
			ocmConsts.CreatorArn,
			creator.AccountID,
			creator.AccountID)
	}
	return filter
}

func (c *Client) HasClusters(creator *aws.Creator) (bool, error) {
	query := getClusterFilter(creator)
	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		List().
		Search(query).
		Page(1).
		Size(1).
		Send()
	if err != nil {
		return false, handleErr(response.Error(), err)
	}

	return response.Total() > 0, nil
}

func (c *Client) CreateCluster(config Spec) (*cmv1.Cluster, error) {
	spec, err := c.createClusterSpec(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create cluster spec: %v", err)
	}

	cluster, err := c.ocm.ClustersMgmt().V1().Clusters().
		Add().
		Parameter("dryRun", *config.DryRun).
		Body(spec).
		Send()
	if config.DryRun != nil && *config.DryRun {
		if cluster.Error() != nil {
			return nil, handleErr(cluster.Error(), err)
		}
		return nil, nil
	}
	if err != nil {
		return nil, handleErr(cluster.Error(), err)
	}

	clusterObject := cluster.Body()

	return clusterObject, nil
}

// Maps between the account role type and the field on the cluster struct that the role
// will be present in so that we can send the correct query in the "search" parameter
// in the GET to the /api/clusters_mgmt/v1/clusters endpoint
var accountRoleTypeFieldMap = map[string]string{
	aws.InstallerAccountRoleType:    "aws.sts.role_arn",
	aws.ControlPlaneAccountRoleType: "aws.sts.instance_iam_roles.master_role_arn",
	aws.SupportAccountRoleType:      "aws.sts.support_role_arn",
	aws.WorkerAccountRoleType:       "aws.sts.instance_iam_roles.worker_role_arn",
}

func getAccountRoleClusterFilter(aws *aws.Creator, role aws.Role) (string, error) {
	query := getClusterFilter(aws)
	accountRoleField := accountRoleTypeFieldMap[role.RoleType]
	if accountRoleField == "" {
		return "",
			fmt.Errorf("unrecognised Role Type '%s' for Account Role with ARN '%s'", role.RoleType, role.RoleARN)
	}

	// Append to our normal query our search for the specific role arn based around the role type
	return fmt.Sprintf("%s AND %s='%s'", query, accountRoleField, role.RoleARN), nil
}

func (c *Client) GetClustersUsingAccountRole(aws *aws.Creator, role aws.Role, count int) ([]*cmv1.Cluster, error) {
	query, err := getAccountRoleClusterFilter(aws, role)
	if err != nil {
		return nil, err
	}

	return c.queryClusters(query, count)
}

func (c *Client) queryClusters(query string, count int) (clusters []*cmv1.Cluster, err error) {

	if count < 0 {
		err = errors.Errorf("Invalid Cluster count")
		return
	}

	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)
	page := 1
	for {
		clusterRequestList := request.Page(page)
		if count > 0 {
			clusterRequestList = clusterRequestList.Size(count)
		}
		response, err := clusterRequestList.Send()
		if err != nil {
			return clusters, err
		}

		response.Items().Each(func(cluster *cmv1.Cluster) bool {
			clusters = append(clusters, cluster)
			return true
		})
		if response.Size() != count {
			break
		}
		page++
	}
	return clusters, nil
}

// Pass 0 to get all clusters
func (c *Client) GetClusters(creator *aws.Creator, count int) (clusters []*cmv1.Cluster, err error) {
	return c.queryClusters(getClusterFilter(creator), count)
}

func (c *Client) GetAllClusters(creator *aws.Creator) (clusters []*cmv1.Cluster, err error) {
	query := getClusterFilter(creator)
	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)
	response, err := request.Send()

	if err != nil {
		return clusters, err
	}
	return response.Items().Slice(), nil
}

// GetCluster gets a cluster key that can be either 'id', 'name' or 'external_id'
func (c *Client) GetCluster(clusterKey string, creator *aws.Creator) (*cmv1.Cluster, error) {
	query := fmt.Sprintf("%s AND (id = '%s' OR name = '%s' OR external_id = '%s')",
		getClusterFilter(creator),
		clusterKey, clusterKey, clusterKey,
	)
	response, err := c.ocm.ClustersMgmt().V1().Clusters().List().
		Search(query).
		Page(1).
		Size(1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	switch response.Total() {
	case 0:
		return nil, errors.NotFound.Errorf("There is no cluster with identifier or name '%s'", clusterKey)
	case 1:
		return response.Items().Slice()[0], nil
	default:
		return nil, fmt.Errorf("there are %d clusters with identifier or name '%s'", response.Total(), clusterKey)
	}
}

func (c *Client) GetSubscriptionBySubscriptionID(id string) (*amv1.Subscription, bool, error) {
	response, err := c.ocm.AccountsMgmt().V1().Subscriptions().Subscription(id).
		Get().
		Send()

	if err != nil {
		return nil, false, err
	}
	if response.Body() == nil {
		return &amv1.Subscription{}, false, nil
	}

	return response.Body(), true, nil
}

func (c *Client) GetClusterByID(clusterKey string, creator *aws.Creator) (*cmv1.Cluster, error) {
	query := fmt.Sprintf("%s AND id = '%s'",
		getClusterFilter(creator),
		clusterKey,
	)
	response, err := c.ocm.ClustersMgmt().V1().Clusters().List().
		Search(query).
		Page(1).
		Size(1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	switch response.Total() {
	case 0:
		return nil, errors.NotFound.Errorf("There is no cluster with identifier '%s'", clusterKey)
	case 1:
		return response.Items().Slice()[0], nil
	default:
		return nil, fmt.Errorf("there are %d clusters with identifier '%s'", response.Total(), clusterKey)
	}
}

func (c *Client) GetClusterUsingSubscription(clusterKey string, creator *aws.Creator) (*amv1.Subscription, error) {
	query := fmt.Sprintf("(plan.id = 'MOA' OR plan.id = 'MOA-HostedControlPlane')"+
		" AND (display_name  = '%s' OR cluster_id = '%s') AND status = 'Deprovisioned'", clusterKey, clusterKey)
	response, err := c.ocm.AccountsMgmt().V1().Subscriptions().List().
		Search(query).
		Page(1).
		Size(1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	switch response.Total() {
	case 0:
		return nil, nil
	case 1:
		return response.Items().Slice()[0], nil
	default:
		return nil, errors.Conflict.Errorf("There are %d clusters with identifier '%s'", response.Total(),
			clusterKey)
	}
}

// Gets only pending non-STS clusters that are installed in the same AWS account
func (c *Client) GetPendingClusterForARN(creator *aws.Creator) (cluster *cmv1.Cluster, err error) {
	query := fmt.Sprintf(
		"state = 'pending' AND product.id = 'rosa' AND aws.sts.role_arn = '' AND properties.%s LIKE '%%:%s:%%'",
		ocmConsts.CreatorArn,
		creator.AccountID,
	)
	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)

	response, err := request.Send()
	if err != nil {
		return cluster, err
	}
	return response.Items().Get(0), nil
}

func (c *Client) HasAClusterUsingOperatorRolesPrefix(prefix string) (bool, error) {
	query := fmt.Sprintf(
		"aws.sts.operator_iam_roles.role_arn like '%%/%s-%%'", prefix,
	)
	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)
	page := 1
	response, err := request.Page(page).Send()
	if err != nil {
		return false, err
	}
	if response.Total() > 0 {
		return true, nil
	}
	return false, nil
}

func (c *Client) HasAClusterUsingOidcProvider(
	issuerUrl string, curAccountId string) (bool, error) {
	query := fmt.Sprintf(
		"aws.sts.oidc_endpoint_url = '%s' AND aws.sts.role_arn like '%%%s%%'",
		issuerUrl, curAccountId,
	)
	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)
	page := 1
	response, err := request.Page(page).Send()
	if err != nil {
		return false, err
	}
	if response.Total() > 0 {
		return true, nil
	}
	return false, nil
}

func (c *Client) HasAClusterUsingOidcEndpointUrl(issuerUrl string) (bool, error) {
	query := fmt.Sprintf(
		"aws.sts.oidc_endpoint_url = '%s'", issuerUrl,
	)
	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)
	page := 1
	response, err := request.Page(page).Send()
	if err != nil {
		return false, err
	}
	if response.Total() > 0 {
		return true, nil
	}
	return false, nil
}

func (c *Client) IsSTSClusterExists(creator *aws.Creator, count int, roleARN string) (exists bool, err error) {
	if count < 1 {
		err = errors.Errorf("Cannot fetch fewer than 1 cluster")
		return
	}
	query := fmt.Sprintf(
		"product.id = 'rosa' AND ("+
			"properties.%s LIKE '%%:%s:%%' OR "+
			"aws.sts.role_arn = '%s' OR "+
			"aws.sts.support_role_arn = '%s' OR "+
			"aws.sts.instance_iam_roles.master_role_arn = '%s' OR "+
			"aws.sts.instance_iam_roles.worker_role_arn = '%s')",
		ocmConsts.CreatorArn,
		creator.AccountID,
		roleARN,
		roleARN,
		roleARN,
		roleARN,
	)
	request := c.ocm.ClustersMgmt().V1().Clusters().List().Search(query)
	page := 1
	response, err := request.Page(page).Size(count).Send()
	if err != nil {
		return false, err
	}
	if response.Total() > 0 {
		return true, nil
	}
	return false, nil
}

func (c *Client) GetClusterState(clusterID string) (cmv1.ClusterState, error) {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		Status().
		Get().
		Send()
	if err != nil || response.Body() == nil {
		return cmv1.ClusterState(""), err
	}
	return response.Body().State(), nil
}

func (c *Client) getClusterNodesBuilder(config Spec) (clusterNodesBuilder *cmv1.ClusterNodesBuilder, updateNodes bool) {

	clusterNodesBuilder = cmv1.NewClusterNodes()
	if config.Autoscaling {
		updateNodes = true
		autoscalingBuilder := cmv1.NewMachinePoolAutoscaling()
		if config.MinReplicas != 0 {
			autoscalingBuilder = autoscalingBuilder.MinReplicas(config.MinReplicas)
		}
		if config.MaxReplicas != 0 {
			autoscalingBuilder = autoscalingBuilder.MaxReplicas(config.MaxReplicas)
		}
		clusterNodesBuilder = clusterNodesBuilder.AutoscaleCompute(autoscalingBuilder)
	} else if config.ComputeNodes != 0 {
		updateNodes = true
		clusterNodesBuilder = clusterNodesBuilder.Compute(config.ComputeNodes)
	}

	if config.ComputeLabels != nil {
		updateNodes = true
		clusterNodesBuilder = clusterNodesBuilder.ComputeLabels(config.ComputeLabels)
	}

	return

}

func (c *Client) UpdateCluster(clusterKey string, creator *aws.Creator, config Spec) error {
	cluster, err := c.GetCluster(clusterKey, creator)
	if err != nil {
		return err
	}

	clusterBuilder := cmv1.NewCluster()

	// Update expiration timestamp
	if !config.Expiration.IsZero() {
		clusterBuilder = clusterBuilder.ExpirationTimestamp(config.Expiration)
	}

	// Update channel group
	if config.ChannelGroup != "" {
		clusterBuilder.Version(cmv1.NewVersion().
			ChannelGroup(config.ChannelGroup),
		)
	}

	// Scale cluster
	clusterNodesBuilder, updateNodes := c.getClusterNodesBuilder(config)
	if updateNodes {
		clusterBuilder = clusterBuilder.Nodes(clusterNodesBuilder)
	}

	// Toggle private mode
	if config.Private != nil {
		if *config.Private {
			clusterBuilder = clusterBuilder.API(
				cmv1.NewClusterAPI().
					Listening(cmv1.ListeningMethodInternal),
			)
		} else {
			clusterBuilder = clusterBuilder.API(
				cmv1.NewClusterAPI().
					Listening(cmv1.ListeningMethodExternal),
			)
		}
	}

	if config.NodeDrainGracePeriodInMinutes != 0 {
		clusterBuilder = clusterBuilder.NodeDrainGracePeriod(
			cmv1.NewValue().
				Value(config.NodeDrainGracePeriodInMinutes).
				Unit("minutes"),
		)
	}

	if config.DisableWorkloadMonitoring != nil {
		clusterBuilder = clusterBuilder.DisableUserWorkloadMonitoring(*config.DisableWorkloadMonitoring)
	}

	// SDN -> OVN Migration
	if config.NetworkType == NetworkTypes[1] {
		// Create a request body for the specific cluster migration.
		requestBuilder := cmv1.ClusterMigrationBuilder{}
		requestBuilder.Type(cmv1.ClusterMigrationTypeSdnToOvn) // Type is required

		if len(config.OvnInternalSubnetConfiguration) > 0 {
			// Create a builder for the specific migration type's configuration if necessary
			sdnToOvnBuilder := &cmv1.SdnToOvnClusterMigrationBuilder{}
			if _, ok := config.OvnInternalSubnetConfiguration[SubnetConfigJoin]; ok {
				sdnToOvnBuilder.JoinIpv4(config.OvnInternalSubnetConfiguration[SubnetConfigJoin])
			}
			if _, ok := config.OvnInternalSubnetConfiguration[SubnetConfigTransit]; ok {
				sdnToOvnBuilder.TransitIpv4(config.OvnInternalSubnetConfiguration[SubnetConfigTransit])
			}
			if _, ok := config.OvnInternalSubnetConfiguration[SubnetConfigMasquerade]; ok {
				sdnToOvnBuilder.MasqueradeIpv4(config.OvnInternalSubnetConfiguration[SubnetConfigMasquerade])
			}
			requestBuilder.SdnToOvn(sdnToOvnBuilder)
		}

		requestBody, err := requestBuilder.Build()
		if err != nil {
			return errors.UserWrapf(err, "Unable to create cluster migration request")
		}

		// Send the request to add a cluster migration.
		response, err := c.ocm.ClustersMgmt().V1().Clusters().
			Cluster(cluster.ID()).Migrations().Add().Body(requestBody).Send()
		if err != nil {
			return handleErr(response.Error(), err)
		}
	}

	if config.HTTPProxy != nil || config.HTTPSProxy != nil || config.NoProxy != nil {
		clusterProxyBuilder := cmv1.NewProxy()
		if config.HTTPProxy != nil {
			clusterProxyBuilder = clusterProxyBuilder.HTTPProxy(*config.HTTPProxy)
		}
		if config.HTTPSProxy != nil {
			clusterProxyBuilder = clusterProxyBuilder.HTTPSProxy(*config.HTTPSProxy)
		}
		if config.NoProxy != nil {
			clusterProxyBuilder = clusterProxyBuilder.NoProxy(*config.NoProxy)
		}
		clusterBuilder = clusterBuilder.Proxy(clusterProxyBuilder)
	}

	if config.AdditionalTrustBundle != nil {
		clusterBuilder = clusterBuilder.AdditionalTrustBundle(*config.AdditionalTrustBundle)
	}

	if config.Hypershift.Enabled {
		hyperShiftBuilder := cmv1.NewHypershift().Enabled(true)
		clusterBuilder.Hypershift(hyperShiftBuilder)
	}

	registryConfigBuilder, err := BuildRegistryConfig(config)
	if err != nil {
		return err
	}
	if registryConfigBuilder != nil {
		clusterBuilder.RegistryConfig(registryConfigBuilder)
	}

	if config.AuditLogRoleARN != nil || config.AdditionalAllowedPrincipals != nil || config.BillingAccount != "" ||
		config.AutoNodeRoleARN != "" {
		awsBuilder := cmv1.NewAWS()
		if config.AdditionalAllowedPrincipals != nil {
			awsBuilder = awsBuilder.AdditionalAllowedPrincipals(config.AdditionalAllowedPrincipals...)
		}
		// Edit audit log role arn
		if config.AuditLogRoleARN != nil {
			auditLogBuiler := cmv1.NewAuditLog().RoleArn(*config.AuditLogRoleARN)
			awsBuilder = awsBuilder.AuditLog(auditLogBuiler)
		}
		if config.BillingAccount != "" {
			awsBuilder.BillingAccountID(config.BillingAccount)
		}
		// Add AutoNode configuration
		if config.AutoNodeRoleARN != "" {
			autoNodeBuilder := cmv1.NewAwsAutoNode().RoleArn(config.AutoNodeRoleARN)
			awsBuilder = awsBuilder.AutoNode(autoNodeBuilder)
		}
		clusterBuilder.AWS(awsBuilder)
	}

	// Set AutoNode mode if specified
	if config.AutoNodeMode != "" {
		autoNodeBuilder := cmv1.NewClusterAutoNode().Mode(config.AutoNodeMode)
		clusterBuilder.AutoNode(autoNodeBuilder)
	}

	clusterSpec, err := clusterBuilder.Build()
	if err != nil {
		return err
	}

	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(cluster.ID()).
		Update().
		Body(clusterSpec).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}

func (c *Client) DeleteCluster(clusterKey string, bestEffort bool,
	creator *aws.Creator) (*cmv1.Cluster, error) {
	cluster, err := c.GetCluster(clusterKey, creator)
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(cluster.ID()).
		Delete().
		BestEffort(bestEffort).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return cluster, nil
}

func (c *Client) UpdateClusterDeletionProtection(clusterId string, deleteProtection *cmv1.DeleteProtection) error {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterId).
		DeleteProtection().
		Update().
		Body(deleteProtection).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

// EnsureNoPendingClusters ensures that no clusters are pending in the account. For non-STS clusters,
// the osdCcsAdmin user credentials are used to create the cluster, and it is required that these credentials
// are rotated between cluster creation. If a user is creating a non-STS cluster, we need to therefore make sure
// no other clusters are pending in the account in order to ensure no race condition occurs.
func (c *Client) EnsureNoPendingClusters(awsCreator *aws.Creator) error {
	reporter := rprtr.CreateReporter()
	/**
	1) Poll the cluster with same arn from ocm
	2) Check the status and if pending enter to a loop until it becomes installing
	3) Do it only for ROSA clusters and before UpsertAccessKey
	*/
	deadline := time.Now().Add(5 * time.Minute)
	for {
		pendingCluster, err := c.GetPendingClusterForARN(awsCreator)
		if err != nil {
			reporter.Errorf("Error getting cluster using ARN '%s'", awsCreator.ARN)
			os.Exit(1)
		}
		if time.Now().After(deadline) {
			reporter.Errorf("Timeout waiting for the cluster '%s' installation. Try again in a few minutes",
				pendingCluster.ID())
			os.Exit(1)
		}
		if pendingCluster == nil {
			break
		} else {
			reporter.Infof("Waiting for cluster '%s' with the same creator ARN to start installing",
				pendingCluster.ID())
			time.Sleep(30 * time.Second)
		}
	}
	return nil
}

func (c *Client) createClusterSpec(config Spec) (*cmv1.Cluster, error) {
	reporter := rprtr.CreateReporter()
	clusterProperties := map[string]string{}

	if config.CustomProperties != nil {
		for key, value := range config.CustomProperties {
			clusterProperties[key] = value
		}
	}

	// Make sure we don't have a custom properties collision
	if _, present := clusterProperties[ocmConsts.CreatorArn]; present {
		return nil, fmt.Errorf(
			"custom properties key %s collides with a property needed by rosa",
			ocmConsts.CreatorArn,
		)
	}

	if _, present := clusterProperties[properties.CLIVersion]; present {
		return nil, fmt.Errorf(
			"custom properties key %s collides with a property needed by rosa",
			properties.CLIVersion,
		)
	}

	if config.AWSCreator == nil {
		return nil, fmt.Errorf("AWS creator metadata is required")
	}

	clusterProperties[ocmConsts.CreatorArn] = config.AWSCreator.ARN
	clusterProperties[properties.CLIVersion] = info.DefaultVersion

	// Create the cluster:
	clusterBuilder := cmv1.NewCluster().
		Name(config.Name).
		MultiAZ(config.MultiAZ).
		Product(
			cmv1.NewProduct().
				ID("rosa"),
		).
		Region(
			cmv1.NewCloudRegion().
				ID(config.Region),
		).
		FIPS(config.FIPS).
		EtcdEncryption(config.EtcdEncryption).
		Properties(clusterProperties)

	if config.DomainPrefix != "" {
		clusterBuilder.DomainPrefix(config.DomainPrefix)
	}

	if config.DisableWorkloadMonitoring != nil {
		clusterBuilder = clusterBuilder.DisableUserWorkloadMonitoring(*config.DisableWorkloadMonitoring)
	}

	registryConfigBuilder, err := BuildRegistryConfig(config)
	if err != nil {
		return nil, err
	}
	if registryConfigBuilder != nil {
		clusterBuilder.RegistryConfig(registryConfigBuilder)
	}

	if config.Flavour != "" {
		clusterBuilder = clusterBuilder.Flavour(
			cmv1.NewFlavour().
				ID(config.Flavour),
		)
		reporter.Debugf("Using cluster flavour '%s'", config.Flavour)
	}

	if config.Version != "" {
		clusterBuilder = clusterBuilder.Version(
			cmv1.NewVersion().
				ID(config.Version).
				ChannelGroup(config.ChannelGroup),
		)

		reporter.Debugf(
			"Using OpenShift version '%s' on channel group '%s'",
			config.Version, config.ChannelGroup)
	}

	if !config.Expiration.IsZero() {
		clusterBuilder = clusterBuilder.ExpirationTimestamp(config.Expiration)
	}

	if config.Hypershift.Enabled {
		hyperShiftBuilder := cmv1.NewHypershift().Enabled(true)
		clusterBuilder.Hypershift(hyperShiftBuilder)
	}

	if config.ExternalAuthProvidersEnabled {
		externalAuthConfigBuilder := cmv1.NewExternalAuthConfig().Enabled(true)
		clusterBuilder.ExternalAuthConfig(externalAuthConfigBuilder)
	}

	if config.ComputeMachineType != "" || config.ComputeNodes != 0 || len(config.AvailabilityZones) > 0 ||
		config.Autoscaling || len(config.ComputeLabels) > 0 {
		clusterNodesBuilder := cmv1.NewClusterNodes()
		if config.ComputeMachineType != "" {
			clusterNodesBuilder = clusterNodesBuilder.ComputeMachineType(
				cmv1.NewMachineType().ID(config.ComputeMachineType),
			)

			reporter.Debugf("Using machine type '%s'", config.ComputeMachineType)
		}
		if machinePoolRootDisk := config.MachinePoolRootDisk; machinePoolRootDisk != nil &&
			machinePoolRootDisk.Size != 0 {
			machineTypeRootVolumeBuilder := cmv1.NewRootVolume().
				AWS(cmv1.NewAWSVolume().
					Size(machinePoolRootDisk.Size))
			clusterNodesBuilder = clusterNodesBuilder.ComputeRootVolume(
				(machineTypeRootVolumeBuilder),
			)
		}
		if config.Autoscaling {
			clusterNodesBuilder = clusterNodesBuilder.AutoscaleCompute(
				cmv1.NewMachinePoolAutoscaling().
					MinReplicas(config.MinReplicas).
					MaxReplicas(config.MaxReplicas))
		} else if config.ComputeNodes != 0 {
			clusterNodesBuilder = clusterNodesBuilder.Compute(config.ComputeNodes)
		}
		if len(config.AvailabilityZones) > 0 {
			clusterNodesBuilder = clusterNodesBuilder.AvailabilityZones(config.AvailabilityZones...)
		}
		if len(config.ComputeLabels) > 0 {
			clusterNodesBuilder = clusterNodesBuilder.ComputeLabels(config.ComputeLabels)
		}
		if config.MasterMachineType != "" {
			clusterNodesBuilder.MasterMachineType(cmv1.NewMachineType().ID(config.MasterMachineType))
		}
		if config.InfraMachineType != "" {
			clusterNodesBuilder.InfraMachineType(cmv1.NewMachineType().ID(config.InfraMachineType))
		}
		clusterBuilder = clusterBuilder.Nodes(clusterNodesBuilder)
	}

	if config.NetworkType != "" || config.NoCni ||
		!IsEmptyCIDR(config.MachineCIDR) ||
		!IsEmptyCIDR(config.ServiceCIDR) ||
		!IsEmptyCIDR(config.PodCIDR) ||
		config.HostPrefix != 0 {
		networkBuilder := cmv1.NewNetwork()
		if config.NetworkType != "" {
			networkBuilder = networkBuilder.Type(config.NetworkType)
		}
		if config.NoCni {
			networkBuilder = networkBuilder.Type("Other")
		}
		if !IsEmptyCIDR(config.MachineCIDR) {
			networkBuilder = networkBuilder.MachineCIDR(config.MachineCIDR.String())
		}
		if !IsEmptyCIDR(config.ServiceCIDR) {
			networkBuilder = networkBuilder.ServiceCIDR(config.ServiceCIDR.String())
		}
		if !IsEmptyCIDR(config.PodCIDR) {
			networkBuilder = networkBuilder.PodCIDR(config.PodCIDR.String())
		}
		if config.HostPrefix != 0 {
			networkBuilder = networkBuilder.HostPrefix(config.HostPrefix)
		}
		clusterBuilder = clusterBuilder.Network(networkBuilder)
	}

	awsBuilder := cmv1.NewAWS().
		AccountID(config.AWSCreator.AccountID)

	if len(config.AdditionalComputeSecurityGroupIds) > 0 {
		awsBuilder = awsBuilder.AdditionalComputeSecurityGroupIds(config.AdditionalComputeSecurityGroupIds...)
	}

	if len(config.AdditionalInfraSecurityGroupIds) > 0 {
		awsBuilder = awsBuilder.AdditionalInfraSecurityGroupIds(config.AdditionalInfraSecurityGroupIds...)
	}

	if len(config.AdditionalControlPlaneSecurityGroupIds) > 0 {
		awsBuilder = awsBuilder.AdditionalControlPlaneSecurityGroupIds(config.AdditionalControlPlaneSecurityGroupIds...)
	}

	if len(config.AdditionalAllowedPrincipals) > 0 {
		awsBuilder = awsBuilder.AdditionalAllowedPrincipals(config.AdditionalAllowedPrincipals...)
	}

	if config.SubnetIds != nil {
		awsBuilder = awsBuilder.SubnetIDs(config.SubnetIds...)
	}

	if config.PrivateLink != nil {
		awsBuilder = awsBuilder.PrivateLink(*config.PrivateLink)
		if *config.PrivateLink {
			*config.Private = true
		}
	}

	if config.BillingAccount != "" {
		awsBuilder = awsBuilder.BillingAccountID(config.BillingAccount)
	}

	if config.Ec2MetadataHttpTokens != "" {
		awsBuilder = awsBuilder.Ec2MetadataHttpTokens(config.Ec2MetadataHttpTokens)
	}

	if config.RoleARN != "" {
		stsBuilder := cmv1.NewSTS().RoleARN(config.RoleARN)
		if config.ExternalID != "" {
			stsBuilder = stsBuilder.ExternalID(config.ExternalID)
		}
		if config.SupportRoleARN != "" {
			stsBuilder = stsBuilder.SupportRoleARN(config.SupportRoleARN)
		}
		if len(config.OperatorIAMRoles) > 0 {
			roles := []*cmv1.OperatorIAMRoleBuilder{}
			for _, role := range config.OperatorIAMRoles {
				roles = append(roles, cmv1.NewOperatorIAMRole().
					Name(role.Name).
					Namespace(role.Namespace).
					RoleARN(role.RoleARN),
				)
			}
			stsBuilder = stsBuilder.OperatorIAMRoles(roles...)
		}
		if config.OidcConfigId != "" {
			stsBuilder = stsBuilder.OidcConfig(cmv1.NewOidcConfig().ID(config.OidcConfigId))
		}
		instanceIAMRolesBuilder := cmv1.NewInstanceIAMRoles()
		if config.ControlPlaneRoleARN != "" {
			instanceIAMRolesBuilder.MasterRoleARN(config.ControlPlaneRoleARN)
		}
		if config.WorkerRoleARN != "" {
			instanceIAMRolesBuilder.WorkerRoleARN(config.WorkerRoleARN)
		}
		stsBuilder = stsBuilder.InstanceIAMRoles(instanceIAMRolesBuilder)

		mode := false
		if config.Mode == "auto" {
			mode = true
		}
		stsBuilder.AutoMode(mode)

		awsBuilder = awsBuilder.STS(stsBuilder)
	} else {
		if config.AWSAccessKey == nil {
			return nil, fmt.Errorf("AWS access key metadata is required for non-STS clusters")
		}
		awsBuilder = awsBuilder.
			AccessKeyID(config.AWSAccessKey.AccessKeyID).
			SecretAccessKey(config.AWSAccessKey.SecretAccessKey)
	}
	if config.KMSKeyArn != "" {
		awsBuilder = awsBuilder.KMSKeyArn(config.KMSKeyArn)
	}
	if len(config.Tags) > 0 {
		awsBuilder = awsBuilder.Tags(config.Tags)
	}

	if config.AuditLogRoleARN != nil {
		auditLogBuiler := cmv1.NewAuditLog().RoleArn(*config.AuditLogRoleARN)
		awsBuilder = awsBuilder.AuditLog(auditLogBuiler)
	}

	// etcd encryption kms key arn
	if config.EtcdEncryptionKMSArn != "" {
		awsBuilder = awsBuilder.EtcdEncryption(cmv1.NewAwsEtcdEncryption().KMSKeyARN(config.EtcdEncryptionKMSArn))
	}

	// shared vpc
	if config.PrivateHostedZoneID != "" {
		awsBuilder = awsBuilder.PrivateHostedZoneID(config.PrivateHostedZoneID)
		awsBuilder = awsBuilder.PrivateHostedZoneRoleARN(config.SharedVPCRoleArn)
	}
	// hcp shared vpc
	if config.VpcEndpointRoleArn != "" {
		awsBuilder = awsBuilder.PrivateHostedZoneID(config.PrivateHostedZoneID)
		awsBuilder = awsBuilder.PrivateHostedZoneRoleARN(config.SharedVPCRoleArn)
		awsBuilder = awsBuilder.VpcEndpointRoleArn(config.VpcEndpointRoleArn)
		awsBuilder = awsBuilder.HcpInternalCommunicationHostedZoneId(config.InternalCommunicationHostedZoneId)
	}
	if config.BaseDomain != "" {
		clusterBuilder = clusterBuilder.DNS(cmv1.NewDNS().BaseDomain(config.BaseDomain))
	}

	clusterBuilder = clusterBuilder.AWS(awsBuilder)

	clusterApiListeningMethod := cmv1.ListeningMethodExternal

	if config.Private != nil {
		if *config.Private {
			clusterBuilder = clusterBuilder.API(
				cmv1.NewClusterAPI().
					Listening(cmv1.ListeningMethodInternal),
			)
			clusterApiListeningMethod = cmv1.ListeningMethodInternal
		} else {
			clusterBuilder = clusterBuilder.API(
				cmv1.NewClusterAPI().
					Listening(cmv1.ListeningMethodExternal),
			)
		}
	}

	if config.DisableSCPChecks != nil && *config.DisableSCPChecks {
		clusterBuilder = clusterBuilder.CCS(cmv1.NewCCS().
			Enabled(true).
			DisableSCPChecks(true),
		)
	}

	if config.HTTPProxy != nil || config.HTTPSProxy != nil {
		proxyBuilder := cmv1.NewProxy()
		if config.HTTPProxy != nil {
			proxyBuilder.HTTPProxy(*config.HTTPProxy)
		}
		if config.HTTPSProxy != nil {
			proxyBuilder.HTTPSProxy(*config.HTTPSProxy)
		}
		if config.NoProxy != nil {
			proxyBuilder.NoProxy(*config.NoProxy)
		}
		clusterBuilder = clusterBuilder.Proxy(proxyBuilder)
	}

	if config.AdditionalTrustBundle != nil {
		clusterBuilder = clusterBuilder.AdditionalTrustBundle(*config.AdditionalTrustBundle)
	}

	if config.ClusterAdminUser != "" {
		hashedPwd, err := idputils.GenerateHTPasswdCompatibleHash(config.ClusterAdminPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to get access keys for user '%s': %v",
				aws.AdminUserName, err)
		}
		htpasswdUsers := []*cmv1.HTPasswdUserBuilder{}
		htpasswdUsers = append(htpasswdUsers, cmv1.NewHTPasswdUser().
			Username(config.ClusterAdminUser).HashedPassword(hashedPwd))
		htpassUserList := cmv1.NewHTPasswdUserList().Items(htpasswdUsers...)
		htPasswdIDP := cmv1.NewHTPasswdIdentityProvider().Users(htpassUserList)
		clusterBuilder = clusterBuilder.Htpasswd(htPasswdIDP)
	}

	// Build default ingress if changes detected in config
	defaultIngress := cmv1.NewIngress().Default(true)
	if len(config.DefaultIngress.RouteSelectors) != 0 {
		defaultIngress.RouteSelectors(config.DefaultIngress.RouteSelectors)
	}
	if len(config.DefaultIngress.ExcludedNamespaces) != 0 {
		defaultIngress.ExcludedNamespaces(config.DefaultIngress.ExcludedNamespaces...)
	}
	if !helper.Contains([]string{"", consts.SkipSelectionOption}, config.DefaultIngress.WildcardPolicy) {
		defaultIngress.RouteWildcardPolicy(cmv1.WildcardPolicy(config.DefaultIngress.WildcardPolicy))
	}
	if !helper.Contains([]string{"", consts.SkipSelectionOption}, config.DefaultIngress.NamespaceOwnershipPolicy) {
		defaultIngress.RouteNamespaceOwnershipPolicy(
			cmv1.NamespaceOwnershipPolicy(config.DefaultIngress.NamespaceOwnershipPolicy))
	}

	// Decide ingress listening method if HCP and not fedramp enabled
	isHcpNotFedramp := !fedramp.Enabled() && config.Hypershift.Enabled
	if isHcpNotFedramp {
		if config.PrivateIngress != nil {
			if *config.PrivateIngress {
				defaultIngress.Listening(cmv1.ListeningMethodInternal)
			} else {
				defaultIngress.Listening(cmv1.ListeningMethodExternal)
			}
		} else {
			defaultIngress.Listening(clusterApiListeningMethod)
		}
	}

	if !reflect.DeepEqual(config.DefaultIngress, NewDefaultIngressSpec()) || isHcpNotFedramp {
		clusterBuilder.Ingresses(cmv1.NewIngressList().Items(defaultIngress))
	}

	if config.AutoscalerConfig != nil {
		clusterBuilder.Autoscaler(BuildClusterAutoscaler(config.AutoscalerConfig))
	}

	if config.Hypershift.Enabled {
		// LogForwarder
		var logForwarderList []*cmv1.LogForwarderBuilder
		if config.S3LogForwarder != nil {
			boundForwarder, err := logforwarding.BindS3LogForwarder(config.S3LogForwarder).Build()
			if err != nil {
				return nil, fmt.Errorf("failed to create S3 log forwarder: %v", err)
			}
			s3LogFwBuilder := BuildLogForwarder(boundForwarder)
			logForwarderList = append(logForwarderList, s3LogFwBuilder)
		}
		if config.CloudWatchLogForwarder != nil {
			boundForwarder, err := logforwarding.BindCloudWatchLogForwarder(config.CloudWatchLogForwarder).Build()
			if err != nil {
				return nil, fmt.Errorf("failed to create CloudWatch log forwarder: %v", err)
			}
			cloudWatchFwBuilder := BuildLogForwarder(boundForwarder)
			logForwarderList = append(logForwarderList, cloudWatchFwBuilder)
		}
		if len(logForwarderList) > 0 {
			clusterBuilder.ControlPlane(cmv1.NewControlPlane().LogForwarders(
				cmv1.NewLogForwarderList().Items(logForwarderList...)))
		}
	}

	clusterSpec, err := clusterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create description of cluster: %v", err)
	}

	return clusterSpec, nil
}

func (c *Client) HibernateCluster(clusterID string) error {
	enabled, err := c.IsCapabilityEnabled(HibernateCapability)
	if err != nil {
		return err
	}
	if !enabled {
		return fmt.Errorf("the '%s' capability is not set for current org", HibernateCapability)
	}
	_, err = c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).Hibernate().Send()
	if err != nil {
		return fmt.Errorf("failed to hibernate the cluster: %v", err)
	}

	return nil
}

func (c *Client) ResumeCluster(clusterID string) error {
	enabled, err := c.IsCapabilityEnabled(HibernateCapability)
	if err != nil {
		return err
	}
	if !enabled {
		return fmt.Errorf("the '%s' capability is not set for current org", HibernateCapability)
	}
	_, err = c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).Resume().Send()
	if err != nil {
		return fmt.Errorf("failed to resume the cluster: %v", err)
	}

	return nil
}

func IsConsoleAvailable(cluster *cmv1.Cluster) bool {
	return cluster.Console() != nil && cluster.Console().URL() != ""
}

func IsHyperShiftCluster(cluster *cmv1.Cluster) bool {
	return cluster != nil && cluster.Hypershift() != nil && cluster.Hypershift().Enabled()
}

func IsOidcConfigReusable(cluster *cmv1.Cluster) bool {
	return cluster != nil &&
		cluster.AWS().STS().OidcConfig() != nil && cluster.AWS().STS().OidcConfig().Reusable()
}

func IsSts(cluster *cmv1.Cluster) bool {
	return cluster != nil && cluster.AWS().STS().RoleARN() != ""
}

func (c *Client) HasLegacyIngressSupport(cluster *cmv1.Cluster) (bool, error) {
	labelList, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(cluster.ID()).ExternalConfiguration().Labels().List().Send()
	if err != nil {
		return true, fmt.Errorf("failed to retrieve external configuration label list: %v", err)
	}

	for _, label := range labelList.Items().Slice() {
		if label.Key() == legacyIngressSupportLabel {
			return strconv.ParseBool(label.Value())
		}
	}
	return true, nil
}
