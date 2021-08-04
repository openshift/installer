terraform {
  required_version = ">= 1.0.0"
  required_providers {
    local = {
      source = "openshift/local/local"
    }
    azurerm = {
      source = "openshift/local/azurerm"
    }
    ignition = {
      source = "openshift/local/ignition"
    }
  }
}

