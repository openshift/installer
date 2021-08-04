terraform {
  required_version = ">= 0.14"
  required_providers {
    azurestack = {
      source = "openshift/local/azurestack"
    }
  }
}
