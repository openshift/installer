terraform {
  required_version = ">= 0.12"
  required_providers {
    local = {
      source = "openshift/local/local"
    }
    azurestack = {
      source = "openshift/local/azurestack"
    }
    ignition = {
      source = "openshift/local/ignition"
    }
  }
}

