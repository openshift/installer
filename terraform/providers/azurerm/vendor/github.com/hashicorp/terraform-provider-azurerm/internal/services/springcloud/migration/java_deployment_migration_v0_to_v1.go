package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudJavaDeploymentV0ToV1 struct{}

func (s SpringCloudJavaDeploymentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"spring_cloud_app_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"environment_variables": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"instance_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  1,
		},

		"jvm_options": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"quota": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"memory": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"runtime_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s SpringCloudJavaDeploymentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudDeploymentIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
