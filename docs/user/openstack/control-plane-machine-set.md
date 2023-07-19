# OpenStack Control plane machine set

The [cluster-control-plane-machine-set-operator](https://github.com/openshift/cluster-control-plane-machine-set-operator) manages the Control plane machines through the `ControlPlaneMachineSet` resource.

The `ControlPlaneMachineSet` (CPMS) contains a template of the Control plane Machine provider Spec. The template is applied uniformly to all Control plane Machines, except for a few properties of the template can change between Machines: these variable properties are defined in the `FailureDomain` stanza of the CPMS spec.

---

## Example 1: availability zones set to the default

The "default" availability zone `""` means that Nova will schedule instances ignoring availability zones. This is the default setting in OpenShift.

This is the CPMS that the Installer generates for a cluster with name `ocp1-2g2xs`, for which no `zones` have been set in `install-config.yaml`:

```yaml
apiVersion: machine.openshift.io/v1
kind: ControlPlaneMachineSet
metadata:
  creationTimestamp: null
  labels:
    machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
  name: cluster
  namespace: openshift-machine-api
spec:
  replicas: 3
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
      machine.openshift.io/cluster-api-machine-role: master
      machine.openshift.io/cluster-api-machine-type: master
  state: Active
  strategy: {}
  template:
    machineType: machines_v1beta1_machine_openshift_io
    machines_v1beta1_machine_openshift_io:
      failureDomains: # <-- Empty property
        platform: ""
      metadata:
        labels:
          machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
          machine.openshift.io/cluster-api-machine-role: master
          machine.openshift.io/cluster-api-machine-type: master
      spec:
        lifecycleHooks: {}
        metadata: {}
        providerSpec:
          value: # <-- The OpenStack providerSpec
            apiVersion: machine.openshift.io/v1alpha1
            cloudName: openstack
            cloudsSecret:
              name: openstack-cloud-credentials
              namespace: openshift-machine-api
            flavor: m1.xlarge
            image: ocp1-2g2xs-rhcos
            kind: OpenstackProviderSpec
            metadata:
              creationTimestamp: null
            networks:
            - filter: {}
              subnets:
              - filter:
                  name: ocp1-2g2xs-nodes
                  tags: openshiftClusterID=ocp1-2g2xs
            securityGroups:
            - filter: {}
              name: ocp1-2g2xs-master
            serverGroupName: ocp1-2g2xs-master
            serverMetadata:
              Name: ocp1-2g2xs-master
              openshiftClusterID: ocp1-2g2xs
            tags:
            - openshiftClusterID=ocp1-2g2xs
            trunk: true
            userDataSecret:
              name: master-user-data
```

In this case, the `spec.template.machines_v1beta1_machine_openshift_io.failureDomains` stanza does not contain the `openstack` property.

Leaving `failureDomains` without a platform value means: do not substitute the values in the `providerSpec`. In this case, the `providerSpec` does not contain the `availabilityZone` property, because it's implicitly set to the empty string.

---

## Example 2: three Compute availability zones

```yaml
apiVersion: machine.openshift.io/v1
kind: ControlPlaneMachineSet
metadata:
  creationTimestamp: null
  labels:
    machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
  name: cluster
  namespace: openshift-machine-api
spec:
  replicas: 3
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
      machine.openshift.io/cluster-api-machine-role: master
      machine.openshift.io/cluster-api-machine-type: master
  state: Active
  strategy: {}
  template:
    machineType: machines_v1beta1_machine_openshift_io
    machines_v1beta1_machine_openshift_io:
      failureDomains:
        openstack:
        - availabilityZone: zone-one
        - availabilityZone: zone-two
        - availabilityZone: zone-three
        platform: OpenStack
      metadata:
        labels:
          machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
          machine.openshift.io/cluster-api-machine-role: master
          machine.openshift.io/cluster-api-machine-type: master
      spec:
        lifecycleHooks: {}
        metadata: {}
        providerSpec:
          value: # <-- The OpenStack providerSpec
            apiVersion: machine.openshift.io/v1alpha1
            cloudName: openstack
            cloudsSecret:
              name: openstack-cloud-credentials
              namespace: openshift-machine-api
            flavor: m1.xlarge
            image: ocp1-2g2xs-rhcos
            kind: OpenstackProviderSpec
            metadata:
              creationTimestamp: null
            networks:
            - filter: {}
              subnets:
              - filter:
                  name: ocp1-2g2xs-nodes
                  tags: openshiftClusterID=ocp1-2g2xs
            securityGroups:
            - filter: {}
              name: ocp1-2g2xs-master
            serverGroupName: ocp1-2g2xs-master
            serverMetadata:
              Name: ocp1-2g2xs-master
              openshiftClusterID: ocp1-2g2xs
            tags:
            - openshiftClusterID=ocp1-2g2xs
            trunk: true
            userDataSecret:
              name: master-user-data
```

When reconciling the Machines, cluster-control-plane-machine-set-operator will match their spec against the template, after substituting the `availabilityZone` properties for each of them.

The three Control plane Machines will each be provisioned on a different availability zone.

---

## Example 3: three Storage availability zones

The storage availability zones apply to the root volume. The `providerSpec` must contain a `rootVolume` property.

```yaml
apiVersion: machine.openshift.io/v1
kind: ControlPlaneMachineSet
metadata:
  creationTimestamp: null
  labels:
    machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
  name: cluster
  namespace: openshift-machine-api
spec:
  replicas: 3
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
      machine.openshift.io/cluster-api-machine-role: master
      machine.openshift.io/cluster-api-machine-type: master
  state: Active
  strategy: {}
  template:
    machineType: machines_v1beta1_machine_openshift_io
    machines_v1beta1_machine_openshift_io:
      failureDomains:
        openstack:
        - rootVolume:
            availabilityZone: cinder-one
        - rootVolume:
            availabilityZone: cinder-two
        - rootVolume:
            availabilityZone: cinder-three
        platform: OpenStack
      metadata:
        labels:
          machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
          machine.openshift.io/cluster-api-machine-role: master
          machine.openshift.io/cluster-api-machine-type: master
      spec:
        lifecycleHooks: {}
        metadata: {}
        providerSpec:
          value: # <-- The OpenStack providerSpec
            apiVersion: machine.openshift.io/v1alpha1
            cloudName: openstack
            cloudsSecret:
              name: openstack-cloud-credentials
              namespace: openshift-machine-api
            flavor: m1.xlarge
            image: ocp1-2g2xs-rhcos
            kind: OpenstackProviderSpec
            metadata:
              creationTimestamp: null
            networks:
            - filter: {}
              subnets:
              - filter:
                  name: ocp1-2g2xs-nodes
                  tags: openshiftClusterID=ocp1-2g2xs
            rootVolume:
              diskSize: 30
              volumeType: performance
            securityGroups:
            - filter: {}
              name: ocp1-2g2xs-master
            serverGroupName: ocp1-2g2xs-master
            serverMetadata:
              Name: ocp1-2g2xs-master
              openshiftClusterID: ocp1-2g2xs
            tags:
            - openshiftClusterID=ocp1-2g2xs
            trunk: true
            userDataSecret:
              name: master-user-data
```

When reconciling the Machines, cluster-control-plane-machine-set-operator will match their spec against the template, after substituting the `availabilityZone` property of the `rootVolume` for each of them.

The three Control plane Machines will each have their `rootVolume` provisioned on a different availability zone.

---

## Example 4: three Compute availability zones, three Storage availability zones

The storage availability zones apply to the root volume. The `providerSpec` must contain a `rootVolume` property.

```yaml
apiVersion: machine.openshift.io/v1
kind: ControlPlaneMachineSet
metadata:
  creationTimestamp: null
  labels:
    machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
  name: cluster
  namespace: openshift-machine-api
spec:
  replicas: 3
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
      machine.openshift.io/cluster-api-machine-role: master
      machine.openshift.io/cluster-api-machine-type: master
  state: Active
  strategy: {}
  template:
    machineType: machines_v1beta1_machine_openshift_io
    machines_v1beta1_machine_openshift_io:
      failureDomains:
        openstack:
        - availabilityZone: nova-one
          rootVolume:
            availabilityZone: cinder-one
        - availabilityZone: nova-two
          rootVolume:
            availabilityZone: cinder-two
        - availabilityZone: nova-three
          rootVolume:
            availabilityZone: cinder-three
        platform: OpenStack
      metadata:
        labels:
          machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
          machine.openshift.io/cluster-api-machine-role: master
          machine.openshift.io/cluster-api-machine-type: master
      spec:
        lifecycleHooks: {}
        metadata: {}
        providerSpec:
          value: # <-- The OpenStack providerSpec
            apiVersion: machine.openshift.io/v1alpha1
            cloudName: openstack
            cloudsSecret:
              name: openstack-cloud-credentials
              namespace: openshift-machine-api
            flavor: m1.xlarge
            image: ocp1-2g2xs-rhcos
            kind: OpenstackProviderSpec
            metadata:
              creationTimestamp: null
            networks:
            - filter: {}
              subnets:
              - filter:
                  name: ocp1-2g2xs-nodes
                  tags: openshiftClusterID=ocp1-2g2xs
            rootVolume:
              diskSize: 30
              volumeType: performance
            securityGroups:
            - filter: {}
              name: ocp1-2g2xs-master
            serverGroupName: ocp1-2g2xs-master
            serverMetadata:
              Name: ocp1-2g2xs-master
              openshiftClusterID: ocp1-2g2xs
            tags:
            - openshiftClusterID=ocp1-2g2xs
            trunk: true
            userDataSecret:
              name: master-user-data
```

When reconciling the Machines, cluster-control-plane-machine-set-operator will match their spec against the template, after substituting `availabilityZone` and `rootVolume.availabilityZone` for each of them.

The three Control plane Machines will all be provisioned on a different availability zone and have their `rootVolume` provisioned on a different availability zone.

---

## Example 5: three Compute availability zones, three Storage types

The storage types apply to the root volume. The `providerSpec` must contain a `rootVolume` property.

```yaml
apiVersion: machine.openshift.io/v1
kind: ControlPlaneMachineSet
metadata:
  creationTimestamp: null
  labels:
    machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
  name: cluster
  namespace: openshift-machine-api
spec:
  replicas: 3
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
      machine.openshift.io/cluster-api-machine-role: master
      machine.openshift.io/cluster-api-machine-type: master
  state: Active
  strategy: {}
  template:
    machineType: machines_v1beta1_machine_openshift_io
    machines_v1beta1_machine_openshift_io:
      failureDomains:
        openstack:
        - availabilityZone: nova-one
          rootVolume:
            volumeType: fastpool-1
        - availabilityZone: nova-two
          rootVolume:
            volumeType: fastpool-2
        - availabilityZone: nova-three
          rootVolume:
            volumeType: fastpool-3
        platform: OpenStack
      metadata:
        labels:
          machine.openshift.io/cluster-api-cluster: ocp1-2g2xs
          machine.openshift.io/cluster-api-machine-role: master
          machine.openshift.io/cluster-api-machine-type: master
      spec:
        lifecycleHooks: {}
        metadata: {}
        providerSpec:
          value: # <-- The OpenStack providerSpec
            apiVersion: machine.openshift.io/v1alpha1
            cloudName: openstack
            cloudsSecret:
              name: openstack-cloud-credentials
              namespace: openshift-machine-api
            flavor: m1.xlarge
            image: ocp1-2g2xs-rhcos
            kind: OpenstackProviderSpec
            metadata:
              creationTimestamp: null
            networks:
            - filter: {}
              subnets:
              - filter:
                  name: ocp1-2g2xs-nodes
                  tags: openshiftClusterID=ocp1-2g2xs
            rootVolume:
              diskSize: 30
            securityGroups:
            - filter: {}
              name: ocp1-2g2xs-master
            serverGroupName: ocp1-2g2xs-master
            serverMetadata:
              Name: ocp1-2g2xs-master
              openshiftClusterID: ocp1-2g2xs
            tags:
            - openshiftClusterID=ocp1-2g2xs
            trunk: true
            userDataSecret:
              name: master-user-data
```

When reconciling the Machines, cluster-control-plane-machine-set-operator will match their spec against the template, after substituting `availabilityZone` and `rootVolume.volumeType` for each of them.

The three Control plane Machines will all be provisioned on a different availability zone and have their `rootVolume` provisioned with a different volume type.