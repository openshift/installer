package bootkube

const (
	// AppVersionTectonicNetwork is the constant to represent contents of App_VersionTectonicNetwork.yaml file
	AppVersionTectonicNetwork = `
---
apiVersion: tco.coreos.com/v1
kind: AppVersion
metadata:
  name: tectonic-network
  namespace: kube-system
  labels:
    managed-by-channel-operator: "true"
spec:
  desiredVersion:
  paused: false
status:
  currentVersion:
  paused: false
upgradereq: 1
upgradecomp: 0
`
)
