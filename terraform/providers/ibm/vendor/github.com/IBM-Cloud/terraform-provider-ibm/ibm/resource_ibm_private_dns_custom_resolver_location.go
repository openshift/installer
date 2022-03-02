// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"

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

func resourceIBMPrivateDNSCRLocation() *schema.Resource {
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
				DiffSuppressFunc: applyOnce,
			},
		},
	}
}
func resourceIBMPrivateDNSLocationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsResolverID).(string)

	mk := "private_dns_resource_custom_resolver_location_" + instanceID + resolverID
	ibmMutexKV.Lock(mk)
	defer ibmMutexKV.Unlock(mk)

	opt := sess.NewAddCustomResolverLocationOptions(instanceID, resolverID)

	if subnetcrn, ok := d.GetOk(pdnsCRLocationSubnetCRN); ok {
		opt.SetSubnetCrn(subnetcrn.(string))
	}
	var enable_loc bool
	if enable, ok := d.GetOkExists(pdnsCRLocationEnable); ok {
		opt.SetEnabled(enable.(bool))
		enable_loc = enable.(bool)
	}
	if _, ok := d.GetOkExists(pdnsCustomReolverEnabled); ok {
		optCr := sess.NewUpdateCustomResolverOptions(instanceID, resolverID)
		optCr.SetEnabled(false)
		resultCr, respCr, errCr := sess.UpdateCustomResolverWithContext(context, optCr)
		if errCr != nil || resultCr == nil {
			return diag.FromErr(fmt.Errorf("Error updating the custom resolver with cr_enable false %s:%s", errCr, respCr))
		}
	}
	result, resp, err := sess.AddCustomResolverLocationWithContext(context, opt)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("Error creating the custom resolver location %s:%s", err, resp))
	}
	d.SetId(convertCisToTfThreeVar(*result.ID, resolverID, instanceID))

	if cr_enable, ok := d.GetOkExists(pdnsCustomReolverEnabled); ok {
		if cr_enable.(bool) && enable_loc {
			_, err = waitForPDNSCustomResolverHealthy(d, meta)
			if err != nil {
				return diag.FromErr(err)
			}
			optCr := sess.NewUpdateCustomResolverOptions(instanceID, resolverID)
			optCr.SetEnabled(cr_enable.(bool))
			resultCr, respCr, errCr := sess.UpdateCustomResolverWithContext(context, optCr)
			if errCr != nil || resultCr == nil {
				return diag.FromErr(fmt.Errorf("Error updating the custom resolver %s:%s", errCr, respCr))
			}
		}

	}
	return resourceIBMPrivateDNSLocationRead(context, d, meta)
}

func resourceIBMPrivateDNSLocationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceIBMPrivateDNSLocationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	locationID, resolverID, instanceID, err := convertTfToCisThreeVar(d.Id())

	mk := "private_dns_resource_custom_resolver_location_" + instanceID + resolverID
	ibmMutexKV.Lock(mk)
	defer ibmMutexKV.Unlock(mk)

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
			return diag.FromErr(fmt.Errorf("Error updating the custom resolver location %s:%s", err, resp))
		}
	}
	return resourceIBMPrivateDNSLocationRead(context, d, meta)
}
func resourceIBMPrivateDNSLocationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	locationID, resolverID, instanceID, err := convertTfToCisThreeVar(d.Id())
	delete_loc := true
	if cr_enable, ok := d.GetOkExists(pdnsCustomReolverEnabled); ok {
		if cr_enable.(bool) {
			// Disable the Cutsom Resolver
			optEnabled := sess.NewUpdateCustomResolverOptions(instanceID, resolverID)
			optEnabled.SetEnabled(false)

			result, resp, errEnabled := sess.UpdateCustomResolverWithContext(context, optEnabled)
			if err != nil || result == nil {
				return diag.FromErr(fmt.Errorf("Error Disable and Update the custom resolver %s:%s", errEnabled, resp))
			}
		} else {
			// Disable the Cutsom Resolver Location
			updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, resolverID, locationID)
			updatelocation.SetEnabled(false)
			result, resp, err := sess.UpdateCustomResolverLocationWithContext(context, updatelocation)
			if err != nil || result == nil {
				return diag.FromErr(fmt.Errorf("Error Disbale and updating the custom resolver location %s:%s", err, resp))
			}
		}
	}
	// Disable Cutsom Resolver Location before deleting
	updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, resolverID, locationID)
	updatelocation.SetEnabled(false)
	result, resp, err := sess.UpdateCustomResolverLocationWithContext(context, updatelocation)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("Error Disbale and updating the custom resolver location %s:%s", err, resp))
	}

	if delete_loc {
		opt := sess.NewGetCustomResolverOptions(instanceID, resolverID)
		result, _, _ := sess.GetCustomResolverWithContext(context, opt)
		if len(result.Locations) > 1 {
			deleteCRlocation := sess.NewDeleteCustomResolverLocationOptions(instanceID, resolverID, locationID)
			resp, errDel := sess.DeleteCustomResolverLocationWithContext(context, deleteCRlocation)
			if errDel != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil
				}
				return diag.FromErr(fmt.Errorf("Error Deleting the custom resolver location %s:%s", errDel, resp))
			}
		}
	}

	d.SetId("")
	return nil
}
