resource "google_compute_address" "cluster_public_ip" {
  name = "${var.cluster_id}-cluster-public-ip"
}

resource "google_compute_http_health_check" "api" {
  name = "${var.cluster_id}-api"

  port         = 6080
  request_path = "/readyz"
}

resource "google_compute_target_pool" "api" {
  name = "${var.cluster_id}-api"

  instances     = var.master_instances
  health_checks = [google_compute_http_health_check.api.self_link]
}

resource "google_compute_forwarding_rule" "api" {
  name = "${var.cluster_id}-api"

  ip_address = google_compute_address.cluster_public_ip.address
  target     = google_compute_target_pool.api.self_link
  port_range = "6443"
}
