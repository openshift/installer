// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMSccTemplateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 3)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             1,
			MaxValueLength:             32,
			Regexp:                     ".*",
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             1,
			MaxValueLength:             256,
			Regexp:                     ".*",
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_name",
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     "^[a-z-]*$",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_template", Schema: validateSchema}
	return &resourceValidator
}
