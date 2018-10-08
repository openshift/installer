package bootkube

import (
	"text/template"
)

var customTmplFuncs = template.FuncMap{
	"add": func(i, j int) int {
		return i + j
	},
}

var (
	// EtcdServiceEndpointsKubeSystem is the constant to represent contents of etcd-service-endpoints.yaml file.
	EtcdServiceEndpointsKubeSystem = template.Must(template.New("etcd-service-endpoints.yaml").Funcs(customTmplFuncs).Parse(`
apiVersion: v1
kind: Endpoints
metadata:
  name: etcd
  namespace: kube-system
  annotations:
    alpha.installer.openshift.io/dns-suffix: {{.EtcdEndpointDNSSuffix}}
subsets:
- addresses:
{{- range $idx, $member := .EtcdEndpointHostnames }}
  - ip: 192.0.2.{{ add $idx 1 }}
    hostname: {{ $member }}
{{- end }}
  ports:
  - name: etcd
    port: 2379
    protocol: TCP
`))
)
