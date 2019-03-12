data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

data "vsphere_compute_cluster" "compute_cluster" {
  name          = "${var.vsphere_cluster}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "datastore" {
  name          = "${var.vsphere_datastore}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "network" {
  name          = "${var.vm_network}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_virtual_machine" "template" {
  name          = "${var.vm_template}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

resource "vsphere_virtual_machine" "vm" {
  name             = "bootstrap-${var.cluster_id}"
  resource_pool_id = "${var.resource_pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  num_cpus         = "${var.num_cpus}"
  memory           = "${var.memory}"
  guest_id         = "rhel7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label            = "disk0"
    unit_number      = 0
    size             = 40
    thin_provisioned = false
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"

    customize {
      linux_options {
        host_name = "bootstrap-${var.cluster_id}"
        domain    = "${var.vm_base_domain}"
      }

      network_interface {
        ipv4_address = "${var.bootstrap_ip}"
        ipv4_netmask = "${element(split("/", var.machine_cidr), 1)}"
      }

      ipv4_gateway = "${cidrhost(var.machine_cidr,1)}"
    }
  }
}
