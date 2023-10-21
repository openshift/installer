package glue

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_glue_user_defined_function")
func ResourceUserDefinedFunction() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceUserDefinedFunctionCreate,
		ReadWithoutTimeout:   resourceUserDefinedFunctionRead,
		UpdateWithoutTimeout: resourceUserDefinedFunctionUpdate,
		DeleteWithoutTimeout: resourceUserDefinedFunctionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"catalog_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"class_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"owner_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"owner_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(glue.PrincipalType_Values(), false),
			},
			"resource_uris": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(glue.ResourceType_Values(), false),
						},
						"uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserDefinedFunctionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)
	catalogID := createCatalogID(d, meta.(*conns.AWSClient).AccountID)
	dbName := d.Get("database_name").(string)
	funcName := d.Get("name").(string)

	input := &glue.CreateUserDefinedFunctionInput{
		CatalogId:     aws.String(catalogID),
		DatabaseName:  aws.String(dbName),
		FunctionInput: expandUserDefinedFunctionInput(d),
	}

	_, err := conn.CreateUserDefinedFunctionWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Glue User Defined Function: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", catalogID, dbName, funcName))

	return append(diags, resourceUserDefinedFunctionRead(ctx, d, meta)...)
}

func resourceUserDefinedFunctionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	catalogID, dbName, funcName, err := ReadUDFID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Glue User Defined Function (%s): %s", d.Id(), err)
	}

	input := &glue.UpdateUserDefinedFunctionInput{
		CatalogId:     aws.String(catalogID),
		DatabaseName:  aws.String(dbName),
		FunctionName:  aws.String(funcName),
		FunctionInput: expandUserDefinedFunctionInput(d),
	}

	if _, err := conn.UpdateUserDefinedFunctionWithContext(ctx, input); err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Glue User Defined Function (%s): %s", d.Id(), err)
	}

	return append(diags, resourceUserDefinedFunctionRead(ctx, d, meta)...)
}

func resourceUserDefinedFunctionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	catalogID, dbName, funcName, err := ReadUDFID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Glue User Defined Function (%s): %s", d.Id(), err)
	}

	input := &glue.GetUserDefinedFunctionInput{
		CatalogId:    aws.String(catalogID),
		DatabaseName: aws.String(dbName),
		FunctionName: aws.String(funcName),
	}

	out, err := conn.GetUserDefinedFunctionWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			log.Printf("[WARN] Glue User Defined Function (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}

		return sdkdiag.AppendErrorf(diags, "reading Glue User Defined Function (%s): %s", d.Id(), err)
	}

	udf := out.UserDefinedFunction

	udfArn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "glue",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("userDefinedFunction/%s/%s", dbName, aws.StringValue(udf.FunctionName)),
	}.String()

	d.Set("arn", udfArn)
	d.Set("name", udf.FunctionName)
	d.Set("catalog_id", catalogID)
	d.Set("database_name", dbName)
	d.Set("owner_type", udf.OwnerType)
	d.Set("owner_name", udf.OwnerName)
	d.Set("class_name", udf.ClassName)
	if udf.CreateTime != nil {
		d.Set("create_time", udf.CreateTime.Format(time.RFC3339))
	}
	if err := d.Set("resource_uris", flattenUserDefinedFunctionResourceURI(udf.ResourceUris)); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Glue User Defined Function (%s): setting resource_uris: %s", d.Id(), err)
	}

	return diags
}

func resourceUserDefinedFunctionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)
	catalogID, dbName, funcName, err := ReadUDFID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Glue User Defined Function (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Glue User Defined Function: %s", d.Id())
	_, err = conn.DeleteUserDefinedFunctionWithContext(ctx, &glue.DeleteUserDefinedFunctionInput{
		CatalogId:    aws.String(catalogID),
		DatabaseName: aws.String(dbName),
		FunctionName: aws.String(funcName),
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Glue User Defined Function (%s): %s", d.Id(), err)
	}
	return diags
}

func ReadUDFID(id string) (catalogID string, dbName string, funcName string, err error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 3 {
		return "", "", "", fmt.Errorf("unexpected format of ID (%q), expected CATALOG-ID:DATABASE-NAME:FUNCTION-NAME", id)
	}
	return idParts[0], idParts[1], idParts[2], nil
}

func expandUserDefinedFunctionInput(d *schema.ResourceData) *glue.UserDefinedFunctionInput {
	udf := &glue.UserDefinedFunctionInput{
		ClassName:    aws.String(d.Get("class_name").(string)),
		FunctionName: aws.String(d.Get("name").(string)),
		OwnerName:    aws.String(d.Get("owner_name").(string)),
		OwnerType:    aws.String(d.Get("owner_type").(string)),
	}

	if v, ok := d.GetOk("resource_uris"); ok && v.(*schema.Set).Len() > 0 {
		udf.ResourceUris = expandUserDefinedFunctionResourceURI(d.Get("resource_uris").(*schema.Set))
	}

	return udf
}

func expandUserDefinedFunctionResourceURI(conf *schema.Set) []*glue.ResourceUri {
	result := make([]*glue.ResourceUri, 0, conf.Len())

	for _, r := range conf.List() {
		uriRaw, ok := r.(map[string]interface{})

		if !ok {
			continue
		}

		uri := &glue.ResourceUri{
			ResourceType: aws.String(uriRaw["resource_type"].(string)),
			Uri:          aws.String(uriRaw["uri"].(string)),
		}

		result = append(result, uri)
	}

	return result
}

func flattenUserDefinedFunctionResourceURI(uris []*glue.ResourceUri) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(uris))

	for _, i := range uris {
		l := map[string]interface{}{
			"resource_type": aws.StringValue(i.ResourceType),
			"uri":           aws.StringValue(i.Uri),
		}

		result = append(result, l)
	}
	return result
}
