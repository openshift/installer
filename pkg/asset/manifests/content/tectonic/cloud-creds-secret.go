package tectonic

import (
	"text/template"
)

var (
	// CloudCredsSecret is the constant to represent contents of corresponding yaml file
	CloudCredsSecret = template.Must(template.New("cloud-creds-secret.yaml").Parse(`
---
kind: Secret
apiVersion: v1
metadata:
  namespace: kube-system
{{- if .CloudCreds.AWS}}
  name: aws-creds
{{- else if .CloudCreds.OpenStack}}
  name: openstack-creds
{{- end}}
data:
{{- if .CloudCreds.AWS}}
  aws_access_key_id: {{.CloudCreds.AWS.Base64encodeAccessKeyID}}
  aws_secret_access_key: {{.CloudCreds.AWS.Base64encodeSecretAccessKey}}
{{- else if .CloudCreds.OpenStack}}
  clouds.yaml: {{.CloudCreds.OpenStack.Base64encodeCloudCreds}}
{{- end}}
`))
)
