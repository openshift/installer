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
	"strings"

	"github.com/go-logr/logr"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/datacenters"
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

	"k8s.io/klog/v2/textlogger"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/cos"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
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
const powerEdgeRouter = "power-edge-router"

// PowerVSClusterScopeParams defines the input parameters used to create a new PowerVSClusterScope.
type PowerVSClusterScopeParams struct {
	Client            client.Client
	Logger            logr.Logger
	Cluster           *capiv1beta1.Cluster
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
}

// PowerVSClusterScope defines a scope defined around a Power VS Cluster.
type PowerVSClusterScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper
	session     *ibmpisession.IBMPISession

	IBMPowerVSClient     powervs.PowerVS
	IBMVPCClient         vpc.Vpc
	TransitGatewayClient transitgateway.TransitGateway
	ResourceClient       resourcecontroller.ResourceController
	COSClient            cos.Cos

	Cluster           *capiv1beta1.Cluster
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
}

// NewPowerVSClusterScope creates a new PowerVSClusterScope from the supplied parameters.
func NewPowerVSClusterScope(params PowerVSClusterScopeParams) (*PowerVSClusterScope, error) { //nolint:gocyclo
	if params.Client == nil {
		err := errors.New("error failed to generate new scope from nil Client")
		return nil, err
	}
	if params.Cluster == nil {
		err := errors.New("error failed to generate new scope from nil Cluster")
		return nil, err
	}
	if params.IBMPowerVSCluster == nil {
		err := errors.New("error failed to generate new scope from nil IBMPowerVSCluster")
		return nil, err
	}
	if params.Logger == (logr.Logger{}) {
		params.Logger = textlogger.NewLogger(textlogger.NewConfig())
	}

	helper, err := patch.NewHelper(params.IBMPowerVSCluster, params.Client)
	if err != nil {
		err = fmt.Errorf("error failed to init patch helper: %w", err)
		return nil, err
	}

	options := powervs.ServiceOptions{
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
			if err := rc.SetServiceURL(rcEndpoint); err != nil {
				params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
				return nil, fmt.Errorf("failed to set resource controller endpoint: %w", err)
			}
		}

		res, _, err := rc.GetResourceInstance(
			&resourcecontrollerv2.GetResourceInstanceOptions{
				ID: core.StringPtr(params.IBMPowerVSCluster.Spec.ServiceInstanceID),
			})
		if err != nil {
			err = fmt.Errorf("failed to get resource instance: %w", err)
			return nil, err
		}
		options.Zone = *res.RegionID
		options.CloudInstanceID = params.IBMPowerVSCluster.Spec.ServiceInstanceID
	} else {
		options.Zone = *params.IBMPowerVSCluster.Spec.Zone
	}

	// Fetch the PowerVS service endpoint.
	powerVSServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.PowerVS), params.ServiceEndpoint)
	if powerVSServiceEndpoint != "" {
		params.Logger.V(3).Info("Overriding the default PowerVS endpoint", "powerVSEndpoint", powerVSServiceEndpoint)
		options.IBMPIOptions.URL = powerVSServiceEndpoint
	}

	// TODO(karhtik-k-n): may be optimize NewService to use the session created here
	powerVSClient, err := powervs.NewService(options)
	if err != nil {
		return nil, fmt.Errorf("error failed to create power vs client %w", err)
	}

	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("error failed to create authenticator %w", err)
	}
	account, err := utils.GetAccount(auth)
	if err != nil {
		return nil, fmt.Errorf("error failed to get account details %w", err)
	}

	sessionOptions := &ibmpisession.IBMPIOptions{
		Authenticator: auth,
		UserAccount:   account,
		Zone:          options.Zone,
		Debug:         params.Logger.V(DEBUGLEVEL).Enabled(),
	}
	if powerVSServiceEndpoint != "" {
		sessionOptions.URL = powerVSServiceEndpoint
	}
	session, err := ibmpisession.NewIBMPISession(sessionOptions)
	if err != nil {
		return nil, fmt.Errorf("error failed to get power vs session %w", err)
	}

	// if powervs.cluster.x-k8s.io/create-infra=true annotation is not set, create only powerVSClient.
	if !genUtil.CheckCreateInfraAnnotation(*params.IBMPowerVSCluster) {
		return &PowerVSClusterScope{
			session:           session,
			Logger:            params.Logger,
			Client:            params.Client,
			patchHelper:       helper,
			Cluster:           params.Cluster,
			IBMPowerVSCluster: params.IBMPowerVSCluster,
			ServiceEndpoint:   params.ServiceEndpoint,
			IBMPowerVSClient:  powerVSClient,
		}, nil
	}

	// if powervs.cluster.x-k8s.io/create-infra=true annotation is set, create necessary clients.
	if params.IBMPowerVSCluster.Spec.VPC == nil || params.IBMPowerVSCluster.Spec.VPC.Region == nil {
		return nil, fmt.Errorf("error failed to generate vpc client as VPC info is nil")
	}

	if params.Logger.V(DEBUGLEVEL).Enabled() {
		core.SetLoggingLevel(core.LevelDebug)
	}

	svcEndpoint := endpoints.FetchVPCEndpoint(*params.IBMPowerVSCluster.Spec.VPC.Region, params.ServiceEndpoint)
	vpcClient, err := vpc.NewService(svcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error failed to create IBM VPC client: %w", err)
	}

	// Create TransitGateway client
	tgOptions := &tgapiv1.TransitGatewayApisV1Options{
		Authenticator: auth,
	}
	// Fetch the TransitGateway service endpoint.
	tgServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.TransitGateway), params.ServiceEndpoint)
	if tgServiceEndpoint != "" {
		params.Logger.V(3).Info("Overriding the default TransitGateway endpoint", "transitGatewayEndpoint", tgServiceEndpoint)
		tgOptions.URL = tgServiceEndpoint
	}

	tgClient, err := transitgateway.NewService(tgOptions)
	if err != nil {
		return nil, fmt.Errorf("error failed to create tranist gateway client: %w", err)
	}

	// Create Resource Controller client.
	serviceOption := resourcecontroller.ServiceOptions{
		ResourceControllerV2Options: &resourcecontrollerv2.ResourceControllerV2Options{
			Authenticator: auth,
		},
	}
	// Fetch the resource controller endpoint.
	rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), params.ServiceEndpoint)
	if rcEndpoint != "" {
		serviceOption.URL = rcEndpoint
		params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
	}
	resourceClient, err := resourcecontroller.NewService(serviceOption)
	if err != nil {
		return nil, fmt.Errorf("error failed to create resource client: %w", err)
	}

	clusterScope := &PowerVSClusterScope{
		session:              session,
		Logger:               params.Logger,
		Client:               params.Client,
		patchHelper:          helper,
		Cluster:              params.Cluster,
		IBMPowerVSCluster:    params.IBMPowerVSCluster,
		ServiceEndpoint:      params.ServiceEndpoint,
		IBMPowerVSClient:     powerVSClient,
		IBMVPCClient:         vpcClient,
		TransitGatewayClient: tgClient,
		ResourceClient:       resourceClient,
	}
	return clusterScope, nil
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

