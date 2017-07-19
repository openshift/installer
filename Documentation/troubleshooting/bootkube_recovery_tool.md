# Using utilities to recover Tectonic clusters

Tectonic provides tools to prevent and recover from the following disruptions:

* Loss of all API servers
* Loss of all schedulers
* Loss of all controller managers
* Loss of the majority of self-hosted etcd nodes

These control plane failures might occur due to internal or external interruptions, such as accidental deletion of master nodes. In this event, Tectonic employs pod `checkpoint` and `bootkube`  utilities for limited disaster avoidance and recovery support. `bootkube` is a program for launching self-hosted Kubernetes clusters. It provides a command, `bootkube recover`, to recover one or more components of a control plane. The pod `checkpoint` utility prevents some of these failures by watching over the API servers.

## Pod checkpoint

Tectonic Installer by default deploys pod `checkpoint` as a daemonset on all the master nodes. In the absence of an API server, the pod `checkpoint` utility recovers critical pods. The utility runs on each master node, periodically checking for locally running [parent pods][parent-pod] that have the `checkpointer.alpha.coreos.com/checkpoint=true` annotation, and writing [checkpoints][checkpoint] of their manifests to a local disk. If the parent pod is detected as down, the utility activates the checkpoint by running a static local copy of the pod. The checkpoint continues running until the parent pod is restarted on the local node, or an API server can be contacted to determine that the parent pod is no longer scheduled to this node.

Consider an example of an API server getting recovered from its checkpoint. When a node restarts, the kubelet must contact an API server to determine which pods to run on that node. However, if that node was running the only API server pod (non-HA scenario) in the Tectonic cluster, the kubelet will have no API server to contact to determine which pods to run, and therefore will not be able restart the API server.

In this case, the `checkpoint` utility automatically recovers the API server as follows:

1. Detects that the API server pod (parent pod) is not running.

2. Starts the checkpointed API server. The kubelet will now have an API server to contact to determine which pods to run. The kubelet will start the parent API server pod.

3. Detects that the parent pod is running again and stops the checkpointed API server.

The pod checkpointer is enabled by default and is highly recommended for all clusters to ensure node reboot resiliency. Disabling the utility may lead to cluster outages during node upgrades or reboots. For more information, see the [Pod Checkpointer][pod-checkpointer] documentation.

## Cleaning up Kubernetes resources

Before recovering a cluster, remove the docker state, etcd state, and active manifests from the failed machine. To do so, run the following on all the master nodes.

1. Back up `/var/log/containers` for later inspection:

   `tar -cvzf /var/log/containers containerlogs.tgz`

2.  Stop the kubelet service:

    `$ sudo systemctl stop kubelet.service`

2. Remove the docker state:

    ```
    $ docker kill $(docker ps -q)
    $ docker rm $(docker ps -a -q)
    ```

3. Remove the active Kubernetes active manifest:

    `$ sudo rm -rf /etc/kubernetes/manifests/*`

4. Remove any trace of pods (logs, state, configuration) that the kubelet was aware of:

   `$ sudo rm -rf /var/lib/kubelet/pods`

5. If self-hosted etcd is enabled:

    1. Back up the associated data.

       `$ sudo cp $(sudo find /var/etcd/ -iname db) /root/backup.db`

    2. Remove the data from the etcd directory:

       `$ sudo rm -rf /var/etcd/*`

6. Restart the kubelet service on one of the master nodes:

    `$ sudo systemctl restart kubelet.service`


## bootkube recover

Use the `bootkube recover` command to re-bootstrap the self-hosted control plane in the event of partial or total self-hosted control plane loss. `bootkube recover` does not recover a cluster directly. Instead, it extracts the control plane configuration from an available source and renders manifests in a format that the `bootkube start` command can use to reboot the cluster.

Use the latest version of `bootkube` when using the `bootkube recover` command, regardless of which version was used to create the cluster. To see available options, run:

`$ bootkube recover --help`

### Running bootkube

1. Download `bootkube` from either of the following:

   * [Container image][container-image]
   * [Release tarballs][release-tarball]

