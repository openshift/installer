// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	pdnsSecondaryZoneID           = "secondary_zone_id"
	pdnsSecondaryZoneZone         = "zone"
	pdnsSecondaryZoneTransferFrom = "transfer_from"
	pdnsSecondaryZoneEnabled      = "enabled"
	pdnsSecondaryZoneDescription  = "description"
	pdnsSecondaryZoneCreatedOn    = "created_on"
	pdnsSecondaryZoneModifiedOn   = "modified_on"
	pdnsResolverID                = "resolver_id"
)

func ResourceIBMPrivateDNSSecondaryZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPrivateDNSSecondaryZoneCreate,
		ReadContext:   resourceIBMPrivateDNSSecondaryZoneRead,
		UpdateContext: resourceIBMPrivateDNSSecondaryZoneUpdate,
		DeleteContext: resourceIBMPrivateDNSSecondaryZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a service instance.",
			},
			pdnsResolverID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a custom resolver.",
			},
			pdnsSecondaryZoneID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary Zone ID",
			},

			pdnsSecondaryZoneZone: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the zone.",
			},

			pdnsSecondaryZoneTransferFrom: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The addresses of DNS servers where the secondary zone data should be transferred from",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			pdnsSecondaryZoneEnabled: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable/Disable the secondary zone",
			},

			pdnsSecondaryZoneDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the secondary zone",
			},

			pdnsSecondaryZoneCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary Zone Creation date",
			},

			pdnsSecondaryZoneModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary Zone Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSSecondaryZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsResolverID).(string)
	description := d.Get(pdnsSecondaryZoneDescription).(string)
	zone := d.Get(pdnsSecondaryZoneZone).(string)
	enabled := d.Get(pdnsSecondaryZoneEnabled).(bool)
	transferFrom := flex.ExpandStringList(d.Get(pdnsSecondaryZoneTransferFrom).([]interface{}))

	createSecondaryZoneOptions := sess.NewCreateSecondaryZoneOptions(instanceID, resolverID, zone, transferFrom)

	createSecondaryZoneOptions.SetDescription(description)
	createSecondaryZoneOptions.SetEnabled(enabled)

	mk := "private_dns_secondary_zone_" + instanceID + resolverID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	resource, response, err := sess.CreateSecondaryZone(createSecondaryZoneOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error creating DNS Services secondary zone:%s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, resolverID, *resource.ID))
	return resourceIBMPrivateDNSSecondaryZoneRead(ctx, d, meta)
}

func resourceIBMPrivateDNSSecondaryZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 3 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/resolverID/secondaryZoneID", d.Id()))
	}
	instanceID := idSet[0]
	resolverID := idSet[1]
	secondaryZoneID := idSet[2]
	getSecondaryZoneOptions := sess.NewGetSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)
	resource, response, err := sess.GetSecondaryZone(getSecondaryZoneOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading DNS Services secondary zone:%s\n%s", err, response))
	}

	transferFrom := []string{}
	for _, value := range resource.TransferFrom {
		values := strings.Split(value, ":")
		transferFrom = append(transferFrom, values[0])
	}
	d.Set(pdnsInstanceID, idSet[0])
	d.Set(pdnsResolverID, idSet[1])
	d.Set(pdnsSecondaryZoneDescription, *resource.Description)
	d.Set(pdnsSecondaryZoneZone, *resource.Zone)
	d.Set(pdnsSecondaryZoneTransferFrom, transferFrom)
	d.Set(pdnsSecondaryZoneID, *resource.ID)
	d.Set(pdnsSecondaryZoneCreatedOn, resource.CreatedOn.String())
	d.Set(pdnsSecondaryZoneModifiedOn, resource.ModifiedOn.String())
	d.Set(pdnsSecondaryZoneEnabled, *resource.Enabled)

	return nil
}

func resourceIBMPrivateDNSSecondaryZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 3 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/resolverID/secondaryZoneID", d.Id()))
	}
	instanceID := idSet[0]
	resolverID := idSet[1]
	secondaryZoneID := idSet[2]

	// Check DNS zone is present
	getZoneOptions := sess.NewGetSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)
	_, response, err := sess.GetSecondaryZone(getZoneOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error fetching secondary zone:%s\n%s", err, response))
	}

	// Update DNS zone if attributes has any change
	if d.HasChange(pdnsSecondaryZoneTransferFrom) ||
		d.HasChange(pdnsSecondaryZoneDescription) ||
		d.HasChange(pdnsSecondaryZoneEnabled) {
		updateSecondaryZoneOptions := sess.NewUpdateSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)
		transferFrom := flex.ExpandStringList(d.Get(pdnsSecondaryZoneTransferFrom).([]interface{}))
		description := d.Get(pdnsSecondaryZoneDescription).(string)
		enabled := d.Get(pdnsSecondaryZoneEnabled).(bool)
		updateSecondaryZoneOptions.SetTransferFrom(transferFrom)
		updateSecondaryZoneOptions.SetDescription(description)
		updateSecondaryZoneOptions.SetEnabled(enabled)

		mk := "private_dns_secondary_zone_" + instanceID + resolverID
		conns.IbmMutexKV.Lock(mk)
		defer conns.IbmMutexKV.Unlock(mk)

		_, response, err := sess.UpdateSecondaryZone(updateSecondaryZoneOptions)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating DNS Services zone:%s\n%s", err, response))
		}
	}

	return resourceIBMPrivateDNSSecondaryZoneRead(ctx, d, meta)
}

func resourceIBMPrivateDNSSecondaryZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 3 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/resolverID/secondaryZoneID", d.Id()))
	}
	instanceID := idSet[0]
	resolverID := idSet[1]
	secondaryZoneID := idSet[2]
	deleteSecondaryZoneOptions := sess.NewDeleteSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)

	mk := "private_dns_secondary_zone_" + instanceID + resolverID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)
	response, err := sess.DeleteSecondaryZone(deleteSecondaryZoneOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading DNS Services secondary zone:%s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
