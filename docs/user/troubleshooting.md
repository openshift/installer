# Installer Troubleshooting

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

## Common Failures

### No Worker Nodes Created

The installer doesn't provision worker nodes directly, like it does with master nodes. Instead, the cluster relies on the Machine API Operator, which is able to scale up and down nodes on supported platforms. If more than fifteen to twenty minutes (depending on the speed of the cluster's Internet connection) have elapsed without any workers, the Machine API Operator needs to be investigated.

The status of the Machine API Operator can be checked by running the following command from the machine used to install the cluster:

```sh
oc --config=${INSTALL_DIR}/auth/kubeconfig --namespace=openshift-cluster-api get pods
```

If the API is unavailable, that will need to be [investigated first](#kubernetes-api-is-unavailable).

The previous command should yield output similar to the following:

```
NAME                                             READY     STATUS    RESTARTS   AGE
clusterapi-manager-controllers-774dc4557-nx5xq   3/3       Running   0          4h
machine-api-operator-7894d8f85-lq2ts             1/1       Running   0          4h
```

The logs for the machine-controller container within the `clusterapi-manager-controllers` pod need to be checked to determine why the workers haven't been created. That can be done with the following (the exact name of the pod will need to be substituted):

```sh
oc --config=${INSTALL_DIR}/auth/kubeconfig --namespace=openshift-cluster-api logs clusterapi-manager-controllers-774dc4557-nx5xq --container=machine-controller
```

### Kubernetes API is Unavailable

When the Kubernetes API is unavailable, the master nodes will need to checked to ensure that they are running the correct components. This requires SSH access so it is necessary to include an administrator's SSH key during the installation.

If SSH access to the master nodes isn't available, that will need to be [investigated next](#unable-to-ssh-into-master-node).

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

In order to SSH into the master nodes, it is necessary to include an administrator's SSH key during the installation. If SSH authentication is failing, ensure that the proper SSH key is being used.

If SSH isn't able to connect to the nodes, they may be waiting on the bootstrap node before they can boot. The initial set of master nodes fetch their boot configuration (the Ignition Config) from the bootstrap node and will not complete until they successfully do so. Check the console output of the nodes to determine if they have successfully booted or if they are waiting for Ignition to fetch the remote config.

Master nodes waiting for Ignition is indicative of problems on the bootstrap node. SSH into the bootstrap node to [investigate further](#troubleshooting-the-bootstrap-node).

### Troubleshooting the Bootstrap Node

If the bootstrap node isn't available, first double check that it hasn't been automatically removed by the installer. If it's not being created in the first place, the installer will need to be [troubleshot](#installer-fails-to-create-resources).

After using SSH to access the bootstrap node, the most important thing to look at is `bootkube.service`. The logs can be viewed with the following command:

```sh
journalctl --unit=bootkube.service
```

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

The easiest way to get more debugging information from the installer is to increase the logging level. This can be done by adding `--log-level=debug` to the command line arguments. Of course, this cannot be retroactively applied, so it won't help to debug an installation that has already failed. The installation will have to be attempted again.

## Generic Troubleshooting

Here are some ideas if none of the [common failures](#common-failures) match your symptoms.
For other generic troubleshooting, see [the Kubernetes documentation][kubernetes-debug].

### Check for Pending or Crashing Pods

This is the generic version of the [*No Worker Nodes Created*](#no-worker-nodes-created) troubleshooting procedure.

```console
$ oc --config=${INSTALL_DIR}/auth/kubeconfig get pods --all-namespaces
NAMESPACE                              NAME                                                              READY     STATUS              RESTARTS   AGE
kube-system                            etcd-member-wking-master-0                                        1/1       Running             0          46s
openshift-cluster-api                  machine-api-operator-586bd5b6b9-bxq9s                             0/1       Pending             0          1m
openshift-cluster-dns-operator         cluster-dns-operator-7f4f6866b9-kzth5                             0/1       Pending             0          2m
...
```

You can investigate any pods listed as `Pending` with:

```sh
oc --config=${INSTALL_DIR}/auth/kubeconfig describe -n openshift-cluster-api pod/machine-api-operator-586bd5b6b9-bxq9s
```

which may show events with warnings like:

```
Warning  FailedScheduling  1m (x10 over 1m)  default-scheduler  0/1 nodes are available: 1 node(s) had taints that the pod didn't tolerate.
```

You can get the image used for a crashing pod with:

```console
$ oc --config=${INSTALL_DIR}/auth/kubeconfig get pod -o "jsonpath={range .status.containerStatuses[*]}{.name}{'\t'}{.state}{'\t'}{.image}{'\n'}{end}" -n openshift-cluster-api machine-api-operator-586bd5b6b9-bxq9s
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

[kubernetes-debug]: https://kubernetes.io/docs/tasks/debug-application-cluster/
