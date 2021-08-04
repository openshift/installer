terraform {
  required_version = ">= 1.0.0"
  required_providers {
    openstack = {
      source = "openshift/local/openstack"
    }
  }
}

