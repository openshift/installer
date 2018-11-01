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

// BootkubeKubeProxyKubeConfig is a template for kube-proxy-kubeconfig secret.
var (
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
)
