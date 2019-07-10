resource "google_compute_address" "cluster_public_ip" {
  name = "${var.cluster_id}-cluster-public-ip"
}

resource "google_compute_http_health_check" "api_external" {
  name = "${var.cluster_id}-api-external"

  port         = 6443
  request_path = "/readyz"
}

resource "google_compute_target_pool" "api_external" {
  name = "${var.cluster_id}-api-external"

  instances     = concat(var.master_instances, [var.bootstrap_instance])
  health_checks = [google_compute_http_health_check.api_external.self_link]
}

resource "google_compute_forwarding_rule" "api_external" {
  name = "${var.cluster_id}-api-external"

  ip_address = google_compute_address.cluster_public_ip.address
  target     = google_compute_target_pool.api_external.self_link
  port_range = "6443"
}
