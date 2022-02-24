terraform {
  required_version = ">= 1.0.0"
  required_providers {
    alicloud = {
      source  = "hashicorp/alicloud"
      version = "1.148.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "3.70.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "2.90.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "4.5.0"
    }
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = "0.6.12"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.1.0"
    }
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "1.46.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.1.0"
    }
  }
}
