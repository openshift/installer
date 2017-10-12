provider "openstack" {
  version = "0.2.1"
}

data "template_file" "etcd_hostname_list" {
  count    = "${var.tectonic_etcd_count}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "kube_certs" {
  source = "../../../modules/tls/kube/self-signed"

  ca_cert_pem        = "${var.tectonic_ca_cert}"
  ca_key_alg         = "${var.tectonic_ca_key_alg}"
  ca_key_pem         = "${var.tectonic_ca_key}"
  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  service_cidr       = "${var.tectonic_service_cidr}"
}

module "etcd_certs" {
  source = "../../../modules/tls/etcd"

  etcd_ca_cert_path     = "${var.tectonic_etcd_ca_cert_path}"
  etcd_cert_dns_names   = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_client_cert_path = "${var.tectonic_etcd_client_cert_path}"
  etcd_client_key_path  = "${var.tectonic_etcd_client_key_path}"
  self_signed           = "${var.tectonic_experimental || var.tectonic_etcd_tls_enabled}"
  service_cidr          = "${var.tectonic_service_cidr}"
}

module "ingress_certs" {
  source = "../../../modules/tls/ingress/self-signed"

  base_address = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  ca_cert_pem  = "${module.kube_certs.ca_cert_pem}"
  ca_key_alg   = "${module.kube_certs.ca_key_alg}"
  ca_key_pem   = "${module.kube_certs.ca_key_pem}"
}

module "identity_certs" {
  source = "../../../modules/tls/identity/self-signed"

  ca_cert_pem = "${module.kube_certs.ca_cert_pem}"
  ca_key_alg  = "${module.kube_certs.ca_key_alg}"
  ca_key_pem  = "${module.kube_certs.ca_key_pem}"
}

module "bootkube" {
  source = "../../../modules/bootkube"

  cloud_provider        = ""
  cloud_provider_config = ""

  cluster_name = "${var.tectonic_cluster_name}"

  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  oidc_issuer_url    = "https://${var.tectonic_cluster_name}.${var.tectonic_base_domain}/identity"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"
  versions         = "${var.tectonic_versions}"

  service_cidr = "${var.tectonic_service_cidr}"
  cluster_cidr = "${var.tectonic_cluster_cidr}"

  advertise_address = "0.0.0.0"
  anonymous_auth    = "false"

  oidc_username_claim = "email"
  oidc_groups_claim   = "groups"
  oidc_client_id      = "tectonic-kubectl"
  oidc_ca_cert        = "${module.ingress_certs.ca_cert_pem}"

  apiserver_cert_pem   = "${module.kube_certs.apiserver_cert_pem}"
  apiserver_key_pem    = "${module.kube_certs.apiserver_key_pem}"
  etcd_ca_cert_pem     = "${module.etcd_certs.etcd_ca_crt_pem}"
  etcd_client_cert_pem = "${module.etcd_certs.etcd_client_crt_pem}"
  etcd_client_key_pem  = "${module.etcd_certs.etcd_client_key_pem}"
  etcd_peer_cert_pem   = "${module.etcd_certs.etcd_peer_crt_pem}"
  etcd_peer_key_pem    = "${module.etcd_certs.etcd_peer_key_pem}"
  etcd_server_cert_pem = "${module.etcd_certs.etcd_server_crt_pem}"
  etcd_server_key_pem  = "${module.etcd_certs.etcd_server_key_pem}"
  kube_ca_cert_pem     = "${module.kube_certs.ca_cert_pem}"
  kubelet_cert_pem     = "${module.kube_certs.kubelet_cert_pem}"
  kubelet_key_pem      = "${module.kube_certs.kubelet_key_pem}"

  etcd_endpoints       = "${module.dns.etcd_a_nodes}"
  experimental_enabled = "${var.tectonic_experimental}"

  master_count = "${var.tectonic_master_count}"

  cloud_config_path = ""
}

module "tectonic" {
  source   = "../../../modules/tectonic"
  platform = "openstack"

  cluster_name = "${var.tectonic_cluster_name}"

  base_address       = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  service_cidr       = "${var.tectonic_service_cidr}"

  # Platform-independent variables wiring, do not modify.
  container_images      = "${var.tectonic_container_images}"
  container_base_images = "${var.tectonic_container_base_images}"
  versions              = "${var.tectonic_versions}"

  license_path     = "${var.tectonic_vanilla_k8s ? "/dev/null" : pathexpand(var.tectonic_license_path)}"
  pull_secret_path = "${var.tectonic_vanilla_k8s ? "/dev/null" : pathexpand(var.tectonic_pull_secret_path)}"

  admin_email    = "${var.tectonic_admin_email}"
  admin_password = "${var.tectonic_admin_password}"

