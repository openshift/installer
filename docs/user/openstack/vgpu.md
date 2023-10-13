# Installing a cluster on OpenStack with vGPU support

If the underlying OpenStack deployment have proper GPU hardware installed and
configured there is a way to pass down vGPU to the pods by using gpu-operator.


## Pre-requisites

The following steps are required to be checked before starting the deployment of
OpenShift.

- Appropriate hardware is installed (like [NVIDIA Tesla
  V100](https://www.nvidia.com/en-gb/data-center/tesla-v100)) on
  the OpenStack compute node
- NVIDIA host drivers installed and nouveau driver removed
- Compute service installed on it and properly configured


## Driver installation

All of the examples assume RHEL8.4 and OSP 16.2 are used.

Given, there is NVIDIA vGPU capable card installed on the machine which intended
to have compute role, which may be confirmed by using a command which should
display similar output:

```console
$ lspci -nn | grep -i nvidia
3b:00.0 3D controller [0302]: NVIDIA Corporation GV100GL [Tesla V100 PCIe 16GB]
[10de:1db4] (rev a1)
```
make sure to remove `nouveau` driver from loading. It might be necessary to add
it to `/etc/modprobe.d/blacklist.conf` and/or change grub config:

```console
$ sudo sed -i 's/console=/rd.driver.blacklist=nouveau console=/' /etc/default/grub
$ sudo grub2-mkconfig -o /boot/grub2/grub.cfg
```

After that install host vGPU NVIDIA drivers (which are available to download for
license purchasers on
[NVIDIA application hub](https://nvid.nvidia.com/dashboard/)):

```console
$ sudo rpm -iv NVIDIA-vGPU-rhel-8.4-510.73.06.x86_64.rpm
```

Note, that drivers version may differ. Be careful to get right RHEL version and
architecture of the drivers to match installed RHEL.

Reboot the machine. After reboot, confirm there are correct drivers used:

```console
$ lsmod | grep nvidia
nvidia_vgpu_vfio       57344  0
nvidia              39055360  11
mdev                   20480  2 vfio_mdev,nvidia_vgpu_vfio
vfio                   36864  3 vfio_mdev,nvidia_vgpu_vfio,vfio_iommu_type1
drm                   569344  4 drm_kms_helper,nvidia,mgag200
```

You can also use `nvidia-smi` tool for displaying device state.


## OpenStack compute node

There should be mediated devices populated by the driver (bus address may vary):

```console
$ ls /sys/class/mdev_bus/0000\:3b\:00.0/mdev_supported_types/
nvidia-105  nvidia-106  nvidia-107  nvidia-108  nvidia-109  nvidia-110
nvidia-111  nvidia-112  nvidia-113  nvidia-114  nvidia-115  nvidia-163
nvidia-217  nvidia-247  nvidia-299  nvidia-300  nvidia-301
```

Depending of the type of workload and purchased license edition, appropriate
types needs to be configured in `nova.conf` for compute node, i.e.:

```ini
...
[devices]
enabled_vgpu_types: nvidia-105

...
```

After compute service restart, placement-api should report additional resources -
command `openstack resource provider list` and `openstack resource provider inventory list <id of the main provider>`
should display VGPU resource class available. For more information
[navigate to OpenStack Nova docs](https://docs.openstack.org/nova/train/admin/virtual-gpu.html).


## OpenStack vGPU flavor

Now, create a flavor, to be used to spin up new vGPU enabled nodes:

```console
$ openstack flavor create --disk 25 --ram 8192 --vcpus 4 \
    --property "resources:VGPU=1" --public <nova_gpu_flavor>
```


## Create vGPU enabled Worker Nodes

Worker nodes can be created by using machine API. To do that,
[create new machineSet in OpenShift](https://docs.openshift.com/container-platform/4.11/machine_management/creating_machinesets/creating-machineset-osp.html).

```console
$ oc get machineset -n openshift-machine-api <machineset_name> -o yaml > vgpu_machineset.yaml
```

Edit yaml file, be sure to have different name, have replicas set to the amount
of your cGPU capacity at maximum and set the right flavor, which would hint
OpenStack about right resources to include into virtual machine (Note, that this
is just an example, yours might be different):

```yaml
apiVersion: machine.openshift.io/v1beta1
kind: MachineSet
metadata:
  annotations:
    machine.openshift.io/memoryMb: "8192"
    machine.openshift.io/vCPU: "4"
  labels:
    machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
    machine.openshift.io/cluster-api-machine-role: <node_role>
    machine.openshift.io/cluster-api-machine-type: <node_role>
  name: <infrastructure_ID>-<node_role>-gpu-0
  namespace: openshift-machine-api
spec:
  replicas: <amount_of_nodes_with_gpu>
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
      machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>-gpu-0
  template:
    metadata:
      labels:
        machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
        machine.openshift.io/cluster-api-machine-role: <node_role>
        machine.openshift.io/cluster-api-machine-type: <node_role>
        machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>-gpu-0
    spec:
      lifecycleHooks: {}
      metadata: {}
      providerSpec:
        value:
          apiVersion: openstackproviderconfig.openshift.io/v1alpha1
          cloudName: openstack
          cloudsSecret:
            name: openstack-cloud-credentials
            namespace: openshift-machine-api
          flavor: <nova_gpu_flavor>
          image: <glance_image_name_or_location>
          kind: OpenstackProviderSpec
          metadata:
            creationTimestamp: null
          networks:
          - filter: {}
            subnets:
            - filter:
                name: <infrastructure_ID>-nodes
                tags: openshiftClusterID=<infrastructure_ID>
          securityGroups:
          - filter: {}
            name: <infrastructure_ID>-<node_role>
          serverGroupName: <infrastructure_ID>-<node_role>
          serverMetadata:
            Name: <infrastructure_ID>-<node_role>
            openshiftClusterID: <infrastructure_ID>
          tags:
          - openshiftClusterID=<infrastructure_ID>
          trunk: true
          userDataSecret:
            name: <node_role>-user-data
```

Save the file, and create machineset:

```console
$ oc create -f vgpu_machineset.yaml
```

And wait for new node to show up. You can examine its presence and state using
`openstack server list` and after VM is ready `oc get nodes`. New node should be
available with status "Ready".


## Discover features and enable GPU

Now it's time to install two operators:

- [Node Feature Discovery](https://docs.openshift.com/container-platform/4.11/hardware_enablement/psap-node-feature-discovery-operator.html)
- [Gpu Operator](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/openshift/contents.html)

### Node Feature Discovery Operator

This operator is needed for labeling nodes with detected hardware features. It
is required by the gpu operator. To install it, follow
[the documentation for nfd operator](https://docs.openshift.com/container-platform/4.11/hardware_enablement/psap-node-feature-discovery-operator.html)

To include NVIDIA card(s) in the NodeFeatureDiscovery instance, following
changes has been made:

```yaml
apiVersion: nfd.kubernetes.io/v1
kind: NodeFeatureDiscovery
metadata:
  name: nfd-instance
  namespace: node-feature-discovery-operator
spec:
  instance: ""
  topologyupdater: false
  operand:
    image: registry.redhat.io/openshift4/ose-node-feature-discovery:v<ocp_version>
    imagePullPolicy: Always
  workerConfig:
    configData: |
      sources:
        pci:
          deviceClassWhitelist:
            - "10de"
          deviceLabelFields:
            - vendor
```

Be sure to replace `<ocp_version>` with correct OCP version.


### GPU Operator

Follow documentation for it on [NVIDIA
site](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/openshift/install-gpu-ocp.html#installing-the-nvidia-gpu-operator-using-the-cli),
which basically take down to following steps:

1. Create namespace and group (save to file an do the `oc create -f filename`):
   ```yaml
   ---
   apiVersion: v1
   kind: Namespace
   metadata:
     name: nvidia-gpu-operator
   ---
   apiVersion: operators.coreos.com/v1
   kind: OperatorGroup
   metadata:
     name: nvidia-gpu-operator-group
     namespace: nvidia-gpu-operator
   spec:
     targetNamespaces:
     - nvidia-gpu-operator
   ```
1. Get the proper channel for gpu-operator:
   ```console
   $ CH=$(oc get packagemanifest gpu-operator-certified \
       -n openshift-marketplace -o jsonpath='{.status.defaultChannel}')
   $ echo $CH
   v22.9
   ```
1. Get right name for the gpu-operator:
   ```console
   $ GPU_OP_NAME=$(oc get packagemanifests/gpu-operator-certified \
       -n openshift-marketplace -o json | jq \
       -r '.status.channels[]|select(.name == "'${CH}'")|.currentCSV')
   $ echo $GPU_OP_NAME
   gpu-operator-certified.v22.9.0
   ```
1. Now, create nvidia-sub.yaml with subscription with the values, which was
   earlier fetched (save to file an do the `oc create -f filename`):
   ```yaml
   apiVersion: operators.coreos.com/v1alpha1
   kind: Subscription
   metadata:
     name: gpu-operator-certified
     namespace: nvidia-gpu-operator
   spec:
     channel: "<channel>"
     installPlanApproval: Manual
     name: gpu-operator-certified
     source: certified-operators
     sourceNamespace: openshift-marketplace
     startingCSV: "<gpu_operator_name>"
   ```
1. Verify if installplan has been created.
   ```console
   $ oc get installplan -n nvidia-gpu-operator
   ```
   In column APPROVED you will see `false`
1. Approve the plan:
   ```console
   $ oc patch installplan.operators.coreos.com/<install_plan_name> \
       -n nvidia-gpu-operator --type merge \
       --patch '{"spec":{"approved":true }}'
   ```

Now, it is needed to build an image which will be used by gpu-operator for
building drivers on the cluster.

Download needed drivers from the [NVIDIA application
hub](https://nvid.nvidia.com/dashboard/), along with vgpuDriverCatalog.yaml
file. The only files needed for vGPU are (at the time of
writing):

- NVIDIA-Linux-x86_64-510.85.02-grid.run
- vgpuDriverCatalog.yaml
- gridd.conf

Note, that drivers which should be used are the **guest** ones, not the host,
which was installed on the OpenStack compute node.

Clone the driver repository and copy all of needed drivers to the
driver/rehel8/drivers directory:

```console
$ git clone https://gitlab.com/nvidia/container-images/driver
$ cd driver rhel8
$ cp /path/to/obtained/drivers/* drivers/
```

Create gridd.conf file and copy it to `drivers` (installation of licensing
server is out of scope for this document):
```
# Description: Set License Server Address
# Data type: string
# Format:  "<address>"
ServerAddress=<licensing_server_address>
```

Go to the driver/rhel8/ path, and prepare image:
```console
$ export PRIVATE_REGISTRY=<registry_name/path>
$ export OS_TAG=<ocp_tag>
$ export VERSION=<version>
$ export VGPU_DRIVER_VERSION=<vgpu_version>
$ export CUDA_VERSION=<cuda_version>
$ export TARGETARCH=<architecture>
$ podman build \
    --build-arg CUDA_VERSION=${CUDA_VERSION} \
    --build-arg DRIVER_TYPE=vgpu \
    --build-arg TARGETARCH=$TARGETARCH \
    --build-arg DRIVER_VERSION=$VGPU_DRIVER_VERSION \
    -t ${PRIVATE_REGISTRY}/driver:${VERSION}-${OS_TAG} .
```

where:

- `PRIVATE_REGISTRY` is a name for private registry where image will be pushed
  to/pulled from, i.e. "quay.io/someuser"
- `OS_TAG` is a proper string matching RHCOS version used for cluster
  installation, i.e. "rhcos4.12"
- `VERSION` may be any string or number, i.e. "1.0.0"
- `VGPU_DRIVER_VERSION` is a substring from drivers. I.e. if there is file for
  building driver like "NVIDIA-Linux-x86_64-510.85.02-grid.run", then the
  version will be "510.85.02-grid".
- `CUDA_VERSION` is the latest supported version of CUDA supported on that
  particular GPU (or any other needed), i.e. "11.7.1".
- `TARGETARCH` is the target architecture which cluster runs on (usually
  "x86_64")


Push image to the registry:
```console
$ podman push ${PRIVATE_REGISTRY}/driver:${VERSION}-${OS_TAG}
```

Create license server configmap:
```console
$ oc create configmap licensing-config \
    -n nvidia-gpu-operator --from-file=drivers/gridd.conf
```

Create secret for connecting to the registry:
```console
$ oc -n nvidia-gpu-operator \
    create secret docker-registry my-registry \
    --docker-server=${PRIVATE_REGISTRY} \
    --docker-username=<username> \
    --docker-password=<pass> \
    --docker-email=<e-mail>
```

Substitute `<username>` `<pass>` and `<e-mail>` with real data. Here,
`my-registry` is used as the name of the secret and also could be changed (it
corresponds with `imagePullSectrets` array in `clusterpolicy` later on).

Get the clusterpolicy:
```console
$ oc get csv -n nvidia-gpu-operator $GPU_OP_NAME \
    -o jsonpath={.metadata.annotations.alm-examples} | \
    jq .[0] > clusterpolicy.json
```

Edit it and add marked in fields:

```json
{
  ...
  "spec": {
    ...
    "driver": {
       ...
     "repository": "<registry_name/path>",
     "image": "driver",
     "imagePullSecrets": ["my-registry"],
     "licensingConfig": {
        "configMapName": "licensing-config",
        "nlsEnabled": true
     },
     "version": "<version>",
     ...
    }
    ...
  }
}
```

Apply changes:
```console
$ oc apply -f clusterpolicy.json
```

Wait for drivers to be built. It may take a while. State of the pods should be
either running or completed.
```console
$ oc get pods -n nvidia-gpu-operator
```

## Run sample app

To verify installation, create simple app (app.yaml):
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cuda-vectoradd
spec:
 restartPolicy: OnFailure
 containers:
 - name: cuda-vectoradd
   image: "nvidia/samples:vectoradd-cuda11.2.1"
   resources:
     limits:
       nvidia.com/gpu: 1
```

Run it:
```console
$ oc apply -f app.yaml
```

Check the logs after pod finish its job:
```console
$ oc logs cuda-vectoradd
[Vector addition of 50000 elements]
Copy input data from the host memory to the CUDA device
CUDA kernel launch with 196 blocks of 256 threads
Copy output data from the CUDA device to the host memory
Test PASSED
Done
```
