package logs

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_cloudwatch_log_data_protection_policy")
func resourceDataProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceDataProtectionPolicyPut,
		ReadWithoutTimeout:   resourceDataProtectionPolicyRead,
		UpdateWithoutTimeout: resourceDataProtectionPolicyPut,
		DeleteWithoutTimeout: resourceDataProtectionPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"log_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validLogGroupName,
			},
			"policy_document": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
		},
	}
}

func resourceDataProtectionPolicyPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	logGroupName := d.Get("log_group_name").(string)

	policy, err := structure.NormalizeJsonString(d.Get("policy_document").(string))

	if err != nil {
		return diag.Errorf("policy (%s) is invalid JSON: %s", policy, err)
	}

	input := &cloudwatchlogs.PutDataProtectionPolicyInput{
		LogGroupIdentifier: aws.String(logGroupName),
		PolicyDocument:     aws.String(policy),
	}

	_, err = conn.PutDataProtectionPolicy(ctx, input)

	if err != nil {
		return diag.Errorf("putting CloudWatch Logs Data Protection Policy (%s): %s", logGroupName, err)
	}

	if d.IsNewResource() {
		d.SetId(logGroupName)
	}

	return resourceDataProtectionPolicyRead(ctx, d, meta)
}

func resourceDataProtectionPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	output, err := FindDataProtectionPolicyByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudWatch Logs Data Protection Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading CloudWatch Logs Data Protection Policy (%s): %s", d.Id(), err)
	}

	d.Set("log_group_name", output.LogGroupIdentifier)

	policyToSet, err := verify.SecondJSONUnlessEquivalent(d.Get("policy_document").(string), aws.ToString(output.PolicyDocument))

	if err != nil {
		return diag.Errorf("while setting policy (%s), encountered: %s", policyToSet, err)
	}

	policyToSet, err = structure.NormalizeJsonString(policyToSet)

	if err != nil {
		return diag.Errorf("policy (%s) is invalid JSON: %s", policyToSet, err)
	}

	d.Set("policy_document", policyToSet)

	return nil
}

func resourceDataProtectionPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	log.Printf("[DEBUG] Deleting CloudWatch Logs Data Protection Policy: %s", d.Id())
	_, err := conn.DeleteDataProtectionPolicy(ctx, &cloudwatchlogs.DeleteDataProtectionPolicyInput{
		LogGroupIdentifier: aws.String(d.Id()),
	})

	if nfe := (*types.ResourceNotFoundException)(nil); errors.As(err, &nfe) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting CloudWatch Logs Data Protection Policy (%s): %s", d.Id(), err)
	}

	return nil
}

func FindDataProtectionPolicyByID(ctx context.Context, conn *cloudwatchlogs.Client, id string) (*cloudwatchlogs.GetDataProtectionPolicyOutput, error) {
	input := &cloudwatchlogs.GetDataProtectionPolicyInput{
		LogGroupIdentifier: aws.String(id),
	}

	output, err := conn.GetDataProtectionPolicy(ctx, input)

	if nfe := (*types.ResourceNotFoundException)(nil); errors.As(err, &nfe) {
		return nil, &retry.NotFoundError{
			LastError:   nfe,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