  update_channel = "${var.tectonic_update_channel}"
  update_app_id  = "${var.tectonic_update_app_id}"
  update_server  = "${var.tectonic_update_server}"

  ca_generated = "${var.tectonic_ca_cert == "" ? false : true}"
  ca_cert      = "${module.kube_certs.ca_cert_pem}"

  ingress_ca_cert_pem = "${module.ingress_certs.ca_cert_pem}"
  ingress_cert_pem    = "${module.ingress_certs.cert_pem}"
  ingress_key_pem     = "${module.ingress_certs.key_pem}"

  identity_client_cert_pem = "${module.identity_certs.client_cert_pem}"
  identity_client_key_pem  = "${module.identity_certs.client_key_pem}"
  identity_server_cert_pem = "${module.identity_certs.server_cert_pem}"
  identity_server_key_pem  = "${module.identity_certs.server_key_pem}"

  console_client_id = "tectonic-console"
  kubectl_client_id = "tectonic-kubectl"
  ingress_kind      = "HostPort"
  experimental      = "${var.tectonic_experimental}"
  master_count      = "${var.tectonic_master_count}"
  stats_url         = "${var.tectonic_stats_url}"

  image_re = "${var.tectonic_image_re}"
}

module "etcd" {
  source = "../../../modules/openstack/etcd"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
${join("\n", formatlist("nameserver %s", var.tectonic_openstack_dns_nameservers))}
EOF

  base_domain           = "${var.tectonic_base_domain}"
  cluster_name          = "${var.tectonic_cluster_name}"
  container_image       = "${var.tectonic_container_images["etcd"]}"
  core_public_keys      = ["${module.secrets.core_public_key_openssh}"]
  tectonic_experimental = "${var.tectonic_experimental}"
  tls_enabled           = "${var.tectonic_etcd_tls_enabled}"

  tls_ca_crt_pem     = "${module.etcd_certs.etcd_ca_crt_pem}"
  tls_server_crt_pem = "${module.etcd_certs.etcd_server_crt_pem}"
  tls_server_key_pem = "${module.etcd_certs.etcd_server_key_pem}"
  tls_client_crt_pem = "${module.etcd_certs.etcd_client_crt_pem}"
  tls_client_key_pem = "${module.etcd_certs.etcd_client_key_pem}"
  tls_peer_crt_pem   = "${module.etcd_certs.etcd_peer_crt_pem}"
  tls_peer_key_pem   = "${module.etcd_certs.etcd_peer_key_pem}"

  instance_count = "${var.tectonic_etcd_count}"

  ign_etcd_dropin_id_list       = "${module.ignition_masters.etcd_dropin_id_list}"
  ign_coreos_metadata_dropin_id = "${module.ignition_masters.coreos_metadata_dropin_id}"
}

module "ignition_masters" {
  source = "../../../modules/ignition"

  base_domain               = "${var.tectonic_base_domain}"
  bootstrap_upgrade_cl      = "${var.tectonic_bootstrap_upgrade_cl}"
  cluster_name              = "${var.tectonic_cluster_name}"
  container_images          = "${var.tectonic_container_images}"
  etcd_advertise_name_list  = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_count                = "${var.tectonic_etcd_count}"
  etcd_initial_cluster_list = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_tls_enabled          = "${var.tectonic_etcd_tls_enabled}"
  image_re                  = "${var.tectonic_image_re}"
  kube_dns_service_ip       = "${module.bootkube.kube_dns_service_ip}"
  kubelet_cni_bin_dir       = "${var.tectonic_networking == "calico" || var.tectonic_networking == "canal" ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  metadata_provider         = "openstack-metadata"
  tectonic_vanilla_k8s      = "${var.tectonic_vanilla_k8s}"
}

module "master_nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
${join("\n", formatlist("nameserver %s", var.tectonic_openstack_dns_nameservers))}
EOF

  cluster_name       = "${var.tectonic_cluster_name}"
  core_public_keys   = ["${module.secrets.core_public_key_openssh}"]
  hostname_infix     = "master"
  instance_count     = "${var.tectonic_master_count}"
  kubeconfig_content = "${module.bootkube.kubeconfig}"

  ign_bootkube_path_unit_id         = "${module.bootkube.systemd_path_unit_id}"
  ign_bootkube_service_id           = "${module.bootkube.systemd_service_id}"
  ign_docker_dropin_id              = "${module.ignition_masters.docker_dropin_id}"
  ign_installer_kubelet_env_id      = "${module.ignition_masters.installer_kubelet_env_id}"
  ign_k8s_node_bootstrap_service_id = "${module.ignition_masters.k8s_node_bootstrap_service_id}"
  ign_kubelet_service_id            = "${module.ignition_masters.kubelet_service_id}"
  ign_locksmithd_service_id         = "${module.ignition_masters.locksmithd_service_id}"
  ign_max_user_watches_id           = "${module.ignition_masters.max_user_watches_id}"
  ign_tectonic_path_unit_id         = "${var.tectonic_vanilla_k8s ? "" : module.tectonic.systemd_path_unit_id}"
  ign_tectonic_service_id           = "${module.tectonic.systemd_service_id}"
}

