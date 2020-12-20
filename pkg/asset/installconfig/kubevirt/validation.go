package kubevirt

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate executes kubevirt specific validation
func Validate(ic *types.InstallConfig, client Client) error {
	allErrs := field.ErrorList{}
	fldPath := field.NewPath("platform", "kubevirt")
	kubevirtPlatform := ic.Platform.Kubevirt

	allErrs = append(allErrs, validateNamespace(kubevirtPlatform.Namespace, client, fldPath.Child("namespace"))...)
	allErrs = append(allErrs, validateStorageClassExistsInInfraCluster(kubevirtPlatform.StorageClass, client, fldPath.Child("storageClass"))...)
	allErrs = append(allErrs, validateNetworkAttachmentDefinitionExistsInNamespace(kubevirtPlatform.NetworkName, kubevirtPlatform.Namespace, client, fldPath.Child("networkName"))...)

	return allErrs.ToAggregate()
}

func validateStorageClassExistsInInfraCluster(name string, client Client, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// StorageClass can be empty string, in this case will use default Storage Class
	// Therefore, make the call to the cluster only if its not empty string
	if name == "" {
		return allErrs
	}
	if _, err := client.GetStorageClass(context.Background(), name); err != nil {
		allErrs = append(
			allErrs,
			field.Invalid(
				fieldPath,
				name,
				fmt.Sprintf("failed to get StorageClass from InfraCluster, with error: %v", err),
			),
		)
	}

	return allErrs
}

// validateNetworkAttachmentDefinitionExistsInNamespace validate the following:
// 1. The namespace exists
// 2. The network-attachment-definition exists
func validateNetworkAttachmentDefinitionExistsInNamespace(name string, namespace string, client Client, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if _, err := client.GetNetworkAttachmentDefinition(context.Background(), name, namespace); err != nil {
		allErrs = append(
			allErrs,
			field.Invalid(
				fieldPath,
				name,
				fmt.Sprintf("failed to get network-attachment-definition from InfraCluster, with error: %v", err),
			),
		)
	}

	return allErrs
}

func validateNamespace(namespace string, client Client, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	ns, err := client.GetNamespace(context.Background(), namespace)
	if err != nil {
		allErrs = append(
			allErrs,
			field.Invalid(
				fieldPath,
				namespace,
				fmt.Sprintf("failed to get namepsace, with error: %v", err),
			),
		)
		return allErrs
	}
	labelRequiredKey := "mutatevirtualmachines.kubemacpool.io"
	labelRequiredVal := "allocate"
	labelVal, ok := ns.Labels[labelRequiredKey]
	if !ok || labelVal != labelRequiredVal {
		allErrs = append(
			allErrs,
			field.Invalid(
				fieldPath,
				namespace,
				fmt.Sprintf("KubeMacPool component is not enabled for the namespace, the namespace must have label \"%s: %s\"", labelRequiredKey, labelRequiredVal),
			),
		)
	}

	return allErrs
}
