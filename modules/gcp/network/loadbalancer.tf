resource "google_compute_target_pool" "master-targetpool" {
  name = "${var.cluster_name}-master-targetpool"
}

resource "google_compute_target_pool" "worker-targetpool" {
  name = "${var.cluster_name}-worker-targetpool"

  health_checks = [
    "${google_compute_http_health_check.worker-hc.name}",
  ]
}

resource "google_compute_http_health_check" "worker-hc" {
  name         = "${var.cluster_name}-worker-hc"
  request_path = "/"

  timeout_sec        = 1
  check_interval_sec = 1
}

resource "google_compute_address" "masters-ip" {
  name = "${var.cluster_name}-masters-ip"
}

resource "google_compute_forwarding_rule" "api-external-fwd-rule" {
  load_balancing_scheme = "EXTERNAL"
  name                  = "${var.cluster_name}-api-external-fwd-rule"
  ip_address            = "${google_compute_address.masters-ip.address}"
  region                = "${var.gcp_region}"
  target                = "${google_compute_target_pool.master-targetpool.self_link}"
  port_range            = "443"
}

resource "google_compute_address" "ingress-ip" {
  name = "${var.cluster_name}-ingress-ip"
}

resource "google_compute_forwarding_rule" "ingress-external-http-fwd-rule" {
  load_balancing_scheme = "EXTERNAL"
  name                  = "${var.cluster_name}-ingress-external-http-fwd-rule"
  ip_address            = "${google_compute_address.ingress-ip.address}"
  region                = "${var.gcp_region}"
  target                = "${google_compute_target_pool.worker-targetpool.self_link}"
  port_range            = "80"
}

resource "google_compute_forwarding_rule" "ingress-external-https-fwd-rule" {
  load_balancing_scheme = "EXTERNAL"
  name                  = "${var.cluster_name}-ingress-external-https-fwd-rule"
  ip_address            = "${google_compute_address.ingress-ip.address}"
  region                = "${var.gcp_region}"
  target                = "${google_compute_target_pool.worker-targetpool.self_link}"
  port_range            = "443"
}
