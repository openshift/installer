package ec2

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ami_copy", name="AMI")
// @Tags(identifierAttribute="id")
func ResourceAMICopy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAMICopyCreate,
		// The remaining operations are shared with the generic aws_ami resource,
		// since the aws_ami_copy resource only differs in how it's created.
		ReadWithoutTimeout:   resourceAMIRead,
		UpdateWithoutTimeout: resourceAMIUpdate,
		DeleteWithoutTimeout: resourceAMIDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(amiRetryTimeout),
			Update: schema.DefaultTimeout(amiRetryTimeout),
			Delete: schema.DefaultTimeout(amiDeleteTimeout),
		},

		// Keep in sync with aws_ami's schema.
		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"boot_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprecation_time": {
				Type:                  schema.TypeString,
				Optional:              true,
				ValidateFunc:          validation.IsRFC3339Time,
				DiffSuppressFunc:      verify.SuppressEquivalentRoundedTime(time.RFC3339, time.Minute),
				DiffSuppressOnRefresh: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_outpost_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},
			// The following block device attributes intentionally mimick the
			// corresponding attributes on aws_instance, since they have the
			// same meaning.
			// However, we don't use root_block_device here because the constraint
			// on which root device attributes can be overridden for an instance to
			// not apply when registering an AMI.
			"ebs_block_device": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"outpost_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"throughput": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["device_name"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["snapshot_id"].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			"ena_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"ephemeral_block_device": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["device_name"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["virtual_name"].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			"hypervisor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_owner_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"imds_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kernel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			// Not a public attribute; used to let the aws_ami_copy and aws_ami_from_instance
			// resources record that they implicitly created new EBS snapshots that we should
			// now manage. Not set by aws_ami, since the snapshots used there are presumed to
			// be independently managed.
			"manage_ebs_snapshots": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_details": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ramdisk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ami_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_ami_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sriov_net_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"tpm_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage_operation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtualization_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceAMICopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	name := d.Get("name").(string)
	sourceImageID := d.Get("source_ami_id").(string)
	input := &ec2.CopyImageInput{
		ClientToken:   aws.String(id.UniqueId()),
		Description:   aws.String(d.Get("description").(string)),
		Encrypted:     aws.Bool(d.Get("encrypted").(bool)),
		Name:          aws.String(name),
		SourceImageId: aws.String(sourceImageID),
		SourceRegion:  aws.String(d.Get("source_ami_region").(string)),
	}

	if v, ok := d.GetOk("destination_outpost_arn"); ok {
		input.DestinationOutpostArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
	}

	output, err := conn.CopyImageWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 AMI (%s) from source EC2 AMI (%s): %s", name, sourceImageID, err)
	}

	d.SetId(aws.StringValue(output.ImageId))
	d.Set("manage_ebs_snapshots", true)

	if err := createTags(ctx, conn, d.Id(), GetTagsIn(ctx)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting EC2 AMI (%s) tags: %s", d.Id(), err)
	}

	if _, err := WaitImageAvailable(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 AMI (%s) from source EC2 AMI (%s): waiting for completion: %s", name, sourceImageID, err)
	}

	if v, ok := d.GetOk("deprecation_time"); ok {
		if err := enableImageDeprecation(ctx, conn, d.Id(), v.(string)); err != nil {
			return sdkdiag.AppendErrorf(diags, "creating EC2 AMI (%s) from source EC2 AMI (%s): %s", name, sourceImageID, err)
		}
	}

	return append(diags, resourceAMIRead(ctx, d, meta)...)
}
