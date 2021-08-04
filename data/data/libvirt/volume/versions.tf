terraform {
  required_version = ">= 1.0.0"
  required_providers {
    libvirt = {
      source = "openshift/local/libvirt"
    }
  }
}

