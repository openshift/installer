module "resource_group" {
  source = "../../modules/azure/resource-group"

  external_rsg_id = "${var.tectonic_azure_external_resource_group}"
  azure_location  = "${var.tectonic_azure_location}"
  cluster_name    = "${var.tectonic_cluster_name}"
}

module "vnet" {
  source = "../../modules/azure/vnet"

  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${module.resource_group.name}"
  cluster_name        = "${var.tectonic_cluster_name}"
  base_domain         = "${var.tectonic_base_domain}"
  vnet_cidr_block     = "${var.tectonic_azure_vnet_cidr_block}"

  etcd_count           = "${var.tectonic_experimental ? 0 : max(var.tectonic_etcd_count, 1)}"
  master_count         = "${var.tectonic_master_count}"
  worker_count         = "${var.tectonic_worker_count}"
  etcd_cidr            = "${module.vnet.etcd_cidr}"
  master_cidr          = "${module.vnet.master_cidr}"
  worker_cidr          = "${module.vnet.worker_cidr}"
  ssh_network_internal = "${var.tectonic_azure_ssh_network_internal}"
  ssh_network_external = "${var.tectonic_azure_ssh_network_external}"

  external_vnet_id          = "${var.tectonic_azure_external_vnet_id}"
  external_master_subnet_id = "${var.tectonic_azure_external_master_subnet_id}"
  external_worker_subnet_id = "${var.tectonic_azure_external_worker_subnet_id}"
  external_nsg_etcd_id      = "${var.tectonic_azure_external_nsg_etcd_id}"
  external_nsg_api_id       = "${var.tectonic_azure_external_nsg_api_id}"
  external_nsg_master_id    = "${var.tectonic_azure_external_nsg_master_id}"
  external_nsg_worker_id    = "${var.tectonic_azure_external_nsg_worker_id}"
}

module "etcd" {
  source = "../../modules/azure/etcd"

  location             = "${var.tectonic_azure_location}"
  resource_group_name  = "${module.resource_group.name}"
  vm_size              = "${var.tectonic_azure_etcd_vm_size}"
  storage_account_type = "${var.tectonic_azure_etcd_storage_account_type}"
  container_image      = "${var.tectonic_container_images["etcd"]}"

  etcd_count            = "${var.tectonic_experimental ? 0 : max(var.tectonic_etcd_count, 1)}"
  base_domain           = "${var.tectonic_base_domain}"
  cluster_name          = "${var.tectonic_cluster_name}"
  public_ssh_key        = "${var.tectonic_azure_ssh_key}"
  network_interface_ids = "${module.vnet.etcd_network_interface_ids}"
  versions              = "${var.tectonic_versions}"
  cl_channel            = "${var.tectonic_cl_channel}"

  tls_enabled        = "${var.tectonic_etcd_tls_enabled}"
  tls_ca_crt_pem     = "${module.bootkube.etcd_ca_crt_pem}"
  tls_server_crt_pem = "${module.bootkube.etcd_server_crt_pem}"
  tls_server_key_pem = "${module.bootkube.etcd_server_key_pem}"
  tls_client_crt_pem = "${module.bootkube.etcd_client_crt_pem}"
  tls_client_key_pem = "${module.bootkube.etcd_client_key_pem}"
  tls_peer_crt_pem   = "${module.bootkube.etcd_peer_crt_pem}"
  tls_peer_key_pem   = "${module.bootkube.etcd_peer_key_pem}"
}

module "masters" {
  source = "../../modules/azure/master-as"

  location             = "${var.tectonic_azure_location}"
  resource_group_name  = "${module.resource_group.name}"
  vm_size              = "${var.tectonic_azure_master_vm_size}"
  storage_account_type = "${var.tectonic_azure_master_storage_account_type}"

  master_count                 = "${var.tectonic_master_count}"
  base_domain                  = "${var.tectonic_base_domain}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  public_ssh_key               = "${var.tectonic_azure_ssh_key}"
  virtual_network              = "${module.vnet.vnet_id}"
  network_interface_ids        = "${module.vnet.master_network_interface_ids}"
  kube_image_url               = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
  kube_image_tag               = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"
  kubeconfig_content           = "${module.bootkube.kubeconfig}"
  tectonic_kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  cloud_provider               = ""
  kubelet_node_label           = "node-role.kubernetes.io/master"
  kubelet_node_taints          = "node-role.kubernetes.io/master=:NoSchedule"
  bootkube_service             = "${module.bootkube.systemd_service}"
  tectonic_service             = "${module.tectonic.systemd_service}"
  tectonic_service_disabled    = "${var.tectonic_vanilla_k8s}"
  versions                     = "${var.tectonic_versions}"
  cl_channel                   = "${var.tectonic_cl_channel}"
}

module "workers" {
  source = "../../modules/azure/worker-as"

  location             = "${var.tectonic_azure_location}"
  resource_group_name  = "${module.resource_group.name}"
  vm_size              = "${var.tectonic_azure_worker_vm_size}"
  storage_account_type = "${var.tectonic_azure_worker_storage_account_type}"

  worker_count                 = "${var.tectonic_worker_count}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  public_ssh_key               = "${var.tectonic_azure_ssh_key}"
  virtual_network              = "${module.vnet.vnet_id}"
  network_interface_ids        = "${module.vnet.worker_network_interface_ids}"
  kube_image_url               = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
  kube_image_tag               = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"
  kubeconfig_content           = "${module.bootkube.kubeconfig}"
  tectonic_kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  cloud_provider               = ""
  kubelet_node_label           = "node-role.kubernetes.io/node"
  versions                     = "${var.tectonic_versions}"
  cl_channel                   = "${var.tectonic_cl_channel}"
}

module "dns" {
  source = "../../modules/dns/azure"

  etcd_count   = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count}"
  master_count = "${var.tectonic_master_count}"
  worker_count = "${var.tectonic_worker_count}"

  etcd_ip_addresses    = "${module.vnet.etcd_endpoints}"
  master_ip_addresses  = "${module.vnet.master_private_ip_addresses}"
  worker_ip_addresses  = "${module.vnet.worker_private_ip_addresses}"
  api_ip_addresses     = "${module.vnet.api_ip_addresses}"
  console_ip_addresses = "${module.vnet.console_ip_addresses}"

  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  location             = "${var.tectonic_azure_location}"
  external_dns_zone_id = "${var.tectonic_azure_external_dns_zone_id}"
}
