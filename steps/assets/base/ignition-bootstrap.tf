module "ignition_bootstrap" {
  source = "../../../modules/ignition"

  cloud_provider       = "${var.cloud_provider}"
  container_images     = "${var.tectonic_container_images}"
  etcd_ca_cert_pem     = "${local.etcd_ca_cert_pem}"
  etcd_count           = "${length(data.template_file.etcd_hostname_list.*.id)}"
  image_re             = "${var.tectonic_image_re}"
  ingress_ca_cert_pem  = "${local.ingress_ca_cert_pem}"
  root_ca_cert_pem     = "${local.root_ca_cert_pem}"
  kube_dns_service_ip  = "${module.bootkube.kube_dns_service_ip}"
  kubelet_debug_config = "${var.tectonic_kubelet_debug_config}"
  kubelet_node_label   = "node-role.kubernetes.io/bootstrap"
  kubelet_node_taints  = "node-role.kubernetes.io/bootstrap=:NoSchedule"
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

data "ignition_file" "kubelet_kubeconfig" {
  filesystem = "root"
  path       = "/var/lib/kubelet/kubeconfig"
  mode       = 0644

  content {
    content = "${module.bootkube.kubeconfig-kubelet}"
  }
}

data "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${var.tectonic_admin_ssh_key}",
  ]
}

data "ignition_config" "bootstrap" {
  files = ["${compact(flatten(list(
    list(
      data.ignition_file.kube-system_cluster_config.id,
      data.ignition_file.tectonic_cluster_config.id,
      data.ignition_file.tnco_config.id,
      data.ignition_file.kco_config.id,
      data.ignition_file.bootstrap_kubeconfig.id,
      data.ignition_file.kubelet_kubeconfig.id,
    ),
    module.ignition_bootstrap.etcd_crt_id_list,
    module.ignition_bootstrap.ignition_file_id_list,
    module.bootkube.ignition_file_id_list,
    module.tectonic.ignition_file_id_list,
    local.ca_certs_ignition_file_id_list,
    local.etcd_certs_ignition_file_id_list,
    local.kube_certs_ignition_file_id_list,
    local.tnc_certs_ignition_file_id_list,
    local.service_account_keys_ignition_file_id_list,
   )))}"]

  systemd = ["${compact(flatten(list(
    list(
      module.bootkube.systemd_service_id,
      module.tectonic.systemd_service_id,
    ),
    module.ignition_bootstrap.ignition_systemd_id_list,
   )))}"]

  users = [
    "${data.ignition_user.core.id}",
  ]
}
