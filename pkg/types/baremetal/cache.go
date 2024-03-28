package baremetal

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// BmhCacheListerWatcher is an object that wraps the listing and wrapping
// functionality for baremetal host resources.
type BmhCacheListerWatcher struct {
	Resource dynamic.ResourceInterface
}

// List returns a list of baremetal hosts as dynamic objects.
func (bc BmhCacheListerWatcher) List(options metav1.ListOptions) (runtime.Object, error) {
	list, err := bc.Resource.List(context.TODO(), options)
	if err != nil {
		if err.Error() == "the server could not find the requested resource" {
			return &unstructured.UnstructuredList{}, nil
		}
	}
	return list, err
}

// Watch starts a watch over baremetal hosts.
func (bc BmhCacheListerWatcher) Watch(options metav1.ListOptions) (watch.Interface, error) {
	w, err := bc.Resource.Watch(context.TODO(), options)
	if err != nil {
		if err.Error() == "the server could not find the requested resource" {
			// We can't use watch.NewEmptyWatch here because it closes too quickly.
			fake := watch.NewFake()
			defer func() {
				time.Sleep(time.Second * 2)
				fake.Stop()
			}()
			return fake, nil
		}
	}
	return w, err
}
