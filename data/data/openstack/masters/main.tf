data "openstack_images_image_v2" "masters_img" {
  name        = "${var.base_image}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "masters_flavor" {
  name = "${var.flavor_name}"
}

data "ignition_config" "master_ignition_config" {
  append {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.user_data_ign)}"
  }

  files = [
    "${data.ignition_file.master_ifcfg.id}",
    "${data.ignition_file.master_hacks_script.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.master_hacks_service.id}",
  ]
}

data "ignition_file" "master_ifcfg" {
    filesystem = "root"
    mode = "420"  // 0644
    path = "/etc/sysconfig/network-scripts/ifcfg-eth0"
    content {
      content = <<EOF
DEVICE="eth0"
BOOTPROTO="dhcp"
ONBOOT="yes"
TYPE="Ethernet"
PERSISTENT_DHCLIENT="yes"
DNS1="${var.service_vm_fixed_ip}"
PEERDNS="no"
NM_CONTROLLED="yes"
EOF
    }
}

data "ignition_file" "master_hacks_script" {
    filesystem = "root"
    mode = "493"  // 0755
    path = "/opt/hacks.sh"
    content {
      content = <<EOF
#!/usr/bin/env bash
set -ex

sed -i '/cloud-provider=openstack/d' /etc/systemd/system/kubelet.service

# NOTE(shadower): this is run before kubelet so we don't need to restart it.
systemctl daemon-reload
EOF
    }
}

data "ignition_systemd_unit" "master_hacks_service" {
    name = "hacks.service"
      content = <<EOF
[Unit]
Description=Run hacks after bootup
Before=kubelet.service

[Service]
Type=oneshot
ExecStart=/opt/hacks.sh
RemainAfterExit=true

[Install]
WantedBy=multi-user.target
EOF
}

resource "openstack_compute_instance_v2" "master_conf" {
  name  = "${var.cluster_name}-master-${count.index}"
  count = "${var.instance_count}"

  flavor_id       = "${data.openstack_compute_flavor_v2.masters_flavor.id}"
  image_id        = "${data.openstack_images_image_v2.masters_img.id}"
  security_groups = ["${var.master_sg_ids}"]
  user_data       = "${data.ignition_config.master_ignition_config.rendered}"

  network = {
    port = "${var.subnet_ids[count.index]}"
  }

  #network = {
  #  name = "openshift"
  #  fixed_ip_v4 = "10.3.0.${count.index == 0? 1: count.index + 2}"
  #}

  metadata {
    Name               = "${var.cluster_name}-master"
    owned              = "kubernetes.io/cluster/${var.cluster_name}"
    tectonicClusterID  = "${var.cluster_id}"
    openshiftClusterID = "${var.cluster_id}"
  }
}
