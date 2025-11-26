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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateImage validates an image.
func ValidateImage(image *Image, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if image == nil {
		// allow empty image as it is defaulted in the AzureMachine controller
		return allErrs
	}

	allErrs = append(allErrs, validateSingleDetailsOnly(image, fldPath)...)

	if image.Marketplace != nil {
		allErrs = append(allErrs, validateMarketplaceImage(image, fldPath)...)
	}
	if image.SharedGallery != nil {
		allErrs = append(allErrs, validateSharedGalleryImage(image, fldPath)...)
	}
	if image.ID != nil {
		allErrs = append(allErrs, validateSpecificImage(image, fldPath)...)
	}
	if image.ComputeGallery != nil {
		allErrs = append(allErrs, validateComputeGalleryImage(image, fldPath)...)
	}

	return allErrs
}

func validateSingleDetailsOnly(image *Image, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	imageDetailsFound := (image.ID != nil)

	if image.Marketplace != nil {
		if imageDetailsFound {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("Marketplace"), "Marketplace cannot be used as an image ID has been specified"))
		} else {
			imageDetailsFound = true
		}
	}

	if image.SharedGallery != nil {
		if imageDetailsFound {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("SharedGallery"), "SharedGallery cannot be used as an image ID. Marketplace or ComputeGallery images has been specified"))
		} else {
			imageDetailsFound = true
		}
	}

	if image.ComputeGallery != nil {
		if imageDetailsFound {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("ComputeGallery"), "ComputeGallery cannot be used as an image ID. Marketplace or SharedGallery images has been specified"))
		} else {
			imageDetailsFound = true
		}
	}

	if !imageDetailsFound {
		allErrs = append(allErrs, field.Required(fldPath, "You must supply an ID, Marketplace or ComputeGallery image details"))
	}

	return allErrs
}

func validateComputeGalleryImage(image *Image, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if image.ComputeGallery.SubscriptionID != nil && image.ComputeGallery.ResourceGroup == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ResourceGroup"), "", "ResourceGroup cannot be empty when SubscriptionID is specified"))
	}
	if image.ComputeGallery.ResourceGroup != nil && image.ComputeGallery.SubscriptionID == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("SubscriptionID"), "", "SubscriptionID cannot be empty when ResourceGroup is specified"))
	}

	return allErrs
}

func validateSharedGalleryImage(image *Image, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if image.SharedGallery.SubscriptionID == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("SubscriptionID"), "", "SubscriptionID cannot be empty when specifying an AzureSharedGalleryImage"))
	}
	if image.SharedGallery.ResourceGroup == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ResourceGroup"), "", "ResourceGroup cannot be empty when specifying an AzureSharedGalleryImage"))
	}
	if image.SharedGallery.Gallery == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Gallery"), "", "Gallery cannot be empty when specifying an AzureSharedGalleryImage"))
	}
	if image.SharedGallery.Name == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Name"), "", "Name cannot be empty when specifying an AzureSharedGalleryImage"))
	}
	if image.SharedGallery.Version == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Version"), "", "Version cannot be empty when specifying an AzureSharedGalleryImage"))
	}

	return allErrs
}

func validateMarketplaceImage(image *Image, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if image.Marketplace.Publisher == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Publisher"), "", "Publisher cannot be empty when specifying an AzureMarketplaceImage"))
	}
	if image.Marketplace.Offer == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Offer"), "", "Offer cannot be empty when specifying an AzureMarketplaceImage"))
	}
	if image.Marketplace.SKU == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("SKU"), "", "SKU cannot be empty when specifying an AzureMarketplaceImage"))
	}
	if image.Marketplace.Version == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Version"), "", "Version cannot be empty when specifying an AzureMarketplaceImage"))
	}
	return allErrs
}

func validateSpecificImage(image *Image, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if *image.ID == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ID"), "", "ID cannot be empty when specifying an AzureImageByID"))
	}

	return allErrs
}