2. Extract the file to the current directory:

    `$ tar -xvzf  bootkube.tar.gz`

3. Change directory to navigate to the respective installer directory.

4. Provide necessary permissions to the group `core` for the etcd backup file,  `backup.db`.

    `$ sudo chown core:core /root/backup.db`

5. Run `bootkube recover` with a flag corresponding to the current state of the cluster:  

    `$ sudo ./bootkube recover --recovery-dir=/home/core/recovered --etcd-backup-file=/root/backup.db --kubeconfig=/etc/kubernetes/kubeconfig`

    Running the `bootkube recover` command creates a `recovery` directory in the home directory. The recovered data is stored in the `recovered` directory. If restoring the kubernetes data is successful, the following message is displayed:

    ```
      Attempting recovery using etcd backup file at "/root/backup.db"...
      Writing asset: /etc/kubernetes/manifests/recovery-etcd.yaml
      Waiting for etcd server to start...
      Writing asset: /home/core/recovered/bootstrap-manifests/bootstrap-kube-apiserver.yaml
      Writing asset: /home/core/recovered/bootstrap-manifests/bootstrap-kube-controller-manager.yaml
      Writing asset: /home/core/recovered/bootstrap-manifests/bootstrap-kube-scheduler.yaml
      Writing asset: /home/core/recovered/tls/etcd-client-ca.crt
      Writing asset: /home/core/recovered/tls/etcd-client.crt
      Writing asset: /home/core/recovered/tls/etcd-client.key
      Writing asset: /home/core/recovered/tls/etcd/peer-ca.crt
      Writing asset: /home/core/recovered/tls/etcd/peer.crt
      Writing asset: /home/core/recovered/tls/etcd/peer.key
      Writing asset: /home/core/recovered/tls/etcd/server-ca.crt
      Writing asset: /home/core/recovered/tls/etcd/server.crt
      Writing asset: /home/core/recovered/tls/etcd/server.key
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/etcd-client-ca.crt
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/etcd-client.crt
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/etcd-client.key
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/service-account.pub
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/apiserver.crt
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/apiserver.key
      Writing asset: /home/core/recovered/tls/secrets/kube-apiserver/ca.crt
      Writing asset: /home/core/recovered/tls/secrets/kube-controller-manager/ca.crt
      Writing asset: /home/core/recovered/tls/secrets/kube-controller-manager/service-account.key
      Writing asset: /home/core/recovered/bootstrap-manifests/bootstrap-etcd.yaml
      Writing asset: /home/core/recovered/etcd/bootstrap-etcd-service.json
      Writing asset: /home/core/recovered/etcd/migrate-etcd-cluster.json
      Writing asset: /home/core/recovered/auth/kubeconfig
      ```  

6. Run `bootkube start` to reboot the cluster.

    `$ sudo ./bootkube start --asset-dir=/home/core/recovered`

Supported cluster states and corresponding recovery methods are given below.

### Recovery with a running API server

If an API server is running but other control plane components are down, preventing cluster functionality, the control plane can be extracted directly from the API server:

```
$ bootkube recover --recovery-dir=recovered --kubeconfig=/etc/kubernetes/kubeconfig
```
Alternatively, use a rescue pod to recover the control plane as explained in [Disaster recovery of scheduler and controller manager pods][controller-recovery]. This method is manual, but does not require `bootkube`.

### Recovery with an external etcd cluster

If an [external etcd][external-etcd] cluster is running, the control plane can be extracted directly from etcd:

```
$ bootkube recover --recovery-dir=recovered --etcd-servers=<etcd-server-ip>:2379 --kubeconfig=/etc/kubernetes/kubeconfig
```

Replace `etcd-server-ip` with the IP address of your etcd cluster.

### Recovery with an external etcd backup

1. Recover the external etcd cluster from the backup.

   For more information, see [Disaster recovery][disaster-recovery-etcd].

2. Recover the control plane manifests:

    ```
    $ bootkube recover --recovery-dir=recovered --etcd-servers=<etcd-server-ip>:2379 --kubeconfig=/etc/kubernetes/kubeconfig
    ```

    Replace `etcd-server-ip` with the IP address of the recovered etcd cluster.

