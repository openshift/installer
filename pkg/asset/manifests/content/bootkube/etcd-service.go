package bootkube

const (
	// EtcdServiceKubeSystem is the constant to represent contents of etcd-service.yaml file
	EtcdServiceKubeSystem = `
---
apiVersion: v1
kind: Service
metadata:
  name: etcd
  namespace: kube-system
  labels:
    # this label is used to indicate that it should be scaped by prometheus
    k8s-app: etcd
spec:
  clusterIP: None
  ports:
  - name: etcd
    port: 2379
    protocol: TCP
`
)
