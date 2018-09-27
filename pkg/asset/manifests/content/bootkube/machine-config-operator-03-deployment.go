package bootkube

import (
	"text/template"
)

var (
	// MachineConfigOperator03Deployment is the constant to represent contents of machine_configoperator03deployment.yaml file
	MachineConfigOperator03Deployment = template.Must(template.New("machine-config-operator-03-deployment.yaml").Parse(`
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: machine-config-operator
  namespace: openshift-machine-config-operator
  labels:
    k8s-app: machine-config-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: machine-config-operator
  template:
    metadata:
      labels:
        k8s-app: machine-config-operator
    spec:
      containers:
        - name: machine-config-operator
          image: {{.MachineConfigOperatorImage}}
          args:
            - "start"
            - "--images-json=/etc/mco/images/images.json"
          resources:
            limits:
              cpu: 20m
              memory: 50Mi
            requests:
              cpu: 20m
              memory: 50Mi
          volumeMounts:
            - name: root-ca
              mountPath: /etc/ssl/kubernetes/ca.crt
            - name: etcd-ca
              mountPath: /etc/ssl/etcd/ca.crt
            - name: images
              mountPath: /etc/mco/images
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
        - name: images
          configMap:
            name: machine-config-operator-images
        - name: etcd-ca
          hostPath:
            path: /etc/ssl/etcd/ca.crt
        - name: root-ca
          hostPath:
            path: /etc/kubernetes/ca.crt
`))
)
