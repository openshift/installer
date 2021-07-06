// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
)

func dataSourceIbmAppConfigEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Environment Id.",
			},
			"expand": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to `true`, returns expanded view of the resource details.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Environment name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Environment description.",
			},
			"tags": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tags associated with the environment.",
			},
			"color_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.",
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

func dataSourceIbmAppConfigEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return err
	}

	options := &appconfigurationv1.GetEnvironmentOptions{}
	options.SetEnvironmentID(d.Get("environment_id").(string))

	if _, ok := d.GetOk("expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}
	result, response, err := appconfigClient.GetEnvironment(options)
	if err != nil {
		log.Printf("GetEnvironment failed %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", guid, *result.EnvironmentID))

	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return fmt.Errorf("error setting name: %s", err)
		}
	}
	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return fmt.Errorf("error setting description: %s", err)
		}
	}
	if result.Tags != nil {
		if err = d.Set("tags", result.Tags); err != nil {
			return fmt.Errorf("error setting tags: %s", err)
		}
	}
	if result.ColorCode != nil {
		if err = d.Set("color_code", result.ColorCode); err != nil {
			return fmt.Errorf("error setting color_code: %s", err)
		}
	}
	if result.CreatedTime != nil {
		if err = d.Set("created_time", result.CreatedTime.String()); err != nil {
			return fmt.Errorf("error setting created_time: %s", err)
		}
	}
	if result.UpdatedTime != nil {
		if err = d.Set("updated_time", result.UpdatedTime.String()); err != nil {
			return fmt.Errorf("error setting updated_time: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return fmt.Errorf("error setting href: %s", err)
		}
	}
	return nil
}
