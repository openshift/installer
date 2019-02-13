# Installer Troubleshooting

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

If you have a Red Hat subscription for OpenShift, see [here][access-article] for support.

## Common Failures

### No Worker Nodes Created

The installer doesn't provision worker nodes directly, like it does with master nodes. Instead, the cluster relies on the Machine API Operator, which is able to scale up and down nodes on supported platforms. If more than fifteen to twenty minutes (depending on the speed of the cluster's Internet connection) have elapsed without any workers, the Machine API Operator needs to be investigated.

The status of the Machine API Operator can be checked by running the following command from the machine used to install the cluster:

```sh
oc --config=${INSTALL_DIR}/auth/kubeconfig --namespace=openshift-machine-api get deployments
```

If the API is unavailable, that will need to be [investigated first](#kubernetes-api-is-unavailable).

The previous command should yield output similar to the following:

```
NAME                             DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cluster-autoscaler-operator      1         1         1            1           1d
clusterapi-manager-controllers   1         1         1            1           1d
machine-api-operator             1         1         1            1           1d
```

Check the machine controller logs with the following command.

```sh
oc --config=${INSTALL_DIR}/auth/kubeconfig --namespace=openshift-machine-api logs deployments/clusterapi-manager-controllers --container=machine-controller
```

### Kubernetes API is Unavailable

When the Kubernetes API is unavailable, the master nodes will need to checked to ensure that they are running the correct components. This requires SSH access so it is necessary to include an administrator's SSH key during the installation.

If SSH access to the master nodes isn't available, that will need to be [investigated next](#unable-to-ssh-into-master-nodes).

The first thing to check is to make sure that etcd is running on each of the masters. The etcd logs can be viewed by running the following on each master node:

```sh
sudo crictl logs $(sudo crictl ps --pod=$(sudo crictl pods --name=etcd-member --quiet) --quiet)
```

If the previous command fails, ensure that the etcd pods have been created by the Kubelet:

```sh
sudo crictl pods --name=etcd-member
```

If no pods are shown, etcd will need to be [investigated](#etcd-is-not-running).

### Unable to SSH into Master Nodes

In order to SSH into the master nodes as user `core`, it is necessary to include an administrator's SSH key during the installation.
The selected key, if any, will be added to the `core` user's `~/.ssh/authorized_keys` via [Ignition](https://github.com/coreos/ignition) and is not configured via platform-specific approaches like [AWS key pairs][aws-key-pairs].
See [here][machine-config-daemon-ssh-keys] for information about managing SSH keys via the machine-config daemon.

If SSH isn't able to connect to the nodes, they may be waiting on the bootstrap node before they can boot. The initial set of master nodes fetch their boot configuration (the Ignition Config) from the bootstrap node and will not complete until they successfully do so. Check the console output of the nodes to determine if they have successfully booted or if they are waiting for Ignition to fetch the remote config.

Master nodes waiting for Ignition is indicative of problems on the bootstrap node. SSH into the bootstrap node to [investigate further](#troubleshooting-the-bootstrap-node).

### Troubleshooting the Bootstrap Node

If the bootstrap node isn't available, first double check that it hasn't been automatically removed by the installer. If it's not being created in the first place, the installer will need to be [troubleshot](#installer-fails-to-create-resources).

The most important thing to look at on the bootstrap node is `bootkube.service`. The logs can be viewed in two different ways:

1. If SSH is available, the following command can be run on the bootstrap node: `journalctl --unit=bootkube.service`
2. Regardless of whether or not SSH is available, the following command can be run: `curl --insecure --cert ${INSTALL_DIR}/tls/journal-gatewayd.crt --key ${INSTALL_DIR}/tls/journal-gatewayd.key 'https://${BOOTSTRAP_IP}:19531/entries?follow&_SYSTEMD_UNIT=bootkube.service'`

### etcd Is Not Running

etcd is started and managed by the Kubelet as a static pod. This requires a newer Kubelet which started shipping with version 47.29 of Red Hat CoreOS. The OS version can be checked using the following command:

```sh
grep OSTREE_VERSION /etc/os-release
```

If an older version of Red Hat CoreOS is in use, it will need to be updated. Try using the version suggested by the OpenShift Installer.

During the bootstrap process, the Kubelet may emit errors like the following:

```
Error signing CSR provided in request from agent: error parsing profile: invalid organization
```

This is safe to ignore and merely indicates that the etcd bootstrapping is still in progress. etcd makes use of the CSR APIs provided by Kubernetes to issue and rotate its TLS assets, but these facilities aren't available before etcd has formed quorum. In order to break this dependency loop, a CSR service is run on the bootstrap node which only signs CSRs for etcd. When the Kubelet attempts to go through its TLS bootstrap, it is initially denied because the service it is communicating with only respects CSRs from etcd. After etcd starts and the control plane begins bootstrapping, an approver is scheduled and the Kubelet CSR requests will succeed.

### Installer Fails to Create Resources

The easiest way to get more debugging information from the installer is to check the log file (`.openshift_install.log`) in the install directory. Regardless of the logging level specified, the installer will write its logs in case they need to be inspected retroactively.

## Generic Troubleshooting

Here are some ideas if none of the [common failures](#common-failures) match your symptoms.
For other generic troubleshooting, see [the Kubernetes documentation][kubernetes-debug].

### Check for Pending or Crashing Pods

This is the generic version of the [*No Worker Nodes Created*](#no-worker-nodes-created) troubleshooting procedure.

```console
$ oc --config=${INSTALL_DIR}/auth/kubeconfig get pods --all-namespaces
NAMESPACE                              NAME                                                              READY     STATUS              RESTARTS   AGE
kube-system                            etcd-member-wking-master-0                                        1/1       Running             0          46s
openshift-machine-api                  machine-api-operator-586bd5b6b9-bxq9s                             0/1       Pending             0          1m
openshift-cluster-dns-operator         cluster-dns-operator-7f4f6866b9-kzth5                             0/1       Pending             0          2m
...
```

You can investigate any pods listed as `Pending` with:

```sh
oc --config=${INSTALL_DIR}/auth/kubeconfig describe -n openshift-machine-api pod/machine-api-operator-586bd5b6b9-bxq9s
```

which may show events with warnings like:

```
Warning  FailedScheduling  1m (x10 over 1m)  default-scheduler  0/1 nodes are available: 1 node(s) had taints that the pod didn't tolerate.
```

You can get the image used for a crashing pod with:

```console
$ oc --config=${INSTALL_DIR}/auth/kubeconfig get pod -o "jsonpath={range .status.containerStatuses[*]}{.name}{'\t'}{.state}{'\t'}{.image}{'\n'}{end}" -n openshift-machine-api machine-api-operator-586bd5b6b9-bxq9s
machine-api-operator	map[running:map[startedAt:2018-11-13T19:04:50Z]]	registry.svc.ci.openshift.org/openshift/origin-v4.0-20181113175638@sha256:c97d0b53b98d07053090f3c9563cfd8277587ce94f8c2400b33e246aa08332c7
```

And you can see where that image comes from with:

```console
$ oc adm release info registry.svc.ci.openshift.org/openshift/origin-release:v4.0-20181113175638 --commits
Name:      v4.0-20181113175638
Digest:    sha256:58196e73cc7bbc16346483d824fb694bf1a73d517fe13f6b5e589a7e0e1ccb5b
Created:   2018-11-13 09:56:46 -0800 PST
OS/Arch:   linux/amd64
Manifests: 121

Images:
  NAME                  REPO                                               COMMIT
  ...
  machine-api-operator  https://github.com/openshift/machine-api-operator  e681e121e15d2243739ad68978113a07aa35c6ae
  ...
```

### One or more nodes are never Ready (Network / CNI issues)

You might see that one or more nodes are never ready, e.g

```console
$ kubectl get nodes
NAME                           STATUS     ROLES     AGE       VERSION
ip-10-0-27-9.ec2.internal      NotReady   master    29m       v1.11.0+d4cacc0
...
```

This usually means that, for whatever reason, networking is not available on the node. You can confirm this by looking at the detailed output of the node:

```console
$ kubectl describe node ip-10-0-27-9.ec2.internal
 ... (lots of output skipped)
'runtime network not ready: NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni config uninitialized'
```

The first thing to determine is the status of the SDN. The SDN deploys three daemonsets:
- *sdn-controller*, a control-plane component
- *sdn*, the node-level networking daemon
- *ovs*, the OpenVSwitch management daemon

All 3 must be healthy (though only a single `sdn-controller` needs to be running). `sdn` and `ovs` must be running on every node, and DESIRED should equal AVAILABLE. On a healthy 2-node cluster you would see:

```console
$ kubectl -n openshift-sdn get daemonsets
NAME             DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                     AGE
ovs              2         2         2       2            2           beta.kubernetes.io/os=linux       2h
sdn              2         2         2       2            2           beta.kubernetes.io/os=linux       2h
sdn-controller   1         1         1       1            1           node-role.kubernetes.io/master=   2h
```

If, instead, you get a different error message:

```console
$ kubectl -n openshift-sdn get daemonsets
No resources found.
```

This means the network-operator didn't run. Skip ahead [to that section](#debugging-the-cluster-network-operator). Otherwise, let's debug the SDN.

#### Debugging the openshift-sdn

On the NotReady node, you need to find out which pods, if any, are in a bad state. Be sure to substitute in the correct `spec.nodeName` (or just remove it).

```console
$ kubectl -n openshift-sdn get pod --field-selector "spec.nodeName=ip-10-0-27-9.ec2.internal"
NAME                   READY   STATUS             RESTARTS   AGE
ovs-dk8bh              1/1     Running            1          52m
sdn-8nl47              1/1     CrashLoopBackoff   3          52m
```

Then, retrieve the logs for the SDN (and the OVS pod, if that is failed):

```sh
kubectl -n openshift-sdn logs sdn-8nl47
```

Some common error messages:
- `Cannot fetch default cluster network`:  This means the `sdn-controller` has failed to run to completion. Retrieve its logs with `kubectl -n openshift-sdn logs -l app=sdn-controller`.
- `warning: Another process is currently listening on the CNI socket, waiting 15s`: Something has gone wrong, and multiple SDN processes are running. SSH to the node in question, capture the out of `ps -faux`. If you just need the cluster up, reboot the node.
- Error messages about ovs or OpenVSwitch: Check that the `ovs-*` pod on the same node is healthy. Retrieve its logs with `kubectl -n openshift-sdn logs ovs-<name>`. Rebooting the node should fix it.
- Any indication that the control plane is unavailable: Check to make sure the apiserver is reachable from the node. You may be able to find useful information via `journalctl -f -u kubelet`.

If you think it's a misconfiguration, file a [network operator](https://github.com/openshift/cluster-network-operator) issue. RH employees can also try #forum-sdn.

#### Debugging the cluster-network-operator
The cluster network operator is responsible for deploying the networking components. It does this in response to a special object created by the installer.

From a deployment perspective, the network operator is often the "canary in the coal mine." It runs very early in the installation process, after the master nodes have come up but before the bootstrap control plane has been torn down. It can be indicative of more subtle installer issues, such as long delays in bringing up master nodes or apiserver communication issues. Nevertheless, it can have other bugs.

First, determine that the network configuration exists:

```console
$ kubectl get network.config.openshift.io cluster -oyaml
apiVersion: config.openshift.io/v1
kind: Network
metadata:
  name: cluster
spec:
  serviceNetwork:
  - 172.30.0.0/16
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  networkType: OpenShiftSDN
```

If it doesn't exist, the installer didn't create it. You'll have to run `openshift-install create manifests` to determine why.

Next, check that the network-operator is running:

```sh
kubectl -n openshift-network-operator get pods
```

And retrieve the logs. Note that, on multi-master systems, the operator will perform leader election and all other operators will sleep:

```sh
kubectl -n openshift-network-operator logs -l "name=network-operator"
```

If appropriate, file a [network operator](https://github.com/openshift/cluster-network-operator) issue. RH employees can also try #forum-sdn.

[access-article]: https://access.redhat.com/articles/3780981#debugging-an-install-1
[aws-key-pairs]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
[kubernetes-debug]: https://kubernetes.io/docs/tasks/debug-application-cluster/
[machine-config-daemon-ssh-keys]: https://github.com/openshift/machine-config-operator/blob/master/docs/Update-SSHKeys.md
