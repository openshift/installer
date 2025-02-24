// Copyright IBM Corp. 2017, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	satLocation = "location"
	sateLocZone = "managed_from"

	isLocationDeleting     = "deleting"
	isLocationDeleteDone   = "done"
	isLocationDeploying    = "deploying"
	isLocationReady        = "action required"
	isLocationDeployFailed = "deploy_failed"
)

func ResourceIBMSatelliteLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMSatelliteLocationCreate,
		Read:   resourceIBMSatelliteLocationRead,
		Update: resourceIBMSatelliteLocationUpdate,
		Delete: resourceIBMSatelliteLocationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				ID := d.Id()
				satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
				if err != nil {
					return nil, err
				}

				getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
					Controller: &ID,
				}

				instance, response, err := satClient.GetSatelliteLocation(getSatLocOptions)
				if err != nil || instance == nil {
					return nil, fmt.Errorf("Error reading satellite location: %s\n%s", err, response)
				}

				d.Set("zones", flex.NewStringSet(schema.HashString, instance.WorkerZones))
				return []*schema.ResourceData{d}, nil
			},
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ImmutableResourceCustomizeDiff([]string{satLocation, sateLocZone, "resource_group_id", "zones"}, diff)
			},
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			satLocation: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique name for the new Satellite location",
			},
			sateLocZone: {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					} else if o[0:3] == n[0:3] {
						return true
					}
					return o == n
				},
				Description: "The IBM Cloud metro from which the Satellite location is managed",
			},
			"physical_address": {
				Type:             schema.TypeString,
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Description:      "An optional physical address of the new Satellite location which is deployed on premise",
			},
			"capabilities": {
				Type:             schema.TypeSet,
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				RequiredWith:     []string{"physical_address"},
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				Description:      "The satellite capabilities attached to the location",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the new Satellite location",
			},
			"coreos_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable Red Hat CoreOS features within the Satellite location",
			},
			"logging_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The account ID for IBM Log Analysis with LogDNA log forwarding",
			},
			"cos_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "COSBucket - IBM Cloud Object Storage bucket configuration details",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cos_credentials": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "COSAuthorization - IBM Cloud Object Storage authorization keys",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HMAC secret access key ID",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HMAC secret access key",
						},
					},
				},
			},
			"zones": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The names of at least three high availability zones to use for the location",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the resource group.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_satellite_location", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags associated with resource instance",
			},
			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the resource group",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location CRN",
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
			"service_subnet": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Custom subnet CIDR to provide private IP addresses for services",
				DiffSuppressFunc: flex.ApplyOnce,
			},
			"pod_subnet": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Custom subnet CIDR to provide private IP addresses for pods",
				DiffSuppressFunc: flex.ApplyOnce,
			},
		},
	}
}

func ResourceIBMSatelliteLocationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "physical_address",
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Optional:                   true,
			MinValueLength:             0,
			MaxValueLength:             400,
		},
		validate.ValidateSchema{
			Identifier:                 "capabilities",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "on-prem",
		},
	)

	ibmSatelliteLocationValidator := validate.ResourceValidator{ResourceName: "ibm_satellite_location", Schema: validateSchema}
	return &ibmSatelliteLocationValidator
}

func resourceIBMSatelliteLocationCreate(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	createSatLocOptions := &kubernetesserviceapiv1.CreateSatelliteLocationOptions{}
	satLocation := d.Get(satLocation).(string)
	createSatLocOptions.Name = &satLocation
	sateLocZone := d.Get(sateLocZone).(string)
	createSatLocOptions.Location = &sateLocZone

	if v, ok := d.GetOk("coreos_enabled"); ok {
		coreosEnabled := v.(bool)
		createSatLocOptions.CoreosEnabled = coreosEnabled
	}
	if v, ok := d.GetOk("cos_config"); ok {
		createSatLocOptions.CosConfig = flex.ExpandCosConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("cos_credentials"); ok {
		createSatLocOptions.CosCredentials = flex.ExpandCosCredentials(v.([]interface{}))
	}

	if v, ok := d.GetOk("logging_account_id"); ok {
		logAccID := v.(string)
		createSatLocOptions.LoggingAccountID = &logAccID
	}

	if v, ok := d.GetOk("physical_address"); ok {
		addr := v.(string)
		createSatLocOptions.PhysicalAddress = &addr
	}

	if v, ok := d.GetOk("capabilities"); ok {
		z := v.(*schema.Set)
		createSatLocOptions.CapabilitiesManagedBySatellite = flex.FlattenSatelliteCapabilities(z)
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		createSatLocOptions.Description = &desc
	}

	if v, ok := d.GetOk("zones"); ok {
		z := v.(*schema.Set)
		createSatLocOptions.Zones = flex.FlattenSatelliteZones(z)
	}

	if v, ok := d.GetOk("resource_group_id"); ok && v != nil {
		pathParamsMap := map[string]string{
			"X-Auth-Resource-Group": v.(string),
		}
		createSatLocOptions.Headers = pathParamsMap
	}

	if v, ok := d.GetOk("pod_subnet"); ok {
		podSubnet := v.(string)
		createSatLocOptions.PodSubnet = &podSubnet
	}

	if v, ok := d.GetOk("service_subnet"); ok {
		serviceSubnet := v.(string)
		createSatLocOptions.ServiceSubnet = &serviceSubnet
	}

	instance, response, err := satClient.CreateSatelliteLocation(createSatLocOptions)
	if err != nil || instance == nil {
		return fmt.Errorf("[ERROR] Error Creating Satellite Location: %s\n%s", err, response)
	}

	d.SetId(*instance.ID)
	log.Printf("[INFO] Created satellite location : %s", satLocation)

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on create of ibm satellite location (%s) tags: %s", d.Id(), err)
		}
	}

	//Wait for location to be in ready state
	_, err = waitForLocationToReady(*instance.ID, d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for location (%s) to reach action required state: %s", *instance.ID, err)
	}

	return resourceIBMSatelliteLocationRead(d, meta)
}

