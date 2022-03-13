package vsphere

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/ovfdeploy"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"

	"github.com/vmware/govmomi/vim25/types"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/vmworkflow"
)

func dataSourceVSphereOvfVMTemplate() *schema.Resource {

	vmConfigSpecSchema := map[string]*schema.Schema{
		"num_cpus": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of virtual processors to assign to this virtual machine.",
		},
		"num_cores_per_socket": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of cores to distribute amongst the CPUs in this virtual machine. If specified, the value supplied to num_cpus must be evenly divisible by this value.",
		},
		"cpu_hot_add_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Allow CPUs to be added to this virtual machine while it is running.",
		},
		"cpu_hot_remove_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Allow CPUs to be added to this virtual machine while it is running.",
		},
		"nested_hv_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable nested hardware virtualization on this virtual machine, facilitating nested virtualization in the guest.",
		},
		"cpu_performance_counters_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable CPU performance counters on this virtual machine.",
		},
		"memory": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The size of the virtual machine's memory, in MB.",
		},
		"memory_hot_add_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Allow memory to be added to this virtual machine while it is running.",
		},
		"swap_placement_policy": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The swap file placement policy for this virtual machine. Can be one of inherit, hostLocal, or vmDirectory.",
		},
		"annotation": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User-provided description of the virtual machine.",
		},
		"guest_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The guest ID for the operating system.",
		},
		"alternate_guest_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The guest name for the operating system when guest_id is other or other-64.",
		},
		"firmware": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The firmware interface to use on the virtual machine. Can be one of bios or EFI.",
		},
		"sata_controller_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of SATA controllers that Terraform manages on this virtual machine. This directly affects the amount of disks you can add to the virtual machine and the maximum disk unit number. Note that lowering this value does not remove controllers.",
		},
		"ide_controller_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of IDE controllers that Terraform manages on this virtual machine. This directly affects the amount of disks you can add to the virtual machine and the maximum disk unit number. Note that lowering this value does not remove controllers.",
		},
		"scsi_controller_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The number of SCSI controllers that Terraform manages on this virtual machine. This directly affects the amount of disks you can add to the virtual machine and the maximum disk unit number. Note that lowering this value does not remove controllers.",
		},
		"scsi_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of SCSI bus this virtual machine will have. Can be one of lsilogic, lsilogic-sas or pvscsi.",
		},
	}
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the virtual machine to create.",
		},
		"resource_pool_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of a resource pool to put the virtual machine in.",
		},

		"host_system_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of an optional host system to pin the virtual machine to.",
		},
		"datastore_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the virtual machine's datastore. The virtual machine configuration is placed here, along with any virtual disks that are created without datastores.",
		},
		"folder": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the folder to locate the virtual machine in.",
			StateFunc:   folder.NormalizePath,
		},
	}
	structure.MergeSchema(s, vmworkflow.VirtualMachineOvfDeploySchema())
	structure.MergeSchema(s, vmConfigSpecSchema)

	return &schema.Resource{
		Read:   dataSourceVSphereOvfVMTemplateRead,
		Schema: s,
	}
}

func NewOvfHelperParamsFromVMDatasource(d *schema.ResourceData) *ovfdeploy.OvfHelperParams {
	ovfParams := &ovfdeploy.OvfHelperParams{
		AllowUnverifiedSSL: d.Get("allow_unverified_ssl_cert").(bool),
		DatastoreId:        d.Get("datastore_id").(string),
		DeploymentOption:   d.Get("deployment_option").(string),
		DiskProvisioning:   d.Get("disk_provisioning").(string),
		FilePath:           d.Get("local_ovf_path").(string),
		Folder:             d.Get("folder").(string),
		HostId:             d.Get("host_system_id").(string),
		IpAllocationPolicy: d.Get("ip_allocation_policy").(string),
		IpProtocol:         d.Get("ip_protocol").(string),
		Name:               d.Get("name").(string),
		NetworkMappings:    d.Get("ovf_network_map").(map[string]interface{}),
		OvfUrl:             d.Get("remote_ovf_url").(string),
		PoolId:             d.Get("resource_pool_id").(string),
	}
	return ovfParams
}

func dataSourceVSphereOvfVMTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	ovfParams := NewOvfHelperParamsFromVMDatasource(d)
	ovfHelper, err := ovfdeploy.NewOvfHelper(client, ovfParams)
	if err != nil {
		return fmt.Errorf("while extracting OVF parameters: %s", err)
	}

	is, err := ovfHelper.GetImportSpec(client)
	if err != nil {
		return fmt.Errorf("while retrieving import spec: %s", err)
	}

	vmConfigSpec := is.ImportSpec.(*types.VirtualMachineImportSpec).ConfigSpec
	_ = d.Set("num_cpus", vmConfigSpec.NumCPUs)
	_ = d.Set("num_cores_per_socket", vmConfigSpec.NumCoresPerSocket)
	_ = d.Set("cpu_hot_add_enabled", vmConfigSpec.CpuHotAddEnabled)
	_ = d.Set("cpu_hot_remove_enabled", vmConfigSpec.CpuHotRemoveEnabled)
	_ = d.Set("nested_hv_enabled", vmConfigSpec.NestedHVEnabled)
	_ = d.Set("memory", vmConfigSpec.MemoryMB)
	_ = d.Set("memory_hot_add_enabled", vmConfigSpec.MemoryHotAddEnabled)
	_ = d.Set("swap_placement_policy", vmConfigSpec.SwapPlacement)
	_ = d.Set("annotation", vmConfigSpec.Annotation)
	_ = d.Set("guest_id", vmConfigSpec.GuestId)
	_ = d.Set("alternate_guest_name", vmConfigSpec.AlternateGuestName)
	_ = d.Set("firmware", vmConfigSpec.Firmware)

	controllers := map[string]int{}
	var scsiType string

	for _, dvc := range vmConfigSpec.DeviceChange {
		dvcSpec := dvc.GetVirtualDeviceConfigSpec()

		switch reflect.TypeOf(dvcSpec.Device) {
		case reflect.TypeOf(&types.VirtualLsiLogicController{}):
			if scsiType == "" {
				scsiType = "lsilogic"
			} else if scsiType != "lsilogic" {
				return fmt.Errorf("multiple scsi controller types are not supported (found %s and %s)", scsiType, "lsilogic")
			}
			controllers["scsi"] = controllers["scsi"] + 1
		case reflect.TypeOf(&types.VirtualLsiLogicSASController{}):
			if scsiType == "" {
				scsiType = "lsilogic-sas"
			} else if scsiType != "lsilogic-sas" {
				return fmt.Errorf("multiple scsi controller types are not supported (found %s and %s)", scsiType, "lsilogic-sas")
			}
			controllers["scsi"] = controllers["scsi"] + 1
		case reflect.TypeOf(&types.ParaVirtualSCSIController{}):
			if scsiType == "" {
				scsiType = "pvscsi"
			} else if scsiType != "pvscsi" {
				return fmt.Errorf("multiple scsi controller types are not supported (found %s and %s)", scsiType, "pvsci")
			}
			controllers["scsi"] = controllers["scsi"] + 1
		case reflect.TypeOf(&types.VirtualSATAController{}):
			controllers["sata"] = controllers["sata"] + 1
		case reflect.TypeOf(&types.VirtualIDEController{}):
			controllers["ide"] = controllers["ide"] + 1
		}
	}

	_ = d.Set("scsi_type", scsiType)
	_ = d.Set("scsi_controller_count", controllers["scsi"])
	_ = d.Set("sata_controller_count", controllers["sata"])
	_ = d.Set("ide_controller_count", controllers["ide"])

	d.SetId(d.Get("name").(string))

	return nil
}
