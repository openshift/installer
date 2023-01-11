package inspector

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourceGroupCreate,
		Read:   resourceResourceGroupRead,
		Delete: resourceResourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"tags": {
				ForceNew: true,
				Required: true,
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn

	req := &inspector.CreateResourceGroupInput{
		ResourceGroupTags: expandResourceGroupTags(d.Get("tags").(map[string]interface{})),
	}
	log.Printf("[DEBUG] Creating Inspector resource group: %#v", req)
	resp, err := conn.CreateResourceGroup(req)

	if err != nil {
		return fmt.Errorf("error creating Inspector resource group: %s", err)
	}

	d.SetId(aws.StringValue(resp.ResourceGroupArn))

	return resourceResourceGroupRead(d, meta)
}

func resourceResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn

	resp, err := conn.DescribeResourceGroups(&inspector.DescribeResourceGroupsInput{
		ResourceGroupArns: aws.StringSlice([]string{d.Id()}),
	})

	if err != nil {
		return fmt.Errorf("error reading Inspector resource group (%s): %s", d.Id(), err)
	}

	if len(resp.ResourceGroups) == 0 {
		if failedItem, ok := resp.FailedItems[d.Id()]; ok {
			failureCode := aws.StringValue(failedItem.FailureCode)
			if failureCode == inspector.FailedItemErrorCodeItemDoesNotExist {
				log.Printf("[WARN] Inspector resource group (%s) not found, removing from state", d.Id())
				d.SetId("")
				return nil
			}

			return fmt.Errorf("error reading Inspector resource group (%s): %s", d.Id(), failureCode)
		}

		return fmt.Errorf("error reading Inspector resource group (%s): %v", d.Id(), resp.FailedItems)
	}

	resourceGroup := resp.ResourceGroups[0]
	d.Set("arn", resourceGroup.Arn)

	//lintignore:AWSR002
	if err := d.Set("tags", flattenResourceGroupTags(resourceGroup.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func expandResourceGroupTags(m map[string]interface{}) []*inspector.ResourceGroupTag {
	var result []*inspector.ResourceGroupTag

	for k, v := range m {
		result = append(result, &inspector.ResourceGroupTag{
			Key:   aws.String(k),
			Value: aws.String(v.(string)),
		})
	}

	return result
}

func flattenResourceGroupTags(tags []*inspector.ResourceGroupTag) map[string]interface{} {
	m := map[string]interface{}{}

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}

	return m
}
