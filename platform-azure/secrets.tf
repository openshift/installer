resource "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${file(var.tectonic_ssh_key)}",
  ]
}
