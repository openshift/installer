package nutanix

import (
	"fmt"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
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
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGiB defines the size of disk in GiB.
	//
	// +optional
	DiskSizeGiB int64 `json:"diskSizeGiB,omitempty"`
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
}

// ValidateConfig validates the MachinePool configuration.
func (p *MachinePool) ValidateConfig(platform *Platform) error {
	nc, err := CreateNutanixClientFromPlatform(platform)
	if err != nil {
		return fmt.Errorf("fail to create nutanix client. %w", err)
	}

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
		switch p.Project.Type {
		case machinev1.NutanixIdentifierName:
			if p.Project.Name == nil || *p.Project.Name == "" {
				errList = append(errList, field.Required(fldPath.Child("project", "name"), "missing projct name"))
			} else {
				projectName := *p.Project.Name
				filter := fmt.Sprintf("name==%s", projectName)
				res, err := nc.V3.ListProject(&nutanixclientv3.DSMetadata{
					Filter: &filter,
				})
				switch {
				case err != nil:
					errMsg = fmt.Sprintf("failed to find project with name %q. error: %v", projectName, err)
					errList = append(errList, field.Invalid(fldPath.Child("project", "name"), projectName, errMsg))
				case len(res.Entities) == 0:
					errMsg = fmt.Sprintf("found no project with name %q.", projectName)
					errList = append(errList, field.Invalid(fldPath.Child("project", "name"), projectName, errMsg))
				case len(res.Entities) > 1:
					errMsg = fmt.Sprintf("found more than one (%v) projects with name %q.", len(res.Entities), projectName)
					errList = append(errList, field.Invalid(fldPath.Child("project", "name"), projectName, errMsg))
				default:
					p.Project.Type = machinev1.NutanixIdentifierUUID
					p.Project.UUID = res.Entities[0].Metadata.UUID
				}
			}
		case machinev1.NutanixIdentifierUUID:
			if p.Project.UUID == nil || *p.Project.UUID == "" {
				errList = append(errList, field.Required(fldPath.Child("project", "uuid"), "missing projct uuid"))
			} else {
				if _, err = nc.V3.GetProject(*p.Project.UUID); err != nil {
					errMsg = fmt.Sprintf("failed to get the project with uuid %s. error: %v", *p.Project.UUID, err)
					errList = append(errList, field.Invalid(fldPath.Child("project", "uuid"), *p.Project.UUID, errMsg))
				}
			}
		default:
			errMsg = fmt.Sprintf("invalid project identifier type, valid types are: %q, %q.", machinev1.NutanixIdentifierName, machinev1.NutanixIdentifierUUID)
			errList = append(errList, field.Invalid(fldPath.Child("project", "type"), p.Project.Type, errMsg))
		}
	}

	// validate categories if configured
	if len(p.Categories) > 0 {
		for _, category := range p.Categories {
			if _, err = nc.V3.GetCategoryValue(category.Key, category.Value); err != nil {
				errMsg = fmt.Sprintf("Failed to find the category with key %q and value %q. error: %v", category.Key, category.Value, err)
				errList = append(errList, field.Invalid(fldPath.Child("categories"), category, errMsg))
			}
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf(errList.ToAggregate().Error())
	}
	return nil
}
