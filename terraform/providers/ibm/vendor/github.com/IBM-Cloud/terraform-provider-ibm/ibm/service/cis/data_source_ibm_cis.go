// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func DataSourceIBMCISInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Resource instance name for example, my cis instance",
				Type:        schema.TypeString,
				Required:    true,
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the resource group in which the cis instance is present",
			},

			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of resource instance",
			},

			"location": {
				Description: "The location or the environment in which cis instance exists",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"service": {
				Description: "The name of the Cloud Internet Services offering, 'internet-svcs'",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"plan": {
				Description: "The plan type of the cis instance",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"status": {
				Description: "The resource instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about the resource",
			},
		},
	}
}

func dataSourceIBMCISInstanceRead(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerAPIV2()
	if err != nil {
		return err
	}
	rsAPI := rsConClient.ResourceServiceInstanceV2()
	name := d.Get("name").(string)

	rsInstQuery := controllerv2.ServiceInstanceQuery{
		Name: name,
	}

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rsInstQuery.ResourceGroupID = rsGrpID.(string)
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return err
		}
		rsInstQuery.ResourceGroupID = defaultRg
	}

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	if service, ok := d.GetOk("service"); ok {

		serviceOff, err := rsCatRepo.FindByName(service.(string), true)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
		}

		rsInstQuery.ServiceID = serviceOff[0].ID
	}

	var instances []models.ServiceInstanceV2

	instances, err = rsAPI.ListInstances(rsInstQuery)
	if err != nil {
		return err
	}
	var filteredInstances []models.ServiceInstanceV2
	var location string

	if loc, ok := d.GetOk("location"); ok {
		location = loc.(string)
		for _, instance := range instances {
			if flex.GetLocation(instance) == location {
				filteredInstances = append(filteredInstances, instance)
			}
		}
	} else {
		filteredInstances = instances
	}

	if len(filteredInstances) == 0 {
		return fmt.Errorf("[ERROR] No resource instance found with name [%s]\nIf not specified please specify more filters like resource_group_id if instance doesn't exists in default group, location or service", name)
	}

	var instance models.ServiceInstanceV2

	if len(filteredInstances) > 1 {
		return fmt.Errorf("[ERROR] More than one resource instance found with name matching [%s]\nIf not specified please specify more filters like resource_group_id if instance doesn't exists in default group, location or service", name)
	}
	instance = filteredInstances[0]

	d.SetId(instance.ID)
	d.Set("status", instance.State)
	d.Set("resource_group_id", instance.ResourceGroupID)
	d.Set("location", instance.RegionID)
	d.Set("guid", instance.Guid)
	serviceOff, err := rsCatRepo.GetServiceName(instance.ServiceID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	d.Set("service", serviceOff)

	servicePlan, err := rsCatRepo.GetServicePlanName(instance.ResourcePlanID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	d.Set("plan", servicePlan)

	d.Set(flex.ResourceName, instance.Name)
	d.Set(flex.ResourceCRN, instance.Crn.String())
	d.Set(flex.ResourceStatus, instance.State)
	d.Set(flex.ResourceGroupName, instance.ResourceGroupName)

	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, rcontroller+"/internet-svcs/"+url.QueryEscape(instance.Crn.String()))

	return nil
}
