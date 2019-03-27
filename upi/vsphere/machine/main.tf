locals {
  mask = "${element(split("/", var.machine_cidr), 1)}"
  gw   = "${cidrhost(var.machine_cidr,1)}"

  ignition_encoded = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition)}"
}

data "vsphere_datastore" "datastore" {
  name          = "${var.datastore}"
  datacenter_id = "${var.datacenter_id}"
}

data "vsphere_network" "network" {
  name          = "${var.network}"
  datacenter_id = "${var.datacenter_id}"
}

data "vsphere_virtual_machine" "template" {
  name          = "${var.template}"
  datacenter_id = "${var.datacenter_id}"
}

data "ignition_file" "hostname" {
  count = "${length(var.ips)}"

  filesystem = "root"
  path       = "/etc/hostname"
  mode       = "420"

  content {
    content = "${var.name}-${count.index}.${var.cluster_domain}"
  }
}

data "ignition_file" "static_ip" {
  count = "${length(var.ips)}"

  filesystem = "root"
  path       = "/etc/sysconfig/network-scripts/ifcfg-eth0"
  mode       = "420"

  content {
    content = <<EOF
TYPE=Ethernet
BOOTPROTO=none
NAME=eth0
DEVICE=eth0
ONBOOT=yes
IPADDR=${var.ips[count.index]}
PREFIX=${local.mask}
GATEWAY=${local.gw}
DNS1=8.8.8.8
EOF
  }
}

data "ignition_systemd_unit" "restart" {
  count = "${length(var.ips)}"

  name = "restart.service"

  content = <<EOF
[Unit]
ConditionFirstBoot=yes
[Service]
Type=idle
ExecStart=/sbin/reboot
[Install]
WantedBy=multi-user.target
EOF
}

data "ignition_user" "extra_users" {
  count = "${length(var.extra_user_names)}"

  name          = "${var.extra_user_names[count.index]}"
  password_hash = "${var.extra_user_password_hashes[count.index]}"
}

data "ignition_config" "ign" {
  count = "${length(var.ips)}"

  append {
    source = "${var.ignition_url != "" ? var.ignition_url : local.ignition_encoded}"
  }

  systemd = [
    "${data.ignition_systemd_unit.restart.*.id[count.index]}",
  ]

  files = [
    "${data.ignition_file.hostname.*.id[count.index]}",
    "${data.ignition_file.static_ip.*.id[count.index]}",
  ]

  users = ["${data.ignition_user.extra_users.*.id}"]
}

resource "vsphere_virtual_machine" "vm" {
  count = "${length(var.ips)}"

  name             = "${var.name}-${count.index}"
  resource_pool_id = "${var.resource_pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  num_cpus         = "4"
  memory           = "8192"
  guest_id         = "other26xLinux64Guest"
  folder           = "${var.folder}"
  enable_disk_uuid = "true"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label            = "disk0"
    size             = 60
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
  }

  vapp {
    properties {
      "guestinfo.ignition.config.data"          = "${base64encode(data.ignition_config.ign.*.rendered[count.index])}"
      "guestinfo.ignition.config.data.encoding" = "base64"
    }
  }
}
