module github.com/openshift/installer/terraform/azureprivatedns

go 1.16

replace (
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
	github.com/hashicorp/go-azure-helpers => github.com/hashicorp/go-azure-helpers v0.16.5
)

require (
	github.com/Azure/azure-sdk-for-go v51.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.17
	github.com/hashicorp/go-azure-helpers v0.10.0
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/terraform-providers/terraform-provider-azurerm v1.44.0
)
