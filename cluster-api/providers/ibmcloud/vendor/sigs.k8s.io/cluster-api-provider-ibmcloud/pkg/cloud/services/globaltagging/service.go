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

package globaltagging

import (
	"fmt"
	"net/http"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
)

// Service holds the IBM Cloud Global Tagging Service specific information.
type Service struct {
	client *globaltaggingv1.GlobalTaggingV1
}

// ServiceOptions holds the IBM Cloud Global Tagging Service Options specific information.
type ServiceOptions struct {
	*globaltaggingv1.GlobalTaggingV1Options
}

// CreateTag creates a new Tag.
func (s *Service) CreateTag(options *globaltaggingv1.CreateTagOptions) (*globaltaggingv1.CreateTagResults, *core.DetailedResponse, error) {
	return s.client.CreateTag(options)
}

// AttachTag will add tag(s) to resource(s).
func (s *Service) AttachTag(options *globaltaggingv1.AttachTagOptions) (*globaltaggingv1.TagResults, *core.DetailedResponse, error) {
	return s.client.AttachTag(options)
}

// GetTagByName returns the Tag with the provided name, if found.
func (s *Service) GetTagByName(tagName string) (*globaltaggingv1.Tag, error) {
	accountID, err := utils.GetAccountID()
	if err != nil {
		return nil, err
	}

	listOptions := s.client.NewListTagsOptions()
	listOptions.SetTagType(globaltaggingv1.AttachTagOptionsTagTypeUserConst)
	listOptions.SetAccountID(accountID)

	result, response, err := s.client.ListTags(listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed listing user tags: %w", err)
	}
	if result == nil || (response != nil && response.StatusCode == http.StatusNotFound) {
		return nil, fmt.Errorf("failed to list tags")
	}
	for _, tag := range result.Items {
		if tag.Name != nil && *tag.Name == tagName {
			return ptr.To(tag), nil
		}
	}
	return nil, nil
}

// NewService returns a new service for the IBM Cloud Global Tagging api client.
func NewService(options ServiceOptions) (*Service, error) {
	if options.GlobalTaggingV1Options == nil {
		options.GlobalTaggingV1Options = &globaltaggingv1.GlobalTaggingV1Options{}
	}
	if options.Authenticator == nil {
		auth, err := authenticator.GetAuthenticator()
		if err != nil {
			return nil, err
		}
		options.Authenticator = auth
	}
	service, err := globaltaggingv1.NewGlobalTaggingV1(options.GlobalTaggingV1Options)
	if err != nil {
		return nil, err
	}
	return &Service{
		client: service,
	}, nil
}
