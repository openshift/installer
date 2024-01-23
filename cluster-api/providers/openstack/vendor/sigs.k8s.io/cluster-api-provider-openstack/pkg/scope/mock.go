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

package scope

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients/mock"
)

// MockScopeFactory implements both the ScopeFactory and ClientScope interfaces. It can be used in place of the default ProviderScopeFactory
// when we want to use mocked service clients which do not attempt to connect to a running OpenStack cloud.
type MockScopeFactory struct {
	ComputeClient *mock.MockComputeClient
	NetworkClient *mock.MockNetworkClient
	VolumeClient  *mock.MockVolumeClient
	ImageClient   *mock.MockImageClient
	LbClient      *mock.MockLbClient

	logger                 logr.Logger
	projectID              string
	clientScopeCreateError error
}

func NewMockScopeFactory(mockCtrl *gomock.Controller, projectID string, logger logr.Logger) *MockScopeFactory {
	computeClient := mock.NewMockComputeClient(mockCtrl)
	volumeClient := mock.NewMockVolumeClient(mockCtrl)
	imageClient := mock.NewMockImageClient(mockCtrl)
	networkClient := mock.NewMockNetworkClient(mockCtrl)
	lbClient := mock.NewMockLbClient(mockCtrl)

	return &MockScopeFactory{
		ComputeClient: computeClient,
		VolumeClient:  volumeClient,
		ImageClient:   imageClient,
		NetworkClient: networkClient,
		LbClient:      lbClient,
		projectID:     projectID,
		logger:        logger,
	}
}

func (f *MockScopeFactory) SetClientScopeCreateError(err error) {
	f.clientScopeCreateError = err
}

func (f *MockScopeFactory) NewClientScopeFromMachine(_ context.Context, _ client.Client, _ *infrav1.OpenStackMachine, _ []byte, _ logr.Logger) (Scope, error) {
	if f.clientScopeCreateError != nil {
		return nil, f.clientScopeCreateError
	}
	return f, nil
}

func (f *MockScopeFactory) NewClientScopeFromCluster(_ context.Context, _ client.Client, _ *infrav1.OpenStackCluster, _ []byte, _ logr.Logger) (Scope, error) {
	if f.clientScopeCreateError != nil {
		return nil, f.clientScopeCreateError
	}
	return f, nil
}

func (f *MockScopeFactory) NewComputeClient() (clients.ComputeClient, error) {
	return f.ComputeClient, nil
}

func (f *MockScopeFactory) NewVolumeClient() (clients.VolumeClient, error) {
	return f.VolumeClient, nil
}

func (f *MockScopeFactory) NewImageClient() (clients.ImageClient, error) {
	return f.ImageClient, nil
}

func (f *MockScopeFactory) NewNetworkClient() (clients.NetworkClient, error) {
	return f.NetworkClient, nil
}

func (f *MockScopeFactory) NewLbClient() (clients.LbClient, error) {
	return f.LbClient, nil
}

func (f *MockScopeFactory) Logger() logr.Logger {
	return f.logger
}

func (f *MockScopeFactory) ProjectID() string {
	return f.projectID
}
