module "bootkube" {
  source         = "../../../modules/bootkube"
  cloud_provider = ""

  cluster_name = "${var.tectonic_cluster_name}"

  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  oidc_issuer_url    = "https://${var.tectonic_cluster_name}.${var.tectonic_base_domain}/identity"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"
  versions         = "${var.tectonic_versions}"

  ca_cert    = "${var.tectonic_ca_cert}"
  ca_key     = "${var.tectonic_ca_key}"
  ca_key_alg = "${var.tectonic_ca_key_alg}"

  service_cidr = "${var.tectonic_service_cidr}"
  cluster_cidr = "${var.tectonic_cluster_cidr}"

  advertise_address = "0.0.0.0"
  anonymous_auth    = "false"

  oidc_username_claim = "email"
  oidc_groups_claim   = "groups"
  oidc_client_id      = "tectonic-kubectl"

  etcd_endpoints   = "${module.dns.etcd_a_nodes}"
  etcd_ca_cert     = "${var.tectonic_etcd_ca_cert_path}"
  etcd_client_cert = "${var.tectonic_etcd_client_cert_path}"
  etcd_client_key  = "${var.tectonic_etcd_client_key_path}"
  etcd_tls_enabled = "${var.tectonic_etcd_tls_enabled}"

  etcd_cert_dns_names = [
    "${var.tectonic_cluster_name}-etcd-0.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-1.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-2.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-3.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-4.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-5.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-6.${var.tectonic_base_domain}",
  ]

  experimental_enabled = "${var.tectonic_experimental}"

  master_count = "${var.tectonic_master_count}"
}

module "tectonic" {
  source   = "../../../modules/tectonic"
  platform = "aws"

  cluster_name = "${var.tectonic_cluster_name}"

  base_address       = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"
  versions         = "${var.tectonic_versions}"

  license_path     = "${var.tectonic_vanilla_k8s ? "/dev/null" : pathexpand(var.tectonic_license_path)}"
  pull_secret_path = "${var.tectonic_vanilla_k8s ? "/dev/null" : pathexpand(var.tectonic_pull_secret_path)}"

  admin_email         = "${var.tectonic_admin_email}"
  admin_password_hash = "${var.tectonic_admin_password_hash}"

  update_channel = "${var.tectonic_update_channel}"
  update_app_id  = "${var.tectonic_update_app_id}"
  update_server  = "${var.tectonic_update_server}"

  ca_generated = "${module.bootkube.ca_cert == "" ? false : true}"
  ca_cert      = "${module.bootkube.ca_cert}"
  ca_key_alg   = "${module.bootkube.ca_key_alg}"
  ca_key       = "${module.bootkube.ca_key}"

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

  tls_ca_crt_pem     = "${module.bootkube.etcd_ca_crt_pem}"
  tls_server_crt_pem = "${module.bootkube.etcd_server_crt_pem}"
  tls_server_key_pem = "${module.bootkube.etcd_server_key_pem}"
  tls_client_crt_pem = "${module.bootkube.etcd_client_crt_pem}"
  tls_client_key_pem = "${module.bootkube.etcd_client_key_pem}"
  tls_peer_crt_pem   = "${module.bootkube.etcd_peer_crt_pem}"
  tls_peer_key_pem   = "${module.bootkube.etcd_peer_key_pem}"

  instance_count = "${var.tectonic_etcd_count}"
}

data "null_data_source" "local" {
  inputs = {
    kube_image_url = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
    kube_image_tag = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"
  }
}

module "ignition_masters" {
  source = "../../../modules/ignition"

  container_images    = "${var.tectonic_container_images}"
  image_re            = "${var.tectonic_image_re}"
  kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  kubelet_cni_bin_dir = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label  = "node-role.kubernetes.io/master"
  kubelet_node_taints = "node-role.kubernetes.io/master=:NoSchedule"
}

