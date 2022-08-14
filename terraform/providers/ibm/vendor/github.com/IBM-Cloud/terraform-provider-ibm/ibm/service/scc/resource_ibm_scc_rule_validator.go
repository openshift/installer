// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMSccRuleValidator() *validate.ResourceValidator {

	validateSchemaList := make([]validate.ValidateSchema, 2)
	validateSchemaList = append(validateSchemaList, validateIBMSccRuleReqConfig())
	resourceValidator := validate.ResourceValidator{
		ResourceName: "ibm_scc_rule",
		Schema:       validateSchemaList,
	}
	return &resourceValidator
}

func validateIBMSccRuleReqConfig() validate.ValidateSchema {
	validateSchema := validate.ValidateSchema{
		Identifier:                 "operator",
		ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
		Type:                       validate.TypeString,
		Required:                   true,
		AllowedValues:              "is_true, is_false, is_empty, is_not_empty, string_equals, string_not_equals, string_match, string_not_match, num_equals, num_not_equals, num_less_than, num_less_than_equals, num_greater_than, num_greater_than_equals",
	}
	return validateSchema
}
