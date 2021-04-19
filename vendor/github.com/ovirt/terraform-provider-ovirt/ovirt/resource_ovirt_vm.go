// Copyright (C) 2017 Battelle Memorial Institute
// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// BlankTemplateID indicates the ID of default blank template in oVirt
const BlankTemplateID = "00000000-0000-0000-0000-000000000000"

func resourceOvirtVM(c *providerContext) *schema.Resource {
	return &schema.Resource{
		Create: c.resourceOvirtVMCreate,
		Read:   resourceOvirtVMRead,
		Update: resourceOvirtVMUpdate,
		Delete: resourceOvirtVMDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  BlankTemplateID,
			},
			"clone": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"lease_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if new memory is not set
					return new == "0"
				},
				Description: "in MB",
			},
			"maximum_memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if new memory is not set
					return new == "0"
				},
				Description: "in MB",
			},
			"cores": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"os": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"vnic_profile_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"boot_devices": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(ovirtsdk4.BOOTDEVICE_CDROM),
						string(ovirtsdk4.BOOTDEVICE_HD),
						string(ovirtsdk4.BOOTDEVICE_NETWORK),
					}, false),
				},
			},
			"block_device": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    false,
							Description: "in GiB",
						},
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"interface": {
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(ovirtsdk4.DISKINTERFACE_IDE),
								string(ovirtsdk4.DISKINTERFACE_SPAPR_VSCSI),
								string(ovirtsdk4.DISKINTERFACE_VIRTIO),
								string(ovirtsdk4.DISKINTERFACE_VIRTIO_SCSI),
							}, false),
							Required: true,
							ForceNew: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"storage_domain": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							Default:  "",
						},
						"logical_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"pass_discard": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"use_scsi_reservation": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
					},
				},
			},
			"auto_start": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"initialization": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_script": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_servers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_search": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nic_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label": {
										Type:     schema.TypeString,
										Required: true,
									},
									"boot_proto": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(ovirtsdk4.BOOTPROTOCOL_AUTOCONF),
											string(ovirtsdk4.BOOTPROTOCOL_DHCP),
											string(ovirtsdk4.BOOTPROTOCOL_NONE),
											string(ovirtsdk4.BOOTPROTOCOL_STATIC),
										}, false),
									},
									"address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"netmask": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gateway": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"on_boot": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"authorized_ssh_key": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: fmt.Sprintf(
					"One of %s, %s, %s",
					ovirtsdk4.VMTYPE_DESKTOP,
					ovirtsdk4.VMTYPE_SERVER,
					ovirtsdk4.VMTYPE_HIGH_PERFORMANCE),
				ValidateFunc: validation.StringInSlice([]string{
					string(ovirtsdk4.VMTYPE_DESKTOP),
					string(ovirtsdk4.VMTYPE_SERVER),
					string(ovirtsdk4.VMTYPE_HIGH_PERFORMANCE),
				}, false),
			},
			"instance_type_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: fmt.Sprintf(
					"The ID of the Instance Type." +
						" Checkout the IDs by requesting ovirt-engine/api/instancetypes" +
						" from APIs or the WebAdmin portal"),
			},
			"auto_pinning_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Description: fmt.Sprintf("The Auto Pinning Policy. One of %s, %s, %s",
					ovirtsdk4.AUTOPINNINGPOLICY_DISABLED,
					ovirtsdk4.AUTOPINNINGPOLICY_EXISTING,
					ovirtsdk4.AUTOPINNINGPOLICY_ADJUST),
				ValidateFunc: validation.StringInSlice([]string{
					string(ovirtsdk4.AUTOPINNINGPOLICY_DISABLED),
					string(ovirtsdk4.AUTOPINNINGPOLICY_EXISTING),
					string(ovirtsdk4.AUTOPINNINGPOLICY_ADJUST),
				}, false),
			},
			"affinity_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The name of the Affinity Groups that the VM will join",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hugepages": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The size of hugepage to use in KiB, One of 2048 or 1048576",
				ValidateFunc: validation.IntInSlice([]int{
					2048,
					1048576,
				}),
			},
		},
	}
}

