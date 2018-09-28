package bootkube

const (
	// MachineAPIOperator is the constant to represent contents of Machine_Api_Operator.yaml file
	MachineAPIOperator = `
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: machine-api-operator
  namespace: kube-system
  labels:
    k8s-app: machine-api-operator
    managed-by-channel-operator: "true"
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: machine-api-operator
  template:
    metadata:
      labels:
        k8s-app: machine-api-operator
        tectonic-app-version-name: machine-api
    spec:
      containers:
        - name: machine-api-operator
          image: quay.io/coreos/machine-api-operator:b6a04c2
          command:
            - "/machine-api-operator"
          resources:
            limits:
              cpu: 20m
              memory: 50Mi
            requests:
              cpu: 20m
              memory: 50Mi
          volumeMounts:
            - name: cluster-config
              mountPath: /etc/mao-config
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
              - key: mao-config
                path: config
`
)
