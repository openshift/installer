package nutanix

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	vmTimeout    = 1 * time.Minute
	vmDelay      = 3 * time.Second
	vmMinTimeout = 3 * time.Second
	IDE          = "IDE"
	useHotAdd    = true
)

func resourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVirtualMachineCreate,
		Read:   resourceNutanixVirtualMachineRead,
		Update: resourceNutanixVirtualMachineUpdate,
		Delete: resourceNutanixVirtualMachineDelete,
		Exists: resourceNutanixVirtualMachineExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixVirtualMachineInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceVirtualMachineInstanceStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"cloud_init_cdrom_uuid": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(
						"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"),
					"please see http://developer.nutanix.com/reference/prism_central/v3/api/models/cluster-reference"),
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nic_list_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_endpoint_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// RESOURCES ARGUMENTS

			"enable_cpu_passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"use_hot_add": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"num_vnuma_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"nic_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_endpoint_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"subnet_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "true",
						},
					},
				},
			},
			"guest_os_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nutanix_guest_tools": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ngt_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"iso_mount_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"available_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guest_os_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vss_snapshot_capable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_reachable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_mobility_drivers_installed": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ngt_credentials": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"ngt_enabled_capability_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_vcpus_per_socket": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"gpu_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frame_buffer_size_mib": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Optional: true,
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
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fraction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"num_virtual_display_heads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"guest_driver_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"parent_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"memory_size_mib": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"boot_device_order_list": {
				Type: schema.TypeList,
				// // remove MaxItems when the issue #28 is fixed
				// MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"boot_device_disk_address": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_index": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"adapter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"boot_device_mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"boot_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"UEFI", "LEGACY", "SECURE_BOOT"}, false),
			},
			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hardware_clock_timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guest_customization_cloud_init_user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"guest_customization_cloud_init_meta_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"guest_customization_cloud_init_custom_key_values": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"guest_customization_is_overridable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"guest_customization_sysprep": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"unattend_xml": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"guest_customization_sysprep_custom_key_values": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"should_fail_on_script_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_script_exec": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"power_state_mechanism": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vga_console_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disk_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_size_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"disk_size_mib": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"storage_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flash_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"storage_container_reference": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"kind": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "storage_container",
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"device_properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"disk_address": {
										Type:     schema.TypeMap,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"device_index": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"adapter_type": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"data_source_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"volume_group_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"serial_port_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API
	setVMTimeout(meta)
	// Prepare request
	request := &v3.VMIntentInput{}
	spec := &v3.VM{}
	metadata := &v3.Metadata{}
	res := &v3.VMResources{}

	// Read Arguments and set request values
	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")
	azr, azrok := d.GetOk("availability_zone_reference")
	clusterUUID, crok := d.GetOk("cluster_uuid")

	if !nok {
		return fmt.Errorf("please provide the required name attribute")
	}
	if err := getMetadataAttributes(d, metadata, "vm"); err != nil {
		return fmt.Errorf("error reading metadata for Virtual Machine %s", err)
	}
	if descok {
		spec.Description = utils.StringPtr(desc.(string))
	}
	if azrok {
		a := azr.(map[string]interface{})
		spec.AvailabilityZoneReference = validateRef(a)
	}
	if crok {
		spec.ClusterReference = buildReference(clusterUUID.(string), "cluster")
	}

	if err := getVMResources(d, res); err != nil {
		return err
	}

	spec.Name = utils.StringPtr(n.(string))
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.CreateVM(request)
	if err != nil {
		return err
	}

	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForState(); errWaitTask != nil {
		return fmt.Errorf("error waiting for vm (%s) to create: %s", uuid, errWaitTask)
	}

	// Wait for IP available
	waitIPConf := &resource.StateChangeConf{
		Pending:    []string{WAITING},
		Target:     []string{"AVAILABLE"},
		Refresh:    waitForIPRefreshFunc(conn, uuid),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	vmIntentResponse, err := waitIPConf.WaitForState()
	if err != nil {
		log.Printf("[WARN] could not get the IP for VM(%s): %s", uuid, err)
	} else {
		vm := vmIntentResponse.(*v3.VMIntentResponse)

		if len(vm.Status.Resources.NicList) > 0 && len(vm.Status.Resources.NicList[0].IPEndpointList) != 0 {
			d.SetConnInfo(map[string]string{
				"type": "ssh",
				"host": *vm.Status.Resources.NicList[0].IPEndpointList[0].IP,
			})
		}
	}

	// Set terraform state id
	d.SetId(uuid)
	return resourceNutanixVirtualMachineRead(d, meta)
}

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API
	setVMTimeout(meta)

	var err error
	// Make request to the API
	resp, err := conn.V3.GetVM(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Virtual Machine %s: %s", d.Id(), err)
	}

	// Added check for deletion. Re-running TF right after VM deletion, can cause an error because the ID is still present in API.
	// Check if name is not present and also resources is not present
	if resp.Status.Name == nil && resp.Status.Resources == nil {
		d.SetId("")
		return nil
	}

	if err = flattenClusterReference(resp.Status.ClusterReference, d); err != nil {
		return fmt.Errorf("error setting cluster information for Virtual Machine %s: %s", d.Id(), err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return fmt.Errorf("error setting metadata for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("categories", c); err != nil {
		return fmt.Errorf("error setting categories for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("project_reference", flattenReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return fmt.Errorf("error setting project_reference for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return fmt.Errorf("error setting owner_reference for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("availability_zone_reference", flattenReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return fmt.Errorf("error setting availability_zone_reference for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("nic_list", flattenNicList(resp.Spec.Resources.NicList)); err != nil {
		return fmt.Errorf("error setting nic_list for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("nic_list_status", flattenNicListStatus(resp.Status.Resources.NicList)); err != nil {
		return fmt.Errorf("error setting nic_list_status for Virtual Machine %s: %s", d.Id(), err)
	}
	flatDiskList, err := flattenDiskListFilterCloudInit(d, resp.Spec.Resources.DiskList)
	if err != nil {
		return fmt.Errorf("error flattening disk list for vm %s: %s", d.Id(), err)
	}
	if err = d.Set("disk_list", flatDiskList); err != nil {
		return fmt.Errorf("error setting disk_list for Virtual Machine %s: %s", d.Id(), err)
	}

	if err = d.Set("serial_port_list", flattenSerialPortList(resp.Status.Resources.SerialPortList)); err != nil {
		return fmt.Errorf("error setting serial_port_list for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("host_reference", flattenReferenceValues(resp.Status.Resources.HostReference)); err != nil {
		return fmt.Errorf("error setting host_reference for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = flattenNutanixGuestTools(d, resp.Status.Resources.GuestTools); err != nil {
		return fmt.Errorf("error setting nutanix_guest_tools for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("gpu_list", flattenGPUList(resp.Status.Resources.GpuList)); err != nil {
		return fmt.Errorf("error setting gpu_list for Virtual Machine %s: %s", d.Id(), err)
	}
	if err = d.Set("parent_reference", flattenReferenceValues(resp.Status.Resources.ParentReference)); err != nil {
		return fmt.Errorf("error setting parent_reference for Virtual Machine %s: %s", d.Id(), err)
	}

	if uha, ok := d.GetOkExists("use_hot_add"); ok {
		useHotAdd = uha.(bool)
	}
	if err = d.Set("use_hot_add", useHotAdd); err != nil {
		return fmt.Errorf("error setting use_hot_add for Virtual Machine %s: %s", d.Id(), err)
	}

	diskAddress := make(map[string]interface{})
	mac := ""
	bootType := ""
	b := make([]string, 0)

	log.Printf("[DEBUG] checking BootConfig %+v", resp.Status.Resources.BootConfig)
	if resp.Status.Resources.BootConfig != nil {
		if resp.Status.Resources.BootConfig.BootDevice != nil {
			if resp.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
				i := strconv.Itoa(int(utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)))
				diskAddress["device_index"] = i
				diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
			}
			mac = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)
		}
		if resp.Status.Resources.BootConfig.BootDeviceOrderList != nil {
			log.Printf("[DEBUG] checking BootConfig.BootDeviceOrderList %+v", utils.StringValueSlice(resp.Status.Resources.BootConfig.BootDeviceOrderList))
			b = utils.StringValueSlice(resp.Status.Resources.BootConfig.BootDeviceOrderList)
		}
		if resp.Status.Resources.BootConfig.BootType != nil {
			bootType = utils.StringValue(resp.Status.Resources.BootConfig.BootType)
		}
	}

	if err = d.Set("boot_device_order_list", b); err != nil {
		return fmt.Errorf("error setting boot_device_order_list %s", err)
	}

	d.Set("boot_device_disk_address", diskAddress)
	d.Set("boot_device_mac_address", mac)
	d.Set("boot_type", bootType)
	d.Set("machine_type", resp.Status.Resources.MachineType)

	cloudInitUser := ""
	cloudInitMeta := ""
	sysprep := make(map[string]interface{})
	sysprepCV := make(map[string]string)
	cloudInitCV := make(map[string]string)
	isOv := false
	if resp.Status.Resources.GuestCustomization != nil {
		isOv = utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)

		if resp.Status.Resources.GuestCustomization.CloudInit != nil {
			cloudInitMeta = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInitUser = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
			if resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues != nil {
				for k, v := range resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues {
					cloudInitCV[k] = v
				}
			}
		}

		if resp.Status.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
			if resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues != nil {
				for k, v := range resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues {
					sysprepCV[k] = v
				}
			}
		}
	}
	if err := d.Set("guest_customization_cloud_init_custom_key_values", cloudInitCV); err != nil {
		return fmt.Errorf("error setting guest_customization_cloud_init_custom_key_values for Virtual Machine %s: %s", d.Id(), err)
	}
	if err := d.Set("guest_customization_sysprep_custom_key_values", sysprepCV); err != nil {
		return fmt.Errorf("error setting guest_customization_sysprep_custom_key_values for Virtual Machine %s: %s", d.Id(), err)
	}
	if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
		return fmt.Errorf("error setting guest_customization_sysprep for Virtual Machine %s: %s", d.Id(), err)
	}

	d.Set("enable_cpu_passthrough", resp.Status.Resources.EnableCPUPassthrough)
	d.Set("guest_customization_cloud_init_user_data", cloudInitUser)
	d.Set("guest_customization_cloud_init_meta_data", cloudInitMeta)
	d.Set("hardware_clock_timezone", utils.StringValue(resp.Status.Resources.HardwareClockTimezone))
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("num_vnuma_nodes", utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes))
	d.Set("guest_os_id", utils.StringValue(resp.Status.Resources.GuestOsID))
	d.Set("power_state", utils.StringValue(resp.Status.Resources.PowerState))
	d.Set("num_vcpus_per_socket", utils.Int64Value(resp.Status.Resources.NumVcpusPerSocket))
	d.Set("num_sockets", utils.Int64Value(resp.Status.Resources.NumSockets))
	d.Set("memory_size_mib", utils.Int64Value(resp.Status.Resources.MemorySizeMib))
	d.Set("guest_customization_is_overridable", isOv)
	d.Set("should_fail_on_script_failure", utils.BoolValue(
		resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure))
	d.Set("enable_script_exec", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec))
	d.Set("power_state_mechanism", utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism))
	d.Set("vga_console_enabled", utils.BoolValue(resp.Status.Resources.VgaConsoleEnabled))
	d.SetId(*resp.Metadata.UUID)
	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	setVMTimeout(meta)
	hotPlugChange := true

	log.Printf("[Debug] Updating VM values %s", d.Id())

	request := &v3.VMIntentInput{}
	metadata := &v3.Metadata{}
	res := &v3.VMResources{}
	spec := &v3.VM{}
	guest := &v3.GuestCustomization{}
	guestTool := &v3.GuestToolsSpec{}
	pw := &v3.VMPowerStateMechanism{}

	response, err := conn.V3.GetVM(d.Id())

	//prefill structs
	preFillResUpdateRequest(res, response)
	preFillGTUpdateRequest(guestTool, response)
	preFillGUpdateRequest(guest, response)
	preFillPWUpdateRequest(pw, response)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if d.HasChange("use_hot_add") {
		useHotAdd = d.Get("use_hot_add").(bool)
	}
	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}
	metadata.OwnerReference = response.Metadata.OwnerReference
	if d.HasChange("owner_reference") {
		_, n := d.GetChange("owner_reference")
		metadata.OwnerReference = validateRef(n.(map[string]interface{}))
	}
	metadata.ProjectReference = response.Metadata.ProjectReference

	if d.HasChange("project_reference") {
		_, n := d.GetChange("project_reference")
		metadata.ProjectReference = validateRef(n.(map[string]interface{}))
	}

	spec.Name = response.Status.Name

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		spec.Name = utils.StringPtr(n.(string))
		hotPlugChange = false
	}

	spec.Description = response.Status.Description
	if d.HasChange("description") {
		_, n := d.GetChange("description")
		spec.Description = utils.StringPtr(n.(string))
	}
	spec.AvailabilityZoneReference = response.Status.AvailabilityZoneReference
	if d.HasChange("availability_zone_reference") {
		_, n := d.GetChange("availability_zone_reference")
		spec.AvailabilityZoneReference = validateRef(n.(map[string]interface{}))
		hotPlugChange = false
	}
	spec.ClusterReference = response.Status.ClusterReference
	if d.HasChange("cluster_uuid") {
		_, n := d.GetChange("cluster_uuid")
		spec.ClusterReference = buildReference(n.(string), "cluster")
		hotPlugChange = false
	}
	if d.HasChange("parent_reference") {
		_, n := d.GetChange("parent_reference")
		res.ParentReference = validateRef(n.(map[string]interface{}))
		hotPlugChange = false
	}
	if d.HasChange("enable_cpu_passthrough") {
		_, n := d.GetChange("enable_cpu_passthrough")
		res.EnableCPUPassthrough = utils.BoolPtr(n.(bool))
		// TODO: Is this correct?
		hotPlugChange = false
	}
	if d.HasChange("num_vnuma_nodes") {
		_, n := d.GetChange("num_vnuma_nodes")
		res.VMVnumaConfig = &v3.VMVnumaConfig{
			NumVnumaNodes: utils.Int64Ptr(int64(n.(int))),
		}
		hotPlugChange = false
	}
	if d.HasChange("guest_os_id") {
		n := d.Get("guest_os_id")
		res.GuestOsID = utils.StringPtr(n.(string))
		hotPlugChange = false
	}
	if d.HasChange("num_vcpus_per_socket") {
		n := d.Get("num_vcpus_per_socket")
		res.NumVcpusPerSocket = utils.Int64Ptr(int64(n.(int)))
		hotPlugChange = false
	}
	if d.HasChange("num_sockets") {
		o, n := d.GetChange("num_sockets")
		res.NumSockets = utils.Int64Ptr(int64(n.(int)))

		//remove cpu sockets
		if n.(int) < o.(int) {
			hotPlugChange = false
		} else if !d.Get("use_hot_add").(bool) {
			hotPlugChange = false
		}
	}

	if d.HasChange("memory_size_mib") {
		o, n := d.GetChange("memory_size_mib")
		res.MemorySizeMib = utils.Int64Ptr(int64(n.(int)))
		//remove memory
		if n.(int) < o.(int) {
			hotPlugChange = false
		} else if !d.Get("use_hot_add").(bool) {
			hotPlugChange = false
		}
	}
	if d.HasChange("hardware_clock_timezone") {
		_, n := d.GetChange("hardware_clock_timezone")
		res.HardwareClockTimezone = utils.StringPtr(n.(string))
		hotPlugChange = false
	}
	if d.HasChange("vga_console_enabled") {
		_, n := d.GetChange("vga_console_enabled")
		res.VgaConsoleEnabled = utils.BoolPtr(n.(bool))
		hotPlugChange = false
	}
	if d.HasChange("guest_customization_is_overridable") {
		_, n := d.GetChange("guest_customization_is_overridable")
		guest.IsOverridable = utils.BoolPtr(n.(bool))
		hotPlugChange = false
	}
	if d.HasChange("power_state_mechanism") {
		_, n := d.GetChange("power_state_mechanism")
		pw.Mechanism = utils.StringPtr(n.(string))
		hotPlugChange = false
	}
	if d.HasChange("power_state_guest_transition_config") {
		_, n := d.GetChange("power_state_guest_transition_config")
		val := n.(map[string]interface{})

		p := &v3.VMGuestPowerStateTransitionConfig{}
		if v, ok := val["enable_script_exec"]; ok {
			p.EnableScriptExec = utils.BoolPtr(v.(bool))
		}
		if v, ok := val["should_fail_on_script_failure"]; ok {
			p.ShouldFailOnScriptFailure = utils.BoolPtr(v.(bool))
		}
		pw.GuestTransitionConfig = p
		hotPlugChange = false
	}

	cloudInit := guest.CloudInit

	if cloudInit == nil {
		cloudInit = &v3.GuestCustomizationCloudInit{}
	}

	if d.HasChange("guest_customization_cloud_init_user_data") {
		_, n := d.GetChange("guest_customization_cloud_init_user_data")
		cloudInit.UserData = utils.StringPtr(n.(string))
		hotPlugChange = false
	}

	if d.HasChange("guest_customization_cloud_init_meta_data") {
		_, n := d.GetChange("guest_customization_cloud_init_meta_data")
		cloudInit.MetaData = utils.StringPtr(n.(string))
		hotPlugChange = false
	}

	if d.HasChange("guest_customization_cloud_init_custom_key_values") {
		_, n := d.GetChange("guest_customization_cloud_init_custom_key_values")
		cloudInit.CustomKeyValues = n.(map[string]string)
		hotPlugChange = false
	}

	if !reflect.DeepEqual(*cloudInit, (v3.GuestCustomizationCloudInit{})) {
		guest.CloudInit = cloudInit
	}

	if d.HasChange("guest_customization_sysprep") {
		_, n := d.GetChange("guest_customization_sysprep")
		a := n.(map[string]interface{})

		guest.Sysprep = &v3.GuestCustomizationSysprep{
			InstallType: validateMapStringValue(a, "install_type"),
			UnattendXML: validateMapStringValue(a, "unattend_xml"),
		}
		hotPlugChange = false
	}
	if d.HasChange("guest_customization_sysprep_custom_key_values") {
		if guest.Sysprep == nil {
			guest.Sysprep = &v3.GuestCustomizationSysprep{}
		}
		_, n := d.GetChange("guest_customization_sysprep_custom_key_values")
		guest.Sysprep.CustomKeyValues = n.(map[string]string)
		hotPlugChange = false
	}
	if d.HasChange("nic_list") {
		res.NicList = expandNicList(d)
	}

	if d.HasChange("disk_list") {
		preCdromCount, err := CountDiskListCdrom(res.DiskList)
		if err != nil {
			return err
		}
		res.DiskList = expandDiskListUpdate(d, response)

		postCdromCount, err := CountDiskListCdrom(res.DiskList)
		if err != nil {
			return err
		}
		if preCdromCount != postCdromCount {
			hotPlugChange = false
		}
	}

	if d.HasChange("serial_port_list") {
		res.SerialPortList = expandSerialPortList(d)
		hotPlugChange = false
	}

	if d.HasChange("nutanix_guest_tools") || d.HasChange("ngt_credentials") || d.HasChange("ngt_enabled_capability_list") {
		res.GuestTools = expandNGT(d)
	}

	if d.HasChange("gpu_list") {
		res.GpuList = expandGPUList(d)
		hotPlugChange = false
	}

	if d.HasChange("machine_type") {
		n := d.Get("machine_type")
		res.MachineType = utils.StringPtr(n.(string))
		hotPlugChange = false
	}

	res.PowerStateMechanism = pw
	if bc, change := bootConfigHasChange(res.BootConfig, d); !reflect.DeepEqual(*bc, v3.VMBootConfig{}) {
		res.BootConfig = bc
		hotPlugChange = change
	}

	if !reflect.DeepEqual(*guestTool, (v3.GuestToolsSpec{})) {
		res.GuestTools = guestTool
	}

	if !reflect.DeepEqual(*guest, (v3.GuestCustomization{})) {
		res.GuestCustomization = guest
	}

	// If there are non-hotPlug changes, then poweroff is needed
	if !hotPlugChange {
		if err := changePowerState(conn, d.Id(), "OFF"); err != nil {
			return fmt.Errorf("internal error: cannot shut down the VM with UUID(%s): %s", d.Id(), err)
		}
		// SpecVersion has changed due previous poweroff
		specVersion, specErr := getVMSpecVersion(conn, d.Id())
		if specErr != nil {
			return fmt.Errorf("error getting spec for Virtual Machine UUID(%s): %s", d.Id(), specErr)
		}
		metadata.SpecVersion = specVersion
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	log.Printf("[DEBUG] Updating Virtual Machine: %s, %s", d.Get("name").(string), d.Id())

	resp, err2 := conn.V3.UpdateVM(d.Id(), request)
	if err2 != nil {
		return fmt.Errorf("error updating Virtual Machine UUID(%s): %s", d.Id(), err2)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for vm (%s) to update: %s", d.Id(), err)
	}

	// Then, Turn On the VM.
	if err := changePowerState(conn, d.Id(), "ON"); err != nil {
		return fmt.Errorf("internal error: cannot turn ON the VM with UUID(%s): %s", d.Id(), err)
	}

	return resourceNutanixVirtualMachineRead(d, meta)
}

