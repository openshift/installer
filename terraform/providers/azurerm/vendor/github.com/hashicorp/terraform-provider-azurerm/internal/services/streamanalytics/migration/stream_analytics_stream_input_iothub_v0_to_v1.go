package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsStreamInputIotHubV0ToV1 struct{}

func (s StreamAnalyticsStreamInputIotHubV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"stream_analytics_job_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"iothub_namespace": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"eventhub_consumer_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"shared_access_policy_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			ForceNew:  true,
			Sensitive: true,
		},

		"shared_access_policy_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"serialization": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"field_delimiter": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"encoding": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func (s StreamAnalyticsStreamInputIotHubV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := inputs.ParseInputIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
