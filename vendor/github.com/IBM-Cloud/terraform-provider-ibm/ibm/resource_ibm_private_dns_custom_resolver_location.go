// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCRLocation           = "ibm_dns_custom_resolver_location"
	pdnsResolverID          = "resolver_id"
	pdnsCRLocationID        = "location_id"
	pdnsCRLocationSubnetCRN = "subnet_crn"
	pdnsCRLocationEnable    = "enabled"
	pdnsCRLocationServerIP  = "dns_server_ip"
)

func resourceIBMPrivateDNSCRLocation() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSLocationCreate,
		Read:     resourceIBMPrivateDNSLocationRead,
		Update:   resourceIBMPrivateDNSLocationUpdate,
		Delete:   resourceIBMPrivateDNSLocationDelete,
		Importer: &schema.ResourceImporter{},
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
		},
	}
}
func resourceIBMPrivateDNSLocationCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsResolverID).(string)

	opt := sess.NewAddCustomResolverLocationOptions(instanceID, resolverID)

	if subnetcrn, ok := d.GetOk(pdnsCRLocationSubnetCRN); ok {
		opt.SetSubnetCrn(subnetcrn.(string))
	}
	if enable, ok := d.GetOkExists(pdnsCRLocationEnable); ok {
		opt.SetEnabled(enable.(bool))
	}
	result, resp, err := sess.AddCustomResolverLocationWithContext(context.TODO(), opt)

	if err != nil || result == nil {
		return fmt.Errorf("Error creating the custom resolver location %s:%s", err, resp)
	}
	d.SetId(convertCisToTfThreeVar(*result.ID, resolverID, instanceID))
	return resourceIBMPrivateDNSLocationRead(d, meta)
}

func resourceIBMPrivateDNSLocationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceIBMPrivateDNSLocationUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	locationID, resolverID, instanceID, err := convertTfToCisThreeVar(d.Id())
	updatelocation := sess.NewUpdateCustomResolverLocationOptions(instanceID, resolverID, locationID)

	if d.HasChange(pdnsCRLocationSubnetCRN) ||
		d.HasChange(pdnsCRLocationEnable) {
		if scrn, ok := d.GetOk(pdnsCRLocationSubnetCRN); ok {
			updatelocation.SetSubnetCrn(scrn.(string))
		}
		if e, ok := d.GetOkExists(pdnsCRLocationEnable); ok {
			updatelocation.SetEnabled(e.(bool))
		}
		result, resp, err := sess.UpdateCustomResolverLocationWithContext(context.TODO(), updatelocation)
		if err != nil || result == nil {
			return fmt.Errorf("Error updating the custom resolver location %s:%s", err, resp)
		}
	}
	return resourceIBMPrivateDNSLocationRead(d, meta)
}
func resourceIBMPrivateDNSLocationDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	locationID, resolverID, instanceID, err := convertTfToCisThreeVar(d.Id())
	deleteCRlocation := sess.NewDeleteCustomResolverLocationOptions(instanceID, resolverID, locationID)
	resp, errDel := sess.DeleteCustomResolverLocationWithContext(context.TODO(), deleteCRlocation)
	if errDel != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error deleting the custom resolver location %s:%s", errDel, resp)
	}
	d.SetId("")
	return nil
}
