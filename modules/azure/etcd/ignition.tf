data "ignition_config" "etcd" {
  count = "${var.etcd_count}"

  systemd = [
    "${data.ignition_systemd_unit.locksmithd.id}",
    "${data.ignition_systemd_unit.etcd3.*.id[count.index]}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]
}

data "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${file(var.public_ssh_key)}",
  ]
}

data "ignition_systemd_unit" "locksmithd" {
  name   = "locksmithd.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\n"
    },
  ]
}

data "ignition_systemd_unit" "etcd3" {
  count  = "${var.etcd_count}"
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
  --name=${azurerm_network_interface.etcd_nic.*.name[count.index]} \
  --advertise-client-urls=http://$${COREOS_AZURE_IPV4_DYNAMIC}:2379 \
  --initial-advertise-peer-urls=http://$${COREOS_AZURE_IPV4_DYNAMIC}:2380 \
  --listen-client-urls=http://0.0.0.0:2379 \
  --listen-peer-urls=http://0.0.0.0:2380 \
  --initial-cluster=${join(",",formatlist("%s=http://%s:2380",azurerm_network_interface.etcd_nic.*.name,azurerm_network_interface.etcd_nic.*.private_ip_address))}
EOF
    },
  ]
}
