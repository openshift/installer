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
	"github.com/IBM-Cloud/power-go-client/power/models"
)

//go:generate ../../../../hack/tools/bin/mockgen -source=./powervs.go -destination=./mock/powervs_generated.go -package=mock
//go:generate /usr/bin/env bash -c "cat ../../../../hack/boilerplate/boilerplate.generatego.txt ./mock/powervs_generated.go > ./mock/_powervs_generated.go && mv ./mock/_powervs_generated.go ./mock/powervs_generated.go"

// PowerVS interface defines methods that a Cluster API IBMCLOUD object should implement.
type PowerVS interface {
	CreateInstance(body *models.PVMInstanceCreate) (*models.PVMInstanceList, error)
	DeleteInstance(id string) error
	GetAllInstance() (*models.PVMInstances, error)
	GetAllImage() (*models.Images, error)
	GetAllNetwork() (*models.Networks, error)
	GetNetworkByID(id string) (*models.Network, error)
	GetInstance(id string) (*models.PVMInstance, error)
	GetImage(id string) (*models.Image, error)
	DeleteImage(id string) error
	CreateCosImage(body *models.CreateCosImageImportJob) (*models.JobReference, error)
	GetCosImages(id string) (*models.Job, error)
	GetJob(id string) (*models.Job, error)
	DeleteJob(id string) error
	GetAllDHCPServers() (models.DHCPServers, error)
	GetDHCPServer(id string) (*models.DHCPServerDetail, error)
	CreateDHCPServer(*models.DHCPServerCreate) (*models.DHCPServer, error)
	DeleteDHCPServer(id string) error
	WithClients(options ServiceOptions) *Service
	GetNetworkByName(networkName string) (*models.NetworkReference, error)
	GetDatacenterCapabilities(zone string) (map[string]bool, error)
}
