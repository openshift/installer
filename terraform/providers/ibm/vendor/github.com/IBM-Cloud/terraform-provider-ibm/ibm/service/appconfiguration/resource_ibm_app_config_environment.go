// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppConfigEnvironment() *schema.Resource {
	return &schema.Resource{
		Read:     resourceEnvironmentRead,
		Create:   resourceEnvironmentCreate,
		Update:   resourceEnvironmentUpdate,
		Delete:   resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Environment name.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Environment Id.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Environment description",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags associated with the environment",
			},
			"color_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Color code to distinguish the environment.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the environment.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the environment data.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Environment URL.",
			},
		},
	}
}

func getAppConfigClient(meta interface{}, guid string) (*appconfigurationv1.AppConfigurationV1, error) {
	appconfigClient, err := meta.(conns.ClientSession).AppConfigurationV1()
	if err != nil {
		return nil, err
	}
	bluemixSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return nil, err
	}
	appConfigURL := fmt.Sprintf("https://%s.apprapp.cloud.ibm.com/apprapp/feature/v1/instances/%s", bluemixSession.Config.Region, guid)
	url := conns.EnvFallBack([]string{"IBMCLOUD_APP_CONFIG_API_ENDPOINT"}, appConfigURL)
	appconfigClient.Service.Options.URL = url
	return appconfigClient, nil
}

func resourceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)
	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}
	options := &appconfigurationv1.CreateEnvironmentOptions{}

	options.SetName(d.Get("name").(string))
	options.SetEnvironmentID(d.Get("environment_id").(string))
	if _, ok := GetFieldExists(d, "description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}
	if _, ok := GetFieldExists(d, "color_code"); ok {
		options.SetColorCode(d.Get("color_code").(string))
	}
	_, response, err := appconfigClient.CreateEnvironment(options)

	if err != nil {
		return flex.FmtErrorf("[ERROR] CreateEnvironment failed %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", guid, *options.EnvironmentID))

	return resourceEnvironmentRead(d, meta)
}

func resourceEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	if ok := d.HasChanges("name", "tags", "color_code", "description"); ok {
		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return nil
		}
		appconfigClient, err := getAppConfigClient(meta, parts[0])
		if err != nil {
			return flex.FmtErrorf(fmt.Sprintf("%s", err))
		}

		options := &appconfigurationv1.UpdateEnvironmentOptions{}

		options.SetName(d.Get("name").(string))
		options.SetEnvironmentID(d.Get("environment_id").(string))
		if _, ok := GetFieldExists(d, "description"); ok {
			options.SetDescription(d.Get("description").(string))
		}
		if _, ok := GetFieldExists(d, "tags"); ok {
			options.SetTags(d.Get("tags").(string))
		}
		if _, ok := GetFieldExists(d, "color_code"); ok {
			options.SetColorCode(d.Get("color_code").(string))
		}

		_, response, err := appconfigClient.UpdateEnvironment(options)
		if err != nil {
			return flex.FmtErrorf("[ERROR] UpdateEnvironment failed %s\n%s", err, response)
		}
		return resourceEnvironmentRead(d, meta)
	}
	return nil
}

func resourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	options := &appconfigurationv1.GetEnvironmentOptions{}

	options.SetExpand(true)
	options.SetEnvironmentID(parts[1])

	result, response, err := appconfigClient.GetEnvironment(options)

	if err != nil {
		return flex.FmtErrorf("[ERROR] GetEnvironment failed %s\n%s", err, response)
	}
	d.Set("guid", parts[0])
	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting name: %s", err)
		}
	}
	if result.EnvironmentID != nil {
		if err = d.Set("environment_id", result.EnvironmentID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting environment_id: %s", err)
		}
	}
	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting description: %s", err)
		}
	}
	if result.Tags != nil {
		if err = d.Set("tags", result.Tags); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting tags: %s", err)
		}
	}
	if result.ColorCode != nil {
		if err = d.Set("color_code", result.ColorCode); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting color_code: %s", err)
		}
	}
	if result.CreatedTime != nil {
		if err = d.Set("created_time", result.CreatedTime.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_time: %s", err)
		}
	}
	if result.UpdatedTime != nil {
		if err = d.Set("updated_time", result.UpdatedTime.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting updated_time: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting href: %s", err)
		}
	}
	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}

	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	options := &appconfigurationv1.DeleteEnvironmentOptions{}
	options.SetEnvironmentID(parts[1])

	response, err := appconfigClient.DeleteEnvironment(options)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return flex.FmtErrorf("[ERROR] DeleteEnvironment failed %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}
