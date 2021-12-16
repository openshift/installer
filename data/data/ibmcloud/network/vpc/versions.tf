terraform {
  required_version = ">= 0.14"
  required_providers {
    ibm = {
      source = "openshift/local/ibm"
    }
  }
}
