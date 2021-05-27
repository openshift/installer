// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMEventStreamsTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMEventStreamsTopicRead,
		Schema: map[string]*schema.Schema{
			"resource_instance_id": &schema.Schema{
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the topic",
				Required:    true,
			},
			"partitions": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The number of partitions of the topic",
				Computed:    true,
			},
			"config": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "The configuration parameters of the topic.",
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMEventStreamsTopicRead(d *schema.ResourceData, meta interface{}) error {
	adminClient, instanceCRN, err := createSaramaAdminClient(d, meta)
	if err != nil {
		log.Printf("[DEBUG]dataSourceIBMEventStreamsTopicRead createSaramaAdminClient err %s", err)
		return err
	}
	topics, err := adminClient.ListTopics()
	if err != nil {
		log.Printf("[DEBUG]dataSourceIBMEventStreamsTopicRead ListTopics err %s", err)
		return err
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
	log.Printf("[DEBUG]dataSourceIBMEventStreamsTopicRead topic %s does not exist", topicName)
	return fmt.Errorf("topic %s does not exist", topicName)
}
