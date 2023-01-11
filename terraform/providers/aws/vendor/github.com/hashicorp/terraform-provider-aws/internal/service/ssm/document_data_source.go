package ssm

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceDocument() *schema.Resource {
	return &schema.Resource{
		Read: dataDocumentRead,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"document_format": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ssm.DocumentFormatJson,
				ValidateFunc: validation.StringInSlice(ssm.DocumentFormat_Values(), false),
			},
			"document_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"document_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataDocumentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SSMConn

	docInput := &ssm.GetDocumentInput{
		Name:           aws.String(d.Get("name").(string)),
		DocumentFormat: aws.String(d.Get("document_format").(string)),
	}

	if docVersion, ok := d.GetOk("document_version"); ok {
		docInput.DocumentVersion = aws.String(docVersion.(string))
	}

	log.Printf("[DEBUG] Reading SSM Document: %s", docInput)
	resp, err := conn.GetDocument(docInput)

	if err != nil {
		return fmt.Errorf("Error reading SSM Document: %w", err)
	}

	name := aws.StringValue(resp.Name)

	d.SetId(name)

	if !strings.HasPrefix(name, "AWS-") {
		arn := arn.ARN{
			Partition: meta.(*conns.AWSClient).Partition,
			Service:   "ssm",
			Region:    meta.(*conns.AWSClient).Region,
			AccountID: meta.(*conns.AWSClient).AccountID,
			Resource:  fmt.Sprintf("document/%s", name),
		}.String()
		d.Set("arn", arn)
	} else {
		d.Set("arn", name)
	}

	d.Set("name", name)
	d.Set("content", resp.Content)
	d.Set("document_version", resp.DocumentVersion)
	d.Set("document_format", resp.DocumentFormat)
	d.Set("document_type", resp.DocumentType)

	return nil
}
