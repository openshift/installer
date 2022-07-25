package nutanix

import (
	"context"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func dataSourceNutanixHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixHostRead,
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_driver_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failover_cluster": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipmi": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cpu_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_nics_id_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_cpu_sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"windows_domain": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gpu_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_virtual_display_heads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"assignable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"license_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"num_vgpus_allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frame_buffer_size_mib": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"numa_node": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_resolution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fraction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"guest_driver_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_capacity_hz": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_capacity_mib": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_disks_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"monitoring_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_cpu_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rackable_unit_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"controller_vm": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"block": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNutanixHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*Client).API

	hostID := d.Get("host_id").(string)

	host, err := conn.V3.GetHost(hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", host.Status.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_driver_version", host.Status.Resources.GPUDriverVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("failover_cluster", flattenFailOverCluster(host.Status.Resources.FailoverCluster)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipmi", flattenIMPI(host.Status.Resources.IPMI)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_model", host.Status.Resources.CPUModel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_nics_id_list", host.Status.Resources.HostNicsIDList); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_cpu_sockets", host.Status.Resources.NumCPUSockets); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("windows_domain", flattenWindowsDomain(host.Status.Resources.WindowsDomain)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_list", flattenGpuList(host.Status.Resources.GPUList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("serial_number", host.Status.Resources.SerialNumber); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_capacity_hz", host.Status.Resources.CPUCapacityHZ); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_capacity_mib", host.Status.Resources.MemoryVapacityMib); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_disks_reference_list", flattenReferenceList(host.Status.Resources.HostDisksReferenceList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("monitoring_state", host.Status.Resources.MonitoringState); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor", flattenHypervisor(host.Status.Resources.Hypervisor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_type", host.Status.Resources.HostType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_cpu_cores", host.Status.Resources.NumCPUCores); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rackable_unit_reference", flattenReference(host.Status.Resources.RackableUnitReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_vm", flattenControllerVM(host.Status.Resources.ControllerVM)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("block", flattenBlock(host.Status.Resources.Block)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", flattenReference(host.Status.ClusterReference)); err != nil {
		return diag.FromErr(err)
	}
	m, c := setRSEntityMetadata(host.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_reference", flattenReferenceValues(host.Metadata.ProjectReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_reference", flattenReferenceValues(host.Metadata.OwnerReference)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hostID)
	return nil
}

func flattenFailOverCluster(failOvercluster *v3.FailoverCluster) map[string]interface{} {
	if failOvercluster != nil {
		return map[string]interface{}{
			"ip":   failOvercluster.IP,
			"name": failOvercluster.Name,
		}
	}
	return map[string]interface{}{}
}

func flattenDomainCredencials(domainCredencial *v3.DomainCredencial) map[string]interface{} {
	if domainCredencial != nil {
		return map[string]interface{}{
			"username": domainCredencial.Username,
			"password": domainCredencial.Password,
		}
	}
	return map[string]interface{}{}
}

func flattenIMPI(impi *v3.IPMI) map[string]interface{} {
	if impi != nil {
		return map[string]interface{}{
			"ip": impi.IP,
		}
	}
	return map[string]interface{}{}
}

func flattenWindowsDomain(windowsDomain *v3.WindowsDomain) map[string]interface{} {
	if windowsDomain != nil {
		return map[string]interface{}{
			"name":                   windowsDomain.Name,
			"name_server_ip":         windowsDomain.NameServerIP,
			"organization_unit_path": windowsDomain.OrganizationUnitPath,
			"name_prefix":            windowsDomain.NamePrefix,
		}
	}
	return map[string]interface{}{}
}

func flattenGpuList(cpuList []*v3.GPU) []map[string]interface{} {
	res := make([]map[string]interface{}, len(cpuList))
	if len(cpuList) > 0 {
		for i, cpu := range cpuList {
			res[i] = map[string]interface{}{
				"status":                    cpu.Status,
				"vendor":                    cpu.Vendor,
				"num_virtual_display_heads": cpu.NumVirtualDisplayHeads,
				"assignable":                cpu.Assignable,
				"license_list":              cpu.LicenseList,
				"num_vgpus_allocated":       cpu.NumVgpusAllocated,
				"pci_address":               cpu.PciAddress,
				"name":                      cpu.Name,
				"frame_buffer_size_mib":     cpu.FrameBufferSizeMib,
				"index":                     cpu.Index,
				"uuid":                      cpu.UUID,
				"numa_node":                 cpu.NumaNode,
				"max_resolution":            cpu.MaxResoution,
				"consumer_reference":        flattenReference(cpu.ConsumerReference),
				"mode":                      cpu.Mode,
				"fraction":                  cpu.Fraction,
				"guest_driver_version":      cpu.GuestDriverVersion,
				"device_id":                 cpu.DeviceID,
			}
		}
	}
	return res
}

func flattenReference(reference *v3.ReferenceValues) map[string]interface{} {
	if reference != nil {
		return map[string]interface{}{
			"kind": reference.Kind,
			"uuid": reference.UUID,
			"name": reference.Name,
		}
	}
	return map[string]interface{}{}
}

func flattenReferenceList(references []*v3.ReferenceValues) []map[string]interface{} {
	res := make([]map[string]interface{}, len(references))
	if len(references) > 0 {
		for i, r := range references {
			res[i] = flattenReference(r)
		}
	}
	return res
}

func flattenExternalNetworkListReference(reference *v3.ReferenceValues) map[string]interface{} {
	if reference != nil {
		return map[string]interface{}{
			"uuid": reference.UUID,
			"name": reference.Name,
		}
	}
	return map[string]interface{}{}
}

func flattenExternalNetworkListReferenceList(references []*v3.ReferenceValues) []map[string]interface{} {
	res := make([]map[string]interface{}, len(references))
	if len(references) > 0 {
		for i, r := range references {
			res[i] = flattenExternalNetworkListReference(r)
		}
	}
	return res
}

func flattenHypervisor(hypervisor *v3.Hypervisor) map[string]interface{} {
	if hypervisor != nil {
		return map[string]interface{}{
			"num_vms":              cast.ToString(hypervisor.NumVms),
			"ip":                   hypervisor.IP,
			"hypervisor_full_name": hypervisor.HypervisorFullName,
		}
	}
	return map[string]interface{}{}
}

func flattenControllerVM(controllerVM *v3.ControllerVM) map[string]interface{} {
	if controllerVM != nil {
		return map[string]interface{}{
			"ip":                           controllerVM.IP,
			"nat_ip":                       controllerVM.NatIP,
			"oplog_usage.opolog_disk_pct":  cast.ToString(controllerVM.OplogUsage.OplogDiskPct),
			"oplog_usage.opolog_disk_size": cast.ToString(controllerVM.OplogUsage.OplogDiskSize),
			"nat_port":                     controllerVM.NatPort,
		}
	}
	return map[string]interface{}{}
}

func flattenBlock(block *v3.Block) map[string]interface{} {
	if block != nil {
		return map[string]interface{}{
			"block_serial_number": block.BlockSerialNumber,
			"block_model":         block.BlockModel,
		}
	}
	return map[string]interface{}{}
}
