package bootkube

const (
	// OpenshiftMachineConfigOperator is the constant to represent contents of Openshift_MachineConfigOperator.yaml file
	OpenshiftMachineConfigOperator = `
apiVersion: v1
kind: Namespace
metadata:
  name: openshift-machine-config-operator
  labels:
    name: openshift-machine-config-operator
    openshift.io/run-level: "1"
`
)
