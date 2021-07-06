// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsPermittedNetworks = "dns_permitted_networks"
)

func dataSourceIBMPrivateDNSPermittedNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPrivateDNSPermittedNetworksRead,

		Schema: map[string]*schema.Schema{

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID",
			},

			pdnsZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID",
			},

			pdnsPermittedNetworks: {

				Type:        schema.TypeList,
				Description: "Collection of permitted networks",
				Computed:    true,
				Elem: &schema.Resource{

					Schema: map[string]*schema.Schema{

						pdnsPermittedNetworkID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network Id",
						},

						pdnsInstanceID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Id",
						},

						pdnsZoneID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone Id",
						},

						pdnsNetworkType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network Type",
						},

						pdnsPermittedNetwork: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "permitted network",
						},

						pdnsPermittedNetworkCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network creation date",
						},

						pdnsPermittedNetworkModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network Modification date",
						},

						pdnsPermittedNetworkState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network status",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPrivateDNSPermittedNetworksRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	instanceID := d.Get(pdnsInstanceID).(string)
	dnsZoneID := d.Get(pdnsZoneID).(string)
	listPermittedNetworkOptions := sess.NewListPermittedNetworksOptions(instanceID, dnsZoneID)
	availablePermittedNetworks, detail, err := sess.ListPermittedNetworks(listPermittedNetworkOptions)
	if err != nil {
		return fmt.Errorf("Error reading list of pdns permitted networks:%s\n%s", err, detail)
	}

	permittedNetworks := make([]map[string]interface{}, 0)

	for _, instance := range availablePermittedNetworks.PermittedNetworks {
		permittedNetwork := map[string]interface{}{}
		permittedNetworkVpcCrn := map[string]interface{}{}
		permittedNetwork[pdnsInstanceID] = instanceID
		permittedNetwork[pdnsPermittedNetworkID] = instance.ID
		permittedNetwork[pdnsPermittedNetworkCreatedOn] = instance.CreatedOn
		permittedNetwork[pdnsPermittedNetworkModifiedOn] = instance.ModifiedOn
		permittedNetwork[pdnsPermittedNetworkState] = instance.State
		permittedNetwork[pdnsNetworkType] = instance.Type
		permittedNetworkVpcCrn[pdnsVpcCRN] = instance.PermittedNetwork.VpcCrn
		permittedNetwork[pdnsPermittedNetwork] = permittedNetworkVpcCrn
		permittedNetwork[pdnsZoneID] = dnsZoneID

		permittedNetworks = append(permittedNetworks, permittedNetwork)
	}
	d.SetId(dataSourceIBMPrivateDNSPermittedNetworkID(d))
	d.Set(pdnsPermittedNetworks, permittedNetworks)
	return nil
}

// dataSourceIBMPrivateDnsPermittedNetworkID returns a reasonable ID for dns permitted network list.
func dataSourceIBMPrivateDNSPermittedNetworkID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