### Recovery with a provisioned etcd cluster

If a [provisioned etcd][provisioned-etcd] cluster is running, the control plane can be extracted directly from etcd:

```
$ bootkube recover --recovery-dir=recovered --etcd-servers=<etcd-server-ip>:2379 --kubeconfig=/etc/kubernetes/kubeconfig
```

Replace `etcd-server-ip` with the IP address of your etcd cluster.

### Recovery with a provisioned etcd backup

1. Recover the provisioned cluster from the backup.

   For more information, see [Disaster recovery][disaster-recovery-etcd].

2. Recover the control plane manifests:

    ```
    $ bootkube recover --recovery-dir=recovered --etcd-servers=<etcd-server-ip>:2379 --kubeconfig=/etc/kubernetes/kubeconfig
    ```

    Replace `etcd-server-ip` with the IP address of the recovered etcd cluster.


### Recovery with a self-hosted etcd backup

If using self-hosted etcd, recovery is supported via reading from an etcd backup file:

```
$ bootkube recover --recovery-dir=recovered --etcd-backup-file=backup --kubeconfig=/etc/kubernetes/kubeconfig
```

In addition to rebooting the control plane, this will also destroy and recreate the self-hosted etcd cluster by using the backup.

## Key concepts

**Self-hosted Kubernetes**: Self-hosted Kubernetes runs all required and optional components of a Kubernetes cluster on top of Kubernetes itself. A kubelet manages itself and all the Kubernetes components by using Kubernetes APIs. Tectonic clusters are self-hosted. For more information, see [self hosted Kubernetes][self-hosted-kubernetes].

**Self-hosted etcd**: A self-hosted etcd cluster runs in containers managed by Kubernetes. The term, self-hosted, itself implies that the cluster is hosted inside Kubernetes. A self-hosted etcd is deployed by using Tectonic Installer and is managed with Tectonic Console.

**Provisioned etcd**: Provisioned etcd clusters are deployed by the Tectonic Installer on a platform of your choice. These clusters are not controlled by Kubernetes, and therefore are not managed with Tectonic Console.

**External etcd**: An external etcd cluster is created and managed by a user outside of Tectonic Installer. Tectonic Installer assumes only network connectivity to the external etcd cluster at the given URL. The cluster cannot be managed with Tectonic Console.

**Parent pod**: The pod that is managed by the API server. The parent pod's manifest is backed up because it has the `checkpointer.alpha.coreos.com/checkpoint=true` metadata annotation.

**Checkpoint of a pod**: A checkpoint is a local on-disk copy of the manifest of a parent pod. A pod checkpoint ensures that existing local pod state can be recovered in the absence of an API server. If a parent pod stops running on the kubelet, though kubelet state indicates that it should still be running, the checkpointer will use the checkpoint manifest to run a temporary pod on the kubelet until the parent pod is up.


[bootkube-test-recovery]: https://github.com/kubernetes-incubator/bootkube/blob/master/hack/multi-node/bootkube-test-recovery
[bootkube-test-recovery-self-hosted-etcd]: https://github.com/kubernetes-incubator/bootkube/blob/master/hack/multi-node/bootkube-test-recovery-self-hosted-etcd
[checkpoint]: #key-concepts
[controller-recovery]: https://coreos.com/tectonic/docs/latest/troubleshooting/controller-recovery.html
[container-image]: https://quay.io/repository/coreos/bootkube
[disaster-recovery-etcd]: https://coreos.com/etcd/docs/latest/op-guide/recovery.html
[external-etcd]: #key-concepts
[pod-checkpointer]: https://github.com/kubernetes-incubator/bootkube/blob/master/cmd/checkpoint/README.md
[parent-pod]: #key-concepts
[provisioned-etcd]: #key-concepts
[release-tarball]: https://github.com/kubernetes-incubator/bootkube/releases
[self-hosted-kubernetes]: https://github.com/kubernetes/community/blob/master/contributors/design-proposals/self-hosted-kubernetes.md#what-is-self-hosted
[self-hosted]: #key-concepts
