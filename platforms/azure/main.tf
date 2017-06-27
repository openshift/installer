module "resource_group" {
  source = "../../modules/azure/resource-group"

  external_rsg_name       = "${var.tectonic_azure_external_rsg_name}"
  tectonic_azure_location = "${var.tectonic_azure_location}"
  tectonic_cluster_name   = "${var.tectonic_cluster_name}"
}

module "vnet" {
  source = "../../modules/azure/vnet"

  location              = "${var.tectonic_azure_location}"
  resource_group_name   = "${module.resource_group.name}"
  tectonic_cluster_name = "${var.tectonic_cluster_name}"
  vnet_cidr_block       = "${var.tectonic_azure_vnet_cidr_block}"

  etcd_count                = "${var.tectonic_etcd_count}"
  etcd_cidr                 = "${module.vnet.etcd_cidr}"
  master_cidr               = "${module.vnet.master_cidr}"
  worker_cidr               = "${module.vnet.worker_cidr}"
  external_vnet_name        = "${var.tectonic_azure_external_vnet_name}"
  external_master_subnet_id = "${var.tectonic_azure_external_master_subnet_id}"
  external_worker_subnet_id = "${var.tectonic_azure_external_worker_subnet_id}"
  ssh_network_internal      = "${var.tectonic_azure_ssh_network_internal}"
  ssh_network_external      = "${var.tectonic_azure_ssh_network_external}"
  external_resource_group   = "${var.tectonic_azure_external_resource_group}"
  external_nsg_etcd         = "${var.tectonic_azure_external_nsg_etcd}"
  external_nsg_api          = "${var.tectonic_azure_external_nsg_api}"
  external_nsg_master       = "${var.tectonic_azure_external_nsg_master}"
  external_nsg_worker       = "${var.tectonic_azure_external_nsg_worker}"
}

module "etcd" {
  source = "../../modules/azure/etcd"

  location             = "${var.tectonic_azure_location}"
  resource_group_name  = "${module.resource_group.name}"
  image_reference      = "${var.tectonic_azure_image_reference}"
  vm_size              = "${var.tectonic_azure_etcd_vm_size}"
  storage_account_type = "${var.tectonic_azure_etcd_storage_account_type}"

  etcd_count            = "${var.tectonic_etcd_count}"
  base_domain           = "${var.tectonic_base_domain}"
  cluster_name          = "${var.tectonic_cluster_name}"
  public_ssh_key        = "${var.tectonic_azure_ssh_key}"
  virtual_network       = "${module.vnet.vnet_id}"
  subnet                = "${module.vnet.master_subnet}"
  endpoints             = "${module.vnet.etcd_private_ips}"
  network_interface_ids = "${module.vnet.etcd_network_interface_ids}"
}

module "masters" {
  source = "../../modules/azure/master-as"

  location             = "${var.tectonic_azure_location}"
  resource_group_name  = "${module.resource_group.name}"
  image_reference      = "${var.tectonic_azure_image_reference}"
  vm_size              = "${var.tectonic_azure_master_vm_size}"
  storage_account_type = "${var.tectonic_azure_master_storage_account_type}"

  master_count                 = "${var.tectonic_master_count}"
  base_domain                  = "${var.tectonic_base_domain}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  public_ssh_key               = "${var.tectonic_azure_ssh_key}"
  virtual_network              = "${module.vnet.vnet_id}"
  subnet                       = "${module.vnet.master_subnet}"
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

  use_custom_fqdn = "${var.tectonic_azure_use_custom_fqdn}"
}

module "workers" {
  source = "../../modules/azure/worker-as"

  location             = "${var.tectonic_azure_location}"
  resource_group_name  = "${module.resource_group.name}"
  image_reference      = "${var.tectonic_azure_image_reference}"
  vm_size              = "${var.tectonic_azure_worker_vm_size}"
  storage_account_type = "${var.tectonic_azure_worker_storage_account_type}"

  worker_count                 = "${var.tectonic_worker_count}"
  base_domain                  = "${var.tectonic_base_domain}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  public_ssh_key               = "${var.tectonic_azure_ssh_key}"
  virtual_network              = "${module.vnet.vnet_id}"
  subnet                       = "${module.vnet.worker_subnet}"
  kube_image_url               = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
  kube_image_tag               = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"
  kubeconfig_content           = "${module.bootkube.kubeconfig}"
  tectonic_kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  cloud_provider               = ""
  kubelet_node_label           = "node-role.kubernetes.io/node"
}

module "dns" {
  source = "../../modules/azure/dns"

  master_ip_addresses = "${module.masters.ip_address}"
  console_ip_address  = "${module.masters.console_ip_address}"
  etcd_ip_addresses   = ["${module.vnet.etcd_public_ip}"]

  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${var.tectonic_azure_dns_resource_group}"

  create_dns_zone = "${var.tectonic_azure_create_dns_zone}"

  // TODO etcd list
  // TODO worker list
}
