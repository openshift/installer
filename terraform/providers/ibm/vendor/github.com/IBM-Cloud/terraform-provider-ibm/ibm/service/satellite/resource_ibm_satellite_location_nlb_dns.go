// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMSatelliteLocationNlbDns() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSatelliteLocationNlbDnsCreate,
		ReadContext:   resourceIbmSatelliteLocationNlbDnsRead,
		DeleteContext: resourceIbmSatelliteLocationNlbDnsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ips": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIbmSatelliteLocationNlbDnsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	bmxSess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}

	registerMultishiftClusterOptions := &kubernetesserviceapiv1.RegisterMultishiftClusterOptions{}

	registerMultishiftClusterOptions.SetXAuthRefreshToken(bmxSess.Config.IAMRefreshToken)
	if controller, ok := d.GetOk("location"); ok {
		registerMultishiftClusterOptions.SetController(controller.(string))
	}
	if _, ok := d.GetOk("ips"); ok {
		ips := []string{}
		for _, segmentsItem := range d.Get("ips").(*schema.Set).List() {
			ips = append(ips, segmentsItem.(string))
		}
		registerMultishiftClusterOptions.SetIps(ips)
	}

	mscRegisterResp, response, err := satClient.RegisterMultishiftClusterWithContext(context, registerMultishiftClusterOptions)
	if err != nil {
		log.Printf("[DEBUG] RegisterMultishiftClusterWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("RegisterMultishiftClusterWithContext failed %s\n%s", err, response))
	}
	d.SetId(*mscRegisterResp.Controller)

	return resourceIbmSatelliteLocationNlbDnsRead(context, d, meta)
}

func resourceIbmSatelliteLocationNlbDnsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ID := d.Id()
	nlbClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	nlbAPI := nlbClient.NlbDns()
	getSatLocationNlbDNSListOptions := &kubernetesserviceapiv1.GetSatLocationNlbDNSListOptions{}
	getSatLocationNlbDNSListOptions.Controller = flex.PtrToString(ID)

	_, err = nlbAPI.GetLocationNLBDNSList(ID)
	if err != nil {
		log.Printf("[DEBUG] GetSatLocationNlbDNSListWithContext failed %s\n", err)
		return diag.FromErr(fmt.Errorf("[ERROR] GetSatLocationNlbDNSListWithContext failed %s", err))
	}

	return nil
}

func resourceIbmSatelliteLocationNlbDnsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
