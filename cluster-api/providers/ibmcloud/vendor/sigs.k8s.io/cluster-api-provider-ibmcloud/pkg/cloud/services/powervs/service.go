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

package powervs

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_images"
	"github.com/IBM-Cloud/power-go-client/power/models"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/accounts"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
)

var _ PowerVS = &Service{}

// Service holds the PowerVS Service specific information.
type Service struct {
	session          *ibmpisession.IBMPISession
	instanceClient   *instance.IBMPIInstanceClient
	networkClient    *instance.IBMPINetworkClient
	imageClient      *instance.IBMPIImageClient
	jobClient        *instance.IBMPIJobClient
	dhcpClient       *instance.IBMPIDhcpClient
	dataCenterClient *instance.IBMPIDatacentersClient
}

// ServiceOptions holds the PowerVS Service Options specific information.
type ServiceOptions struct {
	*ibmpisession.IBMPIOptions
	WorkspaceID string
}

// NewService returns a new, fully initialized service for the PowerVS API client.
func NewService(ctx context.Context, options ServiceOptions) (PowerVS, error) {
	if options.Authenticator == nil {
		auth, err := authenticator.GetAuthenticator()
		if err != nil {
			return nil, err
		}
		options.Authenticator = auth
	}

	if options.UserAccount == "" {
		account, err := accounts.GetAccount(options.Authenticator)
		if err != nil {
			return nil, fmt.Errorf("failed to get user account: %w", err)
		}
		options.IBMPIOptions.UserAccount = account
	}

	session, err := ibmpisession.NewIBMPISession(options.IBMPIOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create PowerVS session: %w", err)
	}

	return &Service{
		session:          session,
		instanceClient:   instance.NewIBMPIInstanceClient(ctx, session, options.WorkspaceID),
		networkClient:    instance.NewIBMPINetworkClient(ctx, session, options.WorkspaceID),
		imageClient:      instance.NewIBMPIImageClient(ctx, session, options.WorkspaceID),
		jobClient:        instance.NewIBMPIJobClient(ctx, session, options.WorkspaceID),
		dhcpClient:       instance.NewIBMPIDhcpClient(ctx, session, options.WorkspaceID),
		dataCenterClient: instance.NewIBMPIDatacenterClient(ctx, session, options.WorkspaceID),
	}, nil
}

// CreateInstance creates the virtual machine in the Power VS service instance.
func (s *Service) CreateInstance(_ context.Context, body *models.PVMInstanceCreate) (*models.PVMInstanceList, error) {
	return s.instanceClient.Create(body)
}

// DeleteInstance deletes the virtual machine in the Power VS service instance.
func (s *Service) DeleteInstance(_ context.Context, id string) error {
	return s.instanceClient.Delete(id)
}

// GetInstance returns the virtual machine in the Power VS service instance.
func (s *Service) GetInstance(_ context.Context, id string) (*models.PVMInstance, error) {
	return s.instanceClient.Get(id)
}

// ListInstances returns all the virtual machine in the Power VS service instance.
func (s *Service) ListInstances(_ context.Context) (*models.PVMInstances, error) {
	return s.instanceClient.GetAll()
}

// GetImage returns the image in the Power VS service instance.
func (s *Service) GetImage(_ context.Context, id string) (*models.Image, error) {
	return s.imageClient.Get(id)
}

// DeleteImage deletes the image in the Power VS service instance.
func (s *Service) DeleteImage(_ context.Context, id string) error {
	return s.imageClient.Delete(id)
}

// ListImages returns all the images in the Power VS service instance.
func (s *Service) ListImages(_ context.Context) (*models.Images, error) {
	return s.imageClient.GetAll()
}

// GetJob returns the import job to in the Power VS service instance.
func (s *Service) GetJob(_ context.Context, id string) (*models.Job, error) {
	return s.jobClient.Get(id)
}

// DeleteJob deletes the image import job in the Power VS service instance.
func (s *Service) DeleteJob(_ context.Context, id string) error {
	return s.jobClient.Delete(id)
}

// CreateCosImage creates a import job to import the image in the Power VS service instance.
func (s *Service) CreateCosImage(_ context.Context, body *models.CreateCosImageImportJob) (*models.JobReference, error) {
	return s.imageClient.CreateCosImage(body)
}

// GetCosImages returns the last import job in the Power VS service instance.
func (s *Service) GetCosImages(_ context.Context, id string) (*models.Job, error) {
	params := p_cloud_images.NewPcloudV1CloudinstancesCosimagesGetParams().WithCloudInstanceID(id)
	resp, err := s.session.Power.PCloudImages.PcloudV1CloudinstancesCosimagesGet(params, s.session.AuthInfo(id))
	if err != nil {
		if _, ok := err.(*p_cloud_images.PcloudV1CloudinstancesCosimagesGetNotFound); ok {
			return nil, nil
		}
		return nil, err
	}
	if resp == nil || resp.Payload == nil {
		return nil, nil
	}
	return resp.Payload, nil
}

// ListNetworks returns all the networks in the Power VS service instance.
func (s *Service) ListNetworks(_ context.Context) (*models.Networks, error) {
	return s.networkClient.GetAll()
}

// GetNetworkByID returns network corresponding to given id.
func (s *Service) GetNetworkByID(_ context.Context, id string) (*models.Network, error) {
	return s.networkClient.Get(id)
}

// GetNetworkByName fetches the network with name. If not found, returns nil.
func (s *Service) GetNetworkByName(ctx context.Context, networkName string) (*models.NetworkReference, error) {
	var network *models.NetworkReference
	networks, err := s.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}
	for _, nw := range networks.Networks {
		if *nw.Name == networkName {
			network = nw
		}
	}

	return network, nil
}

// CreateDHCPServer creates a new DHCP server.
func (s *Service) CreateDHCPServer(_ context.Context, options *models.DHCPServerCreate) (*models.DHCPServer, error) {
	return s.dhcpClient.Create(options)
}

// GetDHCPServer returns the details for DHCP server associated with id.
func (s *Service) GetDHCPServer(_ context.Context, id string) (*models.DHCPServerDetail, error) {
	return s.dhcpClient.Get(id)
}

// DeleteDHCPServer deletes the DHCP server.
func (s *Service) DeleteDHCPServer(_ context.Context, id string) error {
	return s.dhcpClient.Delete(id)
}

// ListDHCPServers returns all the DHCP servers in the Power VS service instance.
func (s *Service) ListDHCPServers(_ context.Context) (models.DHCPServers, error) {
	return s.dhcpClient.GetAll()
}

// GetDatacenterDetails fetches the datacenter capabilities for the given zone.
func (s *Service) GetDatacenterDetails(_ context.Context, zone string) (*models.Datacenter, error) {
	return s.dataCenterClient.Get(zone)
}
