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

package scope

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	cosSession "github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	tgapiv1 "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/cos"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcemanager"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/transitgateway"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	genUtil "sigs.k8s.io/cluster-api-provider-ibmcloud/util"
)

const (
	// DEBUGLEVEL indicates the debug level of the logs.
	DEBUGLEVEL = 5
)

// networkConnectionType represents network connection type in transit gateway.
type networkConnectionType string

var (
	powervsNetworkConnectionType = networkConnectionType("power_virtual_server")
	vpcNetworkConnectionType     = networkConnectionType("vpc")
)

// powerEdgeRouter is identifier for PER.
const (
	powerEdgeRouter = "power-edge-router"
	// vpcSubnetIPAddressCount is the total IP Addresses for the subnet.
	// Support for custom address prefixes will be added at a later time. Currently, we use the ip count for subnet creation.
	vpcSubnetIPAddressCount int64 = 256
)

// PowerVSClusterScopeParams defines the input parameters used to create a new PowerVSClusterScope.
type PowerVSClusterScopeParams struct {
	Client            client.Client
	Logger            logr.Logger
	Cluster           *capiv1beta1.Cluster
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint

	// ClientFactory contains collection of functions to override actual client, which helps in testing.
	ClientFactory
}

// ClientFactory is collection of function used for overriding actual clients to help in testing.
type ClientFactory struct {
	AuthenticatorFactory      func() (core.Authenticator, error)
	PowerVSClientFactory      func() (powervs.PowerVS, error)
	VPCClientFactory          func() (vpc.Vpc, error)
	TransitGatewayFactory     func() (transitgateway.TransitGateway, error)
	ResourceControllerFactory func() (resourcecontroller.ResourceController, error)
	ResourceManagerFactory    func() (resourcemanager.ResourceManager, error)
}

// PowerVSClusterScope defines a scope defined around a Power VS Cluster.
type PowerVSClusterScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	IBMPowerVSClient      powervs.PowerVS
	IBMVPCClient          vpc.Vpc
	TransitGatewayClient  transitgateway.TransitGateway
	ResourceClient        resourcecontroller.ResourceController
	COSClient             cos.Cos
	ResourceManagerClient resourcemanager.ResourceManager

	Cluster           *capiv1beta1.Cluster
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
}

func getTGPowerVSConnectionName(tgName string) string { return fmt.Sprintf("%s-pvs-con", tgName) }

func getTGVPCConnectionName(tgName string) string { return fmt.Sprintf("%s-vpc-con", tgName) }

func dhcpNetworkName(dhcpServerName string) string {
	return fmt.Sprintf("DHCPSERVER%s_Private", dhcpServerName)
}

// NewPowerVSClusterScope creates a new PowerVSClusterScope from the supplied parameters.
func NewPowerVSClusterScope(params PowerVSClusterScopeParams) (*PowerVSClusterScope, error) {
	if params.Client == nil {
		err := errors.New("failed to generate new scope as client is nil")
		return nil, err
	}
	if params.Cluster == nil {
		err := errors.New("failed to generate new scope as cluster is nil")
		return nil, err
	}
	if params.IBMPowerVSCluster == nil {
		err := errors.New("failed to generate new scope IBMPowerVSCluster is nil")
		return nil, err
	}
	if params.Logger == (logr.Logger{}) {
		params.Logger = klog.Background()
	}

	helper, err := patch.NewHelper(params.IBMPowerVSCluster, params.Client)
	if err != nil {
		err = fmt.Errorf("failed to init patch helper: %w", err)
		return nil, err
	}

	// if powervs.cluster.x-k8s.io/create-infra=true annotation is not set, create only powerVSClient.
	if !CheckCreateInfraAnnotation(*params.IBMPowerVSCluster) {
		return &PowerVSClusterScope{
			Logger:            params.Logger,
			Client:            params.Client,
			patchHelper:       helper,
			Cluster:           params.Cluster,
			IBMPowerVSCluster: params.IBMPowerVSCluster,
			ServiceEndpoint:   params.ServiceEndpoint,
		}, nil
	}

	// if powervs.cluster.x-k8s.io/create-infra=true annotation is set, create necessary clients.
	piOptions := powervs.ServiceOptions{
		IBMPIOptions: &ibmpisession.IBMPIOptions{
			Debug: params.Logger.V(DEBUGLEVEL).Enabled(),
		},
	}

	// if Spec.ServiceInstanceID is set fetch zone associated with it or else use Spec.Zone.
	if params.IBMPowerVSCluster.Spec.ServiceInstanceID != "" {
		// Create Resource Controller client.
		var serviceOption resourcecontroller.ServiceOptions
		// Fetch the resource controller endpoint.
		rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), params.ServiceEndpoint)
		if rcEndpoint != "" {
			serviceOption.URL = rcEndpoint
			params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
		}
		rc, err := resourcecontroller.NewService(serviceOption)
		if err != nil {
			return nil, err
		}

		// Fetch the resource controller endpoint.
		if rcEndpoint := endpoints.FetchRCEndpoint(params.ServiceEndpoint); rcEndpoint != "" {
			params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
			if err := rc.SetServiceURL(rcEndpoint); err != nil {
				return nil, fmt.Errorf("failed to set resource controller endpoint: %w", err)
			}
		}

		res, _, err := rc.GetResourceInstance(
			&resourcecontrollerv2.GetResourceInstanceOptions{
				ID: core.StringPtr(params.IBMPowerVSCluster.Spec.ServiceInstanceID),
			})
		if err != nil {
			return nil, fmt.Errorf("failed to get resource instance: %w", err)
		}
		piOptions.Zone = *res.RegionID
		piOptions.CloudInstanceID = params.IBMPowerVSCluster.Spec.ServiceInstanceID
	} else {
		piOptions.Zone = *params.IBMPowerVSCluster.Spec.Zone
	}

	// Get the authenticator.
	auth, err := params.getAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator %w", err)
	}
	piOptions.Authenticator = auth

	// Create PowerVS client.
	powerVSClient, err := params.getPowerVSClient(piOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create PowerVS client %w", err)
	}

	// Create VPC client.
	vpcClient, err := params.getVPCClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC client: %w", err)
	}

	// Create TransitGateway client.
	tgOptions := &tgapiv1.TransitGatewayApisV1Options{
		Authenticator: auth,
	}

	tgClient, err := params.getTransitGatewayClient(tgOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create tranist gateway client: %w", err)
	}

	// Create Resource Controller client.
	serviceOption := resourcecontroller.ServiceOptions{
		ResourceControllerV2Options: &resourcecontrollerv2.ResourceControllerV2Options{
			Authenticator: auth,
		},
	}

	resourceClient, err := params.getResourceControllerClient(serviceOption)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource controller client: %w", err)
	}

	// Create Resource Manager client.
	rcManagerOptions := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: auth,
	}

	rmClient, err := params.getResourceManagerClient(rcManagerOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource manager client: %w", err)
	}

	clusterScope := &PowerVSClusterScope{
		Logger:                params.Logger,
		Client:                params.Client,
		patchHelper:           helper,
		Cluster:               params.Cluster,
		IBMPowerVSCluster:     params.IBMPowerVSCluster,
		ServiceEndpoint:       params.ServiceEndpoint,
		IBMPowerVSClient:      powerVSClient,
		IBMVPCClient:          vpcClient,
		TransitGatewayClient:  tgClient,
		ResourceClient:        resourceClient,
		ResourceManagerClient: rmClient,
	}
	return clusterScope, nil
}

func (params PowerVSClusterScopeParams) getAuthenticator() (core.Authenticator, error) {
	if params.AuthenticatorFactory != nil {
		return params.AuthenticatorFactory()
	}
	return authenticator.GetAuthenticator()
}

func (params PowerVSClusterScopeParams) getPowerVSClient(options powervs.ServiceOptions) (powervs.PowerVS, error) {
	if params.PowerVSClientFactory != nil {
		return params.PowerVSClientFactory()
	}

	// Fetch the PowerVS service endpoint.
	powerVSServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.PowerVS), params.ServiceEndpoint)
	if powerVSServiceEndpoint != "" {
		params.Logger.V(3).Info("Overriding the default PowerVS endpoint", "powerVSEndpoint", powerVSServiceEndpoint)
		options.URL = powerVSServiceEndpoint
	}
	return powervs.NewService(options)
}

func (params PowerVSClusterScopeParams) getVPCClient() (vpc.Vpc, error) {
	if params.Logger.V(DEBUGLEVEL).Enabled() {
		core.SetLoggingLevel(core.LevelDebug)
	}
	if params.VPCClientFactory != nil {
		return params.VPCClientFactory()
	}
	if params.IBMPowerVSCluster.Spec.VPC == nil || params.IBMPowerVSCluster.Spec.VPC.Region == nil {
		return nil, fmt.Errorf("failed to create VPC client as VPC info is nil")
	}
	// Fetch the VPC service endpoint.
	svcEndpoint := endpoints.FetchVPCEndpoint(*params.IBMPowerVSCluster.Spec.VPC.Region, params.ServiceEndpoint)
	return vpc.NewService(svcEndpoint)
}

func (params PowerVSClusterScopeParams) getTransitGatewayClient(options *tgapiv1.TransitGatewayApisV1Options) (transitgateway.TransitGateway, error) {
	if params.TransitGatewayFactory != nil {
		return params.TransitGatewayFactory()
	}
	// Fetch the TransitGateway service endpoint.
	tgServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.TransitGateway), params.ServiceEndpoint)
	if tgServiceEndpoint != "" {
		params.Logger.V(3).Info("Overriding the default TransitGateway endpoint", "transitGatewayEndpoint", tgServiceEndpoint)
		options.URL = tgServiceEndpoint
	}
	return transitgateway.NewService(options)
}

func (params PowerVSClusterScopeParams) getResourceControllerClient(options resourcecontroller.ServiceOptions) (resourcecontroller.ResourceController, error) {
	if params.ResourceControllerFactory != nil {
		return params.ResourceControllerFactory()
	}
	// Fetch the resource controller endpoint.
	rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), params.ServiceEndpoint)
	if rcEndpoint != "" {
		options.URL = rcEndpoint
		params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
	}
	return resourcecontroller.NewService(options)
}

func (params PowerVSClusterScopeParams) getResourceManagerClient(options *resourcemanagerv2.ResourceManagerV2Options) (resourcemanager.ResourceManager, error) {
	if params.ResourceManagerFactory != nil {
		return params.ResourceManagerFactory()
	}
	// Fetch the resource manager endpoint.
	rmEndpoint := endpoints.FetchEndpoints(string(endpoints.RM), params.ServiceEndpoint)
	if rmEndpoint != "" {
		options.URL = rmEndpoint
		params.Logger.V(3).Info("Overriding the default resource manager endpoint", "ResourceManagerEndpoint", rmEndpoint)
	}
	return resourcemanager.NewService(options)
}

// PatchObject persists the cluster configuration and status.
func (s *PowerVSClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMPowerVSCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *PowerVSClusterScope) Close() error {
	return s.PatchObject()
}

// Name returns the CAPI cluster name.
func (s *PowerVSClusterScope) Name() string {
	return s.Cluster.Name
}

// Zone returns the cluster zone.
func (s *PowerVSClusterScope) Zone() *string {
	return s.IBMPowerVSCluster.Spec.Zone
}

// ResourceGroup returns the cluster resource group.
func (s *PowerVSClusterScope) ResourceGroup() *infrav1beta2.IBMPowerVSResourceReference {
	return s.IBMPowerVSCluster.Spec.ResourceGroup
}

// InfraCluster returns the IBMPowerVS infrastructure cluster object name.
func (s *PowerVSClusterScope) InfraCluster() string {
	return s.IBMPowerVSCluster.Name
}

// APIServerPort returns the APIServerPort to use when creating the ControlPlaneEndpoint.
func (s *PowerVSClusterScope) APIServerPort() int32 {
	if s.Cluster.Spec.ClusterNetwork != nil && s.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return *s.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return infrav1beta2.DefaultAPIServerPort
}

// ServiceInstance returns the cluster ServiceInstance.
func (s *PowerVSClusterScope) ServiceInstance() *infrav1beta2.IBMPowerVSResourceReference {
	return s.IBMPowerVSCluster.Spec.ServiceInstance
}

// GetServiceInstanceID returns service instance id set in status field of IBMPowerVSCluster object. If it doesn't exist, returns empty string.
func (s *PowerVSClusterScope) GetServiceInstanceID() string {
	if s.IBMPowerVSCluster.Status.ServiceInstance != nil && s.IBMPowerVSCluster.Status.ServiceInstance.ID != nil {
		return *s.IBMPowerVSCluster.Status.ServiceInstance.ID
	}
	return ""
}

// SetTransitGatewayConnectionStatus sets the connection status of Transit gateway.
func (s *PowerVSClusterScope) SetTransitGatewayConnectionStatus(networkType networkConnectionType, resource *infrav1beta2.ResourceReference) {
	if s.IBMPowerVSCluster.Status.TransitGateway == nil || resource == nil {
		return
	}

	switch networkType {
	case powervsNetworkConnectionType:
		s.IBMPowerVSCluster.Status.TransitGateway.PowerVSConnection = resource
	case vpcNetworkConnectionType:
		s.IBMPowerVSCluster.Status.TransitGateway.VPCConnection = resource
	}
}

// SetTransitGatewayStatus sets the status of Transit gateway.
func (s *PowerVSClusterScope) SetTransitGatewayStatus(id *string, controllerCreated *bool) {
	s.IBMPowerVSCluster.Status.TransitGateway = &infrav1beta2.TransitGatewayStatus{
		ID:                id,
		ControllerCreated: controllerCreated,
	}
}

// TODO: Can we use generic here.