// GetServiceInstanceID get the service instance id.
func (s *PowerVSClusterScope) GetServiceInstanceID() string {
	if s.IBMPowerVSCluster.Spec.ServiceInstanceID != "" {
		return s.IBMPowerVSCluster.Spec.ServiceInstanceID
	}
	if s.IBMPowerVSCluster.Spec.ServiceInstance != nil && s.IBMPowerVSCluster.Spec.ServiceInstance.ID != nil {
		return *s.IBMPowerVSCluster.Spec.ServiceInstance.ID
	}
	if s.IBMPowerVSCluster.Status.ServiceInstance != nil && s.IBMPowerVSCluster.Status.ServiceInstance.ID != nil {
		return *s.IBMPowerVSCluster.Status.ServiceInstance.ID
	}
	return ""
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
	case infrav1beta2.ResourceTypeTransitGateway:
		if s.IBMPowerVSCluster.Status.TransitGateway == nil {
			s.IBMPowerVSCluster.Status.TransitGateway = &resource
			return
		}
		s.IBMPowerVSCluster.Status.TransitGateway.Set(resource)
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

// Network returns the cluster Network.
func (s *PowerVSClusterScope) Network() *infrav1beta2.IBMPowerVSResourceReference {
	return &s.IBMPowerVSCluster.Spec.Network
}

// GetDHCPServerID returns the DHCP id from spec or status of IBMPowerVSCluster object.
func (s *PowerVSClusterScope) GetDHCPServerID() *string {
	if s.IBMPowerVSCluster.Spec.DHCPServer != nil && s.IBMPowerVSCluster.Spec.DHCPServer.ID != nil {
		return s.IBMPowerVSCluster.Spec.DHCPServer.ID
	}
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

// GetVPCID returns the VPC id.
func (s *PowerVSClusterScope) GetVPCID() *string {
	if s.IBMPowerVSCluster.Spec.VPC != nil && s.IBMPowerVSCluster.Spec.VPC.ID != nil {
		return s.IBMPowerVSCluster.Spec.VPC.ID
	}
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
	if s.IBMPowerVSCluster.Status.VPCSubnet == nil {
		return nil
	}
	subnets := []*string{}
	for _, subnet := range s.IBMPowerVSCluster.Status.VPCSubnet {
		subnets = append(subnets, subnet.ID)
	}
	return subnets
}

// SetVPCSubnetID set the VPC subnet id.
func (s *PowerVSClusterScope) SetVPCSubnetID(name string, resource infrav1beta2.ResourceReference) {
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

// TransitGateway returns the cluster Transit Gateway information.
func (s *PowerVSClusterScope) TransitGateway() *infrav1beta2.TransitGateway {
	return s.IBMPowerVSCluster.Spec.TransitGateway
}

// GetTransitGatewayID returns the transit gateway id.
func (s *PowerVSClusterScope) GetTransitGatewayID() *string {
	if s.IBMPowerVSCluster.Spec.TransitGateway != nil && s.IBMPowerVSCluster.Spec.TransitGateway.ID != nil {
		return s.IBMPowerVSCluster.Spec.TransitGateway.ID
	}
	if s.IBMPowerVSCluster.Status.TransitGateway != nil {
		return s.IBMPowerVSCluster.Status.TransitGateway.ID
	}
	return nil
}

// PublicLoadBalancer returns the cluster public loadBalancer information.
func (s *PowerVSClusterScope) PublicLoadBalancer() *infrav1beta2.VPCLoadBalancerSpec {
	// if the user did not specify any loadbalancer then return the public loadbalancer created by the controller.
	if len(s.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		return &infrav1beta2.VPCLoadBalancerSpec{
			Name:   *s.GetServiceName(infrav1beta2.ResourceTypeLoadBalancer),
			Public: ptr.To(true),
		}
	}
	for _, lb := range s.IBMPowerVSCluster.Spec.LoadBalancers {
		if lb.Public != nil && *lb.Public {
			return &lb
		}
	}
	return nil
}

// SetLoadBalancerStatus set the loadBalancer id.
func (s *PowerVSClusterScope) SetLoadBalancerStatus(name string, loadBalancer infrav1beta2.VPCLoadBalancerStatus) {
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

// GetLoadBalancerHostName will return the hostname of load balancer.
func (s *PowerVSClusterScope) GetLoadBalancerHostName(name string) *string {
	if s.IBMPowerVSCluster.Status.LoadBalancers == nil {
		return nil
	}
	if val, ok := s.IBMPowerVSCluster.Status.LoadBalancers[name]; ok {
		return val.Hostname
	}
	return nil
}

// GetResourceGroupID returns the resource group id if it present under spec or statue filed of IBMPowerVSCluster object
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
		return fmt.Errorf("powervs zone is not set")
	}
	// fetch the datacenter capabilities for zone.
	// though the function name is WithDatacenterRegion it takes zone as parameter
	params := datacenters.NewV1DatacentersGetParamsWithContext(context.TODO()).WithDatacenterRegion(*zone)
	datacenter, err := s.session.Power.Datacenters.V1DatacentersGet(params)
	if err != nil {
		return fmt.Errorf("failed to get datacenter details for zone: %s err:%w", *zone, err)
	}
	if datacenter == nil || datacenter.Payload == nil || datacenter.Payload.Capabilities == nil {
		return fmt.Errorf("failed to get datacenter capabilities for zone: %s", *zone)
	}
	// check for the PER support in datacenter capabilities.
	perAvailable, ok := datacenter.Payload.Capabilities[powerEdgeRouter]
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
	s.Info("Fetched resource group id from cloud", "resourceGroupID", resourceGroupID)
	// Set the status of IBMPowerVSCluster object with resource group id.
	s.SetStatus(infrav1beta2.ResourceTypeResourceGroup, infrav1beta2.ResourceReference{ID: &resourceGroupID, ControllerCreated: ptr.To(false)})
	return nil
}

// ReconcilePowerVSServiceInstance reconciles Power VS service instance.
func (s *PowerVSClusterScope) ReconcilePowerVSServiceInstance() error {
	// Verify if service instance id is set in spec or status field of IBMPowerVSCluster object.
	serviceInstanceID := s.GetServiceInstanceID()
	if serviceInstanceID != "" {
		// if serviceInstanceID is set, verify that it exist and in active state.
		s.Info("Service instance id is set", "id", serviceInstanceID)
		serviceInstance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: &serviceInstanceID,
		})
		if err != nil {
			return err
		}
		if serviceInstance == nil {
			return fmt.Errorf("error failed to get service instance with id %s", serviceInstanceID)
		}
		if *serviceInstance.State != string(infrav1beta2.ServiceInstanceStateActive) {
			return fmt.Errorf("service instance not in active state, current state: %s", *serviceInstance.State)
		}
		s.Info("Found service instance and its in active state", "id", serviceInstanceID)
		return nil
	}

	// check PowerVS service instance exist in cloud, if it does not exist proceed with creating the instance.
	serviceInstanceID, err := s.isServiceInstanceExists()
	if err != nil {
		return err
	}
	// Set the status of IBMPowerVSCluster object with serviceInstanceID and ControllerCreated to false as PowerVS service instance is already exist in cloud.
	if serviceInstanceID != "" {
		s.SetStatus(infrav1beta2.ResourceTypeServiceInstance, infrav1beta2.ResourceReference{ID: &serviceInstanceID, ControllerCreated: ptr.To(false)})
		return nil
	}

	// create PowerVS Service Instance
	serviceInstance, err := s.createServiceInstance()
	if err != nil {
		return err
	}
	// Set the status of IBMPowerVSCluster object with serviceInstanceID and ControllerCreated to true as new PowerVS service instance is created.
	s.SetStatus(infrav1beta2.ResourceTypeServiceInstance, infrav1beta2.ResourceReference{ID: serviceInstance.GUID, ControllerCreated: ptr.To(true)})
	return nil
}

