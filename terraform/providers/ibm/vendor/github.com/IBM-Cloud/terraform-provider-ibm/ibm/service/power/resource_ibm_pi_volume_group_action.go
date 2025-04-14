// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud Instance ID - This is the service_instance_id.",
			},
			PIVolumeGroupID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Volume Group ID",
			},
			PIVolumeGroupAction: {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "Performs an action (start stop reset ) on a volume group(one at a time).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ValidateAllowedStringValues([]string{"master", "aux"}),
									},
								},
							},
						},
						"stop": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"reset": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ValidateAllowedStringValues([]string{"available"}),
									},
								},
							},
						},
					},
				},
			},

			// Computed Attributes
			"volume_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group ID",
			},
			"volume_group_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group Status",
			},
			"replication_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group Replication Status",
			},
		},
	}
}

func resourceIBMPIVolumeGroupActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	vgID := d.Get(PIVolumeGroupID).(string)
	vgAction, err := expandVolumeGroupAction(d.Get(PIVolumeGroupAction).([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	body := vgAction

	client := st.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
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

	client := st.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)

	vg, err := client.GetDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("volume_group_name", vg.Name)
	d.Set("volume_group_status", vg.Status)
	d.Set("replication_status", vg.ReplicationStatus)

	return nil
}

func resourceIBMPIVolumeGroupActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete or unset concept for instance action
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

	if v, ok := action["start"]; ok && len(v.([]interface{})) != 0 {
		vgAction.Start = expandVolumeGroupStartAction(action["start"].([]interface{}))
		return &vgAction, nil
	}

	if v, ok := action["stop"]; ok && len(v.([]interface{})) != 0 {
		vgAction.Stop = expandVolumeGroupStopAction(action["stop"].([]interface{}))
		return &vgAction, nil
	}

	if v, ok := action["reset"]; ok && len(v.([]interface{})) != 0 {
		vgAction.Reset = expandVolumeGroupResetAction(action["reset"].([]interface{}))
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
		Source: sl.String(s["source"].(string)),
	}
}

func expandVolumeGroupStopAction(stop []interface{}) *models.VolumeGroupActionStop {
	if len(stop) == 0 {
		return nil
	}

	s := stop[0].(map[string]interface{})

	return &models.VolumeGroupActionStop{
		Access: sl.Bool(s["access"].(bool)),
	}
}

func expandVolumeGroupResetAction(reset []interface{}) *models.VolumeGroupActionReset {
	if len(reset) == 0 {
		return nil
	}

	s := reset[0].(map[string]interface{})

	return &models.VolumeGroupActionReset{
		Status: sl.String(s["status"].(string)),
	}
}
