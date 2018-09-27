package bootkube

const (
	// MachineConfigOperator02Rbac is the constant to represent contents of manifest file machine-config-operator-02-rbac.yaml
	MachineConfigOperator02Rbac = `
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: default-account-openshift-machine-config-operator
subjects:
  - kind: ServiceAccount
    name: default
    namespace: openshift-machine-config-operator
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
`
)
