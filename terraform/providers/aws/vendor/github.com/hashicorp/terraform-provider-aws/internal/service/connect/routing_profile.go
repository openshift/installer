package connect

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	AssociateRoutingProfileQueuesMaxItems    = 10
	DisassociateRoutingProfileQueuesMaxItems = 10
	CreateRoutingProfileQueuesMaxItems       = 10
)

// @SDKResource("aws_connect_routing_profile", name="Routing Profile")
// @Tags(identifierAttribute="arn")
func ResourceRoutingProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRoutingProfileCreate,
		ReadWithoutTimeout:   resourceRoutingProfileRead,
		UpdateWithoutTimeout: resourceRoutingProfileUpdate,
		// Routing profiles do not support deletion today. NoOp the Delete method.
		// Users can rename their routing profiles manually if they want.
		DeleteWithoutTimeout: schema.NoopContext,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_outbound_queue_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 250),
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"media_concurrencies": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(connect.Channel_Values(), false), // Valid values: VOICE | CHAT | TASK
						},
						"concurrency": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 10),
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 127),
			},
			"queue_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(connect.Channel_Values(), false), // Valid values: VOICE | CHAT | TASK
						},
						"delay": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 9999),
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 99),
						},
						"queue_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"routing_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceRoutingProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn(ctx)

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	input := &connect.CreateRoutingProfileInput{
		DefaultOutboundQueueId: aws.String(d.Get("default_outbound_queue_id").(string)),
		Description:            aws.String(d.Get("description").(string)),
		InstanceId:             aws.String(instanceID),
		MediaConcurrencies:     expandRoutingProfileMediaConcurrencies(d.Get("media_concurrencies").(*schema.Set).List()),
		Name:                   aws.String(name),
		Tags:                   GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("queue_configs"); ok && v.(*schema.Set).Len() > 0 && v.(*schema.Set).Len() <= CreateRoutingProfileQueuesMaxItems {
		input.QueueConfigs = expandRoutingProfileQueueConfigs(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] Creating Connect Routing Profile %s", input)
	output, err := conn.CreateRoutingProfileWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating Connect Routing Profile (%s): %s", name, err)
	}

	if output == nil {
		return diag.Errorf("creating Connect Routing Profile (%s): empty output", name)
	}

	// call the batched association API if the number of queues to associate with the routing profile is > CreateRoutingProfileQueuesMaxItems
	if v, ok := d.GetOk("queue_configs"); ok && v.(*schema.Set).Len() > CreateRoutingProfileQueuesMaxItems {
		queueConfigsUpdateRemove := make([]interface{}, 0)
		err = updateQueueConfigs(ctx, conn, instanceID, aws.StringValue(output.RoutingProfileId), v.(*schema.Set).List(), queueConfigsUpdateRemove)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", instanceID, aws.StringValue(output.RoutingProfileId)))

	return resourceRoutingProfileRead(ctx, d, meta)
}

func resourceRoutingProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn(ctx)

	instanceID, routingProfileID, err := RoutingProfileParseID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.DescribeRoutingProfileWithContext(ctx, &connect.DescribeRoutingProfileInput{
		InstanceId:       aws.String(instanceID),
		RoutingProfileId: aws.String(routingProfileID),
	})

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, connect.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Connect Routing Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("getting Connect Routing Profile (%s): %s", d.Id(), err)
	}

	if resp == nil || resp.RoutingProfile == nil {
		return diag.Errorf("getting Connect Routing Profile (%s): empty response", d.Id())
	}

	routingProfile := resp.RoutingProfile

	if err := d.Set("media_concurrencies", flattenRoutingProfileMediaConcurrencies(routingProfile.MediaConcurrencies)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("arn", routingProfile.RoutingProfileArn)
	d.Set("default_outbound_queue_id", routingProfile.DefaultOutboundQueueId)
	d.Set("description", routingProfile.Description)
	d.Set("instance_id", instanceID)
	d.Set("name", routingProfile.Name)

	d.Set("routing_profile_id", routingProfile.RoutingProfileId)

	// getting the routing profile queues uses a separate API: ListRoutingProfileQueues
	queueConfigs, err := getRoutingProfileQueueConfigs(ctx, conn, instanceID, routingProfileID)

	if err != nil {
		return diag.Errorf("finding Connect Routing Profile Queue Configs Summary by Routing Profile ID (%s): %s", routingProfileID, err)
	}

	d.Set("queue_configs", queueConfigs)

	SetTagsOut(ctx, resp.RoutingProfile.Tags)

	return nil
}

func resourceRoutingProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn(ctx)

	instanceID, routingProfileID, err := RoutingProfileParseID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// RoutingProfile has 4 update APIs
	// UpdateRoutingProfileConcurrency: Updates the channels that agents can handle in the Contact Control Panel (CCP) for a routing profile.
	// UpdateRoutingProfileDefaultOutboundQueue: Updates the default outbound queue of a routing profile.
	// UpdateRoutingProfileName: Updates the name and description of a routing profile.
	// UpdateRoutingProfileQueues: Updates the properties associated with a set of queues for a routing profile.

	// updates to concurrency
	inputConcurrency := &connect.UpdateRoutingProfileConcurrencyInput{
		InstanceId:       aws.String(instanceID),
		RoutingProfileId: aws.String(routingProfileID),
	}

	if d.HasChange("media_concurrencies") {
		mediaConcurrencies := expandRoutingProfileMediaConcurrencies(d.Get("media_concurrencies").(*schema.Set).List())
		inputConcurrency.MediaConcurrencies = mediaConcurrencies
		_, err = conn.UpdateRoutingProfileConcurrencyWithContext(ctx, inputConcurrency)
		if err != nil {
			return diag.Errorf("updating RoutingProfile Media Concurrency (%s): %s", d.Id(), err)
		}
	}

	// updates to default outbound queue id
	inputDefaultOutboundQueue := &connect.UpdateRoutingProfileDefaultOutboundQueueInput{
		InstanceId:       aws.String(instanceID),
		RoutingProfileId: aws.String(routingProfileID),
	}

	if d.HasChange("default_outbound_queue_id") {
		inputDefaultOutboundQueue.DefaultOutboundQueueId = aws.String(d.Get("default_outbound_queue_id").(string))
		_, err = conn.UpdateRoutingProfileDefaultOutboundQueueWithContext(ctx, inputDefaultOutboundQueue)

		if err != nil {
			return diag.Errorf("updating RoutingProfile Default Outbound Queue ID (%s): %s", d.Id(), err)
		}
	}

	// updates to name and/or description
	inputNameDesc := &connect.UpdateRoutingProfileNameInput{
		InstanceId:       aws.String(instanceID),
		RoutingProfileId: aws.String(routingProfileID),
	}

	if d.HasChanges("name", "description") {
		inputNameDesc.Name = aws.String(d.Get("name").(string))
		inputNameDesc.Description = aws.String(d.Get("description").(string))
		_, err = conn.UpdateRoutingProfileNameWithContext(ctx, inputNameDesc)

		if err != nil {
			return diag.Errorf("updating RoutingProfile Name (%s): %s", d.Id(), err)
		}
	}

	// updates to queue configs
	// There are 3 APIs for this
	// AssociateRoutingProfileQueues - Associates a set of queues with a routing profile.
	// DisassociateRoutingProfileQueues - Disassociates a set of queues from a routing profile.
	// UpdateRoutingProfileQueues - Updates the properties associated with a set of queues for a routing profile.
	// since the update only updates the existing queues that are associated, we will instead disassociate and associate
	// the respective queues based on the diff detected
	if d.HasChange("queue_configs") {
		o, n := d.GetChange("queue_configs")

		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		queueConfigsUpdateAdd := ns.Difference(os).List()
		queueConfigsUpdateRemove := os.Difference(ns).List()

		err = updateQueueConfigs(ctx, conn, instanceID, routingProfileID, queueConfigsUpdateAdd, queueConfigsUpdateRemove)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRoutingProfileRead(ctx, d, meta)
}

