package srv

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func schemaResourceGroupName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validateResourceGroupName,
	}
}
func validateResourceGroupName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 80 {
		errors = append(errors, fmt.Errorf("%q may not exceed 80 characters in length", k))
	}

	if strings.HasSuffix(value, ".") {
		errors = append(errors, fmt.Errorf("%q may not end with a period", k))
	}

	// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
	if matched := regexp.MustCompile(`^[-\w\._\(\)]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dash, underscores, parentheses and periods", k))
	}

	return warnings, errors
}
