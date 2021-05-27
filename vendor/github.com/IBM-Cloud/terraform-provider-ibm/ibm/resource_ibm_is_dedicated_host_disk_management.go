// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const ()

func resourceIBMISDedicatedHostDiskManagement() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisDedicatedHostDiskManagementCreate,
		Read:     resourceIBMisDedicatedHostDiskManagementRead,
		Update:   resourceIBMisDedicatedHostDiskManagementUpdate,
		Delete:   resourceIBMisDedicatedHostDiskManagementDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"dedicated_host": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the dedicated host for which disks has to be managed",
			},
			"disks": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Disk information that has to be updated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this disk.",
						},
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_is_dedicated_host_disk_management", "name"),
							Description:  "The user-defined name for this disk. The disk will be updated with this new name",
						},
					},
				},
			},
		},
	}
}

func resourceIBMISDedicatedHostDiskManagementValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmISDedicatedHostDiskManagementValidator := ResourceValidator{ResourceName: "ibm_is_dedicated_host_disk_management", Schema: validateSchema}
	return &ibmISDedicatedHostDiskManagementValidator
}

func resourceIBMisDedicatedHostDiskManagementCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	dedicatedhost := d.Get("dedicated_host").(string)
	disks := d.Get("disks")
	diskUpdate := disks.([]interface{})

	for _, disk := range diskUpdate {
		diskItem := disk.(map[string]interface{})
		namestr := diskItem["name"].(string)
		diskid := diskItem["id"].(string)

		updateDedicatedHostDiskOptions := &vpcv1.UpdateDedicatedHostDiskOptions{}
		updateDedicatedHostDiskOptions.SetDedicatedHostID(dedicatedhost)
		updateDedicatedHostDiskOptions.SetID(diskid)
		dedicatedHostDiskPatchModel := &vpcv1.DedicatedHostDiskPatch{
			Name: &namestr,
		}

		dedicatedHostDiskPatch, err := dedicatedHostDiskPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for DedicatedHostDiskPatch: %s", err)
		}
		updateDedicatedHostDiskOptions.SetDedicatedHostDiskPatch(dedicatedHostDiskPatch)

		_, _, err = sess.UpdateDedicatedHostDisk(updateDedicatedHostDiskOptions)
		if err != nil {
			return fmt.Errorf("Error calling UpdateDedicatedHostDisk: %s", err)
		}

	}
	d.SetId(dedicatedhost)
	return resourceIBMisDedicatedHostDiskManagementRead(d, meta)
}

func resourceIBMisDedicatedHostDiskManagementUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange("disks") && !d.IsNewResource() {

		disks := d.Get("disks")
		diskUpdate := disks.([]interface{})

		for _, disk := range diskUpdate {
			diskItem := disk.(map[string]interface{})
			namestr := diskItem["name"].(string)
			diskid := diskItem["id"].(string)
			updateDedicatedHostDiskOptions := &vpcv1.UpdateDedicatedHostDiskOptions{}
			updateDedicatedHostDiskOptions.SetDedicatedHostID(d.Id())
			updateDedicatedHostDiskOptions.SetID(diskid)
			dedicatedHostDiskPatchModel := &vpcv1.DedicatedHostDiskPatch{
				Name: &namestr,
			}

			dedicatedHostDiskPatch, err := dedicatedHostDiskPatchModel.AsPatch()
			if err != nil {
				return fmt.Errorf("Error calling asPatch for DedicatedHostDiskPatch: %s", err)
			}
			updateDedicatedHostDiskOptions.SetDedicatedHostDiskPatch(dedicatedHostDiskPatch)

			_, response, err := sess.UpdateDedicatedHostDisk(updateDedicatedHostDiskOptions)
			if err != nil {
				return fmt.Errorf("Error updating dedicated host disk: %s %s", err, response)
			}

		}

	}
	return resourceIBMisDedicatedHostDiskManagementRead(d, meta)
}

func resourceIBMisDedicatedHostDiskManagementDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")
	return nil
}

func resourceIBMisDedicatedHostDiskManagementRead(d *schema.ResourceData, meta interface{}) error {

	d.Set("dedicated_host", d.Id())

	return nil
}
