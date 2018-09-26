# Cluster Bootstrapping Flow

This is a development document which describes the bootstrapping flow for an OpenShift cluster.

## Overview

All nodes in an OpenShift cluster are booted with a small Ignition config which references a larger dynamically-generated Ignition config, served within the cluster itself. This allows new nodes to always boot with the latest configuration. This does make bootstrap slightly more complex, however, because the master nodes require a cluster which doesn't yet exist in order to boot. This dependency loop is broken by a bootstrap node which temporarily hosts the control plane for the cluster.

On this bootstrap node, the following steps happen on boot:

1. `bootkube.service` is started after `kubelet.service` start
2. a static bootstrapping control-plane is deployed
3. a fully self-hosted control-plane starts (scheduled to the master nodes instead of the bootstrap node) and takes over the previous one
4. `bootkube.service` is completed with success
5. `tectonic.service` is started
6. a self-hosted tectonic control-plane is deployed
7. `tectonic.service` is completed with success

The result of this process is a fully running cluster. At this point, it is safe to remove the bootstrap node.
