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
  count            = "${var.worker_count}"
  name             = "worker-${count.index + 1}"
  resource_pool_id = "${var.resource_pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  num_cpus         = "${var.num_cpus}"
  memory           = "${var.memory}"
  guest_id         = "rhel7_64Guest"
  annotation       = "${var.cluster_id}"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label            = "disk0"
    unit_number      = 0
    size             = 60
    thin_provisioned = false
  }

  disk {
    label            = "disk1"
    unit_number      = 1
    size             = 40
    thin_provisioned = false
  }

  disk {
    label            = "disk2"
    unit_number      = 2
    thin_provisioned = false
    size             = 40
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"

    customize {
      linux_options {
        host_name = "worker-${count.index + 1}"
        domain    = "${var.vm_base_domain}"
      }

      network_interface {
        ipv4_address = "${element(var.worker_ips, count.index)}"
        ipv4_netmask = "${element(split("/", var.machine_cidr), 1)}"
      }

      ipv4_gateway = "${cidrhost(var.machine_cidr,1)}"
    }
  }
}
