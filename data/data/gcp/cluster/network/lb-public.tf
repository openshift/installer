resource "google_compute_address" "cluster_public_ip" {
  count       = var.public_endpoints ? 1 : 0
  description = local.description

  name = "${var.cluster_id}-cluster-public-ip"
}

// Refer to docs/dev/kube-apiserver-health-check.md on how to correctly setup health check probe for kube-apiserver
resource "google_compute_http_health_check" "api" {
  count       = var.public_endpoints ? 1 : 0
  description = local.description

  name = "${var.cluster_id}-api"

  healthy_threshold   = 3
  unhealthy_threshold = 3
  check_interval_sec  = 2
  timeout_sec         = 2

  port         = 6080
  request_path = "/readyz"
}
