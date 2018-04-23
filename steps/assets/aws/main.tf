provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role == "" ? "" : "${var.tectonic_aws_installer_role}"}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

data "aws_availability_zones" "azs" {}

# Terraform doesn't support "inheritance"
# So we have to pass all variables down
module assets_base {
  source = "../base"

  cloud_provider = "aws"
  etcd_count     = "${length(var.tectonic_etcd_servers) > 0 ? length(var.tectonic_etcd_servers) : (var.tectonic_etcd_count > 0 ? var.tectonic_etcd_count : length(data.aws_availability_zones.azs.names) == 5 ? 5 : 3)}"
  ingress_kind   = "NodePort"

  tectonic_base_domain          = "${var.tectonic_base_domain}"
  tectonic_bootstrap_upgrade_cl = "${var.tectonic_bootstrap_upgrade_cl}"
  tectonic_cluster_name         = "${var.tectonic_cluster_name}"
  tectonic_container_images     = "${var.tectonic_container_images}"
  tectonic_custom_ca_pem_list   = "${var.tectonic_custom_ca_pem_list}"
  tectonic_http_proxy_address   = "${var.tectonic_http_proxy_address}"
  tectonic_https_proxy_address  = "${var.tectonic_https_proxy_address}"
  tectonic_no_proxy             = "${var.tectonic_no_proxy}"
  tectonic_image_re             = "${var.tectonic_image_re}"
  tectonic_iscsi_enabled        = "${var.tectonic_iscsi_enabled}"
  tectonic_kubelet_debug_config = "${var.tectonic_kubelet_debug_config}"
  tectonic_etcd_servers         = "${var.tectonic_etcd_servers}"
  tectonic_ca_cert              = "${var.tectonic_ca_cert}"
  tectonic_ca_key_alg           = "${var.tectonic_ca_key_alg}"
  tectonic_ca_key               = "${var.tectonic_ca_key}"
  tectonic_service_cidr         = "${var.tectonic_service_cidr}"
  tectonic_license_path         = "${var.tectonic_license_path}"
  tectonic_pull_secret_path     = "${var.tectonic_pull_secret_path}"
  tectonic_admin_email          = "${var.tectonic_admin_email}"
  tectonic_update_channel       = "${var.tectonic_update_channel}"
  tectonic_platform             = "${var.tectonic_platform}"
  tectonic_versions             = "${var.tectonic_versions}"
  tectonic_admin_password       = "${var.tectonic_admin_password}"
  tectonic_cluster_id           = "${var.tectonic_cluster_id}"
}

# Removing assets is platform-specific
# But it must be installed in /opt/tectonic/rm-assets.sh
data "template_file" "rm_assets_sh" {
  template = "${file("${path.module}/resources/rm-assets.sh")}"

  vars {
    cluster_name       = "${var.tectonic_cluster_name}"
    awscli_image       = "${var.tectonic_container_images["awscli"]}"
    bucket_s3_location = "${var.tectonic_cluster_name}-tnc.${var.tectonic_base_domain}"
  }
}

data "ignition_file" "rm_assets_sh" {
  filesystem = "root"
  path       = "/opt/tectonic/rm-assets.sh"
  mode       = "0700"

  content {
    content = "${data.template_file.rm_assets_sh.rendered}"
  }
}

data "ignition_config" "bootstrap" {
  files = ["${flatten(list(
    list(
      data.ignition_file.rm_assets_sh.id,
    ),
    module.assets_base.ignition_bootstrap_files,
  ))}"]

  systemd = [
    "${module.assets_base.ignition_bootstrap_systemd}",
  ]
}
