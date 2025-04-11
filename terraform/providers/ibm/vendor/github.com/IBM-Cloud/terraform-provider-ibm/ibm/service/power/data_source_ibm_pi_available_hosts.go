// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIAvailableHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIAvailableHostsRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
			Attr_AvailableHosts: {
				Computed:    true,
				Description: "Lists of all availabe hosts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AvailableCores: {
							Computed:    true,
							Description: "Core capacity of the host.",
							Type:        schema.TypeFloat,
						},
						Attr_AvailableMemory: {
							Computed:    true,
							Description: "Memory capacity of the host (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_Count: {
							Computed:    true,
							Description: "How many hosts of such type/capacities are available.",
							Type:        schema.TypeInt,
						},
						Attr_SysType: {
							Computed:    true,
							Description: "System type.",
							Type:        schema.TypeString,
						},
					}},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIAvailableHostsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	hostClient := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hostlist, err := hostClient.GetAvailableHosts()
	if err != nil {
		return diag.FromErr(err)
	}
	availableHosts := []map[string]interface{}{}
	for _, value := range hostlist {
		if value.Capacity != nil {
			availableHosts = append(availableHosts, map[string]interface{}{
				Attr_Count:           int(value.Count),
				Attr_SysType:         value.SysType,
				Attr_AvailableCores:  value.Capacity.Cores.Total,
				Attr_AvailableMemory: value.Capacity.Memory.Total,
			})
		}
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_AvailableHosts, availableHosts)

	return nil
}
