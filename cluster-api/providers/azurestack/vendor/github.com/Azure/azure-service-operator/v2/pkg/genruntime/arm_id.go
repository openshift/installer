/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/pkg/errors"
)

// GetAndParseResourceID gets the ARM ID from the given MetaObject and parses it into its constituent parts
func GetAndParseResourceID(obj ARMMetaObject) (*arm.ResourceID, error) {
	resourceID, hasResourceID := GetResourceID(obj)
	if !hasResourceID {
		return nil, errors.Errorf("cannot find resource id for obj %s/%s", obj.GetNamespace(), obj.GetName())
	}

	return arm.ParseResourceID(resourceID)
}

// TODO: We really want these methods to be on ARMMetaObject itself -- should update code generator to make them at some point
func GetResourceID(obj ARMMetaObject) (string, bool) {
	result, ok := obj.GetAnnotations()[ResourceIDAnnotation]
	return result, ok
}

func GetResourceIDOrDefault(obj ARMMetaObject) string {
	return obj.GetAnnotations()[ResourceIDAnnotation]
}

func SetResourceID(obj ARMMetaObject, id string) {
	AddAnnotation(obj, ResourceIDAnnotation, id)
}

func SetChildResourceIDOverride(obj ARMMetaObject, id string) {
	AddAnnotation(obj, ChildResourceIDOverrideAnnotation, id)
}

func GetChildResourceIDOverride(obj ARMMetaObject) (string, bool) {
	result, ok := obj.GetAnnotations()[ChildResourceIDOverrideAnnotation]
	return result, ok
}

func CheckARMIDMatchesSubscription(subscriptionID string, armID *arm.ResourceID) bool {
	// armIDSub may be empty if there is no subscription
	if armID.SubscriptionID != "" {
		if !strings.EqualFold(armID.SubscriptionID, subscriptionID) {
			return false
		}
	}
	return true
}
