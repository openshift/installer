terraform {
  required_version = ">= 0.14"
  required_providers {
    kubevirt = {
      source = "openshift/local/kubevirt"
    }
  }
}

