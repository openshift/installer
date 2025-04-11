// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/softlayer/softlayer-go/sl"
)

func ResourceIBMPIVolumeGroupAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVolumeGroupActionCreate,
		ReadContext:   resourceIBMPIVolumeGroupActionRead,
		DeleteContext: resourceIBMPIVolumeGroupActionDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeGroupAction: {
				Description: "Performs an action (start stop reset ) on a volume group(one at a time).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Start: {
							Description: "Performs start action on a volume group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Source: {
										Description:  "Indicates the source of the action `master` or `aux`.",
										Required:     true,
										Type:         schema.TypeString,
										ValidateFunc: validate.ValidateAllowedStringValues([]string{Master, Aux}),
									},
								},
							},
							ForceNew: true,
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
						Attr_Stop: {
							Description: "Performs stop action on a volume group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Access: {
										Description: "Indicates the access mode of aux volumes.",
										Required:    true,
										Type:        schema.TypeBool,
									},
								},
							},
							ForceNew: true,
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
						Attr_Reset: {
							Description: "Performs reset action on the volume group to update its status value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Status: {
										Description:  "New status to be set for a volume group.",
										Required:     true,
										Type:         schema.TypeString,
										ValidateFunc: validate.ValidateAllowedStringValues([]string{State_Available}),
									},
								},
							},
							ForceNew: true,
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Type:     schema.TypeList,
			},
			Arg_VolumeGroupID: {
				Description:  "Volume Group ID",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_ReplicationStatus: {
				Computed:    true,
				Description: "Volume Group Replication Status",
				Type:        schema.TypeString,
			},
			Attr_VolumeGroupName: {
				Computed:    true,
				Description: "Volume Group ID",
				Type:        schema.TypeString,
			},
			Attr_VolumeGroupStatus: {
				Computed:    true,
				Description: "Volume Group Status",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIVolumeGroupActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	vgID := d.Get(Arg_VolumeGroupID).(string)
	vgAction, err := expandVolumeGroupAction(d.Get(Arg_VolumeGroupAction).([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	body := vgAction

	client := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	_, err = client.VolumeGroupAction(vgID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, vgID))

	_, err = isWaitForIBMPIVolumeGroupAvailable(ctx, client, vgID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPIVolumeGroupActionRead(ctx, d, meta)
}

func resourceIBMPIVolumeGroupActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)

	vg, err := client.GetDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Attr_VolumeGroupName, vg.Name)
	d.Set(Attr_VolumeGroupStatus, vg.Status)
	d.Set(Attr_ReplicationStatus, vg.ReplicationStatus)

	return nil
}

func resourceIBMPIVolumeGroupActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete or unset concept for volume group action
	d.SetId("")
	return nil
}

// expandVolumeGroupAction retrieve volume group action resource
func expandVolumeGroupAction(data []interface{}) (*models.VolumeGroupAction, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("[ERROR] no pi_volume_group_action received")
	}

	vgAction := models.VolumeGroupAction{}
	action := data[0].(map[string]interface{})

	if v, ok := action[Attr_Start]; ok && len(v.([]interface{})) != 0 {
		vgAction.Start = expandVolumeGroupStartAction(action[Attr_Start].([]interface{}))
		return &vgAction, nil
	}

	if v, ok := action[Attr_Stop]; ok && len(v.([]interface{})) != 0 {
		vgAction.Stop = expandVolumeGroupStopAction(action[Attr_Stop].([]interface{}))
		return &vgAction, nil
	}

	if v, ok := action[Attr_Reset]; ok && len(v.([]interface{})) != 0 {
		vgAction.Reset = expandVolumeGroupResetAction(action[Attr_Reset].([]interface{}))
		return &vgAction, nil
	}
	return nil, fmt.Errorf("[ERROR] no pi_volume_group_action received")
}

func expandVolumeGroupStartAction(start []interface{}) *models.VolumeGroupActionStart {
	if len(start) == 0 {
		return nil
	}

	s := start[0].(map[string]interface{})

	return &models.VolumeGroupActionStart{
		Source: sl.String(s[Attr_Source].(string)),
	}
}

func expandVolumeGroupStopAction(stop []interface{}) *models.VolumeGroupActionStop {
	if len(stop) == 0 {
		return nil
	}

	s := stop[0].(map[string]interface{})

	return &models.VolumeGroupActionStop{
		Access: sl.Bool(s[Attr_Access].(bool)),
	}
}

func expandVolumeGroupResetAction(reset []interface{}) *models.VolumeGroupActionReset {
	if len(reset) == 0 {
		return nil
	}

	s := reset[0].(map[string]interface{})

	return &models.VolumeGroupActionReset{
		Status: sl.String(s[Attr_Status].(string)),
	}
}
