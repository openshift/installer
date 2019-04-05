locals {
  dns_names           = ["${var.dns_names}"]
  number_of_dns_names = "${length(local.dns_names)}"
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

resource "vsphere_virtual_machine" "vm" {
  count = "${local.ips_exist ? var.instance_count : 0}"

  name             = "${var.name}-${count.index}"
  resource_pool_id = "${var.resource_pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  num_cpus         = "4"
  memory           = "8192"
  guest_id         = "other26xLinux64Guest"
  folder           = "${var.folder}"
  enable_disk_uuid = "true"

  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"

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

  depends_on = ["null_resource.dns_resolution"]
}

resource "null_resource" "dns_resolution" {
  count = "${local.ips_exist && var.wait_for_dns_names ? var.instance_count : 0}"

  provisioner "local-exec" {
    command = <<EOF
export dns_name=${element(concat(local.dns_names, list("")), count.index)}
while [[ "$$(dig +short $${dns_name})" != "${local.ip_addresses[count.index]}" ]]
do
  echo waiting for $${dns_name}
  sleep 5
done
EOF
  }
}
