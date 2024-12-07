// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISAPProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISAPProfilesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Profiles: {
				Computed:    true,
				Description: "List of all the SAP Profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Certified: {
							Computed:    true,
							Description: "Has certification been performed on profile.",
							Type:        schema.TypeBool,
						},
						Attr_Cores: {
							Computed:    true,
							Description: "Amount of cores.",
							Type:        schema.TypeInt,
						},
						Attr_Memory: {
							Computed:    true,
							Description: "Amount of memory (in GB).",
							Type:        schema.TypeInt,
						},
						Attr_ProfileID: {
							Computed:    true,
							Description: "SAP Profile ID.",
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Computed:    true,
							Description: "Type of profile.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPISAPProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPISAPInstanceClient(ctx, sess, cloudInstanceID)
	sapProfiles, err := client.GetAllSAPProfiles(cloudInstanceID)
	if err != nil {
		log.Printf("[DEBUG] get all sap profiles failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(sapProfiles.Profiles))
	for _, sapProfile := range sapProfiles.Profiles {
		profile := map[string]interface{}{
			Attr_Certified: *sapProfile.Certified,
			Attr_Cores:     *sapProfile.Cores,
			Attr_Memory:    *sapProfile.Memory,
			Attr_ProfileID: *sapProfile.ProfileID,
			Attr_Type:      *sapProfile.Type,
		}
		result = append(result, profile)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_Profiles, result)

	return nil
}
