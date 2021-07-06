// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isReservedIPProvisioning     = "provisioning"
	isReservedIPProvisioningDone = "done"
	isReservedIP                 = "reserved_ip"
	isReservedIPTarget           = "target"
)

func resourceIBMISReservedIP() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISReservedIPCreate,
		Read:     resourceIBMISReservedIPRead,
		Update:   resourceIBMISReservedIPUpdate,
		Delete:   resourceIBMISReservedIPDelete,
		Exists:   resourceIBMISReservedIPExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
			*/
			isSubNetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet identifier.",
			},
			isReservedIPAutoDelete: {
				Type:        schema.TypeBool,
				Default:     nil,
				Computed:    true,
				Optional:    true,
				Description: "If set to true, this reserved IP will be automatically deleted",
			},
			isReservedIPName: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_subnet_reserved_ip", isReservedIPName),
				Description:  "The user-defined or system-provided name for this reserved IP.",
			},
			isReservedIPTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for target.",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isReservedIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined or system-provided name for this reserved IP.",
			},
			isReservedIP: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the reserved IP.",
			},
			isReservedIPCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the reserved IP was created.",
			},
			isReservedIPhref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this reserved IP.",
			},
			isReservedIPOwner: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of a reserved IP, defining whether it is managed by the user or the provider.",
			},
			isReservedIPType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}
func resourceIBMISSubnetReservedIPValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isReservedIPName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmISSubnetReservedIPCResourceValidator := ResourceValidator{ResourceName: "ibm_is_subnet_reserved_ip", Schema: validateSchema}
	return &ibmISSubnetReservedIPCResourceValidator
}

// resourceIBMISReservedIPCreate Creates a reserved IP given a subnet ID
func resourceIBMISReservedIPCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnetID := d.Get(isSubNetID).(string)
	options := sess.NewCreateSubnetReservedIPOptions(subnetID)

	nameStr := ""
	if name, ok := d.GetOk(isReservedIPName); ok {
		nameStr = name.(string)
	}
	if nameStr != "" {
		options.Name = &nameStr
	}

	autoDeleteBool := d.Get(isReservedIPAutoDelete).(bool)
	options.AutoDelete = &autoDeleteBool
	if t, ok := d.GetOk(isReservedIPTarget); ok {
		targetId := t.(string)
		options.Target = &vpcv1.ReservedIPTargetPrototype{
			ID: &targetId,
		}
	}
	rip, response, err := sess.CreateSubnetReservedIP(options)
	if err != nil || response == nil || rip == nil {
		return fmt.Errorf("Error creating the reserved IP: %s\n%s", err, response)
	}

	// Set id for the reserved IP as combination of subnet ID and reserved IP ID
	d.SetId(fmt.Sprintf("%s/%s", subnetID, *rip.ID))

	return resourceIBMISReservedIPRead(d, meta)
}

func resourceIBMISReservedIPRead(d *schema.ResourceData, meta interface{}) error {
	rip, err := get(d, meta)
	if err != nil {
		return err
	}

	allIDs, err := idParts(d.Id())
	if err != nil {
		return fmt.Errorf("The ID can not be split into subnet ID and reserved IP ID. %s", err)
	}
	subnetID := allIDs[0]

	if rip != nil {
		d.Set(isReservedIPAddress, *rip.Address)
		d.Set(isReservedIP, *rip.ID)
		d.Set(isSubNetID, subnetID)
		d.Set(isReservedIPAutoDelete, *rip.AutoDelete)
		d.Set(isReservedIPCreatedAt, (*rip.CreatedAt).String())
		d.Set(isReservedIPhref, *rip.Href)
		d.Set(isReservedIPName, *rip.Name)
		d.Set(isReservedIPOwner, *rip.Owner)
		d.Set(isReservedIPType, *rip.ResourceType)
		if rip.Target != nil {
			target, ok := rip.Target.(*vpcv1.ReservedIPTarget)
			if ok {
				d.Set(isReservedIPTarget, target.ID)
			}
		}
	}
	return nil
}

func resourceIBMISReservedIPUpdate(d *schema.ResourceData, meta interface{}) error {

	// For updating the name
	nameChanged := d.HasChange(isReservedIPName)
	autoDeleteChanged := d.HasChange(isReservedIPAutoDelete)

	if nameChanged || autoDeleteChanged {
		sess, err := vpcClient(meta)
		if err != nil {
			return err
		}

		allIDs, err := idParts(d.Id())
		if err != nil {
			return err
		}
		subnetID := allIDs[0]
		reservedIPID := allIDs[1]

		options := &vpcv1.UpdateSubnetReservedIPOptions{
			SubnetID: &subnetID,
			ID:       &reservedIPID,
		}

		patch := new(vpcv1.ReservedIPPatch)

		if nameChanged {
			name := d.Get(isReservedIPName).(string)
			patch.Name = core.StringPtr(name)
		}

		if autoDeleteChanged {
			autoDelete := d.Get(isReservedIPAutoDelete).(bool)
			patch.AutoDelete = core.BoolPtr(autoDelete)
		}

		reservedIPPatch, err := patch.AsPatch()
		if err != nil {
			return fmt.Errorf("Error updating the reserved IP %s", err)
		}

		options.ReservedIPPatch = reservedIPPatch

		_, response, err := sess.UpdateSubnetReservedIP(options)
		if err != nil {
			return fmt.Errorf("Error updating the reserved IP %s\n%s", err, response)
		}
	}
	return resourceIBMISReservedIPRead(d, meta)
}

func resourceIBMISReservedIPDelete(d *schema.ResourceData, meta interface{}) error {

	rip, err := get(d, meta)
	if err != nil {
		return err
	}
	if err == nil && rip == nil {
		// If there is no such reserved IP, it can not be deleted
		return nil
	}

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	allIDs, err := idParts(d.Id())
	if err != nil {
		return err
	}
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]
	deleteOptions := sess.NewDeleteSubnetReservedIPOptions(subnetID, reservedIPID)
	response, err := sess.DeleteSubnetReservedIP(deleteOptions)
	if err != nil || response == nil {
		return fmt.Errorf("Error deleting the reserverd ip %s in subnet %s, %s\n%s", reservedIPID, subnetID, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISReservedIPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rip, err := get(d, meta)
	if err != nil {
		return false, err
	}
	if err == nil && rip == nil {
		return false, nil
	}
	return true, nil
}

// get is a generic function that gets the reserved ip given subnet id and reserved ip
func get(d *schema.ResourceData, meta interface{}) (*vpcv1.ReservedIP, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return nil, err
	}
	allIDs, err := idParts(d.Id())
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]
	options := sess.NewGetSubnetReservedIPOptions(subnetID, reservedIPID)
	rip, response, err := sess.GetSubnetReservedIP(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil, nil
		}
		return nil, fmt.Errorf("Error Getting Reserved IP : %s\n%s", err, response)
	}
	return rip, nil
}
