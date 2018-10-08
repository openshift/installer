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
spec:
  clusterIP: None
  ports:
  - name: etcd
    port: 2379
    protocol: TCP
  - name: legacy-etcd
    port: 4001
    targetPort: 2379
    protocol: TCP
`
)
