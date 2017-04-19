// Install CoreOS to disk
// TODO: Match GUI installer, don't match all machines
resource "matchbox_group" "default" {
  name = "default"
  profile = "${matchbox_profile.coreos-install.name}"
  // No selector, matches all nodes
  metadata {
    coreos_channel = "${var.tectonic_coreos_channel}"
    coreos_version = "${var.tectonic_coreos_version}"
    ignition_endpoint = "${var.tectonic_matchbox_http_endpoint}/ignition"
    baseurl = "${var.tectonic_matchbox_http_endpoint}/assets/coreos"
    ssh_authorized_key = "${var.tectonic_ssh_authorized_key}"
  }
}
