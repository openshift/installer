// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventstreams

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/eventstreams-go-sdk/pkg/schemaregistryv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMEventStreamsSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEventStreamsSchemaRead,

		Schema: map[string]*schema.Schema{
			"resource_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID or CRN of the Event Streams service instance",
			},
			"kafka_http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API endpoint for interacting with an Event Streams REST API",
			},
			"schema_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique ID to be assigned to the schema.",
			},
		},
	}
}

func dataSourceIBMEventStreamsSchemaRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schemaregistryClient, err := meta.(conns.ClientSession).ESschemaRegistrySession()
	if err != nil {
		log.Printf("[DEBUG] dataSourceIBMEventStreamsSchemaRead schemaregistryClient err %s", err)
		return diag.FromErr(err)
	}

	adminURL, instanceCRN, err := getInstanceURL(d, meta)
	if err != nil {
		log.Printf("[DEBUG] dataSourceIBMEventStreamsSchemaRead getInstanceURL err %s", err)
		return diag.FromErr(err)
	}
	schemaregistryClient.SetServiceURL(adminURL)

	getLatestSchemaOptions := &schemaregistryv1.GetLatestSchemaOptions{}

	schemaID := d.Get("schema_id").(string)
	getLatestSchemaOptions.SetID(schemaID)

	schema, response, err := schemaregistryClient.GetLatestSchemaWithContext(context, getLatestSchemaOptions)
	if err != nil || schema == nil {
		log.Printf("[DEBUG] GetLatestSchemaWithContext failed with error: %s and response:\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLatestSchemaWithContext failed with error: %s\n and response:%s", err, response))
	}
	uniqueID := getUniqueSchemaID(instanceCRN, schemaID)

	d.SetId(uniqueID)
	d.Set("resource_instance_id", instanceCRN)
	return nil
}
