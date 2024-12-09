// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmCodeEngineSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineSecretCreate,
		ReadContext:   resourceIbmCodeEngineSecretRead,
		UpdateContext: resourceIbmCodeEngineSecretUpdate,
		DeleteContext: resourceIbmCodeEngineSecretDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_secret", "project_id"),
				Description:  "The ID of the project.",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Data container that allows to specify config parameters and their values as a key-value map. Each key field must consist of alphanumeric characters, `-`, `_` or `.` and must not exceed a max length of 253 characters. Each value field can consists of any character and must not exceed a max length of 1048576 characters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_secret", "format"),
				Description:  "Specify the format of the secret.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_secret", "name"),
				Description:  "The name of the secret.",
			},
			"service_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Properties for Service Access Secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_key": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The service credential associated with the secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the service credential associated with the secret.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the service credential associated with the secret.",
									},
								},
							},
						},
						"role": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "A reference to the Role and Role CRN for service binding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CRN of the IAM Role for this service access secret.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Role of the service credential.",
									},
								},
							},
						},
						"service_instance": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The IBM Cloud service instance associated with the secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the IBM Cloud service instance associated with the secret.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of IBM Cloud service associated with the secret.",
									},
								},
							},
						},
						"serviceid": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "A reference to a Service ID.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN value of a Service ID.",
									},
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The ID of the Service ID.",
									},
								},
							},
						},
					},
				},
			},
			"service_operator": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Properties for the IBM Cloud Operator Secret.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apikey_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the apikey associated with the operator secret.",
						},
						"resource_group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The list of resource groups (by ID) that the operator secret can bind services in.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"serviceid": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "A reference to a Service ID.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN value of a Service ID.",
									},
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The ID of the Service ID.",
									},
								},
							},
						},
						"user_managed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether the operator secret is user managed.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the secret instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new secret,  a URL is created identifying the location of the instance.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the secret.",
			},
			"secret_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineSecretValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "project_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "format",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "basic_auth, generic, other, registry, service_access, service_operator, ssh_auth, tls",
			Regexp:                     `^(generic|ssh_auth|basic_auth|tls|service_access|registry|service_operator|other)$`,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_secret", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineSecretCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createSecretOptions := &codeenginev2.CreateSecretOptions{}

	createSecretOptions.SetProjectID(d.Get("project_id").(string))
	createSecretOptions.SetFormat(d.Get("format").(string))
	createSecretOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("data"); ok {
		dataModel, err := resourceIbmCodeEngineSecretMapToSecretData(d.Get("data").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createSecretOptions.SetData(dataModel)
	}
	if _, ok := d.GetOk("service_access"); ok {
		serviceAccessModel, err := resourceIbmCodeEngineSecretMapToServiceAccessSecretPrototypeProps(d.Get("service_access.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createSecretOptions.SetServiceAccess(serviceAccessModel)
	}
	if _, ok := d.GetOk("service_operator"); ok {
		serviceOperatorModel, err := resourceIbmCodeEngineSecretMapToOperatorSecretPrototypeProps(d.Get("service_operator.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createSecretOptions.SetServiceOperator(serviceOperatorModel)
	}

	secret, _, err := codeEngineClient.CreateSecretWithContext(context, createSecretOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretWithContext failed: %s", err.Error()), "ibm_code_engine_secret", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createSecretOptions.ProjectID, *secret.Name))

	return resourceIbmCodeEngineSecretRead(context, d, meta)
}

func resourceIbmCodeEngineSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSecretOptions := &codeenginev2.GetSecretOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "read")
		return tfErr.GetDiag()
	}

	getSecretOptions.SetProjectID(parts[0])
	getSecretOptions.SetName(parts[1])

	secret, response, err := codeEngineClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretWithContext failed: %s", err.Error()), "ibm_code_engine_secret", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", secret.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting project_id: %s", err))
	}
	if !core.IsNil(secret.Data) {
		data := make(map[string]string)
		for k, v := range secret.Data {
			data[k] = string(v)
		}
		if err = d.Set("data", data); err != nil {
			return diag.FromErr(fmt.Errorf("error setting data: %s", err))
		}
	}
	if err = d.Set("format", secret.Format); err != nil {
		return diag.FromErr(fmt.Errorf("error setting format: %s", err))
	}
	if err = d.Set("name", secret.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	if !core.IsNil(secret.ServiceAccess) {
		serviceAccessMap, err := resourceIbmCodeEngineSecretServiceAccessSecretPropsToMap(secret.ServiceAccess)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("service_access", []map[string]interface{}{serviceAccessMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting service_access: %s", err))
		}
	}
	if !core.IsNil(secret.ServiceOperator) {
		serviceOperatorMap, err := resourceIbmCodeEngineSecretOperatorSecretPropsToMap(secret.ServiceOperator)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("service_operator", []map[string]interface{}{serviceOperatorMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting service_operator: %s", err))
		}
	}
	if !core.IsNil(secret.CreatedAt) {
		if err = d.Set("created_at", secret.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
		}
	}
	if err = d.Set("entity_tag", secret.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("error setting entity_tag: %s", err))
	}
	if !core.IsNil(secret.Href) {
		if err = d.Set("href", secret.Href); err != nil {
			return diag.FromErr(fmt.Errorf("error setting href: %s", err))
		}
	}
	if !core.IsNil(secret.ID) {
		if err = d.Set("secret_id", secret.ID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting id: %s", err))
		}
	}
	if !core.IsNil(secret.Region) {
		if err = d.Set("region", secret.Region); err != nil {
			return diag.FromErr(fmt.Errorf("error setting region: %s", err))
		}
	}
	if !core.IsNil(secret.ResourceType) {
		if err = d.Set("resource_type", secret.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("error setting resource_type: %s", err))
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_secret", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineSecretUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	replaceSecretOptions := &codeenginev2.ReplaceSecretOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "update")
		return tfErr.GetDiag()
	}

	replaceSecretOptions.SetProjectID(parts[0])
	replaceSecretOptions.SetName(parts[1])
	replaceSecretOptions.SetFormat(d.Get("format").(string))

	hasChange := false

	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		tfErr := flex.TerraformErrorf(err, errMsg, "ibm_code_engine_secret", "update")
		return tfErr.GetDiag()
	}
	if d.HasChange("format") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "format")
		tfErr := flex.TerraformErrorf(err, errMsg, "ibm_code_engine_secret", "update")
		return tfErr.GetDiag()
	}
	if d.HasChange("data") {
		data, err := resourceIbmCodeEngineSecretMapToSecretData(d.Get("data").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceSecretOptions.SetData(data)
		hasChange = true
	}
	replaceSecretOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		_, _, err = codeEngineClient.ReplaceSecretWithContext(context, replaceSecretOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceSecretWithContext failed: %s", err.Error()), "ibm_code_engine_secret", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineSecretRead(context, d, meta)
}

func resourceIbmCodeEngineSecretDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteSecretOptions := &codeenginev2.DeleteSecretOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_secret", "delete")
		return tfErr.GetDiag()
	}

	deleteSecretOptions.SetProjectID(parts[0])
	deleteSecretOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteSecretWithContext(context, deleteSecretOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSecretWithContext failed: %s", err.Error()), "ibm_code_engine_secret", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmCodeEngineSecretMapToSecretData(modelMap map[string]interface{}) (codeenginev2.SecretDataIntf, error) {
	model := &codeenginev2.SecretData{}

	for key, value := range modelMap {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		model.SetProperty(strKey, &strValue)
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToServiceAccessSecretPrototypeProps(modelMap map[string]interface{}) (*codeenginev2.ServiceAccessSecretPrototypeProps, error) {
	model := &codeenginev2.ServiceAccessSecretPrototypeProps{}
	ResourceKeyModel, err := resourceIbmCodeEngineSecretMapToResourceKeyRefPrototype(modelMap["resource_key"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ResourceKey = ResourceKeyModel
	if modelMap["role"] != nil && len(modelMap["role"].([]interface{})) > 0 {
		RoleModel, err := resourceIbmCodeEngineSecretMapToRoleRefPrototype(modelMap["role"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Role = RoleModel
	}
	ServiceInstanceModel, err := resourceIbmCodeEngineSecretMapToServiceInstanceRefPrototype(modelMap["service_instance"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ServiceInstance = ServiceInstanceModel
	if modelMap["serviceid"] != nil && len(modelMap["serviceid"].([]interface{})) > 0 {
		ServiceidModel, err := resourceIbmCodeEngineSecretMapToServiceIDRef(modelMap["serviceid"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Serviceid = ServiceidModel
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToResourceKeyRefPrototype(modelMap map[string]interface{}) (*codeenginev2.ResourceKeyRefPrototype, error) {
	model := &codeenginev2.ResourceKeyRefPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToRoleRefPrototype(modelMap map[string]interface{}) (*codeenginev2.RoleRefPrototype, error) {
	model := &codeenginev2.RoleRefPrototype{}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.Crn = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToServiceInstanceRefPrototype(modelMap map[string]interface{}) (*codeenginev2.ServiceInstanceRefPrototype, error) {
	model := &codeenginev2.ServiceInstanceRefPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToServiceIDRef(modelMap map[string]interface{}) (*codeenginev2.ServiceIDRef, error) {
	model := &codeenginev2.ServiceIDRef{}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.Crn = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToOperatorSecretPrototypeProps(modelMap map[string]interface{}) (*codeenginev2.OperatorSecretPrototypeProps, error) {
	model := &codeenginev2.OperatorSecretPrototypeProps{}
	if modelMap["resource_group_ids"] != nil {
		resourceGroupIds := []string{}
		for _, resourceGroupIdsItem := range modelMap["resource_group_ids"].([]interface{}) {
			resourceGroupIds = append(resourceGroupIds, resourceGroupIdsItem.(string))
		}
		model.ResourceGroupIds = resourceGroupIds
	}
	if modelMap["serviceid"] != nil && len(modelMap["serviceid"].([]interface{})) > 0 {
		ServiceidModel, err := resourceIbmCodeEngineSecretMapToServiceIDRefPrototype(modelMap["serviceid"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Serviceid = ServiceidModel
	}
	return model, nil
}

func resourceIbmCodeEngineSecretMapToServiceIDRefPrototype(modelMap map[string]interface{}) (*codeenginev2.ServiceIDRefPrototype, error) {
	model := &codeenginev2.ServiceIDRefPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIbmCodeEngineSecretServiceAccessSecretPropsToMap(model *codeenginev2.ServiceAccessSecretProps) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	resourceKeyMap, err := resourceIbmCodeEngineSecretResourceKeyRefToMap(model.ResourceKey)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_key"] = []map[string]interface{}{resourceKeyMap}
	if model.Role != nil {
		roleMap, err := resourceIbmCodeEngineSecretRoleRefToMap(model.Role)
		if err != nil {
			return modelMap, err
		}
		modelMap["role"] = []map[string]interface{}{roleMap}
	}
	serviceInstanceMap, err := resourceIbmCodeEngineSecretServiceInstanceRefToMap(model.ServiceInstance)
	if err != nil {
		return modelMap, err
	}
	modelMap["service_instance"] = []map[string]interface{}{serviceInstanceMap}
	if model.Serviceid != nil {
		serviceidMap, err := resourceIbmCodeEngineSecretServiceIDRefToMap(model.Serviceid)
		if err != nil {
			return modelMap, err
		}
		modelMap["serviceid"] = []map[string]interface{}{serviceidMap}
	}
	return modelMap, nil
}

func resourceIbmCodeEngineSecretResourceKeyRefToMap(model *codeenginev2.ResourceKeyRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	return modelMap, nil
}

func resourceIbmCodeEngineSecretRoleRefToMap(model *codeenginev2.RoleRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Crn != nil {
		modelMap["crn"] = model.Crn
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	return modelMap, nil
}

func resourceIbmCodeEngineSecretServiceInstanceRefToMap(model *codeenginev2.ServiceInstanceRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIbmCodeEngineSecretServiceIDRefToMap(model *codeenginev2.ServiceIDRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Crn != nil {
		modelMap["crn"] = model.Crn
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}

func resourceIbmCodeEngineSecretOperatorSecretPropsToMap(model *codeenginev2.OperatorSecretProps) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["apikey_id"] = model.ApikeyID
	modelMap["resource_group_ids"] = model.ResourceGroupIds
	serviceidMap, err := resourceIbmCodeEngineSecretServiceIDRefToMap(model.Serviceid)
	if err != nil {
		return modelMap, err
	}
	modelMap["serviceid"] = []map[string]interface{}{serviceidMap}
	modelMap["user_managed"] = model.UserManaged
	return modelMap, nil
}
