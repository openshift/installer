locals {
  scheme = "${var.etcd_tls_enabled ? "https" : "http"}"

  // see https://github.com/hashicorp/terraform/issues/9858
  etcd_initial_cluster_list = "${concat(var.etcd_initial_cluster_list, list("dummy"))}"

  metadata_env = "EnvironmentFile=/run/metadata/coreos"

  metadata_deps = <<EOF
Requires=coreos-metadata.service
After=coreos-metadata.service
EOF

  cert_options = <<EOF
--cert-file=/etc/ssl/etcd/server.crt \
  --key-file=/etc/ssl/etcd/server.key \
  --peer-cert-file=/etc/ssl/etcd/peer.crt \
  --peer-key-file=/etc/ssl/etcd/peer.key \
  --peer-trusted-ca-file=/etc/ssl/etcd/ca.crt \
  --peer-client-cert-auth=trueEOF
}

data "template_file" "etcd_names" {
  count    = "${var.etcd_count}"
  template = "${var.cluster_name}-etcd-${count.index}${var.base_domain == "" ? "" : ".${var.base_domain}"}"
}

data "template_file" "advertise_client_urls" {
  count    = "${var.etcd_count}"
  template = "${local.scheme}://${var.etcd_advertise_name_list[count.index]}:2379"
}

data "template_file" "initial_advertise_peer_urls" {
  count    = "${var.etcd_count}"
  template = "${local.scheme}://${var.etcd_advertise_name_list[count.index]}:2380"
}

data "template_file" "initial_cluster" {
  count    = "${length(var.etcd_initial_cluster_list) > 0 ? var.etcd_count : 0}"
  template = "${data.template_file.etcd_names.*.rendered[count.index]}=${local.scheme}://${local.etcd_initial_cluster_list[count.index]}:2380"
}

data "template_file" "etcd" {
  count    = "${var.etcd_count}"
  template = "${file("${path.module}/resources/dropins/40-etcd-cluster.conf")}"

  vars = {
    advertise_client_urls       = "${data.template_file.advertise_client_urls.*.rendered[count.index]}"
    cert_options                = "${var.etcd_tls_enabled ? local.cert_options : ""}"
    container_image             = "${var.container_images["etcd"]}"
    initial_advertise_peer_urls = "${data.template_file.initial_advertise_peer_urls.*.rendered[count.index]}"
    initial_cluster             = "${length(var.etcd_initial_cluster_list) > 0 ? format("--initial-cluster=%s", join(",", data.template_file.initial_cluster.*.rendered)) : ""}"
    metadata_deps               = "${var.use_metadata ? local.metadata_deps : ""}"
    metadata_env                = "${var.use_metadata ? local.metadata_env : ""}"
    name                        = "${data.template_file.etcd_names.*.rendered[count.index]}"
    scheme                      = "${local.scheme}"
  }
}

data "ignition_systemd_unit" "etcd" {
  count  = "${var.etcd_count}"
  name   = "etcd-member.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-cluster.conf"
      content = "${data.template_file.etcd.*.rendered[count.index]}"
    },
  ]
}
