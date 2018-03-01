module "ignition_bootstrap" {
  source = "../../modules/ignition"

  assets_location           = "${local.bucket_name}/${local.bucket_assets_key}"
  base_domain               = "${var.tectonic_base_domain}"
  bootstrap_upgrade_cl      = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider            = "aws"
  cluster_name              = "${var.tectonic_cluster_name}"
  container_images          = "${var.tectonic_container_images}"
  custom_ca_cert_pem_list   = "${var.tectonic_custom_ca_pem_list}"
  etcd_advertise_name_list  = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_ca_cert_pem          = "${module.ca_certs.etcd_ca_cert_pem}"
  etcd_client_crt_pem       = "${module.etcd_certs.etcd_client_cert_pem}"
  etcd_client_key_pem       = "${module.etcd_certs.etcd_client_key_pem}"
  etcd_count                = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_initial_cluster_list = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_peer_crt_pem         = "${module.etcd_certs.etcd_peer_cert_pem}"
  etcd_peer_key_pem         = "${module.etcd_certs.etcd_peer_key_pem}"
  etcd_server_crt_pem       = "${module.etcd_certs.etcd_server_cert_pem}"
  etcd_server_key_pem       = "${module.etcd_certs.etcd_server_key_pem}"
  http_proxy                = "${var.tectonic_http_proxy_address}"
  https_proxy               = "${var.tectonic_https_proxy_address}"
  image_re                  = "${var.tectonic_image_re}"
  ingress_ca_cert_pem       = "${module.ingress_certs.ca_cert_pem}"
  iscsi_enabled             = "${var.tectonic_iscsi_enabled}"
  root_ca_cert_pem          = "${module.ca_certs.root_ca_cert_pem}"
  kube_dns_service_ip       = "${module.bootkube.kube_dns_service_ip}"
  kubelet_debug_config      = "${var.tectonic_kubelet_debug_config}"
  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  no_proxy                  = "${var.tectonic_no_proxy}"
}

data "ignition_config" "bootstrap" {
  files = ["${compact(list(
    data.ignition_file.init_assets.id,
    data.ignition_file.rm_assets.id,
    module.ignition_bootstrap.installer_kubelet_env_id,
    module.ignition_bootstrap.installer_runtime_mappings_id,
    module.ignition_bootstrap.max_user_watches_id,
    module.ignition_bootstrap.s3_puller_id,
    module.ignition_bootstrap.profile_env_id,
    module.ignition_bootstrap.systemd_default_env_id,
   ))}",
    "${module.ignition_bootstrap.ca_cert_id_list}",
  ]

  systemd = ["${compact(list(
    module.ignition_bootstrap.docker_dropin_id,
    module.ignition_bootstrap.locksmithd_service_id,
    module.ignition_bootstrap.kubelet_service_id,
    module.ignition_bootstrap.k8s_node_bootstrap_service_id,
    module.ignition_bootstrap.init_assets_service_id,
    module.ignition_bootstrap.rm_assets_service_id,
    module.ignition_bootstrap.rm_assets_path_unit_id,
    module.bootkube.systemd_service_id,
    module.bootkube.systemd_path_unit_id,
    module.tectonic.systemd_service_id,
    module.tectonic.systemd_path_unit_id,
    module.ignition_bootstrap.update_ca_certificates_dropin_id,
    module.ignition_bootstrap.iscsi_service_id,
   ))}"]
}

data "template_file" "init_assets" {
  template = "${file("${path.module}/resources/init-assets.sh")}"

  vars {
    cluster_name       = "${var.tectonic_cluster_name}"
    awscli_image       = "${var.tectonic_container_images["awscli"]}"
    assets_s3_location = "${local.bucket_name}/${local.bucket_assets_key}"
  }
}

data "ignition_file" "init_assets" {
  filesystem = "root"
  path       = "/opt/init-assets.sh"
  mode       = 0755

  content {
    content = "${data.template_file.init_assets.rendered}"
  }
}

data "template_file" "rm_assets" {
  template = "${file("${path.module}/resources/rm-assets.sh")}"

  vars {
    cluster_name       = "${var.tectonic_cluster_name}"
    awscli_image       = "${var.tectonic_container_images["awscli"]}"
    bucket_s3_location = "${local.bucket_name}"
  }
}

data "ignition_file" "rm_assets" {
  filesystem = "root"
  path       = "/opt/rm-assets.sh"
  mode       = 0755

  content {
    content = "${data.template_file.rm_assets.rendered}"
  }
}
