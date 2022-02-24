// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceActionAvailable = "available"
	isInstanceActionPending   = "pending"
	isInstanceActionFailed    = "failed"
	isInstanceStopType        = "stop_type"
	isInstanceID              = "instance"
	isInstanceActionForce     = "force_action"
)

func resourceIBMISInstanceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISInstanceActionCreate,
		ReadContext:   resourceIBMISInstanceActionRead,
		UpdateContext: resourceIBMISInstanceActionUpdate,
		DeleteContext: resourceIBMISInstanceActionDelete,
		Exists:        resourceIBMISInstanceActionExists,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			isInstanceAction: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_action", isInstanceAction),
				Description:  "This restart/start/stops an instance.",
			},
			isInstanceActionForce: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the action will be forced immediately, and all queued actions deleted. Ignored for the start action.",
			},
			isInstanceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status",
			},

			isInstanceStatusReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceStatusReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason",
						},

						isInstanceStatusReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},
					},
				},
			},
		},
	}
}

func resourceIBMISInstanceActionValidator() *ResourceValidator {

	instanceActions := "start, reboot, stop"
	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceAction,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              instanceActions})
	ibmISInstanceActionResourceValidator := ResourceValidator{ResourceName: "ibm_is_instance_action", Schema: validateSchema}
	return &ibmISInstanceActionResourceValidator
}

func resourceIBMISInstanceActionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	instanceId := ""
	if insId, ok := d.GetOk(isInstanceID); ok {
		instanceId = insId.(string)
	}

	actiontypeIntf := d.Get(isInstanceAction)
	actiontype := actiontypeIntf.(string)

	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &instanceId,
	}
	instance, response, err := sess.GetInstance(getinsOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error Getting Instance (%s): %s\n%s", instanceId, err, response))
	}
	if (actiontype == "stop" || actiontype == "reboot") && *instance.Status != isInstanceStatusRunning {
		d.Set(isInstanceAction, nil)
		return diag.FromErr(fmt.Errorf("Error with stop/reboot action: Cannot invoke stop/reboot action while instance is not in running state"))
	} else if actiontype == "start" && *instance.Status != isInstanceActionStatusStopped {
		d.Set(isInstanceAction, nil)
		return diag.FromErr(fmt.Errorf("Error with start action: Cannot invoke start action while instance is not in stopped state"))
	}
	createinsactoptions := &vpcv1.CreateInstanceActionOptions{
		InstanceID: &instanceId,
		Type:       &actiontype,
	}
	if instanceActionForceIntf, ok := d.GetOk(isInstanceActionForce); ok {
		force := instanceActionForceIntf.(bool)
		createinsactoptions.Force = &force
	}
	_, response, err = sess.CreateInstanceAction(createinsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error Creating Instance Action: %s\n%s", err, response))
	}
	if actiontype == "stop" {
		_, err = isWaitForInstanceActionStop(sess, d.Timeout(schema.TimeoutUpdate), instanceId, d)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if actiontype == "start" || actiontype == "reboot" {
		_, err = isWaitForInstanceActionStart(sess, d.Timeout(schema.TimeoutUpdate), instanceId, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(instanceId)
	return resourceIBMISInstanceActionRead(context, d, meta)
}

func resourceIBMISInstanceActionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()

	options := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := sess.GetInstance(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error getting instance (%s): %s\n%s", id, err, response))
	}

	d.Set(isInstanceStatus, *instance.Status)
	statusReasonsList := make([]map[string]interface{}, 0)
	if instance.StatusReasons != nil {
		for _, sr := range instance.StatusReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isInstanceStatusReasonsCode] = *sr.Code
				currentSR[isInstanceStatusReasonsMessage] = *sr.Message
				statusReasonsList = append(statusReasonsList, currentSR)
			}
		}
	}
	d.Set(isInstanceStatusReasons, statusReasonsList)
	return nil
}

func resourceIBMISInstanceActionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	_, actiontypeIntf := d.GetChange(isInstanceAction)
	actiontype := actiontypeIntf.(string)
	id := d.Id()

	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := sess.GetInstance(getinsOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error Getting Instance (%s): %s\n%s", id, err, response))
	}
	if (actiontype == "stop" || actiontype == "reboot") && *instance.Status != isInstanceStatusRunning {
		d.Set(isInstanceAction, nil)
		return diag.FromErr(fmt.Errorf("Error with stop/reboot action: Cannot invoke stop/reboot action while instance is not in running state"))
	} else if actiontype == "start" && *instance.Status != isInstanceActionStatusStopped {
		d.Set(isInstanceAction, nil)
		return diag.FromErr(fmt.Errorf("Error with start action: Cannot invoke start action while instance is not in stopped state"))
	}
	createinsactoptions := &vpcv1.CreateInstanceActionOptions{
		InstanceID: &id,
		Type:       &actiontype,
	}
	_, response, err = sess.CreateInstanceAction(createinsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error Creating Instance Action: %s\n%s", err, response))
	}
	if actiontype == "stop" {
		_, err = isWaitForInstanceActionStop(sess, d.Timeout(schema.TimeoutUpdate), id, d)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if actiontype == "start" || actiontype == "reboot" {
		_, err = isWaitForInstanceActionStart(sess, d.Timeout(schema.TimeoutUpdate), id, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceIBMISInstanceActionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceIBMISInstanceActionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	id := d.Id()
	getInstanceOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := sess.GetInstance(getInstanceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting instance : %s\n%s", err, response)
	}
	return true, err
}
