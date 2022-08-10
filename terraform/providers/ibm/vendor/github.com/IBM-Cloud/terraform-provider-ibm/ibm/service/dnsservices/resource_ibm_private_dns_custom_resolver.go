// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Description: "Whether High Availability is enabled in custom resolver",
			},
			pdnsCRHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Healthy state of the custom resolver",
			},
			pdnsCustomResolverLocations: {
				Type:        schema.TypeList,
				Description: "Locations on which the custom resolver will be running",
				Optional:    true,
				MaxItems:    3,
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

type location struct {
	locationId string
	subnet     string
	enabled    bool
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

	var loc_enable bool
	cr_enable := d.Get(pdnsCREnabled)

	// Validation
	if _, ok := d.GetOk(pdnsCustomResolverLocations); ok {
		var expandcrLocations []dnssvcsv1.LocationInput
		crLocations := d.Get(pdnsCustomResolverLocations).([]interface{})
		if len(crLocations) > 3 {
			return diag.FromErr(fmt.Errorf("A custom resolver can have a maximum of three locations, either within the same subnet or in different subnets."))
		}
		if cr_highaval && len(crLocations) <= 1 {
			return diag.FromErr(fmt.Errorf("To meet high availability status, configure custom resolvers with a minimum of two resolver locations. A maximum of three locations can be configured within the same subnet location."))
		}
		expandcrLocations, loc_enable = expandPdnsCRLocations(crLocations)
		if cr_enable.(bool) && !loc_enable {
			return diag.FromErr(fmt.Errorf("The Custom resolver cannot be enabled. There should be atleast one enabled location."))
		}
		customResolverOption.SetLocations(expandcrLocations)
	} else {
		if cr_highaval {
			return diag.FromErr(fmt.Errorf("To meet high availability status, configure custom resolvers with a minimum of two resolver locations. A maximum of three locations can be configured within the same subnet location."))
		} else if cr_enable.(bool) {
			return diag.FromErr(fmt.Errorf("The Custom resolver cannot be enabled. There should be atleast one enabled location."))
		}
	}

	// Create a custom resolver
	result, resp, err := sess.CreateCustomResolverWithContext(context, customResolverOption)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %s:%s", err, resp))
	}

	d.SetId(flex.ConvertCisToTfTwoVar(*result.ID, crn))
	d.Set(pdnsCRId, *result.ID)

	// Enable Custom resolver
	if cr_enable.(bool) {
		err := PDNSCustomResolverEnable(meta, crn, *result.ID)
		if err != nil {
			return err
		}
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

	resolverID, instanceID, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var loc_enable, cr_enable, cr_highaval bool

	if enable_cr, ok := d.GetOk(pdnsCREnabled); ok {
		cr_enable = enable_cr.(bool)
	}
	if highaval, ok := d.GetOk(pdnsCRHighAvailability); ok {
		cr_highaval = highaval.(bool)
	}
	var oldRaw, newRaw interface{}

	if d.HasChange(pdnsCRName) ||
		d.HasChange(pdnsCRDescription) ||
		d.HasChange(pdnsCREnabled) ||
		d.HasChange(pdnsCRHighAvailability) {

		// Validation
		if _, ok := d.GetOk(pdnsCustomResolverLocations); ok {
			var expandcrLocations []dnssvcsv1.LocationInput
			crLocations := d.Get(pdnsCustomResolverLocations).([]interface{})
			if len(crLocations) > 3 {
				return diag.FromErr(fmt.Errorf("A custom resolver can have a maximum of three locations, either within the same subnet or in different subnets."))
			}
			if cr_highaval && len(crLocations) <= 1 {
				return diag.FromErr(fmt.Errorf("To meet high availability status, configure custom resolvers with a minimum of two resolver locations .A maximum of three locations can be configured within the same subnet location."))
			}
			expandcrLocations, loc_enable = expandPdnsCRLocations(crLocations)
			if cr_enable && !loc_enable {
				return diag.FromErr(fmt.Errorf("The Custom resolver cannot be enabled. There should be atleast one enabled location."))
			}
			fmt.Print("expandcrLocations", expandcrLocations)
		} else {
			if cr_highaval {
				return diag.FromErr(fmt.Errorf("To meet high availability status, configure custom resolvers with a minimum of two resolver locations. A maximum of three locations can be configured within the same subnet location."))
			} else if cr_enable {
				return diag.FromErr(fmt.Errorf("The Custom resolver cannot be enabled. There should be atleast one enabled location."))
			}
		}

		opt := sess.NewUpdateCustomResolverOptions(instanceID, resolverID)
		if name, ok := d.GetOk(pdnsCRName); ok {
			crName := name.(string)
			opt.SetName(crName)
		}
		if des, ok := d.GetOk(pdnsCRDescription); ok {
			crDescription := des.(string)
			opt.SetDescription(crDescription)
		}
		if !cr_enable {
			opt.SetEnabled(false)
		}
		result, resp, err := sess.UpdateCustomResolverWithContext(context, opt)
		if err != nil || result == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating the custom resolver %s:%s", err, resp))
		}

	}

	if d.HasChange(pdnsCustomResolverLocations) {

		oldRaw, newRaw = d.GetChange(pdnsCustomResolverLocations)

		newState := stateprocess(newRaw)
		oldState := stateprocess(oldRaw)

		// Delete Custom Resolver Location.
		for _, oldloc := range oldState {
			locationIdExists := false
			for _, newloc := range newState {
				if oldloc.locationId == newloc.locationId {
					locationIdExists = true
					break
				}
			}
			if !locationIdExists {
				if oldloc.enabled {
					err := PDNSCustomResolverDisableLocation(meta, instanceID, resolverID, oldloc.locationId)
					if err != nil {
						return err
					}
				}
				err := deleteCRLocation(meta, instanceID, resolverID, oldloc.locationId)
				if err != nil {
					return err
				}
			}
		}

		for _, newLoc := range newState {
			// Add new custom resolver locations
			if strings.Contains(newLoc.locationId, "NEW0") {
				locationID, err := addCRLocation(meta, instanceID, resolverID, newLoc.subnet)
				if err != nil || locationID == "" {
					return err
				}
				if newLoc.enabled {
					err := PDNSCustomResolverEnableLocation(meta, instanceID, resolverID, locationID)
					if err != nil {
						return err
					}
				}
			} else {
				// Update Location
				locationIdExists := false
				for _, oldLoc := range oldState {
					if oldLoc.locationId == newLoc.locationId {
						locationIdExists = true
						if !(oldLoc.subnet == newLoc.subnet) {
							// Update location subnet crn.
							// Disable location before changing the subnet.
							err := PDNSCustomResolverDisableLocation(meta, instanceID, resolverID, newLoc.locationId)
							if err != nil {
								return err
							}
							errSub := updateLocationSubnet(meta, instanceID, resolverID, newLoc.locationId, newLoc.subnet)
							if errSub != nil {
								return errSub
							}
							if newLoc.enabled {
								err := PDNSCustomResolverEnableLocation(meta, instanceID, resolverID, newLoc.locationId)
								if err != nil {
									return err
								}
							}
						} else if newLoc.enabled != oldLoc.enabled {
							// Update location enable/disable
							if newLoc.enabled {
								err := PDNSCustomResolverEnableLocation(meta, instanceID, resolverID, newLoc.locationId)
								if err != nil {
									return err
								}
							} else {
								err := PDNSCustomResolverDisableLocation(meta, instanceID, resolverID, newLoc.locationId)
								if err != nil {
									return err
								}
							}
						}
					}
				}
				if !locationIdExists {
					return diag.FromErr(fmt.Errorf("[ERROR] The custom resolver location %s does not exist anymore: %v", newLoc.locationId, err))
				}
			}
		}
	}

	if d.HasChange(pdnsCREnabled) {
		if cr_enable {
			err := PDNSCustomResolverEnable(meta, instanceID, resolverID)
			if err != nil {
				return err
			}
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

func expandPdnsCRLocations(crLocList []interface{}) (crLocations []dnssvcsv1.LocationInput, loc_enable bool) {
	for _, iface := range crLocList {
		var locOpt dnssvcsv1.LocationInput
		loc := iface.(map[string]interface{})
		locOpt.SubnetCrn = core.StringPtr(loc[pdnsCRLocationSubnetCrn].(string))
		if val, ok := loc[pdnsCRLocationEnabled]; ok {
			if val.(bool) {
				loc_enable = true
			}
			locOpt.Enabled = core.BoolPtr(val.(bool))
		}
		crLocations = append(crLocations, locOpt)
	}
	return crLocations, loc_enable
}

func PDNSCustomResolverEnable(meta interface{}, instanceID string, customResolverID string) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	MaxTimeout := 600
	SleepTime := 20

	for SleepTime < MaxTimeout {
		opt := sess.NewUpdateCustomResolverOptions(instanceID, customResolverID)
		opt.SetEnabled(true)
		result, _, err := sess.UpdateCustomResolver(opt)
		if err != nil || result == nil {
			time.Sleep(20 * time.Second)
			SleepTime = SleepTime + 20
		} else {
			return nil
		}
	}
	return diag.FromErr(fmt.Errorf("[ERROR] Error Enabling the Custom resolver : MaxTimeout"))
}

func PDNSCustomResolverEnableLocation(meta interface{}, instanceID string, customResolverID string, locationID string) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	MaxTimeout := 600
	SleepTime := 20

	for SleepTime < MaxTimeout {
		updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, customResolverID, locationID)
		updatelocation.SetEnabled(true)
		result, _, err := sess.UpdateCustomResolverLocation(updatelocation)
		if err != nil || result == nil {
			time.Sleep(20 * time.Second)
			SleepTime = SleepTime + 20
		} else {
			return nil
		}
	}
	return diag.FromErr(fmt.Errorf("[ERROR] Error Enabling the Custom resolver location : MaxTimeout"))
}

