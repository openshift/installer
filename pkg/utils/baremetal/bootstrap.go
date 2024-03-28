package baremetal

import (
	"context"
	"time"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	clientwatch "k8s.io/client-go/tools/watch"

	"github.com/openshift/installer/pkg/types/baremetal"
)

func WaitForBaremetalBootstrapControlPlane(ctx context.Context, config *rest.Config) error {
	timeout := 30 * time.Minute

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a baremetal client")
	}

	r := client.Resource(baremetalhost.GroupVersion.WithResource("baremetalhosts")).Namespace("openshift-machine-api")
	blw := baremetal.BmhCacheListerWatcher{
		Resource: r,
	}

	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for baremetal control plane to provision...",
		timeout, untilTime.Format(time.Kitchen), timezone)

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	masters := map[string]baremetalhost.BareMetalHost{}

	_, err = clientwatch.UntilWithSync(
		waitCtx,
		blw,
		&unstructured.Unstructured{},
		nil,
		func(event watch.Event) (bool, error) {
			switch event.Type {
			case watch.Added, watch.Modified:
			default:
				return false, nil
			}

			bmh := &baremetalhost.BareMetalHost{}

			unstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(event.Object)
			if err != nil {
				return false, err
			}

			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstr, bmh); err != nil {
				logrus.Error("failed to convert to bmh", err)
				return false, err
			}

			role, found := bmh.Labels["installer.openshift.io/role"]

			if found && role == "control-plane" {
				if bmh.Status.Provisioning.State == baremetalhost.StateNone {
					// StateNone is an empty string
					logrus.Infof("  baremetalhost: %s: uninitialized", bmh.Name)
				} else {
					logrus.Infof("  baremetalhost: %s: %s", bmh.Name, bmh.Status.Provisioning.State)
				}

				if bmh.Status.OperationalStatus == baremetalhost.OperationalStatusError {
					logrus.Infof("  baremetalhost: %s: error: %s", bmh.Name, bmh.Status.ErrorMessage)
				}
				masters[bmh.Name] = *bmh
			}

			if len(masters) == 0 {
				return false, nil
			}

			for _, master := range masters {
				if master.Status.Provisioning.State != baremetalhost.StateProvisioned {
					return false, nil
				}
			}

			return true, nil
		},
	)

	return err
}
