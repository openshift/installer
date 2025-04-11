// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"strings"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"log"
	"time"
)

func ResourceIBMPIInstanceAction() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceIBMPIInstanceActionCreate,
		ReadContext:   resourceIBMPIInstanceActionRead,
		UpdateContext: resourceIBMPIInstanceActionUpdate,
		DeleteContext: resourceIBMPIInstanceActionDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_Action: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Action_HardReboot, Action_ImmediateShutdown, Action_ResetState, Action_Start, Action_Stop, Action_SoftReboot}),
				Description:  "PVM instance action type",
			},
			Arg_CloudInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI Cloud instance id",
			},
			Arg_HealthStatus: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{OK, Warning}),
				Default:      OK,
				Description:  "Set the health status of the PVM instance to connect it faster",
			},
			Arg_InstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PVM instance ID",
			},
			// Attributes
			Attr_HealthStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PVM's health status value",
			},
			Attr_Progress: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The progress of an operation",
			},
			Attr_Status: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the PVM instance",
			},
		},
	}
}

func resourceIBMPIInstanceActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	adiag := takeInstanceAction(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if adiag != nil {
		return adiag
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	id := d.Get(Arg_InstanceID).(string)
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, id))

	return resourceIBMPIInstanceActionRead(ctx, d, meta)
}

func resourceIBMPIInstanceActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, id, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPIInstanceClient(context.Background(), sess, cloudInstanceID)
	powervmdata, err := client.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Attr_Status, powervmdata.Status)
	d.Set(Attr_Progress, powervmdata.Progress)
	if powervmdata.Health != nil {
		d.Set(Attr_HealthStatus, powervmdata.Health.Status)
	}

	return nil
}

func resourceIBMPIInstanceActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange(Arg_Action) {
		adiag := takeInstanceAction(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if adiag != nil {
			return adiag
		}
	}

	return resourceIBMPIInstanceActionRead(ctx, d, meta)
}

func resourceIBMPIInstanceActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete or unset concept for instance action
	d.SetId("")
	return nil
}

func takeInstanceAction(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	id := d.Get(Arg_InstanceID).(string)
	action := d.Get(Arg_Action).(string)
	targetHealthStatus := d.Get(Arg_HealthStatus).(string)

	var targetStatus string
	if action == Action_Stop || action == Action_ImmediateShutdown {
		targetStatus = State_Shutoff
	} else if action == Action_ResetState {
		targetStatus = State_Active
		targetHealthStatus = Critical
	} else {
		// action is "start" or "soft-reboot" or "hard-reboot"
		targetStatus = State_Active
	}

	client := st.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	// special case for action "start", "stop", "immediate-shutdown"
	// skip calling action if instance is already in desired state
	if action == Action_Start || action == Action_Stop || action == Action_ImmediateShutdown {
		pvm, err := client.Get(id)
		if err != nil {
			return diag.FromErr(err)
		}

		if strings.ToLower(*pvm.Status) == targetStatus && pvm.Health != nil && (pvm.Health.Status == targetHealthStatus || pvm.Health.Status == OK) {
			log.Printf("[DEBUG] skipping as action %s not needed on the instance %s", action, id)
			return nil
		}
	}

	body := &models.PVMInstanceAction{Action: &action}
	log.Printf("Calling the IBM PI Action %s on the instance %s", action, id)

	err = client.Action(id, body)
	if err != nil {
		log.Printf("[ERROR] failed to perform the action on the instance %v", err)
		return diag.FromErr(err)
	}

	log.Printf("Executed the action on the instance")

	log.Printf("Calling the check for %s opertion to check for status %s", action, targetStatus)
	_, err = isWaitForPIInstanceActionStatus(ctx, client, id, timeout, targetStatus, targetHealthStatus)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func isWaitForPIInstanceActionStatus(ctx context.Context, client *st.IBMPIInstanceClient, id string, timeout time.Duration, targetStatus, targetHealthStatus string) (interface{}, error) {
	log.Printf("Waiting for the action to be performed on the instance %s", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending},
		Target:     []string{targetStatus, State_Error, ""},
		Refresh:    isPIActionRefreshFunc(client, id, targetStatus, targetHealthStatus),
		Delay:      30 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIActionRefreshFunc(client *st.IBMPIInstanceClient, id, targetStatus, targetHealthStatus string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("Waiting for the target status to be [ %s ]", targetStatus)
		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if strings.ToLower(*pvm.Status) == targetStatus && (pvm.Health.Status == targetHealthStatus || pvm.Health.Status == OK) {
			log.Printf("The health status is now %s", pvm.Health.Status)
			return pvm, targetStatus, nil
		}

		if strings.ToLower(*pvm.Status) == State_Error {
			if pvm.Fault != nil {
				err = fmt.Errorf("failed to perform the action on the instance: %s", pvm.Fault.Message)
			} else {
				err = fmt.Errorf("failed to perform the action on the instance")
			}
			return pvm, *pvm.Status, err
		}

		return pvm, State_Pending, nil
	}
}
