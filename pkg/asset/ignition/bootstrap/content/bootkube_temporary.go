package content

import "text/template"

// KubeProxyBootkubeManifests is a map of manifests needed by kube-proxy to install.
// TODO: This must move to networking operator renderer.
var KubeProxyBootkubeManifests = map[string]string{
	"kube-proxy-kube-system-rbac-role-binding.yaml": bootkubeKubeSystemRBACRoleBinding,
	"kube-proxy-role-binding.yaml":                  bootkubeKubeProxyRoleBinding,
	"kube-proxy-service-account.yaml":               bootkubeKubeProxySA,
	"kube-proxy-daemonset.yaml":                     bootkubeKubeProxyDaemonset,
}

// KubeDNSBootkubeManifests is a map of manifests needed by kube-dns to install.
// TODO: This must move to networking operator renderer.
var KubeDNSBootkubeManifests = map[string]string{
	"kube-dns-deployment.yaml": bootkubeKubeDNSDeployment,
}

// BootkubeKubeDNSService is a template for kube-dns service.
var (
	BootkubeKubeDNSService = template.Must(template.New("bootkube.sh").Parse(`
apiVersion: v1
kind: Service
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: kube-dns
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: KubeDNS
spec:
  selector:
    k8s-app: kube-dns
  clusterIP: {{.ClusterDNSIP}}
  ports:
  - name: dns
    port: 53
    protocol: UDP
    targetPort: 53
  - name: dns-tcp
    port: 53
    protocol: TCP
    targetPort: 53
`))

	BootkubeKubeProxyKubeConfig = template.Must(template.New("kube-proxy-kubeconfig").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: kube-proxy-kubeconfig
  namespace: kube-system
data:
  kubeconfig: {{ .AdminKubeConfigBase64 }}
`))
)

const (
	bootkubeKubeSystemRBACRoleBinding = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:default-sa
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: default
  namespace: kube-system
`

	bootkubeKubeProxySA = `
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: kube-system
  name: kube-proxy
`

	bootkubeKubeProxyDaemonset = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-proxy
    tier: node
  name: kube-proxy
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: kube-proxy
      tier: node
  template:
    metadata:
      labels:
        k8s-app: kube-proxy
        tier: node
    spec:
      containers:
      - command:
        - ./hyperkube
        - proxy
        - --cluster-cidr=10.3.0.0/16
        - --hostname-override=$(NODE_NAME)
        - --kubeconfig=/etc/kubernetes/kubeconfig
        - --proxy-mode=iptables
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        image: quay.io/coreos/hyperkube:v1.9.3_coreos.0
        name: kube-proxy
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /etc/ssl/certs
          name: ssl-certs-host
          readOnly: true
        - mountPath: /etc/kubernetes
          name: kubeconfig
          readOnly: true
      hostNetwork: true
      serviceAccountName: kube-proxy
      tolerations:
      - operator: Exists
      volumes:
      - hostPath:
          path: /etc/ssl/certs
        name: ssl-certs-host
      - name: kubeconfig
        secret:
          defaultMode: 420
          secretName: kube-proxy-kubeconfig
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
`

	bootkubeKubeProxyRoleBinding = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kube-proxy
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:node-proxier # Automatically created system role.
subjects:
- kind: ServiceAccount
  name: kube-proxy
  namespace: kube-system
`
	bootkubeKubeDNSDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: kube-dns
    kubernetes.io/cluster-service: "true"
spec:
  # replicas: not specified here:
  # 1. In order to make Addon Manager do not reconcile this replicas parameter.
  # 2. Default is 1.
  # 3. Will be tuned in real time if DNS horizontal auto-scaling is turned on.
  strategy:
    rollingUpdate:
      maxSurge: 10%
      maxUnavailable: 0
  selector:
    matchLabels:
      k8s-app: kube-dns
  template:
    metadata:
      labels:
        k8s-app: kube-dns
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      containers:
      - name: kubedns
        image: "gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"
        resources:
          # TODO: Set memory limits when we've profiled the container for large
          # clusters, then set request = limit to keep this container in
          # guaranteed class. Currently, this container falls into the
          # "burstable" category so the kubelet doesn't backoff from restarting it.
          limits:
            memory: 170Mi
          requests:
            cpu: 100m
            memory: 70Mi
        livenessProbe:
          httpGet:
            path: /healthcheck/kubedns
            port: 10054
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
            scheme: HTTP
          # we poll on pod startup for the Kubernetes master service and
          # only setup the /readiness HTTP server once that's available.
          initialDelaySeconds: 3
          timeoutSeconds: 5
        args:
        - --domain=cluster.local.
        - --dns-port=10053
        - --config-dir=/kube-dns-config
        - --v=2
        env:
        - name: PROMETHEUS_PORT
          value: "10055"
        ports:
        - containerPort: 10053
          name: dns-local
          protocol: UDP
        - containerPort: 10053
          name: dns-tcp-local
          protocol: TCP
        - containerPort: 10055
          name: metrics
          protocol: TCP
        volumeMounts:
        - name: kube-dns-config
          mountPath: /kube-dns-config
      - name: dnsmasq
        image: "gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"
        livenessProbe:
          httpGet:
            path: /healthcheck/dnsmasq
            port: 10054
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
        args:
        - -v=2
        - -logtostderr
        - -configDir=/etc/k8s/dns/dnsmasq-nanny
        - -restartDnsmasq=true
        - --
        - -k
        - --cache-size=1000
        - --log-facility=-
        - --server=/cluster.local/127.0.0.1#10053
        - --server=/in-addr.arpa/127.0.0.1#10053
        - --server=/ip6.arpa/127.0.0.1#10053
        ports:
        - containerPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          name: dns-tcp
          protocol: TCP
        # see: https://github.com/kubernetes/kubernetes/issues/29055 for details
        resources:
          requests:
            cpu: 150m
            memory: 20Mi
        volumeMounts:
        - name: kube-dns-config
          mountPath: /etc/k8s/dns/dnsmasq-nanny
      - name: sidecar
        image: "gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"
        livenessProbe:
          httpGet:
            path: /metrics
            port: 10054
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
        args:
        - --v=2
        - --logtostderr
        - --probe=kubedns,127.0.0.1:10053,kubernetes.default.svc.cluster.local,5,SRV
        - --probe=dnsmasq,127.0.0.1:53,kubernetes.default.svc.cluster.local,5,SRV
        ports:
        - containerPort: 10054
          name: metrics
          protocol: TCP
        resources:
          requests:
            memory: 20Mi
            cpu: 10m
      dnsPolicy: Default  # Don't use cluster DNS.
      tolerations:
      - key: "CriticalAddonsOnly"
        operator: "Exists"
      - key: "node-role.kubernetes.io/master"
        operator: "Exists"
        effect: "NoSchedule"
      nodeSelector:
        node-role.kubernetes.io/master: ""
      volumes:
      - name: kube-dns-config
        configMap:
          name: kube-dns
          optional: true
`
)
