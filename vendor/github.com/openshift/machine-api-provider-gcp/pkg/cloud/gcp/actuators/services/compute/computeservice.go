package computeservice

import (
	"context"
	"log"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/openshift/machine-api-provider-gcp/pkg/version"
	"google.golang.org/api/compute/v1"
)

// GCPComputeService is a pass through wrapper for google.golang.org/api/compute/v1/compute
// to enable tests to mock this struct and control behavior.
type GCPComputeService interface {
	InstancesDelete(requestId string, project string, zone string, instance string) (*compute.Operation, error)
	InstancesInsert(project string, zone string, instance *compute.Instance) (*compute.Operation, error)
	InstancesGet(project string, zone string, instance string) (*compute.Instance, error)
	ZonesGet(project string, zone string) (*compute.Zone, error)
	ZoneOperationsGet(project string, zone string, operation string) (*compute.Operation, error)
	BasePath() string
	TargetPoolsGet(project string, region string, name string) (*compute.TargetPool, error)
	TargetPoolsAddInstance(project string, region string, name string, instance string) (*compute.Operation, error)
	TargetPoolsRemoveInstance(project string, region string, name string, instance string) (*compute.Operation, error)
	MachineTypesGet(project string, machineType string, zone string) (*compute.MachineType, error)
	RegionGet(project string, region string) (*compute.Region, error)
	GPUCompatibleMachineTypesList(project string, zone string, ctx context.Context) (map[string]GpuInfo, []string)
	AcceleratorTypeGet(project string, zone string, acceleratorType string) (*compute.AcceleratorType, error)
	ImageGet(project string, image string) (*compute.Image, error)
	ImageFamilyGet(project string, zone string, family string) (*compute.ImageFamilyView, error)
	InstanceGroupsListInstances(project string, zone string, instanceGroup string, request *compute.InstanceGroupsListInstancesRequest) (*compute.InstanceGroupsListInstances, error)
	InstanceGroupsAddInstances(project string, zone string, instance string, instanceGroup string) (*compute.Operation, error)
	InstanceGroupsRemoveInstances(project string, zone string, instance string, instanceGroup string) (*compute.Operation, error)
	InstanceGroupInsert(project string, zone string, instanceGroup *compute.InstanceGroup) (*compute.Operation, error)
	InstanceGroupGet(project string, zone string, instanceGroupName string) (*compute.InstanceGroup, error)
	AddInstanceGroupToBackendService(project string, region string, backendServiceName string, backendService *compute.BackendService) (*compute.Operation, error)
	BackendServiceGet(project string, region string, backendServiceName string) (*compute.BackendService, error)
}

type computeService struct {
	service *compute.Service
}

// BuilderFuncType is function type for building gcp client
type BuilderFuncType func(serviceAccountJSON string) (GCPComputeService, error)

