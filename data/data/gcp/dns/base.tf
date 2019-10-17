resource "google_dns_managed_zone" "int" {
  name       = "${var.cluster_id}-private-zone"
  dns_name   = "${var.cluster_domain}."
  visibility = "private"

  private_visibility_config {
    networks {
      network_url = var.network
    }
  }
}

resource "google_dns_record_set" "api_external" {
  count = var.public_endpoints ? 1 : 0

  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = var.public_dns_zone_name
  rrdatas      = [var.api_external_lb_ip]
}

resource "google_dns_record_set" "api_internal" {
  name         = "api-int.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = google_dns_managed_zone.int.name
  rrdatas      = [var.api_internal_lb_ip]
}

resource "google_dns_record_set" "api_external_internal_zone" {
  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = google_dns_managed_zone.int.name
  rrdatas      = [var.api_internal_lb_ip]
}

resource "google_dns_record_set" "etcd_a_nodes" {
  count        = var.etcd_count
  type         = "A"
  ttl          = "60"
  managed_zone = google_dns_managed_zone.int.name
  name         = "etcd-${count.index}.${var.cluster_domain}."
  rrdatas      = [var.etcd_ip_addresses[count.index]]
}

resource "google_dns_record_set" "etcd_cluster" {
  type         = "SRV"
  ttl          = "60"
  managed_zone = google_dns_managed_zone.int.name
  name         = "_etcd-server-ssl._tcp.${var.cluster_domain}."
  rrdatas      = formatlist("0 10 2380 %s", google_dns_record_set.etcd_a_nodes.*.name)
}
