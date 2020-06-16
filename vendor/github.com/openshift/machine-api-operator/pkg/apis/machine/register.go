// Generate deepcopy for apis
//go:generate go run ../../../vendor/sigs.k8s.io/controller-tools/cmd/controller-gen paths=./... object:headerFile=../../../hack/boilerplate.go.txt,year=2019
// Ensure generated code is goimports compliant
//go:generate goimports -w ./v1beta1/zz_generated.deepcopy.go

package machine

const GroupName = "machine.openshift.io"
