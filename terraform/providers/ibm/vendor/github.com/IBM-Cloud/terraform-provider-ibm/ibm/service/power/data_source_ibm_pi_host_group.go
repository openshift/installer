// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIHostGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIHostGroupRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_HostGroupID: {
				Description:  "Host group ID.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
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
	}
}

func dataSourceIBMPIHostGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	hostGroupID := d.Get(Arg_HostGroupID).(string)
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hostGroup, err := client.GetHostGroup(hostGroupID)
	if err != nil {
		log.Printf("[DEBUG] get host group %v", err)
		return diag.FromErr(err)
	}

	d.SetId(hostGroup.ID)

	d.Set(Attr_CreationDate, hostGroup.CreationDate.String())
	d.Set(Attr_Hosts, hostGroup.Hosts)
	d.Set(Attr_Name, hostGroup.Name)
	d.Set(Attr_Primary, hostGroup.Primary)
	d.Set(Attr_Secondaries, hostGroup.Secondaries)
	return nil
}
