package computeservice

import (
	"context"
	"errors"
	"fmt"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	NoMachinesInPool               = "NoMachinesInPool"
	WithMachineInPool              = "WithMachineInPool"
	GroupDoesNotExist              = "groupDoesNotExist"
	EmptyInstanceList              = "emptyInstanceList"
	ErrUnregisteringInstance       = "errUnregisteringInstance"
	ErrRegisteringInstance         = "errRegisteringInstance"
	ErrRegisteringNewInstanceGroup = "errRegisteringNewInstanceGroup"
	ErrPatchingBackendService      = "errPatchingBackendService"
	ErrGettingBackendService       = "errGettingBackendService"
	ErrFailGroupGet                = "errFailGroupGet"
	ErrGroupNotFound               = "errGroupNotFound"
	ErrImageNotFound               = "errImageNotFound"
	PatchBackendService            = "patchBackendService"
	AddGroupSuccessfully           = "addGroupSuccessfully"
	UEFICompatible                 = "UEFI_COMPATIBLE"
)

type GCPComputeServiceMock struct {
	MockGPUCompatibleMachineTypesList func(project string, zone string, ctx context.Context) (map[string]GpuInfo, []string)
	MockInstancesInsert               func(project string, zone string, instance *compute.Instance) (*compute.Operation, error)
	MockMachineTypesGet               func(project string, zone string, machineType string) (*compute.MachineType, error)
	MockRegionGet                     func(project string, region string) (*compute.Region, error)
	mockZoneOperationsGet             func(project string, zone string, operation string) (*compute.Operation, error)
	mockInstancesGet                  func(project string, zone string, instance string) (*compute.Instance, error)
}

func (c *GCPComputeServiceMock) InstancesInsert(project string, zone string, instance *compute.Instance) (*compute.Operation, error) {
	if c.MockInstancesInsert == nil {
		return nil, nil
	}
	return c.MockInstancesInsert(project, zone, instance)
}

func (c *GCPComputeServiceMock) InstancesDelete(requestId string, project string, zone string, instance string) (*compute.Operation, error) {
	return &compute.Operation{
		Status: "DONE",
	}, nil
}

func (c *GCPComputeServiceMock) ZoneOperationsGet(project string, zone string, operation string) (*compute.Operation, error) {
	if c.mockZoneOperationsGet == nil {
		return nil, nil
	}
	return c.mockZoneOperationsGet(project, zone, operation)
}

func (c *GCPComputeServiceMock) InstancesGet(project string, zone string, instance string) (*compute.Instance, error) {
	if c.mockInstancesGet == nil {
		return &compute.Instance{
			Name:         instance,
			Zone:         zone,
			MachineType:  "n1-standard-1",
			CanIpForward: true,
			NetworkInterfaces: []*compute.NetworkInterface{
				{
					NetworkIP: "10.0.0.15",
					AccessConfigs: []*compute.AccessConfig{
						{
							NatIP: "35.243.147.143",
						},
					},
				},
			},
			Status: "RUNNING",
		}, nil
	}
	return c.mockInstancesGet(project, zone, instance)
}

func (c *GCPComputeServiceMock) ZonesGet(project string, zone string) (*compute.Zone, error) {
	return nil, nil
}

func (c *GCPComputeServiceMock) BasePath() string {
	return "path/"
}

func (c *GCPComputeServiceMock) TargetPoolsGet(project string, region string, name string) (*compute.TargetPool, error) {
	if region == NoMachinesInPool {
		return &compute.TargetPool{}, nil
	}
	if region == WithMachineInPool {
		return &compute.TargetPool{
			Instances: []string{
				"https://www.googleapis.com/compute/v1/projects/testProject/zones/zone1/instances/testInstance",
			},
		}, nil
	}
	return nil, nil
}

func (c *GCPComputeServiceMock) TargetPoolsAddInstance(project string, region string, name string, instance string) (*compute.Operation, error) {
	return nil, nil
}

func (c *GCPComputeServiceMock) TargetPoolsRemoveInstance(project string, region string, name string, instance string) (*compute.Operation, error) {
	return nil, nil
}

func (c *GCPComputeServiceMock) MachineTypesGet(project string, zone string, machineType string) (*compute.MachineType, error) {
	if c.MockMachineTypesGet == nil {
		return nil, nil
	}
	return c.MockMachineTypesGet(project, zone, machineType)
}

func NewComputeServiceMock() (*compute.Instance, *GCPComputeServiceMock) {
	var receivedInstance compute.Instance
	computeServiceMock := GCPComputeServiceMock{
		MockInstancesInsert: func(project string, zone string, instance *compute.Instance) (*compute.Operation, error) {
			receivedInstance = *instance
			return &compute.Operation{
				Status: "DONE",
			}, nil
		},
		mockZoneOperationsGet: func(project string, zone string, operation string) (*compute.Operation, error) {
			return &compute.Operation{
				Status: "DONE",
			}, nil
		},
	}
	return &receivedInstance, &computeServiceMock
}

func MockBuilderFuncType(serviceAccountJSON string) (GCPComputeService, error) {
	_, computeSvc := NewComputeServiceMock()
	return computeSvc, nil
}

func MockBuilderFuncTypeNotFound(serviceAccountJSON string) (GCPComputeService, error) {
	_, computeSvc := NewComputeServiceMock()
	computeSvc.mockInstancesGet = func(project string, zone string, instance string) (*compute.Instance, error) {
		return nil, &googleapi.Error{
			Code: 404,
		}
	}
	return computeSvc, nil
}

