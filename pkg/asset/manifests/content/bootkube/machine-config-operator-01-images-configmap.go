package bootkube

const (
	// MachineConfigOperator01ImagesConfigmap is the constant to represent contents of Machine_ConfigOperator01ImagesConfigmap.yaml file
	MachineConfigOperator01ImagesConfigmap = `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: machine-config-operator-images
  namespace: openshift-machine-config-operator
data:
  images.json: |
    {
      "machineConfigController": "docker.io/openshift/origin-machine-config-controller:v4.0.0",
      "machineConfigDaemon": "docker.io/openshift/origin-machine-config-daemon:v4.0.0",
      "machineConfigServer": "docker.io/openshift/origin-machine-config-server:v4.0.0"
    }
`
)
