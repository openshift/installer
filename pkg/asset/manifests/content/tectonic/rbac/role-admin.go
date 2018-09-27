package rbac

const (
	// RoleAdmin  is the variable/constant representing the contents of the respective file
	RoleAdmin = `
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: admin
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
  - nonResourceURLs: ["*"]
    verbs: ["*"]
`
)
