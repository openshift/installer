// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPINetworkPeers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkPeersRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_NetworkPeers: {
				Computed:    true,
				Description: "List of network peers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Description: {
							Computed:    true,
							Description: "Description of the network peer.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "ID of the network peer.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "Name of the network peer.",
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Computed:    true,
							Description: "Type of the network peer.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPINetworkPeersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	networkC := instance.NewIBMPINetworkPeerClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.GetNetworkPeers()
	if err != nil {
		return diag.FromErr(err)
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	networkPeers := []map[string]interface{}{}
	if networkdata.NetworkPeers != nil {
		for _, np := range networkdata.NetworkPeers {
			npMap := dataSourceIBMPINetworkPeersNetworkPeerToMap(np)

			networkPeers = append(networkPeers, npMap)
		}
	}
	d.Set(Attr_NetworkPeers, networkPeers)

	return nil
}

func dataSourceIBMPINetworkPeersNetworkPeerToMap(np *models.NetworkPeer) map[string]interface{} {
	npMap := make(map[string]interface{})
	if np.Description != nil {
		npMap[Attr_Description] = np.Description
	}
	if np.ID != nil {
		npMap[Attr_ID] = np.ID
	}
	if np.Name != nil {
		npMap[Attr_Name] = np.Name
	}
	if np.Type != nil {
		npMap[Attr_Type] = np.Type
	}
	return npMap
}
