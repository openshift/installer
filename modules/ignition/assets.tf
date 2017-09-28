data "template_file" "max_user_watches" {
  template = "${file("${path.module}/resources/sysctl.d/max-user-watches.conf")}"
}

data "ignition_file" "max_user_watches" {
  filesystem = "root"
  path       = "/etc/sysctl.d/10-max-user-watches.conf"
  mode       = 0644

  content {
    content = "${data.template_file.max_user_watches.rendered}"
  }
}

data "template_file" "docker_dropin" {
  template = "${file("${path.module}/resources/dropins/10-dockeropts.conf")}"
}

data "ignition_systemd_unit" "docker_dropin" {
  name   = "docker.service"
  enable = true

  dropin = [
    {
      name    = "10-dockeropts.conf"
      content = "${data.template_file.docker_dropin.rendered}"
    },
  ]
}

data "template_file" "kubelet" {
  template = "${file("${path.module}/resources/services/kubelet.service")}"

  vars {
    cloud_provider        = "${var.cloud_provider != "" ? "--cloud-provider=${var.cloud_provider}" : ""}"
    cloud_provider_config = "${var.cloud_provider_config != "" ? "--cloud-config=/etc/kubernetes/cloud/config" : ""}"
    cluster_dns_ip        = "${var.kube_dns_service_ip}"
    cni_bin_dir_flag      = "${var.kubelet_cni_bin_dir != "" ? "--cni-bin-dir=${var.kubelet_cni_bin_dir}" : ""}"
    kubeconfig_fetch_cmd  = "${var.kubeconfig_fetch_cmd != "" ? "ExecStartPre=${var.kubeconfig_fetch_cmd}" : ""}"
    node_label            = "${var.kubelet_node_label}"
    node_taints_param     = "${var.kubelet_node_taints != "" ? "--register-with-taints=${var.kubelet_node_taints}" : ""}"
  }
}

data "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet.rendered}"
}

data "template_file" "k8s_node_bootstrap" {
  template = "${file("${path.module}/resources/services/k8s-node-bootstrap.service")}"

  vars {
    bootstrap_upgrade_cl     = "${var.bootstrap_upgrade_cl}"
    kubeconfig_fetch_cmd     = "${var.kubeconfig_fetch_cmd != "" ? "ExecStartPre=${var.kubeconfig_fetch_cmd}" : ""}"
    tectonic_torcx_image_url = "${replace(var.container_images["tectonic_torcx"],var.image_re,"$1")}"
    tectonic_torcx_image_tag = "${replace(var.container_images["tectonic_torcx"],var.image_re,"$2")}"
    torcx_skip_setup         = "${var.tectonic_vanilla_k8s ? "true" : "false" }"
    torcx_store_url          = "${var.torcx_store_url}"
  }
}

data "ignition_systemd_unit" "k8s_node_bootstrap" {
  name    = "k8s-node-bootstrap.service"
  enable  = true
  content = "${data.template_file.k8s_node_bootstrap.rendered}"
}

data "ignition_systemd_unit" "init_assets" {
  name    = "init-assets.service"
  enable  = "${var.assets_location != "" ? true : false}"
  content = "${file("${path.module}/resources/services/init-assets.service")}"
}

data "template_file" "s3_puller" {
  template = "${file("${path.module}/resources/bin/s3-puller.sh")}"

  vars {
    awscli_image = "${var.container_images["awscli"]}"
  }
}

data "ignition_file" "s3_puller" {
  filesystem = "root"
  path       = "/opt/s3-puller.sh"
  mode       = 0755

  content {
    content = "${data.template_file.s3_puller.rendered}"
  }
}

data "template_file" "gcs_puller" {
  template = "${file("${path.module}/resources/bin/gcs-puller.sh")}"
}

data "ignition_file" "gcs_puller" {
  filesystem = "root"
  path       = "/opt/gcs-puller.sh"
  mode       = 0755

  content {
    content = "${data.template_file.gcs_puller.rendered}"
  }
}

data "ignition_systemd_unit" "locksmithd" {
  name = "locksmithd.service"
  mask = true
}

data "template_file" "installer_kubelet_env" {
  template = "${file("${path.module}/resources/kubernetes/kubelet.env")}"

  vars {
    kubelet_image_url = "${replace(var.container_images["hyperkube"],var.image_re,"$1")}"
    kubelet_image_tag = "${replace(var.container_images["hyperkube"],var.image_re,"$2")}"
  }
}

data "ignition_file" "installer_kubelet_env" {
  filesystem = "root"
  path       = "/etc/kubernetes/installer/kubelet.env"
  mode       = 0644

  content {
    content = "${data.template_file.installer_kubelet_env.rendered}"
  }
}

data "template_file" "tx_off" {
  template = "${file("${path.module}/resources/services/tx-off.service")}"
}

data "ignition_systemd_unit" "tx_off" {
  name    = "tx-off.service"
  enable  = true
  content = "${data.template_file.tx_off.rendered}"
}

data "template_file" "azure_udev_rules" {
  template = "${file("${path.module}/resources/udev/66-azure-storage.rules")}"
}

data "ignition_file" "azure_udev_rules" {
  filesystem = "root"
  path       = "/etc/udev/rules.d/66-azure-storage.rules"
  mode       = 0644

  content {
    content = "${data.template_file.azure_udev_rules.rendered}"
  }
}

data "template_file" "coreos_metadata" {
  template = "${file("${path.module}/resources/services/coreos-metadata.service")}"

  vars {
    metadata_provider = "${var.metadata_provider}"
  }
}

data "ignition_systemd_unit" "coreos_metadata" {
  name   = "coreos-metadata.service"
  enable = true

  dropin = [
    {
      name    = "10-metadata.conf"
      content = "${data.template_file.coreos_metadata.rendered}"
    },
  ]
}
