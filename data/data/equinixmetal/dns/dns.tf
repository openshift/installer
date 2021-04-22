provider "dns" {
  /**
  // TODO: accept dns update options so DNS can be configured following EM devices
  update {
    server        = "192.168.0.1"
    key_name      = "example.com."
    key_algorithm = "hmac-md5"
    key_secret    = "3VwZXJzZWNyZXQ="
  }
  **/
}

locals {
  basedomain = join(".", var.cluster_name, replace(var.cluster_basedomain, "${var.cluster_name}.", ""))
}

data "dns_a_record_set" "bootstrap" {
  host = "bootstrap.${local.basedomain}"
}

data "dns_a_record_set" "masters" {
  count = var.masters_count
  host  = "master${count.index}.${local.basedomain}"
}

/*
data "dns_a_record_set" "etcd_a" {
  count = var.masters_count
  host  = "etcd-${count.index}.${local.basedomain}"
}

data "dns_srv_record_set" "etcd_srv" {
  // Verifies etcd SRV records have been created
  // TODO: verify that these match the etcd nodes
  service = "_etcd-server-ssl._tcp.${local.basedomain}"
}
*/

/*
data "dns_a_record_set" "workers" {
  count = var.workers_count
  host  = "worker${count.index}.${local.basedomain}"
}
*/

data "dns_a_record_set" "lb" {
  host = "api-int.${local.basedomain}"
}

data "dns_a_record_set" "apps" {
  // TODO: validate that *.apps matches the lb
  // TODO: permit CNAME, dns_cname_record_set is valid 
  host = "*.apps.${local.basedomain}"
}

