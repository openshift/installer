// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	defaultReplicationFactor = 3
	defaultCleanupPolicy     = "delete"
	defaultRetentionBytes    = 1073741824 // 100 MB
	defaultRetentionMs       = 86400000   // 24 hours
	defaultSegmentBytes      = 536870912  // 512 MB
)

var (
	brokerVersion       = sarama.V2_3_0_0
	allowedTopicConfigs = []string{
		"cleanup.policy",
		"retention.ms",
		"retention.bytes",
		"segment.ms",
		"segment.bytes",
		"segment.index.bytes",
	}
	defaultConfigs = map[string]interface{}{
		"cleanup.policy":  defaultCleanupPolicy,
		"retention.ms":    defaultRetentionMs,
		"retention.bytes": defaultRetentionBytes,
		"segment.bytes":   defaultSegmentBytes,
	}
)

func resourceIBMEventStreamsTopic() *schema.Resource {
	return &schema.Resource{
		Exists:   resourceIBMEventStreamsTopicExists,
		Create:   resourceIBMEventStreamsTopicCreate,
		Read:     resourceIBMEventStreamsTopicRead,
		Update:   resourceIBMEventStreamsTopicUpdate,
		Delete:   resourceIBMEventStreamsTopicDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"resource_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The CRN of the Event Streams instance",
				Required:    true,
			},
			"kafka_http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API endpoint for interacting with Event Streams REST API",
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
				Description: "The number of partitions",
				Optional:    true,
				Default:     1,
			},
			"config": {
				Type:        schema.TypeMap,
				Description: "The configuration parameters of a topic",
				Optional:    true,
			},
		},
	}
}

// clientPool maintains Kafka admin client for each instance.
// key is instance's CRN
var clientPool = map[string]sarama.ClusterAdmin{}

func resourceIBMEventStreamsTopicExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	adminClient, _, err := createSaramaAdminClient(d, meta)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicExists createSaramaAdminClient err %s", err)
		return false, err
	}
	topicName := d.Get("name").(string)
	topics, err := adminClient.DescribeTopics([]string{topicName})
	if err != nil || len(topics) != 1 {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicExists DescribeTopics err %s", err)
		return false, err
	}
	log.Printf("[INFO] resourceIBMEventStreamsTopicExists topic %s exists", topicName)
	return true, nil
}

func resourceIBMEventStreamsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	adminClient, instanceCRN, err := createSaramaAdminClient(d, meta)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicCreate createSaramaAdminClient err %s", err)
		return err
	}
	topicName := d.Get("name").(string)
	partitions := d.Get("partitions").(int)
	config := d.Get("config").(map[string]interface{})
	topicDetail := sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(defaultReplicationFactor),
		ConfigEntries:     config2TopicDetail(config),
	}
	err = adminClient.CreateTopic(topicName, &topicDetail, false)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicCreate CreateTopic err %s", err)
	}
	log.Printf("[INFO] resourceIBMEventStreamsTopicCreate CreateTopic: topic is %s, detail is %v", topicName, topicDetail)
	d.SetId(getTopicID(instanceCRN, topicName))
	return resourceIBMEventStreamsTopicRead(d, meta)
}

func resourceIBMEventStreamsTopicRead(d *schema.ResourceData, meta interface{}) error {
	adminClient, instanceCRN, err := createSaramaAdminClient(d, meta)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicRead createSaramaAdminClient err %s", err)
		return err
	}
	topicID := d.Id()
	topicName := getTopicName(topicID)
	topics, err := adminClient.ListTopics()
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicRead ListTopics err %s", err)
		return err
	}
	for name, detail := range topics {
		if name == topicName {
			d.Set("resource_instance_id", instanceCRN)
			d.Set("name", name)
			d.Set("partitions", detail.NumPartitions)
			if config := d.Get("config"); config != nil {
				savedConfig := map[string]*string{}
				for k := range config.(map[string]interface{}) {
					if value, ok := detail.ConfigEntries[k]; ok {
						savedConfig[k] = value
					}
				}
				d.Set("config", topicDetail2Config(savedConfig))
			}
			return nil
		}
	}
	log.Printf("[INFO] resourceIBMEventStreamsTopicRead topic %s does not exist", topicName)
	d.SetId("")
	return nil
}

func resourceIBMEventStreamsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	adminClient, _, err := createSaramaAdminClient(d, meta)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicUpdate createSaramaAdminClient err %s", err)
		return err
	}
	topicName := d.Get("name").(string)
	if d.HasChange("partitions") {
		oi, ni := d.GetChange("partitions")
		oldPartitions := oi.(int)
		newPartitions := ni.(int)
		log.Printf("[INFO]resourceIBMEventStreamsTopicUpdate Updating partitions from %d to %d", oldPartitions, newPartitions)
		err = adminClient.CreatePartitions(topicName, int32(newPartitions), nil, false)
		if err != nil {
			log.Printf("[DEBUG]resourceIBMEventStreamsTopicUpdate CreatePartitions err %s", err)
			return err
		}
		d.Set("partitions", int32(newPartitions))
		log.Printf("[INFO]resourceIBMEventStreamsTopicUpdate partitions is set to %d", newPartitions)
	}
	if d.HasChange("config") {
		config := d.Get("config").(map[string]interface{})
		configEntries := config2TopicDetail(config)
		err = adminClient.AlterConfig(sarama.TopicResource, topicName, configEntries, false)
		if err != nil {
			log.Printf("[DEBUG]resourceIBMEventStreamsTopicUpdate AlterConfig err %s", err)
			return err
		}
		d.Set("config", topicDetail2Config(configEntries))
		log.Printf("[INFO]resourceIBMEventStreamsTopicUpdate config is set to %v", topicDetail2Config(configEntries))
	}
	return resourceIBMEventStreamsTopicRead(d, meta)
}

func resourceIBMEventStreamsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	adminClient, _, err := createSaramaAdminClient(d, meta)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicDelete createSaramaAdminClient err %s", err)
		return err
	}
	topicName := d.Get("name").(string)
	err = adminClient.DeleteTopic(topicName)
	if err != nil {
		log.Printf("[DEBUG] resourceIBMEventStreamsTopicDelete DeleteTopic err %s", err)
		return err
	}
	d.SetId("")
	log.Printf("[INFO]resourceIBMEventStreamsTopicDelete topic %s deleted", topicName)
	return nil
}

func createSaramaAdminClient(d *schema.ResourceData, meta interface{}) (sarama.ClusterAdmin, string, error) {
	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		log.Printf("[DEBUG] createSaramaAdminClient BluemixSession err %s", err)
		return nil, "", err
	}
	apiKey := bxSession.Config.BluemixAPIKey
	if len(apiKey) == 0 {
		log.Printf("[DEBUG] createSaramaAdminClient BluemixAPIKey is empty")
		return nil, "", fmt.Errorf("failed to get IBM cloud API key")
	}
	rsConClient, err := meta.(ClientSession).ResourceControllerAPI()
	if err != nil {
		log.Printf("[DEBUG] createSaramaAdminClient ResourceControllerAPI err %s", err)
		return nil, "", err
	}
	rcAPI := rsConClient.ResourceServiceInstance()
	instanceCRN := d.Get("resource_instance_id").(string)
	if len(instanceCRN) == 0 {
		topicID := d.Id()
		if len(topicID) == 0 || !strings.Contains(topicID, ":") {
			log.Printf("[DEBUG] createSaramaAdminClient resource_instance_id is missing")
			return nil, "", fmt.Errorf("resource_instance_id is required")
		}
		instanceCRN = getInstanceCRN(topicID)
	}
	instance, err := rcAPI.GetInstance(instanceCRN)
	if err != nil {
		log.Printf("[DEBUG] createSaramaAdminClient GetInstance err %s", err)
		return nil, "", err
	}
	if instance.Extensions == nil {
		log.Printf("[DEBUG] createSaramaAdminClient instance %s extension is nil", instance.ID)
		return nil, "", fmt.Errorf("instance %s extension is nil", instance.ID)
	}
	adminURL := instance.Extensions["kafka_http_url"].(string)
	d.Set("kafka_http_url", adminURL)
	log.Printf("[INFO] createSaramaAdminClient kafka_http_url is set to %s", adminURL)
	brokerAddress := expandStringList(instance.Extensions["kafka_brokers_sasl"].([]interface{}))
	d.Set("kafka_brokers_sasl", brokerAddress)
	log.Printf("[INFO] createSaramaAdminClient kafka_brokers_sasl is set to %s", brokerAddress)
	tenantID := strings.TrimPrefix(strings.Split(adminURL, ".")[0], "https://")

	config := sarama.NewConfig()
	config.ClientID, _ = os.Hostname()
	config.Net.SASL.Enable = true
	if tenantID != "" && tenantID != "admin" {
		config.Net.SASL.AuthIdentity = tenantID
	}
	config.Net.SASL.User = "token"
	config.Net.SASL.Password = apiKey
	config.Net.TLS.Enable = true
	config.Version = brokerVersion
	adminClient, err := sarama.NewClusterAdmin(brokerAddress, config)
	if err != nil {
		log.Printf("[DEBUG] createSaramaAdminClient NewClusterAdmin err %s", err)
		return nil, "", err
	}
	clientPool[instanceCRN] = adminClient
	log.Printf("[INFO] createSaramaAdminClient instance %s 's client is initialized", instanceCRN)
	return adminClient, instanceCRN, nil
}

func topicDetail2Config(topicConfigEntries map[string]*string) map[string]*string {
	configs := map[string]*string{}
	for key, value := range topicConfigEntries {
		if indexOf(key, allowedTopicConfigs) != -1 {
			configs[key] = value
		}
	}
	return configs
}

func config2TopicDetail(config map[string]interface{}) map[string]*string {
	configEntries := make(map[string]*string)
	for key, value := range config {
		switch value := value.(type) {
		case string:
			configEntries[key] = &value
		}
	}
	return configEntries
}

func getTopicID(instanceCRN string, topicName string) string {
	crnSegments := strings.Split(instanceCRN, ":")
	crnSegments[8] = "topic"
	crnSegments[9] = topicName
	return strings.Join(crnSegments, ":")
}

func getTopicName(topicID string) string {
	return strings.Split(topicID, ":")[9]
}

func getInstanceCRN(topicID string) string {
	crnSegments := strings.Split(topicID, ":")
	crnSegments[8] = ""
	crnSegments[9] = ""
	return strings.Join(crnSegments, ":")
}
