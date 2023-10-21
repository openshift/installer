package cognitoidp

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_cognito_user_pool_ui_customization")
func ResourceUserPoolUICustomization() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceUserPoolUICustomizationPut,
		ReadWithoutTimeout:   resourceUserPoolUICustomizationRead,
		UpdateWithoutTimeout: resourceUserPoolUICustomizationPut,
		DeleteWithoutTimeout: resourceUserPoolUICustomizationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALL",
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"css": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"css", "image_file"},
			},

			"css_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_file": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"image_file", "css"},
			},

			"image_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_modified_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserPoolUICustomizationPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIDPConn(ctx)

	clientId := d.Get("client_id").(string)
	userPoolId := d.Get("user_pool_id").(string)

	input := &cognitoidentityprovider.SetUICustomizationInput{
		ClientId:   aws.String(clientId),
		UserPoolId: aws.String(userPoolId),
	}

	if v, ok := d.GetOk("css"); ok {
		input.CSS = aws.String(v.(string))
	}

	if v, ok := d.GetOk("image_file"); ok {
		imgFile, err := base64.StdEncoding.DecodeString(v.(string))
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "Base64 decoding image file for Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %s", userPoolId, clientId, err)
		}

		input.ImageFile = imgFile
	}

	_, err := conn.SetUICustomizationWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %s", userPoolId, clientId, err)
	}

	d.SetId(fmt.Sprintf("%s,%s", userPoolId, clientId))

	return append(diags, resourceUserPoolUICustomizationRead(ctx, d, meta)...)
}

func resourceUserPoolUICustomizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIDPConn(ctx)

	userPoolId, clientId, err := ParseUserPoolUICustomizationID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing Cognito User Pool UI customization ID (%s): %s", d.Id(), err)
	}

	uiCustomization, err := FindCognitoUserPoolUICustomization(ctx, conn, userPoolId, clientId)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, cognitoidentityprovider.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s) not found, removing from state", userPoolId, clientId)
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "getting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %s", userPoolId, clientId, err)
	}

	if uiCustomization == nil {
		if d.IsNewResource() {
			return sdkdiag.AppendErrorf(diags, "getting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): not found", userPoolId, clientId)
		}

		log.Printf("[WARN] Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s) not found, removing from state", userPoolId, clientId)
		d.SetId("")
		return diags
	}

	d.Set("client_id", uiCustomization.ClientId)
	d.Set("creation_date", aws.TimeValue(uiCustomization.CreationDate).Format(time.RFC3339))
	d.Set("css", uiCustomization.CSS)
	d.Set("css_version", uiCustomization.CSSVersion)
	d.Set("image_url", uiCustomization.ImageUrl)
	d.Set("last_modified_date", aws.TimeValue(uiCustomization.LastModifiedDate).Format(time.RFC3339))
	d.Set("user_pool_id", uiCustomization.UserPoolId)

	return diags
}

func resourceUserPoolUICustomizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIDPConn(ctx)

	userPoolId, clientId, err := ParseUserPoolUICustomizationID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing Cognito User Pool UI customization ID (%s): %s", d.Id(), err)
	}

	input := &cognitoidentityprovider.SetUICustomizationInput{
		ClientId:   aws.String(clientId),
		UserPoolId: aws.String(userPoolId),
	}

	_, err = conn.SetUICustomizationWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, cognitoidentityprovider.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Cognito User Pool UI customization (UserPoolId: %s, ClientId: %s): %s", userPoolId, clientId, err)
	}

	return diags
}

func ParseUserPoolUICustomizationID(id string) (string, string, error) {
	idParts := strings.SplitN(id, ",", 2)

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return "", "", fmt.Errorf("please make sure ID is in format USER_POOL_ID,CLIENT_ID")
	}

	return idParts[0], idParts[1], nil
}
