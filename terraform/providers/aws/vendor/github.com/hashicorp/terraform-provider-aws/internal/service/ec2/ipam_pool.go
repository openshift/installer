package ec2

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_vpc_ipam_pool", name="IPAM Pool")
// @Tags(identifierAttribute="id")
func ResourceIPAMPool() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIPAMPoolCreate,
		ReadWithoutTimeout:   resourceIPAMPoolRead,
		UpdateWithoutTimeout: resourceIPAMPoolUpdate,
		DeleteWithoutTimeout: resourceIPAMPoolDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ec2.AddressFamily_Values(), false),
			},
			"allocation_default_netmask_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 128),
			},
			"allocation_max_netmask_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 128),
			},
			"allocation_min_netmask_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 128),
			},
			"allocation_resource_tags": tftags.TagsSchema(),
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_import": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aws_service": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ec2.IpamPoolAwsService_Values(), false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_scope_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipam_scope_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.StringInSlice([]string{"None"}, false),
					verify.ValidRegionName,
				),
				Default: "None",
			},
			"pool_depth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_ip_source": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ec2.IpamPoolPublicIpSource_Values(), false),
				// default is byoip when AddressFamily = ipv6
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "byoip" && n == "" {
						return true
					}
					return false
				},
			},
			"publicly_advertisable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceIPAMPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	addressFamily := d.Get("address_family").(string)
	input := &ec2.CreateIpamPoolInput{
		AddressFamily:     aws.String(addressFamily),
		ClientToken:       aws.String(id.UniqueId()),
		IpamScopeId:       aws.String(d.Get("ipam_scope_id").(string)),
		TagSpecifications: getTagSpecificationsIn(ctx, ec2.ResourceTypeIpamPool),
	}

	if v, ok := d.GetOk("allocation_default_netmask_length"); ok {
		input.AllocationDefaultNetmaskLength = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("allocation_max_netmask_length"); ok {
		input.AllocationMaxNetmaskLength = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("allocation_min_netmask_length"); ok {
		input.AllocationMinNetmaskLength = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("allocation_resource_tags"); ok && len(v.(map[string]interface{})) > 0 {
		input.AllocationResourceTags = ipamResourceTags(tftags.New(ctx, v.(map[string]interface{})))
	}

	if v, ok := d.GetOk("auto_import"); ok {
		input.AutoImport = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("locale"); ok && v != "None" {
		input.Locale = aws.String(v.(string))
	}

	if v, ok := d.GetOk("aws_service"); ok {
		input.AwsService = aws.String(v.(string))
	}

	var publicIpSource string
	if v, ok := d.GetOk("public_ip_source"); ok {
		publicIpSource = v.(string)
		input.PublicIpSource = aws.String(publicIpSource)
	}

	// PubliclyAdvertisable must be set if if the AddressFamily is IPv6 and PublicIpSource is byoip.
	// The request can only contain PubliclyAdvertisable if the AddressFamily is IPv6 and PublicIpSource is byoip.
	if addressFamily == ec2.AddressFamilyIpv6 && publicIpSource != ec2.IpamPoolPublicIpSourceAmazon {
		input.PubliclyAdvertisable = aws.Bool(d.Get("publicly_advertisable").(bool))
	}

	if v, ok := d.GetOk("source_ipam_pool_id"); ok {
		input.SourceIpamPoolId = aws.String(v.(string))
	}

	output, err := conn.CreateIpamPoolWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating IPAM Pool: %s", err)
	}

	d.SetId(aws.StringValue(output.IpamPool.IpamPoolId))

	if _, err := WaitIPAMPoolCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for IPAM Pool (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceIPAMPoolRead(ctx, d, meta)...)
}

func resourceIPAMPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	pool, err := FindIPAMPoolByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IPAM Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading IPAM Pool (%s): %s", d.Id(), err)
	}

	d.Set("address_family", pool.AddressFamily)
	d.Set("allocation_resource_tags", KeyValueTags(ctx, tagsFromIPAMAllocationTags(pool.AllocationResourceTags)).Map())
	d.Set("arn", pool.IpamPoolArn)
	d.Set("auto_import", pool.AutoImport)
	d.Set("aws_service", pool.AwsService)
	d.Set("description", pool.Description)
	scopeID := strings.Split(aws.StringValue(pool.IpamScopeArn), "/")[1]
	d.Set("ipam_scope_id", scopeID)
	d.Set("ipam_scope_type", pool.IpamScopeType)
	d.Set("locale", pool.Locale)
	d.Set("pool_depth", pool.PoolDepth)
	d.Set("publicly_advertisable", pool.PubliclyAdvertisable)
	d.Set("public_ip_source", pool.PublicIpSource)
	d.Set("source_ipam_pool_id", pool.SourceIpamPoolId)
	d.Set("state", pool.State)

	SetTagsOut(ctx, pool.Tags)

	return diags
}

func resourceIPAMPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &ec2.ModifyIpamPoolInput{
			IpamPoolId: aws.String(d.Id()),
		}

		if v, ok := d.GetOk("allocation_default_netmask_length"); ok {
			input.AllocationDefaultNetmaskLength = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("allocation_max_netmask_length"); ok {
			input.AllocationMaxNetmaskLength = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("allocation_min_netmask_length"); ok {
			input.AllocationMinNetmaskLength = aws.Int64(int64(v.(int)))
		}

		if d.HasChange("allocation_resource_tags") {
			o, n := d.GetChange("allocation_resource_tags")
			oldTags := tftags.New(ctx, o)
			newTags := tftags.New(ctx, n)

			if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
				input.RemoveAllocationResourceTags = ipamResourceTags(removedTags.IgnoreAWS())
			}

			if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
				input.AddAllocationResourceTags = ipamResourceTags(updatedTags.IgnoreAWS())
			}
		}

		if v, ok := d.GetOk("auto_import"); ok {
			input.AutoImport = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("description"); ok {
			input.Description = aws.String(v.(string))
		}

		_, err := conn.ModifyIpamPoolWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating IPAM Pool (%s): %s", d.Id(), err)
		}

		if _, err := WaitIPAMPoolUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for IPAM Pool (%s) update: %s", d.Id(), err)
		}
	}

	return append(diags, resourceIPAMPoolRead(ctx, d, meta)...)
}

func resourceIPAMPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	log.Printf("[DEBUG] Deleting IPAM Pool: %s", d.Id())
	_, err := conn.DeleteIpamPoolWithContext(ctx, &ec2.DeleteIpamPoolInput{
		IpamPoolId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidIPAMPoolIdNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting IPAM Pool (%s): %s", d.Id(), err)
	}

	if _, err = WaitIPAMPoolDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for IPAM Pool (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func ipamResourceTags(tags tftags.KeyValueTags) []*ec2.RequestIpamResourceTag {
	result := make([]*ec2.RequestIpamResourceTag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &ec2.RequestIpamResourceTag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

func tagsFromIPAMAllocationTags(rts []*ec2.IpamResourceTag) []*ec2.Tag {
	if len(rts) == 0 {
		return nil
	}

	tags := []*ec2.Tag{}
	for _, ts := range rts {
		tags = append(tags, &ec2.Tag{
			Key:   ts.Key,
			Value: ts.Value,
		})
	}

	return tags
}
