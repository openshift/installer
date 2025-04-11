/*
Copyright 2022 The Kubernetes Authors.

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

package resourcecontroller

import (
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
)

const (
	// PowerVSResourceID is Power VS power-iaas service id, can be retrieved using ibmcloud cli
	// ibmcloud catalog service power-iaas.
	PowerVSResourceID = "abd259f0-9990-11e8-acc8-b9f54a8f1661"

	// PowerVSResourcePlanID is Power VS power-iaas plan id, can be retrieved using ibmcloud cli
	// ibmcloud catalog service power-iaas.
	PowerVSResourcePlanID = "f165dd34-3a40-423b-9d95-e90a23f724dd"

	// CosResourceID is IBM COS service id, can be retrieved using ibmcloud cli
	// ibmcloud catalog service cloud-object-storage.
	CosResourceID = "dff97f5c-bc5e-4455-b470-411c3edbe49c"

	// CosResourcePlanID is IBM COS plan id, can be retrieved using ibmcloud cli
	// ibmcloud catalog service cloud-object-storage.
	CosResourcePlanID = "1e4e33e4-cfa6-4f12-9016-be594a6d5f87"
)

// Service holds the IBM Cloud Resource Controller Service specific information.
type Service struct {
	client *resourcecontrollerv2.ResourceControllerV2
}

// ServiceOptions holds the IBM Cloud Resource Controller Service Options specific information.
type ServiceOptions struct {
	*resourcecontrollerv2.ResourceControllerV2Options
}

// SetServiceURL sets the service URL.
func (s *Service) SetServiceURL(url string) error {
	return s.client.SetServiceURL(url)
}

// GetServiceURL will get the service URL.
func (s *Service) GetServiceURL() string {
	return s.client.GetServiceURL()
}

// ListResourceInstances will list all the resource instances.
func (s *Service) ListResourceInstances(listResourceInstancesOptions *resourcecontrollerv2.ListResourceInstancesOptions) (result *resourcecontrollerv2.ResourceInstancesList, response *core.DetailedResponse, err error) {
	return s.client.ListResourceInstances(listResourceInstancesOptions)
}

// GetResourceInstance will get the resource instance.
func (s *Service) GetResourceInstance(getResourceInstanceOptions *resourcecontrollerv2.GetResourceInstanceOptions) (result *resourcecontrollerv2.ResourceInstance, response *core.DetailedResponse, err error) {
	return s.client.GetResourceInstance(getResourceInstanceOptions)
}

// CreateResourceInstance creates the resource instance.
func (s *Service) CreateResourceInstance(options *resourcecontrollerv2.CreateResourceInstanceOptions) (*resourcecontrollerv2.ResourceInstance, *core.DetailedResponse, error) {
	return s.client.CreateResourceInstance(options)
}

// DeleteResourceInstance deletes the resource instance.
func (s *Service) DeleteResourceInstance(options *resourcecontrollerv2.DeleteResourceInstanceOptions) (*core.DetailedResponse, error) {
	return s.client.DeleteResourceInstance(options)
}

// GetServiceInstance returns service instance with given name or id. If not found, returns nil.
// TODO: Combine GetSreviceInstance() and GetInstanceByName().
func (s *Service) GetServiceInstance(id, name string, zone *string) (*resourcecontrollerv2.ResourceInstance, error) {
	var serviceInstancesList []resourcecontrollerv2.ResourceInstance
	f := func(start string) (bool, string, error) {
		listServiceInstanceOptions := &resourcecontrollerv2.ListResourceInstancesOptions{
			ResourceID:     ptr.To(PowerVSResourceID),
			ResourcePlanID: ptr.To(PowerVSResourcePlanID),
		}
		if id != "" {
			listServiceInstanceOptions.GUID = &id
		}
		if name != "" {
			listServiceInstanceOptions.Name = &name
		}
		if start != "" {
			listServiceInstanceOptions.Start = &start
		}

		serviceInstances, _, err := s.client.ListResourceInstances(listServiceInstanceOptions)
		if err != nil {
			return false, "", err
		}
		if serviceInstances != nil {
			if zone != nil && *zone != "" {
				for _, resource := range serviceInstances.Resources {
					if *resource.RegionID == *zone {
						serviceInstancesList = append(serviceInstancesList, resource)
					}
				}
			} else {
				serviceInstancesList = append(serviceInstancesList, serviceInstances.Resources...)
			}

			nextURL, err := serviceInstances.GetNextStart()
			if err != nil {
				return false, "", err
			}
			if nextURL == nil {
				return true, "", nil
			}
			return false, *nextURL, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, fmt.Errorf("error listing service instances %v", err)
	}
	switch len(serviceInstancesList) {
	case 0:
		return nil, nil
	case 1:
		return &serviceInstancesList[0], nil
	default:
		errStr := fmt.Errorf("there exist more than one service instance ID with same name %s, Try setting serviceInstance.ID", name)
		return nil, errStr
	}
}

// GetInstanceByName returns instance with given name, planID and resourceID. If not found, returns nil.
func (s *Service) GetInstanceByName(name, resourceID, planID string) (*resourcecontrollerv2.ResourceInstance, error) {
	var serviceInstancesList []resourcecontrollerv2.ResourceInstance
	f := func(start string) (bool, string, error) {
		listServiceInstanceOptions := &resourcecontrollerv2.ListResourceInstancesOptions{
			Name:           &name,
			ResourceID:     ptr.To(resourceID),
			ResourcePlanID: ptr.To(planID),
		}
		if start != "" {
			listServiceInstanceOptions.Start = &start
		}

		serviceInstances, _, err := s.client.ListResourceInstances(listServiceInstanceOptions)
		if err != nil {
			return false, "", err
		}
		if serviceInstances != nil {
			serviceInstancesList = append(serviceInstancesList, serviceInstances.Resources...)
			nextURL, err := serviceInstances.GetNextStart()
			if err != nil {
				return false, "", err
			}
			if nextURL == nil {
				return true, "", nil
			}
			return false, *nextURL, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, fmt.Errorf("error listing COS instances %v", err)
	}
	switch len(serviceInstancesList) {
	case 0:
		return nil, nil
	case 1:
		return &serviceInstancesList[0], nil
	default:
		errStr := fmt.Errorf("there exist more than one COS instance ID with same name %s, Try setting serviceInstance.ID", name)
		return nil, errStr
	}
}

// CreateResourceKey creates a new resource key.
func (s *Service) CreateResourceKey(options *resourcecontrollerv2.CreateResourceKeyOptions) (*resourcecontrollerv2.ResourceKey, *core.DetailedResponse, error) {
	return s.client.CreateResourceKey(options)
}

// NewService returns a new service for the IBM Cloud Resource Controller api client.
func NewService(options ServiceOptions) (ResourceController, error) {
	if options.ResourceControllerV2Options == nil {
		options.ResourceControllerV2Options = &resourcecontrollerv2.ResourceControllerV2Options{}
	}
	if options.Authenticator == nil {
		auth, err := authenticator.GetAuthenticator()
		if err != nil {
			return nil, err
		}
		options.Authenticator = auth
	}
	service, err := resourcecontrollerv2.NewResourceControllerV2(options.ResourceControllerV2Options)
	if err != nil {
		return nil, err
	}
	return &Service{
		client: service,
	}, nil
}
