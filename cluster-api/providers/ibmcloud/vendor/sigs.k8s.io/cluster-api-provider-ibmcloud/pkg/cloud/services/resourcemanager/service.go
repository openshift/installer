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

package resourcemanager

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
)

// Service holds the IBM Cloud Resource Manager Service specific information.
type Service struct {
	client *resourcemanagerv2.ResourceManagerV2
}

// NewService returns a new service for the resource manager.
func NewService(options *resourcemanagerv2.ResourceManagerV2Options) (ResourceManager, error) {
	if options == nil {
		options = &resourcemanagerv2.ResourceManagerV2Options{}
	}
	if options.Authenticator == nil {
		auth, err := authenticator.GetAuthenticator()
		if err != nil {
			return nil, err
		}
		options.Authenticator = auth
	}
	rmClient, err := resourcemanagerv2.NewResourceManagerV2(options)
	if err != nil {
		return nil, err
	}
	return &Service{
		client: rmClient,
	}, nil
}

// ListResourceGroups lists the resource groups.
func (s *Service) ListResourceGroups(listResourceGroupsOptions *resourcemanagerv2.ListResourceGroupsOptions) (result *resourcemanagerv2.ResourceGroupList, response *core.DetailedResponse, err error) {
	return s.client.ListResourceGroups(listResourceGroupsOptions)
}
