terraform {
  required_version = ">= 0.12"
  required_providers {
    azurestack = {
      source = "openshift/local/azurestack"
    }
  }
}

