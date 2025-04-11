// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package resourcecontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	rg "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func DataSourceIBMResourceInstance() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceIBMResourceInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Resource instance name for example, myobjectstorage",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"identifier", "name"},
			},
			"identifier": {
				Description:   "Resource instance guid",
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{"identifier", "name"},
				ConflictsWith: []string{"resource_group_id", "name", "location", "service"},
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the resource group in which the instance is present",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_resource_instance",
					"resource_group_id"),
			},

			"location": {
				Description: "The location or the environment in which instance exists",
				Optional:    true,
				Type:        schema.TypeString,
				Computed:    true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_resource_instance",
					"location"),
			},

			"service": {
				Description: "The service type of the instance",
				Optional:    true,
				Type:        schema.TypeString,
				Computed:    true,
			},

			"plan": {
				Description: "The plan type of the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"status": {
				Description: "The resource instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Tags of Resource Instance",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Guid of resource instance",
			},

			// ### Modification addded onetime_credentials to Resource scehama
			"onetime_credentials": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "onetime_credentials of resource instance",
			},

			"parameters_json": {
				Description: "Parameters asociated with instance in json string",
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

			"extensions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The extended metadata as a map associated with the resource instance.",
			},
		},
	}
}
func DataSourceIBMResourceInstanceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "resource_group_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_group",
			CloudDataRange:             []string{"resolved_to:id"},
			Optional:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "region",
			Optional:                   true})

	ibmIBMResourceInstanceValidator := validate.ResourceValidator{ResourceName: "ibm_resource_instance", Schema: validateSchema}
	return &ibmIBMResourceInstanceValidator
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

func DataSourceIBMResourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
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
	// ### Modifiction : Setting the onetime credientials
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

	return nil
}
