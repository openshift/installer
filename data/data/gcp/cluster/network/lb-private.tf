resource "google_compute_address" "cluster_ip" {
  name         = "${var.cluster_id}-cluster-ip"
  address_type = "INTERNAL"
  subnetwork   = local.master_subnet
  description  = local.description
}

// Refer to docs/dev/kube-apiserver-health-check.md on how to correctly setup health check probe for kube-apiserver
resource "google_compute_health_check" "api_internal" {
  name        = "${var.cluster_id}-api-internal"
  description = local.description

  healthy_threshold   = 3
  unhealthy_threshold = 3
  check_interval_sec  = 2
  timeout_sec         = 2

  https_health_check {
    port         = 6443
    request_path = "/readyz"
  }
}
