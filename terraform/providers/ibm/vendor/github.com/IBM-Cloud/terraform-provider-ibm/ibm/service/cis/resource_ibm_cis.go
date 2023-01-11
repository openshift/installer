// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	CisInstanceSuccessStatus      = "active"
	CisInstanceProgressStatus     = "in progress"
	cisInstanceProvisioningStatus = "provisioning"
	CisInstanceInactiveStatus     = "inactive"
	CisInstanceFailStatus         = "failed"
	CisInstanceRemovedStatus      = "removed"
	cisInstanceReclamation        = "pending_reclamation"
)

func ResourceIBMCISInstance() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISInstanceCreate,
		Read:     ResourceIBMCISInstanceRead,
		Update:   ResourceIBMCISInstanceUpdate,
		Delete:   ResourceIBMCISInstanceDelete,
		Exists:   ResourceIBMCISInstanceExists,
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
				Computed:    true,
				Description: "The name of the Cloud Internet Services offering",
			},

			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plan type of the service",
			},

			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of resource instance",
			},

			"location": {
				Description: "The location where the instance available",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},

			"resource_group_id": {
				Description: "The resource group id",
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Computed:    true,
			},

			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Arbitrary parameters to pass. Must be a JSON object",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_cis", "tags")},
				Set:      flex.ResourceIBMVPCHash,
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of resource instance",
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

func ResourceIBMCISValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmCISResourceValidator := validate.ResourceValidator{ResourceName: "ibm_cis", Schema: validateSchema}
	return &ibmCISResourceValidator
}

// Replace with func wrapper for resourceIBMResourceInstanceCreate specifying serviceName := "internet-svcs"
func ResourceIBMCISInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	serviceName := "internet-svcs"
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
	deployments, supportedLocations := filterCISDeployments(deployments, location)

	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return fmt.Errorf("[ERROR] No deployment found for service plan %s at location %s.\nValid location(s) are: %q", plan, location, locationList)
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

	if parameters, ok := d.GetOk("parameters"); ok {
		rsInst.Parameters = parameters.(map[string]interface{})
	}

	instance, response, err := rsConClient.CreateResourceInstance(&rsInst)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating resource instance: %s %s", err, response)
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of ibm cis (%s) tags: %s", d.Id(), err)
		}
	}

	// Moved d.SetId(instance.ID) to after waiting for resource to finish creation. Otherwise Terraform initates depedent tasks too early.
	// Original flow had SetId here as its required as input to waitForCISInstanceCreate

	_, err = waitForCISInstanceCreate(d, meta, *instance.ID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for create resource instance (%s) to be succeeded: %s", d.Id(), err)
	}

	d.SetId(*instance.ID)

	return ResourceIBMCISInstanceRead(d, meta)
}

