// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISAPProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISAPProfileRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_SAPProfileID: {
				Description:  "SAP Profile ID",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
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
			Attr_FullSystemProfile: {
				Computed:    true,
				Description: "Requires full system for deployment.",
				Type:        schema.TypeBool,
			},
			Attr_Memory: {
				Computed:    true,
				Description: "Amount of memory (in GB).",
				Type:        schema.TypeInt,
			},
			Attr_SAPS: {
				Computed:    true,
				Description: "SAP Application Performance Standard",
				Type:        schema.TypeInt,
			},
			Attr_SupportedSystems: {
				Computed:    true,
				Description: "List of supported systems.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			Attr_Type: {
				Computed:    true,
				Description: "Type of profile.",
				Type:        schema.TypeString,
			},
			Attr_WorkloadType: {
				Computed:    true,
				Description: "Workload Type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPISAPProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	profileID := d.Get(Arg_SAPProfileID).(string)

	client := instance.NewIBMPISAPInstanceClient(ctx, sess, cloudInstanceID)
	sapProfile, err := client.GetSAPProfile(profileID)
	if err != nil {
		log.Printf("[DEBUG] get sap profile failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(*sapProfile.ProfileID)
	d.Set(Attr_Certified, *sapProfile.Certified)
	d.Set(Attr_Cores, *sapProfile.Cores)
	d.Set(Attr_FullSystemProfile, sapProfile.FullSystemProfile)
	d.Set(Attr_Memory, *sapProfile.Memory)
	d.Set(Attr_SAPS, sapProfile.Saps)
	d.Set(Attr_SupportedSystems, sapProfile.SupportedSystems)
	d.Set(Attr_Type, *sapProfile.Type)
	d.Set(Attr_WorkloadType, *&sapProfile.WorkloadTypes)

	return nil
}
