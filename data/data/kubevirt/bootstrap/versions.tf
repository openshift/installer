terraform {
  required_version = ">= 1.0.0"
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

