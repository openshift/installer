/*
Copyright 2018 The Kubernetes Authors.

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

package networking

import (
	"fmt"
	"sort"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
)

const (
	networkPrefix string = "k8s-clusterapi"
	trunkResource string = "trunks"
	portResource  string = "ports"
)

// Service interfaces with the OpenStack Networking API.
// It will create a network related infrastructure for the cluster, like network, subnet, router, security groups.
type Service struct {
	scope  *scope.WithLogger
	client clients.NetworkClient
}

// NewService returns an instance of the networking service.
func NewService(scope *scope.WithLogger) (*Service, error) {
	networkClient, err := scope.NewNetworkClient()
	if err != nil {
		return nil, err
	}

	return &Service{
		scope:  scope,
		client: networkClient,
	}, nil
}

// replaceAllAttributesTags replaces all tags on a neworking resource.
// the value of resourceType must match one of the allowed constants: trunkResource or portResource.
func (s *Service) replaceAllAttributesTags(eventObject runtime.Object, resourceType string, resourceID string, tags []string) error {
	if len(tags) == 0 {
		s.scope.Logger().Info("No tags provided to replaceAllAttributesTags", "resource", resourceType, "ID", resourceID)
		return nil
	}
	if resourceType != trunkResource && resourceType != portResource {
		record.Warnf(eventObject, "FailedReplaceAllAttributesTags", "Invalid resourceType argument in function call")
		panic(fmt.Errorf("invalid argument: resourceType, %s, does not match allowed arguments: %s or %s", resourceType, trunkResource, portResource))
	}
	// remove duplicate values from tags
	tagsMap := make(map[string]string)
	for _, t := range tags {
		tagsMap[t] = t
	}

	uniqueTags := []string{}
	for k := range tagsMap {
		uniqueTags = append(uniqueTags, k)
	}

	// Sort the tags so that we always get fixed order of tags to make UT easier
	sort.Strings(uniqueTags)

	_, err := s.client.ReplaceAllAttributesTags(resourceType, resourceID, attributestags.ReplaceAllOpts{
		Tags: uniqueTags,
	})
	if err != nil {
		record.Warnf(eventObject, "FailedReplaceAllAttributesTags", "Failed to replace all attributestags, %s: %v", resourceID, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulReplaceAllAttributeTags", "Replaced all attributestags for %s with tags %s", resourceID, uniqueTags)
	return nil
}