func (c *providerContext) resourceOvirtVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	// template with disks attached is conflicted with block_device
	templateID, templateIDOK := d.GetOk("template_id")
	blockDevice, blockDeviceOk := d.GetOk("block_device")

	if !templateIDOK && !blockDeviceOk {
		return fmt.Errorf("one of template_id or block_device must be assigned")
	}

	if templateIDOK {
		tds, err := getTemplateDiskAttachments(templateID.(string), meta)
		if err != nil {
			return err
		}
		if len(tds) > 0 && blockDeviceOk {
			device := blockDevice.([]interface{})
			// check if a disk id was passed to block_device, fail if so.
			if diskID, _ := device[0].(map[string]interface{})["disk_id"].(string); diskID != "" {
				return fmt.Errorf("template_id with disks attached is conflict with block_device")
			}
		}
		if len(tds) == 0 && !blockDeviceOk {
			return fmt.Errorf("template has no disks attached, so block_device must be assigned")
		}
	}

	vmBuilder := ovirtsdk4.NewVmBuilder().
		Name(d.Get("name").(string))

	if memory, ok := d.GetOk("memory"); ok {
		// memory is specified in MB
		vmBuilder.Memory(int64(memory.(int)) * int64(math.Pow(2, 20)))
	}

	memoryPolicyBuilder := ovirtsdk4.NewMemoryPolicyBuilder()
	if maximumMemory, ok := d.GetOk("maximum_memory"); ok {
		// memory is specified in MB
		memoryPolicyBuilder.Max(int64(maximumMemory.(int)) * int64(math.Pow(2, 20)))
	}

	clusterId := d.Get("cluster_id").(string)
	cluster, err := ovirtsdk4.NewClusterBuilder().
		Id(clusterId).Build()
	if err != nil {
		return err
	}
	vmBuilder.Cluster(cluster)

	template, err := ovirtsdk4.NewTemplateBuilder().
		Id(templateID.(string)).Build()
	if err != nil {
		return err
	}
	vmBuilder.Template(template)

	if ha, ok := d.GetOkExists("high_availability"); ok {
		highAvailability, err := ovirtsdk4.NewHighAvailabilityBuilder().
			Enabled(ha.(bool)).Build()

		if err != nil {
			return err
		}
		vmBuilder.HighAvailability(highAvailability)

		if leaseDomain, ok := d.GetOkExists("lease_domain"); ok {
			sdsService := conn.SystemService().StorageDomainsService()
			sdsResp, err := sdsService.List().Search("name=" + leaseDomain.(string)).Send()
			if err != nil {
				return fmt.Errorf("failed to search storage domains, reason: %v", err)
			}
			storageDomains := sdsResp.MustStorageDomains()
			if len(storageDomains.Slice()) == 0 {
				return fmt.Errorf("failed to find storage domain with name=%s", leaseDomain.(string))
			}
			sd := storageDomains.Slice()[0]

			lease, err := ovirtsdk4.NewStorageDomainLeaseBuilder().StorageDomain(sd).Build()
			if err != nil {
				return err
			}

			vmBuilder.Lease(lease)
		}
	}

	cpuTopo := ovirtsdk4.NewCpuTopologyBuilder().
		Cores(int64(d.Get("cores").(int))).
		Threads(int64(d.Get("threads").(int))).
		Sockets(int64(d.Get("sockets").(int))).
		MustBuild()

	cpu, err := ovirtsdk4.NewCpuBuilder().
		Topology(cpuTopo).
		Build()
	if err != nil {
		return err
	}
	vmBuilder.Cpu(cpu)

	if blockDeviceOk {
		if storage_domain, _ := blockDevice.([]interface{})[0].(map[string]interface{})["storage_domain"]; storage_domain != "" && templateIDOK {

			// Get the reference to the service that manages the storage domains
			sdsService := conn.SystemService().StorageDomainsService()

			// Find the storage domain we want to be used for virtual machine disks
			sdsResp, err := sdsService.List().Search("name=" + storage_domain.(string)).Send()
			if err != nil {
				return fmt.Errorf("Failed to search storage domains, reason: %v", err)
			}

			tds, err := getTemplateDiskAttachments(templateID.(string), meta)
			if err != nil {
				return fmt.Errorf("Failed to get Template disks attachments: %v", err)
			}

			if storageDomains, ok := sdsResp.StorageDomains(); ok {
				if len(storageDomains.Slice()) == 0 {
					return fmt.Errorf("Failed to find storage domain with name=%s", storage_domain.(string))
				}
				sd := storageDomains.Slice()[0]

				for i, v := range tds {
					diskIndex := i + 1
					disk := v.MustDisk()
					disk.SetStorageDomain(sd)

					// Gett full information about disk
					diskService := conn.SystemService().DisksService().DiskService(disk.MustId())
					fullDiskInfo := diskService.Get().MustSend().MustDisk()
					diskFormat := fullDiskInfo.MustFormat()

					diskattachment := ovirtsdk4.NewDiskAttachmentBuilder().
						Disk(ovirtsdk4.NewDiskBuilder().
							Id(disk.MustId()).
							Format(diskFormat).
							StorageDomainsOfAny(
								ovirtsdk4.NewStorageDomainBuilder().
									Id(sd.MustId()).
									MustBuild()).
							MustBuild()).
						MustBuild()

					// Define basic disk aliases only if attribute alias defined
					if alias, _ := blockDevice.([]interface{})[0].(map[string]interface{})["alias"]; alias != "" {
						_, diskBotable := disk.Bootable()
						switch diskBotable {
						case true:
							disk.SetAlias(alias.(string))
						case false:
							disk.SetAlias(fmt.Sprintf("%s_Disk%v", d.Get("name").(string), diskIndex))
						}

						var newdiskattachment = diskattachment.MustDisk()
						newdiskattachment.SetAlias(disk.MustAlias())
					}
					vmBuilder.DiskAttachmentsOfAny(diskattachment)
				}
			}
		}
	}
	if cp, ok := d.GetOkExists("custom_properties"); ok {
		customProperties, err := expandOvirtCustomProperties(cp.([]interface{}))
		if err != nil {
			return err
		}
		if len(customProperties) > 0 {
			vmBuilder.CustomPropertiesOfAny(customProperties...)
		}
	}

	os, err := expandOS(d)
	if err != nil {
		return err
	}
	if os != nil {
		vmBuilder.Os(os)
	}

	if v, ok := d.GetOk("initialization"); ok {
		initialization, err := expandOvirtVMInitialization(v.([]interface{}))
		if err != nil {
			return err
		}
		if initialization != nil {
			vmBuilder.Initialization(initialization)
		}
	}

	isHighPerformance := false
	if v, ok := d.GetOk("type"); ok {
		vmType := ovirtsdk4.VmType(fmt.Sprint(v))
		vmBuilder.Type(vmType)
		if vmType == ovirtsdk4.VMTYPE_HIGH_PERFORMANCE {
			isHighPerformance = true

			// disable ballooning
			memoryPolicyBuilder.Ballooning(false)

			// set cpu host-passthrough
			cpu.SetMode(ovirtsdk4.CPUMODE_HOST_PASSTHROUGH)
			vmBuilder.Cpu(cpu)

			// enable serial console
			console, err := ovirtsdk4.NewConsoleBuilder().
				Enabled(true).Build()
			if err != nil {
				return err
			}
			vmBuilder.Console(console)
		}
	}

	isAutoPinning := false
	// TODO: remove the version check when everyone uses engine 4.4.5
	engineVer := ovirtGetEngineVersion(meta)
	versionCompareResult, err := versionCompare(engineVer, ovirtsdk4.NewVersionBuilder().
		Major(4).
		Minor(4).
		Build_(5).
		Revision(0).
		MustBuild())
	if err != nil {
		return err
	}
	if _, ok := d.GetOk("auto_pinning_policy"); ok && versionCompareResult < 0 {
		log.Printf("[WARN] The engine version %d.%d.%d is not supporting the auto pinning feature. "+
			"Please update to 4.4.5 or later.", engineVer.MustMajor(), engineVer.MustMinor(), engineVer.MustBuild())
	} else {
		// Mimic the UI behavior. unless specified set to existing policy for high performance VMs.
		if _, ok := d.GetOk("auto_pinning_policy"); !ok && isHighPerformance {
			err := d.Set("auto_pinning_policy", ovirtsdk4.AUTOPINNINGPOLICY_EXISTING)
			if err != nil {
				return err
			}
		}
		if v, ok := d.GetOk("auto_pinning_policy"); ok {
			isAutoPinning = true
			autoPinningPolicy := ovirtsdk4.AutoPinningPolicy(fmt.Sprint(v))

			// if we have a policy, we need to set the pinning to all the hosts in the cluster.
			if autoPinningPolicy != ovirtsdk4.AUTOPINNINGPOLICY_DISABLED {
				hostsInCluster, err := ovirtGetHostsInCluster(cluster, meta)
				if err != nil {
					return err
				}
				placementPolicyBuilder := ovirtsdk4.NewVmPlacementPolicyBuilder()
				placementPolicy, err := placementPolicyBuilder.Hosts(hostsInCluster).
					Affinity(ovirtsdk4.VMAFFINITY_MIGRATABLE).Build()
				if err != nil {
					return fmt.Errorf("failed to build the placement policy of the vm: %v", err)
				}
				vmBuilder.PlacementPolicy(placementPolicy)
			}
		}
	}

	if v, ok := d.GetOk("hugepages"); ok {
		customProp, err := ovirtsdk4.NewCustomPropertyBuilder().
			Name("hugepages").
			Value(fmt.Sprint(v)).
			Build()
		if err != nil {
			return err
		}
		vmBuilder.CustomPropertiesOfAny(customProp)
	}

	if v, ok := d.GetOk("instance_type_id"); ok {
		vmBuilder.InstanceTypeBuilder(
			ovirtsdk4.NewInstanceTypeBuilder().Id(v.(string)))
	}

	memoryPolicy, err := memoryPolicyBuilder.Build()
	if err != nil {
		return err
	}
	vmBuilder.MemoryPolicy(memoryPolicy)

	vm, err := vmBuilder.Build()
	if err != nil {
		return err
	}

	resp, err := conn.SystemService().
		VmsService().
		Add().
		Vm(vm).
		Clone(d.Get("clone").(bool)).
		Send()

	if err != nil {
		log.Printf("[DEBUG] Error creating the VM (%s)", d.Get("name").(string))
		return err
	}

	newVM, ok := resp.Vm()
	if !ok {
		d.SetId("")
		return nil
	}
	d.SetId(newVM.MustId())

	log.Printf("[DEBUG] VM (%s) is created and wait for ready (status is down)", d.Id())
	downStateConf := &resource.StateChangeConf{
		Target:     []string{string(ovirtsdk4.VMSTATUS_DOWN)},
		Refresh:    VMStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = downStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Failed to wait for VM (%s) to become down: %s", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Newly created VM (%s) is ready (status is down)", d.Id())
	vmService := conn.SystemService().VmsService().VmService(d.Id())

	// Do attach nics
	nics, nicsOk := d.GetOk("nics")
	if nicsOk {
		log.Printf("[DEBUG] Attach nics to VM (%s)", d.Id())
		err = ovirtAttachNics(nics.([]interface{}), d.Id(), meta)
		if err != nil {
			return err
		}
	}

	// remove graphic devices for high performance VM
	if isHighPerformance {
		log.Printf("[DEBUG] High Performance VM (%s) Removing Graphic consoles", d.Id())
		err := ovirtRemoveGraphicsConsoles(d.Id(), meta)
		if err != nil {
			log.Printf("[DEBUG] Error removing graphical devices")
		}
	}

	// update the VM with the auto pinning
	if isAutoPinning {
		autoPinningPolicy := ovirtsdk4.AutoPinningPolicy(fmt.Sprint(d.Get("auto_pinning_policy")))
		if autoPinningPolicy != ovirtsdk4.AUTOPINNINGPOLICY_DISABLED {
			log.Printf("[DEBUG] Setting Auto Pinning Policy to VM (%s).", d.Id())
			err := ovirtSetAutoPinningPolicy(d.Id(), autoPinningPolicy, meta)
			if err != nil {
				return fmt.Errorf("updating the VM (%s) with auto pinning policy failed! %v", d.Id(), err)
			}
		}
	}

	// Do attach disks
	if blockDeviceOk {
		log.Printf("[DEBUG] Attach disk specified by block_device to VM (%s)", d.Id())
		err = ovirtAttachDisks(blockDevice.([]interface{}), d.Id(), meta)
		if err != nil {
			return err
		}
	}

	affinityGroups, affinityGroupsOk := d.GetOk("affinity_groups")
	if affinityGroupsOk {
		agStr, err := convInterfaceArrToStringArr(affinityGroups.(*schema.Set).List())
		if err != nil {
			return err
		}
		ag, err := getAffinityGroups(conn, clusterId, agStr)
		if err != nil {
			return err
		}
		err = c.addVmToAffinityGroups(conn, newVM, clusterId, ag)
		if err != nil {
			return err
		}
	}

	autoStart := d.Get("auto_start").(bool)
	if autoStart {
		// Try to start VM
		log.Printf("[DEBUG] Try to start VM (%s)", d.Id())

		_, initialize := d.GetOk("initialization")
		_, err = vmService.Start().UseInitialization(initialize).Send()
		if err != nil {
			return err
		}
		// Wait until vm is up
		log.Printf("[DEBUG] Wait for VM (%s) status to become up", d.Id())

		upStateConf := &resource.StateChangeConf{
			Target:     []string{string(ovirtsdk4.VMSTATUS_UP)},
			Refresh:    VMStateRefreshFunc(conn, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = upStateConf.WaitForState()
		if err != nil {
			log.Printf("[DEBUG] Error waiting for VM (%s) to become up: %s", d.Id(), err)
			return err
		}

		log.Printf("[DEBUG] VM (%s) status has became to up", d.Id())
	}

	return resourceOvirtVMRead(d, meta)
}

func ovirtRemoveGraphicsConsoles(vmID string, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	vmGraphicConsoleService := conn.SystemService().VmsService().VmService(vmID).GraphicsConsolesService()

	graphics, err := vmGraphicConsoleService.List().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		return fmt.Errorf("error getting VM (%s) graphic consoles before deleting: %s", vmID, err)
	}

	for _, device := range graphics.MustConsoles().Slice() {
		_, err = vmGraphicConsoleService.
			ConsoleService(device.MustId()).
			Remove().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// Wait until NotFoundError raises
				return nil
			}
		}
	}
	return nil
}

