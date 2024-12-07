// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package resourcecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	RsInstanceSuccessStatus       = "active"
	RsInstanceProgressStatus      = "in progress"
	RsInstanceProvisioningStatus  = "provisioning"
	RsInstanceInactiveStatus      = "inactive"
	RsInstanceFailStatus          = "failed"
	RsInstanceRemovedStatus       = "removed"
	RsInstanceReclamation         = "pending_reclamation"
	RsInstanceUpdateSuccessStatus = "succeeded"
)

func ResourceIBMResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMResourceInstanceCreate,
		Read:     ResourceIBMResourceInstanceRead,
		Update:   ResourceIBMResourceInstanceUpdate,
		Delete:   ResourceIBMResourceInstanceDelete,
		Exists:   ResourceIBMResourceInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name for the resource instance",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service offering like cloud-object-storage, kms etc",
			},

			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plan type of the service",
			},

			"location": {
				Description: "The location where the instance available",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				ValidateFunc: validate.InvokeValidator("ibm_resource_instance",
					"location"),
			},

			"resource_group_id": {
				Description: "The resource group id",
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Computed:    true,
				ValidateFunc: validate.InvokeValidator("ibm_resource_instance",
					"resource_group_id"),
			},

			// ### Modification : Adding onetime_credientails into the response scehama
			"onetime_credentials": {
				Description: "A boolean that dictates if the onetime_credentials is true or false.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"parameters": {
				Type:          schema.TypeMap,
				Optional:      true,
				Description:   "Arbitrary parameters to pass. Must be a JSON object",
				ConflictsWith: []string{"parameters_json"},
			},
			"parameters_json": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"parameters"},
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "Arbitrary parameters to pass in Json string format",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_resource_instance", "tags")},
				Set:      flex.ResourceIBMVPCHash,
				ValidateFunc: validate.InvokeValidator("ibm_resource_instance",
					"tags"),
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of resource instance",
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},

			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Guid of resource instance",
			},

			"service_endpoints": {
				Description:  "Types of the service endpoints. Possible values are 'public', 'private', 'public-and-private'.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private", "public-and-private"}),
			},

			"dashboard_url": {
				Description: "Dashboard URL to access resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"plan_history": {
				Description: "The plan history of the instance.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"account_id": {
				Description: "An alpha-numeric value identifying the account ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"resource_group_crn": {
				Description: "The long ID (full CRN) of the resource group",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"resource_id": {
				Description: "The unique ID of the offering",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"resource_plan_id": {
				Description: "The unique ID of the plan associated with the offering",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"target_crn": {
				Description: "The full deployment CRN as defined in the global catalog",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"state": {
				Description: "The current state of the instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"type": {
				Description: "The type of the instance, e.g. service_instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"sub_type": {
				Description: "The sub-type of instance, e.g. cfaas .",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"allow_cleanup": {
				Description: "A boolean that dictates if the resource instance should be deleted (cleaned up) during the processing of a region instance delete call.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"locked": {
				Description: "A boolean that dictates if the resource instance should be deleted (cleaned up) during the processing of a region instance delete call.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"last_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The status of the last operation requested on the instance",
			},

			"resource_aliases_url": {
				Description: "The relative path to the resource aliases for the instance.",
				Type:        schema.TypeString,
				Deprecated:  "Remove this attribute's configuration as it's no longer in use and the attribute will be removed in the upcoming major version of the provider 1.71.0.",
				Computed:    true,
			},

			"resource_bindings_url": {
				Description: "The relative path to the resource bindings for the instance.",
				Type:        schema.TypeString,
				Deprecated:  "Remove this attribute's configuration as it's no longer in use and the attribute will be removed in the upcoming major version of the provider 1.71.0.",
				Computed:    true,
			},

			"resource_keys_url": {
				Description: "The relative path to the resource keys for the instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"created_at": {
				Type:        schema.TypeString,
				Description: "The date when the instance was created.",
				Computed:    true,
			},

			"created_by": {
				Type:        schema.TypeString,
				Description: "The subject who created the instance.",
				Computed:    true,
			},

			"update_at": {
				Type:        schema.TypeString,
				Description: "The date when the instance was last updated.",
				Computed:    true,
			},

			"update_by": {
				Type:        schema.TypeString,
				Description: "The subject who updated the instance.",
				Computed:    true,
			},

			"deleted_at": {
				Type:        schema.TypeString,
				Description: "The date when the instance was deleted.",
				Computed:    true,
			},

			"deleted_by": {
				Type:        schema.TypeString,
				Description: "The subject who deleted the instance.",
				Computed:    true,
			},

			"scheduled_reclaim_at": {
				Type:        schema.TypeString,
				Description: "The date when the instance was scheduled for reclamation.",
				Computed:    true,
			},

			"scheduled_reclaim_by": {
				Type:        schema.TypeString,
				Description: "The subject who initiated the instance reclamation.",
				Computed:    true,
			},

			"restored_at": {
				Type:        schema.TypeString,
				Description: "The date when the instance under reclamation was restored.",
				Computed:    true,
			},

			"restored_by": {
				Type:        schema.TypeString,
				Description: "The subject who restored the instance back from reclamation.",
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

func ResourceIBMResourceInstanceValidator() *validate.ResourceValidator {
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
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "tags",
			Optional:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "region",
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmResourceInstanceResourceValidator := validate.ResourceValidator{ResourceName: "ibm_resource_instance", Schema: validateSchema}
	return &ibmResourceInstanceResourceValidator
}

func ResourceIBMResourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	rsInst := rc.CreateResourceInstanceOptions{
		Name: &name,
	}

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.FindByName(serviceName, true)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	if metadata, ok := serviceOff[0].Metadata.(*models.ServiceResourceMetadata); ok {
		if !metadata.Service.RCProvisionable {
			return fmt.Errorf("%s cannot be provisioned by resource controller", serviceName)
		}
	} else {
		return fmt.Errorf("[ERROR] Cannot create instance of resource %s\nUse 'ibm_service_instance' if the resource is a Cloud Foundry service", serviceName)
	}

	servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	rsInst.ResourcePlanID = &servicePlan

	deployments, err := rsCatRepo.ListDeployments(servicePlan)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving deployment for plan %s : %s", plan, err)
	}
	if len(deployments) == 0 {
		return fmt.Errorf("[ERROR] No deployment found for service plan : %s", plan)
	}
	deployments, supportedLocations := FilterDeployments(deployments, location)

	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return fmt.Errorf("[ERROR] No deployment found for service plan %s at location %s.\nValid location(s) are: %q.\nUse 'ibm_service_instance' if the service is a Cloud Foundry service", plan, location, locationList)
	}

	rsInst.Target = &deployments[0].CatalogCRN

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rg := rsGrpID.(string)
		rsInst.ResourceGroup = &rg
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return err
		}
		rsInst.ResourceGroup = &defaultRg
	}

	params := map[string]interface{}{}

	if serviceEndpoints, ok := d.GetOk("service_endpoints"); ok {
		params["service-endpoints"] = serviceEndpoints.(string)
	}

	if parameters, ok := d.GetOk("parameters"); ok {
		temp := parameters.(map[string]interface{})
		for k, v := range temp {
			if v == "true" || v == "false" {
				b, _ := strconv.ParseBool(v.(string))
				params[k] = b
			} else if strings.HasPrefix(v.(string), "[") && strings.HasSuffix(v.(string), "]") {
				//transform v.(string) to be []string
				arrayString := v.(string)
				result := []string{}
				trimLeft := strings.TrimLeft(arrayString, "[")
				trimRight := strings.TrimRight(trimLeft, "]")
				if len(trimRight) == 0 {
					params[k] = result
				} else {
					array := strings.Split(trimRight, ",")
					for _, a := range array {
						result = append(result, strings.Trim(a, "\""))
					}
					params[k] = result
				}
			} else {
				params[k] = v
			}
		}

	}
	if s, ok := d.GetOk("parameters_json"); ok {
		json.Unmarshal([]byte(s.(string)), &params)
	}

	rsInst.Parameters = params

	//Start to create resource instance
	instance, resp, err := rsConClient.CreateResourceInstance(&rsInst)
	if err != nil {
		log.Printf(
			"Error when creating resource instance: %s, Instance info  NAME->%s, LOCATION->%s, GROUP_ID->%s, PLAN_ID->%s",
			err, *rsInst.Name, *rsInst.Target, *rsInst.ResourceGroup, *rsInst.ResourcePlanID)
		return fmt.Errorf("[ERROR] Error when creating resource instance: %s with resp code: %s", err, resp)
	}

	d.SetId(*instance.ID)

	_, err = waitForResourceInstanceCreate(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for create resource instance (%s) to be succeeded: %s", d.Id(), err)
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}

	return ResourceIBMResourceInstanceRead(d, meta)
}
func ResourceIBMResourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}

	tags, err := flex.GetTagsUsingCRN(meta, *instance.CRN)
	if err != nil {
		log.Printf(
			"Error on get of resource instance tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)
	d.Set("name", instance.Name)
	d.Set("status", instance.State)
	d.Set("resource_group_id", instance.ResourceGroupID)
	if instance.CRN != nil {
		location := strings.Split(*instance.CRN, ":")
		if len(location) > 5 {
			d.Set("location", location[5])
		}
	}
	d.Set("crn", instance.CRN)
	d.Set("dashboard_url", instance.DashboardURL)

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.GetServiceName(*instance.ResourceID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	d.Set("service", serviceOff)

	d.Set(flex.ResourceName, instance.Name)
	d.Set(flex.ResourceCRN, instance.CRN)
	d.Set(flex.ResourceStatus, instance.State)
	d.Set(flex.ResourceGroupName, instance.ResourceGroupCRN)

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
	d.Set("guid", instance.GUID)
	// ### Modificataion : Setting  "onetime_credentials"
	d.Set("onetime_credentials", instance.OnetimeCredentials)
	if instance.Parameters != nil {
		if endpoint, ok := instance.Parameters["service-endpoints"]; ok {
			d.Set("service_endpoints", endpoint)
		}
	}

	if len(instance.Extensions) == 0 {
		d.Set("extensions", instance.Extensions)
	} else {
		d.Set("extensions", flex.Flatten(instance.Extensions))
	}
	d.Set("account_id", instance.AccountID)
	d.Set("restored_by", instance.RestoredBy)
	if instance.RestoredAt != nil {
		d.Set("restored_at", instance.RestoredAt.String())
	}
	d.Set("scheduled_reclaim_by", instance.ScheduledReclaimBy)
	if instance.ScheduledReclaimAt != nil {
		d.Set("scheduled_reclaim_at", instance.ScheduledReclaimAt.String())
	}
	if instance.DeletedAt != nil {
		d.Set("deleted_at", instance.DeletedAt.String())
	}
	d.Set("deleted_by", instance.DeletedBy)
	if instance.UpdatedAt != nil {
		d.Set("update_at", instance.UpdatedAt.String())
	}
	if instance.CreatedAt != nil {
		d.Set("created_at", instance.CreatedAt.String())
	}
	d.Set("update_by", instance.UpdatedBy)
	d.Set("created_by", instance.CreatedBy)
	d.Set("resource_keys_url", instance.ResourceKeysURL)
	d.Set("resource_bindings_url", instance.ResourceBindingsURL)
	d.Set("resource_aliases_url", instance.ResourceAliasesURL)
	if instance.LastOperation != nil {
		operation, err := flex.StructToMap(instance.LastOperation)
		if err == nil {
			d.Set("last_operation", flex.Flatten(operation))
		}
	}
	d.Set("locked", instance.Locked)
	d.Set("allow_cleanup", instance.AllowCleanup)
	d.Set("type", instance.Type)
	d.Set("state", instance.State)
	d.Set("sub_type", instance.SubType)
	d.Set("target_crn", instance.TargetCRN)
	d.Set("resource_plan_id", instance.ResourcePlanID)
	d.Set("resource_id", instance.ResourceID)
	d.Set("resource_group_crn", instance.ResourceGroupCRN)
	if instance.PlanHistory != nil {
		d.Set("plan_history", flattenPlanHistory(instance.PlanHistory))
	}

	return nil
}

func ResourceIBMResourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	instanceID := d.Id()

	resourceInstanceUpdate := rc.UpdateResourceInstanceOptions{
		ID: &instanceID,
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		resourceInstanceUpdate.Name = &name
	}

	if d.HasChange("plan") {
		plan := d.Get("plan").(string)
		service := d.Get("service").(string)
		rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
		if err != nil {
			return err
		}
		rsCatRepo := rsCatClient.ResourceCatalog()

		serviceOff, err := rsCatRepo.FindByName(service, true)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
		}

		servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
		}

		resourceInstanceUpdate.ResourcePlanID = &servicePlan

	}
	params := map[string]interface{}{}

	if d.HasChange("service_endpoints") {
		endpoint := d.Get("service_endpoints").(string)
		params["service-endpoints"] = endpoint
	}

	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	if d.HasChange("parameters") {
		instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
		}

		if parameters, ok := d.GetOk("parameters"); ok {
			temp := parameters.(map[string]interface{})
			for k, v := range temp {
				if v == "true" || v == "false" {
					b, _ := strconv.ParseBool(v.(string))
					params[k] = b
				} else if strings.HasPrefix(v.(string), "[") && strings.HasSuffix(v.(string), "]") {
					//transform v.(string) to be []string
					result := []string{}
					arrayString := v.(string)
					trimLeft := strings.TrimLeft(arrayString, "[")
					trimRight := strings.TrimRight(trimLeft, "]")
					if len(trimRight) == 0 {
						params[k] = result
					} else {
						array := strings.Split(trimRight, ",")
						for _, a := range array {
							result = append(result, strings.Trim(a, "\""))
						}
						params[k] = result
					}

				} else {
					params[k] = v
				}
			}
		}
		if _, ok := params["service-endpoints"]; !ok {
			serviceEndpoints := d.Get("service_endpoints").(string)
			if serviceEndpoints != "" {
				endpoint := d.Get("service_endpoints").(string)
				params["service-endpoints"] = endpoint
			} else if _, ok := instance.Parameters["service-endpoints"]; ok {
				params["service-endpoints"] = instance.Parameters["service-endpoints"]
			}
		}

	}

	if d.HasChange("service_endpoints") || d.HasChange("parameters") {
		resourceInstanceUpdate.Parameters = params
	}
	if d.HasChange("parameters_json") {
		if s, ok := d.GetOk("parameters_json"); ok {
			json.Unmarshal([]byte(s.(string)), &params)
			resourceInstanceUpdate.Parameters = params
		}
	}
	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting resource instance: %s with resp code: %s", err, resp)
	}

	if d.HasChange("tags") {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource instance (%s) tags: %s", d.Id(), err)
		}
	}

	_, resp, err = rsConClient.UpdateResourceInstance(&resourceInstanceUpdate)
	if err != nil {
		return fmt.Errorf("[ERROR] Error updating resource instance: %s with resp code: %s", err, resp)
	}

	_, err = waitForResourceInstanceUpdate(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for update resource instance (%s) to be succeeded: %s", d.Id(), err)
	}

	return ResourceIBMResourceInstanceRead(d, meta)
}

func ResourceIBMResourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	id := d.Id()
	recursive := true
	resourceInstanceDelete := rc.DeleteResourceInstanceOptions{
		ID:        &id,
		Recursive: &recursive,
	}

	resp, error := rsConClient.DeleteResourceInstance(&resourceInstanceDelete)
	if error != nil {
		if resp != nil && resp.StatusCode == 410 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting resource instance: %s with resp code: %s", error, resp)
	}

	_, err = waitForResourceInstanceDelete(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for resource instance (%s) to be deleted: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
func ResourceIBMResourceInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting resource instance: %s with resp code: %s", err, resp)
	}
	if instance != nil && (strings.Contains(*instance.State, "removed") || strings.Contains(*instance.State, RsInstanceReclamation)) {
		log.Printf("[WARN] Removing instance from state because it's in removed or pending_reclamation state")
		d.SetId("")
		return false, nil
	}

	return *instance.ID == instanceID, nil
}

func waitForResourceInstanceCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{RsInstanceProgressStatus, RsInstanceInactiveStatus, RsInstanceProvisioningStatus},
		Target:  []string{RsInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", fmt.Errorf("[ERROR] Get the resource instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance '%s' creation failed: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}

func waitForResourceInstanceUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{RsInstanceProgressStatus, RsInstanceInactiveStatus},
		Target:  []string{RsInstanceSuccessStatus, RsInstanceUpdateSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", fmt.Errorf("[ERROR] Get the resource instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if instance.LastOperation != nil && instance.LastOperation.Async != nil && *instance.LastOperation.Async {
				if *instance.LastOperation.State == RsInstanceFailStatus {
					return instance, *instance.LastOperation.State, fmt.Errorf("[ERROR] The resource instance '%s' update failed: %v", d.Id(), err)
				}
				return instance, *instance.LastOperation.State, nil
			} else {
				if *instance.State == RsInstanceFailStatus {
					return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance '%s' update failed: %v", d.Id(), err)
				}
				return instance, *instance.State, nil
			}
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}

func waitForResourceInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	stateConf := &retry.StateChangeConf{
		Pending: []string{RsInstanceProgressStatus, RsInstanceInactiveStatus, RsInstanceSuccessStatus},
		Target:  []string{RsInstanceRemovedStatus, RsInstanceReclamation},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return instance, RsInstanceSuccessStatus, nil
				}
				return nil, "", fmt.Errorf("[ERROR] Get the resource instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance '%s' deletion failed: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}

func FilterDeployments(deployments []models.ServiceDeployment, location string) ([]models.ServiceDeployment, map[string]bool) {
	supportedDeployments := []models.ServiceDeployment{}
	supportedLocations := make(map[string]bool)
	for _, d := range deployments {
		if d.Metadata.RCCompatible {
			deploymentLocation := d.Metadata.Deployment.Location
			supportedLocations[deploymentLocation] = true
			if deploymentLocation == location {
				supportedDeployments = append(supportedDeployments, d)
			}
		}
	}
	return supportedDeployments, supportedLocations
}

func flattenPlanHistory(keys []rc.PlanHistoryItem) []interface{} {
	var out = make([]interface{}, len(keys))
	for i, k := range keys {
		m := make(map[string]interface{})
		m["resource_plan_id"] = k.ResourcePlanID
		m["start_date"] = k.StartDate.String()
		out[i] = m
	}
	return out
}
