// Install CoreOS to disk
resource "matchbox_group" "coreos_install" {
  count   = "${length(var.tectonic_metal_controller_names) + length(var.tectonic_metal_worker_names)}"
  name    = "${format("coreos-install-%s", element(concat(var.tectonic_metal_controller_names, var.tectonic_metal_worker_names), count.index))}"
  profile = "${matchbox_profile.coreos_install.name}"

  selector {
    mac = "${element(concat(var.tectonic_metal_controller_macs, var.tectonic_metal_worker_macs), count.index)}"
  }

  metadata {
    coreos_channel     = "${var.tectonic_cl_channel}"
    coreos_version     = "${var.tectonic_metal_cl_version}"
    ignition_endpoint  = "${var.tectonic_metal_matchbox_http_url}/ignition"
    baseurl            = "${var.tectonic_metal_matchbox_http_url}/assets/coreos"
    ssh_authorized_key = "${var.tectonic_ssh_authorized_key}"
  }
}

// DO NOT PLACE SECRETS IN USER-DATA

module "ignition_masters" {
  source = "../../modules/ignition"

  container_images    = "${var.tectonic_container_images}"
  image_re            = "${var.tectonic_image_re}"
  kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  kubelet_cni_bin_dir = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label  = "node-role.kubernetes.io/master"
  kubelet_node_taints = "node-role.kubernetes.io/master=:NoSchedule"
}

resource "matchbox_group" "controller" {
  count   = "${length(var.tectonic_metal_controller_names)}"
  name    = "${format("%s-%s", var.tectonic_cluster_name, element(var.tectonic_metal_controller_names, count.index))}"
  profile = "${matchbox_profile.tectonic_controller.name}"

  selector {
    mac = "${element(var.tectonic_metal_controller_macs, count.index)}"
    os  = "installed"
  }

  metadata {
    domain_name        = "${element(var.tectonic_metal_controller_domains, count.index)}"
    ssh_authorized_key = "${var.tectonic_ssh_authorized_key}"
    exclude_tectonic   = "${var.tectonic_vanilla_k8s}"

    etcd_enabled = "${var.tectonic_experimental ? "false" : length(compact(var.tectonic_etcd_servers)) != 0 ? false : "true"}"

    etcd_initial_cluster = "${
      join(",", formatlist(
        var.tectonic_etcd_tls_enabled ? "%s=https://%s:2380" : "%s=http://%s:2380",
        var.tectonic_metal_controller_names,
        var.tectonic_metal_controller_domains
      ))
    }"

    etcd_name        = "${element(var.tectonic_metal_controller_names, count.index)}"
    etcd_scheme      = "${var.tectonic_etcd_tls_enabled ? "https" : "http"}"
    etcd_tls_enabled = "${var.tectonic_etcd_tls_enabled}"

    # extra data
    etcd_image_tag    = "v${var.tectonic_versions["etcd"]}"
    kubelet_image_url = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
    kubelet_image_tag = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"

    ign_bootkube_path_unit_json = "${jsonencode(module.bootkube.systemd_path_unit_rendered)}"
    ign_bootkube_service_json   = "${jsonencode(module.bootkube.systemd_service_rendered)}"
    ign_docker_dropin_json      = "${jsonencode(module.ignition_masters.docker_dropin_rendered)}"
    ign_kubelet_env_json        = "${jsonencode(module.ignition_masters.kubelet_env_rendered)}"
    ign_kubelet_service_json    = "${jsonencode(module.ignition_masters.kubelet_service_rendered)}"
    ign_max_user_watches_json   = "${jsonencode(module.ignition_masters.max_user_watches_rendered)}"
    ign_tectonic_path_unit_json = "${jsonencode(module.tectonic.systemd_path_unit_rendered)}"
    ign_tectonic_service_json   = "${jsonencode(module.tectonic.systemd_service_rendered)}"
  }
}

module "ignition_workers" {
  source = "../../modules/ignition"

  container_images    = "${var.tectonic_container_images}"
  image_re            = "${var.tectonic_image_re}"
  kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  kubelet_cni_bin_dir = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label  = "node-role.kubernetes.io/node"
  kubelet_node_taints = ""
}

resource "matchbox_group" "worker" {
  count   = "${length(var.tectonic_metal_worker_names)}"
  name    = "${format("%s-%s", var.tectonic_cluster_name, element(var.tectonic_metal_worker_names, count.index))}"
  profile = "${matchbox_profile.tectonic_worker.name}"

  selector {
    mac = "${element(var.tectonic_metal_worker_macs, count.index)}"
    os  = "installed"
  }

  metadata {
    domain_name        = "${element(var.tectonic_metal_worker_domains, count.index)}"
    ssh_authorized_key = "${var.tectonic_ssh_authorized_key}"

    # extra data
    kubelet_image_url  = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
    kubelet_image_tag  = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"
    kube_version_image = "${var.tectonic_container_images["kube_version"]}"

    ign_docker_dropin_json       = "${jsonencode(module.ignition_workers.docker_dropin_rendered)}"
    ign_kubelet_env_service_json = "${jsonencode(module.ignition_workers.kubelet_env_service_rendered)}"
    ign_kubelet_service_json     = "${jsonencode(module.ignition_workers.kubelet_service_rendered)}"
    ign_max_user_watches_json    = "${jsonencode(module.ignition_workers.max_user_watches_rendered)}"
  }
}
