terraform {
  required_version = ">= 0.14"
  required_providers {
    azureprivatedns = {
      source = "openshift/local/azureprivatedns"
    }
    azurerm = {
      source = "openshift/local/azurerm"
    }
  }
}