func getVMSpecVersion(conn *v3.Client, vmID string) (*int64, error) {
	response, err := conn.V3.GetVM(vmID)
	if err != nil {
		return nil, err
	}
	if response.Metadata == nil {
		return nil, fmt.Errorf("failed to retrieve metadata for vm with uuid %s", vmID)
	}
	metadata := response.Metadata
	return metadata.SpecVersion, nil
}

func bootConfigHasChange(boot *v3.VMBootConfig, d *schema.ResourceData) (*v3.VMBootConfig, bool) {
	hotPlugChange := false

	if boot == nil {
		boot = &v3.VMBootConfig{}
	}

	if d.HasChange("boot_device_order_list") {
		_, n := d.GetChange("boot_device_order_list")
		boot.BootDeviceOrderList = expandStringList(n.([]interface{}))
		hotPlugChange = false
	}

	if d.HasChange("boot_type") {
		_, n := d.GetChange("boot_type")
		boot.BootType = utils.StringPtr(n.(string))
		hotPlugChange = false
	}

	bd := &v3.VMBootDevice{}
	dska := &v3.DiskAddress{}

	if d.HasChange("boot_device_disk_address") {
		_, n := d.GetChange("boot_device_disk_address")
		dai := n.(map[string]interface{})
		dska = &v3.DiskAddress{
			DeviceIndex: validateMapIntValue(dai, "device_index"),
			AdapterType: validateMapStringValue(dai, "adapter_type"),
		}
		hotPlugChange = false
	}
	if d.HasChange("boot_device_mac_address") {
		_, n := d.GetChange("boot_device_mac_address")
		bd.MacAddress = utils.StringPtr(n.(string))
		hotPlugChange = false
	}
	boot.BootDevice = bd

	if dska.AdapterType == nil && dska.DeviceIndex == nil && bd.MacAddress == nil {
		boot.BootDevice = nil
	}

	return boot, hotPlugChange
}

