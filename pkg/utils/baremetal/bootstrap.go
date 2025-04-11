package baremetal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	clientwatch "k8s.io/client-go/tools/watch"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/infrastructure/baremetal"
)

// WaitForBaremetalBootstrapControlPlane will watch baremetalhost resources on the bootstrap
// and wait for the control plane to finish provisioning.
func WaitForBaremetalBootstrapControlPlane(ctx context.Context, config *rest.Config, dir string) error {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("creating a baremetal client: %w", err)
	}

	r := client.Resource(baremetalhost.GroupVersion.WithResource("baremetalhosts")).Namespace("openshift-machine-api")
	blw := BmhCacheListerWatcher{
		Resource:   r,
		RetryWatch: true,
	}

	logrus.Infof("  Waiting for baremetal control plane to provision...")

	masters := map[string]baremetalhost.BareMetalHost{}

	_, withSyncErr := clientwatch.UntilWithSync(
		ctx,
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
				prev, found := masters[bmh.Name]

				if !found || bmh.Status.Provisioning.State != prev.Status.Provisioning.State {
					if bmh.Status.Provisioning.State == baremetalhost.StateNone {
						// StateNone is an empty string
						logrus.Infof("  baremetalhost: %s: uninitialized", bmh.Name)
					} else {
						logrus.Infof("  baremetalhost: %s: %s", bmh.Name, bmh.Status.Provisioning.State)
					}

					if bmh.Status.OperationalStatus == baremetalhost.OperationalStatusError {
						logrus.Warnf("  baremetalhost: %s: %s: %s", bmh.Name, bmh.Status.ErrorType, bmh.Status.ErrorMessage)
					}
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

	mastersJSON, err := json.Marshal(masters)
	if err != nil {
		return fmt.Errorf("failed to marshal masters: %w", err)
	}

	err = os.WriteFile(filepath.Join(dir, baremetal.MastersFileName), mastersJSON, 0600)
	if err != nil {
		return fmt.Errorf("failed to persist masters file to disk: %w", err)
	}

	if withSyncErr != nil {
		// wrap with ControlPlaneCreationError to trigger bootstrap log bundle gather
		return fmt.Errorf("%s: %w", asset.ControlPlaneCreationError, withSyncErr)
	}
	return nil
}
