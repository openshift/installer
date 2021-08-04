terraform {
  required_version = ">= 0.14"
  required_providers {
    openstack = {
      source = "openshift/local/openstack"
    }
  }
}

