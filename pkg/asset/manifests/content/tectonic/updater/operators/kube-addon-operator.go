package operators

import (
	"text/template"
)

var (
	// KubeAddonOperator  is the variable/constant representing the contents of the respective file
	KubeAddonOperator = template.Must(template.New("kube-addon-operator.yaml").Parse(`
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: kube-addon-operator
  namespace: tectonic-system
  labels:
    k8s-app: kube-addon-operator
    managed-by-channel-operator: "true"
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: kube-addon-operator
  template:
    metadata:
      labels:
        k8s-app: kube-addon-operator
        tectonic-app-version-name: kube-addon
    spec:
      containers:
        - name: kube-addon-operator
          image: {{.KubeAddonOperatorImage}}
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
      imagePullSecrets:
        - name: coreos-pull-secret
      nodeSelector:
        node-role.kubernetes.io/master: ""
      restartPolicy: Always
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      tolerations:
        - key: "node-role.kubernetes.io/master"
          operator: "Exists"
          effect: "NoSchedule"
      volumes:
        - name: cluster-config
          configMap:
            name: cluster-config-v1
            items:
              - key: addon-config
                path: addon-config
`))
)
