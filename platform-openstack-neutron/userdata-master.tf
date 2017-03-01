data "template_file" "userdata-master" {
  count    = "${var.worker_count}"
  template = "${file("${path.module}/userdata-master.yml")}"

  vars {
    kube_config      = "${base64encode(file("${path.cwd}/assets/auth/kubeconfig"))}"
    tectonic_version = "${var.tectonic_version}"
    etcd_fqdn        = "${element(openstack_compute_instance_v2.etcd_node.*.access_ip_v4, 0)}"
    ca               = "${base64encode(file("${path.cwd}/assets/tls/ca.crt"))}"
    client_crt       = "${base64encode(file("${path.cwd}/assets/tls/kubelet.crt"))}"
    client_crt_key   = "${base64encode(file("${path.cwd}/assets/tls/kubelet.key"))}"
    node_hostname    = "${var.cluster_name}-controller-${count.index}"
    base_domain      = "${var.base_domain}"
  }
}