module "master_nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
${join("\n", formatlist("nameserver %s", var.tectonic_openstack_dns_nameservers))}
EOF

  bootkube_service          = "${module.bootkube.systemd_service}"
  cluster_name              = "${var.tectonic_cluster_name}"
  core_public_keys          = ["${module.secrets.core_public_key_openssh}"]
  hostname_infix            = "master"
  instance_count            = "${var.tectonic_master_count}"
  kubeconfig_content        = "${module.bootkube.kubeconfig}"
  tectonic_service          = "${module.tectonic.systemd_service}"
  tectonic_service_disabled = "${var.tectonic_vanilla_k8s}"

  ign_docker_dropin_id      = "${module.ignition_masters.docker_dropin_id}"
  ign_kubelet_env_id        = "${module.ignition_masters.kubelet_env_id}"
  ign_kubelet_service_id    = "${module.ignition_masters.kubelet_service_id}"
  ign_locksmithd_service_id = "${module.ignition_masters.locksmithd_service_id}"
  ign_max_user_watches_id   = "${module.ignition_masters.max_user_watches_id}"
}

module "ignition_workers" {
  source = "../../../modules/ignition"

  container_images    = "${var.tectonic_container_images}"
  image_re            = "${var.tectonic_image_re}"
  kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  kubelet_cni_bin_dir = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label  = "node-role.kubernetes.io/node"
  kubelet_node_taints = ""
}

module "worker_nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
${join("\n", formatlist("nameserver %s", var.tectonic_openstack_dns_nameservers))}
EOF

  bootkube_service          = ""
  cluster_name              = "${var.tectonic_cluster_name}"
  core_public_keys          = ["${module.secrets.core_public_key_openssh}"]
  hostname_infix            = "worker"
  instance_count            = "${var.tectonic_worker_count}"
  kubeconfig_content        = "${module.bootkube.kubeconfig}"
  tectonic_service          = ""
  tectonic_service_disabled = "${var.tectonic_vanilla_k8s}"

  ign_docker_dropin_id      = "${module.ignition_workers.docker_dropin_id}"
  ign_kubelet_env_id        = "${module.ignition_masters.kubelet_env_id}"
  ign_kubelet_service_id    = "${module.ignition_workers.kubelet_service_id}"
  ign_locksmithd_service_id = "${module.ignition_workers.locksmithd_service_id}"
  ign_max_user_watches_id   = "${module.ignition_workers.max_user_watches_id}"
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

  admin_email = "${var.tectonic_admin_email}"

  api_ips          = "${openstack_networking_floatingip_v2.loadbalancer.*.address}"
  etcd_count       = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count}"
  etcd_ips         = "${openstack_networking_port_v2.etcd.*.all_fixed_ips}"
  etcd_tls_enabled = "${var.tectonic_etcd_tls_enabled}"
  master_count     = "${var.tectonic_master_count}"
  master_ips       = "${openstack_networking_port_v2.master.*.all_fixed_ips}"
  worker_count     = "${var.tectonic_worker_count}"
  worker_ips       = "${openstack_networking_floatingip_v2.worker.*.address}"

  tectonic_experimental = "${var.tectonic_experimental}"
  tectonic_vanilla_k8s  = "${var.tectonic_vanilla_k8s}"
}

module "flannel_vxlan" {
  source = "../../../modules/net/flannel-vxlan"

  flannel_image     = "${var.tectonic_container_images["flannel"]}"
  flannel_cni_image = "${var.tectonic_container_images["flannel_cni"]}"
  cluster_cidr      = "${var.tectonic_cluster_cidr}"

  bootkube_id = "${module.bootkube.id}"
}

module "calico_network_policy" {
  source = "../../../modules/net/calico-network-policy"

  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  calico_image       = "${var.tectonic_container_images["calico"]}"
  calico_cni_image   = "${var.tectonic_container_images["calico_cni"]}"
  cluster_cidr       = "${var.tectonic_cluster_cidr}"
  enabled            = "${var.tectonic_calico_network_policy}"

  bootkube_id = "${module.bootkube.id}"
}