func ResourceIBMCISInstanceRead(d *schema.ResourceData, meta interface{}) error {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	instanceID := d.Id()
	rsInst := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, response, err := rsConClient.GetResourceInstance(&rsInst)
	if err != nil {
		if strings.Contains(err.Error(), "Object not found") ||
			strings.Contains(err.Error(), "status code: 404") {
			log.Printf("[WARN] Removing record from state because it's not found via the API")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error retrieving resource instance: %s %s", err, response)
	}
	if strings.Contains(*instance.State, "removed") {
		log.Printf("[WARN] Removing instance from TF state because it's now in removed state")
		d.SetId("")
		return nil
	}
	tags, err := flex.GetTagsUsingCRN(meta, *instance.CRN)
	if err != nil {
		log.Printf(
			"Error on get of ibm cis tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)
	d.Set("name", *instance.Name)
	d.Set("status", *instance.State)
	d.Set("resource_group_id", *instance.ResourceGroupID)
	d.Set("parameters", flex.Flatten(instance.Parameters))
	if instance.CRN != nil {
		location := strings.Split(*instance.CRN, ":")
		if len(location) > 5 {
			d.Set("location", location[5])
		}
	}
	d.Set("guid", *instance.GUID)

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	d.Set("service", "internet-svcs")

	servicePlan, err := rsCatRepo.GetServicePlanName(*instance.ResourcePlanID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	d.Set("plan", servicePlan)

	d.Set(flex.ResourceName, *instance.Name)
	d.Set(flex.ResourceCRN, *instance.CRN)
	d.Set(flex.ResourceStatus, *instance.State)
	d.Set(flex.ResourceGroupName, *instance.ResourceGroupCRN)

	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, rcontroller+"/internet-svcs/"+url.QueryEscape(*instance.CRN))

	return nil
}

func ResourceIBMCISInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	instanceID := d.Id()

	updateReq := rc.UpdateResourceInstanceOptions{
		ID: &instanceID,
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateReq.Name = &name
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

		updateReq.ResourcePlanID = &servicePlan

	}

	if d.HasChange("tags") {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, instanceID)
		if err != nil {
			log.Printf(
				"Error on update of CIS (%s) tags: %s", d.Id(), err)
		}
	}

	_, response, err := rsConClient.UpdateResourceInstance(&updateReq)
	if err != nil {
		return fmt.Errorf("[ERROR] Error updating resource instance: %s %s", err, response)
	}

	_, err = waitForCISInstanceUpdate(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for update resource instance (%s) to be succeeded: %s", d.Id(), err)
	}

	return ResourceIBMCISInstanceRead(d, meta)
}

func ResourceIBMCISInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	id := d.Id()
	recursive := true
	deleteReq := rc.DeleteResourceInstanceOptions{
		ID:        &id,
		Recursive: &recursive,
	}
	response, err := rsConClient.DeleteResourceInstance(&deleteReq)
	if err != nil {
		// If prior delete occurs, instance is not immediately deleted, but remains in "removed" state"
		// RC 410 with "Gone" returned as error
		if strings.Contains(err.Error(), "Gone") ||
			strings.Contains(err.Error(), "status code: 410") {
			log.Printf("[WARN] Resource instance already deleted %s\n %s", err, response)
			err = nil
		} else {
			return fmt.Errorf("[ERROR] Error deleting resource instance: %s %s", err, response)
		}
	}

	_, err = waitForCISInstanceDelete(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for resource instance (%s) to be deleted: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
func ResourceIBMCISInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	rsInst := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, response, err := rsConClient.GetResourceInstance(&rsInst)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting cis instance: %s %s", err, response)
	}
	if instance != nil && (strings.Contains(*instance.State, "removed") || strings.Contains(*instance.State, cisInstanceReclamation)) {
		log.Printf("[WARN] Removing instance from state because it's in removed or pending_reclamation state")
		d.SetId("")
		return false, nil
	}

	return *instance.ID == instanceID, nil
}

func waitForCISInstanceCreate(d *schema.ResourceData, meta interface{}, instanceID string) (interface{}, error) {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	//instanceID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{CisInstanceProgressStatus, CisInstanceInactiveStatus, cisInstanceProvisioningStatus},
		Target:  []string{CisInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			rsInst := rc.GetResourceInstanceOptions{
				ID: &instanceID,
			}
			instance, response, err := rsConClient.GetResourceInstance(&rsInst)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %v %s", d.Id(), err, response)
				}
				return nil, "", err
			}
			if *instance.State == CisInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance %s failed: %v %s", d.Id(), err, response)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForCISInstanceUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{CisInstanceProgressStatus, CisInstanceInactiveStatus},
		Target:  []string{CisInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			rsInst := rc.GetResourceInstanceOptions{
				ID: &instanceID,
			}
			instance, response, err := rsConClient.GetResourceInstance(&rsInst)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %v %s", d.Id(), err, response)
				}
				return nil, "", err
			}
			if *instance.State == CisInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance %s failed: %v %s", d.Id(), err, response)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForCISInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending: []string{CisInstanceProgressStatus, CisInstanceInactiveStatus, CisInstanceSuccessStatus},
		Target:  []string{CisInstanceRemovedStatus, cisInstanceReclamation},
		Refresh: func() (interface{}, string, error) {
			rsInst := rc.GetResourceInstanceOptions{
				ID: &instanceID,
			}
			instance, response, err := rsConClient.GetResourceInstance(&rsInst)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return instance, CisInstanceSuccessStatus, nil
				}
				return nil, "", err
			}
			if *instance.State == CisInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance %s failed to delete: %v %s", d.Id(), err, response)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func filterCISDeployments(deployments []models.ServiceDeployment, location string) ([]models.ServiceDeployment, map[string]bool) {
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
