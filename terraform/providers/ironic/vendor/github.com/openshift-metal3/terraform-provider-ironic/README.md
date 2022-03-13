[![Build Status](https://travis-ci.org/openshift-metal3/terraform-provider-ironic.svg?branch=master)](https://travis-ci.org/openshift-metal3/terraform-provider-ironic) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) ![GitHub release](https://img.shields.io/github/release/openshift-metal3/terraform-provider-ironic.svg)

# Terraform provider for Ironic

This is a terraform provider that lets you provision baremetal servers managed by Ironic.

## Provider

Currently the provider only supports standalone noauth Ironic.  At a
minimum, the Ironic endpoint URL must be specified. The user may also
optionally specify an API microversion.

If you are using Ironic inspector, you may also specify the inspector
URL if you'd like to use the introspection data source.

The timeout option sets the number of seconds that the provider will wait
for the Ironic or Inspector API's to become available. This is useful in cases
where another terraform provider is responsible for bringing up the Ironic
infrastructure.

```terraform
provider "ironic" {
  url          = "http://localhost:6385/v1"
  inspector    = "http://localhost:5050/v1"
  microversion = "1.52"
  timeout      = 900
}
```

## Resources

This provider currently implements a number of native Ironic resources,
described below.

### Nodes

A node describes a hardware resource.  A limited subset of provision
states are supported, you may specify `manage = true` or `available =
true`.  You may also instruct Ironic to inspect (`inspect = true`) or
clean (`clean = true`) the node.  To bring a node to the `active` state,
i.e. deploy the node - use a deployment resource instead.


```terraform
resource "ironic_node_v1" "openshift-master-0" {
  name = "openshift-master-0"

  inspect   = true # Perform inspection
  clean     = true # Clean the node
  available = true # Make the node 'available'

  ports = [
    {
      "address"     = "00:bb:4a:d0:5e:38"
      "pxe_enabled" = "true"
    },
  ]

  properties = {
    "local_gb" = "50"
    "cpu_arch" = "x86_64"
  }

  driver = "ipmi"
  driver_info = {
    "ipmi_port"      = "6230"
    "ipmi_username"  = "admin"
    "ipmi_password"  = "password"
    "ipmi_address"   = "192.168.111.1"
    "deploy_kernel"  = "http://172.22.0.1/images/ironic-python-agent.kernel"
    "deploy_ramdisk" = "http://172.22.0.1/images/ironic-python-agent.initramfs"
  }
}
```

## Ports

Ports may be specified as part of the node resource, or as a separate `ironic_port_v1`
declaration.

```terraform
resource "ironic_port_v1" "openshift-master-0-port-0" {
  node_uuid   = ironic_node_v1.openshift-master-0.id
  pxe_enabled = true
  address     = "00:bb:4a:d0:5e:38"
}
```

## Allocation

The Allocation resource represents a request to find and allocate a Node
for deployment. The microversion must be 1.52 or later.

```terraform
resource "ironic_allocation_v1" "openshift-master-allocation" {
  name  = "master-${count.index}"
  count = 3

  resource_class = "baremetal"

  candidate_nodes = [
    ironic_node_v1.openshift-master-0.id,
    ironic_node_v1.openshift-master-1.id,
    ironic_node_v1.openshift-master-2.id,
  ]

  traits = [
    "CUSTOM_FOO",
  ]
}
```

## Deployment

A deployment will provision a baremetal node, using the information
given in the resource.  The `count` metaparameter can be used to deploy
multiple nodes. Terraform will drive the [Ironic state
machine](https://docs.openstack.org/ironic/latest/contributor/states.html)
to bring the node to an `active` state.  The separation of deployment
from the baremetal node defintion allows a user to manage their hardware
outside of the provider if desired. Destruction of a deployment brings
the baremetal node back to the `available` state.

Users may specify a `node_uuid` directly, or make use of the allocation
resource to dynamically pick a node.


```terraform
resource "ironic_deployment" "masters" {
  count     = 3
  node_uuid = "${element(ironic_allocation_v1.openshift-master-allocation.*.node_uuid, count.index)}"

  instance_info = {
    image_source   = "http://172.22.0.1/images/redhat-coreos-maipo-latest.qcow2"
    image_checksum = "26c53f3beca4e0b02e09d335257826fd"
    capabilities   = "boot_option:local,secure_boot:true"
  }

  user_data    = var.user_data
  network_data = var.network_data
  metadata     = var.metadata
}
```

# Data Sources

## Introspection

When using Ironic inspector, you can use this data source to gather selected information such as network
interface information, number of CPU's and architecture, and memory.

```terraform
data "ironic_introspection" "openshift-master-0" {
  uuid = ironic_node_v1.openshift-master-0.id
}
```

Available data points are:

  - cpu_count
  - cpu_arch
  - memory_mb
  - interfaces, which include:
    - name
    - ip
    - mac

# Development

## Running acceptance tests locally

To run the acceptance tests locally, it's necessary to have a local instance of
ironic and ironic-inspector running.  A similar configuration to that used in
CI can be achieved by running `hack/local_ironic.sh`

# License

Apache 2.0, See LICENSE file
