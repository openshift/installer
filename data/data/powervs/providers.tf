# TODO(mjturek): Remove this file once no longer needed

terraform {
  required_providers {
    ibm = {
      source  = "ibm-cloud/ibm"
      version = "1.25.0"
    }
    ignition = {
      source = "community-terraform-providers/ignition"
      version = "2.1.2"
    }
    presign = {}
  }
}
