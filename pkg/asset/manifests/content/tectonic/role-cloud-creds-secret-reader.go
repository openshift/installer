package tectonic

import (
	"text/template"
)

var (
	// RoleCloudCredsSecretReader is the variable to represent contents of corresponding file
	RoleCloudCredsSecretReader = template.Must(template.New("role-cloud-creds-secret-reader.yaml").Parse(`
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  namespace: kube-system
{{- if .CloudCreds.AWS}}
  name: aws-creds-secret-reader
{{- else if .CloudCreds.OpenStack}}
  name: openstack-creds-secret-reader
{{- end}}
rules:
- apiGroups: [""]
  resources: ["secrets"]
{{- if .CloudCreds.AWS}}
  resourceNames: ["aws-creds"]
{{- else if .CloudCreds.OpenStack}}
  resourceNames: ["openstack-creds"]
{{- end}}
  verbs: ["get"]
`))
)
