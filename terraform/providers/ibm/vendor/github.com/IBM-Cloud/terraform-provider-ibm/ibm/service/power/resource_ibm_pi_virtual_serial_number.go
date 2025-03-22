// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPIVirtualSerialNumber() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVirtualSerialNumberCreate,
		ReadContext:   resourceIBMPIVirtualSerialNumberRead,
		UpdateContext: resourceIBMPIVirtualSerialNumberUpdate,
		DeleteContext: resourceIBMPIVirtualSerialNumberDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description: "This is the Power Instance id that is assigned to the account",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_Description: {
				Description: "Description of virtual serial number.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_InstanceID: {
				Description: "PVM Instance to attach VSN to.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_RetainVirtualSerialNumber: {
				Description:  "Indicates whether to retain virtual serial number after unassigning from PVM instance during deletion.",
				Optional:     true,
				RequiredWith: []string{Arg_InstanceID},
				Type:         schema.TypeBool,
			},
			Arg_Serial: {
				Description:      "Virtual serial number.",
				DiffSuppressFunc: supressVSNDiffAutoAssign,
				ForceNew:         true,
				Required:         true,
				Type:             schema.TypeString,
			},
		},
	}
}

func resourceIBMPIVirtualSerialNumberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIVSNClient(ctx, sess, cloudInstanceID)

	vsnArg := d.Get(Arg_Serial).(string)
	if _, ok := d.GetOk(Arg_InstanceID); !ok && vsnArg == AutoAssign {
		return diag.Errorf("cannot use '%s' unless %s is specified", AutoAssign, Arg_InstanceID)
	}

	serialString := ""
	oldPvmInstanceId := ""
	if vsnArg != AutoAssign {
		vsn, err := client.Get(vsnArg)
		if err != nil {
			return diag.FromErr(err)
		}
		if vsn.PvmInstanceID != nil {
			oldPvmInstanceId = *vsn.PvmInstanceID
		}
		if v, ok := d.GetOk(Arg_Description); ok {
			description := v.(string)
			if description != *vsn.Description {
				if oldPvmInstanceId != "" {
					updateBody := &models.UpdateServerVirtualSerialNumber{
						Description: &description,
					}
					_, err := client.PVMInstanceUpdateVSN(oldPvmInstanceId, updateBody)
					if err != nil {
						return diag.FromErr(err)
					}
				} else {
					updateBody := &models.UpdateVirtualSerialNumber{
						Description: &description,
					}
					_, err := client.Update(vsnArg, updateBody)
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		serialString = vsnArg
	}

	if pvmInstanceId, ok := d.GetOk(Arg_InstanceID); ok {
		pvmInstanceIdArg := pvmInstanceId.(string)
		if oldPvmInstanceId != "" && pvmInstanceIdArg != oldPvmInstanceId {
			return diag.Errorf("please detach virtual serial number from current pvm instance before specifying %s in creation", Arg_InstanceID)
		}
		instanceClient := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
		restartInstance, err := stopLparForVSNChange(ctx, instanceClient, pvmInstanceIdArg, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		if oldPvmInstanceId == "" {
			serialNumber := d.Get(Arg_Serial).(string)
			addBody := &models.AddServerVirtualSerialNumber{
				Serial: &serialNumber,
			}
			if v, ok := d.GetOk(Arg_Description); ok {
				addBody.Description = v.(string)
			}
			_, err = client.PVMInstanceAttachVSN(pvmInstanceIdArg, addBody)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = isWaitForPIInstanceStopped(ctx, instanceClient, pvmInstanceIdArg, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}

			if restartInstance {
				err = startLparAfterVSNChange(ctx, instanceClient, pvmInstanceIdArg, d.Timeout(schema.TimeoutCreate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if vsnArg == AutoAssign {
			vsns, err := client.GetAll(&pvmInstanceIdArg)
			if err != nil {
				return diag.FromErr(err)
			}
			serialString = *vsns[0].Serial

		} else {
			serialString = vsnArg
		}
	}

	id := cloudInstanceID + "/" + serialString
	d.SetId(id)

	return resourceIBMPIVirtualSerialNumberRead(ctx, d, meta)
}

func resourceIBMPIVirtualSerialNumberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	idArr, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := idArr[0]
	serial := idArr[1]
	client := instance.NewIBMPIVSNClient(ctx, sess, cloudInstanceID)

	vsn, err := client.Get(serial)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set(Arg_Description, vsn.Description)
	d.Set(Arg_InstanceID, vsn.PvmInstanceID)
	d.Set(Arg_Serial, vsn.Serial)

	return nil
}

func resourceIBMPIVirtualSerialNumberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	idArr, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := idArr[0]
	client := instance.NewIBMPIVSNClient(ctx, sess, cloudInstanceID)

	if v, ok := d.GetOk(Arg_InstanceID); ok {
		pvmInstanceId := v.(string)
		instanceClient := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
		restartInstance, err := stopLparForVSNChange(ctx, instanceClient, pvmInstanceId, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}

		retainVSN := false
		if v, ok := d.GetOk(Arg_RetainVirtualSerialNumber); ok {
			retainVSN = v.(bool)
		}
		deleteBody := &models.DeleteServerVirtualSerialNumber{
			RetainVSN: retainVSN,
		}
		err = client.PVMInstanceDeleteVSN(pvmInstanceId, deleteBody)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = isWaitForPIInstanceStopped(ctx, instanceClient, pvmInstanceId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		if restartInstance {
			err = startLparAfterVSNChange(ctx, instanceClient, pvmInstanceId, d.Timeout(schema.TimeoutDelete))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	} else {
		serialNumber := d.Get(Arg_Serial).(string)
		err = client.Delete(serialNumber)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}

func resourceIBMPIVirtualSerialNumberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIVSNClient(ctx, sess, cloudInstanceID)

	if d.HasChange(Arg_Description) && !d.HasChange(Arg_InstanceID) {
		newDescription := d.Get(Arg_Description).(string)
		if v, ok := d.GetOk(Arg_InstanceID); ok {
			pvmInstanceId := v.(string)
			updateBody := &models.UpdateServerVirtualSerialNumber{
				Description: &newDescription,
			}

			_, err = client.PVMInstanceUpdateVSN(pvmInstanceId, updateBody)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			updateBody := &models.UpdateVirtualSerialNumber{
				Description: &newDescription,
			}

			vsnArg := d.Get(Arg_Serial).(string)

			_, err = client.Update(vsnArg, updateBody)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange(Arg_InstanceID) {
		oldId, newId := d.GetChange(Arg_InstanceID)
		oldIdString, newIdString := oldId.(string), newId.(string)
		instanceClient := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

		if oldIdString != "" {
			restartInstance, err := stopLparForVSNChange(ctx, instanceClient, oldIdString, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}

			detachBody := &models.DeleteServerVirtualSerialNumber{
				RetainVSN: true,
			}
			err = client.PVMInstanceDeleteVSN(oldIdString, detachBody)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = isWaitForPIInstanceStopped(ctx, instanceClient, oldIdString, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}

			if restartInstance {
				err = startLparAfterVSNChange(ctx, instanceClient, oldIdString, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if newIdString != "" {
			restartInstance, err := stopLparForVSNChange(ctx, instanceClient, newIdString, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}

			serial := d.Get(Arg_Serial).(string)
			addBody := &models.AddServerVirtualSerialNumber{
				Serial: &serial,
			}
			if v, ok := d.GetOk(Arg_Description); ok {
				description := v.(string)
				addBody.Description = description
			}
			_, err = client.PVMInstanceAttachVSN(newIdString, addBody)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = isWaitForPIInstanceStopped(ctx, instanceClient, newIdString, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}

			if restartInstance {
				err = startLparAfterVSNChange(ctx, instanceClient, newIdString, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceIBMPIVirtualSerialNumberRead(ctx, d, meta)
}

func startLparAfterVSNChange(ctx context.Context, client *instance.IBMPIInstanceClient, id string, timeout time.Duration) error {
	body := &models.PVMInstanceAction{
		Action: flex.PtrToString(Action_Start),
	}
	err := client.Action(id, body)
	if err != nil {
		return fmt.Errorf("failed to perform the start action on the pvm instance %v", err)
	}

	_, err = isWaitForPIInstanceAvailable(ctx, client, id, OK, timeout)

	return err
}

func stopLparForVSNChange(ctx context.Context, client *instance.IBMPIInstanceClient, id string, timeout time.Duration) (bool, error) {
	instanceRestart := false
	ins, err := client.Get(id)
	if err != nil {
		return false, fmt.Errorf("failed to get pvm instance (%s): %v", id, err)
	}
	status := *ins.Status
	if strings.ToLower(status) != State_Shutoff {
		body := &models.PVMInstanceAction{
			Action: flex.PtrToString(Action_ImmediateShutdown),
		}
		err := client.Action(id, body)
		if err != nil {
			return false, fmt.Errorf("failed to perform the stop action on the pvm instance %v", err)
		}
		instanceRestart = true
	}

	_, err = isWaitForPIInstanceStopped(ctx, client, id, timeout)

	return instanceRestart, err
}
