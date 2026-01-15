// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package storage

import (
	v20230630s "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1api20230630/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

var _ augmentConversionForFlexibleServer_Spec = &FlexibleServer_Spec{}

// AssignPropertiesFrom implements augmentConversionForFlexibleServer.
func (server *FlexibleServer_Spec) AssignPropertiesFrom(src *v20230630s.FlexibleServer_Spec) error {
	// Clone the existing property bag
	pb := genruntime.NewPropertyBag(server.PropertyBag)

	if src.SourceServerResourceReference != nil {
		armID := src.SourceServerResourceReference.ARMID
		server.SourceServerResourceId = &armID
		pb.Remove("SourceServerResourceReference")
	}

	// Store updated property bag
	server.PropertyBag = pb

	return nil
}

// AssignPropertiesTo implements augmentConversionForFlexibleServer.
func (server *FlexibleServer_Spec) AssignPropertiesTo(dst *v20230630s.FlexibleServer_Spec) error {
	// Get the current property bag
	pb := genruntime.NewPropertyBag(dst.PropertyBag)

	if server.SourceServerResourceId != nil {
		dst.SourceServerResourceReference = &genruntime.ResourceReference{
			ARMID: *server.SourceServerResourceId,
		}

		pb.Remove("SourceServerResourceId")
	}

	// Store updated property bag
	dst.PropertyBag = pb

	return nil
}
