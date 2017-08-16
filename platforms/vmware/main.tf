module "etcd" {
  source         = "../../modules/vmware/etcd"
  instance_count = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count }"

  cluster_name       = "${var.tectonic_cluster_name}"
  core_public_keys   = ["${var.tectonic_vmware_ssh_authorized_key}"]
  container_image    = "${var.tectonic_container_images["etcd"]}"
  base_domain        = "${var.tectonic_base_domain}"
  external_endpoints = ["${compact(var.tectonic_etcd_servers)}"]

  tls_ca_crt_pem     = "${module.bootkube.etcd_ca_crt_pem}"
  tls_server_crt_pem = "${module.bootkube.etcd_server_crt_pem}"
  tls_server_key_pem = "${module.bootkube.etcd_server_key_pem}"
  tls_client_crt_pem = "${module.bootkube.etcd_client_crt_pem}"
  tls_client_key_pem = "${module.bootkube.etcd_client_key_pem}"
  tls_peer_key_pem   = "${module.bootkube.etcd_peer_key_pem}"
  tls_peer_crt_pem   = "${module.bootkube.etcd_peer_crt_pem}"

  hostname   = "${var.tectonic_vmware_etcd_hostnames}"
  dns_server = "${var.tectonic_vmware_node_dns}"
  ip_address = "${var.tectonic_vmware_etcd_ip}"
  gateway    = "${var.tectonic_vmware_etcd_gateway}"

  vmware_datacenter       = "${var.tectonic_vmware_datacenter}"
  vmware_cluster          = "${var.tectonic_vmware_cluster}"
  vm_vcpu                 = "${var.tectonic_vmware_etcd_vcpu}"
  vm_memory               = "${var.tectonic_vmware_etcd_memory}"
  vm_network_label        = "${var.tectonic_vmware_network}"
  vm_disk_datastore       = "${var.tectonic_vmware_etcd_datastore}"
  vm_disk_template        = "${var.tectonic_vmware_vm_template}"
  vm_disk_template_folder = "${var.tectonic_vmware_vm_template_folder}"
  vmware_folder           = "${vsphere_folder.tectonic_vsphere_folder.path}"
}

module "masters" {
  source           = "../../modules/vmware/node"
  instance_count   = "${var.tectonic_master_count}"
  base_domain      = "${var.tectonic_base_domain}"
  core_public_keys = ["${var.tectonic_vmware_ssh_authorized_key}"]
  hostname         = "${var.tectonic_vmware_master_hostnames}"
  dns_server       = "${var.tectonic_vmware_node_dns}"
  ip_address       = "${var.tectonic_vmware_master_ip}"
  gateway          = "${var.tectonic_vmware_master_gateway}"

  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  kubelet_cni_bin_dir       = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kube_dns_service_ip       = "${module.bootkube.kube_dns_service_ip}"
  container_images          = "${var.tectonic_container_images}"
  bootkube_service          = "${module.bootkube.systemd_service}"
  tectonic_service          = "${module.tectonic.systemd_service}"
  tectonic_service_disabled = "${var.tectonic_vanilla_k8s}"
  kube_image_url            = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
  kube_image_tag            = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"

  vmware_datacenter       = "${var.tectonic_vmware_datacenter}"
  vmware_cluster          = "${var.tectonic_vmware_cluster}"
  vm_vcpu                 = "${var.tectonic_vmware_master_vcpu}"
  vm_memory               = "${var.tectonic_vmware_master_memory}"
  vm_network_label        = "${var.tectonic_vmware_network}"
  vm_disk_datastore       = "${var.tectonic_vmware_master_datastore}"
  vm_disk_template        = "${var.tectonic_vmware_vm_template}"
  vm_disk_template_folder = "${var.tectonic_vmware_vm_template_folder}"
  vmware_folder           = "${vsphere_folder.tectonic_vsphere_folder.path}"
  kubeconfig              = "${module.bootkube.kubeconfig}"
  private_key             = "${var.tectonic_vmware_ssh_private_key_path}"
  image_re                = "${var.tectonic_image_re}"
}

module "workers" {
  source           = "../../modules/vmware/node"
  instance_count   = "${var.tectonic_worker_count}"
  base_domain      = "${var.tectonic_base_domain}"
  core_public_keys = ["${var.tectonic_vmware_ssh_authorized_key}"]
  hostname         = "${var.tectonic_vmware_worker_hostnames}"
  dns_server       = "${var.tectonic_vmware_node_dns}"
  ip_address       = "${var.tectonic_vmware_worker_ip}"
  gateway          = "${var.tectonic_vmware_worker_gateway}"

  kubelet_node_label  = "node-role.kubernetes.io/node"
  kubelet_node_taints = ""
  kubelet_cni_bin_dir = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  container_images    = "${var.tectonic_container_images}"
  bootkube_service    = ""
  tectonic_service    = ""
  kube_image_url      = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
  kube_image_tag      = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"

  vmware_datacenter       = "${var.tectonic_vmware_datacenter}"
  vmware_cluster          = "${var.tectonic_vmware_cluster}"
  vm_vcpu                 = "${var.tectonic_vmware_worker_vcpu}"
  vm_memory               = "${var.tectonic_vmware_worker_memory}"
  vm_network_label        = "${var.tectonic_vmware_network}"
  vm_disk_datastore       = "${var.tectonic_vmware_worker_datastore}"
  vm_disk_template        = "${var.tectonic_vmware_vm_template}"
  vm_disk_template_folder = "${var.tectonic_vmware_vm_template_folder}"
  vmware_folder           = "${vsphere_folder.tectonic_vsphere_folder.path}"
  kubeconfig              = "${module.bootkube.kubeconfig}"
  private_key             = "${var.tectonic_vmware_ssh_private_key_path}"
  image_re                = "${var.tectonic_image_re}"
}
