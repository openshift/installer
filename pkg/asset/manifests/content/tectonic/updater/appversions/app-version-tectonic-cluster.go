package appversions

import (
	"text/template"
)

var (
	// AppVersionTectonicCluster  is the variable/constant representing the contents of the respective file
	AppVersionTectonicCluster = template.Must(template.New("app-version-tectonic-cluster.yaml").Parse(`
apiVersion: tco.coreos.com/v1
kind: AppVersion
metadata:
  name: tectonic-cluster
  namespace: tectonic-system
  labels:
    managed-by-channel-operator: "true"
spec:
  desiredVersion: {{.TectonicVersion}}
  paused: false
status:
  currentVersion: {{.TectonicVersion}}
  paused: false
`))
)
