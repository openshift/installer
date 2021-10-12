package alicloud

import (
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMNSTopicSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMNSTopicSubscriptionRead,
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"filter_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"notify_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"notify_content_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMNSTopicSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	topicName := d.Get("topic_name").(string)
	var namePrefix string
	if v, ok := d.GetOk("name_prefix"); ok {
		namePrefix = v.(string)
	}

	var subscriptionAttr []ali_mns.SubscriptionAttribute
	for {
		var nextMaker string
		raw, err := client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
			return subscriptionManager.ListSubscriptionDetailByTopic(nextMaker, 1000, namePrefix)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_topic_subscriptions", "ListSubscriptionDetailByTopic", AliMnsERROR)
		}

		addDebug("ListSubscriptionDetailByTopic", raw)
		subscriptionDetails, _ := raw.(ali_mns.SubscriptionDetails)
		for _, attr := range subscriptionDetails.Attrs {
			subscriptionAttr = append(subscriptionAttr, attr)
		}
		nextMaker = subscriptionDetails.NextMarker
		if nextMaker == "" {
			break
		}
	}
	return mnsTopicSubcriptionDescription(d, subscriptionAttr)
}

func mnsTopicSubcriptionDescription(d *schema.ResourceData, topicAttr []ali_mns.SubscriptionAttribute) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range topicAttr {
		mapping := map[string]interface{}{
			"id":                    item.SubscriptionName,
			"name":                  item.SubscriptionName,
			"topic_name":            item.TopicName,
			"endpoint":              item.Endpoint,
			"filter_tag":            item.FilterTag,
			"notify_strategy":       item.NotifyStrategy,
			"notify_content_format": item.NotifyContentFormat,
		}

		ids = append(ids, item.SubscriptionName)
		names = append(names, item.SubscriptionName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("subscriptions", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