// checkServiceInstance checks PowerVS service instance exist in cloud.
func (s *PowerVSClusterScope) isServiceInstanceExists() (string, error) {
	s.Info("Checking for service instance in cloud")
	// Fetches service instance by name.
	serviceInstance, err := s.getServiceInstance()
	if err != nil {
		s.Error(err, "failed to get service instance")
		return "", err
	}
	if serviceInstance == nil {
		s.Info("Not able to find service instance", "service instance", s.IBMPowerVSCluster.Spec.ServiceInstance)
		return "", nil
	}
	if *serviceInstance.State != string(infrav1beta2.ServiceInstanceStateActive) {
		s.Info("Service instance not in active state", "service instance", s.IBMPowerVSCluster.Spec.ServiceInstance, "state", *serviceInstance.State)
		return "", fmt.Errorf("service instance not in active state, current state: %s", *serviceInstance.State)
	}
	s.Info("Service instance found and its in active state", "id", *serviceInstance.GUID)
	return *serviceInstance.GUID, nil
}

// getServiceInstance return resource instance by name.
func (s *PowerVSClusterScope) getServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	//TODO: Support regular expression
	return s.ResourceClient.GetServiceInstance("", *s.GetServiceName(infrav1beta2.ResourceTypeServiceInstance))
}

// createServiceInstance creates the service instance.
func (s *PowerVSClusterScope) createServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	// fetch resource group id.
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		s.Info("failed to create service instance, failed to fetch resource group id")
		return nil, fmt.Errorf("error getting resource group id for resource group %v, id is empty", s.ResourceGroup())
	}

	// create service instance.
	s.Info("Creating new service instance", "name", s.GetServiceName(infrav1beta2.ResourceTypeServiceInstance))
	zone := s.Zone()
	if zone == nil {
		return nil, fmt.Errorf("error creating new service instance, PowerVS zone is not set")
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
	s.Info("Created new service instance")
	return serviceInstance, nil
}

// ReconcileNetwork reconciles network.
func (s *PowerVSClusterScope) ReconcileNetwork() error {
	if s.GetDHCPServerID() != nil {
		s.Info("DHCP server id is set")
		if err := s.isDHCPServerActive(); err != nil {
			return err
		}
		// if dhcp server exist and in active state, its assumed that dhcp network exist
		// TODO(Phase 2): Verify that dhcp network is exist.
		return nil
		//	TODO(karthik-k-n): If needed set dhcp status here
	}
	// check network exist in cloud
	networkID, err := s.checkNetwork()
	if err != nil {
		return err
	}
	if networkID != nil {
		s.Info("Found network", "id", networkID)
		s.SetStatus(infrav1beta2.ResourceTypeNetwork, infrav1beta2.ResourceReference{ID: networkID, ControllerCreated: ptr.To(false)})
		return nil
	}

	s.Info("Creating DHCP server")
	dhcpServer, err := s.createDHCPServer()
	if err != nil {
		s.Error(err, "Error creating DHCP server")
		return err
	}
	if dhcpServer != nil {
		s.Info("Created DHCP Server", "id", *dhcpServer)
		s.SetStatus(infrav1beta2.ResourceTypeDHCPServer, infrav1beta2.ResourceReference{ID: dhcpServer, ControllerCreated: ptr.To(true)})
		return nil
	}
	return nil
}

// checkNetwork checks the network exist in cloud.
func (s *PowerVSClusterScope) checkNetwork() (*string, error) {
	// get network from cloud.
	networkID, err := s.getNetwork()
	if err != nil {
		s.Error(err, "failed to get network")
		return nil, err
	}
	if networkID == nil {
		s.Info("Not able to find network", "network", s.IBMPowerVSCluster.Spec.Network)
		return nil, nil
	}
	return networkID, nil
}

func (s *PowerVSClusterScope) getNetwork() (*string, error) {
	// fetch the network associated with network id
	if s.IBMPowerVSCluster.Spec.Network.ID != nil {
		network, err := s.IBMPowerVSClient.GetNetworkByID(*s.IBMPowerVSCluster.Spec.Network.ID)
		if err != nil {
			return nil, err
		}
		return network.NetworkID, nil
	}

	// if the user has provided the already existing dhcp server name then there might exist network name
	// with format DHCPSERVER<DHCPServer.Name>_Private , try fetching that
	var networkName string
	if s.DHCPServer() != nil && s.DHCPServer().Name != nil {
		networkName = fmt.Sprintf("DHCPSERVER%s_Private", *s.DHCPServer().Name)
	} else {
		networkName = *s.GetServiceName(infrav1beta2.ResourceTypeNetwork)
	}

	// fetch the network associated with name
	network, err := s.IBMPowerVSClient.GetNetworkByName(networkName)
	if err != nil {
		return nil, err
	}
	if network == nil {
		s.Info("network does not exist", "name", networkName)
		return nil, nil
	}
	return network.NetworkID, nil
	//TODO: Support regular expression
}

