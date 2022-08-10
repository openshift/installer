// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudant

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/IBM-Cloud/terraform-provider-ibm/version"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/go-sdk-core/v5/core"
	iamidentity "github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func ResourceIBMCloudant() *schema.Resource {
	riSchema := resourcecontroller.ResourceIBMResourceInstance().Schema

	riSchema["service"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The service type of the instance",
		Computed:    true,
	}

	riSchema["legacy_credentials"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Use both legacy credentials and IAM for authentication",
		ForceNew:    true,
	}

	riSchema["environment_crn"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "CRN of the IBM Cloudant Dedicated Hardware plan instance",
		ForceNew:    true,
	}

	riSchema["include_data_events"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Include data event types in events sent to IBM Cloud Activity Tracker with LogDNA for the IBM Cloudant instance. By default only emitted events are of \"management\" type.",
	}

	riSchema["capacity"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "A number of blocks of throughput units. A block consists of 100 reads/sec, 50 writes/sec, and 5 global queries/sec of provisioned throughput capacity.",
		ValidateFunc: validation.IntAtLeast(1),
	}

	riSchema["throughput"] = &schema.Schema{
		Type:        schema.TypeMap,
		Computed:    true,
		Description: "Schema for detailed information about throughput capacity with breakdown by specific throughput requests classes.",
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}

	riSchema["enable_cors"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Boolean value to turn CORS on and off.",
	}

	riSchema["cors_config"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// suppress missing cors_config, we set defaults during update
			return k == "cors_config.#" && new == "0"
		},
		Description: "Configuration for CORS.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow_credentials": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Boolean value to allow authentication credentials. If set to true, browser requests must be done by using withCredentials = true.",
				},
				"origins": {
					Type:        schema.TypeList,
					Required:    true,
					Description: "An array of strings that contain allowed origin domains. You have to specify the full URL including the protocol. It is recommended that only the HTTPS protocol is used. Subdomains count as separate domains, so you have to specify all subdomains used.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
		MinItems: 1,
		MaxItems: 1,
	}

	return &schema.Resource{
		Create:   resourceIBMCloudantCreate,
		Read:     resourceIBMCloudantRead,
		Update:   resourceIBMCloudantUpdate,
		Delete:   resourcecontroller.ResourceIBMResourceInstanceDelete,
		Exists:   resourcecontroller.ResourceIBMResourceInstanceExists,
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

		Schema: riSchema,
	}
}

func resourceIBMCloudantCreate(d *schema.ResourceData, meta interface{}) error {
	d.Set("service", "cloudantnosqldb")

	err := validateCloudantInstanceCapacity(d)
	if err != nil {
		return err
	}

	err = validateCloudantInstanceCors(d)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})

	if legacyCredentials, ok := d.GetOkExists("legacy_credentials"); ok {
		params["legacyCredentials"] = fmt.Sprintf("%t", legacyCredentials)
	}

	if environmentCRN, ok := d.GetOk("environment_crn"); ok {
		params["environment_crn"] = environmentCRN
	}

	// copy values from "parameters" to params, unless they are already defined
	parameters, ok := d.GetOk("parameters")
	if ok {
		temp := parameters.(map[string]interface{})
		for k, v := range temp {
			if override, ok := params[k]; ok && override != v {
				log.Printf("[WARN] Overriding %q in 'parameters' to %s", k, override)
				continue
			}
			params[k] = v
		}
	}

	if len(params) > 0 {
		d.Set("parameters", params)
	}

	err = resourcecontroller.ResourceIBMResourceInstanceCreate(d, meta)
	if err != nil {
		return err
	}

	// return original parameters on state
	d.Set("parameters", parameters)

	client, err := getCloudantClient(d, meta)
	if err != nil {
		return err
	}

	// if matches an instance creation default skip request
	if d.Get("include_data_events").(bool) {
		err := updateCloudantActivityTrackerEvents(client, d)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating activity tracker events: %s", err)
		}
	}

	// if matches an instance creation default skip request
	if d.Get("capacity").(int) > 1 {
		err := updateCloudantInstanceCapacity(client, d)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving capacity throughput information: %s", err)
		}
	}

	err = updateCloudantInstanceCors(client, d)
	if err != nil {
		return fmt.Errorf("[ERROR] Error updating CORS settings: %s", err)
	}

	return resourceIBMCloudantRead(d, meta)
}

