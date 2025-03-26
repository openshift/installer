// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIPublicNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIPublicNetworkRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_Name: {
				Computed:    true,
				Description: "The name of the network.",
				Type:        schema.TypeString,
			},
			Attr_Type: {
				Computed:    true,
				Description: "The type of VLAN that the network is connected to.",
				Type:        schema.TypeString,
			},
			Attr_VLanID: {
				Computed:    true,
				Description: "The ID of the VLAN that the network is connected to.",
				Type:        schema.TypeInt,
			},
		},
	}
}

func dataSourceIBMPIPublicNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.GetAllPublic()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(networkdata.Networks) < 1 {
		return diag.Errorf("error getting public network or no public network found in %s", cloudInstanceID)
	}

	d.SetId(*networkdata.Networks[0].NetworkID)
	if networkdata.Networks[0].Crn != "" {
		d.Set(Attr_CRN, networkdata.Networks[0].Crn)
	}
	if networkdata.Networks[0].Name != nil {
		d.Set(Attr_Name, networkdata.Networks[0].Name)
	}
	if networkdata.Networks[0].Type != nil {
		d.Set(Attr_Type, networkdata.Networks[0].Type)
	}
	if networkdata.Networks[0].VlanID != nil {
		d.Set(Attr_VLanID, networkdata.Networks[0].VlanID)
	}

	return nil
}
