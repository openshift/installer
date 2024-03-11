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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
)

// Service holds the IBM Cloud Resource Controller Service specific information.
type Service struct {
	client *resourcecontrollerv2.ResourceControllerV2
}

// ServiceOptions holds the IBM Cloud Resource Controller Service Options specific information.
type ServiceOptions struct {
	*resourcecontrollerv2.ResourceControllerV2Options
}

// ListResourceInstances will list all the resorce instances.
func (s *Service) ListResourceInstances(listResourceInstancesOptions *resourcecontrollerv2.ListResourceInstancesOptions) (result *resourcecontrollerv2.ResourceInstancesList, response *core.DetailedResponse, err error) {
	return s.client.ListResourceInstances(listResourceInstancesOptions)
}

// GetResourceInstance will get the resource instance.
func (s *Service) GetResourceInstance(getResourceInstanceOptions *resourcecontrollerv2.GetResourceInstanceOptions) (result *resourcecontrollerv2.ResourceInstance, response *core.DetailedResponse, err error) {
	return s.client.GetResourceInstance(getResourceInstanceOptions)
}

// SetServiceURL sets the service URL.
func (s *Service) SetServiceURL(url string) error {
	return s.client.SetServiceURL(url)
}

// GetServiceURL will get the service URL.
func (s *Service) GetServiceURL() string {
	return s.client.GetServiceURL()
}

// NewService returns a new service for the IBM Cloud Resource Controller api client.
func NewService(options ServiceOptions) (*Service, error) {
	if options.ResourceControllerV2Options == nil {
		options.ResourceControllerV2Options = &resourcecontrollerv2.ResourceControllerV2Options{}
	}
	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, err
	}
	options.Authenticator = auth
	service, err := resourcecontrollerv2.NewResourceControllerV2(options.ResourceControllerV2Options)
	if err != nil {
		return nil, err
	}
	return &Service{
		client: service,
	}, nil
}
