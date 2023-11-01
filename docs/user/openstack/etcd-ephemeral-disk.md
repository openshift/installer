# Moving etcd to an ephemeral local disk

You can move etcd from a root volume (Cinder) to a dedicated ephemeral local disk to prevent or resolve performance issues.

## Prerequisites

* This migration is currently tested and documented as a day 2 operation.
* An OpenStack cloud where Nova is configured to use local storage for ephemeral disks. The `libvirt.images_type` option in `nova.conf` must not be `rbd`.
* An OpenStack cloud with Cinder being functional and enough available storage to accommodate 3 Root Volumes for the OpenShift control plane.
* OpenShift will be deployed with IPI for now; UPI is not yet documented but technically possible.
* The control-plane machine’s auxiliary storage device, such as /dev/vdb, must match the vdb. Change this reference in all places in the file.

## Procedure

* Create a Nova flavor for the Control Plane which allows 10 GiB of Ephemeral Disk:

```bash
openstack flavor create --ephemeral 10 [...]
```

* We will deploy a cluster with Root Volumes for the Control Plane. Here is an example of `install-config.yaml`:

```yaml
[...]
controlPlane:
  name: master
  platform:
    openstack:
      type: ${CONTROL_PLANE_FLAVOR}
      rootVolume:
        size: 100
        types:
        - ${CINDER_TYPE}
  replicas: 3
[...]
```

* Run openshift-install with the following parameters to create the cluster:

```bash
openshift-install create cluster --dir=install_dir
```

* Once the cluster has been deployed and is healthy, edit the ControlPlaneMachineSet (CPMS) to add the additional block ephemeral device that will be used by etcd:

```bash
oc patch ControlPlaneMachineSet/cluster -n openshift-machine-api --type json -p '[{"op": "add", "path": "/spec/template/machines_v1beta1_machine_openshift_io/spec/providerSpec/value/additionalBlockDevices", "value": [{"name": "etcd", "sizeGiB": 10, "storage": {"type": "Local"}}]}]'
```

> [!NOTE]
> Putting etcd on a block device of type Volume is not supported for performance reasons simply because we don't test it.
> While it's functionally the same as using the root volume, we decided to support local devices only for now.

* Wait for the control-plane to roll out with new Machines. A few commands can be used to check that everything is healthy:

```bash
oc wait --timeout=90m --for=condition=Progressing=false controlplanemachineset.machine.openshift.io -n openshift-machine-api cluster
oc wait --timeout=90m --for=jsonpath='{.spec.replicas}'=3 controlplanemachineset.machine.openshift.io -n openshift-machine-api cluster
oc wait --timeout=90m --for=jsonpath='{.status.updatedReplicas}'=3 controlplanemachineset.machine.openshift.io -n openshift-machine-api cluster
oc wait --timeout=90m --for=jsonpath='{.status.replicas}'=3 controlplanemachineset.machine.openshift.io -n openshift-machine-api cluster
oc wait --timeout=90m --for=jsonpath='{.status.readyReplicas}'=3 controlplanemachineset.machine.openshift.io -n openshift-machine-api cluster
oc wait clusteroperators --timeout=30m --all --for=condition=Progressing=false
```

* Check that we have 3 control plane machines, and that each machine has the additional block device:

```bash
cp_machines=$(oc get machines -n openshift-machine-api --selector='machine.openshift.io/cluster-api-machine-role=master' --no-headers -o custom-columns=NAME:.metadata.name)
if [[ $(echo "${cp_machines}" | wc -l) -ne 3 ]]; then
  exit 1
fi
for machine in ${cp_machines}; do
  if ! oc get machine -n openshift-machine-api "${machine}" -o jsonpath='{.spec.providerSpec.value.additionalBlockDevices}' | grep -q 'etcd'; then
  exit 1
  fi
done
```

* We will use a MachineConfig to handle etcd on local disk. Create a file named `98-var-lib-etcd.yaml` with this content:

```yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: master
  name: 98-var-lib-etcd
spec:
  config:
    ignition:
      version: 3.2.0
    systemd:
      units:
      - contents: |
          [Unit]
          Description=Make File System on /dev/vdb
          DefaultDependencies=no
          BindsTo=dev-vdb.device
          After=dev-vdb.device var.mount
          Before=systemd-fsck@dev-vdb.service

          [Service]
          Type=oneshot
          RemainAfterExit=yes
          ExecStart=/usr/sbin/mkfs.xfs -f /dev/vdb
          TimeoutSec=0

          [Install]
          WantedBy=var-lib-containers.mount
        enabled: true
        name: systemd-mkfs@dev-vdb.service
      - contents: |
          [Unit]
          Description=Mount /dev/vdb to /var/lib/etcd
          Before=local-fs.target
          Requires=systemd-mkfs@dev-vdb.service
          After=systemd-mkfs@dev-vdb.service var.mount

          [Mount]
          What=/dev/vdb
          Where=/var/lib/etcd
          Type=xfs
          Options=defaults,prjquota

          [Install]
          WantedBy=local-fs.target
        enabled: true
        name: var-lib-etcd.mount
      - contents: |
          [Unit]
          Description=Sync etcd data if new mount is empty
          DefaultDependencies=no
          After=var-lib-etcd.mount var.mount
          Before=crio.service

          [Service]
          Type=oneshot
          RemainAfterExit=yes
          ExecCondition=/usr/bin/test ! -d /var/lib/etcd/member
          ExecStart=/usr/sbin/setenforce 0
          ExecStart=/bin/rsync -ar /sysroot/ostree/deploy/rhcos/var/lib/etcd/ /var/lib/etcd/
          ExecStart=/usr/sbin/setenforce 1
          TimeoutSec=0

          [Install]
          WantedBy=multi-user.target graphical.target
        enabled: true
        name: sync-var-lib-etcd-to-etcd.service
      - contents: |
          [Unit]
          Description=Restore recursive SELinux security contexts
          DefaultDependencies=no
          After=var-lib-etcd.mount
          Before=crio.service

          [Service]
          Type=oneshot
          RemainAfterExit=yes
          ExecStart=/sbin/restorecon -R /var/lib/etcd/
          TimeoutSec=0

          [Install]
          WantedBy=multi-user.target graphical.target
        enabled: true
        name: restorecon-var-lib-etcd.service
```

* Apply this file that will create the device and sync the data by entering the following command:

```bash
oc create -f 98-var-lib-etcd.yaml
```

* This will take some time to complete, as the etcd data will be synced from the root volume to the local disk on
the control-plane machines. Run these commands to check whether the cluster is healthy:

```bash
oc wait --timeout=45m --for=condition=Updating=false machineconfigpool/master
oc wait node --selector='node-role.kubernetes.io/master' --for condition=Ready --timeout=30s
oc wait clusteroperators --timeout=30m --all --for=condition=Progressing=false
```


* Once the cluster is healthy, create a file named `etcd-replace.yaml` with this content:

```yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: master
  name: 98-var-lib-etcd
spec:
  config:
    ignition:
      version: 3.2.0
    systemd:
      units:
      - contents: |
          [Unit]
          Description=Mount /dev/vdb to /var/lib/etcd
          Before=local-fs.target
          Requires=systemd-mkfs@dev-vdb.service
          After=systemd-mkfs@dev-vdb.service var.mount

          [Mount]
          What=/dev/vdb
          Where=/var/lib/etcd
          Type=xfs
          Options=defaults,prjquota

          [Install]
          WantedBy=local-fs.target
        enabled: true
        name: var-lib-etcd.mount
```

Apply this file that will remove the logic for creating and syncing the device by entering the following command:

```bash
oc replace -f etcd-replace.yaml
```

* Again we need to wait for the cluster to be healthy. The same commands as above can be used to check that everything is healthy.

* Now etcd is stored on ephemeral local disk. This can be verified by connected to a master nodes with `oc debug node/<master-node-name>` and running the following commands:

```bash
oc debug node/<master-node-name> -- df -T /host/var/lib/etcd
```
