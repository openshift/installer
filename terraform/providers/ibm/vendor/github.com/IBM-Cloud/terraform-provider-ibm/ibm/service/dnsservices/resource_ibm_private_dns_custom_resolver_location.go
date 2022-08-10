// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCRLocation            = "ibm_dns_custom_resolver_location"
	pdnsResolverID           = "resolver_id"
	pdnsCRLocationID         = "location_id"
	pdnsCRLocationSubnetCRN  = "subnet_crn"
	pdnsCRLocationEnable     = "enabled"
	pdnsCRLocationServerIP   = "dns_server_ip"
	pdnsCustomReolverEnabled = "cr_enabled"
)

func ResourceIBMPrivateDNSCRLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPrivateDNSLocationCreate,
		ReadContext:   resourceIBMPrivateDNSLocationRead,
		UpdateContext: resourceIBMPrivateDNSLocationUpdate,
		DeleteContext: resourceIBMPrivateDNSLocationDelete,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID",
			},

			pdnsResolverID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Resolver ID",
			},
			pdnsCRLocationID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRLocation ID",
			},

			pdnsCRLocationSubnetCRN: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CRLocation Subnet CRN",
			},

			pdnsCRLocationEnable: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "CRLocation Enabled",
			},

			pdnsCRLocationHealthy: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "CRLocation Healthy",
			},

			pdnsCRLocationServerIP: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRLocation Server IP",
			},
			pdnsCustomReolverEnabled: {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: flex.ApplyOnce,
			},
		},
		DeprecationMessage: "Resource ibm_dns_custom_resolver_location is deprecated. Using the deprecated resource can cause an outage. If you have used the `ibm_dns_custom_resolver_location` resource, change it to the composite Custom Resolver [ibm_dns_custom_resolver] resource before running terraform apply.",
	}
}
func resourceIBMPrivateDNSLocationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsResolverID).(string)

	mk := "private_dns_resource_custom_resolver_location_" + instanceID + resolverID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	opt := sess.NewAddCustomResolverLocationOptions(instanceID, resolverID)

	if subnetcrn, ok := d.GetOk(pdnsCRLocationSubnetCRN); ok {
		opt.SetSubnetCrn(subnetcrn.(string))
	}
	var enable_loc, cr_enable bool
	if enable, ok := d.GetOkExists(pdnsCRLocationEnable); ok {
		opt.SetEnabled(enable.(bool))
		enable_loc = enable.(bool)
	}
	if enable_cr, ok := d.GetOkExists(pdnsCustomReolverEnabled); ok {
		cr_enable = enable_cr.(bool)
	}
	// if location enabled is false, CR cannot be enabled, fail here
	if !enable_loc && cr_enable {
		return diag.FromErr(fmt.Errorf("[ERROR]The custom resolver location is not enabled, hence cannot add location. Also, cannot enable the custom resolver"))
	}

	if enable_loc {
		// Fetch the Custom Resolver and check the enabled attribute
		optCr := sess.NewGetCustomResolverOptions(instanceID, resolverID)
		cr_result, response, err := sess.GetCustomResolver(optCr)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %v:%v", err, response))
			}
			return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %v:%v", err, response))
		}
		// Disable the Custom Resolver location and add it, otherwise you will get API error:
		// "Not allowed to create enabled location while custom resolver is enabled."
		if *cr_result.Enabled {
			opt.SetEnabled(false)
		}
	}

	result, resp, err := sess.AddCustomResolverLocationWithContext(context, opt)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error creating the custom resolver location %s:%s", err, resp))
	}
	locationID := *result.ID
	d.SetId(flex.ConvertCisToTfThreeVar(locationID, resolverID, instanceID))

	if cr_enable && enable_loc {
		err := PDNSCustomResolverEnableLocation(meta, instanceID, resolverID, locationID)
		if err != nil {
			return err
		}
	}

	if cr_enable && enable_loc {
		err := PDNSCustomResolverEnable(meta, instanceID, resolverID)
		if err != nil {
			return err
		}
	} else if !cr_enable {
		optCr := sess.NewUpdateCustomResolverOptions(instanceID, resolverID)
		optCr.SetEnabled(false)
		resultCr, respCr, errCr := sess.UpdateCustomResolverWithContext(context, optCr)
		if errCr != nil || resultCr == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating the custom resolver with cr_enable false %s:%s", errCr, respCr))
		}
	}
	return resourceIBMPrivateDNSLocationRead(context, d, meta)
}

func resourceIBMPrivateDNSLocationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceIBMPrivateDNSLocationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	locationID, resolverID, instanceID, err := flex.ConvertTfToCisThreeVar(d.Id())

	mk := "private_dns_resource_custom_resolver_location_" + instanceID + resolverID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, resolverID, locationID)

	if d.HasChange(pdnsCRLocationSubnetCRN) ||
		d.HasChange(pdnsCRLocationEnable) {
		if scrn, ok := d.GetOk(pdnsCRLocationSubnetCRN); ok {
			updatelocation.SetSubnetCrn(scrn.(string))
		}
		if e, ok := d.GetOkExists(pdnsCRLocationEnable); ok {
			updatelocation.SetEnabled(e.(bool))
		}
		result, resp, err := sess.UpdateCustomResolverLocationWithContext(context, updatelocation)
		if err != nil || result == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating the custom resolver location %s:%s", err, resp))
		}
	}
	return resourceIBMPrivateDNSLocationRead(context, d, meta)
}

func resourceIBMPrivateDNSLocationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	locationID, resolverID, instanceID, _ := flex.ConvertTfToCisThreeVar(d.Id())
	opt := sess.NewGetCustomResolverOptions(instanceID, resolverID)
	cr_result, response, err := sess.GetCustomResolver(opt)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %v:%v", err, response))
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading the custom resolver %v:%v", err, response))
	}

	crLocations := cr_result.Locations
	locationIdFound := false
	enabledLocation := false
	enabledAnotherLocation := false
	for _, v := range crLocations {
		if *v.ID == locationID {
			locationIdFound = true
			if *v.Enabled {
				enabledLocation = true
			}
		} else {
			if *v.Enabled {
				enabledAnotherLocation = true
			}
		}
	}

	if !locationIdFound {
		d.SetId("")
		return nil
	}
	// Location is enabled
	if enabledLocation {
		if *cr_result.Enabled {
			if !enabledAnotherLocation || (len(crLocations) == 1) {
				// Custom Resolver is Enabled, fail here
				return diag.FromErr(fmt.Errorf("[ERROR] Error Deleting the custom resolver location. Custom resolver is enabled, it needs atleast one enabled location"))
			}
		}
		// Disable the Custom Resolver Location
		updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, resolverID, locationID)
		updatelocation.SetEnabled(false)
		result, resp, err := sess.UpdateCustomResolverLocationWithContext(context, updatelocation)
		if err != nil || result == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Disable and updating the custom resolver location %s:%s", err, resp))
		}
	}
	deleteCRlocation := sess.NewDeleteCustomResolverLocationOptions(instanceID, resolverID, locationID)
	resp, errDel := sess.DeleteCustomResolverLocationWithContext(context, deleteCRlocation)
	if errDel != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error Deleting the custom resolver location %s:%s", errDel, resp))
	}
	d.SetId("")
	return nil
}
