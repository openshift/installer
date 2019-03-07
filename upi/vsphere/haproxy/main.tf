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
data "vsphere_network" "private_network" {
  name          = "${var.vm_private_network}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}
data "vsphere_virtual_machine" "template" {
  name          = "${var.vm_template}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}
# first VM created creates the resource pool
resource "vsphere_resource_pool" "resource_pool" {
  name                    = "${var.resource_pool}"
  parent_resource_pool_id = "${data.vsphere_compute_cluster.compute_cluster.resource_pool_id}"
}
resource "vsphere_virtual_machine" "vm" {
  name             = "haproxy-0"
  resource_pool_id = "${vsphere_resource_pool.resource_pool.id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 2
  memory   = 4096
  guest_id = "rhel7_64Guest"

# multihome as gateway and cluster ingress
  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }
  network_interface {
    network_id = "${data.vsphere_network.private_network.id}"
  }
  disk {
   label = "disk0"
   unit_number = 0
   size             = 60
   thin_provisioned = false
  }
  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {
      linux_options {
        host_name = "haproxy-0"
        domain    = "vmware.devcluster.openshift.com"
      }
        network_interface {
          ipv4_address = "${var.public_ipv4}" 
          ipv4_netmask = "${var.public_netmask}" 
      }
        network_interface {
          ipv4_address = "${var.private_ipv4}" 
          ipv4_netmask = "${var.private_netmask}" 
      }
        ipv4_gateway = "${var.public_ipv4_gw}"
    }
  }
}
