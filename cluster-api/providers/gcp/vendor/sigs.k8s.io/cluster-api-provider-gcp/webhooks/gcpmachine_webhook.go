/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooks

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/utils/strings/slices"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Confidential VM Technology support depends on the configured machine types.
// reference: https://cloud.google.com/compute/confidential-vm/docs/os-and-machine-type#machine-type
var (
	confidentialMachineSeriesSupportingSev    = []string{"n2d", "c2d", "c3d"}
	confidentialMachineSeriesSupportingSevsnp = []string{"n2d"}
	confidentialMachineSeriesSupportingTdx    = []string{"c3"}
)

func (m *GCPMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.GCPMachine{}).
		WithValidator(m).
		WithDefaulter(m).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=gcpmachines,versions=v1beta1,name=validation.gcpmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmachine,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=gcpmachines,versions=v1beta1,name=default.gcpmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// GCPMachine implements a validating and defaulting webhook for GCPMachine.
type GCPMachine struct{}

var (
	_ webhook.CustomValidator = &GCPMachine{}
	_ webhook.CustomDefaulter = &GCPMachine{}
)

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*GCPMachine) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	m, ok := obj.(*infrav1.GCPMachine)
	if !ok {
		return nil, fmt.Errorf("expected an GCPMachine object but got %T", m)
	}

	clusterlog.Info("validate create", "name", m.Name)

	if err := validateConfidentialCompute(m.Spec); err != nil {
		return nil, err
	}
	return nil, validateCustomerEncryptionKey(m.Spec)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*GCPMachine) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	m, ok := newObj.(*infrav1.GCPMachine)
	if !ok {
		return nil, fmt.Errorf("expected an GCPMachine object but got %T", m)
	}

	newGCPMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(m)
	if err != nil {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPMachine").GroupKind(), m.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert new GCPMachine to unstructured object")),
		})
	}
	oldGCPMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObj)
	if err != nil {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPMachine").GroupKind(), m.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert old GCPMachine to unstructured object")),
		})
	}

	newGCPMachineSpec := newGCPMachine["spec"].(map[string]interface{})
	oldGCPMachineSpec := oldGCPMachine["spec"].(map[string]interface{})

	// allow changes to providerID
	delete(oldGCPMachineSpec, "providerID")
	delete(newGCPMachineSpec, "providerID")

	// allow changes to additionalLabels
	delete(oldGCPMachineSpec, "additionalLabels")
	delete(newGCPMachineSpec, "additionalLabels")

	// allow changes to additionalNetworkTags
	delete(oldGCPMachineSpec, "additionalNetworkTags")
	delete(newGCPMachineSpec, "additionalNetworkTags")

	if !reflect.DeepEqual(oldGCPMachineSpec, newGCPMachineSpec) {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPMachine").GroupKind(), m.Name, field.ErrorList{
			field.Forbidden(field.NewPath("spec"), "cannot be modified"),
		})
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*GCPMachine) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Default implements webhookutil.defaulter so a webhook will be registered for the type.
func (*GCPMachine) Default(_ context.Context, _ runtime.Object) error {
	return nil
}

func validateConfidentialCompute(spec infrav1.GCPMachineSpec) error {
	if spec.ConfidentialCompute != nil && *spec.ConfidentialCompute != infrav1.ConfidentialComputePolicyDisabled {
		if spec.OnHostMaintenance == nil || *spec.OnHostMaintenance == infrav1.HostMaintenancePolicyMigrate {
			return fmt.Errorf("ConfidentialCompute require OnHostMaintenance to be set to %s, the current value is: %s", infrav1.HostMaintenancePolicyTerminate, infrav1.HostMaintenancePolicyMigrate)
		}

		machineSeries := strings.Split(spec.InstanceType, "-")[0]
		switch *spec.ConfidentialCompute {
		case infrav1.ConfidentialComputePolicyEnabled, infrav1.ConfidentialComputePolicySEV:
			if !slices.Contains(confidentialMachineSeriesSupportingSev, machineSeries) {
				return fmt.Errorf("ConfidentialCompute %s requires any of the following machine series: %s. %s was found instead", *spec.ConfidentialCompute, strings.Join(confidentialMachineSeriesSupportingSev, ", "), spec.InstanceType)
			}
		case infrav1.ConfidentialComputePolicySEVSNP:
			if !slices.Contains(confidentialMachineSeriesSupportingSevsnp, machineSeries) {
				return fmt.Errorf("ConfidentialCompute %s requires any of the following machine series: %s. %s was found instead", *spec.ConfidentialCompute, strings.Join(confidentialMachineSeriesSupportingSevsnp, ", "), spec.InstanceType)
			}
		case infrav1.ConfidentialComputePolicyTDX:
			if !slices.Contains(confidentialMachineSeriesSupportingTdx, machineSeries) {
				return fmt.Errorf("ConfidentialCompute %s requires any of the following machine series: %s. %s was found instead", *spec.ConfidentialCompute, strings.Join(confidentialMachineSeriesSupportingTdx, ", "), spec.InstanceType)
			}
		default:
			return fmt.Errorf("invalid ConfidentialCompute %s", *spec.ConfidentialCompute)
		}
	}
	return nil
}

func checkKeyType(key *infrav1.CustomerEncryptionKey) error {
	switch key.KeyType {
	case infrav1.CustomerManagedKey:
		if key.ManagedKey == nil || key.SuppliedKey != nil {
			return errors.New("CustomerEncryptionKey KeyType of Managed requires only ManagedKey to be set")
		}
	case infrav1.CustomerSuppliedKey:
		if key.SuppliedKey == nil || key.ManagedKey != nil {
			return errors.New("CustomerEncryptionKey KeyType of Supplied requires only SuppliedKey to be set")
		}
		if len(key.SuppliedKey.RawKey) > 0 && len(key.SuppliedKey.RSAEncryptedKey) > 0 {
			return errors.New("CustomerEncryptionKey KeyType of Supplied requires either RawKey or RSAEncryptedKey to be set, not both")
		}
	default:
		return fmt.Errorf("invalid value for CustomerEncryptionKey KeyType %s", key.KeyType)
	}
	return nil
}

func validateCustomerEncryptionKey(spec infrav1.GCPMachineSpec) error {
	if spec.RootDiskEncryptionKey != nil {
		if err := checkKeyType(spec.RootDiskEncryptionKey); err != nil {
			return err
		}
	}

	for _, disk := range spec.AdditionalDisks {
		if disk.EncryptionKey != nil {
			if err := checkKeyType(disk.EncryptionKey); err != nil {
				return err
			}
		}
	}
	return nil
}
