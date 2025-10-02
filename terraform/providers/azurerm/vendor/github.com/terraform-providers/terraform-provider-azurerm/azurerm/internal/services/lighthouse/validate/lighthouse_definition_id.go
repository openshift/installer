package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/lighthouse/parse"
)

func LighthouseDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.LighthouseDefinitionID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a resource id: %v", k, err))
		return
	}

	return
}
