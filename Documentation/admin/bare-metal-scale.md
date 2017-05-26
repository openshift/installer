# Scaling Tectonic bare metal clusters

This document describes how to add cluster nodes to Tectonic clusters on bare metal.

## Scaling worker nodes

To scale worker nodes, adjust `tectonic_worker_count` in `terraform.vars` and run:

```
$ terraform apply $ terraform plan \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/metal
```
After running `terraform apply`, [power on][power-on] the machines to PXE boot the new nodes and access the cluster.

## Scaling controller nodes

Adding controller nodes to an existing Tectonic cluster on bare metal is not officially supported. This feature is planned for a future release.

### etcd scaling on bare metal

Each controller in a default bare metal Tectonic cluster runs the etcd service. Tectonic Installer will also accept the client endpoint of an existing external etcd v3 cluster to which it should connect instead.


[matchbox-docs]: https://coreos.com/matchbox/docs/latest
[power-on]: ../install/bare-metal/metal-terraform.md#power-on
