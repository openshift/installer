// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"strings"

	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const INSTANCE_ID = "instance_id"

// AddSchemaData will add the Schemas 'instance_id' and 'region' to the resource
func AddSchemaData(resource *schema.Resource) *schema.Resource {
	resource.Schema["instance_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the Security and Compliance Center instance.",
	}
	return resource
}

// getRegionData will check if the field region is defined
func getRegionData(client securityandcompliancecenterapiv3.SecurityAndComplianceCenterApiV3, d *schema.ResourceData) string {
	val, ok := d.GetOk("region")
	if ok {
		return val.(string)
	} else {
		url := client.Service.GetServiceURL()
		return strings.Split(url, ".")[1]
	}
}

// setRegionData will set the field "region" field if the field was previously defined
func setRegionData(d *schema.ResourceData, region string) error {
	if val, ok := d.GetOk("region"); ok {
		return d.Set("region", val.(string))
	}
	return nil
}
