// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package db2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-db2-go-sdk/db2saasv1"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func DataSourceIbmDb2Whitelist() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceIbmDb2WhitelistRead,
		DeprecationMessage: "This service is deprecated and replaced with `ibm_db2_allowlist_ip`",

		Schema: map[string]*schema.Schema{
			"x_deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "CRN deployment id.",
			},
			"ip_addresses": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of IP addresses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address, in IPv4/ipv6 format.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the IP address.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmDb2WhitelistRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db2saasClient, err := meta.(conns.ClientSession).Db2saasV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_whitelist_ip", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDb2SaasWhitelistOptions := &db2saasv1.GetDb2SaasAllowlistOptions{}

	getDb2SaasWhitelistOptions.SetXDeploymentID(d.Get("x_deployment_id").(string))

	successGetWhitelistIPs, _, err := db2saasClient.GetDb2SaasAllowlistWithContext(context, getDb2SaasWhitelistOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDb2SaasAllowlistWithContext failed: %s", err.Error()), "(Data) ibm_db2_whitelist_ip", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmDb2WhitelistID(d))

	ipAddresses := []map[string]interface{}{}
	for _, ipAddressesItem := range successGetWhitelistIPs.IpAddresses {
		ipAddressesItemMap, err := DataSourceIbmDb2WhitelistIpAddressToMap(&ipAddressesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_whitelist_ip", "read", "ip_addresses-to-map").GetDiag()
		}
		ipAddresses = append(ipAddresses, ipAddressesItemMap)
	}
	if err = d.Set("ip_addresses", ipAddresses); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ip_addresses: %s", err), "(Data) ibm_db2_whitelist_ip", "read", "set-ip_addresses").GetDiag()
	}

	return nil
}

// dataSourceIbmDb2SaasWhitelistID returns a reasonable ID for the list.
func dataSourceIbmDb2WhitelistID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmDb2WhitelistIpAddressToMap(model *db2saasv1.IpAddress) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	modelMap["description"] = *model.Description
	return modelMap, nil
}