func PDNSCustomResolverDisableLocation(meta interface{}, instanceID string, customResolverID string, locationID string) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, customResolverID, locationID)
	updatelocation.SetEnabled(false)
	result, resp, err := sess.UpdateCustomResolverLocation(updatelocation)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Disabling the custom resolver location %s:%s", err, resp))
	}
	return nil
}

func deleteCRLocation(meta interface{}, instanceID string, customResolverID string, locationID string) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	deleteCRlocation := sess.NewDeleteCustomResolverLocationOptions(instanceID, customResolverID, locationID)
	resp, errDel := sess.DeleteCustomResolverLocation(deleteCRlocation)
	if errDel != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error Deleting the custom resolver location %s:%s", errDel, resp))
	}
	return nil
}

func addCRLocation(meta interface{}, instanceID string, customResolverID string, subnet string) (string, diag.Diagnostics) {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return "", diag.FromErr(err)
	}
	opt := sess.NewAddCustomResolverLocationOptions(instanceID, customResolverID)
	opt.SetSubnetCrn(subnet)
	opt.SetEnabled(false)
	result, resp, err := sess.AddCustomResolverLocation(opt)
	locationID := *result.ID
	if err != nil || result == nil {
		return "", diag.FromErr(fmt.Errorf("[ERROR] Error creating the custom resolver location %s:%s", err, resp))
	}
	return locationID, nil
}

func updateLocationSubnet(meta interface{}, instanceID string, customResolverID string, locationID string, subnet string) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, customResolverID, locationID)
	updatelocation.SetSubnetCrn(subnet)
	updatelocation.SetEnabled(false)
	result, resp, err := sess.UpdateCustomResolverLocation(updatelocation)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Disable and updating the custom resolver location %s:%s", err, resp))
	}
	return nil
}

func stateprocess(Raw interface{}) (State []location) {
	new_LocationId := 0
	for _, loc := range Raw.([]interface{}) {
		temp_loc := loc.(map[string]interface{})
		newlocationId := (temp_loc["location_id"]).(string)
		newLocation := location{}
		// Add a constant marker for new locations.
		if len(newlocationId) == 0 {
			new_LocationId = new_LocationId + 1
			new_LocationName := "NEW0" + strconv.Itoa(new_LocationId)
			newLocation = location{locationId: new_LocationName, subnet: (temp_loc["subnet_crn"]).(string), enabled: temp_loc["enabled"].(bool)}
		} else {
			newLocation = location{locationId: newlocationId, subnet: (temp_loc["subnet_crn"]).(string), enabled: temp_loc["enabled"].(bool)}
		}
		State = append(State, newLocation)
	}
	return State
}
