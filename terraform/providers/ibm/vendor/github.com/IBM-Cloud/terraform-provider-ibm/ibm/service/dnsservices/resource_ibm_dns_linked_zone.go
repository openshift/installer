// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	//"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DnsLinkedZoneInstanceID             = "instance_id"
	DnsLinkedZoneName                   = "name"
	DnsLinkedZoneDescription            = "description"
	DnsLinkedZoneOwnerInstanceID        = "owner_instance_id"
	DnsLinkedZoneOwnerZoneID            = "owner_zone_id"
	DnsLinkedZoneLinkedTo               = "linked_to"
	DnsLinkedZoneState                  = "state"
	DnsLinkedZoneLabel                  = "label"
	DnsLinkedZoneApprovalRequiredBefore = "approval_required_before"
	DnsLinkedZoneCreatedOn              = "created_on"
	DnsLinkedZoneModifiedOn             = "modified_on"

// DnsLinkedZoneOwnerInstanceID        = "owner_instance_id"
)

func ResourceIBMDNSLinkedZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMDNSLinkedZoneCreate,
		ReadContext:   resourceIBMDNSLinkedZoneRead,
		UpdateContext: resourceIBMDNSLinkedZoneUpdate,
		DeleteContext: resourceIBMDNSLinkedZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			DnsLinkedZoneInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a DNS Linked zone.",
			},
			DnsLinkedZoneName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DNS Linked zone.",
			},
			DnsLinkedZoneDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the DNS Linked zone",
			},
			DnsLinkedZoneOwnerInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the owner DNS instance",
			},
			DnsLinkedZoneOwnerZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the owner DNS zone",
			},
			DnsLinkedZoneLinkedTo: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone that is linked to the DNS Linked zone",
			},
			DnsLinkedZoneState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the DNS Linked zone",
			},
			DnsLinkedZoneLabel: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The label of the DNS Linked zone",
			},
			DnsLinkedZoneApprovalRequiredBefore: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS Linked Approval required before",
			},
			DnsLinkedZoneCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS Linked Zone Creation date",
			},
			DnsLinkedZoneModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS Linked Zone Modification date",
			},
		},
	}
}

func resourceIBMDNSLinkedZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	var (
		description string
		label       string
	)

	instanceID := d.Get(DnsLinkedZoneInstanceID).(string)
	if v, ok := d.GetOk(DnsLinkedZoneDescription); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk(DnsLinkedZoneLabel); ok {
		label = v.(string)
	}
	ownerInstanceID := d.Get(DnsLinkedZoneOwnerInstanceID).(string)
	ownerZoneID := d.Get(DnsLinkedZoneOwnerZoneID).(string)
	createLinkedZoneOptions := sess.NewCreateLinkedZoneOptions(instanceID, ownerInstanceID, ownerZoneID)

	createLinkedZoneOptions.SetDescription(description)
	createLinkedZoneOptions.SetLabel(label)
	mk := "dns_linked_zone_" + instanceID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	resource, response, err := sess.CreateLinkedZone(createLinkedZoneOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error creating DNS Linked zone:%s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, *resource.ID))
	return resourceIBMDNSLinkedZoneRead(ctx, d, meta)
}

func resourceIBMDNSLinkedZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 2 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/linkedDnsZoneID", d.Id()))
	}
	instanceID := idSet[0]
	linkedDnsZoneID := idSet[1]
	getLinkedZoneOptions := sess.NewGetLinkedZoneOptions(instanceID, linkedDnsZoneID)
	resource, response, err := sess.GetLinkedZone(getLinkedZoneOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading DNS Linked zone:%s\n%s", err, response))
	}

	d.Set(DnsLinkedZoneInstanceID, idSet[0])
	d.Set(DnsLinkedZoneDescription, *resource.Description)
	d.Set(DnsLinkedZoneLabel, *resource.Label)
	d.Set(DnsLinkedZoneCreatedOn, resource.CreatedOn.String())
	d.Set(DnsLinkedZoneModifiedOn, resource.ModifiedOn.String())

	return nil
}

func resourceIBMDNSLinkedZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 2 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/resolverID/secondaryZoneID", d.Id()))
	}
	instanceID := idSet[0]
	linkedDnsZoneID := idSet[1]

	// Check DNS zone is present
	getLinkedZoneOptions := sess.NewGetLinkedZoneOptions(instanceID, linkedDnsZoneID)
	_, response, err := sess.GetLinkedZone(getLinkedZoneOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error fetching secondary zone:%s\n%s", err, response))
	}

	// Update DNS Linked zone if attributes has any change
	if d.HasChange(DnsLinkedZoneDescription) ||
		d.HasChange(DnsLinkedZoneLabel) {
		updateLinkedZoneOptions := sess.NewUpdateLinkedZoneOptions(instanceID, linkedDnsZoneID)
		description := d.Get(DnsLinkedZoneDescription).(string)
		label := d.Get(DnsLinkedZoneLabel).(string)
		updateLinkedZoneOptions.SetDescription(description)
		updateLinkedZoneOptions.SetLabel(label)

		mk := "dns_linked_zone_" + instanceID
		conns.IbmMutexKV.Lock(mk)
		defer conns.IbmMutexKV.Unlock(mk)

		_, response, err := sess.UpdateLinkedZone(updateLinkedZoneOptions)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating DNS Services zone:%s\n%s", err, response))
		}
	}

	return resourceIBMDNSLinkedZoneRead(ctx, d, meta)
}
func resourceIBMDNSLinkedZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 2 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/linkedDnsZoneID", d.Id()))
	}
	instanceID := idSet[0]
	linkedDnsZoneID := idSet[1]
	deleteLinkedZoneOptions := sess.NewDeleteLinkedZoneOptions(instanceID, linkedDnsZoneID)

	mk := "linked_dns_zone_" + instanceID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)
	response, err := sess.DeleteLinkedZone(deleteLinkedZoneOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error reading DNS Services secondary zone:%s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
