package lambda

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

// @SDKResource("aws_lambda_provisioned_concurrency_config")
func ResourceProvisionedConcurrencyConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceProvisionedConcurrencyConfigCreate,
		ReadWithoutTimeout:   resourceProvisionedConcurrencyConfigRead,
		UpdateWithoutTimeout: resourceProvisionedConcurrencyConfigUpdate,
		DeleteWithoutTimeout: resourceProvisionedConcurrencyConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceProvisionedConcurrencyConfigV0().CoreConfigSchema().ImpliedType(),
				Upgrade: provisionedConcurrencyConfigStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"provisioned_concurrent_executions": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"qualifier": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"skip_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

const (
	ProvisionedConcurrencyIDPartCount = 2
)

func resourceProvisionedConcurrencyConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).LambdaConn(ctx)
	functionName := d.Get("function_name").(string)
	qualifier := d.Get("qualifier").(string)

	input := &lambda.PutProvisionedConcurrencyConfigInput{
		FunctionName:                    aws.String(functionName),
		ProvisionedConcurrentExecutions: aws.Int64(int64(d.Get("provisioned_concurrent_executions").(int))),
		Qualifier:                       aws.String(qualifier),
	}

	_, err := conn.PutProvisionedConcurrencyConfigWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "putting Lambda Provisioned Concurrency Config (%s,%s): %s", functionName, qualifier, err)
	}

	parts := []string{functionName, qualifier}
	id, err := flex.FlattenResourceId(parts, ProvisionedConcurrencyIDPartCount, false)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting Lambda Provisioned Concurrency Config ID (%s,%s): %s", functionName, qualifier, err)
	}
	d.SetId(id)

	if err := waitForProvisionedConcurrencyConfigStatusReady(ctx, conn, functionName, qualifier, d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Lambda Provisioned Concurrency Config (%s) to be ready: %s", d.Id(), err)
	}

	return append(diags, resourceProvisionedConcurrencyConfigRead(ctx, d, meta)...)
}

func resourceProvisionedConcurrencyConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).LambdaConn(ctx)

	parts, err := flex.ExpandResourceId(d.Id(), ProvisionedConcurrencyIDPartCount, false)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Lambda Provisioned Concurrency Config (%s): %s", d.Id(), err)
	}
	functionName := parts[0]
	qualifier := parts[1]

	input := &lambda.GetProvisionedConcurrencyConfigInput{
		FunctionName: aws.String(functionName),
		Qualifier:    aws.String(qualifier),
	}

	output, err := conn.GetProvisionedConcurrencyConfigWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, lambda.ErrCodeProvisionedConcurrencyConfigNotFoundException) || tfawserr.ErrCodeEquals(err, lambda.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Lambda Provisioned Concurrency Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Lambda Provisioned Concurrency Config (%s): %s", d.Id(), err)
	}

	d.Set("function_name", functionName)
	d.Set("provisioned_concurrent_executions", output.AllocatedProvisionedConcurrentExecutions)
	d.Set("qualifier", qualifier)

	return diags
}

func resourceProvisionedConcurrencyConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).LambdaConn(ctx)

	parts, err := flex.ExpandResourceId(d.Id(), ProvisionedConcurrencyIDPartCount, false)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Lambda Provisioned Concurrency Config (%s): %s", d.Id(), err)
	}
	functionName := parts[0]
	qualifier := parts[1]

	input := &lambda.PutProvisionedConcurrencyConfigInput{
		FunctionName:                    aws.String(functionName),
		ProvisionedConcurrentExecutions: aws.Int64(int64(d.Get("provisioned_concurrent_executions").(int))),
		Qualifier:                       aws.String(qualifier),
	}

	_, err = conn.PutProvisionedConcurrencyConfigWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Lambda Provisioned Concurrency Config (%s): %s", d.Id(), err)
	}

	if err := waitForProvisionedConcurrencyConfigStatusReady(ctx, conn, functionName, qualifier, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Lambda Provisioned Concurrency Config (%s): waiting for completion: %s", d.Id(), err)
	}

	return append(diags, resourceProvisionedConcurrencyConfigRead(ctx, d, meta)...)
}

func resourceProvisionedConcurrencyConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if v, ok := d.GetOk("skip_destroy"); ok && v.(bool) {
		log.Printf("[DEBUG] Retaining Lambda Provisioned Concurrency Config %q", d.Id())
		return diags
	}

	conn := meta.(*conns.AWSClient).LambdaConn(ctx)

	parts, err := flex.ExpandResourceId(d.Id(), ProvisionedConcurrencyIDPartCount, false)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Lambda Provisioned Concurrency Config (%s): %s", d.Id(), err)
	}

	input := &lambda.DeleteProvisionedConcurrencyConfigInput{
		FunctionName: aws.String(parts[0]),
		Qualifier:    aws.String(parts[1]),
	}

	_, err = conn.DeleteProvisionedConcurrencyConfigWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, lambda.ErrCodeProvisionedConcurrencyConfigNotFoundException) || tfawserr.ErrCodeEquals(err, lambda.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Lambda Provisioned Concurrency Config (%s): %s", d.Id(), err)
	}

	return diags
}

func refreshProvisionedConcurrencyConfigStatus(ctx context.Context, conn *lambda.Lambda, functionName, qualifier string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &lambda.GetProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(functionName),
			Qualifier:    aws.String(qualifier),
		}

		output, err := conn.GetProvisionedConcurrencyConfigWithContext(ctx, input)

		if err != nil {
			return "", "", err
		}

		status := aws.StringValue(output.Status)

		if status == lambda.ProvisionedConcurrencyStatusEnumFailed {
			return output, status, fmt.Errorf("status reason: %s", aws.StringValue(output.StatusReason))
		}

		return output, status, nil
	}
}

func waitForProvisionedConcurrencyConfigStatusReady(ctx context.Context, conn *lambda.Lambda, functionName, qualifier string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{lambda.ProvisionedConcurrencyStatusEnumInProgress},
		Target:  []string{lambda.ProvisionedConcurrencyStatusEnumReady},
		Refresh: refreshProvisionedConcurrencyConfigStatus(ctx, conn, functionName, qualifier),
		Timeout: timeout,
		Delay:   5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
