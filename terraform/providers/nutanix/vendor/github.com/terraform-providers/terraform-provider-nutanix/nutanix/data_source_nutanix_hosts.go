package nutanix

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func dataSourceNutanixHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixHostsRead,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain_credencial": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"password": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"ipmi": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name_server_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"organization_unit_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain_credencial": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"password": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_vms": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hypervisor_full_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
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
						"controller_vm": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nat_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"oplog_usage": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"opolog_disk_pct": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"opolog_disk_size": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"nat_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"block": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_serial_number": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"block_model": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_reference": {
							Type:     schema.TypeMap,
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
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_hash": {
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
						"categories": categoriesSchema(),
						"project_reference": {
							Type:     schema.TypeMap,
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
						"owner_reference": {
							Type:     schema.TypeMap,
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
					},
				},
			},
		},
	}
}

func dataSourceNutanixHostsRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllHost()
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}
	if err := d.Set("entities", flattenHostEntities(resp.Entities)); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenHostEntities(hosts []*v3.HostResponse) []map[string]interface{} {
	entities := make([]map[string]interface{}, len(hosts))

	for i, host := range hosts {
		metadata, categories := setRSEntityMetadata(host.Metadata)

		entities[i] = map[string]interface{}{
			"name":                      host.Status.Name,
			"gpu_driver_version":        host.Status.Resources.GPUDriverVersion,
			"failover_cluster":          flattenFailOverCluster(host.Status.Resources.FailoverCluster),
			"ipmi":                      flattenIMPI(host.Status.Resources.IPMI),
			"cpu_model":                 host.Status.Resources.CPUModel,
			"host_nics_id_list":         host.Status.Resources.HostNicsIDList,
			"num_cpu_sockets":           host.Status.Resources.NumCPUSockets,
			"windows_domain":            flattenWindowsDomain(host.Status.Resources.WindowsDomain),
			"gpu_list":                  flattenGpuList(host.Status.Resources.GPUList),
			"serial_number":             host.Status.Resources.SerialNumber,
			"cpu_capacity_hz":           host.Status.Resources.CPUCapacityHZ,
			"memory_capacity_mib":       host.Status.Resources.MemoryVapacityMib,
			"host_disks_reference_list": flattenReferenceList(host.Status.Resources.HostDisksReferenceList),
			"monitoring_state":          host.Status.Resources.MonitoringState,
			"hypervisor":                flattenHypervisor(host.Status.Resources.Hypervisor),
			"host_type":                 host.Status.Resources.HostType,
			"num_cpu_cores":             host.Status.Resources.NumCPUCores,
			"rackable_unit_reference":   flattenReference(host.Status.Resources.RackableUnitReference),
			"controller_vm":             flattenControllerVM(host.Status.Resources.ControllerVM),
			"block":                     flattenBlock(host.Status.Resources.Block),
			"cluster_reference":         flattenReference(host.Status.ClusterReference),
			"metadata":                  metadata,
			"categories":                categories,
			"project_reference":         flattenReferenceValues(host.Metadata.ProjectReference),
			"owner_reference":           flattenReferenceValues(host.Metadata.OwnerReference),
			"api_version":               host.APIVersion,
		}
	}
	return entities
}
