// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	v2 "github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/helpers"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppRoute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppRouteRead,

		Schema: map[string]*schema.Schema{
			"space_guid": {
				Description: "The guid of the space",
				Type:        schema.TypeString,
				Required:    true,
			},
			"domain_guid": {
				Description: "The guid of the domain",
				Type:        schema.TypeString,
				Required:    true,
			},
			"host": {
				Description: "The host of the route",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The path of the route",
				ValidateFunc: validateRoutePath,
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The port of the route",
				ValidateFunc: validateRoutePort,
			},
		},
	}
}

func dataSourceIBMAppRouteRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	spaceAPI := cfClient.Spaces()
	spaceGUID := d.Get("space_guid").(string)
	domainGUID := d.Get("domain_guid").(string)

	params := v2.RouteFilter{
		DomainGUID: domainGUID,
	}

	if host, ok := d.GetOk("host"); ok {
		params.Host = helpers.String(host.(string))
	}

	if port, ok := d.GetOk("port"); ok {
		params.Port = helpers.Int(port.(int))
	}

	if path, ok := d.GetOk("path"); ok {
		params.Path = helpers.String(path.(string))
	}
	route, err := spaceAPI.ListRoutes(spaceGUID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving route: %s", err)
	}
	if len(route) == 0 {
		return fmt.Errorf("No route satifies the given parameters")
	}

	if len(route) > 1 {
		return fmt.Errorf("More than one route satifies the given parameters")
	}

	d.SetId(route[0].GUID)
	return nil

}
