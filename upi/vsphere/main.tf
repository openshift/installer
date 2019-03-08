provider "vsphere" {
  user                 = "${var.vsphere_user}"
  password             = "${var.vsphere_password}"
  vsphere_server       = "${var.vsphere_server}"
  allow_unverified_ssl = true
}

module "haproxy" {
  source = "./haproxy"

  cluster            = "${var.cluster}"
  datacenter         = "${var.datacenter}"
  datastore          = "${var.datastore}"
  public_ipv4        = "${var.public_ipv4}"
  public_ipv4_gw     = "${var.public_ipv4_gw}"
  public_netmask     = "${var.public_netmask}"
  private_ipv4       = "${var.private_ipv4}"
  private_ipv4_gw    = "${var.private_ipv4_gw}"
  private_netmask    = "${var.private_netmask}"
  resource_pool      = "${var.resource_pool}"
  vm_network         = "${var.vm_network}"
  vm_private_network = "${var.vm_private_network}"
  vm_template        = "${var.vm_template}"
}

/*
module "bootstrap" {
  source = "./bootstrap"

}
module "masters" {
  source = "./masters"

}
module "nodes" {
  source = "./masters"

}
*/

