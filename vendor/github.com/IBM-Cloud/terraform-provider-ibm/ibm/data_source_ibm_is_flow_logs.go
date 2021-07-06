// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isFlowLogs = "flow_log_collectors"
)

func dataSourceIBMISFlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISFlowLogsRead,

		Schema: map[string]*schema.Schema{

			isFlowLogs: {
				Type:        schema.TypeList,
				Description: "Collection of flow log collectors",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this flow log collector",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this flow log collector",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this flow log collector",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Log Collector name",
						},
						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group of flow log",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time flow log was created",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the flow log collector",
						},
						"storage_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Object Storage bucket name where the collected flows will be logged",
						},
						"active": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this collector is active",
						},
						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target id that the flow log collector is to collect flow logs",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC this flow log collector is associated with",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISFlowLogsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	start := ""
	allrecs := []vpcv1.FlowLogCollector{}
	for {
		listOptions := &vpcv1.ListFlowLogCollectorsOptions{}
		if start != "" {
			listOptions.Start = &start
		}
		flowlogCollectors, response, err := sess.ListFlowLogCollectors(listOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Flow Logs for VPC %s\n%s", err, response)
		}
		start = GetNext(flowlogCollectors.Next)
		allrecs = append(allrecs, flowlogCollectors.FlowLogCollectors...)
		if start == "" {
			break
		}
	}
	flowlogsInfo := make([]map[string]interface{}, 0)
	for _, flowlogCollector := range allrecs {

		targetIntf := flowlogCollector.Target
		target := targetIntf.(*vpcv1.FlowLogCollectorTarget)

		l := map[string]interface{}{
			"id":              *flowlogCollector.ID,
			"crn":             *flowlogCollector.CRN,
			"href":            *flowlogCollector.Href,
			"name":            *flowlogCollector.Name,
			"resource_group":  *flowlogCollector.ResourceGroup.ID,
			"created_at":      flowlogCollector.CreatedAt.String(),
			"lifecycle_state": *flowlogCollector.LifecycleState,
			"storage_bucket":  *flowlogCollector.StorageBucket.Name,
			"active":          *flowlogCollector.Active,
			"vpc":             *flowlogCollector.VPC.ID,
			"target":          *target.ID,
		}
		flowlogsInfo = append(flowlogsInfo, l)
	}
	d.SetId(dataSourceIBMISFlowLogsID(d))
	d.Set(isFlowLogs, flowlogsInfo)
	return nil
}

// dataSourceIBMISFlowLogsID returns a reasonable ID for a flowlogCollector list.
func dataSourceIBMISFlowLogsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