func updateQueueConfigs(ctx context.Context, conn *connect.Connect, instanceID, routingProfileID string, queueConfigsUpdateAdd, queueConfigsUpdateRemove []interface{}) error {
	// updates to queue configs
	// There are 3 APIs for this
	// AssociateRoutingProfileQueues - Associates a set of queues with a routing profile.
	// DisassociateRoutingProfileQueues - Disassociates a set of queues from a routing profile.
	// UpdateRoutingProfileQueues - Updates the properties associated with a set of queues for a routing profile.
	// since the update only updates the existing queues that are associated, we will instead disassociate and associate
	// the respective queues based on the diff detected

	// disassociate first since Queue and channel type combination cannot be duplicated
	if len(queueConfigsUpdateRemove) > 0 {
		for i := 0; i < len(queueConfigsUpdateRemove); i += DisassociateRoutingProfileQueuesMaxItems {
			j := i + DisassociateRoutingProfileQueuesMaxItems
			if j > len(queueConfigsUpdateRemove) {
				j = len(queueConfigsUpdateRemove)
			}
			_, err := conn.DisassociateRoutingProfileQueuesWithContext(ctx, &connect.DisassociateRoutingProfileQueuesInput{
				InstanceId:       aws.String(instanceID),
				QueueReferences:  expandRoutingProfileQueueReferences(queueConfigsUpdateRemove[i:j]),
				RoutingProfileId: aws.String(routingProfileID),
			})
			if err != nil {
				return fmt.Errorf("updating RoutingProfile Queue Configs, specifically disassociating queues from routing profile (%s): %s", routingProfileID, err)
			}
		}
	}

	if len(queueConfigsUpdateAdd) > 0 {
		for i := 0; i < len(queueConfigsUpdateAdd); i += AssociateRoutingProfileQueuesMaxItems {
			j := i + AssociateRoutingProfileQueuesMaxItems
			if j > len(queueConfigsUpdateAdd) {
				j = len(queueConfigsUpdateAdd)
			}
			_, err := conn.AssociateRoutingProfileQueuesWithContext(ctx, &connect.AssociateRoutingProfileQueuesInput{
				InstanceId:       aws.String(instanceID),
				QueueConfigs:     expandRoutingProfileQueueConfigs(queueConfigsUpdateAdd[i:j]),
				RoutingProfileId: aws.String(routingProfileID),
			})
			if err != nil {
				return fmt.Errorf("updating RoutingProfile Queue Configs, specifically associating queues to routing profile (%s): %s", routingProfileID, err)
			}
		}
	}

	return nil
}

// func resourceRoutingProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.AWSClient).ConnectConn(ctx)

// 	instanceID, routingProfileID, err := RoutingProfileParseID(d.Id())

// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	_, err = conn.DeleteRoutingProfileWithContext(ctx, &connect.DeleteRoutingProfileInput{
// 		InstanceId:       aws.String(instanceID),
// 		RoutingProfileId: aws.String(routingProfileID),
// 	})

// 	if err != nil {
// 		return diag.Errorf("deleting RoutingProfile (%s): %s", d.Id(), err)
// 	}

// 	return nil
// }

func expandRoutingProfileMediaConcurrencies(mediaConcurrencies []interface{}) []*connect.MediaConcurrency {
	if len(mediaConcurrencies) == 0 {
		return nil
	}

	mediaConcurrenciesExpanded := []*connect.MediaConcurrency{}

	for _, mediaConcurrency := range mediaConcurrencies {
		data := mediaConcurrency.(map[string]interface{})
		mediaConcurrencyExpanded := &connect.MediaConcurrency{
			Channel:     aws.String(data["channel"].(string)),
			Concurrency: aws.Int64(int64(data["concurrency"].(int))),
		}
		mediaConcurrenciesExpanded = append(mediaConcurrenciesExpanded, mediaConcurrencyExpanded)
	}

	return mediaConcurrenciesExpanded
}

