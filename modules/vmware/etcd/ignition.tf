data "ignition_config" "etcd" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  users = [
    "${data.ignition_user.core.id}",
  ]

  files = [
    "${data.ignition_file.node_hostname.*.id[count.index]}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.locksmithd.id}",
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
  count = "${length(var.external_endpoints) == 0 ? 1 : 0}"

  name   = "locksmithd.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\n"
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
      Environment="ETCD_IMAGE=docker://${var.container_image}"
      Environment="RKT_RUN_ARGS=--uuid-file-save=/var/lib/coreos/etcd-member-wrapper.uuid --insecure-options=all"
      ExecStart=
      ExecStart=/usr/lib/coreos/etcd-wrapper \
        --name=${var.hostname["${count.index}"]} \
        --advertise-client-urls=http://${var.hostname["${count.index}"]}.${var.base_domain}:2379 \
        --initial-advertise-peer-urls=http://${var.hostname["${count.index}"]}.${var.base_domain}:2380 \
        --listen-client-urls=http://0.0.0.0:2379 \
        --listen-peer-urls=http://0.0.0.0:2380 \
        --initial-cluster="${join("," , data.template_file.etcd-cluster.*.rendered)}" 
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
