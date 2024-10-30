// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIHostGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIHostGroupsRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
			Attr_HostGroups: {
				Computed:    true,
				Description: "List of host groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CreationDate: {
							Computed:    true,
							Description: "Date/Time of host group creation.",
							Type:        schema.TypeString,
						},
						Attr_Hosts: {
							Computed:    true,
							Description: "List of hosts.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type: schema.TypeList,
						},
						Attr_ID: {
							Computed:    true,
							Description: "Host group ID.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "Name of the host group.",
							Type:        schema.TypeString,
						},
						Attr_Primary: {
							Computed:    true,
							Description: "ID of the workspace owning the host group.",
							Type:        schema.TypeString,
						},
						Attr_Secondaries: {
							Computed:    true,
							Description: "IDs of workspaces the host group has been shared with.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Type: schema.TypeList,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIHostGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hostGroups, err := client.GetHostGroups()
	if err != nil {
		return diag.FromErr(err)
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	d.Set(Attr_HostGroups, flattenHostGroups(hostGroups))

	return nil
}
func flattenHostGroups(hostGroupList models.HostGroupList) []map[string]interface{} {
	hostGroups := make([]map[string]interface{}, 0, len(hostGroupList))
	for _, hg := range hostGroupList {
		hostGroup := map[string]interface{}{
			Attr_CreationDate: hg.CreationDate.String(),
			Attr_Hosts:        hg.Hosts,
			Attr_ID:           hg.ID,
			Attr_Name:         hg.Name,
			Attr_Primary:      hg.Primary,
			Attr_Secondaries:  hg.Secondaries,
		}
		hostGroups = append(hostGroups, hostGroup)
	}
	return hostGroups
}
