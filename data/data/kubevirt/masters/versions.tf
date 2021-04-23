terraform {
  required_version = ">= 0.14"
  required_providers {
    ignition = {
      source = "openshift/local/ignition"
    }
    kubernetes = {
      source = "openshift/local/kubernetes"
    }
    kubevirt = {
      source = "openshift/local/kubevirt"
    }
  }
}

