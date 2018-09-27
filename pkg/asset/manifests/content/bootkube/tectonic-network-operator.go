package bootkube

import (
	"text/template"
)

var (
	// TectonicNetworkOperator represents the template variable for tectonic-network-operator.yaml file
	TectonicNetworkOperator = template.Must(template.New("tectonic-network-operator.yaml").Parse(`
---
apiVersion: apps/v1beta2
kind: DaemonSet
metadata:
  name: tectonic-network-operator
  namespace: kube-system
  labels:
    k8s-app: tectonic-network-operator
    managed-by-channel-operator: "true"
spec:
  selector:
    matchLabels:
      k8s-app: tectonic-network-operator
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: tectonic-network-operator
        tectonic-app-version-name: tectonic-network
    spec:
      containers:
        - name: tectonic-network-operator
          image: {{.TectonicNetworkOperatorImage}}
          resources:
            limits:
              cpu: 20m
              memory: 50Mi
            requests:
              cpu: 20m
              memory: 50Mi
          volumeMounts:
            - name: cluster-config
              mountPath: /etc/cluster-config
      hostNetwork: true
      restartPolicy: Always
      imagePullSecrets:
        - name: coreos-pull-secret
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      volumes:
        - name: cluster-config
          configMap:
            name: cluster-config-v1
            items:
              - key: network-config
                path: network-config
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
        - key: "node-role.kubernetes.io/master"
          operator: "Exists"
          effect: "NoSchedule"
    updateStrategy:
      rollingUpdate:
        maxUnavailable: 1
      type: RollingUpdate
`))
)
