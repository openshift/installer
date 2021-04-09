terraform {
  required_version = ">= 0.14"
  required_providers {
    kubernetes = {
      source = "openshift/local/kubernetes"
    }
    kubevirt = {
      source = "openshift/local/kubevirt"
    }
  }
}