// isDHCPServerActive checks if the DHCP server status is active.
func (s *PowerVSClusterScope) isDHCPServerActive() error {
	dhcpID := *s.GetDHCPServerID()
	dhcpServer, err := s.IBMPowerVSClient.GetDHCPServer(dhcpID)
	if err != nil {
		return err
	}

	if *dhcpServer.Status != string(infrav1beta2.DHCPServerStateActive) {
		return fmt.Errorf("error dhcp server state is not active, current state %s", *dhcpServer.Status)
	}
	s.Info("DHCP server is found and its in active state")
	return nil
}

// createDHCPServer creates the DHCP server.
func (s *PowerVSClusterScope) createDHCPServer() (*string, error) {
	var dhcpServerCreateParams models.DHCPServerCreate
	dhcpServerDetails := s.DHCPServer()
	if dhcpServerDetails == nil {
		dhcpServerDetails = &infrav1beta2.DHCPServer{}
	}

	dhcpServerCreateParams.Name = s.GetServiceName(infrav1beta2.ResourceTypeDHCPServer)
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
		return nil, fmt.Errorf("created dhcp server is nil")
	}
	if dhcpServer.Network == nil {
		return nil, fmt.Errorf("created dhcp server network is nil")
	}

	s.Info("DHCP Server network details", "details", *dhcpServer.Network)
	s.SetStatus(infrav1beta2.ResourceTypeNetwork, infrav1beta2.ResourceReference{ID: dhcpServer.Network.ID, ControllerCreated: ptr.To(true)})
	return dhcpServer.ID, nil
}

// ReconcileVPC reconciles VPC.
func (s *PowerVSClusterScope) ReconcileVPC() error {
	// if VPC server id is set means the VPC is already created
	vpcID := s.GetVPCID()
	if vpcID != nil {
		s.Info("VPC id is set", "id", vpcID)
		vpcDetails, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
			ID: vpcID,
		})
		if err != nil {
			return err
		}
		if vpcDetails == nil {
			return fmt.Errorf("error failed to get vpc with id %s", *vpcID)
		}
		s.Info("Found VPC with provided id")
		// TODO(karthik-k-n): Set status here as well
		return nil
	}

	// check vpc exist in cloud
	id, err := s.checkVPC()
	if err != nil {
		return err
	}
	if id != "" {
		s.SetStatus(infrav1beta2.ResourceTypeVPC, infrav1beta2.ResourceReference{ID: &id, ControllerCreated: ptr.To(false)})
		return nil
	}

	// TODO(karthik-k-n): create a generic cluster scope/service and implement common vpc logics, which can be consumed by both vpc and powervs

	// create VPC
	s.Info("Creating a VPC")
	vpcDetails, err := s.createVPC()
	if err != nil {
		return err
	}
	s.Info("Successfully create VPC")
	s.SetStatus(infrav1beta2.ResourceTypeVPC, infrav1beta2.ResourceReference{ID: vpcDetails, ControllerCreated: ptr.To(true)})
	return nil
}

// checkVPC checks VPC exist in cloud.
func (s *PowerVSClusterScope) checkVPC() (string, error) {
	vpcDetails, err := s.getVPCByName()
	if err != nil {
		s.Error(err, "failed to get vpc")
		return "", err
	}
	if vpcDetails == nil {
		s.Info("Not able to find vpc", "vpc", s.IBMPowerVSCluster.Spec.VPC)
		return "", nil
	}
	s.Info("VPC found", "id", *vpcDetails.ID)
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
		s.Info("failed to create vpc, failed to fetch resource group id")
		return nil, fmt.Errorf("error getting resource group id for resource group %v, id is empty", s.ResourceGroup())
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

// ReconcileVPCSubnet reconciles VPC subnet.
func (s *PowerVSClusterScope) ReconcileVPCSubnet() error {
	subnets := make([]infrav1beta2.Subnet, 0)
	// check whether user has set the vpc subnets
	if len(s.IBMPowerVSCluster.Spec.VPCSubnets) == 0 {
		// if the user did not set any subnet, we try to create subnet in all the zones.
		powerVSZone := s.Zone()
		if powerVSZone == nil {
			return fmt.Errorf("error reconicling vpc subnet, powervs zone is not set")
		}
		region := endpoints.ConstructRegionFromZone(*powerVSZone)
		vpcZones, err := genUtil.VPCZonesForPowerVSRegion(region)
		if err != nil {
			return err
		}
		if len(vpcZones) == 0 {
			return fmt.Errorf("error reconicling vpc subnet,error getting vpc zones, no zone found for region %s", region)
		}
		for _, zone := range vpcZones {
			subnet := infrav1beta2.Subnet{
				Name: ptr.To(fmt.Sprintf("%s-%s", *s.GetServiceName(infrav1beta2.ResourceTypeSubnet), zone)),
				Zone: ptr.To(zone),
			}
			subnets = append(subnets, subnet)
		}
	}
	for index, subnet := range s.IBMPowerVSCluster.Spec.VPCSubnets {
		if subnet.Name == nil {
			subnet.Name = ptr.To(fmt.Sprintf("%s-%d", *s.GetServiceName(infrav1beta2.ResourceTypeSubnet), index))
		}
		subnets = append(subnets, subnet)
	}
	for _, subnet := range subnets {
		s.Info("Reconciling vpc subnet", "subnet", subnet)
		var subnetID *string
		if subnet.ID != nil {
			subnetID = subnet.ID
		} else {
			subnetID = s.GetVPCSubnetID(*subnet.Name)
		}
		if subnetID != nil {
			subnetDetails, _, err := s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{
				ID: subnetID,
			})
			if err != nil {
				return err
			}
			if subnetDetails == nil {
				return fmt.Errorf("error failed to get vpc subnet with id %s", *subnetID)
			}
			// check for next subnet
			continue
		}

		// check VPC subnet exist in cloud
		vpcSubnetID, err := s.checkVPCSubnet(*subnet.Name)
		if err != nil {
			s.Error(err, "error checking vpc subnet")
			return err
		}
		if vpcSubnetID != "" {
			s.Info("found vpc subnet", "id", vpcSubnetID)
			s.SetVPCSubnetID(*subnet.Name, infrav1beta2.ResourceReference{ID: &vpcSubnetID, ControllerCreated: ptr.To(false)})
			// check for next subnet
			continue
		}
		subnetID, err = s.createVPCSubnet(subnet)
		if err != nil {
			s.Error(err, "error creating vpc subnet")
			return err
		}
		s.Info("created vpc subnet", "id", subnetID)
		s.SetVPCSubnetID(*subnet.Name, infrav1beta2.ResourceReference{ID: subnetID, ControllerCreated: ptr.To(true)})
	}
	return nil
}

