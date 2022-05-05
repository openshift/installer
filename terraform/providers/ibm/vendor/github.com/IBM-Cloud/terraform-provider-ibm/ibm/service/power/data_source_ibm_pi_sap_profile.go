// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPISAPProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISAPProfileRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PISAPInstanceProfileID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SAP Profile ID",
			},
			// Computed Attributes
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
			PISAPProfileType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of profile",
			},
		},
	}
}

func dataSourceIBMPISAPProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	profileID := d.Get(PISAPInstanceProfileID).(string)

	client := instance.NewIBMPISAPInstanceClient(ctx, sess, cloudInstanceID)
	sapProfile, err := client.GetSAPProfile(profileID)
	if err != nil {
		log.Printf("[DEBUG] get sap profile failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(*sapProfile.ProfileID)
	d.Set(PISAPProfileCertified, *sapProfile.Certified)
	d.Set(PISAPProfileCores, *sapProfile.Cores)
	d.Set(PISAPProfileMemory, *sapProfile.Memory)
	d.Set(PISAPProfileType, *sapProfile.Type)

	return nil
}