func resourceIBMCloudantRead(d *schema.ResourceData, meta interface{}) error {
	err := resourcecontroller.ResourceIBMResourceInstanceRead(d, meta)
	if err != nil {
		return err
	}

	err = setCloudantLegacyCredentials(d, meta)
	if err != nil {
		return err
	}

	err = setCloudantResourceControllerURL(d, meta)
	if err != nil {
		return err
	}

	client, err := getCloudantClient(d, meta)
	if err != nil {
		return err
	}

	err = setCloudantActivityTrackerEvents(client, d)
	if err != nil {
		return err
	}

	err = setCloudantInstanceCapacity(client, d)
	if err != nil {
		return err
	}

	err = setCloudantInstanceCors(client, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceIBMCloudantUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Set("service", "cloudantnosqldb")

	err := validateCloudantInstanceCapacity(d)
	if err != nil {
		return err
	}

	err = validateCloudantInstanceCors(d)
	if err != nil {
		return err
	}

	err = resourcecontroller.ResourceIBMResourceInstanceUpdate(d, meta)
	if err != nil {
		return err
	}

	client, err := getCloudantClient(d, meta)
	if err != nil {
		return err
	}

	if d.HasChange("include_data_events") {
		err := updateCloudantActivityTrackerEvents(client, d)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating activity tracker events: %s", err)
		}
	}

	if d.HasChange("capacity") {
		err := updateCloudantInstanceCapacity(client, d)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving capacity throughput information: %s", err)
		}
	}

	if d.HasChange("enable_cors") {
		err := updateCloudantInstanceCors(client, d)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating CORS settings: %s", err)
		}
	}

	return resourceIBMCloudantRead(d, meta)
}

func setCloudantLegacyCredentials(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instance, response, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil {
		log.Printf("[DEBUG] Error retrieving resource instance: %s\n%s", err, response)
		return err
	}

	parameters := instance.Parameters

	if legacyCredentials, ok := parameters["legacyCredentials"]; ok {
		switch v := legacyCredentials.(type) {
		case bool:
			d.Set("legacy_credentials", v)
		case string:
			d.Set("legacy_credentials", v == "true")
		default:
			d.Set("legacy_credentials", false)
		}
	}

	if environmentCRN, ok := parameters["environment_crn"]; ok {
		v := environmentCRN.(string)
		d.Set("environment_crn", v)
	}

	return nil
}

func setCloudantResourceControllerURL(d *schema.ResourceData, meta interface{}) error {
	crn := d.Get(flex.ResourceCRN).(string)
	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, rcontroller+"/services/cloudantnosqldb/"+url.QueryEscape(crn))

	return nil
}

func getCloudantClient(d *schema.ResourceData, meta interface{}) (*cloudantv1.CloudantV1, error) {

	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return nil, err
	}

	var endpoint string
	extensions := d.Get("extensions").(map[string]interface{})
	if v, ok := extensions["endpoints.public"]; ok {
		endpoint = "https://" + v.(string)
	}

	switch session.Config.Visibility {
	case "private":
		_, ok := extensions["endpoints.private"]
		if !ok {
			return nil, fmt.Errorf("[ERROR] Missing endpoints.private in extensions")
		}
		endpoint = "https://" + extensions["endpoints.private"].(string)
	case "public-and-private":
		if v, ok := extensions["endpoints.private"]; ok {
			endpoint = "https://" + v.(string)
		}
	}

	endpoint = conns.EnvFallBack([]string{"IBMCLOUD_CLOUDANT_ENDPOINT"}, endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] Missing endpoints.public in extensions")
	}

	return GetCloudantClientForUrl(endpoint, meta)
}

