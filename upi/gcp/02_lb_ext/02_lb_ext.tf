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

resource "google_compute_address" "cluster_public_ip" {
  provider = google-beta

  name = "${var.infra_id}-cluster-public-ip"
  region = "${var.region}"
}

resource "google_compute_http_health_check" "api_http_health_check" {
  provider = google-beta

  name = "${var.infra_id}-api-http-health-check"
  port = 6080
  request_path = "/readyz"
}

resource "google_compute_target_pool" "api_target_pool" {
  provider = google-beta

  name = "${var.infra_id}-api-target-pool"
  region = "${var.region}"
  health_checks = [
    google_compute_http_health_check.api_http_health_check.id
  ]
}

resource "google_compute_forwarding_rule" "api_forwarding_rule" {
  provider = google-beta

  name = "${var.infra_id}-api-forwarding-rule"
  ip_address = google_compute_address.cluster_public_ip.address
  port_range = "6443"
  region = "${var.region}"
  target = google_compute_target_pool.api_target_pool.id
}
