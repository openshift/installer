terraform {
  required_version = ">= 0.14"
  required_providers {
    libvirt = {
      source = "openshift/local/libvirt"
    }
  }
}