func GetCloudantClientForUrl(endpoint string, meta interface{}) (*cloudantv1.CloudantV1, error) {
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return nil, err
	}

	var authenticator core.Authenticator
	token := session.Config.IAMAccessToken

	if token != "" {
		token = strings.Replace(token, "Bearer ", "", -1)
		authenticator = &core.BearerTokenAuthenticator{
			BearerToken: token,
		}
	} else {
		apiKey := session.Config.BluemixAPIKey
		region := session.Config.Region
		visibility := session.Config.Visibility
		iamURL := iamidentity.DefaultServiceURL
		if visibility == "private" || visibility == "public-and-private" {
			if region == "us-south" || region == "us-east" {
				iamURL = conns.ContructEndpoint(fmt.Sprintf("private.%s.iam", region), "cloud.ibm.com")
			} else {
				iamURL = conns.ContructEndpoint("private.iam", "cloud.ibm.com")
			}
		}
		authenticator = &core.IamAuthenticator{
			ApiKey: apiKey,
			URL:    conns.EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamURL) + "/identity/token",
		}
	}

	client, err := cloudantv1.NewCloudantV1(&cloudantv1.CloudantV1Options{
		Authenticator: authenticator,
		URL:           endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error occured while configuring Cloudant service: %q", err)
	}
	client.Service.SetUserAgent("cloudant-terraform/" + version.Version)

	return client, nil
}

func setCloudantActivityTrackerEvents(client *cloudantv1.CloudantV1, d *schema.ResourceData) error {
	activityTrackerEvents, err := readCloudantActivityTrackerEvents(client)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving activity tracker events: %s", err)
	}
	if activityTrackerEvents.Types != nil {
		includeDataEvents := false
		for _, t := range activityTrackerEvents.Types {
			if t == "data" {
				includeDataEvents = true
			}
		}
		d.Set("include_data_events", includeDataEvents)
	}
	return nil
}

func readCloudantActivityTrackerEvents(client *cloudantv1.CloudantV1) (*cloudantv1.ActivityTrackerEvents, error) {
	opts := client.NewGetActivityTrackerEventsOptions()

	activityTrackerEvents, response, err := client.GetActivityTrackerEvents(opts)
	if err != nil {
		log.Printf("[DEBUG] Error retrieving activity tracker events: %s\n%s", err, response)
	}
	return activityTrackerEvents, err
}

func updateCloudantActivityTrackerEvents(client *cloudantv1.CloudantV1, d *schema.ResourceData) error {
	auditEventTypes := []string{"management"}
	if d.Get("include_data_events").(bool) {
		auditEventTypes = append(auditEventTypes, "data")
	}

	opts := client.NewPostActivityTrackerEventsOptions(auditEventTypes)

	_, response, err := client.PostActivityTrackerEvents(opts)
	if err != nil {
		log.Printf("[DEBUG] Error updating activity tracker events: %s\n%s", err, response)
	}
	return err
}

func validateCloudantInstanceCapacity(d *schema.ResourceData) error {
	plan := d.Get("plan").(string)
	capacity := d.Get("capacity").(int)
	if capacity > 1 && plan == "lite" {
		return fmt.Errorf("[ERROR] Setting capacity is not supported for your instance's plan")
	}
	return nil
}

func setCloudantInstanceCapacity(client *cloudantv1.CloudantV1, d *schema.ResourceData) error {
	capacityThroughputInformation, err := readCloudantInstanceCapacity(client)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving capacity throughput information: %s", err)
	}

	if capacityThroughputInformation.Current != nil && capacityThroughputInformation.Current.Throughput != nil {
		currentThroughput := capacityThroughputInformation.Current.Throughput
		// lite plan doesn't have "blocks" attr on broker's response
		if d.Get("plan").(string) == "lite" || currentThroughput.Blocks == nil {
			d.Set("capacity", 1)
		} else {
			blocks := int(*currentThroughput.Blocks)
			d.Set("capacity", blocks)
		}
		throughput := map[string]int{
			"query": int(*currentThroughput.Query),
			"read":  int(*currentThroughput.Read),
			"write": int(*currentThroughput.Write),
		}
		d.Set("throughput", throughput)
	}
	return nil
}

