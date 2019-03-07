provider "vsphere" {
  user           = "${var.vsphere_user}"
  password       = "${var.vsphere_password}"
  vsphere_server = "${var.vsphere_server}"
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = "${var.datacenter}"
}
data "vsphere_compute_cluster" "compute_cluster" {
  name          = "${var.cluster}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}
data "vsphere_datastore" "datastore" {
  name          = "${var.datastore}"
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
resource "vsphere_resource_pool" "resource_pool" {
  name                    = "${var.resource_pool}"
  parent_resource_pool_id = "${data.vsphere_compute_cluster.compute_cluster.resource_pool_id}"
}
resource "vsphere_virtual_machine" "vm" {
  count = 4
  name             = "ocp${count.index + 1}"
  resource_pool_id = "${vsphere_resource_pool.resource_pool.id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"

  num_cpus = 4
  memory   = 8192
  guest_id = "rhel7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }
  disk {
   label = "disk0"
   unit_number = 0
   size             = 60
   thin_provisioned = false
  }
  disk {
   label = "disk1"
   unit_number = 1
   size             = 40
   thin_provisioned = false
  }
  disk {
   label = "disk2"
   unit_number = 2
   thin_provisioned = false
   size             = 40
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {
      linux_options {
        host_name = "ocp${count.index + 1}"
        domain    = "vmware.devcluster.openshift.com"
      }
        network_interface {
          ipv4_address = "139.178.89.${130 + count.index}"
          ipv4_netmask = 25
      }
	ipv4_gateway = "139.178.89.129"
    }

  }
}