func flattenRoutingProfileMediaConcurrencies(mediaConcurrencies []*connect.MediaConcurrency) []interface{} {
	mediaConcurrenciesList := []interface{}{}

	for _, mediaConcurrency := range mediaConcurrencies {
		values := map[string]interface{}{
			"channel":     aws.StringValue(mediaConcurrency.Channel),
			"concurrency": aws.Int64Value(mediaConcurrency.Concurrency),
		}

		mediaConcurrenciesList = append(mediaConcurrenciesList, values)
	}
	return mediaConcurrenciesList
}

func expandRoutingProfileQueueConfigs(queueConfigs []interface{}) []*connect.RoutingProfileQueueConfig {
	if len(queueConfigs) == 0 {
		return nil
	}

	queueConfigsExpanded := []*connect.RoutingProfileQueueConfig{}

	for _, queueConfig := range queueConfigs {
		data := queueConfig.(map[string]interface{})
		queueConfigExpanded := &connect.RoutingProfileQueueConfig{
			Delay:    aws.Int64(int64(data["delay"].(int))),
			Priority: aws.Int64(int64(data["priority"].(int))),
		}

		qr := connect.RoutingProfileQueueReference{
			Channel: aws.String(data["channel"].(string)),
			QueueId: aws.String(data["queue_id"].(string)),
		}
		queueConfigExpanded.QueueReference = &qr

		queueConfigsExpanded = append(queueConfigsExpanded, queueConfigExpanded)
	}

	return queueConfigsExpanded
}

func expandRoutingProfileQueueReferences(queueConfigs []interface{}) []*connect.RoutingProfileQueueReference {
	if len(queueConfigs) == 0 {
		return nil
	}

	queueReferencesExpanded := []*connect.RoutingProfileQueueReference{}

	for _, queueConfig := range queueConfigs {
		data := queueConfig.(map[string]interface{})
		queueReferenceExpanded := &connect.RoutingProfileQueueReference{
			Channel: aws.String(data["channel"].(string)),
			QueueId: aws.String(data["queue_id"].(string)),
		}

		queueReferencesExpanded = append(queueReferencesExpanded, queueReferenceExpanded)
	}

	return queueReferencesExpanded
}

func getRoutingProfileQueueConfigs(ctx context.Context, conn *connect.Connect, instanceID, routingProfileID string) ([]interface{}, error) {
	queueConfigsList := []interface{}{}

	input := &connect.ListRoutingProfileQueuesInput{
		InstanceId:       aws.String(instanceID),
		MaxResults:       aws.Int64(ListRoutingProfileQueuesMaxResults),
		RoutingProfileId: aws.String(routingProfileID),
	}

	err := conn.ListRoutingProfileQueuesPagesWithContext(ctx, input, func(page *connect.ListRoutingProfileQueuesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, qc := range page.RoutingProfileQueueConfigSummaryList {
			if qc == nil {
				continue
			}

			values := map[string]interface{}{
				"channel":    aws.StringValue(qc.Channel),
				"delay":      aws.Int64Value(qc.Delay),
				"priority":   aws.Int64Value(qc.Priority),
				"queue_arn":  aws.StringValue(qc.QueueArn),
				"queue_id":   aws.StringValue(qc.QueueId),
				"queue_name": aws.StringValue(qc.QueueName),
			}

			queueConfigsList = append(queueConfigsList, values)
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return queueConfigsList, nil
}

func RoutingProfileParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected instanceID:routingProfileID", id)
	}

	return parts[0], parts[1], nil
}
