package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// MetaRuntimeObject allows for the generic specification of hive objects since all hive objects implement both the meta and runtime object interfaces.
type MetaRuntimeObject interface {
	metav1.Object
	runtime.Object
}
