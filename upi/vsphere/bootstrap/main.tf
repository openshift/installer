data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

data "vsphere_compute_cluster" "compute_cluster" {
  name          = "${var.vsphere_cluster}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "datastore" {
  name          = "${var.vsphere_datastore}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "network" {
  name          = "${var.vm_network}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_virtual_machine" "template" {
  name          = "${var.vm_template}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

resource "vsphere_virtual_machine" "vm" {
  name             = "bootstrap1"
  resource_pool_id = "${var.resource_pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  num_cpus         = "${var.num_cpus}"
  memory           = "${var.memory}"
  guest_id         = "other26xLinux64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label            = "disk0"
    unit_number      = 0
    size             = 40
    thin_provisioned = false
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
  }

  vapp {
    properties {
      "guestinfo.coreos.config.data" = <<EOF
{
  "ignition": {
    "config": {
    },
    "security": {
      "tls": {
      }
    },
    "timeouts": {
    },
    "version": "2.2.0"
  },
  "networkd": {
  },
  "passwd": {
    "users": [
      {
        "groups": [
          "sudo"
        ],
        "name": "core",
        "passwordHash": "$1$sRlwnEsn$T1v1ubQUkyFe2kUzZLMAU.",
        "sshAuthorizedKeys": [
          "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCuY6lBsj95cYTVmzqW2Xci37BJta+ZtNHIee8bFCDskKx/qiH/kbceQGDPpAjREBQBpabhxwd8eUktqdiUI1aQIs7I0Om8XMNUd00Phsz69i8PDzVQyGcLzdn4UPOpS89nFmyto0NJ1V5o4RoR3A7ENbzN7li34g5+zGSBXcVdHFFUErCfOnEtgXw5kltU1byv3+GGKAb+f0CL88BRowp35/NH9sRkP8fzWPS0hl+ofiay5H7xDDd6/nqx6OCd0YHfSEYSTMuRKcM55IFDLLgVNLIumYqwDP6WdCw9dv2aSHPX0bpuLwIEgaMEfGgvTTNi8rZ/rC9Y9OYewfHYOpr9 dphillip@dav1x-m"
        ]
      }
    ]
  },
  "storage": {
    "files": [
      {
        "path": "/etc/hostname",
        "filesystem": "root",
        "mode": 420,
        "contents": {
          "inline": "bootstrap.vmware-demo.acrawford.com"
        }
      },
      {
        "path": "/etc/sysconfig/network-scripts/ifcfg-eth0",
        "filesystem": "root",
        "mode": 420,
        "contents":
        "inline": "TYPE=Ethernet\nPROXY_METHOD=none\nBROWSER_ONLY=no\nBOOTPROTO=none\nDEFROUTE=yes\nIPV4_FAILURE_FATAL=no\nIPV6INIT=yes\nIPV6_AUTOCONF=yes\nIPV6_DEFROUTE=yes\nIPV6_FAILURE_FATAL=no\nIPV6_ADDR_GEN_MODE=stable-privacy\nNAME=eth0\nDEVICE=eth0\nONBOOT=yes\nIPADDR=139.178.89.196\nPREFIX=26\nGATEWAY=139.178.89.193\nDNS1=8.8.8.8\n"
      }
    ]
  },
  "systemd": {
  }
}
EOF
    }
  }
}
