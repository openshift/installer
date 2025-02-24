// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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

// Datasource to list available console languages for an instance
func DataSourceIBMPIInstanceConsoleLanguages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstanceConsoleLanguagesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_InstanceName: {
				Description:  "The unique identifier or name of the instance.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_ConsoleLanguages: {
				Computed:    true,
				Description: "List of all the Console Languages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Code: {
							Computed:    true,
							Description: "Language code.",
							Type:        schema.TypeString,
						},
						Attr_Language: {
							Computed:    true,
							Description: "Language description.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIInstanceConsoleLanguagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	instanceName := d.Get(Arg_InstanceName).(string)

	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	languages, err := client.GetConsoleLanguages(instanceName)
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	if len(languages.ConsoleLanguages) > 0 {
		result := make([]map[string]interface{}, 0, len(languages.ConsoleLanguages))
		for _, language := range languages.ConsoleLanguages {
			l := map[string]interface{}{
				Attr_Code:     *language.Code,
				Attr_Language: language.Language,
			}
			result = append(result, l)
		}
		d.Set(Attr_ConsoleLanguages, result)
	}

	return nil
}
