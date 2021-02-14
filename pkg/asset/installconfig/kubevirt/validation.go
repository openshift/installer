package kubevirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
	authv1 "k8s.io/api/authorization/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

// ValidatePermissions tests that the current user has the required permissions
// Some permissions are required for the installation
// In addition, the current user is used for accessing the kubevirt/platfrom cluster from tenant cluster. E.g. list VMIs
func ValidatePermissions(client Client, ic *types.InstallConfig) error {
	// Prepare requests for permissions check
	reviewObjs := createReviewObjs(ic.Platform.Kubevirt.Namespace)

	// Collection of missing permissions
	var notAllowedObjs []*authv1.SelfSubjectAccessReview

	// Test each permission
	for _, reviewObj := range reviewObjs {
		reviewObjPointer, err := client.CreateSelfSubjectAccessReview(context.Background(), &reviewObj)
		if err != nil {
			return err
		}

		if !reviewObjPointer.Status.Allowed {
			notAllowedObjs = append(notAllowedObjs, reviewObjPointer)
		}
	}

	// Put all missing permissions in one error message
	if len(notAllowedObjs) > 0 {
		var notAllowed []string
		for _, obj := range notAllowedObjs {
			notAllowed = append(notAllowed, fmt.Sprintf("%+v", *obj.Spec.ResourceAttributes))
		}

		return fmt.Errorf("the user is missing the following permissions: %s", strings.Join(notAllowed, ", "))
	}

	return nil
}

// ValidateForProvisioning is called by PlatformProvisionCheck
func ValidateForProvisioning(client Client) error {
	hcUnstructured, err := client.GetHyperConverged(context.Background(), "kubevirt-hyperconverged", "openshift-cnv")
	if err != nil {
		return fmt.Errorf("failed to get resource openshift-cnv/kubevirt-hyperconverged, with error: %v", err)
	}

	enabled, found, err := unstructured.NestedBool(hcUnstructured.Object, "spec", "featureGates", "hotplugVolumes")
	if err != nil {
		return fmt.Errorf("failed to read boolean value 'spec.featureGates.hotplugVolumes' from resource openshift-cnv/kubevirt-hyperconverged, with error: %v", err)
	}

	if !found || !enabled {
		return fmt.Errorf("feature gate hotplugVolumes is either missing or not set to true. Review resource openshift-cnv/kubevirt-hyperconverged. Follow Kubevirt CSI driver documentation for setting the feature gate (https://github.com/openshift/kubevirt-csi-driver)")
	}

	return nil
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

// createReviewObjs creates requests for testing user permissions
func createReviewObjs(namespace string) []authv1.SelfSubjectAccessReview {
	return []authv1.SelfSubjectAccessReview{
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: "openshift-cnv",
					Group:     "hco.kubevirt.io",
					Resource:  "hyperconvergeds",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "",
					Resource:  "secrets",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "",
					Resource:  "secrets",
					Verb:      "list",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "",
					Resource:  "secrets",
					Verb:      "create",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "",
					Resource:  "secrets",
					Verb:      "delete",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "",
					Resource:  "namespaces",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachines",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachines",
					Verb:      "list",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachines",
					Verb:      "create",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachines",
					Verb:      "delete",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachines",
					Verb:      "update",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachineinstances",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "kubevirt.io",
					Resource:  "virtualmachineinstances",
					Verb:      "list",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "cdi.kubevirt.io",
					Resource:  "datavolumes",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "cdi.kubevirt.io",
					Resource:  "datavolumes",
					Verb:      "list",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "cdi.kubevirt.io",
					Resource:  "datavolumes",
					Verb:      "create",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "cdi.kubevirt.io",
					Resource:  "datavolumes",
					Verb:      "delete",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "k8s.cni.cncf.io",
					Resource:  "network-attachment-definitions",
					Verb:      "get",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "subresources.kubevirt.io",
					Resource:  "virtualmachineinstances/addvolume",
					Verb:      "update",
				},
			},
		},
		{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: namespace,
					Group:     "subresources.kubevirt.io",
					Resource:  "virtualmachineinstances/removevolume",
					Verb:      "update",
				},
			},
		},
	}
}
