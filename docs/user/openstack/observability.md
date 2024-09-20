# Observability of OpenShift on OpenStack

OpenShift exposes metrics that establish a correlation between OpenStack
infrastructure resources and their representation in OpenShift.

To map Kubernetes nodes with OpenStack Nova instances:
* in the metric `kube_node_info`:
  * `node` is the Kubernetes node name
  * `provider_id` contains the identifier of the corresponding OpenStack Nova instance

To map Kubernetes volumes with OpenStack Cinder or Manila volumes:
* in the metric `kube_persistentvolume_info`:
  * `persistentvolume` is the Kubernetes volume name
  * `csi_volume_handle` is the Cinder volume identifier

When the metrics produced by the OpenShift cluster-monitoring-operator and the metrics coming from OpenStack's Ceilometer can be queried together, for example by concentrating them on a single Prometheus instance, those metrics can return additional information on the health of the stack.

This query returns the number of OpenShift master nodes per OpenStack host:

```PromQL
sum by (vm_instance) (
  group by (vm_instance,resource) (ceilometer_cpu)
    / on (resource) group_right(vm_instance) (
      group by (node,resource) (
        label_replace(kube_node_info, "resource", "$1", "system_uuid", "(.+)")
      )
    / on (node) group_left group by (node) (
      cluster:master_nodes
    )
  )
)
```

//TODO: Add volume query

//TODO: Add query showing each pod resource consumption (CPU/RAM), with all pod data labeled with its Compute host

## Using remote-write to send RHOSO and OCP metrics to an external storage

<!--

To generate test certificates:

```bash
# Generate a CA if you don't have one already
openssl genrsa -out ca.key 4096
openssl req -batch -new -x509 -key ca.key -out ca.crt

# Generate the client certificates and sign them:
for target in server ocp-client osp-client; do
    openssl genrsa -out "${target}.key" 4096
    openssl req -batch -new -key "${target}.key" -out "${target}.csr"
    openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -in "${target}.csr" -out "${target}.crt"
done
```

-->

### Step 1: Set up the external storage

In this example, we'll assume that the external storage is another Prometheus instance.

We will set up remote-write from both instances, authenticating them with mTLS (mutual TLS). The target Prometheus needs to be configured to

//TODO: configure Prometheus to accept client TLS certificates

### Step 2: Set up remote-write from the OCP cluster-monitoring-operator

Documented [here](https://docs.openshift.com/container-platform/4.17/observability/monitoring/configuring-the-monitoring-stack.html#configuring_remote_write_storage_configuring-the-monitoring-stack).

You may want to label the cluster metrics with a cluster identifier.

Optionally, since metrics will be collected externally, you can set a reduced retention for local metrics.

The resulting `cluster-monitoring-config` ConfigMap should then resemble this:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
data:
  config.yaml: |
    prometheusK8s:
      retention: 1h
      remoteWrite:
      - url: "https://external-prometheus.example/api/v1/write"
        writeRelabelConfigs:
        - sourceLabels:
          - __tmp_openshift_cluster_id__
          targetLabel: cluster_id
          action: replace
        tlsConfig:
          ca:
            secret:
              name: mtls-bundle
              key: ca.crt
          cert:
            secret:
              name: mtls-bundle
              key: ocp-client.crt
          keySecret:
            name: mtls-bundle
            key: ocp-client.key
```

Save it to a file named `cluster-monitoring-config.yaml`. Before applying it,
create the secret containing the HTTPS material generated in the previous step.
In this example, it is called `mtls-bundle`:

```bash
oc --namespace openshift-monitoring \
    create secret generic mtls-bundle \
        --from-file=./ca.crt \
        --from-file=ocp-client.crt \
        --from-file=ocp-client.key
```

```bash
oc apply -f cluster-monitoring-config.yaml
```

### Step 3: Set up remote-write from RHOSO's telemetry-operator

As a prerequisite, Telemetry must be enabled in RHOSO. Refer to the OpenStack
[documentation](https://docs.redhat.com/en/documentation/red_hat_openstack_services_on_openshift/18.0/html/customizing_the_red_hat_openstack_services_on_openshift_deployment/rhoso-observability_custom_dataplane#rhoso-observability_rhoso-observability).

Before configuring remote-write, create the secret containing the HTTPS material generated in the previous step:

```bash
oc --namespace openstack \
    create secret generic mtls-bundle \
        --from-file=./ca.crt \
        --from-file=osp-client.crt \
        --from-file=osp-client.key
```

Then, configure remote-write from RHOSO's telemetry operator:

```bash
oc edit openstackcontrolplane/controlplane
```

Look for the `metricStorage` stanza. It can be found at this path: `.spec.telemetry.template.metricStorage`. Replace it with this:

```yaml
      metricStorage:
        customMonitoringStack:
          alertmanagerConfig:
            disabled: false
          logLevel: info
          prometheusConfig:
            remoteWrite:
            - url: https://external-prometheus.example/api/v1/write
              tlsConfig:
                ca:
                  secret:
                    name: mtls-bundle
                    key: ca.crt
                cert:
                  secret:
                    name: mtls-bundle
                    key: ocp-client.crt
                keySecret:
                  name: mtls-bundle
                  key: ocp-client.key
            replicas: 2
          resourceSelector:
            matchLabels:
              service: metricStorage
          resources:
            limits:
              cpu: 500m
              memory: 512Mi
            requests:
              cpu: 100m
              memory: 256Mi
          retention: 1h
        dashboardsEnabled: false
        dataplaneNetwork: ctlplane
        enabled: true
        prometheusTls: {}
```

## Querying OCP metrics together with RHOSO's

OpenShift exposes a federation endpoint to enable queries originating from
outside the cluster to include OpenShift cluster monitoring data. You can
follow [these instructions][federation] to set it up.

[federation]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.16/html/monitoring/accessing-third-party-monitoring-apis#monitoring-querying-metrics-by-using-the-federation-endpoint-for-prometheus_accessing-monitoring-apis-by-using-the-cli "OpenShift documentation: Querying metrics by using the federation endpoint for Prometheus"



use metrics rather than openstack commands to gather state.
