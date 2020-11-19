// Generate deepcopy for apis
//go:generate go run ../../../vendor/sigs.k8s.io/controller-tools/cmd/controller-gen paths=./... object:headerFile=../../../hack/boilerplate.go.txt,year=2019
// Ensure generated code is goimports compliant
//go:generate goimports -w ./v1beta1/zz_generated.deepcopy.go

// Package apis contains Kubernetes API groups.
package apis

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
