# Node bootstrapping flow

This is a development document which describes the bootstrapping flow for ContainerLinux nodes provisioned by the tectonic-installer as part of a Tectonic cluster.

## Overview

When a cluster node is being bootstrapped from scratch, it goes through several phases in the following order:

1. first-boot OS configuration, via ignition (systemd units, node configuration, etc)
1. provisioning of additional assets (k8s manifests, TLS material), via either of:
   * pushing from terraform file/remote-exec (SSH)
   * pulling from private cloud stores (S3 buckets)
1. system-wide updates via `k8s-node-bootstrap.service`, which includes:
   * determining current kubernetes cluster version (when joining an existing cluster)
   * triggering a ContainerLinux update, via update-engine (optional)
   * downloading and deploying proper docker addon version, via tectonic-torcx
   * writing the `kubelet.env` file
1. if needed, a node reboot is triggered to apply systemd-wide changes and to clean container runtime datadir
1. `kubelet.service` picks up the `kubelet.env` file and actually starts the kubelet as a rkt-fly service.

Additionally, only on one of the master nodes the following kubernetes bootstrapping happens:

1. `bootkube.service` is started after `kubelet.service` start
1. a static bootstrapping control-plane is deployed
1. a fully self-hosted control-plane starts and takes over the previous one
1. `bootkube.service` is completed with success
1. `tectonic.service` is started
1. a self-hosted tectonic control-plane is deployed
1. `tectonic.service` is completed with success

## Systemd units

The following systemd units are deployed to a node by tectonic-installer and take part in the bootstrapping process:

* `k8s-node-bootstrap.service` ensures node and assets freshness. It is automatically started on boot, can crash-loop, and it runs only during bootstrap
* `kubelet.service` is the main kubelet daemon. It is automatically started on boot, it is crash-looping until `kubelet.env` is populated, and it runs on each boot

Additionally, only on one of the master nodes the following kubernetes bootstrapping happens:

* `bootkube.service` deploys the initial bootstrapping control-plane. It is started only after `kubelet.service` _is started_. It is a oneshot unit and cannot crash, and it runs only during bootstrap
* `bootkube.path` waits for bootkube assets/scripts to exist on disk and triggers `bootkube.service`
* `tectonic.service` deploys tectonic control-plane. It is started only after `bootkube.service` _has completed_.  It is a oneshot unit and cannot crash, and it runs only during bootstrap
* `bootkube.path` waits for tectonic assets/scripts to exist on disk and triggers `tectonic.service`

`k8s-node-bootstrap` runs [tectonic-torcx][tectonic-torcx] as a containerized service, thus relying on a container runtime being already on the node.
It currently assumes that Docker is available and working. In case of version changes, a cleanup of the Docker datadir `/var/lib/docker` is scheduled before rebooting.

[tectonic-torcx]: https://github.com/coreos/tectonic-torcx

## Service ordering

Service ordering is enforced via systemd dependencies. This is the rationale for the settings, with relevant snippets:

### `k8s-node-bootstrap.service`

```
ConditionPathExists=!/etc/kubernetes/kubelet.env
Before=kubelet.service
Restart=on-failure
ExecStartPre=[...]
ExecStart=/usr/bin/echo "node components bootstrapped"
WantedBy=multi-user.target kubelet.service
```

This service is enabled by default and can crash-loop until success.
Main logic happens in `Pre`, before the unit is marked as started, to block further services (a synchronous reboot can happen here).

In particular, this blocks kubelet from starting by:
 * a `WantedBy=` and `Before=`
 * writing the actual `kubelet.env` file on success.

It is skipped on further boots, as the condition-path exists.

### `kubelet.service`

```
EnvironmentFile=/etc/kubernetes/kubelet.env
ExecStart=/usr/lib/coreos/kubelet-wrapper [...]
Restart=always
WantedBy=multi-user.target
```

This service is enabled by default and can crash-loop until success.
On first boot, it is initially blocked by `k8s-node-bootstrap.service`.
It crash-loop until the `kubelet.env` file exists.
It is started on every boot.

### `bootkube.path` and `bootkube.service`

```
ConditionPathExists=!/opt/tectonic/init_bootkube.done
Wants=kubelet.service
After=kubelet.service
Type=oneshot
RemainAfterExit=true
ExecStart=/usr/bin/bash /opt/tectonic/bootkube.sh
ExecStartPost=/bin/touch /opt/tectonic/init_bootkube.done
```

Bootkube service unit is not enabled by default. It is instead triggered by a path unit, which waits for assets written synchronously by terraform.

This service waits for kubelet to be *started* via systemd dependency.
It is a oneshot service, thus marked as started only once the script return with success.
It is skipped on further boots, as the condition-path exists.

### `tectonic.path` and `tectonic.service`

```
ConditionPathExists=!/opt/tectonic/init_tectonic.done
Requires=bootkube.service
After=bootkube.service
Type=oneshot
RemainAfterExit=true
ExecStart=/usr/bin/bash /opt/tectonic/tectonic-rkt.sh
ExecStartPost=/bin/touch /opt/tectonic/init_tectonic.done
```

Tectonic service unit is not enabled by default. It is instead triggered by a path unit, which waits for assets written synchronously by terraform.

This service waits for bootkube process to be *completed* via systemd dependency.
It is a oneshot service, thus marked as started only once the script return with success.
It is skipped on further boots, as the condition-path exists.

## Diagram

This is a visual simplified representation of the overall bootstrapping flow.

```bob
Legend:
 * TF    -> terraform provisioner
 * IGN   -> ignition
 * knb.s -> k8s-node-bootstrap.service
 * k.s   -> kubelet.service
 * b.p   -> bootkube.path
 * b.s   -> bootkube.service
 * t.p   -> tectonic.path
 * t.s   -> tectonic.service

.------------------------------------------------------------------------------------------------------------------.
|                                                                                                                  |
|           Provision cloud/userdata                  +----------+                Provision files                  |
|     ,----------------------------------------------o|    TF    |o-----------------.------------------------.     |
|     |                                               +----------+                  |                        |     |
|     |                                                                             v                        v     |
|     |                  +----------+                                            +-----+                 +-------+ |
|     |             .--->| (reboot) |----.                                       | b.p |                 |  t.p  | |
|     |             |    +----------+    |                                       +-----+                 +-------+ |
|     V             |                    |                                          o                        o     |
| +-------+         |                    v  Before   +------------+   Before        | Trigger        Trigger |     |
| |  IGN  |         |                    *---------->|    k.s     |o--------.       |                        |     |
| +-------+         o                    ^           +------------+         |       v                        v     |
|     |       +----------+               |              ^      |            |    +-----+      Before     +-------+ |
|     '------>|   knb.s  |o--------------'              |      v            '--->| b.s |o--------------->|  t.s  | |
|   Enable    +----------+                              '------'                 +-----+                 +-------+ |
|                ^    |                                                                                            |
|                |    v                                                                                            |
|                '----'                        o                            o                                      |
|                                              |                            |                                      |
|               * First boot                   |        * Each boot         |         * First boot                 |
|               * All nodes                    |        * All nodes         |         * Bootkube master            |
|                                              |                            |                                      |
'----------------------------------------------o----------------------------o--------------------------------------'
```
