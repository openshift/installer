terraform {
  required_version = ">= 1.0.0"
  required_providers {
    ironic = {
      source = "openshift/local/ironic"
    }
  }
}

