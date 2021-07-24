terraform {
  required_version = ">= 0.12"
  required_providers {
    random = {
      source = "openshift/local/random"
    }
    azurestack = {
      source = "openshift/local/azurestack"
    }
  }
}

