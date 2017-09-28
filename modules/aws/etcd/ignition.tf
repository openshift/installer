data "ignition_config" "etcd" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  systemd = [
    "${data.ignition_systemd_unit.locksmithd.*.id[count.index]}",
    "${var.ign_etcd_dropin_id_list[count.index]}",
    "${data.ignition_systemd_unit.etcd_unzip_tls.id}",
  ]

  files = [
    "${data.ignition_file.node_hostname.*.id[count.index]}",
    "${data.ignition_file.etcd_tls_zip.id}",
  ]
}

data "ignition_file" "node_hostname" {
  count      = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"
  path       = "/etc/hostname"
  mode       = 0644
  filesystem = "root"

  content {
    content = "${var.cluster_name}-etcd-${count.index}.${var.base_domain}"
  }
}

data "ignition_file" "etcd_tls_zip" {
  path       = "/etc/ssl/etcd/tls.zip"
  mode       = 0400
  uid        = 0
  gid        = 0
  filesystem = "root"

  content {
    mime    = "application/octet-stream"
    content = "${var.tls_zip}"
  }
}

data "ignition_systemd_unit" "etcd_unzip_tls" {
  name   = "etcd-unzip-tls.service"
  enable = true

  content = <<EOF
[Unit]
ConditionPathExists=!/etc/ssl/etcd/ca.crt
[Service]
Type=oneshot
WorkingDirectory=/etc/ssl/etcd
ExecStart=/usr/bin/bash -c 'unzip /etc/ssl/etcd/tls.zip && \
chown etcd:etcd /etc/ssl/etcd/peer.* && \
chown etcd:etcd /etc/ssl/etcd/server.* && \
chmod 0400 /etc/ssl/etcd/peer.* /etc/ssl/etcd/server.* /etc/ssl/etcd/client.*'
[Install]
WantedBy=multi-user.target
RequiredBy=etcd-member.service locksmithd.service
EOF
}

data "ignition_systemd_unit" "locksmithd" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

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
Environment="LOCKSMITHD_ENDPOINT=${var.tls_enabled ? "https" : "http"}://${var.cluster_name}-etcd-${count.index}.${var.base_domain}:2379"
EOF
    },
  ]
}
