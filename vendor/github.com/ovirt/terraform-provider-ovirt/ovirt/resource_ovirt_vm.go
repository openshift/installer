// Copyright (C) 2017 Battelle Memorial Institute
// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
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

func resourceOvirtVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtVMCreate,
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
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
		},
	}
}

func resourceOvirtVMCreate(d *schema.ResourceData, meta interface{}) error {
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

	cluster, err := ovirtsdk4.NewClusterBuilder().
		Id(d.Get("cluster_id").(string)).Build()
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

	if v, ok := d.GetOk("type"); ok {
		vmBuilder.Type(ovirtsdk4.VmType(fmt.Sprint(v)))
	}

	if v, ok := d.GetOk("instance_type_id"); ok {
		vmBuilder.InstanceTypeBuilder(
			ovirtsdk4.NewInstanceTypeBuilder().Id(v.(string)))
	}

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

	// Do attach disks
	if blockDeviceOk {
		log.Printf("[DEBUG] Attach disk specified by block_device to VM (%s)", d.Id())
		err = ovirtAttachDisks(blockDevice.([]interface{}), d.Id(), meta)
		if err != nil {
			return err
		}
	}

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
	return resourceOvirtVMRead(d, meta)
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
	if vm.MustTemplate().MustId() != BlankTemplateID {
		log.Printf("[DEBUG] Set detachOnly flag to false since VM (%s) is based on template (%s)",
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
			newSize := int64(v.(int)) * int64(math.Pow(2, 30))
			disk.SetProvisionedSize(newSize)
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
		resp, err := vmService.NicsService().Add().Nic(
			ovirtsdk4.NewNicBuilder().
				Name(nic["name"].(string)).
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
