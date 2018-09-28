package appversions

const (
	// AppVersionTectonicIngress  is the variable/constant representing the contents of the respective file
	AppVersionTectonicIngress = `
---
apiVersion: tco.coreos.com/v1
kind: AppVersion
metadata:
  name: tectonic-ingress
  namespace: tectonic-system
  labels:
    managed-by-channel-operator: "true"
spec:
  desiredVersion:
  paused: false
status:
  paused: false
upgradereq: 1
upgradecomp: 0
`
)
