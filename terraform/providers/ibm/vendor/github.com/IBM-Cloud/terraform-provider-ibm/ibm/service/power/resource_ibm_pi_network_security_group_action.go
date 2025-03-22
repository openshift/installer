// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMPINetworkSecurityGroupAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkSecurityGroupActionCreate,
		ReadContext:   resourceIBMPINetworkSecurityGroupActionRead,
		UpdateContext: resourceIBMPINetworkSecurityGroupActionUpdate,
		DeleteContext: resourceIBMPINetworkSecurityGroupActionDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_Action: {
				Description:  "Name of the action to take; can be enable to enable NSGs in a workspace or disable to disable NSGs in a workspace.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Disable, Enable}),
			},
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attribute
			Attr_State: {
				Computed:    true,
				Description: "The workspace network security group's state.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPINetworkSecurityGroupActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	action := d.Get(Arg_Action).(string)
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
	wsclient := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	_, err = isWaitForWorkspaceActive(ctx, wsclient, cloudInstanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	body := &models.NetworkSecurityGroupsAction{Action: &action}
	err = nsgClient.Action(body)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = isWaitForNSGStatus(ctx, wsclient, cloudInstanceID, action, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cloudInstanceID)
	return resourceIBMPINetworkSecurityGroupActionRead(ctx, d, meta)
}

func resourceIBMPINetworkSecurityGroupActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	wsclient := instance.NewIBMPIWorkspacesClient(ctx, sess, d.Id())
	ws, err := wsclient.Get(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(Attr_State, ws.Details.NetworkSecurityGroups.State)

	return nil
}
func resourceIBMPINetworkSecurityGroupActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	if d.HasChange(Arg_Action) {
		action := d.Get(Arg_Action).(string)
		nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)

		body := &models.NetworkSecurityGroupsAction{Action: &action}
		err = nsgClient.Action(body)
		if err != nil {
			return diag.FromErr(err)
		}
		wsclient := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
		_, err = isWaitForNSGStatus(ctx, wsclient, cloudInstanceID, action, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPINetworkSecurityGroupActionRead(ctx, d, meta)
}
func resourceIBMPINetworkSecurityGroupActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
func isWaitForWorkspaceActive(ctx context.Context, client *instance.IBMPIWorkspacesClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Provisioning},
		Target:     []string{State_Active},
		Refresh:    isWorkspaceRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
func isWorkspaceRefreshFunc(client *instance.IBMPIWorkspacesClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ws, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if *(ws.Status) == State_Active {
			return ws, State_Active, nil
		}
		if *(ws.Details.NetworkSecurityGroups.State) == State_Provisioning {
			return ws, State_Provisioning, nil
		}
		if *(ws.Details.NetworkSecurityGroups.State) == State_Failed {
			return ws, *ws.Details.NetworkSecurityGroups.State, fmt.Errorf("[ERROR] workspace network security group configuration state is:%s", *ws.Status)
		}

		return ws, State_Configuring, nil
	}
}
func isWaitForNSGStatus(ctx context.Context, client *instance.IBMPIWorkspacesClient, id, action string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Configuring, State_Removing},
		Target:     []string{State_Active, State_Inactive},
		Refresh:    isPERWorkspaceNSGRefreshFunc(client, id, action),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPERWorkspaceNSGRefreshFunc(client *instance.IBMPIWorkspacesClient, id, action string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ws, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if *(ws.Details.NetworkSecurityGroups.State) == State_Active && action == Enable {
			return ws, State_Active, nil
		}
		if *(ws.Details.NetworkSecurityGroups.State) == State_Inactive && action == Disable {
			return ws, State_Inactive, nil
		}
		if *(ws.Details.NetworkSecurityGroups.State) == State_Removing {
			return ws, State_Removing, nil
		}
		if *(ws.Details.NetworkSecurityGroups.State) == State_Error {
			return ws, *ws.Details.NetworkSecurityGroups.State, fmt.Errorf("[ERROR] workspace network security group configuration failed to %s", action)
		}

		return ws, State_Configuring, nil
	}
}
