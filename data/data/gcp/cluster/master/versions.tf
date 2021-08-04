terraform {
  required_version = ">= 0.14"
  required_providers {
    google = {
      source = "openshift/local/google"
    }
  }
}

