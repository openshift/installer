package desktopvirtualization

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Registration - Name is the name of this Service
func (r Registration) Name() string {
	return "Desktop Virtualization"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Desktop Virtualization",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_virtual_desktop_workspace":                               resourceArmDesktopVirtualizationWorkspace(),
		"azurerm_virtual_desktop_host_pool":                               resourceVirtualDesktopHostPool(),
		"azurerm_virtual_desktop_application_group":                       resourceVirtualDesktopApplicationGroup(),
		"azurerm_virtual_desktop_workspace_application_group_association": resourceVirtualDesktopWorkspaceApplicationGroupAssociation(),
	}
}
