package virtualmachineinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func domainSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resources": {
			Type:        schema.TypeList,
			Description: "Resources describes the Compute Resources required by this vmi.",
			MaxItems:    1,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"requests": {
						Type:        schema.TypeMap,
						Description: "Requests is a description of the initial vmi resources.",
						Optional:    true,
					},
					"limits": {
						Type:        schema.TypeMap,
						Description: "Requests is a description of the initial vmi resources.",
						Optional:    true,
					},
					"over_commit_guest_overhead": {
						Type:        schema.TypeBool,
						Description: "Don't ask the scheduler to take the guest-management overhead into account. Instead put the overhead only into the container's memory limit. This can lead to crashes if all memory is in use on a node. Defaults to false.",
						Optional:    true,
					},
				},
			},
		},
		"devices": {
			Type:        schema.TypeList,
			Description: "Devices allows adding disks, network interfaces, ...",
			MaxItems:    1,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk": {
						Type:        schema.TypeList,
						Description: "Disks describes disks, cdroms, floppy and luns which are connected to the vmi.",
						Required:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "Name is the device name",
									Required:    true,
								},
								"disk_device": {
									Type:        schema.TypeList,
									Description: "DiskDevice specifies as which device the disk should be added to the guest.",
									Required:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"disk": {
												Type:        schema.TypeList,
												Description: "Attach a volume as a disk to the vmi.",
												Optional:    true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"bus": {
															Type:        schema.TypeString,
															Description: "Bus indicates the type of disk device to emulate.",
															Required:    true,
														},
														"read_only": {
															Type:        schema.TypeBool,
															Description: "ReadOnly. Defaults to false.",
															Optional:    true,
														},
														"pci_address": {
															Type:        schema.TypeString,
															Description: "If specified, the virtual disk will be placed on the guests pci address with the specifed PCI address. For example: 0000:81:01.10",
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
								"serial": {
									Type:        schema.TypeString,
									Description: "Serial provides the ability to specify a serial number for the disk device.",
									Optional:    true,
								},
							},
						},
					},
					"interface": {
						Type:        schema.TypeList,
						Description: "Interfaces describe network interfaces which are added to the vmi.",
						Required:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "Logical name of the interface as well as a reference to the associated networks.",
									Required:    true,
								},
								"interface_binding_method": {
									Type: schema.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"InterfaceBridge",
										"InterfaceSlirp",
										"InterfaceMasquerade",
										"InterfaceSRIOV",
									}, false),
									Description: "Represents the method which will be used to connect the interface to the guest.",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func domainSpecSchema() *schema.Schema {
	fields := domainSpecFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Specification of the desired behavior of the VirtualMachineInstance on the host."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandDomainSpec(domainSpec []interface{}) (kubevirtapiv1.DomainSpec, error) {
	result := kubevirtapiv1.DomainSpec{}

	if len(domainSpec) == 0 || domainSpec[0] == nil {
		return result, nil
	}

	in := domainSpec[0].(map[string]interface{})

	if v, ok := in["resources"].([]interface{}); ok {
		resources, err := expandResources(v)
		if err != nil {
			return result, err
		}
		result.Resources = resources
	}
	if v, ok := in["devices"].([]interface{}); ok {
		devices, err := expandDevices(v)
		if err != nil {
			return result, err
		}
		result.Devices = devices
	}

	return result, nil
}

func expandResources(resources []interface{}) (kubevirtapiv1.ResourceRequirements, error) {
	result := kubevirtapiv1.ResourceRequirements{}

	if len(resources) == 0 || resources[0] == nil {
		return result, nil
	}

	in := resources[0].(map[string]interface{})

	if v, ok := in["requests"].(map[string]interface{}); ok {
		requests, err := utils.ExpandMapToResourceList(v)
		if err != nil {
			return result, err
		}
		result.Requests = *requests
	}
	if v, ok := in["limits"].(map[string]interface{}); ok {
		limits, err := utils.ExpandMapToResourceList(v)
		if err != nil {
			return result, err
		}
		result.Limits = *limits
	}
	if v, ok := in["over_commit_guest_overhead"].(bool); ok {
		result.OvercommitGuestOverhead = v
	}

	return result, nil
}

func expandDevices(devices []interface{}) (kubevirtapiv1.Devices, error) {
	result := kubevirtapiv1.Devices{}

	if len(devices) == 0 || devices[0] == nil {
		return result, nil
	}

	in := devices[0].(map[string]interface{})

	if v, ok := in["disk"].([]interface{}); ok {
		result.Disks = expandDisks(v)
	}
	if v, ok := in["interface"].([]interface{}); ok {
		result.Interfaces = expandInterfaces(v)
	}

	return result, nil
}

