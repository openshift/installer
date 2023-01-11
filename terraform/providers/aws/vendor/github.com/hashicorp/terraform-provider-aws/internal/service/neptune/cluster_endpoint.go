package neptune

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/neptune"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceClusterEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterEndpointCreate,
		Read:   resourceClusterEndpointRead,
		Update: resourceClusterEndpointUpdate,
		Delete: resourceClusterEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validIdentifier,
			},
			"cluster_endpoint_identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validIdentifier,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"READER", "WRITER", "ANY"}, false),
			},
			"static_members": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"excluded_members": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClusterEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).NeptuneConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &neptune.CreateDBClusterEndpointInput{
		DBClusterEndpointIdentifier: aws.String(d.Get("cluster_endpoint_identifier").(string)),
		DBClusterIdentifier:         aws.String(d.Get("cluster_identifier").(string)),
		EndpointType:                aws.String(d.Get("endpoint_type").(string)),
	}

	if attr := d.Get("static_members").(*schema.Set); attr.Len() > 0 {
		input.StaticMembers = flex.ExpandStringSet(attr)
	}

	if attr := d.Get("excluded_members").(*schema.Set); attr.Len() > 0 {
		input.ExcludedMembers = flex.ExpandStringSet(attr)
	}

	// Tags are currently only supported in AWS Commercial.
	if len(tags) > 0 && meta.(*conns.AWSClient).Partition == endpoints.AwsPartitionID {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	out, err := conn.CreateDBClusterEndpoint(input)
	if err != nil {
		return fmt.Errorf("creating Neptune Cluster Endpoint: %w", err)
	}

	clusterId := aws.StringValue(out.DBClusterIdentifier)
	endpointId := aws.StringValue(out.DBClusterEndpointIdentifier)
	d.SetId(fmt.Sprintf("%s:%s", clusterId, endpointId))

	_, err = WaitDBClusterEndpointAvailable(conn, d.Id())
	if err != nil {
		return fmt.Errorf("waiting for Neptune Cluster Endpoint (%q) to be Available: %w", d.Id(), err)
	}

	return resourceClusterEndpointRead(d, meta)

}

func resourceClusterEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).NeptuneConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	resp, err := FindEndpointByID(conn, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		d.SetId("")
		log.Printf("[DEBUG] Neptune Cluster Endpoint (%s) not found", d.Id())
		return nil
	}

	if err != nil {
		return fmt.Errorf("describing Neptune Cluster Endpoint (%s): %w", d.Id(), err)
	}

	d.Set("cluster_endpoint_identifier", resp.DBClusterEndpointIdentifier)
	d.Set("cluster_identifier", resp.DBClusterIdentifier)
	d.Set("endpoint_type", resp.CustomEndpointType)
	d.Set("endpoint", resp.Endpoint)
	d.Set("excluded_members", flex.FlattenStringSet(resp.ExcludedMembers))
	d.Set("static_members", flex.FlattenStringSet(resp.StaticMembers))

	arn := aws.StringValue(resp.DBClusterEndpointArn)
	d.Set("arn", arn)

	// Tags are currently only supported in AWS Commercial.
	if meta.(*conns.AWSClient).Partition == endpoints.AwsPartitionID {
		tags, err := ListTags(conn, arn)

		if err != nil {
			return fmt.Errorf("listing tags for Neptune Cluster Endpoint (%s): %w", arn, err)
		}

		tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

		//lintignore:AWSR002
		if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
			return fmt.Errorf("setting tags: %w", err)
		}

		if err := d.Set("tags_all", tags.Map()); err != nil {
			return fmt.Errorf("setting tags_all: %w", err)
		}
	} else {
		d.Set("tags", nil)
		d.Set("tags_all", nil)
	}

	return nil
}

func resourceClusterEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).NeptuneConn

	if d.HasChangesExcept("tags", "tags_all") {
		req := &neptune.ModifyDBClusterEndpointInput{
			DBClusterEndpointIdentifier: aws.String(d.Get("cluster_endpoint_identifier").(string)),
		}

		if d.HasChange("endpoint_type") {
			req.EndpointType = aws.String(d.Get("endpoint_type").(string))
		}

		if d.HasChange("static_members") {
			req.StaticMembers = flex.ExpandStringSet(d.Get("static_members").(*schema.Set))
		}

		if d.HasChange("excluded_members") {
			req.ExcludedMembers = flex.ExpandStringSet(d.Get("excluded_members").(*schema.Set))
		}

		_, err := conn.ModifyDBClusterEndpoint(req)
		if err != nil {
			return fmt.Errorf("updating Neptune Cluster Endpoint (%q): %w", d.Id(), err)
		}

		_, err = WaitDBClusterEndpointAvailable(conn, d.Id())
		if err != nil {
			return fmt.Errorf("waiting for Neptune Cluster Endpoint (%q) to be Available: %w", d.Id(), err)
		}
	}

	// Tags are currently only supported in AWS Commercial.
	if d.HasChange("tags_all") && meta.(*conns.AWSClient).Partition == endpoints.AwsPartitionID {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating Neptune Cluster Endpoint (%s) tags: %w", d.Get("arn").(string), err)
		}
	}

	return resourceClusterEndpointRead(d, meta)
}

func resourceClusterEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).NeptuneConn

	endpointId := d.Get("cluster_endpoint_identifier").(string)
	input := &neptune.DeleteDBClusterEndpointInput{
		DBClusterEndpointIdentifier: aws.String(endpointId),
	}

	_, err := conn.DeleteDBClusterEndpoint(input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, neptune.ErrCodeDBClusterEndpointNotFoundFault) ||
			tfawserr.ErrCodeEquals(err, neptune.ErrCodeDBClusterNotFoundFault) {
			return nil
		}
		return fmt.Errorf("Neptune Cluster Endpoint cannot be deleted: %w", err)
	}
	_, err = WaitDBClusterEndpointDeleted(conn, d.Id())
	if err != nil {
		if tfawserr.ErrCodeEquals(err, neptune.ErrCodeDBClusterEndpointNotFoundFault) {
			return nil
		}
		return fmt.Errorf("waiting for Neptune Cluster Endpoint (%q) to be Deleted: %w", d.Id(), err)
	}

	return nil
}
