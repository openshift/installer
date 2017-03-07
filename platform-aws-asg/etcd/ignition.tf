resource "ignition_config" "etcd" {
  count = "${var.node_count}"

  systemd = [
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.etcd3.*.id[count.index]}",
    "${ignition_systemd_unit.etcd2.id}",
    "${ignition_systemd_unit.etcd.id}",
  ]

  files = [
    "${ignition_file.node_hostname.*.id[count.index]}",
  ]
}

resource "ignition_file" "node_hostname" {
  count      = "${var.node_count}"
  path       = "/etc/hostname"
  mode       = 0644
  filesystem = "root"

  content {
    content = "etcd-${count.index}.${var.base_domain}"
  }
}

resource "ignition_systemd_unit" "locksmithd" {
  name   = "locksmithd.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\n"
    },
  ]
}

resource "ignition_systemd_unit" "etcd3" {
  count  = "${var.node_count}"
  name   = "etcd-member.service"
  enable = true

  dropin = [
    {
      name = "40-etcd-cluster.conf"

      content = <<EOF
[Service]
Environment="ETCD_IMAGE_TAG=v3.1.2"
ExecStart=
ExecStart=/usr/lib/coreos/etcd-wrapper \
  --name=etcd \
  --discovery-srv=${var.base_domain} \
  --advertise-client-urls=http://etcd-${count.index}.${var.base_domain}:2379 \
  --initial-advertise-peer-urls=http://etcd-${count.index}.${var.base_domain}:2380 \
  --listen-client-urls=http://0.0.0.0:2379 \
  --listen-peer-urls=http://0.0.0.0:2380
EOF
    },
  ]
}

resource "ignition_systemd_unit" "etcd2" {
  name   = "etcd2.service"
  enable = false
}

resource "ignition_systemd_unit" "etcd" {
  name   = "etcd.service"
  enable = false
}
