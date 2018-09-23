data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.cwd}/assets.tfstate"
  }

  defaults {
    ignition_bootstrap = ""
  }
}

locals {
  ignition_bootstrap = "${var.ignition_bootstrap != "" ? var.ignition_bootstrap : data.terraform_remote_state.assets.ignition_bootstrap}"
}
