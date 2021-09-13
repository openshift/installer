// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/dnssvcsv1"

	"github.com/IBM/go-sdk-core/v5/core"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
)

func resouceIBMPrivateDNSCustomResolver() *schema.Resource {
	return &schema.Resource{
		Create:   resouceIBMPrivateDNSCustomResolverCreate,
		Read:     resouceIBMPrivateDNSCustomResolverRead,
		Update:   resouceIBMPrivateDNSCustomResolverUpdate,
		Delete:   resouceIBMPrivateDNSCustomResolverDelete,
		Exists:   resouceIBMPrivateDNSCustomResolverExists,
		Importer: &schema.ResourceImporter{},

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
			pdnsCRHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Healthy state of the custom resolver",
			},
			pdnsCustomResolverLocations: {
				Type:             schema.TypeSet,
				Description:      "Locations on which the custom resolver will be running",
				Required:         true,
				DiffSuppressFunc: applyOnce,
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

func resouceIBMPrivateDNSCustomResolverCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
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

	crLocations := d.Get(pdnsCustomResolverLocations).(*schema.Set)
	customResolverOption := sess.NewCreateCustomResolverOptions(crn)
	customResolverOption.SetName(crName)
	customResolverOption.SetDescription(crDescription)
	customResolverOption.SetLocations(expandPdnsCRLocations(crLocations))

	result, resp, err := sess.CreateCustomResolverWithContext(context.TODO(), customResolverOption)
	if err != nil || result == nil {
		return fmt.Errorf("Error reading the  custom resolver %s:%s", err, resp)
	}

	d.SetId(convertCisToTfTwoVar(*result.ID, crn))
	d.Set(pdnsCRId, *result.ID)

	_, err = waitForPDNSCustomResolverHealthy(d, meta)
	if err != nil {
		return err
	}

	// Enable Custrom Resolver
	return resouceIBMPrivateDNSCustomResolverUpdate(d, meta)
}

func resouceIBMPrivateDNSCustomResolverRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	customResolverID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}

	opt := sess.NewGetCustomResolverOptions(crn, customResolverID)
	result, response, err := sess.GetCustomResolverWithContext(context.TODO(), opt)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the  custom resolver %s:%s", err, response)
	}
	d.Set(pdnsInstanceID, crn)
	d.Set(pdnsCRId, *result.ID)
	d.Set(pdnsCRName, *result.Name)
	d.Set(pdnsCRDescription, *result.Description)
	d.Set(pdnsCRHealth, *result.Health)
	d.Set(pdnsCREnabled, *result.Enabled)
	d.Set(pdnsCustomResolverLocations, flattenPdnsCRLocations(result.Locations))

	return nil
}

func resouceIBMPrivateDNSCustomResolverUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	customResolverID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
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

		result, resp, err := sess.UpdateCustomResolverWithContext(context.TODO(), opt)
		if err != nil || result == nil {
			return fmt.Errorf("Error updating the  custom resolver %s:%s", err, resp)
		}

	}

	return resouceIBMPrivateDNSCustomResolverRead(d, meta)
}

func resouceIBMPrivateDNSCustomResolverDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	customResolverID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}

	// Disable Cutsom Resolver before deleting
	optEnabled := sess.NewUpdateCustomResolverOptions(crn, customResolverID)
	optEnabled.SetEnabled(false)
	result, resp, errEnabled := sess.UpdateCustomResolverWithContext(context.TODO(), optEnabled)
	if err != nil || result == nil {
		return fmt.Errorf("Error updating the  custom resolver to disable before deleting %s:%s", errEnabled, resp)
	}

	opt := sess.NewDeleteCustomResolverOptions(crn, customResolverID)
	response, err := sess.DeleteCustomResolverWithContext(context.TODO(), opt)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error deleting the  custom resolver %s:%s", err, response)
	}

	d.SetId("")
	return nil
}

func resouceIBMPrivateDNSCustomResolverExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}

	customResolverID, crn, err := convertTftoCisTwoVar(d.Id())
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
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return nil, err
	}
	customResolverID, crn, _ := convertTftoCisTwoVar(d.Id())
	opt := sess.NewGetCustomResolverOptions(crn, customResolverID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{pdnsCustomResolverCritical},
		Target:  []string{pdnsCustomResolverDegraded, pdnsCustomResolverHealthy},
		Refresh: func() (interface{}, string, error) {
			res, detail, err := sess.GetCustomResolver(opt)
			if err != nil {
				if detail != nil && detail.StatusCode == 404 {
					return nil, "", fmt.Errorf("The custom resolver %s does not exist anymore: %v", customResolverID, err)
				}
				return nil, "", fmt.Errorf("Get the custom resolver %s failed with resp code: %s, err: %v", customResolverID, detail, err)
			}
			return res, *res.Health, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}