// checkVPCSubnet checks VPC subnet exist in cloud.
func (s *PowerVSClusterScope) checkVPCSubnet(subnetName string) (string, error) {
	vpcSubnet, err := s.IBMVPCClient.GetVPCSubnetByName(subnetName)
	if err != nil {
		return "", err
	}
	if vpcSubnet == nil {
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
		s.Info("failed to create vpc subnet, failed to fetch resource group id")
		return nil, fmt.Errorf("error getting resource group id for resource group %v, id is empty", s.ResourceGroup())
	}
	var zone string
	if subnet.Zone != nil {
		zone = *subnet.Zone
	} else {
		powerVSZone := s.Zone()
		if powerVSZone == nil {
			return nil, fmt.Errorf("error creating vpc subnet, powervs zone is not set")
		}
		region := endpoints.ConstructRegionFromZone(*powerVSZone)
		vpcZones, err := genUtil.VPCZonesForPowerVSRegion(region)
		if err != nil {
			return nil, err
		}
		// TODO(karthik-k-n): Decide on using all zones or using one zone
		if len(vpcZones) == 0 {
			return nil, fmt.Errorf("error getting vpc zones error: %v", err)
		}
		zone = vpcZones[0]
	}

	// create subnet
	vpcID := s.GetVPCID()
	cidrBlock, err := s.IBMVPCClient.GetSubnetAddrPrefix(*vpcID, zone)
	if err != nil {
		return nil, err
	}
	ipVersion := "ipv4"

	options := &vpcv1.CreateSubnetOptions{}
	options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
		IPVersion:     &ipVersion,
		Ipv4CIDRBlock: &cidrBlock,
		Name:          subnet.Name,
		VPC: &vpcv1.VPCIdentity{
			ID: vpcID,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
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
		return nil, fmt.Errorf("create subnet is nil")
	}
	return subnetDetails.ID, nil
}

// ReconcileTransitGateway reconcile transit gateway.
func (s *PowerVSClusterScope) ReconcileTransitGateway() error {
	if s.GetTransitGatewayID() != nil {
		s.Info("TransitGateway id is set", "id", s.GetTransitGatewayID())
		tg, _, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
			ID: s.GetTransitGatewayID(),
		})
		if err != nil {
			return err
		}
		err = s.checkTransitGatewayStatus(tg.ID)
		if err != nil {
			return err
		}
		return nil
	}

	// check transit gateway exist in cloud
	tgID, err := s.checkTransitGateway()
	if err != nil {
		return err
	}
	if tgID != "" {
		s.SetStatus(infrav1beta2.ResourceTypeTransitGateway, infrav1beta2.ResourceReference{ID: &tgID, ControllerCreated: ptr.To(false)})
		return nil
	}
	// create transit gateway
	transitGatewayID, err := s.createTransitGateway()
	if err != nil {
		return err
	}
	if transitGatewayID != nil {
		s.SetStatus(infrav1beta2.ResourceTypeTransitGateway, infrav1beta2.ResourceReference{ID: transitGatewayID, ControllerCreated: ptr.To(true)})
		return nil
	}
	return nil
}

// checkTransitGateway checks transit gateway exist in cloud.
func (s *PowerVSClusterScope) checkTransitGateway() (string, error) {
	// TODO(karthik-k-n): Support regex
	transitGateway, err := s.TransitGatewayClient.GetTransitGatewayByName(*s.GetServiceName(infrav1beta2.ResourceTypeTransitGateway))
	if err != nil {
		return "", err
	}
	if transitGateway == nil || transitGateway.ID == nil {
		return "", nil
	}
	if err = s.checkTransitGatewayStatus(transitGateway.ID); err != nil {
		return "", err
	}
	return *transitGateway.ID, nil
}

// checkTransitGatewayStatus checks transit gateway status in cloud.
func (s *PowerVSClusterScope) checkTransitGatewayStatus(transitGatewayID *string) error {
	transitGateway, _, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
		ID: transitGatewayID,
	})
	if err != nil {
		return err
	}
	if transitGateway == nil {
		return fmt.Errorf("tranist gateway is nil")
	}
	if *transitGateway.Status != string(infrav1beta2.TransitGatewayStateAvailable) {
		return fmt.Errorf("error tranist gateway %s not in available status, current status: %s", *transitGatewayID, *transitGateway.Status)
	}

	tgConnections, _, err := s.TransitGatewayClient.ListTransitGatewayConnections(&tgapiv1.ListTransitGatewayConnectionsOptions{
		TransitGatewayID: transitGateway.ID,
	})
	if err != nil {
		return fmt.Errorf("error listing transit gateway connections: %w", err)
	}

	if len(tgConnections.Connections) == 0 {
		return fmt.Errorf("no connections are attached to transit gateway")
	}

	vpcCRN, err := s.fetchVPCCRN()
	if err != nil {
		return fmt.Errorf("error failed to fetch VPC CRN: %w", err)
	}

	pvsServiceInstanceCRN, err := s.fetchPowerVSServiceInstanceCRN()
	if err != nil {
		return fmt.Errorf("error failed to fetch powervs service instance CRN: %w", err)
	}

	var powerVSAttached, vpcAttached bool
	for _, conn := range tgConnections.Connections {
		if *conn.NetworkType == string(vpcNetworkConnectionType) && *conn.NetworkID == *vpcCRN {
			if *conn.Status != string(infrav1beta2.TransitGatewayConnectionStateAttached) {
				return fmt.Errorf("error vpc connection not attached to transit gateway, current status: %s", *conn.Status)
			}
			vpcAttached = true
		}
		if *conn.NetworkType == string(powervsNetworkConnectionType) && *conn.NetworkID == *pvsServiceInstanceCRN {
			if *conn.Status != string(infrav1beta2.TransitGatewayConnectionStateAttached) {
				return fmt.Errorf("error powervs connection not attached to transit gateway, current status: %s", *conn.Status)
			}
			powerVSAttached = true
		}
	}
	if !powerVSAttached || !vpcAttached {
		return fmt.Errorf("either one of powervs or vpc transit gateway connections are not attached, PowerVS: %t VPC: %t", powerVSAttached, vpcAttached)
	}
	return nil
}

