module "bootstrapper" {
  source = "../../modules/bootstrap-ssh"

  _dependencies = [
    "${module.masters.instance_group}",
    "${module.etcd.etcd_ip_addresses}",
    "${module.etcd_certs.id}",
    "${module.bootkube.id}",
    "${module.tectonic.id}",
    "${module.flannel_vxlan.id}",
    "${module.calico.id}",
    "${module.canal.id}",
  ]

  bootstrapping_host = "${module.network.ssh_master_ip}"
}
