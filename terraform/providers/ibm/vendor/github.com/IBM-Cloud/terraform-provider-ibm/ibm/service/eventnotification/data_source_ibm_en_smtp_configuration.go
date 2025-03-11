// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnSMTPConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSMTPConfigurationRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"en_smtp_configuration_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for SMTP.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP description.",
			},
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain Name.",
			},
			"config": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Payload describing a SMTP configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dkim": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The DKIM attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"txt_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DMIM text name.",
									},
									"txt_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DMIM text value.",
									},
									"verification": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "dkim verification.",
									},
								},
							},
						},
						"en_authorization": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The en_authorization attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"verification": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "en_authorization verification.",
									},
								},
							},
						},
						"spf": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The SPF attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"txt_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "spf text name.",
									},
									"txt_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "spf text value.",
									},
									"verification": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "spf verification.",
									},
								},
							},
						},
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created time.",
			},
		},
	}
}

func dataSourceIBMEnSMTPConfigurationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSMTPConfigurationOptions := &en.GetSMTPConfigurationOptions{}

	getSMTPConfigurationOptions.SetInstanceID(d.Get("instance_id").(string))
	getSMTPConfigurationOptions.SetID(d.Get("en_smtp_configuration_id").(string))

	smtpConfiguration, _, err := eventNotificationsClient.GetSMTPConfigurationWithContext(context, getSMTPConfigurationOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", *getSMTPConfigurationOptions.InstanceID, *getSMTPConfigurationOptions.ID))

	if err = d.Set("name", smtpConfiguration.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("description", smtpConfiguration.Description); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("domain", smtpConfiguration.Domain); err != nil {
		return diag.FromErr(err)
	}

	config := []map[string]interface{}{}
	if smtpConfiguration.Config != nil {
		modelMap, err := dataSourceIBMEnSMTPConfigurationSMTPConfigToMap(smtpConfiguration.Config)
		if err != nil {
			return diag.FromErr(err)
		}
		config = append(config, modelMap)
	}
	if err = d.Set("config", config); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("updated_at", flex.DateTimeToString(smtpConfiguration.UpdatedAt)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataSourceIBMEnSMTPConfigurationSMTPConfigToMap(model *en.SMTPConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dkim != nil {
		dkimMap, err := dataSourceIBMEnSMTPConfigurationDkimAttributesToMap(model.Dkim)
		if err != nil {
			return modelMap, err
		}
		modelMap["dkim"] = []map[string]interface{}{dkimMap}
	}
	if model.EnAuthorization != nil {
		enAuthorizationMap, err := dataSourceIBMEnSMTPConfigurationEnAuthAttributesToMap(model.EnAuthorization)
		if err != nil {
			return modelMap, err
		}
		modelMap["en_authorization"] = []map[string]interface{}{enAuthorizationMap}
	}
	if model.Spf != nil {
		spfMap, err := dataSourceIBMEnSMTPConfigurationSpfAttributesToMap(model.Spf)
		if err != nil {
			return modelMap, err
		}
		modelMap["spf"] = []map[string]interface{}{spfMap}
	}
	return modelMap, nil
}

func dataSourceIBMEnSMTPConfigurationDkimAttributesToMap(model *en.SmtpdkimAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TxtName != nil {
		modelMap["txt_name"] = model.TxtName
	}
	if model.TxtValue != nil {
		modelMap["txt_value"] = model.TxtValue
	}
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}

func dataSourceIBMEnSMTPConfigurationEnAuthAttributesToMap(model *en.EnAuthAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}

func dataSourceIBMEnSMTPConfigurationSpfAttributesToMap(model *en.SpfAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TxtName != nil {
		modelMap["txt_name"] = model.TxtName
	}
	if model.TxtValue != nil {
		modelMap["txt_value"] = model.TxtValue
	}
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}
