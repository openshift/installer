// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package db2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	rg "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMDb2Instance() *schema.Resource {
	riSchema := resourcecontroller.DataSourceIBMResourceInstance().Schema

	riSchema["high_availability"] = &schema.Schema{
		Description: "If you require high availability, please choose this option",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["instance_type"] = &schema.Schema{
		Description: "Available machine type flavours (default selection will assume smallest configuration)",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["backup_location"] = &schema.Schema{
		Description: "Cross Regional backups can be stored across multiple regions in a zone. Regional backups are stored in only specific region.",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["disk_encryption_instance_crn"] = &schema.Schema{
		Description: "Cross Regional disk encryption crn",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["disk_encryption_crn"] = &schema.Schema{
		Description: "Cross Regional disk encryption crn",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["oracle_compatibility"] = &schema.Schema{
		Description: "Indicates whether is has compatibility for oracle or not",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["subscription_id"] = &schema.Schema{
		Description: "For PerformanceSubscription plans a Subscription ID is required. It is not required for Performance plans.",
		Optional:    true,
		Type:        schema.TypeString,
	}

	return &schema.Resource{
		Read:   dataSourceIBMDb2InstanceRead,
		Schema: riSchema,
	}
}

func getInstancesNext(next *string) (string, error) {
	if reflect.ValueOf(next).IsNil() {
		return "", nil
	}
	u, err := url.Parse(*next)
	if err != nil {
		return "", err
	}
	q := u.Query()
	return q.Get("next_url"), nil
}

func dataSourceIBMDb2InstanceRead(d *schema.ResourceData, meta interface{}) error {
	var instance rc.ResourceInstance
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()
	if _, ok := d.GetOk("name"); ok {
		name := d.Get("name").(string)
		resourceInstanceListOptions := rc.ListResourceInstancesOptions{
			Name: &name,
		}

		if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
			rg := rsGrpID.(string)
			resourceInstanceListOptions.ResourceGroupID = &rg
		}

		if service, ok := d.GetOk("service"); ok {

			serviceOff, err := rsCatRepo.FindByName(service.(string), true)
			if err != nil {
				return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
			}
			resourceId := serviceOff[0].ID
			resourceInstanceListOptions.ResourceID = &resourceId
		}

		next_url := ""
		var instances []rc.ResourceInstance
		for {
			if next_url != "" {
				resourceInstanceListOptions.Start = &next_url
			}
			listInstanceResponse, resp, err := rsConClient.ListResourceInstances(&resourceInstanceListOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
			}
			next_url, err = getInstancesNext(listInstanceResponse.NextURL)
			if err != nil {
				return fmt.Errorf("[DEBUG] ListResourceInstances failed. Error occurred while parsing NextURL: %s", err)

			}
			instances = append(instances, listInstanceResponse.Resources...)
			if next_url == "" {
				break
			}
		}

		var filteredInstances []rc.ResourceInstance
		var location string

		if loc, ok := d.GetOk("location"); ok {
			location = loc.(string)
			for _, instance := range instances {
				if flex.GetLocationV2(instance) == location {
					filteredInstances = append(filteredInstances, instance)
				}
			}
		} else {
			filteredInstances = instances
		}

		if len(filteredInstances) == 0 {
			return fmt.Errorf("[ERROR] No resource instance found with name [%s]\nIf not specified please specify more filters like resource_group_id if instance doesn't exists in default group, location or service", name)
		}
		if len(filteredInstances) > 1 {
			return fmt.Errorf("[ERROR] More than one resource instance found with name matching [%s]\nIf not specified please specify more filters like resource_group_id if instance doesn't exists in default group, location or service", name)
		}
		instance = filteredInstances[0]
	} else if _, ok := d.GetOk("identifier"); ok {
		instanceGUID := d.Get("identifier").(string)
		getResourceInstanceOptions := &rc.GetResourceInstanceOptions{
			ID: &instanceGUID,
		}
		instances, res, err := rsConClient.GetResourceInstance(getResourceInstanceOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] No resource instance found with id [%s\n%v]", instanceGUID, res)
		}
		instance = *instances
		d.Set("name", instance.Name)
	}
	d.SetId(*instance.ID)
	d.Set("status", instance.State)
	d.Set("resource_group_id", instance.ResourceGroupID)
	d.Set("location", instance.RegionID)
	serviceOff, err := rsCatRepo.GetServiceName(*instance.ResourceID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	d.Set("service", serviceOff)

	d.Set(flex.ResourceName, instance.Name)
	d.Set(flex.ResourceCRN, instance.CRN)
	d.Set(flex.ResourceStatus, instance.State)
	// ### Modification : Setting the onetime credentials
	d.Set("onetime_credentials", instance.OnetimeCredentials)
	if instance.Parameters != nil {
		params, err := json.Marshal(instance.Parameters)
		if err != nil {
			return fmt.Errorf("[ERROR] Error marshalling instance parameters: %s", err)
		}
		if err = d.Set("parameters_json", string(params)); err != nil {
			return fmt.Errorf("[ERROR] Error setting instance parameters json: %s", err)
		}
	}
	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return err
	}
	GetResourceGroup := rg.GetResourceGroupOptions{
		ID: instance.ResourceGroupID,
	}
	resourceGroup, resp, err := rMgtClient.GetResourceGroup(&GetResourceGroup)
	if err != nil || resourceGroup == nil {
		log.Printf("[ERROR] Error retrieving resource group: %s %s", err, resp)
	}
	if resourceGroup != nil && resourceGroup.Name != nil {
		d.Set(flex.ResourceGroupName, resourceGroup.Name)
	}
	d.Set("guid", instance.GUID)
	if len(instance.Extensions) == 0 {
		d.Set("extensions", instance.Extensions)
	} else {
		d.Set("extensions", flex.Flatten(instance.Extensions))
	}

	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, rcontroller+"/services/")

	servicePlan, err := rsCatRepo.GetServicePlanName(*instance.ResourcePlanID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	d.Set("plan", servicePlan)
	d.Set("crn", instance.CRN)
	tags, err := flex.GetTagsUsingCRN(meta, *instance.CRN)
	if err != nil {
		log.Printf(
			"Error on get of resource instance tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)

	d.Set("high_availability", instance.Parameters["high_availability"])
	d.Set("instance_type", instance.Parameters["instance_type"])
	d.Set("backup_location", instance.Parameters["backup_location"])
	d.Set("disk_encryption_instance_crn", instance.Parameters["disk_encryption_instance_crn"])
	d.Set("disk_encryption_crn", instance.Parameters["disk_encryption_crn"])
	d.Set("oracle_compatibility", instance.Parameters["oracle_compatibility"])
	d.Set("subscription_id", instance.Parameters["subscription_id"])

	return nil
}
