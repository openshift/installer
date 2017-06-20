# Tectonic bare-metal Smoke Tests on Packet

This README documents how to run Tectonic smoke tests on bare-metal using [Packet](https://www.packet.net/).
This directory has the following configuration:
```
.
├── fake-creds       # Fake credentials used by the bare-metal smoke tests
├── README.md
├── smoke.sh         # Utility script for planning and creating clusters and running tests
├── packet
│   └── main.tf      # A terraform configuration to launch smoke tests on Packet
└── vars
    └── metal.tfvars # Default Tectonic configuration
```

## Packet
Smoke tests are executed using [libvirt](https://libvirt.org/) on a provisioned Linux machine using the [Packet](https://www.packet.net/) provider.

Given a Packet account the following requirements must be met:
- A Packet API key must be configured under "API Keys" in the Packet Web UI.
- A personal SSH public key must be configured under "SSH Keys" in the Packet Web UI.

## Environment
To begin, verify that the following environment variables are set:
- `TF_VAR_auth_token`: The Packet authentication token as set up in the Packet Web UI under "API Keys".
- `TF_VAR_facility`: The Packet facility (location). Most common values would be `ams1` (Amsterdam), `sjc1` (Sunnyvale, CA), `ewr1` (Parsippany, NJ).
- `TF_VAR_tectonic_pull_secret_path` and `TF_VAR_tectonic_license_path`: The local path to the pull secret and Tectonic license file.
- `TF_VAR_smoke_test`: The actual smoke test to execute pointing to a file in the `vars` directory. If left empty, the smoke test will **not** be executed.
- `TF_VAR_project_id`: The Packet project ID under which the launched smoke test machine will be provisioned.

Example:
```sh
$ cd tests/smoke/bare-metal/packet
$ export TF_VAR_project_id=7451e831-241e-411a-bf96-8610b3e0c522
$ export TF_VAR_auth_token=JViTDmqJZ7uBN7g8IJkTBBvjgSpY4abc
$ export TF_VAR_facility=ams1
$ export TF_VAR_tectonic_pull_secret_path=/path/to/pull-secret
$ export TF_VAR_tectonic_license_path=/path/to/tectonic-license
$ export TF_VAR_smoke_test=metal.tfvars
```

## Create and Test Cluster
To create the Packet machine and test the cluster, execute:
```sh
$ cd tests/smoke/bare-metal/packet
$ terraform apply
packet_device.machine: Creating...
  billing_cycle:    "" => "hourly"
  created:          "" => "<computed>"
  facility:         "" => "ams1"
  hostname:         "" => "tf"
  locked:           "" => "<computed>"
  network.#:        "" => "<computed>"
  operating_system: "" => "ubuntu_16_04"
  plan:             "" => "baremetal_1"
  project_id:       "" => "7451e831-241e-411a-bf96-8610b3e0c522"
  state:            "" => "<computed>"
  updated:          "" => "<computed>"
...
null_resource.tectonic: Creation complete (ID: 1801481014396872433)

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path: 

Outputs:

ip = 147.75.205.245
```

## Sanity test cheatsheet
To be able to ssh into the created Packet host machine, use the above IP address or execute the following command:
```sh
$ terraform output ip
147.75.205.245
$ ssh root@147.75.205.245
```

If `TF_VAR_smoke_test` was left empty, execute the following commands to start the smoke test:
```sh
root@tf-ams1-3b7882ac:~# cd $HOME/go/src/github.com/coreos/tectonic-installer
root@tf-ams1-3b7882ac:~# tests/smoke/bare-metal/smoke.sh vars/metal.tfvars
```

To inspect matchbox, execute:
```sh
root@tf-ams1-3b7882ac:~# journalctl -u dev-matchbox
...
Jun 16 15:04:28 tf-ams1-1ac0a730 rkt[8041]: [  418.439084] matchbox[5]: time="2017-06-16T15:04:28Z" level=info msg="HTTP GET /ignition?uuid=efd58a85-849a-46de-b560-2154b5a4db26&mac=52-54-00-c3-61-77&os=installed"
Jun 16 15:04:28 tf-ams1-1ac0a730 rkt[8041]: [  418.439415] matchbox[5]: time="2017-06-16T15:04:28Z" level=debug msg="Matched an Ignition or Container Linux Config template" group=tf-metal-master-1-node3 labels=map[uuid:efd58a85-849a-46de-b
Jun 16 15:04:28 tf-ams1-1ac0a730 rkt[8041]: [  418.459926] matchbox[5]: time="2017-06-16T15:04:28Z" level=info msg="HTTP HEAD /assets/coreos/1298.7.0/coreos_production_image.bin.bz2"
Jun 16 15:04:28 tf-ams1-1ac0a730 rkt[8041]: [  418.464060] matchbox[5]: time="2017-06-16T15:04:28Z" level=info msg="HTTP HEAD /assets/coreos/1298.7.0/coreos_production_image.bin.bz2.sig"
Jun 16 15:04:29 tf-ams1-1ac0a730 rkt[8041]: [  419.481730] matchbox[5]: time="2017-06-16T15:04:29Z" level=info msg="HTTP GET /assets/coreos/1298.7.0/coreos_production_image.bin.bz2.sig"
Jun 16 15:04:29 tf-ams1-1ac0a730 rkt[8041]: [  419.527638] matchbox[5]: time="2017-06-16T15:04:29Z" level=info msg="HTTP GET /assets/coreos/1298.7.0/coreos_production_image.bin.bz2"
```

To verify the nodes are being launched on the host, execute:
```sh
root@tf-ams1-1ac0a730:~# virsh list
 Id    Name                           State
----------------------------------------------------
 1     node1                          running
 2     node2                          running
 3     node3                          running
 4     node4                          running
```

To inspect kernel messages of the master node, execute:
```sh
root@tf-ams1-3b7882ac:~# virsh console node1
Connected to domain node1
Escape character is ^]
[   12.627528] bridge: filtering via arp/ip/ip6tables is no longer available by default. Update your scripts to load br_netfilter if you need this.
[   12.667720] Bridge firewalling registered
[   12.698805] nf_conntrack version 0.5.0 (16384 buckets, 65536 max)
[   13.263150] Initializing XFRM netlink socket
[   13.367960] IPv6: ADDRCONF(NETDEV_UP): docker0: link is not ready


This is node1.example.com (Linux x86_64 4.9.16-coreos-r1) 15:23:55
SSH host key: SHA256:HhrIZw01NIOXipICDdipZMu3P83Ue9slSUUbD3JGVh8 (ED25519)
SSH host key: SHA256:vyCn86z3gmm3rQXLAYQWIf9kkKipR76i+LjLAUI+3IY (DSA)
SSH host key: SHA256:j468CsleP1XxCayXmFOSHzzMe2i20fe0SATqF4WzBwQ (RSA)
SSH host key: SHA256:YxgQu9TB/DpUBf87Ym7ztZlWtT7AS3RG+wiJWdPrDjQ (ECDSA)
ens3: 172.18.0.21 fe80::5054:ff:fea1:9cae

node1 login: 
```

Note that it may take some time until the nodes are provisioned by matchbox.
Press `Ctrl-]` to exit the above console.

To inspect the current smoke test process on the master node, execute:
```sh
root@tf-ams1-3b7882ac:~# ssh -i $HOME/go/src/github.com/coreos/tectonic-installer/matchbox/tests/smoke/fake_rsa core@node1.example.com
Last login: Fri Jun 16 15:24:10 UTC 2017 from 172.18.0.1 on pts/0
Container Linux by CoreOS stable (1298.7.0)
Update Strategy: No Reboots
core@node1 ~ $ 
```

From within `node1` follow the [troubleshooting guide](../../../Documentation/troubleshooting/troubleshooting.md) for master nodes.
