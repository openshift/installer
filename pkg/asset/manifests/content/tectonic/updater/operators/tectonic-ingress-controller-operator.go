package operators

import (
	"text/template"
)

var (
	// TectonicIngressControllerOperator  is the variable/constant representing the contents of the respective file
	TectonicIngressControllerOperator = template.Must(template.New("tectonic-ingress-controller-operator.yaml").Parse(`
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: tectonic-ingress-controller-operator
  namespace: openshift-ingress
  labels:
    k8s-app: tectonic-ingress-controller-operator
    managed-by-channel-operator: "true"
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: tectonic-ingress-controller-operator
  template:
    metadata:
      labels:
        k8s-app: tectonic-ingress-controller-operator
        tectonic-app-version-name: tectonic-ingress
    spec:
      containers:
        - name: tectonic-ingress-controller-operator
          image: {{.TectonicIngressControllerOperatorImage}}
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
      serviceAccount: tectonic-ingress-controller-operator
      tolerations:
        - key: "node-role.kubernetes.io/master"
          operator: "Exists"
          effect: "NoSchedule"
      volumes:
        - name: cluster-config
          configMap:
            name: cluster-config-v1
            items:
              - key: ingress-config
                path: ingress-config
`))
)
