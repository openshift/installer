package wafv2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func DataSourceIPSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPSetRead,

		Schema: map[string]*schema.Schema{
			"addresses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address_version": {
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

func dataSourceIPSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).WAFV2Conn
	name := d.Get("name").(string)

	var foundIpSet *wafv2.IPSetSummary
	input := &wafv2.ListIPSetsInput{
		Scope: aws.String(d.Get("scope").(string)),
		Limit: aws.Int64(100),
	}

	for {
		resp, err := conn.ListIPSets(input)
		if err != nil {
			return fmt.Errorf("Error reading WAFv2 IPSets: %w", err)
		}

		if resp == nil || resp.IPSets == nil {
			return fmt.Errorf("Error reading WAFv2 IPSets")
		}

		for _, ipSet := range resp.IPSets {
			if ipSet != nil && aws.StringValue(ipSet.Name) == name {
				foundIpSet = ipSet
				break
			}
		}

		if resp.NextMarker == nil || foundIpSet != nil {
			break
		}
		input.NextMarker = resp.NextMarker
	}

	if foundIpSet == nil {
		return fmt.Errorf("WAFv2 IPSet not found for name: %s", name)
	}

	resp, err := conn.GetIPSet(&wafv2.GetIPSetInput{
		Id:    foundIpSet.Id,
		Name:  foundIpSet.Name,
		Scope: aws.String(d.Get("scope").(string)),
	})

	if err != nil {
		return fmt.Errorf("Error reading WAFv2 IPSet: %w", err)
	}

	if resp == nil || resp.IPSet == nil {
		return fmt.Errorf("Error reading WAFv2 IPSet")
	}

	d.SetId(aws.StringValue(resp.IPSet.Id))
	d.Set("arn", resp.IPSet.ARN)
	d.Set("description", resp.IPSet.Description)
	d.Set("ip_address_version", resp.IPSet.IPAddressVersion)

	if err := d.Set("addresses", flex.FlattenStringList(resp.IPSet.Addresses)); err != nil {
		return fmt.Errorf("error setting addresses: %w", err)
	}

	return nil
}
