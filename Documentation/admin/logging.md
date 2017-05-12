# Aggregated Logging in Tectonic

Tectonic does not preconfigure any particular aggregated logging stack. Tectonic recommends several example logging configurations that can be customized for site requirements. The recommended logging setup uses Fluentd to retrieve logs on each node and forward them to a log storage backend. The Tectonic examples use Elasticsearch for log storage. Elasticsearch can be replaced by any destination Fluentd supports with an Output plugin. For a list of Fluentd plugins, check [http://www.fluentd.org/plugins/all](http://www.fluentd.org/plugins/all)

If you want to run these examples locally, all of the files mentioned are available in the [Github repo for the Tectonic Installer][logging-config-files].

## Overview

### Prerequisites

- Kubernetes 1.6+
- `kubectl` configured
  - If you need to configure `kubectl`, follow the instructions on [configuring credentials][configuring-credentials].
- An Elasticsearch cluster, or other log storage destination
  - See [customizing log destination][customizing-log-destination] for a list of available log destinations with pre-built images available.

### Fluentd

Fluentd runs as a DaemonSet on all nodes, including masters, and is configured using a ConfigMap of Fluentd config files which define how to collect the logs.

The [provided configuration][fluentd-config] is setup to do the following:

- Read host system logs using the Journal
- Read Kubernetes containers logs
- Read Kubernetes API server audit logs
- Tag logs based on metadata such as container name
- Enable Monitoring (exposed as Prometheus metrics)
- Parse known log formats based on their tag
- Send logs to Elasticsearch

### Elasticsearch

In this setup, we will not be configuring or setting up Elasticsearch, or Kibana. If you are looking for something to get started with, https://github.com/pires/kubernetes-elasticsearch-cluster is a good reference, and has examples of deploying an Elasticsearch cluster on Kubernetes while following best practices.

If you want to customize your Elasticsearch output configuration or look at examples using different storage destinations, see the [customizing log destination][customizing-log-destination] document.

## Deploying Fluentd

First create the `logging` namespace for all of our resources to live in:

```
$ kubectl create ns logging
```

Then setup all the service account and roles that Fluentd needs to query Kubernetes for metadata about the containers logs it's watching:

```
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/fluentd-service-account.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/fluentd-role.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/fluentd-role-binding.yaml
```

Next deploy Fluentd:

```
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/fluentd-configmap.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/fluentd-svc.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/fluentd-ds.yaml
```

Finally, watch for the pods to become ready:

```
$ kubectl get pods --watch --namespace logging
```

Once all the pods are ready, everything should be functioning. To double check, use `kubectl logs` on one of the pods listed above to make sure there aren't any errors, and that Fluentd is able to send logs to where you've configured.


### (Optional) Deploy Prometheus to monitor Fluentd

Tectonic includes the [promtheus-operator][prometheus-operator] in installations by default. This operator can be used to create additional instances of Prometheus to monitor your apps.

If you wish to enable Prometheus monitoring of your Fluentd pods, run the following commands:

```
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/monitoring/prometheus-fluentd-role-binding.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/monitoring/prometheus-fluentd-service-account.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/monitoring/prometheus-fluentd.yaml
$ kubectl create -f https://coreos.com/tectonic/docs/latest/files/logging/monitoring/fluentd-prometheus-service-monitor.yaml
```

## Customization

If you would like to make customizations to your logging setup, view our doc on [customizing Fluentd][customizing-fluentd].


[fluentd-config]: ../files/logging/fluentd-configmap.yaml
[fluentd-ds]: ../files/logging/fluentd-ds.yaml
[configuring-credentials]: ../tutorials/first-app.md#configuring-credentials
[logging-config-files]: https://github.com/coreos/tectonic-installer/tree/master/Documentation/files/logging
[prometheus-operator]: https://github.com/coreos/prometheus-operator
[customizing-fluentd]: logging-customization.md
[customizing-log-destination]: logging-destination.md
