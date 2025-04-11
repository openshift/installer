// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISharedProcessorPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISharedProcessorPoolRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_SharedProcessorPoolID: {
				Description:  "The ID of the shared processor pool.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_AllocatedCores: {
				Computed:    true,
				Description: "The allocated cores in the shared processor pool.",
				Type:        schema.TypeFloat,
			},
			Attr_AvailableCores: {
				Computed:    true,
				Description: "The available cores in the shared processor pool.",
				Type:        schema.TypeFloat,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DedicatedHostID: {
				Computed:    true,
				Description: "The dedicated host ID where the shared processor pool resides.",
				Type:        schema.TypeString,
			},
			Attr_HostID: {
				Computed:    true,
				Description: "The host ID where the shared processor pool resides.",
				Type:        schema.TypeInt,
			},
			Attr_Instances: {
				Computed:    true,
				Description: "List of server instances deployed in the shared processor pool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AvailabilityZone: {
							Computed:    true,
							Description: "Availability zone for the server instances.",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_CPUs: {
							Computed:    true,
							Description: "The amount of cpus for the server instance.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The server instance ID.",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_Memory: {
							Computed:    true,
							Description: "The amount of memory for the server instance.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The server instance name.",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "Status of the instance.",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_Uncapped: {
							Computed:    true,
							Description: "Identifies if uncapped or not.",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						Attr_VCPUs: {
							Computed:    true,
							Description: "The amout of vcpus for the server instance.",
							Optional:    true,
							Type:        schema.TypeFloat,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_Name: {
				Computed:    true,
				Description: "The name of the shared processor pool.",
				Type:        schema.TypeString,
			},
			Attr_ReservedCores: {
				Computed:    true,
				Description: "The amount of reserved cores for the shared processor pool.",
				Type:        schema.TypeInt,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of the shared processor pool.",
				Type:        schema.TypeString,
			},
			Attr_StatusDetail: {
				Computed:    true,
				Description: "The status details of the shared processor pool.",
				Type:        schema.TypeString,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
		},
	}
}

func dataSourceIBMPISharedProcessorPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	poolID := d.Get(Arg_SharedProcessorPoolID).(string)
	client := instance.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(poolID)
	if err != nil || response == nil {
		return diag.Errorf("error fetching the shared processor pool: %v", err)
	}

	d.SetId(*response.SharedProcessorPool.ID)
	d.Set(Attr_AllocatedCores, response.SharedProcessorPool.AllocatedCores)
	d.Set(Attr_AvailableCores, response.SharedProcessorPool.AvailableCores)
	if response.SharedProcessorPool.Crn != "" {
		d.Set(Attr_CRN, response.SharedProcessorPool.Crn)
		if response.SharedProcessorPool.Crn != "" {
			d.Set(Attr_CRN, response.SharedProcessorPool.Crn)
			tags, err := flex.GetGlobalTagsUsingCRN(meta, string(response.SharedProcessorPool.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on get of pi shared_processor_pool (%s) user_tags: %s", *response.SharedProcessorPool.ID, err)
			}
			d.Set(Attr_UserTags, tags)
		}
	}
	d.Set(Attr_DedicatedHostID, response.SharedProcessorPool.DedicatedHostID)
	d.Set(Attr_HostID, response.SharedProcessorPool.HostID)
	d.Set(Attr_Name, response.SharedProcessorPool.Name)
	d.Set(Attr_ReservedCores, response.SharedProcessorPool.ReservedCores)
	d.Set(Attr_Status, response.SharedProcessorPool.Status)
	d.Set(Attr_StatusDetail, response.SharedProcessorPool.StatusDetail)

	serversMap := []map[string]interface{}{}
	if response.Servers != nil {
		for _, s := range response.Servers {
			if s != nil {
				v := map[string]interface{}{
					Attr_AvailabilityZone: s.AvailabilityZone,
					Attr_CPUs:             s.Cpus,
					Attr_ID:               s.ID,
					Attr_Memory:           s.Memory,
					Attr_Name:             s.Name,
					Attr_Status:           s.Status,
					Attr_Uncapped:         s.Uncapped,
					Attr_VCPUs:            s.Vcpus,
				}
				serversMap = append(serversMap, v)
			}
		}
	}
	d.Set(Attr_Instances, serversMap)

	return nil
}
