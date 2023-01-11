package wafv2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceWebACL() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWebACLRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(wafv2.Scope_Values(), false),
			},
		},
	}
}

func dataSourceWebACLRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFV2Conn
	name := d.Get("name").(string)

	var foundWebACL *wafv2.WebACLSummary
	input := &wafv2.ListWebACLsInput{
		Scope: aws.String(d.Get("scope").(string)),
		Limit: aws.Int64(100),
	}

	for {
		resp, err := conn.ListWebACLs(input)
		if err != nil {
			return fmt.Errorf("Error reading WAFv2 WebACLs: %w", err)
		}

		if resp == nil || resp.WebACLs == nil {
			return fmt.Errorf("Error reading WAFv2 WebACLs")
		}

		for _, webACL := range resp.WebACLs {
			if aws.StringValue(webACL.Name) == name {
				foundWebACL = webACL
				break
			}
		}

		if resp.NextMarker == nil || foundWebACL != nil {
			break
		}
		input.NextMarker = resp.NextMarker
	}

	if foundWebACL == nil {
		return fmt.Errorf("WAFv2 WebACL not found for name: %s", name)
	}

	d.SetId(aws.StringValue(foundWebACL.Id))
	d.Set("arn", foundWebACL.ARN)
	d.Set("description", foundWebACL.Description)

	return nil
}
