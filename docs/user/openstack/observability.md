# Observability of OpenShift on OpenStack

This document explains how it is possible to correlate OpenStack and OpenShift
metrics to have a better view of the stack and help troubleshoot issues
affecting your clusters.

This document focuses on Red Hat OpenStack Services on OpenShift (hereinafter
RHOSO, which corresponds to version 18 of the Red Hat OpenStack Platform).

## Make your OpenStack and OpenShift metrics available in the same metric store

The strategy we will be outlining in this document is to make both OpenStack
and OpenShift metrics available in a single Prometheus instance.

There are a number of ways to achieve this goal. Here we document two methods:

* Method A: use the Prometheus feature
  [Remote-Write][prometheus-docs-remote-write] to send both OpenStack and
  OpenShift metrics to an external instance
* Method B: configure the OpenStack prometheus instance to pull certain data
  from the OpenShift federation endpoint allowing data to be combined in the
  single OpenStack prometheus.

[prometheus-docs-remote-write]: https://prometheus.io/docs/specs/remote_write_spec/ "Prometheus Remote-Write Specification"

### Method A: Use Remote-Write to send RHOSO and OCP metrics to an external instance

#### Set up the external storage

In this example, we are using an external Prometheus instance to store the
metrics.

We will set up remote-write from both OpenStack and OpenShift, authenticating
them with mTLS (mutual TLS). The target Prometheus needs to be configured to
[accept client TLS certificates][prometheus-mtls], and
[Remote-Write][prometheus-remote-write-receiver-flag].

[prometheus-mtls]: https://prometheus.io/docs/prometheus/latest/configuration/https/
[prometheus-remote-write-receiver-flag]: https://prometheus.io/docs/prometheus/latest/feature_flags/#remote-write-receiver "Prometheus feature flags: Remote-Write receiver"


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

<!--
For testing purpose, we can do the following to set up basic auth in addition
or in stead of mTLS:

1. Provision a Fedora VM
2. Install `dnf install golang-github-prometheus caddy`
3. Configure prometheus to enable remote write (and limit retention to avoid
   filling up disk space). In `/etc/default/prometheus`, add the following
   line:

```
ARGS='--enable-feature=remote-write-receiver --storage.tsdb.retention.time=1d'
```

1. Enable and restart the Prometheus systemd unit
2. Add a security group rule to allow HTTPS (port 443)
3. Setup Caddy with (`/etc/caddy/Caddyfile`):

```Caddyfile
https://external-prometheus.example {

    basicauth {
        # caddy hash-password
        user hashed-password
    }

    reverse_proxy http://localhost:9090
}
```
-->



We will assume that the external Prometheus is reachable at the URL
`https://external-prometheus.example`.

#### Set up remote-write from RHOSO's telemetry-operator

