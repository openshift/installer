
resource "google_compute_global_address" "master" {
  name = "${var.cluster_name}-master"
}

resource "google_compute_health_check" "masters" {
  name = "${var.cluster_name}-masters"

  https_health_check {
    port = 6443
    request_path = "/healthz"
  }
}

resource "google_compute_backend_service" "masters" {
  name        = "${var.cluster_name}-masters"
  port_name   = "https"
  protocol    = "TCP"
  timeout_sec = 120

  backend {
    group = "${var.bootstrap_instance_group}"
  }
  // gross, in 0.12 this hopefully gets better
  backend {
    group = "${var.master_instance_groups[0]}"
  }
  backend {
    group = "${var.master_instance_groups[1]}"
  }
  backend {
    group = "${var.master_instance_groups[2]}"
  }

  health_checks = ["${google_compute_health_check.masters.self_link}"]
}

resource "google_compute_target_tcp_proxy" "master" {
  name            = "${var.cluster_name}-master-tcp-proxy"
  backend_service = "${google_compute_backend_service.masters.self_link}"
}

resource "google_compute_target_pool" "master-internal" {
  name = "${var.cluster_name}-master-internal"

  instances = ["${concat(list(var.bootstrap_instance),var.master_instances)}"]
}

resource "google_compute_global_forwarding_rule" "master" {
  count      = "${var.public_master_endpoints ? 1 : 0}"
  name       = "${var.cluster_name}-master"
  ip_address = "${google_compute_global_address.master.address}"
  target     = "${google_compute_target_tcp_proxy.master.self_link}"
  port_range = "443"
}

resource "google_compute_forwarding_rule" "master-internal" {
  count      = "${var.private_master_endpoints ? 1 : 0}"
  name       = "${var.cluster_name}-master-internal"
  target     = "${google_compute_target_pool.master-internal.self_link}"
  port_range = "443"
}
