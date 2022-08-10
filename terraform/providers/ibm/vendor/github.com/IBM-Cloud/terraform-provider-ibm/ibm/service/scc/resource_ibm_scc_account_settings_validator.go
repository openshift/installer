// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMSccAccountSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 2)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "location_id",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "us, eu, uk",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_account_settings", Schema: validateSchema}
	return &resourceValidator
}
