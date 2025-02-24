// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISystemPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISystemPoolsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_SystemPools: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available system pools within a particular Datacenter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Capacity: {
							Computed:    true,
							Description: "Advertised capacity cores and memory (GB).",
							Type:        schema.TypeMap,
						},
						Attr_CoreMemoryRatio: {
							Computed:    true,
							Description: "Processor to Memory (GB) Ratio.",
							Type:        schema.TypeFloat,
						},
						Attr_MaxAvailable: {
							Computed:    true,
							Description: "Maximum configurable cores and memory (GB) (aggregated from all hosts).",
							Type:        schema.TypeMap,
						},
						Attr_MaxCoresAvailable: {
							Computed:    true,
							Description: "Maximum configurable cores available combined with available memory of that host.",
							Type:        schema.TypeMap,
						},
						Attr_MaxMemoryAvailable: {
							Computed:    true,
							Description: "Maximum configurable memory available combined with available cores of that host.",
							Type:        schema.TypeMap,
						},
						Attr_SharedCoreRatio: {
							Computed:    true,
							Description: "The min-max-default allocation percentage of shared core per vCPU.",
							Type:        schema.TypeMap,
						},
						Attr_SystemPoolName: {
							Computed:    true,
							Description: "The system pool name",
							Type:        schema.TypeString,
						},
						Attr_Systems: {
							Computed:    true,
							Description: "The Datacenter list of servers and their available resources.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Cores: {
										Computed:    true,
										Description: "The host available Processor units.",
										Type:        schema.TypeString,
									},
									Attr_ID: {
										Computed:    true,
										Description: "The host identifier.",
										Type:        schema.TypeString,
									},
									Attr_Memory: {
										Computed:    true,
										Description: "The host available RAM memory in GiB.",
										Type:        schema.TypeString,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_Type: {
							Computed:    true,
							Description: "Type of system hardware.",
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISystemPoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPISystemPoolClient(ctx, sess, cloudInstanceID)
	sps, err := client.GetSystemPools()
	if err != nil {
		log.Printf("[ERROR] get system pools capacity failed %v", err)
		return diag.FromErr(err)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)

	result := make([]map[string]interface{}, 0, len(sps))
	for s, sp := range sps {
		data := map[string]interface{}{
			Attr_SystemPoolName:     s,
			Attr_Capacity:           flattenMax(sp.Capacity),
			Attr_CoreMemoryRatio:    sp.CoreMemoryRatio,
			Attr_MaxAvailable:       flattenMax(sp.MaxAvailable),
			Attr_MaxCoresAvailable:  flattenMax(sp.MaxCoresAvailable),
			Attr_MaxMemoryAvailable: flattenMax(sp.MaxMemoryAvailable),
			Attr_SharedCoreRatio:    flattenSharedCoreRatio(sp.SharedCoreRatio),
			Attr_Type:               sp.Type,
			Attr_Systems:            flattenSystems(sp.Systems),
		}
		result = append(result, data)
	}

	d.Set(Attr_SystemPools, result)

	return nil
}

func flattenMax(s *models.System) map[string]string {
	ret := map[string]interface{}{
		Attr_Cores:  *s.Cores,
		Attr_Memory: *s.Memory,
	}
	return flex.Flatten(ret)
}

func flattenSystem(s *models.System) map[string]string {
	ret := map[string]interface{}{
		Attr_Cores:  *s.Cores,
		Attr_ID:     s.ID,
		Attr_Memory: *s.Memory,
	}
	return flex.Flatten(ret)
}

func flattenSystems(sl []*models.System) (systems []map[string]string) {
	if sl != nil {
		systems = make([]map[string]string, 0, len(sl))
		for _, s := range sl {
			systems = append(systems, flattenSystem(s))
		}
		return systems
	}
	return
}

func flattenSharedCoreRatio(scr *models.MinMaxDefault) map[string]string {
	ret := map[string]interface{}{
		Attr_Default: scr.Default,
		Attr_Max:     scr.Max,
		Attr_Min:     scr.Min,
	}
	return flex.Flatten(ret)
}
