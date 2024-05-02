package baremetal

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// BmhCacheListerWatcher is an object that wraps the listing and wrapping
// functionality for baremetal host resources.
type BmhCacheListerWatcher struct {
	Resource   dynamic.ResourceInterface
	RetryWatch bool
}

// List returns a list of baremetal hosts as dynamic objects.
func (bc BmhCacheListerWatcher) List(options metav1.ListOptions) (runtime.Object, error) {
	list, err := bc.Resource.List(context.TODO(), options)
	if apierrors.IsNotFound(err) {
		logrus.Debug("    baremetalhost resource not yet available, will retry")
		return &unstructured.UnstructuredList{}, nil
	}

	return list, err
}

// Watch starts a watch over baremetal hosts.
func (bc BmhCacheListerWatcher) Watch(options metav1.ListOptions) (watch.Interface, error) {
	w, err := bc.Resource.Watch(context.TODO(), options)
	if apierrors.IsNotFound(err) && bc.RetryWatch {
		logrus.Debug("    baremetalhost resource not yet available, will retry")
		// When the Resource isn't installed yet, we can encourage the caller to keep
		// retrying by supplying an empty watcher.  In the case of
		// UntilWithSync, the caller also checks how long it takes to create the
		// watch.  To avoid errors, we introduce an artificial delay of one
		// second.
		w := watch.NewEmptyWatch()
		time.Sleep(time.Second)
		return w, nil
	}
	return w, err
}
