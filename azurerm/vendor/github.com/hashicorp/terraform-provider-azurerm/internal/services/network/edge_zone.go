package network

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func expandEdgeZone(input string) *network.ExtendedLocation {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &network.ExtendedLocation{
		Name: utils.String(normalized),
		Type: network.ExtendedLocationTypesEdgeZone,
	}
}

func flattenEdgeZone(input *network.ExtendedLocation) string {
	if input == nil || input.Type != network.ExtendedLocationTypesEdgeZone || input.Name == nil {
		return ""
	}
	return edgezones.NormalizeNilable(input.Name)
}
