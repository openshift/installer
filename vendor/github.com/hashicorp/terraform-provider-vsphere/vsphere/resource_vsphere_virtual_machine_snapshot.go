package vsphere

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVSphereVirtualMachineSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereVirtualMachineSnapshotCreate,
		Read:   resourceVSphereVirtualMachineSnapshotRead,
		Delete: resourceVSphereVirtualMachineSnapshotDelete,

		Schema: map[string]*schema.Schema{
			"virtual_machine_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"memory": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"quiesce": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"remove_children": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"consolidate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVSphereVirtualMachineSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	vm, err := virtualmachine.FromUUID(client, d.Get("virtual_machine_uuid").(string))
	if err != nil {
		return fmt.Errorf("Error while getting the VirtualMachine :%s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout) // This is 5 mins
	defer cancel()
	task, err := vm.CreateSnapshot(ctx, d.Get("snapshot_name").(string), d.Get("description").(string), d.Get("memory").(bool), d.Get("quiesce").(bool))
	if err != nil {
		return fmt.Errorf("Error while creating snapshot: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	taskInfo, err := task.WaitForResult(tctx, nil)
	if err != nil {
		log.Printf("[DEBUG] Error While Creating the Task for Create Snapshot: %v", err)
		return fmt.Errorf(" Error While Creating the Task for Create Snapshot: %s", err)
	}
	log.Printf("[DEBUG] Task created for Create Snapshot: %v", task)
	if err != nil {
		log.Printf("[DEBUG] Error While waiting for the Task for Create Snapshot: %v", err)
		return fmt.Errorf(" Error While waiting for the Task for Create Snapshot: %s", err)
	}
	log.Printf("[DEBUG] Create Snapshot completed %v", d.Get("snapshot_name").(string))
	log.Println("[DEBUG] Managed Object Reference: " + taskInfo.Result.(types.ManagedObjectReference).Value)
	d.SetId(taskInfo.Result.(types.ManagedObjectReference).Value)
	return nil
}

func resourceVSphereVirtualMachineSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	vm, err := virtualmachine.FromUUID(client, d.Get("virtual_machine_uuid").(string))
	if err != nil {
		return fmt.Errorf("Error while getting the VirtualMachine :%s", err)
	}
	resourceVSphereVirtualMachineSnapshotRead(d, meta)
	if d.Id() == "" {
		log.Printf("[DEBUG] Error While finding the Snapshot: %v", err)
		return nil
	}
	log.Printf("[DEBUG] Deleting snapshot with name: %v", d.Get("snapshot_name").(string))
	var consolidatePtr *bool
	var removeChildren bool

	if v, ok := d.GetOk("consolidate"); ok {
		consolidate := v.(bool)
		consolidatePtr = &consolidate
	} else {

		consolidate := true
		consolidatePtr = &consolidate
	}
	if v, ok := d.GetOk("remove_children"); ok {
		removeChildren = v.(bool)
	} else {
		removeChildren = false
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout) // This is 5 mins
	defer cancel()
	task, err := vm.RemoveSnapshot(ctx, d.Id(), removeChildren, consolidatePtr)
	if err != nil {
		log.Printf("[DEBUG] Error While Creating the Task for Delete Snapshot: %v", err)
		return fmt.Errorf("Error While Creating the Task for Delete Snapshot: %s", err)
	}
	log.Printf("[DEBUG] Task created for Delete Snapshot: %v", task)

	err = task.Wait(ctx)
	if err != nil {
		log.Printf("[DEBUG] Error While waiting for the Task of Delete Snapshot: %v", err)
		return fmt.Errorf("Error While waiting for the Task of Delete Snapshot: %s", err)
	}
	log.Printf("[DEBUG] Delete Snapshot completed %v", d.Get("snapshot_name").(string))

	return nil
}

func resourceVSphereVirtualMachineSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	vm, err := virtualmachine.FromUUID(client, d.Get("virtual_machine_uuid").(string))
	if err != nil {
		return fmt.Errorf("Error while getting the VirtualMachine :%s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout) // This is 5 mins
	defer cancel()
	snapshot, err := vm.FindSnapshot(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "no snapshots for this VM") || strings.Contains(err.Error(), "snapshot \""+d.Get("snapshot_name").(string)+"\" not found") {
			log.Printf("[DEBUG] Error While finding the Snapshot: %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Error While finding the Snapshot: %v", err)
		return fmt.Errorf("Error while finding the Snapshot :%s", err)
	}
	log.Printf("[DEBUG] Snapshot found: %v", snapshot)
	return nil
}
