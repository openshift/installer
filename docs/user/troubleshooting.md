# Installer Troubleshooting

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

If you have a Red Hat subscription for OpenShift, see [here][access-article] for support.

## Common Failures

### No Worker Nodes Created

The installer doesn't provision worker nodes directly, like it does with master nodes. Instead, the cluster relies on the Machine API Operator, which is able to scale up and down nodes on supported platforms. If more than fifteen to twenty minutes (depending on the speed of the cluster's Internet connection) have elapsed without any workers, the Machine API Operator needs to be investigated.

The status of the Machine API Operator can be checked by running the following command from the machine used to install the cluster:

```sh
oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig --namespace=openshift-machine-api get deployments
```

If the API is unavailable, that will need to be [investigated first](#kubernetes-api-is-unavailable).

The previous command should yield output similar to the following:

```
NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
cluster-autoscaler-operator   1/1     1            1           86m
machine-api-controllers       1/1     1            1           85m
machine-api-operator          1/1     1            1           86m
```

Check the machine controller logs with the following command.

```sh
oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig --namespace=openshift-machine-api logs deployments/machine-api-controllers --container=machine-controller
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

For added security, SSH isn't available from the Internet by default. There are several options for enabling this functionality:

- Create a bastion host that is accessible from the Internet and has access to the cluster. If the bootstrap machine hasn't been automatically destroyed yet, it can double as a temporary bastion host since it is given a public IP address.
- Configure network peering or a VPN to allow remote access to the private network.

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

The installer can also gather a log bundle from the bootstrap host using SSH as describe in [troubleshooting bootstrap](./troubleshootingbootstrap.md) document.

### etcd Is Not Running

During the bootstrap process, the Kubelet may emit errors like the following:

```
Error signing CSR provided in request from agent: error parsing profile: invalid organization
```

This is safe to ignore and merely indicates that the etcd bootstrapping is still in progress. etcd makes use of the CSR APIs provided by Kubernetes to issue and rotate its TLS assets, but these facilities aren't available before etcd has formed quorum. In order to break this dependency loop, a CSR service is run on the bootstrap node which only signs CSRs for etcd. When the Kubelet attempts to go through its TLS bootstrap, it is initially denied because the service it is communicating with only respects CSRs from etcd. After etcd starts and the control plane begins bootstrapping, an approver is scheduled and the Kubelet CSR requests will succeed.

### Installer Fails to Create Resources

The easiest way to get more debugging information from the installer is to check the log file (`.openshift_install.log`) in the install directory. Regardless of the logging level specified, the installer will write its logs in case they need to be inspected retroactively.

### Installer Fails to Initialize the Cluster

The installer uses the [cluster-version-operator] to create all the components of an OpenShift cluster. When the installer fails to initialize the cluster, the most important information can be fetched by looking at the [ClusterVersion][clusterversion] and [ClusterOperator][clusteroperator] objects:

1. Inspecting the `ClusterVersion` object.

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusterversion -oyaml
    apiVersion: config.openshift.io/v1
    kind: ClusterVersion
    metadata:
      creationTimestamp: 2019-02-27T22:24:21Z
      generation: 1
      name: version
      resourceVersion: "19927"
      selfLink: /apis/config.openshift.io/v1/clusterversions/version
      uid: 6e0f4cf8-3ade-11e9-9034-0a923b47ded4
    spec:
      channel: stable-4.1
      clusterID: 5ec312f9-f729-429d-a454-61d4906896ca
      upstream: https://api.openshift.com/api/upgrades_info/v1/graph
    status:
      availableUpdates: null
      conditions:
      - lastTransitionTime: 2019-02-27T22:50:30Z
        message: Done applying 4.1.1
        status: "True"
        type: Available
      - lastTransitionTime: 2019-02-27T22:50:30Z
        status: "False"
        type: Failing
      - lastTransitionTime: 2019-02-27T22:50:30Z
        message: Cluster version is 4.1.1
        status: "False"
        type: Progressing
      - lastTransitionTime: 2019-02-27T22:24:31Z
        message: 'Unable to retrieve available updates: unknown version 4.1.1
        reason: RemoteFailed
        status: "False"
        type: RetrievedUpdates
      desired:
        image: registry.svc.ci.openshift.org/openshift/origin-release@sha256:91e6f754975963e7db1a9958075eb609ad226968623939d262d1cf45e9dbc39a
        version: 4.1.1
      history:
      - completionTime: 2019-02-27T22:50:30Z
        image: registry.svc.ci.openshift.org/openshift/origin-release@sha256:91e6f754975963e7db1a9958075eb609ad226968623939d262d1cf45e9dbc39a
        startedTime: 2019-02-27T22:24:31Z
        state: Completed
        version: 4.1.1
      observedGeneration: 1
      versionHash: Wa7as_ik1qE=
    ```

    Some of most important [conditions][cluster-operator-conditions] to take note are `Failing`, `Available` and `Progressing`. You can look at the conditions using:

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusterversion version -o=jsonpath='{range .status.conditions[*]}{.type}{" "}{.status}{" "}{.message}{"\n"}{end}'
    Available True Done applying 4.1.1
    Failing False
    Progressing False Cluster version is 4.0.0-0.alpha-2019-02-26-194020
    RetrievedUpdates False Unable to retrieve available updates: unknown version 4.1.1
    ```

2. Inspecting the `ClusterOperator` object.

    You can get the status of all the cluster operators:

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusteroperator
    NAME                                  VERSION   AVAILABLE   PROGRESSING   FAILING   SINCE
    cluster-autoscaler                              True        False         False     17m
    cluster-storage-operator                        True        False         False     10m
    console                                         True        False         False     7m21s
    dns                                             True        False         False     31m
    image-registry                                  True        False         False     9m58s
    ingress                                         True        False         False     10m
    kube-apiserver                                  True        False         False     28m
    kube-controller-manager                         True        False         False     21m
    kube-scheduler                                  True        False         False     25m
    machine-api                                     True        False         False     17m
    machine-config                                  True        False         False     17m
    marketplace-operator                            True        False         False     10m
    monitoring                                      True        False         False     8m23s
    network                                         True        False         False     13m
    node-tuning                                     True        False         False     11m
    openshift-apiserver                             True        False         False     15m
    openshift-authentication                        True        False         False     20m
    openshift-cloud-credential-operator             True        False         False     18m
    openshift-controller-manager                    True        False         False     10m
    openshift-samples                               True        False         False     8m42s
    operator-lifecycle-manager                      True        False         False     17m
    service-ca                                      True        False         False     30m
    ```

    To get detailed information on why an individual cluster operator is `Failing` or not yet `Available`, you can check the status of that individual operator, for example `monitoring`:

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusteroperator monitoring -oyaml
    apiVersion: config.openshift.io/v1
    kind: ClusterOperator
    metadata:
      creationTimestamp: 2019-02-27T22:47:04Z
      generation: 1
      name: monitoring
      resourceVersion: "24677"
      selfLink: /apis/config.openshift.io/v1/clusteroperators/monitoring
      uid: 9a6a5ef9-3ae1-11e9-bad4-0a97b6ba9358
    spec: {}
    status:
      conditions:
      - lastTransitionTime: 2019-02-27T22:49:10Z
        message: Successfully rolled out the stack.
        status: "True"
        type: Available
      - lastTransitionTime: 2019-02-27T22:49:10Z
        status: "False"
        type: Progressing
      - lastTransitionTime: 2019-02-27T22:49:10Z
        status: "False"
        type: Failing
      extension: null
      relatedObjects: null
      version: ""
    ```

    Again, the cluster operators also publish [conditions][cluster-operator-conditions] like `Failing`, `Available` and `Progressing` that can help user provide information on the current state of the operator:

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusteroperator monitoring -o=jsonpath='{range .status.conditions[*]}{.type}{" "}{.status}{" "}{.message}{"\n"}{end}'
    Available True Successfully rolled out the stack
    Progressing False
    Failing False
    ```

    Each clusteroperator also publishes the list of objects owned by the cluster operator. To get that information:

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusteroperator kube-apiserver -o=jsonpath='{.status.relatedObjects}'
    [map[resource:kubeapiservers group:operator.openshift.io name:cluster] map[group: name:openshift-config resource:namespaces] map[group: name:openshift-config-managed resource:namespaces] map[group: name:openshift-kube-apiserver-operator resource:namespaces] map[group: name:openshift-kube-apiserver resource:namespaces]]
    ```

**NOTE:** Failing to initialize the cluster is usually not a fatal failure in terms of cluster creation as the user can look at the failures from `ClusterOperator` to debug failures for a cluster operator and take actions which can allow `cluster-version-operator` to make progress.

### Installer Fails to Fetch Console URL

The installer fetches the URL for OpenShift console using the [route][route-object] in `openshift-console` namespace. If the installer fails the fetch the URL for the console:

1. Check if the console router is `Available` or `Failing`

    ```console
    $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get clusteroperator console -oyaml
    apiVersion: config.openshift.io/v1
    kind: ClusterOperator
    metadata:
      creationTimestamp: 2019-02-27T22:46:57Z
      generation: 1
      name: console
      resourceVersion: "19682"
      selfLink: /apis/config.openshift.io/v1/clusteroperators/console
      uid: 960364aa-3ae1-11e9-bad4-0a97b6ba9358
    spec: {}
    status:
      conditions:
      - lastTransitionTime: 2019-02-27T22:46:58Z
        status: "False"
        type: Failing
      - lastTransitionTime: 2019-02-27T22:50:12Z
        status: "False"
        type: Progressing
      - lastTransitionTime: 2019-02-27T22:50:12Z
        status: "True"
        type: Available
      - lastTransitionTime: 2019-02-27T22:46:57Z
        status: "True"
        type: Upgradeable
      extension: null
      relatedObjects:
      - group: operator.openshift.io
        name: cluster
        resource: consoles
      - group: config.openshift.io
        name: cluster
        resource: consoles
      - group: oauth.openshift.io
        name: console
        resource: oauthclients
      - group: ""
        name: openshift-console-operator
        resource: namespaces
      - group: ""
        name: openshift-console
        resource: namespaces
      versions: null
    ```

2. Manually get the URL for `console`

  ```console
  $ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get route console -n openshift-console -o=jsonpath='{.spec.host}'
  console-openshift-console.apps.adahiya-1.devcluster.openshift.com
  ```

### Installer Fails to Add Default Ingress Certificate to Kubeconfig

The installer adds the default ingress certificate to the list of trusted client certificate authorities in `${INSTALL_DIR}/auth/kubeconfig`. If the installer fails to add the ingress certificate to `kubeconfig`, you can fetch the certificate from the cluster using the following command:

```console
$ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get configmaps default-ingress-cert -n openshift-config-managed -o=jsonpath='{.data.ca-bundle\.crt}'
-----BEGIN CERTIFICATE-----
MIIC/TCCAeWgAwIBAgIBATANBgkqhkiG9w0BAQsFADAuMSwwKgYDVQQDDCNjbHVz
dGVyLWluZ3Jlc3Mtb3BlcmF0b3JAMTU1MTMwNzU4OTAeFw0xOTAyMjcyMjQ2Mjha
Fw0yMTAyMjYyMjQ2MjlaMC4xLDAqBgNVBAMMI2NsdXN0ZXItaW5ncmVzcy1vcGVy
YXRvckAxNTUxMzA3NTg5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
uCA4fQ+2YXoXSUL4h/mcvJfrgpBfKBW5hfB8NcgXeCYiQPnCKblH1sEQnI3VC5Pk
2OfNCF3PUlfm4i8CHC95a7nCkRjmJNg1gVrWCvS/ohLgnO0BvszSiRLxIpuo3C4S
EVqqvxValHcbdAXWgZLQoYZXV7RMz8yZjl5CfhDaaItyBFj3GtIJkXgUwp/5sUfI
LDXW8MM6AXfuG+kweLdLCMm3g8WLLfLBLvVBKB+4IhIH7ll0buOz04RKhnYN+Ebw
tcvFi55vwuUCWMnGhWHGEQ8sWm/wLnNlOwsUz7S1/sW8nj87GFHzgkaVM9EOnoNI
gKhMBK9ItNzjrP6dgiKBCQIDAQABoyYwJDAOBgNVHQ8BAf8EBAMCAqQwEgYDVR0T
AQH/BAgwBgEB/wIBADANBgkqhkiG9w0BAQsFAAOCAQEAq+vi0sFKudaZ9aUQMMha
CeWx9CZvZBblnAWT/61UdpZKpFi4eJ2d33lGcfKwHOi2NP/iSKQBebfG0iNLVVPz
vwLbSG1i9R9GLdAbnHpPT9UG6fLaDIoKpnKiBfGENfxeiq5vTln2bAgivxrVlyiq
+MdDXFAWb6V4u2xh6RChI7akNsS3oU9PZ9YOs5e8vJp2YAEphht05X0swA+X8V8T
C278FFifpo0h3Q0Dbv8Rfn4UpBEtN4KkLeS+JeT+0o2XOsFZp7Uhr9yFIodRsnNo
H/Uwmab28ocNrGNiEVaVH6eTTQeeZuOdoQzUbClElpVmkrNGY0M42K0PvOQ/e7+y
AQ==
-----END CERTIFICATE-----
```

You can then **prepend** that certificate to `client-certificate-authority-data` field in your `${INSTALL_DIR}/auth/kubeconfig`.

## Generic Troubleshooting

Here are some ideas if none of the [common failures](#common-failures) match your symptoms.
For other generic troubleshooting, see [the Kubernetes documentation][kubernetes-debug].

### Check for Pending or Crashing Pods

This is the generic version of the [*No Worker Nodes Created*](#no-worker-nodes-created) troubleshooting procedure.

```console
$ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get pods --all-namespaces
NAMESPACE                              NAME                                                              READY     STATUS              RESTARTS   AGE
kube-system                            etcd-member-wking-master-0                                        1/1       Running             0          46s
openshift-machine-api                  machine-api-operator-586bd5b6b9-bxq9s                             0/1       Pending             0          1m
openshift-cluster-dns-operator         cluster-dns-operator-7f4f6866b9-kzth5                             0/1       Pending             0          2m
...
```

You can investigate any pods listed as `Pending` with:

```sh
oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig describe -n openshift-machine-api pod/machine-api-operator-586bd5b6b9-bxq9s
```

which may show events with warnings like:

```
Warning  FailedScheduling  1m (x10 over 1m)  default-scheduler  0/1 nodes are available: 1 node(s) had taints that the pod didn't tolerate.
```

You can get the image used for a crashing pod with:

```console
$ oc --kubeconfig=${INSTALL_DIR}/auth/kubeconfig get pod -o "jsonpath={range .status.containerStatuses[*]}{.name}{'\t'}{.state}{'\t'}{.image}{'\n'}{end}" -n openshift-machine-api machine-api-operator-586bd5b6b9-bxq9s
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
[cluster-version-operator]: https://github.com/openshift/cluster-version-operator/blob/master/README.md
[clusterversion]: https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusterversion.md
[clusteroperator]: https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusteroperator.md
[cluster-operator-conditions]: https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusteroperator.md#conditions
