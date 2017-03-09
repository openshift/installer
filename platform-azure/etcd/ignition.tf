resource "ignition_config" "etcd" {
  count = "${var.tectonic_etcd_count}"

  systemd = [
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.etcd3.*.id[count.index]}",
    "${ignition_systemd_unit.etcd2.id}",
    "${ignition_systemd_unit.etcd.id}",
  ]
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
  count  = "${var.tectonic_etcd_count}"
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
Environment="ETCD_IMAGE_TAG=v3.1.2"
EnvironmentFile=/run/metadata/coreos
ExecStart=
ExecStart=/usr/lib/coreos/etcd-wrapper \
  --name=etcd \
  --advertise-client-urls=http://$${COREOS_AZURE_IPV4_DYNAMIC}:2379 \
  --initial-advertise-peer-urls=http://$${COREOS_AZURE_IPV4_DYNAMIC}:2380 \
  --listen-client-urls=http://0.0.0.0:2379 \
  --listen-peer-urls=http://0.0.0.0:2380 \
  --initial-cluster=etcd=http://$${COREOS_AZURE_IPV4_DYNAMIC}:2380
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
