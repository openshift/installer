/*
Copyright 2022 The ORC Authors.

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
	"time"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"go.uber.org/mock/gomock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/osclients/mock"
)

// MockScopeFactory implements both the ScopeFactory and ClientScope interfaces. It can be used in place of the default ProviderScopeFactory
// when we want to use mocked service clients which do not attempt to connect to a running OpenStack cloud.
type MockScopeFactory struct {
	AddressScope                *mock.MockAddressScopeClient
	ApplicationCredentialClient *mock.MockApplicationCredentialClient
	ComputeClient               *mock.MockComputeClient
	DomainClient                *mock.MockDomainClient
	EndpointClient              *mock.MockEndpointClient
	GroupClient                 *mock.MockGroupClient
	IdentityClient              *mock.MockIdentityClient
	ImageClient                 *mock.MockImageClient
	KeyPairClient               *mock.MockKeyPairClient
	NetworkClient               *mock.MockNetworkClient
	RoleClient                  *mock.MockRoleClient
	RoleAssignmentClient        *mock.MockRoleAssignmentClient
	ServiceClient               *mock.MockServiceClient
	UserClient                  *mock.MockUserClient
	VolumeClient                *mock.MockVolumeClient
	VolumeTypeClient            *mock.MockVolumeTypeClient
	ShareNetworkClient          *mock.MockShareNetworkClient

	clientScopeCreateError error
}

func NewMockScopeFactory(mockCtrl *gomock.Controller) *MockScopeFactory {
	addressScope := mock.NewMockAddressScopeClient(mockCtrl)
	applicationcredentialClient := mock.NewMockApplicationCredentialClient(mockCtrl)
	computeClient := mock.NewMockComputeClient(mockCtrl)
	domainClient := mock.NewMockDomainClient(mockCtrl)
	endpointClient := mock.NewMockEndpointClient(mockCtrl)
	groupClient := mock.NewMockGroupClient(mockCtrl)
	identityClient := mock.NewMockIdentityClient(mockCtrl)
	imageClient := mock.NewMockImageClient(mockCtrl)
	keypairClient := mock.NewMockKeyPairClient(mockCtrl)
	networkClient := mock.NewMockNetworkClient(mockCtrl)
	roleClient := mock.NewMockRoleClient(mockCtrl)
	roleassignmentClient := mock.NewMockRoleAssignmentClient(mockCtrl)
	serviceClient := mock.NewMockServiceClient(mockCtrl)
	userClient := mock.NewMockUserClient(mockCtrl)
	sharenetworkClient := mock.NewMockShareNetworkClient(mockCtrl)
	volumeClient := mock.NewMockVolumeClient(mockCtrl)
	volumetypeClient := mock.NewMockVolumeTypeClient(mockCtrl)

	return &MockScopeFactory{
		AddressScope:                addressScope,
		ApplicationCredentialClient: applicationcredentialClient,
		ComputeClient:               computeClient,
		DomainClient:                domainClient,
		EndpointClient:              endpointClient,
		GroupClient:                 groupClient,
		IdentityClient:              identityClient,
		ImageClient:                 imageClient,
		KeyPairClient:               keypairClient,
		NetworkClient:               networkClient,
		RoleClient:                  roleClient,
		RoleAssignmentClient:        roleassignmentClient,
		ServiceClient:               serviceClient,
		ShareNetworkClient:          sharenetworkClient,
		UserClient:                  userClient,
		VolumeClient:                volumeClient,
		VolumeTypeClient:            volumetypeClient,
	}
}

func (f *MockScopeFactory) SetClientScopeCreateError(err error) {
	f.clientScopeCreateError = err
}

func (f *MockScopeFactory) NewClientScopeFromObject(_ context.Context, _ client.Client, _ logr.Logger, _ ...orcv1alpha1.CloudCredentialsRefProvider) (Scope, error) {
	if f.clientScopeCreateError != nil {
		return nil, f.clientScopeCreateError
	}
	return f, nil
}

func (f *MockScopeFactory) NewAddressScopeClient() (osclients.AddressScopeClient, error) {
	return f.AddressScope, nil
}

func (f *MockScopeFactory) NewComputeClient() (osclients.ComputeClient, error) {
	return f.ComputeClient, nil
}

func (f *MockScopeFactory) NewImageClient() (osclients.ImageClient, error) {
	return f.ImageClient, nil
}

func (f *MockScopeFactory) NewNetworkClient() (osclients.NetworkClient, error) {
	return f.NetworkClient, nil
}

func (f *MockScopeFactory) NewIdentityClient() (osclients.IdentityClient, error) {
	return f.IdentityClient, nil
}

func (f *MockScopeFactory) NewUserClient() (osclients.UserClient, error) {
	return f.UserClient, nil
}

func (f *MockScopeFactory) NewVolumeClient() (osclients.VolumeClient, error) {
	return f.VolumeClient, nil
}

func (f *MockScopeFactory) NewVolumeTypeClient() (osclients.VolumeTypeClient, error) {
	return f.VolumeTypeClient, nil
}

func (f *MockScopeFactory) NewDomainClient() (osclients.DomainClient, error) {
	return f.DomainClient, nil
}

func (f *MockScopeFactory) NewServiceClient() (osclients.ServiceClient, error) {
	return f.ServiceClient, nil
}

func (f *MockScopeFactory) NewShareNetworkClient() (osclients.ShareNetworkClient, error) {
	return f.ShareNetworkClient, nil
}

func (f *MockScopeFactory) NewKeyPairClient() (osclients.KeyPairClient, error) {
	return f.KeyPairClient, nil
}

func (f *MockScopeFactory) NewGroupClient() (osclients.GroupClient, error) {
	return f.GroupClient, nil
}

func (f *MockScopeFactory) NewRoleClient() (osclients.RoleClient, error) {
	return f.RoleClient, nil
}

func (f *MockScopeFactory) NewRoleAssignmentClient() (osclients.RoleAssignmentClient, error) {
	return f.RoleAssignmentClient, nil
}

func (f *MockScopeFactory) NewEndpointClient() (osclients.EndpointClient, error) {
	return f.EndpointClient, nil
}

func (f *MockScopeFactory) NewApplicationCredentialClient() (osclients.ApplicationCredentialClient, error) {
	return f.ApplicationCredentialClient, nil
}

func (f *MockScopeFactory) ExtractToken() (*tokens.Token, error) {
	return &tokens.Token{ExpiresAt: time.Now().Add(24 * time.Hour)}, nil
}
