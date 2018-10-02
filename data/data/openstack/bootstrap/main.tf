resource "openstack_objectstorage_object_v1" "ignition" {
  container_name = "${var.swift_container}"
  name           = "bootstrap.ign"
  content        = "${var.ignition}"
}

resource "openstack_objectstorage_tempurl_v1" "ignition_tmpurl" {
  container = "${var.swift_container}"
  method    = "get"
  object    = "${openstack_objectstorage_object_v1.ignition.name}"
  ttl       = 3600
}

data "ignition_config" "redirect" {
  replace {
    source = "${openstack_objectstorage_tempurl_v1.ignition_tmpurl.url}"
  }
}

data "openstack_images_image_v2" "bootstrap_image" {
  name        = "${var.image_name}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = "${var.flavor_name}"
}

resource "openstack_compute_instance_v2" "bootstrap" {
  name      = "${var.cluster_name}-bootstrap"
  flavor_id = "${data.openstack_compute_flavor_v2.bootstrap_flavor.id}"
  image_id  = "${data.openstack_images_image_v2.bootstrap_image.id}"

  user_data = "${data.ignition_config.redirect.rendered}"

  network {
    port = "${var.bootstrap_port_id}"
  }

  metadata {
    Name = "${var.cluster_name}-bootstrap"

    # "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    tectonicClusterID = "${var.cluster_id}"
  }
}
