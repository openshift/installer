package nutanix

import (
	"context"
	"fmt"
	"time"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation/field"

	machinev1 "github.com/openshift/api/machine/v1"
)

// MachinePool stores the configuration for a machine pool installed
// on Nutanix.
type MachinePool struct {
	// NumCPUs is the total number of virtual processor cores to assign a vm.
	//
	// +optional
	NumCPUs int64 `json:"cpus,omitempty"`

	// NumCoresPerSocket is the number of cores per socket in a vm. The number
	// of vCPUs on the vm will be NumCPUs times NumCoresPerSocket.
	// For example: 4 CPUs and 4 Cores per socket will result in 16 VPUs.
	// The AHV scheduler treats socket and core allocation exactly the same
	// so there is no benefit to configuring cores over CPUs.
	//
	// +optional
	NumCoresPerSocket int64 `json:"coresPerSocket,omitempty"`

	// Memory is the size of a VM's memory in MiB.
	//
	// +optional
	MemoryMiB int64 `json:"memoryMiB,omitempty"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk,omitempty"`

	// BootType indicates the boot type (Legacy, UEFI or SecureBoot) the Machine's VM uses to boot.
	// If this field is empty or omitted, the VM will use the default boot type "Legacy" to boot.
	// "SecureBoot" depends on "UEFI" boot, i.e., enabling "SecureBoot" means that "UEFI" boot is also enabled.
	// +kubebuilder:validation:Enum="";Legacy;UEFI;SecureBoot
	// +optional
	BootType machinev1.NutanixBootType `json:"bootType,omitempty"`

	// Project optionally identifies a Prism project for the Machine's VM to associate with.
	// +optional
	Project *machinev1.NutanixResourceIdentifier `json:"project,omitempty"`

	// Categories optionally adds one or more prism categories (each with key and value) for
	// the Machine's VM to associate with. All the category key and value pairs specified must
	// already exist in the prism central.
	// +listType=map
	// +listMapKey=key
	// +optional
	Categories []machinev1.NutanixCategory `json:"categories,omitempty"`

	// GPUs is a list of GPU devices to attach to the machine's VM.
	// +listType=set
	// +optional
	GPUs []machinev1.NutanixGPU `json:"gpus"`

	// DataDisks holds information of the data disks to attach to the Machine's VM
	// +listType=set
	// +optional
	DataDisks []DataDisk `json:"dataDisks"`

	// FailureDomains optionally configures a list of failure domain names
	// that will be applied to the MachinePool
	// +listType=set
	// +optional
	FailureDomains []string `json:"failureDomains,omitempty"`
}

// OSDisk defines the system disk for a Machine VM.
type OSDisk struct {
	// DiskSizeGiB defines the size of disk in GiB.
	//
	// +optional
	DiskSizeGiB int64 `json:"diskSizeGiB,omitempty"`
}

// StorageResourceReference holds reference information of a storage resource (storage container, data source image, etc.)
type StorageResourceReference struct {
	// ReferenceName is the identifier of the storage resource configured in the FailureDomain.
	// +optional
	ReferenceName string `json:"referenceName,omitempty"`

	// UUID is the UUID of the storage container resource in the Prism Element.
	// +kubebuilder:validation:Required
	UUID string `json:"uuid"`

	// Name is the name of the storage container resource in the Prism Element.
	// +optional
	Name string `json:"name,omitempty"`
}

// StorageConfig specifies the storage configuration parameters for VM disks.
type StorageConfig struct {
	// diskMode specifies the disk mode.
	// The valid values are Standard and Flash, and the default is Standard.
	// +kubebuilder:default=Standard
	// +kubebuilder:validation:Enum=Standard;Flash
	DiskMode machinev1.NutanixDiskMode `json:"diskMode"`

	// storageContainer refers to the storage_container used by the VM disk.
	// +optional
	StorageContainer *StorageResourceReference `json:"storageContainer,omitempty"`
}

