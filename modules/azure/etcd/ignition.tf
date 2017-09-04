data "ignition_config" "etcd" {
  count = "${var.etcd_count}"

  systemd = [
    "${data.ignition_systemd_unit.locksmithd.*.id[count.index]}",
    "${data.ignition_systemd_unit.etcd3.*.id[count.index]}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]

  files = [
    "${data.ignition_file.etcd_ca.id}",
    "${data.ignition_file.etcd_server_crt.id}",
    "${data.ignition_file.etcd_server_key.id}",
    "${data.ignition_file.etcd_client_crt.id}",
    "${data.ignition_file.etcd_client_key.id}",
    "${data.ignition_file.etcd_peer_crt.id}",
    "${data.ignition_file.etcd_peer_key.id}",
  ]
}

data "ignition_file" "etcd_ca" {
  count = "${var.etcd_count > 0 ? 1 : 0}"

  path       = "/etc/ssl/etcd/ca.crt"
  mode       = 0644
  uid        = 232
  gid        = 232
  filesystem = "root"

  content {
    content = "${var.tls_ca_crt_pem}"
  }
}

data "ignition_file" "etcd_client_key" {
  path       = "/etc/ssl/etcd/client.key"
  mode       = 0400
  uid        = 0
  gid        = 0
  filesystem = "root"

  content {
    content = "${var.tls_client_key_pem}"
  }
}

data "ignition_file" "etcd_client_crt" {
  path       = "/etc/ssl/etcd/client.crt"
  mode       = 0400
  uid        = 0
  gid        = 0
  filesystem = "root"

  content {
    content = "${var.tls_client_crt_pem}"
  }
}

data "ignition_file" "etcd_server_key" {
  count = "${var.etcd_count > 0 ? 1 : 0}"

  path       = "/etc/ssl/etcd/server.key"
  mode       = 0400
  uid        = 232
  gid        = 232
  filesystem = "root"

  content {
    content = "${var.tls_server_key_pem}"
  }
}

data "ignition_file" "etcd_server_crt" {
  count = "${var.etcd_count > 0 ? 1 : 0}"

  path       = "/etc/ssl/etcd/server.crt"
  mode       = 0400
  uid        = 232
  gid        = 232
  filesystem = "root"

  content {
    content = "${var.tls_server_crt_pem}"
  }
}

data "ignition_file" "etcd_peer_key" {
  count = "${var.etcd_count > 0 ? 1 : 0}"

  path       = "/etc/ssl/etcd/peer.key"
  mode       = 0400
  uid        = 232
  gid        = 232
  filesystem = "root"

  content {
    content = "${var.tls_peer_key_pem}"
  }
}

data "ignition_file" "etcd_peer_crt" {
  count = "${var.etcd_count > 0 ? 1 : 0}"

  path       = "/etc/ssl/etcd/peer.crt"
  mode       = 0400
  uid        = 232
  gid        = 232
  filesystem = "root"

  content {
    content = "${var.tls_peer_crt_pem}"
  }
}

data "ignition_user" "core" {
  count = "${var.etcd_count > 0 ? 1 : 0}"

  name = "core"

  ssh_authorized_keys = [
    "${file(var.public_ssh_key)}",
  ]
}

data "ignition_systemd_unit" "locksmithd" {
  count = "${var.etcd_count}"

  name   = "locksmithd.service"
  enable = true

  dropin = [
    {
      name = "40-etcd-lock.conf"

      content = <<EOF
[Service]
Environment=REBOOT_STRATEGY=etcd-lock
${var.tls_enabled ? "Environment=\"LOCKSMITHD_ETCD_CAFILE=/etc/ssl/etcd/ca.crt\"" : ""}
${var.tls_enabled ? "Environment=\"LOCKSMITHD_ETCD_KEYFILE=/etc/ssl/etcd/client.key\"" : ""}
${var.tls_enabled ? "Environment=\"LOCKSMITHD_ETCD_CERTFILE=/etc/ssl/etcd/client.crt\"" : ""}
Environment="LOCKSMITHD_ENDPOINT=${var.tls_enabled ? "https" : "http"}://etcd-${count.index}:2379"
EOF
    },
  ]
}

data "ignition_systemd_unit" "etcd3" {
  count = "${var.etcd_count}"

  name   = "etcd-member.service"
  enable = true

  dropin = [
    {
      name = "40-etcd-cluster.conf"

      content = <<EOF
[Unit]
Requires=coreos-metadata.service
After=coreos-metadata.service

[Service]
Environment="ETCD_IMAGE=${var.container_image}"
EnvironmentFile=/run/metadata/coreos
Environment="RKT_RUN_ARGS=--volume etcd-ssl,kind=host,source=/etc/ssl/etcd \
  --mount volume=etcd-ssl,target=/etc/ssl/etcd"
ExecStart=
ExecStart=/usr/lib/coreos/etcd-wrapper \
  --name=${element(split(";", var.base_domain == "" ? 
          join(";", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)) : 
          join(";", formatlist("%s.${var.base_domain}", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)))), count.index)} \
  --advertise-client-urls=${var.tls_enabled ? "https" : "http"}://$${COREOS_AZURE_IPV4_DYNAMIC}:2379 \
  ${var.tls_enabled
      ? "--cert-file=/etc/ssl/etcd/server.crt --key-file=/etc/ssl/etcd/server.key --peer-cert-file=/etc/ssl/etcd/peer.crt --peer-key-file=/etc/ssl/etcd/peer.key --peer-trusted-ca-file=/etc/ssl/etcd/ca.crt --peer-client-cert-auth=true"
      : ""} \
  --initial-advertise-peer-urls=${var.tls_enabled ? "https" : "http"}://$${COREOS_AZURE_IPV4_DYNAMIC}:2380 \
  --listen-client-urls=${var.tls_enabled ? "https" : "http"}://0.0.0.0:2379 \
  --listen-peer-urls=${var.tls_enabled ? "https" : "http"}://0.0.0.0:2380 \
  --initial-cluster=${
    join(",",
      formatlist("%s=${var.tls_enabled ? "https" : "http"}://%s:2380", 
        split(";", var.base_domain == "" ? 
          join(";", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)) : 
          join(";", formatlist("%s.${var.base_domain}", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)))), 
        split(";", var.base_domain == "" ? 
          join(";", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)) : 
          join(";", formatlist("%s.${var.base_domain}", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count))))
    ))
  }
EOF

      // The abomination of an expression above will go away once we refactor ignition rendering into it's own module.
    },
  ]
}
