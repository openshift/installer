package appversions

const (
	// AppVersionKubeCore  is the variable/constant representing the contents of the respective file
	AppVersionKubeCore = `
---
apiVersion: tco.coreos.com/v1
kind: AppVersion
metadata:
  name: kube-core
  namespace: tectonic-system
  labels:
    managed-by-channel-operator: "true"
spec:
  paused: false
status:
  paused: false
upgradereq: 0
upgradecomp: 0
`
)
