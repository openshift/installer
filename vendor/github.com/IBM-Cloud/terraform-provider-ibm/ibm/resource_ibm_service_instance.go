// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	svcInstanceSuccessStatus  = "succeeded"
	svcInstanceProgressStatus = "in progress"
	svcInstanceFailStatus     = "failed"
)

func resourceIBMServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMServiceInstanceCreate,
		Read:     resourceIBMServiceInstanceRead,
		Update:   resourceIBMServiceInstanceUpdate,
		Delete:   resourceIBMServiceInstanceDelete,
		Exists:   resourceIBMServiceInstanceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name for the service instance",
			},

			"space_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the space in which the instance will be created",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service offering like speech_to_text, text_to_speech etc",
			},

			"credentials": {
				Description: "The service broker-provided credentials to use this service.",
				Type:        schema.TypeMap,
				Sensitive:   true,
				Computed:    true,
			},

			"service_keys": {
				Description: "The service keys asociated with the service instance",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service key name",
						},
						"credentials": {
							Type:        schema.TypeMap,
							Computed:    true,
							Sensitive:   true,
							Description: "The service key credential details like port, username etc",
						},
					},
				},
			},

			"service_plan_guid": {
				Description: "The uniquie identifier of the service offering plan type",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Arbitrary parameters to pass along to the service broker. Must be a JSON object",
			},

			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plan type of the service",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"wait_time_minutes": {
				Description: "Define timeout to wait for the service instances to succeeded/deleted etc.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
			"dashboard_url": {
				Description: "Dashboard URL to access resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceIBMServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)
	name := d.Get("name").(string)
	spaceGUID := d.Get("space_guid").(string)

	svcInst := mccpv2.ServiceInstanceCreateRequest{
		Name:      name,
		SpaceGUID: spaceGUID,
	}

	serviceOff, err := cfClient.ServiceOfferings().FindByLabel(serviceName)
	if err != nil {
		return fmt.Errorf("Error retrieving service offering: %s", err)
	}

	servicePlan, err := cfClient.ServicePlans().FindPlanInServiceOffering(serviceOff.GUID, plan)
	if err != nil {
		return fmt.Errorf("Error retrieving plan: %s", err)
	}
	svcInst.PlanGUID = servicePlan.GUID

	if parameters, ok := d.GetOk("parameters"); ok {
		temp := parameters.(map[string]interface{})
		keyParams := make(map[string]interface{})
		for k, v := range temp {
			if v == "true" || v == "false" {
				b, _ := strconv.ParseBool(v.(string))
				keyParams[k] = b

			} else {
				keyParams[k] = v
			}
		}
		svcInst.Params = keyParams
	}

	if _, ok := d.GetOk("tags"); ok {
		svcInst.Tags = getServiceTags(d)
	}

	service, err := cfClient.ServiceInstances().Create(svcInst)
	if err != nil {
		return fmt.Errorf("Error creating service: %s", err)
	}

	d.SetId(service.Metadata.GUID)

	_, err = waitForServiceInstanceAvailable(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for create service (%s) to be succeeded: %s", d.Id(), err)
	}

	return resourceIBMServiceInstanceRead(d, meta)
}

func resourceIBMServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}

	serviceGUID := d.Id()

	service, err := cfClient.ServiceInstances().Get(serviceGUID, 1)
	if err != nil {
		return fmt.Errorf("Error retrieving service: %s", err)
	}

	servicePlanGUID := service.Entity.ServicePlanGUID
	d.Set("service_plan_guid", servicePlanGUID)
	d.Set("space_guid", service.Entity.SpaceGUID)
	serviceKeys := service.Entity.ServiceKeys
	d.Set("service_keys", flattenServiceInstanceCredentials(serviceKeys))
	d.Set("credentials", Flatten(service.Entity.Credentials))
	d.Set("tags", service.Entity.Tags)
	d.Set("name", service.Entity.Name)
	d.Set("dashboard_url", service.Entity.DashboardURL)

	d.Set("plan", service.Entity.ServicePlan.Entity.Name)

	svcOff, err := cfClient.ServiceOfferings().Get(service.Entity.ServicePlan.Entity.ServiceGUID)
	if err != nil {
		return err
	}
	d.Set("service", svcOff.Entity.Label)

	return nil
}

func resourceIBMServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}

	serviceGUID := d.Id()

	updateReq := mccpv2.ServiceInstanceUpdateRequest{}
	if d.HasChange("name") {
		updateReq.Name = helpers.String(d.Get("name").(string))
	}

	if d.HasChange("plan") {
		plan := d.Get("plan").(string)
		service := d.Get("service").(string)
		serviceOff, err := cfClient.ServiceOfferings().FindByLabel(service)
		if err != nil {
			return fmt.Errorf("Error retrieving service offering: %s", err)
		}

		servicePlan, err := cfClient.ServicePlans().FindPlanInServiceOffering(serviceOff.GUID, plan)
		if err != nil {
			return fmt.Errorf("Error retrieving plan: %s", err)
		}
		updateReq.PlanGUID = helpers.String(servicePlan.GUID)

	}

	if d.HasChange("parameters") {
		updateReq.Params = d.Get("parameters").(map[string]interface{})
	}

	if d.HasChange("tags") {
		tags := getServiceTags(d)
		updateReq.Tags = tags
	}

	_, err = cfClient.ServiceInstances().Update(serviceGUID, updateReq)
	if err != nil {
		return fmt.Errorf("Error updating service: %s", err)
	}

	_, err = waitForServiceInstanceAvailable(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for update service (%s) to be succeeded: %s", d.Id(), err)
	}

	return resourceIBMServiceInstanceRead(d, meta)
}

func resourceIBMServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	id := d.Id()

	err = cfClient.ServiceInstances().Delete(id, true)
	if err != nil {
		return fmt.Errorf("Error deleting service: %s", err)
	}

	_, err = waitForServiceInstanceDelete(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for service (%s) to be deleted: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
func resourceIBMServiceInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	serviceGUID := d.Id()

	service, err := cfClient.ServiceInstances().Get(serviceGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return service.Metadata.GUID == serviceGUID, nil
}

func getServiceTags(d *schema.ResourceData) []string {
	tagSet := d.Get("tags").(*schema.Set)

	if tagSet.Len() == 0 {
		empty := []string{}
		return empty
	}

	tags := make([]string, 0, tagSet.Len())
	for _, elem := range tagSet.List() {
		tag := elem.(string)
		tags = append(tags, tag)
	}
	return tags
}

func waitForServiceInstanceAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	serviceGUID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{svcInstanceProgressStatus},
		Target:  []string{svcInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			service, err := cfClient.ServiceInstances().Get(serviceGUID)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The service instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", err
			}
			if service.Entity.LastOperation.State == svcInstanceFailStatus {
				return service, service.Entity.LastOperation.State, fmt.Errorf("The service instance %s failed: %v", d.Id(), err)
			}
			return service, service.Entity.LastOperation.State, nil
		},
		Timeout:    time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForServiceInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	serviceGUID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending: []string{svcInstanceProgressStatus},
		Target:  []string{svcInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			service, err := cfClient.ServiceInstances().Get(serviceGUID)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return service, svcInstanceSuccessStatus, nil
				}
				return nil, "", err
			}
			if service.Entity.LastOperation.State == svcInstanceFailStatus {
				return service, service.Entity.LastOperation.State, fmt.Errorf("The service instance %s failed to delete: %v", d.Id(), err)
			}
			return service, service.Entity.LastOperation.State, nil
		},
		Timeout:    time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
