package datasync

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/datasync"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_datasync_location_smb", name="Location SMB")
// @Tags(identifierAttribute="id")
func ResourceLocationSMB() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLocationSMBCreate,
		ReadWithoutTimeout:   resourceLocationSMBRead,
		UpdateWithoutTimeout: resourceLocationSMBUpdate,
		DeleteWithoutTimeout: resourceLocationSMBDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_arns": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"domain": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 253),
			},
			"mount_options": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:         schema.TypeString,
							Default:      datasync.SmbVersionAutomatic,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(datasync.SmbVersion_Values(), false),
						},
					},
				},
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 104),
			},
			"server_hostname": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"subdirectory": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
				/*// Ignore missing trailing slash
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "/" {
						return false
					}
					if strings.TrimSuffix(old, "/") == strings.TrimSuffix(new, "/") {
						return true
					}
					return false
				},
				*/
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 104),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceLocationSMBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.CreateLocationSmbInput{
		AgentArns:      flex.ExpandStringSet(d.Get("agent_arns").(*schema.Set)),
		MountOptions:   expandSMBMountOptions(d.Get("mount_options").([]interface{})),
		Password:       aws.String(d.Get("password").(string)),
		ServerHostname: aws.String(d.Get("server_hostname").(string)),
		Subdirectory:   aws.String(d.Get("subdirectory").(string)),
		Tags:           GetTagsIn(ctx),
		User:           aws.String(d.Get("user").(string)),
	}

	if v, ok := d.GetOk("domain"); ok {
		input.Domain = aws.String(v.(string))
	}

	output, err := conn.CreateLocationSmbWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating DataSync Location SMB: %s", err)
	}

	d.SetId(aws.StringValue(output.LocationArn))

	return append(diags, resourceLocationSMBRead(ctx, d, meta)...)
}

func resourceLocationSMBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.DescribeLocationSmbInput{
		LocationArn: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Reading DataSync Location SMB: %s", input)
	output, err := conn.DescribeLocationSmbWithContext(ctx, input)

	if tfawserr.ErrMessageContains(err, "InvalidRequestException", "not found") {
		log.Printf("[WARN] DataSync Location SMB %q not found - removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DataSync Location SMB (%s): %s", d.Id(), err)
	}

	subdirectory, err := SubdirectoryFromLocationURI(aws.StringValue(output.LocationUri))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DataSync Location SMB (%s) tags: %s", d.Id(), err)
	}

	d.Set("agent_arns", flex.FlattenStringSet(output.AgentArns))

	d.Set("arn", output.LocationArn)

	d.Set("domain", output.Domain)

	if err := d.Set("mount_options", flattenSMBMountOptions(output.MountOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting mount_options: %s", err)
	}

	d.Set("subdirectory", subdirectory)
	d.Set("user", output.User)
	d.Set("uri", output.LocationUri)

	return diags
}

func resourceLocationSMBUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	if d.HasChangesExcept("tags_all", "tags") {
		input := &datasync.UpdateLocationSmbInput{
			LocationArn:  aws.String(d.Id()),
			AgentArns:    flex.ExpandStringSet(d.Get("agent_arns").(*schema.Set)),
			MountOptions: expandSMBMountOptions(d.Get("mount_options").([]interface{})),
			Password:     aws.String(d.Get("password").(string)),
			Subdirectory: aws.String(d.Get("subdirectory").(string)),
			User:         aws.String(d.Get("user").(string)),
		}

		if v, ok := d.GetOk("domain"); ok {
			input.Domain = aws.String(v.(string))
		}

		_, err := conn.UpdateLocationSmbWithContext(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating DataSync Location SMB (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceLocationSMBRead(ctx, d, meta)...)
}

func resourceLocationSMBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.DeleteLocationInput{
		LocationArn: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting DataSync Location SMB: %s", input)
	_, err := conn.DeleteLocationWithContext(ctx, input)

	if tfawserr.ErrMessageContains(err, "InvalidRequestException", "not found") {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting DataSync Location SMB (%s): %s", d.Id(), err)
	}

	return diags
}

func flattenSMBMountOptions(mountOptions *datasync.SmbMountOptions) []interface{} {
	if mountOptions == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"version": aws.StringValue(mountOptions.Version),
	}

	return []interface{}{m}
}

func expandSMBMountOptions(l []interface{}) *datasync.SmbMountOptions {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	smbMountOptions := &datasync.SmbMountOptions{
		Version: aws.String(m["version"].(string)),
	}

	return smbMountOptions
}
