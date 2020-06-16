/*
resource "cloudflare_record" "dns_a_cluster_api" {
  zone_id = var.cf_zone_id
  type    = "A"
  name    = "api.${var.cluster_name}.${var.cluster_basedomain}"
  value   = var.node_ips[count.index]
  count   = (var.node_type == "lb" ? length(var.node_ips) : 0)
}

resource "cloudflare_record" "dns_a_cluster_api_int" {
  zone_id = var.cf_zone_id
  type    = "A"
  name    = "api-int.${var.cluster_name}.${var.cluster_basedomain}"
  value   = var.node_ips[count.index]
  count   = (var.node_type == "lb" ? length(var.node_ips) : 0)
}

resource "cloudflare_record" "dns_a_cluster_wildcard_https" {
  zone_id = var.cf_zone_id
  type    = "A"
  name    = "*.apps.${var.cluster_name}.${var.cluster_basedomain}"
  value   = var.node_ips[count.index]
  count   = (var.node_type == "lb" ? length(var.node_ips) : 0)
}

resource "cloudflare_record" "dns_a_node" {
  zone_id = var.cf_zone_id
  type    = "A"
  name    = "${var.node_type}-${count.index}.${var.cluster_name}.${var.cluster_basedomain}"
  value   = var.node_ips[count.index]
  count   = length(var.node_ips)
}

resource "cloudflare_record" "dns_a_etcd" {
  zone_id = var.cf_zone_id
  type    = "A"
  name    = "etcd-${count.index}.${var.cluster_name}.${var.cluster_basedomain}"
  value   = var.node_ips[count.index]
  count   = (var.node_type == "master" ? length(var.node_ips) : 0)
}

resource "cloudflare_record" "dns_srv_etcd" {
  zone_id = var.cf_zone_id
  type    = "SRV"
  name    = "_etcd-server-ssl._tcp"
  count   = (var.node_type == "master" ? length(var.node_ips) : 0)

  data = {
    service  = "_etcd-server-ssl"
    proto    = "_tcp"
    name     = "${var.cluster_name}.${var.cluster_basedomain}"
    priority = 0
    weight   = 10
    port     = 2380
    target   = "etcd-${count.index}.${var.cluster_name}.${var.cluster_basedomain}"
  }

}
*/