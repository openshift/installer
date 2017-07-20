module "bootkube" {
  source = "../../modules/bootkube"

  cloud_provider        = "azure"
  cloud_provider_config = "${jsonencode(data.null_data_source.cloud-provider.inputs)}"

  cluster_name = "${var.tectonic_cluster_name}"

  kube_apiserver_url = "https://${module.vnet.api_external_fqdn}:443"
  oidc_issuer_url    = "https://${module.vnet.ingress_internal_fqdn}/identity"

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

  etcd_endpoints      = "${module.etcd.node_names}"
  etcd_cert_dns_names = "${module.etcd.node_names}"
  etcd_ca_cert        = "${var.tectonic_etcd_ca_cert_path}"
  etcd_client_cert    = "${var.tectonic_etcd_client_cert_path}"
  etcd_client_key     = "${var.tectonic_etcd_client_key_path}"
  etcd_tls_enabled    = "${var.tectonic_etcd_tls_enabled}"

  experimental_enabled = "${var.tectonic_experimental}"

  master_count = "${var.tectonic_master_count}"
}

module "tectonic" {
  source   = "../../modules/tectonic"
  platform = "azure"

  cluster_name = "${var.tectonic_cluster_name}"

  base_address       = "${module.vnet.ingress_internal_fqdn}"
  kube_apiserver_url = "https://${module.vnet.api_external_fqdn}:443"

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
  ingress_kind      = "NodePort"
  experimental      = "${var.tectonic_experimental}"
  master_count      = "${var.tectonic_master_count}"
  stats_url         = "${var.tectonic_stats_url}"
}

module "flannel-vxlan" {
  source = "../../modules/net/flannel-vxlan"

  flannel_image     = "${var.tectonic_container_images["flannel"]}"
  flannel_cni_image = "${var.tectonic_container_images["flannel_cni"]}"
  cluster_cidr      = "${var.tectonic_cluster_cidr}"

  bootkube_id = "${module.bootkube.id}"
}

module "calico-network-policy" {
  source = "../../modules/net/calico-network-policy"

  kube_apiserver_url = "https://${module.vnet.api_external_fqdn}:443"
  calico_image       = "${var.tectonic_container_images["calico"]}"
  calico_cni_image   = "${var.tectonic_container_images["calico_cni"]}"
  cluster_cidr       = "${var.tectonic_cluster_cidr}"
  enabled            = "${var.tectonic_calico_network_policy}"

  bootkube_id = "${module.bootkube.id}"
}

resource "null_resource" "tectonic" {
  depends_on = ["module.vnet", "module.dns", "module.etcd", "module.masters", "module.bootkube", "module.tectonic", "module.flannel-vxlan", "module.calico-network-policy"]

  triggers {
    api-endpoint = "${module.vnet.api_external_fqdn}"
  }

  connection {
    host  = "${module.vnet.api_external_fqdn}"
    user  = "core"
    agent = true
  }

  provisioner "file" {
    source      = "./generated"
    destination = "$HOME/tectonic"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mkdir -p /opt",
      "sudo rm -rf /opt/tectonic",
      "sudo mv /home/core/tectonic /opt/",
      "sudo systemctl start ${var.tectonic_vanilla_k8s ? "bootkube.service" : "tectonic.service"}",
    ]
  }
}
