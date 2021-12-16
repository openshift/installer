terraform {
  required_version = ">= 0.14"
  required_providers {
    alicloud = {
      source = "openshift/local/alicloud"
    }
  }
}
