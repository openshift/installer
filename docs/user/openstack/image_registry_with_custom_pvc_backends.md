# Image Registry With A Custom PVC Backend

Sometimes users need to set the image registry backend volume from a specific availability zone, or with a size other than 100Gi. They can do this as a day2 operation by following next steps as a cluster admin.

1. Create a custom storage class with provided availability zone

```sh
$ export AZ_NAME=volume_availability_zone_name
$ cat <<\EOF | oc apply -f -
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: custom-csi-storageclass
provisioner: cinder.csi.openstack.org
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
parameters:
  availability: $AZ_NAME
EOF
storageclass.storage.k8s.io/custom-csi-storageclass created
```

**Note**: OpenShift doesn't check that the availability zone exists. So, users must verify that they have entered the correct value.

2. Create a pvc using the storageclass in `openshift-image-registry-namespace`, change the size of the volume if necessary.

```sh
$ cat <<EOF | oc apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: csi-pvc-imageregistry
  namespace: openshift-image-registry
  annotations:
    imageregistry.openshift.io: "true"
spec:
  accessModes:
  - ReadWriteOnce
  volumeMode: Block
  resources:
    requests:
      storage: 100Gi
  storageClassName: custom-csi-storageclass
EOF
persistentvolumeclaim/csi-pvc-imageregistry created
```

**Note**: Setting the `imageregistry.openshift.io` annotation is important, because otherwise Cluster Image Registry Operator won't be able to consume this PVC.

3. Replace the original volume claim in the image registry config.

```sh
$ oc patch configs.imageregistry.operator.openshift.io/cluster --type 'json' -p='[{"op": "replace", "path": "/spec/storage/pvc/claim", "value": "csi-pvc-imageregistry"}]'
config.imageregistry.operator.openshift.io/cluster patched
```

It may take several minutes for the operator to update the backend.

4. Verify that the new backend is in-use.

First, check the image registry config status. The claim name should be the same like in the `spec` section.

```sh
$ oc get configs.imageregistry.operator.openshift.io/cluster -oyaml
...
status:
    ...
    managementState: Managed
    pvc:
      claim: csi-pvc-imageregistry
...
```

Second, make sure that the PVC's status is Bound

```sh
$ oc get pvc -n openshift-image-registry csi-pvc-imageregistry
NAME                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS             AGE
csi-pvc-imageregistry  Bound    pvc-72a8f9c9-f462-11e8-b6b6-fa163e18b7b5   100Gi      RWO            custom-csi-storageclass  11m
```
