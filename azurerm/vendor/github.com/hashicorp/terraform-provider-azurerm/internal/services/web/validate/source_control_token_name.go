package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SourceControlTokenName() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"BitBucket",
		"Dropbox",
		"GitHub",
		"OneDrive",
	}, false)
}
