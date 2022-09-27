// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func ResourceIBMCbrZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCbrZoneCreate,
		ReadContext:   resourceIBMCbrZoneRead,
		UpdateContext: resourceIBMCbrZoneUpdate,
		DeleteContext: resourceIBMCbrZoneDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_zone", "name"),
				Description:  "The name of the zone.",
			},
			"account_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_zone", "account_id"),
				Description:  "The id of the account owning this zone.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_zone", "description"),
				Description:  "The description of the zone.",
			},
			"addresses": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of addresses in the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of address.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IP address.",
						},
						"ref": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "A service reference value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The id of the account owning the service.",
									},
									"service_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The service type.",
									},
									"service_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The service name.",
									},
									"service_instance": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The service instance.",
									},
									"location": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The location.",
									},
								},
							},
						},
					},
				},
			},
			"excluded": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of excluded addresses in the zone. Only addresses of type `ipAddress`, `ipRange`, and `subnet` can be excluded.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of address.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IP address.",
						},
					},
				},
			},
			"x_correlation_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_zone", "x_correlation_id"),
				Description:  "The supplied or generated value of this header is logged for a request and repeated in a response header for the corresponding response. The same value is used for downstream requests and retries of those requests. If a value of this headers is not supplied in a request, the service generates a random (version 4) UUID.",
			},
			"transaction_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_zone", "transaction_id"),
				Description:  "The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zone CRN.",
			},
			"address_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of addresses in the zone.",
			},
			"excluded_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of excluded addresses in the zone.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The href link to the resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the resource was created.",
			},
			"created_by_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which created the resource.",
			},
			"last_modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the resource was modified.",
			},
			"last_modified_by_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which modified the resource.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMCbrZoneValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 \-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "account_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9\-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[\x20-\xFE]*$`,
			MinValueLength:             0,
			MaxValueLength:             300,
		},
		validate.ValidateSchema{
			Identifier:                 "x_correlation_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 ,\-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             1024,
		},
		validate.ValidateSchema{
			Identifier:                 "transaction_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 ,\-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             1024,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cbr_zone", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCbrZoneCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createZoneOptions := &contextbasedrestrictionsv1.CreateZoneOptions{}

	if _, ok := d.GetOk("name"); ok {
		createZoneOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("account_id"); ok {
		createZoneOptions.SetAccountID(d.Get("account_id").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createZoneOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("addresses"); ok {
		var addresses []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("addresses").([]interface{}) {
			value := e.(map[string]interface{})
			addressesItem, err := resourceIBMCbrZoneMapToAddress(value)
			if err != nil {
				return diag.FromErr(err)
			}
			addresses = append(addresses, addressesItem)
		}
		createZoneOptions.SetAddresses(addresses)
	}
	if _, ok := d.GetOk("excluded"); ok {
		var excluded []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("excluded").([]interface{}) {
			value := e.(map[string]interface{})
			excludedItem, err := resourceIBMCbrZoneMapToAddress(value)
			if err != nil {
				return diag.FromErr(err)
			}
			excluded = append(excluded, excludedItem)
		}
		createZoneOptions.SetExcluded(excluded)
	}
	if _, ok := d.GetOk("x_correlation_id"); ok {
		createZoneOptions.SetXCorrelationID(d.Get("x_correlation_id").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		createZoneOptions.SetTransactionID(d.Get("transaction_id").(string))
	}

	zone, response, err := contextBasedRestrictionsClient.CreateZoneWithContext(context, createZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId(*zone.ID)

	return resourceIBMCbrZoneRead(context, d, meta)
}

func resourceIBMCbrZoneRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getZoneOptions := &contextbasedrestrictionsv1.GetZoneOptions{}

	getZoneOptions.SetZoneID(d.Id())

	zone, response, err := contextBasedRestrictionsClient.GetZoneWithContext(context, getZoneOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetZoneWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("x_correlation_id", getZoneOptions.XCorrelationID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting x_correlation_id: %s", err))
	}
	if err = d.Set("transaction_id", getZoneOptions.TransactionID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting transaction_id: %s", err))
	}
	if err = d.Set("name", zone.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("account_id", zone.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("description", zone.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	addresses := []map[string]interface{}{}
	if zone.Addresses != nil {
		for _, addressesItem := range zone.Addresses {
			addressesItemMap, err := resourceIBMCbrZoneAddressToMap(addressesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			addresses = append(addresses, addressesItemMap)
		}
	}
	if err = d.Set("addresses", addresses); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting addresses: %s", err))
	}
	excluded := []map[string]interface{}{}
	if zone.Excluded != nil {
		for _, excludedItem := range zone.Excluded {
			excludedItemMap, err := resourceIBMCbrZoneAddressToMap(excludedItem)
			if err != nil {
				return diag.FromErr(err)
			}
			excluded = append(excluded, excludedItemMap)
		}
	}
	if err = d.Set("excluded", excluded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting excluded: %s", err))
	}
	if err = d.Set("crn", zone.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("address_count", flex.IntValue(zone.AddressCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address_count: %s", err))
	}
	if err = d.Set("excluded_count", flex.IntValue(zone.ExcludedCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting excluded_count: %s", err))
	}
	if err = d.Set("href", zone.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(zone.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by_id", zone.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}
	if err = d.Set("last_modified_at", flex.DateTimeToString(zone.LastModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_at: %s", err))
	}
	if err = d.Set("last_modified_by_id", zone.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_by_id: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	return nil
}

func resourceIBMCbrZoneUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceZoneOptions := &contextbasedrestrictionsv1.ReplaceZoneOptions{}

	replaceZoneOptions.SetZoneID(d.Id())
	if _, ok := d.GetOk("name"); ok {
		replaceZoneOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("account_id"); ok {
		replaceZoneOptions.SetAccountID(d.Get("account_id").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		replaceZoneOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("addresses"); ok {
		var addresses []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("addresses").([]interface{}) {
			value := e.(map[string]interface{})
			addressesItem, err := resourceIBMCbrZoneMapToAddress(value)
			if err != nil {
				return diag.FromErr(err)
			}
			addresses = append(addresses, addressesItem)
		}
		replaceZoneOptions.SetAddresses(addresses)
	}
	if _, ok := d.GetOk("excluded"); ok {
		var excluded []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("excluded").([]interface{}) {
			value := e.(map[string]interface{})
			excludedItem, err := resourceIBMCbrZoneMapToAddress(value)
			if err != nil {
				return diag.FromErr(err)
			}
			excluded = append(excluded, excludedItem)
		}
		replaceZoneOptions.SetExcluded(excluded)
	}
	if _, ok := d.GetOk("x_correlation_id"); ok {
		replaceZoneOptions.SetXCorrelationID(d.Get("x_correlation_id").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		replaceZoneOptions.SetTransactionID(d.Get("transaction_id").(string))
	}
	replaceZoneOptions.SetIfMatch(d.Get("version").(string))

	_, response, err := contextBasedRestrictionsClient.ReplaceZoneWithContext(context, replaceZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] ReplaceZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ReplaceZoneWithContext failed %s\n%s", err, response))
	}

	return resourceIBMCbrZoneRead(context, d, meta)
}

func resourceIBMCbrZoneDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteZoneOptions := &contextbasedrestrictionsv1.DeleteZoneOptions{}

	deleteZoneOptions.SetZoneID(d.Id())

	response, err := contextBasedRestrictionsClient.DeleteZoneWithContext(context, deleteZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCbrZoneMapToAddress(modelMap map[string]interface{}) (contextbasedrestrictionsv1.AddressIntf, error) {
	discValue, ok := modelMap["type"]
	if ok {
		if discValue == "ipAddress" {
			return resourceIBMCbrZoneMapToAddressIPAddress(modelMap)
		} else if discValue == "ipRange" {
			return resourceIBMCbrZoneMapToAddressIPAddressRange(modelMap)
		} else if discValue == "subnet" {
			return resourceIBMCbrZoneMapToAddressSubnet(modelMap)
		} else if discValue == "vpc" {
			return resourceIBMCbrZoneMapToAddressVPC(modelMap)
		} else if discValue == "serviceRef" {
			return resourceIBMCbrZoneMapToAddressServiceRef(modelMap)
		} else {
			return nil, fmt.Errorf("unexpected value for discriminator property 'type' found in map: '%s'", discValue)
		}
	} else {
		return nil, fmt.Errorf("discriminator property 'type' not found in map")
	}
}

func resourceIBMCbrZoneMapToServiceRefValue(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.ServiceRefValue, error) {
	model := &contextbasedrestrictionsv1.ServiceRefValue{}
	model.AccountID = core.StringPtr(modelMap["account_id"].(string))
	if modelMap["service_type"] != nil && modelMap["service_type"].(string) != "" {
		model.ServiceType = core.StringPtr(modelMap["service_type"].(string))
	}
	if modelMap["service_name"] != nil && modelMap["service_name"].(string) != "" {
		model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	}
	if modelMap["service_instance"] != nil && modelMap["service_instance"].(string) != "" {
		model.ServiceInstance = core.StringPtr(modelMap["service_instance"].(string))
	}
	if modelMap["location"] != nil && modelMap["location"].(string) != "" {
		model.Location = core.StringPtr(modelMap["location"].(string))
	}
	return model, nil
}

func resourceIBMCbrZoneMapToAddressIPAddress(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.AddressIPAddress, error) {
	model := &contextbasedrestrictionsv1.AddressIPAddress{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMCbrZoneMapToAddressServiceRef(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.AddressServiceRef, error) {
	model := &contextbasedrestrictionsv1.AddressServiceRef{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	RefModel, err := resourceIBMCbrZoneMapToServiceRefValue(modelMap["ref"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Ref = RefModel
	return model, nil
}

func resourceIBMCbrZoneMapToAddressSubnet(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.AddressSubnet, error) {
	model := &contextbasedrestrictionsv1.AddressSubnet{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMCbrZoneMapToAddressIPAddressRange(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.AddressIPAddressRange, error) {
	model := &contextbasedrestrictionsv1.AddressIPAddressRange{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMCbrZoneMapToAddressVPC(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.AddressVPC, error) {
	model := &contextbasedrestrictionsv1.AddressVPC{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMCbrZoneAddressToMap(model contextbasedrestrictionsv1.AddressIntf) (map[string]interface{}, error) {
	if _, ok := model.(*contextbasedrestrictionsv1.AddressIPAddress); ok {
		return resourceIBMCbrZoneAddressIPAddressToMap(model.(*contextbasedrestrictionsv1.AddressIPAddress))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressIPAddressRange); ok {
		return resourceIBMCbrZoneAddressIPAddressRangeToMap(model.(*contextbasedrestrictionsv1.AddressIPAddressRange))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressSubnet); ok {
		return resourceIBMCbrZoneAddressSubnetToMap(model.(*contextbasedrestrictionsv1.AddressSubnet))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressVPC); ok {
		return resourceIBMCbrZoneAddressVPCToMap(model.(*contextbasedrestrictionsv1.AddressVPC))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressServiceRef); ok {
		return resourceIBMCbrZoneAddressServiceRefToMap(model.(*contextbasedrestrictionsv1.AddressServiceRef))
	} else if _, ok := model.(*contextbasedrestrictionsv1.Address); ok {
		modelMap := make(map[string]interface{})
		model := model.(*contextbasedrestrictionsv1.Address)
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		if model.Value != nil {
			modelMap["value"] = model.Value
		}
		if model.Ref != nil {
			refMap, err := resourceIBMCbrZoneServiceRefValueToMap(model.Ref)
			if err != nil {
				return modelMap, err
			}
			modelMap["ref"] = []map[string]interface{}{refMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized contextbasedrestrictionsv1.AddressIntf subtype encountered")
	}
}

func resourceIBMCbrZoneServiceRefValueToMap(model *contextbasedrestrictionsv1.ServiceRefValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["account_id"] = model.AccountID
	if model.ServiceType != nil {
		modelMap["service_type"] = model.ServiceType
	}
	if model.ServiceName != nil {
		modelMap["service_name"] = model.ServiceName
	}
	if model.ServiceInstance != nil {
		modelMap["service_instance"] = model.ServiceInstance
	}
	if model.Location != nil {
		modelMap["location"] = model.Location
	}
	return modelMap, nil
}

func resourceIBMCbrZoneAddressIPAddressToMap(model *contextbasedrestrictionsv1.AddressIPAddress) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIBMCbrZoneAddressServiceRefToMap(model *contextbasedrestrictionsv1.AddressServiceRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	refMap, err := resourceIBMCbrZoneServiceRefValueToMap(model.Ref)
	if err != nil {
		return modelMap, err
	}
	modelMap["ref"] = []map[string]interface{}{refMap}
	return modelMap, nil
}

func resourceIBMCbrZoneAddressSubnetToMap(model *contextbasedrestrictionsv1.AddressSubnet) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIBMCbrZoneAddressIPAddressRangeToMap(model *contextbasedrestrictionsv1.AddressIPAddressRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIBMCbrZoneAddressVPCToMap(model *contextbasedrestrictionsv1.AddressVPC) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}