Telemetry should be enabled in the RHOSO environment. If it is not the case,
refer to the
[documentation](https://docs.redhat.com/en/documentation/red_hat_openstack_services_on_openshift/18.0/html/customizing_the_red_hat_openstack_services_on_openshift_deployment/rhoso-observability_custom_dataplane#rhoso-observability_rhoso-observability).

<!--
Essentially, enabling telemetry boils down to flipping a property of the
openstackcontrolplane object:

```bash
oc -n openstack patch OpenStackControlPlane/controlplane --type merge -p '{"spec":{"telemetry":{"enabled": true, "template":{"ceilometer":{"enabled": true}, "metricStorage":{"enabled": true}}}}}'
```
-->

> [!NOTE]
Make sure you have the Cluster Observability Operator installed in the
OpenShift cluster running the OpenStack control plane, as this is a requirement
for the OpenStack Telemetry Operator. Follow [these
directions](https://github.com/openstack-k8s-operators/architecture/blob/main/examples/dt/uni01alpha/control-plane.md#cluster-observability-operator)
to install it.

To check that the telemetry machinery is correctly installed, issue this
command:

```bash
oc -n openstack get monitoringstacks metric-storage -o yaml
```

The `monitoringstacks` CRD being installed is a good indicator that telemetry
is functional.

Before configuring remote-write in RHOSO's telemetry operator, create a secret
in the `openstack namespace` containing the HTTPS client certificates for
authenticating to Prometheus. We'll call it `mtls-bundle`:

```bash
oc --namespace openstack \
    create secret generic mtls-bundle \
        --from-file=./ca.crt \
        --from-file=osp-client.crt \
        --from-file=osp-client.key
```

Then, edit the controlplane configuration to setup the metric storage:

```bash
oc -n openstack edit openstackcontrolplane/controlplane
```

We will configure RHOSO's telemetry operator to write metrics to our external
Prometheus instance.

Look for the `metricStorage` stanza. It can be found at the
`.spec.telemetry.template.metricStorage` path. We will need to use a
`customMonitoringStack` structure that cannot coexist with the
`monitoringStack` one. Replace the `metricStorage` structure with one that
looks like this:

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
          retention: 1d # Set the desired retention interval
        dashboardsEnabled: false
        dataplaneNetwork: ctlplane
        enabled: true
        prometheusTls: {}
```

After saving the file and letting the change propagate, verify that you receive
OpenStack metrics in the external Prometheus.

#### Set up remote-write from the OCP cluster-monitoring-operator

Refer to the [OpenShift documentation][ocp_docs] for configuring its monitoring stack.

In this example we will [create a cluster monitoring
configuration][create_cluster_monitoring_config], [setup
remote-write][setup_remote_write], and [label the cluster metrics with
a cluster identifier][add_labels].

Optionally, since metrics will be collected externally, you can set a reduced retention for local metrics.

The resulting `cluster-monitoring-config` ConfigMap could then resemble this:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
data:
  config.yaml: |
    prometheusK8s:
      retention: 1d # Set the desired retention interval
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
create the secret containing the HTTPS client certificates, similar to what we
did for RHOSO. We're still calling the secret `mtls-bundle`, but this time in
the `openshift-monitoring` namespace:

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

Let the change propagate and verify that you receive OpenShift metrics in the
external Prometheus.

[ocp_docs]: https://docs.openshift.com/container-platform/4.17/observability/monitoring/configuring-the-monitoring-stack.html#configuring_remote_write_storage_configuring-the-monitoring-stack "Configuring the monitoring stack"
[create_cluster_monitoring_config]: https://docs.openshift.com/container-platform/4.17/observability/monitoring/configuring-the-monitoring-stack.html#creating-cluster-monitoring-configmap_configuring-the-monitoring-stack "Creating a cluster monitoring config map"
[setup_remote_write]: https://docs.openshift.com/container-platform/4.17/observability/monitoring/configuring-the-monitoring-stack.html#configuring-remote-write-storage_configuring-the-monitoring-stack "Configuring remote write storage"
[add_labels]: https://docs.openshift.com/container-platform/4.17/observability/monitoring/configuring-the-monitoring-stack.html#adding-cluster-id-labels-to-metrics_configuring-the-monitoring-stack "Adding cluster ID labels to metrics"

### Method B: Scrape OCP metrics from RHOSO

As opposed to Remote-Write, this solution maintains the traditional direction
of the HTTP calls from the observer to the observed object. In other words, it
complies with the Prometheus "pull" flow.

In the following instructions, instead of using an external arbitrary
Prometheus instance, we will be using RHOSO's Prometheus as the collector of
both OpenShift and OpenStack metrics.

OpenShift exposes a federation endpoint to expose a subset of metrics to an
external scraper. You can follow [these instructions][federation] to get
acquainted to the endpoint.

[federation]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.17/html/monitoring/accessing-third-party-monitoring-apis#monitoring-querying-metrics-by-using-the-federation-endpoint-for-prometheus_accessing-monitoring-apis-by-using-the-cli "OpenShift documentation: Querying metrics by using the federation endpoint for Prometheus"

#### Step 1: Gather credentials and coordinates

While connected to the OpenShift cluster through a username identified by password (as opposed to logging in using the `kubeconfig` file generated by the installer), fetch a token:

```bash
oc whoami -t
```

Then get the Prometheus federation route URL:

```bash
oc -n openshift-monitoring get route prometheus-k8s-federate -ojsonpath={'.status.ingress[].host'}
```

#### Let RHOSO scrape OpenShift's federation endpoint

As stated in the [OpenShift documentation][ocp-federation-docs], it is recommended to limit scraping
to fewer than 1000 samples for each request, and with a maximum frequency of
once every 30 seconds.

[ocp-federation-docs]: https://docs.openshift.com/container-platform/4.17/observability/monitoring/accessing-third-party-monitoring-apis.html#monitoring-querying-metrics-by-using-the-federation-endpoint-for-prometheus_accessing-monitoring-apis-by-using-the-cli

In this example, we will only request three metrics: `kube_node_info`, `kube_persistentvolume_info`
and `cluster:master_nodes` (see the `params.match[]` query below).

While connected to the RHOSO cluster, apply this manifest:

```yaml
apiVersion: monitoring.rhobs/v1alpha1
kind: ScrapeConfig
metadata:
  labels:
    service: metricStorage
  name: sos1-federated
  namespace: openstack
spec:
  params:
    'match[]':
    - '{__name__=~"kube_node_info|kube_persistentvolume_info|cluster:master_nodes"}'
  metricsPath: '/federate'
  authorization:
    type: Bearer
    credentials:
      name: ocp-federated
      key: token
  scheme: HTTPS # or HTTP
  scrapeInterval: 30s
  staticConfigs:
  - targets:
    - prometheus-k8s-federate-openshift-monitoring.apps.openshift.example # This is the URL fetched previously
  # add a tlsConfig stanza in case the endpoint is HTTPS but uses a custom CA
```

Don't forget to make the token available as a secret (in the example above, the name is `ocp-federated`):

```bash
oc -n openstack create secret generic ocp-federated --from-literal=token=<the token fetched previously>
```

Once the new scrapeconfig propagates, the requested OpenShift metrics will be
accessible for querying in RHOSO's OpenShift UI.

## Available mappings

To query metrics and identifying resources across the stack, OpenShift exposes
helper metrics that establish a correlation between OpenStack infrastructure
resources and their representation in OpenShift.

To map **Kubernetes nodes** with **OpenStack Nova instances**:
* in the metric `kube_node_info`:
  * `node` is the Kubernetes node name
  * `provider_id` contains the identifier of the corresponding OpenStack Nova instance

To map **Kubernetes persistent volumes** with **OpenStack Cinder volume or Manila share**:
* in the metric `kube_persistentvolume_info`:
  * `persistentvolume` is the Kubernetes volume name
  * `csi_volume_handle` is the Cinder volume or Manila share identifier

### Example

By default, the Nova VMs backing the OpenShift control plane nodes are created
in a server group with policy "soft-anti-affinity". As a consequence, Nova will
create them on separate hypervisors, on a best effort basis. However, if the
state of the OpenStack cluster doesn't permit it (for example, because only two
hypervisors are available), the VMs will be created anyway.

In combination with the default soft-anti-affinity policy, it might be
interesting to set up an alert firing when a hypervisor hosts more than one
control plane node of a given OpenShift cluster, to highlight the degraded
level of high availability.

This query returns the number of OpenShift master nodes per OpenStack host:

```PromQL
sum by (vm_instance) (
  group by (vm_instance, resource) (ceilometer_cpu)
    / on (resource) group_right(vm_instance) (
      group by (node, resource) (
        label_replace(kube_node_info, "resource", "$1", "system_uuid", "(.+)")
      )
    / on (node) group_left group by (node) (
      cluster:master_nodes
    )
  )
)
```
