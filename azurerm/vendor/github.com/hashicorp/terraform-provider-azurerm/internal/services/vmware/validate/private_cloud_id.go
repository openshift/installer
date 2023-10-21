package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/privateclouds"
)

func PrivateCloudID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := privateclouds.ParsePrivateCloudID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