// SetStatus set the IBMPowerVSCluster status for provided ResourceType.
func (s *PowerVSClusterScope) SetStatus(resourceType infrav1beta2.ResourceType, resource infrav1beta2.ResourceReference) {
	s.V(3).Info("Setting status", "resourceType", resourceType, "resource", resource)
	switch resourceType {
	case infrav1beta2.ResourceTypeServiceInstance:
		if s.IBMPowerVSCluster.Status.ServiceInstance == nil {
			s.IBMPowerVSCluster.Status.ServiceInstance = &resource
			return
		}
		s.IBMPowerVSCluster.Status.ServiceInstance.Set(resource)
	case infrav1beta2.ResourceTypeNetwork:
		if s.IBMPowerVSCluster.Status.Network == nil {
			s.IBMPowerVSCluster.Status.Network = &resource
			return
		}
		s.IBMPowerVSCluster.Status.Network.Set(resource)
	case infrav1beta2.ResourceTypeVPC:
		if s.IBMPowerVSCluster.Status.VPC == nil {
			s.IBMPowerVSCluster.Status.VPC = &resource
			return
		}
		s.IBMPowerVSCluster.Status.VPC.Set(resource)
	case infrav1beta2.ResourceTypeDHCPServer:
		if s.IBMPowerVSCluster.Status.DHCPServer == nil {
			s.IBMPowerVSCluster.Status.DHCPServer = &resource
			return
		}
		s.IBMPowerVSCluster.Status.DHCPServer.Set(resource)
	case infrav1beta2.ResourceTypeCOSInstance:
		if s.IBMPowerVSCluster.Status.COSInstance == nil {
			s.IBMPowerVSCluster.Status.COSInstance = &resource
			return
		}
		s.IBMPowerVSCluster.Status.COSInstance.Set(resource)
	case infrav1beta2.ResourceTypeResourceGroup:
		if s.IBMPowerVSCluster.Status.ResourceGroup == nil {
			s.IBMPowerVSCluster.Status.ResourceGroup = &resource
			return
		}
		s.IBMPowerVSCluster.Status.ResourceGroup.Set(resource)
	}
}

// GetNetworkID returns the Network id from status of IBMPowerVSCluster object. If it doesn't exist, returns nil.
func (s *PowerVSClusterScope) GetNetworkID() *string {
	if s.IBMPowerVSCluster.Status.Network != nil {
		return s.IBMPowerVSCluster.Status.Network.ID
	}
	return nil
}

// Network returns the cluster Network.
func (s *PowerVSClusterScope) Network() *infrav1beta2.IBMPowerVSResourceReference {
	return &s.IBMPowerVSCluster.Spec.Network
}

// GetDHCPServerID returns the DHCP id from status of IBMPowerVSCluster object. If it doesn't exist, returns nil.
func (s *PowerVSClusterScope) GetDHCPServerID() *string {
	if s.IBMPowerVSCluster.Status.DHCPServer != nil {
		return s.IBMPowerVSCluster.Status.DHCPServer.ID
	}
	return nil
}

// DHCPServer returns the DHCP server details.
func (s *PowerVSClusterScope) DHCPServer() *infrav1beta2.DHCPServer {
	return s.IBMPowerVSCluster.Spec.DHCPServer
}

// VPC returns the cluster VPC information.
func (s *PowerVSClusterScope) VPC() *infrav1beta2.VPCResourceReference {
	return s.IBMPowerVSCluster.Spec.VPC
}

// GetVPCID returns the VPC id set in status field of IBMPowerVSCluster object. If it doesn't exist, returns nil.
func (s *PowerVSClusterScope) GetVPCID() *string {
	if s.IBMPowerVSCluster.Status.VPC != nil {
		return s.IBMPowerVSCluster.Status.VPC.ID
	}
	return nil
}

// GetVPCSubnetID returns the VPC subnet id.
func (s *PowerVSClusterScope) GetVPCSubnetID(subnetName string) *string {
	if s.IBMPowerVSCluster.Status.VPCSubnet == nil {
		return nil
	}
	if val, ok := s.IBMPowerVSCluster.Status.VPCSubnet[subnetName]; ok {
		return val.ID
	}
	return nil
}

// GetVPCSubnetIDs returns all the VPC subnet ids.
func (s *PowerVSClusterScope) GetVPCSubnetIDs() []*string {
	subnets := []*string{}
	if s.IBMPowerVSCluster.Status.VPCSubnet == nil {
		return nil
	}
	for _, subnet := range s.IBMPowerVSCluster.Status.VPCSubnet {
		subnets = append(subnets, subnet.ID)
	}
	return subnets
}

// SetVPCSubnetStatus set the VPC subnet id.
func (s *PowerVSClusterScope) SetVPCSubnetStatus(name string, resource infrav1beta2.ResourceReference) {
	s.V(3).Info("Setting status", "name", name, "resource", resource)
	if s.IBMPowerVSCluster.Status.VPCSubnet == nil {
		s.IBMPowerVSCluster.Status.VPCSubnet = make(map[string]infrav1beta2.ResourceReference)
	}
	if val, ok := s.IBMPowerVSCluster.Status.VPCSubnet[name]; ok {
		if val.ControllerCreated != nil && *val.ControllerCreated {
			resource.ControllerCreated = val.ControllerCreated
		}
	}
	s.IBMPowerVSCluster.Status.VPCSubnet[name] = resource
}

// GetVPCSecurityGroupByName returns the VPC security group id and its ruleIDs.
func (s *PowerVSClusterScope) GetVPCSecurityGroupByName(name string) (*string, []*string, *bool) {
	if s.IBMPowerVSCluster.Status.VPCSecurityGroups == nil {
		return nil, nil, nil
	}
	if val, ok := s.IBMPowerVSCluster.Status.VPCSecurityGroups[name]; ok {
		return val.ID, val.RuleIDs, val.ControllerCreated
	}
	return nil, nil, nil
}

// GetVPCSecurityGroupByID returns the VPC security group's ruleIDs.
func (s *PowerVSClusterScope) GetVPCSecurityGroupByID(securityGroupID string) (*string, []*string, *bool) {
	if s.IBMPowerVSCluster.Status.VPCSecurityGroups == nil {
		return nil, nil, nil
	}
	for _, sg := range s.IBMPowerVSCluster.Status.VPCSecurityGroups {
		if *sg.ID == securityGroupID {
			return sg.ID, sg.RuleIDs, sg.ControllerCreated
		}
	}
	return nil, nil, nil
}

// SetVPCSecurityGroupStatus set the VPC security group id.
func (s *PowerVSClusterScope) SetVPCSecurityGroupStatus(name string, resource infrav1beta2.VPCSecurityGroupStatus) {
	s.V(3).Info("Setting VPC security group status", "name", name, "resource", resource)
	if s.IBMPowerVSCluster.Status.VPCSecurityGroups == nil {
		s.IBMPowerVSCluster.Status.VPCSecurityGroups = make(map[string]infrav1beta2.VPCSecurityGroupStatus)
	}
	if val, ok := s.IBMPowerVSCluster.Status.VPCSecurityGroups[name]; ok {
		if val.ControllerCreated != nil && *val.ControllerCreated {
			resource.ControllerCreated = val.ControllerCreated
		}
	}
	s.IBMPowerVSCluster.Status.VPCSecurityGroups[name] = resource
}

// TransitGateway returns the cluster Transit Gateway information.
func (s *PowerVSClusterScope) TransitGateway() *infrav1beta2.TransitGateway {
	return s.IBMPowerVSCluster.Spec.TransitGateway
}

// GetTransitGatewayID returns the transit gateway id set in status field of IBMPowerVSCluster object. If it doesn't exist, returns empty string.
func (s *PowerVSClusterScope) GetTransitGatewayID() *string {
	if s.IBMPowerVSCluster.Status.TransitGateway != nil {
		return s.IBMPowerVSCluster.Status.TransitGateway.ID
	}
	return nil
}

// SetLoadBalancerStatus set the loadBalancer id.
func (s *PowerVSClusterScope) SetLoadBalancerStatus(name string, loadBalancer infrav1beta2.VPCLoadBalancerStatus) {
	s.V(3).Info("Setting status", "name", name, "status", loadBalancer)
	if s.IBMPowerVSCluster.Status.LoadBalancers == nil {
		s.IBMPowerVSCluster.Status.LoadBalancers = make(map[string]infrav1beta2.VPCLoadBalancerStatus)
	}
	if val, ok := s.IBMPowerVSCluster.Status.LoadBalancers[name]; ok {
		if val.ControllerCreated != nil && *val.ControllerCreated {
			loadBalancer.ControllerCreated = val.ControllerCreated
		}
	}
	s.IBMPowerVSCluster.Status.LoadBalancers[name] = loadBalancer
}

// GetLoadBalancerID returns the loadBalancer.
func (s *PowerVSClusterScope) GetLoadBalancerID(loadBalancerName string) *string {
	if s.IBMPowerVSCluster.Status.LoadBalancers == nil {
		return nil
	}
	if val, ok := s.IBMPowerVSCluster.Status.LoadBalancers[loadBalancerName]; ok {
		return val.ID
	}
	return nil
}

// GetLoadBalancerState will return the state for the load balancer.
func (s *PowerVSClusterScope) GetLoadBalancerState(name string) *infrav1beta2.VPCLoadBalancerState {
	if s.IBMPowerVSCluster.Status.LoadBalancers == nil {
		return nil
	}
	if val, ok := s.IBMPowerVSCluster.Status.LoadBalancers[name]; ok {
		return &val.State
	}
	return nil
}

// GetPublicLoadBalancerHostName will return the hostname of the public load balancer.
func (s *PowerVSClusterScope) GetPublicLoadBalancerHostName() (*string, error) {
	if s.IBMPowerVSCluster.Status.LoadBalancers == nil {
		return nil, nil
	}

	var name string
	if len(s.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		name = *s.GetServiceName(infrav1beta2.ResourceTypeLoadBalancer)
	}

	for _, lb := range s.IBMPowerVSCluster.Spec.LoadBalancers {
		if !*lb.Public {
			continue
		}

		if lb.Name != "" {
			name = lb.Name
			break
		}
		if lb.ID != nil {
			loadBalancer, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: lb.ID,
			})
			if err != nil {
				return nil, err
			}
			name = *loadBalancer.Name
			break
		}
	}

	if val, ok := s.IBMPowerVSCluster.Status.LoadBalancers[name]; ok {
		return val.Hostname, nil
	}
	return nil, nil
}

// GetResourceGroupID returns the resource group id if it present under spec or status filed of IBMPowerVSCluster object
// or returns empty string.
func (s *PowerVSClusterScope) GetResourceGroupID() string {
	if s.IBMPowerVSCluster.Spec.ResourceGroup != nil && s.IBMPowerVSCluster.Spec.ResourceGroup.ID != nil {
		return *s.IBMPowerVSCluster.Spec.ResourceGroup.ID
	}
	if s.IBMPowerVSCluster.Status.ResourceGroup != nil && s.IBMPowerVSCluster.Status.ResourceGroup.ID != nil {
		return *s.IBMPowerVSCluster.Status.ResourceGroup.ID
	}
	return ""
}

// IsPowerVSZoneSupportsPER checks whether PowerVS zone supports PER capabilities.
func (s *PowerVSClusterScope) IsPowerVSZoneSupportsPER() error {
	zone := s.Zone()
	if zone == nil {
		return fmt.Errorf("PowerVS zone is not set")
	}
	// fetch the datacenter capabilities for zone.
	datacenterCapabilities, err := s.IBMPowerVSClient.GetDatacenterCapabilities(*zone)
	if err != nil {
		return err
	}
	// check for the PER support in datacenter capabilities.
	perAvailable, ok := datacenterCapabilities[powerEdgeRouter]
	if !ok {
		return fmt.Errorf("%s capability unknown for zone: %s", powerEdgeRouter, *zone)
	}
	if !perAvailable {
		return fmt.Errorf("%s is not available for zone: %s", powerEdgeRouter, *zone)
	}
	return nil
}

// ReconcileResourceGroup reconciles resource group to fetch resource group id.
func (s *PowerVSClusterScope) ReconcileResourceGroup() error {
	// Verify if resource group id is set in spec or status field of IBMPowerVSCluster object.
	if resourceGroupID := s.GetResourceGroupID(); resourceGroupID != "" {
		return nil
	}
	// Try to fetch resource group id from cloud associated with resource group name.
	resourceGroupID, err := s.fetchResourceGroupID()
	if err != nil {
		return err
	}
	s.V(3).Info("Fetched resource group ID from IBM Cloud", "resourceGroupID", resourceGroupID)
	// Set the status of IBMPowerVSCluster object with resource group id.
	s.SetStatus(infrav1beta2.ResourceTypeResourceGroup, infrav1beta2.ResourceReference{ID: &resourceGroupID, ControllerCreated: ptr.To(false)})
	return nil
}

// ReconcilePowerVSServiceInstance reconciles Power VS service instance.
func (s *PowerVSClusterScope) ReconcilePowerVSServiceInstance() (bool, error) {
	// Verify if service instance id is set in status field of IBMPowerVSCluster object.
	serviceInstanceID := s.GetServiceInstanceID()
	if serviceInstanceID != "" {
		s.V(3).Info("PowerVS service instance ID is set, fetching details", "serviceInstanceID", serviceInstanceID)
		// if serviceInstanceID is set, verify that it exist and in active state.
		serviceInstance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: &serviceInstanceID,
		})
		if err != nil {
			return false, err
		}
		if serviceInstance == nil {
			return false, fmt.Errorf("failed to get PowerVS service instance with ID %s", serviceInstanceID)
		}

		requeue, err := s.checkServiceInstanceState(*serviceInstance)
		if err != nil {
			return false, err
		}
		return requeue, nil
	}

	// check PowerVS service instance exist in cloud, if it does not exist proceed with creating the instance.
	serviceInstanceID, requeue, err := s.isServiceInstanceExists()
	if err != nil {
		return false, err
	}
	// Set the status of IBMPowerVSCluster object with serviceInstanceID and ControllerCreated to false as PowerVS service instance is already exist in cloud.
	if serviceInstanceID != "" {
		s.V(3).Info("Found PowerVS service instance in IBM Cloud", "serviceInstanceID", serviceInstanceID)
		s.SetStatus(infrav1beta2.ResourceTypeServiceInstance, infrav1beta2.ResourceReference{ID: &serviceInstanceID, ControllerCreated: ptr.To(false)})
		return requeue, nil
	}

	// create PowerVS Service Instance
	serviceInstance, err := s.createServiceInstance()
	if err != nil {
		return false, fmt.Errorf("failed to create PowerVS service instance: %w", err)
	}
	if serviceInstance == nil {
		return false, fmt.Errorf("created PowerVS service instance is nil")
	}

	s.Info("Created PowerVS service instance", "serviceInstanceID", serviceInstance.GUID)
	// Set the status of IBMPowerVSCluster object with serviceInstanceID and ControllerCreated to true as new PowerVS service instance is created.
	s.SetStatus(infrav1beta2.ResourceTypeServiceInstance, infrav1beta2.ResourceReference{ID: serviceInstance.GUID, ControllerCreated: ptr.To(true)})
	return true, nil
}