func expandDisks(disks []interface{}) []kubevirtapiv1.Disk {
	result := make([]kubevirtapiv1.Disk, len(disks))

	if len(disks) == 0 || disks[0] == nil {
		return result
	}

	for i, condition := range disks {
		in := condition.(map[string]interface{})

		if v, ok := in["name"].(string); ok {
			result[i].Name = v
		}
		if v, ok := in["disk_device"].([]interface{}); ok {
			result[i].DiskDevice = expandDiskDevice(v)
		}
		if v, ok := in["serial"].(string); ok {
			result[i].Serial = v
		}
	}

	return result
}

func expandDiskDevice(diskDevice []interface{}) kubevirtapiv1.DiskDevice {
	result := kubevirtapiv1.DiskDevice{}

	if len(diskDevice) == 0 || diskDevice[0] == nil {
		return result
	}

	in := diskDevice[0].(map[string]interface{})

	if v, ok := in["disk"].([]interface{}); ok {
		result.Disk = expandDiskTarget(v)
	}

	return result
}

func expandDiskTarget(disk []interface{}) *kubevirtapiv1.DiskTarget {
	if len(disk) == 0 || disk[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.DiskTarget{}

	in := disk[0].(map[string]interface{})

	if v, ok := in["bus"].(string); ok {
		result.Bus = v
	}
	if v, ok := in["read_only"].(bool); ok {
		result.ReadOnly = v
	}
	if v, ok := in["pci_address"].(string); ok {
		result.PciAddress = v
	}

	return result
}

func expandInterfaces(interfaces []interface{}) []kubevirtapiv1.Interface {
	result := make([]kubevirtapiv1.Interface, len(interfaces))

	if len(interfaces) == 0 || interfaces[0] == nil {
		return result
	}

	for i, condition := range interfaces {
		in := condition.(map[string]interface{})

		if v, ok := in["name"].(string); ok {
			result[i].Name = v
		}
		if v, ok := in["interface_binding_method"].(string); ok {
			result[i].InterfaceBindingMethod = expandInterfaceBindingMethod(v)
		}
	}

	return result
}

func expandInterfaceBindingMethod(interfaceBindingMethod string) kubevirtapiv1.InterfaceBindingMethod {
	result := kubevirtapiv1.InterfaceBindingMethod{}

	switch interfaceBindingMethod {
	case "InterfaceBridge":
		result.Bridge = &kubevirtapiv1.InterfaceBridge{}
	case "InterfaceSlirp":
		result.Slirp = &kubevirtapiv1.InterfaceSlirp{}
	case "InterfaceMasquerade":
		result.Masquerade = &kubevirtapiv1.InterfaceMasquerade{}
	case "InterfaceSRIOV":
		result.SRIOV = &kubevirtapiv1.InterfaceSRIOV{}
	}

	return result
}

func flattenDomainSpec(in kubevirtapiv1.DomainSpec) []interface{} {
	att := make(map[string]interface{})

	att["resources"] = flattenResources(in.Resources)
	att["devices"] = flattenDevices(in.Devices)

	return []interface{}{att}
}

func flattenResources(in kubevirtapiv1.ResourceRequirements) []interface{} {
	att := make(map[string]interface{})

	att["requests"] = utils.FlattenStringMap(utils.FlattenResourceList(in.Requests))
	att["limits"] = utils.FlattenStringMap(utils.FlattenResourceList(in.Limits))
	att["over_commit_guest_overhead"] = in.OvercommitGuestOverhead

	return []interface{}{att}
}

func flattenDevices(in kubevirtapiv1.Devices) []interface{} {
	att := make(map[string]interface{})

	att["disk"] = flattenDisks(in.Disks)
	att["interface"] = flattenInterfaces(in.Interfaces)

	return []interface{}{att}
}

func flattenDisks(in []kubevirtapiv1.Disk) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})

		c["name"] = v.Name
		c["disk_device"] = flattenDiskDevice(v.DiskDevice)
		c["serial"] = v.Serial

		att[i] = c
	}

	return att
}

func flattenDiskDevice(in kubevirtapiv1.DiskDevice) []interface{} {
	att := make(map[string]interface{})

	if in.Disk != nil {
		att["disk"] = flattenDiskTarget(*in.Disk)
	}

	return []interface{}{att}
}

func flattenDiskTarget(in kubevirtapiv1.DiskTarget) []interface{} {
	att := make(map[string]interface{})

	att["bus"] = in.Bus
	att["read_only"] = in.ReadOnly
	att["pci_address"] = in.PciAddress

	return []interface{}{att}
}

func flattenInterfaces(in []kubevirtapiv1.Interface) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})

		c["name"] = v.Name
		c["interface_binding_method"] = flattenInterfaceBindingMethod(v.InterfaceBindingMethod)

		att[i] = c
	}

	return att
}

func flattenInterfaceBindingMethod(in kubevirtapiv1.InterfaceBindingMethod) string {
	if in.Bridge != nil {
		return "InterfaceBridge"
	}
	if in.Slirp != nil {
		return "InterfaceSlirp"
	}
	if in.Masquerade != nil {
		return "InterfaceMasquerade"
	}
	if in.SRIOV != nil {
		return "InterfaceSRIOV"
	}

	return ""
}
