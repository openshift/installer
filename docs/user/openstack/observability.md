# Observability of OpenShift on OpenStack

This document explains how it is possible to correlate OpenStack and OpenShift
metrics to have a better view of the stack and help troubleshoot issues
affecting your clusters.

## Make your OpenStack and OpenShift metrics available from the same metric store

As a first step, you should setup a metric store instance to have access to both
OpenStack and OpenShift metrics.

There are a number of ways to get your metrics from the same instance. We'll
document two of them: using Prometheus' remote-write to send the metrics to
another instance, and using the federation endpoint to query metrics from
a different instance.

### Using remote-write to send RHOSO and OCP metrics to an external storage

#### Set up the external storage

In this example, we'll use an external storage to store the metrics. We'll
assume that the external storage is another Prometheus instance.

We will set up remote-write from both OpenStack and OpenShift, authenticating them with mTLS
(mutual TLS). The target Prometheus needs to be configured to accept client TLS
certificates, and remote write.

// TODO: point at prometheus doc on how to accept TLS certificate and remote write

<!--
For testing purpose, we can do the following:

1. Provision a Fedora VM
1. Install `dnf install golang-github-prometheus caddy`
1. Configure prometheus to enable remote write. In `/etc/default/prometheus`, add the following line:
```
ARGS='--enable-feature=remote-write-receiver --storage.tsdb.retention.time=1d'
```
1. Start/enable the services
1. Open security group (443)
1. Setup caddy, in `/etc/caddy/Caddyfile`:

```Caddyfile
https://external-prometheus.example {
    # Just an example
    basicauth {
        # caddy hash-password
        user hashed-password
    }

    reverse_proxy http://localhost:9090
}
```

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

Prometheus UI should be available at https://external-prometheus.example
-->

#### Set up remote-write from RHOSO's telemetry-operator

Telemetry should already be enabled in the RHOSO environment. If this not yet
the case, refer to the
[documentation](https://docs.redhat.com/en/documentation/red_hat_openstack_services_on_openshift/18.0/html/customizing_the_red_hat_openstack_services_on_openshift_deployment/rhoso-observability_custom_dataplane#rhoso-observability_rhoso-observability).

> [!NOTE]
Make sure you have the Cluster Observability Operator installed in the OpenShift cluster running the OpenStack control plane, as this is a requirement for the OpenStack Telemetry Operator. Follow [these directions](https://github.com/openstack-k8s-operators/architecture/blob/main/examples/dt/uni01alpha/control-plane.md#cluster-observability-operator) to install it.

To confirm that the Telemetry Operator is correctly running, follow these steps

To check that the telemetry machinery is correctly installed, you can issue this command:

```bash
oc -n openstack get monitoringstacks metric-storage -o yaml
```

The `monitoringstacks` CRD being installed is a good indicator that telemetry is functional.


<!--
```bash
oc patch OpenStackControlPlane/controlplane --type merge -p '{"spec":{"telemetry":{"enabled": true, "template":{"ceilometer":{"enabled": true}}}}}'
```
-->

Before configuring remote-write, create the secret containing the HTTPS client certificates:

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

We will configure RHOSO's telemetry operator to write metrics to our external
Prometheus instance.

Look for the `metricStorage` stanza. It can be found at the
`.spec.telemetry.template.metricStorage` path. We will need to use
a `customMonitoringStack` structure that cannot coexist with the
`monitoringStack` one. Replace the `metricStorage` structure with this:

```yaml
      metricStorage:
        customMonitoringStack:
          alertmanagerConfig:
            disabled: false
          logLevel: info
          prometheusConfig:
            scrapeInterval: 30s
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

After saving the file and letting the change propagate, verify that you receive
metrics in the external Prometheus.

#### Set up remote-write from the OCP cluster-monitoring-operator

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
          targetLabel: <your_cluster_id>
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
create the secret containing the HTTPS client certificates, similar to what we
did for RHOSO.

```bash
oc --namespace openshift-monitoring \
    create secret generic mtls-bundle \
        --from-file=./ca.crt \
        --from-file=ocp-client.crt \
        --from-file=ocp-client.key
```

Once you have created the secret, it's time to apply the cluster-monitoring configuration:

```bash
oc apply -f cluster-monitoring-config.yaml
```

Let the change propagate and verify that you receive metrics in the external
Prometheus.


### Querying OCP metrics together with RHOSO's

As an alternative to write the different metrics to a single store, it's
possible to query remote metric stores instead.

OpenShift exposes a federation endpoint to enable queries originating from
outside the cluster to include OpenShift cluster monitoring data. You can
follow [these instructions][federation] to set it up.

[federation]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.16/html/monitoring/accessing-third-party-monitoring-apis#monitoring-querying-metrics-by-using-the-federation-endpoint-for-prometheus_accessing-monitoring-apis-by-using-the-cli "OpenShift documentation: Querying metrics by using the federation endpoint for Prometheus"


## Available mappings

OpenShift exposes metrics that establish a correlation between OpenStack
infrastructure resources and their representation in OpenShift.

To map Kubernetes nodes with OpenStack Nova instances:
* in the metric `kube_node_info`:
  * `node` is the Kubernetes node name
  * `provider_id` contains the identifier of the corresponding OpenStack Nova instance

To map Kubernetes persistent volumes with OpenStack Cinder volume or Manila share:
* in the metric `kube_persistentvolume_info`:
  * `persistentvolume` is the Kubernetes volume name
  * `csi_volume_handle` is the Cinder volume or Manila share identifier

When the metrics produced by the OpenShift cluster-monitoring-operator and the
metrics coming from OpenStack's Ceilometer can be queried together, for example
by concentrating them on a single Prometheus instance, those metrics can return
additional information on the health of the stack.

### Example queries

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


use metrics rather than openstack commands to gather state.