// checkServiceInstanceState checks the state of a PowerVS service instance.
// If state is provisioning, true is returned indicating a requeue for reconciliation.
// In all other cases, it returns false.
func (s *PowerVSClusterScope) checkServiceInstanceState(instance resourcecontrollerv2.ResourceInstance) (bool, error) {
	s.V(3).Info("Checking the state of PowerVS service instance", "name", *instance.Name)
	switch *instance.State {
	case string(infrav1beta2.ServiceInstanceStateActive):
		s.V(3).Info("PowerVS service instance is in active state")
		return false, nil
	case string(infrav1beta2.ServiceInstanceStateProvisioning):
		s.V(3).Info("PowerVS service instance is in provisioning state")
		return true, nil
	case string(infrav1beta2.ServiceInstanceStateFailed):
		return false, fmt.Errorf("PowerVS service instance is in failed state")
	}
	return false, fmt.Errorf("PowerVS service instance is in %s state", *instance.State)
}

// checkServiceInstance checks PowerVS service instance exist in cloud by ID or name.
func (s *PowerVSClusterScope) isServiceInstanceExists() (string, bool, error) {
	s.V(3).Info("Checking for PowerVS service instance in IBM Cloud")
	var (
		id              string
		err             error
		serviceInstance *resourcecontrollerv2.ResourceInstance
	)

	if s.IBMPowerVSCluster.Spec.ServiceInstanceID != "" {
		id = s.IBMPowerVSCluster.Spec.ServiceInstanceID
	} else if s.IBMPowerVSCluster.Spec.ServiceInstance != nil && s.IBMPowerVSCluster.Spec.ServiceInstance.ID != nil {
		id = *s.IBMPowerVSCluster.Spec.ServiceInstance.ID
	}

	if id != "" {
		// Fetches service instance by ID.
		serviceInstance, _, err = s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: &id,
		})
	} else {
		// Fetches service instance by name.
		serviceInstance, err = s.getServiceInstance()
	}

	if err != nil {
		s.Error(err, "failed to get PowerVS service instance")
		return "", false, err
	}

	if serviceInstance == nil {
		s.V(3).Info("PowerVS service instance with given ID or name does not exist in IBM Cloud")
		return "", false, nil
	}

	requeue, err := s.checkServiceInstanceState(*serviceInstance)
	if err != nil {
		return "", false, err
	}

	return *serviceInstance.GUID, requeue, nil
}

// getServiceInstance return resource instance by name.
func (s *PowerVSClusterScope) getServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	//TODO: Support regular expression
	return s.ResourceClient.GetServiceInstance("", *s.GetServiceName(infrav1beta2.ResourceTypeServiceInstance), s.IBMPowerVSCluster.Spec.Zone)
}

// createServiceInstance creates the service instance.
func (s *PowerVSClusterScope) createServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	// fetch resource group id.
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroup())
	}

	// create service instance.
	s.V(3).Info("Creating new PowerVS service instance", "name", s.GetServiceName(infrav1beta2.ResourceTypeServiceInstance))
	zone := s.Zone()
	if zone == nil {
		return nil, fmt.Errorf("PowerVS zone is not set")
	}
	serviceInstance, _, err := s.ResourceClient.CreateResourceInstance(&resourcecontrollerv2.CreateResourceInstanceOptions{
		Name:           s.GetServiceName(infrav1beta2.ResourceTypeServiceInstance),
		Target:         zone,
		ResourceGroup:  &resourceGroupID,
		ResourcePlanID: ptr.To(resourcecontroller.PowerVSResourcePlanID),
	})
	if err != nil {
		return nil, err
	}
	return serviceInstance, nil
}

// ReconcileNetwork reconciles network
// If only IBMPowerVSCluster.Spec.Network is set, network would be validated and if exists already will get used as cluster’s network or a new network will be created via DHCP service.
// If only IBMPowerVSCluster.Spec.DHCPServer is set, DHCP server would be validated and if exists already, will use DHCP server’s network as cluster network. If not a new DHCP service will be created and it’s network will be used.
// If both IBMPowerVSCluster.Spec.Network & IBMPowerVSCluster.Spec.DHCPServer is set, network and DHCP server would be validated and if both exists already then network is belongs to given DHCP server or not would be validated.
// If both IBMPowerVSCluster.Spec.Network & IBMPowerVSCluster.Spec.DHCPServer is not set, by default DHCP service will be created to setup cluster's network.
func (s *PowerVSClusterScope) ReconcileNetwork() (bool, error) {
	if s.GetNetworkID() != nil {
		// Check the network exists
		if _, err := s.IBMPowerVSClient.GetNetworkByID(*s.GetNetworkID()); err != nil {
			return false, err
		}

		if s.GetDHCPServerID() == nil {
			// If only network is set, return once network is validated to be ok
			return true, nil
		}

		s.V(3).Info("DHCP server ID is set, fetching details", "dhcpServerID", s.GetDHCPServerID())
		active, err := s.isDHCPServerActive()
		if err != nil {
			return false, err
		}
		// DHCP server still not active, skip checking network for now
		if !active {
			return false, nil
		}
		return true, nil
	}
	// check network exist in cloud
	networkID, err := s.checkNetwork()
	if err != nil {
		return false, err
	}
	if networkID != nil {
		s.V(3).Info("Found PowerVS network in IBM Cloud", "networkID", networkID)
		s.SetStatus(infrav1beta2.ResourceTypeNetwork, infrav1beta2.ResourceReference{ID: networkID, ControllerCreated: ptr.To(false)})
	}
	dhcpServerID, err := s.checkDHCPServer()
	if err != nil {
		return false, err
	}
	if dhcpServerID != nil {
		s.V(3).Info("Found DHCP server in IBM Cloud", "dhcpServerID", dhcpServerID)
		s.SetStatus(infrav1beta2.ResourceTypeDHCPServer, infrav1beta2.ResourceReference{ID: dhcpServerID, ControllerCreated: ptr.To(false)})
	}
	if s.GetNetworkID() != nil {
		return true, nil
	}

	dhcpServerID, err = s.createDHCPServer()
	if err != nil {
		s.Error(err, "Error creating DHCP server")
		return false, err
	}

	s.Info("Created DHCP Server", "dhcpServerID", *dhcpServerID)
	s.SetStatus(infrav1beta2.ResourceTypeDHCPServer, infrav1beta2.ResourceReference{ID: dhcpServerID, ControllerCreated: ptr.To(true)})
	return false, nil
}

// checkDHCPServer checks if DHCP server exists in cloud with given DHCPServer's ID or name mentioned in spec.
// If exists and s.IBMPowerVSCluster.Status.Network is not populated will set DHCP server's network as cluster's network.
// If exists and s.IBMPowerVSCluster.Status.Network is populated already will validate the DHCP server's network and cluster networks are matching, if not will throw an error.
func (s *PowerVSClusterScope) checkDHCPServer() (*string, error) {
	if s.DHCPServer() != nil && s.DHCPServer().ID != nil {
		dhcpServer, err := s.IBMPowerVSClient.GetDHCPServer(*s.DHCPServer().ID)
		if err != nil {
			return nil, err
		}
		if s.GetNetworkID() == nil {
			if dhcpServer.Network != nil {
				if _, err := s.IBMPowerVSClient.GetNetworkByID(*dhcpServer.Network.ID); err != nil {
					return nil, err
				}
				s.SetStatus(infrav1beta2.ResourceTypeNetwork, infrav1beta2.ResourceReference{ID: dhcpServer.Network.ID, ControllerCreated: ptr.To(false)})
			} else {
				return nil, fmt.Errorf("found DHCP server with ID `%s`, but network is nil", *s.DHCPServer().ID)
			}
		} else if dhcpServer.Network != nil && *dhcpServer.Network.ID != *s.GetNetworkID() {
			return nil, fmt.Errorf("network details set via spec and DHCP server's network are not matching")
		}
		return dhcpServer.ID, nil
	}

	// if user provides DHCP server name then we can use network name to match the existing DHCP server
	var networkName string
	if s.DHCPServer() != nil && s.DHCPServer().Name != nil {
		networkName = dhcpNetworkName(*s.DHCPServer().Name)
	} else {
		networkName = dhcpNetworkName(s.InfraCluster())
	}

	s.V(3).Info("Checking DHCP server's network list by network name", "name", networkName)
	dhcpServers, err := s.IBMPowerVSClient.GetAllDHCPServers()
	if err != nil {
		return nil, err
	}
	for _, dhcpServer := range dhcpServers {
		if dhcpServer.Network != nil && *dhcpServer.Network.Name == networkName {
			if s.GetNetworkID() == nil {
				if _, err := s.IBMPowerVSClient.GetNetworkByID(*dhcpServer.Network.ID); err != nil {
					return nil, err
				}
				s.SetStatus(infrav1beta2.ResourceTypeNetwork, infrav1beta2.ResourceReference{ID: dhcpServer.Network.ID, ControllerCreated: ptr.To(false)})
			} else if *dhcpServer.Network.ID != *s.GetNetworkID() {
				return nil, fmt.Errorf("error network set via spec and DHCP server's networkID are not matching")
			}
			return dhcpServer.ID, nil
		}
	}

	return nil, nil
}

// checkNetwork checks if network exists in cloud with given network's ID or name mentioned in spec.
func (s *PowerVSClusterScope) checkNetwork() (*string, error) {
	if s.Network().ID != nil {
		s.V(3).Info("Checking if PowerVS network exists in IBM Cloud with ID", "networkID", *s.Network().ID)
		network, err := s.IBMPowerVSClient.GetNetworkByID(*s.Network().ID)
		if err != nil {
			return nil, err
		}
		return network.NetworkID, nil
	}

	if s.Network().Name != nil {
		s.V(3).Info("Checking if PowerVS network exists in IBM Cloud with network name", "name", s.Network().Name)
		network, err := s.IBMPowerVSClient.GetNetworkByName(*s.Network().Name)
		if err != nil {
			return nil, err
		}
		if network == nil || network.NetworkID == nil {
			s.V(3).Info("Unable to find PowerVS network in IBM Cloud", "network", s.IBMPowerVSCluster.Spec.Network)
			return nil, nil
		}
		return network.NetworkID, nil
	}
	return nil, nil
}

// isDHCPServerActive checks if the DHCP server status is active.
func (s *PowerVSClusterScope) isDHCPServerActive() (bool, error) {
	dhcpServer, err := s.IBMPowerVSClient.GetDHCPServer(*s.GetDHCPServerID())
	if err != nil {
		return false, err
	}

	active, err := s.checkDHCPServerStatus(*dhcpServer)
	if err != nil {
		return false, err
	}
	return active, nil
}

// checkDHCPServerStatus checks the state of a DHCP server.
// If state is active, true is returned.
// In all other cases, it returns false.
func (s *PowerVSClusterScope) checkDHCPServerStatus(dhcpServer models.DHCPServerDetail) (bool, error) {
	s.V(3).Info("Checking the status of DHCP server", "dhcpServerID", *dhcpServer.ID)
	switch *dhcpServer.Status {
	case string(infrav1beta2.DHCPServerStateActive):
		s.V(3).Info("DHCP server is in active state")
		return true, nil
	case string(infrav1beta2.DHCPServerStateBuild):
		s.V(3).Info("DHCP server is in build state")
		return false, nil
	case string(infrav1beta2.DHCPServerStateError):
		return false, fmt.Errorf("DHCP server creation failed and is in error state")
	}
	return false, nil
}

// createDHCPServer creates the DHCP server.
func (s *PowerVSClusterScope) createDHCPServer() (*string, error) {
	var dhcpServerCreateParams models.DHCPServerCreate
	dhcpServerDetails := s.DHCPServer()
	if dhcpServerDetails == nil {
		dhcpServerDetails = &infrav1beta2.DHCPServer{}
	}

	dhcpServerCreateParams.Name = s.GetServiceName(infrav1beta2.ResourceTypeDHCPServer)
	s.V(3).Info("Creating a new DHCP server with name", "name", dhcpServerCreateParams.Name)
	if dhcpServerDetails.DNSServer != nil {
		dhcpServerCreateParams.DNSServer = dhcpServerDetails.DNSServer
	}
	if dhcpServerDetails.Cidr != nil {
		dhcpServerCreateParams.Cidr = dhcpServerDetails.Cidr
	}
	if dhcpServerDetails.Snat != nil {
		dhcpServerCreateParams.SnatEnabled = dhcpServerDetails.Snat
	}

	dhcpServer, err := s.IBMPowerVSClient.CreateDHCPServer(&dhcpServerCreateParams)
	if err != nil {
		return nil, err
	}
	if dhcpServer == nil {
		return nil, fmt.Errorf("created DHCP server is nil")
	}
	if dhcpServer.Network == nil {
		return nil, fmt.Errorf("created DHCP server network is nil")
	}

	s.Info("DHCP Server network details", "details", *dhcpServer.Network)
	s.SetStatus(infrav1beta2.ResourceTypeNetwork, infrav1beta2.ResourceReference{ID: dhcpServer.Network.ID, ControllerCreated: ptr.To(true)})
	return dhcpServer.ID, nil
}

