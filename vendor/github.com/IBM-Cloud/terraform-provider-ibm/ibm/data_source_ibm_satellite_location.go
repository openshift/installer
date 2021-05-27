// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/kubernetesservice-go-sdk/kubernetesserviceapiv1"
)

func dataSourceIBMSatelliteLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSatelliteLocationRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique name for the new Satellite location",
			},
			"managed_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IBM Cloud metro from which the Satellite location is managed",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the new Satellite location",
			},
			"logging_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID for IBM Log Analysis with LogDNA log forwarding",
			},
			"zones": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The names of at least three high availability zones to use for the location",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location CRN",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the resource group",
			},
			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource group",
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceIBMVPCHash,
			},
			"host_attached_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The total number of hosts that are attached to the Satellite location.",
			},
			"host_available_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The available number of hosts that can be assigned to a cluster resource in the Satellite location.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created Date",
			},
			"ingress_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceIBMSatelliteLocationRead(d *schema.ResourceData, meta interface{}) error {
	location := d.Get("location").(string)

	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: &location,
	}

	instance, resp, err := satClient.GetSatelliteLocation(getSatLocOptions)
	if err != nil || instance == nil {
		return fmt.Errorf("Error retrieving IBM cloud satellite location %s : %s\n%s", name, err, resp)

	}

	d.SetId(*instance.ID)
	d.Set("location", location)
	d.Set("description", *instance.Description)
	d.Set("zones", instance.WorkerZones)
	d.Set("managed_from", *instance.Datacenter)
	d.Set("crn", *instance.Crn)
	d.Set("resource_group_id", *instance.ResourceGroup)
	d.Set(ResourceGroupName, *instance.ResourceGroupName)
	d.Set("created_on", *instance.CreatedDate)
	if instance.Hosts != nil {
		d.Set("host_attached_count", *instance.Hosts.Total)
		d.Set("host_available_count", *instance.Hosts.Available)
	}

	if instance.Ingress != nil {
		d.Set("ingress_hostname", *instance.Ingress.Hostname)
		d.Set("ingress_secret", *instance.Ingress.SecretName)
	}

	tags, err := GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"An error occured during reading of instance (%s) tags : %s", d.Id(), err)
	}
	d.Set("tags", tags)

	return nil
}