func changePowerState(conn *v3.Client, id string, powerState string) error {
	request := &v3.VMIntentInput{}
	metadata := &v3.Metadata{}
	res := &v3.VMResources{}
	spec := &v3.VM{}
	guest := &v3.GuestCustomization{}
	guestTool := &v3.GuestToolsSpec{}
	boot := &v3.VMBootConfig{}
	pw := &v3.VMPowerStateMechanism{}

	response, err := conn.V3.GetVM(id)
	preFillResUpdateRequest(res, response)
	preFillGTUpdateRequest(guestTool, response)
	preFillGUpdateRequest(guest, response)
	preFillPWUpdateRequest(pw, response)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			return nil
		}
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if !reflect.DeepEqual(*guestTool, (v3.GuestToolsSpec{})) {
		res.GuestTools = guestTool
	}

	if !reflect.DeepEqual(*guest, (v3.GuestCustomization{})) {
		res.GuestCustomization = guest
	}

	if !reflect.DeepEqual(*boot, (v3.VMBootConfig{})) {
		res.BootConfig = boot
	}

	spec.Name = response.Status.Name
	spec.Description = response.Status.Description
	spec.AvailabilityZoneReference = response.Status.AvailabilityZoneReference
	spec.ClusterReference = response.Status.ClusterReference

	res.PowerStateMechanism = pw
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	// Set PowerState OFF
	request.Spec.Resources.PowerState = utils.StringPtr(powerState)

	resp, err2 := conn.V3.UpdateVM(id, request)
	if err2 != nil {
		return fmt.Errorf("error updating Virtual Machine UUID(%s): %s", id, err2)
	}

	// Check update tasks
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for vm (%s) to update: %s", id, err)
	}

	// Check Power State
	stateConfVM := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    taskVMStateRefreshFunc(conn, id, powerState),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, err := stateConfVM.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for vm (%s) to update: %s", id, err)
	}
	return nil
}

