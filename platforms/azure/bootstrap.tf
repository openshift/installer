locals {
  "bootstrapping_host" = "${var.tectonic_azure_private_cluster ? 
    module.vnet.master_private_ip_addresses[0] : 
    module.vnet.api_fqdn}"
}

module "bootstrapper" {
  source = "../../modules/bootstrap-ssh"

  # depends_on         = ["module.etcd_certs", "module.vnet", "module.dns", "module.etcd", "module.masters", "module.bootkube", "module.tectonic", "module.flannel-vxlan", "module.calico-network-policy"]
  bootstrapping_host = "${local.bootstrapping_host}"
}
