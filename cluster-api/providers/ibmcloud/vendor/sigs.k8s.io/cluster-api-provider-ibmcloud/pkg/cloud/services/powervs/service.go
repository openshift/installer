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
	"github.com/IBM-Cloud/power-go-client/power/client/datacenters"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_images"
	"github.com/IBM-Cloud/power-go-client/power/models"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
)

var _ PowerVS = &Service{}

// Service holds the PowerVS Service specific information.
type Service struct {
	session        *ibmpisession.IBMPISession
	instanceClient *instance.IBMPIInstanceClient
	networkClient  *instance.IBMPINetworkClient
	imageClient    *instance.IBMPIImageClient
	jobClient      *instance.IBMPIJobClient
	dhcpClient     *instance.IBMPIDhcpClient
}

// ServiceOptions holds the PowerVS Service Options specific information.
type ServiceOptions struct {
	*ibmpisession.IBMPIOptions

	CloudInstanceID string
}

// NewService returns a new service for the Power VS api client.
func NewService(options ServiceOptions) (PowerVS, error) {
	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, err
	}
	options.Authenticator = auth
	account, err := utils.GetAccount(auth)
	if err != nil {
		return nil, err
	}
	options.IBMPIOptions.UserAccount = account
	session, err := ibmpisession.NewIBMPISession(options.IBMPIOptions)
	if err != nil {
		return nil, err
	}

	return &Service{
		session: session,
	}, nil
}

// WithClients attach the clients to service.
func (s *Service) WithClients(options ServiceOptions) *Service {
	ctx := context.Background()
	s.instanceClient = instance.NewIBMPIInstanceClient(ctx, s.session, options.CloudInstanceID)
	s.networkClient = instance.NewIBMPINetworkClient(ctx, s.session, options.CloudInstanceID)
	s.imageClient = instance.NewIBMPIImageClient(ctx, s.session, options.CloudInstanceID)
	s.jobClient = instance.NewIBMPIJobClient(ctx, s.session, options.CloudInstanceID)
	s.dhcpClient = instance.NewIBMPIDhcpClient(ctx, s.session, options.CloudInstanceID)
	return s
}

// CreateInstance creates the virtual machine in the Power VS service instance.
func (s *Service) CreateInstance(body *models.PVMInstanceCreate) (*models.PVMInstanceList, error) {
	return s.instanceClient.Create(body)
}

// DeleteInstance deletes the virtual machine in the Power VS service instance.
func (s *Service) DeleteInstance(id string) error {
	return s.instanceClient.Delete(id)
}

// GetAllInstance returns all the virtual machine in the Power VS service instance.
func (s *Service) GetAllInstance() (*models.PVMInstances, error) {
	return s.instanceClient.GetAll()
}

// GetInstance returns the virtual machine in the Power VS service instance.
func (s *Service) GetInstance(id string) (*models.PVMInstance, error) {
	return s.instanceClient.Get(id)
}

// GetImage returns the image in the Power VS service instance.
func (s *Service) GetImage(id string) (*models.Image, error) {
	return s.imageClient.Get(id)
}

// GetAllImage returns all the images in the Power VS service instance.
func (s *Service) GetAllImage() (*models.Images, error) {
	return s.imageClient.GetAll()
}

// DeleteImage deletes the image in the Power VS service instance.
func (s *Service) DeleteImage(id string) error {
	return s.imageClient.Delete(id)
}

// CreateCosImage creates a import job to import the image in the Power VS service instance.
func (s *Service) CreateCosImage(body *models.CreateCosImageImportJob) (*models.JobReference, error) {
	return s.imageClient.CreateCosImage(body)
}

// GetCosImages returns the last import job in the Power VS service instance.
func (s *Service) GetCosImages(id string) (*models.Job, error) {
	params := p_cloud_images.NewPcloudV1CloudinstancesCosimagesGetParams().WithCloudInstanceID(id)
	resp, err := s.session.Power.PCloudImages.PcloudV1CloudinstancesCosimagesGet(params, s.session.AuthInfo(id))
	if err != nil || resp.Payload == nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GetJob returns the import job to in the Power VS service instance.
func (s *Service) GetJob(id string) (*models.Job, error) {
	return s.jobClient.Get(id)
}

// DeleteJob deletes the image import job in the Power VS service instance.
func (s *Service) DeleteJob(id string) error {
	return s.jobClient.Delete(id)
}

// GetAllNetwork returns all the networks in the Power VS service instance.
func (s *Service) GetAllNetwork() (*models.Networks, error) {
	return s.networkClient.GetAll()
}

// GetNetworkByID returns network corresponding to given id.
func (s *Service) GetNetworkByID(id string) (*models.Network, error) {
	return s.networkClient.Get(id)
}

// GetAllDHCPServers returns all the DHCP servers in the Power VS service instance.
func (s *Service) GetAllDHCPServers() (models.DHCPServers, error) {
	return s.dhcpClient.GetAll()
}

// GetDHCPServer returns the details for DHCP server associated with id.
func (s *Service) GetDHCPServer(id string) (*models.DHCPServerDetail, error) {
	return s.dhcpClient.Get(id)
}

// CreateDHCPServer creates a new DHCP server.
func (s *Service) CreateDHCPServer(options *models.DHCPServerCreate) (*models.DHCPServer, error) {
	return s.dhcpClient.Create(options)
}

// DeleteDHCPServer deletes the DHCP server.
func (s *Service) DeleteDHCPServer(id string) error {
	return s.dhcpClient.Delete(id)
}

// GetNetworkByName fetches the network with name. If not found, returns nil.
func (s *Service) GetNetworkByName(networkName string) (*models.NetworkReference, error) {
	var network *models.NetworkReference
	networks, err := s.GetAllNetwork()
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

// GetDatacenterCapabilities fetches the datacenter capabilities for the given zone.
func (s *Service) GetDatacenterCapabilities(zone string) (map[string]bool, error) {
	// though the function name is WithDatacenterRegion it takes zone as parameter
	params := datacenters.NewV1DatacentersGetParamsWithContext(context.TODO()).WithDatacenterRegion(zone)
	datacenter, err := s.session.Power.Datacenters.V1DatacentersGet(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get datacenter details for zone: %s err:%w", zone, err)
	}
	if datacenter == nil || datacenter.Payload == nil || datacenter.Payload.Capabilities == nil {
		return nil, fmt.Errorf("failed to get datacenter capabilities for zone: %s", zone)
	}
	return datacenter.Payload.Capabilities, nil
}