func ovirtSetAutoPinningPolicy(vmID string, autoPinningPolicy ovirtsdk4.AutoPinningPolicy, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	vmService := conn.SystemService().VmsService().VmService(vmID)
	optimizeCpuSettings := !(autoPinningPolicy == ovirtsdk4.AUTOPINNINGPOLICY_EXISTING)
	_, err := vmService.AutoPinCpuAndNumaNodes().OptimizeCpuSettings(optimizeCpuSettings).Send()
	if err != nil {
		return fmt.Errorf("failed to set the auto pinning policy on the VM!, %v", err)
	}
	return nil
}

func ovirtGetHostsInCluster(cluster *ovirtsdk4.Cluster, meta interface{}) (*ovirtsdk4.HostSlice, error) {
	conn := meta.(*ovirtsdk4.Connection)
	clusterService := conn.SystemService().ClustersService().ClusterService(cluster.MustId())
	clusterGet, err := clusterService.Get().Send()
	if err != nil {
		return nil, fmt.Errorf("failed to get the cluster: %v", err)
	}
	clusterName := clusterGet.MustCluster().MustName()
	hostsInCluster, err := conn.SystemService().HostsService().List().Search(
		fmt.Sprintf("cluster=%s", clusterName)).Send()
	if err != nil {
		return nil, fmt.Errorf("failed to get the list of hosts in the cluster: %v", err)
	}
	return hostsInCluster.MustHosts(), nil
}

