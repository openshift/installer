// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package hpcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/ibm-hpcs-tke-sdk/tkesdk"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMHPCS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMHPCSCreate,
		ReadContext:   resourceIBMHPCSRead,
		UpdateContext: resourceIBMHPCSUpdate,
		DeleteContext: resourceIBMHPCSDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ImmutableResourceCustomizeDiff([]string{"units", "failover_units", "location", "resource_group_id", "service"}, diff)
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name for the HPCS instance",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plan type of the HPCS Instance",
			},
			"location": {
				Description: "The location where the HPCS instance available",
				Required:    true,
				Type:        schema.TypeString,
			},
			"units": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of operational crypto units for your service instance",
			},
			"failover_units": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of failover crypto units for your service instance",
			},
			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "hs-crypto",
				Description: "The name of the service offering `hs-crypto` ",
			},
			"resource_group_id": {
				Description: "The resource group id",
				Optional:    true,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_endpoints": {
				Description:  "Types of the service endpoints. Possible values are `public-and-private`, `private-only`.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public-and-private", "private-only"}),
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_hpcs", "tags")},
				Set:      flex.ResourceIBMVPCHash,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of HPCS instance",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of HPCS instance",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Guid of HPCS instance",
			},
			"dashboard_url": {
				Description: "Dashboard URL to access resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "The current state of the instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"resource_aliases_url": {
				Description: "The relative path to the resource aliases for the instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"resource_bindings_url": {
				Description: "The relative path to the resource bindings for the instance.",
				Type:        schema.TypeString,
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
			"extensions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The extended metadata as a map associated with the HPCS instance.",
			},
			"signature_server_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of signing service",
			},
			"signature_threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Signature Threshold Value",
			},
			"revocation_threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Revocation Threshold Value",
			},
			"admins": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Crypto Unit Administrators",
				Set:         resourceIBMHPCSAdminHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Admin Name",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The administrator signature key",
						},
						"token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Credential giving access to the administrator signature key",
						},
					},
				},
			},
			"hsm_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "HSM Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"signature_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Signature Threshold Value",
						},
						"revocation_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Revocation Threshold Value",
						},
						"admins": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Crypto Unit Administrators",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Admin Name",
									},
									"ski": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hsm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hsm_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hsm_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_mk_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_mk_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_mkvp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_mkvp": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type HPCSParams struct {
	Units                 int    `json:"units,omitempty"`
	FailoverUnits         int    `json:"failover_units,omitempty"`
	RequiresRecoveryUnits bool   `json:"requires_recovery_units,omitempty"`
	ServiceEndpoints      string `json:"allowed_network,omitempty"`
}

func ResourceIBMHPCSValidator() *validate.ResourceValidator {
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

	ibmResourceInstanceResourceValidator := validate.ResourceValidator{ResourceName: "ibm_hpcs", Schema: validateSchema}
	return &ibmResourceInstanceResourceValidator
}

func resourceIBMHPCSCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	rsInst := rc.CreateResourceInstanceOptions{
		Name: &name,
	}
	// Fetch Service Plan ID using Global Catalog APIs
	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.FindByName(serviceName, true)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving service offering: %s", err))
	}

	if metadata, ok := serviceOff[0].Metadata.(*models.ServiceResourceMetadata); ok {
		if !metadata.Service.RCProvisionable {
			return diag.FromErr(fmt.Errorf("[ERROR] %s cannot be provisioned by resource controller", serviceName))
		}
	} else {
		return diag.FromErr(fmt.Errorf("[ERROR] Cannot create instance of resource %s", serviceName))
	}

	servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving plan: %s", err))
	}
	rsInst.ResourcePlanID = &servicePlan

	// Fetch Catalog CRN using Global Catalog APIs
	deployments, err := rsCatRepo.ListDeployments(servicePlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving deployment for plan %s : %s", plan, err))
	}
	if len(deployments) == 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] No deployment found for service plan : %s", plan))
	}
	deployments, supportedLocations := resourcecontroller.FilterDeployments(deployments, location)
	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return diag.FromErr(fmt.Errorf("[ERROR] No deployment found for service plan %s at location %s.\n valid location(s) are: %q", plan, location, locationList))
	}
	rsInst.Target = &deployments[0].CatalogCRN

	// Get Resource Group ID from User.. If not provided fetch Default resource group ID of the account
	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rg := rsGrpID.(string)
		rsInst.ResourceGroup = &rg
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return diag.FromErr(err)
		}
		rsInst.ResourceGroup = &defaultRg
	}

	// Resource Controller Parameters for HPCS Service
	params := HPCSParams{}
	if units, ok := d.GetOk("units"); ok {
		params.Units = units.(int)
	}
	if failover_units, ok := d.GetOk("failover_units"); ok {
		params.FailoverUnits = failover_units.(int)
	}
	if serviceEndpoint, ok := d.GetOk("service_endpoints"); ok {
		params.ServiceEndpoints = serviceEndpoint.(string)
	}
	params.RequiresRecoveryUnits = true
	// Convert HPCSParams srtuct to map
	parameters, _ := json.Marshal(params)
	var raw map[string]interface{}
	json.Unmarshal(parameters, &raw)
	rsInst.Parameters = raw

	// Create HPCS Instance
	instance, resp, err := rsConClient.CreateResourceInstance(&rsInst)
	if err != nil || instance == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error when creating HPCS instance: %s with resp code: %s", err, resp))
	}
	d.SetId(*instance.ID)                       // Set Resource ID
	_, err = waitForHPCSInstanceCreate(d, meta) // Wait for Instance to be available
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"[ERROR] Error waiting for create HPCS instance (%s) to be succeeded: %s", d.Id(), err))
	}

	// Update Tags for this Resource using Global Tagging APIs
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of HPCS instance (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceIBMHPCSUpdate(context, d, meta)
}
func resourceIBMHPCSRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instance == nil {
		if resp != nil && resp.StatusCode == 404 { // If instance doesnt not exists, empty the statefile
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving HPCS instance: %s with resp code: %s", err, resp))
	}
	if instance != nil && (strings.Contains(*instance.State, "removed") || strings.Contains(*instance.State, resourcecontroller.RsInstanceReclamation)) {
		log.Printf("[WARN] Removing instance from state because it's in removed or pending_reclamation state")
		d.SetId("")
		return nil
	}

	d.Set("name", instance.Name)
	d.Set("status", instance.State)
	d.Set("resource_group_id", instance.ResourceGroupID)
	d.Set("crn", instance.CRN)
	d.Set("dashboard_url", instance.DashboardURL)
	d.Set("guid", instance.GUID)
	d.Set("resource_keys_url", instance.ResourceKeysURL)
	d.Set("resource_bindings_url", instance.ResourceBindingsURL)
	d.Set("resource_aliases_url", instance.ResourceAliasesURL)
	d.Set("state", instance.State)
	d.Set("service", strings.Split(instanceID, ":")[4])
	//Set tags
	tags, err := flex.GetTagsUsingCRN(meta, *instance.CRN)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of HPCS instance tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)
	// Set Location
	if instance.CRN != nil {
		location := strings.Split(*instance.CRN, ":")
		if len(location) > 5 {
			d.Set("location", location[5])
		}
	}
	// Set Service Plan
	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	servicePlan, err := rsCatRepo.GetServicePlanName(*instance.ResourcePlanID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving plan: %s", err))
	}
	d.Set("plan", servicePlan)
	// Set Instance parameters
	if instance.Parameters != nil {
		instanceParameters := flex.Flatten(instance.Parameters)

		if endpoint, ok := instanceParameters["allowed_network"]; ok {
			if endpoint != "private-only" {
				endpoint = "public-and-private"
			}
			d.Set("service_endpoints", endpoint)
		} else {
			d.Set("service_endpoints", "public-and-private")
		}
		if u, ok := instanceParameters["units"]; ok {
			units, err := strconv.Atoi(u)
			if err != nil {
				log.Println("[ERROR] Error converting units from string to integer")
			}
			d.Set("units", units)
		}
		if f, ok := instanceParameters["failover_units"]; ok {
			failover_units, err := strconv.Atoi(f)
			if err != nil {
				log.Println("[ERROR] Error failover_units units from string to integer")
			}
			d.Set("failover_units", failover_units)
		}
	}
	// Set Extensions
	if len(instance.Extensions) == 0 {
		d.Set("extensions", instance.Extensions)
	} else {
		d.Set("extensions", flex.Flatten(instance.Extensions))
	}
	d.Set("restored_by", instance.RestoredBy)
	d.Set("scheduled_reclaim_by", instance.ScheduledReclaimBy)
	d.Set("deleted_by", instance.DeletedBy)
	d.Set("update_by", instance.UpdatedBy)
	d.Set("created_by", instance.CreatedBy)
	if instance.RestoredAt != nil {
		d.Set("restored_at", instance.RestoredAt.String())
	}
	if instance.ScheduledReclaimAt != nil {
		d.Set("scheduled_reclaim_at", instance.ScheduledReclaimAt.String())
	}
	if instance.ScheduledReclaimAt != nil {
		d.Set("deleted_at", instance.DeletedAt.String())
	}
	if instance.UpdatedAt != nil {
		d.Set("update_at", instance.UpdatedAt.String())
	}
	if instance.CreatedAt != nil {
		d.Set("created_at", instance.CreatedAt.String())
	}
	// Bluemix Session to get Oauth tokens
	ci, err := hsmClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	ci.InstanceId = *instance.GUID

	hsmInfo, err := tkesdk.Query(ci)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Quering HSM config: %s", err))
	}
	d.Set("hsm_info", FlattenHSMInfo(hsmInfo))

	if validateHSM(hsmInfo) && !d.IsNewResource() {
		d.Set("admins", nil)
		d.Set("signature_threshold", nil)
		d.Set("revocation_threshold", nil)
	}

	return nil
}

func resourceIBMHPCSUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Id()
	resourceInstanceUpdate := rc.UpdateResourceInstanceOptions{
		ID: &instanceID,
	}
	update := false
	if d.HasChange("name") {
		name := d.Get("name").(string)
		resourceInstanceUpdate.Name = &name
		update = true
	}
	if d.HasChange("plan") {
		plan := d.Get("plan").(string)
		service := d.Get("service").(string)
		rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
		if err != nil {
			return diag.FromErr(err)
		}
		rsCatRepo := rsCatClient.ResourceCatalog()

		serviceOff, err := rsCatRepo.FindByName(service, true)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving service offering: %s", err))
		}

		servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving plan: %s", err))
		}

		resourceInstanceUpdate.ResourcePlanID = &servicePlan
		update = true

	}
	if d.HasChange("service_endpoints") {
		params := HPCSParams{}
		params.ServiceEndpoints = d.Get("service_endpoints").(string)
		parameters, _ := json.Marshal(params)
		var raw map[string]interface{}
		json.Unmarshal(parameters, &raw)
		resourceInstanceUpdate.Parameters = raw
		update = true
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting HPCS instance: %s with resp code: %s", err, resp))
	}
	if d.HasChange("tags") {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"[ERROR] Error on update of HPCS instance (%s) tags: %s", d.Id(), err)
		}
	}
	if update && !d.IsNewResource() { // Update RC API only if its not a new resource
		_, resp, err = rsConClient.UpdateResourceInstance(&resourceInstanceUpdate)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating HPCS instance: %s with resp code: %s", err, resp))
		}

		_, err = waitForHPCSInstanceUpdate(d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for update HPCS instance (%s) to be succeeded: %s", d.Id(), err))
		}
	}
	// Initialise HPCS Crypto Units

	if d.HasChange("signature_threshold") || d.HasChange("revocation_threshold") || d.HasChange("admins") || d.HasChange("signature_server_url") {
		if url, ok := d.GetOk("signature_server_url"); ok {
			serverURL := url.(string)
			err := os.Setenv("TKE_SIGNSERV_URL", serverURL)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		hsm_config := expandHSMConfig(d, meta)
		// Bluemix Session to get Oauth tokens
		ci, err := hsmClient(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		ci.InstanceId = *instance.GUID

		// Check Transitions
		problems, err := tkesdk.CheckTransition(ci, hsm_config)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Checking Transitions: %s", err))
		}
		if len(problems) != 0 {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Checking Transitions: %v", problems))
		}
		// Update / Initialize Crypto Units
		hsmDetails, err := tkesdk.Update(ci, hsm_config)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Updating Crypto Units: %s", err))
		}
		if len(hsmDetails) != 0 {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Updating Crypto Units..One or more problems were found during initial checks: %v", hsmDetails))
		}
	}
	return resourceIBMHPCSRead(context, d, meta)
}
func expandHSMConfig(d *schema.ResourceData, meta interface{}) tkesdk.HsmConfig {
	hsmConfig := tkesdk.HsmConfig{}
	if s, ok := d.GetOk("signature_threshold"); ok {
		hsmConfig.SignatureThreshold = s.(int)
	}
	if r, ok := d.GetOk("revocation_threshold"); ok {
		hsmConfig.RevocationThreshold = r.(int)
	}
	if a, ok := d.GetOk("admins"); ok {
		ads := a.(*schema.Set).List()
		admins := []tkesdk.AdminInfo{}
		for _, a := range ads {
			ad := a.(map[string]interface{})
			admin := tkesdk.AdminInfo{
				Name:  ad["name"].(string),
				Key:   ad["key"].(string),
				Token: ad["token"].(string),
			}
			admins = append(admins, admin)
		}
		hsmConfig.Admins = admins
	}
	return hsmConfig
}
func resourceIBMHPCSDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &id,
	}
	instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting HPCS instance: %s with resp code: %s", err, resp))
	}
	// Bluemix Session to get Oauth tokens
	ci, err := hsmClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	ci.InstanceId = *instance.GUID
	if url, ok := d.GetOk("signature_server_url"); ok {
		serverURL := url.(string)
		err := os.Setenv("TKE_SIGNSERV_URL", serverURL)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// Zeroize Crypto Units
	hsm := expandHSMConfig(d, meta)
	err = tkesdk.Zeroize(ci, hsm)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Zeroizing Crypto Units: %s", err))
	}

	// Delete Instance
	recursive := true
	resourceInstanceDelete := rc.DeleteResourceInstanceOptions{
		ID:        &id,
		Recursive: &recursive,
	}
	resp, error := rsConClient.DeleteResourceInstance(&resourceInstanceDelete)
	if error != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error deleting HPCS instance: %s with resp code: %s", error, resp))
	}
	_, err = waitForHPCSInstanceDelete(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"[ERROR] Error waiting for HPCS instance (%s) to be deleted: %s", d.Id(), err))
	}

	d.SetId("")

	return nil
}
func waitForHPCSInstanceCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{resourcecontroller.RsInstanceProgressStatus, resourcecontroller.RsInstanceInactiveStatus, resourcecontroller.RsInstanceProvisioningStatus},
		Target:  []string{resourcecontroller.RsInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] HPCS instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", fmt.Errorf("[ERROR] Get on HPCS instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == resourcecontroller.RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The status of HPCS instance %s failed: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForHPCSInstanceUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{resourcecontroller.RsInstanceProgressStatus, resourcecontroller.RsInstanceInactiveStatus},
		Target:  []string{resourcecontroller.RsInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] HPCS instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", fmt.Errorf("[ERROR] Get the HPCS instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == resourcecontroller.RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The status of HPCS instance %s failed: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForHPCSInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	stateConf := &resource.StateChangeConf{
		Pending: []string{resourcecontroller.RsInstanceProgressStatus, resourcecontroller.RsInstanceInactiveStatus, resourcecontroller.RsInstanceSuccessStatus},
		Target:  []string{resourcecontroller.RsInstanceRemovedStatus, resourcecontroller.RsInstanceReclamation},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return instance, resourcecontroller.RsInstanceSuccessStatus, nil
				}
				return nil, "", fmt.Errorf("[ERROR] Get on HPCS instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == resourcecontroller.RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] HPCS instance %s failed to delete: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func resourceIBMHPCSAdminHash(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", a["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", a["key"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", a["token"].(string)))

	return conns.String(buf.String())
}
func validateHSM(hsmInfo []tkesdk.HsmInfo) bool {
	update := false
	if len(hsmInfo) == 0 {
		return true
	}
	for _, hsm := range hsmInfo {
		if hsm.CurrentMKStatus != "Valid" {
			update = true
		}
	}
	return update
}
func hsmClient(d *schema.ResourceData, meta interface{}) (tkesdk.CommonInputs, error) {
	ci := tkesdk.CommonInputs{}
	// Bluemix Session to get Oauth tokens
	bluemixSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return ci, err
	}
	err = conns.RefreshToken(bluemixSession)
	if err != nil {
		return ci, fmt.Errorf("[ERROR] Error Refreshing Authentication Token: %s", err)
	}
	var serviceEndpoint string
	if e, ok := d.GetOk("service_endpoints"); ok {
		serviceEndpoint = e.(string)
	}
	ci.Region = d.Get("location").(string)
	ci.ApiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_HPCS_TKE_ENDPOINT"}, "cloud.ibm.com")
	if bluemixSession.Config.Visibility == "private" || bluemixSession.Config.Visibility == "public-and-private" || serviceEndpoint == "private-only" {
		ci.ApiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_HPCS_TKE_ENDPOINT"}, "private.cloud.ibm.com")
	}

	ci.AuthToken = bluemixSession.Config.IAMAccessToken

	return ci, err
}
