package scc

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func ResourceIbmSccInstanceSettings() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		CreateContext: resourceIbmSccInstanceSettingsCreate,
		ReadContext:   resourceIbmSccInstanceSettingsRead,
		UpdateContext: resourceIbmSccInstanceSettingsUpdate,
		DeleteContext: resourceIbmSccInstanceSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"event_notifications": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The Event Notifications settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The Event Notifications instance CRN.",
						},
						"source_description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The description of the source in Event Notifications connected Security and Compliance Center",
						},
						"source_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the Event Notifications source connected Security and Compliance Center instance CRN.",
						},
						"updated_on": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the Event Notifications connection was updated.",
						},
						"source_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Security and Compliance Center instance CRN.",
						},
					},
				},
			},
			"object_storage": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The Cloud Object Storage settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The connected Cloud Object Storage instance CRN.",
						},
						"bucket": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The connected Cloud Object Storage bucket name.",
						},
						"bucket_location": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Cloud Object Storage bucket location.",
						},
						"bucket_endpoint": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Cloud Object Storage bucket endpoint.",
						},
						"updated_on": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the bucket connection was updated.",
						},
					},
				},
			},
		},
	})
}

func ResourceIbmSccInstanceSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 ,\-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             1024,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_instance_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccInstanceSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSettingsOptions := &securityandcompliancecenterapiv3.UpdateSettingsOptions{}
	instance_id := d.Get("instance_id").(string)
	updateSettingsOptions.SetInstanceID(instance_id)

	var eventNotificationsModel *securityandcompliancecenterapiv3.EventNotificationsPrototype
	if _, ok := d.GetOk("event_notifications"); ok {
		eventNotificationsData, err := resourceIbmSccInstanceSettingsMapToEventNotifications(d.Get("event_notifications.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		eventNotificationsModel = eventNotificationsData
	}
	updateSettingsOptions.SetEventNotifications(eventNotificationsModel)

	var objectStorageModel *securityandcompliancecenterapiv3.ObjectStoragePrototype
	if _, ok := d.GetOk("object_storage"); ok {
		objectStorageData, err := resourceIbmSccInstanceSettingsMapToObjectStorage(d.Get("object_storage.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		objectStorageModel = objectStorageData
	}
	updateSettingsOptions.SetObjectStorage(objectStorageModel)

	_, response, err := adminClient.UpdateSettingsWithContext(context, updateSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("UpdateSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(instance_id)

	time.Sleep(5 * time.Second)

	return resourceIbmSccInstanceSettingsRead(context, d, meta)
}

func resourceIbmSccInstanceSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions := &securityandcompliancecenterapiv3.GetSettingsOptions{}
	instance_id := d.Id()
	getSettingsOptions.SetInstanceID(instance_id)

	settings, response, err := adminClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_id", instance_id); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}
	if !core.IsNil(settings.EventNotifications) {
		eventNotificationsMap, err := resourceIbmSccInstanceSettingsEventNotificationsToMap(settings.EventNotifications)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("event_notifications", []map[string]interface{}{eventNotificationsMap}); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting event_notifications: %s", err))
		}
	}
	if !core.IsNil(settings.ObjectStorage) {
		objectStorageMap, err := resourceIbmSccInstanceSettingsObjectStorageToMap(settings.ObjectStorage)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("object_storage", []map[string]interface{}{objectStorageMap}); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting object_storage: %s", err))
		}
	}

	return nil
}

func resourceIbmSccInstanceSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSettingsOptions := &securityandcompliancecenterapiv3.UpdateSettingsOptions{}
	instance_id := d.Get("instance_id").(string)
	updateSettingsOptions.SetInstanceID(instance_id)

	hasChange := false

	if d.HasChange("event_notifications") {
		eventNotifications, err := resourceIbmSccInstanceSettingsMapToEventNotifications(d.Get("event_notifications.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateSettingsOptions.SetEventNotifications(eventNotifications)
		hasChange = true
	}
	if d.HasChange("object_storage") {
		objectStorage, err := resourceIbmSccInstanceSettingsMapToObjectStorage(d.Get("object_storage.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateSettingsOptions.SetObjectStorage(objectStorage)
		hasChange = true
	}

	if hasChange {
		_, response, err := adminClient.UpdateSettingsWithContext(context, updateSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSettingsWithContext failed %s\n%s", err, response)
			return diag.FromErr(flex.FmtErrorf("UpdateSettingsWithContext failed %s\n%s", err, response))
		}
	}

	time.Sleep(5 * time.Second)

	return resourceIbmSccInstanceSettingsRead(context, d, meta)
}

func resourceIbmSccInstanceSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId("")

	return nil
}

func resourceIbmSccInstanceSettingsMapToEventNotifications(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.EventNotificationsPrototype, error) {
	model := &securityandcompliancecenterapiv3.EventNotificationsPrototype{}
	if modelMap["instance_crn"] != nil && modelMap["instance_crn"].(string) != "" {
		model.InstanceCRN = core.StringPtr(modelMap["instance_crn"].(string))
	}
	if modelMap["source_name"] != nil && modelMap["source_name"].(string) != "" {
		model.SourceName = core.StringPtr(modelMap["source_name"].(string))
	}
	if modelMap["source_description"] != nil && modelMap["source_description"].(string) != "" {
		model.SourceDescription = core.StringPtr(modelMap["source_description"].(string))
	}
	if core.StringNilMapper(model.InstanceCRN) != "" && core.StringNilMapper(model.SourceName) == "" {
		return model, errors.New("event_notifications.source_name needs to be defined along with event_notifications.instance_crn")
	}
	return model, nil
}

func resourceIbmSccInstanceSettingsMapToObjectStorage(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ObjectStoragePrototype, error) {
	model := &securityandcompliancecenterapiv3.ObjectStoragePrototype{}
	instanceCrnSet := false
	if modelMap["instance_crn"] != nil && modelMap["instance_crn"].(string) != "" {
		model.InstanceCRN = core.StringPtr(modelMap["instance_crn"].(string))
		instanceCrnSet = true
	}
	if modelMap["bucket"] != nil && modelMap["bucket"].(string) != "" {
		if instanceCrnSet {
			model.Bucket = core.StringPtr(modelMap["bucket"].(string))
		} else {
			return model, errors.New(`object_storage.instance_crn cannot be empty`)
		}
	}
	return model, nil
}

func resourceIbmSccInstanceSettingsEventNotificationsToMap(model *securityandcompliancecenterapiv3.EventNotifications) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InstanceCRN != nil {
		modelMap["instance_crn"] = model.InstanceCRN
	}
	if model.UpdatedOn != nil {
		modelMap["updated_on"] = model.UpdatedOn.String()
	}
	if model.SourceID != nil {
		modelMap["source_id"] = model.SourceID
	}
	if model.SourceDescription != nil {
		modelMap["source_description"] = model.SourceDescription
	}
	if model.SourceName != nil {
		modelMap["source_name"] = model.SourceName
	}
	return modelMap, nil
}

func resourceIbmSccInstanceSettingsObjectStorageToMap(model *securityandcompliancecenterapiv3.ObjectStorage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InstanceCRN != nil {
		modelMap["instance_crn"] = model.InstanceCRN
	}
	if model.Bucket != nil {
		modelMap["bucket"] = model.Bucket
	}
	if model.BucketLocation != nil {
		modelMap["bucket_location"] = model.BucketLocation
	}
	if model.BucketEndpoint != nil {
		modelMap["bucket_endpoint"] = model.BucketEndpoint
	}
	if model.UpdatedOn != nil {
		modelMap["updated_on"] = model.UpdatedOn.String()
	}
	return modelMap, nil
}
