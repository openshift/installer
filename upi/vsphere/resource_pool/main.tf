data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

data "vsphere_compute_cluster" "compute_cluster" {
  name          = "${var.vsphere_cluster}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

resource "vsphere_resource_pool" "resource_pool" {
  name                    = "${var.vsphere_resource_pool}"
  parent_resource_pool_id = "${data.vsphere_compute_cluster.compute_cluster.resource_pool_id}"
}
