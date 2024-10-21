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
	"reflect"

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

	// vpcSubnetIPVersion4 defines the IP v4 string used for VPC Subnet generation.
	vpcSubnetIPVersion4 = "ipv4"
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

func (s *VPCClusterScope) getSecurityGroupIDFromStatus(name string) *string {
	if s.NetworkStatus() != nil && s.NetworkStatus().SecurityGroups != nil {
		if sg, ok := s.NetworkStatus().SecurityGroups[name]; ok {
			return ptr.To(sg.ID)
		}
	}

	// Security Group was not found in Status, return nil.
	return nil
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
	case infrav1beta2.ResourceTypeSubnet:
		// Generate a generic subnet name based off the cluster name, which can be extended as necessary (for Zones).
		return ptr.To(fmt.Sprintf("%s-subnet", s.IBMVPCCluster.Name))
	case infrav1beta2.ResourceTypePublicGateway:
		// Generate a generic public gateway name based off the cluster name, which can be extedned as necessary (for Zone).
		return ptr.To(fmt.Sprintf("%s-pgateway", s.IBMVPCCluster.Name))
	default:
		s.V(3).Info("unsupported resource type", "resourceType", resourceType)
	}
	return nil
}

// GetSubnetID returns the ID of a subnet, provided the name.
func (s *VPCClusterScope) GetSubnetID(name string) (*string, error) {
	// Check Status first
	if s.NetworkStatus() != nil {
		if s.NetworkStatus().ControlPlaneSubnets != nil {
			if subnet, ok := s.NetworkStatus().ControlPlaneSubnets[name]; ok {
				return &subnet.ID, nil
			}
		}
		if s.NetworkStatus().WorkerSubnets != nil {
			if subnet, ok := s.NetworkStatus().WorkerSubnets[name]; ok {
				return &subnet.ID, nil
			}
		}
	}
	// Otherwise, if no Status, or not found, attempt to look it up via IBM Cloud API.
	subnet, err := s.VPCClient.GetVPCSubnetByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed subnet id lookup by name %s: %w", name, err)
	}
	if subnet == nil {
		return nil, nil
	}
	return subnet.ID, nil
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
func (s *VPCClusterScope) SetResourceStatus(resourceType infrav1beta2.ResourceType, resource *infrav1beta2.ResourceStatus) { //nolint:gocyclo
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
	case infrav1beta2.ResourceTypeCustomImage:
		if s.IBMVPCCluster.Status.Image == nil {
			s.IBMVPCCluster.Status.Image = &infrav1beta2.ResourceStatus{
				ID:    resource.ID,
				Name:  resource.Name,
				Ready: resource.Ready,
			}
			return
		}
		s.IBMVPCCluster.Status.Image.Set(*resource)
	case infrav1beta2.ResourceTypeControlPlaneSubnet:
		if s.NetworkStatus() == nil {
			s.IBMVPCCluster.Status.Network = &infrav1beta2.VPCNetworkStatus{}
		}
		if s.NetworkStatus().ControlPlaneSubnets == nil {
			s.IBMVPCCluster.Status.Network.ControlPlaneSubnets = make(map[string]*infrav1beta2.ResourceStatus)
		}
		if subnet, ok := s.NetworkStatus().ControlPlaneSubnets[*resource.Name]; ok {
			subnet.Set(*resource)
		} else {
			s.IBMVPCCluster.Status.Network.ControlPlaneSubnets[*resource.Name] = resource
		}
	case infrav1beta2.ResourceTypeWorkerSubnet:
		if s.NetworkStatus() == nil {
			s.IBMVPCCluster.Status.Network = &infrav1beta2.VPCNetworkStatus{}
		}
		if s.NetworkStatus().WorkerSubnets == nil {
			s.IBMVPCCluster.Status.Network.WorkerSubnets = make(map[string]*infrav1beta2.ResourceStatus)
		}
		if subnet, ok := s.NetworkStatus().WorkerSubnets[*resource.Name]; ok {
			subnet.Set(*resource)
		} else {
			s.IBMVPCCluster.Status.Network.WorkerSubnets[*resource.Name] = resource
		}
	case infrav1beta2.ResourceTypeSecurityGroup:
		if s.NetworkStatus() == nil {
			s.IBMVPCCluster.Status.Network = &infrav1beta2.VPCNetworkStatus{}
		}
		if s.IBMVPCCluster.Status.Network.SecurityGroups == nil {
			s.IBMVPCCluster.Status.Network.SecurityGroups = make(map[string]*infrav1beta2.ResourceStatus)
		}
		if securityGroup, ok := s.IBMVPCCluster.Status.Network.SecurityGroups[*resource.Name]; ok {
			securityGroup.Set(*resource)
		} else {
			s.IBMVPCCluster.Status.Network.SecurityGroups[*resource.Name] = resource
		}
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
	err = s.createVPC()
	if err != nil {
		return false, fmt.Errorf("failed to create vpc: %w", err)
	}

	s.V(3).Info("Successfully created VPC")
	return true, nil
}

func (s *VPCClusterScope) createVPC() error {
	// We use the cluster's Resource Group ID, as we expect to create all resources in that Resource Group.
	resourceGroupID, err := s.GetResourceGroupID()
	if err != nil {
		return fmt.Errorf("failed retreiving resource group id during vpc creation: %w", err)
	} else if resourceGroupID == "" {
		return fmt.Errorf("resource group id is empty cannot create vpc")
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
		return fmt.Errorf("error creating vpc: %w", err)
	} else if vpcDetails == nil {
		return fmt.Errorf("no vpc details after creation")
	}

	// Set the VPC status.
	s.SetResourceStatus(infrav1beta2.ResourceTypeVPC, &infrav1beta2.ResourceStatus{
		ID:   *vpcDetails.ID,
		Name: vpcDetails.Name,
		// We wait for a followup reconcile loop to set as Ready, to confirm the VPC can be found.
		Ready: false,
	})

	// NOTE: This tagging is only attempted once. We may wish to refactor in case this single attempt fails.
	if err = s.TagResource(s.Name(), *vpcDetails.CRN); err != nil {
		return fmt.Errorf("error tagging vpc: %w", err)
	}

	return nil
}

