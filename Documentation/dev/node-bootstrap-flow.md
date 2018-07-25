# Node bootstrapping flow

This is a development document which describes the bootstrapping flow for ContainerLinux nodes provisioned by the openshift-installer as part of a OpenShift cluster.

## Overview

When a cluster node is being bootstrapped from scratch, it goes through several phases in the following order:

1. first-boot OS configuration, via ignition (systemd units, node configuration, etc)
2. provisioning of additional assets (k8s manifests, TLS material), via either of:
   * pushing from terraform file/remote-exec (SSH)
   * pulling from private cloud stores (S3 buckets)
3. if needed, a node reboot is triggered to apply systemd-wide changes and to clean container runtime datadir

Additionally, only on one of the master nodes the following kubernetes bootstrapping happens:

1. `bootkube.service` is started after `kubelet.service` start
2. a static bootstrapping control-plane is deployed
3. a fully self-hosted control-plane starts and takes over the previous one
4. `bootkube.service` is completed with success
5. `tectonic.service` is started
6. a self-hosted tectonic control-plane is deployed
7. `tectonic.service` is completed with success

## Systemd units

The following systemd unit is deployed to a node by openshift-installer and take part in the bootstrapping process:

* `kubelet.service` is the main kubelet daemon. It is automatically started on boot.

Additionally, only on one of the master nodes the following kubernetes bootstrapping happens:

* `bootkube.service` deploys the initial bootstrapping control-plane. It is started only after `kubelet.service` _is started_. It is a oneshot unit and cannot crash, and it runs only during bootstrap
* `tectonic.service` deploys tectonic control-plane. It is started only after `bootkube.service` _has completed_.  It is a oneshot unit and cannot crash, and it runs only during bootstrap

## Service ordering

Service ordering is enforced via systemd dependencies. This is the rationale for the settings, with relevant snippets:

### `kubelet.service`

```
Restart=always
WantedBy=multi-user.target
```

This service is enabled by default and can crash-loop until success.
It is started on every boot.

### `rm-assets.service`

This service waits for the bootkube and tectonic process to be completed.
It is a oneshot service, thus marked as started only once the script returns with success.
This is an optional service only present on platforms which pull assets from block storage.

## Diagram

This is a visual simplified representation of the overall bootstrapping flow.

```bob
Legend:
 * TF    -> terraform provisioner
 * IGN   -> ignition
 * k.s   -> kubelet.service
 * b.s   -> bootkube.service
 * t.s   -> tectonic.service
 * rm.s  -> rm-assets.service

.--------------------------------------------------------------------------------------------------------------------------------+
|                                                                                                                                |
|           Provision cloud/userdata           +----------+                                                                      |
|     ,---------------------------------------o|    TF    |                                                                      |
|     |                                        +----------+                                                                      |
|     |                                                                                                                          |
|     |                                                                                                                          |
|     |                                                                                                                          |
|     |                                                                                                                          |
|     V                                                                                                                          |
| +-------+                          Before   +------------+   Before                                                            |
| |  IGN  |                  .--------------->|    k.s     |o--------.                                                           |
| +-------+                  |                +------------+         |                                                           |
|     |                      |                   ^      |            |    +-----+      Before     +-------+   Before +-----+     |
|     '----------------------'                   |      v            '--->| b.s |o--------------->|  t.s  |--------> |rm.s |     |
|   Enable                                       '------'                 +-----+                 +-------+          +-----+     |
|                                                                                                                                |
|                                                                                                                                |
|                                       o                            o                                                           |
|                                       |                            |                                                           |
|                                       |        * Each boot         |         * First boot                                      |
|                                       |        * All nodes         |         * Bootkube master                                 |
|                                       |                            |                                                           |
'---------------------------------------o----------------------------o-----------------------------------------------------------+
```
