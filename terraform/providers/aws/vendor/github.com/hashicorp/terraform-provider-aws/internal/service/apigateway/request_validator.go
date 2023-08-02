package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_api_gateway_request_validator")
func ResourceRequestValidator() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRequestValidatorCreate,
		ReadWithoutTimeout:   resourceRequestValidatorRead,
		UpdateWithoutTimeout: resourceRequestValidatorUpdate,
		DeleteWithoutTimeout: resourceRequestValidatorDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected REST-API-ID/REQUEST-VALIDATOR-ID", d.Id())
				}
				restApiID := idParts[0]
				requestValidatorID := idParts[1]
				d.Set("rest_api_id", restApiID)
				d.SetId(requestValidatorID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rest_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"validate_request_body": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"validate_request_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceRequestValidatorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	name := d.Get("name").(string)
	input := &apigateway.CreateRequestValidatorInput{
		Name:                      aws.String(name),
		RestApiId:                 aws.String(d.Get("rest_api_id").(string)),
		ValidateRequestBody:       aws.Bool(d.Get("validate_request_body").(bool)),
		ValidateRequestParameters: aws.Bool(d.Get("validate_request_parameters").(bool)),
	}

	output, err := conn.CreateRequestValidatorWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating API Gateway Request Validator (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(output.Id))

	return diags
}

func resourceRequestValidatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	output, err := FindRequestValidatorByTwoPartKey(ctx, conn, d.Id(), d.Get("rest_api_id").(string))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] API Gateway Request Validator (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading API Gateway Request Validator (%s): %s", d.Id(), err)
	}

	d.Set("name", output.Name)
	d.Set("validate_request_body", output.ValidateRequestBody)
	d.Set("validate_request_parameters", output.ValidateRequestParameters)

	return diags
}

func resourceRequestValidatorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	operations := make([]*apigateway.PatchOperation, 0)

	if d.HasChange("name") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/name"),
			Value: aws.String(d.Get("name").(string)),
		})
	}

	if d.HasChange("validate_request_body") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/validateRequestBody"),
			Value: aws.String(fmt.Sprintf("%t", d.Get("validate_request_body").(bool))),
		})
	}

	if d.HasChange("validate_request_parameters") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/validateRequestParameters"),
			Value: aws.String(fmt.Sprintf("%t", d.Get("validate_request_parameters").(bool))),
		})
	}

	input := &apigateway.UpdateRequestValidatorInput{
		RequestValidatorId: aws.String(d.Id()),
		RestApiId:          aws.String(d.Get("rest_api_id").(string)),
		PatchOperations:    operations,
	}

	_, err := conn.UpdateRequestValidatorWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating API Gateway Request Validator (%s): %s", d.Id(), err)
	}

	return append(diags, resourceRequestValidatorRead(ctx, d, meta)...)
}

func resourceRequestValidatorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	log.Printf("[DEBUG] Deleting API Gateway Request Validator: %s", d.Id())
	_, err := conn.DeleteRequestValidatorWithContext(ctx, &apigateway.DeleteRequestValidatorInput{
		RequestValidatorId: aws.String(d.Id()),
		RestApiId:          aws.String(d.Get("rest_api_id").(string)),
	})

	if err != nil {
		// XXX: Figure out a way to delete the method that depends on the request validator first
		// otherwise the validator will be dangling until the API is deleted
		if !strings.Contains(err.Error(), apigateway.ErrCodeConflictException) {
			return sdkdiag.AppendErrorf(diags, "deleting API Gateway Request Validator (%s): %s", d.Id(), err)
		}
	}

	return diags
}

func FindRequestValidatorByTwoPartKey(ctx context.Context, conn *apigateway.APIGateway, requestValidatorID, apiID string) (*apigateway.UpdateRequestValidatorOutput, error) {
	input := &apigateway.GetRequestValidatorInput{
		RequestValidatorId: aws.String(requestValidatorID),
		RestApiId:          aws.String(apiID),
	}

	output, err := conn.GetRequestValidatorWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
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