// ReconcileVPCCustomImage reconciles the VPC Custom Image.
func (s *VPCClusterScope) ReconcileVPCCustomImage() (bool, error) {
	// VPC Custom Image reconciliation is based on the following possibilities.
	// 1. Check Status for ID or Name, from previous lookup in reconciliation loop.
	// 2. If no Image spec is provided, assume the image is managed externally, thus no reconciliation required.
	// 3. If Image name is provided, check if an existing VPC Custom Image exists with that name (unfortunately names may not be unique), checking status of the image, updating appropriately.
	// 4. If Image CRN is provided, parse the ID from the CRN to perform lookup. CRN may be for another account, causing lookup to fail (permissions), may require better safechecks based on other CRN details.
	// 5. If no Image ID has been identified, assume a VPC Custom Image needs to be created, do so.
	var imageID *string
	// Attempt to collect VPC Custom Image info from Status.
	if s.IBMVPCCluster.Status.Image != nil {
		if s.IBMVPCCluster.Status.Image.ID != "" {
			imageID = ptr.To(s.IBMVPCCluster.Status.Image.ID)
		}
	} else if s.IBMVPCCluster.Spec.Image == nil {
		// If no Image spec was defined, we expect it is maintained externally and continue without reconciling. For example, using a Catalog Offering Custom Image, which may be in another account, which means it cannot be looked up, but can be used when creating Instances.
		s.V(3).Info("No VPC Custom Image defined, skipping reconciliation")
		return false, nil
	} else if s.IBMVPCCluster.Spec.Image.Name != nil {
		// Attempt to retrieve the image details via the name, if it already exists
		imageDetails, err := s.VPCClient.GetImageByName(*s.IBMVPCCluster.Spec.Image.Name)
		if err != nil {
			return false, fmt.Errorf("error checking vpc custom image by name: %w", err)
		} else if imageDetails != nil && imageDetails.ID != nil {
			// Prevent relookup (API request) of VPC Custom Image if we already have the necessary data
			requeue := true
			if imageDetails.Status != nil && *imageDetails.Status == string(vpcv1.ImageStatusAvailableConst) {
				requeue = false
			}
			s.SetResourceStatus(infrav1beta2.ResourceTypeCustomImage, &infrav1beta2.ResourceStatus{
				ID:   *imageDetails.ID,
				Name: s.IBMVPCCluster.Spec.Image.Name,
				// Ready status will be invert of the need to requeue.
				Ready: !requeue,
			})
			return requeue, nil
		}
	} else if s.IBMVPCCluster.Spec.Image.CRN != nil {
		// Parse the supplied Image CRN for Id, to perform image lookup.
		imageCRN, err := ParseCRN(*s.IBMVPCCluster.Spec.Image.CRN)
		if err != nil {
			return false, fmt.Errorf("error parsing vpc custom image crn: %w", err)
		}
		// If the value provided isn't a CRN or is missing the Resource ID, raise an error.
		if imageCRN == nil || imageCRN.Resource == "" {
			return false, fmt.Errorf("error parsing vpc custom image crn, missing resource id")
		}
		// If we didn't hit an error during parsing, and Resource was set, set that as the Image ID.
		imageID = ptr.To(imageCRN.Resource)
	}

	// Check status of VPC Custom Image.
	if imageID != nil {
		image, _, err := s.VPCClient.GetImage(&vpcv1.GetImageOptions{
			ID: imageID,
		})
		if err != nil {
			return false, fmt.Errorf("error retrieving vpc custom image by id: %w", err)
		}
		if image == nil {
			return false, fmt.Errorf("error failed to retrieve vpc custom image with id %s", *imageID)
		}
		s.V(3).Info("Found VPC Custom Image with provided id", "imageID", imageID)

		requeue := true
		if image.Status != nil && *image.Status == string(vpcv1.ImageStatusAvailableConst) {
			requeue = false
		}
		s.SetResourceStatus(infrav1beta2.ResourceTypeCustomImage, &infrav1beta2.ResourceStatus{
			ID:   *imageID,
			Name: image.Name,
			// Ready status will be invert of the need to requeue.
			Ready: !requeue,
		})
		return requeue, nil
	}

	// No VPC Custom Image exists or was found, so create the Custom Image.
	s.V(3).Info("Creating a VPC Custom Image")
	err := s.createCustomImage()
	if err != nil {
		return false, fmt.Errorf("error failure trying to create vpc custom image: %w", err)
	}

	s.V(3).Info("Successfully created VPC Custom Image")
	return true, nil
}

// createCustomImage will create a new VPC Custom Image.
func (s *VPCClusterScope) createCustomImage() error {
	// TODO(cjschaef): Remove in favor of webhook validation.
	if s.IBMVPCCluster.Spec.Image.OperatingSystem == nil {
		return fmt.Errorf("error failed to create vpc custom image due to missing operatingSystem")
	}

	// Collect the Resource Group ID.
	var resourceGroupID *string
	// Check Resource Group in Image spec.
	if s.IBMVPCCluster.Spec.Image.ResourceGroup != nil {
		if s.IBMVPCCluster.Spec.Image.ResourceGroup.ID != "" {
			resourceGroupID = ptr.To(s.IBMVPCCluster.Spec.Image.ResourceGroup.ID)
		} else if s.IBMVPCCluster.Spec.Image.ResourceGroup.Name != nil {
			id, err := s.ResourceManagerClient.GetResourceGroupByName(*s.IBMVPCCluster.Spec.Image.ResourceGroup.Name)
			if err != nil {
				return fmt.Errorf("error retrieving resource group by name: %w", err)
			}
			resourceGroupID = id.ID
		}
	} else {
		// Otherwise, we will use the cluster Resource Group ID, as we expect to create all resources in that Resource Group.
		id, err := s.GetResourceGroupID()
		if err != nil {
			return fmt.Errorf("error retrieving resource group id: %w", err)
		}
		resourceGroupID = ptr.To(id)
	}

	// Build the COS Object URL using the ImageSpec
	fileHRef, err := s.buildCOSObjectHRef()
	if err != nil {
		return fmt.Errorf("error building vpc custom image file href: %w", err)
	}

	options := &vpcv1.CreateImageOptions{
		ImagePrototype: &vpcv1.ImagePrototype{
			Name: s.IBMVPCCluster.Spec.Image.Name,
			File: &vpcv1.ImageFilePrototype{
				Href: fileHRef,
			},
			OperatingSystem: &vpcv1.OperatingSystemIdentity{
				Name: s.IBMVPCCluster.Spec.Image.OperatingSystem,
			},
			ResourceGroup: &vpcv1.ResourceGroupIdentity{
				ID: resourceGroupID,
			},
		},
	}

	imageDetails, _, err := s.VPCClient.CreateImage(options)
	if err != nil {
		return fmt.Errorf("error unknown failure creating vpc custom image: %w", err)
	}
	if imageDetails == nil || imageDetails.ID == nil || imageDetails.Name == nil || imageDetails.CRN == nil {
		return fmt.Errorf("error failed creating custom image")
	}

	// Initially populate the Image's status.
	s.SetResourceStatus(infrav1beta2.ResourceTypeCustomImage, &infrav1beta2.ResourceStatus{
		ID:   *imageDetails.ID,
		Name: imageDetails.Name,
		// We must wait for the image to be ready, on followup reconciliation loops.
		Ready: false,
	})

	// NOTE: This tagging is only attempted once. We may wish to refactor in case this single attempt fails.
	if err := s.TagResource(s.Name(), *imageDetails.CRN); err != nil {
		return fmt.Errorf("error failure tagging vpc custom image: %w", err)
	}
	return nil
}

// buildCOSObjectHRef will build the HRef path to a COS Object that can be used for VPC Custom Image creation.
func (s *VPCClusterScope) buildCOSObjectHRef() (*string, error) {
	// TODO(cjschaef): Remove in favor of webhook validation.
	// We need COS details in order to create the Custom Image from.
	if s.IBMVPCCluster.Spec.Image.COSInstance == nil || s.IBMVPCCluster.Spec.Image.COSBucket == nil || s.IBMVPCCluster.Spec.Image.COSObject == nil {
		return nil, fmt.Errorf("error failed to build cos object href, cos details missing")
	}

	// Get COS Bucket Region, defaulting to cluster Region if not specified.
	bucketRegion := s.IBMVPCCluster.Spec.Region
	if s.IBMVPCCluster.Spec.Image.COSBucketRegion != nil {
		bucketRegion = *s.IBMVPCCluster.Spec.Image.COSBucketRegion
	}

	// Expected HRef format:
	//   cos://<bucket_region>/<bucket_name>/<object_name>
	href := fmt.Sprintf("cos://%s/%s/%s", bucketRegion, *s.IBMVPCCluster.Spec.Image.COSBucket, *s.IBMVPCCluster.Spec.Image.COSObject)
	s.V(3).Info("building image ref", "href", href)
	return ptr.To(href), nil
}