// DataDisk defines a data disk for a Machine VM.
type DataDisk struct {
	// diskSize is size (in Quantity format) of the disk to attach to the VM.
	// See https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#Format for the Quantity format and example documentation.
	// The minimum diskSize is 1GB.
	// +kubebuilder:validation:Required
	DiskSize resource.Quantity `json:"diskSize"`

	// deviceProperties are the properties of the disk device.
	// +optional
	DeviceProperties *machinev1.NutanixVMDiskDeviceProperties `json:"deviceProperties,omitempty"`

	// storageConfig are the storage configuration parameters of the VM disks.
	// +optional
	StorageConfig *StorageConfig `json:"storageConfig,omitempty"`

	// dataSource refers to a data source image for the VM disk.
	// +optional
	DataSourceImage *StorageResourceReference `json:"dataSourceImage,omitempty"`
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.NumCPUs != 0 {
		p.NumCPUs = required.NumCPUs
	}

	if required.NumCoresPerSocket != 0 {
		p.NumCoresPerSocket = required.NumCoresPerSocket
	}

	if required.MemoryMiB != 0 {
		p.MemoryMiB = required.MemoryMiB
	}

	if required.OSDisk.DiskSizeGiB != 0 {
		p.OSDisk.DiskSizeGiB = required.OSDisk.DiskSizeGiB
	}

	if len(required.BootType) != 0 {
		p.BootType = required.BootType
	}

	if required.Project != nil {
		p.Project = required.Project
	}

	if len(required.Categories) > 0 {
		p.Categories = required.Categories
	}

	if len(required.FailureDomains) > 0 {
		p.FailureDomains = required.FailureDomains
	}

	if len(required.GPUs) > 0 {
		p.GPUs = required.GPUs
	}

	if len(required.DataDisks) > 0 {
		p.DataDisks = required.DataDisks
	}
}

