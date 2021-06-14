locals {
  description = "Created By OpenShift Installer"

  public_endpoints = var.gcp_publish_strategy == "External" ? true : false
}

provider "google" {
  credentials = var.gcp_service_account
  project     = var.gcp_project_id
  region      = var.gcp_region
}

resource "google_compute_region_backend_service" "api_internal" {
  name        = "${var.cluster_id}-api-internal"
  description = local.description

  load_balancing_scheme = "INTERNAL"
  protocol              = "TCP"
  timeout_sec           = 120

  dynamic "backend" {
    for_each = var.gcp_bootstrap_lb ? concat(var.bootstrap_instance_groups, var.master_instance_groups) : var.master_instance_groups

    content {
      group = backend.value
    }
  }

  health_checks = var.api_internal_health_checks
}

resource "google_compute_forwarding_rule" "api_internal" {
  name        = "${var.cluster_id}-api-internal"
  description = local.description

  ip_address      = var.cluster_ip
  backend_service = google_compute_region_backend_service.api_internal.self_link
  ports           = ["6443", "22623"]
  subnetwork      = var.master_subnet
  network         = var.network

  load_balancing_scheme = "INTERNAL"
}

resource "google_compute_target_pool" "api" {
  count = local.public_endpoints ? 1 : 0

  name        = "${var.cluster_id}-api"
  description = local.description

  instances     = var.gcp_bootstrap_lb ? concat(var.bootstrap_instances, var.master_instances) : var.master_instances
  health_checks = var.api_health_checks
}

resource "google_compute_forwarding_rule" "api" {
  count = local.public_endpoints ? 1 : 0

  name        = "${var.cluster_id}-api"
  description = local.description

  ip_address = var.cluster_public_ip
  target     = google_compute_target_pool.api[0].self_link
  port_range = "6443"
}