// ReconcileSubnets reconciles the VPC Subnet(s).
// For Subnets, we collect all of the required subnets, for each Plane, and reconcile them individually. Requeing if one is missing or just created. Reconciliation is attempted on all subnets each loop, to prevent single subnet creation per reconciliation loop.
func (s *VPCClusterScope) ReconcileSubnets() (bool, error) {
	var subnets []infrav1beta2.Subnet
	var err error
	// If no ControlPlane Subnets were supplied, we default to create one in each availability zone of the region.
	if len(s.IBMVPCCluster.Spec.Network.ControlPlaneSubnets) == 0 {
		subnets, err = s.buildSubnetsForZones()
		if err != nil {
			return false, fmt.Errorf("error failed building control plane subnets: %w", err)
		}
	} else {
		subnets = s.IBMVPCCluster.Spec.Network.ControlPlaneSubnets
	}

	// Reconcile Control Plane subnets.
	requeue := false
	for _, subnet := range subnets {
		if requiresRequeue, err := s.reconcileSubnet(subnet, true); err != nil {
			return false, fmt.Errorf("error failed reconciling control plane subnet: %w", err)
		} else if requiresRequeue {
			// If the reconcile of the subnet requires further reconciliation, plan to requeue entire ReconcileSubnets call, but attempt to further reconcile additional Subnets (attempt all subnet reconciliation).
			requeue = true
		}
	}

	// If no Worker subnets were supplied, attempt to create one in each zone.
	if len(s.IBMVPCCluster.Spec.Network.WorkerSubnets) == 0 {
		// Build subnets for Workers if none were provided, but only if Control Plane subnets were.
		// Otherwise, if neither Control Plane nor Worker subnets were supplied, we rely on both Planes using the same subnet per zone, and we will re-reconcile those subnets below, for IBMVPCCluster Status updates.
		if len(s.IBMVPCCluster.Spec.Network.ControlPlaneSubnets) != 0 {
			subnets, err = s.buildSubnetsForZones()
			if err != nil {
				return false, fmt.Errorf("error failed building worker subnets: %w", err)
			}
		}
	} else {
		subnets = s.IBMVPCCluster.Spec.Network.WorkerSubnets
	}

	// Reconcile Worker subnets.
	for _, subnet := range subnets {
		if requiresRequeue, err := s.reconcileSubnet(subnet, false); err != nil {
			return false, fmt.Errorf("error failed reconciling worker subnet: %w", err)
		} else if requiresRequeue {
			// If the reconcile of the subnet requires further reconciliation, plan to requeue entire ReconcileSubnets call, but attempt to further reconcile additional Subnets (attempt all subnet reconciliation).
			requeue = true
		}
	}

	// Return whether or not one or more subnets required further reconciling after attempting to process all Control Plane and Worker subnets.
	return requeue, nil
}

// reconcileSubnet will attempt to find the existing subnet, or create it if necessary.
// The logic can handle either Control Plane or Worker subnets, but must distinguish between them for Status updates.
func (s *VPCClusterScope) reconcileSubnet(subnet infrav1beta2.Subnet, isControlPlane bool) (bool, error) { //nolint: gocyclo
	// If no ID or name was provided, that is an error to be raised. One or the other must be specified when subnets are supplied.
	if subnet.ID == nil && subnet.Name == nil {
		return false, fmt.Errorf("error subnet has no defined id or name, one is required")
	}

	// Check Status first and update as necessary.
	if s.NetworkStatus() != nil {
		var subnetMap map[string]*infrav1beta2.ResourceStatus
		var subnetID, subnetName *string
		if isControlPlane && s.NetworkStatus().ControlPlaneSubnets != nil {
			subnetMap = s.NetworkStatus().ControlPlaneSubnets
		} else if !isControlPlane && s.NetworkStatus().WorkerSubnets != nil {
			subnetMap = s.NetworkStatus().WorkerSubnets
		}
		// Based on Network Status, setup either the name or ID for lookup of the subnet's current status.
		if subnet.Name != nil {
			if _, ok := subnetMap[*subnet.Name]; ok {
				subnetName = subnet.Name
			}
		} else if subnet.ID != nil {
			for _, statusSubnet := range subnetMap {
				if statusSubnet.ID == *subnet.ID {
					subnetID = subnet.ID
				}
			}
		}

		// Perform current status lookup of subnet, using ID or name if one was found in Network Status.
		if subnetID != nil {
			options := &vpcv1.GetSubnetOptions{
				ID: subnetID,
			}
			subnetDetails, _, err := s.VPCClient.GetSubnet(options)
			if err != nil {
				return false, fmt.Errorf("error retrieving existing subnet by id %s: %w", *subnetID, err)
			} else if subnetDetails == nil {
				return false, fmt.Errorf("error failed to find existing subnet by id %s", *subnetID)
			}
			return s.updateSubnetStatus(subnetDetails, isControlPlane)
		} else if subnetName != nil {
			subnetDetails, err := s.VPCClient.GetVPCSubnetByName(*subnetName)
			if err != nil {
				return false, fmt.Errorf("error retrieving existing subnet by name %s: %w", *subnetName, err)
			} else if subnetDetails == nil {
				return false, fmt.Errorf("error failed to find existing subnet by name: %s", *subnetName)
			}
			return s.updateSubnetStatus(subnetDetails, isControlPlane)
		}
	}

	// Otherwise, if these is an ID or name, attempt to lookup the subnet and update status as necessary.
	if subnet.ID != nil {
		options := &vpcv1.GetSubnetOptions{
			ID: subnet.ID,
		}
		subnetDetails, _, err := s.VPCClient.GetSubnet(options)
		if err != nil {
			return false, fmt.Errorf("error retrieving subnet by id %s: %w", *subnet.ID, err)
		} else if subnetDetails == nil {
			// If the subnet was not found with provided ID, that is an error and a new subnet will not be created.
			return false, fmt.Errorf("error failed to find subnet with id: %s", *subnet.ID)
		}
		return s.updateSubnetStatus(subnetDetails, isControlPlane)
	} else if subnet.Name != nil {
		// Attempt to check if a subnet exists with the name and update status as necessary.
		subnetDetails, err := s.VPCClient.GetVPCSubnetByName(*subnet.Name)
		if err != nil {
			return false, fmt.Errorf("error retrieving subnet by name %s: %w", *subnet.Name, err)
		} else if subnetDetails != nil {
			// Update status if subnet was found.
			return s.updateSubnetStatus(subnetDetails, isControlPlane)
		}
		// If subnet was not found, expect that it needs to be created.
	}

	// If the subnet has not yet been at this point, assume it needs to be created.
	s.V(3).Info("creating subnet", "subnetName", subnet.Name)
	err := s.createSubnet(subnet, isControlPlane)
	if err != nil {
		return false, err
	}
	s.V(3).Info("Successfully created subnet", "subnetName", subnet.Name)

	// Recommend we requeue reconciliation after subnet was successfully created
	return true, nil
}

// buildSubnetsForZones will create a set of Subnets, using default names, for each availability zone within a Region. This is typically used when no subnets were provided, so a set of default subnets gets created.
func (s *VPCClusterScope) buildSubnetsForZones() ([]infrav1beta2.Subnet, error) {
	subnets := make([]infrav1beta2.Subnet, 0)
	zones, err := s.VPCClient.GetVPCZonesByRegion(s.IBMVPCCluster.Spec.Region)
	if err != nil {
		return subnets, fmt.Errorf("error unknown failure retrieving zones for region %s: %w", s.IBMVPCCluster.Spec.Region, err)
	}
	if len(zones) == 0 {
		return subnets, fmt.Errorf("error retrieving subnet zones, no zones found in %s", s.IBMVPCCluster.Spec.Region)
	}
	for _, zone := range zones {
		name := fmt.Sprintf("%s-%s", *s.GetServiceName(infrav1beta2.ResourceTypeSubnet), zone)
		subnets = append(subnets, infrav1beta2.Subnet{
			Name: ptr.To(name),
			Zone: ptr.To(zone),
		})
	}
	return subnets, nil
}