func taskVMStateRefreshFunc(client *v3.Client, vmUUID string, powerState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetVM(vmUUID)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return v, DELETED, nil
			}
			return nil, ERROR, err
		}

		if *v.Status.State == "COMPLETE" && *v.Status.Resources.PowerState == powerState {
			return v, *v.Status.State, nil
		}
		return v, "RUNNING", nil
	}
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	setVMTimeout(meta)
	log.Printf("[DEBUG] Deleting Virtual Machine: %s, %s", d.Get("name").(string), d.Id())
	resp, err := conn.V3.DeleteVM(d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error while deleting Virtual Machine UUID(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for vm (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixVirtualMachineExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*Client).API

	_, err := conn.V3.GetVM(d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("error checking virtual Virtual Machine %s existence: %s", d.Id(), err)
	}
	return true, nil
}

func getVMResources(d *schema.ResourceData, vm *v3.VMResources) error {
	vm.PowerState = utils.StringPtr("ON")

	if v, ok := d.GetOk("num_vnuma_nodes"); ok {
		vm.VMVnumaConfig.NumVnumaNodes = utils.Int64Ptr(v.(int64))
	}

	if v, ok := d.GetOk("guest_os_id"); ok {
		vm.GuestOsID = utils.StringPtr(v.(string))
	}

	vm.NicList = expandNicList(d)
	vm.GpuList = expandGPUList(d)
	vm.GuestTools = expandNGT(d)

	if v, ok := d.GetOk("num_vcpus_per_socket"); ok {
		vm.NumVcpusPerSocket = utils.Int64Ptr(int64(v.(int)))
	}
	if v, ok := d.GetOk("enable_cpu_passthrough"); ok {
		vm.EnableCPUPassthrough = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("num_sockets"); ok {
		vm.NumSockets = utils.Int64Ptr(int64(v.(int)))
	}

	if v, ok := d.GetOk("parent_reference"); ok {
		val := v.(map[string]interface{})
		vm.ParentReference = validateRef(val)
	}

	if v, ok := d.GetOk("memory_size_mib"); ok {
		vm.MemorySizeMib = utils.Int64Ptr(int64(v.(int)))
	}

	vm.BootConfig = &v3.VMBootConfig{}

	if v, ok := d.GetOk("boot_device_order_list"); ok {
		vm.BootConfig.BootDeviceOrderList = expandStringList(v.([]interface{}))
	}

	bd := &v3.VMBootDevice{}
	da := &v3.DiskAddress{}
	if v, ok := d.GetOk("boot_device_disk_address"); ok {
		dai := v.(map[string]interface{})

		if value3, ok3 := dai["device_index"]; ok3 {
			if i, err := strconv.ParseInt(value3.(string), 10, 64); err == nil {
				da.DeviceIndex = utils.Int64Ptr(i)
			}
		}
		if value3, ok3 := dai["adapter_type"]; ok3 {
			da.AdapterType = utils.StringPtr(value3.(string))
		}
		bd.DiskAddress = da
		vm.BootConfig.BootDevice = bd
	}

	if v, ok := d.GetOk("boot_device_mac_address"); ok {
		bdi := v.(string)
		bd.MacAddress = utils.StringPtr(bdi)
		vm.BootConfig.BootDevice = bd
	}

	if v, ok := d.GetOk("boot_type"); ok {
		biosType := v.(string)
		vm.BootConfig.BootType = utils.StringPtr(biosType)
	}

	if v, ok := d.GetOk("machine_type"); ok {
		mtype := v.(string)
		vm.MachineType = utils.StringPtr(mtype)
	}

	if v, ok := d.GetOk("hardware_clock_timezone"); ok {
		vm.HardwareClockTimezone = utils.StringPtr(v.(string))
	}

	guestCustom := &v3.GuestCustomization{}
	cloudInit := &v3.GuestCustomizationCloudInit{}

	if v, ok := d.GetOk("guest_customization_cloud_init_user_data"); ok {
		cloudInit.UserData = utils.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("guest_customization_cloud_init_meta_data"); ok {
		cloudInit.MetaData = utils.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("guest_customization_cloud_init_custom_key_values"); ok {
		cloudInit.CustomKeyValues = utils.ConvertMapString(v.(map[string]interface{}))
	}

	if !reflect.DeepEqual(*cloudInit, (v3.GuestCustomizationCloudInit{})) {
		guestCustom.CloudInit = cloudInit
	}

	if v, ok := d.GetOk("guest_customization_is_overridable"); ok {
		guestCustom.IsOverridable = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("guest_customization_sysprep"); ok {
		guestCustom.Sysprep = &v3.GuestCustomizationSysprep{}
		spi := v.(map[string]interface{})
		if v2, ok2 := spi["install_type"]; ok2 {
			guestCustom.Sysprep.InstallType = utils.StringPtr(v2.(string))
		}
		if v2, ok2 := spi["unattend_xml"]; ok2 {
			guestCustom.Sysprep.UnattendXML = utils.StringPtr(v2.(string))
		}
	}

	if v, ok := d.GetOk("guest_customization_sysprep_custom_key_values"); ok {
		if guestCustom.Sysprep == nil {
			guestCustom.Sysprep = &v3.GuestCustomizationSysprep{}
		}
		guestCustom.Sysprep.CustomKeyValues = v.(map[string]string)
	}

	if !reflect.DeepEqual(*guestCustom, (v3.GuestCustomization{})) {
		vm.GuestCustomization = guestCustom
	}

	if v, ok := d.GetOk("vga_console_enabled"); ok {
		vm.VgaConsoleEnabled = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("power_state_mechanism"); ok {
		if vm.PowerStateMechanism == nil {
			log.Printf("m.PowerStateMechanism was nil, setting correct value!")
			vm.PowerStateMechanism = &v3.VMPowerStateMechanism{}
		}
		vm.PowerStateMechanism.Mechanism = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("should_fail_on_script_failure"); ok {
		vm.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("enable_script_exec"); ok {
		vm.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec = utils.BoolPtr(v.(bool))
	}
	vm.SerialPortList = expandSerialPortList(d)

	vmDiskList := expandDiskList(d)

	vm.DiskList = vmDiskList

	//check if BootConfig was set
	if reflect.DeepEqual(*vm.BootConfig, v3.VMBootConfig{}) {
		vm.BootConfig = nil
	}

	return nil
}

func expandNicList(d *schema.ResourceData) []*v3.VMNic {
	if v, ok := d.GetOk("nic_list"); ok {
		n := v.([]interface{})
		if len(n) > 0 {
			nics := make([]*v3.VMNic, 0)
			for _, nc := range n {
				val := nc.(map[string]interface{})
				nic := &v3.VMNic{}

				if value, ok := val["nic_type"]; ok && value.(string) != "" {
					nic.NicType = utils.StringPtr(value.(string))
				}
				if value, ok := val["uuid"]; ok && value.(string) != "" {
					nic.UUID = utils.StringPtr(value.(string))
				}
				if value, ok := val["network_function_nic_type"]; ok && value.(string) != "" {
					nic.NetworkFunctionNicType = utils.StringPtr(value.(string))
				}
				if value, ok := val["mac_address"]; ok && value.(string) != "" {
					nic.MacAddress = utils.StringPtr(value.(string))
				}
				if value, ok := val["model"]; ok && value.(string) != "" {
					nic.Model = utils.StringPtr(value.(string))
				}
				if value, ok := val["ip_endpoint_list"]; ok {
					nic.IPEndpointList = expandIPAddressList(value.([]interface{}))
				}
				if value, ok := val["network_function_chain_reference"]; ok && len(value.(map[string]interface{})) != 0 {
					v := value.(map[string]interface{})
					nic.NetworkFunctionChainReference = validateRef(v)
				}
				if value, ok := val["subnet_uuid"]; ok {
					v := value.(string)
					nic.SubnetReference = buildReference(v, "subnet")
				}
				if value, ok := val["is_connected"]; ok {
					v := value.(string)
					IsConnected, _ := strconv.ParseBool(v)
					nic.IsConnected = utils.BoolPtr(IsConnected)
				}

				nics = append(nics, nic)
			}
			return nics
		}
	}
	return nil
}

func expandIPAddressList(ipl []interface{}) []*v3.IPAddress {
	if len(ipl) > 0 {
		ip := make([]*v3.IPAddress, len(ipl))
		for k, i := range ipl {
			v := i.(map[string]interface{})
			v3ip := &v3.IPAddress{}

			if ipset, ipsetok := v["ip"]; ipsetok {
				v3ip.IP = utils.StringPtr(ipset.(string))
			}
			if iptype, iptypeok := v["type"]; iptypeok {
				v3ip.Type = utils.StringPtr(iptype.(string))
			}
			ip[k] = v3ip
		}
		return ip
	}
	return nil
}

func expandDiskListUpdate(d *schema.ResourceData, vm *v3.VMIntentResponse) []*v3.VMDisk {
	eDiskList := expandDiskList(d)
	if cloudInitCdromUUIDInt, ok := d.GetOk("cloud_init_cdrom_uuid"); ok {
		cloudInitCdromUUID := cloudInitCdromUUIDInt.(string)
		if cloudInitCdromUUID != "" && vm.Spec != nil && vm.Spec.Resources != nil {
			for _, disk := range vm.Spec.Resources.DiskList {
				if disk.UUID != nil && *disk.UUID == cloudInitCdromUUID {
					eDiskList = append(eDiskList, disk)
				}
			}
		}
	}
	return eDiskList
}

func expandDiskList(d *schema.ResourceData) []*v3.VMDisk {
	if v, ok := d.GetOk("disk_list"); ok {
		dsk := v.([]interface{})
		if len(dsk) > 0 {
			dls := make([]*v3.VMDisk, len(dsk))

			for k, val := range dsk {
				v := val.(map[string]interface{})
				dl := &v3.VMDisk{}

				// uuid
				if v1, ok1 := v["uuid"]; ok1 && v1.(string) != "" {
					dl.UUID = utils.StringPtr(v1.(string))
				}
				// storage_config
				if v, ok1 := v["storage_config"]; ok1 {
					dl.StorageConfig = expandStorageConfig(v.([]interface{}))
				}
				// device_properties
				if v1, ok1 := v["device_properties"]; ok1 {
					dl.DeviceProperties = expandDeviceProperties(v1.([]interface{}))
				}
				// data_source_reference
				if v1, ok := v["data_source_reference"]; ok && len(v1.(map[string]interface{})) != 0 {
					dsref := v1.(map[string]interface{})
					dl.DataSourceReference = validateShortRef(dsref)
				}
				// volume_group_reference
				if v1, ok := v["volume_group_reference"]; ok {
					volgr := v1.(map[string]interface{})
					dl.VolumeGroupReference = validateRef(volgr)
				}
				// disk_size_bytes
				if v1, ok1 := v["disk_size_bytes"]; ok1 && v1.(int) != 0 {
					dl.DiskSizeBytes = utils.Int64Ptr(int64(v1.(int)))
				}
				// disk_size_mib
				if v1, ok := v["disk_size_mib"]; ok && v1.(int) != 0 {
					dl.DiskSizeMib = utils.Int64Ptr(int64(v1.(int)))
				}
				dls[k] = dl
			}
			return dls
		}
	}
	return nil
}

func expandStorageConfig(storageConfig []interface{}) *v3.VMStorageConfig {
	if len(storageConfig) > 0 {
		v := storageConfig[0].(map[string]interface{})
		scr := v["storage_container_reference"].([]interface{})[0].(map[string]interface{})

		return &v3.VMStorageConfig{
			FlashMode: cast.ToString(v["flash_mode"]),
			StorageContainerReference: &v3.StorageContainerReference{
				URL:  cast.ToString(scr["url"]),
				Kind: cast.ToString(scr["kind"]),
				UUID: cast.ToString(scr["uuid"]),
			},
		}
	}
	return nil
}

func expandDeviceProperties(deviceProperties []interface{}) *v3.VMDiskDeviceProperties {
	if len(deviceProperties) > 0 {
		dp := &v3.VMDiskDeviceProperties{}
		d := deviceProperties[0].(map[string]interface{})

		if v, ok := d["device_type"]; ok {
			dp.DeviceType = utils.StringPtr(v.(string))
		}
		if v, ok := d["disk_address"]; ok && len(v.(map[string]interface{})) > 0 {
			da := v.(map[string]interface{})
			v3disk := &v3.DiskAddress{}

			if di, diOk := da["device_index"]; diOk {
				v3disk.DeviceIndex = utils.Int64Ptr(cast.ToInt64(di))
			}
			if at, atOk := da["adapter_type"]; atOk {
				v3disk.AdapterType = utils.StringPtr(at.(string))
			}
			dp.DiskAddress = v3disk
		}
		return dp
	}
	return nil
}

func expandSerialPortList(d *schema.ResourceData) []*v3.VMSerialPort {
	if v, ok := d.GetOk("serial_port_list"); ok {
		spl := v.([]interface{})

		if len(spl) > 0 {
			serialPortList := make([]*v3.VMSerialPort, len(spl))
			for k, val := range spl {
				v1 := val.(map[string]interface{})
				serialPort := &v3.VMSerialPort{}
				if v1, ok1 := v1["index"]; ok1 {
					serialPort.Index = utils.Int64Ptr(int64(v1.(int)))
				}
				if v1, ok1 := v1["is_connected"]; ok1 {
					serialPort.IsConnected = utils.BoolPtr(v1.(bool))
				}
				serialPortList[k] = serialPort
			}
			return serialPortList
		}
	}
	return nil
}

func expandGPUList(d *schema.ResourceData) []*v3.VMGpu {
	if v, ok := d.GetOk("gpu_list"); ok {
		if len(v.([]interface{})) > 0 {
			gpl := make([]*v3.VMGpu, len(v.([]interface{})))

			for k, va := range v.([]interface{}) {
				val := va.(map[string]interface{})
				gpu := &v3.VMGpu{}
				if value, ok1 := val["vendor"]; ok1 {
					gpu.Vendor = utils.StringPtr(value.(string))
				}
				if value, ok1 := val["device_id"]; ok1 {
					gpu.DeviceID = utils.Int64Ptr(int64(value.(int)))
				}
				if value, ok1 := val["mode"]; ok1 {
					gpu.Mode = utils.StringPtr(value.(string))
				}
				gpl[k] = gpu
			}
			return gpl
		}
	}
	return nil
}

func expandNGT(d *schema.ResourceData) *v3.GuestToolsSpec {
	guestTools := &v3.GuestToolsSpec{
		NutanixGuestTools: &v3.NutanixGuestToolsSpec{},
	}

	if v, ok := d.GetOk("nutanix_guest_tools"); ok {
		ngt := v.(map[string]interface{})
		if val, ok2 := ngt["state"]; ok2 {
			guestTools.NutanixGuestTools.State = utils.StringPtr(val.(string))
		}
		if val, ok2 := ngt["version"]; ok2 {
			guestTools.NutanixGuestTools.Version = utils.StringPtr(val.(string))
		}
		if val, ok2 := ngt["ngt_state"]; ok2 {
			guestTools.NutanixGuestTools.NgtState = utils.StringPtr(val.(string))
		}
		if val, ok2 := ngt["iso_mount_state"]; ok2 {
			guestTools.NutanixGuestTools.IsoMountState = utils.StringPtr(val.(string))
		}
	}

	if val, ok2 := d.GetOk("ngt_enabled_capability_list"); ok2 {
		guestTools.NutanixGuestTools.EnabledCapabilityList = expandStringList(val.([]interface{}))
	}

	if val, ok := d.GetOk("ngt_credentials"); ok {
		guestTools.NutanixGuestTools.Credentials = convertMapInterfaceToMapString(val.(map[string]interface{}))
	}

	if reflect.DeepEqual(guestTools.NutanixGuestTools, &v3.NutanixGuestToolsSpec{}) {
		return nil
	}

	return guestTools
}

func preFillResUpdateRequest(res *v3.VMResources, response *v3.VMIntentResponse) {
	res.GuestOsID = response.Spec.Resources.GuestOsID
	res.NumSockets = response.Spec.Resources.NumSockets
	res.PowerState = response.Spec.Resources.PowerState
	res.MemorySizeMib = response.Spec.Resources.MemorySizeMib
	res.VMVnumaConfig = &v3.VMVnumaConfig{NumVnumaNodes: response.Spec.Resources.VMVnumaConfig.NumVnumaNodes}
	res.ParentReference = response.Spec.Resources.ParentReference
	res.NumVcpusPerSocket = response.Spec.Resources.NumVcpusPerSocket
	res.VgaConsoleEnabled = response.Spec.Resources.VgaConsoleEnabled
	res.HardwareClockTimezone = response.Spec.Resources.HardwareClockTimezone
	res.DiskList = response.Spec.Resources.DiskList

	nold := make([]*v3.VMNic, len(response.Spec.Resources.NicList))

	if len(response.Spec.Resources.NicList) > 0 {
		for k, v := range response.Spec.Resources.NicList {
			nold[k] = &v3.VMNic{
				UUID:                          v.UUID,
				Model:                         v.Model,
				NicType:                       v.NicType,
				MacAddress:                    v.MacAddress,
				IPEndpointList:                v.IPEndpointList,
				SubnetReference:               v.SubnetReference,
				NetworkFunctionNicType:        v.NetworkFunctionNicType,
				NetworkFunctionChainReference: v.NetworkFunctionChainReference,
				IsConnected:                   v.IsConnected,
			}
		}
	} else {
		nold = nil
	}
	res.NicList = nold

	var spl []*v3.VMSerialPort
	if len(response.Spec.Resources.SerialPortList) > 0 {
		spl = make([]*v3.VMSerialPort, len(response.Spec.Resources.SerialPortList))
		for k, v := range response.Spec.Resources.SerialPortList {
			spl[k] = &v3.VMSerialPort{
				Index:       v.Index,
				IsConnected: v.IsConnected,
			}
		}
	}
	res.SerialPortList = spl

	gold := make([]*v3.VMGpu, len(response.Spec.Resources.GpuList))
	if len(response.Spec.Resources.GpuList) > 0 {
		for k, v := range response.Spec.Resources.GpuList {
			gold[k] = &v3.VMGpu{
				Mode:     v.Mode,
				Vendor:   v.Vendor,
				DeviceID: v.DeviceID,
			}
		}
	} else {
		gold = nil
	}
	res.GpuList = gold

	if response.Spec.Resources.BootConfig != nil {
		res.BootConfig = response.Spec.Resources.BootConfig
	} else {
		res.BootConfig = nil
	}
}

func preFillGTUpdateRequest(guestTool *v3.GuestToolsSpec, response *v3.VMIntentResponse) {
	if response.Spec.Resources.GuestTools != nil {
		guestTool.NutanixGuestTools = &v3.NutanixGuestToolsSpec{
			EnabledCapabilityList: response.Spec.Resources.GuestTools.NutanixGuestTools.EnabledCapabilityList,
			IsoMountState:         response.Spec.Resources.GuestTools.NutanixGuestTools.IsoMountState,
			State:                 response.Spec.Resources.GuestTools.NutanixGuestTools.State,
		}
	} else {
		guestTool = nil
	}
}

func preFillGUpdateRequest(guest *v3.GuestCustomization, response *v3.VMIntentResponse) {
	if response.Spec.Resources.GuestCustomization != nil {
		guest.CloudInit = response.Spec.Resources.GuestCustomization.CloudInit
		guest.Sysprep = response.Spec.Resources.GuestCustomization.Sysprep
		guest.IsOverridable = response.Spec.Resources.GuestCustomization.IsOverridable
	} else {
		guest = nil
	}
}

func preFillPWUpdateRequest(pw *v3.VMPowerStateMechanism, response *v3.VMIntentResponse) {
	if response.Spec.Resources.PowerStateMechanism != nil {
		pw.Mechanism = response.Spec.Resources.PowerStateMechanism.Mechanism
		pw.GuestTransitionConfig = response.Spec.Resources.PowerStateMechanism.GuestTransitionConfig
	} else {
		pw = nil
	}
}

func waitForIPRefreshFunc(client *v3.Client, vmUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.V3.GetVM(vmUUID)

		log.Printf("[DEBUG] GetVM Response %+v", resp)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return resp, ERROR, nil
			}
			return nil, "", err
		}

		if resp.Status == nil {
			return resp, WAITING, nil
		}

		if resp.Status.Resources == nil {
			return resp, WAITING, nil
		}

		if resp.Status.Resources.NicList != nil && len(resp.Status.Resources.NicList) != 0 {
			for _, v := range resp.Status.Resources.NicList {
				if len(v.IPEndpointList) > 0 {
					for _, v2 := range v.IPEndpointList {
						if v2.IP != nil {
							return resp, "AVAILABLE", nil
						}
					}
				}
			}
		}
		return resp, WAITING, nil
	}
}

func GetCdromDiskList(dl []*v3.VMDisk) []*v3.VMDisk {
	cdList := make([]*v3.VMDisk, 0)
	for _, v := range dl {
		if isCdromDisk(v) {
			cdList = append(cdList, v)
		}
	}
	return cdList
}

func isCdromDisk(d *v3.VMDisk) bool {
	if d.DeviceProperties != nil && *d.DeviceProperties.DeviceType == "CDROM" {
		return true
	}
	return false
}

func CountDiskListCdrom(dl []*v3.VMDisk) (int, error) {
	counter := len(GetCdromDiskList(dl))
	return counter, nil
}

func resourceVirtualMachineInstanceStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceVirtualMachineInstanceStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func setVMTimeout(meta interface{}) {
	client := meta.(*Client)
	if client.WaitTimeout != 0 {
		vmTimeout = time.Duration(client.WaitTimeout) * time.Minute
	}
}

func resourceNutanixVirtualMachineInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
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
			"categories": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(
						"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"),
					"please see http://developer.nutanix.com/reference/prism_central/v3/api/models/cluster-reference"),
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nic_list_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_endpoint_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// RESOURCES ARGUMENTS

			"enable_cpu_passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"num_vnuma_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"nic_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_endpoint_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"subnet_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "true",
						},
					},
				},
			},
			"guest_os_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nutanix_guest_tools": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ngt_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"iso_mount_state": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"available_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guest_os_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vss_snapshot_capable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_reachable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_mobility_drivers_installed": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ngt_credentials": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"ngt_enabled_capability_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_vcpus_per_socket": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"gpu_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frame_buffer_size_mib": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Optional: true,
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
						"pci_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fraction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"num_virtual_display_heads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"guest_driver_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"parent_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"memory_size_mib": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"boot_device_order_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"boot_device_disk_address": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_index": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"adapter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"boot_device_mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hardware_clock_timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guest_customization_cloud_init_user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"guest_customization_cloud_init_meta_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"guest_customization_cloud_init_custom_key_values": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"guest_customization_is_overridable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"guest_customization_sysprep": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"unattend_xml": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"guest_customization_sysprep_custom_key_values": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"should_fail_on_script_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_script_exec": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"power_state_mechanism": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vga_console_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disk_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_size_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"disk_size_mib": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"device_properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"disk_address": {
										Type:     schema.TypeMap,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"device_index": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"adapter_type": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"data_source_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"volume_group_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"serial_port_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
