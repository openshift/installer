resource "google_compute_address" "cluster_public_ip" {
  count = var.public_endpoints ? 1 : 0

  name = "${var.cluster_id}-cluster-public-ip"
}

// Refer to docs/dev/kube-apiserver-health-check.md on how to correctly setup health check probe for kube-apiserver
resource "google_compute_http_health_check" "api" {
  count = var.public_endpoints ? 1 : 0

  name = "${var.cluster_id}-api"

  healthy_threshold   = 3
  unhealthy_threshold = 3
  check_interval_sec  = 2
  timeout_sec         = 2

  port         = 6080
  request_path = "/readyz"
}

resource "google_compute_target_pool" "api" {
  count = var.public_endpoints ? 1 : 0

  name = "${var.cluster_id}-api"

  instances     = var.bootstrap_lb ? concat(var.bootstrap_instances, var.master_instances) : var.master_instances
  health_checks = [google_compute_http_health_check.api[0].self_link]
}

resource "google_compute_forwarding_rule" "api" {
  count = var.public_endpoints ? 1 : 0

  name = "${var.cluster_id}-api"

  ip_address = google_compute_address.cluster_public_ip[0].address
  target     = google_compute_target_pool.api[0].self_link
  port_range = "6443"
}
