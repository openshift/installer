package rbac

const (
	// BindingAdmin  is the variable/constant representing the contents of the respective file
	BindingAdmin = `
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: admin-user
subjects:
  - kind: ServiceAccount
    namespace: tectonic-system
    name: default
  - kind: ServiceAccount
    namespace: openshift-ingress
    name: tectonic-ingress-controller-operator
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
`
)
