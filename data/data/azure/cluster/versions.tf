terraform {
  required_version = ">= 1.0.0"
  required_providers {
    random = {
      source = "openshift/local/random"
    }
    azurerm = {
      source = "openshift/local/azurerm"
    }
    azureprivatedns = {
      source = "openshift/local/azureprivatedns"
    }
  }
}
