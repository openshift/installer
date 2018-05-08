module "ignition_bootstrap" {
  source = "../../../modules/ignition"

  base_domain               = "${var.tectonic_base_domain}"
  bootstrap_upgrade_cl      = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider            = "${var.cloud_provider}"
  cluster_name              = "${var.tectonic_cluster_name}"
  container_images          = "${var.tectonic_container_images}"
  custom_ca_cert_pem_list   = "${var.tectonic_custom_ca_pem_list}"
  etcd_advertise_name_list  = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_ca_cert_pem          = "${module.ca_certs.etcd_ca_cert_pem}"
  etcd_client_crt_pem       = "${module.etcd_certs.etcd_client_cert_pem}"
  etcd_client_key_pem       = "${module.etcd_certs.etcd_client_key_pem}"
  etcd_count                = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_initial_cluster_list = "${data.template_file.etcd_hostname_list.*.rendered}"
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

# The cluster configs written by the install binary external to Terraform.
# Read them in so we can install them via ignition
data "ignition_file" "kube-system_cluster_config" {
  filesystem = "root"
  mode       = "0644"
  path       = "/opt/tectonic/manifests/cluster-config.yaml"

  content {
    content = "${file("./generated/manifests/cluster-config.yaml")}"
  }
}

data "ignition_file" "tectonic_cluster_config" {
  filesystem = "root"
  mode       = "0644"
  path       = "/opt/tectonic/tectonic/cluster-config.yaml"

  content {
    content = "${file("./generated/tectonic/cluster-config.yaml")}"
  }
}

data "ignition_file" "tnco_config" {
  filesystem = "root"
  mode       = "0644"
  path       = "/opt/tectonic/tnco-config.yaml"

  content {
    content = "${file("./generated/tnco-config.yaml")}"
  }
}

data "ignition_file" "kco_config" {
  filesystem = "root"
  mode       = "0644"
  path       = "/opt/tectonic/kco-config.yaml"

  content {
    content = "${file("./generated/kco-config.yaml")}"
  }
}

data "ignition_file" "bootstrap_kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = 0644

  content {
    content = "${module.bootkube.kubeconfig-kubelet}"
  }
}
