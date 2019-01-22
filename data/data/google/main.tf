provider "google" {
  region = "${var.google_region}"
}

module "bootstrap" {
  source = "./bootstrap"

  image_name   = "${var.google_image_name_override}"
  cluster_name = "${var.cluster_name}"
  ignition     = "${var.ignition_bootstrap}"
  zone         = "${module.network.zones[0]}"
  network      = "${module.network.network}"
  subnetwork   = "${module.network.subnetwork}"

  extra_labels = "${merge(map(
      "name", "${var.cluster_name}-bootstrap",
    ), var.google_extra_labels)}"
}

module "masters" {
  source = "./master"

  image_name       = "${var.google_image_name_override}"
  cluster_name     = "${var.cluster_name}"
  instance_type    = "${var.google_master_instance_type}"
  extra_labels     = "${var.google_extra_labels}"
  instance_count   = "${var.master_count}"
  root_volume_size = "${var.google_master_root_volume_size}"
  root_volume_type = "${var.google_master_root_volume_type}"
  zones            = "${module.network.zones}"
  network          = "${module.network.network}"
  subnetwork       = "${module.network.subnetwork}"
  ignition         = "${var.ignition_master}"
}

module "network" {
  source = "./network"

  cidr_block   = "${var.machine_cidr}"
  cluster_name = "${var.cluster_name}"
  region       = "${var.google_region}"

  bootstrap_instance_group = "${module.bootstrap.bootstrap_instance_group}"
  bootstrap_instance       = "${module.bootstrap.bootstrap_instance}"
  master_instance_groups   = "${module.masters.master_instance_groups}"
  master_instances         = "${module.masters.master_instances}"

  extra_labels = "${var.google_extra_labels}"
}
