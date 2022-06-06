// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	ConsoleLanguages    = "console_languages"
	ConsoleLanguageCode = "code"
	ConsoleLanguageDesc = "language"
)

/*
Datasource to get the list of available console languages for an instance
*/
func DataSourceIBMPIInstanceConsoleLanguages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstanceConsoleLanguagesRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PIInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique identifier or name of the instance",
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			ConsoleLanguages: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ConsoleLanguageCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "language code",
						},
						ConsoleLanguageDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "language description",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIInstanceConsoleLanguagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	instanceName := d.Get(helpers.PIInstanceName).(string)

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
				ConsoleLanguageCode: *language.Code,
				ConsoleLanguageDesc: language.Language,
			}
			result = append(result, l)
		}
		d.Set(ConsoleLanguages, result)
	}

	return nil
}