// createTransitGateway create transit gateway.
func (s *PowerVSClusterScope) createTransitGateway() (*string, error) {
	// TODO(karthik-k-n): Verify that the supplied zone supports PER
	// TODO(karthik-k-n): consider moving to clusterscope

	// fetch resource group id
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		s.Info("failed to create transit gateway, failed to fetch resource group id")
		return nil, fmt.Errorf("error getting resource group id for resource group %v, id is empty", s.ResourceGroup())
	}

	vpcRegion := s.getVPCRegion()
	if vpcRegion == nil {
		return nil, fmt.Errorf("failed to get vpc region")
	}

	tgName := s.GetServiceName(infrav1beta2.ResourceTypeTransitGateway)
	tg, _, err := s.TransitGatewayClient.CreateTransitGateway(&tgapiv1.CreateTransitGatewayOptions{
		Location:      vpcRegion,
		Name:          tgName,
		Global:        ptr.To(true),
		ResourceGroup: &tgapiv1.ResourceGroupIdentity{ID: ptr.To(resourceGroupID)},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating transit gateway: %w", err)
	}

	vpcCRN, err := s.fetchVPCCRN()
	if err != nil {
		return nil, fmt.Errorf("error failed to fetch VPC CRN: %w", err)
	}

	if _, _, err = s.TransitGatewayClient.CreateTransitGatewayConnection(&tgapiv1.CreateTransitGatewayConnectionOptions{
		TransitGatewayID: tg.ID,
		NetworkType:      ptr.To(string(vpcNetworkConnectionType)),
		NetworkID:        vpcCRN,
		Name:             ptr.To(fmt.Sprintf("%s-vpc-con", *tgName)),
	}); err != nil {
		return nil, fmt.Errorf("error creating vpc connection in transit gateway: %w", err)
	}

	pvsServiceInstanceCRN, err := s.fetchPowerVSServiceInstanceCRN()
	if err != nil {
		return nil, fmt.Errorf("error failed to fetch powervs service instance CRN: %w", err)
	}

	if _, _, err = s.TransitGatewayClient.CreateTransitGatewayConnection(&tgapiv1.CreateTransitGatewayConnectionOptions{
		TransitGatewayID: tg.ID,
		NetworkType:      ptr.To(string(powervsNetworkConnectionType)),
		NetworkID:        pvsServiceInstanceCRN,
		Name:             ptr.To(fmt.Sprintf("%s-pvs-con", *tgName)),
	}); err != nil {
		return nil, fmt.Errorf("error creating powervs connection in transit gateway: %w", err)
	}
	return tg.ID, nil
}

// ReconcileLoadBalancer reconcile loadBalancer.
func (s *PowerVSClusterScope) ReconcileLoadBalancer() error {
	loadBalancers := make([]infrav1beta2.VPCLoadBalancerSpec, 0)
	if len(s.IBMPowerVSCluster.Spec.LoadBalancers) == 0 {
		loadBalancer := infrav1beta2.VPCLoadBalancerSpec{
			Name:   *s.GetServiceName(infrav1beta2.ResourceTypeLoadBalancer),
			Public: ptr.To(true),
		}
		loadBalancers = append(loadBalancers, loadBalancer)
	}
	for index, loadBalancer := range s.IBMPowerVSCluster.Spec.LoadBalancers {
		if loadBalancer.Name == "" {
			loadBalancer.Name = fmt.Sprintf("%s-%d", *s.GetServiceName(infrav1beta2.ResourceTypeLoadBalancer), index)
		}
		loadBalancers = append(loadBalancers, loadBalancer)
	}

	for _, loadBalancer := range loadBalancers {
		var loadBalancerID *string
		if loadBalancer.ID != nil {
			loadBalancerID = loadBalancer.ID
		} else {
			loadBalancerID = s.GetLoadBalancerID(loadBalancer.Name)
		}
		if loadBalancerID != nil {
			s.Info("LoadBalancer ID is set, fetching loadbalancer details", "loadbalancerid", *loadBalancerID)
			loadBalancer, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
				ID: loadBalancerID,
			})
			if err != nil {
				return err
			}
			if infrav1beta2.VPCLoadBalancerState(*loadBalancer.ProvisioningStatus) != infrav1beta2.VPCLoadBalancerStateActive {
				return fmt.Errorf("loadbalancer is not in active state, current state %s", *loadBalancer.ProvisioningStatus)
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
			return err
		}
		if loadBalancerStatus != nil {
			s.SetLoadBalancerStatus(loadBalancer.Name, *loadBalancerStatus)
			continue
		}
		// create loadBalancer
		loadBalancerStatus, err = s.createLoadBalancer(loadBalancer)
		if err != nil {
			return err
		}
		s.SetLoadBalancerStatus(loadBalancer.Name, *loadBalancerStatus)
	}
	return nil
}

// checkLoadBalancer checks loadBalancer in cloud.
func (s *PowerVSClusterScope) checkLoadBalancer(lb infrav1beta2.VPCLoadBalancerSpec) (*infrav1beta2.VPCLoadBalancerStatus, error) {
	loadBalancer, err := s.IBMVPCClient.GetLoadBalancerByName(lb.Name)
	if err != nil {
		return nil, err
	}
	if loadBalancer == nil {
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
		s.Info("failed to create load balancer, failed to fetch resource group id")
		return nil, fmt.Errorf("error getting resource group id for resource group %v, id is empty", s.ResourceGroup())
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
		return nil, fmt.Errorf("error subnet required for load balancer creation")
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
	if s.COSInstance() == nil || s.COSInstance().Name == "" {
		return nil
	}

	// check COS service instance exist in cloud
	cosServiceInstanceStatus, err := s.checkCOSServiceInstance()
	if err != nil {
		s.Error(err, "error checking cos service instance")
		return err
	}
	if cosServiceInstanceStatus != nil {
		s.SetStatus(infrav1beta2.ResourceTypeCOSInstance, infrav1beta2.ResourceReference{ID: cosServiceInstanceStatus.GUID, ControllerCreated: ptr.To(false)})
	} else {
		// create COS service instance
		cosServiceInstanceStatus, err = s.createCOSServiceInstance()
		if err != nil {
			s.Error(err, "error creating cos service instance")
			return err
		}
		s.SetStatus(infrav1beta2.ResourceTypeCOSInstance, infrav1beta2.ResourceReference{ID: cosServiceInstanceStatus.GUID, ControllerCreated: ptr.To(true)})
	}

	props, err := authenticator.GetProperties()
	if err != nil {
		s.Error(err, "error while fetching service properties")
		return err
	}

	apiKey, ok := props["APIKEY"]
	if !ok {
		return fmt.Errorf("ibmcloud api key is not provided, set %s environmental variable", "IBMCLOUD_API_KEY")
	}
	region := s.IBMPowerVSCluster.Spec.CosInstance.BucketRegion
	// if the bucket region is not set, use vpc region
	if region == "" {
		vpcDetails := s.VPC()
		if vpcDetails == nil || vpcDetails.Region == nil {
			return fmt.Errorf("failed to determine cos bucket region, both buckeet region and vpc region not set")
		}
		region = *vpcDetails.Region
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
		s.Error(err, "error creating cosClient")
		return fmt.Errorf("failed to create cos client: %w", err)
	}
	s.COSClient = cosClient

	// check bucket exist in service instance
	if exist, err := s.checkCOSBucket(); exist {
		return nil
	} else if err != nil {
		s.Error(err, "error checking cos bucket")
		return err
	}

	// create bucket in service instance
	if err := s.createCOSBucket(); err != nil {
		s.Error(err, "error creating cos bucket")
		return err
	}
	return nil
}

func (s *PowerVSClusterScope) checkCOSBucket() (bool, error) {
	if _, err := s.COSClient.GetBucketByName(s.COSInstance().BucketName); err != nil {
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
		Bucket: ptr.To(s.COSInstance().BucketName),
	}
	_, err := s.COSClient.CreateBucket(input)
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return fmt.Errorf("error creating COS bucket %w", err)
	}

	switch aerr.Code() {
	// If bucket already exists, all good.
	case s3.ErrCodeBucketAlreadyOwnedByYou:
		return nil
	case s3.ErrCodeBucketAlreadyExists:
		return nil
	default:
		return fmt.Errorf("error creating COS bucket %w", err)
	}
}

