package kubevirt

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	apilabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"

	ickubevirt "github.com/openshift/installer/pkg/asset/installconfig/kubevirt"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteFunc is the interface a function needs to implement to be delete resources.
type deleteFunc func(ctx context.Context, namespace string, listOpts metav1.ListOptions, kubevirtClient ickubevirt.Client) error

// ClusterUninstaller holds the Metadata info needed to delete the tenantCluster resources from the infraCluster.
type ClusterUninstaller struct {
	Metadata types.ClusterMetadata
	Logger   logrus.FieldLogger
}

// New returns KubeVirt Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Metadata: *metadata,
		Logger:   logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (uninstaller *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	ctx := context.Background()
	namespace := uninstaller.Metadata.Kubevirt.Namespace

	listOpts := metav1.ListOptions{LabelSelector: apilabels.FormatLabels(uninstaller.Metadata.Kubevirt.Labels)}
	kubevirtClient, err := ickubevirt.NewClient()
	if err != nil {
		return nil, err
	}

	deleteFuncs := []deleteFunc{uninstaller.deleteAllVMs, uninstaller.deleteAllDVs, uninstaller.deleteAllSecrets}
	var results = make(chan error, len(deleteFuncs))
	for i, del := range deleteFuncs {
		go func(index int, delFunc deleteFunc) {
			for {
				err := delFunc(ctx, namespace, listOpts, kubevirtClient)
				if err == nil {
					results <- err
					break
				}
				time.Sleep(10 * time.Second)
			}

		}(i, del)
	}

	var resultMsg = ""
	for i := len(deleteFuncs); i > 0; i-- {
		err := <-results
		if err != nil {
			resultMsg = fmt.Sprintf("%s, %s", resultMsg, err.Error())
		}
	}
	if resultMsg != "" {
		return nil, fmt.Errorf("destroy finished with errors: %s", resultMsg)
	}
	return nil, nil
}

func (uninstaller *ClusterUninstaller) deleteAllVMs(ctx context.Context, namespace string, listOpts metav1.ListOptions, kubevirtClient ickubevirt.Client) error {
	vmList, err := kubevirtClient.ListVirtualMachine(ctx, namespace, listOpts)
	if err != nil {
		uninstaller.Logger.Errorf("failed to delete VirtualMachines: %s", err)
		return fmt.Errorf("failed to list VMs")
	}

	if len(vmList.Items) == 0 {
		uninstaller.Logger.Infof("There is no Virtual Machine to delete in namespace %s", namespace)
		return nil
	}
	uninstaller.Logger.Infof("Found %d Virtual Machines to delete in namespace %s", len(vmList.Items), namespace)
	for _, vm := range vmList.Items {
		uninstaller.Logger.Infof("Delete Virtual Machine %s from Namespace %s", vm.Name, namespace)
		if err := kubevirtClient.DeleteVirtualMachine(ctx, namespace, vm.Name); err != nil {
			uninstaller.Logger.Errorf("failed to delete VirtualMachines: %s", err)
			return fmt.Errorf("failed to delete VM")
		}
		if err := uninstaller.exponentialBackoff("Virtual Machine", func() error {
			_, err := kubevirtClient.GetVirtualMachine(ctx, namespace, vm.Name)
			return err
		}); err != nil {
			uninstaller.Logger.Errorf("failed to delete VirtualMachines: %s", err)
			return fmt.Errorf("failed to validate VM deleted")
		}
	}
	return nil
}

func (uninstaller *ClusterUninstaller) deleteAllDVs(ctx context.Context, namespace string, listOpts metav1.ListOptions, kubevirtClient ickubevirt.Client) error {
	dvList, err := kubevirtClient.ListDataVolume(ctx, namespace, listOpts)
	if err != nil {
		uninstaller.Logger.Errorf("failed to delete DataVolumes: %s", err)
		return fmt.Errorf("failed to list DVs")
	}

	if len(dvList.Items) == 0 {
		uninstaller.Logger.Infof("There is no Data Volume to delete in namespace %s", namespace)
		return nil
	}
	uninstaller.Logger.Infof("Found %d Data Volumes to delete in namespace %s", len(dvList.Items), namespace)
	for _, dv := range dvList.Items {
		uninstaller.Logger.Infof("Delete Data Volume %s from Namespace %s", dv.Name, namespace)
		if err := kubevirtClient.DeleteDataVolume(ctx, namespace, dv.Name); err != nil {
			uninstaller.Logger.Errorf("failed to delete DataVolumes: %s", err)
			return fmt.Errorf("failed to delete DV")
		}
		if err := uninstaller.exponentialBackoff("Data Volume", func() error {
			_, err := kubevirtClient.GetDataVolume(ctx, namespace, dv.Name)
			return err
		}); err != nil {
			uninstaller.Logger.Errorf("failed to delete DataVolumes: %s", err)
			return fmt.Errorf("failed to validate DV deleted")
		}
	}
	return nil
}

func (uninstaller *ClusterUninstaller) deleteAllSecrets(ctx context.Context, namespace string, listOpts metav1.ListOptions, kubevirtClient ickubevirt.Client) error {
	secretList, err := kubevirtClient.ListSecret(ctx, namespace, listOpts)
	if err != nil {
		uninstaller.Logger.Errorf("failed to delete Secrets: %s", err)
		return fmt.Errorf("failed to list Secrets")
	}

	if len(secretList.Items) == 0 {
		uninstaller.Logger.Infof("There is no Secret to delete in namespace %s", namespace)
		return nil
	}
	uninstaller.Logger.Infof("Found %d Secrets to delete in namespace %s", len(secretList.Items), namespace)
	for _, secret := range secretList.Items {
		uninstaller.Logger.Infof("Delete Secret %s from Namespace %s", secret.Name, namespace)
		if err := kubevirtClient.DeleteSecret(ctx, namespace, secret.Name); err != nil {
			uninstaller.Logger.Errorf("failed to delete Secrets: %s", err)
			return fmt.Errorf("failed to delete Secret")
		}
		if err := uninstaller.exponentialBackoff("Secret", func() error {
			_, err := kubevirtClient.GetVirtualMachine(ctx, namespace, secret.Name)
			return err
		}); err != nil {
			uninstaller.Logger.Errorf("failed to delete Secrets: %s", err)
			return fmt.Errorf("failed to validate Secret deleted")
		}
	}
	return nil
}

func (uninstaller *ClusterUninstaller) exponentialBackoff(resourceType string, tryGetFunc func() error) error {
	backoff := wait.Backoff{
		Duration: 1 * time.Second,
		Jitter:   0,
		Factor:   2,
		Steps:    5,
	}
	return wait.ExponentialBackoff(backoff, func() (done bool, err error) {
		err = tryGetFunc()
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return true, fmt.Errorf("failed to get %s, with error: %v", resourceType, err)
		}
		return false, nil
	})
}
