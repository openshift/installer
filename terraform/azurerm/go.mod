module github.com/openshift/installer/terraform/azurerm

go 1.16

require (
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/terraform-providers/terraform-provider-azurerm v1.44.1-0.20200911233120-57b2bfc9d42c
)

replace github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.44.1-0.20210224232508-7509319df0f4 // Pin to 2.48.0-openshift
