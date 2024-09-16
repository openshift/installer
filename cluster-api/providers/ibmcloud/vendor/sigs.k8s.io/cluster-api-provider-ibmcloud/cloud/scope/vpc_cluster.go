/*
Copyright 2024 The Kubernetes Authors.

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

	"github.com/go-logr/logr"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
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
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/globaltagging"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcemanager"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

const (
	// LOGDEBUGLEVEL indicates the debug level of the logs.
	LOGDEBUGLEVEL = 5
)

// VPCClusterScopeParams defines the input parameters used to create a new VPCClusterScope.
type VPCClusterScopeParams struct {
	Client          client.Client
	Cluster         *capiv1beta1.Cluster
	IBMVPCCluster   *infrav1beta2.IBMVPCCluster
	Logger          logr.Logger
	ServiceEndpoint []endpoints.ServiceEndpoint

	IBMVPCClient vpc.Vpc
}

// VPCClusterScope defines a scope defined around a VPC Cluster.
type VPCClusterScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	COSClient                cos.Cos
	GlobalTaggingClient      globaltagging.GlobalTagging
	ResourceControllerClient resourcecontroller.ResourceController
	ResourceManagerClient    resourcemanager.ResourceManager
	VPCClient                vpc.Vpc

	Cluster         *capiv1beta1.Cluster
	IBMVPCCluster   *infrav1beta2.IBMVPCCluster
	ServiceEndpoint []endpoints.ServiceEndpoint
}

// NewVPCClusterScope creates a new VPCClusterScope from the supplied parameters.
func NewVPCClusterScope(params VPCClusterScopeParams) (*VPCClusterScope, error) {
	if params.Client == nil {
		err := errors.New("error failed to generate new scope from nil Client")
		return nil, err
	}
	if params.Cluster == nil {
		err := errors.New("error failed to generate new scope from nil Cluster")
		return nil, err
	}
	if params.IBMVPCCluster == nil {
		err := errors.New("error failed to generate new scope from nil IBMVPCCluster")
		return nil, err
	}
	if params.Logger == (logr.Logger{}) {
		params.Logger = textlogger.NewLogger(textlogger.NewConfig())
	}

	helper, err := patch.NewHelper(params.IBMVPCCluster, params.Client)
	if err != nil {
		return nil, fmt.Errorf("error failed to init patch helper: %w", err)
	}

	vpcEndpoint := endpoints.FetchVPCEndpoint(params.IBMVPCCluster.Spec.Region, params.ServiceEndpoint)
	vpcClient, err := vpc.NewService(vpcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error failed to create IBM VPC client: %w", err)
	}

	if params.IBMVPCCluster.Spec.Network == nil || params.IBMVPCCluster.Spec.Region == "" {
		return nil, fmt.Errorf("error failed to generate vpc client as Network or Region is nil")
	}

	if params.Logger.V(LOGDEBUGLEVEL).Enabled() {
		core.SetLoggingLevel(core.LevelDebug)
	}

	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("error failed to create authenticator: %w", err)
	}

	// Create Global Tagging client.
	gtOptions := globaltagging.ServiceOptions{
		GlobalTaggingV1Options: &globaltaggingv1.GlobalTaggingV1Options{
			Authenticator: auth,
		},
	}
	// Override the global tagging endpoint if provided.
	if gtEndpoint := endpoints.FetchEndpoints(string(endpoints.GlobalTagging), params.ServiceEndpoint); gtEndpoint != "" {
		gtOptions.URL = gtEndpoint
		params.Logger.V(3).Info("Overriding the default global tagging endpoint", "GlobaTaggingEndpoint", gtEndpoint)
	}
	globalTaggingClient, err := globaltagging.NewService(gtOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create global tagging client: %w", err)
	}

	// Create Resource Controller client.
	rcOptions := resourcecontroller.ServiceOptions{
		ResourceControllerV2Options: &resourcecontrollerv2.ResourceControllerV2Options{
			Authenticator: auth,
		},
	}
	// Override the resource controller endpoint if provided.
	if rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), params.ServiceEndpoint); rcEndpoint != "" {
		rcOptions.URL = rcEndpoint
		params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
	}
	resourceControllerClient, err := resourcecontroller.NewService(rcOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource controller client: %w", err)
	}

	// Create Resource Manager client.
	rmOptions := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: auth,
	}
	// Override the ResourceManager endpoint if provided.
	if rmEndpoint := endpoints.FetchEndpoints(string(endpoints.RM), params.ServiceEndpoint); rmEndpoint != "" {
		rmOptions.URL = rmEndpoint
		params.Logger.V(3).Info("Overriding the default resource manager endpoint", "ResourceManagerEndpoint", rmEndpoint)
	}
	resourceManagerClient, err := resourcemanager.NewService(rmOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource manager client: %w", err)
	}

	clusterScope := &VPCClusterScope{
		Logger:                   params.Logger,
		Client:                   params.Client,
		patchHelper:              helper,
		Cluster:                  params.Cluster,
		IBMVPCCluster:            params.IBMVPCCluster,
		ServiceEndpoint:          params.ServiceEndpoint,
		GlobalTaggingClient:      globalTaggingClient,
		ResourceControllerClient: resourceControllerClient,
		ResourceManagerClient:    resourceManagerClient,
		VPCClient:                vpcClient,
	}
	return clusterScope, nil
}

// PatchObject persists the cluster configuration and status.
func (s *VPCClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMVPCCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *VPCClusterScope) Close() error {
	return s.PatchObject()
}

// Name returns the CAPI cluster name.
func (s *VPCClusterScope) Name() string {
	return s.Cluster.Name
}

// NetworkSpec returns the VPCClusterScope's Network spec.
func (s *VPCClusterScope) NetworkSpec() *infrav1beta2.VPCNetworkSpec {
	return s.IBMVPCCluster.Spec.Network
}

// NetworkStatus returns the VPCClusterScope's Network status.
func (s *VPCClusterScope) NetworkStatus() *infrav1beta2.VPCNetworkStatus {
	return s.IBMVPCCluster.Status.Network
}

// CheckTagExists checks whether a user tag already exists.
func (s *VPCClusterScope) CheckTagExists(tagName string) (bool, error) {
	exists, err := s.GlobalTaggingClient.GetTagByName(tagName)
	if err != nil {
		return false, fmt.Errorf("failed checking for tag: %w", err)
	}
	return exists != nil, nil
}

// GetNetworkResourceGroupID returns the Resource Group ID for the Network Resources if it is present. Otherwise, it defaults to the cluster's Resource Group ID.
func (s *VPCClusterScope) GetNetworkResourceGroupID() (string, error) {
	// Check if the ID is available from Status first.
	if s.NetworkStatus() != nil && s.NetworkStatus().ResourceGroup != nil && s.NetworkStatus().ResourceGroup.ID != "" {
		return s.NetworkStatus().ResourceGroup.ID, nil
	}

	// If there is no Network Resource Group defined, use the cluster's Resource Group.
	if s.NetworkSpec() == nil || s.NetworkSpec().ResourceGroup == nil || (s.NetworkSpec().ResourceGroup.ID == "" && s.NetworkSpec().ResourceGroup.Name == nil) {
		return s.GetResourceGroupID()
	}

	// Otherwise, collect the Network's Resource Group Id.
	resourceGroupID := s.NetworkSpec().ResourceGroup.ID
	var resourceGroupName *string
	if resourceGroupID != "" {
		// Verify the Resource Group exists, using the provided ID.
		resourceGroupDetails, _, err := s.ResourceManagerClient.GetResourceGroup(&resourcemanagerv2.GetResourceGroupOptions{
			ID: ptr.To(resourceGroupID),
		})
		if err != nil {
			return "", fmt.Errorf("failed to retrieve newtork resource group by id: %w", err)
		} else if resourceGroupDetails == nil || resourceGroupDetails.Name == nil {
			return "", fmt.Errorf("error retrieving network resource group by id: %s", resourceGroupID)
		}
		resourceGroupName = resourceGroupDetails.Name
	} else {
		// Retrieve the Resource Group based on the name (Name must exist if ID is empty).
		resourceGroup, err := s.ResourceManagerClient.GetResourceGroupByName(*s.NetworkSpec().ResourceGroup.Name)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve network resource group id by name: %w", err)
		} else if resourceGroup == nil || resourceGroup.ID == nil {
			return "", fmt.Errorf("error retrieving network resource group by name: %s", *s.NetworkSpec().ResourceGroup.Name)
		}
		resourceGroupID = *resourceGroup.ID
		resourceGroupName = s.NetworkSpec().ResourceGroup.Name
	}

	// Populate the Network Status' Resource Group to shortcut future lookups.
	s.SetResourceStatus(infrav1beta2.ResourceTypeResourceGroup, &infrav1beta2.ResourceStatus{
		ID:    resourceGroupID,
		Name:  resourceGroupName,
		Ready: true,
	})

	return resourceGroupID, nil
}

// GetResourceGroupID returns the Resource Group ID for the cluster.
func (s *VPCClusterScope) GetResourceGroupID() (string, error) {
	// Check if the Resource Group ID is available from Status first.
	if s.IBMVPCCluster.Status.ResourceGroup != nil && s.IBMVPCCluster.Status.ResourceGroup.ID != "" {
		return s.IBMVPCCluster.Status.ResourceGroup.ID, nil
	}

	// If the Resource Group is not defined in Spec, we generate the name based on the cluster name.
	resourceGroupName := s.IBMVPCCluster.Spec.ResourceGroup
	if resourceGroupName == "" {
		resourceGroupName = s.Name()
	}

	// Retrieve the Resource Group based on the name.
	resourceGroup, err := s.ResourceManagerClient.GetResourceGroupByName(resourceGroupName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve resource group by name: %w", err)
	} else if resourceGroup == nil || resourceGroup.ID == nil {
		return "", fmt.Errorf("failed to find resource group by name: %s", resourceGroupName)
	}

	// Populate the Stauts Resource Group to shortcut future lookups.
	s.SetResourceStatus(infrav1beta2.ResourceTypeResourceGroup, &infrav1beta2.ResourceStatus{
		ID:    *resourceGroup.ID,
		Name:  ptr.To(resourceGroupName),
		Ready: true,
	})

	return *resourceGroup.ID, nil
}

// GetServiceName returns the name of a given service type from Spec or generates a name for it.
func (s *VPCClusterScope) GetServiceName(resourceType infrav1beta2.ResourceType) *string {
	switch resourceType {
	case infrav1beta2.ResourceTypeVPC:
		// Generate a name based off cluster name if no VPC defined in Spec, or no VPC name nor ID.
		if s.NetworkSpec().VPC == nil || (s.NetworkSpec().VPC.Name == nil && s.NetworkSpec().VPC.ID == nil) {
			return ptr.To(fmt.Sprintf("%s-vpc", s.Name()))
		}
		if s.NetworkSpec().VPC.Name != nil {
			return s.NetworkSpec().VPC.Name
		}
	default:
		s.V(3).Info("unsupported resource type", "resourceType", resourceType)
	}
	return nil
}

// GetVPCID returns the VPC id, if available.
func (s *VPCClusterScope) GetVPCID() (*string, error) {
	// Check if the VPC ID is available from Status first.
	if s.NetworkStatus() != nil && s.NetworkStatus().VPC != nil {
		return ptr.To(s.NetworkStatus().VPC.ID), nil
	}

	if s.NetworkSpec() != nil && s.NetworkSpec().VPC != nil {
		if s.NetworkSpec().VPC.ID != nil {
			return s.NetworkSpec().VPC.ID, nil
		} else if s.NetworkSpec().VPC.Name != nil {
			vpcDetails, err := s.VPCClient.GetVPCByName(*s.NetworkSpec().VPC.Name)
			if err != nil {
				return nil, fmt.Errorf("failed vpc id lookup: %w", err)
			}

			// Check if the VPC was found and has an ID
			if vpcDetails != nil && vpcDetails.ID != nil {
				// Set VPC ID in Status to shortcut future lookups
				s.SetResourceStatus(infrav1beta2.ResourceTypeVPC, &infrav1beta2.ResourceStatus{
					ID:    *vpcDetails.ID,
					Name:  s.NetworkSpec().VPC.Name,
					Ready: true,
				})
			}
		}
	}
	return nil, nil
}

// SetResourceStatus sets the status for the provided ResourceType.
func (s *VPCClusterScope) SetResourceStatus(resourceType infrav1beta2.ResourceType, resource *infrav1beta2.ResourceStatus) {
	// Ignore attempts to set status without resource.
	if resource == nil {
		return
	}
	s.V(3).Info("Setting status", "resourceType", resourceType, "resource", resource)
	switch resourceType {
	case infrav1beta2.ResourceTypeResourceGroup:
		if s.IBMVPCCluster.Status.ResourceGroup == nil {
			s.IBMVPCCluster.Status.ResourceGroup = resource
			return
		}
		s.IBMVPCCluster.Status.ResourceGroup.Set(*resource)
	case infrav1beta2.ResourceTypeVPC:
		if s.NetworkStatus() == nil {
			s.IBMVPCCluster.Status.Network = &infrav1beta2.VPCNetworkStatus{
				VPC: resource,
			}
			return
		} else if s.NetworkStatus().VPC == nil {
			s.IBMVPCCluster.Status.Network.VPC = resource
		}
		s.NetworkStatus().VPC.Set(*resource)
	default:
		s.V(3).Info("unsupported resource type", "resourceType", resourceType)
	}
}

// TagResource will attach a user Tag to a resource.
func (s *VPCClusterScope) TagResource(tagName string, resourceCRN string) error {
	// Verify the Tag we wish to use exists, otherwise create it.
	exists, err := s.CheckTagExists(tagName)
	if err != nil {
		return fmt.Errorf("failure checking if tag exists: %w", err)
	}

	// Create tag if it doesn't exist.
	if !exists {
		createOptions := &globaltaggingv1.CreateTagOptions{}
		createOptions.SetTagNames([]string{tagName})
		if _, _, err := s.GlobalTaggingClient.CreateTag(createOptions); err != nil {
			return fmt.Errorf("failure creating tag: %w", err)
		}
	}

	// Finally, tag resource.
	tagOptions := &globaltaggingv1.AttachTagOptions{}
	tagOptions.SetResources([]globaltaggingv1.Resource{
		{
			ResourceID: ptr.To(resourceCRN),
		},
	})
	tagOptions.SetTagName(tagName)
	tagOptions.SetTagType(globaltaggingv1.AttachTagOptionsTagTypeUserConst)

	if _, _, err = s.GlobalTaggingClient.AttachTag(tagOptions); err != nil {
		return fmt.Errorf("failure tagging resource: %w", err)
	}

	return nil
}

// ReconcileVPC reconciles the cluster's VPC.
func (s *VPCClusterScope) ReconcileVPC() (bool, error) {
	// If VPC id is set, that indicates the VPC already exists.
	vpcID, err := s.GetVPCID()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve vpc id: %w", err)
	}
	if vpcID != nil {
		s.V(3).Info("VPC id is set", "id", vpcID)
		vpcDetails, _, err := s.VPCClient.GetVPC(&vpcv1.GetVPCOptions{
			ID: vpcID,
		})
		if err != nil {
			return false, fmt.Errorf("failed to retrieve vpc by id: %w", err)
		} else if vpcDetails == nil {
			return false, fmt.Errorf("failed to retrieve vpc with id: %s", *vpcID)
		}
		s.V(3).Info("Found VPC with provided id", "id", vpcID)

		requeue := true
		if vpcDetails.Status != nil && *vpcDetails.Status == string(vpcv1.VPCStatusAvailableConst) {
			requeue = false
		}
		s.SetResourceStatus(infrav1beta2.ResourceTypeVPC, &infrav1beta2.ResourceStatus{
			ID:   *vpcID,
			Name: vpcDetails.Name,
			// Ready status will be invert of the need to requeue.
			Ready: !requeue,
		})

		// After updating the Status of VPC, return with requeue or return as reconcile complete.
		return requeue, nil
	}

	// If no VPC id was found, we need to create a new VPC.
	s.V(3).Info("Creating a VPC")
	vpcDetails, err := s.createVPC()
	if err != nil {
		return false, fmt.Errorf("failed to create vpc: %w", err)
	}

	s.V(3).Info("Successfully created VPC")
	var vpcName *string
	if vpcDetails != nil {
		vpcName = vpcDetails.Name
	}
	s.SetResourceStatus(infrav1beta2.ResourceTypeVPC, &infrav1beta2.ResourceStatus{
		ID:    *vpcDetails.ID,
		Name:  vpcName,
		Ready: false,
	})
	return true, nil
}

func (s *VPCClusterScope) createVPC() (*vpcv1.VPC, error) {
	// We use the cluster's Resource Group ID, as we expect to create all resources in that Resource Group.
	resourceGroupID, err := s.GetResourceGroupID()
	if err != nil {
		return nil, fmt.Errorf("failed retreiving resource group id during vpc creation: %w", err)
	} else if resourceGroupID == "" {
		return nil, fmt.Errorf("resource group id is empty cannot create vpc")
	}
	vpcName := s.GetServiceName(infrav1beta2.ResourceTypeVPC)
	if s.NetworkSpec() != nil && s.NetworkSpec().VPC != nil && s.NetworkSpec().VPC.Name != nil {
		vpcName = s.NetworkSpec().VPC.Name
	}

	// TODO(cjschaef): Look at adding support to specify prefix management
	addressPrefixManagement := "auto"
	vpcOptions := &vpcv1.CreateVPCOptions{
		AddressPrefixManagement: &addressPrefixManagement,
		Name:                    vpcName,
		ResourceGroup:           &vpcv1.ResourceGroupIdentity{ID: &resourceGroupID},
	}
	vpcDetails, _, err := s.VPCClient.CreateVPC(vpcOptions)
	if err != nil {
		return nil, fmt.Errorf("error creating vpc: %w", err)
	} else if vpcDetails == nil {
		return nil, fmt.Errorf("no vpc details after creation")
	}

	// NOTE: This tagging is only attempted once. We may wish to refactor in case this single attempt fails.
	if err = s.TagResource(s.Name(), *vpcDetails.CRN); err != nil {
		return nil, fmt.Errorf("error tagging vpc: %w", err)
	}

	return vpcDetails, nil
}