func resourceIBMSatelliteLocationRead(d *schema.ResourceData, meta interface{}) error {
	ID := d.Id()
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: &ID,
	}

	instance, response, err := satClient.GetSatelliteLocation(getSatLocOptions)
	if err != nil || instance == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set(satLocation, *instance.Name)
	if instance.Description != nil {
		d.Set("description", *instance.Description)
	}

	if instance.PhysicalAddress != nil {
		d.Set("physical_address", *instance.PhysicalAddress)
	}

	if instance.CapabilitiesManagedBySatellite != nil {
		d.Set("capabilities", instance.CapabilitiesManagedBySatellite)
	}

	if instance.CoreosEnabled != nil {
		d.Set("coreos_enabled", *instance.CoreosEnabled)
	}

	if instance.Datacenter != nil {
		d.Set(sateLocZone, *instance.Datacenter)
	}

	if instance.ResourceGroup != nil {
		d.Set("resource_group_id", instance.ResourceGroup)
	}

	tags, err := flex.GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"Error on get of ibm satellite location tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)
	d.Set("crn", *instance.Crn)
	d.Set(flex.ResourceGroupName, *instance.ResourceGroupName)
	if instance.Hosts != nil {
		d.Set("host_attached_count", *instance.Hosts.Total)
		d.Set("host_available_count", *instance.Hosts.Available)
	}
	d.Set("created_on", *instance.CreatedDate)
	if instance.Ingress != nil {
		d.Set("ingress_hostname", *instance.Ingress.Hostname)
		d.Set("ingress_secret", *instance.Ingress.SecretName)
	}

	if instance.PodSubnet != nil {
		d.Set("pod_subnet", *instance.PodSubnet)
	}

	if instance.ServiceSubnet != nil {
		d.Set("service_subnet", *instance.ServiceSubnet)
	}

	return nil
}

func resourceIBMSatelliteLocationUpdate(d *schema.ResourceData, meta interface{}) error {
	ID := d.Id()
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if d.HasChange("tags") || v != "" {
		oldList, newList := d.GetChange("tags")
		getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
			Controller: &ID,
		}

		instance, response, err := satClient.GetSatelliteLocation(getSatLocOptions)
		if err != nil || instance == nil {
			return fmt.Errorf("[ERROR] Error retrieving satellite location: %s\n%s", err, response)
		}
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"An error occured during update of instance (%s) tags: %s", ID, err)
		}
	}
	return nil
}

func resourceIBMSatelliteLocationDelete(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	removeSatLocOptions := &kubernetesserviceapiv1.RemoveSatelliteLocationOptions{}
	name := d.Get(satLocation).(string)
	removeSatLocOptions.Controller = &name

	response, err := satClient.RemoveSatelliteLocation(removeSatLocOptions)
	if err != nil && response.StatusCode != 404 {
		return fmt.Errorf("[ERROR] Error Deleting Satellite Location: %s\n%s", err, response)
	}

	//Wait for location to delete
	_, err = waitForLocationDelete(name, d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for deleting location instance: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForLocationDelete(location string, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isLocationDeleting, ""},
		Target:  []string{isLocationDeleteDone},
		Refresh: func() (interface{}, string, error) {
			getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationsOptions{}
			locations, response, err := satClient.GetSatelliteLocations(getSatLocOptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error Getting locations list to delete : %s\n%s", err, response)
			}

			isExist := false
			if locations != nil {
				for _, loc := range locations {
					if *loc.ID == location || *loc.Name == location {
						isExist = true
						return "", isLocationDeleting, nil
					}
				}
				if !isExist {
					return location, isLocationDeleteDone, nil
				}
			}
			return nil, "", fmt.Errorf("[ERROR] Failed to delete location : %s\n%s", err, response)
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForLocationToReady(loc string, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isLocationDeploying},
		Target:  []string{isLocationReady, isLocationDeployFailed},
		Refresh: func() (interface{}, string, error) {
			getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
				Controller: flex.PtrToString(loc),
			}

			var location *kubernetesserviceapiv1.MultishiftGetController
			var response *core.DetailedResponse
			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				location, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
				if err != nil || location == nil {
					if response != nil && response.StatusCode == 404 {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if conns.IsResourceTimeoutError(err) {
				location, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
			}

			if location != nil && *location.State == isLocationDeployFailed {
				return location, isLocationDeployFailed, fmt.Errorf("[ERROR] The location is in failed state: %s", d.Id())
			}

			if location != nil && *location.State == isLocationReady {
				return location, isLocationReady, nil
			}
			return location, isLocationDeploying, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}
