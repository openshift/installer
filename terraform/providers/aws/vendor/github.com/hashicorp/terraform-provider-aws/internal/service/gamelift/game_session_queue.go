package gamelift

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceGameSessionQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceGameSessionQueueCreate,
		Read:   resourceGameSessionQueueRead,
		Update: resourceGameSessionQueueUpdate,
		Delete: resourceGameSessionQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destinations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"notification_target": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},
			"player_latency_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maximum_individual_player_latency_milliseconds": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"policy_duration_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"timeout_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 600),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceGameSessionQueueCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GameLiftConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := gamelift.CreateGameSessionQueueInput{
		Name:                  aws.String(d.Get("name").(string)),
		Destinations:          expandGameSessionQueueDestinations(d.Get("destinations").([]interface{})),
		PlayerLatencyPolicies: expandGameSessionPlayerLatencyPolicies(d.Get("player_latency_policy").([]interface{})),
		TimeoutInSeconds:      aws.Int64(int64(d.Get("timeout_in_seconds").(int))),
		Tags:                  Tags(tags.IgnoreAWS()),
	}

	if v, ok := d.GetOk("notification_target"); ok {
		input.NotificationTarget = aws.String(v.(string))
	}

	log.Printf("[INFO] Creating GameLift Session Queue: %s", input)
	out, err := conn.CreateGameSessionQueue(&input)
	if err != nil {
		return fmt.Errorf("error creating GameLift Game Session Queue: %s", err)
	}

	d.SetId(aws.StringValue(out.GameSessionQueue.Name))

	return resourceGameSessionQueueRead(d, meta)
}

func resourceGameSessionQueueRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GameLiftConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	log.Printf("[INFO] Describing GameLift Session Queues: %s", d.Id())
	limit := int64(1)
	out, err := conn.DescribeGameSessionQueues(&gamelift.DescribeGameSessionQueuesInput{
		Names: aws.StringSlice([]string{d.Id()}),
		Limit: &limit,
	})
	if err != nil {
		if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
			log.Printf("[WARN] GameLift Session Queues (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading GameLift Game Session Queue (%s): %s", d.Id(), err)
	}
	sessionQueues := out.GameSessionQueues

	if len(sessionQueues) < 1 {
		log.Printf("[WARN] GameLift Session Queue (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if len(sessionQueues) != 1 {
		return fmt.Errorf("expected exactly 1 GameLift Session Queues, found %d under %q",
			len(sessionQueues), d.Id())
	}
	sessionQueue := sessionQueues[0]

	arn := aws.StringValue(sessionQueue.GameSessionQueueArn)
	d.Set("arn", arn)
	d.Set("name", sessionQueue.Name)
	d.Set("notification_target", sessionQueue.NotificationTarget)
	d.Set("timeout_in_seconds", sessionQueue.TimeoutInSeconds)
	if err := d.Set("destinations", flattenGameSessionQueueDestinations(sessionQueue.Destinations)); err != nil {
		return fmt.Errorf("error setting destinations: %s", err)
	}
	if err := d.Set("player_latency_policy", flattenPlayerLatencyPolicies(sessionQueue.PlayerLatencyPolicies)); err != nil {
		return fmt.Errorf("error setting player_latency_policy: %s", err)
	}

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("error listing tags for Game Lift Session Queue (%s): %s", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func flattenGameSessionQueueDestinations(destinations []*gamelift.GameSessionQueueDestination) []interface{} {
	l := make([]interface{}, 0)

	for _, destination := range destinations {
		if destination == nil {
			continue
		}
		l = append(l, aws.StringValue(destination.DestinationArn))
	}

	return l
}

func flattenPlayerLatencyPolicies(playerLatencyPolicies []*gamelift.PlayerLatencyPolicy) []interface{} {
	l := make([]interface{}, 0)
	for _, policy := range playerLatencyPolicies {
		m := map[string]interface{}{
			"maximum_individual_player_latency_milliseconds": aws.Int64Value(policy.MaximumIndividualPlayerLatencyMilliseconds),
			"policy_duration_seconds":                        aws.Int64Value(policy.PolicyDurationSeconds),
		}
		l = append(l, m)
	}
	return l
}

func resourceGameSessionQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GameLiftConn

	log.Printf("[INFO] Updating GameLift Session Queue: %s", d.Id())

	input := gamelift.UpdateGameSessionQueueInput{
		Name:                  aws.String(d.Id()),
		Destinations:          expandGameSessionQueueDestinations(d.Get("destinations").([]interface{})),
		PlayerLatencyPolicies: expandGameSessionPlayerLatencyPolicies(d.Get("player_latency_policy").([]interface{})),
		TimeoutInSeconds:      aws.Int64(int64(d.Get("timeout_in_seconds").(int))),
	}

	if v, ok := d.GetOk("notification_target"); ok {
		input.NotificationTarget = aws.String(v.(string))
	}

	_, err := conn.UpdateGameSessionQueue(&input)
	if err != nil {
		return fmt.Errorf("error updating GameLift Game Session Queue (%s): %s", d.Id(), err)
	}

	arn := d.Get("arn").(string)
	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, arn, o, n); err != nil {
			return fmt.Errorf("error updating Game Lift Session Queue (%s) tags: %s", arn, err)
		}
	}

	return resourceGameSessionQueueRead(d, meta)
}

func resourceGameSessionQueueDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).GameLiftConn
	log.Printf("[INFO] Deleting GameLift Session Queue: %s", d.Id())
	_, err := conn.DeleteGameSessionQueue(&gamelift.DeleteGameSessionQueueInput{
		Name: aws.String(d.Id()),
	})
	if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error deleting GameLift Game Session Queue (%s): %s", d.Id(), err)
	}

	return nil
}

func expandGameSessionQueueDestinations(destinationsMap []interface{}) []*gamelift.GameSessionQueueDestination {
	if len(destinationsMap) < 1 {
		return nil
	}
	var destinations []*gamelift.GameSessionQueueDestination
	for _, destination := range destinationsMap {
		destinations = append(
			destinations,
			&gamelift.GameSessionQueueDestination{
				DestinationArn: aws.String(destination.(string)),
			})
	}
	return destinations
}

func expandGameSessionPlayerLatencyPolicies(destinationsPlayerLatencyPolicyMap []interface{}) []*gamelift.PlayerLatencyPolicy {
	if len(destinationsPlayerLatencyPolicyMap) < 1 {
		return nil
	}
	var playerLatencyPolicies []*gamelift.PlayerLatencyPolicy
	for _, playerLatencyPolicy := range destinationsPlayerLatencyPolicyMap {
		item := playerLatencyPolicy.(map[string]interface{})
		playerLatencyPolicies = append(
			playerLatencyPolicies,
			&gamelift.PlayerLatencyPolicy{
				MaximumIndividualPlayerLatencyMilliseconds: aws.Int64(int64(item["maximum_individual_player_latency_milliseconds"].(int))),
				PolicyDurationSeconds:                      aws.Int64(int64(item["policy_duration_seconds"].(int))),
			})
	}
	return playerLatencyPolicies
}
