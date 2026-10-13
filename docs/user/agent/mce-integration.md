# MultiCluster Engine integration

This document describes how to install the multicluster engine for Kubernetes operator (mce) and how to deploy the cluster zero (hub) using the agent-based installer for Openshift.
The procedure is partially automated, and it will require some manual steps after the initial cluster deployment.

## Installing while disconnected from the greater internet

### Environment Pre-requisites

[Create a local image registry in your environment.](https://docs.openshift.com/container-platform/4.11/installing/disconnected_install/index.html) Make note of the DNS name and port of your local registry.

Download the [oc mirror](https://github.com/openshift/oc-mirror) tool.

#### Agent Installer pre-requisites

For the complete pre-requisites description, please check section 1.2.1 in [[1]](#1).

Create a folder (`<asset dir>`) and place a valid `install-config.yaml` and `agent-config.yaml`.

### Mirror the Openshift release and the MCE operator

Using the oc mirror tool mirror the desired Openshift release (4.11+). Here is an example ImageSetConfiguration for OCP 4.11:

__ocp-mce-imageset.yaml__
```
kind: ImageSetConfiguration
apiVersion: mirror.openshift.io/v1alpha2
archiveSize: 4
storageConfig:
  imageURL: <your-local-registry-dns-name>:<your-local-registry-port>/mirror/oc-mirror-metadata
  skipTLS: true
mirror:
  platform:
    architectures:
      - "amd64"
    channels:
      - name: stable-4.11
        type: ocp
  additionalImages:
    - name: registry.redhat.io/ubi8/ubi:latest
  operators:
    - catalog: registry.redhat.io/redhat/redhat-operator-index:v4.11
      packages:
        - name: multicluster-engine
        - name: local-storage-operator
```

With this file you will be able to use this command to mirror the OCP release and the MCE and LSO operators:

```oc mirror --dest-skip-tls --config ocp-mce-imageset.yaml docker://<your-local-registry-dns-name>:<your-local-registry-port>```

**Q: Why do I need the LSO operator?** A: The MCE operator is a large package and comes packaged with the infrastructure-operator as well. The infrastructure-operator needs a local storage volume to function.

### Update the mirror configurations in install-config.yaml

In your __install-config.yaml__ you will need to update the registry and the certificate. The two configurations to update are `imageContentSources` and `additionalTrustBundle`.
The following file snippet shows how to update the `imageContentSources` this for OCP and the MCE+LSO operators.

```
imageContentSources:
  - source: "quay.io/openshift-release-dev/ocp-release"
    mirrors:
      - "<your-local-registry-dns-name>:<your-local-registry-port>/openshift/release-images"
  - source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev"
    mirrors:
      - "<your-local-registry-dns-name>:<your-local-registry-port>/openshift/release"
  - source: "registry.redhat.io/ubi8"
    mirrors:
      - "<your-local-registry-dns-name>:<your-local-registry-port>/ubi8"
  - source: "registry.redhat.io/multicluster-engine"
    mirrors:
      - "<your-local-registry-dns-name>:<your-local-registry-port>/multicluster-engine"
  - source: "registry.redhat.io/rhel8"
    mirrors:
      - "<your-local-registry-dns-name>:<your-local-registry-port>/rhel8"
  - source: "registry.redhat.io/redhat"
    mirrors:
      - "<your-local-registry-dns-name>:<your-local-registry-port>/redhat"
```

**Note:** This procedure with `oc mirror` can be used to mirror any of the operators in the Red Hat Operator Indexes. After running `oc mirror` there will be a folder called `oc-mirror-workspace` with several outputs. There will be a file called __imageContentSourcePolicy.yaml__ in the most recent results folder that will identify all the mirrors you need for OCP and your selected operators.

You will also need to make sure your certificate is present in the `additionalTrustBundle` field.
This should look like something like this in your __install-config.yaml__:

```
additionalTrustBundle: |
  -----BEGIN CERTIFICATE-----
  .
  .
  .
  .
  .
  .
  -----END CERTIFICATE-----
```

Now you can generate the cluster manifests with the command:

`openshift-install agent create cluster-manifests`

This will update the cluster manifests to include a `mirror` folder with your mirror configuration.

**From here you should be able to follow the steps for installing while connected online.**

## Installing while connected online

### Pre-requisites

For the complete pre-requisites description, please check section 1.2.1 in [[1]](#1).

Create a folder (`<asset dir>`) and place a valid `install-config.yaml` and `agent-config.yaml`.

### Prepare the manifests for the installation

Create a subfolder named `openshift` within the `<asset dir>`. This subfolder will be used to store the extra manifests that will be applied during the installation to further customize the deployed cluster (note that the extra manifests will not be validated by the installer).

#### MCE operator manifests

Save the following manifests in the `<asset dir>/openshift` folder (use distinct files):

__99_01_mce_namespace.yaml__
```
apiVersion: v1
kind: Namespace
metadata:
  labels:
    openshift.io/cluster-monitoring: "true"
  name: multicluster-engine
```

__99_02_mce_operatorgroup.yaml__
```
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: multicluster-engine-operatorgroup
  namespace: multicluster-engine
spec:
  targetNamespaces:
  - multicluster-engine
```

__99_03_mce_subscription.yaml__
```
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: multicluster-engine
  namespace: multicluster-engine
spec:
  channel: "stable-2.1"
  name: multicluster-engine
  source: redhat-operators
  sourceNamespace: openshift-marketplace
```

#### Assisted Installer Service manifests

Since distributed units (DU) will be installed via Assisted Installer Service (AIS), it will be necessary to enable it in the hub cluster (for more details, see [[2]](#2)). The AIS requires at least a couple of persistent volumes (PVs), and they could be installed via the OpenShift Local Storage operator (LSO) (see [[3]](#3)).

Save the following manifests in the `<asset dir>/openshift` folder for the LSO setup (still using a separate file for each manifest):

__99_04_lso_namespace.yaml__
```
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    openshift.io/cluster-monitoring: "true"
  name: openshift-local-storage
```

__99_05_lso_operatorgroup.yaml__
```
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: local-operator-group
  namespace: openshift-local-storage
spec:
  targetNamespaces:
    - openshift-local-storage
```

__99_06_lso_subscription.yaml__
```
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: local-storage-operator
  namespace: openshift-local-storage
spec:
  installPlanApproval: Automatic
  name: local-storage-operator
  source: redhat-operators
  sourceNamespace: openshift-marketplace
```

### Create the agent ISO

At this point your filesystem layout should look like the following:

```
<asset dir>
    ├─ install-config.yaml
    ├─ agent-config.yaml
    └─ /openshift
        ├─ 99_01_mce_namespace.yaml
        ├─ 99_02_mce_operatorgroup.yaml
        ├─ 99_03_mce_subscription.yaml
        ├─ 99_04_lso_namespace.yaml
        ├─ 99_05_lso_operatorgroup.yaml
        └─ 99_06_lso_subscription.yaml
```

Create the agent ISO by running the command:

```
$ openshift-install agent create image --dir <asset dir>
```

### Cluster installation

Once ready, boot the target machine using the agent ISO and wait for the installation to successfully complete. To monitor the ongoing deployment you can use this command:

```
$ openshift-install agent wait-for install-complete --dir <asset dir>
```

### Hub setup

 As soon as the installation is completed it'd be possible then to finalize the setup to have a fully functioning hub cluster.
 The manifests shown in this section are meant to be applied manually. The order is relevant, and where needed the required
 waiting condition will also be illustrated.

#### Local volumes

Use the following manifest to create the required PVs that will be used by the AIS. Note that the `devicePaths` configuration
should match your target machines hardware setup:

__07_localvolumes.yaml__
 ```
 apiVersion: local.storage.openshift.io/v1
kind: LocalVolume
metadata:
  name: assisted-service
  namespace: openshift-local-storage
spec:
  logLevel: Normal
  managementState: Managed
  storageClassDevices:
    - devicePaths:
        - /dev/vda
        - /dev/vdb
      storageClassName: assisted-service
      volumeMode: Filesystem
```
```
$ oc apply -f 07_localvolumes.yaml
```

To wait for the availability of the PVs you can use this command:

```
$ oc wait localvolume -n openshift-local-storage assisted-service --for condition=Available --timeout 10m
```

#### MCE

Create a new multicluster engine instance by applying the following manifest:

__08_mce.yaml__
```
apiVersion: multicluster.openshift.io/v1
kind: MultiClusterEngine
metadata:
  name: multiclusterengine
spec: {}
```
```
$ oc apply -f 08_mce.yaml
```

#### Enable the Assisted Installer service

The AIS could be enabled through this manifest:

__09_agentserviceconfig.yaml__
```
apiVersion: agent-install.openshift.io/v1beta1
kind: AgentServiceConfig
metadata:
  name: agent
  namespace: assisted-installer
spec:
 databaseStorage:
  storageClassName: assisted-service
  accessModes:
  - ReadWriteOnce
  resources:
   requests:
    storage: 10Gi
 filesystemStorage:
  storageClassName: assisted-service
  accessModes:
  - ReadWriteOnce
  resources:
   requests:
    storage: 10Gi
```
```
$ oc apply -f 09_agentserviceconfig.yaml
```

At this stage you can also apply the following manifest, as it will be useful when deploying spoke clusters:

__10_clusterimageset.yaml__
```
apiVersion: hive.openshift.io/v1
kind: ClusterImageSet
metadata:
  name: "4.12"
spec:
  releaseImage: quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64
```
```
$ oc apply -f 10_clusterimageset.yaml
```

#### Auto-importing the hub cluster

As the last step of the finalization procedure, at this point you can import the current cluster (the one that it's hosting the mce operator and the assisted service) as the hub cluster, by applying the following manifest:

__11_autoimport.yaml__
```
apiVersion: cluster.open-cluster-management.io/v1
kind: ManagedCluster
metadata:
 labels:
   local-cluster: "true"
   cloud: auto-detect
   vendor: auto-detect
 name: local-cluster
spec:
 hubAcceptsClient: true
```
```
$ oc apply -f 11_autoimport.yaml
```

And wait for the managed cluster to be created:

```
$ oc wait -n multicluster-engine managedclusters local-cluster --for condition=ManagedClusterJoined=True --timeout 10m
```

If everything goes fine, you should be able to successfully observe the new managed cluster:

```
$ oc get managedcluster
NAME            HUB ACCEPTED   MANAGED CLUSTER URLS             JOINED   AVAILABLE  AGE
local-cluster   true           https://<your cluster url>:6443   True     True       77m
```

# References

* <a id="1">[1] https://access.redhat.com/documentation/en-us/red_hat_advanced_cluster_management_for_kubernetes/2.6/html/multicluster_engine/multicluster_engine_overview#installing-from-the-cli-mce
* <a id="2">[2] https://docs.openshift.com/container-platform/4.9/scalability_and_performance/ztp-deploying-disconnected.html
* <a id="3">[3] https://docs.openshift.com/container-platform/4.11/storage/persistent_storage/persistent-storage-local.html