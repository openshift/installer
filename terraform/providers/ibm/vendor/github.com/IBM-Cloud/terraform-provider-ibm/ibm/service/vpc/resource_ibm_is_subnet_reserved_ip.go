// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isReservedIPProvisioning     = "provisioning"
	isReservedIPProvisioningDone = "done"
	isReservedIP                 = "reserved_ip"
	isReservedIPTarget           = "target"
	isReservedIPLifecycleState   = "lifecycle_state"
)

func ResourceIBMISReservedIP() *schema.Resource {
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
				ValidateFunc: validate.InvokeValidator("ibm_is_subnet_reserved_ip", isReservedIPName),
				Description:  "The user-defined or system-provided name for this reserved IP.",
			},
			isReservedIPTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for target.",
			},
			isReservedIPLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the reserved IP",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isReservedIPAddress: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The address for this reserved IP.",
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
func ResourceIBMISSubnetReservedIPValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isReservedIPName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmISSubnetReservedIPCResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_subnet_reserved_ip", Schema: validateSchema}
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
	addStr := ""
	if address, ok := d.GetOk(isReservedIPAddress); ok {
		addStr = address.(string)
	}
	if addStr != "" {
		options.Address = &addStr
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
		return fmt.Errorf("[ERROR] Error creating the reserved IP: %s\n%s", err, response)
	}

	// Set id for the reserved IP as combination of subnet ID and reserved IP ID
	d.SetId(fmt.Sprintf("%s/%s", subnetID, *rip.ID))
	_, err = isWaitForReservedIpAvailable(sess, subnetID, *rip.ID, d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for the reserved IP to be available: %s", err)
	}

	return resourceIBMISReservedIPRead(d, meta)
}

func resourceIBMISReservedIPRead(d *schema.ResourceData, meta interface{}) error {
	rip, err := get(d, meta)
	if err != nil {
		return err
	}

	allIDs, err := flex.IdParts(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] The ID can not be split into subnet ID and reserved IP ID. %s", err)
	}
	subnetID := allIDs[0]

	if rip != nil {
		d.Set(isReservedIPAddress, *rip.Address)
		d.Set(isReservedIP, *rip.ID)
		d.Set(isSubNetID, subnetID)
		if rip.LifecycleState != nil {
			d.Set(isReservedIPLifecycleState, *rip.LifecycleState)
		}
		d.Set(isReservedIPAutoDelete, *rip.AutoDelete)
		d.Set(isReservedIPCreatedAt, (*rip.CreatedAt).String())
		d.Set(isReservedIPhref, *rip.Href)
		d.Set(isReservedIPName, *rip.Name)
		d.Set(isReservedIPOwner, *rip.Owner)
		d.Set(isReservedIPType, *rip.ResourceType)
		if rip.Target != nil {
			targetIntf := rip.Target
			switch reflect.TypeOf(targetIntf).String() {
			case "*vpcv1.ReservedIPTargetEndpointGatewayReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetEndpointGatewayReference)
					d.Set(isReservedIPTarget, target.ID)
				}
			case "*vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext)
					d.Set(isReservedIPTarget, target.ID)
				}
			case "*vpcv1.ReservedIPTargetLoadBalancerReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetLoadBalancerReference)
					d.Set(isReservedIPTarget, target.ID)
				}
			case "*vpcv1.ReservedIPTargetVPNGatewayReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetVPNGatewayReference)
					d.Set(isReservedIPTarget, target.ID)
				}
			case "*vpcv1.ReservedIPTarget":
				{
					target := targetIntf.(*vpcv1.ReservedIPTarget)
					d.Set(isReservedIPTarget, target.ID)
				}
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

		allIDs, err := flex.IdParts(d.Id())
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
			return fmt.Errorf("[ERROR] Error updating the reserved IP %s", err)
		}

		options.ReservedIPPatch = reservedIPPatch

		_, response, err := sess.UpdateSubnetReservedIP(options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the reserved IP %s\n%s", err, response)
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
	allIDs, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]
	deleteOptions := sess.NewDeleteSubnetReservedIPOptions(subnetID, reservedIPID)
	response, err := sess.DeleteSubnetReservedIP(deleteOptions)
	if err != nil || response == nil {
		return fmt.Errorf("[ERROR] Error deleting the reserverd ip %s in subnet %s, %s\n%s", reservedIPID, subnetID, err, response)
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
	allIDs, err := flex.IdParts(d.Id())
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]
	options := sess.NewGetSubnetReservedIPOptions(subnetID, reservedIPID)
	rip, response, err := sess.GetSubnetReservedIP(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error Getting Reserved IP : %s\n%s", err, response)
	}
	return rip, nil
}

func isWaitForReservedIpAvailable(sess *vpcv1.VpcV1, subnetid, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for reseved ip (%s/%s) to be available.", subnetid, id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"done", "failed", ""},
		Refresh:    isReserveIpRefreshFunc(sess, subnetid, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isReserveIpRefreshFunc(sess *vpcv1.VpcV1, subnetid, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getreservedipOptions := &vpcv1.GetSubnetReservedIPOptions{
			ID:       &id,
			SubnetID: &subnetid,
		}
		rsip, response, err := sess.GetSubnetReservedIP(getreservedipOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting reserved ip(%s/%s) : %s\n%s", subnetid, id, err, response)
		}
		if rsip.LifecycleState != nil {
			d.Set(isReservedIPLifecycleState, *rsip.LifecycleState)
		}
		d.Set(isReservedIPAddress, *rsip.Address)

		if rsip.LifecycleState != nil && *rsip.LifecycleState == "failed" {
			return rsip, "failed", fmt.Errorf("[ERROR] Error Reserved ip(%s/%s) creation failed : %s\n%s", subnetid, id, err, response)
		}
		if rsip.LifecycleState != nil && *rsip.LifecycleState == "stable" {
			return rsip, "done", nil
		}
		return rsip, "pending", nil
	}
}
