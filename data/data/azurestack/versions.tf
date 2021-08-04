terraform {
  required_version = ">= 1.0.0"
  required_providers {
    random = {
      source = "openshift/local/random"
    }
    azurestack = {
      source = "openshift/local/azurestack"
    }
  }
}

