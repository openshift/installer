// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

const (
	cbrZoneAddressIdDefault = ""
)

func ResourceIBMCbrZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCbrZoneCreate,
		ReadContext:   resourceIBMCbrZoneRead,
		UpdateContext: resourceIBMCbrZoneUpdate,
		DeleteContext: resourceIBMCbrZoneDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

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
				Optional:    true,
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

func getZone(cbrClient *contextbasedrestrictionsv1.ContextBasedRestrictionsV1, context context.Context, id string) (result *contextbasedrestrictionsv1.Zone, response *core.DetailedResponse, found bool, err error) {
	getZoneOptions := cbrClient.NewGetZoneOptions(id)

	result, response, err = cbrClient.GetZoneWithContext(context, getZoneOptions)
	found = err == nil
	if found {
		return
	}
	if response != nil && response.StatusCode == 404 {
		err = nil
		return
	}
	log.Printf("[DEBUG] GetZoneWithContext failed %s\n%s", err, response)
	err = fmt.Errorf("GetZoneWithContext failed %s\n%s", err, response)
	return
}

func resourceIBMCbrZoneCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createZoneOptions := contextBasedRestrictionsClient.NewCreateZoneOptions()

	if _, ok := d.GetOk("name"); ok {
		createZoneOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("account_id"); ok {
		createZoneOptions.SetAccountID(d.Get("account_id").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createZoneOptions.SetDescription(d.Get("description").(string))
	}
	addresses := []contextbasedrestrictionsv1.AddressIntf{}
	if _, ok := d.GetOk("addresses"); ok {
		addresses, err = resourceEncodeAddressList(d.Get("addresses").([]interface{}), cbrZoneAddressIdDefault)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	createZoneOptions.SetAddresses(addresses)
	if _, ok := d.GetOk("excluded"); ok {
		var excluded []contextbasedrestrictionsv1.AddressIntf
		excluded, err = resourceEncodeAddressList(d.Get("excluded").([]interface{}), cbrZoneAddressIdDefault)
		if err != nil {
			return diag.FromErr(err)
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

	if err := ResourceIBMCbrZoneSetData(zone, response, d); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone's resource data: %s", err))
	}
	return nil
}

func ResourceIBMCbrZoneSetData(zone *contextbasedrestrictionsv1.Zone, response *core.DetailedResponse, d *schema.ResourceData) error {
	if err := d.Set("name", zone.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("account_id", zone.AccountID); err != nil {
		return fmt.Errorf("Error setting account_id: %s", err)
	}
	if err := d.Set("description", zone.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}

	var addresses []map[string]interface{}
	addresses, err := resourceDecodeAddressList(zone.Addresses, cbrZoneAddressIdDefault)
	if err != nil {
		return err
	}
	if err = d.Set("addresses", addresses); err != nil {
		return fmt.Errorf("Error setting addresses: %s", err)
	}

	var excluded []map[string]interface{}
	excluded, err = resourceDecodeAddressList(zone.Excluded, cbrZoneAddressIdDefault)
	if err != nil {
		return err
	}
	if err = d.Set("excluded", excluded); err != nil {
		return fmt.Errorf("Error setting excluded: %s", err)
	}

	if err = d.Set("crn", zone.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("address_count", flex.IntValue(zone.AddressCount)); err != nil {
		return fmt.Errorf("Error setting address_count: %s", err)
	}
	if err = d.Set("excluded_count", flex.IntValue(zone.ExcludedCount)); err != nil {
		return fmt.Errorf("Error setting excluded_count: %s", err)
	}
	if err = d.Set("href", zone.Href); err != nil {
		return fmt.Errorf("Error setting href: %s", err)
	}
	if err = d.Set("created_at", flex.DateTimeToString(zone.CreatedAt)); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err = d.Set("created_by_id", zone.CreatedByID); err != nil {
		return fmt.Errorf("Error setting created_by_id: %s", err)
	}
	if err = d.Set("last_modified_at", flex.DateTimeToString(zone.LastModifiedAt)); err != nil {
		return fmt.Errorf("Error setting last_modified_at: %s", err)
	}
	if err = d.Set("last_modified_by_id", zone.LastModifiedByID); err != nil {
		return fmt.Errorf("Error setting last_modified_by_id: %s", err)
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}

	return nil
}

func resourceIBMCbrZoneRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	zone, response, found, err := getZone(contextBasedRestrictionsClient, context, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if !found {
		d.SetId("")
		return nil
	}

	if err = ResourceIBMCbrZoneSetData(zone, response, d); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone's resource data: %s", err))
	}

	return nil
}

func resourceIBMCbrZoneUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	zoneId := d.Id()

	// synchronize with zone address update operations
	mutex := zoneMutexKV.get(zoneId)
	mutex.Lock()
	defer mutex.Unlock()

	currentZone, response, found, err := getZone(contextBasedRestrictionsClient, context, zoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	if !found {
		d.SetId("")
		return nil
	}

	replaceZoneOptions := contextBasedRestrictionsClient.NewReplaceZoneOptions(zoneId, response.Headers.Get("Etag"))

	if _, ok := d.GetOk("name"); ok {
		replaceZoneOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("account_id"); ok {
		replaceZoneOptions.SetAccountID(d.Get("account_id").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		replaceZoneOptions.SetDescription(d.Get("description").(string))
	}
	addresses := []contextbasedrestrictionsv1.AddressIntf{}
	if _, ok := d.GetOk("addresses"); ok {
		addresses, err = resourceEncodeAddressList(d.Get("addresses").([]interface{}), cbrZoneAddressIdDefault)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	preservedAddresses := FilterAddressList(currentZone.Addresses, func(id string) bool {
		return id != cbrZoneAddressIdDefault
	})
	if len(preservedAddresses) > 0 {
		addresses = append(addresses, preservedAddresses...)
	}
	replaceZoneOptions.SetAddresses(addresses)
	if _, ok := d.GetOk("excluded"); ok {
		var excluded []contextbasedrestrictionsv1.AddressIntf
		excluded, err = resourceEncodeAddressList(d.Get("excluded").([]interface{}), cbrZoneAddressIdDefault)
		if err != nil {
			return diag.FromErr(err)
		}
		replaceZoneOptions.SetExcluded(excluded)
	}
	if _, ok := d.GetOk("x_correlation_id"); ok {
		replaceZoneOptions.SetXCorrelationID(d.Get("x_correlation_id").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		replaceZoneOptions.SetTransactionID(d.Get("transaction_id").(string))
	}

	zone, response, err := contextBasedRestrictionsClient.ReplaceZoneWithContext(context, replaceZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] ReplaceZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ReplaceZoneWithContext failed %s\n%s", err, response))
	}

	if err := ResourceIBMCbrZoneSetData(zone, response, d); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone's resource data: %s", err))
	}

	return nil
}

func resourceIBMCbrZoneDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteZoneOptions := contextBasedRestrictionsClient.NewDeleteZoneOptions(d.Id())

	response, err := contextBasedRestrictionsClient.DeleteZoneWithContext(context, deleteZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCbrZoneMapToAddress(modelMap map[string]interface{}, addressId string) (contextbasedrestrictionsv1.AddressIntf, error) {
	discValue, ok := modelMap["type"]
	if ok {
		if discValue == "ipAddress" {
			return resourceIBMCbrZoneMapToAddressIPAddress(modelMap, addressId)
		} else if discValue == "ipRange" {
			return resourceIBMCbrZoneMapToAddressIPAddressRange(modelMap, addressId)
		} else if discValue == "subnet" {
			return resourceIBMCbrZoneMapToAddressSubnet(modelMap, addressId)
		} else if discValue == "vpc" {
			return resourceIBMCbrZoneMapToAddressVPC(modelMap, addressId)
		} else if discValue == "serviceRef" {
			return resourceIBMCbrZoneMapToAddressServiceRef(modelMap, addressId)
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

func resourceIBMCbrZoneMapToAddressIPAddress(modelMap map[string]interface{}, addressId string) (*contextbasedrestrictionsv1.AddressIPAddress, error) {
	model := &contextbasedrestrictionsv1.AddressIPAddress{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if addressId != "" {
		model.ID = core.StringPtr(addressId)
	}
	return model, nil
}

func resourceIBMCbrZoneMapToAddressServiceRef(modelMap map[string]interface{}, addressId string) (*contextbasedrestrictionsv1.AddressServiceRef, error) {
	model := &contextbasedrestrictionsv1.AddressServiceRef{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	RefModel, err := resourceIBMCbrZoneMapToServiceRefValue(modelMap["ref"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Ref = RefModel
	if addressId != "" {
		model.ID = core.StringPtr(addressId)
	}
	return model, nil
}

func resourceIBMCbrZoneMapToAddressSubnet(modelMap map[string]interface{}, addressId string) (*contextbasedrestrictionsv1.AddressSubnet, error) {
	model := &contextbasedrestrictionsv1.AddressSubnet{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if addressId != "" {
		model.ID = core.StringPtr(addressId)
	}
	return model, nil
}

func resourceIBMCbrZoneMapToAddressIPAddressRange(modelMap map[string]interface{}, addressId string) (*contextbasedrestrictionsv1.AddressIPAddressRange, error) {
	model := &contextbasedrestrictionsv1.AddressIPAddressRange{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if addressId != "" {
		model.ID = core.StringPtr(addressId)
	}
	return model, nil
}

func resourceIBMCbrZoneMapToAddressVPC(modelMap map[string]interface{}, addressId string) (*contextbasedrestrictionsv1.AddressVPC, error) {
	model := &contextbasedrestrictionsv1.AddressVPC{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if addressId != "" {
		model.ID = core.StringPtr(addressId)
	}
	return model, nil
}

func resourceIBMCbrZoneAddressToMap(model contextbasedrestrictionsv1.AddressIntf) (modelMap map[string]interface{}, addressId string, err error) {
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
		modelMap = make(map[string]interface{})
		address := model.(*contextbasedrestrictionsv1.Address)
		if address.Type != nil {
			modelMap["type"] = address.Type
		}
		if address.Value != nil {
			modelMap["value"] = address.Value
		}
		if address.Ref != nil {
			var refMap map[string]interface{}
			refMap, err = resourceIBMCbrZoneServiceRefValueToMap(address.Ref)
			if err != nil {
				return
			}
			modelMap["ref"] = []map[string]interface{}{refMap}
		}
		if address.ID != nil {
			addressId = *address.ID
		}
	} else {
		err = fmt.Errorf("Unrecognized contextbasedrestrictionsv1.AddressIntf subtype encountered")
	}

	return
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

func resourceIBMCbrZoneAddressIPAddressToMap(model *contextbasedrestrictionsv1.AddressIPAddress) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func resourceIBMCbrZoneAddressServiceRefToMap(model *contextbasedrestrictionsv1.AddressServiceRef) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	modelMap["type"] = model.Type
	var refMap map[string]interface{}
	refMap, err = resourceIBMCbrZoneServiceRefValueToMap(model.Ref)
	if err != nil {
		return
	}
	modelMap["ref"] = []map[string]interface{}{refMap}
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func resourceIBMCbrZoneAddressSubnetToMap(model *contextbasedrestrictionsv1.AddressSubnet) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func resourceIBMCbrZoneAddressIPAddressRangeToMap(model *contextbasedrestrictionsv1.AddressIPAddressRange) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func resourceIBMCbrZoneAddressVPCToMap(model *contextbasedrestrictionsv1.AddressVPC) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func resourceDecodeAddressList(addresses []contextbasedrestrictionsv1.AddressIntf, wantAddressId string) (result []map[string]interface{}, err error) {
	result = make([]map[string]interface{}, 0, len(addresses))
	for _, addr := range addresses {
		var m map[string]interface{}
		var addressId string
		m, addressId, err = resourceIBMCbrZoneAddressToMap(addr)
		if err != nil {
			return
		}
		if addressId == wantAddressId {
			result = append(result, m)
		}
	}
	return
}

func resourceEncodeAddressList(addresses []interface{}, addressId string) (result []contextbasedrestrictionsv1.AddressIntf, err error) {
	result = make([]contextbasedrestrictionsv1.AddressIntf, 0, len(addresses))
	for _, item := range addresses {
		var addr contextbasedrestrictionsv1.AddressIntf
		addr, err = resourceIBMCbrZoneMapToAddress(item.(map[string]interface{}), addressId)
		if err != nil {
			return
		}
		result = append(result, addr)
	}
	return
}

func FilterAddressList(addresses []contextbasedrestrictionsv1.AddressIntf, keep func(id string) bool) (result []contextbasedrestrictionsv1.AddressIntf) {
	result = make([]contextbasedrestrictionsv1.AddressIntf, 0, len(addresses))
	for _, addr := range addresses {
		_, addressId, err := resourceIBMCbrZoneAddressToMap(addr)
		if err == nil && keep(addressId) {
			result = append(result, addr)
		}
	}
	return
}

// Synchronization for zone operations
var zoneMutexKV = newMutexKV()

type mutexKV struct {
	lock  sync.Mutex
	store map[string]*sync.Mutex
}

func (m *mutexKV) get(key string) *sync.Mutex {
	m.lock.Lock()
	defer m.lock.Unlock()
	mutex, ok := m.store[key]
	if !ok {
		mutex = &sync.Mutex{}
		m.store[key] = mutex
	}
	return mutex
}

func newMutexKV() *mutexKV {
	return &mutexKV{
		store: make(map[string]*sync.Mutex),
	}
}
