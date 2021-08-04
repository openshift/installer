terraform {
  required_version = ">= 0.14"
  required_providers {
    vsphere = {
      source = "openshift/local/vsphere"
    }
    vsphereprivate = {
      source = "openshift/local/vsphereprivate"
    }
  }
}

