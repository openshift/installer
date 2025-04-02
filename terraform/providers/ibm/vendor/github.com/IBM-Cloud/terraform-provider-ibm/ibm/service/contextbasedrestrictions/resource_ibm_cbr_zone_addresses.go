// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func ResourceIBMCbrZoneAddresses() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCbrZoneAddressesCreate,
		ReadContext:   resourceIBMCbrZoneAddressesRead,
		UpdateContext: resourceIBMCbrZoneAddressesUpdate,
		DeleteContext: resourceIBMCbrZoneAddressesDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"zone_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_zone_addresses", "zone_id"),
				Description:  "The id of the zone containing the addresses.",
			},
			"addresses": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of addresses added to the zone.",
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
		},
	}
}

func ResourceIBMCbrZoneAddressesValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "zone_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9\-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cbr_zone_addresses", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCbrZoneAddressesCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneId := d.Get("zone_id").(string)
	newUuid, _ := uuid.GenerateUUID()
	addressesId := fmt.Sprintf("TF-%s", newUuid)
	err := resourceReplaceZoneAddresses(context, d, meta, zoneId, addressesId, false)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("%s", err.Error()), "ibm_cbr_zone_addresses", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(composeZoneAddressesId(zoneId, addressesId))
	return nil

}

func resourceIBMCbrZoneAddressesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_zone_addresses", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	zoneId, addressesId := decomposeZoneAddressesId(d.Id())

	var zone *contextbasedrestrictionsv1.Zone
	var found bool
	zone, _, found, err = getZone(contextBasedRestrictionsClient, context, zoneId)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("%s", err.Error()), "ibm_cbr_zone_addresses", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if !found {
		d.SetId("")
		return nil
	}

	if err = d.Set("zone_id", zoneId); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_zone_addresses", "read", "set-zone-id").GetDiag()
	}

	var addresses []map[string]interface{}
	addresses, err = ResourceDecodeAddressList(zone.Addresses, addressesId)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_zone_addresses", "read", "ResourceDecodeAddressList").GetDiag()
	}
	if len(addresses) == 0 {
		d.SetId("")
		return nil
	}

	if err = d.Set("addresses", addresses); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_zone_addresses", "read", "set-addresses").GetDiag()
	}

	return nil
}

func resourceIBMCbrZoneAddressesUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneId, addressesId := decomposeZoneAddressesId(d.Id())
	err := resourceReplaceZoneAddresses(context, d, meta, zoneId, addressesId, false)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("%s", err.Error()), "ibm_cbr_zone_addresses", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	return nil
}

func resourceIBMCbrZoneAddressesDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneId, addressesId := decomposeZoneAddressesId(d.Id())
	err := d.Set("addresses", nil)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_zone_addresses", "delete", "set-addresses").GetDiag()
	}
	err = resourceReplaceZoneAddresses(context, d, meta, zoneId, addressesId, true)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("%s", err.Error()), "ibm_cbr_zone_addresses", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")
	return nil
}

func resourceReplaceZoneAddresses(context context.Context, d *schema.ResourceData, meta interface{}, zoneId string, addressesId string, delete bool) error {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return err
	}

	// synchronize with zone address update operations
	mutex := zoneMutexKV.get(zoneId)
	mutex.Lock()
	defer mutex.Unlock()

	currentZone, response, found, err := getZone(contextBasedRestrictionsClient, context, zoneId)
	if err != nil {
		return err
	}
	if !found {
		d.SetId("")
		if !delete {
			return fmt.Errorf("zone_id %s not found", zoneId)
		}
		return nil
	}

	replaceZoneOptions := contextBasedRestrictionsClient.NewReplaceZoneOptions(zoneId, response.Headers.Get("Etag"))
	replaceZoneOptions.SetName(*currentZone.Name)
	replaceZoneOptions.SetAccountID(*currentZone.AccountID)
	if currentZone.Description != nil {
		replaceZoneOptions.SetDescription(*currentZone.Description)
	}
	if currentZone.Excluded != nil {
		replaceZoneOptions.SetExcluded(currentZone.Excluded)

	}

	addresses := []contextbasedrestrictionsv1.AddressIntf{}
	if _, ok := d.GetOk("addresses"); ok {
		addresses, err = ResourceEncodeAddressList(d.Get("addresses").([]interface{}), addressesId)
		if err != nil {
			return err
		}
	}
	preservedAddresses := FilterAddressList(currentZone.Addresses, func(id string) bool {
		return id != addressesId
	})
	if len(preservedAddresses) > 0 {
		addresses = append(preservedAddresses, addresses...)
	}
	replaceZoneOptions.SetAddresses(addresses)

	_, response, err = contextBasedRestrictionsClient.ReplaceZoneWithContext(context, replaceZoneOptions)
	if err != nil {
		return fmt.Errorf("ReplaceZoneWithContext failed %s\n%s", err, response)
	}

	return nil
}

func composeZoneAddressesId(zoneId, addressesId string) (id string) {
	id = fmt.Sprintf("%s/%s", zoneId, addressesId)
	return
}

func decomposeZoneAddressesId(id string) (zoneId, addressesId string) {
	if index := strings.Index(id, "/"); index >= 0 {
		zoneId = id[:index]
		addressesId = id[index+1:]
	}
	return
}
