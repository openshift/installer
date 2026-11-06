# Converting a vSphere cluster to component-scoped credentials

Clusters installed with `global` vCenter credentials use a single vCenter account for every OpenShift component: the installer stores that account in the `vsphere-creds` Secret in the `kube-system` namespace, and the Cloud Credential Operator (CCO) copies it into each component's credential Secret and keeps the copies synchronized. This procedure converts such a cluster to `component-scoped` credentials, where machine management, storage, the cloud controller manager, and the vSphere problem detector each use a dedicated, least-privilege vCenter account.

See [vCenter credential types](customization.md#vcenter-credential-types) for the credential model and the [candidate component privilege matrix](component-privileges.md) for the privileges each component account requires.

This is a day-2 operation performed entirely with `oc` on the running cluster. The installer is not involved, and the `VSpherePerComponentScopedCreds` feature gate is not relevant here: that gate is only read by the installer when it validates an install config. Component-scoped credentials are a Development Preview feature; the same support status applies to this procedure. Development Preview features are not supported in production clusters and are excluded from upgrade guarantees.

NOTE: Ordering matters. While CCO runs in its default mode on vSphere, which behaves as passthrough, it synchronizes each component Secret from the `kube-system/vsphere-creds` root Secret: a component Secret you replace is reverted to the root Secret's data on the next sync. Switch CCO to Manual mode and remove the root Secret *before* replacing the component Secrets, as described below.

## Prerequisites

* Cluster administrator access to the cluster with `oc`.
* A cluster that uses global credentials: the `vsphere-creds` Secret exists in `kube-system`, and CCO is in its default mode or explicitly in `Passthrough` mode. Clusters that already run `credentialsMode: Manual` with per-component Secrets do not need this procedure.
* The four component accounts created in vCenter for every configured vCenter server, with the privileges described in the [candidate component privilege matrix](component-privileges.md). A privilege held by one account does not satisfy a requirement of another.

## Procedure

### Discover the current configuration

Determine the cluster's credential mode:

```console
$ oc get cloudcredential cluster -o jsonpath='{.spec.credentialsMode}'
```

An empty result means the default mode, which behaves as passthrough on vSphere; `Passthrough` may also be reported. Either is a valid starting point for this procedure.

Determine the configured vCenter servers. The servers are the keys of the existing root Secret, of the form `<vcenter-server>.username` and `<vcenter-server>.password`:

```console
$ oc -n kube-system get secret vsphere-creds -o json | jq -r '.data | keys[]'
"vcenter.example.com.password"
"vcenter.example.com.username"
```

Every component Secret you create later must contain a key pair for every configured vCenter server, using the same server names.

### Switch the Cloud Credential Operator to Manual mode

```console
$ oc patch cloudcredential cluster --type merge -p '{"spec":{"credentialsMode":"Manual"}}'
```

Confirm the change:

```console
$ oc get cloudcredential cluster -o jsonpath='{.spec.credentialsMode}'
Manual
```

In Manual mode, CCO no longer reconciles credentials requests and stops synchronizing the component Secrets from the root Secret.

### Remove the root Secret

```console
$ oc -n kube-system delete secret vsphere-creds
```

This removes the shared global credential from the cluster and ensures no component can fall back to it. The global vCenter account itself is retired at the end of this procedure.

### Replace the component Secrets

Replace each component's Secret with data from its component account. The Secret names, namespaces, and account assignments are:

| Component | Secret | Namespace | Component account |
|---|---|---|---|
| Machine API | `vsphere-cloud-credentials` | `openshift-machine-api` | `machineManagement` |
| Machine Config Operator | `vsphere-cloud-credentials` | `openshift-machine-config-operator` | `machineManagement` |
| vSphere CSI driver | `vmware-vsphere-cloud-credentials` | `openshift-cluster-csi-drivers` | `storage` |
| Cloud Controller Manager | `vsphere-cloud-credentials` | `openshift-cloud-controller-manager` | `cloudControllerManager` |
| vSphere Problem Detector | `vsphere-cloud-credentials` | `openshift-cluster-storage-operator` | `vsphereProblemDetector` |

Each Secret holds one `<vcenter-server>.username` and `<vcenter-server>.password` key pair per configured vCenter, with the values of that component's account. The Machine API and Machine Config Operator Secrets both use the `machineManagement` account.

For example, to replace the Machine API Secret:

```console
$ oc -n openshift-machine-api create secret generic vsphere-cloud-credentials \
    --from-literal=vcenter.example.com.username='ocp-machine-api@vsphere.local' \
    --from-literal=vcenter.example.com.password='machine-management-password' \
    --dry-run=client -o yaml | oc apply -f -
```

Repeat for the remaining Secrets using the table above, substituting each component's account, Secret name, and namespace. If the cluster uses more than one vCenter server, add a `--from-literal` pair for each one; the set of keys must match those of the Secret being replaced.

### Verify the conversion

Confirm the affected operators remain healthy:

```console
$ oc get co cloud-credential machine-api machine-config storage cloud-controller-manager
```

All listed operators must report `Available=True` and `Degraded=False`. The vSphere Problem Detector has no ClusterOperator of its own; it reports the problems it detects as conditions on the `storage` ClusterOperator, which is included in the list above.

Secrets mounted into pods refresh automatically, but not immediately. If a component does not pick up its new credentials, restart the pods that consume the Secret, for example:

```console
$ oc -n openshift-machine-api rollout restart deployment/machine-api-controllers
```

The consuming resources are:

| Component | Resource | Namespace |
|---|---|---|
| Machine API | `deployment/machine-api-controllers` | `openshift-machine-api` |
| Machine Config Operator | `deployment/machine-config-operator` | `openshift-machine-config-operator` |
| vSphere CSI driver | `deployment/vmware-vsphere-controller-manager`, `daemonset/vmware-vsphere-node-driver` | `openshift-cluster-csi-drivers` |
| Cloud Controller Manager | `deployment/vsphere-cloud-controller-manager` | `openshift-cloud-controller-manager` |
| vSphere Problem Detector | `deployment/vsphere-problem-detector-operator` | `openshift-cluster-storage-operator` |

Confirm the resource names on the cluster with `oc -n <namespace> get deploy` before restarting.

Then exercise each component to confirm its account is sufficient:

* Machine API: if the cluster manages machines with machine sets, scale one up and back down.
* Storage: create and delete a `PersistentVolumeClaim` on a vSphere `StorageClass`.
* Cloud Controller Manager: confirm nodes report a `vsphere://` provider ID, for example with `oc get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.spec.providerID}{"\n"}{end}'`.
* Problem Detector: confirm its checks run and pass with the `vsphereProblemDetector` account.

### Retire the global account

Once every component is verified, remove or disable the old global vCenter account. Nothing in the cluster references it anymore.

## Upgrades

In Manual mode, CCO reports `Upgradeable=False` before minor version upgrades until an administrator confirms the component credentials satisfy the release being upgraded to. Before each upgrade, review the [candidate component privilege matrix](component-privileges.md) for privilege changes affecting the component accounts, then acknowledge the upgrade:

```console
$ oc patch cloudcredential cluster --type merge \
    -p '{"metadata":{"annotations":{"cloudcredential.openshift.io/upgradeable-to":"<minor-version>"}}}'
```

Set `<minor-version>` to the minor version being upgraded to, for example `4.20`.

## Rolling back

To return the cluster to global credentials, recreate the root Secret with the global account and return CCO to its default mode:

```console
$ oc -n kube-system create secret generic vsphere-creds \
    --from-literal=vcenter.example.com.username='global-account@vsphere.local' \
    --from-literal=vcenter.example.com.password='global-account-password'
$ oc patch cloudcredential cluster --type merge -p '{"spec":{"credentialsMode":""}}'
```

CCO resumes passthrough propagation and overwrites the five component Secrets with the root Secret's data. Rolling back requires the global account to still exist.
