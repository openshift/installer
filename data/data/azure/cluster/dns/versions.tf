terraform {
  required_version = ">= 1.0.0"
  required_providers {
    azureprivatedns = {
      source = "openshift/local/azureprivatedns"
    }
    azurerm = {
      source = "openshift/local/azurerm"
    }
  }
}

