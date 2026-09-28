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

	"github.com/IBM-Cloud/power-go-client/power/models"
)

//go:generate ../../../../hack/tools/bin/mockgen -source=./powervs.go -destination=./mock/powervs_generated.go -package=mock
//go:generate /usr/bin/env bash -c "cat ../../../../hack/scripts/verify/boilerplate/boilerplate.generatego.txt ./mock/powervs_generated.go > ./mock/_powervs_generated.go && mv ./mock/_powervs_generated.go ./mock/powervs_generated.go"

// PowerVS interface defines methods that a Cluster API IBM Cloud object should implement.
type PowerVS interface {
	// Instances
	CreateInstance(ctx context.Context, body *models.PVMInstanceCreate) (*models.PVMInstanceList, error)
	DeleteInstance(ctx context.Context, id string) error
	GetInstance(ctx context.Context, id string) (*models.PVMInstance, error)
	ListInstances(ctx context.Context) (*models.PVMInstances, error)

	// Images
	GetImage(ctx context.Context, id string) (*models.Image, error)
	DeleteImage(ctx context.Context, id string) error
	ListImages(ctx context.Context) (*models.Images, error)
	GetJob(ctx context.Context, id string) (*models.Job, error)
	DeleteJob(ctx context.Context, id string) error

	// COS Image Jobs
	CreateCosImage(ctx context.Context, body *models.CreateCosImageImportJob) (*models.JobReference, error)
	GetCosImages(ctx context.Context, id string) (*models.Job, error)

	// Networks
	ListNetworks(ctx context.Context) (*models.Networks, error)
	GetNetworkByID(ctx context.Context, id string) (*models.Network, error)
	GetNetworkByName(ctx context.Context, networkName string) (*models.NetworkReference, error)

	// DHCP Servers
	CreateDHCPServer(ctx context.Context, body *models.DHCPServerCreate) (*models.DHCPServer, error)
	GetDHCPServer(ctx context.Context, id string) (*models.DHCPServerDetail, error)
	DeleteDHCPServer(ctx context.Context, id string) error
	ListDHCPServers(ctx context.Context) (models.DHCPServers, error)

	// Datacenter
	GetDatacenterDetails(ctx context.Context, zone string) (*models.Datacenter, error)
}
