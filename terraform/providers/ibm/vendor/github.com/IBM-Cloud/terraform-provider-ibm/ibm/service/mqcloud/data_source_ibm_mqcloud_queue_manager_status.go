// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package mqcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudQueueManagerStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudQueueManagerStatusRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"queue_manager_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the queue manager to retrieve its full details.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deploying and failed states are not queue manager states, they are states which can occur when the request to deploy has been fired, or with that request has failed without producing a queue manager to have any state. The other states map to the queue manager states. State \"ending\" is either quiesing or ending immediately. State \"ended\" is either ended normally or endedimmediately. The others map one to one with queue manager states.",
			},
		},
	}
}

func dataSourceIbmMqcloudQueueManagerStatusRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_queue_manager_status", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read Queue Manager Status failed: %s", err.Error()), "(Data) ibm_mqcloud_queue_manager_status", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getQueueManagerStatusOptions := &mqcloudv1.GetQueueManagerStatusOptions{}

	getQueueManagerStatusOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	getQueueManagerStatusOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))

	queueManagerStatus, _, err := mqcloudClient.GetQueueManagerStatusWithContext(context, getQueueManagerStatusOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetQueueManagerStatusWithContext failed: %s", err.Error()), "(Data) ibm_mqcloud_queue_manager_status", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmMqcloudQueueManagerStatusID(d))

	if err = d.Set("status", queueManagerStatus.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_mqcloud_queue_manager_status", "read", "set-status").GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudQueueManagerStatusID returns a reasonable ID for the list.
func dataSourceIbmMqcloudQueueManagerStatusID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
