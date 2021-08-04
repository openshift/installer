terraform {
  required_version = ">= 1.0.0"
  required_providers {
    vsphere = {
      source = "openshift/local/vsphere"
    }
  }
}