// NewComputeService return a new computeService
func NewComputeService(serviceAccountJSON string) (GCPComputeService, error) {
	ctx := context.TODO()

	creds, err := google.CredentialsFromJSON(ctx, []byte(serviceAccountJSON), compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	service, err := compute.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	service.UserAgent = "gcpprovider.openshift.io/" + version.Version.String()

	return &computeService{
		service: service,
	}, nil
}

// InstancesInsert is a pass through wrapper for compute.Service.Instances.Insert(...)
func (c *computeService) InstancesInsert(project string, zone string, instance *compute.Instance) (*compute.Operation, error) {
	return c.service.Instances.Insert(project, zone, instance).Do()
}

// ZoneOperationsGet is a pass through wrapper for compute.Service.ZoneOperations.Get(...)
func (c *computeService) ZoneOperationsGet(project string, zone string, operation string) (*compute.Operation, error) {
	return c.service.ZoneOperations.Get(project, zone, operation).Do()
}

func (c *computeService) InstancesGet(project string, zone string, instance string) (*compute.Instance, error) {
	return c.service.Instances.Get(project, zone, instance).Do()
}

func (c *computeService) InstancesDelete(requestId string, project string, zone string, instance string) (*compute.Operation, error) {
	return c.service.Instances.Delete(project, zone, instance).RequestId(requestId).Do()
}

func (c *computeService) ZonesGet(project string, zone string) (*compute.Zone, error) {
	return c.service.Zones.Get(project, zone).Do()
}

func (c *computeService) BasePath() string {
	return c.service.BasePath
}

func (c *computeService) TargetPoolsGet(project string, region string, name string) (*compute.TargetPool, error) {
	return c.service.TargetPools.Get(project, region, name).Do()
}

func (c *computeService) TargetPoolsAddInstance(project string, region string, name string, instanceLink string) (*compute.Operation, error) {
	rb := &compute.TargetPoolsAddInstanceRequest{
		Instances: []*compute.InstanceReference{
			{
				Instance: instanceLink,
			},
		},
	}
	return c.service.TargetPools.AddInstance(project, region, name, rb).Do()
}

func (c *computeService) TargetPoolsRemoveInstance(project string, region string, name string, instanceLink string) (*compute.Operation, error) {
	rb := &compute.TargetPoolsRemoveInstanceRequest{
		Instances: []*compute.InstanceReference{
			{
				Instance: instanceLink,
			},
		},
	}
	return c.service.TargetPools.RemoveInstance(project, region, name, rb).Do()
}

func (c *computeService) MachineTypesGet(project string, zone string, machineType string) (*compute.MachineType, error) {
	return c.service.MachineTypes.Get(project, zone, machineType).Do()
}

type GpuInfo struct {
	Count int64
	Type  string
}

// GPUCompatibleMachineTypesList function lists machineTypes available in the zone and return map of A2 and A3 family and slice of N1 family machineTypes
func (c *computeService) GPUCompatibleMachineTypesList(project string, zone string, ctx context.Context) (map[string]GpuInfo, []string) {
	req := c.service.MachineTypes.List(project, zone)
	var (
		a2or3MachineFamily = map[string]GpuInfo{}
		n1MachineFamily    []string
	)
	if err := req.Pages(ctx, func(page *compute.MachineTypeList) error {
		for _, machineType := range page.Items {
			if strings.HasPrefix(machineType.Name, "a2") || strings.HasPrefix(machineType.Name, "a3") {
				a2or3MachineFamily[machineType.Name] = GpuInfo{
					Count: machineType.Accelerators[0].GuestAcceleratorCount,
					Type:  machineType.Accelerators[0].GuestAcceleratorType,
				}
			} else if strings.HasPrefix(machineType.Name, "n1") {
				n1MachineFamily = append(n1MachineFamily, machineType.Name)
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return a2or3MachineFamily, n1MachineFamily
}

func (c *computeService) AcceleratorTypeGet(project string, zone string, acceleratorType string) (*compute.AcceleratorType, error) {
	return c.service.AcceleratorTypes.Get(project, zone, acceleratorType).Do()
}

func (c *computeService) RegionGet(project string, region string) (*compute.Region, error) {
	return c.service.Regions.Get(project, region).Do()
}

func (c *computeService) InstanceGroupsAddInstances(project string, zone string, instance string, instanceGroup string) (*compute.Operation, error) {
	request := &compute.InstanceGroupsAddInstancesRequest{
		Instances: []*compute.InstanceReference{
			{
				Instance: instance,
			},
		},
	}
	return c.service.InstanceGroups.AddInstances(project, zone, instanceGroup, request).Do()
}

func (c *computeService) InstanceGroupsRemoveInstances(project string, zone string, instance string, instanceGroup string) (*compute.Operation, error) {
	request := &compute.InstanceGroupsRemoveInstancesRequest{
		Instances: []*compute.InstanceReference{
			{
				Instance: instance,
			},
		},
	}
	return c.service.InstanceGroups.RemoveInstances(project, zone, instanceGroup, request).Do()
}

func (c *computeService) InstanceGroupsListInstances(project string, zone string, instanceGroup string, request *compute.InstanceGroupsListInstancesRequest) (*compute.InstanceGroupsListInstances, error) {
	return c.service.InstanceGroups.ListInstances(project, zone, instanceGroup, request).Do()
}

func (c *computeService) InstanceGroupInsert(project string, zone string, instanceGroup *compute.InstanceGroup) (*compute.Operation, error) {
	return c.service.InstanceGroups.Insert(project, zone, instanceGroup).Do()
}

func (c *computeService) InstanceGroupGet(project string, zone string, instanceGroupName string) (*compute.InstanceGroup, error) {
	return c.service.InstanceGroups.Get(project, zone, instanceGroupName).Do()
}

func (c *computeService) AddInstanceGroupToBackendService(project string, region string, backendServiceName string, backendService *compute.BackendService) (*compute.Operation, error) {
	return c.service.RegionBackendServices.Update(project, region, backendServiceName, backendService).Do()
}

func (c *computeService) BackendServiceGet(project string, region string, backendServiceName string) (*compute.BackendService, error) {
	return c.service.RegionBackendServices.Get(project, region, backendServiceName).Do()
}

func (c *computeService) ImageGet(project string, image string) (*compute.Image, error) {
	return c.service.Images.Get(project, image).Do()
}

func (c *computeService) ImageFamilyGet(project string, zone string, family string) (*compute.ImageFamilyView, error) {
	return c.service.ImageFamilyViews.Get(project, zone, family).Do()
}