// ReconcileVPC reconciles VPC.
func (s *PowerVSClusterScope) ReconcileVPC() (bool, error) {
	// if VPC server id is set means the VPC is already created
	vpcID := s.GetVPCID()
	if vpcID != nil {
		s.V(3).Info("VPC ID is set, fetching details", "vpcID", *vpcID)
		vpcDetails, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
			ID: vpcID,
		})
		if err != nil {
			return false, err
		}
		if vpcDetails == nil {
			return false, fmt.Errorf("failed to get VPC with ID %s", *vpcID)
		}

		if vpcDetails.Status != nil && *vpcDetails.Status == string(infrav1beta2.VPCStatePending) {
			s.V(3).Info("VPC creation is in pending state")
			return true, nil
		}
		// TODO(karthik-k-n): Set status here as well
		return false, nil
	}

	// check vpc exist in cloud
	id, err := s.checkVPC()
	if err != nil {
		return false, err
	}
	if id != "" {
		s.V(3).Info("VPC found in IBM Cloud", "vpcID", id)
		s.SetStatus(infrav1beta2.ResourceTypeVPC, infrav1beta2.ResourceReference{ID: &id, ControllerCreated: ptr.To(false)})
		return false, nil
	}

	// TODO(karthik-k-n): create a generic cluster scope/service and implement common vpc logics, which can be consumed by both vpc and powervs

	// create VPC
	s.V(3).Info("Creating a VPC")
	vpcID, err = s.createVPC()
	if err != nil {
		return false, fmt.Errorf("failed to create VPC: %w", err)
	}
	s.Info("Created VPC", "vpcID", *vpcID)
	s.SetStatus(infrav1beta2.ResourceTypeVPC, infrav1beta2.ResourceReference{ID: vpcID, ControllerCreated: ptr.To(true)})
	return true, nil
}

// checkVPC checks VPC exist in cloud.
func (s *PowerVSClusterScope) checkVPC() (string, error) {
	var (
		err        error
		vpcDetails *vpcv1.VPC
	)
	if s.IBMPowerVSCluster.Spec.VPC != nil && s.IBMPowerVSCluster.Spec.VPC.ID != nil {
		vpcDetails, _, err = s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
			ID: s.IBMPowerVSCluster.Spec.VPC.ID,
		})
	} else {
		vpcDetails, err = s.getVPCByName()
	}

	if err != nil {
		s.Error(err, "failed to get VPC")
		return "", err
	}
	if vpcDetails == nil {
		s.V(3).Info("VPC not found in IBM Cloud", "vpc", s.IBMPowerVSCluster.Spec.VPC)
		return "", nil
	}
	s.V(3).Info("VPC found in IBM Cloud", "vpcID", *vpcDetails.ID)
	return *vpcDetails.ID, nil
}

func (s *PowerVSClusterScope) getVPCByName() (*vpcv1.VPC, error) {
	vpcDetails, err := s.IBMVPCClient.GetVPCByName(*s.GetServiceName(infrav1beta2.ResourceTypeVPC))
	if err != nil {
		return nil, err
	}
	return vpcDetails, nil
	//TODO: Support regular expression
}

// createVPC creates VPC.
func (s *PowerVSClusterScope) createVPC() (*string, error) {
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroup())
	}
	addressPrefixManagement := "auto"
	vpcOption := &vpcv1.CreateVPCOptions{
		ResourceGroup:           &vpcv1.ResourceGroupIdentity{ID: &resourceGroupID},
		Name:                    s.GetServiceName(infrav1beta2.ResourceTypeVPC),
		AddressPrefixManagement: &addressPrefixManagement,
	}
	vpcDetails, _, err := s.IBMVPCClient.CreateVPC(vpcOption)
	if err != nil {
		return nil, err
	}

	// set security group for vpc
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
		return nil, err
	}
	return vpcDetails.ID, nil
}

// ReconcileVPCSubnets reconciles VPC subnet.
func (s *PowerVSClusterScope) ReconcileVPCSubnets() (bool, error) {
	subnets := make([]infrav1beta2.Subnet, 0)
	vpcZones, err := regionUtil.VPCZonesForVPCRegion(*s.VPC().Region)
	if err != nil {
		return false, err
	}
	if len(vpcZones) == 0 {
		return false, fmt.Errorf("failed to fetch VPC zones, no zone found for region %s", *s.VPC().Region)
	}
	// check whether user has set the vpc subnets
	if len(s.IBMPowerVSCluster.Spec.VPCSubnets) == 0 {
		// if the user did not set any subnet, we try to create subnet in all the zones.
		for _, zone := range vpcZones {
			subnet := infrav1beta2.Subnet{
				Name: ptr.To(fmt.Sprintf("%s-%s", *s.GetServiceName(infrav1beta2.ResourceTypeSubnet), zone)),
				Zone: ptr.To(zone),
			}
			subnets = append(subnets, subnet)
		}
	} else {
		subnets = append(subnets, s.IBMPowerVSCluster.Spec.VPCSubnets...)
	}

	for index, subnet := range subnets {
		s.Info("Reconciling VPC subnet", "subnet", subnet)
		var subnetID *string
		if subnet.ID != nil {
			subnetID = subnet.ID
		} else {
			if subnet.Name == nil {
				subnet.Name = ptr.To(fmt.Sprintf("%s-%d", *s.GetServiceName(infrav1beta2.ResourceTypeSubnet), index))
			}
			subnetID = s.GetVPCSubnetID(*subnet.Name)
		}

		if subnetID != nil {
			s.V(3).Info("VPC subnet ID is set, fetching details", "subnetID", *subnetID)
			subnetDetails, _, err := s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{
				ID: subnetID,
			})
			if err != nil {
				return false, err
			}
			if subnetDetails == nil {
				return false, fmt.Errorf("failed to get VPC subnet with ID %s", *subnetID)
			}
			// check for next subnet
			s.SetVPCSubnetStatus(*subnetDetails.Name, infrav1beta2.ResourceReference{ID: subnetDetails.ID})
			continue
		}

		// check VPC subnet exist in cloud
		vpcSubnetID, err := s.checkVPCSubnet(*subnet.Name)
		if err != nil {
			s.Error(err, "error checking VPC subnet in IBM Cloud")
			return false, err
		}
		if vpcSubnetID != "" {
			s.V(3).Info("Found VPC subnet in IBM Cloud", "subnetID", vpcSubnetID)
			s.SetVPCSubnetStatus(*subnet.Name, infrav1beta2.ResourceReference{ID: &vpcSubnetID, ControllerCreated: ptr.To(false)})
			// check for next subnet
			continue
		}

		if subnet.Zone == nil {
			subnet.Zone = &vpcZones[index%len(vpcZones)]
		}
		s.V(3).Info("Creating VPC subnet")
		subnetID, err = s.createVPCSubnet(subnet)
		if err != nil {
			s.Error(err, "failed to create VPC subnet")
			return false, err
		}
		s.Info("Created VPC subnet", "subnetID", subnetID)
		s.SetVPCSubnetStatus(*subnet.Name, infrav1beta2.ResourceReference{ID: subnetID, ControllerCreated: ptr.To(true)})
		// Requeue only when the creation of all subnets has been triggered.
		if index == len(subnets)-1 {
			return true, nil
		}
	}
	return false, nil
}

// checkVPCSubnet checks if VPC subnet by the given name exists in cloud.
func (s *PowerVSClusterScope) checkVPCSubnet(subnetName string) (string, error) {
	vpcSubnet, err := s.IBMVPCClient.GetVPCSubnetByName(subnetName)
	if err != nil {
		return "", err
	}
	if vpcSubnet == nil {
		s.V(3).Info("VPC subnet not found in IBM Cloud")
		return "", nil
	}
	return *vpcSubnet.ID, nil
}

// createVPCSubnet creates a VPC subnet.
func (s *PowerVSClusterScope) createVPCSubnet(subnet infrav1beta2.Subnet) (*string, error) {
	// TODO(karthik-k-n): consider moving to clusterscope
	// fetch resource group id
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroup())
	}

	// create subnet
	vpcID := s.GetVPCID()
	if vpcID == nil {
		return nil, fmt.Errorf("VPC ID is empty")
	}

	ipVersion := vpcSubnetIPVersion4

	options := &vpcv1.CreateSubnetOptions{}
	options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
		IPVersion:             &ipVersion,
		TotalIpv4AddressCount: ptr.To(vpcSubnetIPAddressCount),
		Name:                  subnet.Name,
		VPC: &vpcv1.VPCIdentity{
			ID: vpcID,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: subnet.Zone,
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroupID,
		},
	})

	subnetDetails, _, err := s.IBMVPCClient.CreateSubnet(options)
	if err != nil {
		return nil, err
	}
	if subnetDetails == nil {
		return nil, fmt.Errorf("create VPC subnet is nil")
	}
	return subnetDetails.ID, nil
}

// ReconcileVPCSecurityGroups reconciles VPC security group.
func (s *PowerVSClusterScope) ReconcileVPCSecurityGroups() error {
	for _, securityGroup := range s.IBMPowerVSCluster.Spec.VPCSecurityGroups {
		var securityGroupID *string
		var securityGroupRuleIDs []*string

		if securityGroup.Name != nil {
			securityGroupID, securityGroupRuleIDs, _ = s.GetVPCSecurityGroupByName(*securityGroup.Name)
		} else {
			securityGroupID, securityGroupRuleIDs, _ = s.GetVPCSecurityGroupByID(*securityGroup.ID)
		}

		if securityGroupID != nil && securityGroupRuleIDs != nil {
			if _, _, err := s.IBMVPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{
				ID: securityGroupID,
			}); err != nil {
				return fmt.Errorf("failed to fetch existing security group '%s': %w", *securityGroupID, err)
			}
			for _, rule := range securityGroupRuleIDs {
				if _, _, err := s.IBMVPCClient.GetSecurityGroupRule(&vpcv1.GetSecurityGroupRuleOptions{
					SecurityGroupID: securityGroupID,
					ID:              rule,
				}); err != nil {
					return fmt.Errorf("failed to fetch rules of existing security group '%s': %w", *securityGroupID, err)
				}
			}
			continue
		}

		sg, ruleIDs, err := s.validateVPCSecurityGroup(securityGroup)
		if err != nil {
			return fmt.Errorf("failed to validate existing security group: %s", err)
		}
		if sg != nil {
			s.V(3).Info("VPC security group already exists", "name", *sg.Name)
			s.SetVPCSecurityGroupStatus(*sg.Name, infrav1beta2.VPCSecurityGroupStatus{
				ID:                sg.ID,
				RuleIDs:           ruleIDs,
				ControllerCreated: ptr.To(false),
			})
			continue
		}

		securityGroupID, err = s.createVPCSecurityGroup(securityGroup)
		if err != nil {
			return fmt.Errorf("failed to create VPC security group: %w", err)
		}
		s.Info("VPC security group created", "name", *securityGroup.Name)
		s.SetVPCSecurityGroupStatus(*securityGroup.Name, infrav1beta2.VPCSecurityGroupStatus{
			ID:                securityGroupID,
			ControllerCreated: ptr.To(true),
		})

		if err := s.createVPCSecurityGroupRulesAndSetStatus(securityGroup.Rules, securityGroupID, securityGroup.Name); err != nil {
			return err
		}
	}

	return nil
}

// createVPCSecurityGroupRule creates a specific rule for a existing security group.
func (s *PowerVSClusterScope) createVPCSecurityGroupRule(securityGroupID, direction, protocol *string, portMin, portMax *int64, remote infrav1beta2.VPCSecurityGroupRuleRemote) (*string, error) {
	setRemote := func(remote infrav1beta2.VPCSecurityGroupRuleRemote, remoteOption *vpcv1.SecurityGroupRuleRemotePrototype) error {
		switch remote.RemoteType {
		case infrav1beta2.VPCSecurityGroupRuleRemoteTypeCIDR:
			cidrSubnet, err := s.IBMVPCClient.GetVPCSubnetByName(*remote.CIDRSubnetName)
			if err != nil {
				return fmt.Errorf("failed to find VPC subnet by name '%s' for fetching CIDR block: %w", *remote.CIDRSubnetName, err)
			}
			if cidrSubnet == nil {
				return fmt.Errorf("VPC subnet by name '%s' does not exist", *remote.CIDRSubnetName)
			}
			s.V(3).Info("Creating VPC security group rule", "securityGroupID", *securityGroupID, "direction", *direction, "protocol", *protocol, "cidrBlockSubnet", *remote.CIDRSubnetName, "cidr", *cidrSubnet.Ipv4CIDRBlock)
			remoteOption.CIDRBlock = cidrSubnet.Ipv4CIDRBlock
		case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAddress:
			s.V(3).Info("Creating VPC security group rule", "securityGroupID", *securityGroupID, "direction", *direction, "protocol", *protocol, "ip", *remote.Address)
			remoteOption.Address = remote.Address
		case infrav1beta2.VPCSecurityGroupRuleRemoteTypeSG:
			sg, err := s.IBMVPCClient.GetSecurityGroupByName(*remote.SecurityGroupName)
			if err != nil {
				return fmt.Errorf("failed to find VPC security group by name '%s', err: %w", *remote.SecurityGroupName, err)
			}
			if sg == nil {
				return fmt.Errorf("VPC security group by name '%s' does not exist", *remote.SecurityGroupName)
			}
			s.V(3).Info("Creating VPC security group rule", "securityGroupID", *securityGroupID, "direction", *direction, "protocol", *protocol, "securityGroup", *remote.SecurityGroupName, "securityGroupCRN", *sg.CRN)
			remoteOption.CRN = sg.CRN
		default:
			s.V(3).Info("Creating VPC security group rule", "securityGroupID", *securityGroupID, "direction", *direction, "protocol", *protocol, "cidr", "0.0.0.0/0")
			remoteOption.CIDRBlock = ptr.To("0.0.0.0/0")
		}

		return nil
	}

	remoteOption := &vpcv1.SecurityGroupRuleRemotePrototype{}
	if err := setRemote(remote, remoteOption); err != nil {
		return nil, fmt.Errorf("failed to set remote option while creating VPC security group rule: %w", err)
	}

	options := vpcv1.CreateSecurityGroupRuleOptions{
		SecurityGroupID: securityGroupID,
	}

	options.SetSecurityGroupRulePrototype(&vpcv1.SecurityGroupRulePrototype{
		Direction: direction,
		Protocol:  protocol,
		PortMin:   portMin,
		PortMax:   portMax,
		Remote:    remoteOption,
	})

	var ruleID *string
	ruleIntf, _, err := s.IBMVPCClient.CreateSecurityGroupRule(&options)
	if err != nil {
		return nil, err
	}

	switch reflect.TypeOf(ruleIntf).String() {
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		rule := ruleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
		ruleID = rule.ID
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		rule := ruleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
		ruleID = rule.ID
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		rule := ruleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
		ruleID = rule.ID
	}
	s.Info("Created VPC security group rule", "ruleID", *ruleID)
	return ruleID, nil
}