// ValidateConfig validates the MachinePool configuration.
func (p *MachinePool) ValidateConfig(platform *Platform, role string) error {
	nc, err := CreateNutanixClientFromPlatform(platform)
	if err != nil {
		return fmt.Errorf("fail to create nutanix client. %w", err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	errList := field.ErrorList{}
	fldPath := field.NewPath("platform", "nutanix")
	var errMsg string

	// validate BootType
	if p.BootType != "" && p.BootType != machinev1.NutanixLegacyBoot &&
		p.BootType != machinev1.NutanixUEFIBoot && p.BootType != machinev1.NutanixSecureBoot {
		errMsg = fmt.Sprintf("valid bootType: \"\", %q, %q, %q.", machinev1.NutanixLegacyBoot, machinev1.NutanixUEFIBoot, machinev1.NutanixSecureBoot)
		errList = append(errList, field.Invalid(fldPath.Child("bootType"), p.BootType, errMsg))
	}

	// validate project if configured
	if p.Project != nil {
		fldErr := p.validateProjectConfig(ctx, nc, fldPath)
		if fldErr != nil {
			errList = append(errList, fldErr)
		}
	}

	// validate categories if configured
	if len(p.Categories) > 0 {
		for _, category := range p.Categories {
			if _, err = nc.V3.GetCategoryValue(ctx, category.Key, category.Value); err != nil {
				errMsg = fmt.Sprintf("Failed to find the category with key %q and value %q. error: %v", category.Key, category.Value, err)
				errList = append(errList, field.Invalid(fldPath.Child("categories"), category, errMsg))
			}
		}
	}

	// validate FailureDomains if configured
	for _, fdName := range p.FailureDomains {
		_, err := platform.GetFailureDomainByName(fdName)
		if err != nil {
			errList = append(errList, field.Invalid(fldPath.Child("failureDomains"), fdName, fmt.Sprintf("The failure domain is not defined: %v", err)))
		}
	}

	// validate GPUs if configured, currently only "worker" machines allow GPUs.
	if len(p.GPUs) > 0 {
		if role == "master" {
			errList = append(errList, field.Forbidden(fldPath.Child("gpus"), "'gpus' are not supported for 'master' nodes, you can only configure it for 'worker' nodes."))
		} else {
			fldErrs := p.validateGPUsConfig(ctx, nc, platform, fldPath)
			for _, fldErr := range fldErrs {
				errList = append(errList, fldErr)
			}
		}
	}

	// validate DataDisks if configured, currently only "worker" machines allow DataDisks.
	if len(p.DataDisks) > 0 {
		if role == "master" {
			errList = append(errList, field.Forbidden(fldPath.Child("gpus"), "'dataDisks' are not supported for 'master' nodes, you can only configure it for 'worker' nodes."))
		} else {
			fldErrs := p.validateDataDisksConfig(ctx, nc, platform, fldPath)
			for _, fldErr := range fldErrs {
				errList = append(errList, fldErr)
			}
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("%s", errList.ToAggregate().Error())
	}
	return nil
}

// validateProjectConfig validates the Project configuration in the machinePool.
func (p *MachinePool) validateProjectConfig(ctx context.Context, nc *nutanixclientv3.Client, fldPath *field.Path) *field.Error {
	if p.Project != nil {
		switch p.Project.Type {
		case machinev1.NutanixIdentifierName:
			if p.Project.Name == nil || *p.Project.Name == "" {
				return field.Required(fldPath.Child("project", "name"), "missing projct name")
			}

			projectName := *p.Project.Name
			filter := fmt.Sprintf("name==%s", projectName)
			res, err := nc.V3.ListProject(ctx, &nutanixclientv3.DSMetadata{
				Filter: &filter,
			})
			switch {
			case err != nil:
				return field.Invalid(fldPath.Child("project", "name"), projectName,
					fmt.Sprintf("failed to find project with name %q. error: %v", projectName, err))
			case len(res.Entities) == 0:
				return field.Invalid(fldPath.Child("project", "name"), projectName,
					fmt.Sprintf("unable to find project with name %q.", projectName))
			case len(res.Entities) > 1:
				return field.Invalid(fldPath.Child("project", "name"), projectName,
					fmt.Sprintf("found more than one (%v) projects with name %q.", len(res.Entities), projectName))
			default:
				p.Project.Type = machinev1.NutanixIdentifierUUID
				p.Project.UUID = res.Entities[0].Metadata.UUID
			}
		case machinev1.NutanixIdentifierUUID:
			if p.Project.UUID == nil || *p.Project.UUID == "" {
				return field.Required(fldPath.Child("project", "uuid"), "missing projct uuid")
			} else {
				if _, err := nc.V3.GetProject(ctx, *p.Project.UUID); err != nil {
					return field.Invalid(fldPath.Child("project", "uuid"), *p.Project.UUID,
						fmt.Sprintf("failed to get the project with uuid %s. error: %v", *p.Project.UUID, err))
				}
			}
		default:
			return field.Invalid(fldPath.Child("project", "type"), p.Project.Type,
				fmt.Sprintf("invalid project identifier type, valid types are: %q, %q.", machinev1.NutanixIdentifierName, machinev1.NutanixIdentifierUUID))
		}
	}

	return nil
}

// validateGPUsConfig validates the GPUs configuration in the machinePool.
func (p *MachinePool) validateGPUsConfig(ctx context.Context, nc *nutanixclientv3.Client, platform *Platform, fldPath *field.Path) (fldErrs []*field.Error) {
	if len(p.GPUs) == 0 {
		return fldErrs
	}

	peUUIDs := []string{}
	for _, fdName := range p.FailureDomains {
		if fd, err := platform.GetFailureDomainByName(fdName); err == nil {
			peUUIDs = append(peUUIDs, fd.PrismElement.UUID)
		}
	}
	if len(peUUIDs) == 0 {
		peUUIDs = append(peUUIDs, platform.PrismElements[0].UUID)
	}

	for _, peUUID := range peUUIDs {
		peGPUs, err := GetGPUsForPE(ctx, nc, peUUID)
		if err != nil || len(peGPUs) == 0 {
			err = fmt.Errorf("no available GPUs found in Prism Element cluster (uuid: %s): %w", peUUID, err)
			fldErrs = append(fldErrs, field.InternalError(fldPath.Child("gpus"), err))
			return fldErrs
		}

		for _, gpu := range p.GPUs {
			switch gpu.Type {
			case machinev1.NutanixGPUIdentifierDeviceID:
				if gpu.DeviceID == nil {
					fldErrs = append(fldErrs, field.Required(fldPath.Child("gpus", "deviceID"), "missing gpu deviceID"))
				} else {
					_, err := GetGPUFromList(ctx, nc, gpu, peGPUs)
					if err != nil {
						fldErrs = append(fldErrs, field.Invalid(fldPath.Child("gpus", "deviceID"), *gpu.DeviceID, err.Error()))
					}
				}
			case machinev1.NutanixGPUIdentifierName:
				if gpu.Name == nil || *gpu.Name == "" {
					fldErrs = append(fldErrs, field.Required(fldPath.Child("gpus", "name"), "missing gpu name"))
				} else {
					_, err := GetGPUFromList(ctx, nc, gpu, peGPUs)
					if err != nil {
						fldErrs = append(fldErrs, field.Invalid(fldPath.Child("gpus", "name"), gpu.Name, err.Error()))
					}
				}
			default:
				errMsg := fmt.Sprintf("invalid gpu identifier type, the valid values: %q, %q.", machinev1.NutanixGPUIdentifierDeviceID, machinev1.NutanixGPUIdentifierName)
				fldErrs = append(fldErrs, field.Invalid(fldPath.Child("gpus", "type"), gpu.Type, errMsg))
			}
		}
	}

	return fldErrs
}

// validateDataDisksConfig validates the DataDisks configuration in the machinePool.
func (p *MachinePool) validateDataDisksConfig(ctx context.Context, nc *nutanixclientv3.Client, platform *Platform, fldPath *field.Path) (fldErrs []*field.Error) {
	var err error
	var errMsg string

	for _, disk := range p.DataDisks {
		// the minimum diskSize is 1Gi bytes
		diskSizeBytes := disk.DiskSize.Value()
		if diskSizeBytes < 1024*1024*1024 {
			fldErrs = append(fldErrs, field.Invalid(fldPath.Child("dataDisks", "diskSize"), fmt.Sprintf("%v bytes", diskSizeBytes), "The minimum diskSize is 1Gi bytes."))
		}

		if disk.DeviceProperties != nil {
			switch disk.DeviceProperties.DeviceType {
			case machinev1.NutanixDiskDeviceTypeDisk:
				switch disk.DeviceProperties.AdapterType {
				case machinev1.NutanixDiskAdapterTypeSCSI, machinev1.NutanixDiskAdapterTypeIDE, machinev1.NutanixDiskAdapterTypePCI, machinev1.NutanixDiskAdapterTypeSATA, machinev1.NutanixDiskAdapterTypeSPAPR:
					// valid configuration
				default:
					// invalid configuration
					fldErrs = append(fldErrs, field.Invalid(fldPath.Child("deviceProperties", "adapterType"), disk.DeviceProperties.AdapterType,
						fmt.Sprintf("invalid adapter type for the %q device type, the valid values: %q, %q, %q, %q, %q.",
							machinev1.NutanixDiskDeviceTypeDisk, machinev1.NutanixDiskAdapterTypeSCSI, machinev1.NutanixDiskAdapterTypeIDE,
							machinev1.NutanixDiskAdapterTypePCI, machinev1.NutanixDiskAdapterTypeSATA, machinev1.NutanixDiskAdapterTypeSPAPR)))
				}
			case machinev1.NutanixDiskDeviceTypeCDROM:
				switch disk.DeviceProperties.AdapterType {
				case machinev1.NutanixDiskAdapterTypeIDE, machinev1.NutanixDiskAdapterTypeSATA:
					// valid configuration
				default:
					// invalid configuration
					fldErrs = append(fldErrs, field.Invalid(fldPath.Child("deviceProperties", "adapterType"), disk.DeviceProperties.AdapterType,
						fmt.Sprintf("invalid adapter type for the %q device type, the valid values: %q, %q.",
							machinev1.NutanixDiskDeviceTypeCDROM, machinev1.NutanixDiskAdapterTypeIDE, machinev1.NutanixDiskAdapterTypeSATA)))
				}
			default:
				fldErrs = append(fldErrs, field.Invalid(fldPath.Child("deviceProperties", "deviceType"), disk.DeviceProperties.DeviceType,
					fmt.Sprintf("invalid device type, the valid types are: %q, %q.", machinev1.NutanixDiskDeviceTypeDisk, machinev1.NutanixDiskDeviceTypeCDROM)))
			}

			if disk.DeviceProperties.DeviceIndex < 0 {
				fldErrs = append(fldErrs, field.Invalid(fldPath.Child("deviceProperties", "deviceIndex"),
					disk.DeviceProperties.DeviceIndex, "invalid device index, the valid values are non-negative integers."))
			}
		}

		if disk.StorageConfig != nil {
			if disk.StorageConfig.DiskMode != machinev1.NutanixDiskModeStandard && disk.StorageConfig.DiskMode != machinev1.NutanixDiskModeFlash {
				fldErrs = append(fldErrs, field.Invalid(fldPath.Child("storageConfig", "diskMode"), disk.StorageConfig.DiskMode,
					fmt.Sprintf("invalid disk mode, the valid values: %q, %q.", machinev1.NutanixDiskModeStandard, machinev1.NutanixDiskModeFlash)))
			}

			storageContainerRef := disk.StorageConfig.StorageContainer
			if storageContainerRef != nil {
				if storageContainerRef.ReferenceName != "" {
					for _, fdName := range p.FailureDomains {
						_, err := platform.GetStorageContainerFromFailureDomain(fdName, storageContainerRef.ReferenceName)
						if err != nil {
							fldErrs = append(fldErrs, field.Invalid(fldPath.Child("storageConfig", "storageContainer", "referenceName"), storageContainerRef.ReferenceName,
								fmt.Sprintf("not found storageContainer with the referenceName in the failureDomain %q configuration.", fdName)))
						}
					}
				} else if storageContainerRef.UUID == "" {
					fldErrs = append(fldErrs, field.Required(fldPath.Child("storageConfig", "storageContainer", "uuid"), "missing storageContainer uuid"))
				}
			}
		}

		if disk.DataSourceImage != nil {
			dsImgRef := disk.DataSourceImage
			if dsImgRef.ReferenceName != "" {
				for _, fdName := range p.FailureDomains {
					_, err = platform.GetDataSourceImageFromFailureDomain(fdName, disk.DataSourceImage.ReferenceName)
					if err != nil {
						fldErrs = append(fldErrs, field.Invalid(fldPath.Child("storageConfig", "dataSourceImage", "referenceName"), disk.DataSourceImage.ReferenceName,
							fmt.Sprintf("not found datasource image with the referenceName in the failureDomain %q configuration.", fdName)))
					}
				}
			} else {
				switch {
				case dsImgRef.UUID != "":
					if _, err = nc.V3.GetImage(ctx, dsImgRef.UUID); err != nil {
						errMsg = fmt.Sprintf("failed to find the dataSource image with uuid %s: %v", dsImgRef.UUID, err)
						fldErrs = append(fldErrs, field.Invalid(fldPath.Child("dataDisks", "dataSourceImage", "uuid"), dsImgRef.UUID, errMsg))
					}
				case dsImgRef.Name != "":
					if dsImgUUID, err := FindImageUUIDByName(ctx, nc, dsImgRef.Name); err != nil {
						errMsg = fmt.Sprintf("failed to find the dataSource image with name %q: %v", dsImgRef.UUID, err)
						fldErrs = append(fldErrs, field.Invalid(fldPath.Child("dataDisks", "dataSourceImage", "name"), dsImgRef.Name, errMsg))
					} else {
						dsImgRef.UUID = *dsImgUUID
					}
				default:
					fldErrs = append(fldErrs, field.Required(fldPath.Child("dataDisks", "dataSourceImage"), "both the dataSourceImage's uuid and name are empty, you need to configure one."))
				}
			}
		}
	}

	return fldErrs
}
