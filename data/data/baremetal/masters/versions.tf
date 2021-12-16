terraform {
  required_version = ">= 0.14"
  required_providers {
    ironic = {
      source = "openshift/local/ironic"
    }
  }
}

