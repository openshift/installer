// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventstreams

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/eventstreams-go-sdk/pkg/schemaregistryv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMEventStreamsSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEventStreamsSchemaCreate,
		ReadContext:   resourceIBMEventStreamsSchemaRead,
		UpdateContext: resourceIBMEventStreamsSchemaUpdate,
		DeleteContext: resourceIBMEventStreamsSchemaDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"resource_instance_id": {
				Type:        schema.TypeString,
				Description: "The ID or the CRN of the Event Streams service instance",
				Required:    true,
				ForceNew:    true,
			},
			"kafka_http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API endpoint for interacting with an Event Streams REST API",
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				ValidateFunc: validateAvroSchema,
				Description:  "The schema in JSON format",
			},
			"schema_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID to be assigned to schema, which must be unique. If this value is not specified, a generated UUID is assigned.",
			},
		},
	}
}

var primitiveTypes = map[string]Type{
	"null":    Null{},
	"boolean": Boolean{},
	"int":     Int{},
	"long":    Long{},
	"float":   Float{},
	"double":  Double{},
	"bytes":   Bytes{},
	"string":  String{},
}

type Null struct{}
type Boolean struct{}
type Int struct{}
type Long struct{}
type Float struct{}
type Double struct{}
type Bytes struct{}

type String struct{}
type Type interface {
	_type() string
}

func (Null) _type() string    { return "null" }
func (Boolean) _type() string { return "boolean" }
func (Int) _type() string     { return "int" }
func (Long) _type() string    { return "long" }
func (Float) _type() string   { return "float" }
func (Double) _type() string  { return "double" }
func (Bytes) _type() string   { return "bytes" }
func (String) _type() string  { return "string" }

func validateAvroSchema(v interface{}, k string) (ws []string, errors []error) {
	if _, err := flex.NormalizeJSONString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	var j interface{}
	err := json.Unmarshal([]byte(v.(string)), &j)
	if err != nil {
		errors = append(errors, err)
	}

	switch avroType := j.(type) {
	case map[string]interface{}:
		validationError := validateAvroObjectSchema(avroType, k)
		if validationError != nil {
			errors = append(errors, validationError)
		}
	default:
		errors = append(errors, fmt.Errorf("unknown type: %q for %q", avroType, k))
	}
	return
}

func validateAvroObjectSchema(avroType map[string]interface{}, k string) error {
	mp := "given avro schema"
	v, ok := avroType["type"]
	if !ok {
		return fmt.Errorf("%s must have 'type' property", mp)
	}
	valueString, ok := v.(string)
	if !ok {
		return fmt.Errorf("%s has non-string 'type' property", mp)
	}
	switch valueString {
	case "record":
		if err := validateName(avroType, valueString); err != nil {
			return err
		}
		if _, ok := avroType["fields"]; !ok {
			return fmt.Errorf("%s of type %s must contain 'fields' property", mp, valueString)
		}
	case "enum":
		if err := validateName(avroType, valueString); err != nil {
			return err
		}
		if _, ok := avroType["symbols"]; !ok {
			return fmt.Errorf("%s of type %s must contain 'symbols' property", mp, valueString)
		}
	case "array":
		if _, ok := avroType["items"]; !ok {
			return fmt.Errorf("%s of type %s must contain 'items'", mp, valueString)
		}
	case "map":
		if _, ok := avroType["values"]; !ok {
			return fmt.Errorf("%s of type %s must contain 'values'", mp, valueString)
		}
	case "fixed":
		if err := validateName(avroType, valueString); err != nil {
			return err
		}
		if _, ok := avroType["size"]; !ok {
			return fmt.Errorf("%s of type %s must contain 'size'", mp, valueString)
		}
	default:
		if _, ok = primitiveTypes[valueString]; !ok {
			return fmt.Errorf("unknown type: %q in schema", valueString)
		}
	}
	return nil
}

func validateName(s map[string]interface{}, t string) error {
	if _, ok := s["name"]; !ok {
		return fmt.Errorf("missing required key-value: given schema of type %s must contain 'name'", t)
	}
	return nil
}

func resourceIBMEventStreamsSchemaCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaregistryClient, err := meta.(conns.ClientSession).ESschemaRegistrySession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaCreate schemaregistryClient: %s", err), "ibm_event_streams_schema", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	adminURL, instanceCRN, err := getInstanceURL(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaCreate getInstanceURL: %s", err), "ibm_event_streams_schema", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	schemaregistryClient.SetServiceURL(adminURL)
	createSchemaOptions := &schemaregistryv1.CreateSchemaOptions{}

	if s, ok := d.GetOk("schema"); ok {
		var schema map[string]interface{}
		json.Unmarshal([]byte(s.(string)), &schema)
		createSchemaOptions.Schema = schema
	}
	if _, ok := d.GetOk("schema_id"); ok {
		createSchemaOptions.SetID(d.Get("schema_id").(string))
	}

	schemaMetadata, response, err := schemaregistryClient.CreateSchemaWithContext(context, createSchemaOptions)
	if err != nil || schemaMetadata == nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaCreate CreateSchemaWithContext failed with error: %s and response:\n%s", err, response), "ibm_event_streams_schema", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	uniqueID := getUniqueSchemaID(instanceCRN, *schemaMetadata.ID)
	d.SetId(uniqueID)

	return resourceIBMEventStreamsSchemaRead(context, d, meta)
}

func resourceIBMEventStreamsSchemaRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaregistryClient, err := meta.(conns.ClientSession).ESschemaRegistrySession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaRead schemaregistryClient: %s", err), "ibm_event_streams_schema", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	adminURL, instanceCRN, err := getInstanceURL(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaRead getInstanceURL: %s", err), "ibm_event_streams_schema", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	schemaregistryClient.SetServiceURL(adminURL)

	getSchemaOptions := &schemaregistryv1.GetLatestSchemaOptions{}

	schemaID := getSchemaID(d.Id())
	getSchemaOptions.SetID(schemaID)

	avroSchema, response, err := schemaregistryClient.GetLatestSchemaWithContext(context, getSchemaOptions)
	if err != nil || avroSchema == nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("[DEBUG] GetSchemaWithContext failed with 404 Not Found error and response: \n%s", response)
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaRead GetLatestSchemaWithContext failed with error: %s and response:\n%s", err, response), "ibm_event_streams_schema", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	s, err := json.Marshal(avroSchema)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaRead marshalling the schema failed with error: %s", err), "ibm_event_streams_schema", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("schema", string(s)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaRead setting the schema failed error: %s", err), "ibm_event_streams_schema", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.Set("resource_instance_id", instanceCRN)
	d.Set("schema_id", schemaID)

	return nil
}

func resourceIBMEventStreamsSchemaUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaregistryClient, err := meta.(conns.ClientSession).ESschemaRegistrySession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaUpdate schemaregistryClient: %s", err), "ibm_event_streams_schema", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	adminURL, _, err := getInstanceURL(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaUpdate getInstanceURL: %s", err), "ibm_event_streams_schema", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	schemaregistryClient.SetServiceURL(adminURL)

	updateSchemaOptions := &schemaregistryv1.UpdateSchemaOptions{}
	schemaID := d.Get("schema_id").(string)
	updateSchemaOptions.SetID(schemaID)

	if d.HasChange("schema") {
		if s, ok := d.GetOk("schema"); ok {
			var schema map[string]interface{}
			json.Unmarshal([]byte(s.(string)), &schema)
			updateSchemaOptions.Schema = schema
		}
		schemaMetadata, response, err := schemaregistryClient.UpdateSchemaWithContext(context, updateSchemaOptions)
		if err != nil || schemaMetadata == nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaUpdate UpdateSchema failed with error: %s and response:\n%s", err, response), "ibm_event_streams_schema", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMEventStreamsSchemaRead(context, d, meta)
}

func resourceIBMEventStreamsSchemaDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaregistryClient, err := meta.(conns.ClientSession).ESschemaRegistrySession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaDelete schemaregistryClient: %s", err), "ibm_event_streams_schema", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	adminURL, _, err := getInstanceURL(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaDelete getInstanceURL: %s", err), "ibm_event_streams_schema", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	schemaregistryClient.SetServiceURL(adminURL)

	deleteSchemaOptions := &schemaregistryv1.DeleteSchemaOptions{}
	schemaID := d.Get("schema_id").(string)
	deleteSchemaOptions.SetID(schemaID)

	setSchemaOptions := &schemaregistryv1.SetSchemaStateOptions{}
	setSchemaOptions.SetID(schemaID)
	setSchemaOptions.SetState("DISABLED")

	// set schema state to disabled before deleting
	response, err := schemaregistryClient.SetSchemaStateWithContext(context, setSchemaOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaDelete SetSchemaState failed with error: %s and response:\n%s", err, response), "ibm_event_streams_schema", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	response, err = schemaregistryClient.DeleteSchemaWithContext(context, deleteSchemaOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMEventStreamsSchemaDelete DeleteSchema failed with error: %s and response:\n%s", err, response), "ibm_event_streams_schema", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func getInstanceURL(d *schema.ResourceData, meta interface{}) (string, string, error) {
	instanceCRN := d.Get("resource_instance_id").(string)
	if len(instanceCRN) == 0 {
		schemaID := d.Id()
		if len(schemaID) == 0 || !strings.Contains(schemaID, ":") {
			log.Printf("[DEBUG] getInstanceURL resource_instance_id is missing")
			return "", "", fmt.Errorf("resource_instance_id is required")
		}
		instanceCRN = getInstanceCRN(schemaID)
	}

	instance, err := getInstanceDetails(instanceCRN, meta)
	if err != nil {
		return "", "", err
	}

	adminURL := instance.Extensions["kafka_http_url"].(string)
	planID := *instance.ResourcePlanID
	valid := strings.Contains(planID, "enterprise")
	if !valid {
		return "", "", fmt.Errorf("schema registry is not supported by the Event Streams %s plan, enterprise plan is expected",
			planID)
	}
	d.Set("kafka_http_url", adminURL)
	log.Printf("[INFO]getInstanceURL kafka_http_url is set to %s", adminURL)
	return adminURL, instanceCRN, nil
}

func getInstanceDetails(crn string, meta interface{}) (*resourcecontrollerv2.ResourceInstance, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		log.Printf("[DEBUG] getInstanceURL ResourceControllerAPI err %s", err)
		return nil, err
	}

	getResourceInstanceOptions := rsConClient.NewGetResourceInstanceOptions(crn)
	instance, response, err := rsConClient.GetResourceInstance(getResourceInstanceOptions)
	if err != nil || instance == nil {
		log.Printf("[DEBUG]getInstanceDetails GetResourceInstance failed with error: %s and response:\n%s", err, response)
		return nil, err
	}
	if instance.Extensions == nil {
		log.Printf("[DEBUG]instance %s extension is nil", *instance.ID)
		return nil, fmt.Errorf("getInstanceDetails instance %s extension is nil", *instance.ID)
	}

	return instance, nil
}

func getUniqueSchemaID(instanceCRN string, schemaID string) string {
	crnSegments := strings.Split(instanceCRN, ":")
	crnSegments[8] = "schema"
	crnSegments[9] = schemaID
	return strings.Join(crnSegments, ":")
}

func getSchemaID(id string) string {
	return strings.Split(id, ":")[9]
}