// createVPCSecurityGroupRules creates rules for a security group.
func (s *PowerVSClusterScope) createVPCSecurityGroupRules(ogSecurityGroupRules []*infrav1beta2.VPCSecurityGroupRule, securityGroupID *string) ([]*string, error) {
	ruleIDs := []*string{}
	s.V(3).Info("Creating VPC security group rules")

	for _, rule := range ogSecurityGroupRules {
		var protocol *string
		var portMax, portMin *int64

		direction := ptr.To(string(rule.Direction))
		switch rule.Direction {
		case infrav1beta2.VPCSecurityGroupRuleDirectionInbound:
			protocol = ptr.To(string(rule.Source.Protocol))
			if rule.Source.PortRange != nil {
				portMin = ptr.To(rule.Source.PortRange.MinimumPort)
				portMax = ptr.To(rule.Source.PortRange.MaximumPort)
			}

			for _, remote := range rule.Source.Remotes {
				id, err := s.createVPCSecurityGroupRule(securityGroupID, direction, protocol, portMin, portMax, remote)
				if err != nil {
					return nil, fmt.Errorf("failed to create VPC security group rule: %v", err)
				}
				ruleIDs = append(ruleIDs, id)
			}
		case infrav1beta2.VPCSecurityGroupRuleDirectionOutbound:
			protocol = ptr.To(string(rule.Destination.Protocol))
			if rule.Destination.PortRange != nil {
				portMin = ptr.To(rule.Destination.PortRange.MinimumPort)
				portMax = ptr.To(rule.Destination.PortRange.MaximumPort)
			}

			for _, remote := range rule.Destination.Remotes {
				id, err := s.createVPCSecurityGroupRule(securityGroupID, direction, protocol, portMin, portMax, remote)
				if err != nil {
					return nil, fmt.Errorf("failed to create VPC security group rule: %v", err)
				}
				ruleIDs = append(ruleIDs, id)
			}
		}
	}

	return ruleIDs, nil
}

// createVPCSecurityGroupRulesAndSetStatus creates VPC security group rules and sets its status.
func (s *PowerVSClusterScope) createVPCSecurityGroupRulesAndSetStatus(ogSecurityGroupRules []*infrav1beta2.VPCSecurityGroupRule, securityGroupID, securityGroupName *string) error {
	ruleIDs, err := s.createVPCSecurityGroupRules(ogSecurityGroupRules, securityGroupID)
	if err != nil {
		return fmt.Errorf("failed to create VPC security group rules: %w", err)
	}
	s.Info("VPC security group rules created", "security group name", *securityGroupName)

	s.SetVPCSecurityGroupStatus(*securityGroupName, infrav1beta2.VPCSecurityGroupStatus{
		ID:                securityGroupID,
		RuleIDs:           ruleIDs,
		ControllerCreated: ptr.To(true),
	})

	return nil
}

// createVPCSecurityGroup creates a VPC security group.
func (s *PowerVSClusterScope) createVPCSecurityGroup(spec infrav1beta2.VPCSecurityGroup) (*string, error) {
	s.V(3).Info("Creating VPC security group", "name", *spec.Name)

	options := &vpcv1.CreateSecurityGroupOptions{
		VPC: &vpcv1.VPCIdentity{
			ID: s.GetVPCID(),
		},
		Name: spec.Name,
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: ptr.To(s.GetResourceGroupID()),
		},
	}

	securityGroup, _, err := s.IBMVPCClient.CreateSecurityGroup(options)
	if err != nil {
		return nil, err
	}
	// To-Do: Add tags to VPC security group, need to implement the client for "github.com/IBM/platform-services-go-sdk/globaltaggingv1".
	return securityGroup.ID, nil
}

// validateVPCSecurityGroupRuleRemote compares a specific security group rule's remote with the spec and existing security group rule's remote.
func (s *PowerVSClusterScope) validateVPCSecurityGroupRuleRemote(originalSGRemote *vpcv1.SecurityGroupRuleRemote, expectedSGRemote infrav1beta2.VPCSecurityGroupRuleRemote) (bool, error) {
	var match bool

	switch expectedSGRemote.RemoteType {
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAny:
		if originalSGRemote.CIDRBlock != nil && *originalSGRemote.CIDRBlock == "0.0.0.0/0" {
			match = true
		}
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAddress:
		if originalSGRemote.Address != nil && *originalSGRemote.Address == *expectedSGRemote.Address {
			match = true
		}
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeCIDR:
		cidrSubnet, err := s.IBMVPCClient.GetVPCSubnetByName(*expectedSGRemote.CIDRSubnetName)
		if err != nil {
			return false, fmt.Errorf("failed to find VPC subnet by name '%s' for fetching CIDR block: %w", *expectedSGRemote.CIDRSubnetName, err)
		}

		if originalSGRemote.CIDRBlock != nil && cidrSubnet != nil && *originalSGRemote.CIDRBlock == *cidrSubnet.Ipv4CIDRBlock {
			match = true
		}
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeSG:
		securityGroup, err := s.IBMVPCClient.GetSecurityGroupByName(*expectedSGRemote.SecurityGroupName)
		if err != nil {
			return false, fmt.Errorf("failed to find ID for resource group '%s': %w", *expectedSGRemote.SecurityGroupName, err)
		}

		if originalSGRemote.CRN != nil && securityGroup.Name != nil && *originalSGRemote.CRN == *securityGroup.CRN {
			match = true
		}
	}

	return match, nil
}

// validateSecurityGroupRule compares a specific security group's rule with the spec and existing security group's rule.
func (s *PowerVSClusterScope) validateSecurityGroupRule(originalSecurityGroupRules []vpcv1.SecurityGroupRuleIntf, direction infrav1beta2.VPCSecurityGroupRuleDirection, rule *infrav1beta2.VPCSecurityGroupRulePrototype, remote infrav1beta2.VPCSecurityGroupRuleRemote) (ruleID *string, match bool, err error) {
	updateError := func(e error) {
		err = fmt.Errorf("failed to validate VPC security group rule's remote: %w", e)
	}

	protocol := string(rule.Protocol)

	for _, ogRuleIntf := range originalSecurityGroupRules {
		switch reflect.TypeOf(ogRuleIntf).String() {
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
			ogRule := ogRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
			ruleID = ogRule.ID

			if *ogRule.Direction == string(direction) && *ogRule.Protocol == protocol {
				ogRemote := ogRule.Remote.(*vpcv1.SecurityGroupRuleRemote)
				match, err = s.validateVPCSecurityGroupRuleRemote(ogRemote, remote)
				if err != nil {
					updateError(err)
					return nil, false, err
				}
			}
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
			portMin := rule.PortRange.MinimumPort
			portMax := rule.PortRange.MaximumPort
			ogRule := ogRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			ruleID = ogRule.ID

			if *ogRule.Direction == string(direction) && *ogRule.Protocol == protocol && *ogRule.PortMax == portMax && *ogRule.PortMin == portMin {
				ogRemote := ogRule.Remote.(*vpcv1.SecurityGroupRuleRemote)
				match, err = s.validateVPCSecurityGroupRuleRemote(ogRemote, remote)
				if err != nil {
					updateError(err)
					return nil, false, err
				}
			}
		case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
			icmpCode := rule.ICMPCode
			icmpType := rule.ICMPType
			ogRule := ogRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
			ruleID = ogRule.ID

			if *ogRule.Direction == string(direction) && *ogRule.Protocol == protocol && *ogRule.Code == *icmpCode && *ogRule.Type == *icmpType {
				ogRemote := ogRule.Remote.(*vpcv1.SecurityGroupRuleRemote)
				match, err = s.validateVPCSecurityGroupRuleRemote(ogRemote, remote)
				if err != nil {
					updateError(err)
					return nil, false, err
				}
			}
		}
		if match {
			return ruleID, match, nil
		}
	}

	return nil, false, nil
}

// validateVPCSecurityGroupRules compares a specific security group rules spec with the existing security group's rules.
func (s *PowerVSClusterScope) validateVPCSecurityGroupRules(originalSecurityGroupRules []vpcv1.SecurityGroupRuleIntf, expectedSecurityGroupRules []*infrav1beta2.VPCSecurityGroupRule) ([]*string, bool, error) {
	ruleIDs := []*string{}
	for _, expectedRule := range expectedSecurityGroupRules {
		direction := expectedRule.Direction

		switch direction {
		case infrav1beta2.VPCSecurityGroupRuleDirectionInbound:
			for _, remote := range expectedRule.Source.Remotes {
				id, match, err := s.validateSecurityGroupRule(originalSecurityGroupRules, direction, expectedRule.Source, remote)
				if err != nil {
					return nil, false, fmt.Errorf("failed to validate VPC security group rule: %w", err)
				}
				if !match {
					return nil, false, nil
				}
				ruleIDs = append(ruleIDs, id)
			}
		case infrav1beta2.VPCSecurityGroupRuleDirectionOutbound:
			for _, remote := range expectedRule.Destination.Remotes {
				id, match, err := s.validateSecurityGroupRule(originalSecurityGroupRules, direction, expectedRule.Destination, remote)
				if err != nil {
					return nil, false, fmt.Errorf("failed to validate VPC security group rule: %v", err)
				}
				if !match {
					return nil, false, nil
				}
				ruleIDs = append(ruleIDs, id)
			}
		}
	}

	return ruleIDs, true, nil
}

// validateVPCSecurityGroup validates the security group and it's rules provided by user via spec.
func (s *PowerVSClusterScope) validateVPCSecurityGroup(securityGroup infrav1beta2.VPCSecurityGroup) (*vpcv1.SecurityGroup, []*string, error) {
	var securityGroupDet *vpcv1.SecurityGroup
	var err error

	if securityGroup.ID != nil {
		securityGroupDet, _, err = s.IBMVPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{
			ID: securityGroup.ID,
		})
		if err != nil {
			return nil, nil, err
		}
		if securityGroupDet == nil {
			return nil, nil, fmt.Errorf("failed to find VPC security group with provided ID '%v'", securityGroup.ID)
		}
	} else {
		securityGroupDet, err = s.IBMVPCClient.GetSecurityGroupByName(*securityGroup.Name)
		if err != nil {
			if _, ok := err.(*vpc.SecurityGroupByNameNotFound); !ok {
				return nil, nil, err
			}
		}
		if securityGroupDet == nil {
			return nil, nil, nil
		}
	}
	if securityGroupDet.VPC == nil || securityGroupDet.VPC.ID == nil || *securityGroupDet.VPC.ID != *s.GetVPCID() {
		return nil, nil, fmt.Errorf("VPC security group by name exists but is not attached to VPC")
	}

	ruleIDs, ok, err := s.validateVPCSecurityGroupRules(securityGroupDet.Rules, securityGroup.Rules)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate VPC security group rules: %v", err)
	}
	if !ok {
		if _, _, controllerCreated := s.GetVPCSecurityGroupByName(*securityGroup.Name); controllerCreated != nil && !*controllerCreated {
			return nil, nil, fmt.Errorf("VPC security group by name exists but rules are not matching")
		}
		return nil, nil, s.createVPCSecurityGroupRulesAndSetStatus(securityGroup.Rules, securityGroupDet.ID, securityGroupDet.Name)
	}

	return securityGroupDet, ruleIDs, nil
}

// ReconcileTransitGateway reconcile transit gateway.
func (s *PowerVSClusterScope) ReconcileTransitGateway() (bool, error) {
	if s.GetTransitGatewayID() != nil {
		s.V(3).Info("Transit gateway ID is set, fetching details", "tgID", s.GetTransitGatewayID())
		tg, _, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
			ID: s.GetTransitGatewayID(),
		})
		if err != nil {
			return false, err
		}
		requeue, err := s.checkAndUpdateTransitGateway(tg)
		if err != nil {
			return false, err
		}
		return requeue, nil
	}

	// check transit gateway exist in cloud
	tg, err := s.isTransitGatewayExists()
	if err != nil {
		return false, err
	}

	// check the status and update the transit gateway's connections if they are not proper
	if tg != nil {
		requeue, err := s.checkAndUpdateTransitGateway(tg)
		if err != nil {
			return false, err
		}
		return requeue, nil
	}

	// create transit gateway
	s.V(3).Info("Creating transit gateway")
	if err := s.createTransitGateway(); err != nil {
		return false, fmt.Errorf("failed to create transit gateway: %v", err)
	}

	return true, nil
}

// isTransitGatewayExists checks transit gateway exist in cloud.
func (s *PowerVSClusterScope) isTransitGatewayExists() (*tgapiv1.TransitGateway, error) {
	// TODO(karthik-k-n): Support regex
	var transitGateway *tgapiv1.TransitGateway
	var err error

	if s.IBMPowerVSCluster.Spec.TransitGateway != nil && s.IBMPowerVSCluster.Spec.TransitGateway.ID != nil {
		transitGateway, _, err = s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
			ID: s.IBMPowerVSCluster.Spec.TransitGateway.ID,
		})
	} else {
		transitGateway, err = s.TransitGatewayClient.GetTransitGatewayByName(*s.GetServiceName(infrav1beta2.ResourceTypeTransitGateway))
	}

	if err != nil {
		return nil, err
	}

	if transitGateway == nil || transitGateway.ID == nil {
		s.V(3).Info("Transit gateway not found in IBM Cloud")
		return nil, nil
	}

	s.SetTransitGatewayStatus(transitGateway.ID, ptr.To(false))

	return transitGateway, nil
}

// checkAndUpdateTransitGateway checks given transit gateway's status and its connections.
// if update is set to true, it updates the transit gateway connections too if it is not exist already.
func (s *PowerVSClusterScope) checkAndUpdateTransitGateway(transitGateway *tgapiv1.TransitGateway) (bool, error) {
	requeue, err := s.checkTransitGatewayStatus(transitGateway)
	if err != nil {
		return false, err
	}
	if requeue {
		return requeue, nil
	}

	return s.checkAndUpdateTransitGatewayConnections(transitGateway)
}

