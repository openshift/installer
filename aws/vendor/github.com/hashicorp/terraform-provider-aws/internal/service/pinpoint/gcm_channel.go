package pinpoint

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_pinpoint_gcm_channel")
func ResourceGCMChannel() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceGCMChannelUpsert,
		ReadWithoutTimeout:   resourceGCMChannelRead,
		UpdateWithoutTimeout: resourceGCMChannelUpsert,
		DeleteWithoutTimeout: resourceGCMChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGCMChannelUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).PinpointConn(ctx)

	applicationId := d.Get("application_id").(string)

	params := &pinpoint.GCMChannelRequest{}

	params.ApiKey = aws.String(d.Get("api_key").(string))
	params.Enabled = aws.Bool(d.Get("enabled").(bool))

	req := pinpoint.UpdateGcmChannelInput{
		ApplicationId:     aws.String(applicationId),
		GCMChannelRequest: params,
	}

	_, err := conn.UpdateGcmChannelWithContext(ctx, &req)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "putting Pinpoint GCM Channel for application %s: %s", applicationId, err)
	}

	d.SetId(applicationId)

	return append(diags, resourceGCMChannelRead(ctx, d, meta)...)
}

func resourceGCMChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).PinpointConn(ctx)

	log.Printf("[INFO] Reading Pinpoint GCM Channel for application %s", d.Id())

	output, err := conn.GetGcmChannelWithContext(ctx, &pinpoint.GetGcmChannelInput{
		ApplicationId: aws.String(d.Id()),
	})
	if err != nil {
		if tfawserr.ErrCodeEquals(err, pinpoint.ErrCodeNotFoundException) {
			log.Printf("[WARN] Pinpoint GCM Channel for application %s not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}

		return sdkdiag.AppendErrorf(diags, "getting Pinpoint GCM Channel for application %s: %s", d.Id(), err)
	}

	d.Set("application_id", output.GCMChannelResponse.ApplicationId)
	d.Set("enabled", output.GCMChannelResponse.Enabled)
	// api_key is never returned

	return diags
}

func resourceGCMChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).PinpointConn(ctx)

	log.Printf("[DEBUG] Deleting Pinpoint GCM Channel for application %s", d.Id())
	_, err := conn.DeleteGcmChannelWithContext(ctx, &pinpoint.DeleteGcmChannelInput{
		ApplicationId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, pinpoint.ErrCodeNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Pinpoint GCM Channel for application %s: %s", d.Id(), err)
	}
	return diags
}
