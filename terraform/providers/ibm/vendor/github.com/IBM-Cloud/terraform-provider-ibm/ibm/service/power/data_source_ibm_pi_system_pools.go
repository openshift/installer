// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	SystemPoolName     = "system_pool_name"
	SystemPools        = "system_pools"
	SystemPool         = "system_pool"
	Capacity           = "capacity"
	CoreMemoryRatio    = "core_memory_ratio"
	MaxAvailable       = "max_available"
	MaxCoresAvailable  = "max_cores_available"
	MaxMemoryAvailable = "max_memory_available"
	SharedCoreRatio    = "shared_core_ratio"
	Type               = "type"
	Systems            = "systems"
	Cores              = "cores"
	ID                 = "id"
	Memory             = "memory"
	Default            = "default"
	Max                = "max"
	Min                = "min"
)

func DataSourceIBMPISystemPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISystemPoolsRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			SystemPools: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available system pools within a particular DataCenter",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SystemPoolName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The system pool name",
						},
						Capacity: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Advertised capacity cores and memory (GB)",
						},
						CoreMemoryRatio: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Processor to Memory (GB) Ratio",
						},
						MaxAvailable: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Maximum configurable cores and memory (GB) (aggregated from all hosts)",
						},
						MaxCoresAvailable: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Maximum configurable cores available combined with available memory of that host",
						},
						MaxMemoryAvailable: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Maximum configurable memory available combined with available cores of that host",
						},
						SharedCoreRatio: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The min-max-default allocation percentage of shared core per vCPU",
						},
						Systems: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The DataCenter list of servers and their available resources",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Cores: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host available Processor units",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host identifier",
									},
									Memory: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host available RAM memory in GiB",
									},
								},
							},
						},
						Type: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of system hardware",
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

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPISystemPoolClient(ctx, sess, cloudInstanceID)
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
			SystemPoolName:     s,
			Capacity:           flattenSystem(sp.Capacity),
			CoreMemoryRatio:    sp.CoreMemoryRatio,
			MaxAvailable:       flattenSystem(sp.MaxAvailable),
			MaxCoresAvailable:  flattenSystem(sp.MaxCoresAvailable),
			MaxMemoryAvailable: flattenSystem(sp.MaxMemoryAvailable),
			SharedCoreRatio:    flattenSharedCoreRatio(sp.SharedCoreRatio),
			Type:               sp.Type,
			Systems:            flattenSystems(sp.Systems),
		}
		result = append(result, data)
	}

	d.Set(SystemPools, result)

	return nil
}

func flattenSystem(s *models.System) map[string]string {
	ret := map[string]interface{}{
		Cores:  s.Cores,
		ID:     s.ID,
		Memory: s.Memory,
	}
	return flex.Flatten(ret)
}

func flattenSystems(sl []*models.System) (systems []map[string]string) {
	if sl != nil {
		systems = make([]map[string]string, len(sl))
		for _, s := range sl {
			systems = append(systems, flattenSystem(s))
		}
		return systems
	}
	return
}

func flattenSharedCoreRatio(scr *models.MinMaxDefault) map[string]string {
	ret := map[string]interface{}{
		Default: scr.Default,
		Max:     scr.Max,
		Min:     scr.Min,
	}
	return flex.Flatten(ret)
}