// checkTransitGatewayStatus checks the state of a transit gateway.
// If state is pending, true is returned indicating a requeue for reconciliation.
// In all other cases, it returns false.
func (s *PowerVSClusterScope) checkTransitGatewayStatus(tg *tgapiv1.TransitGateway) (bool, error) {
	s.V(3).Info("Checking the status of transit gateway", "name", *tg.Name)
	switch *tg.Status {
	case string(infrav1beta2.TransitGatewayStateAvailable):
		s.V(3).Info("Transit gateway is in available state")
	case string(infrav1beta2.TransitGatewayStateFailed):
		return false, fmt.Errorf("failed to create transit gateway, current status: %s", *tg.Status)
	case string(infrav1beta2.TransitGatewayStatePending):
		s.V(3).Info("Transit gateway is in pending state")
		return true, nil
	}

	return false, nil
}

// checkAndUpdateTransitGatewayConnections checks given transit gateway's connections status.
// it also creates the transit gateway connections if it is not exist already.
func (s *PowerVSClusterScope) checkAndUpdateTransitGatewayConnections(transitGateway *tgapiv1.TransitGateway) (bool, error) {
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

	pvsServiceInstanceCRN, err := s.fetchPowerVSServiceInstanceCRN()
	if err != nil {
		return false, fmt.Errorf("failed to fetch PowerVS service instance CRN: %w", err)
	}

	if len(tgConnections.Connections) == 0 {
		s.V(3).Info("Connections not exist on transit gateway, creating them")
		if err := s.createTransitGatewayConnections(transitGateway, pvsServiceInstanceCRN, vpcCRN); err != nil {
			return false, err
		}

		return true, nil
	}

	requeue, powerVSConnStatus, vpcConnStatus, err := s.validateTransitGatewayConnections(tgConnections.Connections, vpcCRN, pvsServiceInstanceCRN)
	if err != nil {
		return false, err
	} else if requeue {
		return requeue, nil
	}

	// return when connections are in attached state.
	if powerVSConnStatus && vpcConnStatus {
		return false, nil
	}

	// update the connections when connection not exist
	if !powerVSConnStatus {
		s.V(3).Info("Only PowerVS connection not exist in transit gateway, creating it")
		if err := s.createTransitGatewayConnection(transitGateway.ID, ptr.To(getTGPowerVSConnectionName(*transitGateway.Name)), pvsServiceInstanceCRN, powervsNetworkConnectionType); err != nil {
			return false, err
		}
	}

	if !vpcConnStatus {
		s.V(3).Info("Only VPC connection not exist in transit gateway, creating it")
		if err := s.createTransitGatewayConnection(transitGateway.ID, ptr.To(getTGVPCConnectionName(*transitGateway.Name)), vpcCRN, vpcNetworkConnectionType); err != nil {
			return false, err
		}
	}

	return true, nil
}

// validateTransitGatewayConnections validates the existing transit gateway connections.
// to avoid returning many return values, connection id will be returned and considered that connection is in attached state.
func (s *PowerVSClusterScope) validateTransitGatewayConnections(connections []tgapiv1.TransitGatewayConnectionCust, vpcCRN, pvsServiceInstanceCRN *string) (bool, bool, bool, error) {
	var powerVSConnStatus, vpcConnStatus bool
	for _, conn := range connections {
		if *conn.NetworkType == string(vpcNetworkConnectionType) && *conn.NetworkID == *vpcCRN {
			if requeue, err := s.checkTransitGatewayConnectionStatus(conn); err != nil {
				return requeue, false, false, err
			} else if requeue {
				return requeue, false, false, nil
			}

			if s.IBMPowerVSCluster.Status.TransitGateway != nil && s.IBMPowerVSCluster.Status.TransitGateway.VPCConnection == nil {
				s.SetTransitGatewayConnectionStatus(vpcNetworkConnectionType, &infrav1beta2.ResourceReference{ID: conn.ID, ControllerCreated: ptr.To(false)})
			}
			vpcConnStatus = true
		}
		if *conn.NetworkType == string(powervsNetworkConnectionType) && *conn.NetworkID == *pvsServiceInstanceCRN {
			if requeue, err := s.checkTransitGatewayConnectionStatus(conn); err != nil {
				return requeue, false, false, err
			} else if requeue {
				return requeue, false, false, nil
			}

			if s.IBMPowerVSCluster.Status.TransitGateway != nil && s.IBMPowerVSCluster.Status.TransitGateway.PowerVSConnection == nil {
				s.SetTransitGatewayConnectionStatus(powervsNetworkConnectionType, &infrav1beta2.ResourceReference{ID: conn.ID, ControllerCreated: ptr.To(false)})
			}
			powerVSConnStatus = true
		}
	}

	return false, powerVSConnStatus, vpcConnStatus, nil
}

// checkTransitGatewayConnectionStatus checks the state of a transit gateway connection.
// If state is pending, true is returned indicating a requeue for reconciliation.
// In all other cases, it returns false.
func (s *PowerVSClusterScope) checkTransitGatewayConnectionStatus(con tgapiv1.TransitGatewayConnectionCust) (bool, error) {
	s.V(3).Info("Checking the status of transit gateway connection", "name", *con.Name)
	switch *con.Status {
	case string(infrav1beta2.TransitGatewayConnectionStateAttached):
		return false, nil
	case string(infrav1beta2.TransitGatewayConnectionStateFailed):
		return false, fmt.Errorf("failed to attach connection to transit gateway, current status: %s", *con.Status)
	case string(infrav1beta2.TransitGatewayConnectionStatePending):
		s.V(3).Info("Transit gateway connection is in pending state")
		return true, nil
	}
	return false, nil
}

// createTransitGatewayConnection creates transit gateway connection and sets the connection status.
func (s *PowerVSClusterScope) createTransitGatewayConnection(transitGatewayID, connName, networkID *string, networkType networkConnectionType) error {
	s.V(3).Info("Creating transit gateway connection", "tgID", transitGatewayID, "type", networkType, "name", connName)
	conn, _, err := s.TransitGatewayClient.CreateTransitGatewayConnection(&tgapiv1.CreateTransitGatewayConnectionOptions{
		TransitGatewayID: transitGatewayID,
		NetworkType:      ptr.To(string(networkType)),
		NetworkID:        networkID,
		Name:             connName,
	})
	if err != nil {
		return err
	}
	s.SetTransitGatewayConnectionStatus(networkType, &infrav1beta2.ResourceReference{ID: conn.ID, ControllerCreated: ptr.To(true)})

	return nil
}

// createTransitGatewayConnections creates PowerVS and VPC connections in the transit gateway.
func (s *PowerVSClusterScope) createTransitGatewayConnections(tg *tgapiv1.TransitGateway, pvsServiceInstanceCRN, vpcCRN *string) error {
	if err := s.createTransitGatewayConnection(tg.ID, ptr.To(getTGPowerVSConnectionName(*tg.Name)), pvsServiceInstanceCRN, powervsNetworkConnectionType); err != nil {
		return fmt.Errorf("failed to create PowerVS connection in transit gateway: %w", err)
	}

	if err := s.createTransitGatewayConnection(tg.ID, ptr.To(getTGVPCConnectionName(*tg.Name)), vpcCRN, vpcNetworkConnectionType); err != nil {
		return fmt.Errorf("failed to create VPC connection in transit gateway: %w", err)
	}

	return nil
}

// createTransitGateway creates transit gateway and sets the transit gateway status.
func (s *PowerVSClusterScope) createTransitGateway() error {
	// TODO(karthik-k-n): Verify that the supplied zone supports PER
	// TODO(karthik-k-n): consider moving to clusterscope

	// fetch resource group id
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroup())
	}

	if s.IBMPowerVSCluster.Status.ServiceInstance == nil || s.IBMPowerVSCluster.Status.VPC == nil {
		return fmt.Errorf("failed to proeceed with transit gateway creation as either one of VPC or PowerVS service instance reconciliation is not successful")
	}

	location, globalRouting, err := genUtil.GetTransitGatewayLocationAndRouting(s.Zone(), s.VPC().Region)
	if err != nil {
		return fmt.Errorf("failed to get transit gateway location and routing: %w", err)
	}

	// throw error when user tries to use local routing where global routing is required.
	// TODO: Add a webhook validation for below condition.
	if s.IBMPowerVSCluster.Spec.TransitGateway != nil && s.IBMPowerVSCluster.Spec.TransitGateway.GlobalRouting != nil && !*s.IBMPowerVSCluster.Spec.TransitGateway.GlobalRouting && *globalRouting {
		return fmt.Errorf("failed to use local routing for transit gateway since powervs and vpc are in different region and requires global routing")
	}
	// setting global routing to true when it is set by user.
	if s.IBMPowerVSCluster.Spec.TransitGateway != nil && s.IBMPowerVSCluster.Spec.TransitGateway.GlobalRouting != nil && *s.IBMPowerVSCluster.Spec.TransitGateway.GlobalRouting {
		globalRouting = ptr.To(true)
	}

	tgName := s.GetServiceName(infrav1beta2.ResourceTypeTransitGateway)
	tg, _, err := s.TransitGatewayClient.CreateTransitGateway(&tgapiv1.CreateTransitGatewayOptions{
		Location:      location,
		Name:          tgName,
		Global:        globalRouting,
		ResourceGroup: &tgapiv1.ResourceGroupIdentity{ID: ptr.To(resourceGroupID)},
	})
	if err != nil {
		return err
	}

	s.SetTransitGatewayStatus(tg.ID, ptr.To(true))

	vpcCRN, err := s.fetchVPCCRN()
	if err != nil {
		return fmt.Errorf("failed to fetch VPC CRN: %w", err)
	}
	pvsServiceInstanceCRN, err := s.fetchPowerVSServiceInstanceCRN()
	if err != nil {
		return fmt.Errorf("failed to fetch PowerVS service instance CRN: %w", err)
	}

	if err := s.createTransitGatewayConnections(tg, pvsServiceInstanceCRN, vpcCRN); err != nil {
		return err
	}

	return nil
}

// ReconcileLoadBalancers reconcile loadBalancer.
func (s *PowerVSClusterScope) ReconcileLoadBalancers() (bool, error) {
	loadBalancers := make([]infrav1beta2.VPCLoadBalancerSpec, 0)
	if len(s.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		loadBalancer := infrav1beta2.VPCLoadBalancerSpec{
			Name:   *s.GetServiceName(infrav1beta2.ResourceTypeLoadBalancer),
			Public: ptr.To(true),
		}
		loadBalancers = append(loadBalancers, loadBalancer)
	} else {
		loadBalancers = append(loadBalancers, s.IBMPowerVSCluster.Spec.LoadBalancers...)
	}

	isAnyLoadBalancerNotReady := false

	for index, loadBalancer := range loadBalancers {
		var loadBalancerID *string
		if loadBalancer.ID != nil {
			loadBalancerID = loadBalancer.ID
		} else {
			if loadBalancer.Name == "" {
				loadBalancer.Name = fmt.Sprintf("%s-%d", *s.GetServiceName(infrav1beta2.ResourceTypeLoadBalancer), index)
			}
			loadBalancerID = s.GetLoadBalancerID(loadBalancer.Name)
		}
		if loadBalancerID != nil {
			s.V(3).Info("LoadBalancer ID is set, fetching loadbalancer details", "loadbalancerid", *loadBalancerID)
			loadBalancer, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: loadBalancerID,
			})
			if err != nil {
				return false, err
			}

			if isReady := s.checkLoadBalancerStatus(*loadBalancer); !isReady {
				s.V(3).Info("LoadBalancer is still not Active", "name", *loadBalancer.Name, "state", *loadBalancer.ProvisioningStatus)
				isAnyLoadBalancerNotReady = true
			}

			loadBalancerStatus := infrav1beta2.VPCLoadBalancerStatus{
				ID:       loadBalancer.ID,
				State:    infrav1beta2.VPCLoadBalancerState(*loadBalancer.ProvisioningStatus),
				Hostname: loadBalancer.Hostname,
			}
			s.SetLoadBalancerStatus(*loadBalancer.Name, loadBalancerStatus)
			continue
		}

		// check VPC load balancer exist in cloud
		loadBalancerStatus, err := s.checkLoadBalancer(loadBalancer)
		if err != nil {
			return false, err
		}
		if loadBalancerStatus != nil {
			s.V(3).Info("Found VPC load balancer in IBM Cloud", "loadBalancerID", *loadBalancerStatus.ID)
			s.SetLoadBalancerStatus(loadBalancer.Name, *loadBalancerStatus)
			continue
		}

		// check loadbalancer port against apiserver port.
		err = s.checkLoadBalancerPort(loadBalancer)
		if err != nil {
			return false, err
		}

		// create loadBalancer
		s.V(3).Info("Creating VPC load balancer")
		loadBalancerStatus, err = s.createLoadBalancer(loadBalancer)
		if err != nil {
			return false, fmt.Errorf("failed to create VPC load balancer: %w", err)
		}
		s.Info("Created VPC load balancer", "loadBalancerID", loadBalancerStatus.ID)
		s.SetLoadBalancerStatus(loadBalancer.Name, *loadBalancerStatus)
		isAnyLoadBalancerNotReady = true
	}
	if isAnyLoadBalancerNotReady {
		return false, nil
	}
	return true, nil
}

// checkLoadBalancerStatus checks the state of a VPC load balancer.
// If state is active, true is returned, in all other cases, it returns false indicating that load balancer is still not ready.
func (s *PowerVSClusterScope) checkLoadBalancerStatus(lb vpcv1.LoadBalancer) bool {
	s.V(3).Info("Checking the status of VPC load balancer", "name", *lb.Name)
	switch *lb.ProvisioningStatus {
	case string(infrav1beta2.VPCLoadBalancerStateActive):
		s.V(3).Info("VPC load balancer is in active state")
		return true
	case string(infrav1beta2.VPCLoadBalancerStateCreatePending):
		s.V(3).Info("VPC load balancer creation is in pending state")
	case string(infrav1beta2.VPCLoadBalancerStateUpdatePending):
		s.V(3).Info("VPC load balancer is in updating state")
	}
	return false
}

func (s *PowerVSClusterScope) checkLoadBalancerPort(lb infrav1beta2.VPCLoadBalancerSpec) error {
	for _, listerner := range lb.AdditionalListeners {
		if listerner.Port == int64(s.APIServerPort()) {
			return fmt.Errorf("port %d for the %s load balancer cannot be used as an additional listener port, as it is already assigned to the API server", listerner.Port, lb.Name)
		}
	}
	return nil
}

