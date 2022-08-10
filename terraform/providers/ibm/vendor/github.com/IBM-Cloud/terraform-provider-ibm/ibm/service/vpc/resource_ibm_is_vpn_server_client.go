// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsVPNServerClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVPNServerClientDisconnect,
		ReadContext:   resourceIBMIsVPNServerClientDisconnect,
		UpdateContext: resourceIBMIsVPNServerClientDisconnect,
		DeleteContext: resourceIBMIsVPNServerClientDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"vpn_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPN server identifier.",
			},
			"vpn_client": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPN Client identifier.",
			},
			"delete": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The delete to use for this VPN client to be deleted or not, when false, client is disconneted and when set to true client is deleted.",
			},
			"status_code": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "status code of the result.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "description of the result.",
			},
		},
	}
}

func resourceIBMIsVPNServerClientDisconnect(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVPNServerClientOptions := &vpcv1.GetVPNServerClientOptions{}

	getVPNServerClientOptions.SetVPNServerID(d.Get("vpn_server").(string))
	getVPNServerClientOptions.SetID(d.Get("vpn_client").(string))

	_, response, err := sess.GetVPNServerClientWithContext(context, getVPNServerClientOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVPNServerClientWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerClientWithContext failed %s\n%s", err, response))
	}

	var flag bool
	if v, ok := d.GetOk("delete"); ok {
		flag = v.(bool)
	}

	if flag == false {

		disconnectVPNServerRouteOptions := &vpcv1.DisconnectVPNClientOptions{}
		disconnectVPNServerRouteOptions.SetVPNServerID(d.Get("vpn_server").(string))
		disconnectVPNServerRouteOptions.SetID(d.Get("vpn_client").(string))

		response, err := sess.DisconnectVPNClientWithContext(context, disconnectVPNServerRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] DisconnectVPNClientWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] DisconnectVPNClientWithContext failed %s\n%s", err, response))
		}

		if err = d.Set("status_code", response.StatusCode); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting status_code: %s", err))
		}

		if err = d.Set("description", "The VPN client disconnection request was accepted."); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}

		d.SetId(fmt.Sprintf("%s/%s/%v", d.Get("vpn_server").(string), d.Get("vpn_client").(string), response.StatusCode))

	} else if flag == true {

		deleteVPNServerClientOptions := &vpcv1.DeleteVPNServerClientOptions{}
		deleteVPNServerClientOptions.SetVPNServerID(d.Get("vpn_server").(string))
		deleteVPNServerClientOptions.SetID(d.Get("vpn_client").(string))

		response, err := sess.DeleteVPNServerClientWithContext(context, deleteVPNServerClientOptions)
		if err != nil {
			log.Printf("[DEBUG] DeleteVPNServerClientWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] DeleteVPNServerClientWithContext failed %s\n%s", err, response))
		}

		if err = d.Set("status_code", response.StatusCode); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting status_code: %s", err))
		}

		if err = d.Set("description", "The VPN client disconnection request was accepted."); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting status_code: %s", err))
		}

		d.SetId(fmt.Sprintf("%s/%s", d.Get("vpn_server").(string), d.Get("vpn_client").(string)))
	}

	if err = d.Set("delete", d.Get("delete")); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting delete: %s", err))
	}
	return nil
}

func resourceIBMIsVPNServerClientDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR]  Failed %s\n%s", "false", err))
	}
	if len(parts) != 2 {
		return diag.FromErr(fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of vpnServer/vpnClient", d.Id()))
	}
	vpnServer := parts[0]
	vpnClient := parts[1]

	getVPNServerClientOptions := &vpcv1.GetVPNServerClientOptions{}

	getVPNServerClientOptions.SetVPNServerID(vpnServer)
	getVPNServerClientOptions.SetID(vpnClient)

	_, response, err := sess.GetVPNServerClientWithContext(context, getVPNServerClientOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVPNServerClientWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerClientWithContext failed %s\n%s", err, response))
	}

	deleteVPNServerClientOptions := &vpcv1.DeleteVPNServerClientOptions{}
	deleteVPNServerClientOptions.SetVPNServerID(d.Get("vpn_server").(string))
	deleteVPNServerClientOptions.SetID(d.Get("vpn_client").(string))

	response, err = sess.DeleteVPNServerClientWithContext(context, deleteVPNServerClientOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteVPNServerClientWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] DeleteVPNServerClientWithContext failed %s\n%s", err, response))
	}

	d.SetId("")
	return nil
}
