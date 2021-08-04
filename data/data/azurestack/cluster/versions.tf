terraform {
  required_version = ">= 1.0.0"
  required_providers {
    azurestack = {
      source = "openshift/local/azurestack"
    }
  }
}