func readCloudantInstanceCapacity(client *cloudantv1.CloudantV1) (*cloudantv1.CapacityThroughputInformation, error) {
	opts := client.NewGetCapacityThroughputInformationOptions()

	capacityThroughputInformation, response, err := client.GetCapacityThroughputInformation(opts)
	if err != nil {
		log.Printf("[DEBUG] Error getting capacity throughput information: %s\n%s", err, response)
	}
	return capacityThroughputInformation, nil
}

func updateCloudantInstanceCapacity(client *cloudantv1.CloudantV1, d *schema.ResourceData) error {
	blocks := int64(d.Get("capacity").(int))

	putOpts := client.NewPutCapacityThroughputConfigurationOptions(blocks)

	_, response, err := client.PutCapacityThroughputConfiguration(putOpts)
	if err != nil {
		log.Printf("[DEBUG] Error updating capacity throughput: %s\n%s", err, response)
		return err
	}

	return isWaitForCapacityUpdated(client)
}

func isWaitForCapacityUpdated(client *cloudantv1.CloudantV1) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry"},
		Target:  []string{"done", "failed"},
		Refresh: func() (interface{}, string, error) {
			capacityThroughputInformation, err := readCloudantInstanceCapacity(client)
			if err != nil {
				return nil, "failed", err
			}

			state := "retry"
			current := *capacityThroughputInformation.Current.Throughput.Blocks
			target := *capacityThroughputInformation.Target.Throughput.Blocks

			if current == target {
				state = "done"
			}

			return current, state, nil
		},
		Timeout:    5 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 2 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func validateCloudantInstanceCors(d *schema.ResourceData) error {
	enableCors := d.Get("enable_cors").(bool)
	corsConfigRaw := d.Get("cors_config").([]interface{})
	if !enableCors && len(corsConfigRaw) > 0 {
		corsConfig := corsConfigRaw[0].(map[string]interface{})
		allowCredentials := corsConfig["allow_credentials"].(bool)
		origins := corsConfig["origins"].([]interface{})
		if !allowCredentials || len(origins) > 0 {
			return fmt.Errorf("[ERROR] Setting \"cors_config\" conflicts with enable_cors set to false")
		}
	}
	return nil
}

func setCloudantInstanceCors(client *cloudantv1.CloudantV1, d *schema.ResourceData) error {
	corsInformation, err := readCloudantInstanceCors(client)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving CORS config: %s", err)
	}
	if corsInformation != nil {
		d.Set("enable_cors", corsInformation.EnableCors)

		if *corsInformation.EnableCors {
			corsConfig := []map[string]interface{}{
				map[string]interface{}{
					"allow_credentials": corsInformation.AllowCredentials,
					"origins":           corsInformation.Origins,
				},
			}
			d.Set("cors_config", corsConfig)
		}
	}
	return nil
}

func readCloudantInstanceCors(client *cloudantv1.CloudantV1) (*cloudantv1.CorsInformation, error) {
	opts := client.NewGetCorsInformationOptions()

	corsInformation, response, err := client.GetCorsInformation(opts)
	if err != nil {
		log.Printf("[DEBUG] Error retrieving CORS config: %s\n%s", err, response)
	}
	return corsInformation, err
}

func updateCloudantInstanceCors(client *cloudantv1.CloudantV1, d *schema.ResourceData) error {
	enableCors := d.Get("enable_cors").(bool)
	allowCredentials := true
	origins := make([]string, 0)
	corsConfigRaw := d.Get("cors_config").([]interface{})
	if enableCors && len(corsConfigRaw) > 0 {
		corsConfig := corsConfigRaw[0].(map[string]interface{})
		allowCredentials = corsConfig["allow_credentials"].(bool)
		origins = flex.ExpandStringList(corsConfig["origins"].([]interface{}))
	}

	opts := client.NewPutCorsConfigurationOptions(origins)
	opts.SetEnableCors(enableCors)
	opts.SetAllowCredentials(allowCredentials)

	_, response, err := client.PutCorsConfiguration(opts)
	if err != nil {
		log.Printf("[DEBUG] Error updating CORS settings: %s\n%s", err, response)
	}
	return err
}
