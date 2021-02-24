# How to set affinity rules for workers at install-time

The Installer automatically creates Control plane nodes with a
"soft-anti-affinity" policy. This means that each server will be created on a
separate host, if enough hosts are available.

The Installer does not set any affinity policy on Compute nodes; the choice is
left to the user.

We are going to create an OpenStack `server group` with an `anti-affinity`
policy. Then we are going to modify the `MachineSet` definitions during the
installation process, so that when `openshift-install create cluster` is run,
the Compute nodes will directly be created with the given policy.

### Step 1: create the Server group

In OpenStack Compute (Nova), scheduler policies are set using Server groups.

The first step is then to create the Server group.

Depending on the OpenStack configuration, there are a maximum of four possible policies:
* affinity
* soft-affinity
* anti-affinity
* soft-anti-affinity

These options are best described in the [OpenStack API reference][openstack-compute-api-docs].

Here we set `anti-affinity`:

```shell
openstack \
	--os-compute-api-version=2.15 \
	server group create \
	--policy anti-affinity \
	my-openshift-worker-group
```

Note that setting the Nova microversion ("os-compute-api-version") is only
mandatory when dealing with the "soft" policies, but does not hurt in any case.

**Take note of the resulting Server group id.**

### Step 2: generate the OpenShift manifests

```shell
openshift-install create manifests
```

This command will interactively prompt for the required information, if an
`install-config.yaml` file is not present in the current working directory.

### Step 3: Edit the MachineSet manifests

The "create manifests" command creates several files, among which the worker
MachineSet definition in
`openshift/99_openshift-cluster-api_worker-machineset-0.yaml`.

We must add a new property to the MachineSet: `serverGroupID`. That property lives under `spec.template.spec.providerSpec.value`.

**Use the Server group id from step 1.**

I have used the example UUID `aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee` in the template below.

```yaml
apiVersion: machine.openshift.io/v1beta1
kind: MachineSet
metadata:
  labels:
    machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
    machine.openshift.io/cluster-api-machine-role: <node_role>
    machine.openshift.io/cluster-api-machine-type: <node_role>
  name: <infrastructure_ID>-<node_role>
  namespace: openshift-machine-api
spec:
  replicas: <number_of_replicas>
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
      machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
  template:
    metadata:
      labels:
        machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
        machine.openshift.io/cluster-api-machine-role: <node_role>
        machine.openshift.io/cluster-api-machine-type: <node_role>
        machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
    spec:
      providerSpec:
        value:
          apiVersion: openstackproviderconfig.openshift.io/v1alpha1
          cloudName: openstack
          cloudsSecret:
            name: openstack-cloud-credentials
            namespace: openshift-machine-api
          flavor: <nova_flavor>
          image: <glance_image_name_or_location>
          serverGroupID: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
          kind: OpenstackProviderSpec
          networks:
          - filter: {}
            subnets:
            - filter:
                name: <subnet_name>
                tags: openshiftClusterID=<infrastructure_ID>
          securityGroups:
          - filter: {}
            name: <infrastructure_ID>-<node_role>
          serverMetadata:
            Name: <infrastructure_ID>-<node_role>
            openshiftClusterID: <infrastructure_ID>
          tags:
          - openshiftClusterID=<infrastructure_ID>
          trunk: true
          userDataSecret:
            name: <node_role>-user-data
          availabilityZone: <optional_openstack_availability_zone>
```

### Step 4: Complete the installation

After saving the MachineSet file in place, use `openshift-install` to complete
the installation.

Note: running the following command has the side effect of deleting both
`openshift` and `manifests` directories from the current working directory.

```shell
openshift-install create cluster
```

The modified `MachineSet` will be used by the Installer to create the Compute
nodes; the `serverGroupID` property make OpenShift create the Compute nodes
within that OpenStack Server Group.

[openstack-compute-api-docs]: https://docs.openstack.org/api-ref/compute/?expanded=create-server-group-detail#server-groups-os-server-groups