// checkLoadBalancer checks if VPC load balancer by the given name exists in cloud.
func (s *PowerVSClusterScope) checkLoadBalancer(lb infrav1beta2.VPCLoadBalancerSpec) (*infrav1beta2.VPCLoadBalancerStatus, error) {
	loadBalancer, err := s.IBMVPCClient.GetLoadBalancerByName(lb.Name)
	if err != nil {
		return nil, err
	}
	if loadBalancer == nil {
		s.V(3).Info("VPC load balancer not found in IBM Cloud")
		return nil, nil
	}
	return &infrav1beta2.VPCLoadBalancerStatus{
		ID:       loadBalancer.ID,
		State:    infrav1beta2.VPCLoadBalancerState(*loadBalancer.ProvisioningStatus),
		Hostname: loadBalancer.Hostname,
	}, nil
}

// createLoadBalancer creates loadBalancer.
func (s *PowerVSClusterScope) createLoadBalancer(lb infrav1beta2.VPCLoadBalancerSpec) (*infrav1beta2.VPCLoadBalancerStatus, error) {
	options := &vpcv1.CreateLoadBalancerOptions{}
	// TODO(karthik-k-n): consider moving resource group id to clusterscope
	// fetch resource group id
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroup())
	}

	var isPublic bool
	if lb.Public != nil && *lb.Public {
		isPublic = true
	}
	options.SetIsPublic(isPublic)
	options.SetName(lb.Name)
	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &resourceGroupID,
	})

	subnetIDs := s.GetVPCSubnetIDs()
	if subnetIDs == nil {
		return nil, fmt.Errorf("no subnets are present for load balancer creation")
	}
	for _, subnetID := range subnetIDs {
		subnet := &vpcv1.SubnetIdentity{
			ID: subnetID,
		}
		options.Subnets = append(options.Subnets, subnet)
	}
	options.SetPools([]vpcv1.LoadBalancerPoolPrototype{
		{
			Algorithm:     core.StringPtr("round_robin"),
			HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{Delay: core.Int64Ptr(5), MaxRetries: core.Int64Ptr(2), Timeout: core.Int64Ptr(2), Type: core.StringPtr("tcp")},
			// Note: Appending port number to the name, it will be referenced to set target port while adding new pool member
			Name:     core.StringPtr(fmt.Sprintf("%s-pool-%d", lb.Name, s.APIServerPort())),
			Protocol: core.StringPtr("tcp"),
		},
	})

	options.SetListeners([]vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{
		{
			Protocol: core.StringPtr("tcp"),
			Port:     core.Int64Ptr(int64(s.APIServerPort())),
			DefaultPool: &vpcv1.LoadBalancerPoolIdentityByName{
				Name: core.StringPtr(fmt.Sprintf("%s-pool-%d", lb.Name, s.APIServerPort())),
			},
		},
	})

	for _, additionalListeners := range lb.AdditionalListeners {
		pool := vpcv1.LoadBalancerPoolPrototype{
			Algorithm:     core.StringPtr("round_robin"),
			HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{Delay: core.Int64Ptr(5), MaxRetries: core.Int64Ptr(2), Timeout: core.Int64Ptr(2), Type: core.StringPtr("tcp")},
			// Note: Appending port number to the name, it will be referenced to set target port while adding new pool member
			Name:     ptr.To(fmt.Sprintf("additional-pool-%d", additionalListeners.Port)),
			Protocol: core.StringPtr("tcp"),
		}
		options.Pools = append(options.Pools, pool)

		listener := vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{
			Protocol: core.StringPtr("tcp"),
			Port:     core.Int64Ptr(additionalListeners.Port),
			DefaultPool: &vpcv1.LoadBalancerPoolIdentityByName{
				Name: ptr.To(fmt.Sprintf("additional-pool-%d", additionalListeners.Port)),
			},
		}
		options.Listeners = append(options.Listeners, listener)
	}

	loadBalancer, _, err := s.IBMVPCClient.CreateLoadBalancer(options)
	if err != nil {
		return nil, err
	}
	lbState := infrav1beta2.VPCLoadBalancerState(*loadBalancer.ProvisioningStatus)
	return &infrav1beta2.VPCLoadBalancerStatus{
		ID:                loadBalancer.ID,
		State:             lbState,
		Hostname:          loadBalancer.Hostname,
		ControllerCreated: ptr.To(true),
	}, nil
}

// COSInstance returns the COS instance reference.
func (s *PowerVSClusterScope) COSInstance() *infrav1beta2.CosInstance {
	return s.IBMPowerVSCluster.Spec.CosInstance
}

// ReconcileCOSInstance reconcile COS bucket.
func (s *PowerVSClusterScope) ReconcileCOSInstance() error {
	// check COS service instance exist in cloud
	cosServiceInstanceStatus, err := s.checkCOSServiceInstance()
	if err != nil {
		return err
	}
	if cosServiceInstanceStatus != nil {
		s.V(3).Info("COS service instance found in IBM Cloud")
		s.SetStatus(infrav1beta2.ResourceTypeCOSInstance, infrav1beta2.ResourceReference{ID: cosServiceInstanceStatus.GUID, ControllerCreated: ptr.To(false)})
	} else {
		// create COS service instance
		s.V(3).Info("Creating COS service instance")
		cosServiceInstanceStatus, err = s.createCOSServiceInstance()
		if err != nil {
			s.Error(err, "failed to create COS service instance")
			return err
		}
		s.Info("Created COS service instance", "cosID", cosServiceInstanceStatus.GUID)
		s.SetStatus(infrav1beta2.ResourceTypeCOSInstance, infrav1beta2.ResourceReference{ID: cosServiceInstanceStatus.GUID, ControllerCreated: ptr.To(true)})
	}

	props, err := authenticator.GetProperties()
	if err != nil {
		s.Error(err, "failed to fetch service properties")
		return err
	}

	apiKey, ok := props["APIKEY"]
	if !ok {
		return fmt.Errorf("IBM Cloud API key is not provided, set %s environmental variable", "IBMCLOUD_API_KEY")
	}

	region := s.bucketRegion()
	if region == "" {
		return fmt.Errorf("failed to determine COS bucket region, both bucket region and VPC region not set")
	}

	serviceEndpoint := fmt.Sprintf("s3.%s.%s", region, cosURLDomain)
	// Fetch the COS service endpoint.
	cosServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.COS), s.ServiceEndpoint)
	if cosServiceEndpoint != "" {
		s.Logger.V(3).Info("Overriding the default COS endpoint", "cosEndpoint", cosServiceEndpoint)
		serviceEndpoint = cosServiceEndpoint
	}

	cosOptions := cos.ServiceOptions{
		Options: &cosSession.Options{
			Config: aws.Config{
				Endpoint: &serviceEndpoint,
				Region:   &region,
			},
		},
	}

	cosClient, err := cos.NewService(cosOptions, apiKey, *cosServiceInstanceStatus.GUID)
	if err != nil {
		return fmt.Errorf("failed to create COS client: %w", err)
	}
	s.COSClient = cosClient

	// check bucket exist in service instance
	if exist, err := s.checkCOSBucket(); exist {
		s.V(3).Info("COS bucket found in IBM Cloud")
		return nil
	} else if err != nil {
		s.Error(err, "failed to check COS bucket")
		return err
	}

	// create bucket in service instance
	if err := s.createCOSBucket(); err != nil {
		return err
	}
	return nil
}

func (s *PowerVSClusterScope) checkCOSBucket() (bool, error) {
	if _, err := s.COSClient.GetBucketByName(*s.GetServiceName(infrav1beta2.ResourceTypeCOSBucket)); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket, "Forbidden", "NotFound":
				// If the bucket doesn't exist that's ok, we'll try to create it
				return false, nil
			default:
				return false, err
			}
		} else {
			return false, err
		}
	}
	return true, nil
}

func (s *PowerVSClusterScope) createCOSBucket() error {
	input := &s3.CreateBucketInput{
		Bucket: ptr.To(*s.GetServiceName(infrav1beta2.ResourceTypeCOSBucket)),
	}
	_, err := s.COSClient.CreateBucket(input)
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return fmt.Errorf("failed to create COS bucket %w", err)
	}

	switch aerr.Code() {
	// If bucket already exists, all good.
	case s3.ErrCodeBucketAlreadyOwnedByYou:
		return nil
	case s3.ErrCodeBucketAlreadyExists:
		return nil
	default:
		return fmt.Errorf("failed to create COS bucket %w", err)
	}
}

func (s *PowerVSClusterScope) checkCOSServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	// check cos service instance
	serviceInstance, err := s.ResourceClient.GetInstanceByName(*s.GetServiceName(infrav1beta2.ResourceTypeCOSInstance), resourcecontroller.CosResourceID, resourcecontroller.CosResourcePlanID)
	if err != nil {
		return nil, err
	}
	if serviceInstance == nil {
		s.V(3).Info("COS service instance is not found in IBM Cloud", "name", *s.GetServiceName(infrav1beta2.ResourceTypeCOSInstance))
		return nil, nil
	}
	if *serviceInstance.State != string(infrav1beta2.ServiceInstanceStateActive) {
		return nil, fmt.Errorf("COS service instance is not in active state, current state: %s", *serviceInstance.State)
	}
	return serviceInstance, nil
}

func (s *PowerVSClusterScope) createCOSServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	// fetch resource group id.
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		return nil, fmt.Errorf("failed to fetch resource group ID for resource group %v, ID is empty", s.ResourceGroup())
	}

	target := "Global"
	// create service instance
	serviceInstance, _, err := s.ResourceClient.CreateResourceInstance(&resourcecontrollerv2.CreateResourceInstanceOptions{
		Name:           s.GetServiceName(infrav1beta2.ResourceTypeCOSInstance),
		Target:         &target,
		ResourceGroup:  &resourceGroupID,
		ResourcePlanID: ptr.To(resourcecontroller.CosResourcePlanID),
	})
	if err != nil {
		return nil, err
	}
	return serviceInstance, nil
}

// fetchResourceGroupID retrieving id of resource group.
func (s *PowerVSClusterScope) fetchResourceGroupID() (string, error) {
	if s.ResourceGroup() == nil || s.ResourceGroup().Name == nil {
		return "", fmt.Errorf("resource group name is not set")
	}

	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return "", err
	}

	account, err := utils.GetAccount(auth)
	if err != nil {
		return "", err
	}

	resourceGroup := s.ResourceGroup().Name
	rmv2ListResourceGroupOpt := resourcemanagerv2.ListResourceGroupsOptions{Name: resourceGroup, AccountID: &account}
	resourceGroupListResult, _, err := s.ResourceManagerClient.ListResourceGroups(&rmv2ListResourceGroupOpt)
	if err != nil {
		return "", err
	}

	if resourceGroupListResult != nil && len(resourceGroupListResult.Resources) > 0 {
		rg := resourceGroupListResult.Resources[0]
		resourceGroupID := *rg.ID
		return resourceGroupID, nil
	}

	err = fmt.Errorf("could not retrieve resource group ID for %s", *resourceGroup)
	return "", err
}

// fetchVPCCRN returns VPC CRN.
func (s *PowerVSClusterScope) fetchVPCCRN() (*string, error) {
	vpcID := s.GetVPCID()
	if vpcID == nil {
		if s.IBMPowerVSCluster.Spec.VPC != nil && s.IBMPowerVSCluster.Spec.VPC.ID != nil {
			vpcID = s.IBMPowerVSCluster.Spec.VPC.ID
		}
	}
	vpcDetails, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
		ID: vpcID,
	})
	if err != nil {
		return nil, err
	}
	return vpcDetails.CRN, nil
}

// fetchPowerVSServiceInstanceCRN returns Power VS service instance CRN.
func (s *PowerVSClusterScope) fetchPowerVSServiceInstanceCRN() (*string, error) {
	serviceInstanceID := s.GetServiceInstanceID()
	if serviceInstanceID == "" {
		if s.IBMPowerVSCluster.Spec.ServiceInstanceID != "" {
			serviceInstanceID = s.IBMPowerVSCluster.Spec.ServiceInstanceID
		} else if s.IBMPowerVSCluster.Spec.ServiceInstance != nil && s.IBMPowerVSCluster.Spec.ServiceInstance.ID != nil {
			serviceInstanceID = *s.IBMPowerVSCluster.Spec.ServiceInstance.ID
		}
	}
	pvsDetails, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &serviceInstanceID,
	})
	if err != nil {
		return nil, err
	}
	return pvsDetails.CRN, nil
}

// TODO(karthik-k-n): Decide on proper naming format for services.

// GetServiceName returns name of given service type from spec or generate a name for it.
func (s *PowerVSClusterScope) GetServiceName(resourceType infrav1beta2.ResourceType) *string { //nolint:gocyclo
	switch resourceType {
	case infrav1beta2.ResourceTypeServiceInstance:
		if s.ServiceInstance() == nil || s.ServiceInstance().Name == nil {
			return ptr.To(fmt.Sprintf("%s-serviceInstance", s.InfraCluster()))
		}
		return s.ServiceInstance().Name
	case infrav1beta2.ResourceTypeNetwork:
		if s.Network() == nil || s.Network().Name == nil {
			return ptr.To(fmt.Sprintf("DHCPSERVER%s_Private", s.InfraCluster()))
		}
		return s.Network().Name
	case infrav1beta2.ResourceTypeVPC:
		if s.VPC() == nil || s.VPC().Name == nil {
			return ptr.To(fmt.Sprintf("%s-vpc", s.InfraCluster()))
		}
		return s.VPC().Name
	case infrav1beta2.ResourceTypeTransitGateway:
		if s.TransitGateway() == nil || s.TransitGateway().Name == nil {
			return ptr.To(fmt.Sprintf("%s-transitgateway", s.InfraCluster()))
		}
		return s.TransitGateway().Name
	case infrav1beta2.ResourceTypeDHCPServer:
		if s.DHCPServer() == nil || s.DHCPServer().Name == nil {
			return ptr.To(s.InfraCluster())
		}
		return s.DHCPServer().Name
	case infrav1beta2.ResourceTypeCOSInstance:
		if s.COSInstance() == nil || s.COSInstance().Name == "" {
			return ptr.To(fmt.Sprintf("%s-cosinstance", s.InfraCluster()))
		}
		return &s.COSInstance().Name
	case infrav1beta2.ResourceTypeCOSBucket:
		if s.COSInstance() == nil || s.COSInstance().BucketName == "" {
			return ptr.To(fmt.Sprintf("%s-cosbucket", s.InfraCluster()))
		}
		return &s.COSInstance().BucketName
	case infrav1beta2.ResourceTypeSubnet:
		return ptr.To(fmt.Sprintf("%s-vpcsubnet", s.InfraCluster()))
	case infrav1beta2.ResourceTypeLoadBalancer:
		return ptr.To(fmt.Sprintf("%s-loadbalancer", s.InfraCluster()))
	}
	return nil
}

