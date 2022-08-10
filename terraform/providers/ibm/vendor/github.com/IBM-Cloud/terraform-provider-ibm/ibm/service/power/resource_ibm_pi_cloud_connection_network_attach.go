// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

const (
	PICloudConnectionNetworkId = "pi_network_id"
)

func ResourceIBMPICloudConnectionNetworkAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPICloudConnectionNetworkAttachCreate,
		ReadContext:   resourceIBMPICloudConnectionNetworkAttachRead,
		DeleteContext: resourceIBMPICloudConnectionNetworkAttachDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Required Attributes
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PI cloud instance ID",
			},
			helpers.PICloudConnectionId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud Connection ID",
			},
			PICloudConnectionNetworkId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network ID to attach to this cloud connection",
			},
		},
	}
}

func resourceIBMPICloudConnectionNetworkAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	cloudConnectionID := d.Get(helpers.PICloudConnectionId).(string)
	networkID := d.Get(PICloudConnectionNetworkId).(string)

	client := st.NewIBMPICloudConnectionClient(ctx, sess, cloudInstanceID)
	jobClient := st.NewIBMPIJobClient(ctx, sess, cloudInstanceID)

	_, jobReference, err := client.AddNetwork(cloudConnectionID, networkID)
	if err != nil {
		log.Printf("[ERROR] attach network to cloud connection failed %v", err)
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, cloudConnectionID, networkID))
	if jobReference != nil {
		_, err = waitForIBMPIJobCompleted(ctx, jobClient, *jobReference.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPICloudConnectionNetworkAttachRead(ctx, d, meta)
}

func resourceIBMPICloudConnectionNetworkAttachRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := parts[0]
	cloudConnectionID := parts[1]
	networkID := parts[2]

	d.Set(helpers.PICloudInstanceId, cloudInstanceID)
	d.Set(helpers.PICloudConnectionId, cloudConnectionID)
	d.Set(PICloudConnectionNetworkId, networkID)

	return nil
}

func resourceIBMPICloudConnectionNetworkAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := parts[0]
	cloudConnectionID := parts[1]
	networkID := parts[2]

	client := st.NewIBMPICloudConnectionClient(ctx, sess, cloudInstanceID)
	jobClient := st.NewIBMPIJobClient(ctx, sess, cloudInstanceID)

	_, jobReference, err := client.DeleteNetwork(cloudConnectionID, networkID)
	if err != nil {
		log.Printf("[DEBUG] detach network from cloud connection failed %v", err)
		return diag.FromErr(err)
	}
	if jobReference != nil {
		_, err = waitForIBMPIJobCompleted(ctx, jobClient, *jobReference.ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}
