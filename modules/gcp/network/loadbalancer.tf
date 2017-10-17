resource "google_compute_target_pool" "tectonic-master-targetpool" {
  name = "${var.cluster_name}-tectonic-master-targetpool"
}

resource "google_compute_target_pool" "tectonic-worker-targetpool" {
  name = "${var.cluster_name}-tectonic-worker-targetpool"

  health_checks = [
    "${google_compute_http_health_check.tectonic-worker-hc.name}",
  ]
}

resource "google_compute_http_health_check" "tectonic-worker-hc" {
  name         = "${var.cluster_name}-tectonic-worker-hc"
  request_path = "/"

  timeout_sec        = 1
  check_interval_sec = 1
}

resource "google_compute_address" "tectonic-masters-ip" {
  name = "${var.cluster_name}-tectonic-masters-ip"
}

resource "google_compute_forwarding_rule" "tectonic-api-external-fwd-rule" {
  load_balancing_scheme = "EXTERNAL"
  name                  = "${var.cluster_name}-tectonic-api-external-fwd-rule"
  ip_address            = "${google_compute_address.tectonic-masters-ip.address}"
  region                = "${var.gcp_region}"
  target                = "${google_compute_target_pool.tectonic-master-targetpool.self_link}"
  port_range            = "443"
}

resource "google_compute_address" "tectonic-ingress-ip" {
  name = "${var.cluster_name}-tectonic-ingress-ip"
}

resource "google_compute_forwarding_rule" "tectonic-ingress-external-http-fwd-rule" {
  load_balancing_scheme = "EXTERNAL"
  name                  = "${var.cluster_name}-tectonic-ingress-external-http-fwd-rule"
  ip_address            = "${google_compute_address.tectonic-ingress-ip.address}"
  region                = "${var.gcp_region}"
  target                = "${google_compute_target_pool.tectonic-worker-targetpool.self_link}"
  port_range            = "80"
}

resource "google_compute_forwarding_rule" "tectonic-ingress-external-https-fwd-rule" {
  load_balancing_scheme = "EXTERNAL"
  name                  = "${var.cluster_name}-tectonic-ingress-external-https-fwd-rule"
  ip_address            = "${google_compute_address.tectonic-ingress-ip.address}"
  region                = "${var.gcp_region}"
  target                = "${google_compute_target_pool.tectonic-worker-targetpool.self_link}"
  port_range            = "443"
}
