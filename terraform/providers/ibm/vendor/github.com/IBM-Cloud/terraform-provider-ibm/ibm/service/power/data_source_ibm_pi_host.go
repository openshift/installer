// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIHostRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_HostID: {
				Description:  "Host ID.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			//Attribute
			Attr_Capacity: {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AvailableCores: {
							Computed:    true,
							Description: "Number of cores currently available.",
							Type:        schema.TypeFloat,
						},
						Attr_AvailableMemory: {
							Computed:    true,
							Description: "Amount of memory currently available (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_ReservedCore: {
							Computed:    true,
							Description: "Number of cores reserved for system use.",
							Type:        schema.TypeFloat,
						},
						Attr_ReservedMemory: {
							Computed:    true,
							Description: "Amount of memory reserved for system use (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_TotalCore: {
							Computed:    true,
							Description: "Total number of cores of the host.",
							Type:        schema.TypeFloat,
						},
						Attr_TotalMemory: {
							Computed:    true,
							Description: "Total amount of memory of the host (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_UsedCore: {
							Computed:    true,
							Description: "Number of cores in use on the host.",
							Type:        schema.TypeFloat,
						},
						Attr_UsedMemory: {
							Computed:    true,
							Description: "Amount of memory used on the host (in GB).",
							Type:        schema.TypeFloat,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_DisplayName: {
				Computed:    true,
				Description: "Name of the host (chosen by the user).",
				Type:        schema.TypeString,
			},
			Attr_HostGroup: {
				Computed:    true,
				Description: "Link to host group resource.",
				Type:        schema.TypeMap,
			},
			Attr_State: {
				Computed:    true,
				Description: "State of the host (up/down).",
				Type:        schema.TypeString,
			},
			Attr_Status: {
				Computed:    true,
				Description: "Status of the host (enabled/disabled).",
				Type:        schema.TypeString,
			},
			Attr_SysType: {
				Computed:    true,
				Description: "System type.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	hostID := d.Get(Arg_HostID).(string)
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	host, err := client.GetHost(hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(host.ID)
	if host.Capacity != nil {
		d.Set(Attr_Capacity, hostCapacityToMap(host.Capacity))
	}
	if host.DisplayName != "" {
		d.Set(Attr_DisplayName, host.DisplayName)
	}
	if host.HostGroup != nil {
		d.Set(Attr_HostGroup, hostGroupToMap(host.HostGroup))
	}
	if host.State != "" {
		d.Set(Attr_State, host.State)
	}
	if host.Status != "" {
		d.Set(Attr_Status, host.Status)
	}
	if host.SysType != "" {
		d.Set(Attr_SysType, host.SysType)
	}

	return nil
}
