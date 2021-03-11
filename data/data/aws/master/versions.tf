terraform {
  required_version = ">= 0.14"
  required_providers {
    aws = {
      source = "openshift/local/aws"
    }
  }
}

