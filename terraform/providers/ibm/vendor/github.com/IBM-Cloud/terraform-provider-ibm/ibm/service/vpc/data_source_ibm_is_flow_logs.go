// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isFlowLogs = "flow_log_collectors"
)

func DataSourceIBMISFlowLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISFlowLogsRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group this flow log belongs to",
			},
			"vpc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc ID this flow log is in",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc name this flow log is in",
			},
			"vpc_crn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc CRN this flow log is in",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the flow log ",
			},
			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The target id of the flow log ",
			},
			"target_resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The target resource type of the flow log ",
			},
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
						isFlowLogTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "Tags for the VPC Flow logs",
						},

						isFlowLogAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access management tags",
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
	listOptions := &vpcv1.ListFlowLogCollectorsOptions{}
	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listOptions.ResourceGroupID = &resGroup
	}
	if nameintf, ok := d.GetOk("name"); ok {
		name := nameintf.(string)
		listOptions.Name = &name
	}
	if vpcIntf, ok := d.GetOk("vpc"); ok {
		vpcid := vpcIntf.(string)
		listOptions.VPCID = &vpcid
	}
	if vpcNameIntf, ok := d.GetOk("vpc_name"); ok {
		vpcName := vpcNameIntf.(string)
		listOptions.VPCName = &vpcName
	}
	if vpcCrnIntf, ok := d.GetOk("vpc_crn"); ok {
		vpcCrn := vpcCrnIntf.(string)
		listOptions.VPCCRN = &vpcCrn
	}
	if targetIntf, ok := d.GetOk("target"); ok {
		target := targetIntf.(string)
		listOptions.TargetID = &target
	}
	if targetTypeIntf, ok := d.GetOk("target_resource_type"); ok {
		targetType := targetTypeIntf.(string)
		listOptions.TargetResourceType = &targetType
	}
	for {

		if start != "" {
			listOptions.Start = &start
		}
		flowlogCollectors, response, err := sess.ListFlowLogCollectors(listOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching Flow Logs for VPC %s\n%s", err, response)
		}
		start = flex.GetNext(flowlogCollectors.Next)
		allrecs = append(allrecs, flowlogCollectors.FlowLogCollectors...)
		if start == "" {
			break
		}
	}
	flowlogsInfo := make([]map[string]interface{}, 0)
	for _, flowlogCollector := range allrecs {

		targetIntf := flowlogCollector.Target
		target := targetIntf.(*vpcv1.FlowLogCollectorTarget)

		tags, err := flex.GetGlobalTagsUsingCRN(meta, *flowlogCollector.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc flow log (%s) tags: %s", d.Id(), err)
		}

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *flowlogCollector.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource VPC Flow Log (%s) access tags: %s", d.Id(), err)
		}

		l := map[string]interface{}{
			"id":                *flowlogCollector.ID,
			"crn":               *flowlogCollector.CRN,
			"href":              *flowlogCollector.Href,
			"name":              *flowlogCollector.Name,
			"resource_group":    *flowlogCollector.ResourceGroup.ID,
			"created_at":        flowlogCollector.CreatedAt.String(),
			"lifecycle_state":   *flowlogCollector.LifecycleState,
			"storage_bucket":    *flowlogCollector.StorageBucket.Name,
			"active":            *flowlogCollector.Active,
			"vpc":               *flowlogCollector.VPC.ID,
			"target":            *target.ID,
			isFlowLogTags:       tags,
			isFlowLogAccessTags: accesstags,
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
