terraform {
  required_version = ">= 0.14"
  required_providers {
    azurerm = {
      source = "openshift/local/azurerm"
    }
  }
}

