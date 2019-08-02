resource "google_compute_address" "cluster_public_ip" {
  name = "${var.cluster_id}-cluster-public-ip"
}

resource "google_compute_global_address" "cluster_private_ip" {
  name = "${var.cluster_id}-cluster-private-ip"
}

resource "google_compute_health_check" "api_internal" {
  name = "${var.cluster_id}-ign-internal"

  tcp_health_check {
    port = "5222"
  }
}

resource "google_compute_backend_service" "api_internal" {
  name = "${var.cluster_id}-api-internal"

  protocol    = "TCP"
  port_name   = "tcp5222"
  timeout_sec = 120

  dynamic "backend" {
    for_each = var.bootstrap_lb ? var.bootstrap_instance_groups : var.master_instance_groups

    content {
      group = backend.value
    }
  }

  health_checks = [google_compute_health_check.api_internal.self_link]
}

resource "google_compute_target_tcp_proxy" "api_internal" {
  name            = "${var.cluster_id}-api-internal"
  backend_service = google_compute_backend_service.api_internal.self_link
}

resource "google_compute_global_forwarding_rule" "api_internal" {
  name                  = "${var.cluster_id}-api-internal"
  target                = google_compute_target_tcp_proxy.api_internal.self_link
  port_range            = "5222"
  load_balancing_scheme = "EXTERNAL"
  ip_address            = google_compute_global_address.cluster_private_ip.address
}

resource "google_compute_http_health_check" "api" {
  name = "${var.cluster_id}-api"

  port         = 6080
  request_path = "/readyz"
}

resource "google_compute_target_pool" "api" {
  name = "${var.cluster_id}-api"

  instances     = var.bootstrap_lb ? var.bootstrap_instances : var.master_instances
  health_checks = [google_compute_http_health_check.api.self_link]
}

resource "google_compute_forwarding_rule" "api" {
  name = "${var.cluster_id}-api"

  ip_address = google_compute_address.cluster_public_ip.address
  target     = google_compute_target_pool.api.self_link
  port_range = "6443"
}

resource "google_compute_http_health_check" "ign" {
  name = "${var.cluster_id}-ign"

  port         = 22624
  request_path = "/healthz"
}

resource "google_compute_target_pool" "ign" {
  name = "${var.cluster_id}-ign"

  instances     = var.bootstrap_lb ? var.bootstrap_instances : var.master_instances
  health_checks = [google_compute_http_health_check.ign.self_link]
}

resource "google_compute_forwarding_rule" "ign" {
  name = "${var.cluster_id}-ign"

  ip_address = google_compute_address.cluster_public_ip.address
  target     = google_compute_target_pool.ign.self_link
  port_range = "22623"
}
