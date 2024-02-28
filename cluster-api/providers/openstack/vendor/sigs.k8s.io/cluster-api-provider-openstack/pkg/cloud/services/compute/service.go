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

package compute

import (
	"fmt"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
)

type Service struct {
	scope              *scope.WithLogger
	_computeClient     clients.ComputeClient
	_volumeClient      clients.VolumeClient
	_imageClient       clients.ImageClient
	_networkingService *networking.Service
}

// NewService returns an instance of the compute service.
func NewService(scope *scope.WithLogger) (*Service, error) {
	return &Service{
		scope: scope,
	}, nil
}

func (s Service) getComputeClient() clients.ComputeClient {
	if s._computeClient == nil {
		computeClient, err := s.scope.NewComputeClient()
		if err != nil {
			return clients.NewComputeErrorClient(err)
		}

		s._computeClient = computeClient
	}

	return s._computeClient
}

func (s Service) getVolumeClient() clients.VolumeClient {
	if s._volumeClient == nil {
		volumeClient, err := s.scope.NewVolumeClient()
		if err != nil {
			return clients.NewVolumeErrorClient(err)
		}

		s._volumeClient = volumeClient
	}

	return s._volumeClient
}

func (s Service) getImageClient() clients.ImageClient {
	if s._imageClient == nil {
		imageClient, err := s.scope.NewImageClient()
		if err != nil {
			return clients.NewImageErrorClient(err)
		}

		s._imageClient = imageClient
	}

	return s._imageClient
}

func (s Service) getNetworkingService() (*networking.Service, error) {
	if s._networkingService == nil {
		networkingService, err := networking.NewService(s.scope)
		if err != nil {
			return nil, fmt.Errorf("failed to create networking service: %v", err)
		}

		s._networkingService = networkingService
	}

	return s._networkingService, nil
}
