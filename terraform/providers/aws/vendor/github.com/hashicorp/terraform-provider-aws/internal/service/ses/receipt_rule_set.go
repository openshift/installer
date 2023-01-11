package ses

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceReceiptRuleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceReceiptRuleSetCreate,
		Read:   resourceReceiptRuleSetRead,
		Delete: resourceReceiptRuleSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_set_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceReceiptRuleSetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SESConn

	ruleSetName := d.Get("rule_set_name").(string)

	createOpts := &ses.CreateReceiptRuleSetInput{
		RuleSetName: aws.String(ruleSetName),
	}

	_, err := conn.CreateReceiptRuleSet(createOpts)
	if err != nil {
		return fmt.Errorf("creating SES rule set: %w", err)
	}

	d.SetId(ruleSetName)

	return resourceReceiptRuleSetRead(d, meta)
}

func resourceReceiptRuleSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SESConn

	input := &ses.DescribeReceiptRuleSetInput{
		RuleSetName: aws.String(d.Id()),
	}

	resp, err := conn.DescribeReceiptRuleSet(input)

	if tfawserr.ErrCodeEquals(err, ses.ErrCodeRuleSetDoesNotExistException) {
		log.Printf("[WARN] SES Receipt Rule Set (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("describing SES Receipt Rule Set (%s): %w", d.Id(), err)
	}

	if resp.Metadata == nil {
		log.Print("[WARN] No Receipt Rule Set found")
		d.SetId("")
		return nil
	}

	name := aws.StringValue(resp.Metadata.Name)
	d.Set("rule_set_name", name)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "ses",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("receipt-rule-set/%s", name),
	}.String()
	d.Set("arn", arn)

	return nil
}

func resourceReceiptRuleSetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SESConn

	log.Printf("[DEBUG] SES Delete Receipt Rule Set: %s", d.Id())
	input := &ses.DeleteReceiptRuleSetInput{
		RuleSetName: aws.String(d.Id()),
	}
	if _, err := conn.DeleteReceiptRuleSet(input); err != nil {
		return fmt.Errorf("deleting SES Receipt Rule Set (%s): %w", d.Id(), err)
	}

	return nil
}
