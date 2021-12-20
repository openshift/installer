package azurestack

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func getArmResourceNameAndGroup(s *terraform.State, name string) (string, string, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return "", "", fmt.Errorf("Not found: %s", name)
	}

	armName, hasName := rs.Primary.Attributes["name"]
	armResourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
	if !hasResourceGroup {
		return "", "", fmt.Errorf("Error: no resource group found in state for resource: %s", name)
	}

	if !hasName {
		return "", "", fmt.Errorf("Error: no name found in state for resource: %s", name)
	}

	return armName, armResourceGroup, nil
}