func (s *PowerVSClusterScope) checkCOSServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	// check cos service instance
	serviceInstance, err := s.ResourceClient.GetInstanceByName(s.COSInstance().Name, resourcecontroller.CosResourceID, resourcecontroller.CosResourcePlanID)
	if err != nil {
		return nil, err
	}
	if serviceInstance == nil {
		s.Info("cos service instance is nil", "name", s.COSInstance().Name)
		return nil, nil
	}
	if *serviceInstance.State != string(infrav1beta2.ServiceInstanceStateActive) {
		s.Info("cos service instance not in active state", "current state", *serviceInstance.State)
		return nil, fmt.Errorf("cos instance not in active state, current state: %s", *serviceInstance.State)
	}
	return serviceInstance, nil
}

func (s *PowerVSClusterScope) createCOSServiceInstance() (*resourcecontrollerv2.ResourceInstance, error) {
	// fetch resource group id.
	resourceGroupID := s.GetResourceGroupID()
	if resourceGroupID == "" {
		s.Info("failed to create COS service instance, failed to fetch resource group id")
		return nil, fmt.Errorf("error getting resource group id for resource group %v, id is empty", s.ResourceGroup())
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
	rmv2, err := resourcemanagerv2.NewResourceManagerV2(&resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: s.session.Options.Authenticator,
	})
	if err != nil {
		return "", err
	}
	if rmv2 == nil {
		return "", fmt.Errorf("unable to get resource controller")
	}
	resourceGroup := s.ResourceGroup().Name
	rmv2ListResourceGroupOpt := resourcemanagerv2.ListResourceGroupsOptions{Name: resourceGroup, AccountID: &s.session.Options.UserAccount}
	resourceGroupListResult, _, err := rmv2.ListResourceGroups(&rmv2ListResourceGroupOpt)
	if err != nil {
		return "", err
	}

	if resourceGroupListResult != nil && len(resourceGroupListResult.Resources) > 0 {
		rg := resourceGroupListResult.Resources[0]
		resourceGroupID := *rg.ID
		return resourceGroupID, nil
	}

	err = fmt.Errorf("could not retrieve resource group id for %s", *resourceGroup)
	return "", err
}

// getVPCRegion returns region associated with VPC zone.
func (s *PowerVSClusterScope) getVPCRegion() *string {
	if s.IBMPowerVSCluster.Spec.VPC != nil {
		return s.IBMPowerVSCluster.Spec.VPC.Region
	}
	// if vpc region is not set try to fetch corresponding region from power vs zone
	zone := s.Zone()
	if zone == nil {
		s.Info("powervs zone is not set")
		return nil
	}
	region := endpoints.ConstructRegionFromZone(*zone)
	vpcRegion, err := genUtil.VPCRegionForPowerVSRegion(region)
	if err != nil {
		s.Error(err, fmt.Sprintf("failed to fetch vpc region associated with powervs region %s", region))
		return nil
	}
	return &vpcRegion
}

// fetchVPCCRN returns VPC CRN.
func (s *PowerVSClusterScope) fetchVPCCRN() (*string, error) {
	vpcDetails, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
		ID: s.GetVPCID(),
	})
	if err != nil {
		return nil, err
	}
	return vpcDetails.CRN, nil
}

// fetchPowerVSServiceInstanceCRN returns Power VS service instance CRN.
func (s *PowerVSClusterScope) fetchPowerVSServiceInstanceCRN() (*string, error) {
	serviceInstanceID := s.GetServiceInstanceID()
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
	case infrav1beta2.ResourceTypeSubnet:
		return ptr.To(fmt.Sprintf("%s-vpcsubnet", s.InfraCluster()))
	case infrav1beta2.ResourceTypeLoadBalancer:
		return ptr.To(fmt.Sprintf("%s-loadbalancer", s.InfraCluster()))
	}
	return nil
}

// DeleteLoadBalancer deletes loadBalancer.
func (s *PowerVSClusterScope) DeleteLoadBalancer() error {
	for _, lb := range s.IBMPowerVSCluster.Status.LoadBalancers {
		if lb.ID == nil || lb.ControllerCreated == nil || !*lb.ControllerCreated {
			continue
		}

		lb, _, err := s.IBMVPCClient.GetLoadBalancer(&vpcv1.GetLoadBalancerOptions{
			ID: lb.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "cannot be found") {
				return nil
			}
			return fmt.Errorf("error fetching the load balancer: %w", err)
		}

		if lb != nil && lb.ProvisioningStatus != nil && *lb.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateDeletePending) {
			if _, err = s.IBMVPCClient.DeleteLoadBalancer(&vpcv1.DeleteLoadBalancerOptions{
				ID: lb.ID,
			}); err != nil {
				s.Error(err, "error deleting the load balancer")
				return err
			}
			s.Info("Load balancer successfully deleted")
		}
	}
	return nil
}

