/*
Copyright 2021 The Kubernetes Authors.

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

package powervs

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	tgapiv1 "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/accounts"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/cos"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcemanager"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/transitgateway"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
	ibmcloudutil "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/util"
)

// vpcSubnetIPVersion4 defines the IP v4 string used for VPC Subnet generation.
var vpcSubnetIPVersion4 = "ipv4"

// ResourceNotFound is the string representing an error when a resource is not found in IBM Cloud.
type ResourceNotFound string

// networkConnectionType represents network connection type in transit gateway.
type networkConnectionType string

var (
	powervsNetworkConnectionType = networkConnectionType("power_virtual_server")
	vpcNetworkConnectionType     = networkConnectionType("vpc")

	// ResourceNotFoundCode indicates the http status code when a resource does not exist.
	ResourceNotFoundCode = 404

	// DHCPServerNotFound is the error returned when a DHCP server is not found.
	DHCPServerNotFound = ResourceNotFound("dhcp server does not exist")
)

// sgRuleKey is a comparable identity for a security group rule.
type sgRuleKey struct {
	direction string
	protocol  string
	portMin   int64
	portMax   int64
	cidrBlock string
	address   string
	crn       string
}

// powerEdgeRouter is identifier for PER.
const (
	// DEBUGLEVEL indicates the debug level of the logs.
	DEBUGLEVEL = 5

	powerEdgeRouter = "power-edge-router"
	// vpcSubnetIPAddressCount is the total IP Addresses for the subnet.
	// Support for custom address prefixes will be added at a later time. Currently, we use the ip count for subnet creation.
	vpcSubnetIPAddressCount int64 = 256

	// cosHMACAccessKeyField and cosHMACSecretKeyField are the keys used in the
	// Kubernetes Secret that stores COS HMAC credentials for pre-signed URL generation.
	cosHMACAccessKeyField = "access_key_id"
	cosHMACSecretKeyField = "secret_access_key"
)

// ClusterScopeParams defines the input parameters used to create a new ClusterScope.
type ClusterScopeParams struct {
	Client            client.Client
	Cluster           *clusterv1.Cluster
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
	ClientBuilder     ClientBuilder
	Recorder          cgrecord.EventRecorder
}

// ClusterScope defines a scope defined around a Power VS Cluster.
type ClusterScope struct {
	Client      client.Client
	patchHelper *patch.Helper

	// IBM Cloud Clients
	IBMPowerVSClient      powervs.PowerVS
	IBMVPCClient          vpc.Vpc
	TransitGatewayClient  transitgateway.TransitGateway
	ResourceClient        resourcecontroller.ResourceController
	COSClient             cos.Cos
	ResourceManagerClient resourcemanager.ResourceManager

	// ClientBuilder is retained so that deferred client construction
	// (e.g. setupCOSClient) can use the same injected builder rather than
	// hard-coding the production implementation.
	ClientBuilder ClientBuilder

	// Recorder is the event recorder for emitting Kubernetes events.
	Recorder cgrecord.EventRecorder

	// Kubernetes Resources
	Cluster           *clusterv1.Cluster
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
}

// NewPowerVSClusterScope creates a new ClusterScope from the supplied parameters.
func NewPowerVSClusterScope(ctx context.Context, params ClusterScopeParams) (*ClusterScope, error) {
	// 1. Validate inputs
	if err := params.validate(); err != nil {
		return nil, err
	}

	// 2. Initialize the patch helper
	helper, err := patch.NewHelper(params.IBMPowerVSCluster, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// 3. Construct the base scope
	clusterScope := &ClusterScope{
		Client:            params.Client,
		patchHelper:       helper,
		Cluster:           params.Cluster,
		IBMPowerVSCluster: params.IBMPowerVSCluster,
		ServiceEndpoint:   params.ServiceEndpoint,
		ClientBuilder:     params.ClientBuilder,
		Recorder:          params.Recorder,
	}

	// 4. Incase of VIP Topology we need not build clients
	if params.IBMPowerVSCluster.Spec.Topology == infrav1.PowerVSVirtualIPTopology {
		return clusterScope, nil
	}

	// 5. Initialize IBM Cloud Clients
	if err := clusterScope.initClients(ctx, &params); err != nil {
		return nil, fmt.Errorf("failed to initialize IBM Cloud clients: %w", err)
	}

	return clusterScope, nil
}

// Validate ensures all required fields are present before scope creation.
func (p *ClusterScopeParams) validate() error {
	if p.Client == nil {
		return errors.New("failed to generate new scope: client is nil")
	}
	if p.Cluster == nil {
		return errors.New("failed to generate new scope: cluster is nil")
	}
	if p.IBMPowerVSCluster == nil {
		return errors.New("failed to generate new scope: IBMPowerVSCluster is nil")
	}
	if p.ClientBuilder == nil {
		return errors.New("failed to generate new scope: ClientBuilder is nil")
	}
	return nil
}

// initClients bootstraps the required IBM Cloud clients based on the current cluster state and topology.
func (s *ClusterScope) initClients(ctx context.Context, params *ClusterScopeParams) error {
	log := ctrl.LoggerFrom(ctx)

	// Build the authenticator
	auth, err := params.ClientBuilder.GetAuthenticator(ctx)
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	// Build the unified ClientOptions
	opts := ClientOptions{
		Authenticator:   auth,
		Zone:            s.IBMPowerVSCluster.Spec.Zone,
		VPCRegion:       s.IBMPowerVSCluster.Spec.VPC.Region,
		WorkspaceID:     s.IBMPowerVSCluster.Status.Workspace.ID,
		ServiceEndpoint: s.ServiceEndpoint,
		Debug:           log.V(DEBUGLEVEL).Enabled(),
	}

	// 1. Build ResourceController & ResourceManager
	s.ResourceClient, err = params.ClientBuilder.GetResourceControllerClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create Resource Controller client: %w", err)
	}
	s.ResourceManagerClient, err = params.ClientBuilder.GetResourceManagerClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create Resource Manager client: %w", err)
	}

	// 2. Build PowerVS client
	// If WorkspaceID is empty, Its needed for Validating PER where workspace ID is not needed.
	s.IBMPowerVSClient, err = params.ClientBuilder.GetPowerVSClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create PowerVS client: %w", err)
	}

	if opts.WorkspaceID == "" {
		log.V(3).Info("Workspace ID not yet populated; PowerVS client initialized in Zone-only mode")
	}

	// 3. Build VPC and TransitGateway client
	s.IBMVPCClient, err = params.ClientBuilder.GetVPCClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create VPC client: %w", err)
	}

	s.TransitGatewayClient, err = params.ClientBuilder.GetTransitGatewayClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create Transit Gateway client: %w", err)
	}

	return nil
}

// PatchObject persists the cluster configuration and status.
func (s *ClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMPowerVSCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.PatchObject()
}

// Name returns the CAPI cluster name.
func (s *ClusterScope) Name() string {
	return s.Cluster.Name
}

// Zone returns the cluster zone.
func (s *ClusterScope) Zone() string {
	return s.IBMPowerVSCluster.Spec.Zone
}

// GetResourceGroupID returns the resource group ID.
// It first checks the Spec (user-provided), then falls back to Status (resolved).
func (s *ClusterScope) GetResourceGroupID() string {
	// Spec takes precedence over Status
	if s.IBMPowerVSCluster.Spec.ResourceGroup.Reference.ID != "" {
		return s.IBMPowerVSCluster.Spec.ResourceGroup.Reference.ID
	}
	return s.IBMPowerVSCluster.Status.ResourceGroup.ID
}

// ResourceGroupName returns the resource group name.
// It first checks the Status (resolved), then falls back to Spec (user-provided).
func (s *ClusterScope) ResourceGroupName() string {
	if s.IBMPowerVSCluster.Status.ResourceGroup.Name != "" {
		return s.IBMPowerVSCluster.Status.ResourceGroup.Name
	}
	return s.IBMPowerVSCluster.Spec.ResourceGroup.Reference.Name
}

// InfraCluster returns the IBMPowerVS infrastructure cluster object name.
func (s *ClusterScope) InfraCluster() string {
	return s.IBMPowerVSCluster.Name
}

// APIServerPort returns the APIServerPort to use when creating the ControlPlaneEndpoint.
func (s *ClusterScope) APIServerPort() int32 {
	if s.Cluster.Spec.ClusterNetwork.APIServerPort > 0 {
		return s.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return infrav1.DefaultAPIServerPort
}

// GetLoadBalancerID returns the cached load balancer ID from the status slice if it exists.
func (s *ClusterScope) GetLoadBalancerID(name string) string {
	for _, lb := range s.IBMPowerVSCluster.Status.LoadBalancers {
		if lb.Name == name {
			return lb.ID
		}
	}
	return ""
}

// GetPublicLoadBalancerHostName will return the hostname of the public load balancer.
func (s *ClusterScope) GetPublicLoadBalancerHostName() (*string, error) {
	cluster := s.IBMPowerVSCluster

	// If no load balancers have been tracked in status yet, there's nothing to return.
	if len(cluster.Status.LoadBalancers) == 0 {
		return nil, nil
	}

	// Case 1: If no load balancers are specified in the spec, it defaults to a single
	// auto-managed public load balancer using the default service name.
	if len(cluster.Spec.LoadBalancers) == 0 {
		defaultName := ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeLBPublic, "")

		for _, lbStatus := range cluster.Status.LoadBalancers {
			if lbStatus.Name == defaultName && lbStatus.Hostname != "" {
				return ptr.To(lbStatus.Hostname), nil
			}
		}
		return nil, nil
	}

	// Case 2: Evaluate explicitly configured load balancers to locate the public one.
	for i, lb := range cluster.Spec.LoadBalancers {
		var targetName string

		switch lb.Type {
		case infrav1.SourceTypeProvision:
			// Default to Public if not explicitly overridden to Private
			if lb.Provision.Type == infrav1.LoadBalancerTypePrivate {
				continue
			}

			targetName = lb.Provision.Name
			if targetName == "" {
				targetName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeLBPublic, strconv.Itoa(i))
			}

		case infrav1.SourceTypeReference:
			// For references, resolve the name via ID or use the explicit name provided
			if lb.Reference.Name != "" {
				targetName = lb.Reference.Name
			} else if lb.Reference.ID != "" {
				// Fallback to a live API fetch if we only have an ID
				lbDetails, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
					ID: ptr.To(lb.Reference.ID),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to fetch referenced load balancer (%s) details: %w", lb.Reference.ID, err)
				}
				if lbDetails != nil && lbDetails.Name != nil {
					targetName = *lbDetails.Name
				}
			}
		}

		// Search status array for a name match
		if targetName != "" {
			for _, lbStatus := range cluster.Status.LoadBalancers {
				if lbStatus.Name == targetName && lbStatus.Hostname != "" {
					return ptr.To(lbStatus.Hostname), nil
				}
			}
		}
	}

	return nil, nil
}

// ValidateZoneSupportsPER checks whether PowerVS zone supports PER capabilities.
func (s *ClusterScope) ValidateZoneSupportsPER(ctx context.Context) error {
	zone := s.IBMPowerVSCluster.Spec.Zone

	if zone == "" {
		return fmt.Errorf("PowerVS zone is required but not set in the spec")
	}

	// Fetch the datacenter details for the specified zone.
	datacenterDetails, err := s.IBMPowerVSClient.GetDatacenterDetails(ctx, zone)
	if err != nil {
		return fmt.Errorf("failed to get datacenter details: %w", err)
	}
	if datacenterDetails == nil || datacenterDetails.Capabilities == nil {
		return fmt.Errorf("failed to get datacenter details for zone: %s", zone)
	}
	// check for the PER support in datacenter capabilities.
	perAvailable, ok := datacenterDetails.Capabilities[powerEdgeRouter]
	if !ok {
		return fmt.Errorf("%s capability unknown for zone %q", powerEdgeRouter, zone)
	}
	if !perAvailable {
		return fmt.Errorf("%s is not available for zone %q", powerEdgeRouter, zone)
	}
	return nil
}

// ReconcileResourceGroup reconciles resource group to fetch resource group id.
func (s *ClusterScope) ReconcileResourceGroup(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)

	// 1. Find the ID: Check Status first (already resolved), then fallback to Spec (user provided).
	resourceGroupID := s.IBMPowerVSCluster.Status.ResourceGroup.ID
	if resourceGroupID == "" {
		resourceGroupID = s.IBMPowerVSCluster.Spec.ResourceGroup.Reference.ID
	}

	// 2. ID exists: Verify it in IBM Cloud and hydrate the Status.
	if resourceGroupID != "" {
		log.V(3).Info("Resource group ID is set in status, fetching details", "resourceGroupID", resourceGroupID)

		resourceGroup, _, err := s.ResourceManagerClient.GetResourceGroup(&resourcemanagerv2.GetResourceGroupOptions{
			ID: &resourceGroupID,
		})
		if err != nil {
			return fmt.Errorf("failed to fetch resource group (id: %s) details: %w", resourceGroupID, err)
		}

		if resourceGroup == nil {
			return fmt.Errorf("resource group not found with ID: %s", resourceGroupID)
		}

		s.IBMPowerVSCluster.Status.ResourceGroup.ID = resourceGroupID
		if resourceGroup.Name != nil {
			s.IBMPowerVSCluster.Status.ResourceGroup.Name = *resourceGroup.Name
		}

		return nil
	}

	// 3. No ID exists anywhere: The user must have provided only a Name in the Spec.
	fetchedID, err := s.resolveResourceGroupIDByName()
	if err != nil {
		return fmt.Errorf("failed to resolve resource group ID by name: %w", err)
	}

	log.Info("Successfully fetched resource group ID from cloud", "resourceGroupID", fetchedID)

	// 4. Save to Status.
	s.IBMPowerVSCluster.Status.ResourceGroup.ID = fetchedID
	s.IBMPowerVSCluster.Status.ResourceGroup.Name = s.IBMPowerVSCluster.Spec.ResourceGroup.Reference.Name

	return nil
}

// resolveResourceGroupIDByName retrieves the ID of the resource group from IBM Cloud using its name.
func (s *ClusterScope) resolveResourceGroupIDByName() (string, error) {
	resourceGroupName := s.IBMPowerVSCluster.Spec.ResourceGroup.Reference.Name
	if resourceGroupName == "" {
		return "", fmt.Errorf("resource group name is not set in the spec")
	}

	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return "", fmt.Errorf("failed to get authenticator: %w", err)
	}

	account, err := accounts.GetAccount(auth)
	if err != nil {
		return "", fmt.Errorf("failed to get account: %w", err)
	}

	rmv2ListResourceGroupOpt := resourcemanagerv2.ListResourceGroupsOptions{
		Name:      &resourceGroupName,
		AccountID: &account,
	}

	resourceGroupListResult, _, err := s.ResourceManagerClient.ListResourceGroups(&rmv2ListResourceGroupOpt)
	if err != nil {
		return "", fmt.Errorf("failed to list resource groups: %w", err)
	}

	if resourceGroupListResult != nil && len(resourceGroupListResult.Resources) > 0 {
		rg := resourceGroupListResult.Resources[0]
		if rg.ID != nil {
			return *rg.ID, nil
		}
	}

	return "", fmt.Errorf("could not retrieve resource group ID for %q", resourceGroupName)
}

// ReconcileWorkspace reconciles PowerVS workspace.
func (s *ClusterScope) ReconcileWorkspace(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// 1. Idempotency & State Check: If we already resolved the Workspace ID, just verify its state.
	workspaceID := cluster.Status.Workspace.ID
	if workspaceID != "" {
		log.V(3).Info("PowerVS workspace ID is set, fetching details", "workspaceID", workspaceID)

		workspace, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: &workspaceID,
		})
		if err != nil {
			return false, fmt.Errorf("failed to fetch workspace (id: %s) details: %w", workspaceID, err)
		}

		if workspace == nil {
			return false, fmt.Errorf("workspace not found with ID: %s", workspaceID)
		}

		requeue, err := s.checkWorkspaceState(ctx, *workspace)
		if err != nil {
			return false, fmt.Errorf("failed to check workspace state: %w", err)
		}
		return requeue, nil
	}

	// 2. We don't have an ID yet. Route logic based strictly on the user's explicit intent.
	log.Info("Resolving PowerVS workspace", "type", cluster.Spec.Workspace.Type)

	switch cluster.Spec.Workspace.Type {
	case infrav1.SourceTypeReference:
		return s.reconcileWorkspaceReference(ctx)

	case infrav1.SourceTypeProvision:
		return s.reconcileWorkspaceProvision(ctx)

	default:
		return false, fmt.Errorf("unknown workspace source type: %s", cluster.Spec.Workspace.Type)
	}
}

// reconcileWorkspaceReference handles the logic when a user brings their own existing workspace.
func (s *ClusterScope) reconcileWorkspaceReference(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	ref := s.IBMPowerVSCluster.Spec.Workspace.Reference

	log.Info("Verifying existing workspace", "reference", ref)

	resourceInstance := resourcecontroller.InstanceFilter{
		ID:             ref.ID,
		Name:           ref.Name,
		Zone:           &s.IBMPowerVSCluster.Spec.Zone,
		ResourceID:     resourcecontroller.PowerVSResourceID,
		ResourcePlanID: resourcecontroller.PowerVSResourcePlanID,
	}

	workspace, err := s.ResourceClient.GetResourceInstanceByFilter(resourceInstance)
	if err != nil {
		return false, fmt.Errorf("failed to fetch workspace by ref %q: %w", ref, err)
	}

	if workspace == nil {
		return false, fmt.Errorf("workspace with reference %q not found in IBM Cloud", ref)
	}

	if workspace.GUID == nil || workspace.Name == nil {
		return false, fmt.Errorf("workspace %q has missing GUID or name", ref)
	}

	log.Info("Successfully verified existing workspace", "workspaceID", *workspace.GUID, "name", workspace.Name)

	s.IBMPowerVSCluster.Status.Workspace = infrav1.ResourceReference{ID: *workspace.GUID, Name: *workspace.Name}

	return true, nil // requeue so that the state of worksapce will be checked in the next reconcile
}

// reconcileWorkspaceProvision handles the logic when the controller must create a new workspace.
func (s *ClusterScope) reconcileWorkspaceProvision(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	provision := s.IBMPowerVSCluster.Spec.Workspace.Provision

	// 1. Determine the name to use
	workspaceName := provision.Name
	if workspaceName == "" {
		workspaceName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeWorkspace, "")
	}

	// 2. Idempotency check: Did we already create this, but crash before saving to Status?
	resourceInstance := resourcecontroller.InstanceFilter{
		Name:           workspaceName,
		Zone:           &s.IBMPowerVSCluster.Spec.Zone,
		ResourceID:     resourcecontroller.PowerVSResourceID,
		ResourcePlanID: resourcecontroller.PowerVSResourcePlanID,
	}

	workspace, err := s.ResourceClient.GetResourceInstanceByFilter(resourceInstance)
	if err != nil {
		return false, fmt.Errorf("failed to check for existing workspace: %w", err)
	}

	if workspace != nil && workspace.GUID != nil {
		log.Info("Recovered previously provisioned workspace", "workspaceID", workspace.GUID)
		s.IBMPowerVSCluster.Status.Workspace = infrav1.ResourceReference{
			ID:   *workspace.GUID,
			Name: workspaceName,
		}
		return true, nil // requeue so that the state of worksapce will be checked in the next reconcile
	}

	// 3. Create the new Workspace
	log.Info("Provisioning new workspace", "name", workspaceName)

	workspace, err = s.createWorkspace(ctx, workspaceName)
	if err != nil {
		return false, fmt.Errorf("failed to provision workspace: %w", err)
	}
	if workspace == nil || workspace.GUID == nil {
		return false, errors.New("provisioned workspace or GUID is nil")
	}

	log.Info("Successfully provisioned workspace", "workspaceID", *workspace.GUID)

	// 4. Save to Status.
	s.IBMPowerVSCluster.Status.Workspace = infrav1.ResourceReference{
		ID:   *workspace.GUID,
		Name: workspaceName,
	}

	return true, nil // Requeue to wait for it to become ACTIVE
}

// checkWorkspaceState checks the state of a PowerVS workspace.
// If state is provisioning, true is returned indicating a requeue for reconciliation.
// In all other cases, it returns false.
func (s *ClusterScope) checkWorkspaceState(ctx context.Context, workspace resourcecontrollerv2.ResourceInstance) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("Checking the state of PowerVS workspace", "name", *workspace.Name)

	switch *workspace.State {
	case string(infrav1.WorkspaceStateActive):
		log.V(3).Info("PowerVS workspace is in active state")
		return false, nil
	case string(infrav1.WorkspaceStateProvisioning):
		log.V(3).Info("PowerVS workspace is in provisioning state")
		return true, nil
	case string(infrav1.WorkspaceStateFailed):
		return false, fmt.Errorf("PowerVS workspace is in failed state")
	}
	return false, fmt.Errorf("PowerVS workspacee is in %s state", *workspace.State)
}

// createServiceInstance creates the service instance.
func (s *ClusterScope) createWorkspace(ctx context.Context, workspaceName string) (*resourcecontrollerv2.ResourceInstance, error) {
	log := ctrl.LoggerFrom(ctx)

	// fetch resource group id.
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroupName())
	}

	zone := s.Zone()
	if zone == "" {
		return nil, fmt.Errorf("PowerVS zone is not set")
	}

	// create worksapce.
	log.V(3).Info("Creating new worksapce", "worksapceName", workspaceName, "zone", zone)

	workspace, _, err := s.ResourceClient.CreateResourceInstance(&resourcecontrollerv2.CreateResourceInstanceOptions{
		Name:           &workspaceName,
		Target:         &zone,
		ResourceGroup:  &resourceGroupID,
		ResourcePlanID: ptr.To(resourcecontroller.PowerVSResourcePlanID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create worksapce: %w", err)
	}
	return workspace, nil
}

// ReconcileNetwork reconciles network.
func (s *ClusterScope) ReconcileNetwork(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	cluster := s.IBMPowerVSCluster

	// 1. Idempotency & State Check: If we already resolved the Network ID, just verify its state.
	networkID := cluster.Status.Network.ID
	if networkID != "" {
		log.V(3).Info("PowerVS network ID is set in status, verifying existence", "networkID", networkID)
		if _, err := s.IBMPowerVSClient.GetNetworkByID(ctx, networkID); err != nil {
			return false, fmt.Errorf("failed to fetch network by ID: %w", err)
		}

		// If we provisioned this network via DHCP, ensure the DHCP server is fully active
		if cluster.Spec.Network.Type == infrav1.SourceTypeProvision {
			dhcpServerID := cluster.Status.Network.DHCPServer.ID
			if dhcpServerID == "" {
				log.Info("Recovering state: Network ID is present but DHCP Server ID is missing in status. Requeuing to resolve", "networkID", networkID)
				return true, nil
			}

			log.V(3).Info("Verifying provisioned DHCP server state", "dhcpServerID", dhcpServerID)
			active, err := s.isDHCPServerActive(ctx)
			if err != nil {
				return false, fmt.Errorf("failed to check if DHCP server is active: %w", err)
			}

			if !active {
				log.V(3).Info("DHCP server is still building")
				return true, nil // requeue and wait
			}
		}

		// Network is resolved and ready!
		return false, nil
	}

	// 2. We don't have a Network ID yet. Route logic based strictly on the user's explicit intent.
	log.Info("Resolving PowerVS network", "type", cluster.Spec.Network.Type)

	switch cluster.Spec.Network.Type {
	case infrav1.SourceTypeReference:
		return s.reconcileNetworkReference(ctx)

	case infrav1.SourceTypeProvision:
		return s.reconcileNetworkProvision(ctx)

	default:
		return false, fmt.Errorf("unknown network source type: %q", cluster.Spec.Network.Type)
	}
}

// reconcileNetworkReference handles the logic when a user brings their own existing network.
func (s *ClusterScope) reconcileNetworkReference(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	ref := s.IBMPowerVSCluster.Spec.Network.Reference

	log.Info("Verifying existing network", "reference", ref)

	var networkID, networkName string

	if ref.ID != "" {
		network, err := s.IBMPowerVSClient.GetNetworkByID(ctx, ref.ID)
		if err != nil {
			return false, fmt.Errorf("failed to fetch network by ID %q: %w", ref.ID, err)
		}
		if network == nil || network.NetworkID == nil || network.Name == nil {
			return false, fmt.Errorf("invalid network payload received from IBM cloud for ID %q: network object, ID, or Name is nil", ref.ID)
		}
		networkID = *network.NetworkID
		networkName = *network.Name
	} else if ref.Name != "" {
		network, err := s.IBMPowerVSClient.GetNetworkByName(ctx, ref.Name)
		if err != nil {
			return false, fmt.Errorf("failed to fetch network by name %q: %w", ref.Name, err)
		}
		if network == nil || network.NetworkID == nil || network.Name == nil {
			return false, fmt.Errorf("invalid network payload received from IBM cloud for name %q: network object, ID, or Name is nil", ref.Name)
		}
		networkID = *network.NetworkID
		networkName = *network.Name
	} else {
		return false, fmt.Errorf("network reference must contain either an ID or a Name")
	}

	log.Info("Successfully verified existing network", "networkID", networkID, "networkName", networkName)

	s.IBMPowerVSCluster.Status.Network.ID = networkID
	s.IBMPowerVSCluster.Status.Network.Name = networkName

	return true, nil // requeue so the fast-path verifies it
}

// reconcileNetworkProvision handles the logic when the controller must create a new DHCP server and Network.
func (s *ClusterScope) reconcileNetworkProvision(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	dhcpSpec := s.IBMPowerVSCluster.Spec.Network.Provision.DHCPServer

	// 1. Determine the exact name to use for the DHCP server
	dhcpName := dhcpSpec.Name
	if dhcpName == "" {
		dhcpName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeDHCP, "")
	}

	// 2. Idempotency check: Did we already create this DHCP server, but crash before saving to Status?
	dhcpServers, err := s.IBMPowerVSClient.ListDHCPServers(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to fetch existing DHCP servers for idempotency check: %w", err)
	}

	expectedNetworkName := dhcpNetworkName(dhcpName)
	for _, server := range dhcpServers {
		// Identify by looking at the network name IBM Cloud generated for it
		if server.Network == nil || server.Network.Name == nil || *server.Network.Name != expectedNetworkName {
			continue
		}
		if server.ID == nil || server.Network.ID == nil {
			log.V(4).Info("Skipping malformed DHCP server record from IBM Cloud (missing ID or Network)")
			continue
		}
		log.Info("Recovered previously provisioned DHCP server", "dhcpServerID", *server.ID, "networkID", *server.Network.ID)

		// Save recovered IDs directly to Status
		s.IBMPowerVSCluster.Status.Network.ID = *server.Network.ID
		s.IBMPowerVSCluster.Status.Network.Name = *server.Network.Name
		s.IBMPowerVSCluster.Status.Network.DHCPServer.ID = *server.ID

		return true, nil // requeue
	}

	// 3. Create the new DHCP Server and Network
	log.Info("Provisioning new DHCP Server and Network", "name", dhcpName)

	dhcpServerID, networkID, err := s.createDHCPServer(ctx, dhcpName)
	if err != nil {
		return false, fmt.Errorf("failed to provision DHCP server: %w", err)
	}

	log.Info("Successfully triggered DHCP Server provision", "dhcpServerID", dhcpServerID, "networkID", networkID)

	// 4. Save directly to the new v1beta3 Status value types
	s.IBMPowerVSCluster.Status.Network.ID = networkID
	s.IBMPowerVSCluster.Status.Network.Name = expectedNetworkName
	s.IBMPowerVSCluster.Status.Network.DHCPServer.ID = dhcpServerID

	return true, nil // Requeue to wait for it to become ACTIVE
}

// isDHCPServerActive checks if the DHCP server status is active.
func (s *ClusterScope) isDHCPServerActive(ctx context.Context) (bool, error) {
	dhcpID := s.IBMPowerVSCluster.Status.Network.DHCPServer.ID

	dhcpServer, err := s.IBMPowerVSClient.GetDHCPServer(ctx, dhcpID)
	if err != nil {
		return false, err
	}

	if dhcpServer == nil {
		return false, fmt.Errorf("DHCP server details are nil for ID: %s", dhcpID)
	}

	return s.checkDHCPServerStatus(ctx, *dhcpServer)
}

// checkDHCPServerStatus checks the state of a DHCP server.
func (s *ClusterScope) checkDHCPServerStatus(ctx context.Context, dhcpServer models.DHCPServerDetail) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	if dhcpServer.Status == nil {
		return false, fmt.Errorf("DHCP server status is nil")
	}
	log.V(3).Info("Checking the status of DHCP server", "state", *dhcpServer.Status)

	switch *dhcpServer.Status {
	case string(infrav1.DHCPServerStateActive):
		return true, nil
	case string(infrav1.DHCPServerStateBuild):
		return false, nil
	case string(infrav1.DHCPServerStateError):
		return false, fmt.Errorf("DHCP server creation failed and is in error state")
	default:
		return false, fmt.Errorf("DHCP server is in an unknown state: %s", *dhcpServer.Status)
	}
}

// createDHCPServer creates the DHCP server and returns its ID and its associated Network ID.
func (s *ClusterScope) createDHCPServer(ctx context.Context, dhcpName string) (string, string, error) {
	log := ctrl.LoggerFrom(ctx)

	dhcpSpec := s.IBMPowerVSCluster.Spec.Network.Provision.DHCPServer

	params := models.DHCPServerCreate{
		Name: &dhcpName,
	}

	if dhcpSpec.CIDR != "" {
		params.Cidr = &dhcpSpec.CIDR
	}
	if dhcpSpec.DNSServer != "" {
		params.DNSServer = &dhcpSpec.DNSServer
	}

	snatEnabled := true // default
	if dhcpSpec.Snat == infrav1.DHCPSnatPolicyDisabled {
		snatEnabled = false
	}
	params.SnatEnabled = &snatEnabled

	dhcpServer, err := s.IBMPowerVSClient.CreateDHCPServer(ctx, &params)
	if err != nil {
		return "", "", err
	}
	if dhcpServer == nil || dhcpServer.ID == nil {
		return "", "", fmt.Errorf("created DHCP server or its ID is nil")
	}
	if dhcpServer.Network == nil || dhcpServer.Network.ID == nil {
		return "", "", fmt.Errorf("created DHCP server network or its ID is nil")
	}

	log.V(3).Info("DHCP Server network details", "details", *dhcpServer.Network)

	return *dhcpServer.ID, *dhcpServer.Network.ID, nil
}

// ReconcileTransitGateway reconcile transit gateway.
func (s *ClusterScope) ReconcileTransitGateway(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	// Skip TG reconciliation entirely if the topology is VirtualIP.
	if s.IBMPowerVSCluster.Spec.Topology == infrav1.PowerVSVirtualIPTopology {
		return false, nil
	}

	var tg *tgapiv1.TransitGateway
	var err error

	// 1. Idempotency & State Check: If we already resolved the TG ID, just verify its state.
	tgID := s.IBMPowerVSCluster.Status.TransitGateway.ID
	if tgID != "" {
		log.V(3).Info("Transit Gateway ID is set in status, fetching details", "tgID", tgID)
		tg, _, err = s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
			ID: ptr.To(tgID),
		})
		if err != nil {
			return false, fmt.Errorf("failed to fetch transit gateway (id: %s) details: %w", tgID, err)
		}
		if tg == nil {
			return false, fmt.Errorf("transit gateway not found with ID: %s", tgID)
		}

		// Check status and update connections
		return s.checkAndUpdateTransitGateway(ctx, tg)
	}

	tgSpec := s.IBMPowerVSCluster.Spec.TransitGateway
	log.Info("Resolving Transit Gateway", "type", tgSpec.Type)

	switch tgSpec.Type {
	case infrav1.SourceTypeReference:
		tg, err = s.resolveTransitGatewayReference(ctx, tgSpec.Reference)
		if err != nil {
			return false, err
		}

		if tg == nil || tg.ID == nil || tg.Name == nil {
			return false, fmt.Errorf("transit gateway reference resolved, but IBM Cloud returned a nil ID or Name")
		}

		s.IBMPowerVSCluster.Status.TransitGateway = infrav1.TransitGatewayStatus{
			ID:   *tg.ID,
			Name: *tg.Name,
		}

		// Check status and update connections
		return s.checkAndUpdateTransitGateway(ctx, tg)

	case infrav1.SourceTypeProvision:
		log.Info("Creating transit gateway")
		tg, err = s.provisionTransitGateway(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to provision transit gateway: %w", err)
		}

		if tg == nil || tg.ID == nil || tg.Name == nil {
			return false, fmt.Errorf("transit gateway reference resolved, but IBM Cloud returned a nil ID or Name")
		}

		s.IBMPowerVSCluster.Status.TransitGateway = infrav1.TransitGatewayStatus{
			ID:   *tg.ID,
			Name: *tg.Name,
		}

		// Newly created TG is not ready for connections.
		return true, nil

	default:
		return false, fmt.Errorf("unknown transit gateway source type: %s", tgSpec.Type)
	}
}

// resolveTransitGatewayReference fetches an existing TG strictly by ID or Name.
func (s *ClusterScope) resolveTransitGatewayReference(_ context.Context, ref infrav1.ResourceIdentifier) (*tgapiv1.TransitGateway, error) {
	if ref.ID != "" {
		tg, _, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{ID: ptr.To(ref.ID)})
		if err != nil {
			return nil, fmt.Errorf("failed to get transit gateway by ID %q: %w", ref.ID, err)
		}
		return tg, nil
	}

	if ref.Name != "" {
		tg, err := s.TransitGatewayClient.GetTransitGatewayByName(ref.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get transit gateway by name %q: %w", ref.Name, err)
		}
		if tg == nil {
			return nil, fmt.Errorf("transit gateway with name %q not found", ref.Name)
		}
		return tg, nil
	}

	return nil, fmt.Errorf("transit gateway reference must have either ID or Name")
}

// provisionTransitGateway creates a new TG if it doesn't already exist.
func (s *ClusterScope) provisionTransitGateway(ctx context.Context) (*tgapiv1.TransitGateway, error) {
	log := ctrl.LoggerFrom(ctx)
	tgSpec := s.IBMPowerVSCluster.Spec.TransitGateway.Provision

	// Determine TG Name
	tgName := tgSpec.Name
	if tgName == "" {
		tgName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeTransitGateway, "")
	}

	// Idempotency: Check if we already created it
	if existingTG, _ := s.TransitGatewayClient.GetTransitGatewayByName(tgName); existingTG != nil && existingTG.ID != nil {
		log.V(3).Info("Transit Gateway already exists", "name", tgName)
		return existingTG, nil
	}

	// Fetch required resource group ID
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroupName())
	}

	if s.IBMPowerVSCluster.Status.Workspace.ID == "" || s.IBMPowerVSCluster.Status.VPC.ID == "" {
		return nil, fmt.Errorf("failed to proceed with transit gateway creation: PowerVS workspace or VPC reconciliation is not yet complete")
	}

	// Determine Routing
	location, sysGlobalRouting, err := ibmcloudutil.GetTransitGatewayLocationAndRouting(ptr.To(s.Zone()), &s.IBMPowerVSCluster.Status.VPC.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get transit gateway location and routing: %w", err)
	}

	// The Webhook already prevents illegal "Local" overrides.
	// We just apply the system default, unless the user explicitly requested Global.
	globalRouting := *sysGlobalRouting
	if tgSpec.GlobalRouting == infrav1.TransitGatewayRoutingGlobal {
		globalRouting = true
	}

	// Create TG
	log.Info("Creating Transit Gateway in IBM Cloud", "name", tgName)
	tg, _, err := s.TransitGatewayClient.CreateTransitGateway(&tgapiv1.CreateTransitGatewayOptions{
		Location:      location,
		Name:          ptr.To(tgName),
		Global:        ptr.To(globalRouting),
		ResourceGroup: &tgapiv1.ResourceGroupIdentity{ID: ptr.To(resourceGroupID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transit gateway: %w", err)
	}

	return tg, nil
}

// checkAndUpdateTransitGateway checks given transit gateway's status and its connections.
func (s *ClusterScope) checkAndUpdateTransitGateway(ctx context.Context, tg *tgapiv1.TransitGateway) (bool, error) {
	requeue, err := s.checkTransitGatewayStatus(ctx, tg)
	if err != nil {
		return false, err
	}
	if requeue {
		return true, nil
	}

	return s.checkAndUpdateTransitGatewayConnections(ctx, tg)
}

// checkTransitGatewayStatus checks the state of a transit gateway safely.
func (s *ClusterScope) checkTransitGatewayStatus(ctx context.Context, tg *tgapiv1.TransitGateway) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	if tg.Status == nil {
		return false, fmt.Errorf("transit gateway returned from IBM Cloud is missing Status")
	}

	log.V(3).Info("Checking the status of transit gateway", "name", *tg.Name)
	switch *tg.Status {
	case string(infrav1.TransitGatewayStateAvailable):
		log.V(3).Info("Transit gateway is in available state")
		return false, nil
	case string(infrav1.TransitGatewayStateFailed):
		return false, fmt.Errorf("failed to create transit gateway, current status: %s", *tg.Status)
	case string(infrav1.TransitGatewayStatePending):
		log.V(3).Info("Transit gateway is in pending state")
		return true, nil
	default:
		return false, fmt.Errorf("transit gateway is in unknown state: %s", *tg.Status)
	}
}

// checkAndUpdateTransitGatewayConnections reconciles connections based on explicit user intent.
func (s *ClusterScope) checkAndUpdateTransitGatewayConnections(ctx context.Context, transitGateway *tgapiv1.TransitGateway) (bool, error) {
	tgConnections, _, err := s.TransitGatewayClient.ListTransitGatewayConnections(&tgapiv1.ListTransitGatewayConnectionsOptions{
		TransitGatewayID: transitGateway.ID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to list transit gateway connections: %w", err)
	}

	vpcCRN, err := s.fetchVPCCRN()
	if err != nil {
		return false, fmt.Errorf("failed to fetch VPC CRN: %w", err)
	}

	pvsServiceInstanceCRN, err := s.fetchPowerVSWorkspaceCRN()
	if err != nil {
		return false, fmt.Errorf("failed to fetch PowerVS service instance CRN: %w", err)
	}

	tgSpec := s.IBMPowerVSCluster.Spec.TransitGateway

	// Reconcile VPC Connection based on intent.
	requeueVPC, err := s.reconcileConnection(ctx, transitGateway, tgSpec.VPCConnection, vpcCRN, vpcNetworkConnectionType, tgConnections.Connections)
	if err != nil {
		return false, fmt.Errorf("failed to reconcile VPC connection: %w", err)
	}

	// Reconcile PowerVS Connection based on intent.
	requeuePVS, err := s.reconcileConnection(ctx, transitGateway, tgSpec.PowerVSConnection, pvsServiceInstanceCRN, powervsNetworkConnectionType, tgConnections.Connections)
	if err != nil {
		return false, fmt.Errorf("failed to reconcile PowerVS connection: %w", err)
	}

	// Return the combined requeue status cleanly.
	return requeueVPC || requeuePVS, nil
}

// reconcileConnection evaluates intent, routes to the appropriate handler, and returns the requeue state.
func (s *ClusterScope) reconcileConnection(ctx context.Context, tg *tgapiv1.TransitGateway, connSpec infrav1.TransitGatewayConnectionSource, networkID *string, netType networkConnectionType, existingConns []tgapiv1.TransitGatewayConnectionCust) (bool, error) {
	switch connSpec.Type {
	case infrav1.SourceTypeReference:
		return s.reconcileConnectionReference(ctx, connSpec.Reference, networkID, netType, existingConns)

	case infrav1.SourceTypeProvision:
		return s.reconcileConnectionProvision(ctx, tg, connSpec.Provision, networkID, netType, existingConns)

	default:
		return false, fmt.Errorf("unknown connection source type: %s", connSpec.Type)
	}
}

// checkTransitGatewayConnectionStatus checks the state of a transit gateway connection safely.
func (s *ClusterScope) checkTransitGatewayConnectionStatus(ctx context.Context, con *tgapiv1.TransitGatewayConnectionCust) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	if con == nil || con.Status == nil || con.Name == nil {
		return false, fmt.Errorf("connection status or name is nil")
	}

	log.V(3).Info("Checking the status of transit gateway connection", "name", *con.Name)
	switch *con.Status {
	case string(infrav1.TransitGatewayConnectionStateAttached):
		return false, nil
	case string(infrav1.TransitGatewayConnectionStateFailed):
		return false, fmt.Errorf("failed to attach connection to transit gateway, current status: %s", *con.Status)
	case string(infrav1.TransitGatewayConnectionStatePending):
		log.V(3).Info("Transit gateway connection is in pending state")
		return true, nil
	default:
		return false, fmt.Errorf("transit gateway connection is in unknown state: %s", *con.Status)
	}
}

// setTransitGatewayConnectionStatus sets the connection status of the Transit Gateway.
func (s *ClusterScope) setTransitGatewayConnectionStatus(networkType networkConnectionType, id, name, connState string) {
	connStatus := infrav1.ResourceConnectionStatus{
		ID:    id,
		Name:  name,
		State: connState,
	}

	switch networkType {
	case powervsNetworkConnectionType:
		s.IBMPowerVSCluster.Status.TransitGateway.PowerVSConnection = connStatus
	case vpcNetworkConnectionType:
		s.IBMPowerVSCluster.Status.TransitGateway.VPCConnection = connStatus
	}
}

// reconcileConnectionReference handles strictly verifying an explicitly referenced connection.
func (s *ClusterScope) reconcileConnectionReference(ctx context.Context, ref infrav1.ResourceIdentifier, networkID *string, netType networkConnectionType, existingConns []tgapiv1.TransitGatewayConnectionCust) (bool, error) {
	foundConn := findConnectionByRef(existingConns, ref)
	if foundConn == nil {
		return false, fmt.Errorf("transit gateway connection reference (ID: %q, Name: %q) not found on Transit Gateway", ref.ID, ref.Name)
	}

	// Ensure the connection they referenced actually points to our cluster's network.
	if foundConn.NetworkID == nil || *foundConn.NetworkID != *networkID {
		return false, fmt.Errorf("referenced transit gateway connection exists, but it connects to the wrong network CRN")
	}

	s.setTransitGatewayConnectionStatus(netType, *foundConn.ID, *foundConn.Name, *foundConn.Status)
	return s.checkTransitGatewayConnectionStatus(ctx, foundConn)
}

// reconcileConnectionProvision handles creating a new connection or verifying an existing one idempotently.
func (s *ClusterScope) reconcileConnectionProvision(ctx context.Context, tg *tgapiv1.TransitGateway, provSpec infrav1.TransitGatewayConnectionProvision, networkID *string, netType networkConnectionType, existingConns []tgapiv1.TransitGatewayConnectionCust) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	// Idempotency check.
	foundConn := findConnectionByNetwork(existingConns, netType, *networkID)
	if foundConn != nil {
		if foundConn.ID == nil || foundConn.Name == nil || foundConn.Status == nil {
			return false, fmt.Errorf("IBM cloud returned nil fields for existing connection")
		}
		s.setTransitGatewayConnectionStatus(netType, *foundConn.ID, *foundConn.Name, *foundConn.Status)
		return s.checkTransitGatewayConnectionStatus(ctx, foundConn)
	}

	log.Info("Creating transit gateway connection", "type", netType)
	connName := provSpec.Name
	if connName == "" {
		if netType == vpcNetworkConnectionType {
			connName = tgVPCConnectionName(*tg.Name)
		} else {
			connName = tgPowerVSConnectionName(*tg.Name)
		}
	}

	newConn, _, err := s.TransitGatewayClient.CreateTransitGatewayConnection(&tgapiv1.CreateTransitGatewayConnectionOptions{
		TransitGatewayID: tg.ID,
		NetworkType:      ptr.To(string(netType)),
		NetworkID:        networkID,
		Name:             ptr.To(connName),
	})
	if err != nil {
		return false, err
	}

	if newConn == nil || newConn.ID == nil || newConn.Name == nil || newConn.Status == nil {
		return false, fmt.Errorf("IBM Cloud returned nil fields for new connection")
	}

	s.setTransitGatewayConnectionStatus(netType, *newConn.ID, *newConn.Name, *newConn.Status)
	return true, nil // Requeue since we just created it
}

// findConnectionByRef searches for an existing connection strictly by the user's ID or Name reference.
func findConnectionByRef(existingConns []tgapiv1.TransitGatewayConnectionCust, ref infrav1.ResourceIdentifier) *tgapiv1.TransitGatewayConnectionCust {
	for i, conn := range existingConns {
		if ref.ID != "" && conn.ID != nil && *conn.ID == ref.ID {
			return &existingConns[i]
		}
		if ref.Name != "" && conn.Name != nil && *conn.Name == ref.Name {
			return &existingConns[i]
		}
	}
	return nil
}

// findConnectionByNetwork searches for an existing connection that matches the target network CRN.
func findConnectionByNetwork(existingConns []tgapiv1.TransitGatewayConnectionCust, netType networkConnectionType, networkID string) *tgapiv1.TransitGatewayConnectionCust {
	for i, conn := range existingConns {
		if conn.NetworkType != nil && *conn.NetworkType == string(netType) && conn.NetworkID != nil && *conn.NetworkID == networkID {
			return &existingConns[i]
		}
	}
	return nil
}

// ReconcileVPC evaluates the user's intent and reconciles the IBM Cloud VPC accordingly.
func (s *ClusterScope) ReconcileVPC(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// 1. Idempotency & State Check: If we already resolved the VPC ID, just verify its state.
	vpcID := cluster.Status.VPC.ID
	if vpcID != "" {
		log.V(3).Info("VPC ID is set, fetching details", "vpcID", vpcID)

		vpcDetails, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
			ID: ptr.To(vpcID),
		})
		if err != nil {
			return false, fmt.Errorf("failed to fetch VPC (id: %s) details: %w", vpcID, err)
		}

		if vpcDetails == nil {
			return false, fmt.Errorf("vpc not found with ID: %s", vpcID)
		}

		if vpcDetails.Status != nil && *vpcDetails.Status == string(infrav1.VPCStatePending) {
			log.V(3).Info("VPC creation is in pending state")
			return true, nil
		}
		return false, nil
	}

	// 2. We don't have an ID yet. Route logic based strictly on the user's explicit intent.
	log.Info("Resolving IBM Cloud VPC", "type", cluster.Spec.VPC.Type)

	switch cluster.Spec.VPC.Type {
	case infrav1.SourceTypeReference:
		return s.reconcileVPCReference(ctx)

	case infrav1.SourceTypeProvision:
		return s.reconcileVPCProvision(ctx)

	default:
		return false, fmt.Errorf("unknown VPC source type: %s", cluster.Spec.VPC.Type)
	}
}

// reconcileVPCReference handles verifying an explicitly referenced VPC.
//
//nolint:revive,unparam // ctx parameter is used for logging context
func (s *ClusterScope) reconcileVPCReference(ctx context.Context) (bool, error) {
	ref := s.IBMPowerVSCluster.Spec.VPC.Reference

	var vpcDetails *vpcv1.VPC
	var err error

	if ref.ID != "" {
		vpcDetails, _, err = s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{ID: ptr.To(ref.ID)})
	} else if ref.Name != "" {
		vpcDetails, err = s.IBMVPCClient.GetVPCByName(ref.Name)
	} else {
		return false, fmt.Errorf("VPC reference must have either ID or Name set")
	}

	if err != nil {
		return false, fmt.Errorf("failed to get referenced VPC: %w", err)
	}
	if vpcDetails == nil || vpcDetails.ID == nil || vpcDetails.Name == nil {
		return false, fmt.Errorf("referenced VPC not found or returned nil fields")
	}

	// Update the Status using the new highly-specific VPCStatus struct and include the shared Region
	s.IBMPowerVSCluster.Status.VPC = infrav1.VPCStatus{
		ID:     *vpcDetails.ID,
		Name:   *vpcDetails.Name,
		Region: s.IBMPowerVSCluster.Spec.VPC.Region,
	}
	return false, nil
}

// reconcileVPCProvision handles idempotently creating a new VPC.
func (s *ClusterScope) reconcileVPCProvision(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	vpcName := s.IBMPowerVSCluster.Spec.VPC.Provision.Name
	if vpcName == "" {
		vpcName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeVPC, "")
	}

	// Check if a VPC with the target name already exists (e.g. from a previous partial run)
	log.Info("Checking whether VPC already exists by name", "name", vpcName)
	vpcDetails, err := s.IBMVPCClient.GetVPCByName(vpcName)
	if err != nil {
		return false, fmt.Errorf("failed to check if VPC exists: %w", err)
	}

	if vpcDetails != nil && vpcDetails.ID != nil && vpcDetails.Name != nil {
		log.V(3).Info("VPC found in cloud by name matching", "vpcID", *vpcDetails.ID)

		// Map matched VPC state to our specific VPCStatus type
		s.IBMPowerVSCluster.Status.VPC = infrav1.VPCStatus{
			ID:     *vpcDetails.ID,
			Name:   *vpcDetails.Name,
			Region: s.IBMPowerVSCluster.Spec.VPC.Region,
		}
		return false, nil
	}

	// Create the VPC if it does not exist
	log.Info("Creating a new VPC", "name", vpcName)
	newVPC, err := s.createVPC(vpcName)
	if err != nil {
		return false, fmt.Errorf("failed to create VPC: %w", err)
	}

	log.Info("Created VPC", "vpcID", *newVPC.ID)

	// Record newly provisioned VPC into Status
	s.IBMPowerVSCluster.Status.VPC = infrav1.VPCStatus{
		ID:     *newVPC.ID,
		Name:   *newVPC.Name,
		Region: s.IBMPowerVSCluster.Spec.VPC.Region,
	}

	return true, nil // Requeue to hit Step 1 and await active status
}

// createVPC provisions a new VPC in IBM Cloud.
func (s *ClusterScope) createVPC(name string) (*vpcv1.VPC, error) {
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroupName())
	}

	addressPrefixManagement := "auto"
	vpcOption := &vpcv1.CreateVPCOptions{
		ResourceGroup:           &vpcv1.ResourceGroupIdentity{ID: ptr.To(resourceGroupID)},
		Name:                    ptr.To(name),
		AddressPrefixManagement: ptr.To(addressPrefixManagement),
	}

	vpcDetails, _, err := s.IBMVPCClient.CreateVPC(vpcOption)
	if err != nil {
		return nil, err
	}

	// set security group rule for vpc
	if vpcDetails.DefaultSecurityGroup != nil && vpcDetails.DefaultSecurityGroup.ID != nil {
		options := &vpcv1.CreateSecurityGroupRuleOptions{}
		options.SetSecurityGroupID(*vpcDetails.DefaultSecurityGroup.ID)
		options.SetSecurityGroupRulePrototype(&vpcv1.SecurityGroupRulePrototype{
			Direction: core.StringPtr("inbound"),
			Protocol:  core.StringPtr("tcp"),
			IPVersion: core.StringPtr("ipv4"),
			PortMin:   core.Int64Ptr(int64(s.APIServerPort())),
			PortMax:   core.Int64Ptr(int64(s.APIServerPort())),
		})
		if _, _, err = s.IBMVPCClient.CreateSecurityGroupRule(options); err != nil {
			return nil, fmt.Errorf("error creating security group rule for VPC: %w", err)
		}
	}

	return vpcDetails, nil
}

// ReconcileVPCSubnets evaluates the user's intent and reconciles all IBM Cloud VPC subnets.
//
//nolint:gocyclo // complexity is acceptable for reconciliation logic
func (s *ClusterScope) ReconcileVPCSubnets(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// If Topology is VirtualIP, subnets are skipped entirely.
	if cluster.Spec.Topology == infrav1.PowerVSVirtualIPTopology {
		return false, nil
	}

	// 1. Gather all available VPC Zones for the region to handle dynamic expansion
	vpcRegion := cluster.Spec.VPC.Region
	if vpcRegion == "" {
		return false, fmt.Errorf("cannot reconcile VPC subnets: VPC region is missing from spec")
	}
	vpcZones, err := regionUtil.VPCZonesForVPCRegion(vpcRegion)
	if err != nil {
		return false, fmt.Errorf("error fetching VPC zones associated with VPC region %s: %w", vpcRegion, err)
	}
	if len(vpcZones) == 0 {
		return false, fmt.Errorf("failed to fetch VPC zones, no zones found for region %s", vpcRegion)
	}

	// 2. Expand desired subnets (Use spec array, or generate default ones if empty)
	var desiredSubnets []infrav1.VPCSubnetSource
	if len(cluster.Spec.VPCSubnets) == 0 {
		log.V(3).Info("VPC subnets not configured in spec, auto-expanding default subnets for all zones", "region", vpcRegion)
		for _, zone := range vpcZones {
			desiredSubnets = append(desiredSubnets, infrav1.VPCSubnetSource{
				Type: infrav1.SourceTypeProvision,
				Zone: zone, // Set at top-level
				Provision: infrav1.VPCSubnetProvision{
					Name: ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeSubnet, zone),
				},
			})
		}
	} else {
		desiredSubnets = cluster.Spec.VPCSubnets
	}

	// Track whether any subnets require a requeue to catch up to an active state
	anyPending := false

	// 3. Process each subnet using our SSA-compliant list matching
	for i, subnetSpec := range desiredSubnets {
		// Determine the calculated or explicit Name for this subnet configuration
		var subnetName string
		if subnetSpec.Type == infrav1.SourceTypeReference {
			subnetName = subnetSpec.Reference.Name
		} else {
			subnetName = subnetSpec.Provision.Name
			if subnetName == "" {
				subnetName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeSubnet, strconv.Itoa(i))
			}
		}

		log.Info("Reconciling VPC subnet entry", "name", subnetName, "type", subnetSpec.Type)

		// 4. Idempotency Check: See if this specific subnet is already populated in Status.VPCSubnets list
		var existingStatus *infrav1.VPCSubnetStatus // Updated to the new specific status type
		for j := range cluster.Status.VPCSubnets {
			if cluster.Status.VPCSubnets[j].Name == subnetName {
				existingStatus = &cluster.Status.VPCSubnets[j]
				break
			}
		}

		// If it's already resolved in status, verify it still exists in the cloud
		if existingStatus != nil && existingStatus.ID != "" {
			log.V(3).Info("VPC subnet ID already tracked in status, verifying runtime presence", "name", subnetName, "id", existingStatus.ID)
			subnetDetails, _, err := s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{ID: ptr.To(existingStatus.ID)})
			if err != nil {
				return false, fmt.Errorf("error verifying active VPC subnet (id: %s): %w", existingStatus.ID, err)
			}
			if subnetDetails == nil {
				return false, fmt.Errorf("tracked VPC subnet (id: %s) was not found in IBM Cloud", existingStatus.ID)
			}

			// Subnet is healthy and verified; carry on to the next list entry
			continue
		}

		// 5. Route to intent-based sub-handlers if status is empty
		var subnetDetails *vpcv1.Subnet
		switch subnetSpec.Type {
		case infrav1.SourceTypeReference:
			subnetDetails, err = s.reconcileSubnetReference(subnetSpec.Reference)
		case infrav1.SourceTypeProvision:
			// Zone is optional in a user-supplied VPCSubnets entry; fall back to a
			// round-robin zone assignment across the region's available zones.
			if subnetSpec.Zone == "" {
				subnetSpec.Zone = vpcZones[i%len(vpcZones)]
			}
			// Pass the resolved zone explicitly to the provisioner
			subnetDetails, err = s.reconcileSubnetProvision(subnetName, subnetSpec.Zone)
			anyPending = true // Newly provisioned resources warrant a requeue loop
		default:
			return false, fmt.Errorf("unknown VPC subnet source type: %s", subnetSpec.Type)
		}

		if err != nil {
			return false, fmt.Errorf("failed resolving subnet %s: %w", subnetName, err)
		}

		// Skip status update if subnet details are not available
		if subnetDetails == nil || subnetDetails.ID == nil || subnetDetails.Name == nil {
			continue
		}

		// Dynamically extract the actual Zone from the cloud payload if available, fallback to spec
		actualZone := subnetSpec.Zone
		if subnetDetails.Zone != nil && subnetDetails.Zone.Name != nil {
			actualZone = *subnetDetails.Zone.Name
		}

		// 6. Update Status array using the new struct
		s.updateSubnetStatusList(infrav1.VPCSubnetStatus{
			ID:   *subnetDetails.ID,
			Name: *subnetDetails.Name,
			Zone: actualZone,
		})
	}

	return anyPending, nil
}

// reconcileSubnetReference verifies a referenced external VPC subnet.
func (s *ClusterScope) reconcileSubnetReference(ref infrav1.ResourceIdentifier) (*vpcv1.Subnet, error) {
	var subnetDetails *vpcv1.Subnet
	var err error

	if ref.ID != "" {
		subnetDetails, _, err = s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{ID: ptr.To(ref.ID)})
	} else if ref.Name != "" {
		subnetDetails, err = s.IBMVPCClient.GetVPCSubnetByName(ref.Name)
	} else {
		return nil, fmt.Errorf("subnet reference configuration must have either ID or Name defined")
	}

	if err != nil {
		return nil, fmt.Errorf("failed fetching referenced subnet: %w", err)
	}
	if subnetDetails == nil || subnetDetails.ID == nil || subnetDetails.Name == nil {
		return nil, fmt.Errorf("referenced subnet could not be located or returned partial results")
	}

	return subnetDetails, nil
}

// reconcileSubnetProvision handles lookup by name or dynamically creating a new VPC subnet.
func (s *ClusterScope) reconcileSubnetProvision(name string, zone string) (*vpcv1.Subnet, error) {
	// Look up by name to enforce idempotency
	subnetDetails, err := s.IBMVPCClient.GetVPCSubnetByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed checking subnet presence by name: %w", err)
	}

	if subnetDetails != nil && subnetDetails.ID != nil && subnetDetails.Name != nil {
		return subnetDetails, nil
	}

	// Create a new subnet if it does not exist
	return s.createVPCSubnet(name, zone)
}

// createVPCSubnet provisions the raw subnet resource within the Cloud VPC.
func (s *ClusterScope) createVPCSubnet(name string, zone string) (*vpcv1.Subnet, error) {
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("resource group ID is empty")
	}

	vpcID := s.IBMPowerVSCluster.Status.VPC.ID
	if vpcID == "" {
		return nil, fmt.Errorf("cannot create subnet: managing VPC ID is not found in status")
	}

	ipVersion := vpcSubnetIPVersion4
	options := &vpcv1.CreateSubnetOptions{}
	options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
		IPVersion:             &ipVersion,
		TotalIpv4AddressCount: ptr.To(vpcSubnetIPAddressCount),
		Name:                  ptr.To(name),
		VPC: &vpcv1.VPCIdentity{
			ID: ptr.To(vpcID),
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: ptr.To(zone),
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroupID,
		},
	})

	subnetDetails, _, err := s.IBMVPCClient.CreateSubnet(options)
	if err != nil {
		return nil, fmt.Errorf("error invoking VPC CreateSubnet API: %w", err)
	}
	if subnetDetails == nil {
		return nil, fmt.Errorf("created VPC subnet details returned nil response")
	}

	return subnetDetails, nil
}

// updateSubnetStatusList handles adding or updating an item within the SSA associative list.
func (s *ClusterScope) updateSubnetStatusList(status infrav1.VPCSubnetStatus) {
	for i, existing := range s.IBMPowerVSCluster.Status.VPCSubnets {
		if existing.Name == status.Name {
			s.IBMPowerVSCluster.Status.VPCSubnets[i] = status
			return
		}
	}
	s.IBMPowerVSCluster.Status.VPCSubnets = append(s.IBMPowerVSCluster.Status.VPCSubnets, status)
}

// ReconcileLoadBalancers reconcile loadBalancer.
//
//nolint:gocyclo // complexity is acceptable for reconciliation logic
func (s *ClusterScope) ReconcileLoadBalancers(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	// Setup defaulting if no load balancers are specified in the spec
	loadBalancers := make([]infrav1.LoadBalancerSource, 0)
	if len(s.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		log.V(3).Info("VPC load balancer is not set, constructing a default provisioned public load balancer")
		loadBalancers = append(loadBalancers, infrav1.LoadBalancerSource{
			Type: infrav1.SourceTypeProvision,
			Provision: infrav1.LoadBalancerProvision{
				Name: ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeLBPublic, ""),
				Type: infrav1.LoadBalancerTypePublic,
			},
		})
	} else {
		loadBalancers = append(loadBalancers, s.IBMPowerVSCluster.Spec.LoadBalancers...)
	}

	isAnyLoadBalancerNotReady := false

	for i, lbSource := range loadBalancers {
		if lbSource.Type == infrav1.SourceTypeReference {
			var loadBalancer *vpcv1.LoadBalancer
			var err error

			if lbSource.Reference.ID != "" {
				loadBalancer, _, err = s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{ID: ptr.To(lbSource.Reference.ID)})
			} else if lbSource.Reference.Name != "" {
				loadBalancer, err = s.IBMVPCClient.GetLoadBalancerByName(lbSource.Reference.Name)
			} else {
				return false, fmt.Errorf("referenced load balancer must have either an ID or Name")
			}

			if err != nil || loadBalancer == nil {
				return false, fmt.Errorf("failed to fetch referenced load balancer details: %w", err)
			}

			lbName := *loadBalancer.Name
			if isReady := s.checkLoadBalancerState(ctx, *loadBalancer); !isReady {
				log.V(3).Info("Referenced LoadBalancer is still not Active", "loadBalancerName", lbName)
				isAnyLoadBalancerNotReady = true
			}

			s.SetLoadBalancerStatus(ctx, lbName, infrav1.LoadBalancerStatus{
				Name:     lbName,
				ID:       *loadBalancer.ID,
				State:    infrav1.LoadBalancerState(*loadBalancer.ProvisioningStatus),
				Hostname: ptr.Deref(loadBalancer.Hostname, ""),
			})
			continue
		}

		// Provision the load balancer
		provision := lbSource.Provision
		lbName := provision.Name
		if lbName == "" {
			lbType := ResourceTypeLBPublic
			if provision.Type == infrav1.LoadBalancerTypePrivate {
				lbType = ResourceTypeLBPrivate
			}
			qualifier := ""
			if i > 0 {
				qualifier = strconv.Itoa(i)
			}
			lbName = ResourceName(s.IBMPowerVSCluster.Name, lbType, qualifier)
		}

		// Check if we already have the ID cached in status
		lbID := s.GetLoadBalancerID(lbName)
		if lbID != "" {
			log.V(3).Info("Load balancer ID is set, fetching load balancer details", "loadBalancerID", lbID)
			loadBalancer, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: ptr.To(lbID),
			})
			if err != nil {
				return false, fmt.Errorf("failed to fetch load balancer details: %w", err)
			}

			if loadBalancer == nil {
				return false, fmt.Errorf("failed to fetch load balancer details: received empty/nil response from cloud API for ID %s", lbID)
			}

			if isReady := s.checkLoadBalancerState(ctx, *loadBalancer); !isReady {
				log.V(3).Info("LoadBalancer is still not Active", "loadBalancerName", lbName, "state", *loadBalancer.ProvisioningStatus)
				isAnyLoadBalancerNotReady = true
			}

			s.SetLoadBalancerStatus(ctx, lbName, infrav1.LoadBalancerStatus{
				Name:     lbName,
				ID:       *loadBalancer.ID,
				State:    infrav1.LoadBalancerState(*loadBalancer.ProvisioningStatus),
				Hostname: ptr.Deref(loadBalancer.Hostname, ""),
			})
			continue
		}

		// Check if load balancer already exists in cloud by name
		lbStatus, err := s.checkLoadBalancer(ctx, lbName)
		if err != nil {
			return false, fmt.Errorf("failed to check if load balancer exists: %w", err)
		}
		if lbStatus != nil {
			log.V(3).Info("Found load balancer in cloud", "loadBalancerID", lbStatus.ID)
			s.SetLoadBalancerStatus(ctx, lbName, *lbStatus)
			continue
		}

		// Pre-flight check on ports
		if err := s.checkLoadBalancerPort(lbName, provision); err != nil {
			return false, fmt.Errorf("failed to check load balancer port: %w", err)
		}

		// Create loadBalancer
		log.Info("Creating load balancer", "name", lbName)
		lbStatus, err = s.createLoadBalancer(ctx, lbName, provision)
		if err != nil {
			return false, fmt.Errorf("failed to create load balancer: %w", err)
		}
		log.Info("Created load balancer", "loadBalancerID", lbStatus.ID)
		s.SetLoadBalancerStatus(ctx, lbName, *lbStatus)
		isAnyLoadBalancerNotReady = true
	}

	if isAnyLoadBalancerNotReady {
		return false, nil
	}
	return true, nil
}

// SetLoadBalancerStatus updates or appends the load balancer status.
func (s *ClusterScope) SetLoadBalancerStatus(ctx context.Context, name string, status infrav1.LoadBalancerStatus) {
	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("Setting load balancer status", "name", name, "status", status)

	// Ensure the slice is initialized
	if s.IBMPowerVSCluster.Status.LoadBalancers == nil {
		s.IBMPowerVSCluster.Status.LoadBalancers = make([]infrav1.LoadBalancerStatus, 0)
	}

	// Check if the entry already exists in the slice to update it in place
	for i, current := range s.IBMPowerVSCluster.Status.LoadBalancers {
		if current.Name == name {
			s.IBMPowerVSCluster.Status.LoadBalancers[i] = status
			return
		}
	}

	s.IBMPowerVSCluster.Status.LoadBalancers = append(s.IBMPowerVSCluster.Status.LoadBalancers, status)
}

func (s *ClusterScope) checkLoadBalancerPort(lbName string, prov infrav1.LoadBalancerProvision) error {
	for _, listener := range prov.AdditionalListeners {
		if listener.Port == int64(s.APIServerPort()) {
			return fmt.Errorf("port %d for the %s load balancer cannot be used as an additional listener port, as it is already assigned to the API server", listener.Port, lbName)
		}
	}
	return nil
}

// checkLoadBalancer checks if VPC load balancer by the given name exists in cloud.
func (s *ClusterScope) checkLoadBalancer(ctx context.Context, name string) (*infrav1.LoadBalancerStatus, error) {
	log := ctrl.LoggerFrom(ctx)
	loadBalancer, err := s.IBMVPCClient.GetLoadBalancerByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch load balancer details: %w", err)
	}
	if loadBalancer == nil {
		log.V(3).Info("VPC load balancer not found in cloud")
		return nil, nil
	}
	return &infrav1.LoadBalancerStatus{
		Name:     name,
		ID:       *loadBalancer.ID,
		State:    infrav1.LoadBalancerState(*loadBalancer.ProvisioningStatus),
		Hostname: ptr.Deref(loadBalancer.Hostname, ""),
	}, nil
}

// createLoadBalancer creates loadBalancer.
func (s *ClusterScope) createLoadBalancer(ctx context.Context, lbName string, prov infrav1.LoadBalancerProvision) (*infrav1.LoadBalancerStatus, error) {
	log := ctrl.LoggerFrom(ctx)
	options := &vpcv1.CreateLoadBalancerOptions{}

	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroupName())
	}

	// Handle the new Enum for visibility
	isPublic := true // Defaults to Public
	if prov.Type == infrav1.LoadBalancerTypePrivate {
		isPublic = false
	}

	options.SetIsPublic(isPublic)
	options.SetName(lbName)
	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &resourceGroupID,
	})

	if len(s.IBMPowerVSCluster.Status.VPCSubnets) == 0 {
		return nil, fmt.Errorf("no VPC subnets are present in cluster status for load balancer creation")
	}

	for _, subnet := range s.IBMPowerVSCluster.Status.VPCSubnets {
		if subnet.ID == "" {
			continue
		}
		subnet := &vpcv1.SubnetIdentity{
			ID: ptr.To(subnet.ID),
		}
		options.Subnets = append(options.Subnets, subnet)
	}

	options.SetPools([]vpcv1.LoadBalancerPoolPrototypeLoadBalancerContext{
		{
			Algorithm:     core.StringPtr("round_robin"),
			HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{Delay: core.Int64Ptr(5), MaxRetries: core.Int64Ptr(2), Timeout: core.Int64Ptr(2), Type: core.StringPtr("tcp")},
			Name:          core.StringPtr(fmt.Sprintf("%s-pool-%d", lbName, s.APIServerPort())),
			Protocol:      core.StringPtr("tcp"),
		},
	})

	options.SetListeners([]vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{
		{
			Protocol: core.StringPtr("tcp"),
			Port:     core.Int64Ptr(int64(s.APIServerPort())),
			DefaultPool: &vpcv1.LoadBalancerPoolIdentityByName{
				Name: core.StringPtr(fmt.Sprintf("%s-pool-%d", lbName, s.APIServerPort())),
			},
		},
	})

	for _, additionalListener := range prov.AdditionalListeners {
		pool := vpcv1.LoadBalancerPoolPrototypeLoadBalancerContext{
			Algorithm:     core.StringPtr("round_robin"),
			HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{Delay: core.Int64Ptr(5), MaxRetries: core.Int64Ptr(2), Timeout: core.Int64Ptr(2), Type: core.StringPtr("tcp")},
			Name:          ptr.To(fmt.Sprintf("additional-pool-%d", additionalListener.Port)),
			Protocol:      core.StringPtr("tcp"),
		}
		options.Pools = append(options.Pools, pool)

		listener := vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{
			Protocol: core.StringPtr("tcp"),
			Port:     core.Int64Ptr(additionalListener.Port),
			DefaultPool: &vpcv1.LoadBalancerPoolIdentityByName{
				Name: ptr.To(fmt.Sprintf("additional-pool-%d", additionalListener.Port)),
			},
		}
		options.Listeners = append(options.Listeners, listener)
	}

	log.V(5).Info("Creating load balancer", "options", options)
	loadBalancer, _, err := s.IBMVPCClient.CreateLoadBalancer(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	lbState := infrav1.LoadBalancerState(*loadBalancer.ProvisioningStatus)
	return &infrav1.LoadBalancerStatus{
		Name:     lbName,
		ID:       *loadBalancer.ID,
		State:    lbState,
		Hostname: ptr.Deref(loadBalancer.Hostname, ""),
	}, nil
}

// checkLoadBalancerState checks the state of a VPC load balancer.
// If state is active, true is returned, in all other cases, it returns false indicating that load balancer is still not ready.
func (s *ClusterScope) checkLoadBalancerState(ctx context.Context, lb vpcv1.LoadBalancer) bool {
	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("Checking the status of VPC load balancer", "loadBalancerName", *lb.Name)
	switch *lb.ProvisioningStatus {
	case string(infrav1.LoadBalancerStateActive):
		log.V(3).Info("Load balancer is in active state")
		return true
	case string(infrav1.LoadBalancerStateCreatePending):
		log.V(3).Info("Load balancer creation is in pending state")
	case string(infrav1.LoadBalancerStateUpdatePending):
		log.V(3).Info("Load balancer is in updating state")
	}
	return false
}

// ReconcileVPCSecurityGroups evaluates user intent and reconciles all VPC security groups.
func (s *ClusterScope) ReconcileVPCSecurityGroups(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	var updatedStatus []infrav1.VPCSecurityGroupStatus

	for _, sgSource := range s.IBMPowerVSCluster.Spec.VPCSecurityGroups {
		var sgStatus *infrav1.VPCSecurityGroupStatus
		var err error

		switch sgSource.Type {
		case infrav1.SourceTypeReference:
			log.Info("Reconciling referenced VPC security group")
			sgStatus, err = s.reconcileVPCSecurityGroupReference(ctx, sgSource.Reference)
		case infrav1.SourceTypeProvision:
			log.Info("Reconciling managed VPC security group", "name", sgSource.Provision.Name)
			sgStatus, err = s.reconcileVPCSecurityGroupProvision(ctx, sgSource.Provision)
		default:
			err = fmt.Errorf("unknown security group source type: %s", sgSource.Type)
		}

		if err != nil {
			return fmt.Errorf("failed to reconcile security group: %w", err)
		}

		if sgStatus != nil {
			updatedStatus = append(updatedStatus, *sgStatus)
		}
	}

	// Overwrite the status completely with the freshly validated state
	s.IBMPowerVSCluster.Status.VPCSecurityGroups = updatedStatus
	return nil
}

// reconcileVPCSecurityGroupReference verifies an existing SG exists and returns its status.
func (s *ClusterScope) reconcileVPCSecurityGroupReference(_ context.Context, ref infrav1.ResourceIdentifier) (*infrav1.VPCSecurityGroupStatus, error) {
	var sg *vpcv1.SecurityGroup
	var err error

	if ref.ID != "" {
		sg, _, err = s.IBMVPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{ID: ptr.To(ref.ID)})
	} else if ref.Name != "" {
		sg, err = s.IBMVPCClient.GetSecurityGroupByName(ref.Name)
	} else {
		return nil, fmt.Errorf("referenced security group must have either ID or Name specified")
	}

	if err != nil || sg == nil {
		return nil, fmt.Errorf("failed to find referenced VPC security group: %w", err)
	}

	return &infrav1.VPCSecurityGroupStatus{
		ID:   *sg.ID,
		Name: *sg.Name,
		// Note: We don't track rules for referenced SGs because CAPI shouldn't manage them.
	}, nil
}

// reconcileVPCSecurityGroupProvision creates or updates a managed SG and its rules.
func (s *ClusterScope) reconcileVPCSecurityGroupProvision(ctx context.Context, prov infrav1.VPCSecurityGroupProvision) (*infrav1.VPCSecurityGroupStatus, error) {
	targetName := prov.Name
	if targetName == "" {
		targetName = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeSG, "")
	}

	// 1. Ensure the Security Group exists
	sg, err := s.IBMVPCClient.GetSecurityGroupByName(targetName)
	if err != nil {
		// Ignore not found errors, return actual API errors
		if _, ok := err.(*vpc.SecurityGroupByNameNotFound); !ok {
			return nil, fmt.Errorf("failed to query VPC security group by name '%s': %w", targetName, err)
		}
	}

	if sg == nil {
		sgID, err := s.createVPCSecurityGroup(ctx, targetName)
		if err != nil {
			return nil, err
		}
		// Fetch the newly created SG object
		sg, _, err = s.IBMVPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{ID: sgID})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch newly created security group '%s': %w", targetName, err)
		}
	}

	// 2. Reconcile Rules
	ruleIDs, err := s.createVPCSecurityGroupRules(ctx, prov.Rules, *sg.ID, sg.Rules)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC security group rules for '%s': %w", targetName, err)
	}

	// 3. Map rule string IDs to our new Status struct array
	var ruleStatus []infrav1.VPCSecurityGroupRuleStatus
	for _, id := range ruleIDs {
		ruleStatus = append(ruleStatus, infrav1.VPCSecurityGroupRuleStatus{ID: id})
	}

	return &infrav1.VPCSecurityGroupStatus{
		ID:    *sg.ID,
		Name:  *sg.Name,
		Rules: ruleStatus,
	}, nil
}

// createVPCSecurityGroup creates a basic VPC security group without rules.
func (s *ClusterScope) createVPCSecurityGroup(ctx context.Context, name string) (*string, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Creating new VPC security group", "name", name)

	options := &vpcv1.CreateSecurityGroupOptions{
		VPC: &vpcv1.VPCIdentity{
			ID: &s.IBMPowerVSCluster.Status.VPC.ID,
		},
		Name: ptr.To(name),
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: ptr.To(s.GetResourceGroupID()),
		},
	}

	securityGroup, _, err := s.IBMVPCClient.CreateSecurityGroup(options)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC security group: %w", err)
	}
	return securityGroup.ID, nil
}

// getVPCSecurityGroupRuleID returns the rule ID from any concrete SG rule type, or "" if unrecognized.
func getVPCSecurityGroupRuleID(ruleIntf vpcv1.SecurityGroupRuleIntf) string {
	switch r := ruleIntf.(type) {
	case *vpcv1.SecurityGroupRule:
		if r != nil && r.ID != nil {
			return *r.ID
		}
	case *vpcv1.SecurityGroupRuleProtocolAny:
		if r != nil && r.ID != nil {
			return *r.ID
		}
	case *vpcv1.SecurityGroupRuleProtocolIcmptcpudp:
		if r != nil && r.ID != nil {
			return *r.ID
		}
	case *vpcv1.SecurityGroupRuleProtocolIndividual:
		if r != nil && r.ID != nil {
			return *r.ID
		}
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp:
		if r != nil && r.ID != nil {
			return *r.ID
		}
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp:
		if r != nil && r.ID != nil {
			return *r.ID
		}
	}
	return ""
}

// validateVPCSecurityGroupRuleKey builds an sgRuleKey from an existing SG rule.
func validateVPCSecurityGroupRuleKey(ctx context.Context, ruleIntf vpcv1.SecurityGroupRuleIntf) (sgRuleKey, bool) {
	if ruleIntf == nil {
		return sgRuleKey{}, false
	}
	var key sgRuleKey
	switch r := ruleIntf.(type) {
	case *vpcv1.SecurityGroupRule:
		key = buildSGRuleKey(r.Direction, r.Protocol, r.PortMin, r.PortMax, r.Remote)
	case *vpcv1.SecurityGroupRuleProtocolAny:
		key = buildSGRuleKey(r.Direction, r.Protocol, nil, nil, r.Remote)
	case *vpcv1.SecurityGroupRuleProtocolIcmptcpudp:
		key = buildSGRuleKey(r.Direction, r.Protocol, nil, nil, r.Remote)
	case *vpcv1.SecurityGroupRuleProtocolIndividual:
		key = buildSGRuleKey(r.Direction, r.Protocol, nil, nil, r.Remote)
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp:
		key = buildSGRuleKey(r.Direction, r.Protocol, nil, nil, r.Remote)
	case *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp:
		key = buildSGRuleKey(r.Direction, r.Protocol, r.PortMin, r.PortMax, r.Remote)
	default:
		ctrl.LoggerFrom(ctx).V(4).Info("Skipping unrecognized security group rule type during deduplication", "type", fmt.Sprintf("%T", ruleIntf))
		return sgRuleKey{}, false
	}
	return key, true
}

// buildSGRuleKey assembles an sgRuleKey from the fields common to all rule types.
func buildSGRuleKey(direction, protocol *string, portMin, portMax *int64, remoteIntf vpcv1.SecurityGroupRuleRemoteIntf) sgRuleKey {
	var key sgRuleKey
	if direction != nil {
		key.direction = *direction
	}
	if protocol != nil {
		key.protocol = *protocol
	}
	if portMin != nil {
		key.portMin = *portMin
	}
	if portMax != nil {
		key.portMax = *portMax
	}
	if remote, ok := remoteIntf.(*vpcv1.SecurityGroupRuleRemote); ok && remote != nil {
		if remote.CIDRBlock != nil {
			key.cidrBlock = *remote.CIDRBlock
		}
		if remote.Address != nil {
			key.address = *remote.Address
		}
		if remote.CRN != nil {
			key.crn = *remote.CRN
		}
	}
	return key
}

// createVPCSecurityGroupRules creates the given rules on the security group.
// Rules already present on the SG are skipped.
func (s *ClusterScope) createVPCSecurityGroupRules(ctx context.Context, rules []infrav1.VPCSecurityGroupRule, securityGroupID string, existingRules []vpcv1.SecurityGroupRuleIntf) ([]string, error) {
	log := ctrl.LoggerFrom(ctx)
	var ruleIDs []string

	log.V(3).Info("Reconciling VPC security group rules", "securityGroupID", securityGroupID, "specRuleCount", len(rules), "existingRuleCount", len(existingRules))

	// Index existing rules by key so we can skip ones already present.
	existingByKey := make(map[sgRuleKey]string, len(existingRules))
	for _, rule := range existingRules {
		key, ok := validateVPCSecurityGroupRuleKey(ctx, rule)
		if !ok {
			continue
		}
		if id := getVPCSecurityGroupRuleID(rule); id != "" {
			existingByKey[key] = id
		}
	}

	for _, rule := range rules {
		// Work on a copy with normalized prototypes so the deprecated 'all' protocol is
		// translated without mutating the cluster spec.
		if src := normalizedVPCSecurityGroupRulePrototype(&rule.Source); src != nil {
			rule.Source = *src
		}
		if dst := normalizedVPCSecurityGroupRulePrototype(&rule.Destination); dst != nil {
			rule.Destination = *dst
		}

		direction := string(rule.Direction)

		var protocol string
		var portMin, portMax int64
		var remotes []infrav1.VPCSecurityGroupRuleRemote

		// Route the extraction based on the direction of the rule
		switch rule.Direction {
		case infrav1.VPCSecurityGroupRuleDirectionInbound:
			protocol = string(rule.Source.Protocol)
			portMin = rule.Source.PortRange.MinimumPort
			portMax = rule.Source.PortRange.MaximumPort
			remotes = rule.Source.Remotes

		case infrav1.VPCSecurityGroupRuleDirectionOutbound:
			protocol = string(rule.Destination.Protocol)
			portMin = rule.Destination.PortRange.MinimumPort
			portMax = rule.Destination.PortRange.MaximumPort
			remotes = rule.Destination.Remotes

		default:
			// Kubernetes CEL validation should prevent this, but it is good practice to catch it.
			return nil, fmt.Errorf("invalid rule direction provided: %s", direction)
		}

		if protocol == "" {
			return nil, fmt.Errorf("security group rule with direction %q has an empty protocol; source/destination must be fully specified", direction)
		}
		if len(remotes) == 0 {
			return nil, fmt.Errorf("security group rule with direction %q has no remotes; at least one remote must be specified", direction)
		}

		// Create a distinct rule in IBM Cloud for every remote target specified
		for _, remote := range remotes {
			id, err := s.createVPCSecurityGroupRule(ctx, securityGroupID, direction, protocol, portMin, portMax, remote, existingByKey)
			if err != nil {
				return nil, fmt.Errorf("failed to create %s VPC security group rule: %w", direction, err)
			}
			ruleIDs = append(ruleIDs, id)
		}
	}

	return ruleIDs, nil
}

// createVPCSecurityGroupRule creates a single rule. If a matching rule already exists in
// existingByKey, it returns that ID without calling the API.
func (s *ClusterScope) createVPCSecurityGroupRule(ctx context.Context, securityGroupID string, direction string, protocol string, portMin, portMax int64, remote infrav1.VPCSecurityGroupRuleRemote, existingByKey map[sgRuleKey]string) (string, error) {
	log := ctrl.LoggerFrom(ctx)

	// 1. Resolve Remote Target
	remoteOption, err := s.resolveRuleRemote(remote)
	if err != nil {
		return "", err
	}

	// 2. Check whether this rule already exists.
	key := sgRuleKey{
		direction: direction,
		protocol:  protocol,
	}
	if (protocol == string(infrav1.VPCSecurityGroupRuleProtocolTCP) || protocol == string(infrav1.VPCSecurityGroupRuleProtocolUDP)) && portMin > 0 {
		key.portMin = portMin
		key.portMax = portMax
	}
	if remoteOption.CIDRBlock != nil {
		key.cidrBlock = *remoteOption.CIDRBlock
	}
	if remoteOption.Address != nil {
		key.address = *remoteOption.Address
	}
	if remoteOption.CRN != nil {
		key.crn = *remoteOption.CRN
	}

	if existingID, found := existingByKey[key]; found {
		log.V(3).Info("Skipping duplicate VPC security group rule", "securityGroupID", securityGroupID, "existingRuleID", existingID)
		return existingID, nil
	}

	// 3. Build Protocol Prototype (Injecting Pointers here for SDK)
	prototype := &vpcv1.SecurityGroupRulePrototype{
		Direction: ptr.To(direction),
		Protocol:  ptr.To(protocol),
		Remote:    remoteOption,
	}

	// Only attach port ranges if it's TCP/UDP and ports were explicitly provided
	if (protocol == string(infrav1.VPCSecurityGroupRuleProtocolTCP) || protocol == string(infrav1.VPCSecurityGroupRuleProtocolUDP)) && portMin > 0 {
		prototype.PortMin = ptr.To(portMin)
		prototype.PortMax = ptr.To(portMax)
	}

	options := &vpcv1.CreateSecurityGroupRuleOptions{
		SecurityGroupID:            ptr.To(securityGroupID),
		SecurityGroupRulePrototype: prototype,
	}

	log.V(3).Info("Creating VPC security group rule", "securityGroupID", securityGroupID, "direction", direction, "protocol", protocol)

	// 4. Execute API Call
	ruleIntf, _, err := s.IBMVPCClient.CreateSecurityGroupRule(options)
	if err != nil {
		return "", fmt.Errorf("failed to execute CreateSecurityGroupRule API: %w", err)
	}

	// 5. Extract Rule ID based on returned interface type
	ruleID := getVPCSecurityGroupRuleID(ruleIntf)
	if ruleID == "" {
		return "", fmt.Errorf("unrecognized rule type returned from API: %T", ruleIntf)
	}

	// Record so identical remotes later in the same spec are also skipped.
	if existingByKey != nil {
		existingByKey[key] = ruleID
	}

	log.Info("Successfully created VPC security group rule", "ruleID", ruleID)
	return ruleID, nil
}

// resolveRuleRemote resolves a CRD remote to the SDK prototype, looking up subnet CIDRs and SG CRNs as needed.
func (s *ClusterScope) resolveRuleRemote(remote infrav1.VPCSecurityGroupRuleRemote) (*vpcv1.SecurityGroupRuleRemotePrototype, error) {
	remoteOption := &vpcv1.SecurityGroupRuleRemotePrototype{}
	switch remote.RemoteType {
	case infrav1.VPCSecurityGroupRuleRemoteTypeCIDR:
		cidrSubnet, err := s.IBMVPCClient.GetVPCSubnetByName(remote.CIDRSubnetName)
		if err != nil || cidrSubnet == nil {
			return nil, fmt.Errorf("failed to find VPC subnet by name '%s': %w", remote.CIDRSubnetName, err)
		}
		remoteOption.CIDRBlock = cidrSubnet.Ipv4CIDRBlock
	case infrav1.VPCSecurityGroupRuleRemoteTypeAddress:
		remoteOption.Address = ptr.To(remote.Address)
	case infrav1.VPCSecurityGroupRuleRemoteTypeSG:
		sg, err := s.IBMVPCClient.GetSecurityGroupByName(remote.SecurityGroupName)
		if err != nil || sg == nil {
			return nil, fmt.Errorf("failed to find VPC security group by name '%s': %w", remote.SecurityGroupName, err)
		}
		remoteOption.CRN = sg.CRN
	default:
		// Any/0.0.0.0 mapping
		remoteOption.CIDRBlock = ptr.To("0.0.0.0/0")
	}
	return remoteOption, nil
}

// ReconcileCOSInstance evaluates the user's intent and reconciles the COS instance and bucket.
func (s *ClusterScope) ReconcileCOSInstance(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// 1. Opt-out check: If Type is empty, the user omitted the COSInstance block entirely.
	if cluster.Spec.COSInstance.Type == "" {
		return nil
	}

	// 2. Resolve the COS instance to a stable (instanceID, instanceName, instanceCRN).
	instanceID, instanceName, instanceCRN, err := s.resolveCOSInstance(ctx)
	if err != nil {
		return err
	}

	// Record the resolved instance identity.
	cluster.Status.COSInstance.ID = instanceID
	cluster.Status.COSInstance.Name = instanceName

	// 3. Resolve the bucket name and region, falling back through spec then VPC region.
	targetBucketName, targetBucketRegion, err := s.resolveCOSBucketParams()
	if err != nil {
		return err
	}

	log.V(3).Info("ReconcileCOSInstance: resolved params",
		"instanceID", instanceID,
		"instanceCRN", instanceCRN,
		"bucketName", targetBucketName,
		"bucketRegion", targetBucketRegion,
	)

	// 4. Setup the COS Client now that we have a guaranteed active instance ID.
	if err := s.setupCOSClient(ctx, instanceID, targetBucketRegion); err != nil {
		return fmt.Errorf("failed to configure COS client: %w", err)
	}

	// 5. Reconcile the Bucket.
	if err := s.reconcileCOSBucket(ctx, targetBucketName); err != nil {
		return fmt.Errorf("failed to reconcile COS bucket: %w", err)
	}

	cluster.Status.COSInstance.BucketName = targetBucketName
	cluster.Status.COSInstance.BucketRegion = targetBucketRegion

	// 6. Reconcile the HMAC key Secret for secure ignition pre-signed URL support.
	if err := s.reconcileCOSHMACKey(ctx, instanceCRN); err != nil {
		return fmt.Errorf("failed to reconcile COS HMAC key: %w", err)
	}

	return nil
}

// resolveCOSInstance returns the (instanceID, instanceName, instanceCRN) for the COS
// instance declared in the cluster spec. It takes two paths:
//   - Fast path: if the ID is already in Status, verify liveness in the cloud.
//   - Slow path: resolve via spec (Reference or Provision) and verify the instance is active.
func (s *ClusterScope) resolveCOSInstance(ctx context.Context) (instanceID, instanceName, instanceCRN string, err error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	if cluster.Status.COSInstance.ID != "" {
		return s.verifyCOSInstanceByID(ctx, cluster.Status.COSInstance.ID)
	}

	log.Info("Resolving IBM Cloud COS Instance", "type", cluster.Spec.COSInstance.Type)
	return s.resolveCOSInstanceFromSpec(ctx)
}

// verifyCOSInstanceByID fetches the COS instance by its known ID from the cloud and
// confirms it is active. Returns the (id, name, crn) triple ready for use.
func (s *ClusterScope) verifyCOSInstanceByID(ctx context.Context, id string) (instanceID, instanceName, instanceCRN string, err error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(3).Info("COS Instance ID is set, verifying presence in cloud", "id", id)

	instance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &id,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch COS instance (id: %s) details: %w", id, err)
	}
	if instance == nil {
		return "", "", "", fmt.Errorf("COS instance not found in cloud with ID: %s", id)
	}
	if instance.State == nil {
		return "", "", "", fmt.Errorf("COS instance state is nil for ID: %s", id)
	}
	if *instance.State != string(infrav1.WorkspaceStateActive) {
		return "", "", "", fmt.Errorf("COS instance is not active, current state: %s", *instance.State)
	}

	if instance.CRN == nil || *instance.CRN == "" {
		return "", "", "", fmt.Errorf("COS instance (id: %s) has no CRN; cannot create HMAC credential", id)
	}
	if instance.Name == nil {
		return "", "", "", fmt.Errorf("COS instance name is nil for ID: %s", id)
	}
	return id, *instance.Name, *instance.CRN, nil
}

// resolveCOSInstanceFromSpec routes to reconcileCOSReference or reconcileCOSProvision
// based on the spec type, verifies the instance is active, and returns the identity triple.
// If the instance is not yet active, it records the ID in Status so the next reconcile
// can poll via the fast path (verifyCOSInstanceByID) instead of calling the API again.
func (s *ClusterScope) resolveCOSInstanceFromSpec(ctx context.Context) (instanceID, instanceName, instanceCRN string, err error) {
	cluster := s.IBMPowerVSCluster
	cosSpec := cluster.Spec.COSInstance

	var instance *resourcecontrollerv2.ResourceInstance

	switch cosSpec.Type {
	case infrav1.SourceTypeReference:
		instance, err = s.reconcileCOSReference(ctx, cosSpec.Reference)
	case infrav1.SourceTypeProvision:
		name := cosSpec.Provision.Name
		if name == "" {
			name = ResourceName(cluster.Name, ResourceTypeCOS, "")
		}
		instance, err = s.reconcileCOSProvision(ctx, name)
	default:
		return "", "", "", fmt.Errorf("unknown COS instance source type: %s", cosSpec.Type)
	}

	if err != nil {
		return "", "", "", fmt.Errorf("failed to resolve COS instance: %w", err)
	}
	if instance == nil {
		return "", "", "", fmt.Errorf("COS instance resolved to nil")
	}
	if instance.GUID == nil {
		return "", "", "", fmt.Errorf("COS instance GUID is nil")
	}

	if instance.State == nil {
		return "", "", "", fmt.Errorf("COS instance (id: %s) has no state", *instance.GUID)
	}

	// A newly provisioned instance may still be provisioning — record the ID so the
	// next reconcile loop can skip the provision step and just poll for active state.
	if *instance.State != string(infrav1.WorkspaceStateActive) {
		cluster.Status.COSInstance.ID = *instance.GUID
		if instance.Name != nil {
			cluster.Status.COSInstance.Name = *instance.Name
		}
		return "", "", "", fmt.Errorf("COS instance is not yet active, current state: %s", *instance.State)
	}

	if instance.CRN == nil || *instance.CRN == "" {
		return "", "", "", fmt.Errorf("COS instance has no CRN; cannot create HMAC credential")
	}
	if instance.Name == nil {
		return "", "", "", fmt.Errorf("COS instance name is nil")
	}
	return *instance.GUID, *instance.Name, *instance.CRN, nil
}

// resolveCOSBucketParams determines the target bucket name and region by walking a
// priority chain: Status (already confirmed) → Spec → derived from VPC region.
func (s *ClusterScope) resolveCOSBucketParams() (bucketName, bucketRegion string, err error) {
	cluster := s.IBMPowerVSCluster
	cosSpec := cluster.Spec.COSInstance

	bucketName = cluster.Status.COSInstance.BucketName
	if bucketName == "" {
		bucketName = cosSpec.BucketName
		if bucketName == "" {
			bucketName = ResourceName(cluster.Name, ResourceTypeCOSBucket, "")
		}
	}

	bucketRegion = cluster.Status.COSInstance.BucketRegion
	if bucketRegion == "" {
		bucketRegion = cosSpec.BucketRegion
		if bucketRegion == "" {
			bucketRegion = cluster.Status.VPC.Region
			if bucketRegion == "" {
				return "", "", fmt.Errorf("failed to determine COS bucket region: both bucket region and VPC region are unset")
			}
		}
	}

	return bucketName, bucketRegion, nil
}

// reconcileCOSHMACKey ensures a COS HMAC service credential exists for the given
// COS instance and stores the access_key_id / secret_access_key in a Kubernetes
// Secret named "<cluster>-cos-hmac" in the cluster's namespace.
func (s *ClusterScope) reconcileCOSHMACKey(ctx context.Context, instanceCRN string) error {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	secretName := ResourceName(cluster.Name, ResourceTypeCOSHMACKey, "")

	// Idempotency: if we already stored the Secret name and the Secret exists, we are done.
	if cluster.Status.COSInstance.HMACSecretName != "" {
		existing := &corev1.Secret{}
		err := s.Client.Get(ctx, types.NamespacedName{
			Namespace: cluster.Namespace,
			Name:      cluster.Status.COSInstance.HMACSecretName,
		}, existing)
		if err == nil {
			log.V(3).Info("COS HMAC Secret already exists, skipping creation", "secret", cluster.Status.COSInstance.HMACSecretName)
			return nil
		}
	}

	log.Info("Creating COS HMAC service credential", "instance", instanceCRN)

	params := &resourcecontrollerv2.ResourceKeyPostParameters{}
	params.SetProperty("HMAC", true)

	keyOpts := &resourcecontrollerv2.CreateResourceKeyOptions{
		Name:       ptr.To(ResourceName(cluster.Name, ResourceTypeCOSHMACKey, "")),
		Source:     ptr.To(instanceCRN),
		Parameters: params,
	}

	resourceKey, _, err := s.ResourceClient.CreateResourceKey(keyOpts)
	if err != nil {
		return fmt.Errorf("failed to create COS HMAC resource key: %w", err)
	}

	// The IBM Resource Controller returns HMAC keys under the "cos_hmac_keys" property
	// of the credential credentials object.
	hmacRaw := resourceKey.Credentials.GetProperty("cos_hmac_keys")
	if hmacRaw == nil {
		return fmt.Errorf("COS HMAC key created but cos_hmac_keys property is missing in response")
	}

	hmacMap, ok := hmacRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("COS HMAC key cos_hmac_keys has unexpected type %T", hmacRaw)
	}

	accessKeyID, _ := hmacMap[cosHMACAccessKeyField].(string)
	secretAccessKey, _ := hmacMap[cosHMACSecretKeyField].(string)
	if accessKeyID == "" || secretAccessKey == "" {
		return fmt.Errorf("COS HMAC key is missing %s or %s", cosHMACAccessKeyField, cosHMACSecretKeyField)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         infrav1.GroupVersion.String(),
					Kind:               "IBMPowerVSCluster",
					Name:               cluster.Name,
					UID:                cluster.UID,
					Controller:         ptr.To(true),
					BlockOwnerDeletion: ptr.To(true),
				},
			},
		},
		Data: map[string][]byte{
			cosHMACAccessKeyField: []byte(accessKeyID),
			cosHMACSecretKeyField: []byte(secretAccessKey),
		},
	}

	if err := s.Client.Create(ctx, secret); err != nil {
		return fmt.Errorf("failed to store COS HMAC Secret %q: %w", secretName, err)
	}

	cluster.Status.COSInstance.HMACSecretName = secretName
	log.Info("COS HMAC Secret created successfully", "secret", secretName)
	return nil
}

// reconcileCOSReference handles verifying an explicitly referenced COS instance.
func (s *ClusterScope) reconcileCOSReference(_ context.Context, ref infrav1.ResourceIdentifier) (*resourcecontrollerv2.ResourceInstance, error) {
	if ref.ID != "" {
		instance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: &ref.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed fetching referenced COS instance by ID: %w", err)
		}
		if instance == nil {
			return nil, fmt.Errorf("referenced COS instance ID %s not found", ref.ID)
		}
		return instance, nil
	} else if ref.Name != "" {
		filter := resourcecontroller.InstanceFilter{
			Name:           ref.Name,
			ResourceID:     resourcecontroller.CosResourceID,
			ResourcePlanID: resourcecontroller.CosResourcePlanID,
		}
		instance, err := s.ResourceClient.GetResourceInstanceByFilter(filter)
		if err != nil {
			return nil, fmt.Errorf("failed fetching referenced COS instance by Name: %w", err)
		}
		if instance == nil {
			return nil, fmt.Errorf("referenced COS instance Name %s not found", ref.Name)
		}
		return instance, nil
	}
	return nil, fmt.Errorf("COS reference must have either ID or Name set")
}

// reconcileCOSProvision handles idempotently creating a new COS instance.
func (s *ClusterScope) reconcileCOSProvision(ctx context.Context, name string) (*resourcecontrollerv2.ResourceInstance, error) {
	log := ctrl.LoggerFrom(ctx)

	// Check if an instance with the target name already exists
	filter := resourcecontroller.InstanceFilter{
		Name:           name,
		ResourceID:     resourcecontroller.CosResourceID,
		ResourcePlanID: resourcecontroller.CosResourcePlanID,
	}
	instance, err := s.ResourceClient.GetResourceInstanceByFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("failed checking for existing COS instance: %w", err)
	}

	if instance != nil {
		log.V(3).Info("COS instance found in cloud by name matching", "name", name, "id", *instance.GUID)
		return instance, nil
	}

	// Create a new instance if it doesn't exist
	log.Info("Creating a new COS service instance", "name", name)
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("resource group ID is empty")
	}

	target := "Global"
	newInstance, _, err := s.ResourceClient.CreateResourceInstance(&resourcecontrollerv2.CreateResourceInstanceOptions{
		Name:           ptr.To(name),
		Target:         ptr.To(target),
		ResourceGroup:  ptr.To(resourceGroupID),
		ResourcePlanID: ptr.To(resourcecontroller.CosResourcePlanID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating COS instance via API: %w", err)
	}

	return newInstance, nil
}

// setupCOSClient builds the COS client for bucket manipulation.
func (s *ClusterScope) setupCOSClient(ctx context.Context, instanceID, bucketRegion string) error {
	// Skip if already initialized during a previous run in this reconciliation loop
	if s.COSClient != nil {
		return nil
	}

	if s.ClientBuilder == nil {
		return errors.New("ClientBuilder is not set; cannot create COS client")
	}

	cosClient, err := s.ClientBuilder.GetCOSClient(ctx, COSClientOptions{
		InstanceID:      instanceID,
		BucketRegion:    bucketRegion,
		ServiceEndpoint: s.ServiceEndpoint,
	})
	if err != nil {
		return fmt.Errorf("failed to create COS client: %w", err)
	}

	s.COSClient = cosClient
	ctrl.Log.V(3).Info("COS client created successfully")
	return nil
}

// reconcileCOSBucket checks for the existence of the bucket and creates it if missing.
func (s *ClusterScope) reconcileCOSBucket(ctx context.Context, bucketName string) error {
	log := ctrl.LoggerFrom(ctx)

	// Attempt to get the bucket
	_, err := s.COSClient.GetBucketByName(bucketName)
	if err == nil {
		log.V(3).Info("COS bucket already exists in cloud", "bucketName", bucketName)
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return fmt.Errorf("failed to check if COS bucket exists: %w", err)
	}

	log.V(3).Info("COS bucket lookup returned error", "code", aerr.Code(), "message", aerr.Message(), "origErr", aerr.OrigErr())

	switch aerr.Code() {
	case s3.ErrCodeNoSuchBucket, "Forbidden", "NotFound":
		log.Info("Creating new COS bucket", "bucketName", bucketName)

		input := &s3.CreateBucketInput{
			Bucket: ptr.To(bucketName),
		}

		out, err := s.COSClient.CreateBucket(input)
		if err != nil {
			if createErr, isAwsErr := err.(awserr.Error); isAwsErr {
				// BucketAlreadyOwnedByYou means we own it
				if createErr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou {
					log.V(3).Info("COS bucket already owned by us (concurrent create), treating as success", "bucketName", bucketName)
					return nil
				}
				// BucketAlreadyExists means the name is taken globally by a DIFFERENT account/instance.
				// The 403 Forbidden on HeadBucket confirms we don't own it. Surface this as a hard error
				// so the operator can surface it and the user can choose a different bucket name.
				if createErr.Code() == s3.ErrCodeBucketAlreadyExists {
					return fmt.Errorf("bucket name %q is already taken in the global IBM COS namespace by a different instance; set a unique bucketName in spec.cosInstance.bucketName: %w", bucketName, err)
				}
			}
			return fmt.Errorf("failed to execute CreateBucket API: %w", err)
		}
		log.V(3).Info("COS bucket created successfully", "bucketName", bucketName, "response", out)
		return nil

	default:
		return fmt.Errorf("unexpected error checking bucket presence: %w", err)
	}
}

// fetchVPCCRN returns VPC CRN.
func (s *ClusterScope) fetchVPCCRN() (*string, error) {
	vpcID := s.IBMPowerVSCluster.Status.VPC.ID
	if vpcID == "" {
		return nil, fmt.Errorf("failed to fetch VPC ID for VPC %v, ID is empty", s.IBMPowerVSCluster.Status.VPC.ID)
	}
	vpcDetails, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
		ID: ptr.To(vpcID),
	})
	if err != nil {
		return nil, err
	}
	return vpcDetails.CRN, nil
}

// fetchPowerVSWorkspaceCRN returns PowerVS workspace CRN.
func (s *ClusterScope) fetchPowerVSWorkspaceCRN() (*string, error) {
	workspaceID := s.IBMPowerVSCluster.Status.Workspace.ID
	workspace, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &workspaceID,
	})
	if err != nil {
		return nil, err
	}
	if workspace.CRN == nil {
		return nil, fmt.Errorf("workspace CRN is empty %s", workspaceID)
	}
	return workspace.CRN, nil
}

// DeleteVPCSecurityGroups deletes managed VPC security groups provisioned by the controller.
func (s *ClusterScope) DeleteVPCSecurityGroups(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)

	// Build a set of security group names that the controller provisioned.
	managedSGs := make(map[string]bool)
	for _, sgSource := range s.IBMPowerVSCluster.Spec.VPCSecurityGroups {
		if sgSource.Type == infrav1.SourceTypeProvision {
			name := sgSource.Provision.Name
			if name == "" {
				name = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeSG, "")
			}
			managedSGs[name] = true
		}
	}

	// Iterate through the status entries. Skip any that are not in the managed set.
	for _, sgStatus := range s.IBMPowerVSCluster.Status.VPCSecurityGroups {
		if !managedSGs[sgStatus.Name] {
			log.Info("Skipping VPC security group deletion as it is referenced, not managed by controller", "securityGroupName", sgStatus.Name)
			continue
		}

		// 3. Verify existence before deletion
		if _, resp, err := s.IBMVPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{
			ID: ptr.To(sgStatus.ID),
		}); err != nil {
			if resp != nil && resp.StatusCode == ResourceNotFoundCode {
				log.Info("VPC security group has already been deleted from cloud", "securityGroupID", sgStatus.ID)
				continue
			}
			return fmt.Errorf("failed to fetch VPC security group '%s' during deletion: %w", sgStatus.ID, err)
		}

		// 4. Execute Deletion
		log.V(3).Info("Issuing delete command for VPC security group", "securityGroupID", sgStatus.ID)
		options := &vpcv1.DeleteSecurityGroupOptions{
			ID: ptr.To(sgStatus.ID),
		}

		if _, err := s.IBMVPCClient.DeleteSecurityGroup(options); err != nil {
			return fmt.Errorf("failed to execute DeleteSecurityGroup API for '%s': %w", sgStatus.ID, err)
		}

		log.Info("VPC security group successfully deleted", "securityGroupID", sgStatus.ID, "securityGroupName", sgStatus.Name)
	}

	return nil
}

// DeleteLoadBalancer deletes provisioned load balancers.
//
//nolint:gocyclo // complexity is acceptable for deletion logic
func (s *ClusterScope) DeleteLoadBalancer(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	var errs []error
	requeue := false

	// Gather all load balancers configured for this cluster
	loadBalancers := s.IBMPowerVSCluster.Spec.LoadBalancers

	// Fallback to the default naming setup if the spec was totally omitted
	if len(loadBalancers) == 0 {
		loadBalancers = []infrav1.LoadBalancerSource{
			{
				Type: infrav1.SourceTypeProvision,
				Provision: infrav1.LoadBalancerProvision{
					Name: ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeLBPublic, ""),
					Type: infrav1.LoadBalancerTypePublic,
				},
			},
		}
	}

	for i, lbSource := range loadBalancers {
		if lbSource.Type == infrav1.SourceTypeReference {
			log.V(3).Info("Skipping load balancer deletion as it is an external reference resource")
			continue
		}

		provision := lbSource.Provision
		lbName := provision.Name
		if lbName == "" {
			lbType := ResourceTypeLBPublic
			if provision.Type == infrav1.LoadBalancerTypePrivate {
				lbType = ResourceTypeLBPrivate
			}
			qualifier := ""
			if i > 0 {
				qualifier = strconv.Itoa(i)
			}
			lbName = ResourceName(s.IBMPowerVSCluster.Name, lbType, qualifier)
		}

		// Look up the cached ID from the status slice
		var lbID string
		for _, lbStatus := range s.IBMPowerVSCluster.Status.LoadBalancers {
			if lbStatus.Name == lbName {
				lbID = lbStatus.ID
				break
			}
		}

		var cloudLB *vpcv1.LoadBalancer
		var err error
		isNotFound := false

		// Efficient Fetch: Only use the ID if we have it, otherwise fallback to Name
		if lbID != "" {
			var resp *core.DetailedResponse
			cloudLB, resp, err = s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: ptr.To(lbID),
			})
			// Check if the cloud returned a 404 (Not Found)
			if err != nil && resp != nil && resp.StatusCode == ResourceNotFoundCode {
				isNotFound = true
			}
		} else {
			cloudLB, err = s.IBMVPCClient.GetLoadBalancerByName(lbName)
			if err == nil && cloudLB == nil {
				isNotFound = true
			}
		}

		// Handle fetch errors (ignoring 404s since we want to delete anyway)
		if err != nil && !isNotFound {
			errs = append(errs, fmt.Errorf("failed to fetch load balancer %s details: %w", lbName, err))
			continue
		}

		// If it's already gone, we have nothing left to do for this iteration!
		if isNotFound || cloudLB == nil {
			log.Info("Load balancer already deleted or never existed", "name", lbName)
			continue
		}

		// Ensure we have the definitive ID and state from the cloud payload
		lbID = ptr.Deref(cloudLB.ID, "")
		statusState := ptr.Deref(cloudLB.ProvisioningStatus, "")

		// If it's already mid-deletion, we must requeue to check again later
		if statusState == string(infrav1.LoadBalancerStateDeletePending) {
			log.V(3).Info("Load balancer is currently being deleted, waiting", "name", lbName)
			requeue = true
			continue
		}

		// Issue the direct deletion command
		log.Info("Sending delete request for load balancer", "name", lbName, "id", lbID)
		if _, err = s.IBMVPCClient.DeleteLoadBalancer(&vpcv1.DeleteLoadBalancerOptions{
			ID: ptr.To(lbID),
		}); err != nil {
			errs = append(errs, fmt.Errorf("failed to trigger deletion for load balancer %s: %w", lbName, err))
			continue
		}

		// Deletion triggered successfully, requeue to verify complete removal on the next tick
		requeue = true
	}

	if len(errs) > 0 {
		return false, kerrors.NewAggregate(errs)
	}
	return requeue, nil
}

// DeleteVPCSubnets handles tearing down IBM Cloud VPC subnets if they were provisioned by the controller.
func (s *ClusterScope) DeleteVPCSubnets(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster
	var errs []error
	requeue := false

	// 1. Determine which subnets are actively managed (Provisioned) by the controller.
	managedSubnetNames := make(map[string]bool)
	if len(cluster.Spec.VPCSubnets) == 0 {
		// If spec is empty, all subnets tracked in status were auto-expanded defaults and thus provisioned.
		for _, ref := range cluster.Status.VPCSubnets {
			managedSubnetNames[ref.Name] = true
		}
	} else {
		// Otherwise, loop through spec and only flag subnets explicitly marked as 'Provision'.
		for i, subnetSpec := range cluster.Spec.VPCSubnets {
			if subnetSpec.Type == infrav1.SourceTypeProvision {
				name := subnetSpec.Provision.Name
				if name == "" {
					name = ResourceName(s.IBMPowerVSCluster.Name, ResourceTypeSubnet, strconv.Itoa(i))
				}
				managedSubnetNames[name] = true
			}
		}
	}

	// Updated to use the new VPCSubnetStatus struct
	var updatedStatus []infrav1.VPCSubnetStatus

	// 2. Process each subnet currently tracked in the status array.
	for _, subnetRef := range cluster.Status.VPCSubnets {
		// If it's a referenced subnet, skip deletion and do not keep it in the updated status
		if !managedSubnetNames[subnetRef.Name] {
			log.Info("Skipping VPC subnet deletion as it is referenced, not managed by controller", "name", subnetRef.Name)
			continue
		}

		if subnetRef.ID == "" {
			continue
		}

		// 3. Fetch current state from IBM Cloud
		net, resp, err := s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{ID: ptr.To(subnetRef.ID)})
		if err != nil {
			// If it's completely gone, we successfully deleted it.
			if resp != nil && resp.StatusCode == ResourceNotFoundCode {
				log.Info("VPC subnet successfully deleted from cloud", "id", subnetRef.ID, "name", subnetRef.Name)
				continue
			}
			errs = append(errs, fmt.Errorf("failed to fetch VPC subnet (id: %s) during deletion: %w", subnetRef.ID, err))
			updatedStatus = append(updatedStatus, subnetRef) // Keep in status to retry on next requeue
			continue
		}

		// 4. Check if deletion is already in progress
		if net != nil && net.Status != nil && *net.Status == string(infrav1.VPCSubnetStateDeleting) {
			log.V(3).Info("VPC subnet deletion is actively in progress", "id", subnetRef.ID)
			requeue = true
			updatedStatus = append(updatedStatus, subnetRef)
			continue
		}

		// 5. Issue the Delete API Call
		log.Info("Issuing delete command for VPC subnet", "id", subnetRef.ID, "name", subnetRef.Name)
		if _, err = s.IBMVPCClient.DeleteSubnet(&vpcv1.DeleteSubnetOptions{ID: net.ID}); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete VPC subnet (id: %s): %w", subnetRef.ID, err))
			updatedStatus = append(updatedStatus, subnetRef)
			continue
		}

		// Requeue to monitor progress until it returns 404
		requeue = true
		updatedStatus = append(updatedStatus, subnetRef)
	}

	// 6. Update the status array to drop subnets that successfully returned 404
	cluster.Status.VPCSubnets = updatedStatus

	if len(errs) > 0 {
		return false, kerrors.NewAggregate(errs)
	}
	return requeue, nil
}

// DeleteVPC handles tearing down the IBM Cloud VPC if it was provisioned by the controller.
func (s *ClusterScope) DeleteVPC(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// 1. Guard check: If the user referenced an existing VPC, do not delete it.
	if cluster.Spec.VPC.Type != infrav1.SourceTypeProvision {
		log.Info("Skipping VPC deletion as it is not managed/provisioned by the controller")
		return false, nil
	}

	// 2. If it's already completely wiped out from our status tracking, nothing to do.
	vpcID := cluster.Status.VPC.ID
	if vpcID == "" {
		log.Info("VPC already removed or untracked in status")
		return false, nil
	}

	// 3. Fetch the current runtime state of the VPC from IBM Cloud.
	vpcDetails, resp, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
		ID: ptr.To(vpcID),
	})

	if err != nil {
		if resp != nil && resp.StatusCode == ResourceNotFoundCode {
			log.Info("VPC successfully deleted from cloud")
			// Clear status using the new VPCStatus struct
			cluster.Status.VPC = infrav1.VPCStatus{}
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch VPC (id: %s) details during deletion: %w", vpcID, err)
	}

	if vpcDetails == nil {
		log.Info("VPC returned nil details, assuming gone")
		// Clear status using the new VPCStatus struct
		cluster.Status.VPC = infrav1.VPCStatus{}
		return false, nil
	}

	// 4. Short-circuit if the VPC deletion operation is currently actively running in the cloud.
	if vpcDetails.Status != nil && *vpcDetails.Status == string(infrav1.VPCStateDeleting) {
		log.V(3).Info("VPC deletion is actively in progress in IBM Cloud")
		return true, nil
	}

	// 5. Issue the actual Delete call to the IBM Cloud VPC API.
	log.Info("Issuing delete command for VPC", "vpcID", vpcID, "name", *vpcDetails.Name)
	if _, err = s.IBMVPCClient.DeleteVPC(&vpcv1.DeleteVPCOptions{
		ID: vpcDetails.ID,
	}); err != nil {
		return false, fmt.Errorf("failed to delete VPC (id: %s): %w", vpcID, err)
	}

	// Requeue immediately to hit step 4 and monitor progress until it returns a 404.
	return true, nil
}

// DeleteTransitGateway deletes transit gateway.
func (s *ClusterScope) DeleteTransitGateway(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	// 1. If we don't have an ID in status, there is nothing to delete/clean up
	tgID := s.IBMPowerVSCluster.Status.TransitGateway.ID
	if tgID == "" {
		return false, nil
	}

	// 2. Fetch the current state from the cloud
	tg, resp, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
		ID: ptr.To(tgID),
	})

	if err != nil {
		if resp != nil && resp.StatusCode == ResourceNotFoundCode {
			log.Info("Transit gateway successfully deleted (not found in cloud)")
			s.IBMPowerVSCluster.Status.TransitGateway = infrav1.TransitGatewayStatus{}
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch transit gateway during deletion: %w", err)
	}

	// 3. Handle pending deletions gracefully
	if tg.Status != nil && *tg.Status == string(infrav1.TransitGatewayStateDeletePending) {
		log.V(3).Info("Transit gateway is currently being deleted, requeuing")
		return true, nil
	}

	// 4. Clean up connections first
	requeue, err := s.deleteTransitGatewayConnections(ctx, tg)
	if err != nil {
		return false, err
	} else if requeue {
		return true, nil
	}

	// 5. Evaluate intent for the Transit Gateway itself
	if s.IBMPowerVSCluster.Spec.TransitGateway.Type == infrav1.SourceTypeReference {
		log.Info("Skipping Transit Gateway deletion because it is explicitly defined as a Reference")
		s.IBMPowerVSCluster.Status.TransitGateway = infrav1.TransitGatewayStatus{}
		return false, nil
	}

	// 6. Intent is Provision, so we issue the deletion command to IBM Cloud
	log.Info("Deleting Transit Gateway", "id", tgID)
	if _, err = s.TransitGatewayClient.DeleteTransitGateway(&tgapiv1.DeleteTransitGatewayOptions{
		ID: ptr.To(tgID),
	}); err != nil {
		return false, fmt.Errorf("failed to issue delete for transit gateway: %w", err)
	}

	return true, nil
}

// deleteTransitGatewayConnections evaluates intent for individual connections and deletes them if owned.
func (s *ClusterScope) deleteTransitGatewayConnections(ctx context.Context, tg *tgapiv1.TransitGateway) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	if tg == nil || tg.ID == nil {
		return false, fmt.Errorf("transit gateway or its ID is nil during connection deletion")
	}

	tgSpec := s.IBMPowerVSCluster.Spec.TransitGateway
	tgStatus := s.IBMPowerVSCluster.Status.TransitGateway

	deleteConnection := func(connID string) (bool, error) {
		if connID == "" {
			return false, nil
		}

		conn, resp, err := s.TransitGatewayClient.GetTransitGatewayConnection(&tgapiv1.GetTransitGatewayConnectionOptions{
			TransitGatewayID: tg.ID,
			ID:               ptr.To(connID),
		})

		if err != nil {
			if resp != nil && resp.StatusCode == ResourceNotFoundCode {
				log.V(3).Info("Connection deleted in transit gateway", "connectionID", connID)
				return false, nil
			}
			return false, fmt.Errorf("failed to get transit gateway connection: %w", err)
		}

		// Check for nil status to prevent panic before dereferencing
		if conn != nil && conn.Status != nil && *conn.Status == string(infrav1.TransitGatewayConnectionStateDeleting) {
			log.V(3).Info("Transit gateway connection is in deleting state", "connectionID", connID)
			return true, nil
		}

		if _, err = s.TransitGatewayClient.DeleteTransitGatewayConnection(&tgapiv1.DeleteTransitGatewayConnectionOptions{
			ID:               ptr.To(connID),
			TransitGatewayID: tg.ID,
		}); err != nil {
			return false, fmt.Errorf("failed to delete transit gateway connection: %w", err)
		}

		return true, nil
	}

	// 1. Delete PowerVS Connection only if user intended to provision it
	if tgSpec.PowerVSConnection.Type == infrav1.SourceTypeProvision && tgStatus.PowerVSConnection.ID != "" {
		log.V(3).Info("Deleting provisioned PowerVS connection in Transit gateway")
		requeue, err := deleteConnection(tgStatus.PowerVSConnection.ID)
		if err != nil {
			return false, err
		}
		if requeue {
			return requeue, nil
		}
	}

	// 2. Delete VPC Connection only if user intended to provision it
	if tgSpec.VPCConnection.Type == infrav1.SourceTypeProvision && tgStatus.VPCConnection.ID != "" {
		log.V(3).Info("Deleting provisioned VPC connection in Transit gateway")
		requeue, err := deleteConnection(tgStatus.VPCConnection.ID)
		if err != nil {
			return false, err
		}
		if requeue {
			return requeue, nil
		}
	}

	return false, nil
}

// DeleteDHCPServer deletes the DHCP server if it was provisioned by the controller.
func (s *ClusterScope) DeleteDHCPServer(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)

	// 1. Check if we provisioned this network
	if s.IBMPowerVSCluster.Spec.Network.Type != infrav1.SourceTypeProvision {
		log.Info("Skipping DHCP server deletion as network is in Reference mode")
		return nil
	}

	// 2. If the controller owns the workspace, deleting the workspace cascades
	// and destroys the DHCP server internally
	if s.IBMPowerVSCluster.Spec.Workspace.Type == infrav1.SourceTypeProvision {
		log.Info("Skipping separate DHCP server deletion as PowerVS workspace is being deleted by the controller (cascading delete)")
		return nil
	}

	// 3. Get DHCP server ID saved in status
	dhcpID := s.IBMPowerVSCluster.Status.Network.DHCPServer.ID
	if dhcpID == "" {
		log.Info("DHCP server ID not found in status, nothing to delete")
		return nil
	}

	// 4. Fetch the server to verify it exists
	server, err := s.IBMPowerVSClient.GetDHCPServer(ctx, dhcpID)
	if err != nil {
		// If it's a 404, we're already done!
		if strings.Contains(err.Error(), string(DHCPServerNotFound)) {
			log.Info("DHCP server no longer exists in IBM Cloud")
			return nil
		}
		return fmt.Errorf("failed to fetch DHCP server: %w", err)
	}

	// 5. Issue the delete command
	log.Info("Deleting provisioned DHCP server", "dhcpServerID", *server.ID)
	if err = s.IBMPowerVSClient.DeleteDHCPServer(ctx, *server.ID); err != nil {
		return fmt.Errorf("failed to delete DHCP server: %w", err)
	}

	return nil
}

// DeleteWorkspace deletes the PowerVS workspace if it was provisioned by the controller.
func (s *ClusterScope) DeleteWorkspace(ctx context.Context) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// 1. Declarative Safety Check: Only delete if the user asked us to provision it.
	if cluster.Spec.Workspace.Type != infrav1.SourceTypeProvision {
		log.Info("Skipping PowerVS workspace deletion as it was not provisioned by the controller")
		return false, nil
	}

	// 2. Check if there is actually an ID to delete.
	workspaceID := cluster.Status.Workspace.ID
	if workspaceID == "" {
		log.Info("PowerVS workspace ID is empty, nothing to delete")
		return false, nil
	}

	// 3. Fetch the current state of the workspace from IBM Cloud.
	workspace, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &workspaceID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to fetch PowerVS workspace: %w", err)
	}

	// 4. Check if it is already removed.
	if workspace != nil && workspace.State != nil && *workspace.State == string(infrav1.WorkspaceStateRemoved) {
		log.Info("PowerVS workspace has been removed")
		return false, nil
	}

	// 5. Trigger the deletion.
	log.Info("Deleting PowerVS workspace", "workspaceID", workspaceID)
	if _, err = s.ResourceClient.DeleteResourceInstance(&resourcecontrollerv2.DeleteResourceInstanceOptions{
		ID: &workspaceID,
	}); err != nil {
		return false, fmt.Errorf("failed to delete PowerVS workspace: %w", err)
	}

	// Return true to requeue so the controller can verify the deletion completed in the next loop.
	return true, nil
}

// DeleteCOSInstance handles tearing down the IBM Cloud COS instance and its contents
// if it was provisioned by the controller.
func (s *ClusterScope) DeleteCOSInstance(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx)
	cluster := s.IBMPowerVSCluster

	// 1. If it was never configured or never resolved into status, there is nothing to delete.
	if cluster.Status.COSInstance.ID == "" {
		return nil
	}
	// 2. We only delete the instance if the Type is explicitly set to Provision.
	if cluster.Spec.COSInstance.Type != infrav1.SourceTypeProvision {
		log.Info("Skipping COS instance deletion as it is referenced, not managed by controller",
			"id", cluster.Status.COSInstance.ID, "name", cluster.Status.COSInstance.Name)
		return nil
	}

	instanceID := cluster.Status.COSInstance.ID

	// 3. Fetch the current state from IBM Cloud
	cosInstance, resp, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &instanceID,
	})
	if err != nil {
		// If it's completely gone, we successfully finished deletion.
		if resp != nil && resp.StatusCode == ResourceNotFoundCode {
			log.Info("COS service instance successfully removed from cloud (returned 404)", "id", instanceID)
			return nil
		}
		return fmt.Errorf("failed to fetch COS service instance (id: %s) during deletion: %w", instanceID, err)
	}

	// 4. Check if deletion/reclamation is already processing in the cloud
	if cosInstance != nil && cosInstance.State != nil {
		state := *cosInstance.State
		if state == "pending_reclamation" || state == string(infrav1.WorkspaceStateRemoved) {
			log.Info("COS service instance deletion is actively in progress or reclaimed", "id", instanceID, "state", state)
			return nil
		}
	}

	// 5. Delete the HMAC Secret from Kubernetes (best-effort; log but do not block on failure).
	if cluster.Status.COSInstance.HMACSecretName != "" {
		hmacSecret := &corev1.Secret{}
		getErr := s.Client.Get(ctx, types.NamespacedName{
			Namespace: cluster.Namespace,
			Name:      cluster.Status.COSInstance.HMACSecretName,
		}, hmacSecret)
		if getErr == nil {
			if delErr := s.Client.Delete(ctx, hmacSecret); delErr != nil {
				log.Error(delErr, "Failed to delete COS HMAC Secret; continuing with instance deletion",
					"secret", cluster.Status.COSInstance.HMACSecretName)
			} else {
				log.Info("Deleted COS HMAC Secret", "secret", cluster.Status.COSInstance.HMACSecretName)
			}
		}
	}

	// 6. Issue the recursive Delete API Call
	log.Info("Issuing recursive delete command for managed COS service instance", "id", instanceID, "name", cluster.Status.COSInstance.Name)
	if _, err = s.ResourceClient.DeleteResourceInstance(&resourcecontrollerv2.DeleteResourceInstanceOptions{
		ID:        ptr.To(instanceID),
		Recursive: ptr.To(true), // Recursive drops associated buckets/bindings inside it
	}); err != nil {
		return fmt.Errorf("failed to execute DeleteResourceInstance API for %s: %w", instanceID, err)
	}

	log.Info("COS service instance delete command accepted successfully")
	return nil
}
