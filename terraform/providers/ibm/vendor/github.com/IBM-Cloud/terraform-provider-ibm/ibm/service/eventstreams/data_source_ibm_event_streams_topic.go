// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventstreams

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMEventStreamsTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEventStreamsTopicRead,
		Schema: map[string]*schema.Schema{
			"resource_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CRN of the Event Streams instance",
			},
			"kafka_http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API endpoint for interacting with Event Streams REST API",
			},
			"kafka_brokers_sasl": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Kafka brokers addresses for interacting with Kafka native API",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the topic",
				Required:    true,
			},
			"partitions": {
				Type:        schema.TypeInt,
				Description: "The number of partitions of the topic",
				Computed:    true,
			},
			"config": {
				Type:        schema.TypeMap,
				Description: "The configuration parameters of the topic.",
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMEventStreamsTopicRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminClient, instanceCRN, err := createSaramaAdminClient(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMEventStreamsTopicRead createSaramaAdminClient: %s", err), "ibm_event_streams_topic", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	topics, err := adminClient.ListTopics()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMEventStreamsTopicRead ListTopics: %s", err), "ibm_event_streams_topic", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	topicName := d.Get("name").(string)
	for name := range topics {
		if name == topicName {
			topicID := getTopicID(instanceCRN, topicName)
			d.SetId(topicID)
			log.Printf("[INFO]dataSourceIBMEventStreamsTopicRead set topic ID to %s", topicID)
			d.Set("resource_instance_id", instanceCRN)
			return nil
		}
	}
	tfErr := flex.TerraformErrorf(fmt.Errorf("topic %s does not exist", topicName),
		fmt.Sprintf("dataSourceIBMEventStreamsTopicRead topic %s does not exist", topicName), "ibm_event_streams_topic", "read")
	log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
	return tfErr.GetDiag()
}
