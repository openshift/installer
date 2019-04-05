package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(p.VirtualCenters) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("virtualCenters"), "must include at least one vCenter"))
	}
	foundServer := false
	vcNames := map[string]bool{}
	for i, vc := range p.VirtualCenters {
		allErrs = append(allErrs, validateVirtualCenter(&vc, fldPath.Child("virtualCenters").Index(i))...)
		if vcNames[vc.Name] {
			allErrs = append(allErrs, field.Duplicate(fldPath.Child("virtualCenters").Index(i), vc.Name))
		}
		vcNames[vc.Name] = true
		if vc.Name == p.Workspace.Server {
			foundDatacenter := false
			for _, dc := range vc.Datacenters {
				if dc == p.Workspace.Datacenter {
					foundDatacenter = true
					break
				}
			}
			if p.Workspace.Datacenter != "" && !foundDatacenter {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("workspace").Child("datacenter"), p.Workspace.Datacenter, "workspace datacenter must be a datacenter in the workspace server"))
			}
			foundServer = true
		}
	}
	if len(p.Workspace.Server) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("workspace").Child("server"), "must specify the workspace server"))
	} else if !foundServer {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("workspace").Child("server"), p.Workspace.Server, "workspace server must be a specified vCenter"))
	}
	if len(p.Workspace.Datacenter) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("workspace").Child("datacenter"), "must specify the workspace datacenter"))
	}
	if len(p.Workspace.DefaultDatastore) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("workspace").Child("defaultDatastore"), "must specify the default datastore"))
	}
	if len(p.Workspace.Folder) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("workspace").Child("folder"), "must specify the VM folder"))
	}
	if len(p.SCSIControllerType) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("scsiControllerType"), "must specify the SCSI controller type"))
	}
	if len(p.PublicNetwork) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("publicNetwork"), "must specify the public VM network"))
	}
	return allErrs
}

func validateVirtualCenter(vc *vsphere.VirtualCenter, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(vc.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "vCenter must have a name"))
	}
	if len(vc.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("username"), "username required for each vCenter"))
	}
	if len(vc.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("password"), "password required for each vCenter"))
	}
	if len(vc.Datacenters) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("datacenters"), "must include at least one datacenter"))
	}
	dcs := map[string]bool{}
	for i, dc := range vc.Datacenters {
		if dcs[dc] {
			allErrs = append(allErrs, field.Duplicate(fldPath.Child("datacenters").Index(i), dc))
		}
		dcs[dc] = true
	}
	return allErrs
}