// DeleteLoadBalancer deletes loadBalancer.
func (s *PowerVSClusterScope) DeleteLoadBalancer() (bool, error) {
	errs := []error{}
	requeue := false
	for _, lb := range s.IBMPowerVSCluster.Status.LoadBalancers {
		if lb.ID == nil || lb.ControllerCreated == nil || !*lb.ControllerCreated {
			s.Info("Skipping VPC load balancer deletion as resource is not created by controller")
			continue
		}

		lb, resp, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
			ID: lb.ID,
		})

		if err != nil {
			if resp != nil && resp.StatusCode == ResourceNotFoundCode {
				s.Info("VPC load balancer successfully deleted")
				continue
			}
			errs = append(errs, fmt.Errorf("failed to fetch VPC load balancer: %w", err))
			continue
		}

		if lb != nil && lb.ProvisioningStatus != nil && *lb.ProvisioningStatus == string(infrav1beta2.VPCLoadBalancerStateDeletePending) {
			s.V(3).Info("VPC load balancer is currently being deleted")
			return true, nil
		}

		if _, err = s.IBMVPCClient.DeleteLoadBalancer(&vpcv1.DeleteLoadBalancerOptions{
			ID: lb.ID,
		}); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete VPC load balancer: %w", err))
			continue
		}
		requeue = true
	}
	if len(errs) > 0 {
		return false, kerrors.NewAggregate(errs)
	}
	return requeue, nil
}

// DeleteVPCSecurityGroups deletes VPC security group.
func (s *PowerVSClusterScope) DeleteVPCSecurityGroups() error {
	for _, securityGroup := range s.IBMPowerVSCluster.Status.VPCSecurityGroups {
		if securityGroup.ControllerCreated == nil || !*securityGroup.ControllerCreated {
			s.Info("Skipping VPC security group deletion as resource is not created by controller", "ID", *securityGroup.ID)
			continue
		}
		if _, resp, err := s.IBMVPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{
			ID: securityGroup.ID,
		}); err != nil {
			if resp != nil && resp.StatusCode == ResourceNotFoundCode {
				s.Info("VPC security group has been already deleted", "securityGroupID", *securityGroup.ID)
				continue
			}
			return fmt.Errorf("failed to fetch VPC security group '%s': %w", *securityGroup.ID, err)
		}

		s.V(3).Info("Deleting VPC security group", "securityGroupID", *securityGroup.ID)
		options := &vpcv1.DeleteSecurityGroupOptions{
			ID: securityGroup.ID,
		}
		if _, err := s.IBMVPCClient.DeleteSecurityGroup(options); err != nil {
			return fmt.Errorf("failed to delete VPC security group '%s': %w", *securityGroup.ID, err)
		}
		s.Info("VPC security group successfully deleted", "securityGroupID", *securityGroup.ID)
	}
	return nil
}

// DeleteVPCSubnet deletes VPC subnet.
func (s *PowerVSClusterScope) DeleteVPCSubnet() (bool, error) {
	errs := []error{}
	requeue := false
	for _, subnet := range s.IBMPowerVSCluster.Status.VPCSubnet {
		if subnet.ID == nil || subnet.ControllerCreated == nil || !*subnet.ControllerCreated {
			s.Info("Skipping VPC subnet deletion as resource is not created by controller")
			continue
		}

		net, resp, err := s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{
			ID: subnet.ID,
		})

		if err != nil {
			if resp != nil && resp.StatusCode == ResourceNotFoundCode {
				s.Info("VPC subnet successfully deleted")
				continue
			}
			errs = append(errs, fmt.Errorf("failed to fetch VPC subnet: %w", err))
			continue
		}

		if net != nil && net.Status != nil && *net.Status == string(infrav1beta2.VPCSubnetStateDeleting) {
			return true, nil
		}

		if _, err = s.IBMVPCClient.DeleteSubnet(&vpcv1.DeleteSubnetOptions{
			ID: net.ID,
		}); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete VPC subnet: %w", err))
			continue
		}
		requeue = true
	}
	if len(errs) > 0 {
		return false, kerrors.NewAggregate(errs)
	}
	return requeue, nil
}

// DeleteVPC deletes VPC.
func (s *PowerVSClusterScope) DeleteVPC() (bool, error) {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeVPC) {
		s.Info("Skipping VPC deletion as resource is not created by controller")
		return false, nil
	}

	if s.IBMPowerVSCluster.Status.VPC.ID == nil {
		return false, nil
	}

	vpc, resp, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
		ID: s.IBMPowerVSCluster.Status.VPC.ID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == ResourceNotFoundCode {
			s.Info("VPC successfully deleted")
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch VPC: %w", err)
	}

	if vpc != nil && vpc.Status != nil && *vpc.Status == string(infrav1beta2.VPCStateDeleting) {
		return true, nil
	}

	if _, err = s.IBMVPCClient.DeleteVPC(&vpcv1.DeleteVPCOptions{
		ID: vpc.ID,
	}); err != nil {
		return false, fmt.Errorf("failed to delete VPC: %w", err)
	}
	return true, nil
}

// DeleteTransitGateway deletes transit gateway.
func (s *PowerVSClusterScope) DeleteTransitGateway() (bool, error) {
	skipTGDeletion := false
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeTransitGateway) {
		s.Info("Skipping transit gateway deletion as resource is not created by controller, but will check if connections are created by the controller.")
		skipTGDeletion = true
	}

	if s.IBMPowerVSCluster.Status.TransitGateway == nil {
		return false, nil
	}

	tg, resp, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
		ID: s.IBMPowerVSCluster.Status.TransitGateway.ID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == ResourceNotFoundCode {
			s.Info("Transit gateway successfully deleted")
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch transit gateway: %w", err)
	}

	if tg.Status != nil && *tg.Status == string(infrav1beta2.TransitGatewayStateDeletePending) {
		s.V(3).Info("Transit gateway is being deleted")
		return true, nil
	}

	requeue, err := s.deleteTransitGatewayConnections(tg)
	if err != nil {
		return false, err
	} else if requeue {
		return true, nil
	}

	if skipTGDeletion {
		return false, nil
	}

	if _, err = s.TransitGatewayClient.DeleteTransitGateway(&tgapiv1.DeleteTransitGatewayOptions{
		ID: s.IBMPowerVSCluster.Status.TransitGateway.ID,
	}); err != nil {
		return false, fmt.Errorf("failed to delete transit gateway: %w", err)
	}
	return true, nil
}

func (s *PowerVSClusterScope) deleteTransitGatewayConnections(tg *tgapiv1.TransitGateway) (bool, error) {
	deleteConnection := func(connID *string) (bool, error) {
		conn, resp, err := s.TransitGatewayClient.GetTransitGatewayConnection(&tgapiv1.GetTransitGatewayConnectionOptions{
			TransitGatewayID: tg.ID,
			ID:               connID,
		})
		if resp.StatusCode == ResourceNotFoundCode {
			s.V(3).Info("Connection deleted in transit gateway", "connectionID", *connID)
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("failed to get transit gateway powervs connection: %w", err)
		}
		if conn.Status != nil && *conn.Status == string(infrav1beta2.TransitGatewayConnectionStateDeleting) {
			s.V(3).Info("Transit gateway connection is in deleting state")
			return true, nil
		}

		if _, err = s.TransitGatewayClient.DeleteTransitGatewayConnection(&tgapiv1.DeleteTransitGatewayConnectionOptions{
			ID:               connID,
			TransitGatewayID: tg.ID,
		}); err != nil {
			return false, fmt.Errorf("failed to delete transit gateway connection: %w", err)
		}

		return true, nil
	}
	if *s.IBMPowerVSCluster.Status.TransitGateway.PowerVSConnection.ControllerCreated {
		s.V(3).Info("Deleting PowerVS connection in Transit gateway")
		requeue, err := deleteConnection(s.IBMPowerVSCluster.Status.TransitGateway.PowerVSConnection.ID)
		if err != nil {
			return false, err
		}
		if requeue {
			return requeue, nil
		}
	}

	if *s.IBMPowerVSCluster.Status.TransitGateway.VPCConnection.ControllerCreated {
		s.V(3).Info("Deleting VPC connection in Transit gateway")
		requeue, err := deleteConnection(s.IBMPowerVSCluster.Status.TransitGateway.VPCConnection.ID)
		if err != nil {
			return false, err
		}
		if requeue {
			return requeue, nil
		}
	}

	return false, nil
}

// DeleteDHCPServer deletes DHCP server.
func (s *PowerVSClusterScope) DeleteDHCPServer() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeDHCPServer) {
		s.Info("Skipping DHP server deletion as resource is not created by controller")
		return nil
	}
	if s.isResourceCreatedByController(infrav1beta2.ResourceTypeServiceInstance) {
		s.Info("Skipping DHCP server deletion as PowerVS service instance is created by controller, will directly delete the PowerVS service instance since it will delete the DHCP server internally")
		return nil
	}

	if s.IBMPowerVSCluster.Status.DHCPServer.ID == nil {
		return nil
	}

	server, err := s.IBMPowerVSClient.GetDHCPServer(*s.IBMPowerVSCluster.Status.DHCPServer.ID)
	if err != nil {
		if strings.Contains(err.Error(), string(DHCPServerNotFound)) {
			s.Info("DHCP server successfully deleted")
			return nil
		}
		return fmt.Errorf("failed to fetch DHCP server: %w", err)
	}

	if err = s.IBMPowerVSClient.DeleteDHCPServer(*server.ID); err != nil {
		return fmt.Errorf("failed to delete DHCP server: %w", err)
	}
	return nil
}

// DeleteServiceInstance deletes service instance.
func (s *PowerVSClusterScope) DeleteServiceInstance() (bool, error) {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeServiceInstance) {
		s.Info("Skipping PowerVS service instance deletion as resource is not created by controller")
		return false, nil
	}

	if s.IBMPowerVSCluster.Status.ServiceInstance.ID == nil {
		return false, nil
	}

	serviceInstance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: s.IBMPowerVSCluster.Status.ServiceInstance.ID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to fetch PowerVS service instance: %w", err)
	}

	if serviceInstance != nil && *serviceInstance.State == string(infrav1beta2.ServiceInstanceStateRemoved) {
		s.Info("PowerVS service instance has been removed")
		return false, nil
	}

	if _, err = s.ResourceClient.DeleteResourceInstance(&resourcecontrollerv2.DeleteResourceInstanceOptions{
		ID: serviceInstance.ID,
	}); err != nil {
		s.Error(err, "failed to delete Power VS service instance")
		return false, err
	}

	return true, nil
}

// DeleteCOSInstance deletes COS instance.
func (s *PowerVSClusterScope) DeleteCOSInstance() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeCOSInstance) {
		s.Info("Skipping COS instance deletion as resource is not created by controller")
		return nil
	}

	if s.IBMPowerVSCluster.Status.COSInstance.ID == nil {
		return nil
	}

	cosInstance, resp, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: s.IBMPowerVSCluster.Status.COSInstance.ID,
	})
	if err != nil {
		if resp != nil && resp.StatusCode == ResourceNotFoundCode {
			return nil
		}
		return fmt.Errorf("failed to fetch COS service instance: %w", err)
	}

	if cosInstance != nil && (*cosInstance.State == "pending_reclamation" || *cosInstance.State == string(infrav1beta2.ServiceInstanceStateRemoved)) {
		s.Info("COS service instance has been removed")
		return nil
	}

	if _, err = s.ResourceClient.DeleteResourceInstance(&resourcecontrollerv2.DeleteResourceInstanceOptions{
		ID:        cosInstance.ID,
		Recursive: ptr.To(true),
	}); err != nil {
		s.Error(err, "failed to delete COS service instance")
		return err
	}
	s.Info("COS service instance successfully deleted")
	return nil
}

// resourceCreatedByController helps to identify resource created by controller or not.
func (s *PowerVSClusterScope) isResourceCreatedByController(resourceType infrav1beta2.ResourceType) bool { //nolint:gocyclo
	switch resourceType {
	case infrav1beta2.ResourceTypeVPC:
		vpcStatus := s.IBMPowerVSCluster.Status.VPC
		if vpcStatus == nil || vpcStatus.ControllerCreated == nil || !*vpcStatus.ControllerCreated {
			return false
		}
		return true
	case infrav1beta2.ResourceTypeServiceInstance:
		serviceInstance := s.IBMPowerVSCluster.Status.ServiceInstance
		if serviceInstance == nil || serviceInstance.ControllerCreated == nil || !*serviceInstance.ControllerCreated {
			return false
		}
		return true
	case infrav1beta2.ResourceTypeTransitGateway:
		transitGateway := s.IBMPowerVSCluster.Status.TransitGateway
		if transitGateway == nil || transitGateway.ControllerCreated == nil || !*transitGateway.ControllerCreated {
			return false
		}
		return true
	case infrav1beta2.ResourceTypeDHCPServer:
		dhcpServer := s.IBMPowerVSCluster.Status.DHCPServer
		if dhcpServer == nil || dhcpServer.ControllerCreated == nil || !*dhcpServer.ControllerCreated {
			return false
		}
		return true
	case infrav1beta2.ResourceTypeCOSInstance:
		cosInstance := s.IBMPowerVSCluster.Status.COSInstance
		if cosInstance == nil || cosInstance.ControllerCreated == nil || !*cosInstance.ControllerCreated {
			return false
		}
		return true
	}
	return false
}

// TODO: duplicate function, optimize it.
func (s *PowerVSClusterScope) bucketRegion() string {
	if s.COSInstance() != nil && s.COSInstance().BucketRegion != "" {
		return s.COSInstance().BucketRegion
	}
	// if the bucket region is not set, use vpc region
	vpcDetails := s.VPC()
	if vpcDetails != nil && vpcDetails.Region != nil {
		return *vpcDetails.Region
	}
	return ""
}
