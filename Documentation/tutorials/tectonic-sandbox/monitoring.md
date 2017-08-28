# Monitoring an application with Prometheus

Prometheus is an open-source metrics and monitoring application that is shipped as an open cloud service in Tectonic. The default Prometheus instance is used to monitor the Tectonic control plane, and a cluster-wide AlertManager aggregates the alerts across multiple Prometheus instances for greater visibility. Prometheus instances can be created through the Tectonic Console or using `kubectl`.

In this tutorial we will:
* Deploy a Prometheus instance to monitor the application created in [Deploying an application on Tectonic][first-app].
* Observe alerts firing as the application is scaled up and down.

## Objects

Prometheus instances are comprised of a few Kubernetes objects:

1. The Prometheus instance.
2. A ServiceMonitor.
3. A ConfigMap to store alerting rules.
4. A ServiceAccount and a ClusterRoleBinding to provide access to `kube-state-metrics`.
5. A Service and an Ingress to access the Prometheus UI.

## Deploying a Prometheus Instance

Deploying a Prometheus instance requires several Kubernetes objects. First, create a `ServiceAccount` and a `ClusterRoleBinding` to ensure that Prometheus can scrape `kube-state-metrics` to determine the number of replicas deployed by an application.

Copy the following YAML into a file named `prometheus-serviceaccount.yaml`:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus-hello
  namespace: default
```

Then, copy the following YAML into a file named `prometheus-crb.yaml`:

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: prometheus-hello
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus-k8s
subjects:
- kind: ServiceAccount
  name: prometheus-hello
  namespace: default
```

Next, use the open cloud service to deploy the Prometheus instance. Prometheus as an open cloud service will manage the underlying `StatefulSet` and `Pod`s for Prometheus.

Copy the following into a file named `prometheus.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: Prometheus
metadata:
  name: hello
  labels:
    prometheus: hello
  namespace: default
spec:
  replicas: 1
  version: v1.7.1
  serviceAccountName: prometheus-hello
  serviceMonitorSelector:
    matchLabels:
      k8s-app: kube-state-metrics
  ruleSelector:
    matchLabels:
      prometheus: hello
  resources:
    requests:
      memory: 400Mi
  alerting:
    alertmanagers:
      - namespace: tectonic-system
        name: alertmanager-main
        port: web
```

To tell this new Prometheus instance what to monitor, create the `ServiceMonitor` and the `ConfigMap`.

First, copy the following into a file named `prometheus-servicemonitor.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  name: kube-state-metrics
  namespace: default
  labels:
    k8s-app: kube-state-metrics
spec:
  endpoints:
    - honorLabels: true
      interval: 30s
      port: http-metrics
      targetPort: 0
  jobLabel: simple-app
  namespaceSelector:
    matchNames:
      - tectonic-system
  selector:
    matchLabels:
      k8s-app: kube-state-metrics
```

Next, copy the following into a file named `prometheus-configmap.yaml`:

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: simple-prom
  labels:
    prometheus: hello
  namespace: default
data:
  alerting.rules: |
    # Alert if deployment missing
    ALERT SimpleDeploymentMissing
      IF kube_deployment_status_replicas{deployment="simple-deployment"} < 3
      FOR 10m
      LABELS {severity="critical"}
      ANNOTATIONS {description="Prometheus could not find the Simple Deployment"}
```

Next, create a means to access the Prometheus UI. Copy the following into a file named `prometheus-service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: prometheus-hello-world
  namespace: default
spec:
  type: LoadBalancer
  ports:
  - name: web
    port: 9090
    protocol: TCP
    targetPort: 9090
  selector:
    prometheus: hello

```

Finally, create `Ingress` to enable access to Prometheus from outside the cluster. Copy the following into a file named `prometheus-ingress.yaml`:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: prom-ingress
  annotations:
    kubernetes.io/ingress.class: "tectonic"
    ingress.kubernetes.io/rewrite-target: /
    ingress.kubernetes.io/ssl-redirect: "true"
    ingress.kubernetes.io/use-port-in-redirects: "true"
spec:
  rules:
    - host: prometheus.ingress.tectonicsandbox.com
      http:
        paths:
          - path: /
            backend:
              serviceName: prometheus-hello-world
              servicePort: 9090
```

Use `kubectl` to push all of these files up to the cluster at the same time. Confirm that the working directory contains only the YAML files described above, then create all required objects with a single command:

```sh
$ kubectl create -f .
configmap "simple-prom" created
clusterrolebinding "prometheus-hello" created
service "prometheus-hello-world" created
ingress "prom-ingress" created
serviceaccount "prometheus-hello" created
servicemonitor "kube-state-metrics" created
prometheus "hello" created
```

## Accessing Prometheus

Now, with the Prometheus instance online, and configured to alert based on the number of replicas of the app, access the Prometheus UI to watch the configured alerts fire.

Use the URL provided by `Ingress` to access the console: [prometheus.ingress.tectonicsandbox.com/alerts][prom-ingress].

## Triggering an alert

The alert configured in Prometheus will fire if the `replicas` for the `simple-deployment` is fewer than 3. Scale down the `simple-deployment`:

```sh
$ kubectl scale deployment/simple-deployment --replicas 2
deployment "simple-deployment" scaled
```

By default, Prometheus scrapes the `kube-state-metrics` every 30 seconds. After scaling the deployment, wait 30 seconds, then refresh [prometheus.ingress.tectonicsandbox.com/alerts][prom-ingress] to see that the alert is now firing, as the `simple-deployment` has fewer than 3 replicas.

Scale the `simple-deployment` back up to 3:

```sh
$ kubectl scale deployment/simple-deployment --replicas 3
deployment "simple-deployment" scaled
```

[prom-ingress]: http://prometheus.ingress.tectonicsandbox.com/alerts
[first-app]: first-app.md
