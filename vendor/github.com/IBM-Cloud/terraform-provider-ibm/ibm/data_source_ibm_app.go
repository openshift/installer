// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the app",
			},
			"space_guid": {
				Description: "Define space guid to which app belongs",
				Type:        schema.TypeString,
				Required:    true,
			},
			"memory": {
				Description: "The amount of memory each instance should have. In megabytes.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"instances": {
				Description: "The number of instances",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"disk_quota": {
				Description: "The maximum amount of disk available to an instance of an app. In megabytes.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"buildpack": {
				Description: "Buildpack to build the app. 3 options: a) Blank means autodetection; b) A Git Url pointing to a buildpack; c) Name of an installed buildpack.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment_json": {
				Description: "Key/value pairs of all the environment variables to run in your app. Does not include any system or service variables.",
				Type:        schema.TypeMap,
				Computed:    true,
			},
			"route_guid": {
				Description: "Define the route guids which should be bound to the application.",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Computed:    true,
			},
			"service_instance_guid": {
				Description: "Define the service instance guids that should be bound to this application.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"package_state": {
				Description: "The state of the application package whether staged, pending etc",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "The state of the application",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"health_check_http_endpoint": {
				Description: "Endpoint called to determine if the app is healthy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"health_check_type": {
				Description: "Type of health check to perform.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"health_check_timeout": {
				Description: "Timeout in seconds for health checking of an staged app when starting up.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	appAPI := cfClient.Apps()
	name := d.Get("name").(string)
	spaceGUID := d.Get("space_guid").(string)

	app, err := appAPI.FindByName(spaceGUID, name)
	if err != nil {
		return err
	}
	d.SetId(app.GUID)
	d.Set("memory", app.Memory)
	d.Set("disk_quota", app.DiskQuota)
	if app.BuildPack != nil {
		d.Set("buildpack", app.BuildPack)
	}
	d.Set("environment_json", Flatten(app.EnvironmentJSON))
	d.Set("package_state", app.PackageState)
	d.Set("state", app.State)
	d.Set("instances", app.Instances)
	d.Set("health_check_type", app.HealthCheckType)
	d.Set("health_check_http_endpoint", app.HealthCheckHTTPEndpoint)
	d.Set("health_check_timeout", app.HealthCheckTimeout)

	route, err := appAPI.ListRoutes(app.GUID)
	if err != nil {
		return err
	}
	if len(route) > 0 {
		d.Set("route_guid", flattenRoute(route))
	}
	svcBindings, err := appAPI.ListServiceBindings(app.GUID)
	if err != nil {
		return err
	}
	if len(svcBindings) > 0 {
		d.Set("service_instance_guid", flattenServiceBindings(svcBindings))
	}
	return nil
}
