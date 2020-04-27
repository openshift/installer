resource "google_compute_address" "cluster_ip" {
  name         = "${var.cluster_id}-cluster-ip"
  address_type = "INTERNAL"
  subnetwork   = local.master_subnet
}

resource "google_compute_health_check" "api_internal" {
  name = "${var.cluster_id}-api-internal"

  // CAUTION: the gcp-routes mechanism must be _faster_ than this value:
  //
  //          Otherwise, local client traffic will go to the GCP LB until the
  //          gcp-routes mechanism is done rerouting, and therefore back to the
  //          local node in 1/3 of cases, which is blackholed (due to missing
  //          hairpinning support).
  healthy_threshold   = 3
  // CAUTION: the gcp-routes mechanism must be _slower_ than this value:
  //
  //          Otherwise, local client traffic  will go to the GCP LB until the LB
  //          endpoint is deactived, and therefore back to the local node in 1/3
  //          of cases, which is blackholed (due to missing hairpinning support).
  unhealthy_threshold = 3
  check_interval_sec  = 2
  timeout_sec         = 2

  https_health_check {
    port         = 6443
    request_path = "/readyz"
  }
}

resource "google_compute_region_backend_service" "api_internal" {
  name = "${var.cluster_id}-api-internal"

  load_balancing_scheme = "INTERNAL"
  protocol              = "TCP"
  timeout_sec           = 120

  dynamic "backend" {
    for_each = var.bootstrap_lb ? concat(var.bootstrap_instance_groups, var.master_instance_groups) : var.master_instance_groups

    content {
      group = backend.value
    }
  }

  health_checks = [google_compute_health_check.api_internal.self_link]
}

resource "google_compute_forwarding_rule" "api_internal" {
  name = "${var.cluster_id}-api-internal"

  ip_address      = google_compute_address.cluster_ip.address
  backend_service = google_compute_region_backend_service.api_internal.self_link
  ports           = ["6443", "22623"]
  subnetwork      = local.master_subnet
  network         = local.cluster_network

  load_balancing_scheme = "INTERNAL"
}
