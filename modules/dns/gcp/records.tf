resource "google_dns_record_set" "api-external" {
  name         = "${var.cluster_name}.api.${var.base_domain}."
  type         = "A"
  ttl          = 300
  managed_zone = "${var.managed_zone_name}"
  rrdatas      = ["${var.tectonic_masters_ip}"]
}

resource "google_dns_record_set" "ingress-external" {
  name         = "${var.cluster_name}.${var.base_domain}."
  type         = "A"
  ttl          = 300
  managed_zone = "${var.managed_zone_name}"
  rrdatas      = ["${var.tectonic_ingress_ip}"]
}

# Etcd
resource "google_dns_record_set" "etcd_srv_discover" {
  count        = "${var.etcd_dns_enabled ? 1 : 0}"
  name         = "${var.tls_enabled ? "_etcd-server-ssl._tcp.${var.base_domain}" : "_etcd-server._tcp.${var.base_domain}"}."
  type         = "SRV"
  managed_zone = "${var.managed_zone_name}"
  rrdatas      = ["${formatlist("0 0 2380 %s", google_dns_record_set.etc_a_node.*.name)}"]
  ttl          = "300"
}

resource "google_dns_record_set" "etcd_srv_client" {
  count        = "${var.etcd_dns_enabled ? 1 : 0}"
  name         = "${var.tls_enabled ? "_etcd-client-ssl._tcp.${var.base_domain}" : "_etcd-client._tcp.${var.base_domain}"}."
  type         = "SRV"
  managed_zone = "${var.managed_zone_name}"
  rrdatas      = ["${formatlist("0 0 2379 %s", google_dns_record_set.etc_a_node.*.name)}"]
  ttl          = "60"
}

resource "google_dns_record_set" "etc_a_node" {
  count        = "${var.etcd_dns_enabled ? var.etcd_instance_count : 0}"
  type         = "A"
  ttl          = "60"
  managed_zone = "${var.managed_zone_name}"
  name         = "${var.cluster_name}-etcd-${count.index}.${var.base_domain}."
  rrdatas      = ["${var.etcd_ip_addresses[count.index]}"]
}
