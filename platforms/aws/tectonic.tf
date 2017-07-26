module "bootkube" {
  source         = "../../modules/bootkube"
  cloud_provider = "aws"

  cluster_name = "${var.tectonic_cluster_name}"

  kube_apiserver_url = "https://${module.masters.api_internal_fqdn}:443"
  oidc_issuer_url    = "https://${module.masters.ingress_internal_fqdn}/identity"

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

  etcd_endpoints   = "${module.etcd.endpoints}"
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

  # The default behavior of Kubernetes's controller manager is to mark a node
  # as Unhealthy after 40s without an update from the node's kubelet. However,
  # AWS ELB's Route53 records have a fixed TTL of 60s. Therefore, when an ELB's
  # node disappears (e.g. scaled down or crashed), kubelet might fail to report
  # for a period of time that exceed the default grace period of 40s and the
  # node might become Unhealthy. While the eviction process won't start until
  # the pod_eviction_timeout is reached, 5min by default, certain operators
  # might already have taken action. This is the case for the etcd operator as
  # of v0.3.3, which removes the likely-healthy etcd pods from the the
  # cluster, potentially leading to a loss-of-quorum as generally all kubelets
  # are affected simultaneously.
  #
  # To cope with this issue, we increase the grace period, and reduce the
  # pod eviction time-out accordingly so pods still get evicted after an total
  # time of 340s after the first post-status failure.
  #
  # Ref: https://github.com/kubernetes/kubernetes/issues/41916
  # Ref: https://github.com/kubernetes-incubator/kube-aws/issues/598
  node_monitor_grace_period = "2m"

  pod_eviction_timeout = "220s"
}

module "tectonic" {
  source   = "../../modules/tectonic"
  platform = "aws"

  cluster_name = "${var.tectonic_cluster_name}"

  base_address       = "${module.masters.ingress_internal_fqdn}"
  kube_apiserver_url = "https://${module.masters.api_internal_fqdn}:443"

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

  image_re = "${var.tectonic_image_re}"
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

  kube_apiserver_url = "https://${module.masters.api_internal_fqdn}:443"
  calico_image       = "${var.tectonic_container_images["calico"]}"
  calico_cni_image   = "${var.tectonic_container_images["calico_cni"]}"
  cluster_cidr       = "${var.tectonic_cluster_cidr}"
  enabled            = "${var.tectonic_calico_network_policy}"

  bootkube_id = "${module.bootkube.id}"
}

data "archive_file" "assets" {
  type       = "zip"
  source_dir = "./generated/"

  # Because the archive_file provider is a data source, depends_on can't be
  # used to guarantee that the tectonic/bootkube modules have generated
  # all the assets on disk before trying to archive them. Instead, we use their
  # ID outputs, that are only computed once the assets have actually been
  # written to disk. We re-hash the IDs (or dedicated module outputs, like module.bootkube.content_hash)
  # to make the filename shorter, since there is no security nor collision risk anyways.
  #
  # Additionally, data sources do not support managing any lifecycle whatsoever,
  # and therefore, the archive is never deleted. To avoid cluttering the module
  # folder, we write it in the TerraForm managed hidden folder `.terraform`.
  output_path = "./.terraform/generated_${sha1("${module.tectonic.id} ${module.bootkube.id} ${module.flannel-vxlan.id} ${module.calico-network-policy.id}")}.zip"
}
