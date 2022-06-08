// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmDNSCustomResolver        = "ibm_dns_custom_resolver"
	pdnsCustomResolvers         = "custom_resolvers"
	pdnsCustomResolverLocations = "locations"
	pdnsCRId                    = "custom_resolver_id"
	pdnsCRName                  = "name"
	pdnsCRDescription           = "description"
	pdnsCRHealth                = "health"
	pdnsCREnabled               = "enabled"
	pdnsCRCreatedOn             = "created_on"
	pdnsCRModifiedOn            = "modified_on"
	pdnsCRLocationId            = "location_id"
	pdnsCRLocationSubnetCrn     = "subnet_crn"
	pdnsCRLocationEnabled       = "enabled"
	pdnsCRLocationHealthy       = "healthy"
	pdnsCRLocationDnsServerIp   = "dns_server_ip"
	pdnsCustomResolverCritical  = "CRITICAL"
	pdnsCustomResolverDegraded  = "DEGRADED"
	pdnsCustomResolverHealthy   = "HEALTHY"
	pdnsCRHighAvailability      = "high_availability"
)

func ResourceIBMPrivateDNSCustomResolver() *schema.Resource {
	return &schema.Resource{
		CreateContext: resouceIBMPrivateDNSCustomResolverCreate,
		ReadContext:   resouceIBMPrivateDNSCustomResolverRead,
		UpdateContext: resouceIBMPrivateDNSCustomResolverUpdate,
		DeleteContext: resouceIBMPrivateDNSCustomResolverDelete,
		Exists:        resouceIBMPrivateDNSCustomResolverExists,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID",
			},

			pdnsCRId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of the custom resolver",
			},
			pdnsCRName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the custom resolver",
			},
			pdnsCRDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the custom resolver.",
			},
			pdnsCREnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the custom resolver is enabled",
			},
			pdnsCRHighAvailability: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Whether High Availability is enabled in custom resolver",
			},
			pdnsCRHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Healthy state of the custom resolver",
			},
			pdnsCustomResolverLocations: {
				Type:             schema.TypeSet,
				Description:      "Locations on which the custom resolver will be running",
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsCRLocationId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location ID",
						},
						pdnsCRLocationSubnetCrn: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet CRN",
						},
						pdnsCRLocationEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the location is enabled for the custom resolver",
						},
						pdnsCRLocationHealthy: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the DNS server in this location is healthy or not.",
						},
						pdnsCRLocationDnsServerIp: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ip address of this dns server",
						},
					},
				},
			},
			pdnsCRForwardRules: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsCRFRRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the forwarding rule.",
						},
						pdnsCRFRDesctiption: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive text of the forwarding rule.",
						},
						pdnsCRFRType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the forwarding rule.",
						},
						pdnsCRFRMatch: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The matching zone or hostname.",
						},
						pdnsCRFRForwardTo: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The upstream DNS servers will be forwarded to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			pdnsCRCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when a custom resolver is created",
			},

			pdnsCRModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The recent time when a custom resolver is modified",
			},
		},
	}
}

func resouceIBMPrivateDNSCustomResolverCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	var crName, crDescription string

	// session options
	crn := d.Get(pdnsInstanceID).(string)
	if name, ok := d.GetOk(pdnsCRName); ok {
		crName = name.(string)
	}
	if des, ok := d.GetOk(pdnsCRDescription); ok {
		crDescription = des.(string)
	}

	customResolverOption := sess.NewCreateCustomResolverOptions(crn)
	customResolverOption.SetName(crName)
	customResolverOption.SetDescription(crDescription)

	cr_highaval := d.Get(pdnsCRHighAvailability).(bool)

	crLocationCreate := false
	if _, ok := d.GetOk(pdnsCustomResolverLocations); ok {
		crLocationCreate = true
		crLocations := d.Get(pdnsCustomResolverLocations).(*schema.Set)
		if cr_highaval && crLocations.Len() <= 1 {
			return diag.FromErr(fmt.Errorf("To meet high availability status, configure custom resolvers with a minimum of two resolver locations. A maximum of four locations can be configured within the same subnet location."))
		}
		customResolverOption.SetLocations(expandPdnsCRLocations(crLocations))
	} else {
		if cr_highaval {
			return diag.FromErr(fmt.Errorf("To meet high availability status, configure custom resolvers with a minimum of two resolver locations. A maximum of four locations can be configured within the same subnet location."))
		}
	}

	result, resp, err := sess.CreateCustomResolverWithContext(context, customResolverOption)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %s:%s", err, resp))
	}

	d.SetId(flex.ConvertCisToTfTwoVar(*result.ID, crn))
	d.Set(pdnsCRId, *result.ID)

	if crLocationCreate {
		_, err = waitForPDNSCustomResolverHealthy(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		return resouceIBMPrivateDNSCustomResolverUpdate(context, d, meta)
	}
	return resouceIBMPrivateDNSCustomResolverRead(context, d, meta)
}

func resouceIBMPrivateDNSCustomResolverRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	customResolverID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	opt := sess.NewGetCustomResolverOptions(crn, customResolverID)
	result, response, err := sess.GetCustomResolverWithContext(context, opt)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %s:%s", err, response))
	}
	fwopt := sess.NewListForwardingRulesOptions(crn, customResolverID)

	fwresult, fwresp, fwerr := sess.ListForwardingRulesWithContext(context, fwopt)
	if fwerr != nil || fwresult == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error listing the forwarding rules %s:%s", fwerr, fwresp))
	}

	forwardRules := make([]interface{}, 0)
	for _, instance := range fwresult.ForwardingRules {
		forwardRule := map[string]interface{}{}
		forwardRule[pdnsCRFRRuleID] = *instance.ID
		forwardRule[pdnsCRFRDesctiption] = *instance.Description
		forwardRule[pdnsCRFRType] = *instance.Type
		forwardRule[pdnsCRFRMatch] = *instance.Match
		forwardRule[pdnsCRFRForwardTo] = instance.ForwardTo
		forwardRules = append(forwardRules, forwardRule)
	}
	d.Set(pdnsInstanceID, crn)
	d.Set(pdnsCRId, *result.ID)
	d.Set(pdnsCRName, *result.Name)
	d.Set(pdnsCRDescription, *result.Description)
	d.Set(pdnsCRHealth, *result.Health)
	d.Set(pdnsCREnabled, *result.Enabled)
	d.Set(pdnsCustomResolverLocations, flattenPdnsCRLocations(result.Locations))
	d.Set(pdnsCRForwardRules, forwardRules)
	return nil
}

func resouceIBMPrivateDNSCustomResolverUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	customResolverID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(pdnsCRName) ||
		d.HasChange(pdnsCRDescription) ||
		d.HasChange(pdnsCREnabled) {

		opt := sess.NewUpdateCustomResolverOptions(crn, customResolverID)
		if name, ok := d.GetOk(pdnsCRName); ok {
			crName := name.(string)
			opt.SetName(crName)
		}
		if des, ok := d.GetOk(pdnsCRDescription); ok {
			crDescription := des.(string)
			opt.SetDescription(crDescription)
		}
		if enabled, ok := d.GetOkExists(pdnsCREnabled); ok {
			crEnabled := enabled.(bool)
			opt.SetEnabled(crEnabled)
		}

		result, resp, err := sess.UpdateCustomResolverWithContext(context, opt)
		if err != nil || result == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating the custom resolver %s:%s", err, resp))
		}

	}

	return resouceIBMPrivateDNSCustomResolverRead(context, d, meta)
}

func resouceIBMPrivateDNSCustomResolverDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	customResolverID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	// Disable Cutsom Resolver before deleting
	optEnabled := sess.NewUpdateCustomResolverOptions(crn, customResolverID)
	optEnabled.SetEnabled(false)
	result, resp, errEnabled := sess.UpdateCustomResolverWithContext(context, optEnabled)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error updating the custom resolver to disable before deleting %s:%s", errEnabled, resp))
	}

	opt := sess.NewDeleteCustomResolverOptions(crn, customResolverID)
	response, err := sess.DeleteCustomResolverWithContext(context, opt)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error deleting the custom resolver %s:%s", err, response))
	}

	d.SetId("")
	return nil
}

func resouceIBMPrivateDNSCustomResolverExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}

	customResolverID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return false, err
	}
	opt := sess.NewGetCustomResolverOptions(crn, customResolverID)
	_, response, err := sess.GetCustomResolver(opt)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Custom Resolver does not exist.")
			return false, nil
		}
		log.Printf("Custom Resolver failed: %v", response)
		return false, err
	}
	return true, nil
}

func flattenPdnsCRLocations(crLocation []dnssvcsv1.Location) interface{} {
	flattened := make([]interface{}, 0)
	for _, v := range crLocation {
		customLocations := map[string]interface{}{
			pdnsCRLocationId:          *v.ID,
			pdnsCRLocationSubnetCrn:   *v.SubnetCrn,
			pdnsCRLocationEnabled:     *v.Enabled,
			pdnsCRLocationHealthy:     *v.Healthy,
			pdnsCRLocationDnsServerIp: *v.DnsServerIp,
		}
		flattened = append(flattened, customLocations)
	}
	return flattened
}

func expandPdnsCRLocations(crLocList *schema.Set) (crLocations []dnssvcsv1.LocationInput) {
	for _, iface := range crLocList.List() {
		var locOpt dnssvcsv1.LocationInput
		loc := iface.(map[string]interface{})
		locOpt.SubnetCrn = core.StringPtr(loc[pdnsCRLocationSubnetCrn].(string))
		if val, ok := loc[pdnsCRLocationEnabled]; ok {
			locOpt.Enabled = core.BoolPtr(val.(bool))
		}
		crLocations = append(crLocations, locOpt)
	}
	return
}

func waitForPDNSCustomResolverHealthy(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return nil, err
	}

	var customResolverID, crn string

	g := strings.SplitN(d.Id(), ":", -1)
	if len(g) > 2 {
		_, customResolverID, crn, _ = flex.ConvertTfToCisThreeVar(d.Id())
	} else {
		customResolverID, crn, _ = flex.ConvertTftoCisTwoVar(d.Id())
	}

	opt := sess.NewGetCustomResolverOptions(crn, customResolverID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{pdnsCustomResolverCritical, "false"},
		Target:  []string{pdnsCustomResolverDegraded, pdnsCustomResolverHealthy, "true"},
		Refresh: func() (interface{}, string, error) {
			res, detail, err := sess.GetCustomResolver(opt)
			if err != nil {
				if detail != nil && detail.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] The custom resolver %s does not exist anymore: %v", customResolverID, err)
				}
				return nil, "", fmt.Errorf("Get the custom resolver %s failed with resp code: %s, err: %v", customResolverID, detail, err)
			}
			return res, *res.Health, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}
