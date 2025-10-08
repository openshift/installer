package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var _ resourceids.Id = AppConfigurationFeatureId{}

type AppConfigurationFeatureId struct {
	ConfigurationStoreId string
	Name                 string
	Label                string
}

func (id AppConfigurationFeatureId) ID() string {
	return fmt.Sprintf("%s/AppConfigurationFeature/%s/Label/%s", id.ConfigurationStoreId, id.Name, id.Label)
}

func (id AppConfigurationFeatureId) String() string {
	components := []string{
		fmt.Sprintf("Configuration Store Id %q", id.ConfigurationStoreId),
		fmt.Sprintf("Label %q", id.Label),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Feature: %s", strings.Join(components, " / "))
}

func FeatureId(input string) (*AppConfigurationFeatureId, error) {
	resourceID, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	keyName := resourceID.Path["AppConfigurationFeature"]
	label := resourceID.Path["Label"]

	appcfgID := AppConfigurationFeatureId{
		Name:  keyName,
		Label: label,
	}

	// Golang's URL parser will translate %00 to \000 (NUL). This will only happen if we're dealing with an empty
	// label, so we set the label to the expected value (empty string) and trim the input string, so we can properly
	// extract the configuration store ID out of it.
	if label == "\000" {
		appcfgID.Label = ""
		input = strings.TrimSuffix(input, "%00")
	}
	appcfgID.ConfigurationStoreId = strings.TrimSuffix(input, fmt.Sprintf("/AppConfigurationFeature/%s/Label/%s", appcfgID.Name, appcfgID.Label))

	return &appcfgID, nil
}
