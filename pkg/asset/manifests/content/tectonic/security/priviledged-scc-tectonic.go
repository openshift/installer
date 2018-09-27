package security

const (
	// PriviledgedSccTectonic  is the variable/constant representing the contents of the respective file
	PriviledgedSccTectonic = `
---
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  annotations:
    kubernetes.io/description: "privileged-tectonic temporarily for running tectonic assets."
  name: privileged-tectonic
allowHostDirVolumePlugin: true
allowHostIPC: true
allowHostNetwork: true
allowHostPID: true
allowHostPorts: true
allowPrivilegedContainer: true
allowedCapabilities:
  - "*"
fsGroup:
  type: RunAsAny
groups:
  - system:serviceaccounts:tectonic-system
  - system:serviceaccounts:openshift-ingress
readOnlyRootFilesystem: false
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: RunAsAny
seccompProfiles:
  - "*"
supplementalGroups:
  type: RunAsAny
users: []
volumes:
  - "*"
`
)
