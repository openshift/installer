// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func DataSourceIBMPISharedProcessorPool() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPISharedProcessorPoolRead,
		Schema: map[string]*schema.Schema{
			Arg_SharedProcessorPoolID: {
				Type:     schema.TypeString,
				Required: true,
			},

			Arg_CloudInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},

			Attr_SharedProcessorPoolName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the shared processor pool",
			},

			Attr_SharedProcessorPoolHostID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The host ID where the shared processor pool resides",
			},

			Attr_SharedProcessorPoolReservedCores: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of reserved cores for the shared processor pool",
			},

			Attr_SharedProcessorPoolAvailableCores: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Shared processor pool available cores",
			},

			Attr_SharedProcessorPoolAllocatedCores: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Shared processor pool allocated cores",
			},

			Attr_SharedProcessorPoolStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the shared processor pool",
			},

			Attr_SharedProcessorPoolStatusDetail: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status details of the shared processor pool",
			},

			Attr_SharedProcessorPoolInstances: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of server instances deployed in the shared processor pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_SharedProcessorPoolInstanceCpus: {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The amount of cpus for the server instance",
						},
						Attr_SharedProcessorPoolInstanceUncapped: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Identifies if uncapped or not",
						},
						Attr_SharedProcessorPoolInstanceAvailabilityZone: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Availability zone for the server instances",
						},
						Attr_SharedProcessorPoolInstanceId: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The server instance ID",
						},
						Attr_SharedProcessorPoolInstanceMemory: {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The amount of memory for the server instance",
						},
						Attr_SharedProcessorPoolInstanceName: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The server instance name",
						},
						Attr_SharedProcessorPoolInstanceStatus: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Status of the server",
						},
						Attr_SharedProcessorPoolInstanceVcpus: {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "The amout of vcpus for the server instance",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISharedProcessorPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	poolID := d.Get(Arg_SharedProcessorPoolID).(string)
	client := st.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(poolID)
	if err != nil || response == nil {
		return diag.Errorf("error fetching the shared processor pool: %v", err)
	}

	d.SetId(*response.SharedProcessorPool.ID)
	d.Set(Attr_SharedProcessorPoolName, response.SharedProcessorPool.Name)
	d.Set(Attr_SharedProcessorPoolReservedCores, response.SharedProcessorPool.ReservedCores)
	d.Set(Attr_SharedProcessorPoolAllocatedCores, response.SharedProcessorPool.AllocatedCores)
	d.Set(Attr_SharedProcessorPoolAvailableCores, response.SharedProcessorPool.AvailableCores)
	d.Set(Attr_SharedProcessorPoolHostID, response.SharedProcessorPool.HostID)
	d.Set(Attr_SharedProcessorPoolStatus, response.SharedProcessorPool.Status)
	d.Set(Attr_SharedProcessorPoolStatusDetail, response.SharedProcessorPool.StatusDetail)

	serversMap := []map[string]interface{}{}
	if response.Servers != nil {
		for _, s := range response.Servers {
			if s != nil {
				v := map[string]interface{}{
					Attr_SharedProcessorPoolInstanceCpus:             s.Cpus,
					Attr_SharedProcessorPoolInstanceUncapped:         s.Uncapped,
					Attr_SharedProcessorPoolInstanceAvailabilityZone: s.AvailabilityZone,
					Attr_SharedProcessorPoolInstanceId:               s.ID,
					Attr_SharedProcessorPoolInstanceMemory:           s.Memory,
					Attr_SharedProcessorPoolInstanceName:             s.Name,
					Attr_SharedProcessorPoolInstanceStatus:           s.Status,
					Attr_SharedProcessorPoolInstanceVcpus:            s.Vcpus,
				}
				serversMap = append(serversMap, v)
			}
		}
	}
	d.Set(Attr_SharedProcessorPoolInstances, serversMap)

	return nil
}