// updateSubnetStatus will check the status of a IBM Cloud Subnet and update the Network Status.
func (s *VPCClusterScope) updateSubnetStatus(subnetDetails *vpcv1.Subnet, isControlPlane bool) (bool, error) {
	requeue := true
	if subnetDetails.Status != nil && *subnetDetails.Status == string(vpcv1.SubnetStatusAvailableConst) {
		requeue = false
	}

	resourceStatus := &infrav1beta2.ResourceStatus{
		ID:   *subnetDetails.ID,
		Name: subnetDetails.Name,
		// Ready status will be invert of the need to requeue
		Ready: !requeue,
	}
	if isControlPlane {
		s.SetResourceStatus(infrav1beta2.ResourceTypeControlPlaneSubnet, resourceStatus)
	} else {
		s.SetResourceStatus(infrav1beta2.ResourceTypeWorkerSubnet, resourceStatus)
	}
	return requeue, nil
}

// createSubnet creates a new VPC subnet.
func (s *VPCClusterScope) createSubnet(subnet infrav1beta2.Subnet, isControlPlane bool) error {
	// TODO(cjschaef): Move to webhook validation.
	if subnet.Zone == nil {
		return fmt.Errorf("error subnet zone must be defined for subnet %s", *subnet.Name)
	}

	// Created resources should be placed in the cluster Resource Group (not Network, if it exists).
	resourceGroupID, err := s.GetResourceGroupID()
	if err != nil {
		return fmt.Errorf("error retrieving resource group id for subnet creation: %w", err)
	} else if resourceGroupID == "" {
		return fmt.Errorf("error retrieving resource group id for resource group %s", s.IBMVPCCluster.Spec.ResourceGroup)
	}

	vpcID, err := s.GetVPCID()
	if err != nil {
		return fmt.Errorf("error retrieving vpc id for subnet creation: %w", err)
	}

	// NOTE(cjschaef): We likely will want to add support to use custom Address Prefixes.
	// For now, we rely on the API to assign us prefixes, as we request via IP count.
	var ipCount int64 = 256
	// We currnetly only support IP v4.
	ipVersion := vpcSubnetIPVersion4

	// Find or create a Public Gateway in this zone for the subnet, only one Public Gateway is required for each zone, for this cluster.
	// NOTE(cjschaef): We may need to add support to not attach Public Gateways to subnets.
	publicGateway, err := s.findOrCreatePublicGateway(*subnet.Zone)
	if err != nil {
		return fmt.Errorf("error failed to find or create public gateway for subnet %s: %w", *subnet.Name, err)
	}

	options := &vpcv1.CreateSubnetOptions{}
	options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
		IPVersion:             ptr.To(ipVersion),
		TotalIpv4AddressCount: ptr.To(ipCount),
		Name:                  subnet.Name,
		VPC: &vpcv1.VPCIdentity{
			ID: vpcID,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: subnet.Zone,
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: ptr.To(resourceGroupID),
		},
		PublicGateway: &vpcv1.PublicGatewayIdentity{
			ID: publicGateway.ID,
		},
	})

	// Create subnet.
	subnetDetails, _, err := s.VPCClient.CreateSubnet(options)
	if err != nil {
		return fmt.Errorf("error unknown failure creating vpc subnet: %w", err)
	}
	if subnetDetails == nil || subnetDetails.ID == nil || subnetDetails.CRN == nil {
		return fmt.Errorf("error failed creating subnet: %s", *subnet.Name)
	}

	// Initially populate subnet's status.
	resourceStatus := &infrav1beta2.ResourceStatus{
		ID:    *subnetDetails.ID,
		Name:  subnetDetails.Name,
		Ready: false,
	}
	if isControlPlane {
		s.SetResourceStatus(infrav1beta2.ResourceTypeControlPlaneSubnet, resourceStatus)
	} else {
		s.SetResourceStatus(infrav1beta2.ResourceTypeWorkerSubnet, resourceStatus)
	}

	// Add a tag to the subnet for the cluster.
	err = s.TagResource(s.IBMVPCCluster.Name, *subnetDetails.CRN)
	if err != nil {
		return fmt.Errorf("error failed to tag subnet %s: %w", *subnetDetails.Name, err)
	}

	return nil
}

