package tectonic

import (
	"text/template"
)

var (
	// KubeCoreOperator  is the variable/constant representing the contents of the respective file
	KubeCoreOperator = template.Must(template.New("kube-core-00-operator.yaml").Parse(`
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: kube-core-operator
  namespace: kube-system
  labels:
    k8s-app: kube-core-operator
    managed-by-channel-operator: "true"
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: kube-core-operator
  template:
    metadata:
      labels:
        k8s-app: kube-core-operator
        tectonic-app-version-name: kube-core
    spec:
      containers:
        - name: kube-core-operator
          image: {{.KubeCoreOperatorImage}}
          imagePullPolicy: Always
          args:
            - --config=/etc/cluster-config/kco-config.yaml
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
              - key: kco-config
                path: kco-config.yaml
`))
)