func (c *GCPComputeServiceMock) RegionGet(project string, region string) (*compute.Region, error) {
	if c.MockRegionGet == nil {
		return &compute.Region{Quotas: nil}, nil
	}

	return c.MockRegionGet(project, region)
}

func (c *GCPComputeServiceMock) GPUCompatibleMachineTypesList(project string, zone string, ctx context.Context) (map[string]GpuInfo, []string) {
	if c.MockGPUCompatibleMachineTypesList == nil {
		var compatibleMachineType = []string{"n1-test-machineType"}
		return nil, compatibleMachineType
	}

	return c.MockGPUCompatibleMachineTypesList(project, zone, ctx)
}

func (c *GCPComputeServiceMock) AcceleratorTypeGet(project string, zone string, acceleratorType string) (*compute.AcceleratorType, error) {
	return nil, nil
}

func (c *GCPComputeServiceMock) InstanceGroupsListInstances(projectID string, zone string, instanceGroup string, request *compute.InstanceGroupsListInstancesRequest) (*compute.InstanceGroupsListInstances, error) {
	if projectID == GroupDoesNotExist {
		return nil, &googleapi.Error{
			Code: 404,
		}
	}
	if projectID == EmptyInstanceList {
		return &compute.InstanceGroupsListInstances{}, nil
	}
	if projectID == ErrUnregisteringInstance {
		return &compute.InstanceGroupsListInstances{
			Items: []*compute.InstanceWithNamedPorts{
				{
					Instance: "https://www.googleapis.com/compute/v1/projects/errUnregisteringInstance/zones/zone1/instances/testInstance",
				},
			},
		}, nil
	}
	instances := &compute.InstanceGroupsListInstances{
		Items: []*compute.InstanceWithNamedPorts{
			{
				Instance: "https://www.googleapis.com/compute/v1/projects/testProject/zones/zone1/instances/testInstance",
			},
		},
	}
	return instances, nil
}

func (c *GCPComputeServiceMock) InstanceGroupsAddInstances(project string, zone string, instance string, instanceGroup string) (*compute.Operation, error) {
	if project == ErrRegisteringInstance {
		return nil, errors.New("a GCP error")
	}
	return &compute.Operation{
		Status: "DONE",
	}, nil
}

func (c *GCPComputeServiceMock) InstanceGroupsRemoveInstances(project string, zone string, instance string, instanceGroup string) (*compute.Operation, error) {
	if project == ErrUnregisteringInstance {
		return nil, errors.New("a GCP error")
	}
	return &compute.Operation{
		Status: "DONE",
	}, nil
}

func (c *GCPComputeServiceMock) InstanceGroupInsert(project string, zone string, instanceGroup *compute.InstanceGroup) (*compute.Operation, error) {
	if project == AddGroupSuccessfully {
		return &compute.Operation{
			Status: "DONE",
		}, nil

	}
	if project == ErrRegisteringNewInstanceGroup {
		return nil, errors.New("failed to register new instanceGroup")
	}
	return nil, nil
}

func (c *GCPComputeServiceMock) InstanceGroupGet(project string, zone string, instanceGroupName string) (*compute.InstanceGroup, error) {
	if project == ErrFailGroupGet {
		return nil, errors.New("instanceGroupGet request failed")
	}
	if project == ErrGroupNotFound {
		return nil, errors.New("instanceGroupGet request failed")
	}
	return nil, nil
}

func (c *GCPComputeServiceMock) AddInstanceGroupToBackendService(project string, region string, backendServiceName string, backendService *compute.BackendService) (*compute.Operation, error) {
	if project == ErrPatchingBackendService {
		return nil, errors.New("failed to add new instanceGroup to backend service")
	}
	return &compute.Operation{
		Status: "DONE",
	}, nil
}

func (c *GCPComputeServiceMock) BackendServiceGet(project string, region string, backendServiceName string) (*compute.BackendService, error) {
	if project == ErrGettingBackendService || project == ErrPatchingBackendService {
		return nil, errors.New("failed to get the regional backend service")
	}
	return &compute.BackendService{
		Backends: []*compute.Backend{
			{
				Group: fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/zone1/instanceGroups/CLUSTERID-master-zone1", project),
			},
		},
	}, nil
}

func (c *GCPComputeServiceMock) ImageGet(project string, image string) (*compute.Image, error) {
	if project == ErrImageNotFound {
		return nil, errors.New("imageGet request failed")
	}

	img := &compute.Image{
		GuestOsFeatures: make([]*compute.GuestOsFeature, 0),
	}

	if image == "uefi-image" {
		img.GuestOsFeatures = append(img.GuestOsFeatures, &compute.GuestOsFeature{Type: UEFICompatible})
	}

	return img, nil
}

func (c *GCPComputeServiceMock) ImageFamilyGet(project, zone, family string) (*compute.ImageFamilyView, error) {
	if project == ErrImageNotFound {
		return nil, errors.New("imageGet request failed")
	}

	imgView := &compute.ImageFamilyView{
		Image: &compute.Image{
			GuestOsFeatures: make([]*compute.GuestOsFeature, 0),
		},
	}

	if family == "uefi-image-family" {
		imgView.Image.GuestOsFeatures = append(imgView.Image.GuestOsFeatures, &compute.GuestOsFeature{Type: UEFICompatible})
	}

	return imgView, nil
}
