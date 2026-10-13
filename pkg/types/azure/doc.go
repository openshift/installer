// Package azure contains Azure-specific structures for installer
// configuration and management.
// +k8s:deepcopy-gen=package
package azure

// Name is the name for the Azure platform.
const Name string = "azure"

// StackTerraformName is the name used for Terraform code when installing to the Azure Stack platform.
const StackTerraformName string = "azurestack"
