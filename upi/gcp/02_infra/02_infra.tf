terraform {
  # Infra manager supports specific Terraform veresions; ensure compatibility
  required_version = ">=1.2.3"
  required_providers {
    google = {
      source = "hashicorp/google"
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

variable "control_subnet" {
  type        = string
  description = "Subnet for the control plane instances."
}

variable "cluster_domain" {
  type        = string
  description = "ClusterName.BaseDomain"
}

variable "cluster_network" {
  type        = string
  description = "Full link to the cluster network."
}

# Terraform handles lists but the infra-manager --input-values only
# supports scalar types.
# If you require more or less zones, you must manually add them below
# as a single variable for each. You must add the zones to the
# locals `zones` list below.
variable "zone_0" {
  type        = string
  description = "Zone 1 for the instance types."
}

variable "zone_1" {
  type        = string
  description = "Zone 2 for the instance types."
}

variable "zone_2" {
  type        = string
  description = "Zone 3 for the instance types."
}

locals {
  zones = ["${var.zone_0}", "${var.zone_1}", "${var.zone_2}"]
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
  instances = [
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

resource "google_compute_address" "cluster_ip" {
  provider = google-beta

  name = "${var.infra_id}-cluster-ip"
  address_type = "INTERNAL"
  region = "${var.region}"
  subnetwork = "${var.control_subnet}"
}

resource "google_compute_health_check" "api_internal_health_check" {
  provider = google-beta

  name = "${var.infra_id}-api-internal-health-check"
  https_health_check {
    port = 6443
  }
}

resource "google_compute_region_backend_service" "api_internal" {
  provider = google-beta

  name = "${var.infra_id}-api-internal"
  timeout_sec = 120
  protocol = "TCP"
  region = "${var.region}"
  load_balancing_scheme = "INTERNAL"
  health_checks = [
    google_compute_health_check.api_internal_health_check.id
  ]

  dynamic "backend" {
    for_each = google_compute_instance_group.master_ig

    content {
      balancing_mode = "CONNECTION"
      group = backend.value.self_link
    }
  }
}

resource "google_compute_forwarding_rule" "api_internal_forwarding_rule" {
  provider = google-beta

  name = "${var.infra_id}-api-internal-forwarding-rule"
  ip_address = google_compute_address.cluster_ip.address
  backend_service = google_compute_region_backend_service.api_internal.id
  load_balancing_scheme = "INTERNAL"
  ports = [
    "6443",
    "22623"
  ]
  region = "${var.region}"
  subnetwork = "${var.control_subnet}"
}

resource "google_compute_instance_group" "master_ig" {
  provider = google-beta

  for_each = toset(local.zones)

  name = "${var.infra_id}-master-${each.key}-ig"
  network = "${var.cluster_network}"
  zone = "${each.key}"
  named_port {
    name = "ignition"
    port = 22623
  }
  named_port {
    name = "https"
    port = 6443
  }
}

