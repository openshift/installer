/*
This is the ignition config for the NCG.
It's currently made available for the NCG as configMap volume.
It is in the assets step to be dropped in the generated folder so the configMap gets deployed
This needs to get consolidated with the NCG ignition chain workflow
TODO: Move ignition generation out of terraform
*/

module "ignition_workers" {
  source = "../../modules/ignition"

  bootstrap_upgrade_cl    = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider          = "aws"
  container_images        = "${var.tectonic_container_images}"
  custom_ca_cert_pem_list = "${var.tectonic_custom_ca_pem_list}"
  etcd_ca_cert_pem        = "${module.etcd_certs.etcd_ca_crt_pem}"
  http_proxy              = "${var.tectonic_http_proxy_address}"
  https_proxy             = "${var.tectonic_https_proxy_address}"
  image_re                = "${var.tectonic_image_re}"
  ingress_ca_cert_pem     = "${module.ingress_certs.ca_cert_pem}"
  iscsi_enabled           = "${var.tectonic_iscsi_enabled}"
  kube_ca_cert_pem        = "${module.kube_certs.ca_cert_pem}"
  kube_dns_service_ip     = "${module.bootkube.kube_dns_service_ip}"
  kubelet_debug_config    = "${var.tectonic_kubelet_debug_config}"
  kubelet_node_label      = "node-role.kubernetes.io/node"
  kubelet_node_taints     = ""
  no_proxy                = "${var.tectonic_no_proxy}"
}

data "ignition_config" "workers" {
  files = ["${compact(list(
    module.ignition_workers.installer_kubelet_env_id,
    module.ignition_workers.installer_runtime_mappings_id,
    module.ignition_workers.max_user_watches_id,
    module.ignition_workers.s3_puller_id,
    module.ignition_workers.profile_env_id,
    module.ignition_workers.systemd_default_env_id,
   ))}",
    "${module.ignition_workers.ca_cert_id_list}",
  ]

  systemd = [
    "${module.ignition_workers.docker_dropin_id}",
    "${module.ignition_workers.k8s_node_bootstrap_service_id}",
    "${module.ignition_workers.kubelet_service_id}",
    "${module.ignition_workers.locksmithd_service_id}",
    "${module.ignition_workers.update_ca_certificates_dropin_id}",
    "${module.ignition_workers.iscsi_service_id}",
  ]
}

module "ignition_masters" {
  source = "../../modules/ignition"

  base_domain               = "${var.tectonic_base_domain}"
  bootstrap_upgrade_cl      = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider            = "aws"
  cluster_name              = "${var.tectonic_cluster_name}"
  container_images          = "${var.tectonic_container_images}"
  custom_ca_cert_pem_list   = "${var.tectonic_custom_ca_pem_list}"
  etcd_advertise_name_list  = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_ca_cert_pem          = "${module.etcd_certs.etcd_ca_crt_pem}"
  etcd_client_crt_pem       = "${module.etcd_certs.etcd_client_crt_pem}"
  etcd_client_key_pem       = "${module.etcd_certs.etcd_client_key_pem}"
  etcd_count                = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_initial_cluster_list = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_peer_crt_pem         = "${module.etcd_certs.etcd_peer_crt_pem}"
  etcd_peer_key_pem         = "${module.etcd_certs.etcd_peer_key_pem}"
  etcd_server_crt_pem       = "${module.etcd_certs.etcd_server_crt_pem}"
  etcd_server_key_pem       = "${module.etcd_certs.etcd_server_key_pem}"
  http_proxy                = "${var.tectonic_http_proxy_address}"
  https_proxy               = "${var.tectonic_https_proxy_address}"
  image_re                  = "${var.tectonic_image_re}"
  ingress_ca_cert_pem       = "${module.ingress_certs.ca_cert_pem}"
  iscsi_enabled             = "${var.tectonic_iscsi_enabled}"
  kube_ca_cert_pem          = "${module.kube_certs.ca_cert_pem}"
  kube_dns_service_ip       = "${module.bootkube.kube_dns_service_ip}"
  kubelet_debug_config      = "${var.tectonic_kubelet_debug_config}"
  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  no_proxy                  = "${var.tectonic_no_proxy}"
}

data "ignition_config" "masters" {
  files = ["${compact(list(
    module.ignition_masters.installer_kubelet_env_id,
    module.ignition_masters.installer_runtime_mappings_id,
    module.ignition_masters.max_user_watches_id,
    module.ignition_masters.profile_env_id,
    module.ignition_masters.systemd_default_env_id,
   ))}",
    "${module.ignition_masters.ca_cert_id_list}",
  ]

  systemd = ["${compact(list(
    module.ignition_masters.docker_dropin_id,
    module.ignition_masters.locksmithd_service_id,
    module.ignition_masters.kubelet_service_id,
    module.ignition_masters.k8s_node_bootstrap_service_id,
    module.ignition_masters.update_ca_certificates_dropin_id,
    module.ignition_masters.iscsi_service_id,
   ))}"]
}

resource "local_file" "ncg" {
  depends_on = ["module.bootkube"]
  content    = "${data.template_file.ncg.rendered}"
  filename   = "./generated/manifests/ncg.yaml"
}

resource "local_file" "ncg_config" {
  depends_on = ["module.bootkube"]
  content    = "${data.template_file.ncg_config.rendered}"
  filename   = "./generated/manifests/ncg-config.yaml"
}

data "template_file" "ncg" {
  template = "${file("${path.module}/resources/ncg/ncg.yaml")}"
}

data "template_file" "ncg_config" {
  template = "${file("${path.module}/resources/ncg/ncg-config.yaml")}"

  vars {
    ncg_config_worker   = "${jsonencode(data.ignition_config.workers.rendered)}"
    ncg_config_master   = "${jsonencode(data.ignition_config.masters.rendered)}"
    kube_dns_service_ip = "${cidrhost(var.tectonic_service_cidr, 10)}"
  }
}