func ovirtGetEngineVersion(meta interface{}) *ovirtsdk4.Version {
	conn := meta.(*ovirtsdk4.Connection)
	engineVersion := conn.SystemService().Get().MustSend().MustApi().MustProductInfo().MustVersion()
	return engineVersion
}

func versionCompare(v *ovirtsdk4.Version, other *ovirtsdk4.Version) (int64, error) {
	if v == nil || other == nil {
		return 5, fmt.Errorf("can't compare nil objects")
	}
	if v == other {
		return 0, nil
	}
	result := v.MustMajor() - other.MustMajor()
	if result == 0 {
		result = v.MustMinor() - other.MustMinor()
		if result == 0 {
			result = v.MustBuild() - other.MustBuild()
			if result == 0 {
				result = v.MustRevision() - other.MustRevision()
			}
		}
	}
	return result, nil
}

func resourceOvirtVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	vmService := conn.SystemService().VmsService().VmService(d.Id())
	vmBuilder := ovirtsdk4.NewVmBuilder()
	attributeUpdated := false

	// Block that update VM Basic parameters:
	// Name, Memory, Cluster, CPU params
	if name, ok := d.GetOk("name"); ok {
		vmBuilder.Name(name.(string))
	}

	if memory, ok := d.GetOk("memory"); ok {
		// memory is specified in MB
		vmBuilder.Memory(int64(memory.(int)) * int64(math.Pow(2, 20)))
	}
	memoryPolicyBuilder := ovirtsdk4.NewMemoryPolicyBuilder()
	if maximumMemory, ok := d.GetOk("maximum_memory"); ok {
		// memory is specified in MB
		memoryPolicyBuilder.Max(int64(maximumMemory.(int)) * int64(math.Pow(2, 20)))
		memoryPolicy, err := memoryPolicyBuilder.Build()
		if err != nil {
			return err
		}
		vmBuilder.MemoryPolicy(memoryPolicy)
	}

	cluster, err := ovirtsdk4.NewClusterBuilder().
		Id(d.Get("cluster_id").(string)).Build()
	if err != nil {
		return err
	}
	vmBuilder.Cluster(cluster)

	if ha, ok := d.GetOkExists("high_availability"); ok {
		highAvailability, err := ovirtsdk4.NewHighAvailabilityBuilder().
			Enabled(ha.(bool)).Build()

		if err != nil {
			return err
		}
		vmBuilder.HighAvailability(highAvailability)

		if leaseDomain, ok := d.GetOkExists("lease_domain"); ok {
			sdsService := conn.SystemService().StorageDomainsService()
			sdsResp, err := sdsService.List().Search("name=" + leaseDomain.(string)).Send()
			if err != nil {
				return fmt.Errorf("failed to search storage domains, reason: %v", err)
			}
			storageDomains := sdsResp.MustStorageDomains()
			if len(storageDomains.Slice()) == 0 {
				return fmt.Errorf("failed to find storage domain with name=%s", leaseDomain.(string))
			}
			sd := storageDomains.Slice()[0]

			lease, err := ovirtsdk4.NewStorageDomainLeaseBuilder().StorageDomain(sd).Build()
			if err != nil {
				return err
			}

			vmBuilder.Lease(lease)
		}
	}

	cpuTopo := ovirtsdk4.NewCpuTopologyBuilder().
		Cores(int64(d.Get("cores").(int))).
		Threads(int64(d.Get("threads").(int))).
		Sockets(int64(d.Get("sockets").(int))).
		MustBuild()

	cpu, err := ovirtsdk4.NewCpuBuilder().
		Topology(cpuTopo).
		Build()
	if err != nil {
		return err
	}
	vmBuilder.Cpu(cpu)

	if cp, ok := d.GetOkExists("custom_properties"); ok {
		customProperties, err := expandOvirtCustomProperties(cp.([]interface{}))
		if err != nil {
			return err
		}
		if len(customProperties) > 0 {
			vmBuilder.CustomPropertiesOfAny(customProperties...)
		}
	}

	//paramVM.Initialization(initialization)
	if v, ok := d.GetOk("initialization"); ok {
		initialization, err := expandOvirtVMInitialization(v.([]interface{}))
		if err != nil {
			return err
		}
		if initialization != nil {
			vmBuilder.Initialization(initialization)
		}
	}

	if os_data, ok := d.GetOk("os"); ok {
		source := os_data.([]interface{})[0].(map[string]interface{})
		if v, ok := source["type"]; ok {
			os := ovirtsdk4.NewOperatingSystemBuilder().
				Type(v.(string)).
				MustBuild()
			vmBuilder.Os(os)
		}
	}

	_, err = vmService.Update().Vm(vmBuilder.MustBuild()).Send()
	if err != nil {
		log.Printf("[DEBUG] Error updating the VM (%s)", d.Get("name").(string))
		return err
	}

	// Check status and Start/Stop VM
	status, statusOK := d.GetOk("status")

	if d.HasChange("status") && statusOK {
		// Try to start VM
		log.Printf("[DEBUG] Try to update runing status for VM (%s)", d.Id())
		var vm_status ovirtsdk4.VmStatus

		switch status {
		case "up":
			vm_status = ovirtsdk4.VMSTATUS_UP
			_, err = vmService.Start().Send()
		case "down":
			vm_status = ovirtsdk4.VMSTATUS_DOWN
			_, err = vmService.Stop().Send()
		}
		if err != nil {
			log.Printf("[DEBUG] Failed to change status for VM (%s)", d.Id())
			return err
		}

		// Wait until vm is update status
		log.Printf("[DEBUG] Wait for VM (%s) status to become %s", d.Id(), vm_status)

		desiredStateConf := &resource.StateChangeConf{
			Target:     []string{string(vm_status)},
			Refresh:    VMStateRefreshFunc(conn, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = desiredStateConf.WaitForState()
		if err != nil {
			log.Printf("[DEBUG] Error waiting for VM (%s) to become %s: %s", d.Id(), vm_status, err)
			return err
		}

		log.Printf("[DEBUG] VM (%s) status has became to %s", d.Id(), vm_status)
	}

	// Update VM initialization parameters
	d.Partial(true)
	// initialization is a built-in attribute of VM that could be changed
	// at any conditions.
	if d.HasChange("initialization") {
		if v, ok := d.GetOk("initialization"); ok {
			initialization, err := expandOvirtVMInitialization(v.([]interface{}))
			if err != nil {
				return err
			}
			vmBuilder.Initialization(initialization)
		}
		attributeUpdated = true
	}
	if d.HasChange("instance_type_id") {
		if v, ok := d.GetOk("instance_type_id"); ok {
			vmBuilder.InstanceTypeBuilder(
				ovirtsdk4.NewInstanceTypeBuilder().Name(fmt.Sprint(v)))
		}
		attributeUpdated = true
	}

	if attributeUpdated {
		_, err := vmService.Update().Vm(vmBuilder.MustBuild()).Send()
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceOvirtVMRead(d, meta)
}

func resourceOvirtVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	getVmresp, err := conn.SystemService().VmsService().
		VmService(d.Id()).Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	vm, ok := getVmresp.Vm()

	if !ok {
		d.SetId("")
		return nil
	}
	d.Set("name", vm.MustName())
	// memory is specified in MB
	d.Set("memory", vm.MustMemory()/int64(math.Pow(2, 20)))
	d.Set("status", vm.MustStatus())
	d.Set("cores", vm.MustCpu().MustTopology().MustCores())
	d.Set("sockets", vm.MustCpu().MustTopology().MustSockets())
	d.Set("threads", vm.MustCpu().MustTopology().MustThreads())
	d.Set("cluster_id", vm.MustCluster().MustId())
	d.Set("maximum_memory", vm.MustMemoryPolicy().MustMax()/int64(math.Pow(2, 20)))

	if it, ok := vm.InstanceType(); ok {
		d.Set("instance_type_id", it.MustId())
	}

	err = d.Set("os", []map[string]interface{}{
		{"type": vm.MustOs().MustType()},
	})
	if err != nil {
		return fmt.Errorf("error setting os type: %s", err)
	}

	if len(d.Get("boot_devices").([]interface{})) != 0 {
		os, err := convertOS(vm.MustOs())
		if err != nil {
			return fmt.Errorf("error setting operating system: %s", err)
		}

		d.Set("boot_devices", os[0]["boot"].(map[string]interface{})["devices"])
	}

	// If the virtual machine is cloned from a template or another virtual machine,
	// the template links to the Blank template, and the original_template is used to track history.
	// Otherwise the template and original_template are the same.
	originalTemplate, originalTemplateOk := vm.OriginalTemplate()
	templateCloned := originalTemplateOk &&
		vm.MustTemplate().MustId() != originalTemplate.MustId() &&
		vm.MustTemplate().MustId() == BlankTemplateID
	if templateCloned {
		d.Set("template_id", originalTemplate.MustId())
	} else {
		d.Set("template_id", vm.MustTemplate().MustId())
	}
	d.Set("clone", templateCloned)

	if v, ok := vm.Initialization(); ok {
		if err = d.Set("initialization", flattenOvirtVMInitialization(v)); err != nil {
			return fmt.Errorf("error setting initialization: %s", err)
		}
	}

	return nil
}

func convertOS(os *ovirtsdk4.OperatingSystem) ([]map[string]interface{}, error) {
	boot := os.MustBoot()
	devices := boot.MustDevices()
	operatingSystems := make([]map[string]interface{}, 1)
	operatingSystem := make(map[string]interface{})
	operatingSystem["boot"] = make(map[string]interface{})
	outBoot := operatingSystem["boot"].(map[string]interface{})
	outBoot["devices"] = devices

	operatingSystems[0] = operatingSystem

	return operatingSystems, nil
}

func resourceOvirtVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	vmService := conn.SystemService().VmsService().VmService(d.Id())

	getVMResp, err := vmService.Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		return fmt.Errorf("Error getting VM (%s) before deleting: %s", d.Id(), err)
	}

	vm, ok := getVMResp.Vm()
	if !ok {
		d.SetId("")
		return nil
	}

	if vm.MustStatus() != ovirtsdk4.VMSTATUS_DOWN {
		log.Printf("[DEBUG] VM (%s) status is %s and now poweroff", d.Id(), vm.MustStatus())
		_, err := vmService.Stop().Send()
		if err != nil {
			return fmt.Errorf("Error powering off VM (%s) before deleting: %s", d.Id(), err)
		}
	}

	log.Printf("[DEBUG] Wait for VM (%s) status to become down", d.Id())

	downStateConf := &resource.StateChangeConf{
		Target:     []string{string(ovirtsdk4.VMSTATUS_DOWN)},
		Refresh:    VMStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = downStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Failed to wait for VM (%s) to become down: %s", d.Id(), err)
		return fmt.Errorf("Error waiting for VM (%s) to be down: %s", d.Id(), err)
	}

	if vm.MustDeleteProtected() {
		log.Printf("[DEBUG] VM (%s) is set as delete_protected and unset it first", d.Id())
		vm.SetDeleteProtected(false)
		_, err := vmService.Update().
			Vm(
				ovirtsdk4.NewVmBuilder().
					DeleteProtected(false).
					MustBuild()).
			Send()
		if err != nil {
			return fmt.Errorf("Error unsetting delete_protected for VM (%s): %s", d.Id(), err)
		}
	}

	// VM created by Template must be remove with detachOnly=false
	detachOnly := true
	log.Printf("[DEBUG] Determine the detachOnly flag before removing VM (%s)", d.Id())
	if vm.MustTemplate().MustId() != BlankTemplateID || d.Get("clone").(bool) {
		log.Printf("[DEBUG] Set detachOnly flag to false since VM (%s) is based on template (%s) or is a clone",
			d.Id(), vm.MustTemplate().MustId())
		detachOnly = false
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Now to remove VM (%s)", d.Id())
		_, err = vmService.Remove().
			DetachOnly(detachOnly).
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// Wait until NotFoundError raises
				log.Printf("[DEBUG] VM (%s) has been removed", d.Id())
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error removing VM (%s): %s", vm.MustTemplate().MustId(), err))
		}
		return resource.RetryableError(fmt.Errorf("VM (%s) is still being removed", vm.MustTemplate().MustId()))
	})
}

// VMStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt VM.
func VMStateRefreshFunc(conn *ovirtsdk4.Connection, vmID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := conn.SystemService().
			VmsService().
			VmService(vmID).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// Sometimes oVirt has consistency issues and doesn't see
				// newly created VM instance. Return an empty state.
				return nil, "", nil
			}
			return nil, "", err
		}

		return r.MustVm(), string(r.MustVm().MustStatus()), nil
	}
}

func expandOS(d *schema.ResourceData) (*ovirtsdk4.OperatingSystem, error) {
	osBuilder := ovirtsdk4.NewOperatingSystemBuilder()

	devicesExists := d.Get("boot_devices").([]interface{})
	if devicesExists != nil {
		devices, err := expandOvirtBootDevices(devicesExists)
		if err != nil {
			return nil, err
		}
		boot, err := ovirtsdk4.NewBootBuilder().
			Devices(devices).
			Build()
		if err != nil {
			return nil, err
		}

		osBuilder.Boot(boot)
	}

	v, ok := d.GetOk("os")
	if ok {
		source := v.([]interface{})[0].(map[string]interface{})
		if v, ok := source["type"]; ok {
			osBuilder.Type(v.(string))
		}
	}

	return osBuilder.Build()
}

func expandOvirtVMInitialization(l []interface{}) (*ovirtsdk4.Initialization, error) {
	if len(l) == 0 {
		return nil, nil
	}
	s := l[0].(map[string]interface{})
	initializationBuilder := ovirtsdk4.NewInitializationBuilder()
	if v, ok := s["host_name"]; ok {
		initializationBuilder.HostName(v.(string))
	}
	if v, ok := s["timezone"]; ok {
		initializationBuilder.Timezone(v.(string))
	}
	if v, ok := s["user_name"]; ok {
		initializationBuilder.UserName(v.(string))
	}
	if v, ok := s["custom_script"]; ok {
		initializationBuilder.CustomScript(v.(string))
	}
	if v, ok := s["authorized_ssh_key"]; ok {
		initializationBuilder.AuthorizedSshKeys(v.(string))
	}
	if v, ok := s["dns_servers"]; ok {
		initializationBuilder.DnsServers(v.(string))
	}
	if v, ok := s["dns_search"]; ok {
		initializationBuilder.DnsSearch(v.(string))
	}
	if v, ok := s["nic_configuration"]; ok {
		ncs, err := expandOvirtVMNicConfigurations(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		if len(ncs) > 0 {
			initializationBuilder.NicConfigurationsOfAny(ncs...)
		}
	}
	return initializationBuilder.Build()
}

func expandOvirtBootDevices(l []interface{}) ([]ovirtsdk4.BootDevice, error) {
	devices := make([]ovirtsdk4.BootDevice, len(l))
	for i, v := range l {
		devices[i] = ovirtsdk4.BootDevice(v.(string))
	}

	return devices, nil
}

func expandOvirtCustomProperties(l []interface{}) ([]*ovirtsdk4.CustomProperty, error) {
	customProperties := make([]*ovirtsdk4.CustomProperty, len(l))
	for i, v := range l {
		vmap := v.(map[string]interface{})
		customPropBuilder := ovirtsdk4.NewCustomPropertyBuilder()
		customPropBuilder.Name(vmap["name"].(string))
		customPropBuilder.Value(vmap["value"].(string))
		customProp, err := customPropBuilder.Build()
		if err != nil {
			return nil, err
		}

		customProperties[i] = customProp
	}
	return customProperties, nil
}

func expandOvirtVMNicConfigurations(l []interface{}) ([]*ovirtsdk4.NicConfiguration, error) {
	nicConfs := make([]*ovirtsdk4.NicConfiguration, len(l))
	for i, v := range l {
		vmap := v.(map[string]interface{})
		ncbuilder := ovirtsdk4.NewNicConfigurationBuilder()
		ncbuilder.Name(vmap["label"].(string))
		ncbuilder.BootProtocol(ovirtsdk4.BootProtocol(vmap["boot_proto"].(string)))
		if v, ok := vmap["on_boot"]; ok {
			ncbuilder.OnBoot(v.(bool))
		}
		address, addressOK := vmap["address"]
		netmask, netmaskOK := vmap["netmask"]
		gateway, gatewayOK := vmap["gateway"]
		if addressOK || netmaskOK || gatewayOK {
			ipBuilder := ovirtsdk4.NewIpBuilder()
			if addressOK {
				ipBuilder.Address(address.(string))
			}
			if netmaskOK {
				ipBuilder.Netmask(netmask.(string))
			}
			if gatewayOK {
				ipBuilder.Gateway(gateway.(string))
			}
			ncbuilder.IpBuilder(ipBuilder)
		}
		nc, err := ncbuilder.Build()
		if err != nil {
			return nil, err
		}
		nicConfs[i] = nc
	}
	return nicConfs, nil
}

func expandOvirtVMDiskAttachment(d interface{}, disk *ovirtsdk4.Disk) (*ovirtsdk4.DiskAttachment, error) {
	dmap := d.(map[string]interface{})
	builder := ovirtsdk4.NewDiskAttachmentBuilder()
	// block_device only support bootable disk
	builder.Bootable(true)
	if disk != nil {
		builder.Disk(disk)
		if v, ok := dmap["size"]; ok {
			if v != 0 {
				newSize := int64(v.(int)) * int64(math.Pow(2, 30))
				disk.SetProvisionedSize(newSize)
			}
		}
		if v, ok := dmap["alias"]; ok {
			if v != "" {
				disk.SetAlias(v.(string))
			}
		}
	}
	if v, ok := dmap["interface"]; ok {
		builder.Interface(ovirtsdk4.DiskInterface(v.(string)))
	}
	if v, ok := dmap["active"]; ok {
		builder.Active(v.(bool))
	}
	if v, ok := dmap["logical_name"]; ok {
		builder.LogicalName(v.(string))
	}
	if v, ok := dmap["pass_discard"]; ok {
		builder.PassDiscard(v.(bool))
	}
	if v, ok := dmap["read_only"]; ok {
		builder.ReadOnly(v.(bool))
	}
	if v, ok := dmap["use_scsi_reservation"]; ok {
		builder.UsesScsiReservation(v.(bool))
	}

	return builder.Build()
}

func ovirtAttachNics(n []interface{}, vmID string, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	vmService := conn.SystemService().VmsService().VmService(vmID)
	for _, v := range n {
		nic := v.(map[string]interface{})
		mac := &ovirtsdk4.Mac{}
		if len(nic["mac"].(string)) != 0 {
			mac.SetAddress(nic["mac"].(string))
		}
		resp, err := vmService.NicsService().Add().Nic(
			ovirtsdk4.NewNicBuilder().
				Name(nic["name"].(string)).
				Mac(mac).
				VnicProfile(
					ovirtsdk4.NewVnicProfileBuilder().
						Id(nic["vnic_profile_id"].(string)).
						MustBuild()).
				MustBuild()).Send()
		if err != nil {
			return err
		}
		_, ok := resp.Nic()
		if !ok {
			return fmt.Errorf("failed to add nic: response does not contain the nic")
		}
	}
	return nil
}

func ovirtAttachDisks(s []interface{}, vmID string, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	vmService := conn.SystemService().VmsService().VmService(vmID)
	for _, v := range s {
		blockDeviceElement, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed getting block_device content %v", blockDeviceElement)
		}
		// try the passed disk_id, if not, detect boot disk
		diskID, ok := blockDeviceElement["disk_id"].(string)
		attachmentExists := false
		if !ok || diskID == "" {
			findID, err := findVMBootDiskAttachmentID(vmService)
			if err != nil {
				return err
			}
			diskID = findID
			attachmentExists = true
		}
		diskService := conn.SystemService().DisksService().
			DiskService(diskID)
		var disk *ovirtsdk4.Disk
		err := resource.Retry(30*time.Second, func() *resource.RetryError {
			getDiskResp, err := diskService.Get().Send()
			if err != nil {
				return resource.RetryableError(err)
			}
			disk = getDiskResp.MustDisk()
			if disk.MustStatus() == ovirtsdk4.DISKSTATUS_LOCKED {
				return resource.RetryableError(fmt.Errorf("disk is locked, wait for next check"))
			}
			return nil
		})
		if err != nil {
			return err
		}

		da, err := expandOvirtVMDiskAttachment(v, disk)
		if err != nil {
			return err
		}

		attachment, err := attachDisk(vmService.DiskAttachmentsService(), da, attachmentExists)
		if err != nil {
			return fmt.Errorf("failed to attach disk: %s", err)
		}
		err = conn.WaitForDisk(attachment.MustId(), ovirtsdk4.DISKSTATUS_OK, 20*time.Minute)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// attachDisk will attach a disk to vm or update an existing one. Returns the
// new attachment or error.
func attachDisk(service *ovirtsdk4.DiskAttachmentsService, attachment *ovirtsdk4.DiskAttachment, update bool) (*ovirtsdk4.DiskAttachment, error) {
	if update {
		r, err := service.
			AttachmentService(attachment.MustDisk().MustId()).
			Update().
			DiskAttachment(attachment).
			Send()
		if err != nil {
			return nil, err
		}
		return r.MustDiskAttachment(), nil
	}
	r, err := service.
		Add().
		Attachment(attachment).
		Send()
	if err != nil {
		return nil, err
	}
	return r.MustAttachment(), nil
}

// findVMBootDiskAttachmentID returns the disk attachment id of
// the bootable disk of a VM
func findVMBootDiskAttachmentID(vmService *ovirtsdk4.VmService) (string, error) {
	r, err := vmService.DiskAttachmentsService().List().Send()
	if err != nil {
		return "", err
	}

	for _, attachment := range r.MustAttachments().Slice() {
		bootable, ok := attachment.Bootable()
		if ok && bootable {
			return attachment.MustId(), nil
		}
	}

	return "", fmt.Errorf("no bootable disk for the VM")
}

func flattenOvirtVMDiskAttachments(configured []*ovirtsdk4.DiskAttachment) []map[string]interface{} {
	diskAttachments := make([]map[string]interface{}, len(configured))
	for i, v := range configured {
		attrs := make(map[string]interface{})
		attrs["disk_attachment_id"] = v.MustId()
		attrs["disk_id"] = v.MustDisk().MustId()
		attrs["interface"] = v.MustInterface()

		if vi, ok := v.Active(); ok {
			attrs["active"] = vi
		}
		if vi, ok := v.Bootable(); ok {
			attrs["bootable"] = vi
		}
		if vi, ok := v.LogicalName(); ok {
			attrs["logical_name"] = vi
		}
		if vi, ok := v.PassDiscard(); ok {
			attrs["pass_discard"] = vi
		}
		if vi, ok := v.ReadOnly(); ok {
			attrs["read_only"] = vi
		}
		if vi, ok := v.UsesScsiReservation(); ok {
			attrs["use_scsi_reservation"] = vi
		}
		diskAttachments[i] = attrs
	}
	return diskAttachments
}

func flattenOvirtVMInitialization(configured *ovirtsdk4.Initialization) []map[string]interface{} {
	if configured == nil {
		initializations := make([]map[string]interface{}, 0)
		return initializations
	}
	initializations := make([]map[string]interface{}, 1)
	initialization := make(map[string]interface{})

	if v, ok := configured.HostName(); ok {
		initialization["host_name"] = v
	}
	if v, ok := configured.Timezone(); ok {
		initialization["timezone"] = v
	}
	if v, ok := configured.UserName(); ok {
		initialization["user_name"] = v
	}
	if v, ok := configured.CustomScript(); ok {
		initialization["custom_script"] = v
	}
	if v, ok := configured.DnsServers(); ok {
		initialization["dns_servers"] = v
	}
	if v, ok := configured.DnsSearch(); ok {
		initialization["dns_search"] = v
	}
	if v, ok := configured.AuthorizedSshKeys(); ok {
		initialization["authorized_ssh_key"] = v
	}
	if v, ok := configured.NicConfigurations(); ok {
		initialization["nic_configuration"] = flattenOvirtVMInitializationNicConfigurations(v.Slice())
	}
	initializations[0] = initialization
	return initializations
}

func flattenOvirtVMInitializationNicConfigurations(configured []*ovirtsdk4.NicConfiguration) []map[string]interface{} {
	ncs := make([]map[string]interface{}, len(configured))
	for i, v := range configured {
		attrs := make(map[string]interface{})
		if name, ok := v.Name(); ok {
			attrs["label"] = name
		}
		attrs["on_boot"] = v.MustOnBoot()
		attrs["boot_proto"] = v.MustBootProtocol()
		if ipAttrs, ok := v.Ip(); ok {
			if ipAddr, ok := ipAttrs.Address(); ok {
				attrs["address"] = ipAddr
			}
			if netmask, ok := ipAttrs.Netmask(); ok {
				attrs["netmask"] = netmask
			}
			if gateway, ok := ipAttrs.Gateway(); ok {
				attrs["gateway"] = gateway
			}
		}
		ncs[i] = attrs
	}
	return ncs
}

func getTemplateDiskAttachments(templateID string, meta interface{}) ([]*ovirtsdk4.DiskAttachment, error) {
	conn := meta.(*ovirtsdk4.Connection)
	getTemplateDiskResp, err := conn.SystemService().
		TemplatesService().
		TemplateService(templateID).
		DiskAttachmentsService().
		List().
		Send()
	if err != nil {
		return nil, err
	}
	if vs, ok := getTemplateDiskResp.Attachments(); ok {
		return vs.Slice(), nil
	}
	return nil, nil
}

func getAffinityGroups(conn *ovirtsdk4.Connection, cID string, agNames []string) (ag []*ovirtsdk4.AffinityGroup, err error) {
	var ags []*ovirtsdk4.AffinityGroup
	var notFoundAGs []string

	res, err := conn.SystemService().ClustersService().
		ClusterService(cID).AffinityGroupsService().
		List().Send()
	if err != nil {
		return nil, err
	}
	agNamesMap := make(map[string]*ovirtsdk4.AffinityGroup)
	for _, af := range res.MustGroups().Slice() {
		agNamesMap[af.MustName()] = af
	}
	for _, agName := range agNames {
		if _, ok := agNamesMap[agName]; !ok {
			notFoundAGs = append(notFoundAGs, agName)
		} else {
			ags = append(ags, agNamesMap[agName])
		}
	}
	if len(notFoundAGs) > 0 {
		return nil, fmt.Errorf("affinity groups %v were not found on cluster %s", notFoundAGs, cID)
	}
	return ags, nil
}

func (c *providerContext) addVmToAffinityGroups(conn *ovirtsdk4.Connection, vm *ovirtsdk4.Vm, cID string, ags []*ovirtsdk4.AffinityGroup) error {
	// TODO: Remove lock once BZ#1950767 is resolved
	c.semaphores.Lock("vm-ag", 1)
	defer c.semaphores.Unlock("vm-ag")

	for _, ag := range ags {
		log.Printf("Adding machine %s to affinity group %s", vm.MustName(), ag.MustName())
		_, err := conn.SystemService().ClustersService().
			ClusterService(cID).AffinityGroupsService().
			GroupService(ag.MustId()).VmsService().Add().Vm(vm).Send()
		// TODO: Remove error handling workaround when BZ#1931932 is resolved and backported
		if err != nil && !errors.Is(err, ovirtsdk4.XMLTagNotMatchError{ActualTag: "action", ExpectedTag: "vm"}) {
			return fmt.Errorf(
				"failed to add VM %s to AffinityGroup %s, error: %v",
				vm.MustName(),
				ag.MustName(),
				err)
		}
	}
	return nil
}
