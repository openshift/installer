package azurestack

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGroupNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validateArmResourceGroupName,
	}
}

func resourceGroupNameDiffSuppressSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: resourceAzureStackResourceGroupNameDiffSuppress,
		ValidateFunc:     validateArmResourceGroupName,
	}
}

func resourceGroupNameForDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
}

func validateArmResourceGroupName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if len(value) > 80 {
		es = append(es, fmt.Errorf("%q may not exceed 80 characters in length", k))
	}

	if strings.HasSuffix(value, ".") {
		es = append(es, fmt.Errorf("%q may not end with a period", k))
	}

	// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
	if matched := regexp.MustCompile(`^[-\w\._\(\)]+$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters, dash, underscores, parentheses and periods", k))
	}

	return
}
