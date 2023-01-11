package ec2

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceTrafficMirrorSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficMirrorSessionCreate,
		Update: resourceTrafficMirrorSessionUpdate,
		Read:   resourceTrafficMirrorSessionRead,
		Delete: resourceTrafficMirrorSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"packet_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"session_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 32766),
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_network_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16777216),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
	}
}

func resourceTrafficMirrorSessionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &ec2.CreateTrafficMirrorSessionInput{
		NetworkInterfaceId:    aws.String(d.Get("network_interface_id").(string)),
		TrafficMirrorFilterId: aws.String(d.Get("traffic_mirror_filter_id").(string)),
		TrafficMirrorTargetId: aws.String(d.Get("traffic_mirror_target_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("packet_length"); ok {
		input.PacketLength = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("session_number"); ok {
		input.SessionNumber = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("virtual_network_id"); ok {
		input.VirtualNetworkId = aws.Int64(int64(v.(int)))
	}

	if len(tags) > 0 {
		input.TagSpecifications = tagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeTrafficMirrorSession)
	}

	out, err := conn.CreateTrafficMirrorSession(input)
	if nil != err {
		return fmt.Errorf("Error creating traffic mirror session %v", err)
	}

	d.SetId(aws.StringValue(out.TrafficMirrorSession.TrafficMirrorSessionId))
	return resourceTrafficMirrorSessionRead(d, meta)
}

func resourceTrafficMirrorSessionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	sessionId := d.Id()
	input := &ec2.ModifyTrafficMirrorSessionInput{
		TrafficMirrorSessionId: &sessionId,
	}

	if d.HasChange("session_number") {
		n := d.Get("session_number")
		input.SessionNumber = aws.Int64(int64(n.(int)))
	}

	if d.HasChange("traffic_mirror_filter_id") {
		n := d.Get("traffic_mirror_filter_id")
		input.TrafficMirrorFilterId = aws.String(n.(string))
	}

	if d.HasChange("traffic_mirror_target_id") {
		n := d.Get("traffic_mirror_target_id")
		input.TrafficMirrorTargetId = aws.String(n.(string))
	}

	var removeFields []*string
	if d.HasChange("description") {
		n := d.Get("description")
		if n != "" {
			input.Description = aws.String(n.(string))
		} else {
			removeFields = append(removeFields, aws.String("description"))
		}
	}

	if d.HasChange("packet_length") {
		n := d.Get("packet_length")
		if nil != n && n.(int) > 0 {
			input.PacketLength = aws.Int64(int64(n.(int)))
		} else {
			removeFields = append(removeFields, aws.String("packet-length"))
		}
	}
	log.Printf("[DEBUG] removeFields %v", removeFields)

	if d.HasChange("virtual_network_id") {
		n := d.Get("virtual_network_id")
		log.Printf("[DEBUG] VNI has change %v", n)
		if nil != n && n.(int) > 0 {
			input.VirtualNetworkId = aws.Int64(int64(n.(int)))
		} else {
			removeFields = append(removeFields, aws.String("virtual-network-id"))
		}
	}

	log.Printf("[DEBUG] removeFields %v", removeFields)
	if len(removeFields) > 0 {
		input.SetRemoveFields(removeFields)
	}

	_, err := conn.ModifyTrafficMirrorSession(input)
	if nil != err {
		return fmt.Errorf("Error updating traffic mirror session %v", err)
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EC2 Traffic Mirror Session (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceTrafficMirrorSessionRead(d, meta)
}

func resourceTrafficMirrorSessionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	sessionId := d.Id()
	input := &ec2.DescribeTrafficMirrorSessionsInput{
		TrafficMirrorSessionIds: []*string{
			&sessionId,
		},
	}

	out, err := conn.DescribeTrafficMirrorSessions(input)

	if tfawserr.ErrCodeEquals(err, "InvalidTrafficMirrorSessionId.NotFound") {
		log.Printf("[WARN] EC2 Traffic Mirror Session (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if nil != err {
		return fmt.Errorf("error describing EC2 Traffic Mirror Session (%s): %w", sessionId, err)
	}

	if len(out.TrafficMirrorSessions) == 0 {
		log.Printf("[WARN] EC2 Traffic Mirror Session (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	session := out.TrafficMirrorSessions[0]
	d.Set("network_interface_id", session.NetworkInterfaceId)
	d.Set("session_number", session.SessionNumber)
	d.Set("traffic_mirror_filter_id", session.TrafficMirrorFilterId)
	d.Set("traffic_mirror_target_id", session.TrafficMirrorTargetId)
	d.Set("description", session.Description)
	d.Set("packet_length", session.PacketLength)
	d.Set("virtual_network_id", session.VirtualNetworkId)

	tags := KeyValueTags(session.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	d.Set("owner_id", session.OwnerId)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: aws.StringValue(session.OwnerId),
		Resource:  fmt.Sprintf("traffic-mirror-session/%s", d.Id()),
	}.String()
	d.Set("arn", arn)

	return nil
}

func resourceTrafficMirrorSessionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	sessionId := d.Id()
	input := &ec2.DeleteTrafficMirrorSessionInput{
		TrafficMirrorSessionId: &sessionId,
	}

	_, err := conn.DeleteTrafficMirrorSession(input)
	if nil != err {
		return fmt.Errorf("error deleting EC2 Traffic Mirror Session (%s): %w", sessionId, err)
	}

	return nil
}
