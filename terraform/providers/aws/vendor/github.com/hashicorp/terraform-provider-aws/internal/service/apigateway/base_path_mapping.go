package apigateway

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const EmptyBasePathMappingValue = "(none)"

func ResourceBasePathMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceBasePathMappingCreate,
		Read:   resourceBasePathMappingRead,
		Update: resourceBasePathMappingUpdate,
		Delete: resourceBasePathMappingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stage_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBasePathMappingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn
	input := &apigateway.CreateBasePathMappingInput{
		RestApiId:  aws.String(d.Get("api_id").(string)),
		DomainName: aws.String(d.Get("domain_name").(string)),
		BasePath:   aws.String(d.Get("base_path").(string)),
		Stage:      aws.String(d.Get("stage_name").(string)),
	}

	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err := conn.CreateBasePathMapping(input)

		if err != nil {
			if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeBadRequestException) {
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(
				fmt.Errorf("Error creating Gateway base path mapping: %s", err),
			)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateBasePathMapping(input)
	}

	if err != nil {
		return fmt.Errorf("Error creating Gateway base path mapping: %s", err)
	}

	id := fmt.Sprintf("%s/%s", d.Get("domain_name").(string), d.Get("base_path").(string))
	d.SetId(id)

	return resourceBasePathMappingRead(d, meta)
}

func resourceBasePathMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	operations := make([]*apigateway.PatchOperation, 0)

	if d.HasChange("stage_name") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/stage"),
			Value: aws.String(d.Get("stage_name").(string)),
		})
	}

	if d.HasChange("api_id") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/restapiId"),
			Value: aws.String(d.Get("api_id").(string)),
		})
	}

	if d.HasChange("base_path") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String("replace"),
			Path:  aws.String("/basePath"),
			Value: aws.String(d.Get("base_path").(string)),
		})
	}

	domainName, basePath, decodeErr := DecodeBasePathMappingID(d.Id())
	if decodeErr != nil {
		return decodeErr
	}

	input := apigateway.UpdateBasePathMappingInput{
		BasePath:        aws.String(basePath),
		DomainName:      aws.String(domainName),
		PatchOperations: operations,
	}

	log.Printf("[INFO] Updating API Gateway base path mapping: %s", input)

	_, err := conn.UpdateBasePathMapping(&input)

	if err != nil {
		return fmt.Errorf("Updating API Gateway base path mapping failed: %w", err)
	}

	if d.HasChange("base_path") {
		id := fmt.Sprintf("%s/%s", d.Get("domain_name").(string), d.Get("base_path").(string))
		d.SetId(id)
	}

	log.Printf("[DEBUG] API Gateway base path mapping updated: %s", d.Id())

	return resourceBasePathMappingRead(d, meta)
}

func resourceBasePathMappingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	domainName, basePath, err := DecodeBasePathMappingID(d.Id())
	if err != nil {
		return err
	}

	mapping, err := conn.GetBasePathMapping(&apigateway.GetBasePathMappingInput{
		DomainName: aws.String(domainName),
		BasePath:   aws.String(basePath),
	})
	if err != nil {
		if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
			log.Printf("[WARN] API Gateway Base Path Mapping (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error reading API Gateway Base Path Mapping (%s): %w", d.Id(), err)
	}

	mappingBasePath := aws.StringValue(mapping.BasePath)

	if mappingBasePath == EmptyBasePathMappingValue {
		mappingBasePath = ""
	}

	d.Set("base_path", mappingBasePath)
	d.Set("domain_name", domainName)
	d.Set("api_id", mapping.RestApiId)
	d.Set("stage_name", mapping.Stage)

	return nil
}

func resourceBasePathMappingDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayConn

	domainName, basePath, err := DecodeBasePathMappingID(d.Id())
	if err != nil {
		return err
	}

	_, err = conn.DeleteBasePathMapping(&apigateway.DeleteBasePathMappingInput{
		DomainName: aws.String(domainName),
		BasePath:   aws.String(basePath),
	})

	if err != nil {
		if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
			return nil
		}

		return err
	}

	return nil
}

func DecodeBasePathMappingID(id string) (string, string, error) {
	idFormatErr := fmt.Errorf("Unexpected format of ID (%q), expected DOMAIN/BASEPATH", id)

	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", "", idFormatErr
	}

	domainName := parts[0]
	basePath := parts[1]

	if domainName == "" {
		return "", "", idFormatErr
	}

	if basePath == "" {
		basePath = EmptyBasePathMappingValue
	}

	return domainName, basePath, nil
}
