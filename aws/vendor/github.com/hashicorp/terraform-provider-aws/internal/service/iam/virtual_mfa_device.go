package iam

import (
	"context"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_iam_virtual_mfa_device", name="Virtual MFA Device")
// @Tags
func ResourceVirtualMFADevice() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVirtualMFADeviceCreate,
		ReadWithoutTimeout:   resourceVirtualMFADeviceRead,
		UpdateWithoutTimeout: resourceVirtualMFADeviceUpdate,
		DeleteWithoutTimeout: resourceVirtualMFADeviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_32_string_seed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/",
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"qr_code_png": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"virtual_mfa_device_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`[\w+=,.@-]+`),
					"must consist of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: _+=,.@-",
				),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceVirtualMFADeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IAMConn(ctx)

	name := d.Get("virtual_mfa_device_name").(string)
	input := &iam.CreateVirtualMFADeviceInput{
		Path:                 aws.String(d.Get("path").(string)),
		Tags:                 GetTagsIn(ctx),
		VirtualMFADeviceName: aws.String(name),
	}

	output, err := conn.CreateVirtualMFADeviceWithContext(ctx, input)

	// Some partitions (e.g. ISO) may not support tag-on-create.
	if input.Tags != nil && errs.IsUnsupportedOperationInPartitionError(conn.PartitionID, err) {
		input.Tags = nil

		output, err = conn.CreateVirtualMFADeviceWithContext(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating IAM Virtual MFA Device (%s): %s", name, err)
	}

	vMFA := output.VirtualMFADevice
	d.SetId(aws.StringValue(vMFA.SerialNumber))

	d.Set("base_32_string_seed", string(vMFA.Base32StringSeed))
	d.Set("qr_code_png", string(vMFA.QRCodePNG))

	// For partitions not supporting tag-on-create, attempt tag after create.
	if tags := GetTagsIn(ctx); input.Tags == nil && len(tags) > 0 {
		err := virtualMFADeviceCreateTags(ctx, conn, d.Id(), tags)

		// If default tags only, continue. Otherwise, error.
		if v, ok := d.GetOk(names.AttrTags); (!ok || len(v.(map[string]interface{})) == 0) && errs.IsUnsupportedOperationInPartitionError(conn.PartitionID, err) {
			return append(diags, resourceVirtualMFADeviceRead(ctx, d, meta)...)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting IAM Virtual MFA Device (%s) tags: %s", d.Id(), err)
		}
	}

	return append(diags, resourceVirtualMFADeviceRead(ctx, d, meta)...)
}

func resourceVirtualMFADeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IAMConn(ctx)

	vMFA, err := FindVirtualMFADeviceBySerialNumber(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IAM Virtual MFA Device (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading IAM Virtual MFA Device (%s): %s", d.Id(), err)
	}

	d.Set("arn", vMFA.SerialNumber)

	// The call above returns empty tags.
	output, err := conn.ListMFADeviceTagsWithContext(ctx, &iam.ListMFADeviceTagsInput{
		SerialNumber: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing IAM Virtual MFA Device (%s) tags: %s", d.Id(), err)
	}

	SetTagsOut(ctx, output.Tags)

	return diags
}

func resourceVirtualMFADeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IAMConn(ctx)

	o, n := d.GetChange("tags_all")

	err := virtualMFADeviceUpdateTags(ctx, conn, d.Id(), o, n)

	// Some partitions (e.g. ISO) may not support tagging.
	if errs.IsUnsupportedOperationInPartitionError(conn.PartitionID, err) {
		return append(diags, resourceVirtualMFADeviceRead(ctx, d, meta)...)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating tags for IAM Virtual MFA Device (%s): %s", d.Id(), err)
	}

	return append(diags, resourceVirtualMFADeviceRead(ctx, d, meta)...)
}

func resourceVirtualMFADeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IAMConn(ctx)

	log.Printf("[INFO] Deleting IAM Virtual MFA Device: %s", d.Id())
	_, err := conn.DeleteVirtualMFADeviceWithContext(ctx, &iam.DeleteVirtualMFADeviceInput{
		SerialNumber: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting IAM Virtual MFA Device (%s): %s", d.Id(), err)
	}

	return diags
}

func FindVirtualMFADeviceBySerialNumber(ctx context.Context, conn *iam.IAM, serialNumber string) (*iam.VirtualMFADevice, error) {
	input := &iam.ListVirtualMFADevicesInput{}
	var output *iam.VirtualMFADevice

	err := conn.ListVirtualMFADevicesPagesWithContext(ctx, input, func(page *iam.ListVirtualMFADevicesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.VirtualMFADevices {
			if v != nil && aws.StringValue(v.SerialNumber) == serialNumber {
				output = v
				return false
			}
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, &retry.NotFoundError{}
	}

	return output, nil
}
