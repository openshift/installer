terraform {
  required_version = ">= 1.0.0"
  required_providers {
    kubevirt = {
      source = "openshift/local/kubevirt"
    }
  }
}

