// Copyright IBM Corp. 2017, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func DataSourceIBMSatelliteLocation() *schema.Resource {
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
			"physical_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An optional physical address of the new Satellite location which is deployed on premise",
			},
			"capabilities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The satellite capabilities attached to the location",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the new Satellite location",
			},
			"coreos_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If Red Hat CoreOS features are enabled within the Satellite location",
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
			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource group",
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      flex.ResourceIBMVPCHash,
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
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
			"service_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom subnet CIDR to provide private IP addresses for services",
			},
			"pod_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom subnet CIDR to provide private IP addresses for pods",
			},
		},
	}
}

func dataSourceIBMSatelliteLocationRead(d *schema.ResourceData, meta interface{}) error {
	location := d.Get("location").(string)

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: &location,
	}

	var instance *kubernetesserviceapiv1.MultishiftGetController
	var response *core.DetailedResponse
	// TO-DO: resource.Retry, resource.RetryError, resource.RetryableError and resource.NonRetryableError
	// seem to be deprecated. This shall be replaced.
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		instance, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
		if err != nil || instance == nil {
			if response != nil && response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		instance, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
	}
	if err != nil || instance == nil {
		return fmt.Errorf("[ERROR] Error retrieving IBM cloud satellite location %s : %s\n%s", location, err, response)
	}

	d.SetId(*instance.ID)
	d.Set("location", location)
	d.Set("description", *instance.Description)
	if instance.PhysicalAddress != nil {
		d.Set("physical_address", *instance.PhysicalAddress)
	}
	if instance.CapabilitiesManagedBySatellite != nil {
		d.Set("capabilities", instance.CapabilitiesManagedBySatellite)
	}
	if instance.CoreosEnabled != nil {
		d.Set("coreos_enabled", *instance.CoreosEnabled)
	}
	d.Set("zones", instance.WorkerZones)
	d.Set("managed_from", *instance.Datacenter)
	d.Set("crn", *instance.Crn)
	d.Set("resource_group_id", *instance.ResourceGroup)
	d.Set(flex.ResourceGroupName, *instance.ResourceGroupName)
	d.Set("created_on", *instance.CreatedDate)
	if instance.Hosts != nil {
		d.Set("host_attached_count", *instance.Hosts.Total)
		d.Set("host_available_count", *instance.Hosts.Available)
	}

	if instance.Ingress != nil {
		d.Set("ingress_hostname", *instance.Ingress.Hostname)
		d.Set("ingress_secret", *instance.Ingress.SecretName)
	}

	getSatHostOptions := &kubernetesserviceapiv1.GetSatelliteHostsOptions{
		Controller: &location,
	}

	hostList, response, err := satClient.GetSatelliteHosts(getSatHostOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving location hosts %s : %s\n%s", location, err, response)
	}
	if hostList != nil {
		d.Set("hosts", flex.FlattenSatelliteHosts(hostList))
	}

	tags, err := flex.GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"An error occured during reading of instance (%s) tags : %s", d.Id(), err)
	}
	d.Set("tags", tags)

	if instance.PodSubnet != nil {
		d.Set("pod_subnet", *instance.PodSubnet)
	}
	if instance.ServiceSubnet != nil {
		d.Set("service_subnet", *instance.ServiceSubnet)
	}

	return nil
}
