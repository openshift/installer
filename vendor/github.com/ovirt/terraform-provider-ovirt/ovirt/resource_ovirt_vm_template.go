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

//// BlankTemplateID indicates the ID of default blank template in oVirt
//const BlankTemplateID = "00000000-0000-0000-0000-000000000000"

func resourceOvirtTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtTemplateCreate,
		Read:   resourceOvirtTemplateRead,
		Update: resourceOvirtTemplateUpdate,
		Delete: resourceOvirtTemplateDelete,

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
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// we are creating a tempate from VM
			"vm_id": {
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
				ForceNew: true,
			},
			"memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
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
				ForceNew: true,
			},
			"sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: true,
			},
			"threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: true,
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
			"block_device": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
		},
	}
}

func resourceOvirtTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	//// template with disks attached is conflicted with block_device
	//templateID, templateIDOK := d.GetOk("template_id")
	//blockDevice, blockDeviceOk := d.GetOk("block_device")
	//
	//if !templateIDOK && !blockDeviceOk {
	//	return fmt.Errorf("one of template_id or block_device must be assigned")
	//}
	//
	//if templateIDOK {
	//	tds, err := getTemplateDiskAttachments(templateID.(string), meta)
	//	if err != nil {
	//		return err
	//	}
	//	if len(tds) > 0 && blockDeviceOk {
	//		return fmt.Errorf("template_id with disks attached is conflict with block_device")
	//	}
	//	if len(tds) == 0 && !blockDeviceOk {
	//		return fmt.Errorf("template has no disks attached, so block_device must be assigned")
	//	}
	//}

	builder := ovirtsdk4.NewTemplateBuilder().
		Name(d.Get("name").(string))

	if memory, ok := d.GetOk("memory"); ok {
		// memory is specified in MB
		builder.Memory(int64(memory.(int)) * int64(math.Pow(2, 20)))
	}

	cluster, err := ovirtsdk4.NewClusterBuilder().
		Id(d.Get("cluster_id").(string)).Build()
	if err != nil {
		return err
	}
	builder.Cluster(cluster)

	vm, err := ovirtsdk4.NewVmBuilder().Id(d.Get("vm_id").(string)).Build()
	if err != nil {
		return err
	}
	builder.Vm(vm)

	if ha, ok := d.GetOkExists("high_availability"); ok {
		highAvailability, err := ovirtsdk4.NewHighAvailabilityBuilder().
			Enabled(ha.(bool)).Build()

		if err != nil {
			return err
		}
		builder.HighAvailability(highAvailability)
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
	builder.Cpu(cpu)

	if v, ok := d.GetOk("initialization"); ok {
		initialization, err := expandOvirtVMInitialization(v.([]interface{}))
		if err != nil {
			return err
		}
		if initialization != nil {
			builder.Initialization(initialization)
		}
	}

	template, err := builder.Build()
	if err != nil {
		return err
	}

	// NOTE: the provider is creating a VM resource and expect it to be up,
	// a down status is not an option, but in order to create a template from
	// vm, the vm must be down. Instead of hacking the VM resource handling
	// this resource handling will take down the VM.
	vmService := conn.SystemService().VmsService().VmService(vm.MustId())
	vmResponse, err := vmService.Get().Send()
	if err != nil {
		return err
	}
	if vmResponse.MustVm().MustStatus() != ovirtsdk4.VMSTATUS_DOWN {
		_, err = vmService.Stop().Send()
		if err != nil {
			return err
		}
		downVmStateConf := &resource.StateChangeConf{
			Target:     []string{string(ovirtsdk4.VMSTATUS_DOWN)},
			Refresh:    VMStateRefreshFunc(conn, vm.MustId()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = downVmStateConf.WaitForState()
		if err != nil {
			log.Printf("[DEBUG] Failed to wait for VM(%s) to become down: %s", vm.MustId(), err)
			return err
		}
	}

	// now that the VM is down create the template
	resp, err := conn.SystemService().
		TemplatesService().
		Add().
		Template(template).
		Send()

	if err != nil {
		log.Printf("[DEBUG] Error creating the Template (%s)", d.Get("name").(string))
		return err
	}

	newTemplate, ok := resp.Template()
	if !ok {
		d.SetId("")
		return nil
	}
	d.SetId(newTemplate.MustId())

	log.Printf("[DEBUG] Template (%s) is created and wait for ready (status is ok)", d.Id())
	downStateConf := &resource.StateChangeConf{
		Target:     []string{string(ovirtsdk4.TEMPLATESTATUS_OK)},
		Refresh:    TemplateStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = downStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Failed to wait for Templte(%s) to become ok: %s", d.Id(), err)
		return err
	}
	log.Printf("[DEBUG] Newly created Template (%s) is ready", d.Id())
	//templateService := conn.SystemService().TemplatesService().TemplateService(d.Id())

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
	//if blockDeviceOk {
	//	log.Printf("[DEBUG] Attach disk specified by block_device to Template (%s)", d.Id())
	//	err = ovirtAttachDisks(blockDevice.([]interface{}), d.Id(), meta)
	//	if err != nil {
	//		return err
	//	}
	//}

	return resourceOvirtTemplateRead(d, meta)
}

func resourceOvirtTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	templateService := conn.SystemService().TemplatesService().TemplateService(d.Id())
	paramTemplate := ovirtsdk4.NewTemplateBuilder()
	attributeUpdated := false

	d.Partial(true)
	// initialization is a built-in attribute of VM that could be changed
	// at any conditions.
	if d.HasChange("initialization") {
		if v, ok := d.GetOk("initialization"); ok {
			initialization, err := expandOvirtVMInitialization(v.([]interface{}))
			if err != nil {
				return err
			}
			paramTemplate.Initialization(initialization)
		}
		attributeUpdated = true
	}

	if attributeUpdated {
		_, err := templateService.Update().Template(paramTemplate.MustBuild()).Send()
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceOvirtVMRead(d, meta)
}

func resourceOvirtTemplateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	response, err := conn.SystemService().TemplatesService().
		TemplateService(d.Id()).Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	template, ok := response.Template()

	if !ok {
		d.SetId("")
		return nil
	}
	d.Set("name", template.MustName())
	// memory is specified in MB
	d.Set("memory", template.MustMemory()/int64(math.Pow(2, 20)))
	d.Set("status", template.MustStatus())
	d.Set("cores", template.MustCpu().MustTopology().MustCores())
	d.Set("sockets", template.MustCpu().MustTopology().MustSockets())
	d.Set("threads", template.MustCpu().MustTopology().MustThreads())
	d.Set("cluster_id", template.MustCluster().MustId())

	if v, ok := template.Initialization(); ok {
		if err = d.Set("initialization", flattenOvirtVMInitialization(v)); err != nil {
			return fmt.Errorf("error setting initialization: %s", err)
		}
	}

	return nil
}

func resourceOvirtTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	templateService := conn.SystemService().TemplatesService().TemplateService(d.Id())

	response, err := templateService.Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		return fmt.Errorf("Error getting VM (%s) before deleting: %s", d.Id(), err)
	}

	template, ok := response.Template()
	if !ok {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Wait for Template (%s) status to become ok", d.Id())

	downStateConf := &resource.StateChangeConf{
		Target:     []string{string(ovirtsdk4.TEMPLATESTATUS_OK)},
		Refresh:    TemplateStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = downStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Failed to wait for Template (%s) to become down: %s", d.Id(), err)
		return fmt.Errorf("Error waiting for Template (%s) to be down: %s", d.Id(), err)
	}

	if template.MustDeleteProtected() {
		log.Printf("[DEBUG] Template (%s) is set as delete_protected and unset it first", d.Id())
		template.SetDeleteProtected(false)
		_, err := templateService.Update().
			Template(
				ovirtsdk4.NewTemplateBuilder().
					DeleteProtected(false).
					MustBuild()).
			Send()
		if err != nil {
			return fmt.Errorf("Error unsetting delete_protected for Template (%s): %s", d.Id(), err)
		}
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Now to remove Template (%s)", d.Id())
		_, err = templateService.Remove().Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// Wait until NotFoundError raises
				log.Printf("[DEBUG] Template (%s) has been removed", d.Id())
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error removing Template (%s): %s", template.MustId(), err))
		}
		return resource.RetryableError(fmt.Errorf("Template (%s) is still being removed", template.MustId()))
	})
}

// TemplateStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt Template.
func TemplateStateRefreshFunc(conn *ovirtsdk4.Connection, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := conn.SystemService().
			TemplatesService().
			TemplateService(id).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				return nil, "", nil
			}
			return nil, "", err
		}

		return r.MustTemplate(), string(r.MustTemplate().MustStatus()), nil
	}
}
