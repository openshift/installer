# Node bootstrapping flow

This is a development document which describes the bootstrapping flow for ContainerLinux nodes provisioned by the tectonic-installer as part of a Tectonic cluster.

## Overview

When a cluster node is being bootstrapped from scratch, it goes through several phases in the following order:

1. first-boot OS configuration, via ignition (systemd units, node configuration, etc)
2. provisioning of additional assets (k8s manifests, TLS material), via either of:
   * pushing from terraform file/remote-exec (SSH)
   * pulling from private cloud stores (S3 buckets)
3. if needed, a node reboot is triggered to apply systemd-wide changes and to clean container runtime datadir
4. `kubelet.service` picks up the `kubelet.env` file and actually starts the kubelet as a rkt-fly service.

Additionally, only on one of the master nodes the following kubernetes bootstrapping happens:

1. `bootkube.service` is started after `kubelet.service` start
2. a static bootstrapping control-plane is deployed
3. a fully self-hosted control-plane starts and takes over the previous one
4. `bootkube.service` is completed with success
5. `tectonic.service` is started
6. a self-hosted tectonic control-plane is deployed
7. `tectonic.service` is completed with success

## Systemd units

The following systemd unit is deployed to a node by tectonic-installer and take part in the bootstrapping process:

* `kubelet.service` is the main kubelet daemon. It is automatically started on boot, it is crash-looping until `kubelet.env` is populated, and it runs on each boot

Additionally, only on one of the master nodes the following kubernetes bootstrapping happens:

* `bootkube.service` deploys the initial bootstrapping control-plane. It is started only after `kubelet.service` _is started_. It is a oneshot unit and cannot crash, and it runs only during bootstrap
* `bootkube.path` waits for bootkube assets/scripts to exist on disk and triggers `bootkube.service`
* `tectonic.service` deploys tectonic control-plane. It is started only after `bootkube.service` _has completed_.  It is a oneshot unit and cannot crash, and it runs only during bootstrap
* `bootkube.path` waits for tectonic assets/scripts to exist on disk and triggers `tectonic.service`

[tectonic-torcx]: https://github.com/coreos/tectonic-torcx

## Service ordering

Service ordering is enforced via systemd dependencies. This is the rationale for the settings, with relevant snippets:

### `kubelet.service`

```
EnvironmentFile=/etc/kubernetes/kubelet.env
ExecStart=/usr/lib/coreos/kubelet-wrapper [...]
Restart=always
WantedBy=multi-user.target
```

This service is enabled by default and can crash-loop until success.
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
ExecStart=/usr/bin/bash /opt/tectonic/tectonic-wrapper.sh
ExecStartPost=/bin/touch /opt/tectonic/init_tectonic.done
```

Tectonic service unit is not enabled by default. It is instead triggered by a path unit, which waits for assets written synchronously by terraform.

This service waits for bootkube process to be *completed* via systemd dependency.
It is a oneshot service, thus marked as started only once the script returns with success.
It is skipped on further boots, as the condition-path exists.

### `rm-assets.path` and `rm-assets.service`

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
 * b.p   -> bootkube.path
 * b.s   -> bootkube.service
 * t.p   -> tectonic.path
 * t.s   -> tectonic.service
 * rm.p  -> rm-assets.path
 * rm.s  -> rm-assets.service

.--------------------------------------------------------------------------------------------------------------------------------+
|                                                                                                                                |
|           Provision cloud/userdata           +----------+                Provision files                                       |
|     ,---------------------------------------o|    TF    |o-----------------.------------------------.-----------------+        |
|     |                                        +----------+                  |                        |                 |        |
|     |                                                                      v                        v                 v        |
|     |                                                                   +-----+                 +-------+         +------+     |
|     |                                                                   | b.p |                 |  t.p  |         | rm.p |     |
|     |                                                                   +-----+                 +-------+         +------+     |
|     V                                                                      o                        o                 o        |
| +-------+                          Before   +------------+   Before        | Trigger        Trigger |         Trigger |        |
| |  IGN  |                  .--------------->|    k.s     |o--------.       |                        |                 |        |
| +-------+                  |                +------------+         |       v                        v                 v        |
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
