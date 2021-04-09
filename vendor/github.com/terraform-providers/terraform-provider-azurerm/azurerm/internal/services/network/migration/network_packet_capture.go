package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func NetworkPacketCaptureV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_watcher_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"maximum_bytes_per_packet": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},

			"maximum_bytes_per_session": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1073741824,
			},

			"maximum_capture_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  18000,
			},

			"storage_location": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_path": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"storage_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"storage_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"local_port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"remote_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"remote_port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func NetworkPacketCaptureV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "/NetworkPacketCaptures/", "/packetCaptures/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