// findOrCreatePublicGateway will attempt to find if there is an existing Public Gateway for a specific zone, for the cluster (in cluster's Resource Group and VPC), or create a new one. Only one Public Gateway is required in each zone, for any subnets in that zone.
func (s *VPCClusterScope) findOrCreatePublicGateway(zone string) (*vpcv1.PublicGateway, error) {
	publicGatewayName := fmt.Sprintf("%s-%s", *s.GetServiceName(infrav1beta2.ResourceTypePublicGateway), zone)
	// We will use the cluster Resource Group ID, as we expect to create all resources (Public Gateways and Subnets) in that Resource Group.
	resourceGroupID, err := s.GetResourceGroupID()
	if err != nil {
		return nil, fmt.Errorf("error unknown failure retrieving resource group id for public gateway: %w", err)
	}
	publicGateway, err := s.VPCClient.GetVPCPublicGatewayByName(publicGatewayName, resourceGroupID)
	if err != nil {
		return nil, fmt.Errorf("error unknown failure retrieving public gateway for zone %s: %w", zone, err)
	}

	// If we found the Public Gateway, with an ID, for the zone, return it.
	// NOTE(cjschaef): We may wish to confirm the PublicGateway, by checking Tags (Global Tagging), but this might be sufficient, as we don't expect to have duplicate PG's or existing PG's, as we wouldn't create subnets and PG's for existing Network Infrastructure.
	if publicGateway != nil && publicGateway.ID != nil {
		return publicGateway, nil
	}

	// Otherwise, create a new Public Gateway for the zone.
	vpcID, err := s.GetVPCID()
	if err != nil {
		return nil, fmt.Errorf("error failed retrieving vpc id for public gateway creation: %w", err)
	}
	if vpcID == nil {
		return nil, fmt.Errorf("error failed to retrieve vpc id for public gateway creation")
	}

	publicGatewayDetails, _, err := s.VPCClient.CreatePublicGateway(&vpcv1.CreatePublicGatewayOptions{
		Name: ptr.To(publicGatewayName),
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: ptr.To(resourceGroupID),
		},
		VPC: &vpcv1.VPCIdentity{
			ID: vpcID,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: ptr.To(zone),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error unknown failure creating public gateway: %w", err)
	}
	if publicGatewayDetails == nil || publicGatewayDetails.ID == nil || publicGatewayDetails.CRN == nil {
		return nil, fmt.Errorf("error failed creating public gateway for zone %s", zone)
	}

	s.V(3).Info("created public gateway", "id", publicGatewayDetails.ID)

	// Add a tag to the public gateway for the cluster
	err = s.TagResource(s.IBMVPCCluster.Name, *publicGatewayDetails.CRN)
	if err != nil {
		return nil, fmt.Errorf("error failed to tag public gateway %s: %w", *publicGatewayDetails.Name, err)
	}

	return publicGatewayDetails, nil
}

// ReconcileSecurityGroups will attempt to reconcile the defined SecurityGroups and their SecurityGroupRules. Our best option is to perform a first set of passes, creating all the SecurityGroups first, then reconcile the SecurityGroupRules after that, as the SecuirtyGroupRules could be dependent on an IBM Cloud Security Group that must be created first.
func (s *VPCClusterScope) ReconcileSecurityGroups() (bool, error) {
	// If no Security Groups were supplied, we have nothing to do.
	if len(s.IBMVPCCluster.Spec.Network.SecurityGroups) == 0 {
		return false, nil
	}

	// Reconcile each Security Group first, process rules later.
	for _, securityGroup := range s.IBMVPCCluster.Spec.Network.SecurityGroups {
		if err := s.reconcileSecurityGroup(securityGroup); err != nil {
			return false, fmt.Errorf("error failed reonciling security groups: %w", err)
		}
	}

	// Reconcile each Security Groups's Rules.
	requeue := false
	for _, securityGroup := range s.IBMVPCCluster.Spec.Network.SecurityGroups {
		if requiresRequeue, err := s.reconcileSecurityGroupRules(securityGroup); err != nil {
			return false, fmt.Errorf("error failed reconciling security group rules: %w", err)
		} else if requiresRequeue {
			s.V(3).Info("requeuing for security group rules")
			requeue = true
		}
	}

	return requeue, nil
}

// reconcileSecurityGroup will attempt to reconcile a defined SecurityGroup. By design, we confirm the IBM Cloud Security Group exists first, before attempting to reconcile the defined SecurityGroupRules.
func (s *VPCClusterScope) reconcileSecurityGroup(securityGroup infrav1beta2.VPCSecurityGroup) error {
	var securityGroupID *string
	// If Security Group already has an ID defined, use that for lookup.
	if securityGroup.ID != nil {
		securityGroupID = securityGroup.ID
	} else {
		if securityGroup.Name == nil {
			return fmt.Errorf("error securityGroup has no name or id")
		}
		// Check the Status if an ID is already available for the Security Group.
		if id := s.getSecurityGroupIDFromStatus(*securityGroup.Name); id != nil {
			securityGroupID = id
		} else {
			// Otherwise, attempt to lookup Security Group by name.
			if securityGroupDetails, err := s.VPCClient.GetSecurityGroupByName(*securityGroup.Name); err != nil {
				// If the Security Group was not found, we expect it doesn't exist yet, otherwise result in an error.
				if _, ok := err.(*vpc.SecurityGroupByNameNotFound); !ok {
					return fmt.Errorf("error failed lookup of security group by name: %w", err)
				}
			} else if securityGroupDetails != nil {
				// If the Security Group was found, update Status with current details.
				// Security Groups do not have a status, so we assume if it exists, it is ready.
				s.SetResourceStatus(infrav1beta2.ResourceTypeSecurityGroup, &infrav1beta2.ResourceStatus{
					ID:    *securityGroupDetails.ID,
					Name:  securityGroupDetails.Name,
					Ready: true,
				})
				return nil
			}
		}
	}

	// If we have an ID for the SecurityGroup, we can check the status.
	if securityGroupID != nil {
		s.V(3).Info("checking security group status", "securityGroupName", securityGroup.Name, "securityGroupID", securityGroupID)
		securityGroupDetails, _, err := s.VPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{
			ID: securityGroupID,
		})
		if err != nil {
			return fmt.Errorf("error failed lookup of security group: %w", err)
		} else if securityGroupDetails == nil {
			// The Security Group cannot be found by ID, it was removed or didn't exist.
			// TODO(cjschaef): We may wish to clear the ID's to get a new Security Group created, but for now we return an error.
			return fmt.Errorf("error could not find security group with id=%s", *securityGroupID)
		}

		// Security Groups do not have a status, so we assume if it exists, it is ready.
		s.SetResourceStatus(infrav1beta2.ResourceTypeSecurityGroup, &infrav1beta2.ResourceStatus{
			ID:    *securityGroupID,
			Name:  securityGroupDetails.Name,
			Ready: true,
		})
		return nil
	}

	// If we don't have an ID at this point, we assume we need to create the Security Group.
	vpcID, err := s.GetVPCID()
	if err != nil {
		return fmt.Errorf("error retrieving vpc id for security group creation: %w", err)
	}
	resourceGroupID, err := s.GetResourceGroupID()
	if err != nil {
		return fmt.Errorf("error retrieving resource id for security group creation: %w", err)
	}
	createOptions := &vpcv1.CreateSecurityGroupOptions{
		Name: securityGroup.Name,
		VPC: &vpcv1.VPCIdentityByID{
			ID: vpcID,
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentityByID{
			ID: ptr.To(resourceGroupID),
		},
	}
	securityGroupDetails, _, err := s.VPCClient.CreateSecurityGroup(createOptions)
	if err != nil {
		s.V(3).Error(err, "error creating security group", "securityGroupName", securityGroup.Name)
		return fmt.Errorf("error failed to create security group: %w", err)
	}
	if securityGroupDetails == nil {
		s.V(3).Info("error failed creating security group", "securityGroupName", securityGroup.Name)
		return fmt.Errorf("error failed creating security group")
	}

	// Security Groups do not have a status, so just assume they are ready immediately after creation.
	s.SetResourceStatus(infrav1beta2.ResourceTypeSecurityGroup, &infrav1beta2.ResourceStatus{
		ID:    *securityGroupDetails.ID,
		Name:  securityGroupDetails.Name,
		Ready: true,
	})

	// NOTE: This tagging is only attempted once. We may wish to refactor in case this single attempt fails.
	// Add a tag to the Security Group for the cluster.
	err = s.TagResource(s.IBMVPCCluster.Name, *securityGroupDetails.CRN)
	if err != nil {
		return fmt.Errorf("error failed to tag security group %s: %w", *securityGroupDetails.CRN, err)
	}

	return nil
}

// reconcile SecurityGroupRules will attempt to reconcile the set of defined SecurityGroupRules for a SecurityGroup, one Rule at a time. Each defined Rule can contain multiple remotes, requiring a unique IBM Cloud Security Group Rule, based on the expected traffic direction, inbound (Source) or outbound (Destination).
func (s *VPCClusterScope) reconcileSecurityGroupRules(securityGroup infrav1beta2.VPCSecurityGroup) (bool, error) {
	// If the SecurityGroup has no rules, we have nothing more to do for this Security Group.
	if len(securityGroup.Rules) == 0 {
		return false, nil
	}

	// Assume that the securityGroup exists in Status, if it doesn't then it should be re-reconciled. Attempt to find it by name and then ID.
	var securityGroupID *string
	if securityGroup.Name != nil {
		securityGroupID = s.getSecurityGroupIDFromStatus(*securityGroup.Name)
	} else if securityGroup.ID != nil {
		// TODO(cjschaef): Since this does not rely on Status, this could become an issue.
		securityGroupID = securityGroup.ID
	}

	if securityGroupID == nil {
		s.V(3).Info("security group not found, requeue", "securityGroup", securityGroup)
		return true, nil
	}

	// Reconcile each SecurityGroupRule in the SecurityGroup.
	for _, securityGroupRule := range securityGroup.Rules {
		s.V(3).Info("reconcile security group rule", "securityGroupID", securityGroupID)
		if err := s.reconcileSecurityGroupRule(*securityGroupID, *securityGroupRule); err != nil {
			return false, fmt.Errorf("error failed to reconcile security group rule: %w", err)
		}
	}

	// Since Security Group Rules have no status, assume all Rules have been reconciled (they exist or were created).
	return false, nil
}

// reconcileSecurityGroupRule will attempt to reconcile a defined SecurityGroupRule, with one or more Remotes, for a SecurityGroup. If the IBM Cloud Security Group contains no Rules, simply attempt to create the defined Rule (via the Remote(s) provided).
func (s *VPCClusterScope) reconcileSecurityGroupRule(securityGroupID string, securityGroupRule infrav1beta2.VPCSecurityGroupRule) error {
	existingSecurityGroupRuleIntfs, _, err := s.VPCClient.ListSecurityGroupRules(&vpcv1.ListSecurityGroupRulesOptions{
		SecurityGroupID: ptr.To(securityGroupID),
	})
	if err != nil {
		return fmt.Errorf("error failed listing security group rules during reconcile of security group id=%s: %w", securityGroupID, err)
	}

	// If the Security Group has no Rules at all, we simply create all the Rules
	if existingSecurityGroupRuleIntfs == nil || len(existingSecurityGroupRuleIntfs.Rules) == 0 {
		s.V(3).Info("Creating security group rules for security group", "securityGroupID", securityGroupID)
		err := s.createSecurityGroupRuleAllRemotes(securityGroupID, securityGroupRule)
		if err != nil {
			return fmt.Errorf("error failed creating all security group rule remotes: %w", err)
		}
		s.V(3).Info("Created security group rules", "securityGroupID", securityGroupID, "securityGroupRule", securityGroupRule)

		// Security Group Rules do not have a Status, so assume they are ready immediately.
		return nil
	}

	// Validate the Security Group Rule(s) exist or were created.
	if err := s.findOrCreateSecurityGroupRule(securityGroupID, securityGroupRule, existingSecurityGroupRuleIntfs); err != nil {
		return fmt.Errorf("error failed to find or create security group rule: %w", err)
	}
	return nil
}

// findOrCreateSecurityGroupRule will attempt to match up the SecurityGroupRule's Remote(s) (multiple Remotes can be supplied per Rule definition), and will create any missing IBM Cloud Security Group Rules based on the SecurityGroupRule and Remote(s). Remotes are defined either by a Destination (outbound) or a Source (inbound), which defines the type of IBM Cloud Security Group Rule that should exist or be created.
func (s *VPCClusterScope) findOrCreateSecurityGroupRule(securityGroupID string, securityGroupRule infrav1beta2.VPCSecurityGroupRule, existingSecurityGroupRules *vpcv1.SecurityGroupRuleCollection) error { //nolint: gocyclo
	// Use either the SecurityGroupRule.Destination or SecurityGroupRule.Source for further details based on SecurityGroupRule.Direction
	var securityGroupRulePrototype infrav1beta2.VPCSecurityGroupRulePrototype
	switch securityGroupRule.Direction {
	case infrav1beta2.VPCSecurityGroupRuleDirectionInbound:
		securityGroupRulePrototype = *securityGroupRule.Source
	case infrav1beta2.VPCSecurityGroupRuleDirectionOutbound:
		securityGroupRulePrototype = *securityGroupRule.Destination
	default:
		return fmt.Errorf("error unsupported SecurityGroupRuleDirection defined")
	}

	s.V(3).Info("checking security group rules for security group", "securityGroupID", securityGroupID)

	// Each defined SecurityGroupRule can have multiple Remotes specified, each signifying a separate Security Group Rule (with the same Action, Direction, etc.)
	for _, remote := range securityGroupRulePrototype.Remotes {
		remoteMatch := false
		for _, existingRuleIntf := range existingSecurityGroupRules.Rules {
			// Perform analysis of the existingRuleIntf, based on its Protocol type, further analysis is performed based on remaining attributes to find if the specific Rule and Remote match
			switch reflect.TypeOf(existingRuleIntf).String() {
			case infrav1beta2.VPCSecurityGroupRuleProtocolAllType:
				// If our Remote doesn't define all Protocols, we don't need further checks, move on to next Rule
				if securityGroupRulePrototype.Protocol != infrav1beta2.VPCSecurityGroupRuleProtocolAll {
					continue
				}
				existingRule := existingRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
				// If the Remote doesn't have the same Direction as the Rule, no further checks are necessary
				if securityGroupRule.Direction != infrav1beta2.VPCSecurityGroupRuleDirection(*existingRule.Direction) {
					continue
				}
				if found, err := s.checkSecurityGroupRuleProtocolAll(securityGroupRulePrototype, remote, existingRule); err != nil {
					return fmt.Errorf("error failure checking security group rule protocol all: %w", err)
				} else if found {
					// If we found the matching IBM Cloud Security Group Rule for the defined SecurityGroupRule and Remote, we can stop checking IBM Cloud Security Group Rules for this remote and move onto the next remote.
					// The expectation is that only one IBM Cloud Security Group Rule will match, but if at least one matches the defined SecurityGroupRule, that is sufficient.
					s.V(3).Info("security group rule all protocol match found")
					remoteMatch = true
					break
				}
			case infrav1beta2.VPCSecurityGroupRuleProtocolIcmpType:
				// If our Remote doesn't define ICMP Protocol, we don't need further checks, move on to next Rule
				if securityGroupRulePrototype.Protocol != infrav1beta2.VPCSecurityGroupRuleProtocolIcmp {
					continue
				}
				existingRule := existingRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
				// If the Remote doesn't have the same Direction as the Rule, no further checks are necessary
				if securityGroupRule.Direction != infrav1beta2.VPCSecurityGroupRuleDirection(*existingRule.Direction) {
					continue
				}
				if found, err := s.checkSecurityGroupRuleProtocolIcmp(securityGroupRulePrototype, remote, existingRule); err != nil {
					return fmt.Errorf("error failure checking security group rule protocol icmp: %w", err)
				} else if found {
					// If we found the matching IBM Cloud Security Group Rule for the defined SecurityGroupRule and Remote, we can stop checking IBM Cloud Security Group Rules for this remote and move onto the next remote.
					s.V(3).Info("security group rule icmp match found")
					remoteMatch = true
					break
				}
			case infrav1beta2.VPCSecurityGroupRuleProtocolTcpudpType:
				// If our Remote doesn't define TCP/UDP Protocol, we don't need further checks, move on to next Rule
				if securityGroupRulePrototype.Protocol != infrav1beta2.VPCSecurityGroupRuleProtocolTCP && securityGroupRulePrototype.Protocol != infrav1beta2.VPCSecurityGroupRuleProtocolUDP {
					continue
				}
				existingRule := existingRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
				// If the Remote doesn't have the same Direction as the Rule, no further checks are necessary
				if securityGroupRule.Direction != infrav1beta2.VPCSecurityGroupRuleDirection(*existingRule.Direction) {
					continue
				}
				if found, err := s.checkSecurityGroupRuleProtocolTcpudp(securityGroupRulePrototype, remote, existingRule); err != nil {
					return fmt.Errorf("error failure checking security group rule protocol tcp-udp: %w", err)
				} else if found {
					// If we found the matching IBM Cloud Security Group Rule for the defined SecurityGroupRule and Remote, we can stop checking IBM Cloud Security Group Rules for this remote and move onto the next remote.
					s.V(3).Info("security group rule tcp/udp match found")
					remoteMatch = true
					break
				}
			default:
				// This is an unexpected IBM Cloud Security Group Rule Prototype, log it and move on
				s.V(3).Info("unexpected security group rule prototype", "securityGroupRulePrototype", reflect.TypeOf(existingRuleIntf).String())
			}
		}

		// If we did not find a matching SecurityGroupRule for this defined Remote, create one now.
		if !remoteMatch {
			err := s.createSecurityGroupRule(securityGroupID, securityGroupRule, remote)
			if err != nil {
				return fmt.Errorf("error failure creating security group rule: %w", err)
			}
		}
	}
	return nil
}

// checkSecurityGroupRuleProtocolAll analyzes an IBM Cloud Security Group Rule designated for 'all' protocols, to verify if the supplied Rule and Remote match the attributes from the existing 'ProtocolAll' Rule.
func (s *VPCClusterScope) checkSecurityGroupRuleProtocolAll(_ infrav1beta2.VPCSecurityGroupRulePrototype, securityGroupRuleRemote infrav1beta2.VPCSecurityGroupRuleRemote, existingRule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll) (bool, error) {
	if exists, err := s.checkSecurityGroupRulePrototypeRemote(securityGroupRuleRemote, existingRule.Remote); err != nil {
		return false, fmt.Errorf("error failed checking security group rule all remote: %w", err)
	} else if exists {
		s.V(3).Info("security group rule all protocols match")
		return true, nil
	}
	return false, nil
}

// checkSecurityGroupRuleProtocolIcmp analyzes an IBM Cloud Security Group Rule designated for 'icmp' protocol, to verify if the supplied Rule and Remote match the attributes from the existing 'ProtocolIcmp' Rule.
func (s *VPCClusterScope) checkSecurityGroupRuleProtocolIcmp(securityGroupRulePrototype infrav1beta2.VPCSecurityGroupRulePrototype, securityGroupRuleRemote infrav1beta2.VPCSecurityGroupRuleRemote, existingRule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp) (bool, error) {
	if exists, err := s.checkSecurityGroupRulePrototypeRemote(securityGroupRuleRemote, existingRule.Remote); err != nil {
		return false, fmt.Errorf("error failed checking security group rule icmp remote: %w", err)
	} else if !exists {
		return false, nil
	}
	// If ICMPCode is set, then ICMPType must also be set, via kubebuilder specifications
	if securityGroupRulePrototype.ICMPCode != nil && securityGroupRulePrototype.ICMPType != nil {
		// If the existingRule Code and Type are both equal to the securityGroupRulePrototype's ICMPType and ICMPCode, the existingRule matches our definition for ICMP in securityGroupRulePrototype.
		if *securityGroupRulePrototype.ICMPCode == *existingRule.Code && *securityGroupRulePrototype.ICMPType == *existingRule.Type {
			s.V(3).Info("security group rule icmp code and type match", "icmpCode", *existingRule.Code, "icmpType", *existingRule.Type)
			return true, nil
		}
	} else if existingRule.Code == nil && existingRule.Type == nil {
		s.V(3).Info("security group rule unset icmp matches")
		return true, nil
	}
	return false, nil
}

// checkSecurityGroupRuleProtocolTcpudp analyzes an IBM Cloud Security Group Rule designated for either 'tcp' or 'udp' protocols, to verify if the supplied Rule and Remote match the attributes from the existing 'ProtocolTcpudp' Rule.
func (s *VPCClusterScope) checkSecurityGroupRuleProtocolTcpudp(securityGroupRulePrototype infrav1beta2.VPCSecurityGroupRulePrototype, securityGroupRuleRemote infrav1beta2.VPCSecurityGroupRuleRemote, existingRule *vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp) (bool, error) {
	// Check the protocol next, either TCP or UDP, to verify it matches
	if securityGroupRulePrototype.Protocol != infrav1beta2.VPCSecurityGroupRuleProtocol(*existingRule.Protocol) {
		return false, nil
	}

	if exists, err := s.checkSecurityGroupRulePrototypeRemote(securityGroupRuleRemote, existingRule.Remote); err != nil {
		return false, fmt.Errorf("error failed checking security group rule tcp-udp remote: %w", err)
	} else if exists {
		// If PortRange is set, verify whether the MinimumPort and MaximumPort match the existingRule's values, if they are set.
		if securityGroupRulePrototype.PortRange != nil {
			if existingRule.PortMin != nil && securityGroupRulePrototype.PortRange.MinimumPort == *existingRule.PortMin && existingRule.PortMax != nil && securityGroupRulePrototype.PortRange.MaximumPort == *existingRule.PortMax {
				s.V(3).Info("security group rule port range matches", "ruleID", *existingRule.ID, "portMin", *existingRule.PortMin, "portMax", *existingRule.PortMax)
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *VPCClusterScope) checkSecurityGroupRulePrototypeRemote(securityGroupRuleRemote infrav1beta2.VPCSecurityGroupRuleRemote, existingRemote vpcv1.SecurityGroupRuleRemoteIntf) (bool, error) { //nolint: gocyclo
	// NOTE(cjschaef): We only currently monitor Remote, not Local, as we don't support defining Local in SecurityGroup/SecurityGroupRule.
	switch securityGroupRuleRemote.RemoteType {
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeCIDR:
		cidrRule := existingRemote.(*vpcv1.SecurityGroupRuleRemote)
		if cidrRule.CIDRBlock == nil {
			return false, nil
		}
		subnetDetails, err := s.VPCClient.GetVPCSubnetByName(*securityGroupRuleRemote.CIDRSubnetName)
		if err != nil {
			return false, fmt.Errorf("error failed getting subnet by name for security group rule: %w", err)
		} else if subnetDetails == nil {
			return false, fmt.Errorf("error failed getting subnet by name for security group rule")
		}
		if *subnetDetails.Ipv4CIDRBlock == *cidrRule.CIDRBlock {
			s.V(3).Info("security group rule remote cidr's match", "remoteCIDR", *cidrRule.CIDRBlock)
			return true, nil
		}
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAddress:
		ipRule := existingRemote.(*vpcv1.SecurityGroupRuleRemote)
		if ipRule.Address == nil {
			return false, nil
		}
		if *securityGroupRuleRemote.Address == *ipRule.Address {
			s.V(3).Info("security group rule remote addresses match", "remoteAddress", *ipRule.Address)
			return true, nil
		}
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeSG:
		sgRule := existingRemote.(*vpcv1.SecurityGroupRuleRemote)
		if sgRule.Name == nil {
			return false, nil
		}

		// We can compare the SecurityGroup details from the securityGroupRemote and SecurityGroupRuleRemoteSecurityGroupReference, if those values are available
		// Option #1. We can compare the Security Group Name (name is manditory for securityGroupRemote)
		// Option #2. We can compare the Security Group ID (may already have securityGroupRemote ID)
		// Option #3. We can compare the Security Group CRN (need ot lookup the CRN for securityGroupRemote)

		// Option #1: If the SecurityGroupRuleRemoteSecurityGroupReference has a name assigned, we can shortcut and simply check that
		if sgRule.Name != nil && *sgRule.Name == *securityGroupRuleRemote.SecurityGroupName {
			s.V(3).Info("security group rule remote security group name matches", "securityGroupRuleRemoteSecurityGroupName", *sgRule.Name)
			return true, nil
		}
		// Try to get the Security Group Id for quick lookup (from Network Status)
		var securityGroupDetails *vpcv1.SecurityGroup
		var err error
		if securityGroupID := s.getSecurityGroupIDFromStatus(*securityGroupRuleRemote.SecurityGroupName); securityGroupID != nil {
			// Option #2: If the SecurityGroupRuleRemoteSecurityGroupReference has an ID assigned, we can shortcut and simply check that
			if sgRule.ID != nil && *securityGroupID == *sgRule.ID {
				s.V(3).Info("security group rule remote security group id matches", "securityGroupRuleRemoteSecurityGroupID", *sgRule.ID)
				return true, nil
			}
			securityGroupDetails, _, err = s.VPCClient.GetSecurityGroup(&vpcv1.GetSecurityGroupOptions{
				ID: securityGroupID,
			})
		} else {
			securityGroupDetails, err = s.VPCClient.GetSecurityGroupByName(*securityGroupRuleRemote.SecurityGroupName)
		}
		if err != nil {
			return false, fmt.Errorf("error failed getting security group by name for security group rule: %w", err)
		} else if securityGroupDetails == nil {
			return false, fmt.Errorf("error failed getting security group by name for security group rule")
		}

		// Option #3: We check the SecurityGroupRuleRemoteSecurityGroupReference's CRN, if the Name and ID were not available
		if *securityGroupDetails.CRN == *sgRule.CRN {
			s.V(3).Info("security group rule remote security group crn matches", "securityGroupRuleRemoteSecurityGroupCRN", *securityGroupDetails.CRN)
			return true, nil
		}
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAny:
		ipRule := existingRemote.(*vpcv1.SecurityGroupRuleRemote)
		if ipRule.Address == nil {
			s.V(3).Info("security group rule remote has no address, defaults to any remote")
			return true, nil
		}
		if *ipRule.Address == infrav1beta2.CIDRBlockAny {
			s.V(3).Info("security group rule remote address matches %s", infrav1beta2.CIDRBlockAny)
			return true, nil
		}
	default:
		s.V(3).Info("unknown security group rule remote")
	}
	return false, nil
}

// createSecurityGroupRuleAllRemotes will create one or more IBM Cloud Security Group Rules for a specific SecurityGroup, based on the provided SecurityGroupRule and Remotes defined in the SecurityGroupRule definition (one or more Remotes can be defined per SecurityGroupRule definition).
func (s *VPCClusterScope) createSecurityGroupRuleAllRemotes(securityGroupID string, securityGroupRule infrav1beta2.VPCSecurityGroupRule) error {
	var remotes []infrav1beta2.VPCSecurityGroupRuleRemote
	switch securityGroupRule.Direction {
	case infrav1beta2.VPCSecurityGroupRuleDirectionInbound:
		remotes = securityGroupRule.Source.Remotes
	case infrav1beta2.VPCSecurityGroupRuleDirectionOutbound:
		remotes = securityGroupRule.Destination.Remotes
	}
	for _, remote := range remotes {
		err := s.createSecurityGroupRule(securityGroupID, securityGroupRule, remote)
		if err != nil {
			return fmt.Errorf("error failed creating security group rule: %w", err)
		}
	}

	return nil
}

// createSecurityGroupRule will create a new IBM Cloud Security Group Rule for a specific Security Group, based on the provided SecurityGroupRule and Remote definitions.
func (s *VPCClusterScope) createSecurityGroupRule(securityGroupID string, securityGroupRule infrav1beta2.VPCSecurityGroupRule, remote infrav1beta2.VPCSecurityGroupRuleRemote) error {
	options := &vpcv1.CreateSecurityGroupRuleOptions{
		SecurityGroupID: &securityGroupID,
	}
	// Setup variables to use for logging details on the resulting IBM Cloud Security Group Rule creation options
	var securityGroupRulePrototype *infrav1beta2.VPCSecurityGroupRulePrototype
	if securityGroupRule.Direction == infrav1beta2.VPCSecurityGroupRuleDirectionInbound {
		securityGroupRulePrototype = securityGroupRule.Source
	} else {
		securityGroupRulePrototype = securityGroupRule.Destination
	}
	prototypeRemote, err := s.createSecurityGroupRuleRemote(remote)
	if err != nil {
		return fmt.Errorf("error failed to create security group rule remote: %w", err)
	}
	switch securityGroupRulePrototype.Protocol {
	case infrav1beta2.VPCSecurityGroupRuleProtocolAll:
		prototype := &vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAll{
			Direction: ptr.To(string(securityGroupRule.Direction)),
			Protocol:  ptr.To(string(securityGroupRulePrototype.Protocol)),
			Remote:    prototypeRemote,
		}
		options.SetSecurityGroupRulePrototype(prototype)
	case infrav1beta2.VPCSecurityGroupRuleProtocolIcmp:
		prototype := &vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolIcmp{
			Direction: ptr.To(string(securityGroupRule.Direction)),
			Protocol:  ptr.To(string(securityGroupRulePrototype.Protocol)),
			Remote:    prototypeRemote,
		}
		// If ICMP Code or Type is specified, both must be, enforced by kubebuilder
		if securityGroupRulePrototype.ICMPCode != nil && securityGroupRulePrototype.ICMPType != nil {
			prototype.Code = securityGroupRulePrototype.ICMPCode
			prototype.Type = securityGroupRulePrototype.ICMPType
		}
		options.SetSecurityGroupRulePrototype(prototype)
	// TCP and UDP use the same Prototype, simply with different Protocols, which is agnostic in code
	case infrav1beta2.VPCSecurityGroupRuleProtocolTCP, infrav1beta2.VPCSecurityGroupRuleProtocolUDP:
		prototype := &vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolTcpudp{
			Direction: ptr.To(string(securityGroupRule.Direction)),
			Protocol:  ptr.To(string(securityGroupRulePrototype.Protocol)),
			Remote:    prototypeRemote,
		}
		if securityGroupRulePrototype.PortRange != nil {
			prototype.PortMin = ptr.To(securityGroupRulePrototype.PortRange.MinimumPort)
			prototype.PortMax = ptr.To(securityGroupRulePrototype.PortRange.MaximumPort)
		}
		options.SetSecurityGroupRulePrototype(prototype)
	default:
		// This should not be possible, provided the strict kubebuilder enforcements
		return fmt.Errorf("error failed creating security group rule, unknown protocol")
	}

	s.V(3).Info("Creating Security Group Rule for Security Group", "securityGroupID", securityGroupID, "direction", securityGroupRule.Direction, "protocol", securityGroupRulePrototype.Protocol, "prototypeRemote", prototypeRemote)
	securityGroupRuleIntfDetails, _, err := s.VPCClient.CreateSecurityGroupRule(options)
	if err != nil {
		return fmt.Errorf("error unexpected failure creating security group rule: %w", err)
	} else if securityGroupRuleIntfDetails == nil {
		return fmt.Errorf("error failed creating security group rule")
	}

	// Typecast the resulting SecurityGroupRuleIntf, to retrieve the ID for logging
	var ruleID *string
	switch reflect.TypeOf(securityGroupRuleIntfDetails).String() {
	case infrav1beta2.VPCSecurityGroupRuleProtocolAllType:
		rule := securityGroupRuleIntfDetails.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
		ruleID = rule.ID
	case infrav1beta2.VPCSecurityGroupRuleProtocolIcmpType:
		rule := securityGroupRuleIntfDetails.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
		ruleID = rule.ID
	case infrav1beta2.VPCSecurityGroupRuleProtocolTcpudpType:
		rule := securityGroupRuleIntfDetails.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
		ruleID = rule.ID
	}
	s.V(3).Info("Created Security Group Rule", "ruleID", ruleID)
	return nil
}

// createSecurityGroupRuleRemote will create an IBM Cloud SecurityGroupRuleRemotePrototype, which defines the Remote details for an IBM Cloud Security Group Rule, provided by the SecurityGroupRuleRemote. Lookups of Security Group CRN's, by Name, or Subnet CIDRBlock's, by Name, allows the use of CAPI created resources to be defined in the SecurityGroupRuleRemote, when the CRN or CIDRBlock are unknown (runtime defined).
func (s *VPCClusterScope) createSecurityGroupRuleRemote(remote infrav1beta2.VPCSecurityGroupRuleRemote) (*vpcv1.SecurityGroupRuleRemotePrototype, error) {
	remotePrototype := &vpcv1.SecurityGroupRuleRemotePrototype{}
	switch remote.RemoteType {
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAny:
		remotePrototype.CIDRBlock = ptr.To(infrav1beta2.CIDRBlockAny)
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeCIDR:
		// As we nned the Subnet CIDR block, we have to perform an IBM Cloud API call either way, so simply make the call using the item we know, the Name
		subnetDetails, err := s.VPCClient.GetVPCSubnetByName(*remote.CIDRSubnetName)
		if err != nil {
			return nil, fmt.Errorf("error failed lookup of subnet during security group rule remote creation: %w", err)
		} else if subnetDetails == nil {
			return nil, fmt.Errorf("error failed lookup of subnet during security group rule remote creation")
		}
		remotePrototype.CIDRBlock = subnetDetails.Ipv4CIDRBlock
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeAddress:
		remotePrototype.Address = remote.Address
	case infrav1beta2.VPCSecurityGroupRuleRemoteTypeSG:
		// As we need the Security Group CRN, we have to perform an IBM Cloud API call either way, so simply make the call using the item we know, the Name
		securityGroupDetails, err := s.VPCClient.GetSecurityGroupByName(*remote.SecurityGroupName)
		if err != nil {
			return nil, fmt.Errorf("error failed lookup of security group during security group rule remote creation: %w", err)
		} else if securityGroupDetails == nil {
			return nil, fmt.Errorf("error failed lookup of security group during security group rule remote creation")
		}
		remotePrototype.CRN = securityGroupDetails.CRN
	default:
		// This should not be possible, given the strict kubebuilder enforcements
		return nil, fmt.Errorf("error failed creating security group rule remote")
	}

	return remotePrototype, nil
}
