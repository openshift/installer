terraform {
  required_version = ">= 1.0.0"
  required_providers {
    ibm = {
      source = "openshift/local/ibm"
    }
  }
}
