package redshift

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_redshift_hsm_configuration", name="HSM Configuration")
// @Tags(identifierAttribute="arn")
func ResourceHSMConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceHSMConfigurationCreate,
		ReadWithoutTimeout:   resourceHSMConfigurationRead,
		UpdateWithoutTimeout: resourceHSMConfigurationUpdate,
		DeleteWithoutTimeout: resourceHSMConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hsm_configuration_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hsm_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hsm_partition_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hsm_partition_password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"hsm_server_public_certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceHSMConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	hsmConfigurationID := d.Get("hsm_configuration_identifier").(string)
	input := &redshift.CreateHsmConfigurationInput{
		Description:                aws.String(d.Get("description").(string)),
		HsmConfigurationIdentifier: aws.String(hsmConfigurationID),
		HsmIpAddress:               aws.String(d.Get("hsm_ip_address").(string)),
		HsmPartitionName:           aws.String(d.Get("hsm_partition_name").(string)),
		HsmPartitionPassword:       aws.String(d.Get("hsm_partition_password").(string)),
		HsmServerPublicCertificate: aws.String(d.Get("hsm_server_public_certificate").(string)),
		Tags:                       GetTagsIn(ctx),
	}

	output, err := conn.CreateHsmConfigurationWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Redshift HSM Configuration (%s): %s", hsmConfigurationID, err)
	}

	d.SetId(aws.StringValue(output.HsmConfiguration.HsmConfigurationIdentifier))

	return append(diags, resourceHSMConfigurationRead(ctx, d, meta)...)
}

func resourceHSMConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	hsmConfiguration, err := FindHSMConfigurationByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Redshift HSM Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Redshift HSM Configuration (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   redshift.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("hsmconfiguration:%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("hsm_configuration_identifier", hsmConfiguration.HsmConfigurationIdentifier)
	d.Set("hsm_ip_address", hsmConfiguration.HsmIpAddress)
	d.Set("hsm_partition_name", hsmConfiguration.HsmPartitionName)
	d.Set("description", hsmConfiguration.Description)
	d.Set("hsm_partition_password", d.Get("hsm_partition_password").(string))
	d.Set("hsm_server_public_certificate", d.Get("hsm_server_public_certificate").(string))

	SetTagsOut(ctx, hsmConfiguration.Tags)

	return diags
}

func resourceHSMConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Tags only.

	return append(diags, resourceHSMConfigurationRead(ctx, d, meta)...)
}

func resourceHSMConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	log.Printf("[DEBUG] Deleting Redshift HSM Configuration: %s", d.Id())
	_, err := conn.DeleteHsmConfigurationWithContext(ctx, &redshift.DeleteHsmConfigurationInput{
		HsmConfigurationIdentifier: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, redshift.ErrCodeHsmConfigurationNotFoundFault) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Redshift HSM Configuration (%s): %s", d.Id(), err)
	}

	return diags
}
