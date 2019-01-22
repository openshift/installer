resource "google_compute_instance" "master" {
  count = "${var.instance_count}"
  name  = "${var.cluster_name}-master-${element(var.zones,count.index)}"

  machine_type   = "${var.instance_type}"
  zone            = "${element(var.zones, count.index)}"

  metadata = {
    user-data = "${var.ignition}"
  }

  tags = ["ocp", "ocp-master"]

  network_interface = {
    network    = "${var.subnetwork != "" ? "" : var.network}"
    subnetwork = "${var.subnetwork}"

    access_config = {    
    }
  }

  boot_disk {
    initialize_params {
      type  = "${var.root_volume_type}"
      size  = "${var.root_volume_size}"
      image = "${var.image_name}"
    }
  }

  labels = "${merge(map(
      "name", "${var.cluster_name}-master-${count.index}",
      "cluster-kubernetes-io", "${var.cluster_name}",
      "clusterid", "${var.cluster_name}"
    ), var.extra_labels)}"
}

# Not ideal, machine API would need to keep membership up to date
resource "google_compute_instance_group" "master-0" {
  name        = "${var.cluster_name}-master-${element(var.zones,count.index)}"
  zone        = "${element(var.zones,count.index)}"
  network     = "${var.network}"

  named_port {
    name = "https"
    port = "6443"
  }

  instances = ["${google_compute_instance.master.*.self_link[0]}"]
}

resource "google_compute_instance_group" "master-1" {
  name        = "${var.cluster_name}-master-${element(var.zones,count.index+1)}"
  zone        = "${element(var.zones,count.index+1)}"
  network     = "${var.network}"

  named_port {
    name = "https"
    port = "6443"
  }

  instances = ["${google_compute_instance.master.*.self_link[(count.index+1) % 3]}"]
}

resource "google_compute_instance_group" "master-2" {
  name        = "${var.cluster_name}-master-${element(var.zones,count.index+2)}"
  zone        = "${element(var.zones,count.index+2)}"
  network     = "${var.network}"

  named_port {
    name = "https"
    port = "6443"
  }

  instances = ["${google_compute_instance.master.*.self_link[(count.index+2) % 3]}"]
}
