// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	v2 "github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/helpers"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppRoute() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMAppRouteCreate,
		Read:     resourceIBMAppRouteRead,
		Update:   resourceIBMAppRouteUpdate,
		Delete:   resourceIBMAppRouteDelete,
		Exists:   resourceIBMAppRouteExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The host portion of the route. Required for shared-domains.",
			},

			"space_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the associated space",
			},

			"domain_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the associated domain",
			},

			"port": {
				Description:  "The port of the route. Supported for domains of TCP router groups only.",
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateRoutePort,
			},

			"path": {
				Description:  "The path for a route as raw text.Paths must be between 2 and 128 characters.Paths must start with a forward slash '/'.Paths must not contain a '?'",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateRoutePath,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIBMAppRouteCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}

	spaceGUID := d.Get("space_guid").(string)
	domainGUID := d.Get("domain_guid").(string)

	params := v2.RouteRequest{
		SpaceGUID:  spaceGUID,
		DomainGUID: domainGUID,
	}

	if host, ok := d.GetOk("host"); ok {
		params.Host = host.(string)
	}

	if port, ok := d.GetOk("port"); ok {
		params.Port = helpers.Int(port.(int))
	}

	if path, ok := d.GetOk("path"); ok {
		params.Path = path.(string)
	}

	route, err := cfClient.Routes().Create(params)
	if err != nil {
		return fmt.Errorf("Error creating route: %s", err)
	}

	d.SetId(route.Metadata.GUID)

	return resourceIBMAppRouteRead(d, meta)
}

func resourceIBMAppRouteRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	routeGUID := d.Id()

	route, err := cfClient.Routes().Get(routeGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving route: %s", err)
	}

	d.Set("host", route.Entity.Host)
	d.Set("space_guid", route.Entity.SpaceGUID)
	d.Set("domain_guid", route.Entity.DomainGUID)
	if route.Entity.Port != nil {
		d.Set("port", route.Entity.Port)
	}
	d.Set("path", route.Entity.Path)

	return nil
}

func resourceIBMAppRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}

	routeGUID := d.Id()
	params := v2.RouteUpdateRequest{}

	if d.HasChange("host") {
		params.Host = helpers.String(d.Get("host").(string))
	}

	if d.HasChange("port") {
		params.Port = helpers.Int(d.Get("port").(int))
	}

	if d.HasChange("path") {
		params.Path = helpers.String(d.Get("path").(string))
	}

	_, err = cfClient.Routes().Update(routeGUID, params)
	if err != nil {
		return fmt.Errorf("Error updating route: %s", err)
	}
	return resourceIBMAppRouteRead(d, meta)
}

func resourceIBMAppRouteDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	routeGUID := d.Id()

	err = cfClient.Routes().Delete(routeGUID, false)
	if err != nil {
		return fmt.Errorf("Error deleting route: %s", err)
	}

	d.SetId("")

	return nil
}
func resourceIBMAppRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	routeGUID := d.Id()

	route, err := cfClient.Routes().Get(routeGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return route.Metadata.GUID == routeGUID, nil
}
