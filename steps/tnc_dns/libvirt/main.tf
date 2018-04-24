# Sets up the libvirt domain name
resource "null_resource" "tnc_dns" {
  provisioner "local-exec" {
    command = "virsh -c qemu:///system net-update ${var.tectonic_libvirt_network_name} add dns-host \"<host ip='${var.tectonic_libvirt_master_ips[0]}'><hostname>${var.tectonic_cluster_name}-api</hostname><hostname>${var.tectonic_cluster_name}-tnc</hostname></host>\" --live --config"
  }
}
