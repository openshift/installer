// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMISReservedIPPatch() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISReservedIPPatchCreate,
		Read:     resourceIBMISReservedIPPatchRead,
		Update:   resourceIBMISReservedIPPatchUpdate,
		Delete:   resourceIBMISReservedIPPatchDelete,
		Exists:   resourceIBMISReservedIPPatchExists,
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
				Type:         schema.TypeBool,
				Default:      nil,
				AtLeastOneOf: []string{isReservedIPAutoDelete, isReservedIPName},
				Computed:     true,
				Optional:     true,
				Description:  "If set to true, this reserved IP will be automatically deleted",
			},
			isReservedIPName: {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{isReservedIPAutoDelete, isReservedIPName},
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_subnet_reserved_ip", isReservedIPName),
				Description:  "The user-defined or system-provided name for this reserved IP.",
			},
			isReservedIPTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for target.",
			},
			isReservedIPTargetCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The crn for target.",
			},
			isReservedIPLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the reserved IP",
			},
			isReservedIPAddress: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The address for this reserved IP.",
			},
			isReservedIP: {
				Type:        schema.TypeString,
				Required:    true,
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

// resourceIBMISReservedIPCreate Creates a reserved IP given a subnet ID
func resourceIBMISReservedIPPatchCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnetID := d.Get(isSubNetID).(string)
	reservedIPID := d.Get(isReservedIP).(string)
	name := d.Get(isReservedIPName).(string)
	reservedIPPatchModel := &vpcv1.ReservedIPPatch{}
	if name != "" {
		reservedIPPatchModel.Name = &name
	}
	if autoDeleteBoolOk, ok := d.GetOkExists(isReservedIPAutoDelete); ok {
		autoDeleteBool := autoDeleteBoolOk.(bool)
		reservedIPPatchModel.AutoDelete = &autoDeleteBool
	}
	reservedIPPatch, err := reservedIPPatchModel.AsPatch()
	if err != nil {
		return fmt.Errorf("[ERROR] Error updating the reserved IP %s", err)
	}

	options := sess.NewUpdateSubnetReservedIPOptions(subnetID, reservedIPID, reservedIPPatch)

	rip, response, err := sess.UpdateSubnetReservedIP(options)
	if err != nil || response == nil || rip == nil {
		return fmt.Errorf("[ERROR] Error updating the reserved ip patch: %s\n%s", err, response)
	}

	// Set id for the reserved IP as combination of subnet ID and reserved IP ID
	d.SetId(fmt.Sprintf("%s/%s", subnetID, *rip.ID))
	return resourceIBMISReservedIPPatchRead(d, meta)
}

func resourceIBMISReservedIPPatchRead(d *schema.ResourceData, meta interface{}) error {
	rip, err := get(d, meta)
	if err != nil {
		return err
	}

	allIDs, err := flex.IdParts(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] The ID can not be split into subnet ID and reserved IP ID in patch. %s", err)
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
					d.Set(isReservedIPTargetCrn, target.CRN)
				}
			case "*vpcv1.ReservedIPTargetGenericResourceReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetGenericResourceReference)
					d.Set(isReservedIPTargetCrn, target.CRN)
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
					d.Set(isReservedIPTargetCrn, target.CRN)
				}
			case "*vpcv1.ReservedIPTargetVPNGatewayReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetVPNGatewayReference)
					d.Set(isReservedIPTarget, target.ID)
					d.Set(isReservedIPTargetCrn, target.CRN)
				}
			case "*vpcv1.ReservedIPTarget":
				{
					target := targetIntf.(*vpcv1.ReservedIPTarget)
					d.Set(isReservedIPTarget, target.ID)
					d.Set(isReservedIPTargetCrn, target.CRN)
				}
			}
		}
	}
	return nil
}

func resourceIBMISReservedIPPatchUpdate(d *schema.ResourceData, meta interface{}) error {

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
			return fmt.Errorf("[ERROR] Error updating the reserved ip patch %s\n%s", err, response)
		}
	}
	return resourceIBMISReservedIPPatchRead(d, meta)
}

func resourceIBMISReservedIPPatchDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceIBMISReservedIPPatchExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rip, err := get(d, meta)
	if err != nil {
		return false, err
	}
	if err == nil && rip == nil {
		return false, nil
	}
	return true, nil
}
