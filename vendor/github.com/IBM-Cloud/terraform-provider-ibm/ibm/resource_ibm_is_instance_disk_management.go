// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const ()

func resourceIBMISInstanceDiskManagement() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisInstanceDiskManagementCreate,
		Read:     resourceIBMisInstanceDiskManagementRead,
		Update:   resourceIBMisInstanceDiskManagementUpdate,
		Delete:   resourceIBMisInstanceDiskManagementDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the instance for which disks has to be managed",
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
							Description: "The unique identifier for this instance disk.",
						},
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_is_instance_disk_management", "name"),
							Description:  "The user-defined name for this disk. The disk will be updated with this new name",
						},
					},
				},
			},
		},
	}
}

func resourceIBMISInstanceDiskManagementValidator() *ResourceValidator {

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

	ibmISInstanceDiskManagementValidator := ResourceValidator{ResourceName: "ibm_is_instance_disk_management", Schema: validateSchema}
	return &ibmISInstanceDiskManagementValidator
}

func resourceIBMisInstanceDiskManagementCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instance := d.Get("instance").(string)
	disks := d.Get("disks")
	diskUpdate := disks.([]interface{})

	for _, disk := range diskUpdate {
		diskItem := disk.(map[string]interface{})

		namestr := diskItem["name"].(string)
		diskid := diskItem["id"].(string)

		updateInstanceDiskOptions := &vpcv1.UpdateInstanceDiskOptions{}
		updateInstanceDiskOptions.SetInstanceID(instance)
		updateInstanceDiskOptions.SetID(diskid)
		instanceDiskPatchModel := &vpcv1.InstanceDiskPatch{
			Name: &namestr,
		}

		instanceDiskPatch, err := instanceDiskPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for InstanceDiskPatch: %s", err)
		}
		updateInstanceDiskOptions.SetInstanceDiskPatch(instanceDiskPatch)

		_, response, err := sess.UpdateInstanceDisk(updateInstanceDiskOptions)
		if err != nil {
			return fmt.Errorf("Error calling UpdateInstanceDisk: %s %s", err, response)
		}

	}
	d.SetId(instance)
	return resourceIBMisInstanceDiskManagementRead(d, meta)
}

func resourceIBMisInstanceDiskManagementUpdate(d *schema.ResourceData, meta interface{}) error {
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

			updateInstanceDiskOptions := &vpcv1.UpdateInstanceDiskOptions{}
			updateInstanceDiskOptions.SetInstanceID(d.Id())
			updateInstanceDiskOptions.SetID(diskid)
			instanceDiskPatchModel := &vpcv1.InstanceDiskPatch{
				Name: &namestr,
			}

			instanceDiskPatch, err := instanceDiskPatchModel.AsPatch()
			if err != nil {
				return fmt.Errorf("Error calling asPatch for InstanceDiskPatch: %s", err)
			}
			updateInstanceDiskOptions.SetInstanceDiskPatch(instanceDiskPatch)

			_, _, err = sess.UpdateInstanceDisk(updateInstanceDiskOptions)
			if err != nil {
				return fmt.Errorf("Error updating instance disk: %s", err)
			}

		}
	}
	return resourceIBMisInstanceDiskManagementRead(d, meta)
}

func resourceIBMisInstanceDiskManagementDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")
	return nil
}

func resourceIBMisInstanceDiskManagementRead(d *schema.ResourceData, meta interface{}) error {

	d.Set("instance", d.Id())

	return nil
}
