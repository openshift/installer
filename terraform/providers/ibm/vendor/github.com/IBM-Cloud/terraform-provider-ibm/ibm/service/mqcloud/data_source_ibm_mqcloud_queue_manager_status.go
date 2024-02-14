// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package mqcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudQueueManagerStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudQueueManagerStatusRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQ on Cloud service instance.",
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
		return diag.FromErr(err)
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Read Queue Manager Status failed %s", err))
	}

	getQueueManagerStatusOptions := &mqcloudv1.GetQueueManagerStatusOptions{}

	getQueueManagerStatusOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	getQueueManagerStatusOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))

	queueManagerStatus, response, err := mqcloudClient.GetQueueManagerStatusWithContext(context, getQueueManagerStatusOptions)
	if err != nil {
		log.Printf("[DEBUG] GetQueueManagerStatusWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetQueueManagerStatusWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmMqcloudQueueManagerStatusID(d))

	if err = d.Set("status", queueManagerStatus.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	return nil
}

// dataSourceIbmMqcloudQueueManagerStatusID returns a reasonable ID for the list.
func dataSourceIbmMqcloudQueueManagerStatusID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
