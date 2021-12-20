// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func resourceIBMCbrZone() *schema.Resource {
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
				ValidateFunc: InvokeValidator("ibm_cbr_zone", "name"),
				Description:  "The name of the zone.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the account owning this zone.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_cbr_zone", "description"),
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
							Optional:    true,
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
										Computed:    true,
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
							Optional:    true,
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
										Computed:    true,
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
								},
							},
						},
					},
				},
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

func resourceIBMCbrZoneValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 \\-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[\x20-\xFE]*$`,
			MinValueLength:             0,
			MaxValueLength:             300,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_cbr_zone", Schema: validateSchema}
	return &resourceValidator
}

func getIBMCbrAccountId(meta interface{}) (string, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()

	if err != nil {
		return "", err
	} else {
		return userDetails.userAccount, nil
	}
}

func resourceIBMCbrZoneCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createZoneOptions := &contextbasedrestrictionsv1.CreateZoneOptions{}

	if _, ok := d.GetOk("name"); ok {
		createZoneOptions.SetName(d.Get("name").(string))
	}

	accountID, err := getIBMCbrAccountId(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	createZoneOptions.SetAccountID(accountID)

	if _, ok := d.GetOk("description"); ok {
		createZoneOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("addresses"); ok {
		var addresses []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("addresses").([]interface{}) {
			value := e.(map[string]interface{})
			addressesItem := resourceIBMCbrZoneMapToAddress(value, accountID)
			addresses = append(addresses, addressesItem)
		}
		createZoneOptions.SetAddresses(addresses)
	}
	if _, ok := d.GetOk("excluded"); ok {
		var excluded []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("excluded").([]interface{}) {
			value := e.(map[string]interface{})
			excludedItem := resourceIBMCbrZoneMapToAddress(value, accountID)
			excluded = append(excluded, excludedItem)
		}
		createZoneOptions.SetExcluded(excluded)
	}

	zone, response, err := contextBasedRestrictionsClient.CreateZoneWithContext(context, createZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId(*zone.ID)

	return resourceIBMCbrZoneRead(context, d, meta)
}

func resourceIBMCbrZoneMapToAddress(addressMap map[string]interface{}, accountID string) contextbasedrestrictionsv1.AddressIntf {
	var address contextbasedrestrictionsv1.AddressIntf
	disc, ok := addressMap["type"]
	if ok {
		switch disc {
		case "ipAddress":
			ipAddress := resourceIBMCbrZoneMapToAddressIPAddress(addressMap)
			address = &ipAddress
		case "ipRange":
			ipAddressRange := resourceIBMCbrZoneMapToAddressIPAddressRange(addressMap)
			address = &ipAddressRange
		case "subnet":
			subnet := resourceIBMCbrZoneMapToAddressSubnet(addressMap)
			address = &subnet
		case "vpc":
			vpc := resourceIBMCbrZoneMapToAddressVPC(addressMap)
			address = &vpc
		case "serviceRef":
			serviceRef := resourceIBMCbrZoneMapToAddressServiceRef(addressMap, accountID)
			address = &serviceRef
		}
	} else {
		log.Println("[DEBUG] 'type' field is missing from 'addresses'")
	}

	return address
}

func resourceIBMCbrZoneMapToServiceRefValue(serviceRefValueMap map[string]interface{}, accountID string) contextbasedrestrictionsv1.ServiceRefValue {
	serviceRefValue := contextbasedrestrictionsv1.ServiceRefValue{}

	serviceRefValue.AccountID = &accountID

	if serviceRefValueMap["service_type"] != nil && serviceRefValueMap["service_type"] != "" {
		serviceRefValue.ServiceType = core.StringPtr(serviceRefValueMap["service_type"].(string))
	}
	if serviceRefValueMap["service_name"] != nil && serviceRefValueMap["service_name"] != "" {
		serviceRefValue.ServiceName = core.StringPtr(serviceRefValueMap["service_name"].(string))
	}
	if serviceRefValueMap["service_instance"] != nil && serviceRefValueMap["service_instance"] != "" {
		serviceRefValue.ServiceInstance = core.StringPtr(serviceRefValueMap["service_instance"].(string))
	}

	return serviceRefValue
}

func resourceIBMCbrZoneMapToAddressIPAddress(addressIPAddressMap map[string]interface{}) contextbasedrestrictionsv1.AddressIPAddress {
	addressIPAddress := contextbasedrestrictionsv1.AddressIPAddress{}

	addressIPAddress.Type = core.StringPtr(addressIPAddressMap["type"].(string))
	addressIPAddress.Value = core.StringPtr(addressIPAddressMap["value"].(string))

	return addressIPAddress
}

func resourceIBMCbrZoneMapToAddressServiceRef(addressServiceRefMap map[string]interface{}, accountID string) contextbasedrestrictionsv1.AddressServiceRef {
	addressServiceRef := contextbasedrestrictionsv1.AddressServiceRef{}

	addressServiceRef.Type = core.StringPtr(addressServiceRefMap["type"].(string))

	if _, ok := addressServiceRefMap["value"]; ok {
		delete(addressServiceRefMap, "value")
	}

	if refSlice, ok := addressServiceRefMap["ref"]; ok {
		ref := refSlice.([]interface{})
		if len(ref) > 0 {
			serviceRefValue := resourceIBMCbrZoneMapToServiceRefValue(ref[0].(map[string]interface{}), accountID)
			addressServiceRef.Ref = &serviceRefValue
		}
	}

	return addressServiceRef
}

func resourceIBMCbrZoneMapToAddressSubnet(addressSubnetMap map[string]interface{}) contextbasedrestrictionsv1.AddressSubnet {
	addressSubnet := contextbasedrestrictionsv1.AddressSubnet{}

	addressSubnet.Type = core.StringPtr(addressSubnetMap["type"].(string))
	addressSubnet.Value = core.StringPtr(addressSubnetMap["value"].(string))

	return addressSubnet
}

func resourceIBMCbrZoneMapToAddressIPAddressRange(addressIPAddressRangeMap map[string]interface{}) contextbasedrestrictionsv1.AddressIPAddressRange {
	addressIPAddressRange := contextbasedrestrictionsv1.AddressIPAddressRange{}

	addressIPAddressRange.Type = core.StringPtr(addressIPAddressRangeMap["type"].(string))
	addressIPAddressRange.Value = core.StringPtr(addressIPAddressRangeMap["value"].(string))

	return addressIPAddressRange
}

func resourceIBMCbrZoneMapToAddressVPC(addressVPCMap map[string]interface{}) contextbasedrestrictionsv1.AddressVPC {
	addressVPC := contextbasedrestrictionsv1.AddressVPC{}

	addressVPC.Type = core.StringPtr(addressVPCMap["type"].(string))
	addressVPC.Value = core.StringPtr(addressVPCMap["value"].(string))

	return addressVPC
}

func resourceIBMCbrZoneRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getZoneOptions := &contextbasedrestrictionsv1.GetZoneOptions{}

	accountID, err := getIBMCbrAccountId(meta)
	if err != nil {
		return diag.FromErr(err)
	}

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

	if err = d.Set("name", zone.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("description", zone.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if zone.Addresses != nil {
		addresses := []map[string]interface{}{}
		for _, addressesItem := range zone.Addresses {
			addressesItemMap := resourceIBMCbrZoneAddressToMap(addressesItem, accountID)
			addresses = append(addresses, addressesItemMap)
		}

		if err = d.Set("addresses", addresses); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting addresses: %s", err))
		}
	}

	if zone.Excluded != nil {
		excluded := []map[string]interface{}{}
		for _, excludedItem := range zone.Excluded {
			excludedItemMap := resourceIBMCbrZoneAddressToMap(excludedItem, accountID)
			excluded = append(excluded, excludedItemMap)
		}
		if err = d.Set("excluded", excluded); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting excluded: %s", err))
		}
	}
	if err = d.Set("crn", zone.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("address_count", intValue(zone.AddressCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address_count: %s", err))
	}
	if err = d.Set("excluded_count", intValue(zone.ExcludedCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting excluded_count: %s", err))
	}
	if err = d.Set("href", zone.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(zone.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by_id", zone.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}
	if err = d.Set("last_modified_at", dateTimeToString(zone.LastModifiedAt)); err != nil {
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

func resourceIBMCbrZoneAddressToMap(address contextbasedrestrictionsv1.AddressIntf, accountID string) map[string]interface{} {
	addressMap := map[string]interface{}{}

	buf, err := json.Marshal(address)

	if err == nil {
		err = json.Unmarshal(buf, &addressMap)
	}

	if err != nil {
		panic(err)
	}

	if addressMap["type"] == "serviceRef" {
		var refArray []interface{}

		refMap := map[string]string{}
		refBuf, err := json.Marshal(addressMap["ref"])
		if err == nil {
			err = json.Unmarshal(refBuf, &refMap)
		}

		if err != nil {
			panic(err)
		}

		delete(refMap, "account_id")

		refArray = append(refArray, refMap)

		delete(addressMap, "ref")
		addressMap["ref"] = refArray

		addressMap["value"] = ""
	}

	return addressMap
}

func resourceIBMCbrZoneServiceRefValueToMap(serviceRefValue contextbasedrestrictionsv1.ServiceRefValue) map[string]interface{} {
	serviceRefValueMap := map[string]interface{}{}

	serviceRefValueMap["account_id"] = serviceRefValue.AccountID
	if serviceRefValue.ServiceType != nil {
		serviceRefValueMap["service_type"] = serviceRefValue.ServiceType
	}
	if serviceRefValue.ServiceName != nil {
		serviceRefValueMap["service_name"] = serviceRefValue.ServiceName
	}
	if serviceRefValue.ServiceInstance != nil {
		serviceRefValueMap["service_instance"] = serviceRefValue.ServiceInstance
	}

	return serviceRefValueMap
}

func resourceIBMCbrZoneAddressIPAddressToMap(addressIPAddress contextbasedrestrictionsv1.AddressIPAddress) map[string]interface{} {
	addressIPAddressMap := map[string]interface{}{}

	addressIPAddressMap["type"] = addressIPAddress.Type
	addressIPAddressMap["value"] = addressIPAddress.Value

	return addressIPAddressMap
}

func resourceIBMCbrZoneAddressServiceRefToMap(addressServiceRef contextbasedrestrictionsv1.AddressServiceRef) map[string]interface{} {
	addressServiceRefMap := map[string]interface{}{}

	addressServiceRefMap["type"] = addressServiceRef.Type
	RefMap := resourceIBMCbrZoneServiceRefValueToMap(*addressServiceRef.Ref)
	addressServiceRefMap["ref"] = []map[string]interface{}{RefMap}

	return addressServiceRefMap
}

func resourceIBMCbrZoneAddressSubnetToMap(addressSubnet contextbasedrestrictionsv1.AddressSubnet) map[string]interface{} {
	addressSubnetMap := map[string]interface{}{}

	addressSubnetMap["type"] = addressSubnet.Type
	addressSubnetMap["value"] = addressSubnet.Value

	return addressSubnetMap
}

func resourceIBMCbrZoneAddressIPAddressRangeToMap(addressIPAddressRange contextbasedrestrictionsv1.AddressIPAddressRange) map[string]interface{} {
	addressIPAddressRangeMap := map[string]interface{}{}

	addressIPAddressRangeMap["type"] = addressIPAddressRange.Type
	addressIPAddressRangeMap["value"] = addressIPAddressRange.Value

	return addressIPAddressRangeMap
}

func resourceIBMCbrZoneAddressVPCToMap(addressVPC contextbasedrestrictionsv1.AddressVPC) map[string]interface{} {
	addressVPCMap := map[string]interface{}{}

	addressVPCMap["type"] = addressVPC.Type
	addressVPCMap["value"] = addressVPC.Value

	return addressVPCMap
}

func resourceIBMCbrZoneUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceZoneOptions := &contextbasedrestrictionsv1.ReplaceZoneOptions{}

	replaceZoneOptions.SetZoneID(d.Id())
	if _, ok := d.GetOk("name"); ok {
		replaceZoneOptions.SetName(d.Get("name").(string))
	}

	accountID, err := getIBMCbrAccountId(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	replaceZoneOptions.SetAccountID(accountID)

	if _, ok := d.GetOk("description"); ok {
		replaceZoneOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("addresses"); ok {
		var addresses []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("addresses").([]interface{}) {
			value := e.(map[string]interface{})
			addressesItem := resourceIBMCbrZoneMapToAddress(value, accountID)
			addresses = append(addresses, addressesItem)
		}
		replaceZoneOptions.SetAddresses(addresses)
	}
	if _, ok := d.GetOk("excluded"); ok {
		var excluded []contextbasedrestrictionsv1.AddressIntf
		for _, e := range d.Get("excluded").([]interface{}) {
			value := e.(map[string]interface{})
			excludedItem := resourceIBMCbrZoneMapToAddress(value, accountID)
			excluded = append(excluded, excludedItem)
		}
		replaceZoneOptions.SetExcluded(excluded)
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
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
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
