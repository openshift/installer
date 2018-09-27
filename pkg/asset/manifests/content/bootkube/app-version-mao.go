package bootkube

const (
	// AppVersionMao is the constant to represent contents of App_VersionMao.yaml file
	AppVersionMao = `
---
apiVersion: tco.coreos.com/v1
kind: AppVersion
metadata:
  name: machine-api
  namespace: tectonic-system
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
