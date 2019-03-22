locals {
  ignition_encoded = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition)}"
}

data "ignition_file" "hostname" {
  count = "${var.instance_count}"

  filesystem = "root"
  path       = "/etc/hostname"
  mode       = "420"

  content {
    content = "${var.name}-${count.index}.${var.cluster_domain}"
  }
}

data "ignition_user" "extra_users" {
  count = "${length(var.extra_user_names)}"

  name          = "${var.extra_user_names[count.index]}"
  password_hash = "${var.extra_user_password_hashes[count.index]}"
}

data "ignition_config" "ign" {
  count = "${var.instance_count}"

  append {
    source = "${var.ignition_url != "" ? var.ignition_url : local.ignition_encoded}"
  }

  files = [
    "${data.ignition_file.hostname.*.id[count.index]}",
  ]

  users = ["${data.ignition_user.extra_users.*.id}"]
}

resource "vsphere_virtual_machine" "vm" {
  count = "${var.instance_count}"

  name             = "${var.name}-${count.index}"
  resource_pool_id = "${var.resource_pool_id}"
  datastore_id     = "${var.datastore_id}"
  num_cpus         = "4"
  memory           = "8192"
  guest_id         = "other26xLinux64Guest"

  wait_for_guest_net_timeout  = 0
  wait_for_guest_net_routable = false

  network_interface {
    network_id = "${var.network_id}"
  }

  disk {
    label = "disk0"
    size  = 60

    // want to change this to thin provisioned. need to change template.
    thin_provisioned = false
  }

  clone {
    template_uuid = "${var.vm_template_id}"
  }

  vapp {
    properties {
      "guestinfo.coreos.config.data" = "${data.ignition_config.ign.*.rendered[count.index]}"
    }
  }
}

/*
			  "networkd": {
			    "units": [
			      {
			        "contents": "[Match]\nName=eth0\n\n[Network]\nDNS=8.8.8.8\nAddress=${var.bootstrap_ip}/${var.network_prefix}\nGateway=${var.gateway}",
			        "name": "00-eth0.network"
			      }
			    ]
			  },
*/
/*
			    {
				  "filesystem":"root",
				  "path":"/etc/sysconfig/network-scripts/ifcfg-eth0",
				  "contents":
				  {
				    "source":"data:,TYPE%3DEthernet%0APROXY_METHOD%3Dnone%0ABROWSER_ONLY%3Dno%0ABOOTPROTO%3Dnone%0ADEFROUTE%3Dyes%0AIPV4_FAILURE_FATAL%3Dno%0AIPV6INIT%3Dyes%0AIPV6_AUTOCONF%3Dyes%0AIPV6_DEFROUTE%3Dyes%0AIPV6_FAILURE_FATAL%3Dno%0AIPV6_ADDR_GEN_MODE%3Dstable-privacy%0ANAME%3Deth0%0AUUID%3Dcc0fcac7-aabd-440a-b0e2-4c98ed3ef8b5%0ADEVICE%3Deth0%0AONBOOT%3Dyes%0AIPADDR%3D${var.bootstrap_ip}%0APREFIX%3D${var.network_prefix}%0AGATEWAY%3D${var.gateway}%0ADNS1%3D8.8.8.8%0A",
					"verification":{}
				  },
				  "mode":420
				}]
*/

