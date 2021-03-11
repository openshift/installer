terraform {
  required_version = ">= 0.14"
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

