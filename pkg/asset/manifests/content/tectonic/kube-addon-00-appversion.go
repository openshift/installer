package tectonic

const (
	// AppVersionKubeAddon  is the variable/constant representing the contents of the respective file
	AppVersionKubeAddon = `
---
apiVersion: tco.coreos.com/v1
kind: AppVersion
metadata:
  name: kube-addon
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
