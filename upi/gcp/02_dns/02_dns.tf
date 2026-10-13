terraform {
  # Infra manager supports specific Terraform versions; ensure compatibility
  required_version = ">=1.2.3"
  required_providers {
    google = {
      source = "hashicorp/google"
      version = ">= 4.0.0"
    }
    google-beta = {
        source = "hashicorp/google-beta",
        version = ">= 4.0.0"
    }
  }
}

provider "google-beta" {
  project = "${var.project}"
  region = "${var.region}"
}

variable "infra_id" {
  type        = string
  description = "OpenShift Installer Infrastructure ID"
}

variable "project" {
  type        = string
  description = "Project ID"
}

variable "region" {
  type        = string
  description = "GCP Region where the resources will be created."
  default     = "us-central1"
}

variable "cluster_domain" {
  type        = string
  description = "ClusterName.BaseDomain"
}

variable "cluster_network" {
  type        = string
  description = "Full link to the cluster network."
}

resource "google_dns_managed_zone" "private_zone" {
  provider = google-beta

  name = "${var.infra_id}-private-zone"
  dns_name = "${var.cluster_domain}."
  description = "OpenShift Installer UPI create private DNS zone."
  visibility = "private"
  private_visibility_config {
    networks {
      network_url = "${var.cluster_network}"
    }
  }

  force_destroy = false
}
