data "ignition_config" "etcd" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  users = [
    "${data.ignition_user.core.id}",
  ]

  files = [
    "${data.ignition_file.node_hostname.*.id[count.index]}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.locksmithd.*.id[count.index]}",
    "${data.ignition_systemd_unit.etcd3.*.id[count.index]}",
    "${data.ignition_systemd_unit.vmtoolsd_member.id}",
  ]

  networkd = [
    "${data.ignition_networkd_unit.vmnetwork.*.id[count.index]}",
  ]
}

data "ignition_user" "core" {
  name                = "core"
  ssh_authorized_keys = ["${var.core_public_keys}"]
}

data "ignition_systemd_unit" "locksmithd" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  name   = "locksmithd.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\n"
      name    = "40-etcd-lock.conf"

      content = <<EOF
[Service] 
Environment=REBOOT_STRATEGY=etcd-lock 
Environment=LOCKSMITHD_ETCD_CAFILE=/etc/ssl/etcd/ca.crt
Environment=LOCKSMITHD_ETCD_KEYFILE=/etc/ssl/etcd/client.key
Environment=LOCKSMITHD_ETCD_CERTFILE=/etc/ssl/etcd/client.crt
Environment=LOCKSMITHD_ENDPOINT=https://${var.hostname["${count.index}"]}.${var.base_domain}:2379
EOF
    },
  ]
}

data "template_file" "etcd-cluster" {
  template = "${file("${path.module}/resources/etcd-cluster")}"
  count    = "${var.instance_count}"

  vars = {
    etcd-name    = "${var.hostname["${count.index}"]}"
    etcd-address = "${var.hostname["${count.index}"]}.${var.base_domain}"
  }
}

data "ignition_systemd_unit" "etcd3" {
  count  = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"
  name   = "etcd-member.service"
  enable = true

  dropin = [
    {
      name = "40-etcd-cluster.conf"

      content = <<EOF
[Service]
Environment="ETCD_IMAGE=${var.container_image}"
Environment="RKT_RUN_ARGS=--volume etcd-ssl,kind=host,source=/etc/ssl/etcd \
  --mount volume=etcd-ssl,target=/etc/ssl/etcd"
ExecStart=
ExecStart=/usr/lib/coreos/etcd-wrapper \
--name=${var.hostname["${count.index}"]} \
--initial-cluster="${join("," , data.template_file.etcd-cluster.*.rendered)}" \
--advertise-client-urls=https://${var.hostname["${count.index}"]}.${var.base_domain}:2379 \
--cert-file=/etc/ssl/etcd/client.crt \
--key-file=/etc/ssl/etcd/client.key \
--peer-cert-file=/etc/ssl/etcd/peer.crt \
--peer-key-file=/etc/ssl/etcd/peer.key \
--peer-trusted-ca-file=/etc/ssl/etcd/ca.crt \
--peer-client-cert-auth=true \
--initial-advertise-peer-urls=https://${var.hostname["${count.index}"]}.${var.base_domain}:2380 \
--listen-client-urls=https://0.0.0.0:2379 \
--listen-peer-urls=https://0.0.0.0:2380 
EOF
    },
  ]
}

data "ignition_file" "node_hostname" {
  count      = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"
  path       = "/etc/hostname"
  mode       = 0644
  filesystem = "root"

  content {
    content = "${var.hostname["${count.index}"]}"
  }
}

data "ignition_networkd_unit" "vmnetwork" {
  count = "${var.instance_count}"
  name  = "00-ens192.network"

  content = <<EOF
  [Match]
  Name=ens192
  [Network]
  DNS=${var.dns_server}
  Address=${var.ip_address["${count.index}"]}
  Gateway=${var.gateway}
  UseDomains=yes
  Domains=${var.base_domain}
EOF
}

data "ignition_systemd_unit" "vmtoolsd_member" {
  name   = "vmtoolsd.service"
  enable = true

  content = <<EOF
  [Unit]
  Description=VMware Tools Agent
  Documentation=http://open-vm-tools.sourceforge.net/
  ConditionVirtualization=vmware
  [Service]
  ExecStartPre=/usr/bin/ln -sfT /usr/share/oem/vmware-tools /etc/vmware-tools
  ExecStart=/usr/share/oem/bin/vmtoolsd
  TimeoutStopSec=5
EOF
}