module "ignition_workers" {
  source = "../../../modules/ignition"

  bootstrap_upgrade_cl = "${var.tectonic_bootstrap_upgrade_cl}"
  container_images     = "${var.tectonic_container_images}"
  image_re             = "${var.tectonic_image_re}"
  kube_dns_service_ip  = "${module.bootkube.kube_dns_service_ip}"
  kubelet_cni_bin_dir  = "${var.tectonic_networking == "calico" || var.tectonic_networking == "canal" ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label   = "node-role.kubernetes.io/node"
  kubelet_node_taints  = ""
  tectonic_vanilla_k8s = "${var.tectonic_vanilla_k8s}"
}

module "worker_nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
${join("\n", formatlist("nameserver %s", var.tectonic_openstack_dns_nameservers))}
EOF

  cluster_name       = "${var.tectonic_cluster_name}"
  core_public_keys   = ["${module.secrets.core_public_key_openssh}"]
  hostname_infix     = "worker"
  instance_count     = "${var.tectonic_worker_count}"
  kubeconfig_content = "${module.bootkube.kubeconfig}"

  ign_docker_dropin_id              = "${module.ignition_workers.docker_dropin_id}"
  ign_installer_kubelet_env_id      = "${module.ignition_workers.installer_kubelet_env_id}"
  ign_k8s_node_bootstrap_service_id = "${module.ignition_workers.k8s_node_bootstrap_service_id}"
  ign_kubelet_service_id            = "${module.ignition_workers.kubelet_service_id}"
  ign_locksmithd_service_id         = "${module.ignition_workers.locksmithd_service_id}"
  ign_max_user_watches_id           = "${module.ignition_workers.max_user_watches_id}"
}

module "secrets" {
  source       = "../../../modules/openstack/secrets"
  cluster_name = "${var.tectonic_cluster_name}"
}

module "secgroups" {
  source                = "../../../modules/openstack/secgroups"
  cluster_name          = "${var.tectonic_cluster_name}"
  cluster_cidr          = "${var.tectonic_openstack_subnet_cidr}"
  tectonic_experimental = "${var.tectonic_experimental}"
}

module "dns" {
  source = "../../../modules/dns/designate"

  cluster_name = "${var.tectonic_cluster_name}"
  base_domain  = "${var.tectonic_base_domain}"

  admin_email               = "${var.tectonic_admin_email}"
  api_ips                   = "${openstack_networking_floatingip_v2.loadbalancer.*.address}"
  etcd_count                = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count}"
  etcd_ips                  = "${flatten(openstack_networking_port_v2.etcd.*.all_fixed_ips)}"
  etcd_tls_enabled          = "${var.tectonic_etcd_tls_enabled}"
  master_count              = "${var.tectonic_master_count}"
  master_ips                = "${flatten(openstack_networking_port_v2.master.*.all_fixed_ips)}"
  worker_count              = "${var.tectonic_worker_count}"
  worker_ips                = "${flatten(openstack_networking_port_v2.worker.*.all_fixed_ips)}"
  worker_public_ips         = "${openstack_networking_floatingip_v2.worker.*.address}"
  worker_public_ips_enabled = "${var.tectonic_openstack_disable_floatingip ? false : true}"

  tectonic_experimental = "${var.tectonic_experimental}"
  tectonic_vanilla_k8s  = "${var.tectonic_vanilla_k8s}"
}

module "flannel-vxlan" {
  source = "../../../modules/net/flannel-vxlan"

  cluster_cidr     = "${var.tectonic_cluster_cidr}"
  enabled          = "${var.tectonic_networking == "flannel"}"
  container_images = "${var.tectonic_container_images}"
}

module "calico" {
  source = "../../../modules/net/calico"

  container_images = "${var.tectonic_container_images}"
  cluster_cidr     = "${var.tectonic_cluster_cidr}"
  enabled          = "${var.tectonic_networking == "calico"}"
}

module "canal" {
  source = "../../../modules/net/canal"

  container_images = "${var.tectonic_container_images}"
  cluster_cidr     = "${var.tectonic_cluster_cidr}"
  enabled          = "${var.tectonic_networking == "canal"}"
}
