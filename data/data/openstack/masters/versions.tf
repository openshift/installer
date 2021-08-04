terraform {
  required_version = ">= 0.14"
  required_providers {
    ignition = {
      source = "openshift/local/ignition"
    }
    openstack = {
      source = "openshift/local/openstack"
    }
  }
}