// DeleteVPCSubnet deletes VPC subnet.
func (s *PowerVSClusterScope) DeleteVPCSubnet() error {
	for _, subnet := range s.IBMPowerVSCluster.Status.VPCSubnet {
		if subnet.ID == nil || subnet.ControllerCreated == nil || !*subnet.ControllerCreated {
			continue
		}

		net, _, err := s.IBMVPCClient.GetSubnet(&vpcv1.GetSubnetOptions{
			ID: subnet.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "Subnet not found") {
				return nil
			}
			return fmt.Errorf("error fetching the subnet: %w", err)
		}

		if _, err = s.IBMVPCClient.DeleteSubnet(&vpcv1.DeleteSubnetOptions{
			ID: net.ID,
		}); err != nil {
			return fmt.Errorf("error deleting VPC subnet: %w", err)
		}
		s.Info("VPC subnet successfully deleted")
	}
	return nil
}

// DeleteVPC deletes VPC.
func (s *PowerVSClusterScope) DeleteVPC() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeVPC) {
		return nil
	}

	if s.IBMPowerVSCluster.Status.VPC.ID == nil {
		return nil
	}

	vpc, _, err := s.IBMVPCClient.GetVPC(&vpcv1.GetVPCOptions{
		ID: s.IBMPowerVSCluster.Status.VPC.ID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "VPC not found") {
			return nil
		}
		return fmt.Errorf("error fetching the VPC: %w", err)
	}

	if _, err = s.IBMVPCClient.DeleteVPC(&vpcv1.DeleteVPCOptions{
		ID: vpc.ID,
	}); err != nil {
		return fmt.Errorf("error deleting VPC: %w", err)
	}
	s.Info("VPC successfully deleted")
	return nil
}

// DeleteTransitGateway deletes transit gateway.
func (s *PowerVSClusterScope) DeleteTransitGateway() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeTransitGateway) {
		return nil
	}

	if s.IBMPowerVSCluster.Status.TransitGateway.ID == nil {
		return nil
	}

	tg, _, err := s.TransitGatewayClient.GetTransitGateway(&tgapiv1.GetTransitGatewayOptions{
		ID: s.IBMPowerVSCluster.Status.TransitGateway.ID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "gateway was not found") {
			return nil
		}
		return fmt.Errorf("error fetching the transit gateway: %w", err)
	}

	tgConnections, _, err := s.TransitGatewayClient.ListTransitGatewayConnections(&tgapiv1.ListTransitGatewayConnectionsOptions{
		TransitGatewayID: tg.ID,
	})
	if err != nil {
		return fmt.Errorf("error listing transit gateway connections: %w", err)
	}

	for _, conn := range tgConnections.Connections {
		if conn.Status != nil && *conn.Status != string(infrav1beta2.TransitGatewayStateDeletePending) {
			_, err := s.TransitGatewayClient.DeleteTransitGatewayConnection(&tgapiv1.DeleteTransitGatewayConnectionOptions{
				ID:               conn.ID,
				TransitGatewayID: tg.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting transit gateway connection: %w", err)
			}
		}
	}

	if _, err = s.TransitGatewayClient.DeleteTransitGateway(&tgapiv1.DeleteTransitGatewayOptions{
		ID: s.IBMPowerVSCluster.Status.TransitGateway.ID,
	}); err != nil {
		return fmt.Errorf("error deleting transit gateway: %w", err)
	}
	s.Info("Transit gateway successfully deleted")
	return nil
}

// DeleteDHCPServer deletes DHCP server.
func (s *PowerVSClusterScope) DeleteDHCPServer() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeDHCPServer) {
		return nil
	}

	if s.IBMPowerVSCluster.Status.DHCPServer.ID == nil {
		return nil
	}

	server, err := s.IBMPowerVSClient.GetDHCPServer(*s.IBMPowerVSCluster.Status.DHCPServer.ID)
	if err != nil {
		if strings.Contains(err.Error(), "dhcp server does not exist") {
			return nil
		}
		return fmt.Errorf("error fetching DHCP server: %w", err)
	}

	if err = s.IBMPowerVSClient.DeleteDHCPServer(*server.ID); err != nil {
		return fmt.Errorf("error deleting the DHCP server: %w", err)
	}
	s.Info("DHCP server successfully deleted")
	return nil
}

// DeleteServiceInstance deletes service instance.
func (s *PowerVSClusterScope) DeleteServiceInstance() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeServiceInstance) {
		return nil
	}

	if s.IBMPowerVSCluster.Status.ServiceInstance.ID == nil {
		return nil
	}

	serviceInstance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: s.IBMPowerVSCluster.Status.ServiceInstance.ID,
	})
	if err != nil {
		return fmt.Errorf("error fetching service instance: %w", err)
	}

	if serviceInstance != nil && *serviceInstance.State == string(infrav1beta2.ServiceInstanceStateRemoved) {
		s.Info("PowerVS service instance has been removed")
		return nil
	}

	servers, err := s.IBMPowerVSClient.GetAllDHCPServers()
	if err != nil {
		return fmt.Errorf("error fetching networks in the service instance: %w", err)
	}

	if len(servers) > 0 {
		return fmt.Errorf("cannot delete service instance as DHCP server is not yet deleted")
	}

	if _, err = s.ResourceClient.DeleteResourceInstance(&resourcecontrollerv2.DeleteResourceInstanceOptions{
		ID: serviceInstance.ID,
		//Recursive: ptr.To(true),
	}); err != nil {
		s.Error(err, "error deleting Power VS service instance")
		return err
	}
	s.Info("Service instance successfully deleted")
	return nil
}

// DeleteCOSInstance deletes COS instance.
func (s *PowerVSClusterScope) DeleteCOSInstance() error {
	if !s.isResourceCreatedByController(infrav1beta2.ResourceTypeCOSInstance) {
		return nil
	}

	if s.IBMPowerVSCluster.Status.COSInstance.ID == nil {
		return nil
	}

	cosInstance, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: s.IBMPowerVSCluster.Status.COSInstance.ID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "COS instance unavailable") {
			return nil
		}
		return fmt.Errorf("error fetching COS instance: %w", err)
	}

	if cosInstance != nil && (*cosInstance.State == "pending_reclamation" || *cosInstance.State == string(infrav1beta2.ServiceInstanceStateRemoved)) {
		return nil
	}

	if _, err = s.ResourceClient.DeleteResourceInstance(&resourcecontrollerv2.DeleteResourceInstanceOptions{
		ID:        cosInstance.ID,
		Recursive: ptr.To(true),
	}); err != nil {
		s.Error(err, "error deleting COS service instance")
		return err
	}
	s.Info("COS instance successfully deleted")
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
