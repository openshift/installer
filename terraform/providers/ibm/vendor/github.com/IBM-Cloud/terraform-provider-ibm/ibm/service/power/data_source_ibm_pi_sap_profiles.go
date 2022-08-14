// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPISAPProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISAPProfilesRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PISAPProfiles: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PISAPProfileCertified: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Has certification been performed on profile",
						},
						PISAPProfileCores: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Amount of cores",
						},
						PISAPProfileMemory: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Amount of memory (in GB)",
						},
						PISAPProfileID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAP Profile ID",
						},
						PISAPProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of profile",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISAPProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := instance.NewIBMPISAPInstanceClient(ctx, sess, cloudInstanceID)
	sapProfiles, err := client.GetAllSAPProfiles(cloudInstanceID)
	if err != nil {
		log.Printf("[DEBUG] get all sap profiles failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(sapProfiles.Profiles))
	for _, sapProfile := range sapProfiles.Profiles {
		profile := map[string]interface{}{
			PISAPProfileCertified: *sapProfile.Certified,
			PISAPProfileCores:     *sapProfile.Cores,
			PISAPProfileMemory:    *sapProfile.Memory,
			PISAPProfileID:        *sapProfile.ProfileID,
			PISAPProfileType:      *sapProfile.Type,
		}
		result = append(result, profile)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PISAPProfiles, result)

	return nil
}
